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

## Background & Context

**What is a rule?** A rule is a piece of code that:
1. **Matches** a taxonomy entry (e.g., "is this a YAML syntax error?")
2. **Renders** helpful output (e.g., show the offending line with a pointer, suggest fixes)

**Why a registry?** Instead of hard-coding help strings in error messages, we centralize all help generation in one place. This allows:
- Multiple rules to contribute complementary guidance (e.g., general schema tips + module-specific examples)
- Rules to be added/modified without changing resolver code
- Different output formats (CLI text, future HTML/JSON) to reuse the same rules

**How does it work?** When an error occurs:
1. Resolver produces a `Taxonomy` entry (see `design-doc/01-taxonomy-schema.md`)
2. Registry evaluates all registered rules, finds matches
3. Matching rules render help cards
4. Results are aggregated and shown to the user

**What are modules?** Modules (e.g., `connectors`, `boxmounts`) define feature types. Each has a schema file (`pkg/yappgen/modules/*/schema.yaml`) and can register module-specific rules.

## Executive Summary

- Provide a structured rule registry that consumes taxonomy entries and emits human-friendly help blocks plus machine-readable actions.
- Support Go-based rules (v1) so module owners can extend guidance without touching the resolver.
- Ensure rules are testable and composable for different surfaces (CLI text, future HTMX).

## Problem Statement

- Even with a rich taxonomy schema, users still see bare error summaries because there is no central place to convert metadata into guidance.
- Existing docs are siloed (schema help, tutorials) and cannot be dynamically selected based on the error’s stage/symptom/path.
- We need a consistent API so CLIs, IDE extensions, and future UIs can all request “give me the best help for this taxonomy entry” and receive deterministic output.

## Proposed Solution

1. **Rule registry core**
   - Add `pkg/resolver/rules` that exposes:
     ```go
     // Go-first API: rules are Go implementations.
     type Renderer interface {
         // Match returns whether the rule applies and an optional score for ordering (higher wins).
         // Score helps prioritize: 100 = high priority (always show), 50 = supplementary (show if space)
         Match(t *errorx.Taxonomy) (ok bool, score int)
         // Render produces a help card for the taxonomy entry.
         // Context is passed for accessing registry state, module schemas, etc.
         Render(ctx context.Context, t *errorx.Taxonomy) (*RuleResult, error)
     }
     type RuleResult struct {
         Headline string        // Short summary (e.g., "YAML syntax error at line 5")
         Body     template.HTML // Detailed help (can include code blocks, tables, links)
         Severity errorx.Severity // error, warning, or info
         Actions  []Action      // e.g., CLI commands, doc links (for future UI buttons)
     }
     type Action struct {
         Label   string   // "View module docs"
         Command string   // "yappctl explain --module connectors"
         Args    []string // Command arguments
     }
     ```
   - **Registry workflow:**
     1. User calls `rules.RenderAll(ctx, taxonomy)`
     2. Registry evaluates all registered rules via `Match()`
     3. Collects matching rules, sorts by `(score desc, severity desc)`
     4. Calls `Render()` on each match
     5. Aggregates results (multiple rules can contribute)
     6. Returns list of `RuleResult` cards
   - Rules register via `Register(renderer Renderer)` (typically in `init()` functions)
   - Optional helpers allow scoping rules to specific `(Stage, Symptom)` for efficiency (skip evaluation if stage doesn't match)

2. **Rule types**
   - **Go rules (v1)**: canonical rule form for initial release; supports complex logic like schema-driven snippets and dependency graphs.
   - **Future: Template rules**: may be added later as thin wrappers around Go renderers for simpler cases.
   - **Future: Module hooks**: module packages can register rules in `init()` alongside their schemas.

3. **Rendering adapters**
   - CLI adapter renders an aggregated list of `RuleResult` cards with optional ANSI colors (multiple rules can contribute).
   - Future HTMX adapter can convert the same aggregated results into cards or modals.

4. **Testing + validation**
   - Provide `rules/testkit` to load fixtures (taxonomy JSON + expected `RuleResult`) ensuring regressions are caught.
   - Add `docmgr doctor` hook verifying every rule references valid stage/symptom codes.

## Design Decisions

- **Multi-match aggregation**: Several rules can add complementary guidance (e.g., schema enum tips + module-specific examples).
  - Example: An enum mismatch might trigger both `EnumSuggestClosestRule` (suggests "did you mean X?") and `ModuleDocEmbedRule` (shows full field table)
  - Rules are sorted by score, so high-priority rules appear first
- **Context-aware templates**: Use typed contexts so templates can safely access fields (compile-time vetting via `go:generate` that runs template parsing with fake data).
  - Instead of `map[string]any`, rules use `t.Context.(*SchemaConstraintContext)` with type assertions
  - This catches typos at compile time (e.g., `sc.ActualValue` vs `sc.Actual`)
- **Action abstraction**: Represent actionable follow-ups as `{ Label, Command, Args }` so UIs can offer buttons or copyable commands.
  - CLI can print "Run: yappctl explain --module connectors"
  - Future HTML UI can render a clickable button
- **Data-driven overrides**: Allow deployment-specific overrides via `ttmp/.../rules/overrides.yaml`, enabling custom messaging without recompiles.
  - Useful for enterprise deployments that want to add internal support links
  - Deferred to future phase (not in v1)
- **Metrics hook**: Optionally emit telemetry (`rule_applied`, `no_rule_found`) for observability.
  - Helps identify which errors are most common, which rules are most helpful
  - Deferred to future phase (not in v1)

## Alternatives Considered

| Option | Why Rejected |
| --- | --- |
| Hard-code help strings inside resolver errors | Couples business logic with UX, prevents localization or downstream customization. |
| Use only templates | Some rules need to compute ranges, generate snippets, or inspect registry state, which templates alone cannot handle cleanly. |
| Prioritize last-match-wins | Makes rule ordering brittle and difficult to reason about; first-match keeps authors honest about specificity. |
| No registry, each consumer implements its own mapping | Leads to divergence between CLI/IDE/web, defeating the goal of consistent help. |

## Implementation Plan

1. **Create `pkg/resolver/rules` package** (new directory)
   - Define `Renderer` interface and `RuleResult` struct
   - Implement `Registry` type with `Register()` and `RenderAll()` methods
   - Add CLI adapter: `RenderToText(results []RuleResult) string` (formats help cards for terminal)

2. **Implement first three rules** (see `design-doc/03-examples-taxonomy-to-rules-mapping.md` for details)
   - `YamlSyntaxPointerRule`: Shows YAML syntax errors with line/column pointers
   - `EnumSuggestClosestRule`: Suggests closest enum value for mismatches
   - `VarsScaffoldRule`: Generates YAML scaffold for missing variables
   - Each rule implements `Match()` (returns score) and `Render()` (returns help card)

3. **Wire into CLI**
   - Update `cmd/yappctl/resolve_command.go` (or similar) to:
     - Catch errors from resolver
     - Extract taxonomy via `errorx.AsTaxonomy(err)`
     - Call `rules.RenderAll(ctx, taxonomy)`
     - Print aggregated results before exiting

4. **Testing**
   - Unit tests for each rule (verify `Match()` logic and `Render()` output)
   - Integration test: feed malformed YAML through CLI, verify help output

5. **Documentation**
   - Add `pkg/resolver/rules/README.md` explaining how to add new rules
   - Include example rule implementation

## Open Questions

- How should we handle multiple applicable rules (e.g., general schema guidance + module-specific tips)? Compose results or present both?
- Do we need localization support on day one, or can templates stay English-only?
- Should rule evaluation be deterministic across releases (e.g., stable sorting), or can we allow randomized suggestions for experimentation?

## Navigation & Related Documents

**Reading order:**
1. `analysis/01-parse-and-validation-error-taxonomy.md` — problem statement
2. `design-doc/01-taxonomy-schema.md` — taxonomy types (what rules consume)
3. This document (rule registry) — how rules work
4. `design-doc/03-examples-taxonomy-to-rules-mapping.md` — concrete examples with code

**Key code files:**
- `pkg/resolver/rules` (to be created) — registry implementation
- `pkg/resolver/errorx` (to be created) — taxonomy types
- `pkg/docs/schema_help.go` — module field tables that rules can embed
- `cmd/yappctl` — CLI that will consume rule output

**Related concepts:**
- **Modules**: `pkg/yappgen/modules/*/schema.yaml` — module schemas that rules reference
- **Schema help**: `pkg/docs/schema_help.go` — renders module docs from schemas
