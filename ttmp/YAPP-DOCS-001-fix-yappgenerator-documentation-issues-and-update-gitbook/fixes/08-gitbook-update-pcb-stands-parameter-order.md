# GitBook Update: PCB Stands — Parameter Order to v3.3.8

Status: proposed
Owner: manuel
LastUpdated: 2025-11-05

## Summary

The PCB Stands page lists parameters in a v3.0-era order. Update the documentation to reflect the v3.3.8 parameter signature and defaults to avoid confusion when copying examples.

## What to change

- Present required parameters first, then optionals grouped by role (geometry, placement, flags).
- Update named flags to include current options and guidance.

## Suggested Structure (v3.3.8)

Required:
- `p(0)` posx
- `p(1)` posy

Optional:
- `p(2)` heightToBottomOfPCB (Default = `standoffHeight`)
- `p(3)` pcbGap (Default = `-1`; resolves to `pcbThickness` for `yappCoordPCB`, else `0`)
- `p(4)` standoffDiameter (Default = `standoffDiameter`)
- `p(5)` standoffPinDiameter (Default = `standoffPinDiameter`)
- `p(6)` standoffHoleSlack (Default = `standoffHoleSlack`)
- `p(7)` filletRadius (`0` = Auto)

Named flags (order-insensitive):
- `n(a)` `{ <yappBoth> | yappLidOnly | yappBaseOnly }`
- `n(b)` `{ <yappPin>, yappHole, yappTopPin }`
- `n(c)` corners `{ yappAllCorners | yappFrontLeft | yappFrontRight | yappBackLeft | yappBackRight }`
- `n(d)` coords `{ <yappCoordPCB> | yappCoordBox | yappCoordBoxInside }`
- `n(e)` `{ yappNoFillet }`
- (optional) `[yappPCBName, "Main"]`

## Before (example excerpt)

```scad
// Older docs showed a different optional ordering and omitted corner/capability flags
```

## After (v3.3.8-aligned example)

```scad
pcbStands = [
  // p(0) posx, p(1) posy, p(2) heightToBottomOfPCB, p(3) pcbGap,
  // p(4) standoffDiameter, p(5) standoffPinDiameter, p(6) standoffHoleSlack, p(7) filletRadius
  [20, 15, -1, -1, 7, 2.4, 0.4, 0, yappCoordPCB, yappBoth]
];
```

Notes:
- Use `yappCoordPCB` unless you explicitly need Box/BoxInside.
- For circular placements near corners, use `corners` flags to target specific corners.

## Validation

- Cross-check against the v3.3.8 implementation and example templates in the repo.
- Render a minimal case to confirm no warnings and expected geometry.

Links:
- GitBook page: `https://mrwheel-docs.gitbook.io/yappgenerator_en/pcb-stands`
- Reference doc: `reference/pcb-stands-parameter-order-verification.md`


