---
Title: YAPP DSL Reference
Slug: yapp-dsl-reference
Short: Field-by-field reference for enclosure YAML files consumed by yappctl/yappgen.
Topics:
  - yapp
  - dsl
  - reference
Commands:
  - yappctl
  - go
IsTemplate: false
IsTopLevel: true
ShowPerDefault: true
SectionType: Tutorial
Order: 31
---

## Purpose

This document is the authoritative field-by-field reference for the YAML DSL consumed by `yappctl` and `yappgen`. The DSL transforms declarative YAML descriptions into parametric OpenSCAD code that generates custom enclosures for your electronics projects. Whether you're defining PCB standoffs, carving cutouts for USB ports, or fine-tuning wall thickness, this reference explains every available field, its type, default behavior, and how it flows through the resolution pipeline. When you're unsure about a parameter name, valid shape enum, or expression syntax, start here. For a friendlier narrative walkthrough with complete worked examples, pair this with the getting started tutorial:

```
go run ./cmd/yappctl help yapp-dsl-getting-started
```

## Document structure at a glance

Every YAML enclosure definition follows a consistent top-level structure that separates project metadata, dimensional specifications, and feature arrays. Understanding this schema helps you locate the right place to define each aspect of your box. The table below summarizes the mandatory and optional root keys; detailed field-level breakdowns follow in subsequent sections.

| Key | Required | Description |
|-----|----------|-------------|
| `project` | ✓ | Human-friendly name; shows up in logs and generated filenames. |
| `version` | ✓ | Integer revision number for tracking design iterations. |
| `units` | ✓ | Measurement system; currently only `mm` is supported. |
| `yapp_version` | ✓ | Matches the target SCAD generator version (e.g., `v3.3.8`, `v3.1`). |
| `vars` | optional | Dictionary of reusable constants and computed expressions (see Variables section). |
| `pcb` | ✓ | PCB bounding box, default standoff parameters, and clearances. |
| `enclosure` | ✓ | Structural settings: wall/base/lid thickness, ridge geometry, padding. |
| `features` | optional | Named arrays for `pcb_stands`, `cutouts`, `connectors`, `snap_joins`, and future `push_buttons`. |
| `tolerances` | optional | Printer compensation values for hole diameters and perimeter expansion. |
| `notes` or `metadata` | optional | Free-form documentation; completely ignored by the code generator. |

**Expression support:** Almost every numeric field can be a literal (`2.4`) or an expression referencing other parts of the document (`pcb.width + 6`, `max(vars.min_base, 1.6)`). The resolver evaluates all expressions to final numeric values before geometry generation begins, allowing you to keep magic numbers in one place and maintain consistency across related dimensions.

## Variables and expressions

- The optional `vars` block is a dictionary evaluated just like the rest of the document. Each entry can reference other `vars.*` values, top-level fields (`pcb.length`), or functions such as `min`, `max`, `round`, `floor`, `ceil`, and `clamp`.
- During resolution, the engine hoists `vars.foo` into the global environment, so you can reference `foo` directly elsewhere (`standoff_height: vars.base_clearance + 2` or simply `standoff_height: base_clearance` once defined).
- Expressions use dot-notation paths to reference other fields. Example: `enclosure.base.thickness: max(vars.min_base, enclosure.wall.thickness)`.
- Resolution runs until a fixed point (up to 16 passes by default). Chains such as `a = 1`, `b = a + 1`, `c = b + 1` are supported; cycles will error after the max iteration count.
- `yappctl resolve --strict` (or `generate --strict`) enables strict mode: unknown keys trigger errors, and unused `vars` entries are flagged so dead code doesn’t accumulate.

This mechanism is ideal for keeping “magic numbers” in one place, computing offsets, or mirroring relationships from CAD drawings without touching OpenSCAD.

## PCB section

The PCB block defines the bounding box the shell wraps around.

| Field | Required | Default | Notes |
|-------|----------|---------|-------|
| `length`, `width`, `thickness` | ✓ | — | Primary footprint measurements. |
| `z_clearance` | optional | `1.0` | Maps to `standoffHeight`; set higher for tall solder joints. |
| `standoffs.diameter` | optional | YAPP default (6.0) | Global fallback for stand features. |
| `standoffs.screw_d` | optional | 2.4 | Used for pin/hole diameters. |
| `standoffs.hole_slack` | optional | inherits `tolerances.holes` | Override if a specific build material needs a different slack. |

## Enclosure section

Walls, base, lid, and ridge geometry live here.

```yaml
enclosure:
  wall:
    thickness: 2.4
    clearance: 1.0
    fillet_radius: 2.0
  base:
    thickness: 1.6
  lid:
    thickness: 1.6
  ridge:
    height: 5.0
    slack: 0.2
    gap: 0.5
```

- `wall.clearance` is applied to every side. If you need asymmetric padding today, compute it into custom fields and wire them in SCAD; a native DSL option is on the roadmap.
- `ridge.height` should be at least `1.8 * wall thickness` for sturdy overlaps, mirroring SCAD recommendations.
- `ridge.gap` fine-tunes how snug the lid rides over the ridge; bigger gaps ease assembly but can introduce wobble.

## Features section

Every array under `features` corresponds to a YAPP data structure. Items are maps with named fields, which the generator later converts into positional SCAD parameters.

### `pcb_stands`

| YAML field | Required | Description |
|------------|----------|-------------|
| `x`, `y` | ✓ | Measured from the PCB origin (`pcb[0,0,0]`). |
| `height` | optional | Overrides `z_clearance` per stand. |
| `pcb_gap` | optional | Additional offset (useful for sandwiching the board between base and lid). |
| `diameter`, `pin_diameter`, `hole_slack`, `fillet_radius` | optional | Per-stand overrides for fine tuning. |

### `cutouts`

| Field | Required | Description |
|-------|----------|-------------|
| `face` | ✓ | `front`, `back`, `left`, `right`, `base`, or `lid`. |
| `from_back`, `from_left` | ✓ | Distances along the face. |
| `width`, `height`, `radius` | ✓ | Always provide values; the emitter zeroes unused dimensions depending on shape. |
| `shape` | ✓ | `rectangle`, `circle`, `rounded_rect`, `circle_with_flats`, `circle_with_key`. |
| `depth` | optional | Defaults to plane thickness. |
| `angle` | optional | Degrees. |
| `mask` / `polygon` | optional | Reserved for future presets (hex grids, arrow cutouts, etc.). |

Coordinate/origin flags (e.g., `coordinate: box`) will align with YAPP’s `yappCoord*` flags once exposed.

### `connectors`

| Field | Required | Description |
|-------|----------|-------------|
| `x`, `y` | ✓ | PCB-relative positions. |
| `stand_height` | ✓ | Height from origin. |
| `screw_d`, `screw_head_d` | ✓ | Hardware sizing. |
| `insert_d`, `outside_d` | ✓ | Sleeve dimensions. |
| `insert_depth`, `pcb_gap`, `fillet_radius` | optional | Tuning for deeper inserts or rounded transitions. |
| `corner` | optional | `front_left`, `back_right`, etc. Mirrors SCAD corner selection flags. |
| `coordinate`, `origin`, `countersink` | optional | Flag-style options slated for DSL exposure soon. |

### `snap_joins`

```yaml
snap_joins:
  - pos: 15          # distance along the chosen side
    width: 18
    side: front      # enum: front/back/left/right
```

The DSL validates `side` and appends the correct `yappFront`, `yappBack`, etc. tokens when building SCAD arrays.

### `push_buttons` (planned)

The schema is under construction in ticket `YAPP-PUSH-BUTTONS-001`. Expect fields that map 1:1 with the SCAD positional parameters:

- `x`, `y`
- `cap.length`, `cap.width`, `cap.radius`
- `lid.protrusion`, `lid.wall`, `lid.plate_thickness`, `lid.slack`, `lid.snap_slack`
- `switch.height`, `switch.travel`, `switch.pole_diameter`, `switch.top_offset`
- `shape`, `angle`, `fillet_radius`
- `coordinate`, `origin`, `polygon` preset (e.g., `shapeArrow`)

Once implemented, the builder will automatically set `printSwitchExtenders = true` when any entry exists.

## Tolerances section

| Field | Default | Description |
|-------|---------|-------------|
| `holes` | 0.4 | Added to all hole diameters; larger values loosen fits. |
| `perimeter` | 1.0 | Added to shell dimensions to keep the PCB from binding. |

Future keys (e.g., `buttons`, `snap_joins`) will let you target specific features without touching global settings.

## The resolver/generator pipeline

1. **Resolve** – `yappctl resolve` reads the YAML, evaluates expressions, fills defaults, and outputs a fully numeric doc. This step catches missing required fields early.
2. **Build model** – `pkg/yappgen.BuildModel` extracts data into a typed structure (`Model`) and translates DSL arrays into normalized maps.
3. **Emit SCAD** – `yappgen.EmitSCAD` plugs the model into `YAPPgenerator_v3.scad`, rewriting include paths as needed.
4. **Render STLs** (optional) – `yappctl generate --stl-*` invokes OpenSCAD with the right `print*` flags.

Understanding where each field flows helps when debugging: if a cutout lands on the wrong face, inspect the resolved YAML first, then the generated SCAD text to confirm whether the issue is in the DSL, mapper, or OpenSCAD.

## Helpful commands

- `go run ./cmd/yappctl help yapp-dsl-getting-started` — narrative tutorial to complement this reference.
- `go run ./cmd/yappctl resolve --input your.yaml --format json` — inspect the resolved document.
- `go run ./cmd/yappctl generate --input your.yaml --scad-out out.scad` — emit geometry; add `--stl-*` for meshes.

## Additional resources

- `examples/*.yaml` mirrors the classic SCAD demos—great starting points.
- Ticket `ttmp/YAPP-ENCL-DSL-001-*` captures the original design debates and rationale.
