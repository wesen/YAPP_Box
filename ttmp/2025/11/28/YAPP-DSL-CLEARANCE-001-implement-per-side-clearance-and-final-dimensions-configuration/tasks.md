# Tasks

## TODO

- [x] Read design document (design/01-per-side-clearance-and-final-dimensions.md) and understand the requirements
- [x] Update Model struct in pkg/yappgen/model.go: Add PaddingFront, PaddingBack, PaddingLeft, PaddingRight, FinalLength, FinalWidth fields
- [x] Implement helper functions: hasPerSideClearance(), hasFinalDimensions(), hasUniformClearance() in pkg/yappgen/model.go
- [x] Add validateEnclosureDimensions() function in pkg/yappgen/validation.go (or new file) to check mutual exclusivity
- [x] Update BuildModel() to handle per-side clearance: Parse enclosure.wall.clearance object and populate PaddingFront/Back/Left/Right
- [x] Implement computePaddingFromDimensions() method: Compute padding backwards from final dimensions with even distribution
- [x] Update BuildModel() to handle final dimensions: Parse enclosure.dimensions and call computePaddingFromDimensions()
- [x] Add unit tests for validation: Test mutual exclusivity, incomplete per-side clearance, negative padding errors
- [x] Add unit tests for computation: Test forward (clearance→dimensions) and reverse (dimensions→padding) calculations
- [x] Add integration test: Generate SCAD from per-side clearance YAML and verify padding variables are correct
- [x] Add integration test: Generate SCAD from final dimensions YAML and verify computed padding produces correct shell dimensions
- [x] Test backward compatibility: Verify existing YAML files with uniform clearance still work unchanged
- [x] Update DSL reference documentation (pkg/docs/tutorials/yapp-dsl-reference.md): Add per-side clearance and final dimensions sections
- [x] Add examples to documentation: Show per-side clearance, final dimensions, and error cases
- [x] Update changelog with implementation summary
