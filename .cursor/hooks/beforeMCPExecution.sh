#!/bin/bash
# Hook script for beforeMCPExecution - logs to SQLite and allows execution

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
JSON_INPUT=$(cat)

# Log to database
echo "$JSON_INPUT" | "$SCRIPT_DIR/log_hook.sh"

# Allow execution (non-blocking hook)
cat <<EOF
{
  "permission": "allow"
}
EOF

exit 0

