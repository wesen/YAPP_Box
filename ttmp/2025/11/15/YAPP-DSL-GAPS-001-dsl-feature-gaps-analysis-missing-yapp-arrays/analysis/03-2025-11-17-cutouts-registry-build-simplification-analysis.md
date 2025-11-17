---
Title: Cutouts Module - registry.Build Simplification Analysis
Ticket: YAPP-DSL-GAPS-001
Status: draft
Intent: analysis
Authors: []
LastUpdated: 2025-11-17
---

## Purpose

Understand what “simplify cutouts module registry.Build” means, why the current design is awkward, and propose concrete simplification options with trade‑offs and a recommended next step.

## Current Architecture (facts)

- The central registry defines a `FeatureModule` interface whose `Build` returns a flat 2D array for SCAD emission.

```35:58:pkg/registry/schema.go
// implementations are typically generated from YAML schema definitions.
type ModuleSchema interface {
	// Name returns the short module identifier (e.g. "push_buttons").
	Name() string
	// Path returns the fully-qualified registry path (e.g. "features.push_buttons").
	Path() string
	// Description returns a human-readable summary of the module.
	Description() string
	// Fields returns a flattened view of the schema fields for documentation.
	Fields() []FieldSpec
	// ValidateStructure performs phase-1 validation (shape/type checking).
	ValidateStructure(path string, data any) error
	// ValidateConstraints performs phase-2 validation (post-resolution checks).
	ValidateConstraints(path string, data any) error
}

// FeatureModule ties together a schema with the logic that converts validated
// DSL entries into OpenSCAD parameter arrays.
type FeatureModule interface {
	// Schema returns the module's schema metadata.
	Schema() ModuleSchema
	// Build consumes validated items and produces OpenSCAD parameter arrays.
	Build(items []map[string]any) ([][]any, error)
}
```

- All modules are auto-registered (including `cutouts`) via codegen:

```16:24:pkg/yappgen/modules_gen.go
func init() {
	registry.Register(boxmounts.NewModule())
	registry.Register(connectors.NewModule())
	registry.Register(cutouts.NewModule())
	registry.Register(lighttubes.NewModule())
	registry.Register(pcbstands.NewModule())
	registry.Register(pushbuttons.NewModule())
	registry.Register(snapjoins.NewModule())
}
```

- The cutouts module’s “real” builder returns a map of arrays keyed by face (legacy YAPP expects per‑face arrays), not a single array:

```13:23:pkg/yappgen/modules/cutouts/module.go
// Build converts DSL cutouts entries into the YAPP array format.
// Note: Cutouts are distributed by face, so this returns a map.
func Build(items []map[string]any) (map[string][][]any, error) {
	byFace := map[string][][]any{
		"cutoutsFront": {},
		"cutoutsBack":  {},
		"cutoutsLeft":  {},
		"cutoutsRight": {},
		"cutoutsLid":   {},
		"cutoutsBase":  {},
	}
```

- The registry-facing `Build` for cutouts is a stub that discards output, because emission is handled specially elsewhere:

```31:42:pkg/yappgen/modules/cutouts/registry.go
func (m *module) Build(items []map[string]any) ([][]any, error) {
	// Note: cutouts return a map by face, not a simple array
	// This is a special case handled in features.go
	byFace, err := Build(items)
	if err != nil {
		return nil, err
	}
	// For now, return empty to satisfy interface
	// The actual cutout handling is in features.go cutoutFeatureModule
	_ = byFace
	return nil, nil
}
```

- The generator has a special `cutoutFeatureModule` that calls the real builder and emits multiple arrays by face:

```26:49:pkg/yappgen/features.go
var featureModules = []FeatureModule{
	newArrayFeatureModule("pcb_stands", "pcbStands",
		func(m *Model) *[]map[string]any { return &m.PcbStands },
		pcbstands.Build, nil),
	// ...
	newArrayFeatureModule("light_tubes", "lightTubes",
		func(m *Model) *[]map[string]any { return &m.LightTubes },
		lighttubes.Build, nil),
	newCutoutFeatureModule(),
}
```

```160:181:pkg/yappgen/features.go
func (m *cutoutFeatureModule) Emit(ctx context.Context, model *Model, b *strings.Builder) error {
	if len(model.Cutouts) == 0 {
		return nil
	}
	byFace, err := cutouts.Build(model.Cutouts)
	if err != nil {
		return errors.Wrap(err, "cutouts")
	}
	keys := make([]string, 0, len(byFace))
	for k := range byFace {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if len(byFace[k]) == 0 {
			continue
		}
		writeArrayDecl(b, k, byFace[k])
		b.WriteString("\n")
	}
	return nil
}
```

## Problem Statement

- Interface mismatch: `registry.FeatureModule.Build` returns `[][]any`, but cutouts inherently produce a map keyed by face. The registry-level cutouts `Build` thus returns `(nil, nil)` after computing and discarding `byFace`.
- Hidden special-case: Emission is done via a dedicated `cutoutFeatureModule`, not via the generic array path. This is non-obvious unless you read `features.go`.
- Potential pitfalls:
  - Silent no‑op `Build` may confuse future contributors or lead to accidental misuse.
  - Unnecessary work in the stub `Build` (computing `byFace` just to drop it).

## What “simplify cutouts module registry.Build” can mean

1) Make the contract explicit (minimal change)
- Change the registry‑facing cutouts `Build` to return a clear error explaining it is unsupported via the generic array path and is handled by `cutoutFeatureModule` instead.
- Stop computing `byFace` in the stub (avoid wasted work and underscore the point).
- Add a short doc comment in `registry.go` and cross‑link to `features.go`.

Pros: Simple, low‑risk, prevents silent misuse.  
Cons: Does not unify the abstraction; still a special case.

2) Introduce an adapter layer (moderate change)
- Define an optional, more general interface (e.g., `FacePartitionedModule` with `BuildByFace`) that the registry can detect via type assertion.
- Keep `FeatureModule.Build` for most modules; for cutouts, consumers (like `features.go`) downcast and use `BuildByFace`.

Pros: Makes the special shape explicit in types; improves discoverability.  
Cons: Extra complexity in registry and call sites; still two shapes to account for.

3) Unify emission shape (bigger refactor)
- Flatten cutouts to a single array and encode face as a parameter in each row; emit once under `cutouts`.
- Update SCAD side to accept a unified `cutouts` array and split internally.

Pros: Removes generator special‑case; single uniform path.  
Cons: Requires downstream SCAD changes; diverges from legacy “six arrays” API; higher risk.

## Recommendation

Adopt Option 1 now (minimal, explicit, safe). It delivers immediate clarity:
- Make `registry.go` cutouts `Build` return a descriptive error and avoid computing `byFace`.
- Add a brief module‑level comment pointing maintainers to `features.go` `cutoutFeatureModule` for emission.

Consider Option 2 later if multiple modules need non‑array outputs.

## Acceptance Criteria

- The cutouts registry `Build` no longer computes and discards `byFace`.
- Calling `Build` via the registry for cutouts returns a clear, actionable error message.
- A short comment documents that cutouts are emitted via `cutoutFeatureModule`.
- No behavior change in generated SCAD (cutouts still emitted per face arrays).

## References (code)

See code excerpts above for:
- `registry.FeatureModule` interface
- `modules_gen.go` registration
- `cutouts/module.go` real builder returning `map[string][][]any`
- `cutouts/registry.go` stub `Build` (target of simplification)
- `features.go` `cutoutFeatureModule` special emission

## Next Steps

- If we choose Option 1:
  - Change `pkg/yappgen/modules/cutouts/registry.go` `Build` to:
    - Not call `Build(items)` from `module.go`
    - Return `nil, errors.New("cutouts module is emitted via cutoutFeatureModule; registry.Build not supported")`
    - Add a docstring noting the special emission path
  - Add a brief note to the module authoring guide about “face‑partitioned” modules.

- If we choose Option 2:
  - Design `FacePartitionedModule` (name TBD)
  - Update generator call sites to prefer `BuildByFace` when available
  - Keep backward compatibility for existing modules

- Defer Option 3 unless we decide to change the legacy SCAD API.

## Context Links

- Ticket index: `index.md` (this directory)
- Task list item: “Simplify cutouts module registry.Build or document special handling”
- Cutouts schema source: `pkg/yappgen/modules/cutouts/schema.yaml`


