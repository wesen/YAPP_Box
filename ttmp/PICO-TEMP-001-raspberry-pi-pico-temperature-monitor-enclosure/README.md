# Raspberry Pi Pico Temperature Monitor Enclosure

**Project**: PICO-TEMP-001  
**Version**: 1.0  
**Date**: 2025-11-05  
**YAPPgenerator Version**: v3.3.8 (2025-10-24)

## Overview

This is a parametric 3D-printable enclosure designed using YAPPgenerator v3 for a Raspberry Pi Pico-based temperature monitoring device. The enclosure accommodates:

- Raspberry Pi Pico mounted on a custom PCB (via socket)
- 1.3" OLED display module (40x30mm PCB, 36x17mm visible area)
- Panel-mount push button (28mm diameter, 45mm threaded barrel)
- Temperature probe/sensor (8mm probe diameter, 6mm cable)

## Box Dimensions

### External Dimensions
- **Length**: 94mm (X-axis, front to back)
- **Width**: 74mm (Y-axis, side to side)
- **Height**: 38mm (Z-axis, total assembly)

### Internal Dimensions
- **Inside Length**: 90mm
- **Inside Width**: 70mm
- **Inside Height**: 35mm

### PCB Specifications
- **PCB Length**: 80mm
- **PCB Width**: 60mm
- **PCB Thickness**: 1.6mm
- **Standoff Height**: 5mm (allows clearance for Pico stack ~12mm total)

### Wall Specifications
- **Wall Thickness**: 2.0mm
- **Base Plane Thickness**: 1.5mm
- **Lid Plane Thickness**: 1.5mm
- **Base Wall Height**: 20mm
- **Lid Wall Height**: 15mm
- **Ridge Height**: 5.0mm (for lid/base overlap)
- **Ridge Slack**: 0.3mm (gap for fit tolerance)
- **Round Radius**: 3.0mm (rounded corners)

## Features

### Mounting
- **PCB Stands**: 4 corner standoffs with pins (7mm diameter, 2.4mm pin diameter)
- **Snap Joins**: 2 pairs on left and right sides for tool-free assembly

### Cutouts

#### Lid Cutouts
- **Display Opening**: 36mm x 17mm rectangular cutout, centered on PCB
  - Position: X = 22mm from back, Y = 21.5mm from left
  - Allows OLED display visibility

#### Base Cutouts
- **Cable Grommet**: 10mm diameter circular cutout for temperature probe
  - Position: Near back-right corner (70mm from back, 50mm from left)
  - Accommodates 8mm probe with 6mm cable

#### Front Panel Cutouts
- **Push Button**: 28mm diameter circular cutout
  - Position: Centered vertically (30mm from bottom), 15mm from base
  - Accommodates 28mm button cap with 45mm threaded barrel

### Labels
- **Lid**: "Pico Temp Monitor" (Liberation Sans Bold, 5pt)
- **Base**: "v1.0" (Liberation Sans, 3pt)

## Files Included

### Design Files
- `pico_temp_monitor_box.scad` - Complete assembly (for preview)
- `pico_temp_monitor_base_only.scad` - Base shell only (for printing)
- `pico_temp_monitor_lid_only.scad` - Lid shell only (for printing)

### STL Files (Ready to Print)
- `pico_temp_monitor_base.stl` - Base shell (1.4MB, 4009 vertices)
- `pico_temp_monitor_lid.stl` - Lid shell (1.9MB, 5480 vertices)

### Preview Images
- `preview_assembly.png` - Complete assembly view
- `preview_base.png` - Base shell view
- `preview_lid.png` - Lid shell view

### Documentation
- `README.md` - This file
- `component-specifications.md` - Detailed component specifications
- `design_rationale.md` - Design decisions and calculations

## Rendering Performance

| Part | Render Time | Vertices | File Size |
|------|-------------|----------|-----------|
| Assembly | 1:12 min | 9,553 | N/A |
| Base | 0:21 sec | 4,009 | 1.4MB |
| Lid | 0:39 sec | 5,480 | 1.9MB |

## Design Rationale

### PCB Sizing
The custom PCB is sized at 80x60mm to accommodate:
- Raspberry Pi Pico (51x21mm) with socket headers
- OLED display module (40x30mm)
- Additional circuitry and connections
- Adequate spacing for wiring and component clearance

### Height Calculations
Total internal height budget:
- Base plane: 1.5mm
- Standoff height: 5.0mm
- PCB thickness: 1.6mm
- Pico stack height: ~12mm (Pico + headers + socket)
- Clearance to lid: ~15mm
- **Total**: 35mm internal height

### Cutout Positioning

#### Display Cutout
- Display visible area: 36mm x 17mm
- Centered on PCB: `X = (80-36)/2 = 22mm`, `Y = (60-17)/2 = 21.5mm`
- Allows display to be mounted from inside with mounting holes

#### Button Cutout
- 28mm diameter for button cap
- Positioned on front panel for easy access
- Centered vertically for ergonomic operation

#### Cable Grommet
- 10mm diameter (oversized for 8mm probe)
- Positioned at back-right corner for cable management
- Allows probe cable to exit cleanly

### Assembly Method
- **Snap joins** provide tool-free assembly
- **Ridge overlap** ensures proper alignment and dust resistance
- **Ridge slack** (0.3mm) compensates for 3D printing tolerances

## Printing Recommendations

### Print Settings
- **Layer Height**: 0.2mm (as designed)
- **Wall Thickness**: 2-3 perimeters (for 2mm walls)
- **Infill**: 20% (adequate for structural integrity)
- **Supports**: None required (designed for printability)
- **Orientation**: 
  - Base: Print upside down (ridge facing up)
  - Lid: Print upside down (flat surface on build plate)

### Material
- **PLA**: Recommended for ease of printing
- **PETG**: For better durability and temperature resistance
- **ABS**: For maximum strength (requires heated chamber)

### Post-Processing
- Remove any stringing or artifacts
- Test fit before final assembly
- Light sanding of snap join areas if too tight

## Assembly Instructions

1. **Prepare Components**
   - Solder Raspberry Pi Pico to custom PCB (via socket headers)
   - Connect OLED display to PCB
   - Attach temperature probe cable

2. **Mount Display**
   - Position OLED module over lid cutout from inside
   - Secure with M2 screws through mounting holes (if designed into PCB)
   - Alternative: Use hot glue or double-sided tape

3. **Install Button**
   - Insert button through front panel cutout
   - Secure with threaded nut from inside
   - Connect button wires to PCB

4. **Mount PCB**
   - Align PCB over standoff pins in base
   - Press down gently until PCB seats on standoffs
   - Pins should protrude through PCB mounting holes

5. **Route Cable**
   - Thread temperature probe cable through base grommet
   - Ensure adequate slack inside for movement

6. **Close Enclosure**
   - Align lid with base
   - Press down until snap joins engage
   - Check that all edges are flush

## Modifications

### Changing Dimensions
To modify the box size, edit these parameters in the `.scad` files:

```openscad
pcbLength           = 80;   // X-axis dimension
pcbWidth            = 60;   // Y-axis dimension
baseWallHeight      = 20;   // Base height
lidWallHeight       = 15;   // Lid height
```

### Adding Features
The YAPPgenerator supports many additional features:
- Light tubes (for LED indicators)
- Additional cutouts (for connectors, switches)
- Hooks (for cable management)
- Connectors (for screw-based assembly)

See `YAPPgenerator_v3.scad` and `YAPP_Template_v3.scad` for full documentation.

## Design Verification

### Checklist
- [x] All components fit within internal dimensions
- [x] Cutouts sized correctly for components
- [x] Adequate clearance for Pico stack height
- [x] Cable routing accommodated
- [x] Snap joins positioned for easy assembly
- [x] Labels readable and positioned correctly
- [x] STL files generated without errors
- [x] Render times acceptable (<2 minutes per part)

### Known Limitations
- Display mounting holes not explicitly modeled (assumes PCB has mounting holes)
- No ventilation slots (add if needed for heat dissipation)
- No mounting feet or wall-mount provisions (can be added)

## License

This design uses YAPPgenerator v3 by Willem Aandewiel and contributors.  
YAPPgenerator is licensed under the Creative Commons - Attribution - Share Alike license.

## References

- [YAPPgenerator Documentation](https://mrwheel-docs.gitbook.io/yappgenerator_en/)
- [YAPPgenerator GitHub Repository](https://github.com/mrWheel/YAPP_Box)
- Component specifications: See `component-specifications.md`
- Design decisions: See `design_rationale.md`

## Changelog

### v1.0 (2025-11-05)
- Initial design
- Verified rendering with YAPPgenerator v3.3.8
- Generated STL files for base and lid
- Created documentation

---

**Questions or Issues?**  
Refer to the YAPPgenerator documentation or check the component specifications document for detailed measurements.
