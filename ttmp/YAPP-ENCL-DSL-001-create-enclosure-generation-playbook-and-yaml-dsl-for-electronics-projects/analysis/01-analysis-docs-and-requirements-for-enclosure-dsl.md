---
Title: 'Analysis: Docs and requirements for enclosure DSL'
Ticket: YAPP-ENCL-DSL-001
Status: active
Topics:
    - playbook
    - electronics
    - 3d-printing
    - yapp
    - openscad
    - automation
DocType: analysis
Intent: long-term
Owners:
    - manuel
RelatedFiles:
    - Path: ttmp/PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/README.md
      Note: Case study and measurements for DSL validation
    - Path: ttmp/PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/design/design-rationale-and-calculations.md
      Note: Rationale and calculations to verify dimensions
    - Path: ttmp/PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/pico_temp_monitor_base.stl
      Note: Printed part artifact for fit verification
    - Path: ttmp/PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/pico_temp_monitor_base_only.scad
      Note: Base-only model for fast iteration
    - Path: ttmp/PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/pico_temp_monitor_box.scad
      Note: Primary OpenSCAD model for mapping DSL outputs
    - Path: ttmp/PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/pico_temp_monitor_lid_only.scad
      Note: Lid-only model to iterate lid features
    - Path: ttmp/PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/preview_base.png
      Note: Rendered base preview for QA
    - Path: ttmp/PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/preview_lid.png
      Note: Rendered lid preview for QA
    - Path: ttmp/PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/various/component-specifications.md
      Note: Component specs driving cutouts and light tubes
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/analysis/README.md
      Note: Doc analysis context for YAPPgenerator
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/analysis/USAGE_EXAMPLES.md
      Note: Usage examples relevant to feature semantics
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/fixes/06-gitbook-update-light-tubes.md
      Note: Light tubes doc updates informing defaults
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/fixes/07-gitbook-update-push-buttons.md
      Note: Push buttons doc updates informing defaults
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/fixes/08-gitbook-update-pcb-stands-parameter-order.md
      Note: PCB stands doc updates and ordering
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/code-examples-yappcircle-yappcenter-corrections.md
      Note: Corrected examples informing naming and validation
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/coordinate-systems-pcb-vs-box-vs-boxinside-comparison.md
      Note: Coordinate systems reference for feature placement mapping
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/default-values-for-optional-parameters.md
      Note: Defaults to seed DSL fallback values
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/pcb-stands-parameter-order-verification.md
      Note: Standoffs naming/order caveats relevant to DSL mapping
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/various/test_light_tubes_center.scad
      Note: SCAD tests for light tube parameters
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/various/test_push_buttons_center.scad
      Note: SCAD tests for push button placement
ExternalSources: []
Summary: Survey of existing docs and requirements to define an enclosure DSL
LastUpdated: 2025-11-08T15:56:18.992952845-05:00
---






# Analysis: Docs and requirements for enclosure DSL

## Purpose and Scope

Define a concise, YAML-based DSL that maps directly to YAPPgenerator/OpenSCAD parameters to generate robust, dimensionally-accurate electronics enclosures. This analysis aggregates prior work, extracts requirements, and outlines the DSL structure and mapping.

## Relevant Documents

- Coordinate Systems — PCB vs Box vs BoxInside: [link](../../YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/coordinate-systems-pcb-vs-box-vs-boxinside-comparison.md)
- Default Values for Optional Parameters: [link](../../YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/default-values-for-optional-parameters.md)
- PCB Stands — Parameter Order Verification: [link](../../YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/pcb-stands-parameter-order-verification.md)
- Code Examples — yappCircle/yappCenter Corrections: [link](../../YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/code-examples-yappcircle-yappcenter-corrections.md)
- Case Study: PICO-TEMP-001 — README: [link](../../PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/README.md)
- PICO-TEMP-001 — OpenSCAD models: [box.scad](../../PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/pico_temp_monitor_box.scad), [base_only.scad](../../PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/pico_temp_monitor_base_only.scad), [lid_only.scad](../../PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/pico_temp_monitor_lid_only.scad)

## Key Findings and Technical Insights

- Coordinate systems: Prefer placing interior features using BoxInside coordinates; convert when necessary. Keep origin choices explicit (pcb, box, boxinside).
- Defaults: Many features have sensible defaults (e.g., fillets, light tubes) that should be encoded in the DSL and surfaced as overrides.
- Parameter order pitfalls: Functions like PCB stands have ordering and naming caveats; the DSL must be name-oriented (keyed) rather than positional.
- Example inconsistencies: Some example corrections inform validation and naming (e.g., yappCircle/yappCenter). DSL should standardize names and validate common pitfalls.
- Real-world checks: PICO-TEMP-001 provides concrete dimensions and printed results; use it as the first verification target for the DSL.

## DSL Goals

- Readable, minimal YAML focused on enclosure intent (pcb size, walls, lid, features) over raw function calls.
- Deterministic mapping to YAPPgenerator with version pinning for reproducibility.
- Safe defaults with clear units (mm) and tolerances for print fit.
- Extensible features list (holes, cutouts, standoffs, vents, light tubes, push buttons, connectors).

## Variables and Expressions

- Define variables at the top level under `vars:`. Values can be numbers or expressions.
- Any numeric field may be either:
  - a number (e.g., `2.0`), or
  - an expression written as a plain scalar (no `=`, quotes optional).
- Path references are absolute dotted paths from the document root (no `$`). Examples: `pcb.length`, `enclosure.wall.thickness`, `vars.wall_clearance`.
- Allowed operators/functions (numeric only): `+ - * /`, parentheses, and `min(x,y)`, `max(x,y)`, `round(x)`, `floor(x)`, `ceil(x)`, `clamp(x,lo,hi)`.
- Units: all numeric values are millimeters unless otherwise stated. Expressions must be unit-consistent.

Resolution semantics (fixed-point):
1. Initialize with all literal numbers known; mark expressions as unresolved.
2. Iteratively attempt to evaluate unresolved expressions using currently known values.
3. Stop when no values change between passes (fixed point) or after a maximum of 16 iterations.
4. If unresolved remain after max iterations, report each unresolved path and its missing dependencies (error).
5. Immediate self-references or invalid paths error early; division-by-zero and non-numeric results error.
6. Implementations should cache evaluated results and may detect cycles to fail fast; still follow the fixed-point rule above.

## Proposed DSL Outline (v0)

```yaml
project: pico-temp-001
version: 0
units: mm
yapp_version: v3.3.8

vars:
  wall_clearance: tolerances.perimeter
  standoff_height: pcb.z_clearance + 4.0
  min_base_thickness: 2.0

pcb:
  length: 51.0
  width: 21.0
  thickness: 1.6
  z_clearance: 2.0         # space between pcb bottom and base
  holes:
    - x: 3.0;  y: 3.0;  d: 2.5
    - x: 48.0; y: 3.0;  d: 2.5
  standoffs:
    type: round            # round | hex | none
    diameter: 5.0
    height: vars.standoff_height
    screw_d: 2.2

enclosure:
  wall:
    thickness: 2.0
    clearance: vars.wall_clearance   # internal clearance from pcb outline
    fillet_radius: 2.0
  base:
    thickness: max(vars.min_base_thickness, enclosure.wall.thickness)
  lid:
    type: screws           # screws | snap
    screws:
      count: 4
      positions: corners   # corners | custom
      head_clearance: 0.6

features:
  holes:                   # circular through-holes on faces
    - face: top
      x: 10.0; y: 10.0; d: 6.0
  cutouts:                 # rectangular cutouts
    - face: side_x+
      x: 40.0; z: 12.0
      width: 10.0; height: 6.0; fillet: 1.0
  light_tubes:
    - face: top
      x: 15.0; y: 8.0
      lens_d: 3.0; tube_d: 4.0; depth: 2.5

coordinates:
  origin: boxinside        # pcb | box | boxinside
  reference_plane: base    # base | lid

tolerances:
  perimeter: 0.5
  holes: 0.3
  mating: 0.2
```

## Mapping to YAPPgenerator (high-level)

- pcb.length/width/thickness → PCB parameters
- enclosure.wall.thickness → box wall thickness
- enclosure.wall.fillet_radius → corner fillet
- enclosure.wall.clearance → boxInside clearance
- pcb.standoffs.* → yapp stands function (named args)
- features.holes[*] → yappHole
- features.cutouts[*] → yappCutoutRect (with fillet)
- features.light_tubes[*] → light tube helpers per docs
- coordinates.origin/reference_plane → placement transforms
- Expressions resolve before mapping; resolved numbers feed the YAPP parameters

## Validation Ideas

- Schema: keys required (pcb.length/width, wall.thickness), permissible enums with fallbacks.
- Geometry checks: ensure all feature coordinates lie within box limits given origin.
- Printability: minimum wall thickness, hole sizes, and fillet constraints.
- Expression checks: detect unresolved paths, cycles, non-numeric results; cap iterations (16) and report on failure.

## Next Steps

- Finalize v0 schema and enums; document defaults
- Build YAML example for PICO-TEMP-001 and generate STL
- Validate fit on hardware; adjust tolerances/defaults
- Write mapping table (DSL → YAPP functions/params) in the playbook
- Add troubleshooting and QA checklist
