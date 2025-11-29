#!/bin/bash
# Hook script for beforeShellExecution - logs to SQLite and allows execution

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
JSON_INPUT=$(cat)

# Log to database
echo "$JSON_INPUT" | "$SCRIPT_DIR/log_hook.sh" shell_executions

# Allow execution (non-blocking hook)
cat <<EOF
{
  "permission": "allow"
}
EOF

exit 0

