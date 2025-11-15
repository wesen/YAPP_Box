from lib.display import Display
from lib.input import Input
import time
from machine import Pin


def render_pattern(display: Display, pattern_idx: int) -> str:
    """
    Draw one of several diagnostic patterns to verify orientation and width.
    Returns the pattern name for logging.
    """
    dev = display.dev
    name = ""
    display.clear()
    if pattern_idx == 0:
        name = "Text grid 0..F"
        display.text_at(0, 0, "0123456789ABCDEF")
        display.text_at(2, 0, "Cols:16 Rows:8")
        display.text_at(4, 0, "Top-Left origin")
        display.text_at(6, 0, "BTN1: next test")
    elif pattern_idx == 1:
        name = "Vertical lines (8px)"
        # Vertical lines at every 8px
        for x in range(0, 128, 8):
            dev.line(x, 0, x, 63, 1)
        display.text_at(7, 0, "V lines 0..120 step 8")
    elif pattern_idx == 2:
        name = "Horizontal lines (8px)"
        # Horizontal lines at every 8px
        for y in range(0, 64, 8):
            dev.line(0, y, 127, y, 1)
        display.text_at(7, 0, "H lines 0..56 step 8")
    elif pattern_idx == 3:
        name = "Checkerboard 8x8"
        # 8x8 checkerboard
        for y in range(0, 64, 8):
            for x in range(0, 128, 8):
                if ((x // 8) ^ (y // 8)) & 1:
                    dev.fill_rect(x, y, 8, 8, 1)
        display.text_at(7, 0, "Checker 8x8")
    elif pattern_idx == 4:
        name = "Borders + corners"
        dev.rect(0, 0, 128, 64, 1)
        # Corner pixels
        dev.pixel(0, 0, 1)
        dev.pixel(127, 0, 1)
        dev.pixel(0, 63, 1)
        dev.pixel(127, 63, 1)
        display.text_at(6, 0, "Rect (0,0)-(127,63)")
    else:
        name = "Text edges"
        display.text_at(0, 0, "Start at (0,0)")
        display.text_at(1, 8, "Col 8")
        display.text_at(2, 15, "15")
        display.text_at(4, 0, "End cols ->")
        display.text_at(5, 0, "0123456789ABCDEF")
        display.text_at(7, 0, "BTN1 cycles tests")
    t0 = time.ticks_ms()
    display.show()
    dt = time.ticks_diff(time.ticks_ms(), t0)
    print("Rendered pattern:", name, "in", dt, "ms")
    return name


def main() -> None:
    # Hardware test: OLED patterns + button→LED mapping
    display = Display()
    input_dev = Input()

    # LEDs on GP10/11/12 per README pin map
    LED_PINS = (10, 11, 12)
    leds = [Pin(p, Pin.OUT, value=0) for p in LED_PINS]
    led_state = [0, 0, 0]

    # Initial diagnostics
    print("Starting OLED diagnostics")
    print("Logical buffer: 128x64; rotation=cw; mirror=False")
    print("Grid: 16 columns x 8 rows (8x8 font)")

    # Start with pattern 0
    pattern_idx = 0
    render_pattern(display, pattern_idx)

    while True:
        # Toggle corresponding LED on button press (active-low edge)
        for e in input_dev.read():
            print("Button event:", e)
            idx = e - 1
            if 0 <= idx < 3:
                led_state[idx] ^= 1
                leds[idx].value(led_state[idx])
            # BTN1: switch to next pattern
            if e == 1:
                pattern_idx = (pattern_idx + 1) % 6
                # Explicitly clear and push to display between tests
                print("Clearing screen before next pattern")
                display.clear()
                display.show()
                time.sleep_ms(20)
                print("Switching to pattern index:", pattern_idx)
                render_pattern(display, pattern_idx)
        time.sleep_ms(20)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        pass


