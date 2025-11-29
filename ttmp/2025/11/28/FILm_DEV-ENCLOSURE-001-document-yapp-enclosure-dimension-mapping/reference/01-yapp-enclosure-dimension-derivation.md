---
Title: YAPP enclosure dimension derivation
Ticket: FILm_DEV-ENCLOSURE-001
Status: active
Topics:
    - yapp
    - enclosure
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2025-11-28T18:59:45.079962144-05:00
---

# YAPP enclosure dimension derivation

## Goal

Summarize how `YAPPgenerator_v3.scad` converts PCB geometry and enclosure settings into the inside/outer dimensions printed by YAPP, so we can reason about required clearances (for example, the 40 mm buttons in the film-developer project).

## Context

`yappctl` resolves the YAML DSL into a generated `.scad` stub (e.g. `projects/film-developer/enclosure.scad`) that assigns values such as `pcbLength`, `paddingFront`, `baseWallHeight`, and then `include`s `YAPPgenerator_v3.scad`. The generator consumes those assignments through helper arrays/functions, grows the bounding box to cover every PCB present, adds wall clearance and thickness, and finally builds the OpenSCAD solids used for STL export.

## Quick Reference

### PCB inputs

```95:108:pkg/yappgen/assets/YAPPgenerator_v3.scad
//  Parameters:
//   Required:
//    p(0) = name
//    p(1) = length
//    p(2) = width
//    p(3) = posx
//    p(4) = posy
//    p(5) = Thickness
//    p(6) = standoff_Height = Height to bottom of PCB from the inside of the base
//    p(7) = standoff_Diameter
//    p(8) = standoff_PinDiameter
//   Optional:
//    p(9) = standoff_HoleSlack (default to 0.4)
```

- The `pcb` array can host multiple boards; the generator considers each board’s size (`length`, `width`) and its offsets (`posx`, `posy`) when it computes the envelope.

### Dimension pipeline

```324:355:pkg/yappgen/assets/YAPPgenerator_v3.scad
boxLength = maxLength(pcb);
boxWidth = maxWidth(pcb);
…
shellInsideWidth  = boxWidth+paddingLeft+paddingRight;
shellInsideLength = boxLength+paddingFront+paddingBack;
shellInsideHeight = baseWallHeight+lidWallHeight;

shellWidth        = shellInsideWidth+(wallThickness*2);
shellLength       = shellInsideLength+(wallThickness*2);
shellHeight       = basePlaneThickness+shellInsideHeight+lidPlaneThickness;
```

- `boxLength`/`boxWidth` hold the extreme reach of the PCB array before padding.  
- `shellInside*` applies the uniform clearances injected by the DSL (`padding*`).  
- `shell*` applies wall thickness back onto both sides to produce the printed exterior dimensions.  
- Vertical space is determined entirely by the YAML-provided wall heights and plane thicknesses:

```131:136:pkg/yappgen/assets/YAPPgenerator_v3.scad
//-- Total height of box = lidPlaneThickness 
//                       + lidWallHeight 
//                       + baseWallHeight 
//                       + basePlaneThickness
//-- space between pcb and lidPlane :=
//--      (bottonWallHeight+lidWallHeight) - (standoff_Height+pcb_Thickness)
```

### Helpers behind `boxLength` / `boxWidth`

```324:326:pkg/yappgen/assets/YAPPgenerator_v3.scad
boxLength = maxLength(pcb);
boxWidth = maxWidth(pcb);
```

```5692:5694:pkg/yappgen/assets/YAPPgenerator_v3.scad
function maxLength(v, i = 0, r = 0) = i < len(v) ? maxLength(v, i + 1, max(r, v[i][1] + v[i][3])) : r;
function maxWidth(v, i = 0, r = 0) = i < len(v) ? maxWidth(v, i + 1, max(r, v[i][2] + v[i][4])) : r;
```

- Each helper iterates over the `pcb` array and tracks `length + posx` or `width + posy`, so offset boards automatically expand the interior span.

### Padding injected by the DSL

```88:95:projects/film-developer/enclosure.scad
// paddingFront (source: enclosure.wall.clearance) = 1.5 (literal)
paddingFront = 1.5;
// paddingBack (source: enclosure.wall.clearance) = 1.5 (literal)
paddingBack = 1.5;
// paddingLeft (source: enclosure.wall.clearance) = 1.5 (literal)
paddingLeft = 1.5;
// paddingRight (source: enclosure.wall.clearance) = 1.5 (literal)
paddingRight = 1.5;
```

- The Go emitter maps `enclosure.wall.clearance` onto all four padding variables. Asymmetric padding is not exposed in the DSL today; it would require editing the generated SCAD.
- Padding can be zero, but any cutouts or PCB stands near the wall must still respect wall thickness; a zero clearance box means the PCB edge sits flush against the inner wall.

### “PCB envelope vs. shellInside vs. shell” cheat sheet

The generator composes three spans:
- PCB envelope: the far reach of all PCBs and their offsets → `boxLength` × `boxWidth`
- shellInside: add per-side padding to the PCB envelope → `shellInsideLength` × `shellInsideWidth`
- shell (outer): add wall thickness on both sides → `shellLength` × `shellWidth`

```
Top view (length × width)

          <-------- shellLength -------->
        +--------------------------------+  ← shell outer (adds 2×wallThickness)
        |   wallThickness                |
        |   <---- shellInsideLength ---> |  
        |   +--------------------------+ |  ← shellInside = PCB envelope + padding
        |   |  paddingLeft      padding| |
        |   |  +--------------------+  | |  
        |   |  |   PCB envelope     |  | |  ← boxLength (X) / boxWidth (Y)
        |   |  |   (max length+posx)|  | |     (may offset via posx/posy)
        |   |  +--------------------+  | |
        |   |paddingRight             | |
        |   +--------------------------+ |
        |   wallThickness                |
        +--------------------------------+
```

- **Where is boxLength?** It is the horizontal span of the PCB envelope (inner rectangle): the maximum of `(pcb length + posx)` across all PCBs. `shellInsideLength = boxLength + paddingFront + paddingBack`, and `shellLength = shellInsideLength + 2×wallThickness`. The same relationship holds for width.

Vertical cross-section:

```
          lidPlaneThickness
        ┌───────────────────────┐  ← top of shell
        │       lid walls       │  ← lidWallHeight
        │                       │
        │   (interior cavity)   │  ← shellInsideHeight = baseWallHeight + lidWallHeight
        │       base walls      │  ← baseWallHeight
        ├───────────────────────┤  ← basePlaneThickness sits below
```

- The PCB sits `standoffHeight` above the base plane; space above the PCB is `(baseWallHeight + lidWallHeight) - (standoffHeight + pcbThickness)`.
- There is no automatic ridge compensation—ensure `lidWallHeight` ≥ `ridgeHeight` (generator asserts this elsewhere).

### Can we configure per-side padding?

- DSL today exposes only `enclosure.wall.clearance`, applied symmetrically in `model.go`.  
- To vary per side you must edit the emitted SCAD (set `paddingFront`, etc.) or extend the Go resolver/emitter. Any such change would also need a DSL schema update.

### Configuring shellLength / shellWidth in SCAD (no YAML)

Formulas:
- `shellLength  = (maxLength(pcb)) + paddingFront + paddingBack + 2*wallThickness`
- `shellWidth   = (maxWidth(pcb))  + paddingLeft  + paddingRight + 2*wallThickness`

```324:355:pkg/yappgen/assets/YAPPgenerator_v3.scad
boxLength = maxLength(pcb);
boxWidth  = maxWidth(pcb);
shellInsideWidth  = boxWidth+paddingLeft+paddingRight;
shellInsideLength = boxLength+paddingFront+paddingBack;
shellWidth        = shellInsideWidth+(wallThickness*2);
shellLength       = shellInsideLength+(wallThickness*2);
```

Controls in SCAD:
- Set the four paddings (asymmetric allowed):

```71:74:pkg/yappgen/assets/YAPPgenerator_v3.scad
paddingFront        = 1;
paddingBack         = 1;
paddingRight        = 1;
paddingLeft         = 1;
```

- Set `wallThickness`.
- Override the `pcb` vector to adjust the PCB envelope (`boxLength`/`boxWidth`) via size or offsets:

```118:122:pkg/yappgen/assets/YAPPgenerator_v3.scad
pcb = 
[
  ["Main",              pcbLength,pcbWidth,    0,0,    pcbThickness,  standoffHeight, standoffDiameter, standoffPinDiameter, standoffHoleSlack]
];
```

Example (SCAD only):

```scad
// after include <YAPPgenerator_v3.scad>
paddingFront = 2.0;
paddingBack  = 3.0;
paddingLeft  = 1.0;
paddingRight = 4.0;
wallThickness = 3.0;

// Push PCB 5 mm in X so the envelope grows by 5
pcb = [
  ["Main", 90, 70, 5, 0, 1.6, 3, 6, 3, 0.4]
];

YAPPgenerate();
```

## Usage Examples

### Deriving enclosure size for the film developer PCB

Inputs from `projects/film-developer/enclosure.yaml`:

- `pcb.length = 90`, `pcb.width = 70`, `enclosure.wall.clearance = 1.5`.
- `baseWallHeight = 6.6`, `lidWallHeight = 66`, `basePlaneThickness = 2`, `lidPlaneThickness = 2`, `standoffHeight = 3`, `pcbThickness = 1.6`.

Step-by-step:

1. `boxLength = 90` and `boxWidth = 70` (single PCB with no offsets).  
2. `shellInsideLength = 90 + 1.5 + 1.5 = 93`; `shellInsideWidth = 70 + 1.5 + 1.5 = 73`.  
3. `shellInsideHeight = 6.6 + 66 = 72.6`.  
4. Outer dimensions:  
   - `shellLength = 93 + 2 * 2.4 = 97.8` mm  
   - `shellWidth  = 73 + 2 * 2.4 = 77.8` mm  
   - `shellHeight = 2 + 72.6 + 2 = 76.6` mm
5. Vertical clearance above the PCB: `72.6 - (3 + 1.6) = 68` mm, matching the `top_clearance` encoded in YAML for the 40 mm button stack.

These numbers match the OpenSCAD echo output during `yappctl generate`, confirming the mapping between YAML inputs and printed dimensions.

## Related

- `pkg/yappgen/assets/YAPPgenerator_v3.scad`
- Example output: `projects/film-developer/enclosure.scad`
