import time
from machine import Pin


class LedSnake:
    """
    Simple three-LED "snake" indicator on GP10/11/12 by default.
    Pattern cycles 1 -> 2 -> 3 -> 2 -> 1 ...
    Speed varies with remaining time thresholds.
    """

    def __init__(self, pins=(10, 11, 12), debug: bool = False) -> None:
        self._pins = [Pin(p, Pin.OUT, value=0) for p in pins]
        # Sequence of active index positions
        self._seq = [0, 1, 2, 1]
        self._pos = 0
        self._last_ms = time.ticks_ms()
        self._period_ms = 3500
        self._last_period_ms = self._period_ms
        self._debug = debug
        self._blink_on = False

    def _all_off(self) -> None:
        for led in self._pins:
            led.value(0)

    def _step(self) -> None:
        # Turn on only the current index; others off
        idx = self._seq[self._pos]
        for i, led in enumerate(self._pins):
            led.value(1 if i == idx else 0)
        self._pos = (self._pos + 1) % len(self._seq)

    def _update_period(self, remaining_sec: int) -> None:
        # Default slow "idle" during normal run
        if remaining_sec > 60:
            self._period_ms = 3500  # ~3.5s
        elif remaining_sec > 30:
            self._period_ms = 500
        elif remaining_sec > 10:
            self._period_ms = 200
        else:
            self._period_ms = 100
        if self._period_ms != self._last_period_ms:
            if self._debug:
                print("leds: period_ms ->", self._period_ms, "remain_s=", remaining_sec)
            self._last_period_ms = self._period_ms

    def update(self, remaining_sec: int, active: bool, overtime_sec: int = 0) -> None:
        """
        Call frequently. If active, advance snake based on time.
        When remaining reaches 0 or inactive, turn LEDs off.
        """
        if not active:
            self._all_off()
            return
        # Overtime: blink all LEDs in unison (fast)
        if overtime_sec and overtime_sec > 0:
            # Fixed fast blink in overtime
            period = 100
            now = time.ticks_ms()
            if time.ticks_diff(now, self._last_ms) >= period:
                self._last_ms = now
                self._blink_on = not self._blink_on
                if self._debug:
                    print("leds: overtime blink ->", int(self._blink_on))
                for led in self._pins:
                    led.value(1 if self._blink_on else 0)
            return
        self._update_period(remaining_sec)
        now = time.ticks_ms()
        if time.ticks_diff(now, self._last_ms) >= self._period_ms:
            self._last_ms = now
            self._step()


