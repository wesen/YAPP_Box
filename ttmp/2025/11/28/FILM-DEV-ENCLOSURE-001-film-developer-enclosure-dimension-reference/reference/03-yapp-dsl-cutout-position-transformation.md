---
Title: YAPP DSL Cutout Position Transformation
Ticket: FILM-DEV-ENCLOSURE-001
Status: active
Topics:
    - film-developer
    - enclosure
    - yapp
    - dsl
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Explains how YAPP DSL face-relative coordinates transform into SCAD array positions, and how the origin flag affects positioning
LastUpdated: "2025-11-28T20:00:00.000000000-05:00"
---

# YAPP DSL Cutout Position Transformation

## Goal

Explain how the YAPP DSL transforms face-relative coordinates (`from_face_left`, `from_face_bottom`, `from_face_back`) into SCAD array positions (`cutOut[0]`, `cutOut[1]`), and how the `origin: center` flag affects the final positioning. This complements the SCAD-side analysis in `02-cutout-offset-computation.md` by documenting the DSL-to-SCAD transformation layer.

## Context

The YAPP DSL provides face-relative coordinate fields that are intuitive for users (e.g., "10mm from the left edge of the back face"), but the underlying SCAD generator expects positional arrays with global coordinate system values. The `pkg/yappgen/modules/cutouts` module performs this transformation, converting DSL fields into SCAD array format before the SCAD generator applies its own coordinate transformations and shape offsets.

## DSL Coordinate Fields

The YAPP DSL uses three face-relative coordinate fields:

| Field | Required For | Description |
|-------|--------------|-------------|
| `from_face_left` | **All faces** | Horizontal position along the face from the left edge (mm). Meaning varies by face type. |
| `from_face_bottom` | **Side faces only** (front, back, left, right) | Vertical position from the bottom edge of the face (mm). Height from base. |
| `from_face_back` | **Horizontal faces only** (base, lid) | Depth position from the back edge of the face (mm). X coordinate in box space. |

**Important:** The DSL documentation may reference `from_back` and `from_left`, but the actual schema uses `from_face_left`, `from_face_bottom`, and `from_face_back`. The field names in the code and schema are authoritative.

## Transformation: DSL Fields → SCAD Array Positions

The transformation happens in `pkg/yappgen/modules/cutouts/module.go` in the `Build` function:

```47:70:pkg/yappgen/modules/cutouts/module.go
		// Determine position values based on face and new field names
		// All faces use from_face_left for horizontal position
		// Side faces use from_face_bottom for vertical position
		// Base/lid use from_face_back for depth position
		faceLower := strings.ToLower(strings.TrimSpace(item.Face))
		isSideFace := faceLower == "front" || faceLower == "back" || faceLower == "left" || faceLower == "right"
		isHorizFace := faceLower == "base" || faceLower == "lid" || faceLower == "top" || faceLower == "bottom"

		var pos0, pos1 float64

		if isSideFace {
			// Side faces: from_face_left → horizontal (pos0), from_face_bottom → vertical (pos1)
			pos0 = item.FromFaceLeft
			if item.FromFaceBottom != nil {
				pos1 = *item.FromFaceBottom
			}
		} else if isHorizFace {
			// Horizontal faces: from_face_left → Y, from_face_back → X
			// But YAPP arrays expect [X, Y] order, so we swap
			if item.FromFaceBack != nil {
				pos0 = *item.FromFaceBack // from_face_back becomes pos0 (X/back-to-front)
			}
			pos1 = item.FromFaceLeft // from_face_left becomes pos1 (Y/left-to-right)
		}
```

### Transformation Rules by Face Type

#### Side Faces (front, back, left, right)

**DSL fields:**
- `from_face_left` → **pos0** (horizontal position along the face)
- `from_face_bottom` → **pos1** (vertical position, height from base)

**SCAD array:**
```scad
cutOuts = [
  [pos0, pos1, width, length, radius, shapeFlag, ...]
];
```

**Example:**
```yaml
cutouts:
  - face: back
    from_face_left: 20.0      # → pos0 = 20.0
    from_face_bottom: 10.0    # → pos1 = 10.0
    width: 9.0
    length: 3.5
    shape: rounded_rect
```

**Generated SCAD:**
```scad
cutoutsBack = [
  [20.0, 10.0, 9.0, 3.5, 0.5, yappRoundedRect, ...]
];
```

#### Horizontal Faces (base, lid)

**DSL fields:**
- `from_face_back` → **pos0** (X coordinate, back-to-front)
- `from_face_left` → **pos1** (Y coordinate, left-to-right)

**Note:** The order is swapped because YAPP arrays expect `[X, Y]` order, but the DSL uses `from_face_back` (X) and `from_face_left` (Y).

**SCAD array:**
```scad
cutOuts = [
  [pos0, pos1, width, length, radius, shapeFlag, ...]
];
```

**Example:**
```yaml
cutouts:
  - face: lid
    from_face_left: 30.0      # → pos1 = 30.0 (Y coordinate)
    from_face_back: 20.0      # → pos0 = 20.0 (X coordinate)
    width: 25.0
    length: 3.0
    shape: rounded_rect
```

**Generated SCAD:**
```scad
cutoutsLid = [
  [20.0, 30.0, 25.0, 3.0, 1.5, yappRoundedRect, ...]
];
```

## Complete Transformation Pipeline

The full pipeline from DSL YAML to final cutout position involves multiple stages:

### Stage 1: DSL YAML → SCAD Array (This Document)

**Input:** DSL face-relative coordinates
```yaml
cutouts:
  - face: lid
    from_face_left: 40.0
    from_face_back: 20.0
    radius: 10.0
    shape: circle
```

**Output:** SCAD array positions
```scad
cutoutsLid = [
  [20.0, 40.0, 0, 0, 10.0, yappCircle, ...]
];
```

**Transformation:** `from_face_back` → `pos0`, `from_face_left` → `pos1`

### Stage 2: SCAD Array → Box Coordinates (SCAD Generator)

**Input:** SCAD array `cutOut[0] = 20.0`, `cutOut[1] = 40.0`

**Process:** `processCutoutList_Face` extracts positions and transforms them:
```scad
theX = translate2Box_X(cutOut[0], face, coordSystem);  // 20.0 → transformed X
theY = translate2Box_Y(cutOut[1], face, coordSystem);  // 40.0 → transformed Y
```

**Output:** Box coordinate system positions (`theX`, `theY`)

### Stage 3: Box Coordinates → Shape Position (SCAD Generator)

**Process:** `processCutoutList_Shape` positions the shape:
```scad
translate([pos_X, pos_Y, ...])
```

Where `pos_X = base_pos_H` and `pos_Y = base_pos_V` (derived from `theX` and `theY`).

### Stage 4: Shape Offset Application (SCAD Generator)

**Process:** `generateShape` applies shape-specific offsets:
- **Circle:** `translate([Radius, Radius, 0])` if `useCenter = false`
- **Rectangle:** `translate([Width/2, Length/2, 0])` if `useCenter = false`
- **No offset** if `useCenter = true` (when `origin: center` is set)

**Final position:** Shape center (or corner, depending on `origin` flag)

## The `origin: center` Flag

The `origin` field in the DSL controls whether the position specifies a corner or the center of the cutout.

### DSL Field Mapping

```yaml
cutouts:
  - face: lid
    from_face_left: 50.0
    from_face_back: 30.0
    radius: 10.0
    shape: circle
    origin: center    # Sets yappCenter flag
```

**Transformation in `encodeFlags`:**
```276:287:pkg/yappgen/modules/cutouts/module.go
func originFlag(value string) (scad.Raw, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "global":
		return "", nil // Default, don't emit
	case "center":
		return scad.Raw("yappCenter"), nil
	case "alt":
		return scad.Raw("yappAltOrigin"), nil
	default:
		return "", errors.Errorf("invalid origin: %s", value)
	}
}
```

**Generated SCAD:**
```scad
cutoutsLid = [
  [30.0, 50.0, 0, 0, 10.0, yappCircle, ..., yappCenter]
];
```

### Effect on Positioning

**Without `origin: center` (default):**
- DSL position specifies **corner** of bounding box
- SCAD generator applies shape-specific offset:
  - Circle: `+Radius` in X and Y
  - Rectangle: `+Width/2` in X, `+Length/2` in Y
- **Final center** = DSL position + offset

**With `origin: center`:**
- DSL position specifies **center** directly
- SCAD generator applies **no offset**
- **Final center** = DSL position (exactly)

### Example: Circle Cutout Positioning

**Scenario:** Place a 10mm radius circle centered at (50, 30) on the lid.

**Option 1: Without `origin: center` (corner-based)**
```yaml
cutouts:
  - face: lid
    from_face_left: 40.0      # 50 - 10 (center - radius)
    from_face_back: 20.0      # 30 - 10 (center - radius)
    radius: 10.0
    shape: circle
```

**Transformation:**
1. DSL → SCAD: `[20.0, 40.0, ...]`
2. SCAD → Box coords: `(20.0, 40.0)` (after coordinate transformation)
3. Shape offset: `+Radius` = `+10.0` in both X and Y
4. **Final center:** `(30.0, 50.0)` ✅

**Option 2: With `origin: center` (center-based)**
```yaml
cutouts:
  - face: lid
    from_face_left: 50.0      # Direct center position
    from_face_back: 30.0      # Direct center position
    radius: 10.0
    shape: circle
    origin: center
```

**Transformation:**
1. DSL → SCAD: `[30.0, 50.0, ..., yappCenter]`
2. SCAD → Box coords: `(30.0, 50.0)` (after coordinate transformation)
3. Shape offset: **None** (because `yappCenter` flag is set)
4. **Final center:** `(30.0, 50.0)` ✅

## Coordinate System Flags

The DSL also supports coordinate system selection via the `coordinate` field:

| DSL Value | SCAD Flag | Description |
|-----------|-----------|-------------|
| `pcb` (default) | `yappCoordPCB` | Positions relative to PCB origin |
| `box` | `yappCoordBox` | Positions relative to box outer dimensions |
| `box_inside` | `yappCoordBoxInside` | Positions relative to box inner dimensions |

**Transformation:**
```263:274:pkg/yappgen/modules/cutouts/module.go
func coordinateFlag(value string) (scad.Raw, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "pcb":
		return "", nil // Default, don't emit
	case "box":
		return scad.Raw("yappCoordBox"), nil
	case "box_inside":
		return scad.Raw("yappCoordBoxInside"), nil
	default:
		return "", errors.Errorf("invalid coordinate: %s", value)
	}
}
```

**Example:**
```yaml
cutouts:
  - face: lid
    from_face_left: 50.0
    from_face_back: 30.0
    radius: 10.0
    shape: circle
    coordinate: box_inside    # Use inner box dimensions
```

**Generated SCAD:**
```scad
cutoutsLid = [
  [30.0, 50.0, 0, 0, 10.0, yappCircle, ..., yappCoordBoxInside]
];
```

The coordinate system flag affects how `translate2Box_X` and `translate2Box_Y` transform the positions in Stage 2 of the pipeline.

## Summary Table: DSL → SCAD Transformation

| Face Type | DSL Field 1 | DSL Field 2 | → | SCAD pos0 | SCAD pos1 |
|-----------|-------------|-------------|---|-----------|-----------|
| **Side faces** (front, back, left, right) | `from_face_left` | `from_face_bottom` | → | `from_face_left` | `from_face_bottom` |
| **Horizontal faces** (base, lid) | `from_face_back` | `from_face_left` | → | `from_face_back` | `from_face_left` |

**Notes:**
- For side faces, `from_face_left` is horizontal along the face, `from_face_bottom` is vertical (height).
- For horizontal faces, `from_face_back` is X (back-to-front), `from_face_left` is Y (left-to-right).
- The order matches SCAD array `[pos0, pos1]` format.

## Key Insights

1. **Face-relative coordinates:** The DSL uses intuitive face-relative names (`from_face_left`, `from_face_bottom`, `from_face_back`) that map differently depending on face orientation.

2. **Two-stage transformation:**
   - **DSL → SCAD:** Face-relative coordinates → SCAD array positions (`pos0`, `pos1`)
   - **SCAD → Final:** SCAD positions → Box coordinates → Shape positioning with offsets

3. **Origin flag controls offset:** The `origin: center` flag determines whether the DSL position is interpreted as a corner (default) or center (with flag), affecting whether shape-specific offsets are applied.

4. **Coordinate system selection:** The `coordinate` field selects which coordinate system the positions reference (PCB, box outer, box inner), affecting the transformation in Stage 2.

5. **Shape offsets are SCAD-side:** The shape-specific offsets (Radius for circles, Width/2 for rectangles) are applied by the SCAD generator, not the DSL transformer. The DSL only controls whether offsets are applied via the `origin` flag.

## Related

- `pkg/yappgen/modules/cutouts/module.go` - DSL-to-SCAD transformation code
- `pkg/yappgen/modules/cutouts/schema.yaml` - DSL schema definition
- `reference/02-cutout-offset-computation.md` - SCAD-side offset computation analysis
- `pkg/docs/tutorials/yapp-dsl-reference.md` - Complete DSL reference documentation
