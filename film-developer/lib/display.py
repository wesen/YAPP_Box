from machine import Pin, SPI
from .sh1107 import SH1107


class Display:
    """
    Text-mode display helper targeting a 16x8 character grid with the built-in 8x8 font.
    Note: Using a 6x8 font would allow 21x8; current implementation uses 8x8 (16 columns).
    Logical buffer is 128x64; rotated 90° CW for physical 64x128 SH1107 panels.
    """

    def __init__(self, spi_baudrate: int = 1_000_000) -> None:
        self.spi_baudrate = spi_baudrate
        # SPI0: SCK=GP18, MOSI=GP19 (no MISO used)
        self.spi = SPI(
            0,
            baudrate=self.spi_baudrate,
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

    def reset(self) -> None:
        self.dev.reset()

    def reset_spi(self, spi_baudrate: int | None = None) -> None:
        """
        Fully reset the SPI bus and reattach it to the SH1107 device.
        Optionally change baudrate at the same time.
        """
        try:
            self.spi.deinit()
        except Exception:
            # Some ports may not implement deinit; ignore.
            pass
        if spi_baudrate is not None:
            self.spi_baudrate = spi_baudrate
        # Recreate SPI0
        self.spi = SPI(
            0,
            baudrate=self.spi_baudrate,
            polarity=0,
            phase=0,
            sck=Pin(18),
            mosi=Pin(19),
        )
        # Reattach to device and keep current framebuffer contents
        self.dev.spi = self.spi


