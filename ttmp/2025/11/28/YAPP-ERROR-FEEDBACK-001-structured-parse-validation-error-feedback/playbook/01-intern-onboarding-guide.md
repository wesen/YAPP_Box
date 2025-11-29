---
Title: Intern onboarding guide
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
    - Path: analysis/01-parse-and-validation-error-taxonomy.md
      Note: Start here - explains the problem
    - Path: cmd/yappctl/resolve_command.go
      Note: CLI integration
    - Path: design-doc/01-taxonomy-schema.md
      Note: Taxonomy schema design
    - Path: design-doc/02-rule-registry.md
      Note: Rule registry design
    - Path: design-doc/03-examples-taxonomy-to-rules-mapping.md
      Note: Concrete examples with code
    - Path: pkg/cli/resolvercli/resolver.go
      Note: YAML parsing and error production
    - Path: pkg/resolver/errorx
      Note: Core taxonomy types and contexts
    - Path: pkg/resolver/rules
      Note: Rule registry and example rules
ExternalSources: []
Summary: 'Updated guide covering recent progress: position integration, enhanced schema validation, unit tests, and ModuleDocEmbedRule'
LastUpdated: 2025-11-29T12:16:26.779084446-05:00
---



# Intern Onboarding Guide: Structured Error Feedback System

## Purpose

This guide helps new team members (especially interns) understand the structured error feedback system we've built, what's been completed, what remains, and how to continue the work.

## Context: What Problem Are We Solving?

**The Problem:** When users write YAML files for YAPP (our DSL for 3D-printed enclosures), errors are currently just plain strings like `"missing required field"`. Users can't tell:
- Where in the file the error occurred (line/column)
- What type of error it is (YAML syntax? Schema validation? Missing variable?)
- How to fix it (what values are allowed? what variables are missing?)

**Our Solution:** We've built a **taxonomy system** that classifies every error with structured metadata (stage, symptom, path, context), and a **rule registry** that converts taxonomy entries into helpful, actionable guidance.

## What We've Built (Completed)

### 1. Taxonomy Schema (`pkg/resolver/errorx/`)

**What it is:** A structured way to classify errors with metadata.

**Key files:**
- `taxonomy.go` - Core types: `StageCode`, `SymptomCode`, `Severity`, `Taxonomy`
- `contexts.go` - Typed context structs for each error stage:
  - `YAMLIngestContext` - YAML parsing errors
  - `SchemaConstraintContext` - Schema validation errors (enum mismatches, etc.)
  - `SchemaStructureContext` - Missing required fields, type mismatches
  - `ExprDependencyContext` - Missing variables in expressions
  - `ExprSyntaxContext` - Expression syntax errors
  - `StrictModeContext` - Strict-mode violations
- `constructors.go` - Factory functions to create taxonomy entries
- `format.go` - Functions to display taxonomy (text/JSON)

**Key concepts:**
- **Stage** = pipeline phase (ingest, schema validation, expression resolution)
- **Symptom** = specific failure type (syntax, missing_required, enum_mismatch, etc.)
- **Path** = DSL path like `features.connectors[0].corner`
- **Context** = stage-specific data (line numbers, allowed values, missing refs, etc.)

### 2. Error Producers (Updated to Emit Taxonomy)

**Files modified:**
- `pkg/cli/resolvercli/resolver.go` - YAML parsing now wraps errors with `YAMLIngestContext`
- `pkg/resolver/resolver.go` - Expression errors now use `ExprDependencyContext` / `ExprSyntaxContext`
- `pkg/resolver/validation.go` - Schema validation errors use `SchemaConstraintContext` / `SchemaStructureContext`
- `pkg/resolver/strict.go` - Strict-mode violations use `StrictModeContext`

**How it works:** When an error occurs, instead of returning a plain string, we create a `Taxonomy` entry with structured metadata and wrap it with `errors.Wrap()` so it can be unwrapped later.

### 3. Rule Registry (`pkg/resolver/rules/`)

**What it is:** A system that matches taxonomy entries to helpful guidance.

**Key files:**
- `registry.go` - Core `Renderer` interface and `Registry` type
- `cli.go` - Text formatter for CLI output
- `default.go` - Default registry with all built-in rules registered
- `yaml_syntax.go` - Rule for YAML syntax errors (shows snippet with pointer)
- `enum_suggest.go` - Rule for enum mismatches (suggests closest valid value)
- `vars_scaffold.go` - Rule for missing variables (generates YAML scaffold)

**How it works:**
1. Registry evaluates all registered rules via `Match(taxonomy)` 
2. Matching rules render help cards via `Render(taxonomy)`
3. Results are sorted by score/severity and aggregated
4. CLI adapter formats them as text

### 4. CLI Integration (`cmd/yappctl/resolve_command.go`)

**What it does:** When errors occur, extracts taxonomy and either:
- Shows raw taxonomy structure (`--show-taxonomy` flag)
- Executes rules and shows helpful guidance (default)

**New flag:** `--show-taxonomy` - Prints raw taxonomy in YAML or JSON format (useful for debugging)

### 5. Position Tracking (`pkg/cli/resolvercli/positions.go`)

**What it is:** Infrastructure to preserve line/column positions from YAML source through the pipeline.

**Status:** Infrastructure created (`PositionMap`, `DecodeResult` struct), but not yet fully integrated into taxonomy entries. See "Next Steps" below.

## Documents to Read (In Order)

1. **Start here:** `analysis/01-parse-and-validation-error-taxonomy.md`
   - Explains the problem and current gaps
   - Lists all error surfaces we need to handle

2. **Taxonomy design:** `design-doc/01-taxonomy-schema.md`
   - Explains the taxonomy schema structure
   - Documents all stage/symptom codes
   - Shows typed context structs

3. **Rule registry design:** `design-doc/02-rule-registry.md`
   - Explains how rules work
   - Documents the registry API
   - Shows multi-match aggregation

4. **Examples:** `design-doc/03-examples-taxonomy-to-rules-mapping.md`
   - Concrete examples with code
   - Shows how taxonomy entries map to rules
   - Implementation sketches

## Code Structure: Where to Look

### Core Taxonomy Types
```
pkg/resolver/errorx/
├── taxonomy.go          # Core types (StageCode, SymptomCode, Taxonomy)
├── contexts.go          # Typed context structs
├── constructors.go      # Factory functions
└── format.go           # Display functions
```

### Rule System
```
pkg/resolver/rules/
├── registry.go         # Registry and Renderer interface
├── cli.go              # CLI text formatter
├── default.go          # Default registry setup
├── yaml_syntax.go      # YAML syntax error rule
├── enum_suggest.go     # Enum suggestion rule
└── vars_scaffold.go    # Missing vars scaffold rule
```

### Error Producers (Where Errors Are Created)
```
pkg/cli/resolvercli/
├── resolver.go         # YAML parsing (creates YAMLIngestContext)
└── positions.go        # Position tracking infrastructure

pkg/resolver/
├── resolver.go         # Expression errors (ExprDependencyContext, ExprSyntaxContext)
├── validation.go      # Schema validation (SchemaConstraintContext, SchemaStructureContext)
└── strict.go          # Strict mode (StrictModeContext)
```

### CLI Integration
```
cmd/yappctl/
└── resolve_command.go  # CLI command with --show-taxonomy flag
```

## Next Steps / Remaining Tasks

### Task 1: Complete Position Integration ⚠️ HIGH PRIORITY

**What:** Use the `PositionMap` we built to populate `Line`/`Column` fields in taxonomy contexts.

**Why:** Currently, only YAML syntax errors have line numbers. Schema validation and expression errors show `Line: 0, Column: 0`.

**How:**
1. Pass `PositionMap` through resolver pipeline (currently only in `LoadResult`)
2. Look up positions when creating taxonomy entries in `validation.go` and `resolver.go`
3. Update constructors to accept line/column parameters

**Files to modify:**
- `pkg/resolver/resolver.go` - Accept `PositionMap` parameter, look up positions for expression errors
- `pkg/resolver/validation.go` - Accept `PositionMap` parameter, look up positions for schema errors
- `pkg/resolver/errorx/constructors.go` - Add line/column parameters to constructors

**Example:**
```go
// In validation.go
func validateStructure(doc map[string]any, positions PositionMap) error {
    // ... existing code ...
    if err := schema.ValidateStructure(path, value); err != nil {
        line, column := positions.getPosition(path)
        taxonomy := errorx.NewSchemaStructureTaxonomy(path, key, fieldPath, "", "", true, line, column)
        return errors.Wrapf(taxonomy, "validate %s", path)
    }
}
```

### Task 2: Add Unit Tests

**What:** Write tests for taxonomy and rules.

**Files to create:**
- `pkg/resolver/errorx/taxonomy_test.go` - Test taxonomy constructors, `AsTaxonomy()` unwrapping
- `pkg/resolver/rules/registry_test.go` - Test registry matching, sorting, aggregation
- `pkg/resolver/rules/yaml_syntax_test.go` - Test YAML syntax rule
- `pkg/resolver/rules/enum_suggest_test.go` - Test enum suggestion rule
- `pkg/resolver/rules/vars_scaffold_test.go` - Test vars scaffold rule

**Test cases to cover:**
- Taxonomy constructors create valid entries
- `AsTaxonomy()` unwraps nested errors correctly
- Rules match correct taxonomy entries
- Rules render helpful output
- Registry aggregates multiple matching rules correctly
- Position lookup works correctly

### Task 3: Enhance Schema Validation Taxonomy

**What:** Extract more information from schema validation errors (enum values, min/max constraints).

**Current state:** Schema validation errors are basic - we don't extract enum values or constraint ranges.

**How:**
- Parse error messages to extract allowed enum values
- Query schema metadata for min/max constraints
- Populate `SchemaConstraintContext.Allowed`, `Min`, `Max` fields

**Files to modify:**
- `pkg/resolver/validation.go` - Extract enum/constraint info from schema or error messages

### Task 4: Add More Rules

**What:** Implement additional helpful rules.

**Ideas:**
- `ModuleDocEmbedRule` - Embed module field tables from `pkg/docs/schema_help.go`
- `DependencyGraphRule` - Show dependency chains for missing variables
- `YamlKnownFieldsRule` - Suggest known fields when unknown keys detected

**Files to create:**
- `pkg/resolver/rules/module_doc.go`
- `pkg/resolver/rules/dependency_graph.go`
- `pkg/resolver/rules/yaml_known_fields.go`

## How to Test the Functionality

### Prerequisites

```bash
# Build the CLI
go build ./cmd/yappctl

# Or run directly
go run ./cmd/yappctl --help
```

### Test 1: YAML Syntax Errors

**Create test file:** `test-syntax.yaml`
```yaml
features:
  connectors:
    - x: 10
      y 20  # Missing colon - syntax error
```

**Test commands:**
```bash
# Show rule-based help (default)
go run ./cmd/yappctl resolve -i test-syntax.yaml

# Show raw taxonomy structure
go run ./cmd/yappctl resolve -i test-syntax.yaml --show-taxonomy

# Show raw taxonomy as JSON
go run ./cmd/yappctl resolve -i test-syntax.yaml --show-taxonomy --format json
```

**Expected output:**
- Default: Shows YAML syntax error rule output with snippet and pointer
- `--show-taxonomy`: Shows structured taxonomy with stage, symptom, path, context

### Test 2: Schema Validation Errors

**Create test file:** `test-schema.yaml`
```yaml
project: test
units: mm
yapp_version: 3.0

features:
  connectors:
    - x: 10
      y: 20
      # Missing required fields: stand_height, screw_d, etc.
```

**Test commands:**
```bash
go run ./cmd/yappctl resolve -i test-schema.yaml --show-taxonomy
```

**Expected output:** Taxonomy entry with `Stage: schema.structure`, `Symptom: missing_required`, path pointing to the missing field.

### Test 3: Missing Variables

**Create test file:** `test-vars.yaml`
```yaml
project: test
units: mm
yapp_version: 3.0

features:
  cutouts:
    - face: front
      from_face_left: 10
      from_face_top: 10
      width: max(vars.missing_height, 10)
      height: max(vars.missing_width, 20)
      shape: round
```

**Test commands:**
```bash
# Should show vars scaffold rule
go run ./cmd/yappctl resolve -i test-vars.yaml

# Show raw taxonomy
go run ./cmd/yappctl resolve -i test-vars.yaml --show-taxonomy --format json
```

**Expected output:** 
- Default: Shows `VarsScaffoldRule` output with YAML scaffold for missing variables
- `--show-taxonomy`: Shows `ExprDependencyContext` with missing refs list

### Test 4: Enum Mismatch (Future)

**Note:** This requires Task 3 (enhanced schema validation) to work properly.

**Create test file:** `test-enum.yaml`
```yaml
features:
  connectors:
    - x: 10
      y: 20
      stand_height: 5
      screw_d: 3
      screw_head_d: 6
      insert_d: 4
      outside_d: 8
      corner: diagonal  # Invalid - should be single, all, front_left, etc.
```

**Expected:** `EnumSuggestClosestRule` should suggest closest valid value.

## Debugging Tips

### Check if Taxonomy is Being Created

```bash
# Always use --show-taxonomy to see raw structure
go run ./cmd/yappctl resolve -i your-file.yaml --show-taxonomy --format json | jq
```

### Check Rule Matching

Add debug logging in `pkg/resolver/rules/registry.go`:
```go
func (r *Registry) RenderAll(ctx context.Context, taxonomy *errorx.Taxonomy) ([]*RuleResult, error) {
    for _, rule := range r.rules {
        ok, score := rule.Match(taxonomy)
        if ok {
            log.Printf("Rule %T matched with score %d", rule, score)
        }
    }
    // ... rest of function
}
```

### Verify Position Map

Add logging in `pkg/cli/resolvercli/positions.go`:
```go
func buildPositionMap(node *yaml.Node, path string, positions PositionMap) {
    if node.Kind == yaml.ScalarNode && path != "" {
        log.Printf("Path %s -> Line %d, Column %d", path, node.Line, node.Column)
    }
    // ... rest of function
}
```

## Common Pitfalls

1. **Forgetting to unwrap errors:** Use `errorx.AsTaxonomy(err)` to extract taxonomy from wrapped errors
2. **Path format:** Use resolver-style dot paths (e.g., `features.connectors[0].corner`)
3. **Context type assertions:** Always check `ok` when asserting context types: `ctx, ok := t.Context.(*SchemaConstraintContext)`
4. **Position lookup:** Remember that positions are only available if `PositionMap` is passed through the pipeline

## Getting Help

- **Design questions:** Consult the design documents in `design-doc/`
- **Code questions:** Check the examples in `design-doc/03-examples-taxonomy-to-rules-mapping.md`
- **Testing:** Use `--show-taxonomy` flag to inspect raw taxonomy structures
- **Related code:** See `pkg/docs/schema_help.go` for module documentation that rules can embed

## Summary Checklist

Before starting work, make sure you:
- [ ] Read all design documents (analysis + 3 design docs)
- [ ] Understand the taxonomy schema (stages, symptoms, contexts)
- [ ] Understand how rules work (Match + Render)
- [ ] Can run the CLI and test with `--show-taxonomy`
- [ ] Know where each error type is created (resolvercli, resolver, validation, strict)
- [ ] Understand the position tracking infrastructure

Good luck! 🚀
