from micropython import const
import framebuf
from machine import Pin, SPI


_SET_DISP = const(0xAE)
_SET_DISP_ON = const(0xAF)
_SET_COL_LOW = const(0x00)     # lower nibble of column start address
_SET_COL_HIGH = const(0x10)    # higher nibble of column start address
_SET_PAGE = const(0xB0)        # page start address (0..7)
_SET_START_LINE = const(0xDC)  # set display start line (0..63)
_SET_CONTRAST = const(0x81)
_SEG_REMAP = const(0xA0)       # segment remap
_SET_ENTIRE_ON = const(0xA4)
_SET_NORM_INV = const(0xA6)    # normal/inverse display
_SET_MUX_RATIO = const(0xA8)
_SET_DISP_OFFSET = const(0xD3)
_SET_DISP_CLK_DIV = const(0xD5)
_SET_PRECHARGE = const(0xD9)
_SET_VCOM_DESEL = const(0xDB)
_SET_CHARGE_PUMP = const(0x8D)


class SH1107:
    """
    Minimal SH1107 SPI driver sufficient for text-mode rendering via framebuf.
    Target: 128x64 panels (W=128, H=64).
    """

    def __init__(self, width: int, height: int, spi: SPI, dc: Pin, cs: Pin, rst: Pin) -> None:
        self.width = width
        self.height = height
        self.spi = spi
        self.dc = dc
        self.cs = cs
        self.rst = rst

        self.dc.init(Pin.OUT, value=0)
        self.cs.init(Pin.OUT, value=1)
        self.rst.init(Pin.OUT, value=1)

        self.buffer = bytearray(self.width * self.height // 8)
        self.fb = framebuf.FrameBuffer(self.buffer, self.width, self.height, framebuf.MONO_VLSB)

        self._reset()
        self._init()
        self.fill(0)
        self.show()

    # FrameBuffer API passthrough
    def fill(self, c: int) -> None:
        self.fb.fill(c)

    def text(self, s: str, x: int, y: int, c: int = 1) -> None:
        self.fb.text(s, x, y, c)

    def pixel(self, x: int, y: int, c: int = 1) -> None:
        self.fb.pixel(x, y, c)

    def line(self, x0: int, y0: int, x1: int, y1: int, c: int = 1) -> None:
        self.fb.line(x0, y0, x1, y1, c)

    def rect(self, x: int, y: int, w: int, h: int, c: int = 1) -> None:
        self.fb.rect(x, y, w, h, c)

    def fill_rect(self, x: int, y: int, w: int, h: int, c: int = 1) -> None:
        self.fb.fill_rect(x, y, w, h, c)

    def show(self) -> None:
        # SH1107 expects pages (8px tall), 0..7
        for page in range(self.height // 8):
            self._write_cmd(_SET_PAGE | page)
            # Column 0..127; many SH1107 modules use column offset 0
            self._write_cmd(_SET_COL_LOW | 0x00)
            self._write_cmd(_SET_COL_HIGH | 0x00)
            start = page * self.width
            end = start + self.width
            self._write_data(self.buffer[start:end])

    # Low-level
    def _reset(self) -> None:
        self.rst.value(0)
        for _ in range(1000):
            pass
        self.rst.value(1)

    def _init(self) -> None:
        self._write_cmd(_SET_DISP)  # display off
        self._write_cmd(_SET_START_LINE)
        self._write_cmd(0x00)  # start line = 0
        self._write_cmd(_SET_DISP_OFFSET)
        self._write_cmd(0x60)  # typical offset for SH1107 64px height
        self._write_cmd(_SET_NORM_INV)  # normal display
        self._write_cmd(_SET_ENTIRE_ON)  # follow RAM content
        self._write_cmd(_SEG_REMAP | 0x01)  # segment remap
        self._write_cmd(_SET_MUX_RATIO)
        self._write_cmd(0x3F)  # 1/64 duty
        self._write_cmd(_SET_DISP_CLK_DIV)
        self._write_cmd(0x51)  # osc freq/div
        self._write_cmd(_SET_PRECHARGE)
        self._write_cmd(0x22)
        self._write_cmd(_SET_VCOM_DESEL)
        self._write_cmd(0x35)
        self._write_cmd(_SET_CONTRAST)
        self._write_cmd(0x5F)
        self._write_cmd(_SET_DISP_ON)  # display on

    def _write_cmd(self, cmd: int) -> None:
        self.cs.value(0)
        self.dc.value(0)
        self.spi.write(bytearray([cmd & 0xFF]))
        self.cs.value(1)

    def _write_data(self, buf: bytes) -> None:
        self.cs.value(0)
        self.dc.value(1)
        self.spi.write(buf)
        self.cs.value(1)


