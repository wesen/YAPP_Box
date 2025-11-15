---
Title: PushButtons DSL gap analysis
Ticket: YAPP-PUSH-BUTTONS-001
Status: active
Topics:
  - yapp
  - dsl
  - buttons
DocType: analysis
Intent: short-term
Owners: []
RelatedFiles:
  - Path: examples/YAPP_Demo_buttons_v31.scad
    Note: Baseline reference for two-button lid layout
  - Path: examples/YAPP_Demo_buttons2_v31.scad
    Note: Shape variations (circle, polygon, rounded rectangle)
  - Path: examples/yapp-demo-buttons.yaml
    Note: Current DSL facade missing push_buttons node
  - Path: pkg/yappgen/schema.go
    Note: Where new push button schema will be defined
  - Path: pkg/yappgen/map.go
    Note: Build helpers that turn DSL nodes into YAPP arrays
ExternalSources: []
Summary: "Captures the SCAD-level pushButtons capabilities, the gaps in the current YAML DSL, and the plan for mapping YAML into the YAPP parameter arrays."
LastUpdated: 2025-11-15T16:10:00-05:00
---

# PushButtons DSL gap analysis

## 1. What the SCAD examples do today

Both `examples/YAPP_Demo_buttons_v31.scad` and `examples/YAPP_Demo_buttons2_v31.scad` drive the legacy `pushButtons` array directly. The comments in those files highlight the positional parameters:

| Index | Field | Meaning | Notes |
|-------|-------|---------|-------|
| 0 | posx | From PCB origin (yappCoordPCB) | Aligns with DSL `x` |
| 1 | posy | From PCB origin | Aligns with DSL `y` |
| 2 | capLength | Button length or diameter (for circular buttons) |
| 3 | capWidth | Button width (unused for circles) |
| 4 | capRadius | Corner radius or circle radius |
| 5 | capAboveLid | How far the cap protrudes (+) or sinks (-) relative to lid top |
| 6 | switchHeight | Physical switch height |
| 7 | switchTravel | Travel distance |
| 8 | poleDiameter | Center post diameter |
| 9 | heightToPCB | Optional; defaults to `standoffHeight + pcbThickness` |
| 10 | shape flag | `yappRectangle`, `yappCircle`, `yappPolygon`, `yappRoundedRect`, `yappCircleWithFlats`, `yappCircleWithKey` |
| 11 | angle | Rotation |
| 12 | filletRadius | Overrides auto fillet |
| 13 | buttonWall | Lid wall around opening |
| 14 | buttonPlateThickness | Plate thickness |
| 15 | buttonSlack | Clearance between cap and guide |
| 16 | snapSlack | Clearance for snap joists (optional in v31) |
| n(a) | coord flag | `yappCoordPCB` (default) / `yappCoordBox` / `yappCoordBoxInside` |
| n(b) | origin flag | `yappGlobalOrigin` / `yappLeftOrigin` |

`YAPP_Demo_buttons2_v31.scad` mixes multiple shapes: circle, polygon with `shapeArrow`, rounded rect, standard rectangle. The YAML DSL currently cannot express any of these, which is why `examples/yapp-demo-buttons.yaml` carries a TODO.

## 2. Current DSL coverage

The DSL supports:

- `pcb_stands`, `connectors`, `cutouts`, `snap_joins`, etc.
- Basic box measurements, tolerances, and doc metadata.

Missing pieces for parity with the SCAD reference:

1. **`push_buttons` node** under `features`.
2. Schema describing the positional numeric fields and enumerated flags (shape, coordinate, origin).
3. Mapping code in `pkg/yappgen` to convert YAML into the `pushButtons` array.
4. Validation (e.g., shape-specific dimension usage, optional defaults).
5. Derived settings: enabling `printSwitchExtenders = true` when at least one push button exists, plus gating show flags.

## 3. Proposed DSL structure

```yaml
features:
  push_buttons:
    - name: reset
      x: 15
      y: 14
      cap:
        length: 8
        width: 6
        radius: 0
      lid:
        protrusion: 3          # capAboveLid
        wall: 2.0
        plate_thickness: 2.5
        slack: 0.25
        snap_slack: 0.1
      switch:
        height: 5.5
        travel: 1.0
        pole_diameter: 3.5
        top_offset: null       # -> undef in SCAD
      shape: rectangle         # enum
      angle: 0
      fillet_radius: null
      coordinate: pcb          # pcb|box|box_inside
      origin: global           # global|left
      polygon: arrow           # optional mask if shape = polygon
```

Mapping rules:

- Null → `undef` so YAPP uses defaults.
- `shape` enumerates to the SCAD flags; when `polygon`, require `polygon` field containing a known symbol (e.g., `shapeArrow`).
- Coordinates/origin become raw tokens appended to the parameter slice (mirrors connectors/cutouts conventions).
- Absent `snap_slack` maintains compatibility with v31 defaults (0.10).

## 4. Implementation steps

1. **Schema + builder**: extend `pkg/yappgen/schema.go` with `pushButtonsSchema` and add a `buildPushButtons` helper similar to `buildPcbStands`.
2. **Model wiring**: update `Model` struct (likely `pkg/yappgen/model.go`) so `BuildModel` emits `pushButtons` arr when DSL node present; toggle `printSwitchExtenders = true`.
3. **DSL validation**: update resolver to ensure `features.push_buttons` is optional, arrays default to empty, and shape-specific dimension rules are enforced before transformation.
4. **Examples**: add `push_buttons` block to `examples/yapp-demo-buttons.yaml` and create a second YAML referencing the `buttons2` SCAD file to prove parity.
5. **Docs**: produce user-facing instructions (CLI help page in `pkg/docs/tutorials/`) describing the new YAML block (requested in this ticket).

## 5. Open questions

- Should we support custom polygon definitions via inline point arrays, or only alias the built-in shapes (arrow, triangle, etc.)? Initial recommendation: alias only built-ins to avoid exposing raw polygon DSL in MVP.
- How to expose shape presets in YAML? E.g., `shape_preset: arrow` vs `shape: polygon` + `preset: arrow`.
- What defaults should we pick for `buttonWall`, `buttonPlateThickness`, `buttonSlack`, `snapSlack`? We can mirror SCAD defaults (2.0 / 2.5 / 0.25 / 0.1).

These questions need resolution before coding, but the structure above reflects the SCAD capabilities we must match.
