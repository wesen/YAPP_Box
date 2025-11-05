//-----------------------------------------------------------------------
// Yet Another Parameterized Projectbox generator
//
//  This is a box for Raspberry Pi Pico Temperature Monitor
//
//  Version: 1.0
//  Date: 2025-11-05
//
// This design is parameterized based on the size of a PCB.
//-----------------------------------------------------------------------

include <../../YAPPgenerator_v3.scad>

//---------------------------------------------------------
// This design is parameterized based on the size of a PCB.
//---------------------------------------------------------
// Note: length/lengte refers to X axis, 
//       width/breedte refers to Y axis,
//       height/hoogte refers to Z axis

/*
      padding-back|<------pcb length --->|<padding-front
                            RIGHT
        0    X-axis ---> 
        +----------------------------------------+   ---
        |                                        |    ^
        |                                        |   padding-right 
      Y |                                        |    v
      | |    -5,y +----------------------+       |   ---              
 B    a |         | 0,y              x,y |       |     ^              F
 A    x |         |                      |       |     |              R
 C    i |         |                      |       |     | pcb width    O
 K    s |         |                      |       |     |              N
        |         | 0,0              x,0 |       |     v              T
      ^ |    -5,0 +----------------------+       |   ---
      | |                                        |    padding-left
      0 +----------------------------------------+   ---
        0    X-as --->
                          LEFT
*/

//-- which part(s) do you want to print?
printBaseShell        = true;
printLidShell         = true;
printSwitchExtenders  = false;
printDisplayClips     = false;

// ********************************************************************
// PCB Dimensions - Based on custom PCB that accepts Raspberry Pi Pico
// The Pico is 51x21mm, but we'll use a slightly larger custom PCB

pcbLength           = 80;  // X-axis: front to back
pcbWidth            = 60;  // Y-axis: side to side
pcbThickness        = 1.6; 
standoffHeight      = 5.0; // Clearance for components and Pico stack (12mm total height)
standoffDiameter    = 7;
standoffPinDiameter = 2.4;
standoffHoleSlack   = 0.4;

//-------------------------------------------------------------------                            
//-- padding between pcb and inside wall
paddingFront        = 5;
paddingBack         = 5;
paddingRight        = 5;
paddingLeft         = 5;

//-- Edit these parameters for your own box dimensions
wallThickness       = 2.0;
basePlaneThickness  = 1.5;
lidPlaneThickness   = 1.5;

//-- Total height of box
//-- Space needed: PCB (1.6mm) + standoff (5mm) + Pico stack (12mm) + clearance (5mm) = ~24mm
baseWallHeight      = 20;
lidWallHeight       = 15;

//-- ridge where base and lid of box can overlap
ridgeHeight         = 5.0;
ridgeSlack          = 0.3;
ridgeGap            = 0.5;
roundRadius         = 3.0;

boxType             = 0; // All edges rounded

// Set the layer height of your printer
printerLayerHeight  = 0.2;

//---------------------------
//--     C O N T R O L     --
//---------------------------
// -- Render --
renderQuality             = 8;
previewQuality            = 5;

// -- Preview --
showSideBySide            = true;
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

//===================================================================
// *** PCB Stands ***
//-------------------------------------------------------------------
//  Default origin = yappCoordPCB: PCB[0,0,0]
//
//  Parameters:
//   Required:
//    p(0) = posx
//    p(1) = posy
//   Optional:
//    n(a) = { yappBoth | yappLidOnly | yappBaseOnly }
//    n(b) = { yappHole, <yappPin> }  -- Baseplate support
//    n(c) = { <yappAllCorners>, yappFrontLeft | yappFrontRight | yappBackLeft | yappBackRight }
//    n(d) = { yappAddFillet }

pcbStands = 
[
  // Four corners with pins
  [5, 5, yappAllCorners, yappPin, yappBoth]
];

//===================================================================
// *** Snap Joins ***
//-------------------------------------------------------------------
//  Default origin = yappCoordBox: box[0,0,0]
//
//  Parameters:
//   Required:
//    p(0) = posx | posy
//    p(1) = width
//   Optional:
//    n(a) = yappLeft / yappRight / yappFront / yappBack (one or more)
//    n(b) = { <yappSymmetric>, yappRectangle } 

snapJoins = 
[
  // Left and right sides
  [15, 5, yappLeft, yappRight],
  [pcbLength - 15, 5, yappLeft, yappRight]
];

//===================================================================
// *** Cutouts ***
//-------------------------------------------------------------------
//  Default origin = yappCoordPCB: PCB[0,0,0]
//
//  Parameters:
//   Required:
//    p(0) = from Back
//    p(1) = from Left
//    p(2) = width
//    p(3) = length
//    p(4) = radius
//    p(5) = shape
//  Optional:
//    p(6) = depth
//    p(7) = angle
//    n(a) = { yappPolygonDef }
//    n(b) = { yappMaskDef }
//    n(c) = { [yappMaskDef, hOffset, vOffset, rotation] }
//    n(d) = { <yappCoordPCB> | yappCoordBox | yappCoordBoxInside }
//    n(e) = { <yappOrigin>, yappCenter }

// Display cutout: 36x17mm visible area
// Center on PCB: X = (80-36)/2 = 22mm, Y = (60-17)/2 = 21.5mm

cutoutsLid = 
[
  // Display cutout (rectangular, 36x17mm) - centered on PCB
  [(pcbLength - 36) / 2, (pcbWidth - 17) / 2, 36, 17, 0, yappRectangle]
];

//===================================================================
// *** Cutouts in Base ***
//-------------------------------------------------------------------
// Temperature probe cable grommet (8mm probe + 6mm cable)
// Position at back-right corner for cable management

cutoutsBase = 
[
  // Cable grommet for temperature probe (10mm diameter for 8mm probe)
  // Positioned near back-right corner
  [pcbLength - 10, pcbWidth - 10, 0, 0, 5, yappCircle, 0, 0, yappCenter]
];

//===================================================================
// *** Cutouts in Front ***
//-------------------------------------------------------------------
// Button cutout: 28mm diameter, threaded barrel 45mm
// The button needs to go through the front panel
// Center vertically: Y = pcbWidth/2 = 30mm

cutoutsFront = 
[
  // Push button cutout (28mm diameter)
  // Position: centered vertically, 15mm from bottom
  [pcbWidth / 2, 15, 0, 0, 14, yappCircle, 0, 0, yappCenter]
];

//===================================================================
// *** Labels ***
//-------------------------------------------------------------------
//  Default origin = yappCoordBox: box[0,0,0]
//
//  Parameters:
//   p(0) = posx
//   p(1) = posy/z
//   p(2) = rotation degrees CCW
//   p(3) = depth : positive values go into case (Remove) negative values are raised (Add)
//   p(4) = plane {yappLeft, yappRight, yappFront, yappBack, yappLid, yappBase}
//   p(5) = font
//   p(6) = size
//   p(7) = "label text"

labelsPlane = 
[
  // Project label on lid
  [10, 10, 0, 0.5, yappLid, "Liberation Sans:style=Bold", 5, "Pico Temp Monitor"],
  
  // Version on base
  [10, 10, 0, 0.5, yappBase, "Liberation Sans", 3, "v1.0"]
];

//===================================================================
// *** Light Tubes ***
//-------------------------------------------------------------------
// None needed for this project

lightTubes = [];

//===================================================================
// *** Push Buttons ***
//-------------------------------------------------------------------
// Not using push buttons - using panel-mount button with cutout

pushButtons = [];

//===================================================================
// *** Connectors ***
//-------------------------------------------------------------------
// None needed - using cutouts for all external connections

connectors = [];

//===================================================================
// *** Generate the Box ***
//===================================================================
YAPPgenerate();
