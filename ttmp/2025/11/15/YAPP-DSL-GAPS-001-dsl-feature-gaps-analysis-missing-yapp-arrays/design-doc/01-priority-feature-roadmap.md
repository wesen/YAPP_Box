---
Title: Priority Feature Roadmap
Ticket: YAPP-DSL-GAPS-001
Status: active
Topics:
  - yapp
  - dsl
  - roadmap
  - features
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
  - Path: ../reference/01-yapp-feature-array-inventory.md
    Note: Complete YAPP feature catalog
  - Path: ../analysis/01-user-case-study-buttons2-demo-discrepancies.md
    Note: User's specific gaps
ExternalSources: []
Summary: Prioritized roadmap for adding missing YAPP features to DSL
LastUpdated: 2025-11-16
---

# Priority Feature Roadmap

## Executive Summary

The YAPP DSL module system (YAPP-MODULE-SYSTEM-001) successfully implemented 5 core feature modules with schema-driven validation and code generation. However, analysis reveals significant gaps preventing parity with real YAPP examples.

**Current coverage:** 5 of 13+ YAPP feature arrays (38%)
**Missing critical features:** Corner placement flags, shell part control, boxMounts, ridgeExt, lightTubes, labels, displayMounts

This document prioritizes missing features based on:
- Frequency in YAPP examples
- User requirements
- Implementation complexity
- Dependency chains

## Priority 1: Flag Extensions (Weeks 1-2)

### 1.1 Corner Placement Flags

**Arrays affected:** pcb_stands, connectors

**Missing flags:**
- `yappAllCorners` - Auto-place at all 4 corners
- `yappFrontLeft`, `yappFrontRight`, `yappBackLeft`, `yappBackRight` - Specific corners

**Impact:** Currently requires 4 manual definitions; with flags, 1 definition generates 4 standoffs

**Example SCAD:**
```openscad
pcbStands = [
  [5, 5, yappAllCorners]  // Generates 4 standoffs automatically
];
```

**DSL implementation:**
```yaml
fields:
  corner_placement:
    type: string
    enum: [all_corners, front_left, front_right, back_left, back_right, none]
    desc: Auto-place at specified corners (combines with x/y for offset)
```

**Builder logic:**
- If `corner_placement` is set, generate 4 entries with computed positions
- Base positions on pcb dimensions minus offsets
- Append appropriate yapp flag to each entry

**Effort:** 2-3 days (schema update, builder logic, tests)

### 1.2 Shell Part Control

**Arrays affected:** pcb_stands, connectors

**Missing flags:**
- `yappBoth` (default) - Feature in both base and lid
- `yappLidOnly` - Feature only in lid
- `yappBaseOnly` - Feature only in base

**Impact:** Can't control which shell part gets the feature

**Example use case:**
```openscad
pcbStands = [
  [5, 5, yappBaseOnly]  // Standoff only in base, no lid pin
];
```

**DSL implementation:**
```yaml
fields:
  shell_part:
    type: string
    enum: [both, lid_only, base_only]
    default: both
    desc: Which shell part receives this feature
```

**Builder logic:**
- Map enum to YAPP flag
- Append to params array after positional parameters

**Effort:** 1 day (straightforward enum mapping)

### 1.3 Standoff Treatment Flags

**Array affected:** pcb_stands

**Missing flags:**
- `yappPin` (default) - Pin on base, hole on lid
- `yappHole` - Hole on both parts
- `yappTopPin` - Hole on base, pin on lid

**Impact:** Can't customize standoff pin/hole configuration

**DSL implementation:**
```yaml
fields:
  treatment:
    type: string
    enum: [pin, hole, top_pin]
    default: pin
    desc: Standoff pin/hole configuration
```

**Effort:** 1 day

### 1.4 Snap Join Positioning Flags

**Array affected:** snap_joins

**Missing flags:**
- `yappOrigin` (default) vs `yappCenter` - Position reference
- `yappSymmetric` - Mirror on opposite side
- `yappRectangle` - Diamond-shaped snap (vs default round)

**DSL implementation:**
```yaml
fields:
  position_ref:
    type: string
    enum: [origin, center]
    default: origin
    desc: Position reference point
  
  symmetric:
    type: bool
    default: false
    desc: Mirror snap on opposite side
  
  snap_shape:
    type: string
    enum: [round, rectangle]
    default: round
    desc: Snap join shape
```

**Effort:** 1-2 days

**Total Priority 1 effort:** 1-2 weeks

## Priority 2: New Feature Arrays (Weeks 3-6)

### 2.1 boxMounts (External Mounting Tabs)

**Frequency:** Common in production boxes (wall mounts, DIN rail)

**Parameters:**
- pos (position along wall, can be vector [pos, offset])
- screwDiameter
- width (0 = circular hole)
- height
- filletRadius (optional)

**Flags:**
- sides: yappLeft/yappRight/yappFront/yappBack (multiple)
- shell_part: yappBase/yappLid
- yappCenter, yappNoFillet, yappAltOrigin

**DSL schema structure:**
```yaml
module: box_mounts
fields:
  pos: { type: number, required: true }
  screw_diameter: { type: number, required: true }
  width: { type: number, required: true }
  height: { type: number, required: true }
  fillet_radius: { type: number }
  sides:
    type: array
    required: true
    desc: Which sides get the mount
  shell_part:
    type: string
    enum: [base, lid]
    default: base
```

**Complexity:** Medium (array of sides, position vectors)

**Effort:** 3-4 days

### 2.2 lightTubes (LED Light Pipes)

**Frequency:** Common for status indicators

**Parameters:**
- posx, posy
- tubeLength, tubeWidth, tubeWall
- gapAbovePcb
- tubeType (yappCircle/yappRectangle)
- lensThickness (optional)
- heightToPCB (optional)
- filletRadius (optional)

**Flags:**
- coordinate system, origin, yappNoFillet, yappPCBName

**Complexity:** Medium (similar to push_buttons)

**Effort:** 2-3 days

### 2.3 labelsPlane (Text Labels)

**Frequency:** Very common (product names, port labels)

**Parameters:**
- posx, posy/z
- rotation
- depth (positive = engrave, negative = emboss)
- plane (yappLeft/Right/Front/Back/Lid/Base)
- font
- size
- text
- expand, direction, alignment, spacing (optional)

**Complexity:** HIGH (text rendering, font handling, alignment)

**Effort:** 5-7 days (most complex feature)

### 2.4 ridgeExt* (Ridge Extensions)

**Frequency:** Uncommon but enables advanced designs

**Arrays:** ridgeExtLeft, ridgeExtRight, ridgeExtFront, ridgeExtBack

**Parameters:**
- pos
- width
- height (seam relocation)

**Flags:**
- yappOrigin/yappCenter, coordinate system, yappAltOrigin, yappPCBName

**Complexity:** Medium (4 separate arrays, coordinate math)

**Effort:** 3-4 days

### 2.5 displayMounts (Display Mounting)

**Frequency:** Common for OLED/LCD projects

**Parameters:** 17 parameters (most complex feature)
- Display dimensions, window dimensions, offsets
- Pin insets, diameters, overhangs
- PCB thickness, wall gap, bevel

**Flags:**
- yappOrigin/yappCenter, coordinate system, yappAltOrigin, yappPCBName, yappSelfThreading

**Complexity:** VERY HIGH (most parameters of any feature)

**Effort:** 7-10 days

**Total Priority 2 effort:** 3-6 weeks

## Priority 3: Advanced Features (Future)

### 3.1 Multi-PCB Support

**Current limitation:** DSL assumes single "Main" PCB

**YAPP supports:**
```openscad
pcb = [
  ["Main", 100, 50, 0, 0, 1.6, 5, 7, 2.4, 0.4],
  ["Display", 40, 30, 20, 10, 1.0, 15, 5, 2.0, 0.3]
];
```

**DSL would need:**
```yaml
pcbs:
  - name: Main
    length: 100
    width: 50
    # ...
  - name: Display
    length: 40
    width: 30
    # ...

features:
  pcb_stands:
    - x: 5
      y: 5
      pcb: Display  # Reference specific PCB
```

**Complexity:** HIGH (affects all features, coordinate calculations)

**Effort:** 2-3 weeks

### 3.2 Mask Support (Ventilation Patterns)

**Masks:** maskHoneycomb, maskHexCircles, maskBars, maskOffsetBars

**Example:**
```openscad
cutoutsLid = [
  [20, 20, 30, 30, 0, yappRectangle, maskHoneycomb]
];
```

**DSL implementation:**
```yaml
cutouts:
  - face: lid
    from_left: 20
    from_back: 20
    width: 30
    length: 30
    shape: rectangle
    mask: honeycomb  # New field
```

**Complexity:** Medium (mask parameter passing)

**Effort:** 2-3 days

### 3.3 Custom Polygons

**Current:** Only preset polygons (arrow, hexagon, etc.)

**YAPP supports:**
```openscad
myShape = [yappPolygonDef, [[-0.5,-0.5], [0,0.5], [0.5,-0.5]]];

cutoutsLid = [
  [20, 20, 10, 10, 0, yappPolygon, myShape]
];
```

**DSL would need:**
```yaml
custom_shapes:
  my_triangle:
    vertices:
      - [-0.5, -0.5]
      - [0, 0.5]
      - [0.5, -0.5]

cutouts:
  - shape: polygon
    polygon_def: my_triangle
```

**Complexity:** HIGH (shape definition DSL, validation)

**Effort:** 1-2 weeks

### 3.4 Hook Functions

**Purpose:** Custom 3D objects via OpenSCAD modules

**Functions:** hookLidInside, hookLidOutside, hookBaseInside, hookBaseOutside, hookBasePre, hookLidPre

**Challenge:** DSL can't express arbitrary 3D geometry

**Possible approach:** Allow raw OpenSCAD snippets
```yaml
hooks:
  lid_inside: |
    translate([10, 10, 0])
      cube([5, 5, 2]);
```

**Complexity:** VERY HIGH (security, validation, escaping)

**Effort:** 3-4 weeks

## Implementation Strategy

### Phase 1: Flag Extensions (Priority 1)

**Goal:** Achieve parity with common YAPP patterns

**Deliverables:**
1. Corner placement flags (pcb_stands, connectors)
2. Shell part flags (pcb_stands, connectors)
3. Standoff treatment flags (pcb_stands)
4. Snap positioning flags (snap_joins)

**Timeline:** 2 weeks

**Success criteria:**
- Can replicate YAPP_Demo_RealBox_v31.scad with DSL
- Single pcb_stand definition can generate 4 corners
- Can control base-only vs lid-only standoffs

### Phase 2: Common Arrays (Priority 2)

**Goal:** Support 80%+ of production use cases

**Deliverables:**
1. boxMounts module
2. lightTubes module
3. labelsPlane module (basic text only)
4. ridgeExt modules

**Timeline:** 4-6 weeks

**Success criteria:**
- Can add wall mounting tabs
- Can add LED indicators
- Can add basic text labels
- Can create split openings

### Phase 3: Advanced Features (Priority 3)

**Goal:** Complete YAPP feature coverage

**Deliverables:**
1. displayMounts module
2. Multi-PCB support
3. Mask support for cutouts
4. Custom polygon definitions

**Timeline:** 6-8 weeks

**Success criteria:**
- Can mount OLED displays
- Can define multiple PCBs
- Can add honeycomb ventilation
- Can use custom polygon shapes

## Effort Summary

| Priority | Features | Effort | Timeline |
|----------|----------|--------|----------|
| P1 | Flag extensions | 2 weeks | Weeks 1-2 |
| P2 | 4 new modules | 4-6 weeks | Weeks 3-8 |
| P3 | Advanced features | 6-8 weeks | Weeks 9-16 |
| **Total** | **Full parity** | **12-16 weeks** | **4 months** |

## Recommendations

### For Immediate User Needs

1. **Workaround for 4 pillars:** Use expressions in DSL
   ```yaml
   pcb_stands:
     - x: 3
       y: 3
     - x: pcb_length - 3
       y: 3
     - x: 3
       y: pcb_width - 3
     - x: pcb_length - 3
       y: pcb_width - 3
   ```

2. **For boxMounts:** Manually edit generated SCAD to add boxMounts array (temporary)

3. **For "hinges":** Verify if user actually needs ridgeExt or if it's visual misinterpretation

### For Project Planning

1. **Quick win:** Implement Priority 1 flags (2 weeks) - Biggest usability improvement
2. **boxMounts next:** If user needs external mounting (3-4 days)
3. **Defer advanced:** Multi-PCB, custom polygons, hooks until user demand proven

### For Module System Validation

Current module system architecture handles all these features well:
- ✅ Schema-driven validation scales to complex parameters
- ✅ Code generation works for nested objects (proven with push_buttons)
- ✅ Flag handling pattern established (shape, coordinate, origin)
- ✅ Auto-registration works

**No architectural changes needed** - Just add more schemas and builders following established patterns.

## Next Actions

1. **Confirm user needs:** Get user's actual DSL YAML and compare with SCAD output
2. **Prioritize flags:** Start with corner placement (highest ROI)
3. **Create tickets:** One ticket per Priority 1 flag group
4. **Implement incrementally:** Test each flag addition with real examples
5. **Update docs:** Module authoring guide already covers patterns needed

## Open Questions

1. **Corner placement math:** How does YAPP compute corner positions from single [x,y] + yappAllCorners?
   - Need to study YAPPgenerator_v3.scad implementation
   - Likely: x,y become offsets from each corner

2. **boxMounts position vectors:** Parameter p(0) can be scalar or vector [pos, offset]
   - How to represent in YAML? Single number or array?
   - Suggest: `pos: 10` or `pos: [10, 2]` (YAML allows both)

3. **ridgeExt coordinate systems:** Height parameter meaning changes based on coordinate flag
   - Need careful documentation
   - Might need validation (height must be positive for yappCoordBox)

4. **Label fonts:** YAPP accepts font names as strings
   - Need to document available fonts
   - OpenSCAD font support varies by platform

5. **Mask offsets:** yappMaskDef can have [mask, hOffset, vOffset, rotation]
   - Complex parameter structure
   - Might need nested object in schema
