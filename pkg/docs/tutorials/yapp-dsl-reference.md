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

The YAML DSL includes a powerful expression engine that eliminates hard-coded "magic numbers" and lets you define relationships between dimensions declaratively. Instead of scattering `2.4` or `pcb_width + 6` throughout your file and manually keeping them in sync, you define constants and computed values once in the `vars` block, then reference them anywhere. The resolver evaluates expressions in multiple passes, so variable `b` can depend on `a`, and `c` can depend on `b`, as long as there are no circular dependencies.

### How it works

**Basic syntax:**
- **Literals:** Any numeric field can be a plain number: `thickness: 2.4`.
- **Expressions:** Reference other fields using dot-notation: `thickness: enclosure.wall.thickness + 0.5`.
- **Variables:** Define reusable constants in the `vars` block: `vars.min_wall: 2.0`, then use them elsewhere: `thickness: vars.min_wall`.
- **Functions:** Built-in math functions include `min`, `max`, `round`, `floor`, `ceil`, and `clamp`: `thickness: max(vars.min_wall, pcb.thickness * 2)`.

**Resolution passes:**
The resolver makes multiple passes through the document (up to 16 by default) to handle chains of dependencies. For example:

```yaml
vars:
  base_wall: 2.0
  derived_thickness: base_wall + 0.4
  final_value: derived_thickness * 1.5

enclosure:
  wall:
    thickness: vars.final_value  # Resolves to 3.6
```

On the first pass, `base_wall` resolves to `2.0`. On the second pass, `derived_thickness` computes to `2.4`. On the third pass, `final_value` becomes `3.6`, and `wall.thickness` gets its final value. Circular references (e.g., `a: b`, `b: a`) cause an error after exhausting the max iteration count.

**Variable hoisting:**
For convenience, `vars` entries are hoisted into the global scope, so you can omit the `vars.` prefix in most contexts:

```yaml
vars:
  pcb_clearance: 1.0

enclosure:
  wall:
    clearance: pcb_clearance  # Same as vars.pcb_clearance
```

**Strict mode:**
Run `yappctl resolve --strict` (or `generate --strict`) to enable stricter validation:
- Unknown top-level keys produce errors (catches typos like `pcb.widht`).
- Unused `vars` entries are flagged, preventing dead code accumulation.
- Expression errors become fatal rather than warnings.

### Practical examples

**Keeping related dimensions in sync:**

```yaml
vars:
  connector_spacing: 10.0
  num_connectors: 4

features:
  connectors:
    - x: connector_spacing
      y: 10
      # ... other fields
    - x: connector_spacing * 2
      y: 10
    - x: connector_spacing * 3
      y: 10
```

**Computing enclosure dimensions from PCB size:**

```yaml
vars:
  padding: 3.0

pcb:
  length: 50.0
  width: 30.0

enclosure:
  wall:
    clearance: padding
```

**Printer-specific tuning:**

```yaml
vars:
  # My printer runs tight on holes
  hole_compensation: 0.5

tolerances:
  holes: hole_compensation
```

This mechanism is ideal for parameterizing designs, experimenting with different configurations (swap one `vars` block for another), and keeping derived values consistent as you iterate on your enclosure geometry.

## PCB section

The `pcb` block defines the physical dimensions of the printed circuit board your enclosure will house. These measurements establish the core bounding box that all other features reference. The origin sits at the PCB's bottom-left corner when viewed from above, with `x` running along the length, `y` along the width, and `z` perpendicular (height). Getting these dimensions right is critical: too tight and the board won't fit; too loose and you waste material or introduce rattling.

| Field | Required | Default | Notes |
|-------|----------|---------|-------|
| `length` | ✓ | — | PCB dimension along the X axis (mm). Measure edge-to-edge. |
| `width` | ✓ | — | PCB dimension along the Y axis (mm). |
| `thickness` | ✓ | — | Board substrate thickness (typically `1.6` for standard FR4). |
| `z_clearance` | optional | `1.0` | Vertical gap between PCB bottom and enclosure base; maps to SCAD `standoffHeight`. Increase this if you have tall solder joints or through-hole components on the underside. |
| `standoffs.diameter` | optional | `6.0` | Default outer diameter for standoff columns. Individual `pcb_stands` entries can override this. |
| `standoffs.screw_d` | optional | `2.4` | Default screw/pin hole diameter. Applies to both threaded inserts and friction-fit pins unless a specific stand overrides it. |
| `standoffs.hole_slack` | optional | inherits `tolerances.holes` | Extra clearance added to standoff holes. Useful if your standoff material (e.g., metal inserts vs. printed pegs) needs tighter or looser fits than the global tolerance. |

### Measurement tips

**Finding PCB dimensions:**
- Use calipers for accurate edge-to-edge measurements.
- If the board has irregular edges (mounting ears, cutouts), measure the rectangular bounding box that fully encloses it.
- For thickness, standard boards are `1.6 mm`, but thin designs may use `0.8` or `1.0`, and thick ones (high-current) can be `2.0` or more.

**Setting `z_clearance`:**
- Inspect the bottom of your PCB. Measure the tallest feature (solder blob, through-hole lead, surface-mount component).
- Add at least `0.5 mm` safety margin. A common safe default is `1.5` to `2.0 mm`.
- If you have very tall components underneath (e.g., barrel jacks, screw terminals), you might need `3.0 mm` or more.

### Example

```yaml
pcb:
  length: 65.0          # Measured with calipers
  width: 45.0
  thickness: 1.6        # Standard FR4
  z_clearance: 2.0      # Room for solder joints + safety margin
  
  standoffs:
    diameter: 6.0       # Sturdy enough for M3 inserts
    screw_d: 2.5        # Slightly larger than M2.5 for easier insertion
    hole_slack: 0.4     # Inherits from global tolerances
```

This configuration creates a shell sized for a `65 × 45 mm` board with standard thickness, lifting it `2 mm` above the base, and preparing standoff geometry compatible with M2.5 or M3 hardware.

## Enclosure section

The `enclosure` block controls the structural shell that wraps around your PCB: wall thickness, base and lid plates, and the interlocking ridge that holds the two halves together. These parameters directly impact printability, strength, and how easily the box opens and closes. Too thin and the walls may warp or crack; too thick and you waste filament. The ridge geometry is especially important: it must be tall enough for a secure friction fit yet toleranced correctly so the lid doesn't bind.

### Wall subsection

| Field | Required | Default | Notes |
|-------|----------|---------|-------|
| `thickness` | ✓ | — | Wall material thickness (mm). Typical range: `1.8` to `3.0`. Thinner walls save material but reduce rigidity; thicker walls improve durability. |
| `clearance` | optional | `1.0` | Uniform padding between the PCB edge and the inner wall surface on all four sides. Increase if components overhang the board edge. |
| `fillet_radius` | optional | `2.0` | Rounds internal corners where walls meet the base. Larger radii improve printability and strength but consume more interior space. |

**Design note:** `wall.clearance` applies symmetrically to front/back/left/right. If you need asymmetric padding (e.g., extra clearance on one side for a connector), compute custom offsets in `vars` and adjust cutout positions accordingly. Native per-side clearance is on the roadmap.

### Base and lid subsections

| Field | Required | Default | Notes |
|-------|----------|---------|-------|
| `base.thickness` | ✓ | — | Bottom plate thickness (mm). Common values: `1.2` to `2.0`. Must be thick enough to support standoffs without flexing. |
| `lid.thickness` | ✓ | — | Top plate thickness (mm). Often matches `base.thickness` for visual consistency. |

**Strength considerations:**
- If your base hosts many standoffs or heavy components, use at least `1.6 mm`.
- Very thin bases (`< 1.2 mm`) can work for lightweight projects but may warp during printing.
- Lid thickness can be thinner than the base if it doesn't bear loads, saving print time.

### Ridge subsection

The ridge is the raised lip on one half (typically the base) that nests into a matching groove in the other half (the lid), creating a friction-fit closure.

| Field | Required | Default | Notes |
|-------|----------|---------|-------|
| `height` | ✓ | — | Vertical extent of the ridge (mm). YAPP recommends at least `1.8 × wall.thickness` for a secure fit. Taller ridges improve strength but increase material use. |
| `slack` | optional | `0.2` | Horizontal gap between ridge and groove. Larger values ease assembly; smaller values create a tighter seal. |
| `gap` | optional | `0.5` | Vertical clearance that prevents the lid from bottoming out before the ridge engages. Too small causes binding; too large allows wobble. |

**Tuning tips:**
- Start with the defaults (`height: 5.0`, `slack: 0.2`, `gap: 0.5`).
- If the lid is too tight, increase `slack` by `0.1 mm` increments.
- If the lid rattles, decrease `slack` or increase `height`.
- If the lid doesn't sit flush, adjust `gap`.

### Complete example

```yaml
enclosure:
  wall:
    thickness: 2.4        # Sturdy without being bulky
    clearance: 1.5        # Extra room for edge-mounted components
    fillet_radius: 2.0    # Smooth internal corners
  
  base:
    thickness: 1.8        # Thick enough to anchor standoffs
  
  lid:
    thickness: 1.6        # Slightly thinner to save material
  
  ridge:
    height: 5.0           # 2.08 × wall.thickness, well above minimum
    slack: 0.25           # Slightly loose for easy open/close
    gap: 0.5              # Standard clearance
```

This configuration produces a robust enclosure with comfortable hand-feel and reliable closure mechanics. Adjust `clearance` and `ridge.slack` based on your printer's tolerances and the intended use case (frequent access vs. permanent seal).

## Features section

The `features` block is where your enclosure comes to life. This is where you define PCB standoffs, carve holes for connectors, add snap joints, and specify other mechanical elements. Each feature type lives in its own named array (e.g., `features.pcb_stands`, `features.cutouts`), and each entry is a map of named fields. The generator translates these human-friendly names into the positional SCAD arrays that `YAPPgenerator_v3.scad` expects, automatically computing coordinate transformations, setting flags, and filling defaults.

Think of this section as your bill of materials: every standoff screw location, every USB cutout, every snap-fit clip gets its own declarative entry. Changes are easy—add a new cutout by appending a map to the array—and the expressiveness of YAML (anchors, aliases, comments) helps you stay organized as designs grow complex.

### `pcb_stands`

PCB stands are the vertical columns that lift and secure your board above the enclosure base. Each stand can be a friction-fit peg, a threaded insert receiver, or a simple spacer. You specify `x` and `y` coordinates relative to the PCB origin (bottom-left corner, viewed from above), and the generator creates cylindrical geometry with optional screw holes and fillets.

| YAML field | Required | Description |
|------------|----------|-------------|
| `x`, `y` | ✓ | Position on the PCB (mm), measured from the board's origin at `[0, 0]`. Typically align with mounting holes in your PCB layout. |
| `height` | optional | Stand height from base to PCB bottom (mm). Defaults to `pcb.z_clearance`. Increase if this particular area needs extra clearance for a tall component. |
| `pcb_gap` | optional | Extra vertical offset (mm). Useful when sandwiching boards or lifting sections higher than the default clearance. |
| `diameter` | optional | Outer diameter of the standoff column (mm). Defaults to `pcb.standoffs.diameter`. Make it larger if you need more material around the screw hole for strength. |
| `pin_diameter` | optional | Inner hole diameter (mm). Defaults to `pcb.standoffs.screw_d`. Set to `0` for solid pegs (friction-fit into PCB holes), or size for M2, M2.5, M3 screws/inserts. |
| `hole_slack` | optional | Added to `pin_diameter` to loosen the fit (mm). Defaults to `pcb.standoffs.hole_slack` or `tolerances.holes`. Increase if your printer runs tight. |
| `fillet_radius` | optional | Radius of the fillet where the stand meets the base (mm). Larger fillets improve strength and printability but consume floor space. |

**Coordinate tips:**
- Measure standoff positions from your PCB CAD files or use calipers on a physical board.
- If your PCB has four corner mounting holes at `(3, 3)`, `(62, 3)`, `(62, 42)`, `(3, 42)` for a `65 × 45 mm` board, create four stands at those coordinates.
- Offset inward from the PCB edge by at least `2–3 mm` to avoid interfering with enclosure walls.

**Example: Four-corner standoffs with M2.5 inserts**

```yaml
features:
  pcb_stands:
    - x: 3.0
      y: 3.0
      diameter: 6.0
      pin_diameter: 2.5    # M2.5 insert
      fillet_radius: 1.5
    
    - x: 62.0
      y: 3.0
      diameter: 6.0
      pin_diameter: 2.5
      fillet_radius: 1.5
    
    - x: 62.0
      y: 42.0
      diameter: 6.0
      pin_diameter: 2.5
      fillet_radius: 1.5
    
    - x: 3.0
      y: 42.0
      diameter: 6.0
      pin_diameter: 2.5
      fillet_radius: 1.5
```

This creates a stable four-point suspension compatible with M2.5 threaded inserts, ready for assembly with M2.5 × 8mm screws.

### `cutouts`

Cutouts are openings in the enclosure walls, base, or lid for cables, connectors, buttons, ventilation, or displays. You specify which face of the box to cut, where along that face to position the opening, the dimensions, and the shape profile. The generator automatically handles the 3D geometry and coordinate transformations so you can think in terms of "10mm from the left edge of the back face" rather than wrestling with OpenSCAD coordinates.

| Field | Required | Description |
|-------|----------|-------------|
| `face` | ✓ | Which surface to cut: `front`, `back`, `left`, `right`, `base`, or `lid`. |
| `from_back` | ✓ | Distance along the face's horizontal axis (mm), measured from the back edge. For vertical faces, this runs left-right; for base/lid, it's the Y coordinate. |
| `from_left` | ✓ | Distance along the face's vertical or depth axis (mm), measured from the left edge. For vertical faces, this is the Z height; for base/lid, it's the X coordinate. |
| `width`, `height` | ✓ | Cutout dimensions (mm). Depending on `shape`, one or both may be used. Always specify both; the emitter will zero unused dimensions. |
| `radius` | ✓ | For circular or rounded shapes (mm). Ignored by pure rectangles. |
| `shape` | ✓ | Cutout profile: `rectangle`, `circle`, `rounded_rect`, `circle_with_flats`, `circle_with_key`. See shape guide below. |
| `depth` | optional | How far the cutout penetrates (mm). Defaults to the thickness of the target face. Reduce if you want a partial recess instead of a through-hole. |
| `angle` | optional | Rotation in degrees. Useful for angled USB ports or displays. |
| `mask` | optional | Reserved for future preset patterns (hex grids for ventilation, arrow shapes, etc.). |
| `polygon` | optional | Custom polygon definition for advanced shapes. Planned feature. |

**Shape reference:**
- `rectangle` – Simple rectangular hole. Dimensions controlled by `width` and `height`.
- `circle` – Round hole. Size controlled by `radius`; `width`/`height` ignored.
- `rounded_rect` – Rectangle with rounded corners. `radius` sets corner curvature; `width` and `height` set overall bounds.
- `circle_with_flats` – Circle with flattened top/bottom edges (like a D-sub connector profile). 
- `circle_with_key` – Circle with a keyway notch for anti-rotation (e.g., locking barrel jacks).

**Positioning tips:**
- `from_back` and `from_left` reference the **center** of the cutout, not an edge.
- For a USB-C port `15mm` wide and `8mm` tall, centered `10mm` from the left edge of the back face and `5mm` up from the base, use:
  ```yaml
  face: back
  from_left: 10.0
  from_back: 5.0
  width: 15.0
  height: 8.0
  shape: rounded_rect
  radius: 1.5
  ```
- Measure connector positions on your PCB, then add `enclosure.wall.clearance` and `enclosure.wall.thickness` to translate PCB coordinates to face coordinates.

**Example: USB-C and power jack cutouts**

```yaml
features:
  cutouts:
    # USB-C port on back face
    - face: back
      from_left: 20.0      # 20mm from left edge
      from_back: 8.0       # 8mm up from bottom
      width: 9.0           # Standard USB-C width
      height: 3.5          # Standard USB-C height
      radius: 0.5          # Slight corner rounding
      shape: rounded_rect
    
    # Barrel jack on left face
    - face: left
      from_left: 15.0      # 15mm from front
      from_back: 10.0      # 10mm up from bottom
      width: 8.0
      height: 8.0
      radius: 4.0          # Circle for round jack
      shape: circle
    
    # Ventilation slot on lid
    - face: lid
      from_left: 30.0
      from_back: 20.0
      width: 25.0
      height: 3.0
      radius: 1.5
      shape: rounded_rect
```

These cutouts provide access for power and data while allowing passive cooling through the lid slot.

#### Cutouts: face-wise coordinate mapping and new options

- Coordinate mapping by face:
  - front/back: `from_back` → posy (left↔right), `from_left` → posz (height from base)
  - left/right: `from_back` → posx (front↔back), `from_left` → posz (height from base)
  - base/lid: `from_back` → posx (from back edge), `from_left` → posy (from left edge)
- New optional `pos_z` (side faces only): specify vertical position explicitly. If provided, it overrides `from_left` for `front/back/left/right`.
- New `shape: polygon` support with presets:
  - Presets: `hexagon`, `arrow`, `6pt_star`, `iso_triangle`, `iso_triangle2`, `triangle`, `triangle2`
  - Example:
    ```yaml
    face: base
    from_back: 25
    from_left: 20
    width: 20
    length: 20
    shape: polygon
    polygon: hexagon
    ```

#### Faces and axes 101 (what are posx/posy/posz and why do they map differently?)

The underlying YAPP SCAD generator works in a 3D box coordinate system with three axes:
- posx: left ↔ right (length direction)
- posy: back ↔ front (width direction)
- posz: base ↔ lid (height direction)

When you cut a hole on a particular face, that face has its own local “horizontal” and “vertical” directions. YAPP reuses the same two numbers for all faces, but internally remaps them onto the global axes so the cut ends up on the correct wall and orientation.

To avoid making you think in YAPP’s internal axis names, the DSL uses face-aware names:
- `from_back` → “how far along the horizontal direction of that face”
- `from_left` → “how far along the vertical direction of that face”

Because the face rotates relative to the global axes, the same two numbers land on different global axes:
- On front/back faces, “horizontal” is posy and “vertical” is posz.
- On left/right faces, “horizontal” is posx and “vertical” is posz.
- On base/lid, “horizontal” is posx and “vertical” is posy.

This is why the table above looks “weird”: it’s simply documenting how each face’s local directions map onto the global axes.

Why `pos_z` as an override?
- On vertical faces (front/back/left/right) the second coordinate controls height. Historically we called it `from_left` (matching base/lid usage), which can be confusing because it actually means “height from the base” on those faces.
- We introduced optional `pos_z` so you can write the vertical height explicitly when working on side faces. If `pos_z` is set, it takes precedence over `from_left` on side faces. This keeps backward compatibility while making intent obvious.
- Could we have `from_bottom` / `from_top`? Yes—these aliases are reasonable ergonomically. For now we keep the schema stable and offer `pos_z` as a clear, explicit option. If you want these aliases, open a ticket and we can add them as synonyms.

What do “front/back/left/right” mean?
- front: the wall you’re “looking at”
- back: the opposite wall
- left/right: the side walls when looking at the front
- base: the bottom plate
- lid (top): the top plate

Quick examples:
- Front face (center-left window 11 mm above base):
  ```yaml
  face: front
  from_back: 15   # along the wall, left↔right
  pos_z: 11       # 11 mm up from base (overrides from_left)
  width: 24
  length: 8
  shape: rectangle
  ```
- Left face (near top):
  ```yaml
  face: left
  from_back: 18   # along the wall, front↔back
  pos_z: 18       # height from base
  width: 18
  length: 8
  shape: rectangle
  ```
- Lid (top plane):
  ```yaml
  face: lid
  from_back: 30   # X from back edge
  from_left: 25   # Y from left edge
  radius: 6
  shape: circle
  ```

### `light_tubes`

Light tubes guide LED light from the PCB through the lid. Define tube geometry and the generator creates the internal tunnel and lid opening (respecting `lens_thickness`).

| Field | Required | Description |
|-------|----------|-------------|
| `x`, `y` | ✓ | Position on the PCB (mm), relative to PCB origin. |
| `tube_length` | ✓ | Tube length (mm). |
| `tube_width` | ✓ | Tube width/diameter (mm). |
| `tube_wall` | ✓ | Wall thickness of the tube (mm). |
| `gap_above_pcb` | ✓ | Gap between PCB top and tube start (mm). |
| `shape` | ✓ | `circle` or `rectangle`. |
| `lens_thickness` | optional | Material left in the lid above the tube (mm). `0` = open hole. |
| `height` | optional | Height to the top of the PCB (mm). Defaults to `standoffHeight + pcbThickness`. |
| `fillet_radius` | optional | Fillet radius at the tube base (mm). |
| `coordinate` | optional | Coordinate system: `pcb` (default), `box`, `box_inside`. |
| `origin` | optional | Origin mode: `global` (default) or `alt`. |
| `no_fillet` | optional | If `true`, disables auto fillets. |
| `pcb_name` | optional | For multi-board projects, target a named PCB. |

**Example: circle + rectangle with lens**
```yaml
features:
  light_tubes:
    - x: 15
      y: 10
      tube_length: 5
      tube_width: 6
      tube_wall: 1
      gap_above_pcb: 0.1
      shape: circle

    - x: 15
      y: 30
      tube_length: 1.5
      tube_width: 5
      tube_wall: 1
      gap_above_pcb: 0.1
      shape: rectangle
      lens_thickness: 0.5
```

### `connectors`

Connectors are specialized standoffs designed to anchor the enclosure to an external surface (e.g., wall-mounting an IoT sensor box, bolting a controller to a machine frame). They differ from regular `pcb_stands` in that they include provisions for screw heads, optional countersinking, and configurable insert depths. Each connector creates a cylindrical boss on the enclosure exterior with a through-hole for a machine screw, plus an internal sleeve to maintain structural integrity.

| Field | Required | Description |
|-------|----------|-------------|
| `x`, `y` | ✓ | Position on the PCB plane (mm), measured from the board origin. Typically placed near corners or structural strong-points. |
| `stand_height` | ✓ | Height of the connector column from the enclosure base to the PCB bottom (mm). Usually matches `pcb.z_clearance` unless this connector needs to rise higher. |
| `screw_d` | ✓ | Diameter of the mounting screw shaft (mm). Common values: `3.0` for M3, `4.0` for M4. |
| `screw_head_d` | ✓ | Diameter of the screw head (mm). Determines the countersink size. Measure your actual hardware; M3 heads are typically `5.5–6.0 mm`. |
| `insert_d` | ✓ | Inner diameter of the connector sleeve where the screw passes through (mm). Should be slightly larger than `screw_d` to account for print tolerances (e.g., `3.2` for M3 screws). |
| `outside_d` | ✓ | Outer diameter of the connector boss (mm). Must provide enough wall thickness around `insert_d` for strength (typical: `insert_d + 3–4 mm`). |
| `insert_depth` | optional | How far the insert sleeve penetrates into the enclosure (mm). Defaults to base thickness. Increase for very thick bases or to reinforce attachment points. |
| `pcb_gap` | optional | Extra vertical spacing (mm). Use if this connector needs to lift the board higher than `stand_height`. |
| `fillet_radius` | optional | Fillet radius where the boss meets the enclosure surface (mm). Improves printability and distributes stress. |
| `corner` | optional | Named corner position: `front_left`, `front_right`, `back_left`, `back_right`. Automatically computes `x`, `y` from enclosure dimensions, overriding manual coordinates. |
| `countersink` | optional | Boolean. If `true`, creates a conical recess for flush-mount screws. Depth and angle match `screw_head_d`. |
| `coordinate`, `origin` | optional | Advanced coordinate system overrides. Planned for future DSL exposure to match YAPP's `yappCoord*` flags. |

**Design considerations:**
- Connectors add significant mass to the enclosure corners or edges. Ensure your print bed adhesion is strong.
- For wall-mounting, place connectors symmetrically (e.g., two top corners or four corners) to distribute load.
- If the enclosure will carry weight, use at least M3 screws and ensure `outside_d` provides a robust boss (≥ 7mm for M3).

**Example: Four-corner wall-mount connectors**

```yaml
features:
  connectors:
    # Front-left corner
    - x: 5.0
      y: 5.0
      stand_height: 2.0
      screw_d: 3.0
      screw_head_d: 6.0
      insert_d: 3.5       # M3 + tolerance
      outside_d: 8.0      # Robust boss
      fillet_radius: 2.0
      countersink: true
    
    # Front-right corner
    - x: 60.0
      y: 5.0
      stand_height: 2.0
      screw_d: 3.0
      screw_head_d: 6.0
      insert_d: 3.5
      outside_d: 8.0
      fillet_radius: 2.0
      countersink: true
    
    # Back-left corner
    - x: 5.0
      y: 40.0
      stand_height: 2.0
      screw_d: 3.0
      screw_head_d: 6.0
      insert_d: 3.5
      outside_d: 8.0
      fillet_radius: 2.0
      countersink: true
    
    # Back-right corner
    - x: 60.0
      y: 40.0
      stand_height: 2.0
      screw_d: 3.0
      screw_head_d: 6.0
      insert_d: 3.5
      outside_d: 8.0
      fillet_radius: 2.0
      countersink: true
```

This configuration creates a box that can be flush-mounted to a wall with four M3 screws, with countersunk heads sitting below the enclosure surface for a clean look.

### `snap_joins`

Snap joins are mechanical clips that hold the base and lid together without screws. They work by creating flexible cantilever tabs on one half that snap into matching pockets on the other half. Snap joins are ideal for prototypes, frequently opened enclosures (battery compartments), or designs where you want a clean exterior without visible fasteners. The trade-off is reduced structural rigidity compared to screwed assemblies.

| Field | Required | Description |
|-------|----------|-------------|
| `pos` | ✓ | Distance along the chosen edge (mm), measured from the edge's origin. For `front`/`back` sides, this is the X coordinate; for `left`/`right`, it's the Y coordinate. |
| `width` | ✓ | Width of the snap tab (mm). Wider tabs are stronger but harder to release. Typical range: `10–20 mm`. |
| `side` | ✓ | Which edge hosts the snap: `front`, `back`, `left`, or `right`. |

**How they work:**
- The generator creates a flexible tab (typically on the lid) that deflects inward during assembly.
- A corresponding pocket in the mating part captures the tab's locking feature.
- To open, you flex the tab outward until it clears the pocket.

**Placement strategy:**
- Use at least two snap joins per enclosure for stability (one on opposite sides).
- Distribute them symmetrically to avoid twisting during assembly.
- Avoid placing snaps too close to corners where wall stiffness is higher; aim for mid-span positions.
- For long edges (> 80mm), consider adding a third snap joint for even clamping force.

**Tuning tips:**
- If snaps are too hard to engage, reduce `width` or increase `tolerances.perimeter`.
- If the lid pops open easily, increase `width` or reduce tolerance.
- Print orientation matters: tabs should flex perpendicular to layer lines for maximum strength.

**Example: Front and back snap joins**

```yaml
features:
  snap_joins:
    # Front edge, centered
    - pos: 32.5        # Half of PCB length (65 / 2)
      width: 15.0
      side: front
    
    # Back edge, centered
    - pos: 32.5
      width: 15.0
      side: back
```

This adds two opposing snaps that secure the lid with firm but releasable clips, allowing tool-free access to the PCB for battery swaps or debugging.

**Example: Four-corner snap arrangement**

```yaml
features:
  snap_joins:
    - pos: 10.0
      width: 12.0
      side: front
    
    - pos: 55.0
      width: 12.0
      side: front
    
    - pos: 10.0
      width: 12.0
      side: back
    
    - pos: 55.0
      width: 12.0
      side: back
```

This distributes four snaps for maximum stability on a longer enclosure.

### `push_buttons`

**Status:** Available. The DSL now mirrors the entire `pushButtons` array from `YAPPgenerator_v3.scad`, including polygon presets, coordinate/origin flags, and optional tolerances. When at least one entry exists, `yappctl generate` automatically sets `printSwitchExtenders = true` in the emitted SCAD *and* passes `-D printSwitchExtenders=true` to OpenSCAD so STLs include the extenders.

**Top-level fields:**

| Field | Type | Notes |
|-------|------|-------|
| `name` | optional string | Free-form label used in errors and documentation. |
| `x`, `y` | required number | Button position relative to the PCB origin (mm). Expressions (e.g., `pcb.width / 2 - 20`) are encouraged. |
| `shape` | optional enum | `rectangle` (default), `circle`, `rounded_rect`, `circle_with_flats`, `circle_with_key`, or polygon presets (`arrow`, `triangle`, `hexagon`, `star6`, etc.). |
| `polygon` / `polygon_preset` / `shape_preset` | optional string | Required when `shape: polygon`. Accepts friendly names (`arrow`) or raw tokens (`shapeArrow`). |
| `angle` | optional number | Rotation in degrees (default `0`). |
| `fillet_radius` | optional number | Internal fillet at the button shaft. `lid.no_fillet: true` appends `yappNoFillet`. |
| `coordinate` | optional enum | `pcb`, `box`, or `box_inside`. Maps to `yappCoord*` flags. |
| `origin` | optional enum | `global` (default) or `left`, mapping to `yappGlobalOrigin` / `yappLeftOrigin`. |

**Nested blocks:**

| Block | Fields | Description |
|-------|--------|-------------|
| `cap` | `length`, `width`, `radius` (all required) | Defines the lid opening/cap footprint. For circles, set `length`/`width` to the diameter and `radius` to half of that value. |
| `lid` | `protrusion` (required), `wall`, `plate_thickness`, `slack`, `snap_slack`, `no_fillet` | Controls how the extender cap interfaces with the lid. Slack settings mirror SCAD’s `buttonSlack` and `snapSlack`. |
| `switch` | `height`, `travel`, `pole_diameter` (required), `top_height` (optional) | Matches tactile switch dimensions. `top_height` overrides the default `standoffHeight + pcbThickness` reference plane. |

**Polygon presets:**

- Friendly names: `arrow`, `triangle`, `triangle2`, `iso_triangle`, `iso_triangle2`, `hexagon`, `star6`.
- Raw tokens: `shapeArrow`, `shapeTriangle`, `shapeTriangle2`, `shapeIsoTriangle`, `shapeIsoTriangle2`, `shapeHexagon`, `shape6ptStar`.
- Custom SCAD tokens that start with `shape` are passed through unchanged for advanced use.

**Automatic behavior:**

- The generator converts missing optional fields to `undef`, preserving YAPP defaults.
- Each entry appends shape/coordinate/origin flags in the order the SCAD expects, so hand-tuning is unnecessary.
- `RenderSTLs` now forwards the `printSwitchExtenders` flag directly to OpenSCAD, keeping the SCAD preview and STLs in sync.

**Example (from `examples/yapp-demo-buttons2.yaml`):**

```yaml
features:
  push_buttons:
    - name: arrow-rotated
      x: 40
      y: pcb.width / 2 - 20
      shape: polygon
      polygon: arrow
      angle: 90
      cap:
        length: 8
        width: 8
        radius: 4
      lid:
        protrusion: 0
      switch:
        height: 5
        travel: 0.5
        pole_diameter: 3

    - name: circle-offset
      x: 20
      y: pcb.width / 2 - 20
      shape: circle
      cap:
        length: 8
        width: 8
        radius: 4
      lid:
        protrusion: 0
        wall: 2.0
        plate_thickness: 2.5
        slack: 0.5
      switch:
        height: 5
        travel: 0.5
        pole_diameter: 3
```

For a simpler two-button lid, study `examples/yapp-demo-buttons.yaml` (matches `examples/YAPP_Demo_buttons_v31.scad`). For a shape buffet (polygon presets, rounded rectangles, and multiple tolerances), see `examples/yapp-demo-buttons2.yaml` (matches `examples/YAPP_Demo_buttons2_v31.scad`).

> **Tip:** Use `yappctl resolve --strict` before generating to catch typos like `polygon: arrrow` or missing `lid.protrusion`—the resolver reports the offending `features.push_buttons[i]` entry by name if provided.


## Tolerances section

3D printing introduces dimensional variations due to nozzle diameter, material shrinkage, temperature settings, and printer calibration. The `tolerances` block lets you compensate globally for these factors so holes aren't too tight, parts don't bind, and mechanical features work reliably. Think of tolerance values as safety margins: larger tolerances make assembly easier but reduce precision; smaller tolerances demand better printer calibration but yield tighter fits.

| Field | Default | Description |
|-------|---------|-------------|
| `holes` | `0.4` | Added to all hole diameters (mm). Applied to standoff pins, connector inserts, and screw clearances. Increase if screws or inserts won't fit; decrease if they're too loose. Typical range: `0.3–0.6`. |
| `perimeter` | `1.0` | Expands shell outer dimensions (mm). Creates breathing room between the PCB edge and enclosure walls, preventing the board from jamming during insertion. Increase for boards with rough-cut edges or large component overhangs. Typical range: `0.5–2.0`. |

**Planned expansion:**
Future DSL versions will add feature-specific tolerance keys:
- `buttons` – Slack adjustment for button shafts and extenders.
- `snap_joins` – Fine-tune snap tab engagement.
- `ridge` – Independent control over ridge/groove fit beyond `enclosure.ridge.slack`.

**Tuning strategy:**

1. **Start with defaults** (`holes: 0.4`, `perimeter: 1.0`) for your first print.
2. **Test critical features:** Insert a screw into a standoff hole. Does it thread smoothly? If too tight, increase `holes` by `0.1 mm`. If rattling, decrease by `0.1 mm`.
3. **Check PCB fit:** Slide the board into the printed base. It should drop in easily without forcing. If it binds, increase `perimeter`.
4. **Document your printer profile:** Once dialed in, keep a YAML snippet with your tested tolerance values for reuse across projects.

**Printer-specific considerations:**
- **FDM (Fused Deposition Modeling):** Standard PLA/PETG benefits from `holes: 0.4`. High-temp materials (ABS, Nylon) may need `0.5+` due to greater shrinkage.
- **Resin (SLA/DLP):** Resin prints are dimensionally accurate but can stick in tight clearances. Start with `holes: 0.2–0.3`.
- **Nozzle size:** 0.4mm nozzles work well with default tolerances. Larger nozzles (0.6mm+) may need increased values.

**Example: Tight-tolerance resin build**

```yaml
tolerances:
  holes: 0.25        # Resin accuracy allows smaller clearance
  perimeter: 0.8     # Board fit is snug but manageable
```

**Example: Forgiving FDM build for quick prototyping**

```yaml
tolerances:
  holes: 0.5         # Ensures screws always fit, even with mediocre calibration
  perimeter: 1.5     # Extra space for rough PCB edges
```

By tuning these two parameters, you can adapt a single design to different printers, materials, and quality requirements without modifying the core geometry.

## The resolver/generator pipeline

Understanding the transformation stages from YAML to 3D mesh helps you debug issues efficiently and leverage the toolchain's full power. The pipeline has four distinct phases, each with specific responsibilities. When something goes wrong—a cutout in the wrong place, a missing standoff, an expression error—knowing which phase to inspect saves time.

### Phase 1: Resolve

**Command:** `yappctl resolve [--input your.yaml] [--output resolved.yaml] [--strict]`

**Purpose:** Transform the human-friendly YAML (with expressions, variables, references) into a fully evaluated, numeric-only document where every field has a concrete value.

**What happens:**
- Load the input YAML.
- Hoist `vars.*` entries into the global scope.
- Iterate up to 16 times, evaluating expressions like `enclosure.wall.thickness + 2` or `max(vars.a, vars.b)` until all values are resolved to numbers.
- Fill missing optional fields with defaults (e.g., `pcb.z_clearance: 1.0`, `tolerances.holes: 0.4`).
- Validate required fields exist (e.g., `pcb.length`, `enclosure.wall.thickness`).
- Detect circular dependencies and expression errors.

**Output:** A fully resolved YAML file (or JSON with `--format json`) ready for the next stage.

**Use cases:**
- **Preview final values:** See what `vars.computed_offset` resolves to before generating geometry.
- **Catch errors early:** Missing required fields or typos in expressions fail here, not during STL rendering.
- **Generate documentation:** The resolved output can serve as a "flattened" spec for external tools.

**Example:**

```bash
go run ./cmd/yappctl resolve --input mybox.yaml --output resolved.yaml --strict
```

The `--strict` flag enables extra validation: unknown keys error out, unused `vars` are flagged, and expression warnings become fatal.

### Phase 2: Build model

**Internal API:** `pkg/yappgen.BuildModel(resolvedData)`

**Purpose:** Parse the resolved YAML into a strongly-typed Go struct (`yappgen.Model`) that organizes data into convenient in-memory representations.

**What happens:**
- Extract top-level metadata (`project`, `version`, `units`).
- Parse `pcb` and `enclosure` blocks into nested structs.
- Transform feature arrays (`pcb_stands`, `cutouts`, `connectors`, `snap_joins`) from YAML maps into typed Go structs.
- Apply coordinate transformations (e.g., convert face-relative cutout positions into 3D SCAD coordinates).
- Normalize shape enums (`rectangle` → `yappRectangle`).

**Output:** A `Model` object ready for SCAD code generation.

**Error scenarios:**
- Invalid shape name (e.g., `shape: oval` instead of `circle`).
- Out-of-range coordinates (cutout positioned outside the enclosure bounds).
- Conflicting options (e.g., a connector with `corner: front_left` but also explicit `x`, `y`).

### Phase 3: Emit SCAD

**Internal API:** `yappgen.EmitSCAD(model)`

**Purpose:** Translate the `Model` into OpenSCAD source code compatible with `YAPPgenerator_v3.scad`.

**What happens:**
- Rewrite the `include <...>` path to point to the local YAPP generator script.
- Emit global variables: `pcbLength`, `pcbWidth`, `wallThickness`, etc.
- Convert feature structs into positional SCAD arrays:
  ```scad
  pcbStands = [
    [3.0, 3.0, 6.0, 2.5, 0.4, 1.5],
    [62.0, 3.0, 6.0, 2.5, 0.4, 1.5],
    // ...
  ];
  ```
- Append YAPP flag tokens (e.g., `yappRectangle`, `yappFront`, `yappCoordBox`).
- Set `printSwitchExtenders = true` if `push_buttons` are present.

**Output:** A `.scad` file ready to open in OpenSCAD or render headlessly.

**Debugging tip:** If a cutout or standoff appears in the wrong location, open the generated `.scad` file and inspect the raw arrays. This tells you whether the problem is in the DSL→SCAD translation or the SCAD→geometry rendering.

### Phase 4: Render STLs (optional)

**Command:** `yappctl generate --input mybox.yaml --stl-base --stl-lid --scad-out mybox.scad`

**Purpose:** Invoke OpenSCAD's command-line interface to generate printable mesh files (`.stl`).

**What happens:**
- Run phases 1–3 to produce the `.scad` file.
- Call `openscad` with appropriate `print*` flags:
  - `--stl-base` → `-D printBaseShell=true`
  - `--stl-lid` → `-D printLidShell=true`
  - `--stl-base-insert` → `-D printBaseInserts=true` (optional insert plate for heat-set inserts)
- Wait for OpenSCAD to complete rendering (can take 30 seconds to several minutes depending on geometry complexity).
- Write `.stl` files to the output directory.

**Use cases:**
- **Direct-to-print workflow:** `yappctl generate` → slice → print, no OpenSCAD GUI needed.
- **CI/CD integration:** Automate STL generation for every YAML commit.
- **Batch processing:** Generate base and lid in one command.

**Example:**

```bash
go run ./cmd/yappctl generate --input mybox.yaml \
  --scad-out output/mybox.scad \
  --stl-base output/mybox-base.stl \
  --stl-lid output/mybox-lid.stl
```

### Debugging flowchart

When a feature doesn't work as expected:

1. **Run `resolve`:** Does the YAML resolve without errors? Are all required fields present? Do expressions evaluate to sensible numbers?
   - If no: Fix YAML syntax, missing fields, or expression errors.
2. **Inspect resolved output:** Open `resolved.yaml`. Are the final numeric values what you expect?
   - If no: Check expression logic, `vars` definitions, or default values.
3. **Generate SCAD:** Run `generate --scad-out`. Open the `.scad` file and find the feature array (e.g., `cutOuts`, `pcbStands`). Are the positional parameters correct?
   - If no: Bug in the DSL-to-SCAD mapper; report it.
4. **Render in OpenSCAD GUI:** Open the `.scad` file in OpenSCAD. Does the feature appear in the 3D view?
   - If no: YAPP generator bug or OpenSCAD version incompatibility; check YAPP release notes.

This layered approach isolates problems quickly: YAML syntax → expression evaluation → model construction → SCAD generation → geometry rendering.

## Helpful commands

A quick reference for common `yappctl` workflows. All commands assume you're in the project root and have Go installed.

### Basic resolution and validation

```bash
# Resolve expressions and fill defaults
go run ./cmd/yappctl resolve --input mybox.yaml --output resolved.yaml

# Strict mode: catch typos and unused vars
go run ./cmd/yappctl resolve --input mybox.yaml --strict

# Output as JSON for external tools
go run ./cmd/yappctl resolve --input mybox.yaml --format json > mybox.json
```

### SCAD generation

```bash
# Generate OpenSCAD file only (no STL rendering)
go run ./cmd/yappctl generate --input mybox.yaml --scad-out mybox.scad

# Generate and open in OpenSCAD GUI (if installed)
go run ./cmd/yappctl generate --input mybox.yaml --scad-out mybox.scad
openscad mybox.scad &
```

### STL rendering

```bash
# Generate base STL only
go run ./cmd/yappctl generate --input mybox.yaml --stl-base output/base.stl

# Generate lid STL only
go run ./cmd/yappctl generate --input mybox.yaml --stl-lid output/lid.stl

# Generate both base and lid
go run ./cmd/yappctl generate --input mybox.yaml \
  --stl-base output/base.stl \
  --stl-lid output/lid.stl

# Generate everything including insert plate and keep SCAD source
go run ./cmd/yappctl generate --input mybox.yaml \
  --scad-out output/mybox.scad \
  --stl-base output/base.stl \
  --stl-lid output/lid.stl \
  --stl-base-insert output/insert-plate.stl
```

### Documentation and help

```bash
# Show getting started tutorial
go run ./cmd/yappctl help yapp-dsl-getting-started

# Show this reference document
go run ./cmd/yappctl help yapp-dsl-reference

# List all available help topics
go run ./cmd/yappctl help
```

## Complete example

This section shows a full, production-ready YAML file that combines all the concepts from this reference. It defines an IoT sensor box with standoffs, connectors, cutouts for USB-C and a barrel jack, ventilation, snap joins, and printer-tuned tolerances.

```yaml
---
# Project metadata
project: IoT Sensor Box
version: 2
units: mm
yapp_version: v3.3.8

# Reusable constants
vars:
  # Derived dimensions
  mid_length: pcb.length / 2
  mid_width: pcb.width / 2
  
  # Mounting hardware
  mount_screw_d: 2.5
  standoff_diameter: 6.0
  
  # Printer tuning (calibrated for my Prusa i3)
  hole_tolerance: 0.4
  shell_expansion: 1.2

# PCB specifications
pcb:
  length: 65.0
  width: 45.0
  thickness: 1.6
  z_clearance: 2.0
  
  standoffs:
    diameter: vars.standoff_diameter
    screw_d: vars.mount_screw_d

# Enclosure structure
enclosure:
  wall:
    thickness: 2.4
    clearance: 1.5
    fillet_radius: 2.0
  
  base:
    thickness: 1.8
  
  lid:
    thickness: 1.6
  
  ridge:
    height: 5.0
    slack: 0.25
    gap: 0.5

# Features: standoffs, cutouts, connectors, snap joins
features:
  # PCB mounting standoffs (four corners)
  pcb_stands:
    - x: 3.5
      y: 3.5
      diameter: vars.standoff_diameter
      pin_diameter: vars.mount_screw_d
      fillet_radius: 1.5
    
    - x: 61.5
      y: 3.5
      diameter: vars.standoff_diameter
      pin_diameter: vars.mount_screw_d
      fillet_radius: 1.5
    
    - x: 61.5
      y: 41.5
      diameter: vars.standoff_diameter
      pin_diameter: vars.mount_screw_d
      fillet_radius: 1.5
    
    - x: 3.5
      y: 41.5
      diameter: vars.standoff_diameter
      pin_diameter: vars.mount_screw_d
      fillet_radius: 1.5
  
  # Cutouts for connectors and ventilation
  cutouts:
    # USB-C port on back face
    - face: back
      from_left: 20.0
      from_back: 8.0
      width: 9.0
      height: 3.5
      radius: 0.5
      shape: rounded_rect
    
    # Barrel jack (power input) on left face
    - face: left
      from_left: 15.0
      from_back: 10.0
      width: 8.0
      height: 8.0
      radius: 4.0
      shape: circle
    
    # Status LED window on front face
    - face: front
      from_left: vars.mid_length
      from_back: 12.0
      width: 5.0
      height: 5.0
      radius: 2.5
      shape: circle
    
    # Ventilation slots on lid (three slots)
    - face: lid
      from_left: 15.0
      from_back: 22.5
      width: 25.0
      height: 3.0
      radius: 1.5
      shape: rounded_rect
    
    - face: lid
      from_left: 15.0
      from_back: 28.0
      width: 25.0
      height: 3.0
      radius: 1.5
      shape: rounded_rect
    
    - face: lid
      from_left: 15.0
      from_back: 33.5
      width: 25.0
      height: 3.0
      radius: 1.5
      shape: rounded_rect
  
  # Wall-mount connectors (two top corners)
  connectors:
    - x: 5.0
      y: 5.0
      stand_height: 2.0
      screw_d: 3.0
      screw_head_d: 6.0
      insert_d: 3.5
      outside_d: 8.0
      fillet_radius: 2.0
      countersink: true
    
    - x: 60.0
      y: 5.0
      stand_height: 2.0
      screw_d: 3.0
      screw_head_d: 6.0
      insert_d: 3.5
      outside_d: 8.0
      fillet_radius: 2.0
      countersink: true
  
  # Snap joins for tool-free assembly
  snap_joins:
    - pos: vars.mid_length
      width: 15.0
      side: front
    
    - pos: vars.mid_length
      width: 15.0
      side: back

# Printer-specific tolerances
tolerances:
  holes: vars.hole_tolerance
  perimeter: vars.shell_expansion

# Documentation (ignored by generator)
notes: |
  This enclosure houses an ESP32-based environmental sensor.
  PCB revision: v2.1
  Assembly: Insert M2.5 heat-set inserts into standoffs, secure PCB with M2.5×8mm screws.
  Wall mounting: Use M3×20mm screws through top connectors.
```

**What this example demonstrates:**
- **Variables:** `vars` block computes midpoints and standardizes screw sizes.
- **Expression evaluation:** `mid_length: pcb.length / 2` keeps geometry centered automatically.
- **Multiple feature types:** Standoffs, cutouts, connectors, and snap joins in one design.
- **Practical details:** Ventilation slots, LED window, printer-specific tolerances.
- **Documentation:** `notes` field tracks assembly instructions without cluttering the generator.

To use this example:

```bash
# Save as iot-sensor-box.yaml
go run ./cmd/yappctl generate --input iot-sensor-box.yaml \
  --scad-out output/iot-sensor-box.scad \
  --stl-base output/base.stl \
  --stl-lid output/lid.stl
```

## Additional resources

### Example files

The `examples/` directory contains YAML conversions of the classic YAPP SCAD demos:
- `YAPP_Demo_buttons_v31.yaml` – Push button example (planned feature)
- `YAPP_Demo_cutouts.yaml` – Comprehensive cutout shapes
- `YAPP_Demo_standoffs.yaml` – PCB mounting techniques

These are excellent starting points. Copy one, modify dimensions and features, and generate your own design.

### Design documentation

- **Ticket `ttmp/YAPP-ENCL-DSL-001-*`** – Original DSL design rationale, schema debates, and implementation notes.
- **Ticket `ttmp/YAPP-PUSH-BUTTONS-001-*`** – Push button feature design in progress.
- **Getting started tutorial:** `go run ./cmd/yappctl help yapp-dsl-getting-started` – Narrative walkthrough with step-by-step examples.

### YAPP upstream resources

- **YAPP GitHub repository:** [github.com/mrWheel/YAPP_Box](https://github.com/mrWheel/YAPP_Box)
- **YAPP documentation:** `YAPPgenerator_v3.scad` header comments contain detailed parameter explanations.
- **Community forum:** Discussions, tips, and shared designs.

### Troubleshooting

**Problem:** "Required field missing" error during resolve.
- **Solution:** Check that all required fields in `pcb` and `enclosure` sections are present. Use `--strict` mode to catch typos.

**Problem:** Cutout appears on wrong face or wrong position.
- **Solution:** Run `resolve` first, check the numeric coordinates. Verify `from_back`/`from_left` are measured from the correct edge for that face. Inspect the generated SCAD `cutOuts` array.

**Problem:** Standoff holes too tight for screws.
- **Solution:** Increase `tolerances.holes` by `0.1 mm` increments. Alternatively, override `hole_slack` for specific standoffs.

**Problem:** PCB won't fit into printed enclosure.
- **Solution:** Increase `enclosure.wall.clearance` or `tolerances.perimeter`. Verify PCB measurements are accurate (use calipers).

**Problem:** OpenSCAD render fails or produces strange geometry.
- **Solution:** Check `yapp_version` matches your local YAPP generator version. Update YAPP or adjust the version field. Inspect the generated `.scad` file for malformed arrays.

For additional support, consult the project's issue tracker or the YAPP community forums.
