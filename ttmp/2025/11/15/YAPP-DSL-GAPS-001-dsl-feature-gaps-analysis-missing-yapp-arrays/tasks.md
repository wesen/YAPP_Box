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
- [x] Implement cutout mask support (yappMaskDef + mask presets)
- [x] Add cutout coordinate/origin flags (yappCoordBox, yappCenter, yappAltOrigin)
- [x] Add push_buttons missing flags (yappAltOrigin, yappPCBName)
- [ ] Implement ValidateConstraints() for cutouts module to check enum values (shape field) - currently stubbed out, invalid values like 'polygon' only caught during Build phase
- [ ] Implement ValidateStructure() for all modules - currently stubbed out with TODO, should validate required fields and types
- [x] Fix yapp-demo-lighttubes.yaml cutouts structure - use flat list with 'face' field (not nested by face name) and change polygon to rounded_rect
- [ ] Remove ShapeFlag and SnapSideFlag in pkg/yappgen/schema.go if unused
- [x] Remove ParamSpec schemas (pcbStands/connectors/snapJoins/cutouts) if unused
- [ ] Migrate remaining tests away from ParamSpec/buildParams helpers
- [ ] Simplify cutouts module registry.Build or document special handling
- [x] Create polygon cutouts comparison (DSL+SCAD) and generate STLs to /tmp/yapp_compare/cutouts_polygons
- [x] Create lighttubes comparison (DSL+SCAD) and generate STLs to /tmp/yapp_compare/lighttubes
- [x] Face-aware cutout mapping: clearer pos axes for side faces (front/back=posy,posz; left/right=posx,posz); update docs and examples to match legacy semantics
- [x] Adjust lighttubes example cutouts to match legacy heights (posz); verify with STL compare
- [x] yappctl generate: auto-copy embedded YAPPgenerator_v3.scad; remove default render timeout
- [x] Add --stl-all (single STL) and render push-button extenders only in lid when split
- [x] Update yapp-demo-buttons-v30.yaml: base polygon+mask (hexagon + hex_circles), centered
- [ ] Document cutouts.mask usage and presets in DSL reference
- [x] Compare DSL vs legacy STLs for buttons v30 and record findings
