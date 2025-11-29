# Changelog

## 2025-11-28

- Initial workspace created


## 2025-11-28

Added detailed reference document explaining cutout offset computation, especially for circle cutouts. Key finding: circle cutouts specify bottom-left corner of bounding box, not center, unless yappCenter flag is used.


## 2025-11-28

### Related Files

- reference/02-cutout-offset-computation.md — New reference document


## 2025-11-28

Extended cutout offset analysis to cover all shape types: circle-based shapes (Circle, Ring, CircleWithFlats, CircleWithKey) use Radius offset, rectangle-based shapes (Rectangle, RoundedRect, Polygon) use Width/2 and Length/2 offsets. Added comprehensive examples and summary table.


## 2025-11-28

Added DSL-side cutout position transformation document explaining how face-relative coordinates (from_face_left, from_face_bottom, from_face_back) transform into SCAD array positions, and how the origin flag affects positioning.

