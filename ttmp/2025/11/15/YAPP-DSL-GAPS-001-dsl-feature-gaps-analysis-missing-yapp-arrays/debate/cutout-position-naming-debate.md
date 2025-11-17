# Debate: Cutout Position Field Naming

**Topic:** Should we rename `from_back`/`from_left` to `from_face_left`/`from_face_bottom`?

**Date:** 2025-11-16

**Moderator:** Technical Documentation Lead

## Debate Participants

**Dr. Sarah Chen "The API Designer"**
- Senior API/DSL designer, 12 years experience
- Philosophy: "Clarity over cleverness—APIs should be self-documenting"
- Tools: API ergonomics analysis, naming patterns in successful DSLs

**Professor Alex Kim "The Coordinate Systems Expert"**
- Computer graphics/CAD background, coordinate system design specialist  
- Philosophy: "Good coordinate systems feel natural in 3D space"
- Tools: Compares with OpenSCAD, FreeCAD, Fusion360, other 3D tools

**`pkg/yappgen/modules/cutouts/module.go` "The Mapper"**
- Personified code entity (169 lines, handles 6 faces)
- Perspective: "I translate your names to YAPP coordinates—rename me carefully"
- Tools: Shows actual mapping code, traces coordinate transforms

**Jamie "The YAML User"**
- Embedded systems engineer, new to YAPP, learning from examples
- Philosophy: "I just want my USB port hole in the right place"
- Tools: Reads examples, tries to write YAML, hits errors

## Pre-Debate Research

### Current Field Usage Analysis

**The Mapper ran grep analysis on examples:**

```bash
grep -r "from_back\|from_left" examples/ --include="*.yaml" | wc -l
# Result: 60 matches across 9 YAML files
```

**Current cutouts schema fields:**
- `from_back`: required number (Distance along face's horizontal axis)
- `from_left`: required number (Distance along face's vertical/secondary axis)
- `pos_z`: optional number (Override for side faces, vertical from bottom)

**Usage breakdown:**
- `from_back:` appears 28 times across examples
- `from_left:` appears 32 times across examples  
- `pos_z:` appears 6 times (recently added feature)

### Other Module Naming Patterns

**The Mapper checked other feature modules:**

```bash
grep "^  x:\|^  y:" pkg/yappgen/modules/*/schema.yaml
```

**Results:**
- `pcbstands`: uses `x`, `y` (PCB coordinates)
- `connectors`: uses `x`, `y` (PCB coordinates)
- `lighttubes`: uses `x`, `y` (PCB coordinates)
- `pushbuttons`: uses `x`, `y` (PCB coordinates)
- `cutouts`: uses `from_back`, `from_left` (face-relative coordinates)

**Cutouts is the ONLY module using face-relative naming.**

### Current Mapping Logic

**The Mapper's internal translation code:**

```go
// pkg/yappgen/modules/cutouts/module.go lines 67-78
pos0 := item.FromBack   // Maps to posy for front/back, posx for left/right/base/lid
pos1 := item.FromLeft   // Maps to posz for sides, posy for base/lid

// Side face override
faceLower := strings.ToLower(strings.TrimSpace(item.Face))
isSideFace := faceLower == "front" || faceLower == "back" || faceLower == "left" || faceLower == "right"
if isSideFace && item.PosZ != nil {
    pos1 = *item.PosZ  // Override from_left with explicit vertical position
}
```

**Face-specific mappings:**
- Base/Lid: `from_back`→posx, `from_left`→posy
- Front/Back: `from_back`→posy, `from_left`→posz
- Left/Right: `from_back`→posx, `from_left`→posz

---

## Round 1: Should we rename from_back/from_left to from_face_left/from_face_bottom?

**QUESTION:** Should we replace `from_back`/`from_left` with `from_face_left`/`from_face_bottom` (removing `pos_z` as redundant)?

### Opening Statements

#### Dr. Sarah Chen (API Designer)

*[Pulls up ergonomics checklist]*

YES. Absolutely rename these fields. The current naming is a cognitive trap.

Look at the actual user pain: When Jamie writes `from_left: 5` for a front-face cutout, they're thinking "5mm from the left edge of the box"—but they're actually specifying a **vertical height**. That's not "from left," that's "from bottom"!

The names lie about what they mean:
- `from_back` sometimes means horizontal-along-face, sometimes means X-from-back-edge
- `from_left` sometimes means height, sometimes means Y-from-left-edge

Three naming principles violated:
1. **Name should match semantics** — "from_left" doesn't mean "height"
2. **Consistent meaning** — same field name shouldn't change meaning by context
3. **No hidden pivots** — users shouldn't need to memorize face-specific mappings

My proposal:
- `from_face_left`: horizontal position along the face (always means "left edge of THIS face")
- `from_face_bottom`: vertical position on the face (always means "bottom edge of THIS face")
- Delete `pos_z` — it's just an alias for `from_face_bottom`

This makes the fields **self-documenting**. A user reading the YAML understands immediately: "Ah, this cutout is 10mm from the face's left edge and 5mm from the face's bottom."

No mental model required. No mapping table. No confusion.

#### Professor Alex Kim (Coordinate Systems Expert)

*[Opens OpenSCAD and Fusion360 side-by-side]*

I'm torn. Let me think through the spatial reasoning.

**The case FOR face-relative naming:**

When you're standing in front of the box thinking "I need a USB port here," you're reasoning in **face-local coordinates**. Your brain doesn't think "posy=15, posz=8"—you think "centered horizontally, mid-height." Face-local names (`from_face_left`, `from_face_bottom`) match that mental model perfectly.

Industry precedent: CAD tools use face-local coordinates for feature placement. In Fusion360, when you sketch on a face, you get U/V coordinates relative to that face's local frame. OpenSCAD's `projection()` works similarly.

**The case AGAINST:**

But here's the problem: **Every other YAPP feature uses global PCB coordinates** (`x`, `y`). Cutouts would be the outlier. That creates an inconsistency:

```yaml
features:
  pcb_stands:
    - x: 5    # PCB coordinate
      y: 10   # PCB coordinate
  
  cutouts:
    - from_face_left: 5   # Face-local coordinate
      from_face_bottom: 10  # Face-local coordinate
```

One uses global frame, one uses local frame. Users have to context-switch mentally.

**My position:** If we're doing face-local naming, we need to justify why cutouts are special. Is it because cutouts span 6 different faces (base/lid/front/back/left/right) while other features live on the PCB plane? That's a legitimate difference—but it needs to be **clearly documented**.

I lean toward renaming IF we can articulate why cutouts deserve face-local coordinates while everything else stays global.

#### The Mapper (module.go)

*[Displays internal mapping code]*

Let me show you what I actually DO with these names:

```go
pos0 := item.FromBack
pos1 := item.FromLeft

// For front/back faces:
//   pos0 becomes posy (horizontal along wall)
//   pos1 becomes posz (vertical from base)
//
// For left/right faces:
//   pos0 becomes posx (horizontal along wall)  
//   pos1 becomes posz (vertical from base)
//
// For base/lid:
//   pos0 becomes posx (from back edge)
//   pos1 becomes posy (from left edge)
```

Current field names: `from_back`, `from_left`  
Proposed names: `from_face_left`, `from_face_bottom`

**My translation table changes:**

| Face | Old pos0 name | Old pos1 name | New pos0 name | New pos1 name |
|------|---------------|---------------|---------------|---------------|
| Front | from_back | from_left | from_face_left | from_face_bottom |
| Back | from_back | from_left | from_face_left | from_face_bottom |
| Left | from_back | from_left | from_face_left | from_face_bottom |
| Right | from_back | from_left | from_face_left | from_face_bottom |
| Base | from_back | from_left | from_face_left | from_face_bottom |
| Lid | from_back | from_left | from_face_left | from_face_bottom |

**Wait. The proposed names are WRONG for base/lid.**

On the base face:
- `from_face_left` means "from left edge" — that's posy (correct)
- `from_face_bottom` means "from bottom edge" — but the "bottom" of the base IS the back edge, so that's posx

But my current mapping for base is:
- `from_back` → posx (distance from back edge of box)
- `from_left` → posy (distance from left edge of box)

The proposed `from_face_bottom` doesn't match. The "bottom" of the base face (when viewed from above) is the **back edge**, not the front edge. So `from_face_bottom` would map to posx (from back), which IS what `from_back` currently does.

Actually... *[recalculates]*... if "bottom" means "back edge" for base, then:
- `from_face_left` → posy (left edge, matches current `from_left`)  
- `from_face_bottom` → posx (back edge, matches current `from_back`)

**But that's the OPPOSITE field order from the current schema!**

Current: `from_back` (pos0), `from_left` (pos1)  
Proposed: `from_face_left` (pos1?), `from_face_bottom` (pos0?)

The field order would flip. That's... confusing. Or we'd need to reorder the YAPP array positions, which breaks everything.

**Conclusion:** The proposed names work great for vertical faces but create ambiguity for horizontal faces (base/lid). And they imply a field-order reversal that would break my internal mapping.

I need clarification: Does "bottom" for base mean "back edge" or "front edge"? And are we keeping array position order (pos0, pos1) or flipping it?

#### Jamie (YAML User)

*[Opens examples/yapp-demo-lighttubes.yaml]*

Okay, I'm confused. Let me show you my actual experience.

I'm adding a USB-C port to the back face. I know it should be:
- 20mm from the left edge of the back wall
- 5mm up from the bottom

So I write:

```yaml
cutouts:
  - face: back
    from_back: ???   # Is this the 20mm or the 5mm?
    from_left: ???   # Is this the 20mm or the 5mm?
```

I **guess** that `from_left` means "from the left edge" so that's 20mm. And `from_back` must be... the 5mm vertical? But that doesn't make sense—why is vertical position called "from_back"?

I look at the examples. I see:

```yaml
- face: front
  from_back: 5
  from_left: 2
```

And a comment: `# from_back = posy (horizontal along front face)`

Wait, WHAT? `from_back` is horizontal? But it's called "from back"! How is horizontal "from back"?

*[Reads documentation section on faces and axes]*

Oh. OH. The "back" in `from_back` doesn't mean "the back face." It means "the back edge of the box, globally." So for a front-face cutout, `from_back` is actually measuring horizontally along the front wall, counting from the box's back edge.

That's... really unintuitive.

**My vote:** YES, rename to `from_face_left` and `from_face_bottom`. Those names tell me immediately what they mean **relative to the face I'm editing**. No global mental model required.

But I agree with The Mapper's concern: what does "bottom" mean for the base face? When I'm looking down at the base, which edge is "bottom"? 

Maybe we need `from_face_left` and `from_face_right` for one axis, and `from_face_top` and `from_face_bottom` for the other? Or is that overkill?

---

### Rebuttals (Round 1)

#### Dr. Sarah Chen (API Designer)

*[Responds to The Mapper's field-order concern]*

The Mapper raises a valid implementation concern: field order matters for array emission. But that's an **implementation detail**, not a user-facing API decision.

Here's the fix: We choose field names based on **user mental model**, then The Mapper adapts internally. If the natural semantic order for face-local coordinates is (horizontal, vertical) = (`from_face_left`, `from_face_bottom`), then that's what users should write.

If YAPP's internal arrays expect (vertical, horizontal) for some faces, The Mapper swaps them during translation. The user never sees that complexity.

Regarding base/lid ambiguity: You're right that "bottom" is ambiguous when the face IS the bottom. But we can clarify with documentation:
- For vertical faces: `from_face_bottom` = "from bottom edge" (vertical position from floor)
- For horizontal faces: `from_face_bottom` = "from back edge" (the edge toward the back of the box)

Or we use different names for horizontal faces: `from_face_back` and `from_face_left` for base/lid. Then it's explicit.

**Key point:** We're naming fields for human comprehension, not for implementation convenience.

#### Professor Alex Kim (Coordinate Systems Expert)

*[Draws coordinate frame diagrams]*

I'm now leaning AGAINST the rename. Here's why:

The proposed names (`from_face_left`, `from_face_bottom`) assume a **face-local coordinate frame** where each face has its own left/right/top/bottom. That's elegant for vertical faces but breaks down for horizontal faces, as The Mapper noted.

But there's a deeper issue: **Face-local frames are arbitrary**. When you look at the "front" face, is "left" your left (facing the box) or the box's left (from inside looking out)? Different CAD tools make different choices.

OpenSCAD uses **global coordinates** everywhere. YAPP inherits that convention. The current `from_back`/`from_left` naming, while confusing, actually reflects the YAPP/OpenSCAD reality: positions are always measured from fixed global reference edges (back edge, left edge of the entire box).

**Alternative proposal:** Keep the current names BUT add aliases:
- `from_back` also accepts `horizontal` or `x_pos` for clarity
- `from_left` also accepts `vertical` or `z_pos` for side faces, `y_pos` for base/lid

This preserves existing YAML files while offering clearer alternatives. No field-order flip required.

#### The Mapper (module.go)

*[Runs test suite]*

I appreciate Sarah's point about implementation vs. API, but there's a practical concern: **Every example file breaks** if we rename.

Current examples:
- 28 uses of `from_back:`
- 32 uses of `from_left:`
- 6 uses of `pos_z:` (which would be deleted)

If we rename, we'd need to:
1. Update all 9 example YAML files
2. Regenerate all comparison STLs
3. Update 3 documentation files
4. Adjust error messages in validation code

And then there's the semantic question: For the base face, does `from_face_bottom: 10` mean:
- 10mm from the back edge (toward the front), OR
- 10mm from the front edge (toward the back)?

If users interpret "bottom" as "the front edge" (the user-facing edge), that's the OPPOSITE of current `from_back` behavior.

**My position:** If we rename, we need:
1. Unambiguous definitions for all 6 faces
2. Clear migration path (or keep old names as aliases)
3. Updated validation with helpful error messages ("Did you mean from_face_left?")

Otherwise, we're trading one confusion (global reference edges) for another (ambiguous face-local semantics).

#### Jamie (YAML User)

*[Tries to write a cutout for the lid]*

Okay, let me test this. I want a ventilation slot on the lid:
- 30mm from the back edge of the lid
- 20mm from the left edge of the lid

Using current names:

```yaml
- face: lid
  from_back: 30   # I think this is X? From back edge?
  from_left: 20   # This is Y? From left edge?
```

Using proposed names:

```yaml
- face: lid
  from_face_bottom: 30  # Wait, "bottom" of the lid? That's... the back edge?
  from_face_left: 20    # This makes sense
```

**I'm more confused with the new names for horizontal faces.**

Maybe the real issue is that we're trying to use the same field names for both vertical faces (where "left/bottom" is intuitive) and horizontal faces (where "left/bottom" is ambiguous).

What if we have **different fields per face type?**
- Vertical faces: `along_face`, `height_from_base`
- Horizontal faces: `from_back_edge`, `from_left_edge`

Then there's no ambiguity. The names explicitly match the face type.

---

## Round 2: Should All Features Use Consistent Coordinate Naming?

**QUESTION:** Should cutouts use the same coordinate system as pcb_stands, light_tubes, connectors, etc. (all currently use `x`, `y`)?

### Opening Statements

#### Professor Alex Kim (Coordinate Systems Expert)

*[Creates comparison table]*

Let me show the inconsistency:

**Current naming across modules:**

| Module | Field 1 | Field 2 | Reference Frame |
|--------|---------|---------|-----------------|
| pcb_stands | `x` | `y` | PCB global (posx, posy) |
| connectors | `x` | `y` | PCB global (posx, posy) |
| light_tubes | `x` | `y` | PCB global (posx, posy) |
| push_buttons | `x` | `y` | PCB global (posx, posy) |
| cutouts | `from_back` | `from_left` | **Face-relative** |

**Cutouts are the outlier.** Every other feature says "give me x/y" (PCB coordinates). Only cutouts say "give me face-relative positions."

This creates two mental models:
1. "Features on the PCB" → think in x/y
2. "Features on enclosure faces" → think in from_back/from_left

**My argument:** We should unify. Either:

**Option A: Make cutouts use x/y (for side faces, add a z coordinate)**
```yaml
cutouts:
  - face: front
    y: 20   # Horizontal position along front wall
    z: 5    # Vertical position from base
```

**Option B: Make ALL features face-relative (breaking change for everything)**
```yaml
pcb_stands:
  - face: base    # Standoffs are on the base
    along_face: 10
    across_face: 15
```

Option A is less disruptive. Option B is more consistent but requires rewriting every module.

**My recommendation:** Convert cutouts to use x/y/z coordinates matching the other modules. Side faces would use y/z, base/lid use x/y. Consistent with the rest of YAPP.

#### The Mapper (module.go)

*[Shows module interfaces]*

Let me explain why cutouts are different:

**Other modules (pcb_stands, light_tubes, etc.):**
- Live on a **single plane** (the PCB)
- Position is always relative to PCB origin
- Two coordinates (x, y) are sufficient

**Cutouts:**
- Live on **six different planes** (front/back/left/right/base/lid)
- Each face has its own orientation in 3D space
- Two coordinates describe position **on that face**

The reason cutouts use face-relative naming is because **you can't describe a front-face cutout with just PCB x/y**. A cutout at "front face, 10mm from left, 5mm from bottom" is fundamentally different from a cutout at "PCB x=10, y=5."

If we switched cutouts to x/y/z:
- Base/lid cutouts: x, y (like PCB features)
- Front/back cutouts: y, z (horizontal along wall, vertical)
- Left/right cutouts: x, z (horizontal along wall, vertical)

That's actually closer to the YAPP internal model (posx/posy/posz), but now the **field names change meaning by face**:
- For base: x=forward/back, y=left/right
- For front: y=left/right, z=up/down

Users would need to memorize which coordinates are active for each face. That's arguably MORE confusing than face-relative names.

**My position:** Cutouts are a unique case. They span multiple faces with different orientations, so face-relative naming makes sense. The inconsistency with other modules is justified by the different problem domain.

#### Dr. Sarah Chen (API Designer)

*[Reviews consistency principle]*

Consistency is valuable, but **consistency for its own sake** creates bad APIs.

Alex is right that cutouts are the outlier—but that's because cutouts are DIFFERENT. They're not positioned on a single plane like PCB features; they're positioned on six different faces with six different orientations.

Consider the user workflow:

**For a PCB standoff:**
1. Look at PCB
2. Find mounting hole
3. Read x/y from PCB design file
4. Write `x: 10, y: 15` in YAML

**For a cutout:**
1. Look at enclosure face (e.g., back wall)
2. Measure where connector sits
3. Think: "20mm from left edge, 8mm from bottom"
4. Write... `from_back: 20, from_left: 8`? No wait, that's wrong...

The mental model is different. PCB features are placed from a 2D CAD file. Cutouts are placed by looking at a 3D object.

**Forcing x/y/z on cutouts** makes users translate face-local measurements into global coordinates. That's cognitive overhead.

**My recommendation:** Keep cutouts face-relative, but improve the naming to be self-explanatory (e.g., `along_face`, `across_face` or similar). Inconsistency is fine when the problem domains differ.

#### Jamie (YAML User)

*[Tries both approaches]*

Let me test this with a real example: USB-C port on the back wall.

**With x/y/z approach:**

```yaml
cutouts:
  - face: back
    y: 20   # 20mm from left edge
    z: 5    # 5mm from bottom
```

I have to remember that for back face, y=horizontal and z=vertical. But for base face, y is... the other horizontal? And z doesn't exist?

**With face-relative approach:**

```yaml
cutouts:
  - face: back
    horizontal: 20  # 20mm from left
    vertical: 5     # 5mm from bottom
```

This is clearer. I don't care about x/y/z—I'm measuring on the actual face.

**My vote:** Keep cutouts face-relative, but use clearer names than `from_back`/`from_left`. Something like:
- `horizontal` or `left_offset`
- `vertical` or `bottom_offset`

Let other modules keep x/y since they're on a single plane. Cutouts are different enough to justify different naming.

---

### Rebuttals (Round 2)

#### Professor Alex Kim (Coordinate Systems Expert)

*[Considers the arguments]*

I'm convinced. Jamie and Sarah are right: **cutouts are fundamentally different** from PCB features.

A PCB standoff lives in a 2D plane. A cutout lives on a 3D surface. Forcing the same coordinate system onto both creates a leaky abstraction.

If we used x/y/z for cutouts, users would need a lookup table:
- Base: use x, y
- Lid: use x, y  
- Front: use y, z
- Back: use y, z
- Left: use x, z
- Right: use x, z

That's absurd. It forces users to understand the internal YAPP coordinate system.

**New position:** Cutouts should use face-relative naming. But we need clearer field names that work for all 6 faces without ambiguity.

Proposal:
- **For vertical faces** (front/back/left/right): `horizontal`, `vertical`
- **For horizontal faces** (base/lid): `x_offset`, `y_offset` (or keep `from_back`, `from_left` since those are unambiguous)

Or use a single pair of names that work everywhere: `position_a`, `position_b` with clear docs on what each means per face. But that's too abstract.

#### The Mapper (module.go)

*[Calculates refactor impact]*

If we're keeping face-relative coordinates, then the question becomes: **What names minimize confusion?**

Current problems:
- `from_back` doesn't mean "from the back face"—it means "from the back edge of the box"
- `from_left` doesn't mean "from the left edge of the face"—it means "from the left edge of the box" (except for side faces where it's vertical)

Proposed solutions from the debate:

1. **`from_face_left`, `from_face_bottom`** — Clear for vertical faces, ambiguous for horizontal faces
2. **`horizontal`, `vertical`** — Clear for vertical faces, wrong for horizontal faces (base/lid aren't "vertical")
3. **`along_face`, `across_face`** — Abstract but consistent
4. **`x_offset`, `y_offset` (or z_offset)** — Brings back global coordinates

I lean toward **option 3: `along_face`, `across_face`**.

Definition:
- `along_face`: Position along the face's primary axis (left-to-right when looking at the face)
- `across_face`: Position along the face's secondary axis (top-to-bottom for vertical faces, front-to-back for base, back-to-front for lid)

Wait, that's still ambiguous for base/lid. *[Rethinks]*

Okay, here's a clearer pair:
- `u`: First coordinate (horizontal for walls, x-axis for base/lid)
- `v`: Second coordinate (vertical for walls, y-axis for base/lid)

That's standard UV mapping terminology from 3D graphics. Every face has a (u, v) coordinate system. No ambiguity, consistent across all faces, and familiar to anyone who's done 3D modeling.

#### Dr. Sarah Chen (API Designer)

*[Reviews options]*

The Mapper's UV proposal is clever but **too technical**. YAPP users aren't 3D graphics engineers—they're embedded systems developers who need to cut holes for USB ports.

Here's my final proposal:

**Use different field names for vertical vs. horizontal faces:**

For vertical faces (front/back/left/right):
```yaml
- face: front
  horizontal_position: 20  # or just 'horizontal'
  height_from_base: 5      # or just 'height'
```

For horizontal faces (base/lid):
```yaml
- face: base
  x_from_back: 30
  y_from_left: 20
```

**Rationale:** This matches how users think:
- When placing a cutout on a wall, they think "horizontal position" and "height"
- When placing a cutout on the base/lid, they think "x and y from edges"

The field names literally change by face type in the schema. No ambiguity, no translation tables, completely self-documenting.

**Implementation:** The Mapper checks the face type and expects different fields. Users get validation errors if they use wall fields for base/lid or vice versa.

#### Jamie (YAML User)

*[Tests the UV approach]*

Wait, UV coordinates? I've heard of those but I have no idea what they mean. Are they like X/Y but for textures?

If I see:

```yaml
- face: back
  u: 20
  v: 5
```

I have absolutely no idea what u and v represent. Is u horizontal or vertical? Do I need to look it up every time?

**Vote against UV naming.** Too abstract.

Sarah's proposal (different fields per face type) makes more sense, but isn't that overcomplicated? Now I need to remember which fields are valid for which faces.

What if we just go with the simplest option: **Keep the current names** (`from_back`, `from_left`, `pos_z`) but improve the documentation to explain them clearly?

The current names aren't perfect, but everyone who's used YAPP understands them. If we just add better docs, examples, and maybe validation hints, that's less disruptive than a full rename.

---

## Moderator Summary

### Key Findings

**1. Consensus: Cutouts deserve face-relative coordinates**

All participants agreed that cutouts should NOT switch to x/y/z like other modules. Cutouts span 6 faces with different orientations—face-relative naming is appropriate.

**2. Current naming problems identified:**

- `from_back` is ambiguous: doesn't mean "from back face," means "from back edge of box globally"
- `from_left` is overloaded: means vertical for side faces, horizontal for base/lid
- `pos_z` is a workaround for `from_left` confusion on side faces

**3. Proposed solutions (no consensus):**

| Proposal | Pros | Cons | Support |
|----------|------|------|---------|
| `from_face_left`/`from_face_bottom` | Clear for vertical faces | Ambiguous for base/lid, requires field order change | Sarah, Jamie (partial) |
| `horizontal`/`vertical` | Intuitive for walls | Wrong for base/lid | Jamie |
| `along_face`/`across_face` | Consistent abstraction | Still ambiguous for base/lid | - |
| `u`/`v` (UV mapping) | 3D industry standard | Too technical for target users | The Mapper (reluctant) |
| Different fields per face type | No ambiguity, self-documenting | Complex schema, validation overhead | Sarah (final proposal) |
| Keep current names, improve docs | No breaking changes | Doesn't fix root confusion | Jamie (practical fallback) |

**4. Deep disagreement on base/lid naming:**

Vertical faces (front/back/left/right) all benefit from clearer naming. But horizontal faces (base/lid) don't map naturally to "left/bottom" terminology. No proposed naming scheme cleanly handles all 6 faces.

**5. Implementation concerns:**

- The Mapper noted field-order implications (pos0/pos1 mapping to YAPP arrays)
- 60 existing usages across 9 example files would need updates
- Validation, error messages, and documentation would need rewrites

### Unresolved Questions

1. **What does "bottom" mean for the base face?** Back edge or front edge? Ambiguity prevents consensus on `from_face_bottom`.

2. **Should we have one naming scheme for all 6 faces or different schemes for vertical vs. horizontal?** Single scheme is simpler; split scheme is clearer but more complex.

3. **Is the cognitive load of current naming high enough to justify breaking changes?** Jamie suggested improving docs might be sufficient.

### Recommendations for Decision-Maker

**Option A: Minimal Change (Improve Documentation)**
- Keep `from_back`, `from_left`, `pos_z`
- Add comprehensive docs explaining face-by-face mappings
- Add validation hints ("For side faces, consider using pos_z for height")
- **Pros:** No breaking changes, least implementation work
- **Cons:** Doesn't fix root naming confusion

**Option B: Face-Type-Specific Fields (Most Explicit)**
- Vertical faces: `horizontal_position`, `height_from_base`
- Horizontal faces: `x_from_back`, `y_from_left`
- **Pros:** Zero ambiguity, self-documenting
- **Cons:** Complex schema, different fields per context

**Option C: Hybrid (Add Aliases, Don't Deprecate Old Names)**
- Accept both old and new names during a transition period
- `from_back` also accepts `horizontal` or `along_wall`
- `from_left` also accepts `vertical`, `height`, or `up_from_base` (for side faces)
- **Pros:** Gradual adoption, no forced migration
- **Cons:** Two ways to do the same thing

**Option D: Complete Redesign (Long-Term)**
- Define a clear, unambiguous coordinate system for all faces
- Possibly UV-based or explicit 3D positioning
- Document thoroughly with diagrams
- **Pros:** Foundation for future clarity
- **Cons:** High implementation cost, user retraining

### My Recommendation

**Start with Option C (Hybrid Aliases), plan for Option A (Better Docs).**

**Immediate actions:**
1. Keep current field names (`from_back`, `from_left`, `pos_z`)
2. Add comprehensive face-by-face documentation with diagrams
3. Accept aliases like `horizontal`, `vertical`, `height` alongside current names
4. Improve validation error messages ("Did you mean pos_z for height?")

**Future consideration:**
If user feedback shows continued confusion after improved docs, revisit Option B (face-type-specific fields) in a future DSL version.

**Rationale:** The debate revealed that no proposed naming scheme is clearly superior to all participants. Rather than forcing a disruptive change, we should first **improve the explanation** of the current system. If that's insufficient, we have a clear path forward (Option B).

---

## References

- Current implementation: `pkg/yappgen/modules/cutouts/module.go` (lines 67-90)
- Usage examples: `examples/yapp-demo-lighttubes.yaml`, `examples/compare/cutouts_all_faces.yaml`
- Documentation: `pkg/docs/tutorials/yapp-dsl-reference.md` (Face and Axes section)
- Related modules: `pkg/yappgen/modules/{pcbstands,lighttubes,connectors}/schema.yaml`

