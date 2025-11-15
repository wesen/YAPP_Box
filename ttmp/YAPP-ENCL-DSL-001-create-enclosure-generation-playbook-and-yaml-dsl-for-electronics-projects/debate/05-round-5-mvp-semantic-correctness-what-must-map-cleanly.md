---
Title: Round 5 — MVP Semantic Correctness: What Must Map Cleanly?
Ticket: YAPP-ENCL-DSL-001
Status: active
Topics:
    - yapp
    - openscad
    - dsl
DocType: debate
Intent: long-term
Owners:
    - manuel
RelatedFiles:
    - Path: analysis/03-debate-synthesis-dsl-to-yapp-translation-evaluation.md
      Note: Previous debate synthesis
    - Path: reference/01-enclosure-dsl-language-reference.md
      Note: DSL spec
ExternalSources: []
Summary: "Round 5 debate on MVP semantic issues with relaxed constraints"
LastUpdated: 2025-11-09
---

# Round 5 — MVP Semantic Correctness: What Must Map Cleanly?

**Context**: Previous rounds identified the DSL as incomplete and semantically incorrect. New constraints:
- **MVP scope**: pcbStands, connectors, snapJoins, cutouts (basic features only)
- **Single coordinate system acceptable**: No need for per-feature overrides
- **Maintenance/drift not a concern**: Breaking changes OK as we iterate
- **Critical requirement**: Things that don't map cleanly must be fixed NOW to avoid future trouble

**Question**: Given MVP scope and relaxed constraints, what semantic issues MUST be fixed from the start?

**Primary Candidates**: Dr. Sarah (Architect), YAPPgenerator, The Template

---

## Pre-Debate Research

### YAPP Default Coordinate Systems by Feature

```bash
grep -B2 "Default origin.*yappCoord" YAPP_Template_v3.scad
```

**Result**: YAPP features have **different default coordinate systems**:

| Feature | Default Coordinate System | Reason |
|---------|--------------------------|--------|
| `pcbStands` | **yappCoordPCB** | Standoffs positioned relative to PCB |
| `connectors` | **yappCoordPCB** | Screw connectors positioned relative to PCB |
| `lightTubes` | **yappCoordPCB** | Light tubes aligned to PCB LEDs |
| `cutouts*` | **yappCoordBox** | Cutouts positioned relative to box edges |
| `snapJoins` | **yappCoordBox** | Snap joints positioned relative to box edges |
| `boxMounts` | **yappCoordBox** | External mounts positioned relative to box |

**Critical finding**: Even in MVP scope, features have **two different default coordinate systems**:
- PCB-relative: pcbStands, connectors, lightTubes
- Box-relative: cutouts, snapJoins

### DSL Current Coordinate System Design

From `reference/01-enclosure-dsl-language-reference.md`:

```yaml
coordinates:
  origin: pcb|box|boxinside      # GLOBAL setting
```

**Problem**: DSL has one global coordinate system, but YAPP MVP features need two defaults.

### Parameter Structure Analysis

**pcbStands** (9 positional params + flags):
```openscad
[posx, posy, height, pcbGap, diameter, pinDiameter, holeSlack, filletRadius, pinLength, 
 yappBoth|yappLidOnly|yappBaseOnly, yappPin|yappHole|yappTopPin, 
 yappAllCorners|yappFrontLeft|..., yappCoordPCB|yappCoordBox|yappCoordBoxInside, ...]
```

**connectors** (10 positional params + flags):
```openscad
[posx, posy, standHeight, screwDiameter, screwHeadDiameter, insertDiameter, outsideDiameter,
 insertDepth, pcbGap, filletRadius,
 yappAllCorners|yappFrontLeft|..., yappCoordPCB|..., yappNoFillet, yappCountersink, ...]
```

**snapJoins** (3 positional params + flags):
```openscad
[pos, width, yappLeft|yappRight|yappFront|yappBack,
 yappOrigin|yappCenter, yappSymmetric, yappRectangle]
```

**Observation**: All three have **optional positional parameters with defaults** and **multi-valued flag parameters**.

---

## Opening Statements

### Dr. Sarah "The Architect" Martinez

Alright, let's focus on what MUST work for MVP. I've identified **three critical semantic issues** that will cause problems if not fixed now:

### Issue 1: Mixed Default Coordinate Systems (CRITICAL)

Even with single global coordinate system, you have a problem:

**YAPP defaults**:
- `pcbStands` default to `yappCoordPCB`
- `connectors` default to `yappCoordPCB`
- `cutouts` default to `yappCoordBox`
- `snapJoins` default to `yappCoordBox`

**DSL approach**: Set global `coordinates.origin: pcb`

**What happens**:
1. pcbStands and connectors generate correctly (match YAPP default)
2. Cutouts and snapJoins generate with `yappCoordPCB` flag (WRONG — they default to `yappCoordBox`)

**Result**: Cutouts and snap joints positioned incorrectly.

**Example**:
```yaml
coordinates:
  origin: pcb  # Global setting

features:
  cutouts:
    - face: front
      x: 10  # User thinks: 10mm from PCB left edge
      z: 5
      # DSL generates: [..., yappCoordPCB]
      # YAPP interprets: 10mm from PCB (correct)
      # But YAPP default is yappCoordBox!
```

If user doesn't set `coordinates.origin`, what's the default? If it's `pcb`, cutouts are wrong. If it's `box`, standoffs are wrong. **You can't win with a single global setting.**

**Fix required**: DSL must either:
1. **Option A**: Have separate coordinate settings per feature type
   ```yaml
   coordinates:
     pcb_features: pcb      # for pcbStands, connectors
     box_features: box      # for cutouts, snapJoins
   ```
2. **Option B**: Always generate explicit coordinate flags (override YAPP defaults)
   ```openscad
   pcbStands = [[10, 10, yappCoordPCB]];  // Explicit even if default
   cutouts = [[10, 5, 15, 10, 0, yappRectangle, yappCoordBox]];  // Explicit
   ```

Option B is safer — no ambiguity, always explicit.

### Issue 2: Optional Positional Parameters with Defaults (MEDIUM)

Look at pcbStands parameters:
```
p(0) = posx           REQUIRED
p(1) = posy           REQUIRED
p(2) = height         OPTIONAL (default = standoffHeight)
p(3) = pcbGap         OPTIONAL (default = -1, which means pcbThickness)
p(4) = diameter       OPTIONAL (default = standoffDiameter)
...
```

**Problem**: If user wants to set `diameter` but use default `height`, they must write:
```openscad
[10, 10, undef, undef, 7]  // Use undef for optional params
```

**DSL approach**: Named fields
```yaml
pcb_stands:
  - x: 10
    y: 10
    diameter: 7  # Skip height, use default
```

**Mapping challenge**: How do you generate the SCAD array?
- If you omit parameters: `[10, 10, 7]` — WRONG (7 becomes height, not diameter)
- If you use `undef`: `[10, 10, undef, undef, 7]` — CORRECT but verbose

**Fix required**: DSL generator must:
1. Know parameter order for each feature type
2. Insert `undef` for skipped optional parameters
3. Maintain correct positional mapping

This is tedious but solvable. Not a blocker, just implementation work.

### Issue 3: Multi-Valued Flag Parameters (LOW)

Some flags are multi-valued:
```openscad
pcbStands = [
  [10, 10, yappBoth, yappPin, yappFrontLeft]
  //       ^^^^^^^^  ^^^^^^^  ^^^^^^^^^^^^^
  //       placement  type     corner (optional)
];
```

**DSL approach**:
```yaml
pcb_stands:
  - x: 10
    y: 10
    placement: both      # yappBoth
    type: pin            # yappPin
    corner: front_left   # yappFrontLeft (optional)
```

**Mapping**: Straightforward enum translation. Not a semantic issue, just naming.

**Fix required**: Define DSL enum values that map cleanly to YAPP flags. Document the mapping.

---

### `YAPPgenerator_v3.scad` — "The Generator"

Dr. Martinez is right about Issue 1 (mixed defaults). Let me explain why this matters from YAPP's perspective.

**Design rationale for different defaults**:

- **pcbStands default to yappCoordPCB** because standoffs are physically attached to PCB holes. If you change `paddingLeft`, standoffs should move with the PCB.

- **cutouts default to yappCoordBox** because cutouts are for external connectors (USB, power jack, etc.) that mount to the box, not the PCB. If you change padding, cutouts should stay at box edges.

This isn't arbitrary — it's the **correct semantic behavior** for each feature type.

**If DSL uses global coordinate system**, you're forcing all features into one semantic model. That's wrong.

**My recommendation**: Option B from Dr. Martinez — always generate explicit coordinate flags. This way:
1. DSL can have global `coordinates.origin` as a **default**
2. Generator always outputs explicit `yappCoordPCB` or `yappCoordBox` based on feature type
3. No ambiguity, no reliance on YAPP defaults

**Example**:
```yaml
coordinates:
  origin: pcb  # User's default preference

features:
  pcb_stands:
    - x: 10
      y: 10
      # Generator outputs: [10, 10, yappCoordPCB]  (respects feature semantics)
  
  cutouts:
    - face: front
      x: 50
      z: 10
      width: 15
      height: 10
      # Generator outputs: [50, 10, 15, 10, 0, yappRectangle, yappCoordBox]
      # (uses yappCoordBox despite global origin: pcb)
```

This respects YAPP's semantic model while giving users a global default for convenience.

---

### `YAPP_Template_v3.scad` — "The Template"

*[Checks parameter order documentation]*

I want to emphasize Issue 2 (optional positional parameters). This is a **real pain point** for SCAD users.

**Current SCAD experience**:
```openscad
// I want a standoff at (10, 10) with 7mm diameter, but default height
pcbStands = [
  [10, 10, undef, undef, 7]  // Must count positions, use undef
];
```

**Better DSL experience**:
```yaml
pcb_stands:
  - x: 10
    y: 10
    diameter: 7  # Just specify what I want
```

This is a **huge usability win** for the DSL. But the generator must handle the mapping correctly.

**Implementation requirement**: Generator needs a **parameter schema** for each feature:

```python
PCBSTANDS_SCHEMA = [
  ("posx", REQUIRED),
  ("posy", REQUIRED),
  ("height", OPTIONAL, "standoffHeight"),
  ("pcbGap", OPTIONAL, -1),
  ("diameter", OPTIONAL, "standoffDiameter"),
  ("pinDiameter", OPTIONAL, "standoffPinDiameter"),
  ("holeSlack", OPTIONAL, "standoffHoleSlack"),
  ("filletRadius", OPTIONAL, 0),
  ("pinLength", OPTIONAL, 0),
]
```

Then when generating:
```python
def generate_pcbstand(stand_dict):
    params = []
    for (name, required, default) in PCBSTANDS_SCHEMA:
        if name in stand_dict:
            params.append(stand_dict[name])
        elif required:
            raise Error(f"Missing required parameter: {name}")
        else:
            params.append("undef")  # Use undef for skipped optional params
    return params
```

**This is tedious but essential.** If you get parameter order wrong, YAPP will silently misinterpret values (e.g., diameter becomes height).

---

## Rebuttals

### Dr. Sarah "The Architect" Martinez

*[Responds to Generator's Option B proposal]*

Yes, I agree — **always generate explicit coordinate flags** is the right approach. Let me formalize this:

**DSL Semantic Rule**: The global `coordinates.origin` is a **user preference hint**, not a direct mapping. The generator uses feature-specific semantic rules:

| Feature Type | Generated Coordinate Flag | Rationale |
|--------------|--------------------------|-----------|
| `pcb_stands` | Always `yappCoordPCB` | Attached to PCB |
| `connectors` | Always `yappCoordPCB` | Attached to PCB |
| `light_tubes` | Always `yappCoordPCB` | Aligned to PCB LEDs |
| `cutouts` | Use `coordinates.origin` | User-specified |
| `snap_joins` | Always `yappCoordBox` | Attached to box edges |

This way:
- PCB-attached features always use PCB coordinates (semantically correct)
- Box-attached features always use box coordinates (semantically correct)
- Cutouts use user preference (flexible, since they could be either)

**This must be documented clearly** in the DSL spec, because it's not obvious that `coordinates.origin` doesn't apply to all features.

### `YAPPgenerator_v3.scad` — "The Generator"

*[Nods in agreement]*

That's a good compromise. Let me add one more requirement:

**Issue 4: Cutout Shape Parameters (NEW)**

Cutouts have shape-dependent parameters:

```openscad
//  yappRectangle       | width, length         | radius        |
//  yappCircle          | radius                | width, length |
//  yappRoundedRect     | width, length, radius |               |
```

For `yappRectangle`, you provide `width` and `length`. For `yappCircle`, you provide `radius` (and `width`/`length` are ignored).

**DSL must handle this**:
```yaml
cutouts:
  - face: front
    x: 10
    z: 5
    shape: rectangle
    width: 15
    height: 10
    # Generates: [10, 5, 15, 10, 0, yappRectangle]
  
  - face: front
    x: 30
    z: 5
    shape: circle
    radius: 7
    # Generates: [30, 5, 0, 0, 7, yappCircle]
    # (width=0, length=0 because unused)
```

**Fix required**: DSL generator must know which parameters are used for each shape type and set unused params to 0.

This is similar to Issue 2 (parameter mapping) but shape-specific.

---

## Moderator Summary

### Critical Semantic Issues (MUST FIX NOW)

#### 1. Mixed Default Coordinate Systems ⚠️ CRITICAL

**Problem**: YAPP features have different default coordinate systems (pcbStands→PCB, cutouts→Box). DSL global `coordinates.origin` can't match both.

**Solution**: Always generate explicit coordinate flags based on feature semantics:
- `pcb_stands`, `connectors`, `light_tubes` → Always `yappCoordPCB`
- `snap_joins` → Always `yappCoordBox`
- `cutouts` → Use `coordinates.origin` (user choice)

**Implementation**: Generator must override YAPP defaults explicitly.

**Documentation**: DSL spec must explain that `coordinates.origin` doesn't apply to all features.

#### 2. Optional Positional Parameters ⚠️ MEDIUM

**Problem**: YAPP uses positional arrays with optional params. Skipping a param requires `undef` placeholder.

**Solution**: Generator must:
1. Maintain parameter schema for each feature type (order, required/optional, defaults)
2. Insert `undef` for skipped optional parameters
3. Validate required parameters are present

**Implementation**: Tedious but straightforward. Build parameter schema tables.

#### 3. Shape-Dependent Parameters ⚠️ MEDIUM

**Problem**: Cutout parameters depend on shape (rectangle uses width/length, circle uses radius).

**Solution**: Generator must know which parameters are used for each shape and set unused params to 0.

**Implementation**: Shape-specific parameter mapping tables.

### Non-Critical Issues (Can Defer)

#### 4. Multi-Valued Flag Parameters ✓ LOW

**Problem**: YAPP flags like `yappBoth, yappPin, yappFrontLeft` are positional.

**Solution**: DSL uses named enums (`placement: both`, `type: pin`, `corner: front_left`).

**Implementation**: Simple enum translation. Not a semantic issue.

### Consensus Points

- **All candidates agree**: Issue 1 (mixed coordinate defaults) is critical and must be fixed
- **All candidates agree**: Always generating explicit coordinate flags is the right solution
- **All candidates agree**: Issue 2 (optional params) is tedious but solvable
- **Generator and Sarah agree**: Feature-specific semantic rules (pcbStands always PCB-relative) are correct

### Implementation Checklist

To fix semantic issues for MVP:

- [ ] **Define feature-specific coordinate rules** in DSL spec
  - pcb_stands, connectors, light_tubes → always yappCoordPCB
  - snap_joins → always yappCoordBox
  - cutouts → use coordinates.origin

- [ ] **Build parameter schema tables** for each feature type
  - Parameter order
  - Required vs optional
  - Default values
  - Shape-specific variations (for cutouts)

- [ ] **Generator: Always emit explicit coordinate flags**
  - Don't rely on YAPP defaults
  - Override with feature-specific rules

- [ ] **Generator: Handle optional parameters correctly**
  - Insert `undef` for skipped params
  - Maintain correct positional order

- [ ] **Generator: Handle shape-specific parameters**
  - Set unused params to 0 based on shape type

- [ ] **Document semantic rules clearly**
  - Explain why `coordinates.origin` doesn't apply to all features
  - Show examples of generated SCAD for each feature type

### Open Questions

1. Should `light_tubes` also use `coordinates.origin`, or always `yappCoordPCB`? (They're PCB-aligned, so probably always PCB)

2. For cutouts, should DSL allow per-cutout coordinate override, or only global? (Global is simpler for MVP)

3. Should DSL validate that required parameters are present, or let YAPP error? (DSL validation is better UX)

### Verdict

**Three semantic issues must be fixed for MVP**:
1. **Mixed coordinate defaults** (critical) → Always generate explicit flags
2. **Optional positional parameters** (medium) → Build parameter schemas, use `undef`
3. **Shape-dependent parameters** (medium) → Shape-specific mapping tables

These are **solvable** and **not blockers**, but must be implemented correctly from the start to avoid generating incorrect SCAD.

**Recommendation**: Implement these fixes before generating any SCAD. The tedious part is building parameter schema tables, but it's one-time work that ensures correctness.

**Next steps**: 
1. Update DSL spec with feature-specific coordinate rules
2. Build parameter schema tables (pcbStands, connectors, snapJoins, cutouts)
3. Implement generator with explicit coordinate flags and proper parameter mapping
4. Test against real YAPP examples to validate correctness
