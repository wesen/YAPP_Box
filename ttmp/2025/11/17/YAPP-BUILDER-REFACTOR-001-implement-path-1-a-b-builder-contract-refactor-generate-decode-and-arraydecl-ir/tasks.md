# Tasks

## TODO

- [x] Create shared decode helpers package (pkg/yappgen/decode/helpers.go) with GetFloat, GetString, GetBool, GetObject, GetArray and tests
- [x] Update schemagen templates to generate Decode() function in schema_gen.go.tmpl (~80 lines, handles nested objects and arrays)
- [x] Update schemagen templates to generate Decode() tests in schema_gen_test.go.tmpl (TestDecode_ValidInput, TestDecode_MissingRequired, TestDecode_WrongType)
- [x] Regenerate all 7 modules with new Decode() functions (run schemagen discover)
- [x] Update all module Build() functions to use generated Decode() instead of yaml.Marshal/Unmarshal (pcbstands, connectors, boxmounts, snapjoins, lighttubes, cutouts, pushbuttons)
- [x] Validate Step A: run tests and compare SCAD output with pre-refactor baseline
- [x] Define ArrayDecl type in pkg/registry/schema.go
- [x] Change registry.FeatureModule.Build interface to return []ArrayDecl instead of [][]any
- [x] Add ArrayDecl wrappers to all module registry.go files (6 single-array modules + 1 multi-array cutouts)
- [x] Update module Build() signatures to take typed items ([]PcbStandsItem, etc.) instead of []map[string]any
- [x] Create multiArrayFeatureModule helper in pkg/yappgen/features.go for cutouts/ridgeExt pattern
- [x] Update arrayFeatureModule to handle []ArrayDecl output
- [x] Update feature module registration in features.go (replace cutoutFeatureModule with multiArrayFeatureModule)
- [x] Validate Step B: run tests and compare SCAD output
- [ ] Add linter rules (.golangci.yml) to forbid yaml.Marshal/Unmarshal in module.go files
- [x] Generate Build wrapper tests (TestBuild_UsesGeneratedDecode, TestBuild_ReturnsArrayDecl)
- [x] Update module authoring guide (remove marshal/unmarshal pattern, show new Decode() pattern)
- [x] Create migration guide (pkg/docs/migrations/path-1-a-b-refactor.md) documenting the refactor
- [ ] Run complete test suite and end-to-end validation with all example YAMLs
- [ ] Test new module creation following updated guide
- [ ] Verify enforcement (linter catches violations, tests catch missing Decode())
- [ ] Update ticket changelog and relate all changed files

## Notes

**Estimated time:** 3-4 days

**Key references:**
- Architecture guide: [design/01-architecture-implementation-guide-path-1-a-b-builder-contract-refactor.md](../../../2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/design/01-architecture-implementation-guide-path-1-a-b-builder-contract-refactor.md)
- Debate decisions: See debate/ directory in parent ticket

**Success criteria:**
- All tests pass
- SCAD output identical to pre-refactor
- Linter passes
- No yaml.Marshal/Unmarshal in modules
- Documentation updated
