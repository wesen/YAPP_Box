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

