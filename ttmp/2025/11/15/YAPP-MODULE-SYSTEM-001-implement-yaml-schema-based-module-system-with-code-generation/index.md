---
Title: Implement YAML Schema-Based Module System with Code Generation
Ticket: YAPP-MODULE-SYSTEM-001
Status: active
Topics:
    - yapp
    - dsl
    - codegen
    - architecture
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: cmd/schemagen/main.go
      Note: discover wiring
    - Path: examples/test-push-buttons.yaml
      Note: end-to-end test case
    - Path: pkg/docs/docs.go
      Note: loads module help into CLI
    - Path: pkg/docs/schema_help.go
      Note: auto-generates help from schemas
    - Path: pkg/docs/tutorials/yapp-module-authoring-guide.md
      Note: step-by-step guide for adding modules
    - Path: pkg/registry/registry.go
      Note: Register/Get/All and SetTestRegistry implementations
    - Path: pkg/registry/registry_test.go
      Note: verified via go test ./...
    - Path: pkg/registry/schema.go
      Note: Defines FieldSpec
    - Path: pkg/resolver/resolver.go
      Note: integrates registry schema validation
    - Path: pkg/resolver/strict.go
      Note: strict mode validation helpers
    - Path: pkg/resolver/validation.go
      Note: two-phase schema validation
    - Path: pkg/schemagen/codegen.go
      Note: generates ModuleSchema + registry init
    - Path: pkg/schemagen/codegen_test.go
      Note: codegen tests
    - Path: pkg/schemagen/context.go
      Note: snippet rendering + source context
    - Path: pkg/schemagen/discover.go
      Note: schema discovery and validation
    - Path: pkg/schemagen/discover_test.go
      Note: discovery tests
    - Path: pkg/schemagen/errors.go
      Note: ValidationError snippet/hint support
    - Path: pkg/schemagen/module_path.go
      Note: detects module path for codegen
    - Path: pkg/schemagen/schema_doc.go
      Note: ModuleSchema metadata
    - Path: pkg/schemagen/templates/modules_gen.go.tmpl
      Note: registry template
    - Path: pkg/schemagen/templates/schema_gen.go.tmpl
      Note: struct template
    - Path: pkg/schemagen/templates/schema_gen_test.go.tmpl
      Note: test template
    - Path: pkg/schemagen/validate.go
      Note: enhanced error reporting
    - Path: pkg/schemagen/validate_test.go
      Note: tests for snippet/hint + syntax errors
    - Path: pkg/yappgen/modules/connectors/module.go
      Note: connectors typed builder
    - Path: pkg/yappgen/modules/connectors/registry.go
      Note: connectors module registration
    - Path: pkg/yappgen/modules/connectors/schema.yaml
      Note: connectors module schema
    - Path: pkg/yappgen/modules/connectors/schema_gen.go
      Note: generated connectors structs
    - Path: pkg/yappgen/modules/cutouts/module.go
      Note: cutouts typed builder
    - Path: pkg/yappgen/modules/cutouts/registry.go
      Note: cutouts module registration
    - Path: pkg/yappgen/modules/cutouts/schema.yaml
      Note: cutouts module schema
    - Path: pkg/yappgen/modules/cutouts/schema_gen.go
      Note: generated cutouts structs
    - Path: pkg/yappgen/modules/pcbstands/module.go
      Note: typed builder implementation
    - Path: pkg/yappgen/modules/pcbstands/registry.go
      Note: pcb_stands module registration
    - Path: pkg/yappgen/modules/pcbstands/schema.yaml
      Note: pcb_stands schema definition
    - Path: pkg/yappgen/modules/pcbstands/schema_gen.go
      Note: generated pcb_stands structs
    - Path: pkg/yappgen/modules/pushbuttons/module.go
      Note: builder now uses generated structs
    - Path: pkg/yappgen/modules/pushbuttons/registry.go
      Note: push_buttons module registration
    - Path: pkg/yappgen/modules/pushbuttons/schema.yaml
      Note: added shape_preset
    - Path: pkg/yappgen/modules/pushbuttons/schema_gen.go
      Note: generated push_buttons structs
    - Path: pkg/yappgen/modules/snapjoins/module.go
      Note: snap_joins typed builder
    - Path: pkg/yappgen/modules/snapjoins/registry.go
      Note: snap_joins module registration
    - Path: pkg/yappgen/modules/snapjoins/schema.yaml
      Note: snap_joins module schema
    - Path: pkg/yappgen/modules/snapjoins/schema_gen.go
      Note: generated snap_joins structs
    - Path: pkg/yappgen/modules_gen.go
      Note: auto-generated registry init
ExternalSources:
    - https://github.com/deepmap/oapi-codegen
    - https://gqlgen.com/
Summary: 'Build a schema-driven module system for YAPP DSL: YAML schemas define validation and structure, code generation creates typed Go structs and tests, automatic discovery registers modules without core file edits, and documentation auto-generates from schemas.'
LastUpdated: 2025-11-15T00:00:00Z
---





















# Implement YAML Schema-Based Module System with Code Generation

## Overview

This ticket implements a comprehensive module system for the YAPP DSL that eliminates validation duplication, enables zero-touch module addition, and auto-generates documentation from schemas. The work was designed through three debate rounds exploring validation flow, registration mechanics, and prototype implementation, culminating in a detailed implementation guide for newcomers.

**Core problem:** Currently, adding a DSL feature requires editing 4+ files with duplicated validation logic scattered across resolver, model builder, and feature helpers. This creates maintenance burden, merge conflicts, and inconsistent error messages.

**Solution:** YAML schemas as single source of truth + code generation + discovery-based registration. Module authors write one schema.yaml file defining fields, validation rules, and test cases. The `schemagen` tool generates typed Go structs, tests, and registration code. No core file edits required.

## Architecture (Design Consensus)

From debate rounds (see YAPP-PUSH-BUTTONS-001 ticket for full debates):

### Key Design Decisions

1. **All modules use YAML schemas** - No hybrid approach with struct tags. Uniform architecture simplifies tooling.

2. **Code generation via `schemagen` tool**:
   - `schemagen validate schema.yaml` - Validates schema structure
   - `schemagen discover` - Finds all modules, generates registry code
   - Generates: typed Go structs, ApplyDefaults(), test cases, registration code

3. **Top-level `pkg/registry` package** - Breaks circular dependencies between resolver, yappgen, and modules. Contains shared interfaces (ModuleSchema, FeatureModule, FieldSpec).

4. **Two-phase validation**:
   - Phase 1 (pre-resolution): Structure checks (types, required fields)
   - Phase 2 (post-resolution): Constraint checks (min/max, enums)
   - Necessary because expressions like `x: base_x + 10` can't be validated until resolved

5. **Discovery-based registration** - Drop schema.yaml in modules/{name}/, run `go generate`, module auto-registers. No core file edits.

### Data Flow

```
User writes YAML → Resolver validates (Phase 1) → Resolve expressions → 
Validator checks constraints (Phase 2) → Builder converts to OpenSCAD → 
SCAD generation → STL rendering
```

## Implementation Roadmap

**Phase 1: Core Infrastructure (Weeks 1-2)**
- Build `pkg/registry` package (interfaces, registration API)
- Build `cmd/schemagen` tool (validate + discover commands)
- Test with minimal example module

**Phase 2: Pilot (Week 3)**
- Convert push_buttons to YAML schema (most complex module)
- Test full pipeline end-to-end
- Fix issues with generated code

**Phase 3: Rollout (Weeks 4-5)**
- Convert pcb_stands, connectors, snap_joins, cutouts to YAML schemas
- Test each conversion
- Document any patterns discovered

**Phase 4: Documentation & Polish (Week 6)**
- Integrate schema → help page rendering
- CI enforcement (go generate checks)
- Pre-commit hook templates
- Module authoring guide

## Essential Documents

**Start here:**
- `playbook/module-system-implementation-guide.md` - Complete MVP guide for newcomers

**Design rationale:**
- See YAPP-PUSH-BUTTONS-001 ticket:
  - `reference/debate-round-01-validation-parsing-flow.md` - Validation architecture
  - `reference/debate-round-02-registration-mechanics.md` - Registration API
  - `reference/debate-round-03-prototype-implementation.md` - Prototype findings
  - `design/feature-module-registry.md` - Registry design doc

## Success Criteria

Implementation is complete when:

1. ✅ `pkg/registry` package exists with interfaces
2. ✅ `schemagen` tool validates and discovers schemas
3. ✅ Push buttons converted to YAML schema
4. ✅ Generated code compiles and tests pass
5. ✅ Resolver validates with two phases
6. ✅ `go run ./cmd/yappctl help push_buttons` shows auto-generated docs
7. ✅ All 5 modules converted (push_buttons, pcb_stands, connectors, snap_joins, cutouts)
8. ✅ CI enforces `go generate` discipline
9. ✅ Error messages include line numbers and helpful suggestions
10. ✅ Module authoring guide complete

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field
- **Parent Ticket**: YAPP-PUSH-BUTTONS-001 (contains design debates)

## Status

Current status: **active**

## Topics

- yapp
- dsl
- codegen
- architecture

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- playbooks/ - Implementation guide (START HERE)
- design/ - Architecture and design documents (if needed)
- reference/ - API specs, contracts (if needed)
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
