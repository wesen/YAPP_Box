#!/bin/bash
# Dashboard script to monitor Cursor hooks database in real-time

DB_PATH="${CURSOR_HOOKS_DB:-$HOME/.cursor/hooks.db}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Check if database exists
if [ ! -f "$DB_PATH" ]; then
    echo -e "${RED}Error: Database not found at $DB_PATH${NC}"
    echo "Waiting for first hook event to create database..."
    exit 1
fi

# Function to clear screen and show header
show_header() {
    clear
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${CYAN}  Cursor Hooks Dashboard${NC}"
    echo -e "${CYAN}  Database: $DB_PATH${NC}"
    echo -e "${CYAN}  Last Update: $(date '+%Y-%m-%d %H:%M:%S')${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    echo
}

# Function to show statistics
show_stats() {
    echo -e "${YELLOW}📊 Statistics${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    # Count by hook type
    echo -e "${GREEN}Hook Events by Type:${NC}"
    sqlite3 "$DB_PATH" <<EOF | while IFS='|' read -r table count; do
SELECT 
    CASE 
        WHEN table_name = 'shell_executions' THEN 'Shell Executions'
        WHEN table_name = 'mcp_executions' THEN 'MCP Executions'
        WHEN table_name = 'file_operations' THEN 'File Operations'
        WHEN table_name = 'prompt_submissions' THEN 'Prompt Submissions'
        WHEN table_name = 'agent_stops' THEN 'Agent Stops'
        WHEN table_name = 'agent_responses' THEN 'Agent Responses'
        WHEN table_name = 'agent_thoughts' THEN 'Agent Thoughts'
        WHEN table_name = 'tab_file_reads' THEN 'Tab File Reads'
        WHEN table_name = 'tab_file_edits' THEN 'Tab File Edits'
        ELSE table_name
    END as display_name,
    COUNT(*) as count
FROM (
    SELECT 'shell_executions' as table_name FROM shell_executions
    UNION ALL SELECT 'mcp_executions' FROM mcp_executions
    UNION ALL SELECT 'file_operations' FROM file_operations
    UNION ALL SELECT 'prompt_submissions' FROM prompt_submissions
    UNION ALL SELECT 'agent_stops' FROM agent_stops
    UNION ALL SELECT 'agent_responses' FROM agent_responses
    UNION ALL SELECT 'agent_thoughts' FROM agent_thoughts
    UNION ALL SELECT 'tab_file_reads' FROM tab_file_reads
    UNION ALL SELECT 'tab_file_edits' FROM tab_file_edits
) GROUP BY table_name ORDER BY count DESC;
EOF
        if [ -n "$table" ]; then
            printf "  %-25s %s\n" "$table" "$count"
        fi
    done
    
    # Total events
    TOTAL=$(sqlite3 "$DB_PATH" <<EOF
SELECT 
    (SELECT COUNT(*) FROM shell_executions) +
    (SELECT COUNT(*) FROM mcp_executions) +
    (SELECT COUNT(*) FROM file_operations) +
    (SELECT COUNT(*) FROM prompt_submissions) +
    (SELECT COUNT(*) FROM agent_stops) +
    (SELECT COUNT(*) FROM agent_responses) +
    (SELECT COUNT(*) FROM agent_thoughts) +
    (SELECT COUNT(*) FROM tab_file_reads) +
    (SELECT COUNT(*) FROM tab_file_edits);
EOF
)
    echo -e "\n${GREEN}Total Events:${NC} $TOTAL"
    echo
}

# Function to show recent events
show_recent_events() {
    echo -e "${YELLOW}🕐 Recent Events (Last 10)${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    sqlite3 -header -column "$DB_PATH" <<EOF
SELECT 
    timestamp as Time,
    CASE 
        WHEN source = 'shell_executions' THEN 'Shell'
        WHEN source = 'mcp_executions' THEN 'MCP'
        WHEN source = 'file_operations' THEN 'File'
        WHEN source = 'prompt_submissions' THEN 'Prompt'
        WHEN source = 'agent_stops' THEN 'Stop'
        WHEN source = 'agent_responses' THEN 'Response'
        WHEN source = 'agent_thoughts' THEN 'Thought'
        WHEN source = 'tab_file_reads' THEN 'Tab Read'
        WHEN source = 'tab_file_edits' THEN 'Tab Edit'
        ELSE source
    END as Type,
    conversation_id as Conversation,
    model as Model
FROM (
    SELECT timestamp, 'shell_executions' as source, conversation_id, model FROM shell_executions
    UNION ALL SELECT timestamp, 'mcp_executions', conversation_id, model FROM mcp_executions
    UNION ALL SELECT timestamp, 'file_operations', conversation_id, model FROM file_operations
    UNION ALL SELECT timestamp, 'prompt_submissions', conversation_id, model FROM prompt_submissions
    UNION ALL SELECT timestamp, 'agent_stops', conversation_id, model FROM agent_stops
    UNION ALL SELECT timestamp, 'agent_responses', conversation_id, model FROM agent_responses
    UNION ALL SELECT timestamp, 'agent_thoughts', conversation_id, model FROM agent_thoughts
    UNION ALL SELECT timestamp, 'tab_file_reads', conversation_id, model FROM tab_file_reads
    UNION ALL SELECT timestamp, 'tab_file_edits', conversation_id, model FROM tab_file_edits
) ORDER BY timestamp DESC LIMIT 10;
EOF
    echo
}

# Function to show recent agent responses (most relevant for current setup)
show_recent_responses() {
    echo -e "${YELLOW}💬 Recent Agent Responses (Last 5)${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    sqlite3 "$DB_PATH" <<EOF | while IFS='|' read -r timestamp conversation_id text_preview; do
SELECT 
    timestamp,
    conversation_id,
    substr(text, 1, 100) || '...' as text_preview
FROM agent_responses
ORDER BY timestamp DESC
LIMIT 5;
EOF
        if [ -n "$timestamp" ]; then
            echo -e "${BLUE}[$timestamp]${NC} ${GREEN}Conv: $conversation_id${NC}"
            echo "  ${text_preview}"
            echo
        fi
    done
}

# Function to show activity timeline
show_timeline() {
    echo -e "${YELLOW}📈 Activity Timeline (Last Hour)${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    sqlite3 "$DB_PATH" <<EOF | while IFS='|' read -r time_bucket count; do
SELECT 
    strftime('%H:%M', timestamp) as time_bucket,
    COUNT(*) as count
FROM (
    SELECT timestamp FROM shell_executions WHERE timestamp > datetime('now', '-1 hour')
    UNION ALL SELECT timestamp FROM mcp_executions WHERE timestamp > datetime('now', '-1 hour')
    UNION ALL SELECT timestamp FROM file_operations WHERE timestamp > datetime('now', '-1 hour')
    UNION ALL SELECT timestamp FROM prompt_submissions WHERE timestamp > datetime('now', '-1 hour')
    UNION ALL SELECT timestamp FROM agent_stops WHERE timestamp > datetime('now', '-1 hour')
    UNION ALL SELECT timestamp FROM agent_responses WHERE timestamp > datetime('now', '-1 hour')
    UNION ALL SELECT timestamp FROM agent_thoughts WHERE timestamp > datetime('now', '-1 hour')
    UNION ALL SELECT timestamp FROM tab_file_reads WHERE timestamp > datetime('now', '-1 hour')
    UNION ALL SELECT timestamp FROM tab_file_edits WHERE timestamp > datetime('now', '-1 hour')
) GROUP BY time_bucket ORDER BY time_bucket DESC LIMIT 10;
EOF
        if [ -n "$time_bucket" ]; then
            printf "  %s: %s events\n" "$time_bucket" "$count"
        fi
    done
    echo
}

# Main loop
if [ "$1" = "--once" ]; then
    # Run once and exit
    show_header
    show_stats
    show_recent_events
    show_recent_responses
    show_timeline
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
else
    # Continuous monitoring
    while true; do
        show_header
        show_stats
        show_recent_events
        show_recent_responses
        show_timeline
        echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
        echo -e "${MAGENTA}Press Ctrl+C to exit${NC}"
        sleep 2
    done
fi

