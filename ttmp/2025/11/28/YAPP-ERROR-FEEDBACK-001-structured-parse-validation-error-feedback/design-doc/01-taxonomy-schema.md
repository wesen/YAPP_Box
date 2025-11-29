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

## Executive Summary

- Establish a first-class taxonomy schema so every resolver/validator failure emits structured, machine-readable metadata.
- Define enumerations for stages and symptoms, typed context payloads per stage, and a serialization contract consumed by CLI/IDE tooling.
- Introduce helper interfaces so future components (e.g., IDE adapters) can rely on compile-time coverage instead of ad-hoc `map[string]any` payloads.

## Problem Statement

- Errors today are plain strings built from `github.com/pkg/errors`, making it impossible to differentiate ingest, schema, or expression failures without string parsing.
- The analysis doc already outlines desired stages and typed contexts, but we need a concrete schema (Go structs + JSON shape) to unblock implementation.
- Without a shared schema, rule registration, telemetry, and UI rendering will diverge across teams.

## Proposed Solution

1. **Core types**
   - `StageCode` and `SymptomCode` are typed string enums (e.g., `StageIngestYAMLSyntax`, `SymptomMissingRequired`), validated via exhaustive `switch` statements.
   - `Taxonomy` struct wraps stage, symptom, canonical DSL path, severity, and a `TaxonomyContext` interface that is implemented per stage.
   - Each context type embeds stage-specific data; e.g., `YAMLIngestContext` carries file/line/column/snippet, `SchemaConstraintContext` provides module, expected type, allowed values, and actual value.
   - All contexts include optional `Raw` JSON to future-proof additions without schema churn.

2. **Serialization contract**
   - Provide `MarshalJSON`/`UnmarshalJSON` so contexts are serialized as `{ "stage": "schema.structure", "symptom": "missing_required", "path": "...", "context": { "type": "SchemaConstraintContext", ... } }`.
   - Keep the wire format stable so external tooling (GitHub Actions, IDEs) can consume it; publish a JSON Schema under `pkg/docs`.

3. **Constructors + helpers**
   - Add factory functions (`NewYAMLIngestTaxonomy`, `NewSchemaConstraintTaxonomy`, etc.) ensuring required fields are present.
   - Implement `errors.As` helpers (`AsTaxonomy(err)`) to unwrap nested errors and return the structured payload.

4. **Versioning**
   - Embed a `SchemaVersion` constant so future breaking changes can be negotiated.
   - Expose compatibility helpers (`SupportsContext(type) bool`) for clients to down-level gracefully.

## Design Decisions

- **Typed contexts vs. generic maps**: Enforces compile-time coverage, catches missing fields during unit tests, and makes IDE/CLI renderers straightforward.
- **String-based enums**: We need human-readable codes across boundaries (CLI output, JSON logs); string constants are friendlier than numeric enums without runtime lookups.
- **Path canonicalization**: All taxonomy entries must use resolver-style dot paths (e.g., `features.cutouts[0].width`) so rule lookups and UI highlights align; implement a shared `CanonicalPath` helper.
- **Severity scale**: Introduce `Severity` enum (`error`, `warning`, `info`) to represent strict-mode linting vs hard failures.
- **Extensibility**: Context structs reserve `map[string]any Extra` fields for module-specific metadata, preventing schema churn when new concepts (e.g., geometry diffs) appear.

## Alternatives Considered

| Option | Why Rejected |
| --- | --- |
| Single struct with optional fields | Quickly becomes unreadable and unmaintainable; every new stage adds more nullable fields. |
| `map[string]any` context | Makes downstream code untyped, requires reflection or string keys, and prevents Go vet/staticcheck from helping. |
| Relying on protobuf/Cap’n Proto | Overkill for the immediate CLI use case; JSON is sufficient and already flows through docmgr + tooling. |
| No severity enum | Would force rule engine to infer severity heuristically; better to encode intent at the source. |

## Implementation Plan

1. Define `StageCode`, `SymptomCode`, `Severity`, `Taxonomy`, and `TaxonomyContext` interfaces under `pkg/resolver/errorx`.
2. Implement context structs (`YAMLIngestContext`, `SchemaConstraintContext`, `ExprDependencyContext`, `StrictModeContext`, etc.) plus constructors and JSON marshalers.
3. Expose helper functions (`NewTaxonomy`, `AsTaxonomy(err)`) and update resolver/schemagen to return typed errors.
4. Publish a JSON Schema + Markdown reference documenting every field; link it from CLI help.
5. Add unit tests ensuring serialization round-trips, unknown context types panic during development, and exhaustive `switch` statements cover all enums.

## Open Questions

- Do we need localized/translated symptom text, or will English strings suffice for now?
- Should severity default to `error` unless explicitly downgraded (e.g., strict mode warnings)?
- How do we model composite errors that aggregate multiple taxonomy entries (e.g., batch schema violations) without overwhelming the CLI?

## References

- `analysis/01-parse-and-validation-error-taxonomy.md` — taxonomy overview and context rationale.
- `pkg/resolver/resolver.go`, `pkg/resolver/validation.go`, `pkg/cli/resolvercli/resolver.go` — primary producers/consumers of taxonomy errors.
