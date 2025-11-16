import micropython
import time
from machine import Pin

# Ensure exceptions in IRQ context can be reported
micropython.alloc_emergency_exception_buf(100)


class GhostDebouncedButton:
    """
    Minimal falling-edge button with ghost-window debounce.
    - Accepts first falling-edge (active-low wiring)
    - Ignores subsequent edges for ghost_ms
    - Schedules on_press(button_id) via micropython.schedule
    """

    def __init__(self, pin_num: int, ghost_ms: int = 25, on_press=None, button_id: int = 1):
        self.btn = Pin(pin_num, Pin.IN, Pin.PULL_UP)
        self.ghost_ms = int(ghost_ms)
        self.on_press = on_press
        self.button_id = button_id
        self._suppress_until_ms = 0
        self._enabled = True
        self.btn.irq(trigger=Pin.IRQ_FALLING, handler=self._irq)

    def _irq(self, pin):
        if not self._enabled:
            return
        now = time.ticks_ms()
        if time.ticks_diff(now, self._suppress_until_ms) < 0:
            return
        # Extend suppression window
        try:
            self._suppress_until_ms = time.ticks_add(now, self.ghost_ms)
        except AttributeError:
            self._suppress_until_ms = now + self.ghost_ms
        # Schedule user callback
        if self.on_press is not None:
            try:
                micropython.schedule(self._sched, 0)
            except Exception:
                pass

    def _sched(self, _):
        try:
            self.on_press(self.button_id)
        except Exception:
            # Swallow exceptions from user code in scheduled context
            pass

    def set_ghost_ms(self, ghost_ms: int) -> None:
        self.ghost_ms = int(ghost_ms)

    def disable(self) -> None:
        self._enabled = False
        try:
            self.btn.irq(handler=None)
        except Exception:
            pass

    def enable(self) -> None:
        self._enabled = True
        self.btn.irq(trigger=Pin.IRQ_FALLING, handler=self._irq)


