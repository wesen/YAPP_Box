#!/bin/bash
# Hook script for beforeSubmitPrompt - logs to SQLite and allows submission

echo "beforeSubmitPrompt" >> log-hooks.txt

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
JSON_INPUT=$(cat)

# Log to database
echo "$JSON_INPUT" | "$SCRIPT_DIR/log_hook.sh"

# Allow submission (non-blocking hook)
cat <<EOF
{
  "continue": true
}
EOF

exit 0

