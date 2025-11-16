---
Title: PCB Stand Flag Support
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
LastUpdated: 2025-11-15T22:38:08.760963346-05:00
---


# PCB Stand Flag Support

## Executive Summary

PCB standoffs in OpenSCAD gain tremendous leverage from the YAPP corner-selection and shell-part flags, yet our DSL currently emits only the bare positional tuple. As a result, users must define four nearly identical entries anytime they need a stand in each corner and cannot mirror lid/base-only pushdowns at all. This document defines the schema and builder changes required to expose the SCAD-level flags so a single DSL entry can replicate the canonical `pcbHolders()` behavior and stay in parity with `pcbPushdowns()`/`shellConnectors()`.

## Problem Statement

The legacy SCAD logic mirrors stand coordinates into other corners whenever `yappAllCorners` or one of the directional flags is set, and suppresses lid/base geometry by checking `yappLidOnly`/`yappBaseOnly`. Because our DSL lacks any representation of those `n(a)`–`n(g)` selectors, DSL authors must manually duplicate YAML entries (error-prone, verbose) and still cannot toggle lid-only pushdowns or self-threading holes. This breaks parity with examples like `YAPP_Demo_buttons2_v31.scad` and directly caused the “extra pillars” discrepancy identified in the ticket.

```2063:2083:YAPPgenerator_v3.scad
allCorners = (isTrue(yappAllCorners, stand)) ? true : false;
primeOrigin = (!isTrue(yappBackLeft, stand) && !isTrue(yappFrontLeft, stand) && !isTrue(yappFrontRight, stand) && !isTrue(yappBackRight, stand) && !isTrue(yappAllCorners, stand) ) ? true : false;
if (!isTrue(yappLidOnly, stand))
{
  if (primeOrigin || allCorners || isTrue(yappBackLeft, stand))
    translate([offsetX+connX, offsetY + connY, basePlaneThickness])
      pcbStandoff(...);
  if (allCorners || isTrue(yappFrontLeft, stand))
    translate([offsetX + lengthX - connX, offsetY + connY, basePlaneThickness])
      pcbStandoff(...);
  // remaining corners omitted
}
```

## Proposed Solution

1. **Schema extensions**
   - Add a `corner` enum field with values `single` (default), `all`, `front_left`, `front_right`, `back_left`, `back_right`. This maps directly to the `yappAllCorners`/`yappFrontLeft`… flags.
   - Add a `shell_part` enum with `both` (default), `lid_only`, `base_only` that controls the mutual-exclusion checks in `pcbHolders()`/`pcbPushdowns()`.
   - Add `treatment` enum (`pin`, `hole`, `top_pin`) mirroring `yappPin`, `yappHole`, `yappTopPin`.
   - Add `coord_system` enum to select `yappCoordPCB` (default), `yappCoordBox`, or `yappCoordBoxInside`.
   - Add optional booleans/strings for the remaining named flags: `no_fillet`, `self_threading`, `pcb_name`.

2. **Builder updates**
   - Extend `PcbStandsItem` (generated) with the new fields and defaulting logic.
   - Teach `Build` to append a `map[string]any` of named flags after the positional slice (e.g., `[]any{..., flagSet}`) using a helper that converts enum selections into the appropriate `scad.Symbol`.
   - Reuse the same helper for connectors in a follow-up change so both modules stay aligned.

3. **Validation**
   - Enforce `pcb_name` only when present; default to `Main`.
   - Ensure `self_threading` cannot be true unless `shell_part` includes the lid (mirrors SCAD expectation of paired stand/hole).

## Design Decisions

- Model each SCAD flag as a strongly-typed enum/boolean instead of free-form symbol arrays. This keeps YAML readable while still allowing schemagen to emit meaningful Go types.
- Default corner mode to `single` to preserve today’s behavior for existing YAML.
- Keep the “named parameters” order identical to the SCAD documentation (`n(a)` through `n(g)`) so diffing generated SCAD stays predictable.

## Alternatives Considered

- **Literal flag arrays**: Allow users to specify `flags: ["yappAllCorners","yappLidOnly"]`. Rejected because it leaks SCAD internals into the DSL and offers no validation.
- **Automatic inference**: Attempt to detect repeated coordinates and convert them into `yappAllCorners`. Rejected because it is brittle and would still leave shell/treatment flags unaddressed.

## Implementation Plan

1. Update `schema.yaml` and regenerate schemagen output to include the new fields.
2. Add `ApplyDefaults` logic + custom validation for enum values requiring coordination (e.g., `self_threading` guard).
3. Extend `module.go` with a `flagEncoder` that emits the appropriate `scad.Symbol` entries.
4. Add fixture tests covering `corner=all`, `shell_part=lid_only`, `treatment=hole`, etc., ensuring the builder output matches expected SCAD arrays.
5. Mirror the same schema/builder primitives into `connectors` (separate PR but same helper).

## Open Questions

- Should `pcb_name` be required once multi-PCB support lands, or do we keep defaulting to `Main` indefinitely?
- Do we need to expose `yappNoFillet` separately for lid vs base, or is the single boolean sufficient?

## References

- `pkg/yappgen/modules/pcbstands/schema.yaml`
- `pkg/yappgen/modules/pcbstands/module.go`
- `YAPPgenerator_v3.scad` (`pcbHolders`, `pcbPushdowns`, `shellConnectors`)
