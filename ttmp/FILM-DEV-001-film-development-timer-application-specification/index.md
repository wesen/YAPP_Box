---
Title: Film Development Timer Application Specification
Ticket: FILM-DEV-001
Status: review
Topics:
    - embedded
    - ui
    - hardware
    - timer
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/code/others/YAPP_Box/film-developer/README.md
      Note: Top-level wiring summary including all 3 buttons
    - Path: /home/manuel/code/others/YAPP_Box/film-developer/lib/display.py
      Note: Display wrapper for SH1107 16×8 text grid (8×8 font on 128×64 logical buffer)
    - Path: /home/manuel/code/others/YAPP_Box/film-developer/lib/input.py
      Note: Button input with debounce (BTN1/2/3)
    - Path: /home/manuel/code/others/YAPP_Box/film-developer/lib/sh1107.py
      Note: Minimal SH1107 SPI driver for text mode
    - Path: /home/manuel/code/others/YAPP_Box/film-developer/lib/temp.py
      Note: DS18B20 temperature helper
    - Path: /home/manuel/code/others/YAPP_Box/film-developer/lib/ui.py
      Note: UI controller (Splash/Main/Menu/Timer stub)
    - Path: /home/manuel/code/others/YAPP_Box/film-developer/main.py
      Note: UI MVP entry point (loop)
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/FILM-DEV-001-film-development-timer-application-specification/scripts/ui-mvp/lib/display.py
      Note: Display wrapper for SH1107 21x8 text grid
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/FILM-DEV-001-film-development-timer-application-specification/scripts/ui-mvp/lib/input.py
      Note: Button input with debounce (BTN1/2/3)
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/FILM-DEV-001-film-development-timer-application-specification/scripts/ui-mvp/lib/sh1107.py
      Note: Minimal SH1107 SPI driver for text mode
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/FILM-DEV-001-film-development-timer-application-specification/scripts/ui-mvp/lib/temp.py
      Note: DS18B20 temperature helper
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/FILM-DEV-001-film-development-timer-application-specification/scripts/ui-mvp/lib/ui.py
      Note: UI controller (Splash/Main/Menu/Timer stub)
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/FILM-DEV-001-film-development-timer-application-specification/scripts/ui-mvp/main.py
      Note: UI MVP entry point (loop)
    - Path: ttmp/FILM-DEV-001-*/reference/05-raspberry-pi-pico-w-complete-pinout-reference.md
      Note: Comprehensive Pico W pinout lookup for all GPIO planning
    - Path: ttmp/FILM-DEV-001-film-development-timer-application-specification/reference/05-raspberry-pi-pico-w-hardware-and-micropython-reference.md
      Note: Complete Pico W hardware pinout and MicroPython API reference
ExternalSources: []
Summary: Comprehensive specification for Raspberry Pi Pico W film development timer with multi-stage timing, temperature monitoring, and session logging
LastUpdated: 2025-11-08T17:08:05.30832586-05:00
---





















# Film Development Timer Application Specification

## Overview

This specification defines an MVP film development timer application for Raspberry Pi Pico W. The application provides precise timing control for analog film development processes including developer, stop bath, fixer, and wash stages. It features a 128x64 OLED display, 3-button interface with LED indicators, temperature monitoring, and comprehensive session logging.

**Status**: ✅ **Specification Complete - Ready for Implementation**

This specification is ready to be handed off to a developer for implementation. All core features, hardware wiring, UI flows, and data structures are fully documented.

## MVP Features

- ✅ Multi-stage development timer (Developer → Stop Bath → Fixer → Wash)
- ✅ Pause/resume functionality with overtime tracking
- ✅ Film database with developer combinations and push/pull processing
- ✅ Temperature monitoring and display (no automatic compensation in MVP)
- ✅ Session logging with unique run IDs and complete audit trails
- ✅ Favourites system for quick access to common combinations
- ✅ Real-time timer adjustments during development
- ✅ Visual and LED alerts for stage completion
- ✅ 3-button navigation with LED indicators
- ✅ 128x64 OLED display with clear layouts

## Hardware Requirements

### Components
- Raspberry Pi Pico W
- Waveshare SH1107 64x128 OLED (1.3", SPI interface)
- DS18B20 waterproof temperature sensor with 4.7kΩ pull-up resistor
- 3x tactile push buttons (with internal pull-ups)
- 3x LEDs with 220Ω resistors
- Breadboard and jumper wires

### Wiring
Complete pin assignments documented in [Architecture Overview](./reference/01-application-architecture-and-system-overview.md#pin-assignments):
- OLED on SPI: GP17-21 (CS, CLK, DIN, DC, RST)
- Temperature sensor: GP22 (1-Wire)
- Buttons: GP14-16 (to GND when pressed)
- LEDs: GP10-12 (PWM capable)

## Key Documents

### Core Specification
1. **[Application Architecture and System Overview](./reference/01-application-architecture-and-system-overview.md)**
   - Hardware platform and pin assignments
   - System architecture with 4 core components
   - MicroPython implementation details
   - Boot sequence and error handling

2. **[User Interface Specification and Display System](./reference/02-user-interface-specification-and-display-system.md)**
   - Complete UI flows with 19 different screens
   - 3-button navigation system
   - Display layouts for 21x8 character grid
   - Interaction patterns and visual feedback

3. **[Timer System and Development Process Management](./reference/03-timer-system-and-development-process-management.md)**
   - Multi-stage timer engine state machine
   - Film database structure with examples
   - Real-time timer adjustments
   - Favourites system
   - Alert and notification patterns

4. **[Data Management and Logging System](./reference/04-data-management-and-logging-system.md)**
   - Session logging with run IDs (YYYY-MMDD-NNN format)
   - Configuration and film database persistence
   - Atomic file operations and backup strategy
   - Error recovery and data integrity

### Future Enhancements
5. **[Temperature Compensation and Advanced Features](./future-ideas/01-temperature-compensation-and-advanced-features.md)**
   - Automatic temperature compensation algorithms (deferred from MVP)
   - Implementation considerations for future versions

### Supporting Docs
6. **[UI MVP Analysis and MicroPython References](./analysis/01-ui-mvp-analysis-and-micropython-references.md)**
   - Scope, constraints, APIs, wiring for display/buttons/temp
7. **[UI MVP Design and Implementation Plan](./design-doc/01-ui-mvp-design-and-implementation-plan.md)**
   - File layout, pseudocode, and test plan for UI MVP
8. **[docmgr Workflow Playbook for This Project](./playbook/01-docmgr-workflow-playbook-for-this-project.md)**
   - Project-specific docmgr usage (doc types, commands, workflow)

## Implementation Guidance

### Getting Started
1. Set up MicroPython on Raspberry Pi Pico W
2. Wire hardware according to pin assignments in Architecture doc
3. Implement hardware abstraction layer (GPIO, SPI, 1-Wire, PWM)
4. Build UI controller with SH1107 display driver
5. Implement timer engine with state machine
6. Create film database and data persistence layer
7. Test each component individually before integration

### Development Order (Recommended)
1. **Hardware Layer** - Test display, buttons, LEDs, temperature sensor independently
2. **UI Screens** - Implement static screens first, then navigation
3. **Timer Engine** - Build basic countdown timer, then add stages
4. **Data System** - Configuration loading, then logging
5. **Integration** - Connect all components together
6. **Testing** - Full end-to-end development workflow

### Key Implementation Notes
- Display is 21 characters × 8 rows (6x8 pixel font)
- Buttons use internal pull-ups (active low)
- LEDs support PWM for brightness control
- Temperature readings every 30 seconds during operation
- Session logs stored as JSON files in `logs/runs/`
- All times stored as strings in "MM:SS" format

## Status

Current status: **review** - Specification complete, ready for implementation handoff

## Topics

- embedded
- ui
- hardware
- timer

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- design/ - Architecture and design documents
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
