---
Title: User Case Study - buttons2 Demo Discrepancies
Ticket: YAPP-DSL-GAPS-001
Status: active
Topics:
  - yapp
  - dsl
  - analysis
  - debugging
DocType: analysis
Intent: short-term
Owners: []
RelatedFiles:
  - Path: ../../../examples/YAPP_Demo_buttons2_v31.scad
    Note: Original SCAD file user is trying to replicate
  - Path: ../../../examples/test-push-buttons.yaml
    Note: User's DSL attempt
ExternalSources: []
Summary: Root cause analysis of visual differences between DSL-generated and original SCAD enclosures
LastUpdated: 2025-11-16
---

# User Case Study - buttons2 Demo Discrepancies

## User's Observations

When rendering DSL-generated SCAD vs. original `YAPP_Demo_buttons2_v31.scad`:

1. **"I have hinges"** - User sees hinge-like structures
2. **"The side of the box are open"** - Walls missing or split
3. **"Little plates sticking out to fasten it"** - External mounting tabs
4. **"Holes in one side"** - Cutouts visible
5. **"4 stalactites from top to pcb pillars (of which there are 4), but the scad one has only 2"** - Standoff count mismatch

## Root Cause Analysis

### Issue 1: Standoff Count Mismatch (4 vs 1)

**SCAD file defines:**
```openscad
pcbStands = [
  [5, 5]
];
```

**Only 1 standoff explicitly defined at position [5, 5]**

**User sees:** 4 standoffs

**Possible causes:**

**A. YAPP Auto-Corner Generation (Most Likely)**

YAPPgenerator v3 has logic that auto-generates corner standoffs when certain conditions are met. Checking the SCAD file parameters:

```openscad
pcbLength = 100;
pcbWidth = 100;
standoffHeight = 3.0;
standoffDiameter = 6;
```

**Hypothesis:** YAPP might have a feature where defining a single standoff triggers auto-corner placement, OR there's a global setting that enables this.

**Action needed:** Search `YAPPgenerator_v3.scad` for auto-corner logic

**B. User's DSL Has 4 Explicit Entries**

The user's test YAML (`examples/test-push-buttons.yaml`) defines 4 pcb_stands:
```yaml
pcb_stands:
  - x: 3, y: 3
  - x: pcb_length - 3, y: 3
  - x: 3, y: pcb_width - 3
  - x: pcb_length - 3, y: pcb_width - 3
```

**If this is the case:** The DSL is working correctly, but the SCAD file achieves the same result with less code using corner placement flags.

**Missing DSL feature:** `yappAllCorners` flag that auto-generates 4 standoffs from 1 definition

### Issue 2: "Hinges" and "Sides Are Open"

**SCAD file has NO ridgeExt definitions** (all empty arrays)

**SCAD file HAS cutouts:**
```openscad
cutoutsFront = [
  [3, 2, shellWidth-6, shellHeight-4, 2, yappRoundedRect]
];

cutoutsBack = [
  [5, 2, shellWidth-10, shellHeight-4, 3, yappRoundedRect]
];
```

**Analysis:**

The cutouts use `shellWidth-6` and `shellHeight-4`, making them nearly full-wall cutouts. These create large openings that might look like "open sides."

**"Hinges" observation:** Could be:
1. **Ridge itself** - The ridge between base and lid can look like a hinge mechanism
2. **Snap joins** - If snap_joins were defined (they're not in this file)
3. **Misinterpretation** - Large cutouts creating visual appearance of split/hinged design

**Conclusion:** NOT actual hinges. The "open sides" are the large cutouts. True hinges would require ridgeExt arrays.

### Issue 3: "Little Plates Sticking Out"

**SCAD file has NO boxMounts defined** (empty array)

**User observation:** Seeing mounting plates

**Possible causes:**

**A. Default YAPP Behavior**

YAPPgenerator might add default mounting features when certain conditions are met.

**B. User Added boxMounts to DSL**

If user manually added boxMounts to their YAML, the DSL would reject it (not implemented).

**C. Visual Misinterpretation**

Could be:
- **Snap join clips** (if defined)
- **Connector standoffs** protruding
- **Ridge overhang** appearing as tabs

**Most likely:** User's DSL YAML doesn't match the SCAD file exactly, OR YAPP has default features we're not aware of.

**Action needed:** Compare user's actual DSL YAML with SCAD file line-by-line

### Issue 4: "Holes in One Side"

**SCAD file defines cutouts on front and back:**
```openscad
cutoutsFront = [[3, 2, shellWidth-6, shellHeight-4, 2, yappRoundedRect]];
cutoutsBack = [[5, 2, shellWidth-10, shellHeight-4, 3, yappRoundedRect]];
```

**DSL should handle this** - cutouts module is implemented

**If user sees holes but didn't define them:** Check if DSL YAML has cutouts section

## Feature Gaps Preventing Parity

### Critical Gaps

1. **Corner Placement Flags** (yappAllCorners, yappFrontLeft, etc.)
   - **Impact:** Can't auto-generate 4 standoffs from 1 definition
   - **Workaround:** Manually define all 4 positions with expressions
   - **Priority:** HIGH - Common pattern in YAPP examples

2. **Shell Part Flags** (yappBoth, yappLidOnly, yappBaseOnly)
   - **Impact:** Can't control which shell part gets feature
   - **Use case:** Standoffs in base only, connectors through both
   - **Priority:** HIGH - Essential for proper standoff behavior

3. **boxMounts Array**
   - **Impact:** Can't add external mounting tabs
   - **Use case:** Wall mounting, DIN rail, external fastening
   - **Priority:** MEDIUM - Not in buttons2 demo, but user sees them

4. **ridgeExt Arrays**
   - **Impact:** Can't create split openings or hinged sections
   - **Use case:** Cable pass-throughs, split-level enclosures
   - **Priority:** MEDIUM - Would explain "hinge" observation

5. **Standoff Treatment Flags** (yappPin, yappHole, yappTopPin)
   - **Impact:** Can't control pin vs hole on each shell part
   - **Default:** yappPin (pin on base, hole on lid)
   - **Priority:** MEDIUM - Needed for advanced standoff configurations

### Secondary Gaps

6. **lightTubes** - LED light pipes
7. **labelsPlane** - Text labels
8. **displayMounts** - Display module mounting
9. **Mask support** - Honeycomb/bar ventilation patterns
10. **Multi-PCB support** - pcb array + yappPCBName flags

## Recommended Next Steps

### Immediate (to match buttons2 demo):

1. **Add corner placement flags to pcb_stands schema:**
   ```yaml
   corner_placement:
     type: string
     enum: [all_corners, front_left, front_right, back_left, back_right]
     desc: Auto-place at specified corners
   ```

2. **Add shell part flag to pcb_stands:**
   ```yaml
   shell_part:
     type: string
     enum: [both, lid_only, base_only]
     default: both
     desc: Which shell part gets the standoff
   ```

3. **Add standoff treatment flag:**
   ```yaml
   treatment:
     type: string
     enum: [pin, hole, top_pin]
     default: pin
     desc: Pin on base/hole on lid, holes on both, or reversed
   ```

4. **Implement boxMounts module** (if user actually needs it)

### Short-term (common patterns):

5. Add yappCenter/yappOrigin flags to snap_joins
6. Add yappSymmetric flag to snap_joins
7. Add yappAltOrigin support to cutouts/push_buttons

### Medium-term (advanced features):

8. Implement lightTubes module
9. Implement labelsPlane module
10. Implement ridgeExt modules (4 arrays)
11. Implement displayMounts module
12. Add mask support to cutouts (yappMaskDef)

## Verification Steps

To confirm root causes:

1. **Check user's actual DSL YAML:**
   ```bash
   cat <user's yaml file>
   ```
   Compare with YAPP_Demo_buttons2_v31.scad

2. **Render both and compare:**
   ```bash
   openscad -o /tmp/scad-original.stl examples/YAPP_Demo_buttons2_v31.scad
   openscad -o /tmp/dsl-generated.stl <user's generated scad>
   ```

3. **Check standoff count in generated SCAD:**
   ```bash
   grep -A 10 "^pcbStands" <user's generated scad>
   ```

4. **Search for boxMounts in generated output:**
   ```bash
   grep "boxMounts" <user's generated scad>
   ```

5. **Check for ridgeExt:**
   ```bash
   grep "ridgeExt" <user's generated scad>
   ```

## Conclusion

The DSL can generate the **push buttons** from buttons2 demo, but is missing:

1. **Corner placement automation** - Causing manual 4-standoff definitions vs elegant 1-with-flags
2. **boxMounts** - External mounting tabs (if user needs them)
3. **ridgeExt** - Split openings (if that's what "hinges" refers to)
4. **Shell part control** - yappBoth/yappLidOnly/yappBaseOnly flags

**Next action:** User should share their actual DSL YAML and generated SCAD so we can pinpoint exact gaps.
