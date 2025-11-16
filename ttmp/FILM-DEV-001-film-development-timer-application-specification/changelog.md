# Changelog

## 2025-11-08

- Initial workspace created


## 2025-11-08

Created comprehensive specification with 4 reference documents covering architecture, UI, timer system, and data management


## 2025-11-08

Updated specification to MVP focus: removed performance requirements, future enhancements, and optimization sections. Updated hardware wiring to match existing Pico W setup with SH1107 OLED, DS18B20, buttons on GP14-16, and LEDs on GP10-12


## 2025-11-08

Moved temperature compensation to future-ideas document. MVP will display temperature but not automatically adjust development times


## 2025-11-08

Specification complete and ready for implementation handoff. All MVP features documented with hardware wiring, UI flows, timer logic, and data management


## 2025-11-08

Moved UI MVP analysis and design docs from PICO-TEMP-001 to FILM-DEV-001; updated frontmatter and topics


## 2025-11-08

Added docmgr Workflow Playbook and linked from index


## 2025-11-08

Expanded docmgr Workflow Playbook with context, guidelines, examples for onboarding


## 2025-11-08

Added UI MVP code skeleton (main, display, input, temp, ui, sh1107)


## 2025-11-08

Moved UI MVP code to root film-developer/ via git mv


## 2025-11-08

Added film-developer/README.md with 3-button wiring and quickstart


## 2025-11-08

Updated film-developer/README.md with USB power notes, full ASCII wiring diagram, and Pico W pinout references


## 2025-11-08

Added detailed ASCII Pico W pinout section to film-developer/README.md (USB up orientation)


## 2025-11-08

Added full ASCII pinout diagram to film-developer/README.md showing all 40 pins with used connections labeled


## 2025-11-08

Added comprehensive Raspberry Pi Pico W pinout reference (05) with all GPIO, alternate functions, power specs, and wiring examples


## 2025-11-08

Expanded reference doc 05 with comprehensive MicroPython section (installation, APIs, WiFi, patterns, debugging); renamed to hardware-and-micropython-reference


## 2025-11-15

UI grid decision: hardware is 128×64; adopting 8×8 font → 16×8 chars. Adjust UI labels to fit and truncate safely where needed.

### Related Files

- /home/manuel/code/others/YAPP_Box/film-developer/lib/display.py — Display helper enforces 16×8 text grid (8×8 font)
- /home/manuel/code/others/YAPP_Box/ttmp/FILM-DEV-001-film-development-timer-application-specification/reference/02-user-interface-specification-and-display-system.md — UI spec source; note grid differs when using 8×8 font


## 2025-11-15

Correction: 21×8 is achievable with a 6×8 font. Current code uses framebuf’s 8×8 font → 16×8 grid. Decision: keep 8×8 for now; evaluate adding a 6×8 font renderer if wider labels are needed.

### Related Files

- /home/manuel/code/others/YAPP_Box/film-developer/lib/display.py — Current helper assumes 8×8 font (16 cols); 6×8 would permit 21 cols.
- /home/manuel/code/others/YAPP_Box/ttmp/FILM-DEV-001-film-development-timer-application-specification/reference/02-user-interface-specification-and-display-system.md — UI layouts can map to 16×8 (8×8 font) or 21×8 (6×8 font)


## 2025-11-15

Implemented UI navigation state machine with core screens (splash, main, menu, system info, timer running/paused, stage advance, done). Switched app entry to run UI loop by default.

### Related Files

- /home/manuel/code/others/YAPP_Box/film-developer/lib/ui.py — State machine + renderers for 16x8 grid; stubs for extended screens
- /home/manuel/code/others/YAPP_Box/film-developer/main.py — Main now runs UI loop by default (tests retained)


## 2025-11-15

Added serial debug for button input (edge/bounce logs) and UI state transitions to diagnose non-registering clicks.

### Related Files

- /home/manuel/code/others/YAPP_Box/film-developer/lib/input.py — Debug logs for press/release/bounce; toggle via constructor
- /home/manuel/code/others/YAPP_Box/film-developer/lib/ui.py — State transition logs on events for navigation


## 2025-11-15

UI input diagnostics: print every button event and state; reduced Input debounce to 20ms to improve responsiveness.

### Related Files

- /home/manuel/code/others/YAPP_Box/film-developer/lib/ui.py — Log each event regardless of transition
- /home/manuel/code/others/YAPP_Box/film-developer/main.py — Instantiate Input(debounce_ms=20
- debug — True)


## 2025-11-15

Added Minimal UI mode for performance testing (no temp/timer); redraws only on change to avoid blocking from DS18B20. Main runs Minimal UI by default for now.

### Related Files

- /home/manuel/code/others/YAPP_Box/film-developer/lib/ui_min.py — Non-blocking input/display test UI
- /home/manuel/code/others/YAPP_Box/film-developer/main.py — Switch to MinimalUI loop for debugging responsiveness


## 2025-11-15

Added IRQ-backed AsyncInput to avoid missed presses under slower UI loops; Minimal UI now prefers AsyncInput (ghost=15ms), falls back to polling if unavailable.

### Related Files

- /home/manuel/code/others/YAPP_Box/film-developer/lib/input_irq.py — Event-queued
- /home/manuel/code/others/YAPP_Box/film-developer/main.py — Prefer AsyncInput for Minimal UI; fallback maintained


## 2025-11-15

Switched back to main UI for testing; disabled temperature reads (enable_temp=False) to prevent DS18B20 blocking. Prefer IRQ-backed input for responsiveness.

### Related Files

- /home/manuel/code/others/YAPP_Box/film-developer/lib/ui.py — UI now supports enable_temp flag to skip temp reads
- /home/manuel/code/others/YAPP_Box/film-developer/main.py — Runs UI with enable_temp=False and AsyncInput(ghost=40ms)


## 2025-11-15

Implemented functional multi-stage timer: start/pause/resume/next integrated with UI. Shows remaining time and progress bar; temperature disabled for now to avoid blocking.

### Related Files

- /home/manuel/code/others/YAPP_Box/film-developer/lib/timer.py — New TimerEngine with non-blocking tick and multi-stage support
- /home/manuel/code/others/YAPP_Box/film-developer/lib/ui.py — UI wires TimerEngine into states


## 2025-11-15

Implemented non-blocking DS18B20 temperature updates: periodic conversion with cooperative polling; UI displays latest reading without blocking.

### Related Files

- /home/manuel/code/others/YAPP_Box/film-developer/lib/temp.py — Added TempNonBlocking with poll()+get_c()
- /home/manuel/code/others/YAPP_Box/film-developer/lib/ui.py — Polls temperature each loop; renders via get_c()
- /home/manuel/code/others/YAPP_Box/film-developer/main.py — Runs UI with TempNonBlocking
- enable_temp — True


## 2025-11-15

Wired temperature into Timer/Paused screens using non-blocking TempNonBlocking; UI now shows Temp on main and during runs.

### Related Files

- /home/manuel/code/others/YAPP_Box/film-developer/lib/ui.py — Enabled _render_temp_line on timer screens


## 2025-11-15

Added presets system: JSON-backed presets (data/presets.json) and Presets menu in UI. Default preset: Kodak TMax 400 + Xtol 1+1 @ 20C; applies developer time to timer.

### Related Files

- /home/manuel/code/others/YAPP_Box/film-developer/data/presets.json — Seed presets for testing (TMX400/TRX400 with XTOL 1+1)
- /home/manuel/code/others/YAPP_Box/film-developer/lib/presets.py — Loader for JSON presets
- /home/manuel/code/others/YAPP_Box/film-developer/lib/ui.py — Presets menu


## 2025-11-15

Added marquee scrolling for long labels in Presets screen (wrap-around, ~200ms step) so selected setting scrolls like a banner when too long.

### Related Files

- /home/manuel/code/others/YAPP_Box/film-developer/lib/ui.py — Marquee helper and integration on Presets selected item


## 2025-11-15

Presets UI now hierarchical: Developer → Film → Target ISO. Seeded XTOL push+2 entries for Tri‑X 400 (11:45 @20C, 1+1) and HP5+ (18:00 @20C, 1+1). Selection applies developer time to timer.

### Related Files

- /home/manuel/code/others/YAPP_Box/film-developer/data/presets.json — Added XTOL ISO1600 entries for Tri‑X 400 and HP5+
- /home/manuel/code/others/YAPP_Box/film-developer/lib/presets.py — Grouping + best-preset selection (prefer 1+1 @20C)
- /home/manuel/code/others/YAPP_Box/film-developer/lib/ui.py — Added DEV/FILM/ISO states and logic

