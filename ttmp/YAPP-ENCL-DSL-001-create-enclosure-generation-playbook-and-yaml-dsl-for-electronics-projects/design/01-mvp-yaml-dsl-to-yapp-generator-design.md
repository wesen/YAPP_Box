---
Title: MVP Design — YAML DSL to YAPP OpenSCAD Generator
Ticket: YAPP-ENCL-DSL-001
Status: review
Topics:
    - yapp
    - openscad
    - dsl
    - generator
DocType: design-doc
Intent: long-term
Owners:
    - manuel
RelatedFiles:
    - Path: ../analysis/04-mvp-path-forward-semantic-fixes-required.md
      Note: MVP plan and required semantic fixes
    - Path: ../reference/01-enclosure-dsl-language-reference.md
      Note: DSL spec (to be updated per analysis/04)
    - Path: ../debate/05-round-5-mvp-semantic-correctness-what-must-map-cleanly.md
      Note: Rationale for semantic decisions
    - Path: ../../../../YAPP_Template_v3.scad
      Note: Canonical parameter orders and defaults for YAPP
    - Path: ../../../../YAPPgenerator_v3.scad
      Note: Upstream generator include
ExternalSources: []
Summary: 'MVP YAML→YAPP generator: schema-driven mapping, shape-aware cutouts, defaults-based coords'
LastUpdated: 2025-11-08T19:46:32.849969578-05:00
---




# MVP Design — YAML DSL to YAPP OpenSCAD Generator

## 1) Overview and Scope

Goal: Implement a YAML-to-OpenSCAD generator that maps a resolved Enclosure DSL document to YAPP template parameters for a basic enclosure. MVP focuses on four feature types:
- pcbStands
- connectors
- cutouts
- snapJoins

Non-goals for MVP:
- Advanced features (lightTubes, pushButtons, boxMounts, labels, ridgeExt, displayMounts, polygon/ring/sphere cutout shapes)
- Per-feature coordinate overrides and escape hatches for raw SCAD
- Multi-PCB and version compatibility checks

Primary references:
- See RelatedFiles frontmatter: analysis/04 (implementation plan), debate/05 (semantic rationale), language reference, and `YAPP_Template_v3.scad`.


## 2) Inputs and Outputs

- Input: Resolved YAML DSL document (all expressions evaluated to numeric literals). We will use `pkg/resolver.Resolve` when integrating end-to-end; the generator itself expects a fully numeric model.
- Output: A single `.scad` file that:
  - Includes `YAPPgenerator_v3.scad`
  - Emits base/box/global parameters and feature arrays in the exact YAPP order
  - Ends with `YAPPgenerate();`

Naming: `PROJECT_NAME_yapp.scad` (caller decides final path).


## 3) Key Decisions (from analysis/04 and debate/05)

- No `coordinates.origin` in DSL for MVP: rely on YAPP defaults per feature type. The generator does not emit coordinate flags for MVP; feature semantics follow YAPP defaults.
- MVP scope is strictly: pcbStands, connectors, cutouts, snapJoins.
- Optional positional parameters: insert `undef` for skipped optional slots to preserve YAPP positional meaning.
- Cutouts have shape-dependent parameters: set unused parameters to 0 based on shape, and validate required shape fields.
- Breaking changes allowed during iteration.

Note: The debate considered “always emit explicit coordinate flags.” For MVP, we adopt the simpler rule from analysis/04: rely on YAPP defaults per feature (no flags). We document this and keep the door open to emit explicit flags later if drift/ambiguity becomes an issue.


## 4) Semantic Rules

- Coordinates:
  - pcbStands, connectors → interpreted in PCB coordinates (YAPP defaults)
  - cutouts, snapJoins → interpreted in Box coordinates (YAPP defaults)
- Optional parameters:
  - If the user specifies later optional params but omits earlier ones, the generator inserts `undef` at the omitted positions.
- Cutout shapes:
  - rectangle: width, length used; radius = 0
  - circle: radius used; width = 0, length = 0
  - rounded_rect: width, length, radius used
  - circle_with_flats: width (flat distance), radius, length (distance between flats) used


## 5) Parameter Schemas (positional order)

These schemas encode YAPP positional parameter order, required/optional status, and default semantics. The generator uses them to build arrays with correct `undef` placeholders.

- pcbStands (9 positional)
  1) x (required)
  2) y (required)
  3) height (optional; default = `standoffHeight`)
  4) pcb_gap (optional; default = -1)
  5) diameter (optional; default = `standoffDiameter`)
  6) pin_diameter (optional; default = `standoffPinDiameter`)
  7) hole_slack (optional; default = `standoffHoleSlack`)
  8) fillet_radius (optional; default = 0)
  9) pin_length (optional; default = 0)

- connectors (10 positional)
  1) x (required)
  2) y (required)
  3) stand_height (required)
  4) screw_d (required)
  5) screw_head_d (required)
  6) insert_d (required)
  7) outside_d (required)
  8) insert_depth (optional; default = entire connector)
  9) pcb_gap (optional; default = pcbThickness if PCB coords else 0)
  10) fillet_radius (optional; default = 0/auto)

- snapJoins (2 positional + side flag)
  1) pos (required; along wall axis)
  2) width (required)
  - side: enum flag(s) → `yappLeft|yappRight|yappFront|yappBack`

- cutouts (7 positional + shape flag)
  1) from_back (required)
  2) from_left (required)
  3) width (required for rectangle/rounded_rect/circle_with_flats/circle_with_key)
  4) length (required for rectangle/rounded_rect/circle_with_key; 0 for circle; defined for circle_with_flats)
  5) radius (required for circle/rounded_rect/circle_with_flats/circle_with_key; 0 for rectangle)
  6) shape_flag (required; maps from DSL `shape`)
  7) depth (optional; default = 0/auto)
  8) angle (optional; default = 0)


## 6) DSL Field Names and Enums (MVP)

- Faces for cutouts:
  - `front`, `back`, `left`, `right`, `top` (lid), `bottom` (base)
  - Maps to YAPP arrays: `cutoutsFront`, `cutoutsBack`, `cutoutsLeft`, `cutoutsRight`, `cutoutsLid`, `cutoutsBase`

- Snap side flag:
  - DSL: `side: left|right|front|back`
  - YAPP flags: `yappLeft|yappRight|yappFront|yappBack`

- Connector/cutout shape:
  - DSL: `rectangle`, `circle`, `rounded_rect`, `circle_with_flats`, `circle_with_key`
  - YAPP flags: `yappRectangle`, `yappCircle`, `yappRoundedRect`, `yappCircleWithFlats`, `yappCircleWithKey`

- Placement/type/corner flags for pcbStands/connectors:
  - MVP: treat as optional top-level enums under each item when needed; translate to YAPP flags after positional params.


## 7) Architecture

- Language: Go (aligns with repo; reuse `pkg/resolver`).
- CLI: `cmd/yapp-gen` (Cobra), input YAML → resolve expressions → generate `.scad`.
- Library: `pkg/yappgen/`
  - `model.go`: typed structs for resolved DSL subset (pcb, enclosure basics, features)
  - `schema.go`: parameter schemas and shape maps
  - `emit.go`: SCAD emitter (header, globals, feature arrays, footer)
  - `map.go`: mapping from DSL to YAPP positional arrays (insert `undef`)
  - `validate.go`: required fields and enum checks
- Integration:
  - Use `pkg/resolver` to produce a resolved `map[string]any`
  - Convert to typed `pkg/yappgen` model, then emit SCAD

Concurrency:
- For batch emission (many features), we can parallelize per-face cutout emission; gate with `errgroup` and context.

Errors:
- Use `github.com/pkg/errors` for wrapping.


## 8) Generation Flow (pseudocode)

```go
// cmd/yapp-gen/main.go (sketch)
func run(ctx context.Context, inPath, outPath string) error {
    raw := readYAML(inPath)
    resolved, err := resolver.Resolve(ctx, raw, resolver.Options{MaxIterations: 16, Strict: false})
    if err != nil { return errors.Wrap(err, "resolve") }

    m, err := yappgen.BuildModel(resolved)     // validate + normalize
    if err != nil { return errors.Wrap(err, "model") }

    scad, err := yappgen.EmitSCAD(ctx, m)      // returns []byte
    if err != nil { return errors.Wrap(err, "emit") }

    return os.WriteFile(outPath, scad, 0o644)
}
```

```go
// pkg/yappgen/map.go (sketch)
func buildParams(item map[string]any, schema []ParamSpec) ([]any, error) {
    out := make([]any, 0, len(schema))
    for _, p := range schema {
        v, has := item[p.Name]
        if has {
            out = append(out, v)
            continue
        }
        if p.Required {
            return nil, errors.Errorf("missing required: %s", p.Name)
        }
        out = append(out, "undef") // preserve positional order for optional
    }
    return out, nil
}
```

```go
// pkg/yappgen/emit.go (sketch)
func EmitSCAD(ctx context.Context, m *Model) ([]byte, error) {
    var b strings.Builder
    b.WriteString("include <./YAPPgenerator_v3.scad>\n\n")
    // Emit global dims (pcb, wall, base/lid), respecting template names
    // Emit pcbStands = [...]; connectors = [...];
    // Emit cutoutsXxx = [...]; snapJoins = [...];
    b.WriteString("YAPPgenerate();\n")
    return []byte(b.String()), nil
}
```


## 9) Emission Details

- File preamble:
  - `include <./YAPPgenerator_v3.scad>`
  - Assign base parameters: `pcbLength`, `pcbWidth`, `pcbThickness`, paddings, wall/base/lid thicknesses, heights, fillets, etc. Use names from `YAPP_Template_v3.scad`.
- Features:
  - `pcbStands = [ [x, y, height|undef, pcb_gap|undef, diameter|undef, ...], ... ];`
  - `connectors = [ [x, y, stand_height, screw_d, screw_head_d, insert_d, outside_d, insert_depth|undef, pcb_gap|undef, fillet_radius|undef], ... ];`
  - `cutouts{Front|Back|Left|Right|Lid|Base} = [ [from_back, from_left, width|0, length|0, radius|0, shapeFlag, depth|undef, angle|undef], ... ];`
  - `snapJoins = [ [pos, width, sideFlag], ... ];`
- Footer: `YAPPgenerate();`

Flags:
- MVP: Do not emit coordinate flags; rely on YAPP defaults (documented).
- Shape flags: emit the correct `yappXXX` symbol based on DSL `shape`.
- Side flags for snapJoins: map 1:1 from DSL.


## 10) Validation

- Required fields present per schema.
- Enum values valid (`face`, `shape`, `side`).
- Cutout shape params: for unused params, emit 0; required shape params must be present.
- Numeric sanity: all lengths > 0; wall/base/lid thickness minimums (reuse/reference rules from reference doc).


## 11) Testing Plan

- Unit tests for:
  - Schema mapping with `undef` insertion (pcbStands/connectors)
  - Cutout shape param handling (zeros for unused)
  - Per-face cutout routing (front/back/left/right/top/bottom)
  - Snap side flags mapping
- Golden tests: small YAML inputs → generated `.scad` compared to expected strings.
- Manual OpenSCAD validation:
  - Compile generated `.scad` for combinations covering the four MVP features
  - Visual check for coordinate semantics (PCB-relative vs Box-relative) by modifying paddings


## 12) Risks and Mitigations

- Positional order drift in YAPP upstream:
  - Mitigation: Centralize schemas in one file; add comment references to `YAPP_Template_v3.scad` parameter docs; add unit tests.
- Coordinate semantics confusion:
  - Mitigation: Keep defaults-only approach for MVP; document prominently; consider explicit flags in a follow-up if needed.
- Shape coverage creep:
  - Mitigation: Keep to the enumerated shapes; defer others post-MVP.


## 13) Implementation Checklist (MVP)

- [ ] Update DSL spec: remove `coordinates` section; document YAPP default semantics
- [ ] Define DSL enums and field names for the four features
- [ ] Build parameter schemas (pcbStands, connectors, snapJoins, cutouts)
- [ ] Implement schema-based positional array builder (with `undef`)
- [ ] Implement shape-specific mapping for cutouts (set unused to 0)
- [ ] Emit SCAD: header, globals, feature arrays, footer
- [ ] Validation: required params and enums
- [ ] Unit tests and goldens
- [ ] Generate sample `.scad` files for visual validation


## 14) Acceptance Criteria

- Generated SCAD compiles without syntax errors.
- Visual behavior matches YAPP semantics when paddings change:
  - PCB-relative features move with the PCB
  - Box-relative features stay anchored to box faces/edges
- Optional parameter skipping produces correct `undef` placeholders.
- Cutout shape semantics are correct (unused params → 0).


## 15) Notes and References

- See `YAPP_Template_v3.scad` for definitive parameter orders and defaults.
- The generator intentionally avoids non-MVP YAPP features to minimize initial complexity.


