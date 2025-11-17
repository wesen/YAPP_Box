# Changelog

## 2025-11-17

- Initial workspace created


## 2025-11-17

Created ticket with architecture guide, debate references, and task list for Path 1 (A→B) refactor implementation

### Related Files

- /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/design/01-architecture-implementation-guide-path-1-a-b-builder-contract-refactor.md — Implementation guide with pseudocode
- /home/manuel/code/others/YAPP_Box/ttmp/2025/11/17/YAPP-BUILDER-REFACTOR-001-implement-path-1-a-b-builder-contract-refactor-generate-decode-and-arraydecl-ir/tasks.md — High-level task list


## 2025-11-17

Created intern handoff playbook with reading order, testing strategy, and timeline

### Related Files

- /home/manuel/code/others/YAPP_Box/ttmp/2025/11/17/YAPP-BUILDER-REFACTOR-001-implement-path-1-a-b-builder-contract-refactor-generate-decode-and-arraydecl-ir/playbook/01-intern-handoff-getting-started-with-path-1-a-b-refactor.md — Handoff guide


## 2025-11-17

Added shared decode helper package with tests

### Related Files

- /home/manuel/code/others/YAPP_Box/pkg/yappgen/decode/helpers.go — Implements field extraction helpers
- /home/manuel/code/others/YAPP_Box/pkg/yappgen/decode/helpers_test.go — Covers helper behaviors


## 2025-11-17

Updated schemagen to emit Decode functions and tests

### Related Files

- /home/manuel/code/others/YAPP_Box/pkg/schemagen/codegen.go — Provides template inputs for Decode
- /home/manuel/code/others/YAPP_Box/pkg/schemagen/templates/schema_gen.go.tmpl — Generates decode helpers per module
- /home/manuel/code/others/YAPP_Box/pkg/schemagen/templates/schema_gen_test.go.tmpl — Generates Decode validation tests


## 2025-11-17

Regenerated module schemas/tests and updated all module Build() implementations to consume the new Decode() output

### Related Files

- /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/cutouts/module.go — Translates decoded items into per-face arrays
- /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/pcbstands/module.go — Uses typed Decode() slice
- /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/pushbuttons/module.go — Removed bespoke yaml decode logic

