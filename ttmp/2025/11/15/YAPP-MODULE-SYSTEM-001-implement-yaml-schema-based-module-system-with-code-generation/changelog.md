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


## 2025-11-15

Added push_buttons schema.yaml and validator bool support

### Related Files

- pkg/schemagen/validate.go — bool field support
- pkg/yappgen/modules/pushbuttons/schema.yaml — push button schema


## 2025-11-15

schemagen discover now generates Go structs/tests and modules registry

### Related Files

- cmd/schemagen/main.go — discover wiring
- pkg/schemagen/codegen.go — struct/test generation
- pkg/schemagen/schema_doc.go — schema loader
- pkg/yappgen/modules/pushbuttons/schema_gen.go — generated output


## 2025-11-15

Auto-registered schemagen modules and rewrote push_buttons builder to use generated structs

### Related Files

- pkg/schemagen/codegen.go — registry imports + module path
- pkg/schemagen/module_path.go — go.mod parser
- pkg/yappgen/modules/pushbuttons/module.go — typed builder
- pkg/yappgen/modules_gen.go — generated registry init


## 2025-11-15

schemagen now emits ModuleSchema implementations and resolver validates via registry

### Related Files

- pkg/resolver/resolver.go — calls registry schema validation
- pkg/schemagen/codegen.go — generated ModuleSchema code
- pkg/schemagen/schema_doc.go — field metadata for ModuleSchema


## 2025-11-15

Refactored schemagen to use text/template for cleaner codegen

### Related Files

- pkg/schemagen/codegen.go — template execution logic
- pkg/schemagen/templates/modules_gen.go.tmpl — registry template
- pkg/schemagen/templates/schema_gen.go.tmpl — struct template
- pkg/schemagen/templates/schema_gen_test.go.tmpl — test template


## 2025-11-15

Converted pcb_stands to schema-based builder

### Related Files

- pkg/yappgen/modules/pcbstands/module.go — typed builder
- pkg/yappgen/modules/pcbstands/registry.go — module registration
- pkg/yappgen/modules/pcbstands/schema.yaml — pcb_stands schema


## 2025-11-15

Wired resolver two-phase validation and confirmed end-to-end YAML to SCAD works

### Related Files

- examples/test-push-buttons.yaml — test case with expressions
- pkg/resolver/strict.go — strict mode helpers
- pkg/resolver/validation.go — schema validation hooks

