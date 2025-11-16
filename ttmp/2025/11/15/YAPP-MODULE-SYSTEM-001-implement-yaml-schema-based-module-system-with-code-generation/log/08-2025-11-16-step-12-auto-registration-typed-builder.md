---
Title: 2025-11-16 Step 12 - auto registration + typed builder
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
Summary: Registry init + push_buttons builder uses generated structs
LastUpdated: 2025-11-15T20:38:06.991085717-05:00
---


# 2025-11-16 Step 12 - auto registration + typed builder

<!-- Log entries in reverse chronological order (newest first) -->

## 2025-11-16 - Auto registration + typed builder

- `schemagen discover` now detects the module path from `go.mod`, generates `pkg/yappgen/modules_gen.go`, and registers each module with `pkg/registry`.
- Added `pkg/yappgen/modules/pushbuttons/registry.go` exposing `NewModule()` so generated registry wiring compiles.
- Rewrote `pushbuttons.Build` to decode DSL maps into the generated `PushButtonsItem` struct (via YAML), using typed access for cap/lid/switch settings and optional flags.
- Extended the push_buttons schema with `shape_preset`, regenerated struct/tests, and reran `go test ./...`.
