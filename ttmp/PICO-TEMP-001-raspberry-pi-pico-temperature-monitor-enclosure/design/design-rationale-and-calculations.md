---
Title: Design Rationale and Calculations
Ticket: PICO-TEMP-001
Status: active
Topics:
    - enclosures
    - yappgenerator
    - hardware
DocType: design-doc
Intent: long-term
Owners:
    - manuel
RelatedFiles:
    - path: pico_temp_monitor_box.scad
      note: Main design file with complete assembly
    - path: pico_temp_monitor_base_only.scad
      note: Base shell for 3D printing
    - path: pico_temp_monitor_lid_only.scad
      note: Lid shell for 3D printing
ExternalSources:
    - https://mrwheel-docs.gitbook.io/yappgenerator_en/
    - https://github.com/mrWheel/YAPP_Box
Summary: "Design rationale, calculations, and decision documentation for Raspberry Pi Pico temperature monitor enclosure"
LastUpdated: 2025-11-05T16:20:44.898980713-05:00
---

# Design Rationale and Calculations

## Executive Summary

This document explains the design decisions and calculations for a parametric 3D-printable enclosure for a Raspberry Pi Pico-based temperature monitoring device. The enclosure was designed using YAPPgenerator v3.3.8, a parametric OpenSCAD library for creating project boxes based on PCB dimensions.

**Key Design Goals:**
- Accommodate all components with adequate clearance
- Provide tool-free snap-fit assembly
- Enable easy access to display and button
- Allow temperature probe cable routing
- Maintain printability without supports

## Problem Statement

The project requires an enclosure that houses:
1. Raspberry Pi Pico (51x21mm) mounted on a custom PCB via socket
2. 1.3" OLED display module (40x30mm PCB, 36x17mm visible area)
3. Panel-mount push button (28mm diameter, 45mm threaded barrel)
4. Temperature probe with 8mm probe diameter and 6mm cable

**Constraints:**
- Pico + socket stack height: ~12mm (1mm PCB + 8.5mm headers + socket depth)
- Display must be visible through lid
- Button must be accessible from front panel
- Temperature probe cable must exit cleanly
- Box must be 3D printable without supports

## Proposed Solution

### Overall Approach

Use YAPPgenerator v3 to create a parametric box design based on a custom PCB footprint. The PCB serves as the reference coordinate system for all features.

**Box Type:** Type 0 (all edges rounded) for professional appearance and improved printability.

**Assembly Method:** Snap joins for tool-free assembly, with ridge overlap for alignment and dust resistance.

## Design Decisions

### 1. PCB Dimensions

**Decision:** Custom PCB sized at 80mm x 60mm

**Rationale:**
- Raspberry Pi Pico: 51mm x 21mm
- OLED Display Module: 40mm x 30mm
- Need space for:
  - Component placement and routing
  - Mounting holes at corners (inset 5mm)
  - Wire routing and connections
  - Clearance from walls (5mm padding on all sides)

**Calculation:**
```
Minimum PCB length = Pico length + margin = 51 + 29 = 80mm
Minimum PCB width = Display width + margin = 40 + 20 = 60mm
```

The 80x60mm size provides adequate space without being oversized.

### 2. Box Height

**Decision:** Base wall height = 20mm, Lid wall height = 15mm

**Rationale:**
Total internal height needed:
```
Base plane thickness:        1.5mm
Standoff height:            5.0mm
PCB thickness:              1.6mm
Pico stack height:         ~12.0mm (Pico + headers + socket)
Clearance to lid:           ~5.0mm (for wiring)
Lid plane thickness:        1.5mm
--------------------------------
Total required:            ~26.6mm

Actual provided:
Base (20mm) + Lid (15mm) - Ridge overlap (5mm) = 30mm internal
30mm - 1.5mm (base plane) - 1.5mm (lid plane) = 27mm clear height
```

This provides adequate clearance with a small safety margin.

### 3. Wall Thickness

**Decision:** 2.0mm walls, 1.5mm base/lid planes

**Rationale:**
- **2.0mm walls:** Standard for 3D printed enclosures
  - Provides adequate strength
  - Prints well with 0.4mm nozzle (5 perimeters at 0.4mm)
  - Not too thick (saves material and print time)
  
- **1.5mm planes:** Thinner than walls
  - Reduces print time
  - Adequate rigidity for this size box
  - Allows cutouts without excessive depth

### 4. Standoff Configuration

**Decision:** 4 corner standoffs with pins, 7mm diameter, 2.4mm pin diameter

**Rationale:**
- **Corner placement:** Using `yappAllCorners` ensures symmetry
- **5mm inset:** Provides clearance from rounded corners
- **7mm diameter:** Standard for M2.5 mounting holes with adequate strength
- **2.4mm pins:** Fits M2.5 holes (2.5mm) with 0.4mm slack for print tolerance
- **5mm height:** Provides clearance for solder joints and bottom components

### 5. Display Cutout

**Decision:** 36mm x 17mm rectangular cutout, centered on PCB

**Calculation:**
```
Display visible area: 36mm x 17mm
PCB dimensions: 80mm x 60mm

Center X = (80 - 36) / 2 = 22mm from back
Center Y = (60 - 17) / 2 = 21.5mm from left
```

**Rationale:**
- Cutout matches visible display area exactly
- Centered positioning is aesthetically pleasing
- Display PCB (40x30mm) can be mounted from inside with mounting holes
- Display bezel will cover any gaps

### 6. Button Cutout

**Decision:** 28mm diameter circular cutout on front panel, centered vertically

**Position:**
```
Y-position (vertical): pcbWidth / 2 = 30mm (centered)
Z-position (height): 15mm from base
```

**Rationale:**
- 28mm matches button cap diameter
- Centered vertically for ergonomic access
- 15mm from base provides clearance for threaded barrel (45mm total length)
- Front panel placement makes button easily accessible

### 7. Cable Grommet

**Decision:** 10mm diameter circular cutout in base, back-right corner

**Position:**
```
X = pcbLength - 10 = 70mm from back
Y = pcbWidth - 10 = 50mm from left
```

**Rationale:**
- 10mm diameter accommodates 8mm probe with clearance
- Back-right corner position keeps cable out of the way
- Low profile exit (in base) prevents cable interference with lid
- Using `yappCenter` flag for intuitive positioning

### 8. Snap Joins

**Decision:** 2 pairs of snap joins on left and right sides

**Position:**
```
Pair 1: 15mm from back
Pair 2: 65mm from back (pcbLength - 15)
```

**Rationale:**
- Left/right placement distributes force evenly
- 2 pairs provide adequate holding force
- Positioned away from corners for strength
- 5mm width provides good engagement

### 9. Ridge Parameters

**Decision:** 
- Ridge height: 5.0mm
- Ridge slack: 0.3mm
- Ridge gap: 0.5mm

**Rationale:**
- **5.0mm height:** Adequate overlap for alignment and dust resistance
- **0.3mm slack:** Compensates for 3D printing tolerances (typical FDM tolerance)
- **0.5mm gap:** Prevents base ridge from bottoming out in lid

### 10. Rounded Corners

**Decision:** 3.0mm corner radius (boxType = 0)

**Rationale:**
- Professional appearance
- Improved printability (no sharp corners)
- Easier to handle (no sharp edges)
- Standard radius for this size enclosure

## Alternatives Considered

### Alternative 1: Smaller PCB (51x30mm)

**Rejected because:**
- Insufficient space for OLED display (40mm wide)
- No room for additional circuitry
- Cramped layout would make assembly difficult

### Alternative 2: Taller Box (baseWallHeight = 25mm)

**Rejected because:**
- Unnecessary height (current design has 5mm clearance)
- Would waste material and increase print time
- Makes box less compact

### Alternative 3: Screw-Based Assembly

**Rejected because:**
- Requires screws and assembly tools
- More complex design (screw bosses, countersinks)
- Snap joins provide adequate holding force for this application

### Alternative 4: Push Button Guides (tactile switches)

**Rejected because:**
- User specified panel-mount button with threaded barrel
- Push button guides are for PCB-mounted tactile switches
- Would add complexity and render time

### Alternative 5: Display Mounting Clips

**Rejected because:**
- Display PCB has mounting holes (can use M2 screws)
- Clips would add complexity
- Hot glue or double-sided tape is simpler alternative

## Implementation Plan

### Phase 1: Design Verification ✓
- [x] Create parametric design in OpenSCAD
- [x] Verify all dimensions and clearances
- [x] Test render for errors
- [x] Check coordinate system calculations

### Phase 2: File Generation ✓
- [x] Create separate base_only.scad file
- [x] Create separate lid_only.scad file
- [x] Generate preview images
- [x] Generate STL files

### Phase 3: Documentation ✓
- [x] Create README with specifications
- [x] Document design rationale
- [x] Create component specifications document
- [x] Update ticket metadata

### Phase 4: Validation (Future)
- [ ] 3D print prototype
- [ ] Test fit all components
- [ ] Verify snap join strength
- [ ] Check button accessibility
- [ ] Validate cable routing

### Phase 5: Iteration (If Needed)
- [ ] Adjust dimensions based on physical testing
- [ ] Modify cutout sizes if needed
- [ ] Fine-tune snap join parameters
- [ ] Add ventilation if heat is an issue

## Key Calculations Summary

### Box External Dimensions
```
Length = pcbLength + paddingFront + paddingBack + 2*wallThickness
       = 80 + 5 + 5 + 2*2 = 94mm

Width = pcbWidth + paddingLeft + paddingRight + 2*wallThickness
      = 60 + 5 + 5 + 2*2 = 74mm

Height = basePlaneThickness + baseWallHeight + lidWallHeight + lidPlaneThickness
       = 1.5 + 20 + 15 + 1.5 = 38mm
```

### Internal Clearances
```
PCB to base plane: standoffHeight = 5.0mm
PCB top to lid plane: baseWallHeight - standoffHeight - pcbThickness
                    = 20 - 5 - 1.6 = 13.4mm
                    (adequate for 12mm Pico stack)
```

### Cutout Centering
```
Display X: (80 - 36) / 2 = 22mm
Display Y: (60 - 17) / 2 = 21.5mm

Button Y: 60 / 2 = 30mm (centered)
```

## OpenSCAD Include Order Lessons

**Critical Discovery:** OpenSCAD evaluates code in textual order, and YAPPgenerator executes code at `include` time.

**Correct Pattern:**
1. `include <YAPPgenerator_v3.scad>` FIRST
2. Define all parameters AFTER the include
3. Calculated values can use expressions directly in arrays (e.g., `pcbWidth / 2`)

**What Doesn't Work:**
- Defining variables BEFORE the include causes them to be overwritten
- The library sets defaults for many parameters at include time

**Reference:** See `examples/ESP32-CAM-USB_v30.scad` for working pattern.

## Open Questions

### Resolved
- ✓ What is the correct OpenSCAD include order? → Include first, then parameters
- ✓ How to center cutouts? → Use calculated expressions directly in arrays
- ✓ What coordinate system to use? → `yappCoordPCB` (default) for PCB-relative features

### Remaining
- Display mounting method (screws vs. adhesive) - depends on final PCB design
- Need for ventilation slots - depends on power dissipation testing
- Label text and positioning - can be customized by user

## References

### YAPPgenerator Documentation
- [Official Documentation](https://mrwheel-docs.gitbook.io/yappgenerator_en/)
- [GitHub Repository](https://github.com/mrWheel/YAPP_Box)
- [LLM Guidelines](../../yapp-llm-guidelines.md) - Comprehensive guide for working with YAPP

### Component Specifications
- [Component Specifications Document](../various/component-specifications.md)

### Design Files
- `pico_temp_monitor_box.scad` - Main assembly
- `pico_temp_monitor_base_only.scad` - Base for printing
- `pico_temp_monitor_lid_only.scad` - Lid for printing

### Generated Files
- `pico_temp_monitor_base.stl` - Base STL (1.4MB, 4009 vertices)
- `pico_temp_monitor_lid.stl` - Lid STL (1.9MB, 5480 vertices)
- `preview_assembly.png` - Assembly preview
- `preview_base.png` - Base preview
- `preview_lid.png` - Lid preview

---

**Version:** 1.0  
**Last Updated:** 2025-11-05  
**Status:** Complete - Ready for prototyping
