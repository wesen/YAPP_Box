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
    sqlite3 "$DB_PATH" <<EOF | while IFS='|' read -r hook_name count; do
SELECT 
    hook_event_name,
    COUNT(*) as count
FROM hook_invocations
GROUP BY hook_event_name
ORDER BY count DESC;
EOF
        if [ -n "$hook_name" ]; then
            printf "  %-30s %s\n" "$hook_name" "$count"
        fi
    done
    
    # Total events
    TOTAL=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM hook_invocations;")
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
    hook_event_name as Hook,
    conversation_id as Conversation,
    model as Model,
    pwd as PWD
FROM hook_invocations
ORDER BY timestamp DESC
LIMIT 10;
EOF
    echo
}

# Function to show recent agent responses
show_recent_responses() {
    echo -e "${YELLOW}💬 Recent Agent Responses (Last 5)${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    sqlite3 "$DB_PATH" <<EOF | while IFS='|' read -r timestamp conversation_id text_preview; do
SELECT 
    timestamp,
    conversation_id,
    json_extract(raw_json, '$.text') as text_preview
FROM hook_invocations
WHERE hook_event_name = 'afterAgentResponse'
ORDER BY timestamp DESC
LIMIT 5;
EOF
        if [ -n "$timestamp" ]; then
            echo -e "${BLUE}[$timestamp]${NC} ${GREEN}Conv: $conversation_id${NC}"
            # Truncate text preview to 100 chars
            text_preview_short=$(echo "$text_preview" | cut -c1-100)
            echo "  ${text_preview_short}..."
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
FROM hook_invocations
WHERE timestamp > datetime('now', '-1 hour')
GROUP BY time_bucket
ORDER BY time_bucket DESC
LIMIT 10;
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
