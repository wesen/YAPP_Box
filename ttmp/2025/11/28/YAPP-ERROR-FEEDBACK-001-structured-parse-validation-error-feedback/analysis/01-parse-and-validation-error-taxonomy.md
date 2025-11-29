---
Title: Parse and validation error taxonomy
Ticket: YAPP-ERROR-FEEDBACK-001
Status: active
Topics:
    - yapp
    - dx
    - errors
DocType: analysis
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/cli/resolvercli/resolver.go
      Note: CLI entry point where YAML parse errors bubble up
    - Path: pkg/resolver/resolver.go
      Note: Core resolution loop emitting expression and dependency failures
    - Path: pkg/resolver/strict.go
      Note: Strict-mode checks for unknown keys and unused vars
    - Path: pkg/resolver/validation.go
      Note: Phase 1/2 schema enforcement hooks for structure and constraints
    - Path: pkg/schemagen/errors.go
      Note: Existing structured validation error type we can align to
ExternalSources: []
Summary: "Defines a staged taxonomy and rule engine so parse/resolution errors can emit targeted guidance via the CLI/UI."
LastUpdated: 2025-11-28T17:43:18-05:00
---


# Parse and validation error taxonomy

## Overview

- Provide consistent, actionable error feedback for the YAPP DSL by enriching every failure with machine-readable context.
- `pkg/cli/resolvercli/resolver.go` now relies on `gopkg.in/yaml.v3`, so we have line/column/snippet metadata available at decode time, but we still drop that context before bubbling errors up.
- We want a taxonomy + rule engine so the CLI (and future UI surfaces) can map each failure to curated help, including dynamic snippets or schema references.

## Current error surfaces & gaps

### 1. YAML ingest + CLI plumbing
- `pkg/cli/resolvercli/resolver.go` uses `yaml.v3`, but we still unmarshal directly into `map[string]any`. That path throws away the node-level `Line`, `Column`, and `Kind` information `yaml.Node` exposes.
- As a result, both file I/O errors and YAML syntax issues still surface as a single plain string without pointing to the exact row/column or showing the problematic snippet, even though v3 gives us those hooks for free.
- There is no flag to hint whether the user should rerun generation, fix whitespace, or simply ensure the file exists.

### 2. Schema structure / constraints
- Phase 1/2 validation in `pkg/resolver/validation.go` surfaces only the first error per module and wraps it as `validate <path>: ...`, masking the `schemagen.ValidationError` fields (`Path`, `Message`, `Hint`, `Snippet`) defined in `pkg/schemagen/errors.go`.
- Constraint failures (min/max/enum) do not mention allowed ranges or how to inspect module docs.

### 3. Expression resolution + dependencies
- The resolver loop (`pkg/resolver/resolver.go`) can emit syntax errors (`isSyntaxError`), runtime type errors, or the aggregated `unresolved expressions` message when dependencies never converge.
- Missing variables are already collected via `findMissingDependencies`, but that list is only embedded inside the error string. The CLI cannot highlight which `vars.*` entries are needed or suggest scaffolding.

### 4. Strict-mode enforcement
- `pkg/resolver/strict.go` treats unknown top-level keys and unused vars as fatal errors, even though these are closer to lint violations. The user receives only the offending key list without remediation ideas (e.g., “move this field under `features`”).

### 5. Module help (untapped)
- `pkg/docs/schema_help.go` renders module reference docs from schema metadata, but we never link a failed field (e.g., `features.cutouts[0].shape`) back to its module page or example YAML. This is a missed opportunity for dynamic assistance.

## Proposed taxonomy

| Stage | Example sources | Primary dimensions | Example help |
| --- | --- | --- | --- |
| `ingest.yaml.syntax` | `yaml.Unmarshal` | file, line, column, snippet | Show offending lines + YAML tips |
| `ingest.fs` | `os.ReadFile` | path, syscall errno | Suggest checking path/permissions |
| `schema.structure` | `Schema.ValidateStructure` | module, field path, expected type, required flag | Render module reference + minimal valid block |
| `schema.constraints` | `Schema.ValidateConstraints` | min/max/enum metadata, actual value | Print allowed range and nearest valid values |
| `expr.syntax` | `expr.Compile` | expression text, token | Highlight token + link to expression grammar |
| `expr.runtime` | `expr.Run`, numeric casting | operand types, function names | Recommend casting/literals |
| `expr.dependency.missing` | `collectUnresolved`, `findMissingDependencies` | unresolved path, missing refs | Show which `vars.*` keys are absent + snippet to add |
| `strict.unknown-key` | `validateTopLevelKeys` | list of top-level offenders | Provide canonical key list + pointer to docs |
| `strict.unused-var` | `validateUnusedVars` | unused variables | Suggest removing or referencing them |

Every taxonomy entry should expose:
1. **Stage** – pipeline phase used for severity + grouping.
2. **Symptom** – short code (`missing_required`, `enum_mismatch`, `syntax`, etc.).
3. **Subject path** – canonical DSL path for downstream highlighting.
4. **Evidence** – structured payload (line number, actual value, list of missing refs).
5. **Remediation hooks** – references to schema docs, code examples, or commands.

## Rule + help system

- Introduce a `Taxonomy` struct returned by all validators:
  ```go
  // Stage-specific context payloads implement TaxonomyContext.
  type Taxonomy struct {
      Stage   StageCode
      Symptom SymptomCode
      Path    string
      Context TaxonomyContext
  }

  type TaxonomyContext interface {
      // Stage returns the owning StageCode so switch statements can stay type-safe.
      Stage() StageCode
  }

  type YAMLIngestContext struct {
      File   string
      Line   int
      Column int
      Snippet string
  }

  type SchemaConstraintContext struct {
      Module       string
      FieldPath    string
      ExpectedType string
      Allowed      []any
      Actual       any
  }

  type ExprDependencyContext struct {
      Expression string
      MissingRefs []string
      Iterations  int
  }
  ```
- Build a rule registry keyed by `(Stage, Symptom)` that emits:
  - A concise summary (`"cutouts[0].shape must be one of round, square"`).
  - A remediation body (Markdown) that can embed schema table rows or dynamic YAML scaffolds.
  - Optional “actions” (e.g., `yappctl explain --path features.cutouts`) for future interactive flows.
- Rules switch on the typed `TaxonomyContext` (e.g., `*YAMLIngestContext`, `*SchemaConstraintContext`) instead of downcasting `map[string]any`, ensuring compile-time coverage for new fields and enabling helpers like `Context.Line`.
- For taxonomy entries referencing specific modules, pull doc snippets from `pkg/docs/schema_help.go` so users see the relevant table or example automatically.
- For dependency errors, auto-generate YAML snippets showing how to declare missing variables, and include the computed dependency graph so advanced users can debug cycles.

## Implementation considerations

1. **Typed errors end-to-end**
   - Decode into a `yaml.Node` (or run `yaml.NewDecoder` with `KnownFields(true)`) so we can populate `IngestError` with exact `Line`, `Column`, and snippet data that v3 now provides.
   - Update schema validation to pass through `schemagen.ValidationErrors` so we retain structured fields.
   - In the resolver loop, replace generic `errors.Errorf` calls with constructors like `NewExprSyntaxError(path, exprText, token)` and `NewMissingDependencyError(path, missingRefs)`.
2. **Central classifier**
   - Create `pkg/resolver/errorx` (or similar) housing constructors + a `Classify(error) (*Taxonomy, bool)` helper that unwraps `github.com/pkg/errors`.
   - Expose typed interfaces (e.g., `interface{ Taxonomy() *Taxonomy }`) so non-resolver packages can emit compatible errors.
3. **Rule-driven renderer**
   - Implement `pkg/docs/helprules` that maps taxonomy entries to Markdown/JSON help blocks.
   - Allow data-driven rules (YAML templates) so non-Go contributors can expand guidance without releases.
4. **CLI + future UI integrations**
   - Enhance `cmd/yappctl` to detect taxonomy-aware errors:
     - Default output: summary + “Help” block.
     - `--json`: emit `{ stage, symptom, path, context, help }` for IDE adapters.
   - Reserve room for HTMX/Bootstrap UI cards (aligned with goGuidelines) where dynamic content (examples, code blocks) can be rendered later.
5. **Quality gates**
   - Add unit tests covering classification for each Stage/Symptom pair.
   - Add integration tests that feed malformed DSL files through `yappctl` and assert the structured JSON payload contains the right taxonomy fields.

## Next steps

1. Draft the canonical taxonomy schema (enum definitions, required context keys) and publish it under `pkg/docs`.
2. Refactor schema + resolver layers to emit typed errors and keep existing string messages as fallbacks for backward compatibility.
3. Implement the rule registry + renderer, starting with the three highest-impact categories (missing vars, enum mismatch, YAML syntax).
4. Update `yappctl` CLI output modes (text + JSON) to consume taxonomy metadata and render help blocks.
5. Collect user feedback, then iterate on additional rules (e.g., hints for face selectors, constraint visualizations) and consider integration with docs/search surfaces.
