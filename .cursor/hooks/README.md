# Cursor Hooks Dashboard

## Usage

### Live Monitoring (default)
```bash
.cursor/hooks/dashboard.sh
```

This will continuously update the dashboard every 2 seconds, showing:
- Statistics by hook type
- Recent events (last 10)
- Recent agent responses (last 5)
- Activity timeline (last hour)

Press `Ctrl+C` to exit.

### Single Run
```bash
.cursor/hooks/dashboard.sh --once
```

This runs the dashboard once and exits (useful for scripting or one-time checks).

## Database Location

The SQLite database is stored at:
- Default: `~/.cursor/hooks.db`
- Custom: Set `CURSOR_HOOKS_DB` environment variable

## Features

- **Real-time updates**: Automatically refreshes every 2 seconds
- **Color-coded output**: Easy to read with color highlighting
- **Multiple views**: Statistics, recent events, responses, and timeline
- **Non-intrusive**: Only reads data, never modifies the database

## Example Output

```
═══════════════════════════════════════════════════════════════
  Cursor Hooks Dashboard
  Database: /home/user/.cursor/hooks.db
  Last Update: 2025-11-29 13:15:30
═══════════════════════════════════════════════════════════════

📊 Statistics
─────────────────────────────────────────────────────────────
Hook Events by Type:
  Agent Responses           15
  Shell Executions          8
  File Operations           3
  ...

Total Events: 26

🕐 Recent Events (Last 10)
─────────────────────────────────────────────────────────────
Time                 Type      Conversation                    Model
2025-11-29 13:15:25  Response  abc123-def456                   gpt-4
2025-11-29 13:15:20  Shell     abc123-def456                   gpt-4
...

💬 Recent Agent Responses (Last 5)
─────────────────────────────────────────────────────────────
[2025-11-29 13:15:25] Conv: abc123-def456
  I've created the hooks system with SQLite logging...

📈 Activity Timeline (Last Hour)
─────────────────────────────────────────────────────────────
  13:15: 3 events
  13:14: 5 events
  13:13: 2 events
```

