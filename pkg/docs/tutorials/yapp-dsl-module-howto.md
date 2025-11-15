---
Title: Add a New YAPP DSL Module
Slug: yapp-dsl-module-howto
Short: Step-by-step guide for extending the enclosure DSL via the feature-module registry and Glazed docs.
Topics:
  - yapp
  - dsl
  - developer
Commands:
  - yappctl
  - go
IsTemplate: false
IsTopLevel: true
ShowPerDefault: true
SectionType: Tutorial
Order: 45
---

## Why this guide exists

YAPP now exposes a feature-module registry (`pkg/yappgen/features.go`) so adding a DSL module is no longer a scavenger hunt across half the codebase. This tutorial walks you through the entire lifecycle—from defining the schema to lighting up `yappctl help`—and highlights the code touch points you must hit for each new module.

You should be comfortable with Go, YAML, and OpenSCAD. When in doubt, clone an existing feature (e.g., `push_buttons`) and compare it to this checklist.

## Prerequisites

- Working Go toolchain (`go version` ≥ 1.21).
- This repo checked out, with `openscad` on your PATH if you plan to render STLs.
- Familiarity with the resolver (`pkg/resolver`) and how `yappctl generate` flows from YAML → SCAD.
- Optional but recommended: read `pkg/docs/tutorials/yapp-dsl-reference.md` for user-facing terminology.

## Step 1 — Sketch the DSL schema

Before writing code, specify what the YAML should look like. Capture it in the ticket workspace (`ttmp/.../various/*.md`) and the DSL reference doc.

Questions to answer:

- Where does the module live? Most features belong under `features.<module_name>`.
- Do you need nested objects (like `cap`, `lid`, `switch`) or a flat map?
- Which fields are required vs optional? What are the defaults?
- Does the feature toggle any SCAD globals (e.g., `printSwitchExtenders`)?

Documenting the schema first keeps resolver churn minimal and provides a contract for review.

## Step 2 — Implement builder helpers

Every DSL module eventually becomes an OpenSCAD array. Builder helpers in `pkg/yappgen/map.go` convert the user-friendly maps into positional arrays.

1. **Add a schema** (optional but common): `[]ParamSpec` describing the positional parameters.
2. **Write the builder**: a function `func([]map[string]any) ([][]any, error)` that enforces required fields, inserts `Undef` defaults, and appends flags.
3. **Cover edge cases** with unit tests in `pkg/yappgen/yappgen_test.go`.

Example snippet (abridged from `buildPushButtons`):

```go
func buildPushButtons(items []map[string]any) ([][]any, error) {
    for idx, it := range items {
        label := fmt.Sprintf("push_buttons[%d]", idx)
        capMap, err := getMapField(it, "cap", label, true)
        // ...snip...
        params := []any{x, y, capLength, capWidth /* etc */}
        if coordStr != "" {
            flag, err := pushButtonCoordinateFlag(coordStr)
            if err != nil {
                return nil, errors.Wrapf(err, "%s.coordinate", label)
            }
            params = append(params, flag)
        }
        out = append(out, params)
    }
    return out, nil
}
```

Keep builder logic in `pkg/yappgen` so it can be shared by both the CLI and any future renderers.

## Step 3 — Register the module

With the registry in place, you only need to describe your feature once. Add an entry to `featureModules` in `pkg/yappgen/features.go`.

```go
var featureModules = []FeatureModule{
    // Existing modules...
    newArrayFeatureModule(
        "display_mounts",           // DSL key: features.display_mounts
        "displayMounts",            // SCAD array name
        func(m *Model) *[]map[string]any { return &m.DisplayMounts },
        buildDisplayMounts,         // builder helper you wrote
        func(m *Model, items []map[string]any) {
            m.PrintDisplayClips = len(items) > 0
        },
    ),
}
```

Under the hood the helper:

- Pulls `features.display_mounts` from the resolved YAML.
- Normalizes the array-of-maps and stores it on the `Model`.
- Runs your builder when emitting SCAD and writes the resulting array.
- Runs the optional `afterCollect` hook so you can toggle globals (like `PrintSwitchExtenders`).

Need more control (e.g., `cutouts` routing to six faces)? Implement `FeatureModule` directly—see `cutoutFeatureModule` for a template. For a full working example, browse `pkg/yappgen/modules/pushbuttons`: it exports a pure builder that the registry wires up via `newArrayFeatureModule`.

## Step 4 — Update the model and SCAD emission (automatic now!)

Thanks to the registry:

- `BuildModel` already calls `collectFeatureModules`, so your module is automatically populated.
- `EmitSCAD` calls `emitFeatureModules`, so as soon as your builder works you’ll see a new SCAD array.

You no longer touch these files when adding features—just confirm your module was registered correctly.

## Step 5 — Add examples, docs, and tests

1. **Examples**: Add or update YAML in `examples/` that exercises the new module. Mirror any existing SCAD reference if parity is required.
2. **Tests**:
   - Unit tests for the builder helper (`pkg/yappgen/yappgen_test.go`).
   - Resolver tests if you introduce new expression semantics (`pkg/resolver`).
   - Optional integration test: run `yappctl generate` in CI to ensure SCAD output contains the new array.
3. **Docs**:
   - Update `pkg/docs/tutorials/yapp-dsl-reference.md` with the new schema section.
   - If applicable, mention workflows in `pkg/docs/tutorials/yapp-dsl-getting-started.md`.
   - Log the work in the ticket workspace via `docmgr` (index + changelog).

## Step 6 — Verify with yappctl

Run both `resolve` and `generate` to smoke-test the experience:

```bash
go run ./cmd/yappctl resolve \
  --input examples/your-module.yaml \
  --out-file /tmp/your-module-resolved.yaml \
  --strict

go run ./cmd/yappctl generate \
  --input examples/your-module.yaml \
  --scad-out build/your-module.scad \
  --stl-lid build/your-module-lid.stl \
  --render-timeout 1m
```

Opening the SCAD in OpenSCAD should show the new geometry or debugging output from YAPP. If the STL render fails, revisit your builder defaults—`undef` placeholders are usually the culprit.

## Bonus — When to implement a custom module

Use `newArrayFeatureModule` for 80% of cases. Reach for a custom `FeatureModule` when:

- The feature data isn’t a simple array (e.g., needs per-face routing).
- Emission order or formatting diverges from the default `writeArrayDecl`.
- You need extra emit-time behavior (e.g., multiple SCAD arrays per module).

Implementing the interface manually still benefits from the registry: errors are wrapped with the module name, and the rest of the pipeline stays untouched.

## Checklist

- [ ] Document the schema in the ticket + DSL reference.
- [ ] Implement builder helpers + unit tests.
- [ ] Register the module via `newArrayFeatureModule` (or a custom implementation).
- [ ] Add YAML examples and run `yappctl resolve/generate`.
- [ ] Update docs/playbooks and log the change through docmgr.

Follow this loop any time you add a YAPP SCAD module to the DSL. The registry keeps additions surgical, and these steps keep onboarding consistent for future contributors (human or LLM).
