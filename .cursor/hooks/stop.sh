#!/bin/bash
# Hook script for stop - logs to SQLite

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
JSON_INPUT=$(cat)

# Log to database
echo "$JSON_INPUT" | "$SCRIPT_DIR/log_hook.sh" agent_stops

# No output needed for stop hook (optional followup_message can be added if needed)
exit 0

