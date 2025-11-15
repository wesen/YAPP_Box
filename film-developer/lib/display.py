from machine import Pin, SPI
from .sh1107 import SH1107


class Display:
    """
    Text-mode display helper targeting a ~21x8 character grid.
    Logical buffer is 128x64; rotated 90° CW for physical 64x128 SH1107 panels.
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
        # rotation="cw" maps 128x64 logical buffer to 64x128 physical panel.
        # mirror=False to avoid mirrored characters; set True if hardware wiring mirrors X.
        self.dev = SH1107(width=128, height=64, spi=self.spi, dc=Pin(20), cs=Pin(17), rst=Pin(21), rotation="cw", mirror=False)
        self.clear()

    def clear(self) -> None:
        self.dev.fill(0)

    def text_at(self, row: int, col: int, s: str, color: int = 1) -> None:
        """
        Draw text aligned to the 8x8 builtin font grid.
        Row range: 0..7, Col range: 0..15 (16 chars)
        """
        if row < 0 or row > 7:
            return
        if col < 0:
            col = 0
        # Ensure we do not overflow visible width (128px / 8px = 16 chars)
        s = s[: max(0, 16 - col)]
        x = col * 8
        y = row * 8
        self.dev.text(s, x, y, color)

    def show(self) -> None:
        self.dev.show()


