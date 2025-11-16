import onewire
import ds18x20
import time
from machine import Pin


class Temp:
    """
    DS18B20 temperature helper.
    """

    def __init__(self, pin: int = 22) -> None:
        self.ds = ds18x20.DS18X20(onewire.OneWire(Pin(pin)))
        self.roms = self.ds.scan()
        self.last_c = None
        self._last_convert_ms: int = 0

    def read_c(self):
        """
        Trigger a conversion if needed, then read the first sensor.
        Returns None if no sensor found.
        """
        if not self.roms:
            return None

        now = time.ticks_ms()
        # Conversion time ~750ms for 12-bit; only trigger if stale
        if time.ticks_diff(now, self._last_convert_ms) > 1000:
            self.ds.convert_temp()
            time.sleep_ms(750)
            self._last_convert_ms = now

        c = self.ds.read_temp(self.roms[0])
        self.last_c = c
        return c


class TempNonBlocking:
    """
    Non-blocking DS18B20 reader.
    Call poll() frequently (e.g., every loop). Use get_c() to read latest value.
    """
    def __init__(self, pin: int = 22, interval_ms: int = 1000) -> None:
        self.ds = ds18x20.DS18X20(onewire.OneWire(Pin(pin)))
        self.roms = self.ds.scan()
        self.last_c = None
        self._interval_ms = interval_ms
        self._state = "idle"  # idle | converting
        self._t_ms = 0

    def poll(self) -> None:
        if not self.roms:
            return
        now = time.ticks_ms()
        if self._state == "idle":
            # Start a new conversion if interval elapsed
            if self._t_ms == 0 or time.ticks_diff(now, self._t_ms) >= self._interval_ms:
                try:
                    self.ds.convert_temp()
                    self._t_ms = now
                    self._state = "converting"
                except Exception:
                    # Sensor error; stay idle and try later
                    self._state = "idle"
        elif self._state == "converting":
            # Check if conversion window (~750ms) passed
            if time.ticks_diff(now, self._t_ms) >= 750:
                try:
                    c = self.ds.read_temp(self.roms[0])
                    self.last_c = c
                except Exception:
                    pass
                # Schedule next cycle
                self._t_ms = now
                self._state = "idle"

    def get_c(self):
        return self.last_c

