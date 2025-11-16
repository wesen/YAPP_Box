# Tasks

## TODO


- [x] Week 1: Create pkg/registry package with ModuleSchema, FeatureModule, and FieldSpec interfaces
- [x] Week 1: Create Registry type with Register(), Get(), All() functions and ordering support
- [x] Week 1: Write unit tests for registry (registration, lookup, ordering)
- [x] Week 1: Create cmd/schemagen skeleton with CLI structure and command parsing
- [x] Week 2: Implement schemagen validate command with YAML parsing and line number tracking
- [x] Week 2: Add schema structure validation (field types, nested objects, test cases)
- [x] Week 2: Implement helpful error messages with context and suggestions
- [x] Week 2: Implement schemagen discover command with filesystem scanning
- [x] Week 2: Add code generation templates for Go structs from YAML schemas
- [x] Week 2: Add test code generation from schema test cases
- [x] Week 2: Implement modules_gen.go generation with imports and registration
- [x] Week 2: Validate generated Go code compiles (go/format, go/parser)
- [ ] Week 3: Create push_buttons schema.yaml with all fields, nested objects, and test cases
- [ ] Week 3: Run schemagen to generate schema_gen.go for push_buttons
- [ ] Week 3: Refactor push_buttons builder.go to use generated types
- [ ] Week 3: Implement CustomValidate() for conditional push_buttons validation
- [ ] Week 3: Test full YAML → validation → SCAD pipeline with push_buttons examples
- [ ] Week 3: Fix any issues discovered with generated code or validation flow
- [x] Week 4: Create pcb_stands schema.yaml and convert module
- [x] Week 4: Create connectors schema.yaml and convert module
- [x] Week 4: Create snap_joins schema.yaml and convert module
- [ ] Week 5: Create cutouts schema.yaml and convert module
- [ ] Week 5: Run full test suite across all converted modules
- [ ] Week 5: Verify all existing YAML examples still work with new validation
- [x] Week 6: Implement pkg/resolver/validation.go with validateAgainstSchemas() and validateConstraints()
- [ ] Week 6: Integrate two-phase validation into Resolve() function
- [ ] Week 6: Implement pkg/docs/schema_help.go for schema → markdown rendering
- [ ] Week 6: Integrate LoadModuleHelp() with Glazed help system
- [ ] Week 6: Add CI workflow to check go generate is up-to-date
- [ ] Week 6: Create pre-commit hook template for auto-running go generate
- [ ] Week 6: Write module authoring guide explaining how to add new modules
- [ ] Week 6: Verify all 10 success criteria are met
