# Changelog

## 2025-11-28

- Initial workspace created


## 2025-11-29

Implemented taxonomy schema (pkg/resolver/errorx) with stage/symptom enums and typed contexts. Updated error producers (resolvercli, resolver, validation, strict) to emit taxonomy entries. Created rule registry (pkg/resolver/rules) with three example rules (YAML syntax, enum suggest, vars scaffold). Wired into CLI with --show-taxonomy flag for raw output.

### Related Files

- cmd/yappctl/resolve_command.go — Added --show-taxonomy flag and rule rendering
- pkg/cli/resolvercli/resolver.go — Updated to emit YAML ingest taxonomy
- pkg/resolver/errorx/constructors.go — Factory functions for creating taxonomy entries
- pkg/resolver/errorx/contexts.go — Typed context structs for each stage
- pkg/resolver/errorx/format.go — Formatting functions for taxonomy display
- pkg/resolver/errorx/taxonomy.go — Core taxonomy types and enums
- pkg/resolver/resolver.go — Updated to emit expression dependency taxonomy
- pkg/resolver/rules/cli.go — CLI text renderer for rule results
- pkg/resolver/rules/default.go — Default registry with built-in rules
- pkg/resolver/rules/enum_suggest.go — Enum mismatch suggestion rule
- pkg/resolver/rules/registry.go — Rule registry with Renderer interface
- pkg/resolver/rules/vars_scaffold.go — Missing variable scaffold rule
- pkg/resolver/rules/yaml_syntax.go — YAML syntax error rule
- pkg/resolver/strict.go — Updated to emit strict-mode taxonomy
- pkg/resolver/validation.go — Updated to emit schema validation taxonomy


## 2025-11-29

Created intern onboarding guide documenting context, code structure, next steps, and testing procedures

### Related Files

- playbook/01-intern-onboarding-guide.md — Comprehensive guide for continuing work


## 2025-11-29

Completed position integration: Moved PositionMap to resolver package, updated taxonomy constructors to accept line/column parameters, passed PositionMap through resolver pipeline, and updated format functions to display positions. Schema validation and expression errors now include accurate line/column information.

### Related Files

- pkg/cli/resolvercli/positions.go — Updated to use resolver.PositionMap
- pkg/cli/resolvercli/resolver.go — Updated to use resolver.PositionMap and pass to ResolveResult
- pkg/resolver/errorx/constructors.go — Updated constructors to accept line/column parameters
- pkg/resolver/errorx/format.go — Added line/column output to text and JSON formats
- pkg/resolver/positions.go — New PositionMap type moved from resolvercli
- pkg/resolver/resolver.go — Pass PositionMap through ResolveResult and resolvePass
- pkg/resolver/validation.go — Accept PositionMap and look up positions for errors


## 2025-11-29

Added comprehensive unit tests for taxonomy package: Test all constructors create valid entries with correct fields, test AsTaxonomy() unwraps nested errors correctly, test context types implement TaxonomyContext interface, test Error() method. All 15 tests passing.

### Related Files

- pkg/resolver/errorx/taxonomy_test.go — New comprehensive test suite for taxonomy constructors and AsTaxonomy unwrapping


## 2025-11-29

Added comprehensive unit tests for rules package: Test registry matching/sorting/aggregation, test YAML syntax rule, test enum suggest rule with Levenshtein distance, test vars scaffold rule. All 24 tests passing.

### Related Files

- pkg/resolver/rules/enum_suggest_test.go — Tests for enum suggestion rule with distance calculation
- pkg/resolver/rules/registry_test.go — Tests for registry matching
- pkg/resolver/rules/vars_scaffold_test.go — Tests for vars scaffold rule
- pkg/resolver/rules/yaml_syntax_test.go — Tests for YAML syntax pointer rule


## 2025-11-29

Enhanced schema validation taxonomy to extract enum values and min/max constraints from error messages and schema metadata. Added extractConstraintInfo() function that parses error messages for allowed enum values, queries schema Fields() for min/max constraints, and extracts actual values. SchemaConstraintContext now properly populated with Allowed, Min, Max, and Actual fields.

### Related Files

- pkg/resolver/validation.go — Updated validateConstraints to use extractConstraintInfo
- pkg/resolver/validation_extract.go — New extraction functions for enum values and constraints


## 2025-11-29

Implemented ModuleDocEmbedRule: Rule that matches schema validation errors and embeds module field tables from schema.yaml files. Shows complete field reference table with types, required flags, defaults, and descriptions. Helps users understand available fields when validation errors occur.

### Related Files

- pkg/resolver/rules/default.go — Registered ModuleDocEmbedRule in default registry
- pkg/resolver/rules/module_doc.go — New rule that embeds module field tables


## 2025-11-29

Implemented ModuleDocEmbedRule: Rule that matches schema validation errors and embeds module field tables from schema documentation. Shows complete field reference with types, required flags, defaults, and descriptions. Helps users understand available fields when validation errors occur.

### Related Files

- pkg/resolver/rules/default.go — Registered ModuleDocEmbedRule in default registry
- pkg/resolver/rules/module_doc.go — New rule that embeds module field documentation


## 2025-11-29

Fixed duplicate rule results in registry: Added deduplication logic to prevent same rule from appearing multiple times in CLI output. ModuleDocEmbedRule now appears once per error, improving readability.

### Related Files

- pkg/resolver/rules/registry.go — Added deduplication by headline+body to prevent duplicate results


## 2025-11-29

Created continuation guide for interns: Comprehensive document covering all recent progress including position integration, enhanced schema validation, unit tests, ModuleDocEmbedRule, and registry deduplication. Includes testing procedures, debugging tips, and next steps for remaining tasks.

### Related Files

- playbook/02-continuation-guide.md — Comprehensive continuation guide for interns


## 2025-11-29

Implemented DependencyGraphRule: Shows dependency chains and resolution order for missing variables

### Related Files

- pkg/resolver/rules/dependency_graph.go — New rule implementation


## 2025-11-29

Completed DependencyGraphRule implementation. Note: Limited by YAPP-BUG-001 bug where expressions cannot be extracted from array paths


## 2025-11-29

DependencyGraphRule now fully functional after YAPP-BUG-001 fix. Can extract expressions from array paths.


## 2025-11-29

Updated DependencyGraphRule task: Now fully functional after YAPP-BUG-001 fix. Removed limitation note.


## 2025-11-29

Updated continuation guide: Marked DependencyGraphRule as complete and fully functional

