---
Title: 2025-11-16 Step 15 - template refactor
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
Summary: Migrated schemagen to text/template for readability
LastUpdated: 2025-11-15T20:55:26.308021464-05:00
---


# 2025-11-16 Step 15 - template refactor

<!-- Log entries in reverse chronological order (newest first) -->

## 2025-11-16 - Migrated schemagen to text/template

Refactored `pkg/schemagen/codegen.go` to use embedded `text/template` files instead of manual `fmt.Fprintf` string building:

- Created `pkg/schemagen/templates/schema_gen.go.tmpl` for struct definitions + ApplyDefaults/CustomValidate stubs
- Created `pkg/schemagen/templates/schema_gen_test.go.tmpl` for generated test cases
- Created `pkg/schemagen/templates/modules_gen.go.tmpl` for registry init code
- Embedded templates via `//go:embed` directives in `codegen.go`
- Refactored `generateModuleCode`, `generateModuleTests`, and `generateModulesRegistry` to execute templates
- Added `buildDefaultAssignment` helper for default value code generation

**Benefits:**
- Template files are easier to read and maintain than string concatenation
- Changes to output format don't require recompiling schemagen
- Template syntax makes structure clearer

**Testing:**
- `go test ./pkg/schemagen` passes
- `go run ./cmd/schemagen discover` regenerates identical output
- `go test ./...` confirms all packages build
