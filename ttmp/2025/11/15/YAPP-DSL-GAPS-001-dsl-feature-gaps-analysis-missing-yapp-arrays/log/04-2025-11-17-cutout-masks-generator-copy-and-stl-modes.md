---
Title: 2025-11-17 Cutout masks, generator copy, and STL modes
Ticket: YAPP-DSL-GAPS-001
Status: active
Topics:
    - yapp
    - dsl
    - analysis
    - features
DocType: log
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Implemented cutout masks; added yappctl generator copy and STL modes; matched buttons v30 base mesh
LastUpdated: 2025-11-16T20:48:03.102814945-05:00
---


# 2025-11-17 Cutout masks, generator copy, and STL modes

## What changed
- Cutout masks implemented in cutouts module (yappMaskDef presets + optional offsets/rotation).
- Resolver updated to treat `features.*.mask.preset` as a literal string (no expression eval).
- yappctl generate now:
  - Auto-copies embedded `YAPPgenerator_v3.scad` to the output dir when rendering any STL (or with `--copy-generator`).
  - Removes the default render timeout (no timeout unless `--render-timeout` is provided).
  - Adds `--stl-all` to render a single STL (base+lid, includes push-button extenders if present).
  - Ensures split STLs render push-button extenders only in the lid STL (not in base).
- Updated `examples/yapp-demo-buttons-v30.yaml` to use a base polygon+mask (hexagon + hex_circles), angle 30, centered, to better match legacy.

## Validation
- Masks demo:
  - SCAD/Render OK; example at `/home/manuel/code/others/YAPP_Box/examples/yapp-demo-masks.yaml`
  - STLs: `/tmp/yapp_compare/masks2/base.stl`, `/tmp/yapp_compare/masks2/lid.stl`
- Buttons v30 comparison:
  - DSL vs legacy STLs written to `/tmp/yapp_compare/buttons_v30/`
  - DSL: `dsl-base.stl`, `dsl-lid.stl`
  - Legacy: `legacy-base.stl`, `legacy-lid.stl`
  - Observations: legacy mesh isn’t perfectly centered and hex-circle rows clip near edges; our mask + offsets reproduces this visually close enough.
- Split vs single STL output:
  - Split: `/tmp/yapp_compare/buttons_split/base.stl`, `/tmp/yapp_compare/buttons_split/lid.stl`
  - Single: `/tmp/yapp_compare/buttons_all/all.stl`

## Files touched
- Cutouts:
  - `pkg/yappgen/modules/cutouts/schema.yaml` (added `mask` object)
  - `pkg/yappgen/modules/cutouts/module.go` (mask emission + helpers)
  - `pkg/resolver/resolver.go` (string-field rule for `.mask.preset`)
- CLI / rendering:
  - `pkg/cli/generatorcli/generator.go` (copy embedded generator, STL modes; extenders only in lid for split)
  - `cmd/yappctl/generate_command.go` (`--stl-all`, generator copy auto-enable, no default timeout)
- Examples:
  - `examples/yapp-demo-masks.yaml` (new)
  - `examples/yapp-demo-buttons-v30.yaml` (base polygon+mask)

## Next
- Document `cutouts.mask` usage and presets in DSL reference (`pkg/docs/tutorials/yapp-dsl-reference.md`).
- Consider optional helpers to fine-tune mask offsets per face to better match legacy clipping patterns when needed.
