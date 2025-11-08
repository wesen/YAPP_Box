---
Title: User Interface Specification and Display System
Ticket: FILM-DEV-001
Status: active
Topics:
    - ui
    - embedded
    - hardware
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Complete UI flows, display layouts, and interaction patterns for the film development timer
LastUpdated: 2025-11-08T16:51:17.999387756-05:00
---

# User Interface Specification and Display System

## Display Specifications

### Hardware Constraints
- **Resolution**: 128x64 pixels
- **Font Size**: 6x8 pixels per character
- **Usable Area**: 21 characters × 8 rows
- **Display Technology**: OLED (high contrast, fast refresh)

### Layout Framework
```
┌─────────────────────┐  ← 21 characters wide
│ Header Row          │  ← Row 1: Title/Status
│                     │  ← Row 2: Spacer
│ Content Area        │  ← Rows 3-6: Main content
│ Content Area        │
│ Content Area        │
│ Content Area        │
│                     │  ← Row 7: Spacer
│ Button Labels       │  ← Row 8: Button functions
└─────────────────────┘
```

## Button System

### Physical Layout
```
[●] [●] [●]  ← LED indicators
[1] [2] [3]  ← Physical buttons
```

### Button Functions by Context

| Context | Button 1 | Button 2 | Button 3 |
|---------|----------|----------|----------|
| Main Screen | START | MENU | (disabled) |
| Configuration | SELECT | NAVIGATE | EXIT |
| Timer Running | PAUSE | (disabled) | NEXT |
| Timer Paused | RESUME | (disabled) | NEXT |
| Edit Mode | APPLY | CYCLE | CANCEL |

### LED Indicators
- **Solid**: Function available
- **Off**: Function disabled
- **Slow Blink (500ms)**: Attention/paused state
- **Fast Blink (200ms)**: Alert/completion
- **All LEDs Fast Blink**: Critical alert (timer complete + overtime)

## Screen Definitions

### 1. Main Screen (Ready State)
```
┌─────────────────────┐
│   FILM DEV TIMER    │
│                     │
│ Temp: 20.1°C        │
│ Target: 20.0°C      │
│                     │
│ Ready to start      │
│                     │
│[START][MENU][    ]  │
└─────────────────────┘
```

**Elements**:
- Title: Application name
- Temperature: Current probe reading
- Target: Expected temperature
- Status: System state message
- Buttons: START (active), MENU (active), third disabled

### 2. Main Screen (Timer Loaded)
```
┌─────────────────────┐
│   FILM DEV TIMER    │
│                     │
│ TRX+2 D76 22°C      │
│ Timer: 14:30        │
│ Temp: 22.1°C ✓      │
│ Ready to develop    │
│                     │
│[START][MENU][    ]  │
└─────────────────────┘
```

**Elements**:
- Configuration: Film type, developer, temperature
- Timer: Configured duration
- Temperature: Current with status indicator (✓ = good, ! = warning)
- Status: Ready message

### 3. Configuration Menu
```
┌─────────────────────┐
│   CONFIGURATION     │
│                     │
│ ► Timer Presets     │
│   Temperature Cal   │
│   System Info       │
│                     │
│                     │
│ [SEL][NAV][EXIT]    │
└─────────────────────┘
```

**Elements**:
- Title: Current menu name
- Menu Items: List with selection indicator (►)
- Navigation: SELECT, NAVIGATE, EXIT functions

### 4. Timer Presets Menu
```
┌─────────────────────┐
│   TIMER PRESETS     │
│                     │
│ ► Favourites        │
│   Film Database     │
│   Manual Entry      │
│   Recent Timers     │
│                     │
│ [SEL][NAV][BACK]    │
└─────────────────────┘
```

### 5. Film Database Navigation
```
┌─────────────────────┐
│   FILM DATABASE     │
│                     │
│ ► Kodak Tri-X       │
│   Ilford HP5+       │
│   Fuji Acros        │
│   Kodak T-Max 400   │
│                     │
│ [SEL][NAV][BACK]    │
└─────────────────────┘
```

### 6. ISO/Push-Pull Selection
```
┌─────────────────────┐
│  KODAK TRI-X 400    │
│                     │
│ ► Normal (ISO 400)  │
│   Push +1 (ISO 800) │
│   Push +2 (ISO1600) │
│   Pull -1 (ISO 200) │
│                     │
│ [SEL][NAV][BACK]    │
└─────────────────────┘
```

### 7. Developer Selection
```
┌─────────────────────┐
│  TRI-X PUSH +2      │
│                     │
│ ► D-76 1+1   14:30  │
│   HC-110 B    8:00  │
│   Rodinal 1+50 11:00│
│   Xtol 1+1   12:30  │
│                     │
│ [SEL][NAV][BACK]    │
└─────────────────────┘
```

### 8. Timer Configuration Review
```
┌─────────────────────┐
│ TRI-X +2 D76 22°C   │
│                     │
│ Film: Tri-X Push+2  │
│ Dev: D-76 1+1       │
│ Temp: 22°C          │
│ Time: 14:30         │
│                     │
│ [START][EDIT][BACK] │
└─────────────────────┘
```

### 9. Timer Running (Developer Stage)
```
┌─────────────────────┐
│    DEVELOPING...    │
│                     │
│ TRX+2 D76 22°C      │
│ 14:30 ███████████▒▒ │
│ 12:30 remaining     │
│ Temp: 22.1°C ✓      │
│ Run: #2024-1108-001 │
│[PAUS][EDIT][NEXT]   │
└─────────────────────┘
```

**Elements**:
- Stage: Current development stage
- Configuration: Quick reference
- Progress Bar: Visual time remaining (█ = elapsed, ▒ = remaining)
- Time: Remaining countdown
- Temperature: Current with status
- Run ID: Session identifier
- Buttons: PAUSE, EDIT, NEXT

### 10. Timer Paused
```
┌─────────────────────┐
│   DEVELOPER PAUSED  │
│                     │
│ TRX+2 D76 22°C      │
│ 14:30 ███████████▒▒ │
│ 12:30 remaining     │
│ Temp: 22.1°C ✓      │
│ Run: #2024-1108-001 │
│[RESU][EDIT][NEXT]   │
└─────────────────────┘
```

**Visual Indicators**:
- Screen: Slow blink (500ms)
- LEDs: Slow blink pattern
- Title: Shows "PAUSED" status

### 11. Timer Edit Mode
```
┌─────────────────────┐
│   EDIT DEV TIME     │
│                     │
│ Current: 14:30      │
│ Adjust: ► +0:30     │
│ New: 15:00          │
│                     │
│                     │
│[APPL][ +- ][CANC]   │
└─────────────────────┘
```

**Adjustment Options**:
- +0:30, +1:00, +2:00, +5:00
- -0:30, -1:00, -2:00, -5:00

### 12. Timer Completion Alert
```
┌─────────────────────┐
│     TIMER DONE!     │
│                     │
│ TRX+2 D76 22°C      │
│ 14:30 ████████████▓ │
│ COMPLETE            │
│ Temp: 22.0°C ✓      │
│                     │
│[    ][    ][NEXT]   │
└─────────────────────┘
```

**Visual Alerts**:
- Screen: Fast blink (200ms)
- LEDs: All fast blink
- Progress bar: Complete (all █)

### 13. Overtime Warning
```
┌─────────────────────┐
│     TIMER DONE!     │
│                     │
│ TRX+2 D76 22°C      │
│ 14:30 ████████████▓ │
│ +00:10 OVERTIME!    │
│ Temp: 22.0°C ✓      │
│                     │
│[    ][    ][NEXT]   │
└─────────────────────┘
```

**Critical Alerts**:
- Screen: Very fast blink (100ms)
- LEDs: All very fast blink
- Overtime: Red text indication (if color available)

### 14. Stop Bath Stage
```
┌─────────────────────┐
│    STOP BATH        │
│                     │
│ Standard Stop       │
│ 01:00 ████████████▒ │
│ 00:45 remaining     │
│ Temp: 22.1°C        │
│ Run: #2024-1108-001 │
│[PAUS][EDIT][NEXT]   │
└─────────────────────┘
```

### 15. Fixer Stage
```
┌─────────────────────┐
│      FIXER          │
│                     │
│ Standard Fixer      │
│ 10:00 ████████████▒ │
│ 09:55 remaining     │
│ Temp: 22.0°C        │
│ Run: #2024-1108-001 │
│[PAUS][EDIT][NEXT]   │
└─────────────────────┘
```

### 16. Wash Stage
```
┌─────────────────────┐
│       WASH          │
│                     │
│ Final Wash          │
│ 15:00 ████████████▒ │
│ 14:58 remaining     │
│ Temp: 21.9°C        │
│ Run: #2024-1108-001 │
│[PAUS][EDIT][NEXT]   │
└─────────────────────┘
```

### 17. Development Complete
```
┌─────────────────────┐
│  DEVELOPMENT DONE!  │
│                     │
│ Run: #2024-1108-001 │
│ Total: 31:25        │
│ Stages: 4/4 ✓       │
│ Temp Avg: 22.1°C    │
│                     │
│[    ][MENU][    ]   │
└─────────────────────┘
```

### 18. System Information
```
┌─────────────────────┐
│    SYSTEM INFO      │
│                     │
│ Firmware: v1.2.3    │
│ Uptime: 00:02:15    │
│ Temp Sensor: OK     │
│ Last Run: 001       │
│                     │
│[    ][    ][BACK]   │
└─────────────────────┘
```

### 19. Favourites List
```
┌─────────────────────┐
│    FAVOURITES       │
│                     │
│ ► TRX+2_D76_22      │
│   HP5_HC110_20      │
│   ACROS_XTOL_21     │
│   (empty)           │
│                     │
│ [SEL][NAV][BACK]    │
└─────────────────────┘
```

## Navigation Flow

### Main Navigation Tree
```
Main Screen
├── START → Timer Running
├── MENU → Configuration
│   ├── Timer Presets
│   │   ├── Favourites → Timer Config
│   │   ├── Film Database → Film Selection → ISO Selection → Developer Selection → Timer Config
│   │   ├── Manual Entry → Timer Config
│   │   └── Recent Timers → Timer Config
│   ├── Temperature Cal → Calibration Screen
│   └── System Info → Information Display
└── (disabled)

Timer Running
├── PAUSE → Timer Paused
├── EDIT → Edit Mode → Timer Running
└── NEXT → Next Stage / Complete

Timer Paused
├── RESUME → Timer Running
├── EDIT → Edit Mode → Timer Paused
└── NEXT → Next Stage / Complete
```

## Interaction Patterns

### Button Press Types
- **Short Press** (< 500ms): Standard action
- **Long Press** (> 2000ms): Special actions (enter config mode)
- **Double Press**: Not implemented (too complex for 3-button interface)

### Visual Feedback
- **Button Press**: LED flash (100ms)
- **Menu Selection**: Immediate highlight change
- **Timer Actions**: Progress bar update
- **Alerts**: Screen/LED blinking patterns

### Error States
- **Temperature Sensor Error**: "Temp: ERROR" display
- **Configuration Error**: "CONFIG ERROR" with error code
- **System Error**: "SYSTEM ERROR" with restart prompt

## Display Characteristics

### Visual Design
- High contrast OLED display (white on black)
- 6x8 pixel fonts for readability
- Status symbols (✓, !, ERROR)
- Progress bars for time visualization
- LED indicators for button states

### Interaction Design
- Physical button press with 50ms debouncing
- Immediate visual feedback on button press
- Clear state transitions between screens
- LED blink patterns for alerts