# GitBook Update: Push Buttons — yappCircle should use yappCenter

Status: proposed
Owner: manuel
LastUpdated: 2025-11-05

## Summary

On the Push Buttons page, when the cap shape is circular (`yappCircle`), enable center-based positioning with `yappCenter`. This simplifies placement and matches v3.3.8 flags.

## What to change

- In the `pushButtons` examples where shape is `yappCircle`, add `yappCenter` in the named flags.
- Keep existing coordinate flags (e.g., `yappCoordPCB`) unchanged.

## Before

```scad
pushButtons = [
  [84.2, 30.7, 8, 8, 0, 2, 1, 3.5, yappCircle],
  [85,   13.5, 8, 5, 3, 2, 1, 3.5, yappRectangle]
];
```

## After (v3.3.8)

```scad
pushButtons = [
  [84.2, 30.7, 8, 8, 0, 2, 1, 3.5, yappCircle, yappCenter],
  [85,   13.5, 8, 5, 3, 2, 1, 3.5, yappRectangle]
];
```

Notes:
- If other named flags (e.g., `yappNoFillet`) are present, add `yappCenter` alongside them.

## Validation

- Local minimal test renders cleanly with `yappCenter` for circular caps.


