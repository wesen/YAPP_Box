---
Title: BoxMounts Module Support
Ticket: YAPP-DSL-GAPS-001
Status: active
Topics:
    - yapp
    - dsl
    - analysis
    - features
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2025-11-15T23:03:27.802709516-05:00
---


# BoxMounts Module Support

## Executive Summary

External mounting tabs (`boxMounts`) are the user-visible “little plates” sticking out of `YAPP_Demo_buttons2_v31.scad`, but the DSL still has no representation for them, blocking parity with wall-mountable examples and forcing hand-edits in SCAD. This doc captures the SCAD contract and the schema/builder work needed to expose the feature through the YAML DSL so a single configuration can describe screw diameter, slot width, wall selection, and per-shell output.

## Problem Statement

SCAD expects each mount entry to carry positional vectors, screw dimensions, fillet overrides, and a set of flags describing which wall(s), origin, and shell part to affect.

```686:707:YAPPgenerator_v3.scad
//    p(0) = pos : position along the wall : [pos,offset]
//    p(1) = screwDiameter
//    p(2) = width of opening in addition to screw diameter
//    p(3) = height
//    n(a) = { yappLeft | yappRight | yappFront | yappBack }
//    p(4) = filletRadius
//    n(b) = { yappNoFillet }
//    n(c) = { <yappBase>, yappLid }
//    n(d) = { yappCenter }
//    n(e) = { <yappGlobalOrigin>, yappAltOrigin }
boxMounts = [];
```

During generation the mounts are subtracted from both base and lid shells whenever present and mirrored based on the chosen faces, so missing DSL support means users cannot specify even a single tab outside SCAD. `YAPP_Demo_buttons2_v31.scad` also proves the mounts are part of pre-cut geometry across both shells, so omission creates glaring discrepancies.`1756:1807:YAPPgenerator_v3.scad`

## Proposed Solution

1. **Schema**
   - Fields: `position` (object with `pos` + optional `offset`), `screw_d`, `slot_width`, `height`, `fillet_radius`.
   - Enums/flags:
     - `faces`: array enum `[left,right,front,back]` (at least one).
     - `shell_part`: `base`/`lid`/`both` (maps to `{yappBase,yappLid}` flags).
     - `alignment`: `edge` (default) or `center` → `yappCenter`.
     - `origin`: `global`/`alt` → `yappGlobalOrigin` / `yappAltOrigin`.
     - `no_fillet`: boolean.
   - Coordinate system defaults to `box` per SCAD; no override needed yet.

2. **Builder**
   - Emit numeric tuple `[pos, screw_d, slot_width, height, fillet_radius]` followed by wall flags and optional selectors in the exact documented order.
   - Accept either scalar `pos` or `{pos, offset}` object to match SCAD’s `vector` support.
   - Add validation to ensure faces array non-empty and slot width non-negative (0 means circular).

3. **Testing**
   - Schema fixtures verifying faces requirement and invalid enums.
   - Builder unit tests for edge-aligned mount on left/right and center-aligned mount on back/lid.

## Design Decisions

- **Structured position field**: we encode `[pos, offset]` as a YAML object with `pos` & `offset` numbers; a plain number is accepted for backwards compatibility and converted during build.
- **Face array**: YAML will expose `faces: [left, right]` instead of multiple booleans; builder maps each entry to `yappLeft` etc.
- **Shell selection**: explicit `shell_part` flag prevents silent base-only behavior; default `both` matches SCAD default subtraction of both shells when mounts exist.

## Alternatives Considered

- **Separate modules per face**: rejected because SCAD expects a single combined array and runs face selection internally; splitting would complicate mirroring.
- **Implicit origin detection**: auto-choosing `yappAltOrigin` based on face would hide an explicit SCAD flag and make it hard to reproduce legacy designs.

## Implementation Plan

1. Update `schema.yaml` with the fields/enums above plus fixtures (`minimal_mount`, `invalid_face`, `centered_back_mount`).
2. Rerun `schemagen discover` to pick up new structs + defaults.
3. Implement `Build` with helpers: `parsePosition`, `faceFlags`, `shellFlag`, `originFlag`.
4. Write Go tests:
   - default left mount only (no extra flags).
   - multi-face, center-aligned, lid-only example verifying emitted `yappCenter`, `yappLid`, multiple face flags.
5. Document usage in `examples/` and add to priority roadmap/checklist.

## Open Questions

- Do we also need a `length` parameter to represent extra reinforcement (not mentioned in v3 docs)?
- Should mounts support per-face offsets differently (SCAD uses `[pos, offset]` but no axis-specific offset), confirm with user once DSL parity achieved.

## References

- `YAPPgenerator_v3.scad` (`boxMounts` parameter block, `oneMount` loop, `minkowskiBox` pre-cut integration)
- `examples/YAPP_Demo_buttons2_v31.scad` (real-world use case)
