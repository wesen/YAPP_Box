---
Title: 2025-11-16 Step 9 - real schema validation
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
Summary: First schemagen run against repo modules
LastUpdated: 2025-11-15T19:57:55.818682396-05:00
---


# 2025-11-16 Step 9 - real schema validation

<!-- Log entries in reverse chronological order (newest first) -->

## 2025-11-16 - First CLI validation attempt

- Ran `go run ./cmd/schemagen validate pkg/yappgen/modules/pushbuttons/schema.yaml`
- Command failed because `schema.yaml` does not exist yet for pushbuttons (still using legacy builder)
- Documented the missing schema so we know the next actionable step is writing real schema files before running discover end-to-end
