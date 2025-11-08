---
Title: 'Generate Electronics Enclosures: End-to-end Playbook'
Ticket: YAPP-ENCL-DSL-001
Status: active
Topics:
    - playbook
    - electronics
    - 3d-printing
    - yapp
    - openscad
    - automation
DocType: playbook
Intent: long-term
Owners:
    - manuel
RelatedFiles: []
ExternalSources: []
Summary: Step-by-step flow from YAML to printed enclosure using OpenSCAD/YAPP
LastUpdated: 2025-11-08T15:56:18.923657506-05:00
---




# Generate Electronics Enclosures: End-to-end Playbook

## Purpose

Define and execute the end-to-end flow from a YAML DSL (with variables and expressions) to a parametric OpenSCAD enclosure and printable STL, with validation and QA checkpoints.

## Environment Assumptions

- OpenSCAD installed
- YAPPgenerator v3.3.8 available and version pinned in the DSL (`yapp_version`)
- Resolver tool available to expand expressions (fixed-point) before mapping to OpenSCAD

## Variables and Expressions

- Define variables under `vars:` at the top level.
- Any numeric field may be a literal number or an expression written as a plain scalar (no leading `=`, quotes optional).
- Expressions reference values via absolute dotted paths from the root (e.g., `pcb.length`, `vars.wall_clearance`) and support `+ - * /`, parentheses, and `min/max/round/floor/ceil/clamp`.

Example excerpt:

```yaml
vars:
  wall_clearance: tolerances.perimeter
  standoff_height: pcb.z_clearance + 4.0

enclosure:
  wall:
    thickness: 2.0
    clearance: vars.wall_clearance

pcb:
  z_clearance: 2.0
  standoffs:
    height: vars.standoff_height
```

## Commands

<!-- List of commands to execute -->

```bash
# 1) Validate schema and expand expressions (fixed-point)
# dsl-resolve input.yaml > resolved.yaml

# 2) Map resolved YAML to OpenSCAD (YAPPgenerator parameters/functions)
# dsl-to-openscad resolved.yaml > model.scad

# 3) Render previews for QA
# openscad -o preview_base.png -D part=\"base\" model.scad
# openscad -o preview_lid.png  -D part=\"lid\"  model.scad

# 4) Export STL for print
# openscad -o enclosure_base.stl -D part=\"base\" model.scad
# openscad -o enclosure_lid.stl  -D part=\"lid\"  model.scad
```

## Exit Criteria

- Expressions fully resolved (no unresolved paths; reached fixed point)
- Geometry validates (feature bounds within box; minimum thickness and hole sizes)
- Visual QA passes (previews match expectations)
- Printed parts fit target hardware within tolerances

## Notes

- Resolution semantics:
  - Iterate expression evaluation up to 16 passes or until no values change.
  - Error on unresolved dependencies, invalid paths, or non-numeric results.
  - Prefer BoxInside coordinates for interior features; set `coordinates.origin` explicitly.
