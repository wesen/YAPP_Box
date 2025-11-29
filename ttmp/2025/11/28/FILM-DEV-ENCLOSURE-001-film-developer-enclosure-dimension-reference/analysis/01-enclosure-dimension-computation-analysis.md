---
Title: YAPP Enclosure Dimension Computation - Complete Analysis
Ticket: FILM-DEV-ENCLOSURE-001
Status: active
Topics:
    - yapp
    - enclosure
    - analysis
DocType: analysis
Intent: long-term
Owners: []
RelatedFiles:
    - Path: examples/YAPP_Demo_buttons_v30.scad
      Note: Example SCAD file showing manual dimension configuration
    - Path: pkg/yappgen/assets/YAPPgenerator_v3.scad
      Note: |-
        Main YAPP generator SCAD file containing all dimension computation logic
        Main YAPP generator SCAD file containing dimension computation logic (lines 324-355
    - Path: projects/film-developer/enclosure-rect-90x70.scad
      Note: |-
        Generated SCAD file showing actual dimension assignments
        Generated SCAD file showing actual dimension assignments from YAML
    - Path: projects/film-developer/enclosure-rect-90x70.yaml
      Note: |-
        Source YAML configuration for film developer enclosure
        Source YAML configuration demonstrating dimension parameters
ExternalSources: []
Summary: Comprehensive analysis of how YAPPgenerator_v3.scad computes final enclosure dimensions from PCB geometry and configuration parameters
LastUpdated: 2025-11-28T19:30:00-05:00
---


# YAPP Enclosure Dimension Computation - Complete Analysis

## Executive Summary

This document provides a comprehensive, verified analysis of how `YAPPgenerator_v3.scad` computes final enclosure dimensions. It traces the complete pipeline from PCB geometry inputs through intermediate calculations to final outer shell dimensions, with detailed explanations of all configuration parameters and their effects.

## Table of Contents

1. [Input Parameters](#input-parameters)
2. [PCB Array Structure](#pcb-array-structure)
3. [Dimension Computation Pipeline](#dimension-computation-pipeline)
4. [Configuration Parameters](#configuration-parameters)
5. [Multiple PCB Handling](#multiple-pcb-handling)
6. [Vertical Dimension Computation](#vertical-dimension-computation)
7. [Practical Examples](#practical-examples)
8. [Configuration Methods](#configuration-methods)

---

## Input Parameters

### PCB Geometry Parameters

The generator accepts PCB dimensions through individual variables or via the `pcb` array:

```scad
pcbLength           = 120;  // X-axis (front to back)
pcbWidth            = 50;   // Y-axis (side to side)
pcbThickness        = 1.6;  // Z-axis (PCB board thickness)
standoffHeight      = 1.0;  // Height from base plane to PCB bottom
standoffDiameter    = 7;    // Standoff post diameter
standoffPinDiameter = 2.4;  // Screw/insert hole diameter
standoffHoleSlack   = 0.4;  // Tolerance for standoff holes
```

### Enclosure Configuration Parameters

```scad
// Wall and padding
wallThickness       = 2.8;  // Thickness of enclosure walls
paddingFront        = 1;    // Clearance between PCB front edge and inner wall
paddingBack         = 1;    // Clearance between PCB back edge and inner wall
paddingRight        = 1;    // Clearance between PCB right edge and inner wall
paddingLeft         = 1;    // Clearance between PCB left edge and inner wall

// Base configuration
basePlaneThickness  = 1.6;  // Thickness of base bottom plane
baseWallHeight      = 10;   // Height of base walls (from base plane)

// Lid configuration
lidPlaneThickness   = 1.6;  // Thickness of lid top plane
lidWallHeight       = 10;   // Height of lid walls (from lid plane)

// Ridge configuration (for snap-fit)
ridgeHeight         = 5.0;  // Height of overlapping ridge
ridgeSlack          = 0.2;  // Gap between lid inner wall and base outer wall
ridgeGap            = 0.5;  // Gap between base ridge bottom and lid bottom

// Aesthetic
roundRadius         = wallThickness + 1;  // Corner fillet radius
```

---

## PCB Array Structure

### Array Format

PCBs are defined in an array where each element represents one PCB:

```scad
pcb = [
  // Format: [name, length, width, posx, posy, thickness, standoffHeight, 
  //          standoffDiameter, standoffPinDiameter, standoffHoleSlack]
  ["Main", pcbLength, pcbWidth, 0, 0, pcbThickness, 
   standoffHeight, standoffDiameter, standoffPinDiameter, standoffHoleSlack]
];
```

**Array Element Indices:**
- `[0]` = PCB name (string, e.g., "Main")
- `[1]` = length (X-axis dimension)
- `[2]` = width (Y-axis dimension)
- `[3]` = posx (X-axis offset from origin)
- `[4]` = posy (Y-axis offset from origin)
- `[5]` = thickness (Z-axis dimension)
- `[6]` = standoffHeight (height from base plane)
- `[7]` = standoffDiameter
- `[8]` = standoffPinDiameter
- `[9]` = standoffHoleSlack (optional, defaults to 0.4)

### Coordinate System

The PCB coordinate system uses:
- **Origin**: Bottom-left corner of the PCB (when posx=0, posy=0)
- **X-axis**: Front to back (length)
- **Y-axis**: Left to right (width)
- **Z-axis**: Bottom to top (height)

**Important**: The `posx` and `posy` offsets are relative to the PCB coordinate origin, not the box coordinate origin. Positive `posx` moves the PCB forward (toward front), positive `posy` moves it right.

---

## Dimension Computation Pipeline

### Step 1: Compute PCB Envelope

The generator first computes the bounding box that contains all PCBs:

```324:326:pkg/yappgen/assets/YAPPgenerator_v3.scad
boxLength = maxLength(pcb);
boxWidth = maxWidth(pcb);
```

**Implementation:**

```5692:5693:pkg/yappgen/assets/YAPPgenerator_v3.scad
function maxLength(v, i = 0, r = 0) = i < len(v) ? maxLength(v, i + 1, max(r, v[i][1] + v[i][3])) : r;
function maxWidth(v, i = 0, r = 0) = i < len(v) ? maxWidth(v, i + 1, max(r, v[i][2] + v[i][4])) : r;
```

**Explanation:**
- `maxLength()` finds the maximum `(length + posx)` across all PCBs
- `maxWidth()` finds the maximum `(width + posy)` across all PCBs
- These represent the extreme reach of all PCBs in the X and Y directions

**Example:**
- PCB1: length=90, posx=0 → extent = 90
- PCB2: length=50, posx=60 → extent = 110
- Result: `boxLength = 110` (covers both PCBs)

### Step 2: Compute Inner Shell Dimensions

Add padding to create the inner cavity dimensions:

```349:351:pkg/yappgen/assets/YAPPgenerator_v3.scad
shellInsideWidth  = boxWidth+paddingLeft+paddingRight;
shellInsideLength = boxLength+paddingFront+paddingBack;
shellInsideHeight = baseWallHeight+lidWallHeight;
```

**Formulas:**
- `shellInsideLength = boxLength + paddingFront + paddingBack`
- `shellInsideWidth = boxWidth + paddingLeft + paddingRight`
- `shellInsideHeight = baseWallHeight + lidWallHeight`

These represent the **inner dimensions** of the enclosure cavity (the space inside the walls).

### Step 3: Compute Outer Shell Dimensions

Add wall thickness to compute the final outer dimensions:

```353:355:pkg/yappgen/assets/YAPPgenerator_v3.scad
shellWidth        = shellInsideWidth+(wallThickness*2);
shellLength       = shellInsideLength+(wallThickness*2);
shellHeight       = basePlaneThickness+shellInsideHeight+lidPlaneThickness;
```

**Formulas:**
- `shellLength = shellInsideLength + (wallThickness × 2)`
- `shellWidth = shellInsideWidth + (wallThickness × 2)`
- `shellHeight = basePlaneThickness + shellInsideHeight + lidPlaneThickness`

**Complete Formula Chain:**

```
shellLength = (maxLength(pcb) + paddingFront + paddingBack) + (wallThickness × 2)
shellWidth  = (maxWidth(pcb) + paddingLeft + paddingRight) + (wallThickness × 2)
shellHeight  = basePlaneThickness + (baseWallHeight + lidWallHeight) + lidPlaneThickness
```

---

## Configuration Parameters

### Padding Parameters

**Purpose**: Control clearance between PCB edges and inner walls.

**Parameters:**
- `paddingFront` - Clearance at front (positive X direction)
- `paddingBack` - Clearance at back (negative X direction)
- `paddingLeft` - Clearance at left (negative Y direction)
- `paddingRight` - Clearance at right (positive Y direction)

**Default Values:** Typically 1.0 mm, but can be set independently.

**Effect on Dimensions:**
- Each padding value directly adds to the corresponding inner dimension
- Asymmetric padding is fully supported in SCAD
- Zero padding means PCB edge sits flush against inner wall (not recommended)

**YAML Mapping:**
- The YAML DSL currently maps `enclosure.wall.clearance` to all four padding variables symmetrically
- To achieve asymmetric padding, edit the generated SCAD file directly

### Wall Thickness

**Purpose**: Controls the thickness of all enclosure walls.

**Parameter:** `wallThickness`

**Default Value:** 2.8 mm (in template), but commonly 1.4-3.0 mm depending on material and printer.

**Effect on Dimensions:**
- Adds `wallThickness × 2` to both length and width outer dimensions
- Does not affect inner dimensions
- Critical for structural integrity and snap-fit features

**Ridge Constraint:**
- For snap-fit joints, `ridgeHeight` must be ≥ `wallThickness × 1.8` (see `wallToRidgeRatio`)

### Base Configuration

**Parameters:**
- `basePlaneThickness` - Thickness of the bottom plane (typically 1.5-2.5 mm)
- `baseWallHeight` - Height of base walls from base plane

**Effect:**
- `basePlaneThickness` adds to total `shellHeight`
- `baseWallHeight` contributes to `shellInsideHeight`

**Typical Values:**
- `basePlaneThickness`: 1.5-2.5 mm (structural support)
- `baseWallHeight`: Varies based on PCB clearance needs

### Lid Configuration

**Parameters:**
- `lidPlaneThickness` - Thickness of the top plane (typically 1.0-2.0 mm)
- `lidWallHeight` - Height of lid walls from lid plane

**Effect:**
- `lidPlaneThickness` adds to total `shellHeight`
- `lidWallHeight` contributes to `shellInsideHeight`

**Vertical Clearance Calculation:**
The space above the PCB top is:
```
clearance_above_pcb = shellInsideHeight - (standoffHeight + pcbThickness)
```

Or equivalently:
```
clearance_above_pcb = (baseWallHeight + lidWallHeight) - (standoffHeight + pcbThickness)
```

---

## Multiple PCB Handling

### How Multiple PCBs Affect Dimensions

When multiple PCBs are defined in the `pcb` array, the generator computes the **union bounding box** that contains all PCBs.

**Example with Two PCBs:**

```scad
pcb = [
  ["Main", 90, 70, 0, 0, 1.6, 3, 7, 3, 0.4],
  ["Secondary", 50, 40, 30, 20, 1.6, 3, 7, 3, 0.4]
];
```

**Computation:**
- Main PCB: length=90, posx=0 → X extent = 90
- Secondary PCB: length=50, posx=30 → X extent = 80
- `boxLength = max(90, 80) = 90`

- Main PCB: width=70, posy=0 → Y extent = 70
- Secondary PCB: width=40, posy=20 → Y extent = 60
- `boxWidth = max(70, 60) = 70`

**Result:** The enclosure is sized to fit both PCBs, with the Main PCB determining the dimensions in this case.

### PCB Offset Behavior

**Positive `posx`**: Moves PCB forward (toward front wall)
- Increases `boxLength` if `posx + length` exceeds other PCBs
- Example: PCB at posx=10 with length=80 creates extent=90

**Positive `posy`**: Moves PCB right (toward right wall)
- Increases `boxWidth` if `posy + width` exceeds other PCBs
- Example: PCB at posy=5 with width=65 creates extent=70

**Negative offsets**: Not typically used, but would move PCB toward origin (back/left)

---

## Vertical Dimension Computation

### Complete Height Formula

```131:136:pkg/yappgen/assets/YAPPgenerator_v3.scad
//-- Total height of box = lidPlaneThickness 
//                       + lidWallHeight 
//                       + baseWallHeight 
//                       + basePlaneThickness
//-- space between pcb and lidPlane :=
//--      (bottonWallHeight+lidWallHeight) - (standoff_Height+pcb_Thickness)
```

**Total Shell Height:**
```
shellHeight = basePlaneThickness + baseWallHeight + lidWallHeight + lidPlaneThickness
```

**Vertical Clearance Above PCB:**
```
clearance_above_pcb = (baseWallHeight + lidWallHeight) - (standoffHeight + pcbThickness)
```

### Standoff Height Behavior

**Positive `standoffHeight`**: Measured from base plane upward
- Standard case: PCB sits above base plane
- Example: `standoffHeight = 3` means PCB bottom is 3mm above base plane

**Negative `standoffHeight`**: Measured from lid plane downward
- Special case: PCB mounted from lid
- The generator converts this internally:
  ```scad
  standoffHeight = shellInsideHeight - pcbThickness + (negative_value)
  ```

### Ridge Interaction

The ridge system creates an overlapping region between base and lid:

```
ridgeHeight         = 5.0;  // Overlap distance
ridgeSlack          = 0.2;  // Gap between lid inner wall and base outer wall
ridgeGap            = 0.5;  // Gap at bottom of ridge
```

**Constraint:** `ridgeHeight` must be ≥ `lidWallHeight` OR ≥ `wallThickness × 1.8` (for snap-fit)

---

## Practical Examples

### Example 1: Film Developer Enclosure

**Input Configuration:**
```yaml
pcb:
  length: 90.0
  width: 70.0
  thickness: 1.6
  z_clearance: 3.0

enclosure:
  wall:
    thickness: 2.4
    clearance: 1.5
  base:
    thickness: 2.0
    wall_height: 8.6  # = 3.0 + 1.6 + 4.0
  lid:
    thickness: 2.0
    wall_height: 24   # = 28.0 + 3.0 + 1.6 - 8.6
```

**Generated SCAD Variables:**
```scad
pcbLength = 90;
pcbWidth = 70;
paddingFront = 1.5;
paddingBack = 1.5;
paddingLeft = 1.5;
paddingRight = 1.5;
wallThickness = 2.4;
basePlaneThickness = 2;
lidPlaneThickness = 2;
baseWallHeight = 8.6;
lidWallHeight = 24;
```

**Computation Steps:**

1. **PCB Envelope:**
   - `boxLength = maxLength(pcb) = 90` (single PCB, no offset)
   - `boxWidth = maxWidth(pcb) = 70`

2. **Inner Dimensions:**
   - `shellInsideLength = 90 + 1.5 + 1.5 = 93`
   - `shellInsideWidth = 70 + 1.5 + 1.5 = 73`
   - `shellInsideHeight = 8.6 + 24 = 32.6`

3. **Outer Dimensions:**
   - `shellLength = 93 + (2.4 × 2) = 97.8`
   - `shellWidth = 73 + (2.4 × 2) = 77.8`
   - `shellHeight = 2 + 32.6 + 2 = 36.6`

**Verification:**
- Vertical clearance above PCB: `32.6 - (3 + 1.6) = 28.0` ✓ (matches `vars.top_clearance`)

### Example 2: Demo Buttons Enclosure

**Input Configuration:**
```scad
pcbLength = 30;
pcbWidth = 40;
paddingFront = 1;
paddingBack = 1;
paddingRight = 1;
paddingLeft = 1;
wallThickness = 1.4;
basePlaneThickness = 1.5;
lidPlaneThickness = 1.0;
baseWallHeight = 10;
lidWallHeight = 10;
```

**Computation:**

1. **PCB Envelope:**
   - `boxLength = 30`
   - `boxWidth = 40`

2. **Inner Dimensions:**
   - `shellInsideLength = 30 + 1 + 1 = 32`
   - `shellInsideWidth = 40 + 1 + 1 = 42`
   - `shellInsideHeight = 10 + 10 = 20`

3. **Outer Dimensions:**
   - `shellLength = 32 + (1.4 × 2) = 34.8`
   - `shellWidth = 42 + (1.4 × 2) = 44.8`
   - `shellHeight = 1.5 + 20 + 1.0 = 22.5`

### Example 3: Multiple PCBs with Offsets

**Configuration:**
```scad
pcb = [
  ["Main", 90, 70, 0, 0, 1.6, 3, 7, 3, 0.4],
  ["Sensor", 30, 20, 50, 10, 1.6, 3, 7, 3, 0.4]
];
paddingFront = 2;
paddingBack = 2;
paddingLeft = 2;
paddingRight = 2;
wallThickness = 2.0;
```

**Computation:**

1. **PCB Envelope:**
   - Main: X extent = 90 + 0 = 90
   - Sensor: X extent = 30 + 50 = 80
   - `boxLength = max(90, 80) = 90`
   
   - Main: Y extent = 70 + 0 = 70
   - Sensor: Y extent = 20 + 10 = 30
   - `boxWidth = max(70, 30) = 70`

2. **Inner Dimensions:**
   - `shellInsideLength = 90 + 2 + 2 = 94`
   - `shellInsideWidth = 70 + 2 + 2 = 74`

3. **Outer Dimensions:**
   - `shellLength = 94 + (2.0 × 2) = 98`
   - `shellWidth = 74 + (2.0 × 2) = 78`

**Note:** The Sensor PCB's offset doesn't affect dimensions because it fits within the Main PCB's extent.

---

## Configuration Methods

### Method 1: Direct SCAD Variable Assignment

Set variables directly in your SCAD file before including `YAPPgenerator_v3.scad`:

```scad
pcbLength = 90;
pcbWidth = 70;
paddingFront = 1.5;
paddingBack = 1.5;
paddingLeft = 1.5;
paddingRight = 1.5;
wallThickness = 2.4;
basePlaneThickness = 2;
lidPlaneThickness = 2;
baseWallHeight = 8.6;
lidWallHeight = 24;

include <YAPPgenerator_v3.scad>
```

**Advantages:**
- Full control over all parameters
- Can use asymmetric padding
- Can override PCB array directly

### Method 2: YAML DSL (via yappctl)

Define configuration in YAML, then generate SCAD:

```yaml
pcb:
  length: 90
  width: 70
  thickness: 1.6
  z_clearance: 3.0

enclosure:
  wall:
    thickness: 2.4
    clearance: 1.5  # Maps to all four padding variables symmetrically
  base:
    thickness: 2.0
    wall_height: 8.6
  lid:
    thickness: 2.0
    wall_height: 24
```

**Advantages:**
- Declarative and readable
- Expression evaluation support
- Version control friendly

**Limitations:**
- Currently only supports symmetric padding (`enclosure.wall.clearance` applies to all four sides)
- Cannot set final outer dimensions directly (computed from PCB + padding + wall thickness)
- Cannot set per-side padding (planned for future)
- Must regenerate SCAD after YAML changes

### Method 3: Hybrid Approach

Use YAML for initial generation, then edit SCAD for fine-tuning:

1. Generate SCAD from YAML: `yappctl generate projects/film-developer/enclosure.yaml`
2. Edit generated SCAD to adjust padding or other parameters
3. Use edited SCAD directly

**Use Case:** When you need asymmetric padding or other features not yet supported in YAML DSL.

---

## YAML DSL Dimension Configuration Capabilities

### What You CAN Configure

✅ **PCB Dimensions:**
- `pcb.length` - X-axis dimension
- `pcb.width` - Y-axis dimension
- `pcb.thickness` - Board thickness

✅ **Uniform Padding:**
- `enclosure.wall.clearance` - Applies symmetrically to all four sides (front, back, left, right)

✅ **Wall Thickness:**
- `enclosure.wall.thickness` - Thickness of all walls

✅ **Base/Lid Configuration:**
- `enclosure.base.thickness` - Base plane thickness
- `enclosure.base.wall_height` - Base wall height
- `enclosure.lid.thickness` - Lid plane thickness
- `enclosure.lid.wall_height` - Lid wall height

✅ **Computed Dimensions (via `vars`):**
You can compute final dimensions in `vars` for reference, but they don't directly control the generator:

```yaml
vars:
  shell_length: pcb.length + 2*enclosure.wall.clearance + 2*enclosure.wall.thickness
  shell_width: pcb.width + 2*enclosure.wall.clearance + 2*enclosure.wall.thickness
```

### What You CANNOT Configure (Currently)

❌ **Final Outer Dimensions:** 
- Cannot set `shellLength` or `shellWidth` directly
- These are computed from: `(PCB size + padding) + (wall thickness × 2)`

❌ **Per-Side Padding:**
- Cannot set `paddingFront`, `paddingBack`, `paddingLeft`, `paddingRight` independently
- `enclosure.wall.clearance` applies the same value to all four sides
- **Status:** Per-side clearance is planned but not yet implemented (see reference documentation)

### Implementation Details

**How Padding Works in the DSL:**

```138:148:pkg/yappgen/model.go
	// Map a single clearance value to all paddings if present
	if v, ok := getFloat(resolved, "enclosure.wall.clearance"); ok {
		m.PaddingFront = v
		m.PaddingBack = v
		m.PaddingLeft = v
		m.PaddingRight = v
		m.Provenance.AddScalar("paddingFront", "enclosure.wall.clearance")
		m.Provenance.AddScalar("paddingBack", "enclosure.wall.clearance")
		m.Provenance.AddScalar("paddingLeft", "enclosure.wall.clearance")
		m.Provenance.AddScalar("paddingRight", "enclosure.wall.clearance")
	}
```

The code explicitly maps a single `enclosure.wall.clearance` value to all four padding variables. There's no code path for per-side configuration.

**Workaround for Asymmetric Padding:**

While per-side padding isn't supported, you can work around it by:

1. **Using `vars` to compute offsets:**
   ```yaml
   vars:
     extra_front_clearance: 3.0
     base_padding: 1.5
   
   enclosure:
     wall:
       clearance: vars.base_padding  # Applied to all sides
   
   features:
     cutouts:
       # Position cutouts accounting for extra clearance needed
       - face: front
         from_back: vars.extra_front_clearance + vars.base_padding
         # ... other fields
   ```

2. **Editing generated SCAD:**
   - Generate SCAD from YAML
   - Manually edit `paddingFront`, `paddingBack`, `paddingLeft`, `paddingRight` in the SCAD file
   - Use the edited SCAD directly

**Reference Documentation:**

The DSL reference explicitly documents this limitation:

> **Design note:** `wall.clearance` applies symmetrically to front/back/left/right. If you need asymmetric padding (e.g., extra clearance on one side for a connector), compute custom offsets in `vars` and adjust cutout positions accordingly. Native per-side clearance is on the roadmap.

---

## Key Insights and Best Practices

### Dimension Relationships

1. **Inner dimensions** = PCB envelope + padding
2. **Outer dimensions** = Inner dimensions + (wall thickness × 2)
3. **Total height** = Base plane + Base walls + Lid walls + Lid plane

### Padding Guidelines

- **Minimum padding**: 0.5-1.0 mm for tight fits
- **Recommended padding**: 1.0-2.0 mm for comfortable clearance
- **Large components**: Increase padding where connectors/switches protrude

### Wall Thickness Guidelines

- **Thin walls** (1.0-1.5 mm): Lightweight, but may be fragile
- **Standard walls** (1.5-2.5 mm): Good balance of strength and material use
- **Thick walls** (2.5-3.5 mm): Maximum strength, but uses more material

### Vertical Clearance Planning

When designing for components that extend above the PCB:

1. Calculate required clearance: `required_clearance = component_height - pcb_thickness`
2. Set `lidWallHeight` to provide this clearance:
   ```
   lidWallHeight = required_clearance + standoffHeight + pcbThickness - baseWallHeight
   ```
3. Verify: `(baseWallHeight + lidWallHeight) - (standoffHeight + pcbThickness) >= required_clearance`

### Multiple PCB Considerations

- Place PCBs to minimize total envelope size
- Consider component placement when positioning PCBs
- Use offsets to optimize space usage
- Remember: The generator uses the **union** of all PCB bounding boxes

---

## Verification and Debugging

### Echo Output

The generator can output computed dimensions (when `printMessages = true`):

```scad
printMessages = true;
include <YAPPgenerator_v3.scad>
```

### Manual Verification

Check computed values match expectations:

```scad
echo("boxLength =", boxLength);
echo("boxWidth =", boxWidth);
echo("shellInsideLength =", shellInsideLength);
echo("shellInsideWidth =", shellInsideWidth);
echo("shellLength =", shellLength);
echo("shellWidth =", shellWidth);
echo("shellHeight =", shellHeight);
```

### Common Issues

1. **Dimensions too small**: Check padding values, may need to increase
2. **Components don't fit**: Verify vertical clearance calculation
3. **Unexpected large dimensions**: Check for PCB offsets extending beyond main PCB
4. **Ridge not working**: Ensure `ridgeHeight >= lidWallHeight` or `>= wallThickness × 1.8`

---

## Conclusion

The YAPP generator uses a straightforward but powerful dimension computation pipeline:

1. **PCB Envelope**: Compute bounding box from PCB array (handles multiple PCBs and offsets)
2. **Inner Shell**: Add padding to create inner cavity dimensions
3. **Outer Shell**: Add wall thickness to create final outer dimensions
4. **Height**: Sum base plane, walls, and lid plane

All dimensions are configurable through SCAD variables, with the YAML DSL providing a convenient way to generate these assignments. Understanding this pipeline enables precise control over enclosure dimensions for any project.

---

## References

- **Main Generator**: `pkg/yappgen/assets/YAPPgenerator_v3.scad` (lines 324-355, 5692-5693)
- **Example Configurations**: `projects/film-developer/enclosure-rect-90x70.yaml`
- **Generated Output**: `projects/film-developer/enclosure-rect-90x70.scad`
- **Demo Examples**: `examples/YAPP_Demo_buttons_v30.scad`

