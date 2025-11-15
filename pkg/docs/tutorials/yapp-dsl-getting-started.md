---
Title: Get Started with the YAPP DSL
Slug: yapp-dsl-getting-started
Short: Walkthrough for writing your first enclosure YAML and generating SCAD/STLs with yappctl.
Topics:
  - yapp
  - dsl
  - getting-started
Commands:
  - yappctl
  - go
IsTemplate: false
IsTopLevel: true
ShowPerDefault: true
SectionType: Tutorial
Order: 30
---

## Why a YAML DSL?

YAPP started life as “copy this SCAD template and tweak dozens of variables.” The DSL replaces that brittle workflow with declarative YAML that `yappctl` can resolve, validate, and turn into real geometry. You describe the PCB, the walls, and the features you need; the tool handles the algebra, padding, and knobs.

## Prerequisites

- Go installed (so you can run `yappctl` straight from source).
- OpenSCAD on your PATH if you want automatic STL renders.
- The repo cloned locally.

## Step 1 — Establish the skeleton

Every enclosure file shares the same backbone. Start with metadata and the board dimensions; these values drive almost every downstream calculation.

```yaml
project: demo-buttons
version: 0
units: mm
yapp_version: v3.3.8

pcb:
  length: 30
  width: 40
  thickness: 1.6
  z_clearance: 3.0
```

Why it matters:

- `project`/`version` show up in logs, output filenames, and future provenance.
- `pcb.*` is the heart of the model—the generator expands the shell outward from these numbers.

### Optional: keep formulas in `vars`

Add a `vars` block near the top when you need reusable expressions or derived constants:

```yaml
vars:
  wall_clearance: tolerances.perimeter
  standoff_height: pcb.z_clearance + 4
  min_base: 2
```

Anything defined here can reference other `vars.*` entries or dotted paths like `enclosure.wall.thickness`. During resolution, those values are hoisted into the global environment, so you can later write `standoffHeight: standoff_height` without repeating the formula. It keeps magic numbers in one tidy location and mirrors the flexibility of SCAD macros.

## Step 2 — Describe the shell

Walls, base, and lid thicknesses determine how beefy the enclosure feels and how much room you have for solder joints.

```yaml
enclosure:
  wall:
    thickness: 2.4
    clearance: 1.0        # padding between PCB and inner wall
    fillet_radius: 2.0
  base:
    thickness: 1.6
  lid:
    thickness: 1.6
  ridge:
    height: 5.0
    slack: 0.2
```

The DSL uses a *clearance* value instead of separate paddings. When resolved, the generator maps it to the four `padding*` fields that the SCAD expects. Ridge settings control how the lid overlaps the base; bump them up for beefier joints.

## Step 3 — Add features

The fun part. Each subarray under `features` corresponds to a YAPP array—`pcb_stands`, `cutouts`, `connectors`, `snap_joins`, and soon `push_buttons`.

```yaml
features:
  pcb_stands:
    - x: 5
      y: 5
    - x: 25
      y: 35
  cutouts:
    - face: front
      from_back: 3
      from_left: 2
      width: 30
      height: 14
      radius: 2
      shape: rounded_rect
  connectors: []
  snap_joins: []
```

- Stands pin the PCB in place; each entry mirrors the SCAD positional parameters but with names instead of indices.
- Cutouts always include a `face` plus coordinates relative to the box. The DSL converts the `shape` strings into the `yappRectangle`/`yappCircle` flags automatically.
- Leave arrays empty if you don’t need them; the generator skips them.

## Step 4 — Set tolerances and defaults

```yaml
tolerances:
  holes: 0.4
  perimeter: 1.0
```

These values tweak global clearances for your printer. Holes get a little bigger so screws fit; the perimeter slack gives breathing room between the PCB and walls. If omitted, YAPP falls back to its built-in defaults.

## Step 5 — Resolve before you render

The resolver expands expressions, fills in defaults, and verifies that required fields exist. It also ensures that numeric math is done outside OpenSCAD, so you get faster iterations.

```bash
go run ./cmd/yappctl resolve \
  --input examples/yapp-demo-buttons.yaml \
  --out-file /tmp/demo-buttons-resolved.yaml \
  --format yaml \
  --strict
```

Use this step when you’re experimenting with formulas (`width: pcb.width + 4`), or when you want a “flattened” view for debugging.

- `--strict` warns about unused `vars` definitions and unknown keys, which is handy when refactoring.
- Complex expression chains resolve automatically; if you ever end up with circular references, the resolver tells you which paths couldn’t settle.

## Step 6 — Generate SCAD and STLs

```bash
go run ./cmd/yappctl generate \
  --input examples/yapp-demo-buttons.yaml \
  --scad-out build/demo-buttons.scad \
  --stl-base build/demo-buttons-base.stl \
  --stl-lid build/demo-buttons-lid.stl \
  --render-timeout 45s
```

Tips:

- Omit `--stl-*` to skip OpenSCAD (faster iteration).
- Adjust `--render-timeout` when printing complex lids.
- `yappctl` automatically loads the embedded tutorials, so `yappctl help tutorial yapp-dsl-reference` works offline.

## Keep iterating

1. Tweak `features.*` arrays to add stands, cutouts, or (soon) push buttons.
2. Run `yappctl resolve` to catch typos early.
3. Run `yappctl generate` and preview in OpenSCAD.
4. Repeat until you like the enclosure.

Now that you’ve seen the flow, jump over to the reference page for exhaustive tables of every field, shape enum, and optional flag: `go run ./cmd/yappctl help yapp-dsl-reference`.
