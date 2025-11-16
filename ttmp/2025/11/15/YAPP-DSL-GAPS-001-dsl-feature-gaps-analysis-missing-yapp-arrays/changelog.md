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

