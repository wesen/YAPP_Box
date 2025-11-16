//-----------------------------------------------------------------------
// Yet Another Parameterized Projectbox generator
//
//    Comparison demo for polygon cutouts (DSL vs legacy SCAD)
//
//    Version 3.x
//-----------------------------------------------------------------------

include <../YAPPgenerator_v3.scad>

// Which part(s) to print
printBaseShell        = true;
printLidShell         = true;
printSwitchExtenders  = false;

//-- pcb dimensions
pcbLength           = 30;
pcbWidth            = 40;
pcbThickness        = 1.6;

//-- padding between pcb and inside wall
paddingFront        = 1;
paddingBack         = 1;
paddingRight        = 1;
paddingLeft         = 1;

//-- Box dimensions
wallThickness       = 2.0;
basePlaneThickness  = 1.0;
lidPlaneThickness   = 1.0;
baseWallHeight      = 8;
lidWallHeight       = 13;
ridgeHeight         = 3.6;
ridgeSlack          = 0.2;
roundRadius         = 2.0;

//-- PCB mounting (defaults used elsewhere)
standoffHeight      = 7.0;
standoffPinDiameter = 2.4;
standoffHoleSlack   = 0.4;
standoffDiameter    = 7;

//===================================================================
//  *** Cutouts (Polygons for comparison) ***
//-------------------------------------------------------------------
//  Default origin = yappCoordBox: box[0,0,0]
//
// Base: hexagon; Front: 6pt star; Lid: arrow — matches DSL YAML
//-------------------------------------------------------------------
cutoutsBase =
[
  [15, 15, 25, 25, 0, yappPolygon, 0, 30, shapeHexagon]
];

cutoutsFront =
[
  [5, 15, 15, 15, 0, yappPolygon, undef, undef, shape6ptStar]
];

cutoutsLid =
[
  [10, 20, 20, 20, 0, yappPolygon, undef, undef, shapeArrow]
];

//---- This is where the magic happens ----
YAPPgenerate();


