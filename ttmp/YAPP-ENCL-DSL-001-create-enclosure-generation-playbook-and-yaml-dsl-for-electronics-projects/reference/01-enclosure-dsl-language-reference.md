---
Title: 'Enclosure DSL: Language Reference'
Ticket: YAPP-ENCL-DSL-001
Status: active
Topics:
    - playbook
    - electronics
    - 3d-printing
    - yapp
    - openscad
    - automation
DocType: reference
Intent: long-term
Owners:
    - manuel
RelatedFiles:
    - Path: ../analysis/04-mvp-path-forward-semantic-fixes-required.md
      Note: MVP semantic decisions (no coordinates block, schema-based params)
    - Path: ../design/01-mvp-yaml-dsl-to-yapp-generator-design.md
      Note: Generator architecture that consumes this DSL
ExternalSources: []
Summary: "Authoritative reference for the MVP enclosure DSL (metadata, PCB, enclosure, tolerances, and the four supported feature groups). Coordinates block removed; feature semantics rely on YAPP defaults."
LastUpdated: 2025-11-08T16:18:37.598375185-05:00
---


# Enclosure DSL: Language Reference

## Goal

Define a concise, complete reference for the Enclosure DSL used to describe electronics enclosures, including structure, types, expressions, and constraints. This document is authoritative for writers and implementers.

## Context

This DSL is consumed by a resolver that evaluates variables and expressions to produce a fully numeric configuration before mapping to YAPPgenerator/OpenSCAD. Expressions are plain scalars using dotted paths (no `=`, no `$`), and are resolved to a fixed point.

## Structure Overview

The MVP DSL is organized into six top-level sections:

1. **Metadata** — `project`, `version`, `units`, `yapp_version`
2. **Variables** — `vars` (reusable computed values)
3. **PCB** — `pcb` (dimensions, z-clearance, optional hole metadata)
4. **Enclosure** — `enclosure` (wall/base/lid thicknesses, fillets)
5. **Features** — `features` (pcb_stands, connectors, cutouts, snap_joins)
6. **Tolerances** — `tolerances` (print fit adjustments)

## Top-Level Keys

### Metadata

```yaml
project: string          # Project identifier (e.g., "pico-temp-001")
version: integer         # DSL schema version (currently 0)
units: mm                # Always "mm" (millimeters)
yapp_version: string     # Target YAPPgenerator version (e.g., "v3.3.8")
```

### Variables

```yaml
vars:
  <name>: <number|expression>
```

Define reusable values that can reference other fields. Variables are resolved before final output.

**Example:**
```yaml
vars:
  wall_clearance: tolerances.perimeter
  standoff_height: pcb.z_clearance + 4.0
  min_base_thickness: 2.0
```

### PCB

```yaml
pcb:
  length: number           # PCB length (X dimension)
  width: number            # PCB width (Y dimension)
  thickness: number        # PCB thickness (typically 1.6mm)
  z_clearance: number      # Space between PCB bottom and enclosure base
  holes:                   # Mounting holes (optional)
    - x: number
      y: number
      d: number            # Hole diameter
  standoffs:               # PCB standoffs/supports
    type: round|hex|none   # Standoff shape
    diameter: number       # Standoff outer diameter
    height: number|expr    # Standoff height (PCB bottom to base)
    screw_d: number        # Screw hole diameter
```

**Notes:**
- `holes` coordinates are relative to PCB origin (typically bottom-left corner)
- `standoffs.height` often references `pcb.z_clearance` plus additional height
- Use `type: none` to disable standoffs

### Enclosure

```yaml
enclosure:
  wall:
    thickness: number           # Wall thickness
    clearance: number|expr      # Internal clearance from PCB outline
    fillet_radius: number       # Corner fillet radius
  base:
    thickness: number|expr      # Base plate thickness
  lid:
    type: screws|snap           # Lid attachment method
    screws:                     # Only if type=screws
      count: integer            # Number of screws
      positions: corners|custom # Screw placement strategy
      head_clearance: number    # Countersink depth for screw heads
```

**Notes:**
- `wall.clearance` defines space between PCB edge and inner wall
- `base.thickness` can use expressions like `max(vars.min_base, wall.thickness)`
- `lid.screws` is only required when `lid.type: screws`

### Features

MVP supports four feature groups that map cleanly to YAPP arrays. Each item is a map with named fields; omitted optional fields use generator defaults.

```yaml
features:
  pcb_stands:
    - x: number                 # Required (PCB X coordinate)
      y: number                 # Required (PCB Y coordinate)
      height: number            # Optional; default uses standoffHeight
      pcb_gap: number           # Optional; default -1 (auto)
      diameter: number          # Optional; default standoffDiameter
      pin_diameter: number      # Optional; default standoffPinDiameter
      hole_slack: number        # Optional; default tolerances.holes
      fillet_radius: number     # Optional (0 = disable)
      pin_length: number        # Optional (0 = disable)

  connectors:
    - x: number                 # Required (PCB X coordinate)
      y: number                 # Required (PCB Y coordinate)
      stand_height: number      # Required
      screw_d: number           # Required (shaft diameter)
      screw_head_d: number      # Required
      insert_d: number          # Required (heat-set insert diameter)
      outside_d: number         # Required (boss outer diameter)
      insert_depth: number      # Optional
      pcb_gap: number           # Optional; default depends on PCB thickness
      fillet_radius: number     # Optional

  cutouts:
    - face: front|back|left|right|top|bottom   # Required YAPP face
      shape: rectangle|circle|rounded_rect|circle_with_flats|circle_with_key
      from_back: number        # Required; alias: x
      from_left: number        # Required; alias: y/z depending on face
      width: number            # Required for most shapes (set to 0 when unused)
      length: number           # Required for rectangle/rounded/key (0 when unused)
      radius: number           # Required for circle/rounded/flats/key (0 when unused)
      depth: number            # Optional (defaults to plane thickness)
      angle: number            # Optional rotation

  snap_joins:
    - pos: number              # Required; distance along wall
      width: number            # Required; snap tab width
      side: front|back|left|right   # Required; maps to yappFront/etc.
```

#### Feature Coordinate Semantics

The DSL no longer exposes a global `coordinates` block. Instead, each feature uses the YAPP default coordinate space for that array:

| Feature        | Coordinates relative to | Behavior when padding changes                  |
|----------------|-------------------------|------------------------------------------------|
| `pcb_stands`   | PCB origin (bottom-left of PCB top surface) | Moves with the PCB when you change padding |
| `connectors`   | PCB origin                               | Moves with the PCB                          |
| `cutouts`      | Box origin (outer bottom-left-back corner) | Stays anchored to the enclosure walls       |
| `snap_joins`   | Box origin                               | Stays anchored to the enclosure walls       |

The generator inserts the correct `undef` placeholders and relies on YAPP's positional defaults; no explicit coordinate flags are emitted for MVP.

### Tolerances

```yaml
tolerances:
  perimeter: number     # Clearance around PCB perimeter
  holes: number         # Clearance for mounting holes
  mating: number        # Clearance between base and lid mating surfaces
```

**Typical values:**
- `perimeter: 0.5` — 0.5mm clearance around PCB
- `holes: 0.3` — 0.3mm larger than nominal hole size
- `mating: 0.2` — 0.2mm gap for lid fit

## Usage Examples

Minimal MVP example with expressions and all four supported feature groups:

```yaml
project: pico-temp-001
version: 0
units: mm
yapp_version: v3.3.8

vars:
  wall_clearance: tolerances.perimeter
  standoff_height: pcb.z_clearance + 4.0
  min_base_thickness: 1.8
  snap_width: 7

pcb:
  length: 51.0
  width: 21.0
  thickness: 1.6
  z_clearance: 2.0
  standoffs:
    type: round
    diameter: 5.0
    height: standoff_height
    screw_d: 2.2

enclosure:
  wall:
    thickness: 2.0
    clearance: wall_clearance
    fillet_radius: 2.0
  base:
    thickness: max(vars.min_base_thickness, enclosure.wall.thickness)

features:
  pcb_stands:
    - x: 8
      y: 8
      diameter: 6.5
    - x: 43
      y: 13
      height: standoff_height
  connectors:
    - x: 15
      y: 8
      stand_height: 4.0
      screw_d: 2.0
      screw_head_d: 4.2
      insert_d: 3.2
      outside_d: 7.5
  cutouts:
    - face: back
      shape: rounded_rect
      from_back: 25
      from_left: 12
      width: 16
      length: 10
      radius: 2
    - face: right
      shape: circle
      from_back: 30
      from_left: 14
      radius: 4
  snap_joins:
    - pos: 20
      width: vars.snap_width
      side: front
    - pos: 25
      width: vars.snap_width
      side: back

tolerances:
  perimeter: 0.5
  holes: 0.3
  mating: 0.2
```

## Expressions

Any numeric field can be an expression instead of a literal number. Expressions are plain YAML scalars (no special prefix like `=` or `$`).

### Syntax

**Literals:**
- Numbers: `2.0`, `51`, `0.5`

**Operators:**
- Arithmetic: `+`, `-`, `*`, `/`
- Parentheses: `(pcb.length + 10.0) / 2`

**Functions:**
- `min(a, b)` — Minimum of two values
- `max(a, b)` — Maximum of two values
- `round(x)` — Round to nearest integer
- `floor(x)` — Round down
- `ceil(x)` — Round up
- `clamp(x, lo, hi)` — Constrain x between lo and hi

**Dotted Paths:**
- Reference any field from the root: `pcb.length`, `enclosure.wall.thickness`, `vars.wall_clearance`
- No `$` prefix needed
- Paths are resolved during fixed-point evaluation

### Examples

```yaml
# Simple reference
clearance: vars.wall_clearance

# Arithmetic
standoff_height: pcb.z_clearance + 4.0

# Functions
base_thickness: max(2.0, enclosure.wall.thickness)

# Complex expression
total_height: (pcb.z_clearance + pcb.thickness + 10.0) * 1.1
```

### Resolution Rules

1. **Fixed-point evaluation:** Resolver iterates up to 16 passes, resolving expressions that reference already-resolved values
2. **Error conditions:**
   - Unresolved paths after 16 passes
   - Circular dependencies (A depends on B, B depends on A)
   - Non-numeric results
3. **Order independence:** You can reference fields defined later in the file; resolver handles ordering

## Type Reference

| Type | Description | Example |
|------|-------------|---------|
| `number` | Literal number or expression | `2.0`, `pcb.length + 5` |
| `integer` | Whole number | `4`, `round(vars.count)` |
| `string` | Text value | `"pico-temp-001"` |
| `enum(...)` | One of specified values | `round`, `hex`, `none` |
| `face_enum` | Face identifier | `front`, `back`, `left`, `right`, `top`, `bottom` |
| `expression` | Numeric expression | `max(a, b)`, `x + y` |

## Validation Rules

- **Required fields:** `project`, `version`, `units`, `yapp_version`, `pcb.length`, `pcb.width`, `pcb.thickness`, `enclosure.wall.thickness`
- **Units:** Always `mm` (no other units supported)
- **Minimum values:** wall/base/lid thickness ≥ 1.0mm; all lengths/dimensions > 0
- **Feature entries:** each required parameter listed in the feature definitions must be present; optional parameters may be omitted
- **Enum values:** `face`, `shape`, and `side` must match the supported sets exactly (case-sensitive)
- **Expressions:** must evaluate to numeric scalars after resolver runs; unresolved expressions are errors

## Related

- Playbook: ../playbook/01-generate-electronics-enclosures-end-to-end-playbook.md
- Analysis: ../analysis/01-analysis-docs-and-requirements-for-enclosure-dsl.md
