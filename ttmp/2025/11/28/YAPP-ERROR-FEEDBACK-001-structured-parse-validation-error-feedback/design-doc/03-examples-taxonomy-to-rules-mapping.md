---
Title: 'Examples: taxonomy to rules mapping'
Ticket: YAPP-ERROR-FEEDBACK-001
Status: active
Topics:
    - yapp
    - dx
    - errors
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/cli/resolvercli/resolver.go
      Note: YAML v3 parse surface for ingest errors
    - Path: pkg/docs/schema_help.go
      Note: Field table embedding for module docs
    - Path: pkg/resolver/resolver.go
      Note: Expression errors and missing dependency reporting
    - Path: pkg/resolver/validation.go
      Note: Schema constraint failures feeding enum mismatch example
    - Path: pkg/yappgen/modules/connectors/schema.yaml
      Note: Enum and required fields referenced in examples
ExternalSources: []
Summary: Concrete examples mapping taxonomy entries to Go rules with aggregated outputs and implementation sketches.
LastUpdated: 2025-11-29T10:45:58.712170791-05:00
---


# Examples: taxonomy to rules mapping

## Background & Context

**What is this document?** This shows concrete, working examples of how taxonomy entries map to rules and what help output users see. Read this after understanding the taxonomy schema (`design-doc/01-taxonomy-schema.md`) and rule registry (`design-doc/02-rule-registry.md`).

**Key concepts:**
- **Taxonomy entry**: Structured error metadata (stage, symptom, path, context) produced by resolver/validator
- **Rule**: Code that matches a taxonomy entry and renders help
- **Aggregation**: Multiple rules can match the same error and contribute complementary guidance

**Why these examples?** These three categories cover ~80% of user errors:
1. YAML syntax errors (file can't be parsed)
2. Schema enum mismatches (user typed invalid value)
3. Missing variables (expression references undefined `vars.*`)

## Executive Summary

- Provide concrete examples of taxonomy entries, the Go rules that match them, and the aggregated help output.
- Focus on high-impact categories: YAML syntax (v3), schema enum mismatch (connectors), and missing variables during expression resolution.
- Sketch Go implementations leveraging typed contexts; multiple rules can match and contribute guidance.

## Examples

### 1) YAML ingest syntax error (yaml.v3)

**Scenario:** User has a YAML syntax error (e.g., missing colon, wrong indentation). The `yaml.v3` parser fails during unmarshaling.

**Taxonomy entry produced:**
```go
&errorx.Taxonomy{
    Stage:   errorx.StageIngestYAMLSyntax,
    Symptom: errorx.SymptomSyntax,
    Path:    "",  // file-level error, no specific path
    Context: &errorx.YAMLIngestContext{
        File:    "enclosure.yaml",
        Line:    12,
        Column:  5,
        Snippet: "features:\n  connectors:\n    - x: 10\n      y 20  # missing colon\n",
    },
}
```

**Matching Go rules:**
1. **`YamlSyntaxPointerRule`** (score 100 — high priority)
   - **Match logic**: `t.Stage == StageIngestYAMLSyntax && t.Symptom == SymptomSyntax`
   - **Headline**: "YAML syntax error at enclosure.yaml:12:5"
   - **Body**: Shows offending snippet with visual pointer (caret `^`) at column 5, suggests checking indentation and trailing commas
   - **Why score 100**: This is the primary error; always show it

2. **`YamlKnownFieldsRule`** (score 60 — supplementary, optional)
   - **Match logic**: Only if `yaml.v3` decoder has `KnownFields(true)` enabled and reports unknown keys
   - **Headline**: "Unknown field detected"
   - **Body**: Lists unknown keys at that line, suggests checking spelling or module docs
   - **Why score 60**: Helpful but secondary to syntax error

**Aggregated output:**
- Two stacked help cards: first shows syntax error with snippet, second mentions unknown fields if applicable
- User sees both because multiple rules matched and contributed guidance

### 2) Schema enum mismatch (connectors.corner)

**Scenario:** User typed `corner: diagonal` but `connectors` module only allows `["single", "all", "front_left", "front_right", "back_left", "back_right"]`. Phase 2 validation (after expression resolution) catches this.

**Taxonomy entry produced:**
```go
&errorx.Taxonomy{
    Stage:   errorx.StageSchemaConstraints,
    Symptom: errorx.SymptomEnumMismatch,
    Path:    "features.connectors[0].corner",
    Context: &errorx.SchemaConstraintContext{
        Module:    "connectors",
        FieldPath: "corner",
        Allowed:   []string{"single", "all", "front_left", "front_right", "back_left", "back_right"},
        Actual:    "diagonal",  // what user typed
    },
}
```

**Matching Go rules:**
1. **`EnumSuggestClosestRule`** (score 90 — high priority)
   - **Match logic**: `t.Stage == StageSchemaConstraints && t.Symptom == SymptomEnumMismatch`
   - **Headline**: "`corner` must be one of: single, all, front_left, front_right, back_left, back_right"
   - **Body**: 
     - Computes Levenshtein distance between `"diagonal"` and each allowed value
     - Finds closest match: `"front_left"` (or similar)
     - Shows: "Did you mean `front_left`?"
   - **Why score 90**: Primary guidance — shows what's wrong and suggests fix

2. **`ModuleDocEmbedRule`** (score 70 — supplementary)
   - **Match logic**: `t.Stage == StageSchemaConstraints && t.Context.(*SchemaConstraintContext).Module != ""`
   - **Headline**: "Connectors module reference"
   - **Body**: 
     - Calls `pkg/docs/schema_help.go` to render connectors module field table
     - Highlights the `corner` row (shows all enum values, description, default)
   - **Why score 70**: Helpful context but secondary to direct suggestion

**Aggregated output:**
- First card: Shows error + suggestion ("Did you mean `front_left`?")
- Second card: Shows full connectors module field table with `corner` row highlighted
- User gets both immediate fix suggestion and full documentation context

### 3) Missing variables for expression resolution

**Scenario:** User has expression `"max(vars.lid_thickness, 2)"` but `vars.lid_thickness` is never defined. After 16 resolution iterations, resolver gives up and reports unresolved dependencies.

**Taxonomy entry produced:**
```go
&errorx.Taxonomy{
    Stage:   errorx.StageExprDependencyMissing,
    Symptom: errorx.SymptomDependencyMissing,
    Path:    "enclosure.lid.thickness",  // where the expression lives
    Context: &errorx.ExprDependencyContext{
        Expression:  "max(vars.lid_thickness, 2)",
        MissingRefs: []string{"vars.lid_thickness"},
        Iterations:  16,  // how many resolution passes were attempted
    },
}
```

**Matching Go rules:**
1. **`VarsScaffoldRule`** (score 100 — high priority)
   - **Match logic**: `t.Stage == StageExprDependencyMissing`
   - **Headline**: "Declare missing variables"
   - **Body**: 
     - Extracts all `vars.*` references from `MissingRefs`
     - Generates YAML scaffold with example defaults:
       ```yaml
       vars:
         lid_thickness: 2  # example default (from expression fallback)
       ```
     - User can copy-paste this into their YAML file
   - **Why score 100**: Primary fix — shows exactly what to add

2. **`DependencyGraphRule`** (score 50 — supplementary)
   - **Match logic**: `t.Stage == StageExprDependencyMissing && len(t.Context.(*ExprDependencyContext).MissingRefs) > 1`
   - **Headline**: "Dependency resolution order"
   - **Body**: 
     - Analyzes expression to find dependency chain
     - Shows textual graph: `enclosure.lid.thickness` → `vars.lid_thickness` → (missing)
     - Suggests resolving upstream dependencies first
   - **Why score 50**: Helpful for complex cases with multiple missing vars, but secondary to scaffold

## Go Rule Sketches

These are simplified implementations showing the core structure. Full implementations would include error handling, edge cases, and helper functions.

```go
// YamlSyntaxPointerRule matches YAML syntax errors and shows the offending line with a pointer.
type YamlSyntaxPointerRule struct{}

// Match checks if this rule applies to the taxonomy entry.
// Returns (true, score) if it matches, (false, 0) otherwise.
// Score 100 = high priority (always show this rule's output).
func (r *YamlSyntaxPointerRule) Match(t *errorx.Taxonomy) (bool, int) {
    // Only match YAML syntax errors (not schema or expression errors)
    return t.Stage == StageIngestYAMLSyntax && t.Symptom == SymptomSyntax, 100
}

// Render produces the help card for this taxonomy entry.
func (r *YamlSyntaxPointerRule) Render(ctx context.Context, t *errorx.Taxonomy) (*RuleResult, error) {
    // Type assertion: we know from Match() that context is YAMLIngestContext
    yc, _ := t.Context.(*YAMLIngestContext)
    
    // Helper function (not shown) that formats snippet with visual pointer
    // Example output:
    //   features:
    //     connectors:
    //       - x: 10
    //         y 20  # missing colon
    //            ^
    body := renderSnippetWithPointer(yc.Snippet, yc.Line, yc.Column)
    
    return &RuleResult{
        Headline: fmt.Sprintf("YAML syntax error at %s:%d:%d", yc.File, yc.Line, yc.Column),
        Body:     body,
        Severity: errorx.SeverityError,
        Actions:  []Action{{Label: "Open file", Command: "open", Args: []string{yc.File}}},
    }, nil
}

// EnumSuggestClosestRule matches enum mismatches and suggests the closest valid value.
type EnumSuggestClosestRule struct{}

func (r *EnumSuggestClosestRule) Match(t *errorx.Taxonomy) (bool, int) {
    // Only match schema constraint errors with enum mismatches
    if t.Stage != StageSchemaConstraints || t.Symptom != SymptomEnumMismatch {
        return false, 0
    }
    // Verify context is the right type (defensive check)
    _, ok := t.Context.(*SchemaConstraintContext)
    return ok, 90  // Score 90 = high priority but slightly lower than syntax errors
}

func (r *EnumSuggestClosestRule) Render(ctx context.Context, t *errorx.Taxonomy) (*RuleResult, error) {
    sc := t.Context.(*SchemaConstraintContext)
    
    // Helper function (not shown) that computes Levenshtein distance
    // Example: nearestEnum(["front_left", "back_right"], "diagonal") → "front_left"
    suggestion := nearestEnum(sc.Allowed, fmt.Sprint(sc.Actual))
    
    // Helper function (not shown) that formats enum help message
    // Shows: "Did you mean `front_left`?" with list of all allowed values
    body := enumHelp(sc.FieldPath, sc.Allowed, suggestion)
    
    return &RuleResult{
        Headline: fmt.Sprintf("%s must be one of: %s", sc.FieldPath, strings.Join(toStrings(sc.Allowed), ", ")),
        Body:     body,
        Severity: errorx.SeverityError,
    }, nil
}

// VarsScaffoldRule matches missing variable errors and generates YAML scaffold.
type VarsScaffoldRule struct{}

func (r *VarsScaffoldRule) Match(t *errorx.Taxonomy) (bool, int) {
    // Match any missing dependency error (always applies)
    return t.Stage == StageExprDependencyMissing, 100  // Score 100 = always show
}

func (r *VarsScaffoldRule) Render(ctx context.Context, t *errorx.Taxonomy) (*RuleResult, error) {
    ec := t.Context.(*ExprDependencyContext)
    
    // Helper function (not shown) that generates YAML scaffold
    // Input: []string{"vars.lid_thickness"}
    // Output: "vars:\n  lid_thickness: 2  # example default"
    // Tries to infer default from expression (e.g., "max(vars.x, 2)" → default 2)
    yaml := scaffoldVars(ec.MissingRefs)
    
    return &RuleResult{
        Headline: "Declare missing variables",
        Body:     wrapYAML(yaml),  // Formats YAML in code block
        Severity: errorx.SeverityError,
    }, nil
}
```

## Aggregation Strategy

- Evaluate all registered rules; collect those where `Match` returns true.
- Sort by `(score desc, severity desc)`, deduplicate headlines, then render all.
- The CLI prints all resulting cards; adapters may collapse low-severity tips.

## Implementation Plan

1. Implement the three example rules above in `pkg/resolver/rules` with tests.
2. Wire `rules.RenderAll(ctx, taxonomy)` into `yappctl` error printing.
3. Add helpers for schema doc embedding (focus specific field rows).
4. Iterate based on feedback; extend with module-specific rules for other arrays (e.g., `boxmounts`).

## Navigation & Related Documents

**Reading order:**
1. `analysis/01-parse-and-validation-error-taxonomy.md` — problem statement and current gaps
2. `design-doc/01-taxonomy-schema.md` — taxonomy types (what errors produce)
3. `design-doc/02-rule-registry.md` — rule registry architecture
4. This document (examples) — concrete implementations

**Key code files referenced:**
- `pkg/cli/resolvercli/resolver.go` — YAML v3 ingest surface (produces `YAMLIngestContext`)
- `pkg/resolver/resolver.go` — expression resolution (produces `ExprDependencyContext`)
- `pkg/resolver/validation.go` — schema validation (produces `SchemaConstraintContext`)
- `pkg/docs/schema_help.go` — module field table renderer (used by `ModuleDocEmbedRule`)
- `pkg/yappgen/modules/connectors/schema.yaml` — example module schema with enum fields

**Where to implement:**
- Create `pkg/resolver/rules/` directory
- Add rule implementations: `yaml_syntax.go`, `enum_suggest.go`, `vars_scaffold.go`
- Register rules in `rules/registry.go` via `init()` functions
