---
Title: 2025-11-16 Step 7 - schemagen error UX
Ticket: YAPP-MODULE-SYSTEM-001
Status: active
Topics:
    - yapp
    - dsl
    - codegen
    - architecture
DocType: log
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Improved validator snippets and hints
LastUpdated: 2025-11-15T19:50:12.676228532-05:00
---


# 2025-11-16 Step 7 - schemagen error UX

<!-- Log entries in reverse chronological order (newest first) -->

## 2025-11-16 - Added contextual schemagen errors

- Introduced `pkg/schemagen/context.go` to render code snippets with caret indicators for each validation error
- Enhanced `ValidationError` to include snippets and hints; CLI now surfaces multi-line messages
- Added YAML parser wrapper that recognizes colon-related syntax errors and provides actionable hints, plus regression tests for snippets/hints
