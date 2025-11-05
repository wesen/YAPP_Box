// Test: Light Tubes with yappCenter (center-based positioning)
include <../../../YAPPgenerator_v3.scad>

// Rely on generator defaults; define only the feature under test
lightTubes =
[
  //  p(0) posx, p(1) posy, p(2) tubeLength, p(3) tubeWidth,
  //  p(4) tubeWall, p(5) gapAbovePcb, p(6) tubeType, p(7) lensThickness,
  //  n(a) coord,    placement flag
  [15, 20, 5, 6, 1, 0.5, yappCircle, 0.5, yappCoordPCB, yappCenter]
];

YAPPgenerate();


