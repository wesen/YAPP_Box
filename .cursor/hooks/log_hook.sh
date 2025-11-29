#!/bin/bash
# Generic hook logging script that stores hook data in SQLite
# Usage: log_hook.sh
# This script reads JSON from stdin and inserts it into the hook_invocations table

DB_PATH="${CURSOR_HOOKS_DB:-$HOME/.cursor/hooks.db}"

# Ensure database is initialized
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
"$SCRIPT_DIR/init_db.sh" > /dev/null 2>&1

# Read JSON input from stdin
JSON_INPUT=$(cat)

# Get current working directory
PWD=$(pwd)

# Extract common fields from JSON for indexing
CONVERSATION_ID=$(echo "$JSON_INPUT" | jq -r '.conversation_id // empty')
GENERATION_ID=$(echo "$JSON_INPUT" | jq -r '.generation_id // empty')
MODEL=$(echo "$JSON_INPUT" | jq -r '.model // empty')
HOOK_EVENT_NAME=$(echo "$JSON_INPUT" | jq -r '.hook_event_name // empty')
CURSOR_VERSION=$(echo "$JSON_INPUT" | jq -r '.cursor_version // empty')
WORKSPACE_ROOTS=$(echo "$JSON_INPUT" | jq -r '.workspace_roots | if type == "array" then join(",") else . end // empty')
USER_EMAIL=$(echo "$JSON_INPUT" | jq -r '.user_email // empty')

# Escape single quotes for SQL
escape_sql() {
    echo "$1" | sed "s/'/''/g"
}

CONVERSATION_ID_ESC=$(escape_sql "$CONVERSATION_ID")
GENERATION_ID_ESC=$(escape_sql "$GENERATION_ID")
MODEL_ESC=$(escape_sql "$MODEL")
HOOK_EVENT_NAME_ESC=$(escape_sql "$HOOK_EVENT_NAME")
CURSOR_VERSION_ESC=$(escape_sql "$CURSOR_VERSION")
WORKSPACE_ROOTS_ESC=$(escape_sql "$WORKSPACE_ROOTS")
USER_EMAIL_ESC=$(escape_sql "$USER_EMAIL")
PWD_ESC=$(escape_sql "$PWD")
RAW_JSON_ESC=$(escape_sql "$JSON_INPUT")

# Insert into the generic hook_invocations table
sqlite3 "$DB_PATH" <<EOF
INSERT INTO hook_invocations (
    hook_event_name, conversation_id, generation_id, model,
    cursor_version, workspace_roots, user_email, pwd, raw_json
) VALUES (
    '$HOOK_EVENT_NAME_ESC', '$CONVERSATION_ID_ESC', '$GENERATION_ID_ESC', '$MODEL_ESC',
    '$CURSOR_VERSION_ESC', '$WORKSPACE_ROOTS_ESC', '$USER_EMAIL_ESC', '$PWD_ESC', '$RAW_JSON_ESC'
);
EOF

exit 0
