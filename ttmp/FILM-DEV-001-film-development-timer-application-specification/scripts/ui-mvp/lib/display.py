from machine import Pin, SPI
from .sh1107 import SH1107


class Display:
    """
    Text-mode display helper targeting a 21x8 character grid on a 128x64 SH1107.
    """

    def __init__(self) -> None:
        # SPI0: SCK=GP18, MOSI=GP19 (no MISO used)
        self.spi = SPI(
            0,
            baudrate=1_000_000,
            polarity=0,
            phase=0,
            sck=Pin(18),
            mosi=Pin(19),
        )
        # SH1107 pins: DC=GP20, CS=GP17, RST=GP21
        self.dev = SH1107(width=128, height=64, spi=self.spi, dc=Pin(20), cs=Pin(17), rst=Pin(21))
        self.clear()

    def clear(self) -> None:
        self.dev.fill(0)

    def text_at(self, row: int, col: int, s: str, color: int = 1) -> None:
        """
        Draw text aligned to the 6x8 font grid.
        Row range: 0..7, Col range: 0..20 (21 chars)
        """
        if row < 0 or row > 7:
            return
        if col < 0:
            col = 0
        # Ensure we do not overflow visible width
        s = s[: max(0, 21 - col)]
        x = col * 6
        y = row * 8
        self.dev.text(s, x, y, color)

    def show(self) -> None:
        self.dev.show()


