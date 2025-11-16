import uasyncio as asyncio
import time
from lib.display import Display
from lib.async_buttons import GhostDebouncedButton


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


async def run_async_button_demo(spi_baudrate: int = 1_000_000, ghost_ms: int = 25):
    display = Display(spi_baudrate=spi_baudrate)

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

    _btn = GhostDebouncedButton(14, ghost_ms=ghost_ms, on_press=on_press, button_id=1)

    await asyncio.gather(
        display_task(display, state),
    )

