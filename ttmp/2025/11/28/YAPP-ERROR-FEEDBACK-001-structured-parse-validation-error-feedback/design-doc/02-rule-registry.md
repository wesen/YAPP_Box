---
Title: Rule registry
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
      Note: CLI surface that will render rule results
    - Path: pkg/docs/schema_help.go
      Note: Source for module metadata that template rules will embed
    - Path: pkg/resolver/resolver.go
      Note: Produces taxonomy errors consumed by the registry
ExternalSources: []
Summary: Outlines a plugin-style rule registry that maps taxonomy entries to rendered help, examples, and actions.
LastUpdated: 2025-11-28T18:44:14.680892022-05:00
---


# Rule registry

## Executive Summary

- Provide a structured rule registry that consumes taxonomy entries and emits human-friendly help blocks plus machine-readable actions.
- Support data-driven (YAML) rule definitions and Go plugins so module owners can extend guidance without touching the resolver.
- Ensure rules are versioned, testable, and hot-swappable for different surfaces (CLI text, JSON, HTMX).

## Problem Statement

- Even with a rich taxonomy schema, users still see bare error summaries because there is no central place to convert metadata into guidance.
- Existing docs are siloed (schema help, tutorials) and cannot be dynamically selected based on the error’s stage/symptom/path.
- We need a consistent API so CLIs, IDE extensions, and future UIs can all request “give me the best help for this taxonomy entry” and receive deterministic output.

## Proposed Solution

1. **Rule registry core**
   - Add `pkg/resolver/rules` that exposes:
     ```go
     type Renderer interface {
         Render(ctx context.Context, t *errorx.Taxonomy) (*RuleResult, bool, error)
     }
     type RuleResult struct {
         Headline string
         Body     template.HTML
         Severity errorx.Severity
         Actions  []Action // e.g., CLI commands, doc links
         Data     map[string]any // for JSON output
     }
     ```
   - Registry stores ordered rules; first rule returning `ok=true` wins. Rules can be registered via `func Register(stage errorx.StageCode, symptom errorx.SymptomCode, renderer Renderer)`.

2. **Rule types**
   - **Go rules**: full power for complex logic (e.g., generating snippets from `SchemaConstraintContext`).
   - **Template rules**: YAML/Markdown definitions loaded at startup, referencing `{{ .Context.FieldPath }}` etc. Use `text/template` with safe helpers.
   - **Module hooks**: modules can ship `rules.yaml` alongside schema docs to offer domain-specific guidance; loader automatically scopes them to relevant features.

3. **Rendering adapters**
   - CLI adapter renders `RuleResult` as text with optional ANSI colors.
   - JSON adapter emits structured data for IDEs.
   - Future HTMX adapter can convert the same `RuleResult` into cards or modals.

4. **Testing + validation**
   - Provide `rules/testkit` to load fixtures (taxonomy JSON + expected `RuleResult`) ensuring regressions are caught.
   - Add `docmgr doctor` hook verifying every rule references valid stage/symptom codes.

## Design Decisions

- **First-match wins** to keep evaluation predictable and avoid conflicting advice.
- **Context-aware templates**: Use typed contexts so templates can safely access fields (compile-time vetting via `go:generate` that runs template parsing with fake data).
- **Action abstraction**: Represent actionable follow-ups as `{ Label, Command, Args }` so UIs can offer buttons or copyable commands.
- **Data-driven overrides**: Allow deployment-specific overrides via `ttmp/.../rules/overrides.yaml`, enabling custom messaging without recompiles.
- **Metrics hook**: Optionally emit telemetry (`rule_applied`, `no_rule_found`) for observability.

## Alternatives Considered

| Option | Why Rejected |
| --- | --- |
| Hard-code help strings inside resolver errors | Couples business logic with UX, prevents localization or downstream customization. |
| Use only templates | Some rules need to compute ranges, generate snippets, or inspect registry state, which templates alone cannot handle cleanly. |
| Prioritize last-match-wins | Makes rule ordering brittle and difficult to reason about; first-match keeps authors honest about specificity. |
| No registry, each consumer implements its own mapping | Leads to divergence between CLI/IDE/web, defeating the goal of consistent help. |

## Implementation Plan

1. Scaffold `pkg/resolver/rules` with registry, renderer interface, and CLI adapter.
2. Port the top-priority rules (YAML syntax, schema enum mismatch, missing vars) as Go renderers.
3. Introduce template/YAML loader + validation pipeline; add sample module rule file.
4. Wire `cmd/yappctl` to call `rules.Render(ctx, taxonomy)` and print the result (text + JSON).
5. Document authoring workflow (how to add a new rule, required tests) and integrate with CI.

## Open Questions

- How should we handle multiple applicable rules (e.g., general schema guidance + module-specific tips)? Compose results or present both?
- Do we need localization support on day one, or can templates stay English-only?
- Should rule evaluation be deterministic across releases (e.g., stable sorting), or can we allow randomized suggestions for experimentation?

## References

- `design-doc/01-taxonomy-schema.md` — defines the taxonomy payload consumed by the registry.
- `analysis/01-parse-and-validation-error-taxonomy.md` — motivation and desired stages/symptoms.
- `pkg/docs/schema_help.go` — source of module metadata referenced by template rules.
