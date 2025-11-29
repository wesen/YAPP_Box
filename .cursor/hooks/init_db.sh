#!/bin/bash
# Shared script to initialize SQLite database for Cursor hooks logging
# This script ensures the database and tables exist before any hook script runs

DB_PATH="${CURSOR_HOOKS_DB:-$HOME/.cursor/hooks.db}"

# Create database directory if it doesn't exist
mkdir -p "$(dirname "$DB_PATH")"

# Initialize database with a single generic table for all hook invocations
sqlite3 "$DB_PATH" <<EOF
-- Single generic table for all hook invocations
CREATE TABLE IF NOT EXISTS hook_invocations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hook_event_name TEXT NOT NULL,
    conversation_id TEXT,
    generation_id TEXT,
    model TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    cursor_version TEXT,
    workspace_roots TEXT,
    user_email TEXT,
    pwd TEXT,  -- Current working directory where hook was invoked
    raw_json TEXT NOT NULL  -- Complete raw JSON input from Cursor
);

-- Create indexes for common queries
CREATE INDEX IF NOT EXISTS idx_hook_invocations_timestamp ON hook_invocations(timestamp);
CREATE INDEX IF NOT EXISTS idx_hook_invocations_conversation ON hook_invocations(conversation_id);
CREATE INDEX IF NOT EXISTS idx_hook_invocations_event_name ON hook_invocations(hook_event_name);
CREATE INDEX IF NOT EXISTS idx_hook_invocations_model ON hook_invocations(model);
EOF

exit 0
