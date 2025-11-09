# Film Developer UI MVP (Pico W)

Minimal on-device UI for a film development timer on Raspberry Pi Pico W with:
- 128×64 SH1107 OLED (SPI)
- 3 buttons (with internal pull-ups)
- DS18B20 temperature probe

## Quick Start

```text
film-developer/
  main.py
  lib/
    display.py   # 21×8 text grid wrapper
    sh1107.py    # minimal SH1107 SPI driver
    input.py     # 3-button debounce/events
    temp.py      # DS18B20 helper
    ui.py        # Splash/Main/Menu/Timer stub
```

Run loop entrypoint:
```bash
python3 film-developer/main.py
```
(On-device you’ll deploy these as `main.py` + `lib/` to MicroPython.)

## Wiring (Pins)

Follow the full pin map in the ticket reference doc `01-application-architecture-and-system-overview.md`. The essentials are summarized below.

### Power (USB)
- Power the Pico W over USB (VBUS → onboard 3V3 regulator). Do not feed 5V into 3V3.
- Power peripherals from 3V3 Out (≤ ~300 mA recommended).
- Common GND is required across OLED, DS18B20, and all buttons.

### OLED (SH1107, SPI0)
- VCC → 3V3
- GND → GND
- SCK → GP18
- MOSI (DIN) → GP19
- CS → GP17
- DC → GP20
- RST → GP21

### Temperature Sensor (DS18B20)
- DATA → GP22 (with 4.7kΩ pull‑up to 3V3)
- VCC → 3V3
- GND → GND

### Buttons (3×, Internal Pull‑Up)
Active‑low wiring: each button connects its GPIO to GND when pressed. No external pull‑ups needed (enabled in software).

| Button | GPIO | Electrical |
|--------|------|------------|
| BTN1   | GP14 | GP14 → (button) → GND, internal pull‑up enabled |
| BTN2   | GP15 | GP15 → (button) → GND, internal pull‑up enabled |
| BTN3   | GP16 | GP16 → (button) → GND, internal pull‑up enabled |

Notes:
- All three buttons share common GND.
- Debounce is handled in `lib/input.py` (50 ms default).

## Button Roles (MVP)
- BTN1: Start/Select (contextual)
- BTN2: Menu/Navigate (contextual)
- BTN3: Back/Stop (contextual)

## Full Wiring Diagram (ASCII)

Orientation: Pico W with micro‑USB at the top.

```text
USB 5V → VBUS (onboard regulator) → 3V3 Out ─┬─> OLED VCC
                                             ├─> DS18B20 VCC
                                             └─> (optional) other 3V3 peripherals

GND (common) ─────────────────────────────────┼─> OLED GND
                                              ├─> DS18B20 GND
                                              └─> BTN1/BTN2/BTN3 to GND when pressed

Pico W GPIOs:
  GP18 (SPI0 SCK)  ────────────────→ OLED SCK
  GP19 (SPI0 MOSI) ────────────────→ OLED DIN
  GP17 (SPI0 CSn)  ────────────────→ OLED CS
  GP20 (DC)        ────────────────→ OLED DC
  GP21 (RST)       ────────────────→ OLED RST

  GP22 ──┬────────→ DS18B20 DATA
         └─[4.7k]─→ 3V3   (mandatory pull‑up)

  GP14 ──/ BTN1 /──→ GND  (internal pull‑up enabled in software)
  GP15 ──/ BTN2 /──→ GND
  GP16 ──/ BTN3 /──→ GND
```

Notes:
  - Buttons are active‑low; internal pull‑ups are enabled in `lib/input.py`.
  - DS18B20 requires the 4.7kΩ pull‑up from DATA to 3V3.

## Pico W Pinout Diagram (ASCII)

USB port at top. Only pins used in this project are labeled with their function.

```text
                    ┌─────────────┐
                    │   USB Port  │
                    └─────────────┘
         ┌───────────────────────────────┐
         │  Raspberry Pi Pico W          │
         │                               │
    GP0  │ 1 ●                      ● 40 │ VBUS (5V USB)
    GP1  │ 2 ●                      ● 39 │ VSYS
    GND  │ 3 ●                      ● 38 │ GND
    GP2  │ 4 ●                      ● 37 │ 3V3_EN
    GP3  │ 5 ●                      ● 36 │ 3V3 Out ──→ OLED, DS18B20
    GP4  │ 6 ●                      ● 35 │ ADC_VREF
    GP5  │ 7 ●                      ● 34 │ GP28
    GND  │ 8 ●                      ● 33 │ GND ──→ OLED, DS18B20, BTNs
    GP6  │ 9 ●                      ● 32 │ GP27
    GP7  │10 ●                      ● 31 │ GP26
    GP8  │11 ●                      ● 30 │ RUN
    GP9  │12 ●                      ● 29 │ GP22 ──→ DS18B20 DATA (+4.7kΩ)
    GND  │13 ●                      ● 28 │ GND
GP10/LED1│14 ●                      ● 27 │ GP21 ──→ OLED RST
GP11/LED2│15 ●                      ● 26 │ GP20 ──→ OLED DC
GP12/LED3│16 ●                      ● 25 │ GP19 ──→ OLED DIN (MOSI)
   GP13  │17 ●                      ● 24 │ GP18 ──→ OLED SCK
    GND  │18 ●                      ● 23 │ GND
GP14/BTN1│19 ●                      ● 22 │ GP17 ──→ OLED CS
GP15/BTN2│20 ●                      ● 21 │ GP16/BTN3
         │                               │
         └───────────────────────────────┘

Legend:
  BTN1/2/3: Buttons to GND (internal pull-ups)
  LED1/2/3: LEDs with 220Ω resistors to GND (PWM capable)
  OLED:     SH1107 128×64 on SPI0
  DS18B20:  Temperature sensor on 1-Wire (GP22 + 4.7kΩ pull-up to 3V3)
```

## Full Pinout (Reference)
- Interactive Pico W pinout: https://picow.pinout.xyz/
- Official Pico W A4 pinout PDF: https://datasheets.raspberrypi.com/picow/PicoW-A4-Pinout.pdf
  - Use these diagrams to confirm orientation and alternate functions (SPI/I2C/UART/PWM).

## References
- Ticket: `FILM-DEV-001` (see `reference/01-application-architecture-and-system-overview.md`)
- UI spec: `reference/02-user-interface-specification-and-display-system.md`
- Timer system: `reference/03-timer-system-and-development-process-management.md`

