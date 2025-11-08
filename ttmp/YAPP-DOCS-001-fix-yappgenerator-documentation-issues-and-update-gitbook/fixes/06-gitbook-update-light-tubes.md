# GitBook Update: Light Tubes — yappCircle should use yappCenter

Status: proposed
Owner: manuel
LastUpdated: 2025-11-05

## Summary

On the Light Tubes page, examples that use `yappCircle` should enable center-based positioning with `yappCenter`. This avoids manual radius offsets and aligns with v3.3.8 placement flags.

## What to change

- In example arrays where `p(6) = yappCircle`, add `yappCenter` in the named flag positions.
- Ensure coordinate system is explicit via `n(a)` (e.g., `yappCoordPCB`).

## Before

```scad
lightTubes = [
  // p(0) posx, p(1) posy, p(2) tubeLength, p(3) tubeWidth,
  // p(4) tubeWall, p(5) gapAbovePcb, p(6) tubeType, p(7) lensThickness
  [15, 20, 1.5, 5, 1, 2, yappCircle, 0.5],
  [15, 30, 1.5, 5, 1, 2, yappRectangle]
];
```

## After (v3.3.8)

```scad
lightTubes = [
  // ... same positional params ...
  [15, 20, 1.5, 5, 1, 2, yappCircle, 0.5, yappCoordPCB, yappCenter],
  [15, 30, 1.5, 5, 1, 2, yappRectangle, 0.5, yappCoordPCB]
];
```

Notes:
- Keep any existing flags like `yappNoFillet` and coordinate choices intact.
- If the example omits lens thickness (`p(7)`), append `yappCoordPCB, yappCenter` directly after the existing entries.

## Rationale

- `yappCenter` reduces errors from manual radius offsets for circular features.
- Matches the placement flag pattern used by other features (cutouts, pushButtons) in v3.3.8.

## Validation

- Local minimal tests compile and render with `yappCenter` enabled.
- Examples pass OpenSCAD preview export in CI when flags are present.


