# Tasks

## TODO

- [ ] Add tasks here

- [x] Research YAPPgenerator_v3.scad corner placement logic (2025-11-16)
- [ ] Get user's actual DSL YAML for comparison
- [x] Implement corner placement flags for pcb_stands (2025-11-16)
- [x] Implement shell part flags (yappBoth/LidOnly/BaseOnly) (2025-11-16)
- [x] Implement connectors flag parity (2025-11-16)
- [x] Implement boxMounts module
- [x] Implement lightTubes module
- [ ] Implement labelsPlane module
- [ ] Implement ridgeExt modules
- [x] Implement snap_joins flags (alignment/symmetric/diamond)
- [x] Add vertical dimension support (baseWallHeight, lidWallHeight, ridgeHeight)
- [x] Implement cutout polygon shapes (yappPolygon + shape presets)
- [ ] Implement cutout mask support (yappMaskDef + mask presets)
- [ ] Add cutout coordinate/origin flags (yappCoordBox, yappCenter, yappAltOrigin)
- [ ] Add push_buttons missing flags (yappAltOrigin, yappPCBName)
- [ ] Implement ValidateConstraints() for cutouts module to check enum values (shape field) - currently stubbed out, invalid values like 'polygon' only caught during Build phase
- [ ] Implement ValidateStructure() for all modules - currently stubbed out with TODO, should validate required fields and types
- [x] Fix yapp-demo-lighttubes.yaml cutouts structure - use flat list with 'face' field (not nested by face name) and change polygon to rounded_rect
- [ ] Remove ShapeFlag and SnapSideFlag in pkg/yappgen/schema.go if unused
- [ ] Remove ParamSpec schemas (pcbStands/connectors/snapJoins/cutouts) if unused
- [ ] Migrate remaining tests away from ParamSpec/buildParams helpers
- [ ] Simplify cutouts module registry.Build or document special handling
