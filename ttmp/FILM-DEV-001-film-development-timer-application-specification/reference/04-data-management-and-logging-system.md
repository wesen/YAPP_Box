---
Title: Data Management and Logging System
Ticket: FILM-DEV-001
Status: active
Topics:
    - embedded
    - timer
    - hardware
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Session logging, data persistence, and configuration management for the film development timer
LastUpdated: 2025-11-08T16:51:17.999387756-05:00
---

# Data Management and Logging System

## Data Architecture Overview

### Storage Hierarchy
```
Flash Storage (2MB total)
├── System (500KB)
│   ├── MicroPython Runtime (300KB)
│   ├── Application Code (150KB)
│   └── System Libraries (50KB)
├── Configuration (50KB)
│   ├── config.json (10KB)
│   ├── film_database.json (30KB)
│   └── favourites.json (10KB)
├── Logs (1.4MB)
│   ├── session_logs/ (1MB)
│   ├── temperature_logs/ (300KB)
│   └── system_logs/ (100KB)
└── Reserved (50KB)
    └── Emergency recovery space
```

### Data Flow Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Application   │───▶│  Data Manager   │───▶│  Flash Storage  │
│     Layer       │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       ▼                       │
         │              ┌─────────────────┐              │
         └─────────────▶│   RAM Cache     │◀─────────────┘
                        │   (Active Data) │
                        └─────────────────┘
```

## Session Logging System

### Run ID Generation
```python
def generate_run_id():
    """
    Generate unique run ID: YYYY-MMDD-NNN
    Example: 2024-1108-001
    """
    from time import localtime
    
    now = localtime()
    date_str = f"{now[0]}-{now[1]:02d}{now[2]:02d}"
    
    # Find next sequence number for today
    existing_runs = get_runs_for_date(date_str)
    next_seq = len(existing_runs) + 1
    
    return f"{date_str}-{next_seq:03d}"
```

### Session Log Structure
```json
{
  "run_id": "2024-1108-001",
  "start_time": "2024-11-08T14:32:15Z",
  "end_time": "2024-11-08T15:03:42Z",
  "settings": {
    "film": {
      "name": "Kodak Tri-X 400",
      "processing": "Push +2 (ISO 1600)",
      "base_iso": 400,
      "effective_iso": 1600
    },
    "developer": {
      "name": "D-76 1+1",
      "dilution": "1+1",
      "base_time_20c": "16:00"
    },
    "target_temperature": 22.0,
    "planned_total_time": "31:30"
  },
  "stages": [
    {
      "stage": "developer",
      "planned_duration": "14:30",
      "actual_duration": "14:40",
      "start_time": "2024-11-08T14:32:15Z",
      "end_time": "2024-11-08T14:46:55Z",
      "adjustments": [
        {
          "timestamp": "2024-11-08T14:40:00Z",
          "adjustment": "+0:30",
          "reason": "user_edit",
          "new_total": "15:00"
        }
      ],
      "pause_events": [
        {
          "pause_start": "2024-11-08T14:42:30Z",
          "pause_end": "2024-11-08T14:43:15Z",
          "duration": "00:45"
        }
      ],
      "temperature_stats": {
        "min": 21.8,
        "max": 22.3,
        "avg": 22.05,
        "readings_count": 29
      }
    }
  ],
  "log_entries": [
    {
      "timestamp": "2024-11-08T14:32:15Z",
      "action": "session_start",
      "stage": null,
      "data": {
        "temperature": 22.1,
        "run_id": "2024-1108-001"
      }
    },
    {
      "timestamp": "2024-11-08T14:32:15Z", 
      "action": "timer_start",
      "stage": "developer",
      "data": {
        "planned_duration": "14:30",
        "temperature": 22.1
      }
    }
  ],
  "summary": {
    "total_duration": "31:27",
    "stages_completed": 4,
    "temperature_avg": 22.05,
    "adjustments_made": 1,
    "pause_count": 1,
    "total_pause_time": "00:45"
  }
}
```

### Log Entry Types

#### Core Actions
```python
class LogAction:
    SESSION_START = "session_start"
    SESSION_END = "session_end"
    TIMER_START = "timer_start"
    TIMER_PAUSE = "timer_pause"
    TIMER_RESUME = "timer_resume"
    TIMER_COMPLETE = "timer_complete"
    TIMER_ADJUST = "timer_adjust"
    STAGE_TRANSITION = "stage_transition"
    TEMPERATURE_READING = "temperature_reading"
    USER_ACTION = "user_action"
    SYSTEM_EVENT = "system_event"
    ERROR_EVENT = "error_event"
```

#### Temperature Logging
```json
{
  "timestamp": "2024-11-08T14:35:00Z",
  "action": "temperature_reading",
  "stage": "developer",
  "data": {
    "temperature": 22.0,
    "sensor_status": "ok",
    "reading_time_ms": 750
  }
}
```

#### User Actions
```json
{
  "timestamp": "2024-11-08T14:40:00Z",
  "action": "user_action",
  "stage": "developer", 
  "data": {
    "button": "middle",
    "action_type": "timer_adjust",
    "adjustment": "+0:30",
    "new_remaining": "10:30"
  }
}
```

## Configuration Management

### System Configuration (config.json)
```json
{
  "version": "1.0.0",
  "hardware": {
    "display": {
      "width": 128,
      "height": 64,
      "i2c_address": "0x3C",
      "contrast": 255
    },
    "temperature_sensor": {
      "pin": 2,
      "calibration_offset": 0.0,
      "reading_interval_ms": 30000
    },
    "buttons": {
      "pin1": 16,
      "pin2": 17, 
      "pin3": 18,
      "debounce_ms": 50
    },
    "leds": {
      "pin1": 19,
      "pin2": 20,
      "pin3": 21
    }
  },
  "application": {
    "default_target_temp": 20.0,
    "temp_warning_threshold": 2.0,
    "temp_critical_threshold": 3.0,
    "auto_save_interval_ms": 10000,
    "max_log_entries": 1000,
    "log_retention_days": 90
  },
  "alerts": {
    "final_warning_seconds": 30,
    "completion_blink_rate_ms": 100,
    "warning_blink_rate_ms": 500,
    "overtime_blink_rate_ms": 50
  }
}
```

### Film Database (film_database.json)
```json
{
  "version": "1.0.0",
  "last_updated": "2024-11-08T12:00:00Z",
  "films": {
    "kodak_tri_x_400": {
      "name": "Kodak Tri-X 400",
      "manufacturer": "Kodak",
      "base_iso": 400,
      "film_type": "black_white",
      "processing_options": {
        "normal_iso_400": {
          "name": "Normal (ISO 400)",
          "push_pull": 0,
          "effective_iso": 400,
          "developers": {
            "d76_1_1": {
              "name": "D-76 1+1",
              "base_time_20c": "11:00",
              "temp_compensation": "standard",
              "notes": "General purpose developer"
            }
          }
        }
      }
    }
  },
  "developers": {
    "d76_1_1": {
      "name": "D-76 1+1",
      "manufacturer": "Kodak",
      "type": "general_purpose",
      "dilution": "1+1",
      "working_solution_life_hours": 8,
      "temperature_range": {
        "min": 18,
        "max": 24,
        "optimal": 20
      }
    }
  },
  "standard_stages": {
    "stop_bath": {
      "name": "Stop Bath",
      "default_time": "00:30",
      "temperature_dependent": false,
      "optional": false
    },
    "fixer": {
      "name": "Fixer",
      "default_time": "10:00", 
      "temperature_dependent": false,
      "optional": false
    },
    "wash": {
      "name": "Final Wash",
      "default_time": "15:00",
      "temperature_dependent": false,
      "optional": false
    }
  }
}
```

### Favourites (favourites.json)
```json
{
  "version": "1.0.0",
  "max_favourites": 10,
  "favourites": [
    {
      "id": "fav_001",
      "name": "TRX+2_D76_22",
      "created": "2024-11-08T14:47:02Z",
      "last_used": "2024-11-08T14:47:02Z",
      "use_count": 3,
      "configuration": {
        "film": "kodak_tri_x_400",
        "processing": "push_2_iso_1600",
        "developer": "d76_1_1",
        "target_temperature": 22.0,
        "custom_adjustments": {
          "developer_time": null,
          "stop_bath_time": null,
          "fixer_time": null,
          "wash_time": null
        }
      }
    }
  ]
}
```

## Data Persistence Layer

### File Operations
```python
class DataManager:
    def __init__(self):
        self.config_file = "config.json"
        self.film_db_file = "film_database.json"
        self.favourites_file = "favourites.json"
        self.log_dir = "logs/sessions/"
        
    def save_json(self, filename, data):
        """Atomic save with backup"""
        backup_file = f"{filename}.backup"
        temp_file = f"{filename}.tmp"
        
        try:
            # Write to temporary file
            with open(temp_file, 'w') as f:
                json.dump(data, f, indent=2)
            
            # Backup existing file
            if file_exists(filename):
                rename(filename, backup_file)
            
            # Move temp to final location
            rename(temp_file, filename)
            
            # Remove backup on success
            if file_exists(backup_file):
                remove(backup_file)
                
        except Exception as e:
            # Restore from backup on failure
            if file_exists(backup_file):
                rename(backup_file, filename)
            raise e
```

### Memory Management
```python
class SessionLogger:
    def __init__(self, max_memory_entries=100):
        self.memory_buffer = []
        self.max_memory_entries = max_memory_entries
        self.current_run_id = None
        
    def log_entry(self, action, stage, data):
        entry = {
            "timestamp": get_iso_timestamp(),
            "action": action,
            "stage": stage,
            "data": data
        }
        
        self.memory_buffer.append(entry)
        
        # Flush to disk when buffer is full
        if len(self.memory_buffer) >= self.max_memory_entries:
            self.flush_to_disk()
    
    def flush_to_disk(self):
        if not self.memory_buffer:
            return
            
        log_file = f"logs/sessions/{self.current_run_id}.json"
        
        # Load existing log or create new
        existing_log = self.load_session_log(self.current_run_id)
        existing_log["log_entries"].extend(self.memory_buffer)
        
        # Save updated log
        self.save_json(log_file, existing_log)
        
        # Clear memory buffer
        self.memory_buffer = []
```

## Data Retention and Cleanup

### Automatic Cleanup Policy
```python
def cleanup_old_logs():
    """Remove logs older than retention period"""
    retention_days = get_config("application.log_retention_days", 90)
    cutoff_date = time.time() - (retention_days * 24 * 3600)
    
    log_files = list_files("logs/sessions/")
    
    for log_file in log_files:
        file_stat = stat(log_file)
        if file_stat.st_mtime < cutoff_date:
            remove(log_file)
            print(f"Removed old log: {log_file}")
```

### Storage Monitoring
```python
def check_storage_space():
    """Monitor flash storage usage"""
    total_space = get_flash_size()
    used_space = get_used_space()
    free_space = total_space - used_space
    
    # Warning at 90% full
    if free_space < (total_space * 0.1):
        return StorageStatus.WARNING
    
    # Critical at 95% full  
    if free_space < (total_space * 0.05):
        return StorageStatus.CRITICAL
        
    return StorageStatus.OK
```

## Data Export and Import

### Session Data Export Format
```json
{
  "export_info": {
    "version": "1.0.0",
    "exported_at": "2024-11-08T16:00:00Z",
    "device_id": "pico_w_001",
    "total_sessions": 25
  },
  "sessions": [
    {
      "run_id": "2024-1108-001",
      "summary": {
        "film": "Kodak Tri-X 400",
        "processing": "Push +2",
        "developer": "D-76 1+1",
        "total_time": "31:27",
        "temperature_avg": 22.05
      },
      "detailed_log": "..." 
    }
  ]
}
```

### Configuration Backup
```python
def create_config_backup():
    """Create complete configuration backup"""
    backup_data = {
        "backup_info": {
            "created_at": get_iso_timestamp(),
            "firmware_version": get_firmware_version(),
            "device_id": get_device_id()
        },
        "config": load_json("config.json"),
        "film_database": load_json("film_database.json"), 
        "favourites": load_json("favourites.json")
    }
    
    backup_filename = f"backup_{get_date_string()}.json"
    save_json(backup_filename, backup_data)
    
    return backup_filename
```

## Error Handling and Recovery

### Data Corruption Recovery
```python
def recover_corrupted_file(filename):
    """Attempt to recover from corrupted JSON file"""
    backup_file = f"{filename}.backup"
    
    if file_exists(backup_file):
        print(f"Recovering {filename} from backup")
        copy_file(backup_file, filename)
        return True
    
    # Try to load default configuration
    if filename == "config.json":
        create_default_config()
        return True
    elif filename == "film_database.json":
        create_default_film_database()
        return True
    elif filename == "favourites.json":
        create_empty_favourites()
        return True
    
    return False
```

### System State Recovery
```python
def recover_system_state():
    """Recover system state after unexpected restart"""
    
    # Check for active session
    active_session_file = "active_session.json"
    
    if file_exists(active_session_file):
        try:
            active_session = load_json(active_session_file)
            
            # Determine if recovery is possible
            elapsed_since_last_update = time.time() - active_session["last_update"]
            
            if elapsed_since_last_update < 300:  # 5 minutes
                # Offer to resume session
                return SessionRecovery.RESUMABLE, active_session
            else:
                # Session too old, mark as interrupted
                return SessionRecovery.INTERRUPTED, active_session
                
        except Exception:
            # Corrupted active session file
            remove(active_session_file)
    
    return SessionRecovery.NONE, None
```

## Data Management Notes

### Simple Caching
- Load film database once at startup
- Keep favourites in memory during operation
- Reload configuration only when modified

### Basic Validation
- Check for required fields in JSON files
- Use default values for missing configuration
- Create default files if none exist

### Data Integrity
- Use atomic file writes (write to temp, then rename)
- Keep backup of previous file during save
- Restore from backup on write failure