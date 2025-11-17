---
Title: 2025-11-17 ValidateStructure requirements analysis
Ticket: YAPP-DSL-GAPS-001
Status: active
Topics:
    - yapp
    - dsl
    - analysis
    - features
DocType: analysis
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Phase-1 (pre-resolution) structure validation requirements for all modules; per-field type/required checks and boundaries between structure vs constraints.
LastUpdated: 2025-11-16T21:04:43.137468196-05:00
---

# 2025-11-17 ValidateStructure requirements analysis

## Purpose and scope

This document defines what ValidateStructure must check for all feature modules, consistently, during Phase‑1 (pre‑resolution) validation. It clarifies the boundary between structure checks (this phase) and constraint checks (Phase‑2, post‑resolution), and lists the per‑module requirements derived from each module’s `schema.yaml` and the DSL reference.

References:
- Registry interface: `/home/manuel/code/others/YAPP_Box/pkg/registry/schema.go`
- Resolver flow: `/home/manuel/code/others/YAPP_Box/pkg/resolver/resolver.go` (calls `validateStructure` before expression resolution)
- DSL reference: `/home/manuel/code/others/YAPP_Box/pkg/docs/tutorials/yapp-dsl-reference.md`
- Module schemas:
  - `/home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/pcbstands/schema.yaml`
  - `/home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/connectors/schema.yaml`
  - `/home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/snapjoins/schema.yaml`
  - `/home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/cutouts/schema.yaml`
  - `/home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/lighttubes/schema.yaml`
  - `/home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/boxmounts/schema.yaml`
  - `/home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/pushbuttons/schema.yaml`

## Phase‑1 vs Phase‑2 responsibilities

- ValidateStructure (Phase‑1, pre‑resolution):
  - Ensure the feature value is an array (where applicable)
  - Ensure each entry is an object/map
  - Ensure required fields are present
  - Ensure field shapes and primitive kinds are correct:
    - number fields: allow numeric literals or string expressions (strings get resolved later)
    - string fields: must be strings
    - bool fields: must be booleans (do not allow string expressions)
    - object fields: must be maps; recursively check required child fields and their kinds
    - array fields: must be arrays/slices with item kind checks
  - Ignore unknown fields for forward compatibility (strict unknown‑key checks are handled elsewhere)
  - Do NOT enforce cardinality or cross‑field rules in Phase‑1 (those belong to Phase‑2)

- ValidateConstraints (Phase‑2, post‑resolution):
  - Enum membership and synonym normalization (e.g., hyphen ↔ underscore)
  - Cross‑field rules that depend on concrete values or choices (e.g., “if shape=polygon, require polygon preset” for cutouts; “for side faces, require vertical coordinate”)
  - Min/max numeric ranges and interdependencies

Rationale: numeric expressions are strings until resolved, so Phase‑1 must not enforce numeric ranges; it only guards structure and kinds.

## Container‑level expectations (all modules)

- If `features.<module>` exists, it must be an array (`[]any`).
- Each array element must be a map (`map[string]any`).
- Do not error if the feature is absent (treat as “no items”).

## Type policy (applies everywhere)

- number fields: accept `int|int64|float32|float64` or `string` (expression or numeric string). Reject arrays/objects/bools.
- string fields (including enum‑like): accept `string` only. Reject numbers/bools/arrays/objects.
- bool fields: accept `bool` only. Reject strings/numbers/arrays/objects.
- object fields: accept `map[string]any`. Reject arrays/bools/strings/numbers.
- array fields: accept `[]any`; validate item kinds per schema.

## Per‑module structure requirements

Below, “Req” means required presence; “Kind” is the expected primitive kind at Phase‑1. Nested object required fields are recursively listed.

### pcb_stands
Source: `pkg/yappgen/modules/pcbstands/schema.yaml`
- Req: `x:number|string`, `y:number|string`
- Optional number|string: `height`, `pcb_gap`, `diameter`, `pin_diameter`, `hole_slack`, `fillet_radius`, `pin_length`
- Optional array `corners`: items are `string`
- Optional string: `shell_part`, `treatment`, `corner`, `coordinate`, `pcb_name`
- Optional bool: `no_fillet`, `self_threading`

Notes:
- Enum membership for `corner/shell_part/treatment/coordinate` is Phase‑2.

### connectors
Source: `pkg/yappgen/modules/connectors/schema.yaml`
- Req number|string: `x`, `y`, `stand_height`, `screw_d`, `screw_head_d`, `insert_d`, `outside_d`
- Optional number|string: `insert_depth`, `pcb_gap`, `fillet_radius`
- Optional string: `corner`, `coordinate`, `pcb_name`
- Optional bool: `no_fillet`, `countersink`, `through_lid`, `self_threading`, `no_internal_fillet`

Notes:
- Enum for `corner/coordinate` is Phase‑2.

### snap_joins
Source: `pkg/yappgen/modules/snapjoins/schema.yaml`
- Req number|string: `pos`, `width`
- Req string: `side`
- Optional string: `alignment`
- Optional bool: `symmetric`, `diamond`

Notes:
- Enum for `side/alignment` is Phase‑2.

### cutouts
Source: `pkg/yappgen/modules/cutouts/schema.yaml`
- Req string: `face`, `shape`
- Req number|string: `from_face_left`, `width`, `length`, `radius`
- Optional number|string: `from_face_bottom`, `from_face_back`, `depth`, `angle`
- Optional string: `polygon`
- Optional object `mask`:
  - Optional string: `preset` (Phase‑2 requires presence if `mask` given)
  - Optional number|string: `offset_x`, `offset_y`, `rotation`
- Optional string: `coordinate`, `origin`

Notes:
- Face‑dependent presence (“side faces require `from_face_bottom`, base/lid require `from_face_back`”) is Phase‑2.
- Enum membership for `face/shape/polygon/mask.preset/coordinate/origin` is Phase‑2.

### light_tubes
Source: `pkg/yappgen/modules/lighttubes/schema.yaml`
- Req number|string: `x`, `y`, `tube_length`, `tube_width`, `tube_wall`, `gap_above_pcb`
- Req string: `shape`
- Optional number|string: `lens_thickness`, `height`, `fillet_radius`
- Optional string: `coordinate`, `origin`, `pcb_name`
- Optional bool: `no_fillet`

Notes:
- Enum for `shape/coordinate/origin` is Phase‑2.

### box_mounts
Source: `pkg/yappgen/modules/boxmounts/schema.yaml`
- Req number|string: `pos`, `screw_d`, `slot_width`, `height`
- Optional number|string: `offset`, `fillet_radius`
- Optional string: `shell_part`, `alignment`, `origin`
- Optional bool: `no_fillet`
- Req object `faces`:
  - Optional bools: `left`, `right`, `front`, `back`
  - Phase‑1: only validate `faces` is an object and its children (if present) are bools
  - Phase‑2: enforce cardinality (e.g., at least one of `left/right/front/back` is true)

Notes:
- Enums for `shell_part/alignment/origin` are Phase‑2.

### push_buttons
Source: `pkg/yappgen/modules/pushbuttons/schema.yaml`
- Req number|string: `x`, `y`
- Optional string: `name`, `shape`, `polygon`, `polygon_preset`, `shape_preset`, `coordinate`, `origin`, `pcb_name`
- Optional number|string: `angle`, `fillet_radius`
- Optional bool: `no_fillet`
- Req object `cap`:
  - Req number|string: `length`, `width`, `radius`
- Req object `lid`:
  - Req number|string: `protrusion`
  - Optional number|string: `wall`, `plate_thickness`, `slack`, `snap_slack`
  - Optional bool: `no_fillet`
- Req object `switch`:
  - Req number|string: `height`, `travel`, `pole_diameter`
  - Optional number|string: `top_height`

Notes:
- Enum membership and polygon preset rules are Phase‑2.

## Implementation guidance

1) Add a shared helper package (or local helpers per module) to check kinds:
   - `isNumberLike(v any) bool` → `int|int64|float32|float64|string`
   - `isBool(v any) bool`, `isString(v any) bool`, `isMap(v any) bool`, `isArray(v any) bool`
   - `require(path, item, key)` to assert presence
   - `asMap(v)` / `asArray(v)` casting helpers

2) In each module’s `ValidateStructure(path string, data any)`:
   - If `data` is not an array, return nil (module absent) or type error if present but wrong
   - For each `[i]`:
     - Assert map
     - Check required fields (presence)
     - Check kinds per bullets above
     - For nested objects/arrays, recurse checks
     - Special case for `box_mounts.faces`: ensure at least one of `left|right|front|back` is `true`

3) Do not check enum values, numeric ranges, or face‑dependent requirements here; leave to Phase‑2.

4) Error messages:
   - Prefix with `features.<module>[i].<field>` for consistency with resolver
   - Use clear phrases like “missing required field”, “expected number|string”, “expected bool”, “expected object”

## Acceptance criteria

- Missing required fields produce Phase‑1 errors with precise paths
- Wrong kinds (e.g., array where number expected) are rejected in Phase‑1
- Features absent are silently accepted (treated as zero items)
- Enum and cross‑field rules handled in Phase‑2 continue to work (no duplication)

## Files referenced (link via docmgr)

- `/home/manuel/code/others/YAPP_Box/pkg/registry/schema.go`
- `/home/manuel/code/others/YAPP_Box/pkg/resolver/validation.go`
- `/home/manuel/code/others/YAPP_Box/pkg/resolver/resolver.go`
- `/home/manuel/code/others/YAPP_Box/pkg/docs/tutorials/yapp-dsl-reference.md`
- Module schemas:
  - `/home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/pcbstands/schema.yaml`
  - `/home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/connectors/schema.yaml`
  - `/home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/snapjoins/schema.yaml`
  - `/home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/cutouts/schema.yaml`
  - `/home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/lighttubes/schema.yaml`
  - `/home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/boxmounts/schema.yaml`
  - `/home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/pushbuttons/schema.yaml`

## Next steps

- Implement `ValidateStructure` for each module following this spec
- Optionally add a small shared validator utility under `pkg/registry` or alongside modules
- Add unit tests per module using the embedded `tests:` from each `schema.yaml` (schemagen) or write dedicated tests exercising structure errors
- Update the ticket `tasks.md` when each module’s structure validation is complete

## Option 2: Generate ValidateStructure with schemagen (detailed plan)

### Goals

- Auto‑generate `ValidateStructure(path string, data any) error` from each module’s `schema.yaml`.
- Keep module authors focused on schema and builder logic; reduce drift between schema and validators.
- Generate methods onto the existing `moduleSchema` type without editing `registry.go`.

### What we will generate

- A new file per module: `pkg/yappgen/modules/<module>/schema_validate.go` containing:
  - `func (s *moduleSchema) ValidateStructure(path string, data any) error`
  - Private helpers (per package) for kind checks (number‑like, string, bool, object, array)
  - For nested object fields, small helper validators (e.g., `validate<StructName><Child>()`)

### Data model parsed from schema.yaml

- For each field we already capture:
  - `type`: number|string|bool|object|array
  - `required`: true|false
  - `fields` (for object children)
  - `items.type` (for array item kind)
  - Optional enums/ranges exist but are Phase‑2 (constraints), not enforced here.
  
No hints: the generator derives Phase‑1 validation solely from schema types and `required` flags. Cross‑field/cardinality rules are deferred to Phase‑2.

### Generated code shape (per module)

- Pseudocode for the generated validator:

```go
func (s *moduleSchema) ValidateStructure(path string, data any) error {
    // features.<module> must be an array if present
    arr, ok := data.([]any)
    if !ok {
        // Allow absent or wrong top-type only if data is nil; otherwise type error
        if data == nil {
            return nil
        }
        return errors.Errorf("%s: expected array", path)
    }
    for i, raw := range arr {
        item, ok := raw.(map[string]any)
        if !ok {
            return errors.Errorf("%s[%d]: expected object", path, i)
        }
        // For each field in schema:
        // - if required: check presence
        // - check kind per type policy
        // - for object: recurse into children
        // - for array: check slice and item kinds
        // - apply hints (e.g., any_true on boolean children)
        if err := validate<Module>Item(path, i, item); err != nil {
            return err
        }
    }
    return nil
}
```

- Example (partial) for `cutouts` field checks:

```go
func validateCutoutsItem(path string, idx int, m map[string]any) error {
    fp := func(field string) string { return fmt.Sprintf("%s[%d].%s", path, idx, field) }
    // required strings
    if _, ok := m["face"]; !ok {
        return errors.Errorf("%s: missing required field", fp("face"))
    }
    if !isString(m["face"]) { return errors.Errorf("%s: expected string", fp("face")) }
    if _, ok := m["shape"]; !ok {
        return errors.Errorf("%s: missing required field", fp("shape"))
    }
    if !isString(m["shape"]) { return errors.Errorf("%s: expected string", fp("shape")) }
    // required number|string
    if err := expectNumberLike(m, "from_face_left", fp("from_face_left")); err != nil { return err }
    if err := expectNumberLike(m, "width",          fp("width"));          err != nil { return err }
    if err := expectNumberLike(m, "length",         fp("length"));         err != nil { return err }
    if err := expectNumberLike(m, "radius",         fp("radius"));         err != nil { return err }
    // optional number|string
    if has(m, "from_face_bottom") { if err := expectNumberLike(m, "from_face_bottom", fp("from_face_bottom")); err != nil { return err } }
    if has(m, "from_face_back")   { if err := expectNumberLike(m, "from_face_back",   fp("from_face_back"));   err != nil { return err } }
    if has(m, "depth")            { if err := expectNumberLike(m, "depth",            fp("depth"));            err != nil { return err } }
    if has(m, "angle")            { if err := expectNumberLike(m, "angle",            fp("angle"));            err != nil { return err } }
    // optional string
    if has(m, "polygon")    && !isString(m["polygon"])    { return errors.Errorf("%s: expected string", fp("polygon")) }
    if has(m, "coordinate") && !isString(m["coordinate"]) { return errors.Errorf("%s: expected string", fp("coordinate")) }
    if has(m, "origin")     && !isString(m["origin"])     { return errors.Errorf("%s: expected string", fp("origin")) }
    // optional object mask
    if has(m, "mask") {
        mm, ok := m["mask"].(map[string]any)
        if !ok { return errors.Errorf("%s: expected object", fp("mask")) }
        if has(mm, "preset")   && !isString(mm["preset"])   { return errors.Errorf("%s: expected string", fp("mask.preset")) }
        if has(mm, "offset_x") && !isNumberLike(mm["offset_x"]) { return errors.Errorf("%s: expected number|string", fp("mask.offset_x")) }
        if has(mm, "offset_y") && !isNumberLike(mm["offset_y"]) { return errors.Errorf("%s: expected number|string", fp("mask.offset_y")) }
        if has(mm, "rotation") && !isNumberLike(mm["rotation"]) { return errors.Errorf("%s: expected number|string", fp("mask.rotation")) }
    }
    return nil
}
```

### Template work

- Add new template: `pkg/schemagen/templates/schema_validate.go.tmpl`
  - Inputs: parsed schema (fields tree), module name, package name
  - Output: methods + helpers as above
- Do NOT modify `registry.go`; define the method with receiver `(s *moduleSchema)` in the generated file.

Suggested template outline:

```gotemplate
// Code generated by schemagen; DO NOT EDIT.
package {{.Package}}

import (
    "fmt"
    "github.com/pkg/errors"
)

func (s *moduleSchema) ValidateStructure(path string, data any) error {
    arr, ok := data.([]any)
    if !ok {
        if data == nil { return nil }
        return errors.Errorf("%s: expected array", path)
    }
    for i, raw := range arr {
        item, ok := raw.(map[string]any)
        if !ok { return errors.Errorf("%s[%d]: expected object", path, i) }
        if err := validate{{.RootStruct}}Item(path, i, item); err != nil { return err }
    }
    return nil
}

// One validator function per object struct
{{range .ObjectValidators}}
func validate{{.Name}}(path string, idx int, m map[string]any) error {
    // generated field checks …
    return nil
}
{{end}}

// Shared helpers (generated once per package)
func has(m map[string]any, k string) bool { _, ok := m[k]; return ok }
func isString(v any) bool { _, ok := v.(string); return ok }
func isBool(v any) bool   { _, ok := v.(bool); return ok }
func isNumberLike(v any) bool {
    switch v.(type) {
    case int, int64, float32, float64, string:
        return true
    default:
        return false
    }
}
func expectNumberLike(m map[string]any, k, fp string) error {
    v, ok := m[k]
    if !ok { return errors.Errorf("%s: missing required field", fp) }
    if !isNumberLike(v) { return errors.Errorf("%s: expected number|string", fp) }
    return nil
}
```

### Integration steps

1) Extend schemagen:
   - Load schema YAML (already done) into an AST that preserves type/required/children/items and optional `hints`.
   - Add template execution for `schema_validate.go.tmpl` during `discover`.
2) Generate for all modules:
   - `go run ./cmd/schemagen discover`
3) Remove TODO stubs in each `registry.go` (no code change needed if method defined in generated file).
4) Run `go test ./...` and add module tests that call `ValidateStructure` with invalid shapes/kinds.

### Testing strategy

- Unit tests per module (generated optional):
  - Valid “minimal” examples from `tests:` must pass structure
  - Invalid structure cases:
    - Missing required field → “missing required field”
    - Wrong kind (string where bool expected) → “expected bool”
    - Non‑object for nested object → “expected object”
- Cross‑field/cardinality (e.g., `box_mounts.faces` “at least one true”) covered by Phase‑2 constraint tests
- Resolver integration tests:
  - Feed YAML documents; ensure `validateStructure` errors surface with `features.<module>[i].<field>` paths

### Schema extensions

None required for Phase‑1. Keep schema as‑is; handle enums, cross‑field and cardinality in Phase‑2.

### Risks and mitigations

- Over‑eager validation: Ensure we treat number fields as “number|string” to allow expressions.
- Enum duplication: Keep enums for Phase‑2 only to avoid double maintenance.
- Unknown keys: Default to allowing unknown keys (forward‑compatible). Consider a future strict mode.

### Milestones

- M1: Template + generator update; compile validators for all modules
- M2: Module unit tests for structure errors
- M3: Documentation updates in `yapp-module-authoring-guide` to describe generated validators and Phase‑1 vs Phase‑2 split
