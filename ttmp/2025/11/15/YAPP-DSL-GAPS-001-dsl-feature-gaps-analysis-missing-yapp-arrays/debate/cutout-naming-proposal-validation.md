# Debate Round: Specific Naming Proposal with Validation

**Proposal:** 
- `from_face_left` (all faces) for horizontal positioning
- `from_face_bottom` (front/back/left/right only) for vertical positioning  
- `from_face_back` (base/lid only) for depth positioning
- **Strict validation:** Using wrong fields for a face type produces an error

**Date:** 2025-11-16

## Participants (Same as Previous Debate)

- Dr. Sarah Chen (API Designer)
- Professor Alex Kim (Coordinate Systems Expert)
- The Mapper (module.go implementation)
- Jamie (YAML User)

## Pre-Debate Research

**Current usage to update:**
- 28 instances of `from_back` → become `from_face_left` (all) or `from_face_back` (base/lid)
- 32 instances of `from_left` → become `from_face_left` (horizontal) or `from_face_bottom` (side faces)
- 6 instances of `pos_z` → become `from_face_bottom`

**Schema changes required:**

```yaml
# Current (single schema for all faces)
fields:
  from_back: {type: number, required: true}
  from_left: {type: number, required: true}
  pos_z: {type: number, optional: true}

# Proposed (conditional fields by face type)
fields:
  from_face_left: {type: number, required: true}  # All faces
  from_face_bottom: {type: number, required: true, condition: "face in [front,back,left,right]"}
  from_face_back: {type: number, required: true, condition: "face in [base,lid]"}
```

**The Mapper's translation would become:**

```go
// All faces
horizontal := item.FromFaceLeft

// Side faces
if isSideFace(item.Face) {
    vertical := item.FromFaceBottom  // Required, validated
    if item.FromFaceBack != nil {
        return error("from_face_back not valid for side faces")
    }
}

// Base/Lid faces
if isHorizontalFace(item.Face) {
    depth := item.FromFaceBack  // Required, validated
    if item.FromFaceBottom != nil {
        return error("from_face_bottom not valid for base/lid faces")
    }
}
```

---

## Opening Statements

### Dr. Sarah Chen (API Designer)

*[Reviews the proposal]*

**I LOVE this.** This is exactly the right balance of clarity and enforceability.

Let me break down why this works:

**1. Naming is self-documenting:**

```yaml
# Side face (front) - crystal clear
- face: front
  from_face_left: 20    # Horizontal along the face
  from_face_bottom: 5   # Vertical from bottom edge
```

When you read `from_face_bottom`, you immediately know: "This is the vertical position from the bottom of THIS face." No mental translation required.

```yaml
# Base face - also clear
- face: base
  from_face_left: 30    # Left-to-right
  from_face_back: 20    # Back-to-front (depth)
```

`from_face_back` tells you: "This is the position from the back edge of the base." Much clearer than the old `from_back` which didn't specify "of what?"

**2. Validation prevents errors:**

The beauty of strict validation is that users get immediate feedback:

```yaml
# User tries this (WRONG)
- face: front
  from_face_left: 20
  from_face_back: 5      # ERROR!
```

Error message: `"from_face_back is not valid for face 'front'. Did you mean from_face_bottom?"`

This is **teaching through errors**. Users learn the correct model by trying and getting helpful corrections. Much better than silent failures or subtle positioning bugs.

**3. Consistency within face types:**

All side faces use the same fields (`from_face_left`, `from_face_bottom`). All horizontal faces use the same fields (`from_face_left`, `from_face_back`). The naming pattern is consistent within each category.

**4. Eliminates ambiguity:**

No more "what does from_left mean for this face?" The field name explicitly states its purpose for the face type. `from_face_bottom` only exists for faces that HAVE a bottom edge (sides). `from_face_back` only exists for faces that HAVE a back edge (base/lid).

**My vote: STRONG YES.** Implement this immediately.

### Professor Alex Kim (Coordinate Systems Expert)

*[Draws face diagrams with the new coordinate labels]*

This proposal solves the base/lid ambiguity problem from the previous round.

**What I like:**

**Explicit face-type awareness:** By having different second fields for vertical vs. horizontal faces, we acknowledge that these are fundamentally different geometry problems. A vertical face has "bottom" (the floor edge). A horizontal face has "back" (the rear edge). The naming reflects physical reality.

**Spatial intuition:** When I imagine standing in front of the box:
- Looking at the front wall: `from_face_left` (left edge) and `from_face_bottom` (floor) make immediate sense
- Looking down at the base: `from_face_left` (left edge) and `from_face_back` (back edge) make immediate sense

The names map to how you'd naturally describe positions in 3D space.

**What concerns me:**

**Validation complexity:** The proposal requires conditional field requirements. The schema says "from_face_bottom is required IF face is a side, and FORBIDDEN if face is base/lid." That's more complex than a simple required/optional flag.

Can the schema validation system handle this? Or does The Mapper need custom validation logic?

**Learning curve:** New users need to learn:
1. `from_face_left` is universal
2. `from_face_bottom` for sides
3. `from_face_back` for base/lid

That's three concepts instead of two. But the tradeoff is worth it because each concept is self-explanatory once learned.

**My position:** YES, with the caveat that we need robust validation implementation. The error messages must be extremely clear about which fields are valid for which faces.

**Suggested error format:**

```
Error: Field 'from_face_back' is not valid for face 'front' (vertical face).
For vertical faces (front/back/left/right), use:
  - from_face_left (horizontal position)
  - from_face_bottom (vertical position)

For horizontal faces (base/lid), use:
  - from_face_left (horizontal position)  
  - from_face_back (depth position)
```

### The Mapper (module.go)

*[Calculates implementation changes]*

Let me work through the implementation implications.

**Current code structure:**

```go
// Single pair of fields, complex mapping logic
pos0 := item.FromBack   // Means different things per face
pos1 := item.FromLeft   // Means different things per face
// Then: face-specific remapping to posx/posy/posz
```

**Proposed code structure:**

```go
// Universal horizontal
horizontal := item.FromFaceLeft

// Face-type-specific vertical/depth
var secondary float64
faceLower := strings.ToLower(item.Face)

switch {
case faceLower == "front" || faceLower == "back" || faceLower == "left" || faceLower == "right":
    // Side faces: vertical position required
    if item.FromFaceBottom == nil {
        return error("from_face_bottom is required for side faces")
    }
    if item.FromFaceBack != nil {
        return error("from_face_back is not valid for side faces")
    }
    secondary = *item.FromFaceBottom

case faceLower == "base" || faceLower == "lid":
    // Horizontal faces: depth position required
    if item.FromFaceBack == nil {
        return error("from_face_back is required for base/lid faces")
    }
    if item.FromFaceBottom != nil {
        return error("from_face_bottom is not valid for base/lid faces")
    }
    secondary = *item.FromFaceBack
}

// Map to YAPP coordinates
pos0, pos1 := mapToYAPPCoords(faceLower, horizontal, secondary)
```

**Changes required:**

1. **Schema:** Add `from_face_left`, `from_face_bottom`, `from_face_back` fields (all optional at schema level)
2. **Validation:** Add `CustomValidate()` method to enforce face-type-specific requirements
3. **Mapping:** Update coordinate translation logic (straightforward)
4. **Tests:** Update 60+ usages across 9 example files
5. **Docs:** Update reference documentation

**Effort estimate:** 4-6 hours of focused work.

**My concerns:**

1. **Optional fields at schema level, required in code:** The schema says all three fields are optional (because different faces need different ones), but the code enforces "this one is required for this face type." That's a validation pattern our current schema system doesn't express declaratively.

2. **Error message quality:** We need to provide excellent error messages because users will hit these errors often when learning. The message needs to explain WHY the field is wrong and WHAT to use instead.

3. **Example migration:** All our example files break. We need to update them atomically and regenerate all comparison STLs to verify nothing changes visually.

**My vote:** YES, but with thorough testing. The validation logic is more complex than current code, so we need comprehensive test coverage for all face types and error conditions.

**Test cases needed:**
- Each face type with correct fields (6 cases)
- Each face type with wrong fields (12+ error cases)
- Boundary cases (missing fields, extra fields, typos)

### Jamie (YAML User)

*[Tries writing cutouts with the new naming]*

Let me test this with my real use case: USB-C port on back wall.

**Old way:**

```yaml
- face: back
  from_back: 20   # Wait, is this horizontal or vertical?
  from_left: 5    # I always confuse this
```

**New way:**

```yaml
- face: back
  from_face_left: 20      # Left edge of the back wall - makes sense!
  from_face_bottom: 5     # Bottom edge of the back wall - obvious!
```

**THIS IS SO MUCH BETTER.** The names tell me exactly what they mean.

Let me try a base cutout (ventilation):

**Old way:**

```yaml
- face: base
  from_back: 30   # Is this X or Y? I always forget
  from_left: 20
```

**New way:**

```yaml
- face: base
  from_face_left: 20      # Left edge - clear
  from_face_back: 30      # Back edge - clear
```

Also better! `from_face_back` explicitly tells me this is measuring from the back edge of the base.

**Now let me try making a mistake on purpose:**

```yaml
- face: front
  from_face_left: 20
  from_face_back: 5    # Oops, wrong field
```

If this gives me an error like: `"from_face_back is not valid for face 'front'. Use from_face_bottom for vertical position."` — then I'll immediately know what to fix.

**What I love:**
1. Field names explain themselves
2. Errors teach me the right model
3. No more lookup tables needed

**What I'm worried about:**
1. **Will I get confused about when to use `from_face_back` vs `from_face_bottom`?** 
   - Probably not, because the error messages will teach me
   - And once I learn "sides use bottom, base/lid use back" it'll stick

2. **Is this too much validation?** 
   - No, I actually WANT strict validation. Better to get an error than to silently place my cutout in the wrong spot and not realize until I print the box.

**My vote: STRONG YES.** This naming makes the DSL actually learnable. I can write cutouts without constantly referring to docs.

---

## Quick Rebuttals

### Dr. Sarah Chen (API Designer)

The Mapper raises a good point about schema-level validation vs. code-level validation. Here's how we handle it:

**Schema declaration:**

```yaml
fields:
  from_face_left:
    type: number
    required: true
    desc: "Horizontal position along the face (left edge), used by all faces"
  
  from_face_bottom:
    type: number
    desc: "Vertical position from bottom edge (front/back/left/right faces only)"
  
  from_face_back:
    type: number  
    desc: "Depth position from back edge (base/lid faces only)"
```

All three fields are declared, but only `from_face_left` is marked `required: true` at schema level. The conditional requirements live in `CustomValidate()`:

```go
func (x *CutoutsItem) CustomValidate() error {
    face := strings.ToLower(x.Face)
    isSide := face == "front" || face == "back" || face == "left" || face == "right"
    isHoriz := face == "base" || face == "lid"
    
    if isSide {
        if x.FromFaceBottom == nil {
            return fmt.Errorf("from_face_bottom is required for face '%s'", x.Face)
        }
        if x.FromFaceBack != nil {
            return fmt.Errorf("from_face_back is not valid for face '%s' (vertical face). Use from_face_bottom instead", x.Face)
        }
    }
    
    if isHoriz {
        if x.FromFaceBack == nil {
            return fmt.Errorf("from_face_back is required for face '%s'", x.Face)
        }
        if x.FromFaceBottom != nil {
            return fmt.Errorf("from_face_bottom is not valid for face '%s' (horizontal face). Use from_face_back instead", x.Face)
        }
    }
    
    return nil
}
```

This pattern is already used in the codebase (e.g., polygon cutouts require `polygon` field to be set when `shape=polygon`). We're just extending it.

**Effort is worth it:** Yes, we need to update 9 example files and regenerate STLs. But that's a one-time cost. The ongoing benefit of clearer naming affects every future user.

### Professor Alex Kim (Coordinate Systems Expert)

Jamie's enthusiasm is telling. When the end user says "this is SO MUCH BETTER," that's strong validation.

My remaining concern about learning curve is addressed by the error-driven learning model. Users don't need to memorize the rules—they'll learn by doing:

1. Try to add cutout
2. Get error if wrong fields used
3. Error message tells them exactly what to fix
4. Next time they remember

This is how good APIs teach their users. The validation becomes documentation.

**Implementation recommendation:** Add examples to error messages:

```
Error: from_face_back is not valid for face 'front'.

For vertical faces (front/back/left/right), use:
  - from_face_left: 20
  - from_face_bottom: 5

Example:
  - face: front
    from_face_left: 20
    from_face_bottom: 5
    width: 10
    length: 8
```

Include a minimal working example in the error message itself. Users can copy-paste it as a starting template.

### The Mapper (module.go)

I'm convinced. The validation complexity is manageable, and the benefits are clear.

**Implementation plan:**

1. **Add new fields to schema** (keep old ones temporarily for migration)
2. **Implement CustomValidate()** with clear error messages
3. **Update Build() method** to use new fields
4. **Migrate all examples** atomically in one commit
5. **Regenerate STLs** and verify no visual regressions
6. **Update docs** with before/after examples
7. **Remove old fields** after confirming everything works

**Migration path for users:**
- Short transition period where both old and new fields are accepted
- Deprecation warnings: "from_back is deprecated, use from_face_left"
- Eventually: hard error for old fields

Actually, scratch that migration path. The user said "we're not interested in backwards compatibility." So we can go straight to the new naming without a transition period. Update all examples, update docs, done.

**Time estimate:** 
- Schema + validation: 1 hour
- Build() method update: 1 hour
- Example migration: 2 hours
- Testing + STL regen: 2 hours
- **Total: ~6 hours**

### Jamie (YAML User)

One more thing I realized: **The naming helps me visualize the box orientation.**

When I see:
```yaml
- face: front
  from_face_left: 20
  from_face_bottom: 5
```

I immediately picture the front wall with a left edge and a bottom edge, and I can place the cutout mentally. The names create a mental model of the 3D object.

With the old names:
```yaml
- face: front  
  from_back: 20
  from_left: 5
```

I had to think: "wait, back? But I'm on the front face. What does back mean here?" It broke my mental model.

**New names = better spatial reasoning.**

---

## Moderator Summary

### Unanimous Consensus: IMPLEMENT THE PROPOSAL

All four participants voted YES with strong enthusiasm.

**Key strengths identified:**

1. **Self-documenting naming:** Field names explicitly state what they measure relative to the face
2. **Eliminates ambiguity:** `from_face_back` only exists for faces with a back edge (base/lid)
3. **Error-driven learning:** Validation errors teach users the correct model
4. **Spatial intuition:** Names map to how humans naturally describe 3D positions
5. **Face-type consistency:** All sides use same fields, all horizontal faces use same fields

**Implementation considerations:**

1. **Conditional validation required:** CustomValidate() enforces face-type-specific field requirements
2. **Error messages are critical:** Must be clear, helpful, and include examples
3. **Example migration:** 60+ usages across 9 files need updating
4. **Testing coverage:** Need comprehensive tests for all face types and error cases
5. **Documentation:** Update reference docs with before/after examples

**No concerns that block implementation.** All identified concerns have clear solutions.

### Recommendation to Decision-Maker

**APPROVE and IMPLEMENT immediately.**

**Specific naming:**
- `from_face_left` (required, all faces) — horizontal position along face
- `from_face_bottom` (required for front/back/left/right) — vertical position from bottom edge  
- `from_face_back` (required for base/lid) — depth position from back edge

**Validation:**
- Strict checking: wrong field for face type produces error
- Error messages include explanation + example
- No backwards compatibility (clean break)

**Implementation checklist:**
- [ ] Update schema.yaml with three fields
- [ ] Implement CustomValidate() with face-type checking
- [ ] Update Build() method to use new fields
- [ ] Add comprehensive test suite (18+ test cases)
- [ ] Migrate all 9 example files
- [ ] Regenerate comparison STLs
- [ ] Update pkg/docs/tutorials/yapp-dsl-reference.md
- [ ] Update face/axes documentation section

**Estimated effort:** 6 hours of focused development + testing

---

## Final Vote

- **Dr. Sarah Chen (API Designer):** STRONG YES
- **Professor Alex Kim (Coordinate Expert):** YES
- **The Mapper (module.go):** YES with confidence
- **Jamie (YAML User):** STRONG YES

**Result: 4-0, unanimous approval with enthusiasm.**

---

## References

- Previous debate: `cutout-position-naming-debate.md`
- Current implementation: `pkg/yappgen/modules/cutouts/module.go`
- Example files: `examples/*.yaml` (60+ usages to update)
- Documentation: `pkg/docs/tutorials/yapp-dsl-reference.md`

