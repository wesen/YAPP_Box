// Test: Push Buttons with yappCenter (center-based positioning)
include <../../../YAPPgenerator_v3.scad>

// Minimal test: rely on defaults; specify push button with circular cap
pushButtons =
[
  //  p(0) posx, p(1) posy, p(2) capLength, p(3) capWidth, p(4) capRadius,
  //  p(5) capAboveLid, p(6) switchHeight, p(7) switchTravel, p(8) poleDiameter,
  //  p(9) heightToTopOfPCB (optional), p(10) shape, placement flag
  [84, 30, 8, 8, 0, 2, 1, 3.5, 2, standoffHeight + pcbThickness, yappCircle, yappCenter]
];

YAPPgenerate();


