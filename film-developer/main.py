from lib.display import Display
from lib.async_buttons import GhostDebouncedButton
import uasyncio as asyncio
from machine import Pin
import time
import onewire
import ds18x20


BTN_PIN = 14  # Single button (active-low)

# Small fixed ring buffer for IRQ-safe event enqueue
_Q = [0] * 16
_Q_MASK = 15
_head = 0
_tail = 0
_flag = asyncio.ThreadSafeFlag()

GHOST_MS = 25      # ignore further IRQs for this many ms after first edge

DEBOUNCE_PRESS_MS = 3
DEBOUNCE_RELEASE_MS = 3


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
    # Deprecated: replaced by GhostDebouncedButton.on_press callback
    pass
    s["pending_draws"] = s.get("pending_draws", 0) + 1


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
        if state.get("dirty") or state.get("pending_draws", 0) > 0:
            display.clear()
            display.text_at(0, 0, "ASYNC BUTTON DEMO")
            display.text_at(2, 0, "BTN1 on GP14")
            display.text_at(4, 0, "Presses: {}".format(state["presses"]))
            # Clipped redraw counter (wrap at 100000)
            cnt = (state.get("draw_count", 0) + 1) % 100000
            state["draw_count"] = cnt
            display.text_at(5, 0, "Draws: {}".format(cnt))
            display.show()
            state["dirty"] = False
            if state.get("pending_draws", 0) > 0:
                state["pending_draws"] -= 1
        await asyncio.sleep_ms(10)

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

    state = {
        "presses": 0,
        "dirty": True,
        "draw_count": 0,
        "pending_draws": 0,
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

    # Button press callback (scheduled context safe)
    def on_press(_button_id: int) -> None:
        t_now = time.ticks_ms()
        last_press = state["stats"]["last_press_ms"]
        if last_press is not None:
            gap = time.ticks_diff(t_now, last_press)
            st = state["stats"]
            st["gap_sum"] += gap
            st["gap_cnt"] += 1
            if st["gap_min"] is None or gap < st["gap_min"]:
                st["gap_min"] = gap
            if gap > st["gap_max"]:
                st["gap_max"] = gap
        state["stats"]["last_press_ms"] = t_now
        state["presses"] += 1
        state["stats"]["accepted"] += 1
        state["dirty"] = True
        state["pending_draws"] = state.get("pending_draws", 0) + 1

    # Configure single button (active-low) with ghosted debounce
    _btn = GhostDebouncedButton(BTN_PIN, ghost_ms=GHOST_MS, on_press=on_press, button_id=1)
    await asyncio.gather(
        # button presses are handled via scheduled callback; no button task needed
        display_task(display, state),
        # stats_task(state),
    )


def main() -> None:
    asyncio.run(amain())

# -------------------------------
# Integrated test: Display + Temp + 3 Buttons (debounced) + LEDs
# -------------------------------

async def display_task_integ(display, state):
    while True:
        if state.get("dirty") or state.get("pending_draws", 0) > 0:
            display.clear()
            display.text_at(0, 0, "FILM DEV TEST")
            temp_line = "--.-C" if state.get("temp_c") is None else "{:.1f}C".format(state.get("temp_c"))
            display.text_at(2, 0, "Temp: {}".format(temp_line))
            display.text_at(3, 0, "BTN1/2/3 toggle LED")
            display.text_at(4, 0, "P1:{} P2:{} P3:{}".format(state.get("p1", 0), state.get("p2", 0), state.get("p3", 0)))
            cnt = (state.get("draw_count", 0) + 1) % 100000
            state["draw_count"] = cnt
            display.text_at(5, 0, "Draws: {}".format(cnt))
            display.show()
            state["dirty"] = False
            if state.get("pending_draws", 0) > 0:
                state["pending_draws"] -= 1
        await asyncio.sleep_ms(10)

async def temp_task_integ(state):
    ow = onewire.OneWire(Pin(22))
    ds = ds18x20.DS18X20(ow)
    roms = ds.scan()
    rom = roms[0] if roms else None
    while True:
        if rom is not None:
            ds.convert_temp()
            await asyncio.sleep_ms(750)
            c = ds.read_temp(rom)
            state["temp_c"] = c
            state["dirty"] = True
            state["pending_draws"] = state.get("pending_draws", 0) + 1
        else:
            state["temp_c"] = None
        await asyncio.sleep_ms(250)

async def amain2():
    display = Display(spi_baudrate=1_000_000)
    # LEDs on GP10/11/12
    LED_PINS = (10, 11, 12)
    leds = [Pin(p, Pin.OUT, value=0) for p in LED_PINS]
    led_state = [0, 0, 0]

    state = {
        "dirty": True,
        "draw_count": 0,
        "pending_draws": 0,
        "temp_c": None,
        "p1": 0,
        "p2": 0,
        "p3": 0,
    }

    def on_press(button_id: int) -> None:
        # Toggle corresponding LED
        if 1 <= button_id <= 3:
            idx = button_id - 1
            led_state[idx] ^= 1
            leds[idx].value(led_state[idx])
        # Update per-button counters
        if button_id == 1:
            state["p1"] = state.get("p1", 0) + 1
        elif button_id == 2:
            state["p2"] = state.get("p2", 0) + 1
        elif button_id == 3:
            state["p3"] = state.get("p3", 0) + 1
        state["dirty"] = True
        state["pending_draws"] = state.get("pending_draws", 0) + 1

    # Three debounced buttons
    GhostDebouncedButton(14, ghost_ms=GHOST_MS, on_press=on_press, button_id=1)
    GhostDebouncedButton(15, ghost_ms=GHOST_MS, on_press=on_press, button_id=2)
    GhostDebouncedButton(16, ghost_ms=GHOST_MS, on_press=on_press, button_id=3)

    await asyncio.gather(
        display_task_integ(display, state),
        temp_task_integ(state),
    )

def main() -> None:
    # Main UI with non-blocking temperature
    try:
        from lib.input_irq import AsyncInput as InputLike
    except Exception:
        from lib.input import Input as InputLike
    from lib.temp import TempNonBlocking
    from lib.ui import UI
    d = Display(spi_baudrate=1_000_000)
    try:
        # Prefer IRQ-backed input for reliable short presses
        i = InputLike(ghost_ms=40, debug=True)
    except TypeError:
        # Fallback to polling signature
        i = InputLike(debounce_ms=40, debug=True)
    t = TempNonBlocking()
    ui = UI(d, i, t, enable_temp=True)
    while True:
        ui.handle()
        ui.render()
        time.sleep_ms(20)

        
def main3():
    import uasyncio as asyncio
    from lib.test_async_button_demo import run_async_button_demo3
    asyncio.run(run_async_button_demo3(ghost_ms=40))

def main4():
    from lib.display import Display
    from lib.input import Input
    from lib.test_input import run_input_demo
    run_input_demo(Display(), Input(debounce_ms=40))

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        pass


