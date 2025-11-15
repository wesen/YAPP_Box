---
Title: MVP Path Forward — Semantic Fixes Required
Ticket: YAPP-ENCL-DSL-001
Status: active
Topics:
    - yapp
    - openscad
    - dsl
DocType: analysis
Intent: long-term
Owners:
    - manuel
RelatedFiles:
    - Path: debate/05-round-5-mvp-semantic-correctness-what-must-map-cleanly.md
      Note: Round 5 debate identifying critical fixes
    - Path: reference/01-enclosure-dsl-language-reference.md
      Note: DSL spec to be updated
ExternalSources: []
Summary: MVP implementation path with 2 critical semantic fixes and ~2 day effort estimate
LastUpdated: 2025-11-08T19:40:13.076998328-05:00
---


# MVP Path Forward — Semantic Fixes Required

## Executive Summary

Round 5 debate identified **three semantic issues** for MVP implementation (pcbStands, connectors, snapJoins, cutouts only). One has been **resolved by simplification**, two require careful implementation.

**Good news**: All issues are **solvable** and **not blockers**. They require careful implementation but are straightforward.

**MVP Scope** (acceptable):
- ✅ Use YAPP default coordinate systems (no DSL coordinate setting)
- ✅ Breaking changes OK as we iterate
- ✅ No maintenance/drift concerns
- ✅ Limited feature set (4 feature types)

**Critical Fixes** (required):
1. ~~Mixed coordinate system defaults~~ → **RESOLVED: Remove `coordinates.origin`, use YAPP defaults**
2. Optional positional parameters → Build parameter schemas, use `undef`
3. Shape-dependent parameters → Shape-specific mapping tables

---

## Critical Issue 1: Coordinate System Defaults (RESOLVED)

### The Problem

YAPP features have **different default coordinate systems**:

| Feature | YAPP Default | Semantic Reason |
|---------|--------------|-----------------|
| `pcbStands` | `yappCoordPCB` | Attached to PCB holes |
| `connectors` | `yappCoordPCB` | Attached to PCB |
| `cutouts` | `yappCoordBox` | Aligned to box edges (for external connectors) |
| `snapJoins` | `yappCoordBox` | Attached to box edges |

**Original DSL design**: Global `coordinates.origin: pcb|box|boxinside`

**Problem**: Creates confusion about which features use which coordinate system.

### The Solution

**Remove `coordinates.origin` entirely from DSL.** Use YAPP's default coordinate systems for each feature type:

```yaml
# User writes (NO coordinates section):
features:
  pcb_stands:
    - x: 10      # Interpreted as PCB coordinates (YAPP default)
      y: 10

  cutouts:
    - face: front
      x: 50      # Interpreted as Box coordinates (YAPP default)
      z: 10
      width: 15
      height: 10
```

```openscad
// Generator outputs (using YAPP defaults, no explicit flags needed):
pcbStands = [
  [10, 10]  // Uses yappCoordPCB (YAPP default for pcbStands)
];

cutoutsFront = [
  [50, 10, 15, 10, 0, yappRectangle]  // Uses yappCoordBox (YAPP default for cutouts)
];
```

### Feature-Specific Coordinate Defaults

| DSL Feature | YAPP Default Coordinate System | User Interprets Coordinates As |
|-------------|-------------------------------|-------------------------------|
| `pcb_stands` | `yappCoordPCB` | Relative to PCB origin (bottom-left of PCB) |
| `connectors` | `yappCoordPCB` | Relative to PCB origin |
| `cutouts` | `yappCoordBox` | Relative to box origin (bottom-left-back of box) |
| `snap_joins` | `yappCoordBox` | Relative to box origin |

### Implementation

1. **DSL spec update**: Remove `coordinates` section entirely
2. **Generator logic**: Rely on YAPP defaults, don't emit coordinate flags
   ```python
   # No coordinate flag generation needed - YAPP defaults are correct
   def generate_pcbstand(stand_dict):
       return [stand_dict['x'], stand_dict['y'], ...]  # No yappCoordPCB flag
   
   def generate_cutout(cutout_dict):
       return [cutout_dict['x'], cutout_dict['z'], ...]  # No yappCoordBox flag
   ```
3. **Simplification**: No need to track global coordinate setting or emit explicit flags

### Documentation Required

Add to DSL spec:

> **Coordinate System Semantics**
> 
> The DSL uses YAPP's default coordinate systems for each feature type. You do not specify a coordinate system in the DSL - it is determined by the feature type:
> 
> - **PCB-attached features** (`pcb_stands`, `connectors`): Coordinates are relative to the PCB origin (bottom-left corner of PCB, top surface). If you change `paddingLeft` or `paddingBack`, these features move with the PCB.
> 
> - **Box-attached features** (`cutouts`, `snap_joins`): Coordinates are relative to the box origin (bottom-left-back corner of box, outside surface). If you change padding, these features stay at the box edges.
> 
> This matches YAPP's semantic model and ensures correct behavior when box dimensions change.
> 
> **Example**: A PCB standoff at `x: 10, y: 10` is positioned 10mm from the left edge and 10mm from the back edge **of the PCB**. A cutout at `x: 50, z: 10` is positioned 50mm from the back edge and 10mm from the bottom **of the box**.

---

## Critical Issue 2: Optional Positional Parameters

### The Problem

YAPP uses positional arrays with many optional parameters:

```openscad
// pcbStands parameters:
// p(0) = posx           REQUIRED
// p(1) = posy           REQUIRED
// p(2) = height         OPTIONAL (default = standoffHeight)
// p(3) = pcbGap         OPTIONAL (default = -1)
// p(4) = diameter       OPTIONAL (default = standoffDiameter)
// p(5) = pinDiameter    OPTIONAL (default = standoffPinDiameter)
// p(6) = holeSlack      OPTIONAL (default = standoffHoleSlack)
// p(7) = filletRadius   OPTIONAL (default = 0)
// p(8) = pinLength      OPTIONAL (default = 0)
```

**If user wants to set `diameter` but use default `height`**:

```openscad
// Must use undef for skipped params:
[10, 10, undef, undef, 7]
//       ^^^^^  ^^^^^
//       height pcbGap (skipped, use defaults)
```

**DSL advantage**: Named fields, skip what you don't need:

```yaml
pcb_stands:
  - x: 10
    y: 10
    diameter: 7  # Skip height and pcbGap
```

**Generator challenge**: Must produce `[10, 10, undef, undef, 7]`, not `[10, 10, 7]` (which would make 7 the height).

### The Solution

**Build parameter schema tables** for each feature type:

```python
PCBSTANDS_PARAMS = [
    ('x', REQUIRED, None),
    ('y', REQUIRED, None),
    ('height', OPTIONAL, 'standoffHeight'),  # Use global var as default
    ('pcb_gap', OPTIONAL, -1),
    ('diameter', OPTIONAL, 'standoffDiameter'),
    ('pin_diameter', OPTIONAL, 'standoffPinDiameter'),
    ('hole_slack', OPTIONAL, 'standoffHoleSlack'),
    ('fillet_radius', OPTIONAL, 0),
    ('pin_length', OPTIONAL, 0),
]

CONNECTORS_PARAMS = [
    ('x', REQUIRED, None),
    ('y', REQUIRED, None),
    ('stand_height', REQUIRED, None),
    ('screw_diameter', REQUIRED, None),
    ('screw_head_diameter', REQUIRED, None),
    ('insert_diameter', REQUIRED, None),
    ('outside_diameter', REQUIRED, None),
    ('insert_depth', OPTIONAL, 'entire_connector'),
    ('pcb_gap', OPTIONAL, 'pcbThickness_or_0'),  # Depends on coord system
    ('fillet_radius', OPTIONAL, 0),
]

SNAPJOINS_PARAMS = [
    ('pos', REQUIRED, None),
    ('width', REQUIRED, None),
    # side is a flag, not positional param
]

CUTOUTS_PARAMS = [
    ('from_back', REQUIRED, None),
    ('from_left', REQUIRED, None),
    ('width', REQUIRED, None),
    ('length', REQUIRED, None),
    ('radius', REQUIRED, None),
    # shape is a flag, not positional param
    ('depth', OPTIONAL, 0),  # 0 = auto (plane thickness)
    ('angle', OPTIONAL, 0),
]
```

**Generator logic**:

```python
def generate_params(feature_dict, param_schema):
    params = []
    for (name, required, default) in param_schema:
        if name in feature_dict:
            params.append(feature_dict[name])
        elif required == REQUIRED:
            raise ValueError(f"Missing required parameter: {name}")
        else:
            params.append('undef')  # Use undef for skipped optional params
    return params
```

### Implementation

1. **Create parameter schema tables** for pcbStands, connectors, snapJoins, cutouts
2. **Generator uses schemas** to build positional arrays
3. **Validate required parameters** are present
4. **Insert `undef`** for skipped optional parameters

### Testing

Test that skipped parameters work correctly:

```yaml
pcb_stands:
  - x: 10
    y: 10
    diameter: 7  # Skip height, pcbGap
```

Should generate:

```openscad
[10, 10, undef, undef, 7, undef, undef, undef, undef]
```

Not:

```openscad
[10, 10, 7]  // WRONG - 7 becomes height
```

---

## Critical Issue 3: Shape-Dependent Parameters

### The Problem

Cutout parameters depend on shape type:

```openscad
//  yappRectangle       | width, length         | radius        |
//  yappCircle          | radius                | width, length |
//  yappRoundedRect     | width, length, radius |               |
```

For `yappRectangle`: width and length are used, radius is ignored (set to 0)
For `yappCircle`: radius is used, width and length are ignored (set to 0)

### The Solution

**Shape-specific parameter mapping**:

```python
CUTOUT_SHAPE_PARAMS = {
    'rectangle': {
        'width': REQUIRED,
        'length': REQUIRED,
        'radius': 0,  # Unused, set to 0
    },
    'circle': {
        'width': 0,  # Unused, set to 0
        'length': 0,  # Unused, set to 0
        'radius': REQUIRED,
    },
    'rounded_rect': {
        'width': REQUIRED,
        'length': REQUIRED,
        'radius': REQUIRED,
    },
    'circle_with_flats': {
        'width': REQUIRED,
        'radius': REQUIRED,
        'length': REQUIRED,  # Distance between flats
    },
}
```

**Generator logic**:

```python
def generate_cutout_params(cutout_dict):
    shape = cutout_dict['shape']
    shape_params = CUTOUT_SHAPE_PARAMS[shape]
    
    width = cutout_dict.get('width', shape_params['width'])
    length = cutout_dict.get('length', shape_params['length'])
    radius = cutout_dict.get('radius', shape_params['radius'])
    
    # Validate required params are present
    if width == REQUIRED and 'width' not in cutout_dict:
        raise ValueError(f"Shape {shape} requires width parameter")
    
    return [from_back, from_left, width, length, radius, shape_flag, ...]
```

### Implementation

1. **Create shape parameter mapping table**
2. **Generator validates** shape-specific required parameters
3. **Set unused parameters to 0** based on shape type

### Testing

Test that unused parameters are set to 0:

```yaml
cutouts:
  - face: front
    x: 10
    z: 5
    shape: circle
    radius: 7
    # width and length not specified
```

Should generate:

```openscad
[10, 5, 0, 0, 7, yappCircle, ...]
//      ^  ^
//      width=0, length=0 (unused for circle)
```

---

## Implementation Checklist

### Phase 1: Update DSL Spec

- [ ] Remove `coordinates` section from DSL spec
- [ ] Document YAPP default coordinate systems for each feature type
- [ ] Add clear explanation of PCB-relative vs Box-relative coordinates
- [ ] Define DSL field names for each feature type
- [ ] Define DSL enum values (placement, type, corner, shape, etc.)
- [ ] Add examples showing coordinate interpretation for each feature

### Phase 2: Build Parameter Schemas

- [ ] Create `pcbStands` parameter schema (9 positional params)
- [ ] Create `connectors` parameter schema (10 positional params)
- [ ] Create `snapJoins` parameter schema (2 positional params)
- [ ] Create `cutouts` parameter schema (7 positional params)
- [ ] Create cutout shape-specific parameter mapping

### Phase 3: Implement Generator

- [ ] Implement positional parameter generation (use schemas, insert `undef`)
- [ ] Implement shape-specific parameter handling (cutouts)
- [ ] Implement flag generation (placement, type, corner, shape, etc.)
- [ ] Add validation (required parameters present, valid enum values)
- [ ] Do NOT emit coordinate flags (rely on YAPP defaults)

### Phase 4: Test Against Real Examples

- [ ] Generate SCAD for simple box (pcbStands only)
- [ ] Generate SCAD for box with connectors
- [ ] Generate SCAD for box with cutouts (multiple shapes)
- [ ] Generate SCAD for box with snapJoins
- [ ] Validate OpenSCAD compilation (no syntax errors)
- [ ] Validate OpenSCAD rendering (visual inspection)

---

## Effort Estimate

| Task | Effort | Notes |
|------|--------|-------|
| Update DSL spec | 1-2 hours | Remove coordinates section, document YAPP defaults |
| Build parameter schemas | 4-6 hours | Tedious but straightforward |
| Implement generator | 6-10 hours | Core logic, validation, testing (simpler without coordinate flags) |
| Test against examples | 2-4 hours | Generate SCAD, validate rendering |
| **Total** | **13-22 hours** | ~2 days of focused work |

---

## Success Criteria

MVP is successful if:

1. ✅ DSL can express basic boxes with pcbStands, connectors, cutouts, snapJoins
2. ✅ Generated SCAD compiles in OpenSCAD without errors
3. ✅ Generated SCAD renders correctly (visual inspection)
4. ✅ Coordinate systems use YAPP defaults (no explicit flags in generated SCAD)
5. ✅ PCB-relative features (pcbStands, connectors) move with PCB when padding changes
6. ✅ Box-relative features (cutouts, snapJoins) stay at box edges when padding changes
7. ✅ Optional parameters work correctly (skipped params use `undef`)
8. ✅ Shape-specific parameters work correctly (unused params set to 0)

---

## What's NOT in MVP

Explicitly out of scope (can add later):

- ❌ Per-feature coordinate overrides (global is fine)
- ❌ lightTubes, pushButtons, boxMounts, labels, ridgeExt, displayMounts
- ❌ Advanced cutout shapes (polygon, ring, sphere)
- ❌ Advanced flags (yappNoFillet, yappSymmetric, yappCountersink, etc.)
- ❌ Multiple PCBs
- ❌ Hook functions (custom SCAD modules)
- ❌ Version compatibility checking
- ❌ SCAD-to-DSL reverse translation
- ❌ Escape hatch for raw SCAD

---

## Next Steps

1. **Update DSL spec** — Remove coordinates section, document YAPP defaults (1-2 hours)
2. **Build parameter schemas** for 4 feature types (4-6 hours)
3. **Implement generator** with semantic fixes (6-10 hours)
4. **Test with real examples** (2-4 hours)

**Total: ~2 days of focused work to MVP**

After MVP:
- Iterate based on real usage
- Add more features as needed
- Consider advanced features (per-feature coord overrides, escape hatch, etc.)

---

## Conclusion

The three semantic issues identified in Round 5 have been addressed:

1. **Mixed coordinate defaults** → **RESOLVED: Remove `coordinates.origin`, use YAPP defaults**
2. **Optional positional parameters** → Build parameter schemas, use `undef`
3. **Shape-dependent parameters** → Shape-specific mapping tables

**Key simplification**: By removing the `coordinates` section from the DSL and relying on YAPP's default coordinate systems, we eliminate the most complex semantic issue. Users learn one simple rule: "PCB features use PCB coordinates, box features use box coordinates."

With these fixes, MVP DSL-to-YAPP translation will be **semantically correct** and ready for real use.

**Recommendation**: Implement these fixes before generating any SCAD. The effort is ~2 days, and it ensures correctness from day one.
