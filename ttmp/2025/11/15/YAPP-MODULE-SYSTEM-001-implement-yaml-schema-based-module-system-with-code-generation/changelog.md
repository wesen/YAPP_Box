# Changelog

## 2025-11-15

- Initial workspace created


## 2025-11-15

Created ticket for module system implementation. Moved implementation guide from YAPP-PUSH-BUTTONS-001. Created 24 tasks covering 6-week roadmap: registry package, schemagen tool, module conversions, resolver integration, documentation, and CI.

### Related Files

- ttmp/2025/11/15/YAPP-MODULE-SYSTEM-001-implement-yaml-schema-based-module-system-with-code-generation/playbook/module-system-implementation-guide.md


## 2025-11-15

Defined registry schema interfaces and field metadata scaffolding

### Related Files

- pkg/registry/schema.go — Field definitions and schema contracts


## 2025-11-15

Implemented registry Register/Get/All plus SetTestRegistry

### Related Files

- pkg/registry/registry.go — concurrency-safe registry map with ordering


## 2025-11-15

Added registry unit tests covering ordering, duplicates, SetTestRegistry

### Related Files

- pkg/registry/registry_test.go — new tests


## 2025-11-15

Ran go test ./... to confirm registry scaffolding integrates cleanly


## 2025-11-15

Created cmd/schemagen CLI skeleton with validate/discover stubs

### Related Files

- cmd/schemagen/main.go — Cobra root and placeholder commands


## 2025-11-15

Implemented schemagen validate with YAML parser, nested field checks, and CLI integration

### Related Files

- cmd/schemagen/main.go — wired validate subcommand
- pkg/schemagen/validate.go — validation walker
- pkg/schemagen/validate_test.go — tests for validator


## 2025-11-15

Upgraded schemagen error UX with context snippets, hints, and tests

### Related Files

- pkg/schemagen/context.go — snippet renderer
- pkg/schemagen/errors.go — enhanced ValidationError
- pkg/schemagen/validate.go — parser wrapper and helpful messages
- pkg/schemagen/validate_test.go — error UX tests


## 2025-11-15

Attempted schemagen validate on pushbuttons schema; blocked because schema.yaml not present yet

