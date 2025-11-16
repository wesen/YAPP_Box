//-----------------------------------------------------------------------
// Yet Another Parameterized Projectbox generator
//
//  This is a YAPP_Test_buttons_v30 test box - COMPUTED VALUES VERSION
//
//  Version for debugging: all shellWidth/shellHeight expressions replaced
//  with their computed values to isolate differences from DSL output
//
//-----------------------------------------------------------------------

//-- Bambu Lab X1C 0.4mm Nozzle XT-Copolyester
insertDiam  = 3.8 + 0.5;
screwDiam   = 2.5 + 0.5;
  

include <../YAPPgenerator_v3.scad>

// Note: length/lengte refers to X axis, 
//       width/breedte to Y, 
//       height/hoogte to Z

//-- which part(s) do you want to print?
printBaseShell        = true;
printLidShell         = true;
printSwitchExtenders  = true;

//-- pcb dimensions -- very important!!!
pcbLength           = 30;
pcbWidth            = 40;
pcbThickness        = 1.6;
                            
//-- padding between pcb and inside wall
paddingFront        = 1;
paddingBack         = 1;
paddingRight        = 1;
paddingLeft         = 1;

//-- Edit these parameters for your own box dimensions
wallThickness       = 1.4;
basePlaneThickness  = 1.5;
lidPlaneThickness   = 1.0;

//-- Total height of box = basePlaneThickness + lidPlaneThickness 
//--                     + baseWallHeight + lidWallHeight
//-- space between pcb and lidPlane :=
//--      (bottonWallHeight+lidWallHeight) - (standoffHeight+pcbThickness)
baseWallHeight      = 10;
lidWallHeight       = 10;

//-- ridge where base and lid off box can overlap
//-- Make sure this isn't less than lidWallHeight
ridgeHeight         = 3.0;  //-> at least 1.8 * wallThickness
ridgeSlack          = 0.2;
roundRadius         = 2.0;

//-- How much the PCB needs to be raised from the base
//-- to leave room for solderings and whatnot
standoffHeight      = 3.0;
standoffPinDiameter = 2.4;
standoffHoleSlack   = 0.4;
standoffDiameter    = 6;


//-- C O N T R O L -------------//-> Default ---------
showSideBySide      = false;     //-> true
previewQuality      = 5;        //-> from 1 to 32, Default = 5
renderQuality       = 5;        //-> from 1 to 32, Default = 8
onLidGap            = 0;
shiftLid            = 1;
hideLidWalls        = false;    //-> false
hideBaseWalls       = false;    //-> false
colorBase           = "yellow";
alphaBase           = 0.8;//0.2;   
colorLid            = "silver";
alphaLid            = 0.8;//0.2;   
showOrientation     = true;
showPCB             = true;
showSwitches        = true;
showPCBmarkers      = false;
showShellZero       = false;
showCenterMarkers   = false;
inspectX            = 0;        //-> 0=none (>0 from Back)
inspectY            = 0;        //-> 0=none (>0 from Right)
inspectZ            = 0;        //-> 0=none (>0 from Base)
inspectXfromBack    = false;     //-> View from the inspection cut foreward
inspectYfromLeft    = true;     //-> View from the inspection cut to the right
inspectZfromTop     = true;     //-> View from the inspection cut down
//-- C O N T R O L ---------------------------------------


//===================================================================
// *** PCB Supports ***
// Pin and Socket standoffs 
//-------------------------------------------------------------------
//  COMPUTED VALUES VERSION - using literal values instead of expressions
//-------------------------------------------------------------------
pcbStands =
[
    //-- 0, 1,
        [5, 5, yappBaseOnly, yappFrontLeft, yappBackRight] 
      , [5, 5, yappBoth, yappBackLeft, yappFrontRight]
];


//===================================================================
//  *** Cutouts ***
//  COMPUTED VALUES VERSION - shellWidth=44.8, shellHeight=22.5
//-------------------------------------------------------------------
cutoutsBase =   
[
    // Original: [shellLength/2,shellWidth/2 ,25,25, 5, yappPolygon ,0 ,30, yappCenter, shapeHexagon, [maskHexCircles,0,5]]
    // Using rounded_rect placeholder since polygon+mask not in DSL yet
    [17.4, 22.4, 25, 25, 5, yappRoundedRect, 0, 0, yappCenter]
];

// (0) = posy
// (1) = posz
cutoutsFront =  
[
//-- Original: [3, 2, shellWidth-6, shellHeight-4, 2, yappRoundedRect]
//-- Computed: shellWidth-6 = 44.8-6 = 38.8, shellHeight-4 = 22.5-4 = 18.5
    [3, 2, 38.8, 18.5, 2, yappRoundedRect]
];

// (0) = posy
// (1) = posz
cutoutsBack =   
[
//-- Original: [5, 2, shellWidth-10, shellHeight-4, 3, yappRoundedRect]
//-- Computed: shellWidth-10 = 44.8-10 = 34.8, shellHeight-4 = 22.5-4 = 18.5
    [5, 2, 34.8, 18.5, 3, yappRoundedRect]
];


cutoutsLeft =  
[
];


//===================================================================
//  *** Snap Joins ***
//-------------------------------------------------------------------
//  Original used shellLength/2 expressions with yappCenter/yappSymmetric
//  Not yet implemented in DSL
//-------------------------------------------------------------------
snapJoins   =   
[
    // Original: [(shellLength/2)-10, 4, yappLeft, yappCenter, yappSymmetric]
    // Computed: 17.4-10 = 7.4
    [7.4, 4, yappLeft, yappCenter, yappSymmetric]
   // Original: [(shellLength/2)-12, 4, yappRight, yappCenter, yappRectangle, yappSymmetric]
   // Computed: 17.4-12 = 5.4
   ,[5.4, 4, yappRight, yappCenter, yappRectangle, yappSymmetric]
];


//===================================================================
//  *** Box Mounts ***
//  Original had these commented out
//-------------------------------------------------------------------
boxMounts   =  
[
 // [(shellLength/2)-0, 3, 6, 2.5, yappLeft, yappRight, yappCenter]
 //,[(shellLength/2)-0, 3, -2, 2.5, yappLeft, yappRight, yappLid, yappCenter]
];
                                

//===================================================================
//  *** Push Buttons ***
//-------------------------------------------------------------------
pushButtons = 
[
 //-- 0,  1, 2, 3, 4, 5,   6, 7,   8
    [15, 30, 0, 0, 4, 0, 3,   1, 3.5, undef, yappCircle]
   ,[15, 10, 8, 6, 0, 3, 5.5, 1, 3.5, undef, yappRectangle]
];     
             


//========= MAIN CALL's ===========================================================


//---- This is where the magic happens ----
YAPPgenerate();

