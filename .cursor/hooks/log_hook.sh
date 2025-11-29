#!/bin/bash
# Generic hook logging script that stores hook data in SQLite
# Usage: log_hook.sh <table_name> <json_input>
# This script reads JSON from stdin and inserts it into the specified table

TABLE_NAME="$1"
DB_PATH="${CURSOR_HOOKS_DB:-$HOME/.cursor/hooks.db}"

# Ensure database is initialized
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
"$SCRIPT_DIR/init_db.sh" > /dev/null 2>&1

# Read JSON input from stdin
JSON_INPUT=$(cat)

# Extract common fields from JSON
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

# Insert into appropriate table based on hook type
case "$TABLE_NAME" in
    shell_executions)
        COMMAND=$(echo "$JSON_INPUT" | jq -r '.command // empty')
        CWD=$(echo "$JSON_INPUT" | jq -r '.cwd // empty')
        OUTPUT=$(echo "$JSON_INPUT" | jq -r '.output // empty')
        DURATION=$(echo "$JSON_INPUT" | jq -r '.duration // empty')
        
        COMMAND_ESC=$(escape_sql "$COMMAND")
        CWD_ESC=$(escape_sql "$CWD")
        OUTPUT_ESC=$(escape_sql "$OUTPUT")
        
        sqlite3 "$DB_PATH" <<EOF
INSERT INTO shell_executions (
    hook_event_name, conversation_id, generation_id, model,
    command, cwd, output, duration_ms,
    cursor_version, workspace_roots, user_email
) VALUES (
    '$HOOK_EVENT_NAME_ESC', '$CONVERSATION_ID_ESC', '$GENERATION_ID_ESC', '$MODEL_ESC',
    '$COMMAND_ESC', '$CWD_ESC', '$OUTPUT_ESC', $DURATION,
    '$CURSOR_VERSION_ESC', '$WORKSPACE_ROOTS_ESC', '$USER_EMAIL_ESC'
);
EOF
        ;;
    
    mcp_executions)
        TOOL_NAME=$(echo "$JSON_INPUT" | jq -r '.tool_name // empty')
        TOOL_INPUT=$(echo "$JSON_INPUT" | jq -c '.tool_input // empty')
        RESULT_JSON=$(echo "$JSON_INPUT" | jq -c '.result_json // empty')
        DURATION=$(echo "$JSON_INPUT" | jq -r '.duration // empty')
        URL=$(echo "$JSON_INPUT" | jq -r '.url // empty')
        COMMAND=$(echo "$JSON_INPUT" | jq -r '.command // empty')
        
        TOOL_NAME_ESC=$(escape_sql "$TOOL_NAME")
        TOOL_INPUT_ESC=$(escape_sql "$TOOL_INPUT")
        RESULT_JSON_ESC=$(escape_sql "$RESULT_JSON")
        URL_ESC=$(escape_sql "$URL")
        COMMAND_ESC=$(escape_sql "$COMMAND")
        
        sqlite3 "$DB_PATH" <<EOF
INSERT INTO mcp_executions (
    hook_event_name, conversation_id, generation_id, model,
    tool_name, tool_input, result_json, duration_ms, url, command,
    cursor_version, workspace_roots, user_email
) VALUES (
    '$HOOK_EVENT_NAME_ESC', '$CONVERSATION_ID_ESC', '$GENERATION_ID_ESC', '$MODEL_ESC',
    '$TOOL_NAME_ESC', '$TOOL_INPUT_ESC', '$RESULT_JSON_ESC', $DURATION, '$URL_ESC', '$COMMAND_ESC',
    '$CURSOR_VERSION_ESC', '$WORKSPACE_ROOTS_ESC', '$USER_EMAIL_ESC'
);
EOF
        ;;
    
    file_operations)
        FILE_PATH=$(echo "$JSON_INPUT" | jq -r '.file_path // empty')
        CONTENT=$(echo "$JSON_INPUT" | jq -r '.content // empty')
        EDITS=$(echo "$JSON_INPUT" | jq -c '.edits // []')
        
        FILE_PATH_ESC=$(escape_sql "$FILE_PATH")
        CONTENT_ESC=$(escape_sql "$CONTENT")
        EDITS_ESC=$(escape_sql "$EDITS")
        
        sqlite3 "$DB_PATH" <<EOF
INSERT INTO file_operations (
    hook_event_name, conversation_id, generation_id, model,
    file_path, content, edits,
    cursor_version, workspace_roots, user_email
) VALUES (
    '$HOOK_EVENT_NAME_ESC', '$CONVERSATION_ID_ESC', '$GENERATION_ID_ESC', '$MODEL_ESC',
    '$FILE_PATH_ESC', '$CONTENT_ESC', '$EDITS_ESC',
    '$CURSOR_VERSION_ESC', '$WORKSPACE_ROOTS_ESC', '$USER_EMAIL_ESC'
);
EOF
        ;;
    
    prompt_submissions)
        PROMPT=$(echo "$JSON_INPUT" | jq -r '.prompt // empty')
        ATTACHMENTS=$(echo "$JSON_INPUT" | jq -c '.attachments // []')
        
        PROMPT_ESC=$(escape_sql "$PROMPT")
        ATTACHMENTS_ESC=$(escape_sql "$ATTACHMENTS")
        
        sqlite3 "$DB_PATH" <<EOF
INSERT INTO prompt_submissions (
    hook_event_name, conversation_id, generation_id, model,
    prompt, attachments,
    cursor_version, workspace_roots, user_email
) VALUES (
    '$HOOK_EVENT_NAME_ESC', '$CONVERSATION_ID_ESC', '$GENERATION_ID_ESC', '$MODEL_ESC',
    '$PROMPT_ESC', '$ATTACHMENTS_ESC',
    '$CURSOR_VERSION_ESC', '$WORKSPACE_ROOTS_ESC', '$USER_EMAIL_ESC'
);
EOF
        ;;
    
    agent_stops)
        STATUS=$(echo "$JSON_INPUT" | jq -r '.status // empty')
        LOOP_COUNT=$(echo "$JSON_INPUT" | jq -r '.loop_count // 0')
        
        STATUS_ESC=$(escape_sql "$STATUS")
        
        sqlite3 "$DB_PATH" <<EOF
INSERT INTO agent_stops (
    hook_event_name, conversation_id, generation_id, model,
    status, loop_count,
    cursor_version, workspace_roots, user_email
) VALUES (
    '$HOOK_EVENT_NAME_ESC', '$CONVERSATION_ID_ESC', '$GENERATION_ID_ESC', '$MODEL_ESC',
    '$STATUS_ESC', $LOOP_COUNT,
    '$CURSOR_VERSION_ESC', '$WORKSPACE_ROOTS_ESC', '$USER_EMAIL_ESC'
);
EOF
        ;;
    
    agent_responses)
        TEXT=$(echo "$JSON_INPUT" | jq -r '.text // empty')
        
        TEXT_ESC=$(escape_sql "$TEXT")
        
        sqlite3 "$DB_PATH" <<EOF
INSERT INTO agent_responses (
    hook_event_name, conversation_id, generation_id, model,
    text,
    cursor_version, workspace_roots, user_email
) VALUES (
    '$HOOK_EVENT_NAME_ESC', '$CONVERSATION_ID_ESC', '$GENERATION_ID_ESC', '$MODEL_ESC',
    '$TEXT_ESC',
    '$CURSOR_VERSION_ESC', '$WORKSPACE_ROOTS_ESC', '$USER_EMAIL_ESC'
);
EOF
        ;;
    
    agent_thoughts)
        TEXT=$(echo "$JSON_INPUT" | jq -r '.text // empty')
        DURATION=$(echo "$JSON_INPUT" | jq -r '.duration_ms // empty')
        
        TEXT_ESC=$(escape_sql "$TEXT")
        
        sqlite3 "$DB_PATH" <<EOF
INSERT INTO agent_thoughts (
    hook_event_name, conversation_id, generation_id, model,
    text, duration_ms,
    cursor_version, workspace_roots, user_email
) VALUES (
    '$HOOK_EVENT_NAME_ESC', '$CONVERSATION_ID_ESC', '$GENERATION_ID_ESC', '$MODEL_ESC',
    '$TEXT_ESC', $DURATION,
    '$CURSOR_VERSION_ESC', '$WORKSPACE_ROOTS_ESC', '$USER_EMAIL_ESC'
);
EOF
        ;;
    
    tab_file_reads)
        FILE_PATH=$(echo "$JSON_INPUT" | jq -r '.file_path // empty')
        CONTENT=$(echo "$JSON_INPUT" | jq -r '.content // empty')
        
        FILE_PATH_ESC=$(escape_sql "$FILE_PATH")
        CONTENT_ESC=$(escape_sql "$CONTENT")
        
        sqlite3 "$DB_PATH" <<EOF
INSERT INTO tab_file_reads (
    hook_event_name, conversation_id, generation_id, model,
    file_path, content,
    cursor_version, workspace_roots, user_email
) VALUES (
    '$HOOK_EVENT_NAME_ESC', '$CONVERSATION_ID_ESC', '$GENERATION_ID_ESC', '$MODEL_ESC',
    '$FILE_PATH_ESC', '$CONTENT_ESC',
    '$CURSOR_VERSION_ESC', '$WORKSPACE_ROOTS_ESC', '$USER_EMAIL_ESC'
);
EOF
        ;;
    
    tab_file_edits)
        FILE_PATH=$(echo "$JSON_INPUT" | jq -r '.file_path // empty')
        EDITS=$(echo "$JSON_INPUT" | jq -c '.edits // []')
        
        FILE_PATH_ESC=$(escape_sql "$FILE_PATH")
        EDITS_ESC=$(escape_sql "$EDITS")
        
        sqlite3 "$DB_PATH" <<EOF
INSERT INTO tab_file_edits (
    hook_event_name, conversation_id, generation_id, model,
    file_path, edits,
    cursor_version, workspace_roots, user_email
) VALUES (
    '$HOOK_EVENT_NAME_ESC', '$CONVERSATION_ID_ESC', '$GENERATION_ID_ESC', '$MODEL_ESC',
    '$FILE_PATH_ESC', '$EDITS_ESC',
    '$CURSOR_VERSION_ESC', '$WORKSPACE_ROOTS_ESC', '$USER_EMAIL_ESC'
);
EOF
        ;;
    
    *)
        echo "Unknown table: $TABLE_NAME" >&2
        exit 1
        ;;
esac

exit 0

