---
Title: 2025-11-16 Implementation Diary: Polygon Cutout Shapes
Ticket: YAPP-DSL-GAPS-001
Status: active
Topics:
    - yapp
    - dsl
    - implementation
    - cutouts
    - polygon
    - shapes
DocType: log
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/yappgen/modules/cutouts/schema.yaml
      Note: Added polygon to shape enum and polygon preset field
    - Path: pkg/yappgen/modules/cutouts/module.go
      Note: Added polygonPresetFlag function and polygon preset emission logic
    - Path: pkg/yappgen/schema.go
      Note: Added polygon support to ShapeFlag function for legacy code path
    - Path: pkg/yappgen/map.go
      Note: Added polygon preset support to buildCutoutParams function
    - Path: examples/yapp-demo-polygon-test.yaml
      Note: Test YAML file with various polygon presets
ExternalSources: []
Summary: Complete implementation diary documenting polygon cutout shape implementation with preset support
LastUpdated: 2025-11-16
---

# Implementation Diary: Polygon Cutout Shapes

**Date:** 2025-11-16  
**Task:** Implement polygon cutout shapes with preset support (Task #13 from YAPP-DSL-GAPS-001)  
**Status:** ✅ Complete  
**Time:** ~1.5 hours

## Objective

Implement support for polygon cutout shapes (`yappPolygon`) with preset shape definitions (hexagon, arrow, 6pt_star, iso_triangle, etc.). This enables users to create custom-shaped cutouts beyond basic rectangles and circles.

## Reference Materials Used

1. **YAPP_Template_v3.scad** - Parameter definitions for cutouts
2. **YAPPgenerator_v3.scad** (lines 376-392) - Preset shape definitions
3. **YAPP_Demo_buttons_v30.scad** - Example using `yappPolygon` with `shapeHexagon`
4. **YAPP_Reference_Shapes_v30.scad** - Examples of all polygon presets
5. **Existing cutouts module** - Studied current implementation patterns

## Step-by-Step Implementation Process

### Step 1: Study the YAPP Specification

**What I did:**
- Read YAPP documentation for polygon cutouts
- Analyzed `YAPPgenerator_v3.scad` to find preset shape definitions
- Identified available presets: hexagon, arrow, 6pt_star, iso_triangle, iso_triangle2, triangle, triangle2

**Key findings:**
- Polygon cutouts require `yappPolygon` flag + a preset shape reference
- Format: `[fromBack, fromLeft, width, length, radius, yappPolygon, depth, angle, shapePreset, ...flags]`
- Preset comes after depth/angle but before other flags
- Available presets defined in YAPPgenerator_v3.scad lines 376-392

**What worked:** Clear examples in YAPP_Reference_Shapes_v30.scad showed exact usage patterns.

**What didn't work:** N/A - documentation was clear.

### Step 2: Update Schema YAML

**What I did:**
- Added `polygon` to the shape enum in `pkg/yappgen/modules/cutouts/schema.yaml`
- Added `polygon` field (optional string) for preset name
- Added enum values for all 7 preset shapes
- Added test cases for hexagon and arrow polygons

**Schema changes:**
```yaml
shape:
  enum: [rectangle, circle, rounded_rect, circle_with_flats, circle_with_key, polygon]

polygon:
  type: string
  enum: [hexagon, arrow, 6pt_star, iso_triangle, iso_triangle2, triangle, triangle2]
```

**What worked:** Schema validation passed immediately.

**What didn't work:** Initial YAML syntax error - description had colon that needed quoting.

**Fix applied:** Quoted the description string.

**What I learned:** Always quote YAML strings that contain special characters like colons.

### Step 3: Generate Code with schemagen

**What I did:**
```bash
go run ./cmd/schemagen validate pkg/yappgen/modules/cutouts/schema.yaml
go run ./cmd/schemagen discover
```

**What worked:** 
- Schema validation passed after fixing YAML syntax
- Code generation created `schema_gen.go` with `Polygon *string` field

**What didn't work:** N/A - code generation worked smoothly.

**What I learned:** Schema-driven code generation makes adding fields straightforward.

### Step 4: Update Module Build Function

**What I did:**
- Updated `cutoutShapeFlag()` to handle `polygon` shape
- Added `polygonPresetFlag()` function to map DSL preset names to YAPP constants
- Updated `Build()` function to emit polygon preset after depth/angle

**Builder changes:**
```go
// For polygon shapes, add the preset shape after depth/angle
if strings.ToLower(strings.TrimSpace(item.Shape)) == "polygon" {
    if item.Polygon == nil || *item.Polygon == "" {
        return nil, errors.Errorf("%s: polygon preset is required when shape=polygon", label)
    }
    presetFlag, err := polygonPresetFlag(*item.Polygon)
    if err != nil {
        return nil, errors.Wrapf(err, "%s", label)
    }
    params = append(params, presetFlag)
}
```

**What worked:** Module code compiled and logic was straightforward.

**What didn't work:** N/A - implementation worked on first try.

**What I learned:** Following existing patterns (like lightTubes module) made implementation straightforward.

### Step 5: Update Legacy Code Path

**What I did:**
- Discovered that cutouts use `buildCutoutParams()` in `map.go`, not the module Build function
- Updated `ShapeFlag()` in `schema.go` to support polygon
- Added `polygonPresetFlag()` function to `map.go`
- Updated `buildCutoutParams()` to emit polygon preset

**Critical discovery:** Cutouts are processed by `distributeCutouts()` which calls `buildCutoutParams()` in `map.go`, NOT the module's Build function. This is a legacy code path that still needs to be maintained.

**What worked:** After finding the correct code path, implementation was straightforward.

**What didn't work:** 
- Initially only updated module Build function, but that wasn't being called
- Error message "unsupported cutout shape: polygon" led me to discover `ShapeFlag()` in `schema.go`
- Had to trace through code to find `buildCutoutParams()` in `map.go`

**Fix applied:** Updated both code paths (module Build function and legacy `buildCutoutParams()`).

**What I learned:** 
- Always trace through the actual code execution path, not just the module code
- Legacy code paths may still be in use even when module system exists
- Error messages can help locate where code is actually executing

### Step 6: Test End-to-End

**What I did:**
- Created `examples/yapp-demo-polygon-test.yaml` with hexagon, arrow, and 6pt_star polygons
- Tested resolve → generate pipeline
- Verified SCAD output includes polygon presets

**Test YAML:**
```yaml
cutouts:
  - face: base
    from_back: 15
    from_left: 15
    width: 25
    length: 25
    radius: 5
    shape: polygon
    polygon: hexagon
    depth: 0
    angle: 30
```

**What worked:**
- Resolution succeeded
- Generation succeeded
- Output matched YAPP format:
  ```scad
  cutoutsBase = [
    [15, 15, 25, 25, 0, yappPolygon, 0, 30, shapeHexagon]
  ];
  ```

**What didn't work:** N/A - end-to-end test passed.

**What I learned:** End-to-end testing is critical - catches integration issues that unit tests might miss.

## Step 7: Legacy Removal and Migration to Module Build

**What I did:**
- Migrated cutouts emission to use the module `Build()` directly from `features.go` (removed the special-case legacy pipeline)
- Changed `Model` to store `Cutouts []map[string]any` (removed the `Cutout` struct)
- Removed legacy functions from `pkg/yappgen/map.go`: `buildCutoutParams`, `distributeCutouts`, and polygon preset helpers
- Updated tests in `pkg/yappgen/yappgen_test.go` to target module `Build()` and fixed type assertions (float vs int)

**Why:**
- Single authoritative code path per feature (module Build) reduces drift and bugs
- Simplifies maintenance; no special-case for cutouts in the core layer

**Outcome:**
- All tests pass (`go test ./pkg/yappgen/...`)
- SCAD output unchanged for existing examples

**Follow-ups (tracked as tasks):**
- Evaluate removing `ShapeFlag`/`SnapSideFlag` and ParamSpec schemas in `schema.go` if no longer used
- Migrate any remaining tests away from ParamSpec-based helpers and then remove them

## Challenges Encountered

### Challenge 1: Finding the Correct Code Path

**Problem:** Updated module Build function but polygon presets weren't being emitted.

**Root cause:** Cutouts use legacy `buildCutoutParams()` in `map.go`, not the module Build function.

**Solution:** Updated both code paths - module Build function (for future) and legacy `buildCutoutParams()` (for current).

**Prevention:** 
- Trace through actual execution path before implementing
- Check which FeatureModule is registered for cutouts
- Look at `features.go` to see how cutouts are processed

**Future improvement:** Consider migrating cutouts fully to module system to avoid maintaining two code paths.

### Challenge 2: YAML Schema Syntax Error

**Problem:** Schema validation failed with "mapping values are not allowed in this context".

**Root cause:** Description string contained colon without quotes.

**Solution:** Quoted the description string.

**Prevention:** Always quote YAML strings containing special characters.

**Future improvement:** Consider improving schemagen to handle this automatically or provide better error messages.

### Challenge 3: LightTubes Schema Type Mismatch

**Problem:** Compilation error in lightTubes module: `cannot use &v (value of type *int) as *float64`.

**Root cause:** Generated code had `v := 0` (int) but field type is `*float64`.

**Solution:** Manually fixed generated code to use `0.0` instead of `0`.

**Prevention:** This is a known schemagen limitation - always use `0.0` for numeric defaults.

**Future improvement:** Fix schemagen to infer float64 from field type for defaults.

## What Worked Well

1. **Following existing patterns** - Studying lightTubes module made implementation straightforward
2. **Schema-driven approach** - Having schema.yaml define structure made validation smooth
3. **Clear YAPP documentation** - Examples in YAPP_Reference_Shapes_v30.scad showed exact usage
4. **End-to-end testing** - Validated the complete pipeline

## What Didn't Work Well

1. **Dual code paths** - Having to maintain both module Build function and legacy `buildCutoutParams()` is confusing
2. **Schemagen limitations** - Type inference for defaults needs improvement
3. **Error messages** - "unsupported cutout shape" didn't immediately point to the right location

## Key Learnings

### Technical Learnings

1. **Dual code paths:** Cutouts use legacy `buildCutoutParams()` in `map.go`, not module Build function
2. **Preset emission:** Polygon presets come after depth/angle but before other flags
3. **Shape flag location:** Shape flag is at position 5 in the array (after fromBack, fromLeft, width, length, radius)
4. **Preset mapping:** DSL uses lowercase with underscores (hexagon, 6pt_star) while YAPP uses camelCase (shapeHexagon, shape6ptStar)

### Process Learnings

1. **Trace execution paths:** Always check which code is actually being executed, not just what exists
2. **Test incrementally:** Test after each step, not just at the end
3. **Check legacy code:** Module system doesn't always replace legacy code paths immediately

### Documentation Learnings

1. **Code path documentation needed** - Should document which modules use legacy code paths
2. **Error message improvements** - Error messages should point to actual execution location
3. **Migration guide** - Need guide for migrating features from legacy to module system

## Recommendations for Future Implementations

### Before Starting

1. ✅ Check which code path is actually used (module Build vs legacy functions)
2. ✅ Study YAPP examples for exact parameter order
3. ✅ Review existing similar implementations

### During Implementation

1. ✅ Update schema.yaml first
2. ✅ Generate code and verify types
3. ✅ Update both code paths if legacy code exists
4. ✅ Test after each change

### Integration Checklist

- [x] Schema created and validated
- [x] Code generated (`schemagen discover`)
- [x] Module Build function updated (if used)
- [x] Legacy code paths updated (if used)
- [x] End-to-end test passes
- [x] SCAD output verified

### After Implementation

1. ✅ Create working DSL example
2. ✅ Test end-to-end (resolve → generate → verify SCAD output)
3. ✅ Update ticket with docmgr (changelog, relate files, check tasks)
4. ✅ Document any dual code path issues

## Files Created/Modified

### Created
- `examples/yapp-demo-polygon-test.yaml` - Test YAML with polygon examples

### Modified
- `pkg/yappgen/modules/cutouts/schema.yaml` - Added polygon shape and preset field
- `pkg/yappgen/modules/cutouts/module.go` - Added polygonPresetFlag and preset emission
- `pkg/yappgen/modules/cutouts/schema_gen.go` - Auto-generated with Polygon field
- `pkg/yappgen/schema.go` - Added polygon to ShapeFlag function
- `pkg/yappgen/map.go` - Added polygon preset support to buildCutoutParams
- `pkg/yappgen/modules/lighttubes/schema_gen.go` - Fixed type mismatch (0.0 instead of 0)

## Verification

**Unit tests:** ✅ `go test ./pkg/yappgen/modules/cutouts` passes

**Schema validation:** ✅ `go run ./cmd/schemagen validate` passes

**Code generation:** ✅ `go run ./cmd/schemagen discover` succeeds

**End-to-end:** ✅ Resolve → Generate → Verify SCAD output works

**SCAD output matches YAPP format:**
```scad
cutoutsBase = [
  [15, 15, 25, 25, 0, yappPolygon, 0, 30, shapeHexagon]
];
cutoutsFront = [
  [5, 15, 15, 15, 0, yappPolygon, undef, undef, shape6ptStar]
];
cutoutsLid = [
  [10, 20, 20, 20, 0, yappPolygon, undef, undef, shapeArrow]
];
```

## Conclusion

The polygon cutout shape implementation was successful but revealed an important architectural issue: cutouts still use a legacy code path (`buildCutoutParams()` in `map.go`) rather than the module Build function. This required updating both code paths to ensure polygon presets are emitted correctly.

The main challenges were:
1. Finding the correct code execution path (legacy vs module system)
2. YAML schema syntax issues (easily fixed)
3. LightTubes schema type mismatch (known schemagen limitation)

The implementation demonstrates that the module system is working well for new features, but some legacy code paths still need to be maintained. The key to success was tracing through the actual execution path and updating both code locations.

**Time estimate:** 2-3 days (as estimated in ticket) was accurate. Actual time: ~1.5 hours for implementation, but could take longer for someone unfamiliar with the dual code path issue.

**Next steps:** Continue with Priority 2 features (cutout masks, labelsPlane, ridgeExt) following the same pattern, but always check which code path is actually used.

---

## Comparison Runs and Findings (2025-11-16)

### What I did
- Created comparison inputs and generated STL pairs for DSL vs legacy:
  - `examples/compare/cutouts_polygons.yaml` ⇄ `examples/YAPP_Compare_cutouts_polygons_v3.scad`
  - `examples/yapp-demo-lighttubes.yaml` ⇄ `examples/YAPP_Demo_lightTubes_v30.scad`
- Outputs written to:
  - `/tmp/yapp_compare/cutouts_polygons/{dsl-base.stl,dsl-lid.stl,legacy-base.stl,legacy-lid.stl}`
  - `/tmp/yapp_compare/lighttubes/{dsl-base.stl,dsl-lid.stl,legacy-base.stl,legacy-lid.stl}`
- Added all-faces cutouts coverage: `examples/compare/cutouts_all_faces.yaml` with outputs in `/tmp/yapp_compare/cutouts_all_faces/`

### Findings
- Polygon cutouts: DSL and legacy STL pairs match (hexagon base, 6pt star front, arrow lid).
- Lighttubes example: Side-face cutouts appear lower than legacy — root cause is axis naming:
  - Legacy front/back use `[posy, posz, ...]`. Our DSL maps `from_back → posy`, `from_left → posz`. In the current example, `from_left` is small (e.g., 5/2), so openings sit near the bottom. Setting `from_left` to the legacy `posz` value yields identical results.
- This is an ergonomics gap (not a functional one). We should add face-aware mapping or clearer keys for side faces.

### Follow-ups
- Added tasks:
  - Face-aware cutout mapping (pos axes clarity for side faces) and docs update
  - Align lighttubes example cutout heights (posz) with legacy and re-compare
