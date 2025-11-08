---
Title: Application Architecture and System Overview
Ticket: FILM-DEV-001
Status: active
Topics:
    - embedded
    - hardware
    - timer
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Hardware platform, system architecture, and core components for the film development timer
LastUpdated: 2025-11-08T16:51:17.999387756-05:00
---

# Application Architecture and System Overview

## Hardware Platform

### Raspberry Pi Pico W Specifications
- **Microcontroller**: RP2040 dual-core ARM Cortex-M0+ @ 133MHz
- **Memory**: 264KB SRAM, 2MB Flash storage
- **Wireless**: 802.11n WiFi (2.4GHz) via CYW43439
- **GPIO**: 26 multi-function GPIO pins
- **ADC**: 3x 12-bit ADC channels + 1x internal temperature sensor
- **Power**: 1.8-5.5V DC input, 3.3V logic levels

### Hardware Components

#### Display System
- **Display**: Waveshare SH1107 64x128 pixel OLED (1.3")
- **Interface**: SPI (4-wire)
- **Font**: 6x8 pixels per character
- **Usable Area**: 21 characters × 8 rows
- **Features**: High contrast monochrome display

#### Input System
- **Buttons**: 3x tactile push buttons with internal pull-ups
- **LEDs**: 3x PWM-controlled LEDs with 220Ω resistors
- **Layout**: [Button 1] [Button 2] [Button 3]
- **Debouncing**: Software debouncing with 50ms delay

#### Temperature Monitoring
- **Sensor**: DS18B20 digital temperature sensor (waterproof probe)
- **Interface**: 1-Wire protocol on GPIO 22
- **Pull-up**: 4.7kΩ resistor between DATA and VCC
- **Accuracy**: ±0.5°C
- **Conversion Time**: 750ms for 12-bit resolution

### System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
├─────────────────────────────────────────────────────────────┤
│  UI Controller  │  Timer Engine  │  Data Manager  │  Config │
├─────────────────────────────────────────────────────────────┤
│                    Hardware Abstraction Layer              │
├─────────────────────────────────────────────────────────────┤
│  Display Driver │  Button Handler │ Temp Sensor │  Storage │
├─────────────────────────────────────────────────────────────┤
│                    MicroPython Runtime                     │
├─────────────────────────────────────────────────────────────┤
│                    Raspberry Pi Pico W                     │
└─────────────────────────────────────────────────────────────┘
```

## Core Components

### 1. UI Controller
**Responsibility**: Manages all user interface interactions and display updates

**Key Functions**:
- Screen rendering and layout management
- Button event handling and LED control
- Menu navigation and state transitions
- Alert and notification display

**State Management**:
- Current screen/menu context
- Button press handling with debouncing
- LED status indicators
- Display refresh coordination

### 2. Timer Engine
**Responsibility**: Manages all timing operations and development process control

**Key Functions**:
- Multi-stage timer management (Developer, Stop Bath, Fixer, Wash)
- Pause/resume functionality
- Timer adjustments during operation
- Stage transition automation
- Alert generation (visual/audio)

**Timer States**:
- `READY` - Timer configured, ready to start
- `RUNNING` - Timer actively counting down
- `PAUSED` - Timer paused, can be resumed
- `COMPLETE` - Timer finished, awaiting next action
- `OVERTIME` - Timer exceeded planned duration

### 3. Data Manager
**Responsibility**: Handles all data persistence, logging, and configuration

**Key Functions**:
- Session logging with run IDs
- Film database management
- Favourites storage and retrieval
- Configuration persistence
- Temperature data logging

**Data Structures**:
- Run logs with timestamps and actions
- Film/developer combination database
- User favourites and presets
- System configuration settings

### 4. Configuration System
**Responsibility**: Manages application settings and calibration data

**Key Functions**:
- Temperature sensor calibration
- Display settings and preferences
- System information and diagnostics
- Factory reset capabilities

## Development Framework

### MicroPython Implementation
- **Language**: MicroPython 1.20+ for Pico W
- **Libraries**: 
  - `machine` - Hardware control (GPIO, SPI, PWM)
  - `time` - Timing and delays
  - `json` - Data serialization
  - `sh1107` - SH1107 OLED display driver
  - `onewire`/`ds18x20` - Temperature sensor

### Pin Assignments
| Component | GPIO | Physical Pin | Notes |
|-----------|------|--------------|-------|
| **OLED Display (SPI)** |
| VCC | 3V3 | Pin 36 | Power |
| GND | GND | Pin 38 | Ground |
| DIN (MOSI) | GP19 | Pin 25 | Data |
| CLK (SCK) | GP18 | Pin 24 | Clock |
| CS | GP17 | Pin 22 | Chip Select |
| DC | GP20 | Pin 26 | Data/Command |
| RST | GP21 | Pin 27 | Reset |
| **Temperature Sensor** |
| DATA | GP22 | Pin 29 | 1-Wire with 4.7kΩ pull-up |
| VCC | 3V3 | Pin 36 | Shared with OLED |
| GND | GND | Pin 33 | Ground |
| **Buttons (Internal Pull-up)** |
| Button 1 | GP14 | Pin 19 | To GND when pressed |
| Button 2 | GP15 | Pin 20 | To GND when pressed |
| Button 3 | GP16 | Pin 21 | To GND when pressed |
| Common GND | GND | Pin 33 | Shared |
| **LEDs (with 220Ω resistors)** |
| LED 1 | GP10 | Pin 14 | PWM capable |
| LED 2 | GP11 | Pin 15 | PWM capable |
| LED 3 | GP12 | Pin 16 | PWM capable |
| Common GND | GND | Pin 28 | Shared |

### File Structure
```
/
├── main.py              # Application entry point
├── config.json          # System configuration
├── film_database.json   # Film/developer combinations
├── lib/
│   ├── ui_controller.py # UI management
│   ├── timer_engine.py  # Timer logic
│   ├── data_manager.py  # Data persistence
│   ├── hardware.py      # Hardware abstraction
│   └── utils.py         # Utility functions
└── logs/
    └── runs/            # Session log files
```

## System Initialization

### Boot Sequence
1. **Hardware Initialization**
   - GPIO pin configuration
   - SPI bus setup for display
   - 1-Wire bus setup for temperature sensor
   - Button configuration with internal pull-ups
   - PWM setup for LEDs

2. **Display Initialization**
   - SH1107 OLED controller setup
   - Clear screen and show ready state

3. **Data System Initialization**
   - Configuration file loading
   - Film database loading
   - Log system initialization

4. **Application Ready**
   - Main screen display
   - Temperature reading start
   - Ready for user input

### Error Handling
- **Hardware Failures**: Display error message on screen
- **Data Corruption**: Fallback to default configurations
- **Sensor Failures**: Continue operation with "Temp: ERROR" display
- **File System Issues**: Create default configuration files