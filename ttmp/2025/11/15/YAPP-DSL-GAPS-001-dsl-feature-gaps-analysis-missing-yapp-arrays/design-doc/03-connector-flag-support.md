---
Title: Connector Flag Support
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
LastUpdated: 2025-11-15T22:44:49.98527837-05:00
---


# Connector Flag Support

## Executive Summary

Connectors (screw standoffs) in YAPP rely on the same corner mirroring, coordinate, and shell-part semantics as `pcbStands`, plus additional flags for countersinks, through-lid inversion, internal fillet suppression, and PCB-specific routing. The DSL currently emits only the raw numeric tuple, forcing users to duplicate entries and preventing parity with examples like `YAPP_Demo_RealBox_v31.scad`. This document spells out the schema/builder changes needed so a single DSL entry can produce the exact SCAD structure expected by `shellConnectors()` and expose options such as self-threading or reversing screw direction.

## Problem Statement

The SCAD interface accepts nine positional numbers followed by a series of optional named flags, but the DSL exposes only the numerics, omitting everything below line 570 in the generator comments.

```552:577:YAPPgenerator_v3.scad
//    n(a) = { yappAllCorners, yappFrontLeft | <yappBackLeft> | yappFrontRight | yappBackRight }
//    n(b) = { <yappCoordPCB> | yappCoordBox | yappCoordBoxInside }
//    n(c) = { yappNoFillet } : Don't add fillets
//    n(d) = { yappCountersink }
//    n(e) = [yappPCBName, "XXX"] : Specify a PCB. Defaults to [yappPCBName, "Main"]
//    n(f) = { yappThroughLid = changes the screwhole to the lid and the socket to the base}
//    n(g) = {yappSelfThreading}
//    n(h) = { yappNoInternalFillet }
```

At runtime `shellConnectors()` mirrors coordinates per-corner and conditionally inverts the shell part when `yappThroughLid` is set, meaning parity requires the DSL to emit the same sentinel symbols.

```4062:4124:YAPPgenerator_v3.scad
allCorners = (isTrue(yappAllCorners, conn)) ? true : false;
primeOrigin = (!isTrue(yappBackLeft, conn) && !isTrue(yappFrontLeft, conn) && !isTrue(yappFrontRight, conn) && !isTrue(yappBackRight, conn) && !isTrue(yappAllCorners, conn) ) ? true : false;
shellPart = (isTrue(yappThroughLid, conn)) ? ((shellPartRaw==yappPartBase) ? yappPartLid : yappPartBase) : shellPartRaw;
if (primeOrigin || allCorners || isTrue(yappBackLeft, conn))
  connectorNew(...);
// ... repeats for other corners
```

Without DSL-level access to these flags, users cannot:

- Generate four connectors from a single definition (`yappAllCorners`)
- Flip screw direction (`yappThroughLid`)
- Call out alternate PCBs (`[yappPCBName, "Aux"]`)
- Opt into countersinks, internal fillet suppression, or self-threading

## Proposed Solution

1. **Schema extensions**
   - Add enums/booleans mirroring each `n(a)`–`n(h)` flag:
     - `corner`: `single`, `all`, `front_left`, `front_right`, `back_left`, `back_right`
     - `coordinate`: `pcb`, `box`, `box_inside`
     - `no_fillet`, `countersink`, `through_lid`, `self_threading`, `no_internal_fillet`
     - `pcb_name` string bucket (emits `[yappPCBName, value]`)

2. **Builder helpers**
   - Share the same flag-encoding helpers introduced for `pcb_stands` (moved into a reusable package or duplicated temporarily).
   - Append flags in the documented YAPP order to keep diffs predictable.
   - Ensure `through_lid` flips shell parts the same way SCAD does by emitting `yappThroughLid` plus requested corners.

3. **Validation/tests**
   - Schema tests for enum acceptance + invalid values.
   - Builder tests covering default behavior and a fully flagged entry so we verify `[..., yappAllCorners, yappCoordBoxInside, yappCountersink, ...]` ordering.

## Design Decisions

- **Typed enums over raw arrays:** Keeps YAML approachable and lets schemagen enforce valid values.
- **Helper reuse:** Mirror the shell-part/corner logic from `pcb_stands` to avoid divergence between modules that SCAD treats identically.
- **Ordered append:** Respect the SCAD documentation order for optional flags so generated `.scad` stays diff-able against canonical examples.

## Alternatives Considered

- **Single `flags` array:** Similar to SCAD but would require users to memorize symbol names and would bypass validation.
- **Auto-mirroring connectors when four definitions detected:** Too brittle and still fails to expose countersink / through-lid features.

## Implementation Plan

1. Update `pkg/yappgen/modules/connectors/schema.yaml` with the new fields + tests.
2. Regenerate schemagen output to add typed fields/defaults.
3. Introduce a `flagEncoder` in `module.go` (either via shared helper or local copy).
4. Add unit tests verifying both defaults (no extra flags) and an “everything turned on” sample.
5. Run `go test ./...` and update ticket docs/tasks.

## Open Questions

- Should we emit both `yappThroughLid` and a derived shell-part flag, or rely solely on the SCAD inversion logic?
- Do we need to expose `yappCoordPCBInside` (if added upstream) preemptively, or wait until SCAD supports it?

## References

- `YAPPgenerator_v3.scad` (`shellConnectors`, connector parameter docs)
- `pkg/yappgen/modules/connectors/schema.yaml`
- `pkg/yappgen/modules/connectors/module.go`
- `design-doc/02-pcb-stand-flag-support.md` (shared helper rationale)
