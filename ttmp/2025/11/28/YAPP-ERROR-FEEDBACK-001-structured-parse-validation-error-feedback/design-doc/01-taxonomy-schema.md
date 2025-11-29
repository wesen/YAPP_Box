---
Title: Taxonomy schema
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
      Note: CLI consumer that renders taxonomy payloads
    - Path: pkg/resolver/resolver.go
      Note: Emits structural and expression errors that must be classified
    - Path: pkg/resolver/validation.go
      Note: Schema validation phases produce taxonomy entries
ExternalSources: []
Summary: Defines the canonical taxonomy schema (stage/symptom enums + typed contexts) used to classify resolver and validation errors.
LastUpdated: 2025-11-28T18:44:12.260450537-05:00
---


# Taxonomy schema

## Background & Context

**What is YAPP?** YAPP (Yet Another Parametric Parts) is a DSL (Domain-Specific Language) for defining 3D-printed enclosures. Users write YAML files describing features like cutouts, connectors, and mounting points, and YAPP generates OpenSCAD code.

**What is the resolver?** The `pkg/resolver` package processes YAML documents in multiple phases:
1. **Ingest**: Read and parse YAML into Go data structures (`pkg/cli/resolvercli/resolver.go`)
2. **Structure validation** (Phase 1): Check that fields exist and have correct types before expression evaluation (`pkg/resolver/validation.go`)
3. **Expression resolution**: Evaluate expressions like `"max(vars.height, 10)"` iteratively until all dependencies resolve (`pkg/resolver/resolver.go`)
4. **Constraint validation** (Phase 2): After resolution, check min/max/enum constraints against resolved numeric values (`pkg/resolver/validation.go`)

**What are modules?** Modules (e.g., `connectors`, `boxmounts`) define reusable feature types. Each module has a schema (`pkg/yappgen/modules/*/schema.yaml`) that specifies required fields, types, enums, and constraints.

**Why taxonomy?** Currently, errors are plain strings. We want structured metadata (stage, symptom, path, context) so the CLI can show targeted help, examples, and remediation steps.

## Executive Summary

- Establish a first-class taxonomy schema so every resolver/validator failure emits structured, machine-readable metadata.
- Define enumerations for stages and symptoms, typed context payloads per stage, consumed by CLI/IDE tooling.
- Introduce helper interfaces so future components (e.g., IDE adapters) can rely on compile-time coverage instead of ad-hoc `map[string]any` payloads.

## Problem Statement

- Errors today are plain strings built from `github.com/pkg/errors`, making it impossible to differentiate ingest, schema, or expression failures without string parsing.
- The analysis doc already outlines desired stages and typed contexts, but we need a concrete schema (Go structs + JSON shape) to unblock implementation.
- Without a shared schema, rule registration, telemetry, and UI rendering will diverge across teams.

## Proposed Solution

1. **Core types**
   - `StageCode` and `SymptomCode` are typed string enums (e.g., `StageIngestYAMLSyntax`, `SymptomMissingRequired`), validated via exhaustive `switch` statements.
     - **Stage** = pipeline phase (ingest, schema validation, expression resolution, strict mode)
     - **Symptom** = specific failure type (syntax error, missing required field, enum mismatch, dependency missing)
   - `Taxonomy` struct wraps stage, symptom, canonical DSL path, severity, and a `TaxonomyContext` interface that is implemented per stage.
     - **Path** = dot-notation DSL path like `features.connectors[0].corner` (used for highlighting in UI)
     - **Severity** = `error`, `warning`, or `info` (strict-mode violations are warnings)
   - Each context type embeds stage-specific data; e.g., `YAMLIngestContext` carries file/line/column/snippet, `SchemaConstraintContext` provides module, expected type, allowed values, and actual value.
     - **Why typed contexts?** Instead of `map[string]any`, each stage gets a concrete struct so Go's type system catches missing fields at compile time.
   - All contexts include optional `Raw` JSON to future-proof additions without schema churn.

2. **Go-first API (no wire format in v1)**
   - Render help directly in the CLI via the rule registry; no JSON output or external schema is required for the initial release.
   - If external consumers emerge (CI/IDE), we can add an optional `--json` output later without changing core Go types.

3. **Constructors + helpers**
   - Add factory functions (`NewYAMLIngestTaxonomy`, `NewSchemaConstraintTaxonomy`, etc.) ensuring required fields are present.
   - Implement `errors.As` helpers (`AsTaxonomy(err)`) to unwrap nested errors and return the structured payload.

4. **Versioning**
   - Embed a `SchemaVersion` constant so future breaking changes can be negotiated.
   - Expose compatibility helpers (`SupportsContext(type) bool`) for clients to down-level gracefully.

## Design Decisions

- **Typed contexts vs. generic maps**: Enforces compile-time coverage, catches missing fields during unit tests, and makes IDE/CLI renderers straightforward.
  - Example: If a rule tries to access `yc.Line` on a `SchemaConstraintContext`, the compiler catches it immediately.
- **String-based enums**: We need human-readable codes across boundaries (CLI output, JSON logs); string constants are friendlier than numeric enums without runtime lookups.
  - Example: `StageIngestYAMLSyntax` is clearer than `Stage(1)` when debugging or logging.
- **Path canonicalization**: All taxonomy entries must use resolver-style dot paths (e.g., `features.cutouts[0].width`) so rule lookups and UI highlights align; implement a shared `CanonicalPath` helper.
  - The resolver already uses dot paths internally; we reuse that convention so rules can match on path patterns.
- **Severity scale**: Introduce `Severity` enum (`error`, `warning`, `info`) to represent strict-mode linting vs hard failures.
  - Strict-mode violations (unknown keys, unused vars) are warnings; parse/schema errors are hard failures.
- **Extensibility**: Context structs reserve `map[string]any Extra` fields for module-specific metadata, preventing schema churn when new concepts (e.g., geometry diffs) appear.
  - Modules can attach custom hints without changing core taxonomy types.

## Alternatives Considered

| Option | Why Rejected |
| --- | --- |
| Single struct with optional fields | Quickly becomes unreadable and unmaintainable; every new stage adds more nullable fields. |
| `map[string]any` context | Makes downstream code untyped, requires reflection or string keys, and prevents Go vet/staticcheck from helping. |
| Relying on protobuf/Cap’n Proto | Overkill for the immediate CLI use case; JSON is sufficient and already flows through docmgr + tooling. |
| No severity enum | Would force rule engine to infer severity heuristically; better to encode intent at the source. |

## Implementation Plan

1. **Create `pkg/resolver/errorx` package** (new directory)
   - Define `StageCode`, `SymptomCode`, `Severity` enums as typed string constants
   - Define `Taxonomy` struct and `TaxonomyContext` interface
   - Implement context structs: `YAMLIngestContext`, `SchemaConstraintContext`, `ExprDependencyContext`, `StrictModeContext`
   - Add constructors: `NewYAMLIngestTaxonomy(...)`, `NewSchemaConstraintTaxonomy(...)`, etc.

2. **Update error producers** (where errors are created)
   - `pkg/cli/resolvercli/resolver.go`: Wrap `yaml.Unmarshal` errors with `YAMLIngestContext` (preserve line/column from `yaml.Node`)
   - `pkg/resolver/validation.go`: Convert `schemagen.ValidationError` to `SchemaConstraintContext` taxonomy entries
   - `pkg/resolver/resolver.go`: Convert expression errors to `ExprDependencyContext` taxonomy entries
   - `pkg/resolver/strict.go`: Convert strict-mode violations to `StrictModeContext` taxonomy entries

3. **Add helper functions**
   - `AsTaxonomy(err error) (*Taxonomy, bool)`: Unwrap nested errors to find taxonomy payload (similar to `errors.As`)
   - `CanonicalPath(path string) string`: Normalize paths to resolver-style dot notation

4. **Documentation**
   - Add Markdown reference in `pkg/resolver/errorx/README.md` listing all stages, symptoms, and context fields
   - Link from CLI help (`yappctl help errors`)

5. **Testing**
   - Unit tests for each context type (verify required fields are set)
   - Test `AsTaxonomy` unwrapping through nested error chains
   - Exhaustive `switch` tests ensuring all enum values are handled

## Open Questions

- Do we need localized/translated symptom text, or will English strings suffice for now?
- Should severity default to `error` unless explicitly downgraded (e.g., strict mode warnings)?
- How do we model composite errors that aggregate multiple taxonomy entries (e.g., batch schema violations) without overwhelming the CLI?

## Navigation & Related Documents

**Start here if you're new:**
1. Read `analysis/01-parse-and-validation-error-taxonomy.md` first — it explains the problem and desired stages/symptoms
2. Then read this document (taxonomy schema) to understand the Go types
3. Next: `design-doc/02-rule-registry.md` — how rules consume taxonomy entries
4. Finally: `design-doc/03-examples-taxonomy-to-rules-mapping.md` — concrete examples

**Key code files to understand:**
- `pkg/resolver/resolver.go` — main resolution loop (produces expression/dependency errors)
- `pkg/resolver/validation.go` — schema validation phases (produces constraint errors)
- `pkg/cli/resolvercli/resolver.go` — CLI entry point (produces YAML parse errors)
- `pkg/schemagen/errors.go` — existing validation error type we're aligning with

**Where taxonomy will be used:**
- `pkg/resolver/rules` (to be created) — rule registry consumes taxonomy entries
- `cmd/yappctl` — CLI renders taxonomy via rules
