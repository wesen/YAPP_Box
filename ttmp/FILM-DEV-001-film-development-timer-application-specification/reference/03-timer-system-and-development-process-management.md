---
Title: Timer System and Development Process Management
Ticket: FILM-DEV-001
Status: active
Topics:
    - timer
    - embedded
    - hardware
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Timer logic, film database, and development workflows for the film development timer
LastUpdated: 2025-11-08T16:51:17.999387756-05:00
---

# Timer System and Development Process Management

## Timer Engine Architecture

### Core Timer States
```
READY ──START──> RUNNING ──PAUSE──> PAUSED
  ↑                 │                  │
  │                 │                  │
  │              COMPLETE           RESUME
  │                 │                  │
  │                 ↓                  ↓
  └──────────── COMPLETE ←──────── RUNNING
                    │
                    │
                 NEXT_STAGE
                    │
                    ↓
                 READY (next stage)
```

### Timer State Definitions

#### READY State
- Timer configured with duration and settings
- Temperature monitoring active
- Waiting for START command
- All parameters can be edited

#### RUNNING State  
- Timer actively counting down
- Progress bar updating
- Temperature logging active
- PAUSE and NEXT commands available
- Real-time adjustments possible via EDIT

#### PAUSED State
- Timer stopped but retains position
- Visual indicators show paused state
- RESUME returns to exact position
- EDIT adjustments still possible
- Pause duration tracked separately

#### COMPLETE State
- Timer reached zero
- Alert notifications active
- Overtime tracking if not acknowledged
- NEXT advances to next stage

### Multi-Stage Development Process

#### Standard Development Sequence
1. **Developer** - Primary development stage (film-specific timing)
2. **Stop Bath** - Halts development (typically 30-60 seconds)
3. **Fixer** - Makes image permanent (5-15 minutes)
4. **Wash** - Removes chemicals (10-30 minutes)

#### Stage Transition Logic
```python
class DevelopmentStage:
    DEVELOPER = "developer"
    STOP_BATH = "stop_bath"
    FIXER = "fixer"
    WASH = "wash"
    COMPLETE = "complete"

stage_sequence = [
    DEVELOPER,
    STOP_BATH, 
    FIXER,
    WASH,
    COMPLETE
]
```

## Film Database System

### Database Structure
```json
{
  "films": {
    "kodak_tri_x_400": {
      "name": "Kodak Tri-X 400",
      "base_iso": 400,
      "processing_options": {
        "normal_iso_400": {
          "name": "Normal (ISO 400)",
          "developers": {
            "d76_1_1": {
              "name": "D-76 1+1",
              "time": "11:00"
            },
            "hc110_b": {
              "name": "HC-110 Dilution B", 
              "time": "6:30"
            }
          }
        },
        "push_1_iso_800": {
          "name": "Push +1 (ISO 800)",
          "developers": {
            "d76_1_1": {
              "name": "D-76 1+1",
              "time": "13:30"
            }
          }
        }
      }
    }
  },
  "standard_stages": {
    "stop_bath": {
      "name": "Stop Bath",
      "default_time": "00:30",
      "temperature_dependent": false
    },
    "fixer": {
      "name": "Fixer", 
      "default_time": "10:00",
      "temperature_dependent": false
    },
    "wash": {
      "name": "Final Wash",
      "default_time": "15:00", 
      "temperature_dependent": false
    }
  }
}
```

### Temperature Monitoring

For MVP, the system monitors and displays temperature but does not automatically adjust development times. Users can:
- View current developer temperature
- Select presets appropriate for their temperature
- Manually adjust timer duration if needed

**Note**: Automatic temperature compensation is documented as a future enhancement in the future-ideas document.

### Film Database Examples

#### Kodak Tri-X 400 (at 20°C)
```json
{
  "kodak_tri_x_400": {
    "normal_iso_400": {
      "d76_1_1": "11:00",
      "hc110_b": "6:30", 
      "rodinal_1_50": "9:00",
      "xtol_1_1": "10:30"
    },
    "push_1_iso_800": {
      "d76_1_1": "13:30",
      "hc110_b": "8:00",
      "rodinal_1_50": "11:00", 
      "xtol_1_1": "12:30"
    },
    "push_2_iso_1600": {
      "d76_1_1": "16:00",
      "hc110_b": "10:00",
      "rodinal_1_50": "13:30"
    },
    "pull_1_iso_200": {
      "d76_1_1": "9:00",
      "hc110_b": "5:00",
      "xtol_1_1": "8:30"
    }
  }
}
```

**Note**: Times are for 20°C developer temperature. Users should adjust manually for different temperatures or select appropriate presets.

#### Ilford HP5+ 400 (at 20°C)
```json
{
  "ilford_hp5_400": {
    "normal_iso_400": {
      "id11_1_1": "11:00",
      "hc110_b": "7:00",
      "rodinal_1_50": "10:00",
      "xtol_1_1": "11:30"
    },
    "push_1_iso_800": {
      "id11_1_1": "14:00", 
      "hc110_b": "9:00",
      "rodinal_1_50": "13:00"
    }
  }
}
```

## Timer Adjustment System

### Real-Time Adjustments
During timer operation, users can adjust the current stage timing:

#### Adjustment Options
- **+0:30** - Add 30 seconds
- **+1:00** - Add 1 minute  
- **+2:00** - Add 2 minutes
- **+5:00** - Add 5 minutes
- **-0:30** - Subtract 30 seconds
- **-1:00** - Subtract 1 minute
- **-2:00** - Subtract 2 minutes
- **-5:00** - Subtract 5 minutes

#### Adjustment Logic
```python
def adjust_timer(current_remaining_seconds, adjustment_seconds):
    """
    Adjust timer while running, with safety bounds
    """
    new_remaining = current_remaining_seconds + adjustment_seconds
    
    # Safety bounds: minimum 10 seconds, maximum 60 minutes
    new_remaining = max(10, min(3600, new_remaining))
    
    return new_remaining
```

#### Session-Only Changes
- Adjustments apply only to current development session
- Original presets remain unchanged in database
- All adjustments logged with timestamps and reasons

## Favourites System

### Favourite Entry Structure
```json
{
  "favourites": [
    {
      "id": "trx_d76_push2_22c",
      "name": "TRX+2_D76_22",
      "film": "kodak_tri_x_400",
      "processing": "push_2_iso_1600", 
      "developer": "d76_1_1",
      "target_temperature": 22.0,
      "custom_time": null,
      "created": "2024-11-08T14:32:15Z",
      "last_used": "2024-11-08T14:47:02Z",
      "use_count": 3
    }
  ]
}
```

### Favourites Management
- **Maximum**: 10 favourite combinations
- **Auto-naming**: Film+Processing+Developer+Temp format
- **Custom naming**: User can override auto-generated names
- **Usage tracking**: Last used date and frequency
- **Quick access**: Direct selection from main presets menu

## Alert and Notification System

### Timer Completion Alerts

#### Visual Alerts
- **30 seconds remaining**: Screen and LEDs start blinking (500ms)
- **10 seconds remaining**: Faster blinking (200ms)
- **Timer complete**: Very fast blinking (100ms)
- **Overtime**: Critical fast blinking (50ms)

#### Alert Escalation
```python
def get_alert_pattern(seconds_remaining, overtime_seconds):
    if overtime_seconds > 0:
        return AlertPattern.CRITICAL_FAST  # 50ms
    elif seconds_remaining <= 0:
        return AlertPattern.COMPLETE_FAST  # 100ms
    elif seconds_remaining <= 10:
        return AlertPattern.WARNING_FAST   # 200ms
    elif seconds_remaining <= 30:
        return AlertPattern.WARNING_SLOW   # 500ms
    else:
        return AlertPattern.NONE
```

### Temperature Alerts
- **Optimal Range**: ±1°C from target (green ✓)
- **Warning Range**: ±2°C from target (yellow !)
- **Critical Range**: >±3°C from target (red ERROR)

## Process Automation

### Stage Transitions
- **Manual**: User presses NEXT to advance stages
- **Automatic**: Optional auto-advance after completion (future feature)
- **Skip**: User can skip stages if needed
- **Repeat**: User can restart current stage

### Development Workflow
```python
class DevelopmentSession:
    def __init__(self, film_config):
        self.stages = [
            TimerStage("developer", film_config.dev_time),
            TimerStage("stop_bath", "00:30"),
            TimerStage("fixer", "10:00"), 
            TimerStage("wash", "15:00")
        ]
        self.current_stage = 0
        
    def next_stage(self):
        if self.current_stage < len(self.stages) - 1:
            self.current_stage += 1
            return self.stages[self.current_stage]
        else:
            return None  # Development complete
```

## Timer Implementation Notes

### Basic Requirements
- **Resolution**: 1 second countdown accuracy
- **Temperature monitoring**: Display current temperature during operation
- **Adjustment**: Real-time modifications during operation
- **State persistence**: Save timer state to flash storage

### Error Handling
- **Sensor failures**: Display "Temp: ERROR" and continue with last known temperature
- **File system errors**: Use default configurations
- **Invalid adjustments**: Clamp to reasonable bounds (10s minimum, 60min maximum)