from lib.display import Display
import uasyncio as asyncio
import micropython
from machine import Pin
import time


BTN_PIN = 14  # Single button (active-low)
micropython.alloc_emergency_exception_buf(100)

# Small fixed ring buffer for IRQ-safe event enqueue
_Q = [0] * 16
_Q_MASK = 15
_head = 0
_tail = 0
_flag = asyncio.ThreadSafeFlag()

_irq_total = 0
_irq_dropped = 0
_irq_peak_depth = 0
_irq_suppress = 0
_state_ref = None  # set in amain so scheduled callbacks can update state
GHOST_MS = 25      # ignore further IRQs for this many ms after first edge
_suppress_until_ms = 0

DEBOUNCE_PRESS_MS = 3
DEBOUNCE_RELEASE_MS = 3


def _btn_irq(pin):
    # Simple ghosting debounce: accept first edge, suppress further for GHOST_MS
    # Keep ISR tiny; schedule userland update via micropython.schedule
    global _irq_total, _suppress_until_ms
    now = time.ticks_ms()
    if time.ticks_diff(now, _suppress_until_ms) < 0:
        _irq_total += 1
        return
    _irq_total += 1
    # Extend suppression window
    _suppress_until_ms = time.ticks_add(now, GHOST_MS) if hasattr(time, "ticks_add") else (now + GHOST_MS)
    try:
        micropython.schedule(_sched_click, 0)
    except Exception:
        pass


def _pop_events():
    global _tail
    count = 0
    while _tail != _head:
        _tail = (_tail + 1) & 15
        count += 1
    return count

def _pop_one():
    global _tail
    if _tail != _head:
        _tail = (_tail + 1) & 15
        return True
    return False

def _sched_click(_):
    # Runs in soft-IRQ context (scheduled), safe to touch Python state
    s = _state_ref
    if s is None:
        return
    t_now = time.ticks_ms()
    last_press = s["stats"]["last_press_ms"]
    if last_press is not None:
        gap = time.ticks_diff(t_now, last_press)
        st = s["stats"]
        st["gap_sum"] += gap
        st["gap_cnt"] += 1
        if st["gap_min"] is None or gap < st["gap_min"]:
            st["gap_min"] = gap
        if gap > st["gap_max"]:
            st["gap_max"] = gap
    s["stats"]["last_press_ms"] = t_now
    s["presses"] += 1
    s["stats"]["accepted"] += 1
    s["dirty"] = True


async def buttons_task(state, btn):
    # Debounce and one-event-per-press in task context (no time gate)
    while True:
        if _tail == _head:
            await _flag.wait()
            continue
        # Consume exactly one falling-edge event
        if not _pop_one():
            continue
        t_down = time.ticks_ms()
        # Gap stats
        last_press = state["stats"]["last_press_ms"]
        if last_press is not None:
            gap = time.ticks_diff(t_down, last_press)
            s = state["stats"]
            s["gap_sum"] += gap
            s["gap_cnt"] += 1
            if s["gap_min"] is None or gap < s["gap_min"]:
                s["gap_min"] = gap
            if gap > s["gap_max"]:
                s["gap_max"] = gap
        state["stats"]["last_press_ms"] = t_down
        # Accept press immediately on first edge
        state["presses"] += 1
        state["stats"]["accepted"] += 1
        state["dirty"] = True
        # Suppress further IRQ enqueues during this press; drop any queued bounce now
        global _irq_suppress
        _irq_suppress = 1
        _pop_events()
        # Wait for release (pin high)
        while btn.value() == 0:
            await asyncio.sleep_ms(1)
        t_up = time.ticks_ms()
        # Hold duration stats
        hold = time.ticks_diff(t_up, t_down)
        s = state["stats"]
        s["hold_sum"] += hold
        s["hold_cnt"] += 1
        if s["hold_min"] is None or hold < s["hold_min"]:
            s["hold_min"] = hold
        if hold > s["hold_max"]:
            s["hold_max"] = hold
        # Ensure stable high for release debounce
        stable_start = time.ticks_ms()
        while True:
            if btn.value() == 0:
                # went low again; wait high again and restart stability window
                while btn.value() == 0:
                    await asyncio.sleep_ms(1)
                stable_start = time.ticks_ms()
            if time.ticks_diff(time.ticks_ms(), stable_start) >= DEBOUNCE_RELEASE_MS:
                break
            await asyncio.sleep_ms(1)
        # Re-enable IRQ enqueues for next press
        _irq_suppress = 0


async def display_task(display, state):
    while True:
        if state.get("dirty"):
            display.clear()
            display.text_at(0, 0, "ASYNC BUTTON DEMO")
            display.text_at(2, 0, "BTN1 on GP14")
            display.text_at(4, 0, "Presses: {}".format(state["presses"]))
            display.show()
            state["dirty"] = False
        await asyncio.sleep_ms(50)

async def stats_task(state):
    global _irq_total, _irq_dropped, _irq_peak_depth
    while True:
        await asyncio.sleep_ms(2000)
        s = state["stats"]
        gap_avg = (s["gap_sum"] // s["gap_cnt"]) if s["gap_cnt"] else 0
        hold_avg = (s["hold_sum"] // s["hold_cnt"]) if s["hold_cnt"] else 0
        print(
            "stats: irq_total={}, dropped={}, peak_depth={}, accepted={}, ignored={}, "
            "gap_ms(avg/min/max)={}/{}/{} hold_ms(avg/min/max)={}/{}/{}".format(
                _irq_total, _irq_dropped, _irq_peak_depth,
                s["accepted"], s["debounce_ignored"],
                gap_avg, s["gap_min"] if s["gap_min"] is not None else -1, s["gap_max"],
                hold_avg, s["hold_min"] if s["hold_min"] is not None else -1, s["hold_max"],
            )
        )
        _irq_total = 0
        _irq_dropped = 0
        _irq_peak_depth = 0
        s["accepted"] = 0
        s["debounce_ignored"] = 0
        s["gap_sum"] = 0
        s["gap_cnt"] = 0
        s["gap_min"] = None
        s["gap_max"] = 0
        s["hold_sum"] = 0
        s["hold_cnt"] = 0
        s["hold_min"] = None
        s["hold_max"] = 0

async def amain():
    # Keep SPI very low for stability during bring-up
    display = Display(spi_baudrate=1_000_000)

    # Configure single button (active-low) with IRQ
    btn = Pin(BTN_PIN, Pin.IN, Pin.PULL_UP)
    btn.irq(trigger=Pin.IRQ_FALLING, handler=_btn_irq)

    state = {
        "presses": 0,
        "dirty": True,
        "stats": {
            "accepted": 0,
            "debounce_ignored": 0,
            "last_press_ms": None,
            "gap_sum": 0,
            "gap_cnt": 0,
            "gap_min": None,
            "gap_max": 0,
            "hold_sum": 0,
            "hold_cnt": 0,
            "hold_min": None,
            "hold_max": 0,
        },
    }
    # expose to scheduled callback
    global _state_ref
    _state_ref = state
    await asyncio.gather(
        # button presses are handled via scheduled callback; no button task needed
        display_task(display, state),
        stats_task(state),
    )


def main() -> None:
    asyncio.run(amain())


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        pass


