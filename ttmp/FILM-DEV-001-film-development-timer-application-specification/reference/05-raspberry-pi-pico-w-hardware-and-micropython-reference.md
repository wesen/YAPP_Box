---
Title: Raspberry Pi Pico W Hardware and MicroPython Reference
Ticket: FILM-DEV-001
Status: active
Topics:
    - embedded
    - hardware
    - micropython
    - reference
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources:
    - https://picow.pinout.xyz/
    - https://datasheets.raspberrypi.com/picow/PicoW-A4-Pinout.pdf
    - https://datasheets.raspberrypi.com/picow/pico-w-datasheet.pdf
    - https://docs.micropython.org/en/latest/rp2/quickref.html
    - https://micropython.org/download/RPI_PICO_W/
Summary: Complete hardware pinout and MicroPython programming reference for Raspberry Pi Pico W
LastUpdated: 2025-11-08
---

# Raspberry Pi Pico W Hardware and MicroPython Reference

## Overview

This document provides a comprehensive hardware and software reference for the Raspberry Pi Pico W microcontroller board. Use this for GPIO planning, alternate function lookup, MicroPython API reference, and troubleshooting.

**Board**: Raspberry Pi Pico W (RP2040 + CYW43439 wireless)  
**Microcontroller**: RP2040 (dual-core ARM Cortex-M0+)  
**GPIO Count**: 26 user-accessible GPIO pins (GP0-GP22, GP26-GP28)  
**Special**: 3 additional GPIOs reserved for wireless (GP23-GP25, GP29)

## Pin Orientation

All pin numbers and diagrams assume **USB port at the top** (standard orientation).

```text
                    ┌─────────────┐
                    │   USB Port  │
                    └─────────────┘
         ┌───────────────────────────────┐
         │  Raspberry Pi Pico W          │
         │                               │
         │ 1 ●                      ● 40 │
         │ 2 ●                      ● 39 │
         │ 3 ●                      ● 38 │
         │...                        ...│
         │20 ●                      ● 21 │
         │                               │
         └───────────────────────────────┘
```

## Complete Pin Table

| Pin | Name      | GPIO | Default Pull | SPI         | I2C      | UART       | PWM    | ADC | Notes |
|-----|-----------|------|--------------|-------------|----------|------------|--------|-----|-------|
| 1   | GP0       | 0    | Down         | SPI0 RX     | I2C0 SDA | UART0 TX   | PWM0 A | -   | |
| 2   | GP1       | 1    | Down         | SPI0 CSn    | I2C0 SCL | UART0 RX   | PWM0 B | -   | |
| 3   | GND       | -    | -            | -           | -        | -          | -      | -   | Ground |
| 4   | GP2       | 2    | Down         | SPI0 SCK    | I2C1 SDA | UART0 CTS  | PWM1 A | -   | |
| 5   | GP3       | 3    | Down         | SPI0 TX     | I2C1 SCL | UART0 RTS  | PWM1 B | -   | |
| 6   | GP4       | 4    | Down         | SPI0 RX     | I2C0 SDA | UART1 TX   | PWM2 A | -   | |
| 7   | GP5       | 5    | Down         | SPI0 CSn    | I2C0 SCL | UART1 RX   | PWM2 B | -   | |
| 8   | GND       | -    | -            | -           | -        | -          | -      | -   | Ground |
| 9   | GP6       | 6    | Down         | SPI0 SCK    | I2C1 SDA | UART1 CTS  | PWM3 A | -   | |
| 10  | GP7       | 7    | Down         | SPI0 TX     | I2C1 SCL | UART1 RTS  | PWM3 B | -   | |
| 11  | GP8       | 8    | Down         | SPI1 RX     | I2C0 SDA | UART1 TX   | PWM4 A | -   | |
| 12  | GP9       | 9    | Down         | SPI1 CSn    | I2C0 SCL | UART1 RX   | PWM4 B | -   | |
| 13  | GND       | -    | -            | -           | -        | -          | -      | -   | Ground |
| 14  | GP10      | 10   | Down         | SPI1 SCK    | I2C1 SDA | UART1 CTS  | PWM5 A | -   | |
| 15  | GP11      | 11   | Down         | SPI1 TX     | I2C1 SCL | UART1 RTS  | PWM5 B | -   | |
| 16  | GP12      | 12   | Down         | SPI1 RX     | I2C0 SDA | UART0 TX   | PWM6 A | -   | |
| 17  | GP13      | 13   | Down         | SPI1 CSn    | I2C0 SCL | UART0 RX   | PWM6 B | -   | |
| 18  | GND       | -    | -            | -           | -        | -          | -      | -   | Ground |
| 19  | GP14      | 14   | Down         | SPI1 SCK    | I2C1 SDA | UART0 CTS  | PWM7 A | -   | |
| 20  | GP15      | 15   | Down         | SPI1 TX     | I2C1 SCL | UART0 RTS  | PWM7 B | -   | |
| 21  | GP16      | 16   | Down         | SPI0 RX     | I2C0 SDA | UART0 TX   | PWM0 A | -   | |
| 22  | GP17      | 17   | Down         | SPI0 CSn    | I2C0 SCL | UART0 RX   | PWM0 B | -   | |
| 23  | GND       | -    | -            | -           | -        | -          | -      | -   | Ground |
| 24  | GP18      | 18   | Down         | SPI0 SCK    | I2C1 SDA | UART0 CTS  | PWM1 A | -   | |
| 25  | GP19      | 19   | Down         | SPI0 TX     | I2C1 SCL | UART0 RTS  | PWM1 B | -   | |
| 26  | GP20      | 20   | Down         | SPI0 RX     | I2C0 SDA | UART1 TX   | PWM2 A | -   | |
| 27  | GP21      | 21   | Down         | SPI0 CSn    | I2C0 SCL | UART1 RX   | PWM2 B | -   | |
| 28  | GND       | -    | -            | -           | -        | -          | -      | -   | Ground |
| 29  | GP22      | 22   | Down         | SPI0 SCK    | I2C1 SDA | UART1 CTS  | PWM3 A | -   | |
| 30  | RUN       | -    | Up (ext)     | -           | -        | -          | -      | -   | Reset (active low) |
| 31  | GP26      | 26   | Down         | SPI1 SCK    | I2C1 SDA | UART1 CTS  | PWM5 A | 0   | ADC0 |
| 32  | GP27      | 27   | Down         | SPI1 TX     | I2C1 SCL | UART1 RTS  | PWM5 B | 1   | ADC1 |
| 33  | GND       | -    | -            | -           | -        | -          | -      | -   | Ground (ADC) |
| 34  | GP28      | 28   | Down         | SPI1 RX     | I2C0 SDA | UART0 TX   | PWM6 A | 2   | ADC2 |
| 35  | ADC_VREF  | -    | -            | -           | -        | -          | -      | -   | ADC reference (optional) |
| 36  | 3V3(OUT)  | -    | -            | -           | -        | -          | -      | -   | 3.3V output (max 300mA) |
| 37  | 3V3_EN    | -    | -            | -           | -        | -          | -      | -   | 3.3V enable (pull low to disable) |
| 38  | GND       | -    | -            | -           | -        | -          | -      | -   | Ground |
| 39  | VSYS      | -    | -            | -           | -        | -          | -      | -   | System voltage (1.8-5.5V in, 5V out from USB) |
| 40  | VBUS      | -    | -            | -           | -        | -          | -      | -   | USB 5V input (when connected) |

## Reserved/Special Pins (Not User-Accessible)

| GPIO  | Function                | Notes |
|-------|-------------------------|-------|
| GP23  | Wireless Power On       | CYW43439 control |
| GP24  | Wireless SPI Data       | CYW43439 SPI MISO/MOSI |
| GP25  | Wireless SPI CS         | CYW43439 chip select; also connected to onboard LED |
| GP29  | Wireless SPI Clock / VSYS Sense | CYW43439 SPI CLK; can read VSYS/3 via ADC3 |

**Important**: GP25 controls the onboard LED on Pico W (via CYW43439), unlike Pico (non-W) where it's directly connected.

## Power Pins

| Pin | Name      | Voltage | Max Current | Notes |
|-----|-----------|---------|-------------|-------|
| 40  | VBUS      | 5V      | -           | USB 5V input (only when USB connected) |
| 39  | VSYS      | 1.8-5.5V| -           | System input (when powered externally) or 5V output (from USB via diode) |
| 36  | 3V3(OUT)  | 3.3V    | 300mA       | Regulated 3.3V output from RT6150B-33GQW |
| 37  | 3V3_EN    | -       | -           | Pull low to disable 3.3V regulator |
| 35  | ADC_VREF  | 3.3V    | -           | Optional external ADC reference (leave floating for internal) |
| 3,8,13,18,23,28,33,38 | GND | 0V | - | Ground (all connected internally) |

### Power Notes
- **USB powered**: VBUS (5V) → VSYS (via diode) → 3V3 regulator → 3V3(OUT)
- **Battery powered**: Connect to VSYS (1.8-5.5V) → 3V3 regulator → 3V3(OUT)
- **3V3(OUT) load**: Keep total load under 300mA (includes RP2040 + CYW43439 + peripherals)
- **VBUS sensing**: Use GP24 (WL_GPIO2) or read VSYS via GP29/ADC3 to detect USB power

## Alternate Function Groups

### SPI (Serial Peripheral Interface)

**SPI0** (can be mapped to multiple pin sets):
- **Set 1**: GP16(RX), GP17(CSn), GP18(SCK), GP19(TX)
- **Set 2**: GP0(RX), GP1(CSn), GP2(SCK), GP3(TX)
- **Set 3**: GP4(RX), GP5(CSn), GP6(SCK), GP7(TX)
- **Set 4**: GP20(RX), GP21(CSn), GP22(SCK) [partial]

**SPI1** (can be mapped to multiple pin sets):
- **Set 1**: GP8(RX), GP9(CSn), GP10(SCK), GP11(TX)
- **Set 2**: GP12(RX), GP13(CSn), GP14(SCK), GP15(TX)
- **Set 3**: GP26(SCK), GP27(TX), GP28(RX) [partial]

**Notes**:
- Each SPI peripheral can only use one set at a time
- RX = MISO, TX = MOSI in standard terminology
- CSn can be managed as GPIO for multiple devices

### I2C (Inter-Integrated Circuit)

**I2C0**:
- **Set 1**: GP0(SDA), GP1(SCL)
- **Set 2**: GP4(SDA), GP5(SCL)
- **Set 3**: GP8(SDA), GP9(SCL)
- **Set 4**: GP12(SDA), GP13(SCL)
- **Set 5**: GP16(SDA), GP17(SCL)
- **Set 6**: GP20(SDA), GP21(SCL)
- **Set 7**: GP28(SDA) [partial]

**I2C1**:
- **Set 1**: GP2(SDA), GP3(SCL)
- **Set 2**: GP6(SDA), GP7(SCL)
- **Set 3**: GP10(SDA), GP11(SCL)
- **Set 4**: GP14(SDA), GP15(SCL)
- **Set 5**: GP18(SDA), GP19(SCL)
- **Set 6**: GP22(SDA) [partial]
- **Set 7**: GP26(SDA), GP27(SCL)

**Notes**:
- Both I2C0 and I2C1 support standard (100kHz) and fast (400kHz) modes
- Internal pull-ups available but external 4.7kΩ recommended for reliability

### UART (Universal Asynchronous Receiver-Transmitter)

**UART0**:
- **Set 1**: GP0(TX), GP1(RX), GP2(CTS), GP3(RTS)
- **Set 2**: GP12(TX), GP13(RX), GP14(CTS), GP15(RTS)
- **Set 3**: GP16(TX), GP17(RX), GP18(CTS), GP19(RTS)
- **Set 4**: GP28(TX) [partial]

**UART1**:
- **Set 1**: GP4(TX), GP5(RX), GP6(CTS), GP7(RTS)
- **Set 2**: GP8(TX), GP9(RX), GP10(CTS), GP11(RTS)
- **Set 3**: GP20(TX), GP21(RX), GP22(CTS) [partial]
- **Set 4**: GP24(TX), GP25(RX) [reserved for wireless]

**Notes**:
- CTS/RTS are optional hardware flow control
- Default baud rates: 115200 typical, up to 3Mbps supported
- UART0 is often used for USB CDC (serial over USB)

### PWM (Pulse Width Modulation)

All GPIO pins support PWM. There are 8 PWM slices (0-7), each with 2 channels (A/B):

| Slice | Channel A | Channel B |
|-------|-----------|-----------|
| PWM0  | GP0, GP16 | GP1, GP17 |
| PWM1  | GP2, GP18 | GP3, GP19 |
| PWM2  | GP4, GP20 | GP5, GP21 |
| PWM3  | GP6, GP22 | GP7       |
| PWM4  | GP8       | GP9       |
| PWM5  | GP10, GP26| GP11, GP27|
| PWM6  | GP12, GP28| GP13      |
| PWM7  | GP14      | GP15      |

**Notes**:
- Each slice shares a counter; both channels (A/B) have the same frequency
- Channels A and B can have independent duty cycles
- Typical uses: LED dimming, servo control, audio (via low-pass filter)

### ADC (Analog-to-Digital Converter)

| ADC Channel | GPIO  | Pin | Notes |
|-------------|-------|-----|-------|
| ADC0        | GP26  | 31  | 12-bit, 0-3.3V |
| ADC1        | GP27  | 32  | 12-bit, 0-3.3V |
| ADC2        | GP28  | 34  | 12-bit, 0-3.3V |
| ADC3        | GP29  | -   | Internal: VSYS/3 sense (reserved) |
| ADC4        | -     | -   | Internal: temperature sensor |

**Notes**:
- 12-bit resolution (0-4095 for 0-3.3V)
- 500 ksps max sample rate
- ADC_VREF (pin 35) can provide external reference (default: internal 3.3V)
- GP26-28 lose digital I/O capability when used as ADC

## GPIO Electrical Characteristics

| Parameter | Min | Typ | Max | Unit | Notes |
|-----------|-----|-----|-----|------|-------|
| Input voltage (3V3 I/O) | -0.3 | - | 3.63 | V | Do not exceed 3.6V |
| Output current (per pin) | - | - | 12 | mA | Recommended max |
| Output current (all pins) | - | - | 50 | mA | Total GPIO sink/source |
| Input leakage | - | - | ±1 | μA | With pull-up/down disabled |
| Pull-up/down resistor | 50 | 60 | 75 | kΩ | Internal, configurable |

**Important**: All GPIO are 3.3V logic. Use level shifters for 5V devices.

## Common Wiring Patterns

### Button Input (Active Low with Internal Pull-Up)
```python
from machine import Pin
button = Pin(14, Pin.IN, Pin.PULL_UP)  # GP14, reads 1 when open, 0 when pressed
```

### LED Output (with external resistor)
```python
from machine import Pin
led = Pin(10, Pin.OUT)  # GP10 → 220Ω resistor → LED → GND
led.value(1)  # Turn on
```

### PWM LED (dimming)
```python
from machine import Pin, PWM
led = PWM(Pin(10))  # GP10 (PWM5 A)
led.freq(1000)  # 1kHz
led.duty_u16(32768)  # 50% brightness (0-65535)
```

### I2C Device
```python
from machine import Pin, I2C
i2c = I2C(0, scl=Pin(1), sda=Pin(0), freq=400000)  # I2C0 on GP0/GP1
devices = i2c.scan()  # Returns list of 7-bit addresses
```

### SPI Device
```python
from machine import Pin, SPI
spi = SPI(0, baudrate=1000000, polarity=0, phase=0,
          sck=Pin(18), mosi=Pin(19), miso=Pin(16))  # SPI0
cs = Pin(17, Pin.OUT)
cs.value(1)  # Deselect
```

### 1-Wire (e.g., DS18B20)
```python
import onewire, ds18x20
from machine import Pin
ow = onewire.OneWire(Pin(22))  # GP22 with 4.7kΩ pull-up to 3.3V
ds = ds18x20.DS18X20(ow)
roms = ds.scan()
ds.convert_temp()
temp_c = ds.read_temp(roms[0])
```

## Pin Assignment Best Practices

1. **Group by function**: Keep SPI/I2C/UART pins together for cleaner wiring
2. **Use hardware peripherals**: Prefer hardware SPI/I2C over bit-banging for reliability and speed
3. **Check conflicts**: Ensure alternate functions don't overlap (e.g., SPI0 SCK on GP18 vs GP2)
4. **PWM planning**: Group PWM outputs on different slices if independent frequencies needed
5. **ADC isolation**: Keep analog inputs away from noisy digital signals
6. **Power budget**: Calculate total current draw; keep 3V3(OUT) load under 300mA
7. **Ground strategy**: Use multiple GND pins for low-impedance return paths
8. **Reserved pins**: Avoid GP23-25, GP29 (used by wireless module)

## Debugging Tips

- **Pin not working?** Check if it's configured correctly (IN/OUT, pull-up/down)
- **SPI/I2C issues?** Verify pin mapping matches the peripheral (SPI0 vs SPI1, I2C0 vs I2C1)
- **PWM conflict?** Check if both channels on the same slice need different frequencies
- **Voltage problems?** Measure 3V3(OUT) under load; ensure it's stable at 3.3V
- **Wireless not working?** Don't use GP23-25 or GP29 for other purposes

## MicroPython for Pico W

### What is MicroPython?

MicroPython is a lean implementation of Python 3 optimized for microcontrollers. It includes:
- Most Python 3 standard library features
- Hardware-specific modules (`machine`, `rp2`, `network`)
- Interactive REPL (Read-Eval-Print Loop) over USB serial
- ~200KB firmware footprint, leaving ~1.8MB for user code on Pico W

**Key advantages**:
- Rapid prototyping (no compile cycle)
- Interactive debugging via REPL
- Rich ecosystem of Python libraries
- Easy file management over USB

### Installation

1. **Download firmware**: Get the latest `.uf2` file from https://micropython.org/download/RPI_PICO_W/
2. **Enter bootloader mode**: Hold BOOTSEL button while plugging in USB (or press BOOTSEL + RUN if already powered)
3. **Flash firmware**: Pico W appears as USB mass storage device; drag `.uf2` file onto it
4. **Reboot**: Device automatically reboots into MicroPython

**Verify installation**:
```bash
# Linux/Mac
screen /dev/ttyACM0 115200

# Windows (use Thonny IDE or PuTTY)
```

You should see the MicroPython REPL prompt: `>>>`

### Development Environment

**Recommended IDEs**:
- **Thonny** (beginner-friendly, built-in REPL, file manager): https://thonny.org/
- **VS Code** with Pico-W-Go extension
- **rshell** (command-line file transfer): `pip install rshell`
- **mpremote** (official tool): `pip install mpremote`

**Thonny setup**:
1. Install Thonny
2. Tools → Options → Interpreter
3. Select "MicroPython (Raspberry Pi Pico)"
4. Choose correct COM/serial port
5. Click "Install or update MicroPython" if needed

### Core MicroPython Modules

#### machine — Hardware Access

```python
from machine import Pin, SPI, I2C, PWM, ADC, Timer
import time

# Digital I/O
led = Pin(25, Pin.OUT)  # Note: Use Pin('LED') for onboard LED on Pico W
led.on()
led.off()
led.toggle()

button = Pin(14, Pin.IN, Pin.PULL_UP)
if button.value() == 0:  # Pressed (active low)
    print("Button pressed")

# PWM (all GPIO support PWM)
pwm = PWM(Pin(10))
pwm.freq(1000)  # 1kHz
pwm.duty_u16(32768)  # 50% duty cycle (0-65535)

# ADC (GP26-28 only)
adc = ADC(Pin(26))  # ADC0
voltage = adc.read_u16() * 3.3 / 65535  # Convert to volts

# Temperature sensor (internal)
temp_sensor = ADC(4)  # ADC channel 4
temp_c = 27 - (temp_sensor.read_u16() * 3.3 / 65535 - 0.706) / 0.001721

# Timers
def callback(timer):
    print("Timer fired")

tim = Timer()
tim.init(period=1000, mode=Timer.PERIODIC, callback=callback)  # Every 1s
```

#### SPI

```python
from machine import Pin, SPI

# SPI0 on GP16-19
spi = SPI(0, baudrate=1_000_000, polarity=0, phase=0,
          sck=Pin(18), mosi=Pin(19), miso=Pin(16))

cs = Pin(17, Pin.OUT)
cs.value(1)  # Deselect

# Write and read
cs.value(0)
spi.write(b'\x01\x02\x03')
data = spi.read(4)  # Read 4 bytes
cs.value(1)

# Write-then-read (common pattern)
cs.value(0)
result = spi.write_readinto(b'\xAA', bytearray(1))
cs.value(1)
```

#### I2C

```python
from machine import Pin, I2C

# I2C0 on GP0/GP1
i2c = I2C(0, scl=Pin(1), sda=Pin(0), freq=400_000)  # 400kHz

# Scan for devices
devices = i2c.scan()  # Returns list of 7-bit addresses
print(f"Found devices: {[hex(d) for d in devices]}")

# Read/write
i2c.writeto(0x3C, b'\x00\x01')  # Write to address 0x3C
data = i2c.readfrom(0x3C, 2)    # Read 2 bytes

# Read register (common pattern)
i2c.writeto(0x3C, b'\x00')      # Set register address
value = i2c.readfrom(0x3C, 1)   # Read register value
```

#### UART

```python
from machine import Pin, UART

# UART0 on GP0/GP1
uart = UART(0, baudrate=115200, tx=Pin(0), rx=Pin(1))

# Write
uart.write("Hello\n")

# Read (non-blocking)
if uart.any():  # Check if data available
    data = uart.read()  # Read all available
    # or: data = uart.readline()  # Read until \n

# Read (blocking with timeout)
uart.init(baudrate=115200, timeout=1000)  # 1s timeout
data = uart.read(10)  # Read up to 10 bytes or timeout
```

#### network — WiFi (Pico W specific)

```python
import network
import time

# Connect to WiFi
wlan = network.WLAN(network.STA_IF)
wlan.active(True)
wlan.connect('SSID', 'PASSWORD')

# Wait for connection
max_wait = 10
while max_wait > 0:
    if wlan.status() < 0 or wlan.status() >= 3:
        break
    max_wait -= 1
    print('Waiting for connection...')
    time.sleep(1)

# Check status
if wlan.status() != 3:
    raise RuntimeError('Network connection failed')
else:
    print('Connected')
    status = wlan.ifconfig()
    print(f'IP: {status[0]}')

# Scan networks
wlan.active(True)
networks = wlan.scan()
for ssid, bssid, channel, RSSI, authmode, hidden in networks:
    print(f"SSID: {ssid.decode()}, RSSI: {RSSI}")

# Access Point mode
ap = network.WLAN(network.AP_IF)
ap.config(essid='PicoW-AP', password='12345678')
ap.active(True)
print(f'AP IP: {ap.ifconfig()[0]}')
```

#### rp2 — RP2040 Specific Features

```python
import rp2

# PIO (Programmable I/O) - advanced topic
@rp2.asm_pio(set_init=rp2.PIO.OUT_LOW)
def blink():
    wrap_target()
    set(pins, 1)   [31]
    nop()          [31]
    set(pins, 0)   [31]
    nop()          [31]
    wrap()

# State machine
sm = rp2.StateMachine(0, blink, freq=2000, set_base=Pin(25))
sm.active(1)

# Flash memory access (read-only)
import rp2
flash_size = rp2.Flash().size  # 2MB on Pico W

# Unique board ID
import machine
uid = machine.unique_id()
print(f"UID: {uid.hex()}")
```

### MicroPython File System

Files are stored in the internal flash (LittleFS filesystem):

```python
# Write file
with open('config.txt', 'w') as f:
    f.write('Setting=Value\n')

# Read file
with open('config.txt', 'r') as f:
    content = f.read()

# List files
import os
os.listdir('/')  # Root directory
os.mkdir('logs')
os.remove('old_file.txt')
os.rename('old.txt', 'new.txt')

# File info
os.stat('config.txt')  # Returns (mode, ino, dev, nlink, uid, gid, size, atime, mtime, ctime)
```

### Boot Sequence

MicroPython executes files in this order:
1. **boot.py** (if exists) — runs first, typically for hardware init
2. **main.py** (if exists) — your main application code

**Example boot.py**:
```python
# boot.py - runs on every boot (even during REPL)
import machine
import time

# Optional: disable WiFi to save power
# import network
# network.WLAN(network.STA_IF).active(False)

print("Boot complete")
```

**Example main.py**:
```python
# main.py - main application
from lib.display import Display
from lib.ui import UI

print("Starting film developer app...")
ui = UI()
ui.run()
```

### Memory Management

```python
import gc
import micropython

# Garbage collection
gc.collect()  # Force collection
gc.mem_free()  # Free heap memory
gc.mem_alloc()  # Allocated memory

# Memory info
micropython.mem_info()  # Detailed memory stats

# Optimize memory usage
gc.threshold(gc.mem_free() // 4 + gc.mem_alloc())  # Adjust GC threshold
```

### Common Patterns

#### Debounced Button

```python
from machine import Pin
import time

class Button:
    def __init__(self, pin, debounce_ms=50):
        self.pin = Pin(pin, Pin.IN, Pin.PULL_UP)
        self.debounce = debounce_ms
        self.last_state = 1
        self.last_time = 0
    
    def pressed(self):
        """Returns True if button was pressed (falling edge)"""
        state = self.pin.value()
        now = time.ticks_ms()
        
        if state != self.last_state and time.ticks_diff(now, self.last_time) > self.debounce:
            self.last_state = state
            self.last_time = now
            if state == 0:  # Pressed (active low)
                return True
        return False

btn = Button(14)
while True:
    if btn.pressed():
        print("Button pressed!")
    time.sleep_ms(10)
```

#### Non-Blocking Delay

```python
import time

class Timer:
    def __init__(self, interval_ms):
        self.interval = interval_ms
        self.last = time.ticks_ms()
    
    def expired(self):
        """Returns True if interval has elapsed"""
        now = time.ticks_ms()
        if time.ticks_diff(now, self.last) >= self.interval:
            self.last = now
            return True
        return False

# Usage
led_timer = Timer(1000)  # 1 second
while True:
    if led_timer.expired():
        led.toggle()
    # Do other work here
```

#### Persistent Configuration

```python
import json

def save_config(filename, config):
    """Save config dict to JSON file"""
    with open(filename, 'w') as f:
        json.dump(config, f)

def load_config(filename, default=None):
    """Load config from JSON file, return default if not found"""
    try:
        with open(filename, 'r') as f:
            return json.load(f)
    except (OSError, ValueError):
        return default if default else {}

# Usage
config = load_config('settings.json', {'wifi_ssid': '', 'temp_unit': 'C'})
config['temp_unit'] = 'F'
save_config('settings.json', config)
```

### Debugging Tips

**Print debugging**:
```python
print(f"Variable x = {x}")  # Outputs to USB serial (REPL)
```

**Exception handling**:
```python
try:
    risky_operation()
except Exception as e:
    import sys
    sys.print_exception(e)  # Print full traceback
```

**REPL over UART** (free up USB for other uses):
```python
# In boot.py
from machine import UART, Pin
import os

uart = UART(0, baudrate=115200, tx=Pin(0), rx=Pin(1))
os.dupterm(uart)  # Redirect REPL to UART0
```

**Soft reset**: `Ctrl+D` in REPL (runs boot.py and main.py again)  
**Hard reset**: `machine.reset()` or press RUN button  
**Enter REPL**: `Ctrl+C` (interrupts running program)

### Performance Considerations

- **Native code**: Use `@micropython.native` decorator for ~2x speedup on compute-heavy functions
- **Viper**: Use `@micropython.viper` for ~10x speedup (restricted Python subset)
- **Avoid allocations**: Reuse buffers instead of creating new objects in loops
- **Use `const()`**: For constants, saves RAM: `from micropython import const; LED_PIN = const(25)`
- **Frozen modules**: Pre-compile modules into firmware for faster import and less RAM

### Wireless Features (Pico W)

**Onboard LED** (controlled via CYW43439):
```python
from machine import Pin
led = Pin('LED', Pin.OUT)  # or Pin('WL_GPIO0')
led.on()
```

**Bluetooth** (experimental, requires specific firmware):
```python
import bluetooth
ble = bluetooth.BLE()
ble.active(True)
# BLE functionality is still evolving in MicroPython
```

### Troubleshooting

| Issue | Solution |
|-------|----------|
| REPL not responding | Press Ctrl+C to interrupt, or soft reset with Ctrl+D |
| Import error | Check file is in root or `/lib` directory |
| Out of memory | Call `gc.collect()`, reduce buffer sizes, use generators |
| WiFi won't connect | Check SSID/password, ensure 2.4GHz network (Pico W doesn't support 5GHz) |
| Pin not working | Verify pin number, check if it's reserved (GP23-25) |
| SPI/I2C not working | Verify pin mapping matches peripheral (SPI0 vs SPI1, I2C0 vs I2C1) |

## External References

### Hardware
- **Interactive pinout**: https://picow.pinout.xyz/
- **Official pinout PDF**: https://datasheets.raspberrypi.com/picow/PicoW-A4-Pinout.pdf
- **Pico W datasheet**: https://datasheets.raspberrypi.com/picow/pico-w-datasheet.pdf
- **RP2040 datasheet**: https://datasheets.raspberrypi.com/rp2040/rp2040-datasheet.pdf

### MicroPython
- **Quick reference**: https://docs.micropython.org/en/latest/rp2/quickref.html
- **RP2 library docs**: https://docs.micropython.org/en/latest/library/rp2.html
- **Firmware downloads**: https://micropython.org/download/RPI_PICO_W/
- **MicroPython forum**: https://forum.micropython.org/
- **Getting started guide**: https://www.raspberrypi.com/documentation/microcontrollers/micropython.html

## Revision History

| Date       | Change |
|------------|--------|
| 2025-11-08 | Initial reference document created for FILM-DEV-001 |
| 2025-11-08 | Added comprehensive MicroPython section with API examples, patterns, and troubleshooting |
