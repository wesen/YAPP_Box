---
Title: UI MVP Analysis and MicroPython References
Ticket: FILM-DEV-001
Status: active
Topics:
    - embedded
    - ui
    - hardware
    - timer
DocType: analysis
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Scope, constraints, and references for UI MVP on Pico W (OLED + buttons + temp display)
LastUpdated: 2025-11-08T17:15:37.429673131-05:00
---


# UI MVP Analysis and MicroPython References

## Purpose

Define a minimal, shippable UI subset for the Pico W-based film development timer to validate display readability and button navigation: splash, main status screen, configuration shell, and timer-running stub. Focus only on UI surfaces; no WiFi/CLI.

## Scope (UI MVP)
- Display: Waveshare SH1107 64x128 OLED (monochrome), text-only using 6x8 font
- Layout: 21 characters × 8 rows (rotated usage for readability)
- Screens:
  1) Splash/boot (1s)
  2) Main: shows loaded combo (e.g., TRX+2 D76 22°C) and Ready
  3) Config shell: Timer Presets (navigate only), System Info (stub)
  4) Timer Running stub: shows time total and remaining; STOP only
- Inputs: 3 buttons with internal pull-ups (GP14, GP15, GP16)
- LEDs: 3 LEDs (GP10, GP11, GP12) for simple status; optional in MVP

## Constraints
- Memory-limited environment (RP2040; prefer compact buffers and static strings)
- Keep redraws minimal: only update changed rows
- Sampling: DS18B20 ~1 Hz (or slower) to keep UX stable and reduce flicker
- Orientation: SH1107 is 64x128; map to 21×8 row/col grid with consistent origin

## ASCII UI (target grid 21×8)

Splash (1s):
```
┌─────────────────────┐
│   PICO TEMP MON     │
│                     │
│   Booting...        │
│                     │
│                     │
│                     │
│                     │
│[    ][ MENU][    ]  │
└─────────────────────┘
```

Main:
```
┌─────────────────────┐
│   FILM DEV TIMER    │
│                     │
│ TRX+2 D76 22°C      │
│ Timer: 14:30        │
│ Temp: 22.1°C ✓      │
│                     │
│[    ][ MENU][    ]  │
└─────────────────────┘
```

Config (shell):
```
┌─────────────────────┐
│       MENU          │
│                     │
│ ► Timer Presets     │
│   System Info       │
│   Back              │
│                     │
│[ SEL][ NAV][ BACK]  │
└─────────────────────┘
```

## Wiring
- DS18B20 DATA: GP22 with 4.7kΩ pull-up to 3V3
- OLED SPI: CS=GP17, CLK=GP18, DIN=GP19, DC=GP20, RST=GP21, VCC=3V3, GND
- Buttons: GP14/15/16 → GND when pressed (Pin.PULL_UP)
- LEDs: GP10/11/12 → 220Ω → LED anode → GND

## MicroPython APIs (authoritative references)
- machine.Pin — GPIO, pull-ups, IRQ
  - https://docs.micropython.org/en/latest/library/machine.Pin.html
- machine.SPI — SPI for SH1107
  - https://docs.micropython.org/en/latest/library/machine.SPI.html
- framebuf — FrameBuffer for text rendering
  - https://docs.micropython.org/en/latest/library/framebuf.html
- onewire/ds18x20 — DS18B20 temperature sensor
  - https://docs.micropython.org/en/latest/library/onewire.html
  - https://docs.micropython.org/en/latest/library/ds18x20.html

## Third-party driver reference
- SH1107 MicroPython driver (e.g., Waveshare/ports/libs)
  - Example community driver: https://github.com/robert-hh/SH1107-SSD1306 (verify pin mapping)

## Risks / Open Questions
- SH1107 coordinate system/rotation differences vs SSD1306
- Font size/spacing on 64×128: confirm 6×8 mapping to 21×8 characters
- Debounce strategy: 50 ms is typical; adjust based on hardware
- Timer screen layout fits within 21×8 without truncation

## Next Steps → Design
- Lock grid and coordinates per screen (Main, Config shell, Timer stub)
- Define UI controller API (render header/body/footer; invalidate rows)
- Define button event model (short press only for MVP)
- Draft minimal modules layout (ui.py, input.py, display.py, temp.py, main.py)
