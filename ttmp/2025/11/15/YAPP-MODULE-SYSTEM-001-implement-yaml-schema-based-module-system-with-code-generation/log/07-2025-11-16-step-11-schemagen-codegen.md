---
Title: 2025-11-16 Step 11 - schemagen codegen
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
Summary: Added schema parser, Go struct generation, modules registry
LastUpdated: 2025-11-15T20:29:44.159546606-05:00
---


# 2025-11-16 Step 11 - schemagen codegen

<!-- Log entries in reverse chronological order (newest first) -->

## 2025-11-16 - Generated Go structs/tests/registry

- Added `pkg/schemagen/schema_doc.go` to parse schemas into ordered field metadata and expose test cases.
- Implemented `pkg/schemagen/codegen.go` to emit struct definitions, positive YAML tests, and a `pkg/yappgen/modules_gen.go` metadata table; covered by `codegen_test.go`.
- Updated `schemagen discover` to load docs, call the generator, and drop the “code generation coming soon” error; `go run ./cmd/schemagen discover` now produces `schema_gen.go`, `_test.go`, and `modules_gen.go`.
