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


## 2025-11-17

Updated 04-features example/resolved fixtures to satisfy new cutout/light tube schemas and regenerated SCAD baseline

### Related Files

- /home/manuel/code/others/YAPP_Box/examples/04-features.resolved.yaml — Mirror resolved data for validation
- /home/manuel/code/others/YAPP_Box/examples/04-features.yaml — Uses new cutout + light_tube fields


## 2025-11-17

Step A validation complete: go test ./... and regenerated 04-features + MVP SCAD outputs with zero diffs vs /tmp/refactor-baseline

### Related Files

- /home/manuel/code/others/YAPP_Box/examples/04-features.yaml — Source for regenerated SCAD
- /home/manuel/code/others/YAPP_Box/examples/yapp-mvp-pcbstands-cutouts.yaml — Second baseline used for diff


## 2025-11-17

Implemented Step B scaffolding: added ArrayDecl IR, updated FeatureModule interface, rewrote module builders/registries to emit typed ArrayDecls, and refreshed feature wiring + tests (go test ./... passes)

### Related Files

- /home/manuel/code/others/YAPP_Box/pkg/registry/schema.go — ArrayDecl type + interface change
- /home/manuel/code/others/YAPP_Box/pkg/yappgen/features.go — ArrayDecl-aware array/multi-array helpers
- /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules — All module build/registry updates


## 2025-11-17

Validated Step B outputs: go test ./..., regenerated SCAD for 04-features + yapp-mvp-pcbstands-cutouts, diffs vs /tmp/refactor-baseline clean

### Related Files

- /home/manuel/code/others/YAPP_Box/ttmp/2025/11/17/YAPP-BUILDER-REFACTOR-001-implement-path-1-a-b-builder-contract-refactor-generate-decode-and-arraydecl-ir/various/2025-11-17-step-2-arraydecl-implementation-diary.md — Recorded validation run details


## 2025-11-17

Generated Build wrapper tests via schemagen (TestBuild_UsesGeneratedDecode + TestBuild_ReturnsArrayDecl) and regenerated all module schema tests; verified with go test ./pkg/yappgen/modules/... and go build ./...

### Related Files

- /home/manuel/code/others/YAPP_Box/pkg/schemagen/codegen.go — Supplies template data for new Build wrapper tests
- /home/manuel/code/others/YAPP_Box/pkg/schemagen/templates/schema_gen_test.go.tmpl — Emits the new test cases
- /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules — Regenerated schema_gen_test.go files include the new coverage


## 2025-11-17

Generated schemagen Build wrapper tests (TestBuild_UsesGeneratedDecode + TestBuild_ReturnsArrayDecl), regenerated module schema tests, and verified with go test ./pkg/yappgen/modules/... + go build ./...

### Related Files

- /home/manuel/code/others/YAPP_Box/pkg/schemagen/codegen.go — Supplies template data for new wrapper tests
- /home/manuel/code/others/YAPP_Box/pkg/schemagen/templates/schema_gen_test.go.tmpl — Emits new Build wrapper test cases
- /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules — Regenerated schema_gen_test.go files now enforce ArrayDecl + Decode usage


## 2025-11-17

Refreshed module authoring guide Step 5/6 to show typed Build() + ArrayDecl registry wrappers (no yaml.Marshal) and logged the update in the diary; go test ./pkg/yappgen/modules/... && go build ./... still clean

### Related Files

- /home/manuel/code/others/YAPP_Box/pkg/docs/tutorials/yapp-module-authoring-guide.md — Builder/registry instructions now match Decode + ArrayDecl architecture
- /home/manuel/code/others/YAPP_Box/ttmp/2025/11/17/YAPP-BUILDER-REFACTOR-001-implement-path-1-a-b-builder-contract-refactor-generate-decode-and-arraydecl-ir/various/2025-11-17-step-2-arraydecl-implementation-diary.md — Diary entry for Task 17


## 2025-11-17

Added pkg/docs/migrations/path-1-a-b-refactor.md with the official Decode()+ArrayDecl migration checklist, updated the diary, and reran go test ./pkg/yappgen/modules/... && go build ./...

### Related Files

- /home/manuel/code/others/YAPP_Box/pkg/docs/migrations/path-1-a-b-refactor.md — New migration guide
- /home/manuel/code/others/YAPP_Box/ttmp/2025/11/17/YAPP-BUILDER-REFACTOR-001-implement-path-1-a-b-builder-contract-refactor-generate-decode-and-arraydecl-ir/various/2025-11-17-step-2-arraydecl-implementation-diary.md — Logged migration guide work

