#!/bin/bash
# Hook script for beforeReadFile - logs to SQLite and allows read

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
JSON_INPUT=$(cat)

# Log to database
echo "$JSON_INPUT" | "$SCRIPT_DIR/log_hook.sh" file_operations

# Allow read (non-blocking hook)
cat <<EOF
{
  "permission": "allow"
}
EOF

exit 0

