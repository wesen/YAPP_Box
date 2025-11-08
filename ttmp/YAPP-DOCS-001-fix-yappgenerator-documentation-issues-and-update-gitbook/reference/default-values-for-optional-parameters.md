# Default Values for Optional Parameters

Status: draft
Owner: manuel
LastUpdated: 2025-11-05

## Overview

This reference aggregates default values for commonly used optional parameters in YAPPgenerator features. Extracted from current docs (v3.3.8) and example snippets.

## Light Tubes

- `p(7)` lensThickness: Default = `0` (open)
- `p(8)` heightToTopOfPCB: Default = `standoffHeight + pcbThickness`
- `p(9)` filletRadius: Default = `0` (Auto)
- `n(a)` coordinate system: Default = `yappCoordPCB`
- `n(b)` flags: `yappNoFillet` (if provided)

## Push Buttons

- `p(9)` heightToTopOfPCB: Default = `standoffHeight + pcbThickness`
- `p(10)` shape: Default = `yappRectangle`
- `p(11)` angle: Default = `0`
- `p(12)` filletRadius: Default = `0` (Auto)
- `p(13)` buttonWall: Default = `2.0`
- `p(14)` buttonPlateThickness: Default = `2.5`
- `p(15)` buttonSlack: Default = `0.25`
- `p(16)` snapSlack: Default = `0.10`
- `n(d)` `[yappPCBName, "Main"]` by default

## Cutouts

- `p(6)` depth: Default = `0` (Auto → plane thickness)
- `p(7)` angle: Default = `0`
- `n(d)` coordinate system: `{ <yappCoordPCB> | yappCoordBox | yappCoordBoxInside }` (default in docs is PCB for many examples; set explicitly)
- `n(e)` origin/center: `{ <yappOrigin>, yappCenter }` (default is origin)

## Box Mounts

- `p(4)` filletRadius: Default = `0` (Auto)
- `n(a)` faces: `{ yappLeft | yappRight | yappFront | yappBack }`
- `n(b)` flags: may include `yappNoFillet`
- `n(c)` base/lid: `{ <yappBase>, yappLid }` (default base)
- `n(d)` center flag: optional `{ yappCenter }`

## Connectors

- `p(7)` insertDepth: Default = entire connector
- `p(8)` pcbGap: Default depends on coord (`yappCoordPCB` → `pcbThickness`, else `0`)
- `p(9)` filletRadius: Default = `0` (Auto)
- `n(d)` countersink: optional `{ yappCountersink }`
- `n(g)` self threading: optional `{ yappSelfThreading }`
- `n(d/e/...)` coordinate/PCB name: `{ <yappCoordPCB> | yappCoordBox | yappCoordBoxInside }`, `[yappPCBName, "Main"]`

## Labels

- Direction: Default = `yappTextLeftToRight`
- Horizontal alignment: Default = `yappTextHAlignLeft`
- Vertical alignment: Default = `yappTextVAlignBottom`
- Spacing: Default = `1.0`

## Ridge Extensions

- Default coordinate system: `yappCoordBox`
- `n(a)` origin/center: `{ <yappOrigin>, yappCenter }` (default origin)
- `n(b)` left/global origin: `{ <yappGlobalOrigin>, yappAltOrigin }` (default global)
- Notes: `height` is measured relative to the ridge; negative values move into lid

## PCB Stands

- `p(2)` heightToBottomOfPCB: Default = `standoffHeight`
- `p(3)` pcbGap: Default = `-1` (resolves to `pcbThickness` for `yappCoordPCB`, else `0`)
- `p(4)` standoffDiameter: Default = `standoffDiameter`
- `p(5)` standoffPinDiameter: Default = `standoffPinDiameter`
- `p(6)` standoffHoleSlack: Default = `standoffHoleSlack`
- `p(7)` filletRadius: `0` (Auto)
- Named flags: see `PCB Stands: Parameter Order Verification`

## Notes

- For each feature, keep defaults listed near the parameter description.
- When examples rely on defaults, add a comment to make the assumption visible.

## To Do (Data Backfill)

- Cross-check against upstream OpenSCAD modules for authoritative defaults.
- Extend with other features (imagesPlane, vents, labels) as we verify.

---
Title: Default Values for Optional Parameters
Ticket: YAPP-DOCS-001
Status: active
Topics:
    - yappgenerator
    - documentation
DocType: reference
Intent: long-term
Owners:
    - manuel
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2025-11-05T17:17:37.102264573-05:00
---


# Default Values for Optional Parameters

## Goal

<!-- What is the purpose of this reference document? -->

## Context

<!-- Provide background context needed to use this reference -->

## Quick Reference

<!-- Provide copy/paste-ready content, API contracts, or quick-look tables -->

## Usage Examples

<!-- Show how to use this reference in practice -->

## Related

<!-- Link to related documents or resources -->
