---
Title: UI MVP Design and Implementation Plan
Ticket: FILM-DEV-001
Status: active
Topics:
    - embedded
    - ui
    - hardware
    - timer
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Minimal UI plan for film dev timer (text UI 16×8 with 8×8 font; 21×8 possible with 6×8 font)
LastUpdated: 2025-11-15T20:37:07.980942745-05:00
---




# UI MVP Design and Implementation Plan

## Executive Summary

Deliver a minimal UI for the Pico W film development timer: splash → main status screen → config shell → timer-running stub. Implement display driver wiring, a small UI controller on a 21×8 text grid, and button navigation with debouncing. Temperature reading is displayed live (DS18B20); no persistence or WiFi.

## Problem Statement

We need a working on-device UI to validate enclosure ergonomics, display readability, and button layout. The UI should run on MicroPython, draw text to the SH1107 OLED, and respond to 3 buttons. This MVP de-risks hardware and establishes the baseline code structure.

## Proposed Solution

Implement a micro UI stack:
- `display.py`: Initialize SPI + SH1107 driver and expose a line-based text API on a 21×8 grid using `framebuf` text.
- `input.py`: Button handling with internal pull-ups + 50 ms debounce; emit simple events (BTN1, BTN2, BTN3).
- `temp.py`: Poll DS18B20 ~1 Hz and cache latest reading.
- `ui.py`: Screen registry (Splash, Main, Menu), render functions, and event routing.
- `main.py`: System init, loop scheduling (cooperative), and state transitions.

## Design Decisions

- Text-only UI (6×8 font) for simplicity and speed
- Fixed grid (21×8) to guarantee fit on SH1107 64×128 with rotation
- Debounce in software (50 ms) to avoid extra components
- One-second temperature sampling cadence to prevent flicker
- No persistence/WiFi for MVP; keep boot time short

## Alternatives Considered

- Graphics-based UI: deferred; text-only is sufficient now
- Interrupt-driven UI loop: simple polling is easier and reliable for MVP
- SSD1306 I2C displays: SH1107 SPI chosen to match available module and wiring

## Implementation Plan

### File Structure
```
/
├── main.py
├── lib/
│   ├── display.py
│   ├── input.py
│   ├── temp.py
│   └── ui.py
```

### Pseudocode

display.py
```
from machine import Pin, SPI
from sh1107 import SH1107

class Display:
    def __init__(self):
        self.spi = SPI(0, baudrate=1_000_000, polarity=0, phase=0,
                       sck=Pin(18), mosi=Pin(19))
        self.dev = SH1107(64, 128, self.spi, Pin(20), Pin(17), Pin(21))
        self.clear()

    def clear(self):
        self.dev.fill(0)

    def text_at(self, row, col, s):
        x = col * 6
        y = row * 8
        self.dev.text(s[:21], x, y, 1)

    def show(self):
        self.dev.show()
```

input.py
```
from machine import Pin
import time

BTN_PINS = (14, 15, 16)

class Input:
    def __init__(self, debounce_ms=50):
        self.btns = [Pin(p, Pin.IN, Pin.PULL_UP) for p in BTN_PINS]
        self.last = [1, 1, 1]
        self.t_last = [0, 0, 0]
        self.debounce = debounce_ms

    def read(self):
        now = time.ticks_ms()
        events = []
        for i, b in enumerate(self.btns):
            v = b.value()
            if v != self.last[i] and time.ticks_diff(now, self.t_last[i]) > self.debounce:
                self.t_last[i] = now
                self.last[i] = v
                if v == 0:
                    events.append(i+1)  # BTN1/2/3
        return events
```

temp.py
```
import onewire, ds18x20, time
from machine import Pin

class Temp:
    def __init__(self, pin=22):
        self.ds = ds18x20.DS18X20(onewire.OneWire(Pin(pin)))
        self.roms = self.ds.scan()
        self.last_c = None

    def read_c(self):
        if not self.roms:
            return None
        self.ds.convert_temp(); time.sleep_ms(750)
        c = self.ds.read_temp(self.roms[0])
        self.last_c = c
        return c
```

ui.py
```
class UI:
    def __init__(self, display, input, temp):
        self.d, self.i, self.t = display, input, temp
        self.screen = 'splash'
        self.t0 = time.ticks_ms()

    def render(self):
        self.d.clear()
        if self.screen == 'splash':
            self.d.text_at(0, 2, 'PICO TEMP MON')
            self.d.text_at(2, 2, 'Booting...')
            if time.ticks_diff(time.ticks_ms(), self.t0) > 1000:
                self.screen = 'main'
        elif self.screen == 'main':
            c = self.t.read_c()
            temp_line = 'Temp: --.-C' if c is None else f'Temp: {c:.1f}C'
            self.d.text_at(0, 1, 'FILM DEV TIMER')
            self.d.text_at(2, 0, 'TRX+2 D76 22C')
            self.d.text_at(3, 0, 'Timer: 14:30')
            self.d.text_at(4, 0, temp_line)
            self.d.text_at(6, 0, '[    ][ MENU][    ]')
        elif self.screen == 'menu':
            self.d.text_at(0, 4, 'MENU')
            self.d.text_at(2, 0, '\u25BA Timer Presets')
            self.d.text_at(3, 0, '  System Info')
            self.d.text_at(6, 0, '[ SEL][ NAV][ BACK]')
        elif self.screen == 'timer':
            self.d.text_at(0, 3, 'DEVELOPING...')
            self.d.text_at(2, 0, 'TRX+2 D76 22C')
            self.d.text_at(3, 0, '14:30 [#####     ]')
            self.d.text_at(4, 0, '12:30 remaining')
            self.d.text_at(6, 0, '[    ][    ][STOP]')
        self.d.show()

    def handle(self):
        for e in self.i.read():
            if self.screen == 'main' and e == 2:  # BTN2 = MENU
                self.screen = 'menu'
            elif self.screen == 'menu' and e == 3:  # BTN3 = BACK
                self.screen = 'main'
            elif self.screen == 'main' and e == 1:  # BTN1 = start timer stub
                self.screen = 'timer'
            elif self.screen == 'timer' and e == 3:  # BTN3 = STOP
                self.screen = 'main'
```

main.py
```
from lib.display import Display
from lib.input import Input
from lib.temp import Temp
from lib.ui import UI
import time

d = Display(); i = Input(); t = Temp(); ui = UI(d, i, t)
while True:
    ui.handle()
    ui.render()
    time.sleep_ms(50)
```

### Test Plan
- Power-on shows Splash for ~1s, then Main
- Sensor disconnected: shows Temp: --.-C (no crash)
- BTN2 opens Menu; BTN3 returns to Main
- BTN1 from Main goes to timer stub; BTN3 stops/returns to Main
- Temperature value updates within 1–2 seconds

## Film Dev Timer UI Subset (MVP)

Screens:
- Splash → Main → (Menu) → Timer stub → Main

Buttons:
- BTN1: Start (from Main), no-op elsewhere (except future pause)
- BTN2: MENU (from Main)
- BTN3: BACK/STOP (contextual)

Rendering:
- 21×8 char grid; truncate strings beyond 21 chars; avoid multi-line wraps

## Open Questions

- Confirm SH1107 rotation for 21×8 text grid
- Confirm button/LED GPIOs match enclosure harness
- Decide Units toggle behavior (persist or session-only?)

## References

- Analysis: UI MVP Analysis and MicroPython References
- MicroPython docs: machine.Pin, machine.SPI, framebuf, onewire, ds18x20
- Enclosure README for wiring and constraints
