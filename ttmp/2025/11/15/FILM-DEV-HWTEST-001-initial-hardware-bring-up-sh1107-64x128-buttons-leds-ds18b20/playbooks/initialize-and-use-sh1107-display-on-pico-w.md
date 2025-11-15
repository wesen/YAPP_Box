---
Title: Initialize and Use SH1107 Display on Pico W
Ticket: FILM-DEV-HWTEST-001
Status: active
Topics:
    - hardware
    - oled
    - sh1107
    - buttons
    - leds
    - temperature
    - pico
DocType: playbook
Intent: long-term
Owners:
    - manuel
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2025-11-15T18:29:03.278105234-05:00
---


# Initialize and Use SH1107 Display on Pico W

## Purpose

Bring up a 64x128 SH1107 OLED on Raspberry Pi Pico W using SPI0 and our driver in `lib/sh1107.py`, then render text via the `lib/display.py` helper. Covers rotation/mirroring for portrait-wired panels and basic diagnostics.

## Environment Assumptions

- MicroPython v1.26.1 on Raspberry Pi Pico W (others may work)
- Wiring per `README.md`:
  - SPI0: SCK=GP18, MOSI=GP19; CS=GP17; DC=GP20; RST=GP21
- Panel is 64x128 physically; we use a 128x64 logical buffer rotated 90 deg CW
- Deploy `film-developer/` as `main.py` + `lib/` to the board

## Commands

- Copy files and open REPL:

```bash
# Copy project to Pico (example using mpremote)
mpremote cp -r film-developer/* :

# Open REPL to watch logs
mpremote repl
```

## Initialization Sequence (Code)

- Use `rotation="cw"` to map 128x64 logical buffer to 64x128 physical panel.
- Flip `mirror` to True if characters appear mirrored horizontally.

```12:28:/home/manuel/code/others/YAPP_Box/film-developer/lib/display.py
        self.dev = SH1107(width=128, height=64, spi=self.spi, dc=Pin(20), cs=Pin(17), rst=Pin(21), rotation="cw", mirror=False)
```

Mirror control in the driver:

```12:24:/home/manuel/code/others/YAPP_Box/film-developer/lib/sh1107.py
        if self.mirror:
            self._write_cmd(_SEG_REMAP | 0x01)
        else:
            self._write_cmd(_SEG_REMAP | 0x00)
```

## Exit Criteria

- Text is readable left-to-right with 16 columns (8x8 font) and 8 rows.
- Pattern demo shows correct borders, checkerboard, and line spacing.
- No excessive flicker; `display.show()` only called when needed.

## Notes

- Diagnostics: use the pattern demo to verify orientation and addressing.

```1:36:/home/manuel/code/others/YAPP_Box/film-developer/lib/test_display.py
def run_pattern_demo(display, input_dev):
    ...
```

Minimal harness (if you want to run the demo from main):

```python
from lib.display import Display
from lib.input import Input
from lib.test_display import run_pattern_demo

def main():
    display = Display()
    input_dev = Input()
    run_pattern_demo(display, input_dev)
```

## Troubleshooting (EMI/Signal Integrity)

- Lower SPI baudrate: pass `spi_baudrate=200_000` to `Display()` and test stability.
- Keep wires short; route SCK/MOSI/CS/DC with a close ground return (twist with GND if possible).
- Add decoupling at OLED VCC: at least 0.1uF ceramic + 4.7–10uF bulk near the module.
- Add small series resistors (22–100 ohm) at SCK/MOSI near the Pico to tame ringing.
- Ensure solid common ground; avoid loose jumpers; reseat connectors.
