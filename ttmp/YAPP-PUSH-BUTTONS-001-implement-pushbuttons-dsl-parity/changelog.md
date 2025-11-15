# Changelog

## 2025-11-15

- Initial workspace created
- Captured pushButtons SCAD parameter inventory, DSL gap analysis, and proposed YAML structure in `various/pushbuttons-dsl-analysis.md`; populated tasks/index with goals and references.
- Added general DSL documentation to the help system: `pkg/docs/tutorials/yapp-dsl-getting-started.md` (tutorial) and `pkg/docs/tutorials/yapp-dsl-reference.md` (reference), which future pushButton docs will build upon.
- Landed push button DSL support end-to-end: added schema/builder/helpers in `pkg/yappgen`, taught the model about `pushButtons` + `printSwitchExtenders`, updated `generatorcli`/`yappctl` to propagate the new flag to OpenSCAD, and covered the feature with fresh unit tests.
- Authored YAML equivalents for the SCAD references (`examples/yapp-demo-buttons.yaml`, `examples/yapp-demo-buttons2.yaml`) and generated SCAD/STL artifacts via `yappctl generate` to prove parity.
- Expanded the DSL tutorial/reference sections with in-depth push button guidance (field tables, polygon presets, resolver tips) so `yappctl help …` now explains how to use the feature.
- Introduced a registry-driven feature module system (`pkg/yappgen/features.go`) powering collection + emission for pcb_stands, connectors, snap_joins, push_buttons, and cutouts; updated `BuildModel`/`EmitSCAD` to rely on the registry and documented the architecture in `design/feature-module-registry.md`.
