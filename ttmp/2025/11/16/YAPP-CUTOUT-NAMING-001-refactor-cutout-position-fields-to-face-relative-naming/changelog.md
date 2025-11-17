# Changelog

## 2025-11-16

- Initial workspace created


## 2025-11-16

Verified complete implementation: schema, validation, Build(), tests, examples, STLs, and docs all working correctly with new from_face_* field names

### Related Files

- /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/cutouts/module.go — Build method using new fields
- /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/cutouts/schema.yaml — face-relative fields
- /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/cutouts/schema_gen.go — CustomValidate with helpful errors


## 2025-11-16

Fixed lighttubes example: cutout dimensions were hardcoded (20x9, 24x9) instead of computed expressions (36x17, 40x17). Added vars for shell_width and shell_height

### Related Files

- /home/manuel/code/others/YAPP_Box/examples/yapp-demo-lighttubes.yaml — use computed expressions for cutout sizes


## 2025-11-16

Added coordinate and origin flag support to cutouts module (yappCoordBox, yappCoordPCB, yappCoordBoxInside, yappCenter, yappAltOrigin). Updated lighttubes example to use coordinate: box for correct positioning

### Related Files

- /home/manuel/code/others/YAPP_Box/examples/yapp-demo-lighttubes.yaml — add coordinate:box to cutouts
- /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/cutouts/module.go — encodeFlags
- /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/cutouts/schema.yaml — add coordinate/origin fields

