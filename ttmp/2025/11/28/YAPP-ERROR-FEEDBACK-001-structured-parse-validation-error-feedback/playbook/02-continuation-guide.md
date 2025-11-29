---
Title: Continuation Guide - Recent Progress
Ticket: YAPP-ERROR-FEEDBACK-001
Status: active
Topics:
    - yapp
    - dx
    - errors
DocType: playbook
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/resolver/errorx/taxonomy_test.go
      Note: Comprehensive taxonomy tests
    - Path: pkg/resolver/positions.go
      Note: |-
        PositionMap type moved here
        PositionMap type for position tracking
    - Path: pkg/resolver/rules/enum_suggest_test.go
      Note: Enum suggest rule tests
    - Path: pkg/resolver/rules/module_doc.go
      Note: ModuleDocEmbedRule implementation
    - Path: pkg/resolver/rules/registry_test.go
      Note: Registry tests
    - Path: pkg/resolver/rules/vars_scaffold_test.go
      Note: Vars scaffold rule tests
    - Path: pkg/resolver/rules/yaml_syntax_test.go
      Note: YAML syntax rule tests
    - Path: pkg/resolver/validation_extract.go
      Note: |-
        New extraction functions for enum/constraints
        Enum and constraint extraction functions
    - Path: playbook/01-intern-onboarding-guide.md
      Note: Read this first for full context
ExternalSources: []
Summary: Updated guide covering recent progress - position integration, enhanced schema validation, unit tests, and ModuleDocEmbedRule
LastUpdated: 2025-11-29T18:00:00-05:00
---


# Continuation Guide: Recent Progress on YAPP-ERROR-FEEDBACK-001

## Purpose

This document updates the original onboarding guide (`01-intern-onboarding-guide.md`) with all the work completed since it was written. Read the original guide first for full context, then use this document to understand what's been added.

## What's Been Completed Since the Original Guide

### 1. Position Integration ✅ COMPLETE

**What:** All error types now include accurate line/column positions from the source YAML file.

**Files Created/Modified:**
- `pkg/resolver/positions.go` - **NEW**: PositionMap type moved here from resolvercli to avoid circular dependencies
- `pkg/resolver/errorx/constructors.go` - Updated all constructors to accept `line, column int` parameters
- `pkg/resolver/resolver.go` - Updated `ResolveResult()` and `resolvePass()` to accept and use `PositionMap`
- `pkg/resolver/validation.go` - Updated to accept `PositionMap` and look up positions
- `pkg/resolver/errorx/format.go` - Updated to display line/column in text and JSON output
- `pkg/cli/resolvercli/resolver.go` - Updated to use `resolver.PositionMap` and pass it through
- `pkg/cli/resolvercli/positions.go` - Updated to use `resolver.PositionMap`

**Key Changes:**
- `PositionMap` moved from `pkg/cli/resolvercli` to `pkg/resolver` (shared type)
- All taxonomy constructors now accept line/column: `NewSchemaConstraintTaxonomy(..., line, column int)`
- Position lookup happens when creating taxonomy entries: `positions.GetPosition(path)`
- Positions are displayed in both text and JSON format output

**Testing:**
```bash
# Test with missing required field
go run ./cmd/yappctl resolve -i test-missing-var.yaml --show-taxonomy
# Should show Line: 2, Column: 3 (or appropriate position)

# Test with enum error
go run ./cmd/yappctl resolve -i /tmp/test-enum-cutouts.yaml --show-taxonomy --format json | jq '.context.line, .context.column'
# Should show accurate line/column numbers
```

### 2. Enhanced Schema Validation Taxonomy ✅ COMPLETE

**What:** Schema validation errors now extract enum values, min/max constraints, and actual values from error messages and schema metadata.

**Files Created:**
- `pkg/resolver/validation_extract.go` - **NEW**: Extraction functions for parsing error messages and querying schema

**Files Modified:**
- `pkg/resolver/validation.go` - Updated `validateConstraints()` to use `extractConstraintInfo()`

**Key Functions:**
- `extractConstraintInfo()` - Main coordinator function
- `extractAllowedValues()` - Parses "(allowed: value1, value2, ...)" from error messages
- `extractActualValue()` - Extracts invalid value from error message or data
- `findFieldConstraints()` - Queries schema Fields() for min/max constraints
- `extractFieldPathFromError()` - Extracts field path from error messages

**How It Works:**
1. When `ValidateConstraints()` fails, it calls `extractConstraintInfo()`
2. Function parses error message for enum values (pattern: `(allowed: ...)`)
3. Queries schema `Fields()` to find field and get min/max constraints
4. Extracts actual value from error message or data structure
5. Populates `SchemaConstraintContext` with all available information

**Testing:**
```bash
# Create test file with invalid enum
cat > /tmp/test-enum.yaml << 'EOF'
features:
  cutouts:
    - face: front
      from_face_left: 10
      from_face_top: 10
      width: 20
      height: 30
      shape: invalid_shape
EOF

# Should show allowed enum values in taxonomy
go run ./cmd/yappctl resolve -i /tmp/test-enum.yaml --show-taxonomy --format json | jq '.context.allowed'
# Output: ["rectangle", "circle", "rounded_rect", ...]
```

### 3. Comprehensive Unit Tests ✅ COMPLETE

**What:** Full test coverage for taxonomy and rules packages.

**Files Created:**
- `pkg/resolver/errorx/taxonomy_test.go` - **NEW**: 15 tests covering all constructors and AsTaxonomy()
- `pkg/resolver/rules/registry_test.go` - **NEW**: 5 tests for registry matching/sorting
- `pkg/resolver/rules/yaml_syntax_test.go` - **NEW**: 5 tests for YAML syntax rule
- `pkg/resolver/rules/enum_suggest_test.go` - **NEW**: 7 tests for enum suggestion rule
- `pkg/resolver/rules/vars_scaffold_test.go` - **NEW**: 6 tests for vars scaffold rule

**Test Coverage:**
- **Taxonomy tests:** All 6 constructors, context type assertions, AsTaxonomy() unwrapping (direct, wrapped, double-wrapped), Error() method, edge cases
- **Registry tests:** Registration, no matches, single match, multiple matches, sorting by score and severity
- **Rule tests:** Match() methods, Render() methods, edge cases (empty snippets, no allowed values, etc.)

**Running Tests:**
```bash
# Run all taxonomy tests
go test ./pkg/resolver/errorx/... -v

# Run all rules tests
go test ./pkg/resolver/rules/... -v

# Run specific test
go test ./pkg/resolver/errorx/... -run TestAsTaxonomy_DoubleWrapped
```

### 4. ModuleDocEmbedRule ✅ COMPLETE

**What:** Rule that embeds complete module field tables in error messages, showing users all available fields when validation errors occur.

**Files Created:**
- `pkg/resolver/rules/module_doc.go` - **NEW**: ModuleDocEmbedRule implementation

**Files Modified:**
- `pkg/resolver/rules/default.go` - Registered ModuleDocEmbedRule
- `pkg/resolver/rules/registry.go` - Added deduplication to prevent duplicate results

**Key Features:**
- Matches both structure and constraint validation errors
- Loads schema documentation using `schemagen.LoadSchemaDoc()`
- Renders markdown table with all fields: name, type, required, default, description
- Handles nested objects recursively (e.g., `mask.preset`)
- Shows enum values inline in type column
- Provides action link to full module documentation

**How It Works:**
1. Rule matches `StageSchemaStructure` or `StageSchemaConstraints`
2. Extracts module name from taxonomy context
3. Loads schema.yaml file using `schemagen.LoadSchemaDoc()`
4. Renders fields table using `renderFieldsTable()` (similar to `schema_help.go`)
5. Returns RuleResult with severity Info (informational, not error)

**Testing:**
```bash
# Test with enum error - should show module doc
go run ./cmd/yappctl resolve -i /tmp/test-enum-cutouts.yaml
# Should show "Module documentation: cutouts" with full field table

# Test with missing required field - should also show module doc
go run ./cmd/yappctl resolve -i /tmp/test-missing-field.yaml
# Should show module documentation once (deduplicated)
```

### 5. Registry Deduplication ✅ COMPLETE

**What:** Fixed issue where same rule could appear multiple times in CLI output.

**Files Modified:**
- `pkg/resolver/rules/registry.go` - Added deduplication logic in `RenderAll()`

**Implementation:**
- Uses headline+body as deduplication key
- Prevents same rule from appearing twice even if it matches multiple times
- Maintains sorting by score and severity

## Code Structure: New Files

### Core Infrastructure
```
pkg/resolver/
├── positions.go              # NEW: PositionMap type (moved from resolvercli)
├── validation_extract.go      # NEW: Enum/constraint extraction functions
└── validation.go              # MODIFIED: Uses extractConstraintInfo()
```

### Taxonomy Tests
```
pkg/resolver/errorx/
└── taxonomy_test.go           # NEW: 15 comprehensive tests
```

### Rules
```
pkg/resolver/rules/
├── module_doc.go              # NEW: ModuleDocEmbedRule
├── registry.go                # MODIFIED: Added deduplication
├── registry_test.go           # NEW: Registry tests
├── yaml_syntax_test.go        # NEW: YAML syntax rule tests
├── enum_suggest_test.go       # NEW: Enum suggest rule tests
└── vars_scaffold_test.go      # NEW: Vars scaffold rule tests
```

## How to Test Everything

### Prerequisites
```bash
# Build the CLI
go build ./cmd/yappctl

# Or run directly
go run ./cmd/yappctl --help
```

### Test 1: Position Integration

**Create test file:** `test-positions.yaml`
```yaml
features:
  cutouts:
    - face: front
      from_face_left: 10
      # Missing required fields
```

**Test commands:**
```bash
# Should show line/column in taxonomy
go run ./cmd/yappctl resolve -i test-positions.yaml --show-taxonomy

# Should show line/column in JSON
go run ./cmd/yappctl resolve -i test-positions.yaml --show-taxonomy --format json | jq '.context.line, .context.column'
```

**Expected:** Line and column numbers should be > 0 (not 0,0)

### Test 2: Enhanced Schema Validation

**Create test file:** `test-enum-extraction.yaml`
```yaml
features:
  cutouts:
    - face: front
      from_face_left: 10
      from_face_top: 10
      width: 20
      height: 30
      shape: invalid_shape  # Invalid enum
```

**Test commands:**
```bash
# Should show allowed enum values
go run ./cmd/yappctl resolve -i test-enum-extraction.yaml --show-taxonomy --format json | jq '.context.allowed'

# Should show enum suggest rule with closest match
go run ./cmd/yappctl resolve -i test-enum-extraction.yaml
```

**Expected:** 
- Taxonomy should include `allowed: ["rectangle", "circle", ...]`
- Enum suggest rule should suggest closest valid value
- Module doc rule should show field table

### Test 3: ModuleDocEmbedRule

**Test commands:**
```bash
# Test with enum error
go run ./cmd/yappctl resolve -i /tmp/test-enum-cutouts.yaml

# Test with missing required field
go run ./cmd/yappctl resolve -i /tmp/test-missing-field.yaml
```

**Expected:**
- Should show "Module documentation: cutouts" section
- Should include complete field table with all fields
- Should appear only once (deduplicated)
- Should include action link to full docs

### Test 4: All Rules Together

**Create test file:** `test-comprehensive.yaml`
```yaml
features:
  cutouts:
    - face: front
      from_face_left: 10
      width: max(vars.missing_height, 10)  # Missing variable
      height: max(vars.missing_width, 20)
      shape: invalid_shape  # Invalid enum
```

**Test commands:**
```bash
go run ./cmd/yappctl resolve -i test-comprehensive.yaml
```

**Expected:** Should show multiple rules:
- VarsScaffoldRule (for missing variables)
- EnumSuggestClosestRule (for invalid enum)
- ModuleDocEmbedRule (for module documentation)
- All should be deduplicated and properly sorted

### Test 5: Unit Tests

```bash
# Run all tests
go test ./pkg/resolver/errorx/... ./pkg/resolver/rules/... -v

# Run with coverage
go test ./pkg/resolver/errorx/... ./pkg/resolver/rules/... -cover

# Run specific test
go test ./pkg/resolver/errorx/... -run TestAsTaxonomy
```

## Remaining Tasks

### Task 1: Implement DependencyGraphRule ⚠️ NEXT PRIORITY

**What:** Show dependency chains for missing variables to help users understand resolution order.

**Why:** When variables are missing, users need to understand which variables depend on which, and in what order they should be defined.

**How:**
1. Extract dependency graph from resolver trace
2. Build dependency chain visualization
3. Show which variables are missing and what depends on them
4. Suggest resolution order

**Files to create:**
- `pkg/resolver/rules/dependency_graph.go`

**Example output:**
```
Missing variables: vars.height, vars.width

Dependency chain:
  vars.height → features.cutouts[0].width
  vars.width → features.cutouts[0].height

Suggested resolution order:
  1. Define vars.height first
  2. Then vars.width
```

**Where to look:**
- `pkg/resolver/trace.go` - Trace recording infrastructure
- `pkg/resolver/resolver.go` - How dependencies are tracked
- `pkg/resolver/rules/vars_scaffold.go` - Similar rule for reference

### Task 2: Implement YamlKnownFieldsRule

**What:** Suggest known fields when unknown keys are detected.

**Why:** Users often make typos in field names. This rule would suggest valid field names.

**Requirements:**
- Requires yaml.v3 KnownFields support (already enabled in resolvercli)
- Need to extract unknown keys from error messages
- Use Levenshtein distance (like EnumSuggestClosestRule) to suggest closest field names

**Files to create:**
- `pkg/resolver/rules/yaml_known_fields.go`

**Where to look:**
- `pkg/cli/resolvercli/resolver.go` - KnownFields is already enabled
- `pkg/resolver/rules/enum_suggest.go` - Similar pattern for suggestions
- `pkg/resolver/strict.go` - Strict mode validation (might have unknown key detection)

## Key Files to Understand

### Position Tracking
- `pkg/resolver/positions.go` - PositionMap type and GetPosition() method
- `pkg/cli/resolvercli/positions.go` - buildPositionMap() function that creates the map
- `pkg/cli/resolvercli/resolver.go` - Where PositionMap is created and passed to resolver

### Schema Validation Enhancement
- `pkg/resolver/validation_extract.go` - All extraction logic
- `pkg/resolver/validation.go` - Where extraction is called
- `pkg/registry/schema.go` - ModuleSchema interface and Fields() method
- `pkg/schemagen/schema_doc.go` - SchemaDoc structure for loading schema.yaml files

### Rules System
- `pkg/resolver/rules/registry.go` - Registry with deduplication
- `pkg/resolver/rules/module_doc.go` - ModuleDocEmbedRule implementation
- `pkg/resolver/rules/cli.go` - Text rendering for CLI output
- `pkg/resolver/rules/default.go` - Default registry setup

### Taxonomy
- `pkg/resolver/errorx/taxonomy.go` - Core types and AsTaxonomy()
- `pkg/resolver/errorx/constructors.go` - All constructor functions
- `pkg/resolver/errorx/contexts.go` - Context structs with Line/Column fields
- `pkg/resolver/errorx/format.go` - Display functions

## Debugging Tips

### Check Position Extraction
```bash
# Add debug logging in validation.go
log.Printf("Looking up position for path: %s", path)
line, column := positions.GetPosition(path)
log.Printf("Found: line=%d, column=%d", line, column)
```

### Check Enum Extraction
```bash
# Test extraction function directly
go test ./pkg/resolver/... -run TestExtractAllowedValues -v
```

### Check Rule Matching
```bash
# Use --show-taxonomy to see raw taxonomy
go run ./cmd/yappctl resolve -i your-file.yaml --show-taxonomy --format json | jq

# Check which rules match
# Add logging in registry.go RenderAll():
log.Printf("Rule %T matched with score %d", rule, score)
```

### Verify Module Doc Loading
```bash
# Check if schema file exists
ls pkg/yappgen/modules/cutouts/schema.yaml

# Test LoadSchemaDoc directly
go test ./pkg/schemagen/... -run TestLoadSchemaDoc -v
```

## Common Pitfalls

1. **Position lookup returns 0,0:** Make sure PositionMap is passed through the entire pipeline. Check that `positions != nil` before calling `GetPosition()`.

2. **Enum extraction fails:** Error message format might have changed. Check the actual error message format in `pkg/yappgen/modules/*/schema_validate.go` to see how errors are formatted.

3. **Module doc not loading:** Schema file might not exist for legacy modules. Rule should handle this gracefully (returns error, registry continues).

4. **Duplicate results:** Deduplication uses headline+body as key. If same rule produces different bodies, they'll both appear (which might be correct).

5. **Field path extraction:** Field paths can be complex (`features.cutouts[0].shape`). Make sure extraction handles array indices correctly.

## Getting Help

- **Design questions:** See `design-doc/` directory
- **Code examples:** See `design-doc/03-examples-taxonomy-to-rules-mapping.md`
- **Testing:** Use `--show-taxonomy` flag extensively
- **Schema docs:** See `pkg/docs/schema_help.go` for how module docs are rendered
- **Registry patterns:** Look at existing rules in `pkg/resolver/rules/`

## Summary Checklist

Before starting work, make sure you:
- [ ] Read `playbook/01-intern-onboarding-guide.md` for full context
- [ ] Understand position tracking (positions.go, how PositionMap flows through pipeline)
- [ ] Understand schema validation enhancement (validation_extract.go, how enum/constraints are extracted)
- [ ] Understand rules system (registry.go, how rules match and render)
- [ ] Can run all tests successfully
- [ ] Can test functionality with CLI using various error scenarios
- [ ] Know where to add new rules (rules/ directory, register in default.go)
- [ ] Understand taxonomy structure (errorx/ directory)

Good luck! 🚀

