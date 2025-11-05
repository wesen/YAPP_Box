//-----------------------------------------------------------------------
// Raspberry Pi Pico Temperature Monitor - LID ONLY
//-----------------------------------------------------------------------

include <../../YAPPgenerator_v3.scad>

//-- Print only the lid shell
printBaseShell        = false;
printLidShell         = true;
printSwitchExtenders  = false;
printDisplayClips     = false;

// PCB Dimensions
pcbLength           = 80;
pcbWidth            = 60;
pcbThickness        = 1.6; 
standoffHeight      = 5.0;
standoffDiameter    = 7;
standoffPinDiameter = 2.4;
standoffHoleSlack   = 0.4;

paddingFront        = 5;
paddingBack         = 5;
paddingRight        = 5;
paddingLeft         = 5;

wallThickness       = 2.0;
basePlaneThickness  = 1.5;
lidPlaneThickness   = 1.5;

baseWallHeight      = 20;
lidWallHeight       = 15;

ridgeHeight         = 5.0;
ridgeSlack          = 0.3;
ridgeGap            = 0.5;
roundRadius         = 3.0;

boxType             = 0;
printerLayerHeight  = 0.2;

renderQuality             = 8;
previewQuality            = 5;

showSideBySide            = false;
onLidGap                  = 5;
shiftLid                  = 5;
hideLidWalls              = false;
colorLid                  = "YellowGreen";
hideBaseWalls             = false;
colorBase                 = "BurlyWood";
showPCB                   = true;
showMarkers               = false;
inspectX                  = 0;
inspectY                  = 0;
inspectLightTubes         = 0;
inspectButtons            = 0;

pcbStands = 
[
  [5, 5, yappAllCorners, yappPin, yappBoth]
];

snapJoins = 
[
  [15, 5, yappLeft, yappRight],
  [pcbLength - 15, 5, yappLeft, yappRight]
];

cutoutsLid = 
[
  [(pcbLength - 36) / 2, (pcbWidth - 17) / 2, 36, 17, 0, yappRectangle]
];

labelsPlane = 
[
  [10, 10, 0, 0.5, yappLid, "Liberation Sans:style=Bold", 5, "Pico Temp Monitor"]
];

lightTubes = [];
pushButtons = [];
connectors = [];

YAPPgenerate();
