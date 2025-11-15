---
Title: Feature Module Registry
Ticket: YAPP-PUSH-BUTTONS-001
Status: active
Topics:
    - yapp
    - cli
    - dsl
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/yappgen/features.go
      Note: Central registry + FeatureModule interface
    - Path: pkg/yappgen/model.go
      Note: BuildModel entry point now delegates feature collection to registry
    - Path: pkg/yappgen/emit.go
      Note: Feature emitter shares registry plumbing
    - Path: pkg/yappgen/map.go
      Note: Builder helpers reused by registry modules
ExternalSources: []
Summary: "Describes the shared feature-module registry that standardizes how new DSL modules are collected and emitted so future YAPP features can be added by configuration instead of bespoke plumbing."
LastUpdated: 2025-11-15T16:40:12-05:00
---


# Feature Module Registry

## Executive Summary

Integrating a new DSL feature previously required copy/pasting bespoke code across `BuildModel`, `EmitSCAD`, and helper files—easy to miss steps and hard to reason about dependencies like `printSwitchExtenders`. The new registry centralizes that plumbing. Each feature implements a lightweight `FeatureModule` interface that knows how to (a) collect data from the resolved DSL and populate the `Model`, and (b) emit the resulting SCAD arrays. Built-in helpers cover the common “array-of-maps” modules so most additions are declarative: specify the DSL path, the SCAD array name, and the builder function.

## Problem Statement

- Adding `push_buttons` involved touching four files (model, emit, builder glue, CLI toggles) and duplicating boilerplate for array extraction.
- Future modules in the YAPP library will repeat that churn, increasing merge conflicts and divergence across feature implementations.
- Behavior such as toggling `printSwitchExtenders` was wired manually, so forgetting to add similar hooks for upcoming modules (e.g., `display_mounts`) is likely.

## Proposed Solution

Introduce a registry at `pkg/yappgen/features.go` with a `FeatureModule` interface:

```go
type FeatureModule interface {
    Name() string
    Collect(resolved map[string]any, features map[string]any, model *Model) error
    Emit(ctx context.Context, model *Model, b *strings.Builder) error
}
```

Key pieces:

1. **Array module helper** – `arrayFeatureModule` takes a DSL key (e.g., `push_buttons`), the SCAD array name, a getter/setter closure for the `Model` slice, the builder function (`buildPushButtons`), and an optional `afterCollect` hook. Registration is as simple as:

```go
newArrayFeatureModule(
    "push_buttons",
    "pushButtons",
    func(m *Model) *[]map[string]any { return &m.PushButtons },
    buildPushButtons,
    func(m *Model, items []map[string]any) {
        m.PrintSwitchExtenders = len(items) > 0
    },
)
```

2. **Specialized modules** – When a feature needs custom data (e.g., `cutouts` with face routing), it implements `FeatureModule` directly (`cutoutFeatureModule`). This pattern remains contained within `features.go`.

3. **Build/emit integration** – `BuildModel` now calls `collectFeatureModules`, and `EmitSCAD` calls `emitFeatureModules`. New features only update the registry; no more scattered edits.

4. **Error context** – Registry functions wrap errors with the feature name so CLI output points at the offending module.

5. **Extensibility** – Additional hooks (pre/post emit, CLI toggles) can extend the module struct without editing call sites.

## Design Decisions

- Kept the registry in `pkg/yappgen` rather than `pkg/cli` so both SCAD emission and any future renderers can reuse it.
- Opted for closures returning field pointers to avoid reflection and keep assignments type-safe.
- Maintained deterministic emission order by registering modules in the desired sequence.
- Accepted a specialized module type for `cutouts` instead of forcing it through the generic array helper, keeping the helper simple.

## Alternatives Considered

- **Code generation/template per feature** – Too heavy-weight; still requires editing multiple files.
- **Map-driven reflection** – Would avoid closures, but reflection-based setters add runtime costs and obscure compile-time checks.
- **Embedding feature logic in the resolver** – Resolver should stay focused on expression evaluation; mixing generation tasks would blur responsibilities.

## Implementation Plan

1. Implement `FeatureModule` interface and helper types in `pkg/yappgen/features.go`.
2. Register existing features (stands, connectors, snap joins, push buttons, cutouts).
3. Refactor `BuildModel` and `EmitSCAD` to defer to the registry.
4. Update unit tests (already covering push buttons) and add new docs describing how to register modules.
5. Provide a developer tutorial (via the Glazed help system) that walks through building/registering a new module.

### Reference implementation

The push button feature now ships as an external module under `pkg/yappgen/modules/pushbuttons`. It exports a pure builder (`pushbuttons.Build`) that the registry wires up via `newArrayFeatureModule`. Because the module only depends on the shared `scad` types, it proves that future DSL capabilities can live in their own sub-packages without creating circular imports. This is the pattern new modules should follow until we formalize a full plugin API.

## Open Questions

- Should modules support emit-time ordering hints beyond declaration order (e.g., allow grouping by SCAD dependency)?
- Do we need lifecycle hooks for validation or CLI-only behavior (e.g., `yappctl resolve --explain`), or is `Collect` sufficient?
- How should we expose registry-driven introspection to tooling (e.g., auto-generated docs of available features)?

## References

- `pkg/yappgen/features.go` – source of record for built-in modules.
- `pkg/docs/tutorials/yapp-dsl-reference.md` – user-facing DSL guide referencing push buttons and future modules.
