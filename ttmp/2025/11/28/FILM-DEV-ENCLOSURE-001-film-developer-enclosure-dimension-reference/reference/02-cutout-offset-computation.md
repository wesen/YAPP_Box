---
Title: Cutout Offset Computation in YAPP Generator
Ticket: FILM-DEV-ENCLOSURE-001
Status: active
Topics:
    - film-developer
    - enclosure
    - yapp
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Detailed explanation of how cutout positions and offsets are computed, especially for circle cutouts
LastUpdated: "2025-11-28T19:30:00.000000000-05:00"
---

# Cutout Offset Computation in YAPP Generator

## Goal

Explain in detail how `YAPPgenerator_v3.scad` computes cutout positions and offsets, with special focus on circle cutouts. This is critical because **circle cutouts do NOT specify the distance to the center of the circle**—they specify the bottom-left corner of the bounding box, and the generator applies an internal offset to position the circle correctly.

## Key Finding: Circle Cutouts Use Corner-Based Positioning

**Important:** When you specify a circle cutout position `(x, y)`, you are **NOT** specifying the center of the circle. Instead, you are specifying the **bottom-left corner** of the square bounding box that would contain the circle. The generator then internally translates the circle by `Radius` in both X and Y directions to position it correctly.

## Cutout Processing Pipeline

### Step 1: Extract Position and Dimensions (`processCutoutList_Face`)

```2334:2377:pkg/yappgen/assets/YAPPgenerator_v3.scad
module processCutoutList_Face(face, cutoutList, casePart, swapXY, swapWH, invertZ, rot_X, rot_Y, rot_Z, offset_x, offset_y, offset_z, wallDepth)
{
  for ( cutOut = cutoutList )
  {
    //-- Get the desired coordinate system    
    theCoordSystem = getCoordSystem(cutOut, yappCoordPCB);   
   
    theX = translate2Box_X (cutOut[0], face, theCoordSystem);
    theY = translate2Box_Y (cutOut[1], face, theCoordSystem);
    theWidth = cutOut[2];
    theLength = cutOut[3];
    theRadius = cutOut[4];
    theShape = cutOut[5];
    theDepth = getParamWithDefault(cutOut[6],0);
    theAngle = getParamWithDefault(cutOut[7],0);

    useCenterCoordinates = isTrue(yappCenter, cutOut);
    
    if (printMessages) echo("useCenterCoordinates", useCenterCoordinates);
    if (printMessages) echo("processCutoutList_Face", cutOut);

    //-- Calc H&W if only Radius is given
    tempWidth = (theShape == yappCircle) ?theRadius*2 : theWidth;
    tempLength = (theShape == yappCircle) ? theRadius*2 : theLength;
    
    base_width  = (swapWH) ? tempLength : tempWidth;
    base_height = (swapWH) ? tempWidth : tempLength;
    
    base_pos_H  = ((!swapXY) ? theY : theX);
    base_pos_V  = ((!swapXY) ? theX : theY);
  
		
    //-- Add 0.04 to the depth - we will shift by 0.02 later to center it on the wall
    base_depth  = (theDepth == 0) ? wallDepth + 0.04 : abs(theDepth) + 0.04;
    base_angle  = theAngle;

		//--Check for negative depth
		zAdjustForCutFromInside = !isTrue(yappFromInside, cutOut) ? 0 : wallDepth - base_depth;

    if (printMessages) echo ("---Box---");
    pos_X = base_pos_H;
    pos_Y = base_pos_V;
    
		processCutoutList_Mask(cutOut, rot_X, rot_Y, rot_Z, offset_x, offset_y, offset_z, wallDepth, base_pos_H, base_pos_V, base_width, base_height, base_depth, base_angle, pos_X, pos_Y, invertZ, zAdjustForCutFromInside);
		
  } //for ( cutOut = cutoutList )
} //-- processCutoutList_Face()
```

**Key observations:**

1. **Position extraction:** `theX` and `theY` are extracted from `cutOut[0]` and `cutOut[1]`, then transformed through `translate2Box_X` and `translate2Box_Y` to convert from the specified coordinate system (PCB, box, or boxInside) to box coordinates.

2. **Circle dimension handling:** For circles (`theShape == yappCircle`):
   - `tempWidth = theRadius * 2`
   - `tempLength = theRadius * 2`
   - This creates a square bounding box with side length equal to the circle diameter.

3. **Position assignment:** 
   - `base_pos_H` and `base_pos_V` are assigned from `theX` and `theY` (possibly swapped based on `swapXY`)
   - `pos_X = base_pos_H` and `pos_Y = base_pos_V`
   - **These positions are used directly without any offset for circles**

### Step 2: Shape Generation (`processCutoutList_Shape`)

```2281:2329:pkg/yappgen/assets/YAPPgenerator_v3.scad
module processCutoutList_Shape(cutOut, rot_X, rot_Y, rot_Z, offset_x, offset_y, offset_z, wallDepth,base_pos_H, base_pos_V, base_width, base_height, base_depth, base_angle, pos_X, pos_Y, invertZ, zAdjustForCutFromInside)
{
  theRadius = cutOut[4];
  theShape = cutOut[5];
  theAngle = getParamWithDefault(cutOut[7],0);
  
  zShift = invertZ ? -base_depth - zAdjustForCutFromInside : zAdjustForCutFromInside;
  
  //-- Output all of the current parameters
  if (printMessages) echo("base_pos_H",base_pos_H);
  if (printMessages) echo("base_pos_V",base_pos_V);
  if (printMessages) echo("base_width",base_width);
  if (printMessages) echo("base_height",base_height);
  if (printMessages) echo("base_depth",base_depth);
  if (printMessages) echo("wallDepth",wallDepth);

  if (printMessages) echo ("rot_X", rot_X); 
  if (printMessages) echo ("rot_Y", rot_Y); 
  if (printMessages) echo ("rot_Z", rot_Z); 
  if (printMessages) echo ("offset_x", offset_x); 
  if (printMessages) echo ("offset_y", offset_y); 
  if (printMessages) echo ("offset_z", offset_z); 
  if (printMessages) echo ("pos_X", pos_X); 
  if (printMessages) echo ("pos_Y", pos_Y); 
  if (printMessages) echo ("base_depth", base_depth); 
  if (printMessages) echo ("base_angle", base_angle);
  if (printMessages) echo ("invertZ", invertZ); 
  if (printMessages) echo ("zShift", zShift); 
  
  thePolygon = getVector(yappPolygonDef, cutOut);
  if (printMessages) echo("Polygon Definition", thePolygon=thePolygon);

  translate([offset_x, offset_y, offset_z]) 
  {
    rotate([rot_X, rot_Y, rot_Z])
    {
      translate([pos_X, pos_Y, wallDepth + zShift - 0.02]) 
      {
        if (printMessages) echo("Drawing cutout shape");
        // Draw the shape
          color("Fuchsia")
            generateShape (theShape,(isTrue(yappCenter, cutOut)), base_width, base_height, base_depth + 0.04, theRadius, theAngle, thePolygon);
      } //translate
    }// rotate
  } //translate
  
  if (printMessages) echo ("------------------------------");
    
} //-- processCutoutList_Shape()
```

**Key observations:**

1. **Translation chain:** The shape is positioned through a series of transformations:
   - First: `translate([offset_x, offset_y, offset_z])` - face-level offsets
   - Second: `rotate([rot_X, rot_Y, rot_Z])` - face-level rotation
   - Third: `translate([pos_X, pos_Y, wallDepth + zShift - 0.02])` - **This is where the position is applied**

2. **Shape generation:** `generateShape` is called with:
   - `useCenter` flag: `isTrue(yappCenter, cutOut)`
   - `base_width` and `base_height`: For circles, these are `radius * 2`
   - `theRadius`: The actual circle radius

### Step 3: Circle Drawing (`generateShape`)

```4592:4622:pkg/yappgen/assets/YAPPgenerator_v3.scad
module generateShape (Shape, useCenter, Width, Length, Thickness, Radius, Rotation, Polygon, expand=0)
// Creates a shape centered at 0,0 in the XY and from 0-thickness in the Z
{ 
  rotate([0,0,Rotation])
  {
    //-- Sphere cutout handled as a 3d object not a 2d Extruded 
    if (Shape == yappSphere) {
      //translate([0,0,(Thickness-0.08)/2]) 
      {
        intersection() 
        {
          translate([0,0,(Thickness/2)+.02]) // adjust to center 
        //  translate([0,0,0.04])
            cube([Radius*3,Radius*3,Thickness], center=true);
          
          //translate([0,0,Width+(Radius/2)-((Thickness-0.08)/2)])
          translate([0,0,Width+(Thickness/2)])
            sphere(r=Radius);
          
        } //intersection
      } //translate
    } else {
      linear_extrude(height = Thickness)
      { 
        offset(expand)
        { 
          if (Shape == yappCircle)
          {
            translate([(useCenter) ? 0 : Radius,(useCenter) ? 0 : Radius,0])
            circle(r=Radius);
          }
```

**Critical insight:** For circle cutouts, the translation is:

```scad
translate([(useCenter) ? 0 : Radius,(useCenter) ? 0 : Radius,0])
circle(r=Radius);
```

**This means:**

- **When `useCenter` is `false` (default):** The circle is translated by `Radius` in both X and Y directions. Since the circle is drawn centered at `(0,0)`, this translation moves it to `(Radius, Radius)`, which is the **center** of a square that starts at `(0,0)` and has side length `2*Radius`.

- **When `useCenter` is `true`:** No translation is applied, so the circle remains centered at `(0,0)`, meaning the position `(pos_X, pos_Y)` directly specifies the center.

## Complete Position Calculation for Circle Cutouts

### Default Behavior (useCenter = false)

1. **Input:** `cutOut[0] = x`, `cutOut[1] = y`, `cutOut[4] = radius`

2. **Coordinate transformation:**
   - `theX = translate2Box_X(x, face, coordSystem)`
   - `theY = translate2Box_Y(y, face, coordSystem)`

3. **Position assignment:**
   - `pos_X = theY` (or `theX` if swapped)
   - `pos_Y = theX` (or `theY` if swapped)

4. **Shape positioning:**
   - The shape is translated to `[pos_X, pos_Y, ...]`
   - Inside `generateShape`, the circle is further translated by `[Radius, Radius, 0]`
   - **Final circle center:** `(pos_X + Radius, pos_Y + Radius)`

5. **Conclusion:** The input position `(x, y)` represents the **bottom-left corner** of the bounding box, not the center.

### With yappCenter Flag (useCenter = true)

1. **Input:** Same as above, but `yappCenter` flag is set in the cutout definition

2. **Position assignment:** Same as above

3. **Shape positioning:**
   - The shape is translated to `[pos_X, pos_Y, ...]`
   - Inside `generateShape`, **no additional translation** is applied
   - **Final circle center:** `(pos_X, pos_Y)`

4. **Conclusion:** The input position `(x, y)` directly represents the **center** of the circle.

## Complete Analysis of All Cutout Shapes

The YAPP generator supports multiple cutout shapes, and they fall into two categories based on how they handle positioning offsets:

### Category 1: Circle-Based Shapes (Radius Offset)

These shapes use `Radius` as the offset value, positioning from the bottom-left corner of the bounding square:

#### 1. Circle (`yappCircle`)

```4618:4621:pkg/yappgen/assets/YAPPgenerator_v3.scad
          if (Shape == yappCircle)
          {
            translate([(useCenter) ? 0 : Radius,(useCenter) ? 0 : Radius,0])
            circle(r=Radius);
```

**Positioning:**
- **Default (useCenter=false):** Translated by `(Radius, Radius)` 
- **With useCenter=true:** No translation, position directly specifies center
- **Bounding box:** Square with side length `2*Radius`
- **Input position meaning:** Bottom-left corner of bounding square (default) or center (with flag)

#### 2. Ring (`yappRing`)

```4623:4642:pkg/yappgen/assets/YAPPgenerator_v3.scad
          else if (Shape == yappRing)
          {
            connectorCount=(Width==0) ? 0 : (Width>0) ? 1 : 2; 
            connectorWidth=abs(Width);
            translate([(useCenter) ? 0 : Radius,(useCenter) ? 0 : Radius,0])
              difference() {
                  difference() {
                      circle(r=Radius);
                      circle(r=Length);
                  }
                  if (connectorCount>0) 
                  {
                    square([connectorWidth, Radius*2], center=true);
                    if (connectorCount>1) 
                    {
                      rotate([0,0,90])
                      square([connectorWidth, Radius*2], center=true);
                    }
                  }
              }
```

**Positioning:**
- **Default (useCenter=false):** Translated by `(Radius, Radius)`
- **With useCenter=true:** No translation
- **Bounding box:** Square with side length `2*Radius` (outer radius)
- **Parameters:** `Radius` = outer radius, `Length` = inner radius, `Width` = connector width
- **Input position meaning:** Bottom-left corner of outer bounding square (default) or center (with flag)

#### 3. Circle with Flats (`yappCircleWithFlats`)

```4665:4674:pkg/yappgen/assets/YAPPgenerator_v3.scad
          else if (Shape == yappCircleWithFlats)
          {
            translate([(useCenter) ? 0 : Radius,(useCenter) ? 0 : Radius,0])
            {
              intersection()
              { 
                circle(r=Radius);    
                square([Width, Radius*2], center=true);
              }
            }
```

**Positioning:**
- **Default (useCenter=false):** Translated by `(Radius, Radius)`
- **With useCenter=true:** No translation
- **Bounding box:** Square with side length `2*Radius`
- **Parameters:** `Radius` = circle radius, `Width` = width of flat section
- **Shape:** Circle intersected with a horizontal rectangle to create flat sides
- **Input position meaning:** Bottom-left corner of bounding square (default) or center (with flag)

#### 4. Circle with Key (`yappCircleWithKey`)

```4676:4707:pkg/yappgen/assets/YAPPgenerator_v3.scad
          else if (Shape == yappCircleWithKey)
          {
            translate([(useCenter) ? 0 : Radius,(useCenter) ? 0 : Radius,0])
            {
              intersect = Radius - sqrt(Radius^2 - (Width/2)^2);   
              depth = Length;
              //--Add the Actual Key
              if (depth <= 0) 
              {
                //-- Create the circle with the flat for the key
                difference()
                {
                  circle(r=Radius); 
                    translate ([Radius ,0,0]) 
                      square([intersect*2, Width ], center=true);
                }
                //-- Add the outer cut
               translate ([Radius - intersect + 0,0,0]) 
                  square([abs(depth*2), Width ], center=true);
              }
              else if (depth > 0) 
              {
                //-- Create the circle with the flat for the key
                difference()
                {
                  circle(r=Radius);  
                  //-- Remove the flat
                  translate ([Radius - depth/2,0,0]) 
                    square([intersect*2 + depth, Width ], center=true);
                }
              }
            }
```

**Positioning:**
- **Default (useCenter=false):** Translated by `(Radius, Radius)`
- **With useCenter=true:** No translation
- **Bounding box:** Square with side length `2*Radius`
- **Parameters:** `Radius` = circle radius, `Width` = key width, `Length` = key depth
- **Shape:** Circle with a keyway (notch) cut out, positioned at the right side of the circle
- **Input position meaning:** Bottom-left corner of bounding square (default) or center (with flag)

### Category 2: Rectangle-Based Shapes (Half-Dimension Offset)

These shapes use `Width/2` and `Length/2` as offsets, positioning from the bottom-left corner of the rectangle:

#### 5. Rectangle (`yappRectangle`)

```4644:4649:pkg/yappgen/assets/YAPPgenerator_v3.scad
          else if (Shape == yappRectangle)
          {
            translate([(useCenter) ? 0 : Width/2,(useCenter) ? 0 : Length/2,0])
            {
              square([Width,Length], center=true); 
            }
```

**Positioning:**
- **Default (useCenter=false):** Translated by `(Width/2, Length/2)` to move from corner to center
- **With useCenter=true:** No translation, position directly specifies center
- **Bounding box:** Rectangle with dimensions `Width × Length`
- **Input position meaning:** Bottom-left corner of rectangle (default) or center (with flag)

#### 6. Rounded Rectangle (`yappRoundedRect`)

```4651:4656:pkg/yappgen/assets/YAPPgenerator_v3.scad
          else if (Shape == yappRoundedRect)
          {
            {
              translate([(useCenter) ? 0 : Width/2,(useCenter) ? 0 : Length/2,0])
              roundedRectangle2D(Width,Length,Radius);
            }
```

**Positioning:**
- **Default (useCenter=false):** Translated by `(Width/2, Length/2)`
- **With useCenter=true:** No translation
- **Bounding box:** Rectangle with dimensions `Width × Length`
- **Parameters:** `Width` = width, `Length` = length, `Radius` = corner radius
- **Shape:** Rectangle with rounded corners
- **Input position meaning:** Bottom-left corner of bounding rectangle (default) or center (with flag)

#### 7. Polygon (`yappPolygon`)

```4658:4663:pkg/yappgen/assets/YAPPgenerator_v3.scad
          else if (Shape == yappPolygon)
          {
            translate([(useCenter) ? 0 : Width/2,(useCenter) ? 0 : Length/2,0])
            scale([Width,Length,0]){
              polygon(Polygon);
            }
          }
```

**Positioning:**
- **Default (useCenter=false):** Translated by `(Width/2, Length/2)`
- **With useCenter=true:** No translation
- **Bounding box:** Rectangle with dimensions `Width × Length` (used for scaling)
- **Parameters:** `Width` = X-axis scale factor, `Length` = Y-axis scale factor, `Polygon` = polygon definition (normalized coordinates)
- **Shape:** Scaled polygon shape (hexagon, arrow, star, triangle, etc.)
- **Input position meaning:** Bottom-left corner of bounding rectangle (default) or center (with flag)
- **Note:** The polygon coordinates are normalized (typically -0.5 to +0.5), then scaled by `Width` and `Length`

### Category 3: 3D Shapes (Special Handling)

#### 8. Sphere (`yappSphere`)

```4597:4612:pkg/yappgen/assets/YAPPgenerator_v3.scad
    if (Shape == yappSphere) {
      //translate([0,0,(Thickness-0.08)/2]) 
      {
        intersection() 
        {
          translate([0,0,(Thickness/2)+.02]) // adjust to center 
        //  translate([0,0,0.04])
            cube([Radius*3,Radius*3,Thickness], center=true);
          
          //translate([0,0,Width+(Radius/2)-((Thickness-0.08)/2)])
          translate([0,0,Width+(Thickness/2)])
            sphere(r=Radius);
          
        } //intersection
      } //translate
```

**Positioning:**
- **Special handling:** This is a 3D object, not a 2D extruded shape
- **No XY translation:** The sphere is positioned at `(0,0)` in XY plane
- **Z positioning:** Complex Z-axis positioning based on `Width` and `Thickness`
- **Bounding box:** Cube with side length `Radius*3` used for intersection
- **Input position meaning:** The `(x, y)` position directly specifies the center (no offset applied in XY)
- **Note:** This shape behaves differently from 2D shapes and doesn't follow the same offset pattern

## Summary Table: Positioning Offsets by Shape

| Shape | Default Offset (X, Y) | With useCenter | Bounding Box | Category |
|-------|---------------------|----------------|--------------|----------|
| Circle | `(Radius, Radius)` | `(0, 0)` | `2*Radius × 2*Radius` | Circle-based |
| Ring | `(Radius, Radius)` | `(0, 0)` | `2*Radius × 2*Radius` | Circle-based |
| CircleWithFlats | `(Radius, Radius)` | `(0, 0)` | `2*Radius × 2*Radius` | Circle-based |
| CircleWithKey | `(Radius, Radius)` | `(0, 0)` | `2*Radius × 2*Radius` | Circle-based |
| Rectangle | `(Width/2, Length/2)` | `(0, 0)` | `Width × Length` | Rectangle-based |
| RoundedRect | `(Width/2, Length/2)` | `(0, 0)` | `Width × Length` | Rectangle-based |
| Polygon | `(Width/2, Length/2)` | `(0, 0)` | `Width × Length` | Rectangle-based |
| Sphere | `(0, 0)` | `(0, 0)` | N/A (3D) | Special |

## Key Insights

1. **Two offset patterns:**
   - **Circle-based shapes:** Use `Radius` offset (because they're conceptually squares with side `2*Radius`)
   - **Rectangle-based shapes:** Use `Width/2` and `Length/2` offsets

2. **Consistent behavior:** All 2D shapes follow the same pattern:
   - Default: Position specifies bottom-left corner of bounding box
   - With `yappCenter`: Position directly specifies center

3. **Sphere is special:** As a 3D object, it doesn't follow the 2D offset pattern and always uses center positioning in XY.

4. **Dimension handling:** For circle-based shapes, `base_width` and `base_height` are set to `radius*2` in `processCutoutList_Face`, creating a square bounding box.

## Practical Implications and Examples

### Example 1: Circle Cutout (Circle-Based Shape)

**Scenario:** You want a circle with radius 10mm centered at position (50, 30) on the lid.

**Option 1: Without yappCenter (default)**
```yaml
cutouts:
  - face: lid
    from_face_left: 40      # 50 - 10 (center - radius)
    from_face_back: 20      # 30 - 10 (center - radius)
    radius: 10
    shape: circle
```

**Option 2: With yappCenter**
```yaml
cutouts:
  - face: lid
    from_face_left: 50      # Direct center position
    from_face_back: 30      # Direct center position
    radius: 10
    shape: circle
    origin: center           # Sets yappCenter flag
```

### Example 2: Rectangle Cutout (Rectangle-Based Shape)

**Scenario:** You want a rectangle 20mm wide × 15mm long centered at position (50, 30) on the lid.

**Option 1: Without yappCenter (default)**
```yaml
cutouts:
  - face: lid
    from_face_left: 40      # 50 - 20/2 (center - width/2)
    from_face_back: 22.5    # 30 - 15/2 (center - length/2)
    width: 20
    length: 15
    shape: rectangle
```

**Option 2: With yappCenter**
```yaml
cutouts:
  - face: lid
    from_face_left: 50      # Direct center position
    from_face_back: 30      # Direct center position
    width: 20
    length: 15
    shape: rectangle
    origin: center           # Sets yappCenter flag
```

### Example 3: Rounded Rectangle Cutout

**Scenario:** You want a rounded rectangle 30mm wide × 20mm long with 5mm corner radius, centered at (50, 30).

**Without yappCenter (default):**
```yaml
cutouts:
  - face: lid
    from_face_left: 35      # 50 - 30/2
    from_face_back: 20      # 30 - 20/2
    width: 30
    length: 20
    radius: 5
    shape: rounded_rect
```

**With yappCenter:**
```yaml
cutouts:
  - face: lid
    from_face_left: 50      # Direct center
    from_face_back: 30      # Direct center
    width: 30
    length: 20
    radius: 5
    shape: rounded_rect
    origin: center
```

### Example 4: Circle with Flats

**Scenario:** You want a circle with radius 10mm and flat width 8mm, centered at (50, 30).

**Without yappCenter (default):**
```yaml
cutouts:
  - face: lid
    from_face_left: 40      # 50 - 10 (center - radius)
    from_face_back: 20      # 30 - 10 (center - radius)
    radius: 10
    width: 8                 # Flat width
    shape: circle_with_flats
```

**With yappCenter:**
```yaml
cutouts:
  - face: lid
    from_face_left: 50      # Direct center
    from_face_back: 30      # Direct center
    radius: 10
    width: 8
    shape: circle_with_flats
    origin: center
```

### Example 5: Polygon Cutout

**Scenario:** You want a hexagon scaled to 20mm × 20mm, centered at (50, 30).

**Without yappCenter (default):**
```yaml
cutouts:
  - face: lid
    from_face_left: 40      # 50 - 20/2
    from_face_back: 20      # 30 - 20/2
    width: 20               # X-axis scale
    length: 20              # Y-axis scale
    shape: polygon
    polygon: hexagon
```

**With yappCenter:**
```yaml
cutouts:
  - face: lid
    from_face_left: 50      # Direct center
    from_face_back: 30      # Direct center
    width: 20
    length: 20
    shape: polygon
    polygon: hexagon
    origin: center
```

### Why This Matters

1. **Documentation clarity:** Users need to understand whether they're specifying a corner or center position.

2. **Migration:** When converting from corner-based to center-based positioning:
   - **Circle-based shapes:** Add `radius` to both X and Y coordinates
   - **Rectangle-based shapes:** Add `width/2` to X and `length/2` to Y coordinates

3. **Precision:** For precise placement, using `origin: center` is more intuitive and less error-prone, especially when aligning multiple cutouts.

4. **Consistency:** Understanding the offset pattern helps ensure consistent positioning across different shape types.

## Mask Positioning (Special Case)

When masks are used, there's an additional offset calculation:

```2258:2266:pkg/yappgen/assets/YAPPgenerator_v3.scad
      centeroffsetH = (isTrue(yappCenter, cutOut)) ? 0 : base_width / 2;
      centeroffsetV = (isTrue(yappCenter, cutOut)) ? 0 : base_height / 2;
      zShift = invertZ ? -base_depth - zAdjustForCutFromInside : zAdjustForCutFromInside;
			
      translate([offset_x, offset_y, offset_z]) 
      {
        rotate([rot_X, rot_Y, rot_Z])
        {
           translate([base_pos_H + centeroffsetH, base_pos_V+centeroffsetV, wallDepth + zShift - 0.02])
```

For masks:
- **Without yappCenter:** Mask is positioned at `(base_pos_H + base_width/2, base_pos_V + base_height/2)` - the center of the bounding box
- **With yappCenter:** Mask is positioned at `(base_pos_H, base_pos_V)` - directly at the specified position

For circles, `base_width = base_height = radius * 2`, so `centeroffsetH = centeroffsetV = radius`, which aligns the mask center with the circle center.

## Summary

**Cutout positioning summary:**

### Two Offset Patterns

1. **Circle-based shapes** (Circle, Ring, CircleWithFlats, CircleWithKey):
   - **Default:** Input position `(x, y)` = bottom-left corner of bounding square
   - **Offset applied:** `(Radius, Radius)`
   - **Final center:** `(x + Radius, y + Radius)`
   - **Bounding box:** Square with side `2*Radius`

2. **Rectangle-based shapes** (Rectangle, RoundedRect, Polygon):
   - **Default:** Input position `(x, y)` = bottom-left corner of rectangle
   - **Offset applied:** `(Width/2, Length/2)`
   - **Final center:** `(x + Width/2, y + Length/2)`
   - **Bounding box:** Rectangle with dimensions `Width × Length`

### With yappCenter Flag

For **all shapes**, when `yappCenter` flag is set:
- Input position `(x, y)` = center directly
- No offset applied
- Final center = `(x, y)`

### Key Points

1. **Consistent behavior:** All 2D shapes follow the same pattern - default uses corner-based positioning, `yappCenter` uses center-based positioning.

2. **The offset happens inside `generateShape`:** The translation is applied after the shape is positioned at `(pos_X, pos_Y)`, moving it from the corner to the center.

3. **Sphere is special:** As a 3D object, it doesn't follow the 2D offset pattern and always uses center positioning in XY.

4. **This is NOT documented in the cutout position specification** - the position appears to be "just a coordinate" but the actual meaning depends on:
   - The shape type (determines offset pattern)
   - The `yappCenter` flag (determines whether offset is applied)

5. **Practical recommendation:** Use `origin: center` in YAML for more intuitive positioning, especially when aligning multiple cutouts or when precise center placement is required.

## Related

- `pkg/yappgen/assets/YAPPgenerator_v3.scad` - Main generator file
- `pkg/yappgen/modules/cutouts/module.go` - Go code that generates cutout arrays
- `pkg/yappgen/modules/cutouts/schema.yaml` - YAML schema for cutouts
- Reference document: `01-enclosure-dimension-derivation.md` - General dimension derivation

