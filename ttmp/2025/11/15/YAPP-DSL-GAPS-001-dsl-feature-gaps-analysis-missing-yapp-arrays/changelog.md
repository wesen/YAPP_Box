# Changelog

## 2025-11-15

- Initial workspace created


## 2025-11-15

Created comprehensive analysis of missing YAPP features and prioritization roadmap


## 2025-11-15

Implemented pcb_stands and connectors flag parity (corner/shell/countersink/self-threading).

### Related Files

- pkg/yappgen/modules/connectors/module.go — builder flag encoding
- pkg/yappgen/modules/connectors/schema.yaml — schema flags
- pkg/yappgen/modules/pcbstands/module.go — builder flag encoding
- pkg/yappgen/modules/pcbstands/schema.yaml — schema flags


## 2025-11-15

Implemented box_mounts module (schema, builder, tests) and registered module for schemagen.

### Related Files

- pkg/yappgen/modules/boxmounts/module.go — builder
- pkg/yappgen/modules/boxmounts/module_test.go — unit tests
- pkg/yappgen/modules/boxmounts/registry.go — module wiring
- pkg/yappgen/modules/boxmounts/schema.yaml — schema definition
- pkg/yappgen/modules_gen.go — schemagen registry
- ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/design-doc/04-boxmounts-module-support.md — design


## 2025-11-15

Coverage audit against YAPP_Demo_buttons_v30.scad:\n- pcb_stands: DSL now emits corner flags (yappFront*/Back*), yappBaseOnly/yappBoth, and supports treatment/coord/no_fillet/pcb_name/self_threading (not used by v30). Verified SCAD shows flags in pcbStands entries.\n- cutouts: DSL supports rectangle/circle/rounded_rect on all faces; v30 also uses yappPolygon + mask (maskHexCircles) on base which is NOT yet implemented (polygon+mask flags and objects).\n- action items: add polygon/mask support to cutouts; add snap_joins flags (yappCenter/yappSymmetric/yappRectangle) to mirror v30 demo patterns.

### Related Files

- examples/YAPP_Demo_buttons_v30.scad — reference features used
- examples/yapp-demo-buttons-v30.yaml — DSL parity example
- pkg/resolver/resolver.go — treat enum fields (corner/corners/etc.) as literals
- pkg/yappgen/features.go — wire modules to emit flags
- pkg/yappgen/modules/cutouts/schema.yaml — current shapes (polygon/mask missing)
- pkg/yappgen/modules/pcbstands/module.go — emits shell/corners/treatment/coord flags
- pkg/yappgen/modules/pcbstands/schema.yaml — adds corners[]


## 2025-11-16

DSL parity improvements for v30 demo: added vars/expressions to YAML, derived shellWidth/shellHeight equivalents, added cutoutsBase/front/back, pcb_stands corners/shell flags now emitted via module builders; regenerated SCAD/STLs.

### Related Files

- examples/yapp-demo-buttons-v30.yaml — parameterized with vars and expressions
- pkg/yappgen/features.go — now uses module builders for stands/boxMounts
- pkg/yappgen/modules/pcbstands/module.go — corners[] support and flags
- tmp/yapp-demo-buttons-v30.scad — generated SCAD reflecting changes


## 2025-11-16

Implemented snap_joins flag support: alignment (yappOrigin/yappCenter), symmetric (yappSymmetric), diamond (yappRectangle). Wired module builder to features.go. Generated and tested v30 example with snap joins.

### Related Files

- examples/yapp-demo-buttons-v30.yaml — includes snap_joins with flags
- pkg/yappgen/features.go — wired snapjoins.Build
- pkg/yappgen/modules/snapjoins/module.go — flag encoding
- pkg/yappgen/modules/snapjoins/module_test.go — flag tests
- pkg/yappgen/modules/snapjoins/schema.yaml — added alignment/symmetric/diamond flags


## 2025-11-16

Added vertical dimension support: baseWallHeight, lidWallHeight, ridgeHeight now configurable via enclosure.base.wall_height, enclosure.lid.wall_height, enclosure.ridge.height. Fixed ridge mismatch (was defaulting to 5mm, now respects YAML value). Generated v30 STLs now match original dimensions.

### Related Files

- examples/yapp-demo-buttons-v30.yaml — includes wall heights and ridge
- pkg/yappgen/emit.go — emit wall heights and ridge
- pkg/yappgen/model.go — added wall height fields


## 2025-11-16

Implemented lightTubes module - LED light pipes with circle/rectangle shapes, lens thickness, coordinate flags, and PCB name support

### Related Files

- examples/yapp-demo-lighttubes.yaml — test example matching YAPP_Demo_lightTubes_v30.scad
- pkg/yappgen/modules/lighttubes/module.go — builder implementation
- pkg/yappgen/modules/lighttubes/registry.go — module registration
- pkg/yappgen/modules/lighttubes/schema.yaml — schema definition


## 2025-11-16

Added test case file with resolver errors (yapp-demo-lighttubes-with-errors.yaml) for testing resolver error message improvements

### Related Files

- examples/yapp-demo-lighttubes-with-errors.yaml — test case demonstrating resolver error when snap_joins uses 'sides' array instead of 'side' string


## 2025-11-16

Added prominent links to module system implementation guide in index.md - marked as CRITICAL reference for all feature implementations


## 2025-11-16

Created comprehensive implementation diary for lightTubes module - documents step-by-step process, challenges, learnings, and recommendations for future implementations


## 2025-11-16

Updated implementation diary with STL rendering section - documented YAML structure fixes, STL generation process, and file comparison results


## 2025-11-16

Identified validation gap: cutouts schema defines enum for shape field but ValidateStructure/ValidateConstraints are stubbed out (TODO) - invalid enum values like 'polygon' only caught during Build phase, not resolve phase

### Related Files

- pkg/yappgen/modules/cutouts/registry.go — ValidateStructure and ValidateConstraints return nil without checking enum values


## 2025-11-16

Added tasks for validation gap: ValidateConstraints/ValidateStructure implementation needed for enum checking. Fixed yapp-demo-lighttubes.yaml cutouts structure and shape value

