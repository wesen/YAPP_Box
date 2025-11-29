# Changelog

## 2025-11-29

- Initial workspace created


## 2025-11-29

Created hooks system with SQLite logging. Enabled only afterAgentResponse hook for initial testing. All hook scripts created and ready for future enablement.


## 2025-11-29

Added dashboard.sh script for live monitoring of hooks database with statistics, recent events, and activity timeline


## 2025-11-29

Enabled additional hooks: afterAgentThought, afterFileEdit, and afterShellExecution for better visibility into agent activity


## 2025-11-29

Enabled all 12 hook types with project-relative paths (.cursor/hooks/). All hooks now active for comprehensive logging.


## 2025-11-29

Refactored database schema to use single generic hook_invocations table. Stores raw JSON, pwd, and common indexed fields. Simplified logging script and updated dashboard.

