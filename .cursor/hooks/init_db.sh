#!/bin/bash
# Shared script to initialize SQLite database for Cursor hooks logging
# This script ensures the database and tables exist before any hook script runs

DB_PATH="${CURSOR_HOOKS_DB:-$HOME/.cursor/hooks.db}"

# Create database directory if it doesn't exist
mkdir -p "$(dirname "$DB_PATH")"

# Initialize database with all required tables
sqlite3 "$DB_PATH" <<EOF
-- Table for beforeShellExecution and afterShellExecution hooks
CREATE TABLE IF NOT EXISTS shell_executions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hook_event_name TEXT NOT NULL,
    conversation_id TEXT,
    generation_id TEXT,
    model TEXT,
    command TEXT,
    cwd TEXT,
    output TEXT,
    duration_ms INTEGER,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    cursor_version TEXT,
    workspace_roots TEXT,
    user_email TEXT
);

-- Table for beforeMCPExecution and afterMCPExecution hooks
CREATE TABLE IF NOT EXISTS mcp_executions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hook_event_name TEXT NOT NULL,
    conversation_id TEXT,
    generation_id TEXT,
    model TEXT,
    tool_name TEXT,
    tool_input TEXT,
    result_json TEXT,
    duration_ms INTEGER,
    url TEXT,
    command TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    cursor_version TEXT,
    workspace_roots TEXT,
    user_email TEXT
);

-- Table for beforeReadFile and afterFileEdit hooks
CREATE TABLE IF NOT EXISTS file_operations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hook_event_name TEXT NOT NULL,
    conversation_id TEXT,
    generation_id TEXT,
    model TEXT,
    file_path TEXT,
    content TEXT,
    edits TEXT,  -- JSON array of edits
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    cursor_version TEXT,
    workspace_roots TEXT,
    user_email TEXT
);

-- Table for beforeSubmitPrompt hook
CREATE TABLE IF NOT EXISTS prompt_submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hook_event_name TEXT NOT NULL,
    conversation_id TEXT,
    generation_id TEXT,
    model TEXT,
    prompt TEXT,
    attachments TEXT,  -- JSON array of attachments
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    cursor_version TEXT,
    workspace_roots TEXT,
    user_email TEXT
);

-- Table for stop hook
CREATE TABLE IF NOT EXISTS agent_stops (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hook_event_name TEXT NOT NULL,
    conversation_id TEXT,
    generation_id TEXT,
    model TEXT,
    status TEXT,
    loop_count INTEGER,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    cursor_version TEXT,
    workspace_roots TEXT,
    user_email TEXT
);

-- Table for afterAgentResponse hook
CREATE TABLE IF NOT EXISTS agent_responses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hook_event_name TEXT NOT NULL,
    conversation_id TEXT,
    generation_id TEXT,
    model TEXT,
    text TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    cursor_version TEXT,
    workspace_roots TEXT,
    user_email TEXT
);

-- Table for afterAgentThought hook
CREATE TABLE IF NOT EXISTS agent_thoughts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hook_event_name TEXT NOT NULL,
    conversation_id TEXT,
    generation_id TEXT,
    model TEXT,
    text TEXT,
    duration_ms INTEGER,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    cursor_version TEXT,
    workspace_roots TEXT,
    user_email TEXT
);

-- Table for beforeTabFileRead hook
CREATE TABLE IF NOT EXISTS tab_file_reads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hook_event_name TEXT NOT NULL,
    conversation_id TEXT,
    generation_id TEXT,
    model TEXT,
    file_path TEXT,
    content TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    cursor_version TEXT,
    workspace_roots TEXT,
    user_email TEXT
);

-- Table for afterTabFileEdit hook
CREATE TABLE IF NOT EXISTS tab_file_edits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hook_event_name TEXT NOT NULL,
    conversation_id TEXT,
    generation_id TEXT,
    model TEXT,
    file_path TEXT,
    edits TEXT,  -- JSON array of edits with range, old_line, new_line
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    cursor_version TEXT,
    workspace_roots TEXT,
    user_email TEXT
);

-- Create indexes for common queries
CREATE INDEX IF NOT EXISTS idx_shell_executions_timestamp ON shell_executions(timestamp);
CREATE INDEX IF NOT EXISTS idx_shell_executions_conversation ON shell_executions(conversation_id);
CREATE INDEX IF NOT EXISTS idx_mcp_executions_timestamp ON mcp_executions(timestamp);
CREATE INDEX IF NOT EXISTS idx_mcp_executions_conversation ON mcp_executions(conversation_id);
CREATE INDEX IF NOT EXISTS idx_file_operations_timestamp ON file_operations(timestamp);
CREATE INDEX IF NOT EXISTS idx_file_operations_conversation ON file_operations(conversation_id);
CREATE INDEX IF NOT EXISTS idx_file_operations_path ON file_operations(file_path);
CREATE INDEX IF NOT EXISTS idx_prompt_submissions_timestamp ON prompt_submissions(timestamp);
CREATE INDEX IF NOT EXISTS idx_prompt_submissions_conversation ON prompt_submissions(conversation_id);
CREATE INDEX IF NOT EXISTS idx_agent_stops_timestamp ON agent_stops(timestamp);
CREATE INDEX IF NOT EXISTS idx_agent_stops_conversation ON agent_stops(conversation_id);
CREATE INDEX IF NOT EXISTS idx_agent_responses_timestamp ON agent_responses(timestamp);
CREATE INDEX IF NOT EXISTS idx_agent_responses_conversation ON agent_responses(conversation_id);
CREATE INDEX IF NOT EXISTS idx_agent_thoughts_timestamp ON agent_thoughts(timestamp);
CREATE INDEX IF NOT EXISTS idx_agent_thoughts_conversation ON agent_thoughts(conversation_id);
CREATE INDEX IF NOT EXISTS idx_tab_file_reads_timestamp ON tab_file_reads(timestamp);
CREATE INDEX IF NOT EXISTS idx_tab_file_reads_conversation ON tab_file_reads(conversation_id);
CREATE INDEX IF NOT EXISTS idx_tab_file_edits_timestamp ON tab_file_edits(timestamp);
CREATE INDEX IF NOT EXISTS idx_tab_file_edits_conversation ON tab_file_edits(conversation_id);
EOF

exit 0

