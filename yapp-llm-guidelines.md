# YAPPgenerator LLM Playbook

**Purpose**: Comprehensive guide for LLMs to efficiently work with YAPPgenerator for creating parametric project boxes.

**Last Updated**: 2025-11-02  
**YAPPgenerator Version Tested**: v3.3.8 (2025-10-24)  
**OpenSCAD Version**: 2021.01-4build1

---

## Table of Contents

1. [Quick Start Workflow](#quick-start-workflow)
2. [Critical Information](#critical-information)
3. [Documentation Sources](#documentation-sources)
4. [Common Pitfalls & Solutions](#common-pitfalls--solutions)
5. [Feature Selection Guide](#feature-selection-guide)
6. [Coordinate Systems Explained](#coordinate-systems-explained)
7. [Rendering & STL Generation](#rendering--stl-generation)
8. [Parameter Syntax Reference](#parameter-syntax-reference)
9. [What Worked Well](#what-worked-well)
10. [What Didn't Work](#what-didnt-work)
11. [Documentation Issues](#documentation-issues)
12. [Best Practices](#best-practices)
13. [Code Examples](#code-examples)
14. [Troubleshooting](#troubleshooting)
15. [Performance Optimization](#performance-optimization)

---

## Quick Start Workflow

### Step-by-Step Process (30-45 minutes total)

1. **Install Dependencies** (5 min)
   ```bash
   sudo apt-get update
   sudo apt-get install -y openscad xvfb
   ```

2. **Clone Repository** (1 min)
   ```bash
   cd /home/ubuntu
   git clone https://github.com/mrWheel/YAPP_Box.git
   ```

3. **Review Documentation** (10 min)
   - Browse online docs: https://mrwheel-docs.gitbook.io/yappgenerator_en/
   - Read template: `/home/ubuntu/YAPP_Box/YAPP_Template_v3.scad`
   - Check examples: `/home/ubuntu/YAPP_Box/examples/`
   - **DO NOT** rely solely on online docs - they may be outdated

4. **Create Design File** (10 min)
   - Start from template or example
   - Define PCB dimensions first
   - Add features incrementally
   - Test render frequently

5. **Generate Preview** (2-3 min)
   ```bash
   xvfb-run -a openscad --render --imgsize=1920,1440 \
     --colorscheme=Tomorrow --camera=0,0,0,60,0,30,250 \
     -o preview.png design.scad
   ```

6. **Generate STL Files** (2-4 min total)
   - Create separate files for base and lid
   - Render each independently
   ```bash
   xvfb-run -a openscad -o base.stl base_only.scad
   xvfb-run -a openscad -o lid.stl lid_only.scad
   ```

7. **Create Documentation** (10 min)
   - README with specifications
   - Technical report explaining choices
   - Verification against docs

---

## Critical Information

### Version Discrepancy Alert ⚠️

**IMPORTANT**: The online documentation claims to be for v3.0, but the actual repository version is **v3.3.8**. There have been significant changes between versions.

**Always verify against the actual source code**, not just the documentation.

### Version Check Command
```bash
grep "^Version=" /path/to/YAPPgenerator_v3.scad
```

### File Locations

| File | Path | Purpose |
|------|------|---------|
| Main library | `YAPPgenerator_v3.scad` | Core generator code |
| Template | `YAPP_Template_v3.scad` | Boilerplate with all parameters |
| Examples | `examples/*.scad` | Working examples |
| Online docs | https://mrwheel-docs.gitbook.io/yappgenerator_en/ | Reference (may be outdated) |

---

## Documentation Sources

### Primary Sources (In Order of Reliability)

1. **Source Code** (Most Reliable)
   - File: `YAPPgenerator_v3.scad`
   - Lines 1-500: Parameter definitions and comments
   - Search for feature names to find implementation

2. **Template File** (Very Reliable)
   - File: `YAPP_Template_v3.scad`
   - Contains complete parameter documentation
   - Shows default values
   - Lines 98-600: All feature arrays with comments

3. **Example Files** (Highly Reliable)
   - Directory: `examples/`
   - Real working code
   - Best for understanding syntax
   - Key examples:
     - `YAPP_Demo_buttons_v31.scad` - Push buttons
     - `YAPP_Demo_cutouts_v30.scad` - Cutouts
     - `YAPP_Demo_DisplayMount_v31.scad` - Display mounting

4. **Online Documentation** (Use with Caution)
   - URL: https://mrwheel-docs.gitbook.io/yappgenerator_en/
   - Claims to be v3.0 but repository is v3.3.8
   - Good for concepts, verify syntax in code
   - **DO NOT** trust parameter order blindly

### How to Verify Documentation

```bash
# Find a feature in source code
grep -n "cutoutsLid" /home/ubuntu/YAPP_Box/YAPPgenerator_v3.scad

# Find examples using a feature
grep -r "yappCircle" /home/ubuntu/YAPP_Box/examples/

# Check template for parameter documentation
grep -A 20 "*** Cutouts ***" /home/ubuntu/YAPP_Box/YAPP_Template_v3.scad
```

---

## Common Pitfalls & Solutions

### 1. Push Buttons vs. Cutouts Confusion

**Pitfall**: User asks for "buttons" - unclear if they want:
- Push button guides (for PCB-mounted tactile switches)
- Simple cutouts (for panel-mount buttons)

**Solution**: 
- **Always clarify** the use case first
- Push button guides: Use `pushButtons` array
- Panel-mount buttons: Use `cutoutsLid` array with `yappCircle`

**Example Clarification**:
```
User: "I need 3 buttons"
LLM: "I can create either:
1. Push button guides (for tactile switches on PCB)
2. Simple cutouts (for panel-mount buttons)
Which type do you need?"
```

### 2. Coordinate System Confusion

**Pitfall**: Mixing coordinate systems or using wrong origin.

**Solution**:
- **Default is `yappCoordPCB`** - PCB origin at back-left corner
- Use `yappCenter` for circles to specify center point
- Document which coordinate system you're using

**Wrong**:
```openscad
cutoutsLid = [
  [20, 25, 0, 0, 9, yappCircle]  // Position is corner, not center!
];
```

**Correct**:
```openscad
cutoutsLid = [
  [20, 25, 0, 0, 9, yappCircle, 0, 0, yappCenter]  // Center-based
];
```

### 3. Parameter Order Errors

**Pitfall**: Online docs may show outdated parameter order.

**Solution**:
- **Always check template file** for current parameter order
- Count parameters carefully
- Use `undef` for optional parameters you want to skip

**Example from Template** (lines 310-333):
```openscad
//  Parameters:
//   Required:
//    p(0) = from Back
//    p(1) = from Left
//    p(2) = width
//    p(3) = length
//    p(4) = radius
//    p(5) = shape
//  Optional:
//    p(6) = depth
//    p(7) = angle
//    n(a) = { yappPolygonDef }
//    n(b) = { yappMaskDef }
//    n(c) = { [yappMaskDef, hOffset, vOffset, rotation] }
//    n(d) = { <yappCoordPCB> | yappCoordBox | yappCoordBoxInside }
//    n(e) = { <yappOrigin>, yappCenter }
```

### 4. Named Parameters (n) vs. Positional Parameters (p)

**Pitfall**: Treating named parameters as positional.

**Solution**:
- **Positional parameters (p)**: Must be in order, use `undef` to skip
- **Named parameters (n)**: Can appear anywhere after positional params
- Named params are flags/enums, not indexed

**Wrong**:
```openscad
cutoutsLid = [
  [36, 11, 28, 28, 0, yappRectangle, yappCoordPCB]  // yappCoordPCB as p(6)
];
```

**Correct**:
```openscad
cutoutsLid = [
  [36, 11, 28, 28, 0, yappRectangle, 0, 0, yappCoordPCB]  // As named param
  // OR just omit since yappCoordPCB is default:
  [36, 11, 28, 28, 0, yappRectangle]
];
```

### 5. Rendering Without Virtual Display

**Pitfall**: OpenSCAD fails with "Can't create OpenGL OffscreenView" in headless environment.

**Solution**:
```bash
# Install xvfb
sudo apt-get install -y xvfb

# Use xvfb-run wrapper
xvfb-run -a openscad -o output.stl input.scad
```

### 6. Incomplete STL Generation

**Pitfall**: Generating assembly STL instead of separate base/lid.

**Solution**:
- Create separate `.scad` files for base and lid
- Set print flags appropriately:
  ```openscad
  // base_only.scad
  printBaseShell = true;
  printLidShell = false;
  printSwitchExtenders = false;
  
  // lid_only.scad
  printBaseShell = false;
  printLidShell = true;
  printSwitchExtenders = true;  // Include if using pushButtons
  ```

### 7. Centering Calculations

**Pitfall**: Incorrect centering math for cutouts.

**Solution**:
```openscad
// For rectangle cutout (origin at corner):
centerX = (pcbLength - cutoutLength) / 2;
centerY = (pcbWidth - cutoutWidth) / 2;

// For circle cutout with yappCenter:
centerX = pcbLength / 2;  // Direct center
centerY = pcbWidth / 2;

// Example:
// PCB: 100x50mm, Display: 28x28mm
displayX = (100 - 28) / 2;  // = 36mm
displayY = (50 - 28) / 2;   // = 11mm
```

---

## Feature Selection Guide

### Decision Tree for Common Requirements

```
User needs to mount component on/through lid?
│
├─ Display/Screen?
│  └─ Use: cutoutsLid with yappRectangle
│     Position: Calculate center
│     Size: Exact display dimensions
│
├─ Panel-mount buttons (with housing)?
│  └─ Use: cutoutsLid with yappCircle + yappCenter
│     Diameter: Button housing diameter
│     Position: Button center coordinates
│
├─ Tactile switches on PCB (need button caps)?
│  └─ Use: pushButtons array
│     Parameters: Cap size, switch specs, pole diameter
│     Note: More complex, generates button mechanisms
│
├─ Connector (USB, HDMI, etc.)?
│  └─ Use: cutoutsFront/Back/Left/Right
│     Shape: Match connector (yappRectangle, yappRoundedRect)
│     Position: Align with PCB connector
│
└─ Ventilation/Speaker grille?
   └─ Use: cutoutsBase/Lid with yappMaskDef
      Mask: maskHoneycomb, maskBars, etc.
```

### Feature Comparison Table

| Requirement | Feature | Complexity | File Size Impact | Use When |
|-------------|---------|------------|------------------|----------|
| Display opening | `cutoutsLid` | Low | Minimal | Display has bezel |
| Panel buttons | `cutoutsLid` | Low | Minimal | Buttons have housing |
| PCB buttons | `pushButtons` | High | Large (+4MB) | Tactile switches on PCB |
| PCB mounting | `pcbStands` | Low | Small | Always needed |
| Lid attachment | `snapJoins` | Low | Small | Tool-free assembly |
| Screw mounting | `connectors` | Medium | Medium | Permanent assembly |

---

## Coordinate Systems Explained

### Three Coordinate Systems

YAPPgenerator uses three coordinate systems. Understanding when to use each is critical.

#### 1. `yappCoordPCB` (Default)

**Origin**: PCB back-left corner at [0, 0, 0]

**Use for**: 
- Components on/above PCB
- Cutouts aligned with PCB features
- Standoffs
- Push buttons

**Diagram**:
```
        Back (X=0)
    ┌─────────────────┐
    │ 0,0         x,0 │
L   │                 │  R
e   │      PCB        │  i
f   │                 │  g
t   │ 0,y         x,y │  h
    └─────────────────┘  t
        Front (X=max)
```

**Example**:
```openscad
// Display centered on 100x50mm PCB
cutoutsLid = [
  [36, 11, 28, 28, 0, yappRectangle]  // Uses yappCoordPCB by default
];
// 36 = (100-28)/2, 11 = (50-28)/2
```

#### 2. `yappCoordBox` (Box Exterior)

**Origin**: Box back-left-bottom corner (outside)

**Use for**:
- Box mounts
- Labels on exterior
- Features independent of PCB

**Conversion**:
```
boxX = pcbX + paddingBack + wallThickness
boxY = pcbY + paddingLeft + wallThickness
```

#### 3. `yappCoordBoxInside` (Box Interior)

**Origin**: Box back-left-bottom corner (inside)

**Use for**:
- Internal features
- Hooks
- Cable management

**Conversion**:
```
boxInsideX = pcbX + paddingBack
boxInsideY = pcbY + paddingLeft
```

### Origin Modifiers

#### `yappOrigin` (Default)
Position measured from back-left corner.

#### `yappCenter`
Position measured from center point.

**Critical for circles**:
```openscad
// Without yappCenter (position is bounding box corner):
[20-9, 25-9, 0, 0, 9, yappCircle]  // Must subtract radius

// With yappCenter (position is circle center):
[20, 25, 0, 0, 9, yappCircle, 0, 0, yappCenter]  // Intuitive
```

---

## Rendering & STL Generation

### Rendering Performance

| Operation | Typical Time | Vertices | Notes |
|-----------|--------------|----------|-------|
| Preview (no render) | 5-10s | N/A | Fast, low quality |
| Full render (simple) | 15-30s | 5,000-10,000 | Base or lid only |
| Full render (complex) | 1-3 min | 15,000-25,000 | With push buttons |
| STL export (simple) | 20-40s | 5,000-10,000 | Includes render time |
| STL export (complex) | 2-4 min | 15,000-25,000 | With push buttons |

### Rendering Commands

#### Preview Image (Fast)
```bash
xvfb-run -a openscad \
  --render \
  --imgsize=1920,1440 \
  --colorscheme=Tomorrow \
  --camera=0,0,0,60,0,30,250 \
  -o preview.png \
  design.scad
```

**Parameters**:
- `--render`: Full CGAL render (vs. preview)
- `--imgsize=W,H`: Output resolution
- `--colorscheme`: BeforeDawn, Tomorrow, Cornfield, etc.
- `--camera=X,Y,Z,RotX,RotY,RotZ,Dist`: Camera position

#### Orthographic Top View
```bash
xvfb-run -a openscad \
  --render \
  --imgsize=1600,1200 \
  --colorscheme=Tomorrow \
  --projection=ortho \
  --viewall \
  --camera=0,0,0,0,0,0,150 \
  -o top_view.png \
  lid_only.scad
```

**Key differences**:
- `--projection=ortho`: No perspective distortion
- `--viewall`: Auto-frame entire model
- `--camera=0,0,0,0,0,0,150`: Top-down view

#### STL Export
```bash
xvfb-run -a openscad -o output.stl input.scad
```

**Simple and fast**. No image parameters needed.

### Camera Angles Reference

```
Camera: X,Y,Z, RotX,RotY,RotZ, Distance

Common views:
- Isometric:  0,0,0, 60,0,30, 250
- Top:        0,0,0, 0,0,0, 150
- Bottom:     0,0,0, 0,0,180, 150
- Front:      0,0,0, 90,0,0, 150
- Back:       0,0,0, 90,0,180, 150
- Left:       0,0,0, 90,0,90, 150
- Right:      0,0,0, 90,0,270, 150
```

### Monitoring Render Progress

```bash
# Render with progress output
xvfb-run -a openscad -o output.stl input.scad 2>&1 | \
  grep -E "(rendering|Vertices|Facets|error)"

# Expected output:
# Total rendering time: 0:01:23.456
#    Vertices:    12345
#    Facets:      6789
```

---

## Parameter Syntax Reference

### Array Syntax Patterns

#### Positional + Named Parameters
```openscad
featureArray = [
  [p0, p1, p2, p3, namedFlag1, namedFlag2],
  [p0, p1, p2, p3, p4, p5, namedFlag1]
];
```

#### Using `undef` to Skip Optional Parameters
```openscad
pushButtons = [
  [20, 25, 0, 0, 9, 2, 5, 0.5, 3.5, undef, yappCircle]
  // p(9) = undef means use default height
];
```

#### Multiple Named Parameters
```openscad
pcbStands = [
  [5, 5, yappAllCorners, yappPin, yappBoth]
  // All are named parameters, order doesn't matter
];
```

### Common Named Parameters

| Parameter | Values | Meaning |
|-----------|--------|---------|
| Coordinate system | `yappCoordPCB`, `yappCoordBox`, `yappCoordBoxInside` | Origin reference |
| Origin type | `yappOrigin`, `yappCenter` | Position reference point |
| Sides | `yappLeft`, `yappRight`, `yappFront`, `yappBack` | Which wall(s) |
| Corners | `yappAllCorners`, `yappFrontLeft`, `yappBackRight`, etc. | Which corner(s) |
| Placement | `yappBoth`, `yappLidOnly`, `yappBaseOnly` | Which shell |
| Pin type | `yappPin`, `yappHole`, `yappTopPin` | Standoff type |
| Shape | `yappRectangle`, `yappCircle`, `yappRoundedRect`, `yappPolygon` | Cutout shape |

---

## What Worked Well

### 1. Incremental Development ✓

**Approach**: Build design step-by-step, testing each feature.

**Workflow**:
1. Define PCB dimensions
2. Add standoffs → render → verify
3. Add cutouts → render → verify
4. Add snap joins → render → verify
5. Generate final STL

**Why it worked**: Caught errors early, easy to debug.

### 2. Separate Base/Lid Files ✓

**Approach**: Create three files:
- `design.scad` - Full assembly for preview
- `base_only.scad` - Base shell only
- `lid_only.scad` - Lid shell only

**Why it worked**:
- Faster rendering (20s vs 2min)
- Easier to verify individual parts
- Separate STL files for printing

### 3. Using Examples as Templates ✓

**Approach**: Copy working example, modify for needs.

**Best examples**:
- `YAPP_Demo_buttons_v31.scad` - Button syntax
- `YAPP_Demo_cutouts_v30.scad` - Cutout shapes
- `YAPP_Template_v3.scad` - Complete reference

**Why it worked**: Guaranteed correct syntax, less trial-and-error.

### 4. Verification Against Source Code ✓

**Approach**: When docs unclear, check source code.

```bash
grep -A 30 "cutoutsLid" YAPPgenerator_v3.scad
```

**Why it worked**: Source code is always correct, docs may lag.

### 5. Documenting Calculations ✓

**Approach**: Show math in comments.

```openscad
// Display: 28x28mm, PCB: 100x50mm
// Center X: (100-28)/2 = 36mm
// Center Y: (50-28)/2 = 11mm
cutoutsLid = [
  [36, 11, 28, 28, 0, yappRectangle]
];
```

**Why it worked**: Easy to verify, modify, and explain.

### 6. Using `yappCenter` for Circles ✓

**Approach**: Always use `yappCenter` with circular cutouts.

```openscad
[20, 25, 0, 0, 9, yappCircle, 0, 0, yappCenter]
```

**Why it worked**: Intuitive positioning, no manual offset calculations.

### 7. Creating Technical Documentation ✓

**Approach**: Explain *why* each feature was chosen, not just *what*.

**Sections**:
- Feature selection rationale
- Alternative approaches considered
- Design calculations
- Trade-offs

**Why it worked**: User understands design, can modify confidently.

---

## What Didn't Work

### 1. Trusting Online Documentation Blindly ✗

**Problem**: Online docs claim v3.0, but repository is v3.3.8.

**What failed**:
- Parameter order differed
- Some features not documented
- Syntax examples outdated

**Lesson**: Always verify against template and source code.

### 2. Initial Misunderstanding of User Requirements ✗

**Problem**: Created push button guides when user wanted simple cutouts.

**What failed**:
- Assumed "buttons" meant tactile switches
- Generated complex button mechanisms
- Wasted 10+ minutes rendering

**Lesson**: Clarify use case before implementing. Ask:
- "Do you have tactile switches on PCB or panel-mount buttons?"
- "Does the display have its own mounting bezel?"

### 3. Rendering Full Assembly for STL ✗

**Problem**: First attempt rendered both base and lid in one STL.

**What failed**:
- Can't print both parts simultaneously
- Harder to orient for printing
- Larger file size

**Lesson**: Always create separate base/lid files from the start.

### 4. Not Using `yappCenter` Initially ✗

**Problem**: Positioned circles using corner coordinates.

**What failed**:
```openscad
// Wrong approach:
[20-9, 25-9, 0, 0, 9, yappCircle]  // Had to calculate offset
```

**Lesson**: Use `yappCenter` for all circular features.

### 5. Assuming Default Coordinate System ✗

**Problem**: Didn't explicitly verify which coordinate system was default.

**What failed**: Initially thought `yappCoordBox` might be default.

**Lesson**: Template clearly states `yappCoordPCB` is default (line 287). Always check.

---

## Documentation Issues

### Issues Found in Online Documentation

#### Issue 1: Version Mismatch

**Location**: https://mrwheel-docs.gitbook.io/yappgenerator_en/

**Claim**: "This documentation applies to v3.0 (from Februari 2024)"

**Reality**: Repository is v3.3.8 (2025-10-24)

**Impact**: Parameter order and features may differ.

**Workaround**: Always check template file for current syntax.

#### Issue 2: Incomplete `yappCenter` Documentation

**Location**: Cutouts section

**Problem**: Doesn't clearly explain when to use `yappCenter` vs. `yappOrigin`.

**Reality**: 
- `yappOrigin` (default): Position is corner/edge
- `yappCenter`: Position is center point
- **Critical for circles** to avoid offset calculations

**Workaround**: Check examples like `YAPP_Demo_cutouts_v30.scad`.

#### Issue 3: Named vs. Positional Parameter Confusion

**Location**: Multiple sections

**Problem**: Uses `p(n)` and `n(a)` notation but doesn't clearly explain:
- `p(n)` = positional parameters (must be in order)
- `n(a)` = named parameters (can be anywhere after positional)

**Reality**: Named parameters are flags/enums, not indexed.

**Workaround**: Count parameters carefully, use template as reference.

#### Issue 4: Missing Default Values

**Location**: Various parameter descriptions

**Problem**: Doesn't always state default values.

**Example**: What's the default for `standoffHeight`?

**Reality**: Check template file:
```openscad
standoffHeight = 1.0;  // Line 83
```

**Workaround**: Always reference template for defaults.

#### Issue 5: Coordinate System Diagrams

**Location**: Coordinate Systems page

**Problem**: Diagrams are helpful but don't show all three systems side-by-side.

**Reality**: Need to mentally convert between systems.

**Workaround**: Create your own conversion formulas:
```openscad
// PCB to Box conversion:
boxX = pcbX + paddingBack + wallThickness;
boxY = pcbY + paddingLeft + wallThickness;
```

### Corrections Needed in Documentation

| Section | Line/Topic | Issue | Correction |
|---------|-----------|-------|------------|
| Version | Header | Claims v3.0 | Should be v3.3.8+ |
| Cutouts | yappCenter | Unclear when to use | Add: "Use yappCenter for circles to specify center point" |
| Parameters | p(n) vs n(a) | Confusing notation | Add: "p(n) are positional, n(a) are named flags" |
| Defaults | All sections | Missing defaults | Add default values for all optional parameters |
| Examples | Various | Some outdated | Update to v3.3.8 syntax |

---

## Best Practices

### Design Process

1. **Start with PCB dimensions** - Everything else scales from this
2. **Use `yappAllCorners`** for standoffs - Ensures symmetry
3. **Calculate centering** - Show math in comments
4. **Use `yappCenter`** for circles - Simpler positioning
5. **Test incrementally** - Render after each feature addition
6. **Create separate files** - base_only, lid_only, full assembly
7. **Document decisions** - Explain why, not just what

### Code Organization

```openscad
//-----------------------------------------------------------------------
// Project: [Name]
// Description: [Purpose]
// Version: [Number]
//-----------------------------------------------------------------------

include <../YAPP_Box/YAPPgenerator_v3.scad>

//-- Print settings
printBaseShell = true;
printLidShell = true;
printSwitchExtenders = false;

//-- PCB dimensions
pcbLength = 100;  // X-axis
pcbWidth = 50;    // Y-axis
pcbThickness = 1.6;
standoffHeight = 5.0;

//-- Box dimensions
paddingFront = 5;
paddingBack = 5;
paddingRight = 5;
paddingLeft = 5;
wallThickness = 2.0;
baseWallHeight = 15;
lidWallHeight = 12;

//-- PCB definition
pcb = [
  ["Main", pcbLength, pcbWidth, 0, 0, pcbThickness, 
   standoffHeight, standoffDiameter, standoffPinDiameter, standoffHoleSlack]
];

//-- Features (in order of complexity)
pcbStands = [ /* ... */ ];
snapJoins = [ /* ... */ ];
cutoutsLid = [ /* ... */ ];
labelsPlane = [ /* ... */ ];

//-- Generate
YAPPgenerate();
```

### Parameter Selection

| Parameter | Recommended Value | Reasoning |
|-----------|------------------|-----------|
| `paddingFront/Back/Left/Right` | 5mm | Standard clearance for wiring |
| `wallThickness` | 2.0mm | Balance strength/print time |
| `basePlaneThickness` | 1.5mm | Fast printing, adequate rigidity |
| `lidPlaneThickness` | 1.5mm | Sufficient for cutouts |
| `standoffHeight` | 5.0mm | Clearance for solder joints |
| `standoffDiameter` | 7mm | Adequate strength |
| `standoffPinDiameter` | 2.4mm | Fits M2.5 holes with slack |
| `standoffHoleSlack` | 0.4mm | Easy insertion, maintains alignment |
| `ridgeHeight` | 5.0mm | Good overlap for alignment |
| `ridgeSlack` | 0.3mm | Compensates for print tolerances |
| `roundRadius` | 3.0mm | Professional appearance |

### Rendering Strategy

1. **Preview during development** - Fast feedback
   ```bash
   openscad design.scad  # GUI preview
   ```

2. **Test render before STL** - Catch errors
   ```bash
   xvfb-run -a openscad --render -o test.png design.scad
   ```

3. **Generate STLs separately** - Faster, cleaner
   ```bash
   xvfb-run -a openscad -o base.stl base_only.scad
   xvfb-run -a openscad -o lid.stl lid_only.scad
   ```

4. **Create multiple views** - Better documentation
   - Isometric assembly
   - Top view of lid
   - Bottom view of base

---

## Code Examples

### Complete Minimal Box

```openscad
//-----------------------------------------------------------------------
// Minimal YAPPgenerator Box
//-----------------------------------------------------------------------

include <../YAPP_Box/YAPPgenerator_v3.scad>

printBaseShell = true;
printLidShell = true;
printSwitchExtenders = false;

pcbLength = 50;
pcbWidth = 30;
pcbThickness = 1.6;
standoffHeight = 5.0;
standoffDiameter = 7;
standoffPinDiameter = 2.4;
standoffHoleSlack = 0.4;

pcb = [
  ["Main", pcbLength, pcbWidth, 0, 0, pcbThickness, 
   standoffHeight, standoffDiameter, standoffPinDiameter, standoffHoleSlack]
];

paddingFront = 5;
paddingBack = 5;
paddingRight = 5;
paddingLeft = 5;

wallThickness = 2.0;
basePlaneThickness = 1.5;
lidPlaneThickness = 1.5;
baseWallHeight = 10;
lidWallHeight = 8;

ridgeHeight = 5.0;
ridgeSlack = 0.3;
roundRadius = 3.0;

pcbStands = [
  [5, 5, yappAllCorners]
];

snapJoins = [
  [10, 5, yappLeft, yappRight]
];

YAPPgenerate();
```

### Box with Display and Buttons (Cutouts)

```openscad
//-----------------------------------------------------------------------
// Box with Display and Button Cutouts
//-----------------------------------------------------------------------

include <../YAPP_Box/YAPPgenerator_v3.scad>

printBaseShell = true;
printLidShell = true;
printSwitchExtenders = false;

pcbLength = 100;
pcbWidth = 50;
pcbThickness = 1.6;
standoffHeight = 5.0;
standoffDiameter = 7;
standoffPinDiameter = 2.4;
standoffHoleSlack = 0.4;

pcb = [
  ["Main", pcbLength, pcbWidth, 0, 0, pcbThickness, 
   standoffHeight, standoffDiameter, standoffPinDiameter, standoffHoleSlack]
];

paddingFront = 5;
paddingBack = 5;
paddingRight = 5;
paddingLeft = 5;

wallThickness = 2.0;
basePlaneThickness = 1.5;
lidPlaneThickness = 1.5;
baseWallHeight = 15;
lidWallHeight = 12;

ridgeHeight = 5.0;
ridgeSlack = 0.3;
roundRadius = 3.0;

pcbStands = [
  [5, 5, yappAllCorners]
];

snapJoins = [
  [10, 5, yappLeft, yappRight],
  [pcbLength-10, 5, yappLeft, yappRight]
];

// Display: 28x28mm centered
// X: (100-28)/2 = 36mm
// Y: (50-28)/2 = 11mm
cutoutsLid = [
  // Display cutout
  [36, 11, 28, 28, 0, yappRectangle],
  
  // Button cutouts: 18mm diameter, centered at Y=25mm
  [20, 25, 0, 0, 9, yappCircle, 0, 0, yappCenter],
  [50, 25, 0, 0, 9, yappCircle, 0, 0, yappCenter],
  [80, 25, 0, 0, 9, yappCircle, 0, 0, yappCenter]
];

YAPPgenerate();
```

### Box with Push Buttons (Button Guides)

```openscad
//-----------------------------------------------------------------------
// Box with Push Button Guides
//-----------------------------------------------------------------------

include <../YAPP_Box/YAPPgenerator_v3.scad>

printBaseShell = true;
printLidShell = true;
printSwitchExtenders = true;  // Important!

pcbLength = 100;
pcbWidth = 50;
pcbThickness = 1.6;
standoffHeight = 5.0;
standoffDiameter = 7;
standoffPinDiameter = 2.4;
standoffHoleSlack = 0.4;

pcb = [
  ["Main", pcbLength, pcbWidth, 0, 0, pcbThickness, 
   standoffHeight, standoffDiameter, standoffPinDiameter, standoffHoleSlack]
];

paddingFront = 5;
paddingBack = 5;
paddingRight = 5;
paddingLeft = 5;

wallThickness = 2.0;
basePlaneThickness = 1.5;
lidPlaneThickness = 1.5;
baseWallHeight = 15;
lidWallHeight = 12;

ridgeHeight = 5.0;
ridgeSlack = 0.3;
roundRadius = 3.0;

pcbStands = [
  [5, 5, yappAllCorners]
];

snapJoins = [
  [10, 5, yappLeft, yappRight],
  [pcbLength-10, 5, yappLeft, yappRight]
];

// Push buttons for tactile switches on PCB
// Parameters: posx, posy, capLength, capWidth, capRadius, 
//             capAboveLid, switchHeight, switchTravel, poleDiameter
pushButtons = [
  [20, 25, 0, 0, 9, 2, 5, 0.5, 3.5, undef, yappCircle],
  [50, 25, 0, 0, 9, 2, 5, 0.5, 3.5, undef, yappCircle],
  [80, 25, 0, 0, 9, 2, 5, 0.5, 3.5, undef, yappCircle]
];

YAPPgenerate();
```

### Separate Base/Lid Files

**base_only.scad**:
```openscad
include <../YAPP_Box/YAPPgenerator_v3.scad>

printBaseShell = true;
printLidShell = false;
printSwitchExtenders = false;

// ... (copy all parameters from main file)

YAPPgenerate();
```

**lid_only.scad**:
```openscad
include <../YAPP_Box/YAPPgenerator_v3.scad>

printBaseShell = false;
printLidShell = true;
printSwitchExtenders = true;  // If using pushButtons

// ... (copy all parameters from main file)

YAPPgenerate();
```

---

## Troubleshooting

### Common Errors and Solutions

#### Error: "Can't create OpenGL OffscreenView"

**Cause**: No display available in headless environment.

**Solution**:
```bash
sudo apt-get install -y xvfb
xvfb-run -a openscad -o output.stl input.scad
```

#### Error: Rendering takes forever (>5 minutes)

**Cause**: Complex geometry, especially with push buttons.

**Solution**:
- Reduce `renderQuality` during development:
  ```openscad
  renderQuality = 5;  // Instead of 8
  ```
- Use preview mode instead of render for testing
- Render base and lid separately

#### Error: STL file is huge (>10 MB)

**Cause**: High polygon count from complex features.

**Solution**:
- Check if you're rendering both base and lid
- Reduce `renderQuality` slightly
- Simplify geometry (e.g., use cutouts instead of push buttons)

#### Warning: "Normalized tree is null"

**Cause**: Geometry error, usually overlapping features.

**Solution**:
- Check for duplicate entries in arrays
- Verify cutout positions don't overlap
- Ensure standoffs don't intersect with cutouts

#### Issue: Cutout appears in wrong location

**Cause**: Coordinate system mismatch.

**Solution**:
- Verify you're using `yappCoordPCB` (default)
- Check if you need `yappCenter` for circles
- Recalculate centering math

#### Issue: Lid doesn't fit on base

**Cause**: `ridgeSlack` too tight or too loose.

**Solution**:
- Increase `ridgeSlack` to 0.3-0.4mm for looser fit
- Decrease to 0.2mm for tighter fit
- Check printer calibration

#### Issue: Snap joins too tight/loose

**Cause**: Print tolerances or design issue.

**Solution**:
- Adjust snap join width (p(1) parameter)
- Check printer calibration
- Increase/decrease `ridgeSlack`

---

## Performance Optimization

### Rendering Speed Tips

1. **Use lower quality during development**
   ```openscad
   renderQuality = 5;  // Fast
   previewQuality = 3;  // Very fast
   ```

2. **Render parts separately**
   - Base: ~20 seconds
   - Lid: ~30 seconds
   - Both: ~2 minutes

3. **Use preview mode for testing**
   ```bash
   openscad design.scad  # Opens GUI, instant preview
   ```

4. **Disable features during testing**
   ```openscad
   printSwitchExtenders = false;  // Skip button caps
   showPCB = false;  // Skip PCB visualization
   ```

### File Size Optimization

| Feature | Impact on STL Size | Alternative |
|---------|-------------------|-------------|
| Push buttons | +4-5 MB | Use cutouts instead |
| High `renderQuality` | +1-2 MB | Use quality=6-7 instead of 8-10 |
| Complex cutouts | +0.5-1 MB | Simplify shapes |
| Labels | Minimal | Keep |
| Snap joins | Minimal | Keep |

### Workflow Optimization

**Efficient Development Cycle**:

1. **Initial design** (10 min)
   - Define PCB dimensions
   - Add standoffs
   - Add snap joins
   - Quick preview

2. **Add features** (5 min each)
   - Add one feature at a time
   - Preview after each
   - Verify position/size

3. **Final render** (5 min)
   - Set `renderQuality = 8`
   - Generate preview images
   - Create documentation

4. **STL generation** (5 min)
   - Create base_only.scad
   - Create lid_only.scad
   - Generate both STLs

**Total time**: 30-40 minutes for complete project

---

## Summary Checklist

### Before Starting

- [ ] Install OpenSCAD and xvfb
- [ ] Clone YAPP_Box repository
- [ ] Check actual version in source code
- [ ] Review template file
- [ ] Find relevant examples

### During Design

- [ ] Clarify user requirements (cutouts vs. push buttons)
- [ ] Start from template or example
- [ ] Define PCB dimensions first
- [ ] Calculate centering for cutouts
- [ ] Use `yappCenter` for circles
- [ ] Test render after each feature
- [ ] Document calculations in comments

### Before Generating STLs

- [ ] Create separate base/lid files
- [ ] Set correct print flags
- [ ] Verify all positions visually
- [ ] Check for overlapping features
- [ ] Set final `renderQuality`

### Documentation

- [ ] Create README with specifications
- [ ] Explain feature selection rationale
- [ ] Show calculation formulas
- [ ] Include preview images
- [ ] Verify against YAPP documentation

### Delivery

- [ ] base.stl (ready to print)
- [ ] lid.stl (ready to print)
- [ ] Source .scad files
- [ ] Preview images (multiple angles)
- [ ] README.md
- [ ] Technical documentation

---

## Quick Reference

### Essential Commands

```bash
# Install dependencies
sudo apt-get install -y openscad xvfb

# Clone repository
git clone https://github.com/mrWheel/YAPP_Box.git

# Check version
grep "^Version=" YAPPgenerator_v3.scad

# Generate preview
xvfb-run -a openscad --render --imgsize=1920,1440 \
  --colorscheme=Tomorrow --camera=0,0,0,60,0,30,250 \
  -o preview.png design.scad

# Generate STL
xvfb-run -a openscad -o output.stl input.scad

# Search for feature
grep -n "featureName" YAPPgenerator_v3.scad
grep -r "featureName" examples/
```

### Key Files

```
YAPP_Box/
├── YAPPgenerator_v3.scad    # Main library
├── YAPP_Template_v3.scad    # Complete template
├── examples/
│   ├── YAPP_Demo_buttons_v31.scad
│   ├── YAPP_Demo_cutouts_v30.scad
│   └── YAPP_Demo_DisplayMount_v31.scad
└── README.md
```

### Common Patterns

```openscad
// PCB definition
pcb = [["Main", 100, 50, 0, 0, 1.6, 5.0, 7, 2.4, 0.4]];

// Standoffs (4 corners)
pcbStands = [[5, 5, yappAllCorners]];

// Snap joins
snapJoins = [[10, 5, yappLeft, yappRight]];

// Rectangle cutout (centered)
cutoutsLid = [[(pcbLength-width)/2, (pcbWidth-height)/2, width, height, 0, yappRectangle]];

// Circle cutout (centered)
cutoutsLid = [[x, y, 0, 0, radius, yappCircle, 0, 0, yappCenter]];

// Push button
pushButtons = [[x, y, 0, 0, radius, 2, 5, 0.5, 3.5, undef, yappCircle]];
```

---

## Version History

- **v1.0** (2025-11-02): Initial playbook based on YAPPgenerator v3.3.8

---

## Contributing to This Playbook

When you encounter new issues or discover better approaches:

1. Document the problem clearly
2. Explain what didn't work and why
3. Provide the solution that worked
4. Include code examples
5. Update relevant sections

This playbook is a living document. Keep it updated with new learnings!

