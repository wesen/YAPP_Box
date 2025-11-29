#!/bin/bash
# Hook script for afterShellExecution - logs to SQLite

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
JSON_INPUT=$(cat)

# Log to database
echo "$JSON_INPUT" | "$SCRIPT_DIR/log_hook.sh" shell_executions

exit 0

