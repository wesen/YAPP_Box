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


