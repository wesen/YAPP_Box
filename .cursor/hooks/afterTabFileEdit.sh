#!/bin/bash
# Hook script for afterTabFileEdit - logs to SQLite

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
JSON_INPUT=$(cat)

# Log to database
echo "$JSON_INPUT" | "$SCRIPT_DIR/log_hook.sh" tab_file_edits

exit 0

