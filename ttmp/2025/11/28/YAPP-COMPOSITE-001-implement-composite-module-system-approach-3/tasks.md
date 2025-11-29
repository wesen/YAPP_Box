# Tasks

## Phase 1: Core Infrastructure (Week 1, Days 1-3)

- [ ] Create pkg/composite/ package with CompositeModule interface
- [ ] Implement Registry for composite module management
- [ ] Implement ProcessCompositeModules function
- [ ] Implement type conversion helpers (ToFloat, ToString, ToArray, etc.)
- [ ] Implement map extraction helpers (GetFloat, GetString, GetMap, etc.)
- [ ] Implement provenance tracking system
- [ ] Add 4-line integration hook in pkg/yappgen/model.go
- [ ] Write unit tests for registry (register, lookup, duplicates)
- [ ] Write unit tests for helpers (type conversions, edge cases)
- [ ] Write unit tests for provenance tracking
- [ ] Verify: go test ./pkg/composite/ passes

## Phase 2: Tooling and Documentation (Week 1, Days 4-5)

- [ ] Implement schemagen docs command for array format reference
- [ ] Implement schemagen builders command for array construction helpers
- [ ] Generate pkg/composite/builders.go with CreateXXXEntry functions
- [ ] Generate pkg/docs/composite-array-formats.md reference
- [ ] Review and verify generated documentation is comprehensive
- [ ] Verify: All array types have builder helpers

## Phase 3: First Composite Module - LCD (Week 2, Days 1-2)

- [ ] Create pkg/composite/modules/lcd/ package
- [ ] Implement LCD Module PostProcess function
- [ ] Add display cutout generation logic
- [ ] Add mounting hole generation logic
- [ ] Implement error handling and validation
- [ ] Write unit tests for LCD module
- [ ] Create examples/composite/test-lcd-basic.yaml
- [ ] Create examples/composite/test-lcd-complex.yaml
- [ ] Verify: End-to-end test YAML → SCAD works
- [ ] Verify: Generated SCAD renders in OpenSCAD

## Phase 4: Documentation and Examples (Week 2, Days 3-5)

- [ ] Write composite module authoring guide
- [ ] Create multiple example YAML files (simple, complex, with user entries)
- [ ] Write troubleshooting guide and FAQ
- [ ] Add integration tests with multiple composite modules
- [ ] Document migration path to Approach 1
- [ ] Update main documentation with composite module references
- [ ] Verify: New developer can follow guide to create second composite module
