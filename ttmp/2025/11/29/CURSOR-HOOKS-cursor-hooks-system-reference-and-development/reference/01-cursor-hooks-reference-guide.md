---
Title: Cursor Hooks Reference Guide
Ticket: CURSOR-HOOKS
Status: active
Topics:
    - cursor
    - automation
    - monitoring
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: .cursor/hooks.json
      Note: Hooks configuration file
    - Path: .cursor/hooks/dashboard.sh
      Note: Dashboard for viewing hook data
    - Path: .cursor/hooks/init_db.sh
      Note: Database initialization script
    - Path: .cursor/hooks/log_hook.sh
      Note: Shared logging script
ExternalSources: []
Summary: Complete reference guide for writing Cursor hooks based on real-world usage and database analysis
LastUpdated: 2025-11-29T13:48:57.822792195-05:00
---


# Cursor Hooks Reference Guide

## Goal

This reference provides a complete guide to writing Cursor hooks based on actual usage patterns and real data from production hooks. It covers setup, input/output formats, hook types, and best practices.

## Context

Cursor hooks allow you to observe, control, and extend the agent loop using custom scripts. Hooks are spawned processes that communicate over stdio using JSON in both directions. They run before or after defined stages of the agent loop and can observe, block, or modify behavior.

This guide is based on analysis of real hook invocations stored in `~/.cursor/hooks.db`, providing concrete examples and patterns that work in practice.

## Quick Reference

### Setup

**1. Create hooks configuration file:**

Place `hooks.json` at project root (`.cursor/hooks.json`) or home directory (`~/.cursor/hooks.json`):

```json
{
  "version": 1,
  "hooks": {
    "afterAgentResponse": [
      {
        "command": ".cursor/hooks/afterAgentResponse.sh"
      }
    ]
  }
}
```

**2. Create hook script:**

```bash
#!/bin/bash
# Hook script for afterAgentResponse - logs to SQLite

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
JSON_INPUT=$(cat)

# Log to database
echo "$JSON_INPUT" | "$SCRIPT_DIR/log_hook.sh"

exit 0
```

**3. Make executable:**

```bash
chmod +x .cursor/hooks/*.sh
```

**4. Restart Cursor** to load hooks configuration.

### Common Input Schema

All hooks receive a base set of fields in addition to hook-specific fields:

```json
{
  "conversation_id": "string",           // Stable ID across conversation turns
  "generation_id": "string",            // Changes with every user message
  "model": "string",                    // Model name (e.g., "composer-1")
  "hook_event_name": "string",         // Which hook is being run
  "cursor_version": "string",          // Cursor version (e.g., "2.1.39")
  "workspace_roots": ["string"],        // Array of workspace root paths
  "user_email": "string | null"         // Authenticated user email
}
```

### Hook Types and Their Inputs/Outputs

#### 1. `afterAgentResponse`

**When:** After agent completes an assistant message.

**Input:**
```json
{
  "conversation_id": "1649388c-47c9-4687-850b-deff0700d2bc",
  "generation_id": "4ca49744-abe1-446c-9682-c68e3642de0e",
  "model": "composer-1",
  "text": "Full agent response text here...",
  "hook_event_name": "afterAgentResponse",
  "cursor_version": "2.1.39",
  "workspace_roots": ["/path/to/workspace"],
  "user_email": "user@example.com"
}
```

**Output:** None required (logging hook).

**Use Cases:** Logging agent responses, analytics, content analysis.

---

#### 2. `afterAgentThought`

**When:** After agent completes a thinking block.

**Input:**
```json
{
  "conversation_id": "1649388c-47c9-4687-850b-deff0700d2bc",
  "generation_id": "4ca49744-abe1-446c-9682-c68e3642de0e",
  "model": "composer-1",
  "text": "Thinking text content...",
  "duration_ms": 406,
  "hook_event_name": "afterAgentThought",
  "cursor_version": "2.1.39",
  "workspace_roots": ["/path/to/workspace"],
  "user_email": "user@example.com"
}
```

**Output:** None required.

**Use Cases:** Observing agent reasoning, performance monitoring.

---

#### 3. `beforeShellExecution` / `afterShellExecution`

**When:** Before/after shell commands are executed.

**Input (beforeShellExecution):**
```json
{
  "conversation_id": "1649388c-47c9-4687-850b-deff0700d2bc",
  "generation_id": "4ca49744-abe1-446c-9682-c68e3642de0e",
  "model": "composer-1",
  "command": "docmgr changelog update --ticket YAPP-HOOKS-001",
  "cwd": "",
  "hook_event_name": "beforeShellExecution",
  "cursor_version": "2.1.39",
  "workspace_roots": ["/path/to/workspace"],
  "user_email": "user@example.com"
}
```

**Input (afterShellExecution):**
```json
{
  "conversation_id": "1649388c-47c9-4687-850b-deff0700d2bc",
  "generation_id": "4ca49744-abe1-446c-9682-c68e3642de0e",
  "model": "composer-1",
  "command": "docmgr changelog update --ticket YAPP-HOOKS-001",
  "output": "Changelog updated: /path/to/changelog.md\n...",
  "duration": 1234,
  "hook_event_name": "afterShellExecution",
  "cursor_version": "2.1.39",
  "workspace_roots": ["/path/to/workspace"],
  "user_email": "user@example.com"
}
```

**Output (beforeShellExecution):**
```json
{
  "permission": "allow" | "deny" | "ask",
  "user_message": "Optional message shown to user",
  "agent_message": "Optional message sent to agent"
}
```

**Output (afterShellExecution):** None required.

**Use Cases:** Command auditing, blocking dangerous commands, redirecting to safer alternatives.

---

#### 4. `beforeMCPExecution` / `afterMCPExecution`

**When:** Before/after MCP (Model Context Protocol) tool usage.

**Input (beforeMCPExecution):**
```json
{
  "conversation_id": "string",
  "generation_id": "string",
  "model": "string",
  "tool_name": "string",
  "tool_input": "json params string",
  "url": "string" | null,              // If MCP server uses URL
  "command": "string" | null,           // If MCP server uses command
  "hook_event_name": "beforeMCPExecution",
  "cursor_version": "string",
  "workspace_roots": ["string"],
  "user_email": "string | null"
}
```

**Input (afterMCPExecution):**
```json
{
  "conversation_id": "string",
  "generation_id": "string",
  "model": "string",
  "tool_name": "string",
  "tool_input": "json params string",
  "result_json": "tool result json string",
  "duration": 1234,
  "url": "string" | null,
  "command": "string" | null,
  "hook_event_name": "afterMCPExecution",
  "cursor_version": "string",
  "workspace_roots": ["string"],
  "user_email": "string | null"
}
```

**Output (beforeMCPExecution):**
```json
{
  "permission": "allow" | "deny" | "ask",
  "user_message": "Optional message",
  "agent_message": "Optional message"
}
```

**Output (afterMCPExecution):** None required.

**Use Cases:** MCP tool auditing, access control, result validation.

---

#### 5. `beforeReadFile` / `afterFileEdit`

**When:** Before file reads / after file edits by Agent.

**Input (beforeReadFile):**
```json
{
  "conversation_id": "string",
  "generation_id": "string",
  "model": "string",
  "file_path": "/absolute/path/to/file",
  "content": "file contents",
  "hook_event_name": "beforeReadFile",
  "cursor_version": "string",
  "workspace_roots": ["string"],
  "user_email": "string | null"
}
```

**Input (afterFileEdit):**
```json
{
  "conversation_id": "1649388c-47c9-4687-850b-deff0700d2bc",
  "generation_id": "4ca49744-abe1-446c-9682-c68e3642de0e",
  "model": "composer-1",
  "file_path": "/path/to/file.sh",
  "edits": [
    {
      "old_string": "original content",
      "new_string": "new content"
    }
  ],
  "hook_event_name": "afterFileEdit",
  "cursor_version": "2.1.39",
  "workspace_roots": ["/path/to/workspace"],
  "user_email": "user@example.com"
}
```

**Output (beforeReadFile):**
```json
{
  "permission": "allow" | "deny"
}
```

**Output (afterFileEdit):** None required.

**Use Cases:** File access control, automatic formatting, secret redaction, change tracking.

---

#### 6. `beforeSubmitPrompt`

**When:** Right after user hits send but before backend request.

**Input:**
```json
{
  "conversation_id": "8047dda4-e4ae-4c61-9d96-8a5a71fd6a8e",
  "generation_id": "13d8cce8-a704-46bb-ba23-005d05fe9115",
  "model": "composer-1",
  "prompt": "Continue",
  "attachments": [
    {
      "type": "file",
      "file_path": "/absolute/path/to/file.md"
    },
    {
      "type": "rule",
      "file_path": "/absolute/path/to/rule.md"
    }
  ],
  "hook_event_name": "beforeSubmitPrompt",
  "cursor_version": "2.1.39",
  "workspace_roots": ["/path/to/workspace"],
  "user_email": "user@example.com"
}
```

**Output:**
```json
{
  "continue": true | false,
  "user_message": "Message shown to user when blocked"
}
```

**Use Cases:** Prompt validation, PII scanning, prompt blocking, prompt modification.

---

#### 7. `stop`

**When:** When agent loop ends.

**Input:**
```json
{
  "conversation_id": "1649388c-47c9-4687-850b-deff0700d2bc",
  "generation_id": "4ca49744-abe1-446c-9682-c68e3642de0e",
  "model": "composer-1",
  "status": "completed" | "aborted" | "error",
  "loop_count": 0,
  "hook_event_name": "stop",
  "cursor_version": "2.1.39",
  "workspace_roots": ["/path/to/workspace"],
  "user_email": "user@example.com"
}
```

**Output (optional):**
```json
{
  "followup_message": "Optional message to auto-submit as next user message"
}
```

**Use Cases:** Cleanup, final logging, triggering follow-up workflows.

---

#### 8. `beforeTabFileRead` / `afterTabFileEdit`

**When:** Before Tab (inline completions) reads files / after Tab edits files.

**Input (beforeTabFileRead):**
```json
{
  "conversation_id": "string",
  "generation_id": "string",
  "model": "string",
  "file_path": "/absolute/path/to/file",
  "content": "file contents",
  "hook_event_name": "beforeTabFileRead",
  "cursor_version": "string",
  "workspace_roots": ["string"],
  "user_email": "string | null"
}
```

**Input (afterTabFileEdit):**
```json
{
  "conversation_id": "string",
  "generation_id": "string",
  "model": "string",
  "file_path": "/absolute/path/to/file",
  "edits": [
    {
      "old_string": "original",
      "new_string": "new",
      "range": {
        "start_line_number": 10,
        "start_column": 5,
        "end_line_number": 10,
        "end_column": 20
      },
      "old_line": "original line content",
      "new_line": "new line content"
    }
  ],
  "hook_event_name": "afterTabFileEdit",
  "cursor_version": "string",
  "workspace_roots": ["string"],
  "user_email": "string | null"
}
```

**Output (beforeTabFileRead):**
```json
{
  "permission": "allow" | "deny"
}
```

**Output (afterTabFileEdit):** None required.

**Use Cases:** Tab-specific access control, formatting Tab edits, auditing Tab usage.

## Usage Examples

### Example 1: Simple Logging Hook

```bash
#!/bin/bash
# afterAgentResponse.sh - Log agent responses

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
JSON_INPUT=$(cat)

# Log to database
echo "$JSON_INPUT" | "$SCRIPT_DIR/log_hook.sh"

exit 0
```

### Example 2: Permission-Based Hook

```bash
#!/bin/bash
# beforeShellExecution.sh - Block dangerous commands

JSON_INPUT=$(cat)
COMMAND=$(echo "$JSON_INPUT" | jq -r '.command')

# Block git commands
if [[ "$COMMAND" =~ ^git ]]; then
    cat <<EOF
{
  "permission": "deny",
  "user_message": "Git commands are blocked. Use GitHub CLI (gh) instead.",
  "agent_message": "Git command blocked. Please use 'gh' tool instead."
}
EOF
    exit 0
fi

# Allow other commands
cat <<EOF
{
  "permission": "allow"
}
EOF
exit 0
```

### Example 3: Auto-Formatting Hook

```bash
#!/bin/bash
# afterFileEdit.sh - Auto-format edited files

JSON_INPUT=$(cat)
FILE_PATH=$(echo "$JSON_INPUT" | jq -r '.file_path')

# Format Python files
if [[ "$FILE_PATH" == *.py ]]; then
    black "$FILE_PATH" 2>/dev/null
fi

# Format Go files
if [[ "$FILE_PATH" == *.go ]]; then
    gofmt -w "$FILE_PATH" 2>/dev/null
fi

exit 0
```

### Example 4: Prompt Validation Hook

```bash
#!/bin/bash
# beforeSubmitPrompt.sh - Validate prompts

JSON_INPUT=$(cat)
PROMPT=$(echo "$JSON_INPUT" | jq -r '.prompt')

# Block prompts containing secrets
if echo "$PROMPT" | grep -qiE "(api[_-]?key|password|secret|token)"; then
    cat <<EOF
{
  "continue": false,
  "user_message": "Prompt contains potentially sensitive information. Please remove secrets before submitting."
}
EOF
    exit 0
fi

# Allow submission
cat <<EOF
{
  "continue": true
}
EOF
exit 0
```

## Best Practices

### 1. Path Configuration

- **Project-level hooks:** Use relative paths from project root: `.cursor/hooks/script.sh`
- **Global hooks:** Use absolute paths or paths relative to `~/.cursor/`
- **Priority:** Project hooks override global hooks

### 2. Error Handling

- Always exit with `exit 0` unless you want to signal failure
- Log errors to stderr: `echo "Error message" >&2`
- Handle missing tools gracefully (check if command exists before using)

### 3. Performance

- Keep hooks fast (< 100ms ideally)
- Avoid blocking operations
- Use background processes for heavy work if needed

### 4. Security

- Validate all inputs
- Escape shell commands properly
- Don't trust file paths - validate them
- Be careful with `eval` or command substitution

### 5. Logging

- Store raw JSON for future analysis
- Include `pwd` (current working directory) in logs
- Use structured logging (SQLite, JSON files, etc.)

### 6. Testing

- Test hooks with sample JSON inputs
- Use `--once` mode for dashboard to verify logging
- Test permission responses (allow/deny/ask)

## Database Schema

Our implementation uses a single `hook_invocations` table:

```sql
CREATE TABLE hook_invocations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hook_event_name TEXT NOT NULL,
    conversation_id TEXT,
    generation_id TEXT,
    model TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    cursor_version TEXT,
    workspace_roots TEXT,
    user_email TEXT,
    pwd TEXT,                    -- Current working directory
    raw_json TEXT NOT NULL       -- Complete raw JSON input
);
```

**Indexes:**
- `idx_hook_invocations_timestamp` - For time-based queries
- `idx_hook_invocations_conversation` - For conversation tracking
- `idx_hook_invocations_event_name` - For filtering by hook type
- `idx_hook_invocations_model` - For model-specific analysis

## Querying Hook Data

### View recent events
```bash
sqlite3 ~/.cursor/hooks.db "SELECT timestamp, hook_event_name, conversation_id FROM hook_invocations ORDER BY timestamp DESC LIMIT 10;"
```

### Count by hook type
```bash
sqlite3 ~/.cursor/hooks.db "SELECT hook_event_name, COUNT(*) FROM hook_invocations GROUP BY hook_event_name;"
```

### Extract specific field from JSON
```bash
sqlite3 ~/.cursor/hooks.db "SELECT json_extract(raw_json, '$.command') FROM hook_invocations WHERE hook_event_name = 'beforeShellExecution';"
```

### Use dashboard script
```bash
.cursor/hooks/dashboard.sh        # Live monitoring
.cursor/hooks/dashboard.sh --once # Single run
```

## Related

- Cursor Hooks Documentation: https://cursor.com/docs/agent/hooks
- Project hooks implementation: `.cursor/hooks/`
- Database location: `~/.cursor/hooks.db` (or `$CURSOR_HOOKS_DB`)
