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
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2025-11-08T16:18:37.598375185-05:00
---


# Enclosure DSL: Language Reference

## Goal

Define a concise, complete reference for the Enclosure DSL used to describe electronics enclosures, including structure, types, expressions, and constraints. This document is authoritative for writers and implementers.

## Context

This DSL is consumed by a resolver that evaluates variables and expressions to produce a fully numeric configuration before mapping to YAPPgenerator/OpenSCAD. Expressions are plain scalars using dotted paths (no `=`, no `$`), and are resolved to a fixed point.

## Structure Overview

The DSL is organized into seven top-level sections:

1. **Metadata** — `project`, `version`, `units`, `yapp_version`
2. **Variables** — `vars` (reusable computed values)
3. **PCB** — `pcb` (dimensions, mounting holes, standoffs)
4. **Enclosure** — `enclosure` (walls, base, lid)
5. **Features** — `features` (holes, cutouts, light tubes, etc.)
6. **Coordinates** — `coordinates` (origin and reference plane)
7. **Tolerances** — `tolerances` (print fit adjustments)

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

```yaml
features:
  holes:                        # Circular through-holes
    - face: face_enum
      x: number
      y: number
      d: number                 # Diameter
  cutouts:                      # Rectangular cutouts
    - face: face_enum
      x: number
      z: number                 # Z coordinate for side faces
      width: number
      height: number
      fillet: number            # Corner fillet radius
  light_tubes:                  # LED light pipes
    - face: face_enum
      x: number
      y: number
      lens_d: number            # Lens diameter
      tube_d: number            # Tube outer diameter
      depth: number             # Tube depth into enclosure
```

**Face enum values:**
- `top` — Lid top surface
- `bottom` — Base bottom surface
- `side_x+` — Side along positive X axis
- `side_x-` — Side along negative X axis
- `side_y+` — Side along positive Y axis
- `side_y-` — Side along negative Y axis

**Notes:**
- Coordinates depend on `coordinates.origin` setting
- For side faces, use `x` or `y` (parallel to face) and `z` (height from base)
- `light_tubes` typically placed on `top` face aligned with PCB LEDs

### Coordinates

```yaml
coordinates:
  origin: pcb|box|boxinside      # Coordinate system origin
  reference_plane: base|lid      # Z-axis reference
```

**Origin options:**
- `pcb` — Origin at PCB bottom-left corner
- `box` — Origin at enclosure outer bottom-left corner
- `boxinside` — Origin at enclosure inner bottom-left corner (recommended)

**Reference plane:**
- `base` — Z=0 at base bottom surface
- `lid` — Z=0 at lid top surface

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

Minimal example with expressions:

```yaml
project: pico-temp-001
version: 0
units: mm
yapp_version: v3.3.8

vars:
  wall_clearance: tolerances.perimeter
  standoff_height: pcb.z_clearance + 4.0

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
| `face_enum` | Face identifier | `top`, `side_x+` |
| `expression` | Numeric expression | `max(a, b)`, `x + y` |

## Validation Rules

- **Required fields:** `project`, `version`, `units`, `yapp_version`, `pcb.length`, `pcb.width`, `enclosure.wall.thickness`
- **Units:** Always `mm` (no other units supported)
- **Minimum values:**
  - `enclosure.wall.thickness >= 1.0`
  - `enclosure.base.thickness >= 1.0`
  - All dimensions > 0
- **Enum values:** Must match exactly (case-sensitive)
- **Coordinates:** Must be numeric after resolution

## Related

- Playbook: ../playbook/01-generate-electronics-enclosures-end-to-end-playbook.md
- Analysis: ../analysis/01-analysis-docs-and-requirements-for-enclosure-dsl.md
