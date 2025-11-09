from machine import Pin
import time


BTN_PINS = (14, 15, 16)  # BTN1, BTN2, BTN3


class Input:
    """
    Simple debounced edge detector for three buttons with internal pull-ups.
    Emits events on press (active-low -> 0).
    """

    def __init__(self, debounce_ms: int = 50) -> None:
        self.buttons = [Pin(p, Pin.IN, Pin.PULL_UP) for p in BTN_PINS]
        self.last_value = [1, 1, 1]
        self.last_change_ms = [0, 0, 0]
        self.debounce_ms = debounce_ms

    def read(self) -> list[int]:
        """
        Return a list of button indices (1..3) pressed since last read.
        """
        now = time.ticks_ms()
        events: list[int] = []

        for idx, btn in enumerate(self.buttons):
            value = btn.value()
            if value != self.last_value[idx]:
                if time.ticks_diff(now, self.last_change_ms[idx]) > self.debounce_ms:
                    self.last_change_ms[idx] = now
                    self.last_value[idx] = value
                    if value == 0:
                        events.append(idx + 1)
        return events


