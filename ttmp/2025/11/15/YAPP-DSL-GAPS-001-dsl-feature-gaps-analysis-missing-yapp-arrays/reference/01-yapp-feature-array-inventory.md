---
Title: YAPP Feature Array Inventory
Ticket: YAPP-DSL-GAPS-001
Status: active
Topics:
  - yapp
  - dsl
  - analysis
  - features
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
  - Path: ../../../YAPP_Template_v3.scad
    Note: Canonical YAPP feature array definitions
  - Path: ../../../examples/YAPP_Demo_buttons2_v31.scad
    Note: Example file user is trying to replicate
ExternalSources: []
Summary: Complete inventory of YAPP v3 feature arrays with parameter specifications
LastUpdated: 2025-11-16
---

# YAPP Feature Array Inventory

## Overview

This document catalogs all feature arrays supported by YAPPgenerator v3, extracted from `YAPP_Template_v3.scad`. Each array is documented with its parameter order, flags, and coordinate system defaults.

## Currently Implemented in DSL

✅ **pcb_stands** (pcbStands)
- Position: x, y
- Optional: height, pcb_gap, diameter, pin_diameter, hole_slack, fillet_radius, pin_length
- Flags: yappBoth/yappLidOnly/yappBaseOnly, yappPin/yappHole/yappTopPin, corner placement, coordinate system, yappNoFillet, yappPCBName, yappSelfThreading

✅ **connectors** (connectors)
- Position: x, y
- Required: stand_height, screw_d, screw_head_d, insert_d, outside_d
- Optional: insert_depth, pcb_gap, fillet_radius
- Flags: corner placement, coordinate system, yappNoFillet, yappCountersink, yappPCBName, yappThroughLid, yappSelfThreading, yappNoInternalFillet

✅ **snap_joins** (snapJoins)
- Position: pos (along edge)
- Required: width, side (yappLeft/yappRight/yappFront/yappBack)
- Flags: yappOrigin/yappCenter, yappSymmetric, yappRectangle

✅ **cutouts** (cutoutsBase/Lid/Front/Back/Left/Right)
- Position: from_back, from_left
- Required: width, length, radius, shape
- Optional: depth, angle
- Flags: yappPolygonDef, yappMaskDef, coordinate system, yappOrigin/yappCenter, yappGlobalOrigin/yappAltOrigin, yappPCBName, yappFromInside

✅ **push_buttons** (pushButtons)
- Position: x, y
- Required: cap (length, width, radius), lid (protrusion), switch (height, travel, pole_diameter)
- Optional: angle, fillet_radius, shape, polygon presets, coordinate, origin, no_fillet
- Flags: coordinate system, origin, yappNoFillet, yappPCBName

## Missing from DSL

### ❌ boxMounts (External Mounting Tabs)

**Purpose:** Mounting tabs on the outside of the box for wall/surface mounting

**Parameters:**
- p(0) = pos : position along wall (can be vector [pos, offset])
- p(1) = screwDiameter
- p(2) = width of opening (0 = circular hole)
- p(3) = height
- p(4) = filletRadius (optional)

**Flags:**
- n(a) = { yappLeft | yappRight | yappFront | yappBack } : one or more sides
- n(b) = { yappNoFillet }
- n(c) = { <yappBase>, yappLid } : which shell part
- n(d) = { yappCenter } : center position
- n(e) = { <yappGlobalOrigin>, yappAltOrigin }

**Use case:** Wall mounting, DIN rail clips, external fastening points

**User observation:** "little plates sticking out to fasten it" - This is likely boxMounts!

### ❌ lightTubes (LED Light Pipes)

**Purpose:** Light pipes from PCB LEDs through the lid

**Parameters:**
- p(0) = posx
- p(1) = posy
- p(2) = tubeLength
- p(3) = tubeWidth
- p(4) = tubeWall
- p(5) = gapAbovePcb
- p(6) = { yappCircle | yappRectangle } : tubeType
- p(7) = lensThickness (optional, 0 = open hole)
- p(8) = Height to top of PCB (optional)
- p(9) = filletRadius (optional)

**Flags:**
- n(a) = { <yappCoordPCB> | yappCoordBox | yappCoordBoxInside }
- n(b) = { <yappGlobalOrigin>, yappAltOrigin }
- n(c) = { yappNoFillet }
- n(d) = [yappPCBName, "XXX"]

**Use case:** Status LEDs, indicator lights visible through enclosure

### ❌ labelsPlane (Text Labels)

**Purpose:** Embossed or engraved text labels on enclosure surfaces

**Parameters:**
- p(0) = posx
- p(1) = posy/z
- p(2) = rotation degrees CCW
- p(3) = depth : positive = engrave (remove), negative = emboss (add)
- p(4) = { yappLeft, yappRight, yappFront, yappBack, yappLid, yappBase } : plane
- p(5) = font
- p(6) = size
- p(7) = "label text"
- p(8) = Expand (optional, makes text bolder)
- p(9) = Direction (optional, text orientation)
- p(10) = Horizontal alignment (optional)
- p(11) = Vertical alignment (optional)
- p(12) = Character spacing multiplier (optional)

**Use case:** Product names, port labels, version numbers, warnings

### ❌ ridgeExt* (Ridge Extensions)

**Purpose:** Extension from lid into case for split openings at various heights

**Arrays:** ridgeExtLeft, ridgeExtRight, ridgeExtFront, ridgeExtBack

**Parameters:**
- p(0) = pos
- p(1) = width
- p(2) = height : Where to relocate the seam
  - yappCoordPCB = Above (positive) the PCB
  - yappCoordBox = Above (positive) the bottom of shell (outside)

**Flags:**
- n(a) = { <yappOrigin>, yappCenter }
- n(b) = { <yappCoordPCB> | yappCoordBox | yappCoordBoxInside }
- n(c) = { yappAltOrigin, <yappGlobalOrigin> }
- n(d) = [yappPCBName, "XXX"]

**Note:** Snaps should not be placed on ridge extensions

**Use case:** Split-level enclosures, hinged openings, cable pass-throughs

**User observation:** "hinges" and "side of the box are open" - This might be ridgeExt creating split openings!

### ❌ displayMounts (Display Module Mounting)

**Purpose:** Cutout in lid with mounting posts for LCD/OLED displays

**Parameters:**
- p(0) = posx
- p(1) = posy
- p(2) = displayWidth
- p(3) = displayHeight
- p(4) = pinInsetH (horizontal inset of mounting hole)
- p(5) = pinInsetV (vertical inset)
- p(6) = pinDiameter
- p(7) = postOverhang
- p(8) = walltoPCBGap
- p(9) = pcbThickness
- p(10) = windowWidth
- p(11) = windowHeight
- p(12) = windowOffsetH
- p(13) = windowOffsetV
- p(14) = bevel (45° bevel on opening)
- p(15) = rotation (optional)
- p(16) = snapDiameter (optional)
- p(17) = lidThickness (optional)

**Flags:**
- n(a) = { <yappOrigin>, yappCenter }
- n(b) = { <yappCoordBox> | yappCoordPCB | yappCoordBoxInside }
- n(c) = { <yappGlobalOrigin>, yappAltOrigin }
- n(d) = [yappPCBName, "XXX"]
- n(e) = {yappSelfThreading}

**Use case:** OLED displays, LCD screens, touchscreens

### ❌ pcb (Multi-PCB Support)

**Purpose:** Define multiple PCBs in one enclosure

**Parameters:**
- p(0) = name
- p(1) = length
- p(2) = width
- p(3) = posx
- p(4) = posy
- p(5) = thickness
- p(6) = standoff_Height (from base inside or lid inside if negative)
- p(7) = standoff_Diameter
- p(8) = standoff_PinDiameter
- p(9) = standoff_HoleSlack (optional)

**Use case:** Multiple boards, stacked PCBs, separate control/power boards

**Current DSL limitation:** Only supports single "Main" PCB via global pcb.* fields

## Missing Flags on Implemented Features

### pcb_stands Missing Flags

Current DSL schema missing:
- ❌ **yappBoth/yappLidOnly/yappBaseOnly** : Which shell part gets the standoff
- ❌ **yappPin/yappHole/yappTopPin** : Standoff treatment (pin vs hole)
- ❌ **yappAllCorners/yappFrontLeft/etc** : Auto-place at corners
- ❌ **yappPCBName** : Multi-PCB support
- ❌ **yappSelfThreading** : Self-threading holes

**User observation:** "4 stalactites coming from top to pcb pillars (of which there are 4), but the scad one has only 2"

This suggests the SCAD file is using corner placement flags (yappAllCorners or specific corners) to auto-generate 4 standoffs from 1 definition, OR using yappBoth to create standoffs in both base and lid.

### connectors Missing Flags

Current DSL schema missing:
- ❌ **yappAllCorners/corner placement** : Auto-place connectors
- ❌ **yappCountersink** : Countersunk screw heads
- ❌ **yappThroughLid** : Reverse screw direction
- ❌ **yappNoInternalFillet** : Fillet control
- ❌ **yappPCBName** : Multi-PCB support

### cutouts Missing Flags

Current DSL schema missing:
- ❌ **yappPolygonDef** : Custom polygon shapes
- ❌ **yappMaskDef** : Honeycomb/bar masks for ventilation
- ❌ **yappRing** : Ring-shaped cutouts
- ❌ **yappSphere** : Spherical cutouts
- ❌ **yappAltOrigin** : Alternate origin for specific faces
- ❌ **yappFromInside** : Cut direction
- ❌ **yappPCBName** : Multi-PCB support

### snap_joins Missing Flags

Current DSL schema missing:
- ❌ **yappOrigin/yappCenter** : Position reference
- ❌ **yappSymmetric** : Mirror on opposite side
- ❌ **yappRectangle** : Diamond-shaped snaps

### push_buttons Missing Flags

Current DSL schema missing:
- ❌ **yappAltOrigin** : Alternate origin
- ❌ **yappPCBName** : Multi-PCB support
- ❌ **snapSlack parameter** : Already in schema but not documented in YAPP Template v3.scad comments (v3.1 addition)

## Analysis of YAPP_Demo_buttons2_v31.scad

**File defines:**
- 1 pcbStand at [5, 5] with no flags
- 2 cutouts (front, back) using shellWidth/shellHeight expressions
- 6 pushButtons with various shapes

**User sees:**
- 4 PCB pillars (but file only defines 1!)
- Hinges/open sides
- Little plates sticking out
- Holes in one side

**Hypothesis:**

1. **4 pillars from 1 definition:** YAPP might auto-generate corner standoffs when only one is defined (legacy behavior), OR the user's DSL YAML has 4 explicit pcb_stands entries

2. **"Hinges":** Likely **ridgeExt*** creating split openings that look like hinge points

3. **"Little plates sticking out":** Definitely **boxMounts** - external mounting tabs

4. **"Holes in one side":** The cutouts defined in cutoutsFront/cutoutsBack

5. **"Sides are open":** Could be:
   - Ridge extensions creating splits
   - Missing walls due to incorrect baseWallHeight/lidWallHeight
   - Cutouts that are too large

## Critical Gaps for Production Use

### Priority 1: Essential for Functional Enclosures

1. **boxMounts** - External mounting (wall mount, DIN rail)
2. **Corner placement flags** - Auto-generate standoffs at corners (yappAllCorners, yappFrontLeft, etc.)
3. **Shell part flags** - Control which part gets feature (yappBoth, yappLidOnly, yappBaseOnly)
4. **Multi-PCB support** - pcb array + yappPCBName references

### Priority 2: Common Features

5. **lightTubes** - LED indicators
6. **labelsPlane** - Text labels
7. **ridgeExt*** - Split openings, hinged sections
8. **displayMounts** - LCD/OLED mounting

### Priority 3: Advanced Features

9. **Mask support** - Honeycomb ventilation (yappMaskDef)
10. **Custom polygons** - yappPolygonDef for cutouts
11. **Advanced cutout shapes** - yappRing, yappSphere
12. **Hook functions** - Custom 3D objects (hookLidInside, etc.)

## Impact on User's Use Case

**User wants to replicate YAPP_Demo_buttons2_v31.scad:**

**What works:**
- ✅ Push buttons (6 buttons with various shapes)
- ✅ Basic cutouts (front/back)
- ✅ PCB stands (but missing corner auto-placement)

**What's missing:**
- ❌ boxMounts (the "little plates")
- ❌ ridgeExt (potential "hinges")
- ❌ Corner placement flags (causing 4 vs 1 pillar discrepancy)
- ❌ Shell part flags (yappBoth/yappLidOnly/yappBaseOnly)

**Recommendation:** Add boxMounts and corner placement flags as next priority to achieve parity with common YAPP examples.
