# Tasks

## Specification Tasks (Complete)

- [x] Define hardware architecture and pin assignments
- [x] Document complete UI flows and screen layouts
- [x] Specify timer system and multi-stage workflow
- [x] Design data management and logging system
- [x] Move temperature compensation to future-ideas
- [x] Update specification to MVP focus

## Implementation Tasks (Ready for Developer)

### Phase 1: Hardware Layer
- [ ] Set up MicroPython development environment for Pico W
- [ ] Implement SPI driver for SH1107 OLED display (64x128)
- [ ] Implement 1-Wire driver for DS18B20 temperature sensor
- [ ] Configure GPIO for 3 buttons with internal pull-ups (GP14-16)
- [ ] Configure PWM for 3 LEDs (GP10-12)
- [ ] Create hardware abstraction layer module
- [ ] Test each hardware component independently

### Phase 2: UI Controller
- [ ] Implement display rendering for 16x8 character grid (128x64, 8x8 font)
- [ ] Create screen layout templates
- [ ] Implement button handler with 50ms debouncing
- [ ] Implement LED control with blink patterns
- [x] Build navigation state machine
- [ ] Implement all 19 screen layouts from UI spec
- [ ] Test button navigation and screen transitions

### Phase 3: Timer Engine
- [ ] Implement timer state machine (READY/RUNNING/PAUSED/COMPLETE)
- [ ] Create multi-stage workflow (Developer/Stop/Fixer/Wash)
- [x] Implement pause/resume functionality
- [ ] Add real-time timer adjustment (+/- time)
- [ ] Implement overtime tracking
- [ ] Create alert system (visual and LED blink patterns)
- [ ] Test timer accuracy and stage transitions

### Phase 4: Data Management
- [ ] Create film database JSON structure
- [ ] Implement configuration file loading
- [ ] Build session logging with run ID generation (YYYY-MMDD-NNN)
- [ ] Implement favourites storage and retrieval
- [ ] Create atomic file write operations with backup
- [x] Add error recovery for corrupted files
- [x] Test data persistence across reboots

### Phase 5: Integration
- [ ] Connect UI controller to timer engine
- [x] Integrate temperature monitoring into displays
- [ ] Wire up film database to timer configuration
- [ ] Connect session logging to timer events
- [ ] Implement favourites add/select workflow
- [ ] Test complete end-to-end development workflow

### Phase 6: Testing & Polish
- [ ] Test full film development workflow (all 4 stages)
- [ ] Verify session logging captures all events
- [ ] Test temperature sensor error handling
- [ ] Verify button debouncing and LED feedback
- [ ] Test with multiple film/developer combinations
- [ ] Validate favourites system
- [ ] Document any issues or edge cases discovered
- [ ] Create UI MVP analysis doc with MicroPython references (framebuf, SPI, ds18x20)
- [ ] Draft UI MVP design and implementation plan (display/input/temp/ui/main modules)
- [ ] Evaluate adopting 6x8 font to enable 21x8 grid; implement renderer if chosen
