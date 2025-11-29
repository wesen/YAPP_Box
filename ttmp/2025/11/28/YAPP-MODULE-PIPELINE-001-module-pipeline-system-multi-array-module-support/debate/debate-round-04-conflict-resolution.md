---
Title: Debate Round 4 - Conflict Resolution: What Happens When Modules Collide?
Ticket: YAPP-MODULE-PIPELINE-001
Status: active
Topics:
  - dsl
  - modules
  - architecture
  - debate
DocType: debate
Intent: long-term
Owners: []
RelatedFiles: []
Summary: Fourth debate round exploring how to handle conflicts when composite modules generate DSL that overlaps with user-defined entries
LastUpdated: 2025-11-28
---

# Debate Round 4: Conflict Resolution - What Happens When Modules Collide?

## Question

When a composite module generates DSL entries that conflict with user-defined entries (e.g., both define `cutouts`), what should happen: append, replace, error, or namespace?

## Primary Candidates

- Jordan "The Module Author" Kim (wants predictable behavior)
- Alex "The Pragmatist" Chen (wants least friction)
- Dr. Sarah "The Architect" Martinez (wants clear semantics)
- `pkg/yappgen/modules/cutouts/` "The Multi-Array Pioneer" (has experience with arrays)

---

## Pre-Debate Research

### Research 1: Potential Conflict Scenarios

**Scenario 1: User defines cutouts, composite module also generates cutouts**
```yaml
features:
  cutouts:
    - face: front
      from_face_left: 10
      radius: 5
      shape: circle
  
  lcd:
    - display:
        face: front
        position: [50, 30]
```
→ LCD module generates cutout for display window

**Conflict:** Both `cutouts` and `lcd` want to add to `cutoutsFront` array

**Possible Outcomes:**
- **Append:** Both cutouts exist (user's circle + LCD's rectangle) ✅
- **Replace:** LCD overwrites user's cutouts ❌
- **Error:** Fail with "cutouts already defined" ❌
- **Namespace:** Separate arrays (cutoutsFront vs lcdCutoutsFront) ❌

---

**Scenario 2: User defines pcb_stands for PCB mounting, composite module generates pcb_stands for LCD mounting**
```yaml
features:
  pcb_stands:
    - x: 5
      y: 5
      diameter: 7
  
  lcd:
    - mounting:
        type: pcb_stands
```
→ LCD module generates 4 pcb_stands for LCD mounting

**Conflict:** Both want to add to `pcbStands` array

**Possible Outcomes:**
- **Append:** PCB stands + LCD stands (both exist) ✅
- **Replace:** LCD stands overwrite PCB stands ❌
- **Error:** Fail ❌

---

**Scenario 3: Multiple composite modules generate same array type**
```yaml
features:
  lcd:
    - display: {...}
  
  sensor_module:
    - sensor: {...}
```
→ Both generate cutouts

**Conflict:** Multiple composite modules want to add cutouts

**Possible Outcomes:**
- **Append:** All cutouts merged ✅
- **Namespace:** lcdCutouts vs sensorCutouts ❌

---

### Research 2: Array Semantics in Current System

**Analysis of `pkg/yappgen/model.go`:**

**Model fields are slices:**
```go
PcbStands   []map[string]any
Cutouts     []map[string]any
```

**Semantics:**
- **Additive** - Each entry adds a feature (stand, cutout, etc.)
- **Order independent** - Array order doesn't matter for most features
- **No conflicts** - Multiple entries can coexist

**Example:** If user defines 4 pcb_stands and LCD adds 4 more, result is 8 stands total.

**This suggests: Append is the natural semantics.**

---

### Research 3: YAML Merge Strategies

**Common YAML merge strategies:**

**1. Deep Merge (maps)**
```yaml
base:
  a: 1
  b: 2

override:
  b: 3
  c: 4

result:
  a: 1
  b: 3  # Overridden
  c: 4  # Added
```

**2. Array Append**
```yaml
base:
  items: [1, 2]

override:
  items: [3, 4]

result:
  items: [1, 2, 3, 4]  # Appended
```

**3. Array Replace**
```yaml
base:
  items: [1, 2]

override:
  items: [3, 4]

result:
  items: [3, 4]  # Replaced
```

**For feature arrays:** Array Append makes sense (additive semantics)

---

### Research 4: Error Handling Analysis

**If we error on conflicts:**

**User writes:**
```yaml
features:
  cutouts:
    - face: front
      ...
  
  lcd:
    - display: ...
```

**Error:**
```
Error: composite module 'lcd' conflicts with user-defined 'cutouts'
```

**User must:**
1. Remove their cutouts? (loses functionality)
2. Rename their cutouts? (not possible)
3. Not use lcd module? (defeats purpose)

**Erroring on conflicts is user-hostile.** Forces choice: user features OR composite module.

---

## Opening Statements (Round 1)

### Jordan "The Module Author" Kim

*[Reviews conflict scenarios, user experience]*

After analyzing conflict scenarios, the answer is clear: **APPEND**.

**Here's why:**

**1. Feature Arrays are Additive**
- PCB stands: Each entry adds a stand
- Cutouts: Each entry adds a cutout
- Connectors: Each entry adds a connector

**The semantics are naturally additive.** More entries = more features.

**If LCD module generates cutouts** and **user defines cutouts**, the natural result is: **both exist**.

**Example:**
```yaml
features:
  cutouts:
    - name: power-jack
      face: back
      ...
  
  lcd:
    - display:
        face: front
        ...
```

**Result:** Power jack cutout (back) + LCD display cutout (front) = **both cutouts**.

**This is what users expect.**

**2. Replace is Destructive**
If composite modules **replace** user entries, users lose their features. That's wrong.

**3. Error is User-Hostile**
If we error on conflicts, users must choose: their features OR composite module. That defeats the purpose.

**4. Namespace is Overcomplicated**
Separate arrays (lcdCutouts vs cutouts) add complexity for no benefit.

**My position: APPEND (always).**

Composite-generated entries are **merged** with user-defined entries. No conflicts, no errors, natural semantics.

---

### Alex "The Pragmatist" Chen

*[Considers user workflows, friction points]*

Jordan, I agree with APPEND. But let me add a **practical** consideration.

**User workflow:**
```yaml
features:
  # User manually adds cutouts
  cutouts:
    - face: front
      ...
  
  # Later, user discovers lcd module
  lcd:
    - display:
        face: front
        ...
```

**With APPEND:** Both cutouts exist. **No problem** ✅

**But what if they overlap?**

**Scenario:** User manually created LCD cutout, then adds lcd module.

**Result:** Two overlapping cutouts (user's + module's).

**This is a problem.** Overlapping cutouts might cause rendering issues.

**Options:**

**Option A: Trust users**
- Append always
- Users are responsible for avoiding overlaps
- **Simple** ✅

**Option B: Detect overlaps**
- Check if generated cutouts overlap with existing
- Warn or error if overlap detected
- **Complex** ❌

**Option C: Let users opt-out**
- Composite modules generate entries
- Users can explicitly disable (e.g., `lcd.disable_cutout: true`)
- **Flexible** ✅

**I prefer Option A (trust users).** Keep it simple.

**But:** Provide good error messages if rendering fails due to overlaps.

**My position: APPEND with clear documentation** about potential overlaps.

---

### Dr. Sarah "The Architect" Martinez

*[Analyzes merge semantics, architectural implications]*

I agree with APPEND, but let me formalize the semantics.

**Merge Strategy:**

**For arrays (features.cutouts, features.pcb_stands, etc.):**
```
result = user_entries ++ composite_entries
```

**Order:**
- User entries come **first**
- Composite entries **appended**
- **Deterministic order**

**For maps (non-arrays):**
- Composite modules should **not** generate non-array entries
- Stick to feature arrays only

**Why this works:**
1. **Additive semantics** - More entries = more features
2. **User control** - User entries take precedence (rendered first)
3. **Predictable** - Order is deterministic
4. **Simple** - No complex merging logic

**But there's a subtlety:** What if **multiple composite modules** generate the same array?

**Example:**
```yaml
features:
  lcd:
    - display: ...  # Generates cutouts
  
  sensor_module:
    - sensor: ...   # Also generates cutouts
```

**Merge order:**
```
result = user_cutouts ++ lcd_cutouts ++ sensor_cutouts
```

**This works.** All cutouts appended.

**My position: APPEND with deterministic order** (user first, then composites in registration order).

---

### `pkg/yappgen/modules/cutouts/` - "The Multi-Array Pioneer"

*[Reflects on array management experience]*

I've managed arrays for a while. Let me share what I've learned.

**Array conflicts in my world:**
User writes:
```yaml
features:
  cutouts:
    - face: front
      ...
    - face: back
      ...
```

I distribute to:
- `cutoutsFront`: [front entry]
- `cutoutsBack`: [back entry]

**No conflicts.** Each entry goes to its own array (by face).

**For composite modules:**
User writes:
```yaml
features:
  cutouts:
    - face: front
      ...  # Entry 1
  
  lcd:
    - display:
        face: front  # Also generates cutoutsFront entry
```

Result:
- `cutoutsFront`: [user entry, lcd entry]

**This is fine.** Both entries exist. Both rendered.

**But here's my concern:** **Overlap validation.**

**What if:**
- User entry: Circle at (10, 10), radius 5
- LCD entry: Rectangle at (12, 12), size 20x20

**They overlap.** OpenSCAD will render both, creating a weird merged cutout.

**Should we:**
- Detect this? (complex geometry check)
- Warn user? (how?)
- Ignore it? (user problem)

**I think:** Ignore it. **Users are responsible for their geometry.**

**If they define overlapping cutouts, that's their choice.** Whether manual or composite-generated.

**My position: APPEND. No overlap detection.**

Simple, predictable, user-controlled.

---

## Rebuttals (Round 2)

### Alex "The Pragmatist" Chen → Sarah

Sarah, you said: "User entries take precedence (rendered first)."

**Why does order matter?** Are you thinking about overlap handling?

**If user entry is rendered first**, and LCD entry overlaps, does one override the other?

**I think:** In OpenSCAD, overlapping cutouts **union**. Order doesn't matter.

**So "user first" is just convention**, not semantic meaning.

**That's fine.** But don't oversell it as "user control."

---

### Dr. Sarah "The Architect" Martinez → Alex

Alex, you're right - in OpenSCAD, overlapping cutouts union, so order doesn't affect geometry.

**But order DOES matter for other reasons:**

**1. Error messages:**
If validation fails, error message shows: "features.cutouts[0]" vs "features.cutouts[5]"

**User entries first** means user's entries have predictable indices.

**2. Debugging:**
When inspecting generated SCAD, user sees their entries first, then composite-generated.

**Clear separation** in the array.

**3. Provenance:**
We track where each entry came from. Order helps: "entries 0-2 from user, 3-6 from lcd module."

**So order matters for UX**, even if not for geometry.

---

### Jordan "The Module Author" Kim → Cutouts

Cutouts, you said: "No overlap detection."

**I disagree.** Overlap detection would be valuable.

**Scenario:**
```yaml
features:
  cutouts:
    - face: front
      from_face_left: 50
      radius: 10
  
  lcd:
    - display:
        face: front
        position: [50, 30]
        size: [40, 20]
```

**User's cutout** (circle at 50) **overlaps with LCD display** (rectangle at 50).

**This is probably a mistake.** User might not realize LCD module generates a cutout.

**Should we:**
- **Warn:** "LCD module generates cutout at (50, 30) which may overlap with cutouts[0]"
- **Error:** "Overlap detected, please resolve"
- **Ignore:** User's problem

**I think: Warn, but don't error.**

**Rationale:**
- Warnings don't block users
- Users might intentionally overlap (merge cutouts)
- But we help catch mistakes

**Implementation:**
- Simple bounding box check
- If composite-generated cutout overlaps user cutout, warn
- ~50 lines of code

**Worth it for better UX.**

---

### `pkg/yappgen/modules/cutouts/` → Jordan

Jordan, overlap detection sounds good in theory. But **bounding box check is insufficient**.

**Example:**
- Circle at (50, 50), radius 10: Bounding box [40-60, 40-60]
- Rectangle at (45, 45), size 5x5: Bounding box [42.5-47.5, 42.5-47.5]

**Bounding boxes don't overlap**, but **geometries DO** (rectangle inside circle).

**Proper overlap detection needs:**
- Shape-aware geometry checks
- Circle-circle intersection
- Rectangle-circle intersection
- Polygon-polygon intersection
- **100+ lines of code, complex math**

**Not worth it.**

**Also:** Users might **intentionally** create overlapping cutouts:
- Merge shapes
- Create complex cutouts
- Boolean operations

**If we warn on overlaps, we annoy users who know what they're doing.**

**My position: No overlap detection.** Users control their geometry.

**Document it clearly:** "Composite modules append to arrays. Check for overlaps manually if needed."

---

### Alex "The Pragmatist" Chen → Jordan

Jordan, I respect the UX thinking. But **overlap detection is scope creep.**

**We're here to answer:** Append, replace, error, or namespace?

**The answer is: APPEND.**

**Overlap detection is a separate feature.** Don't bundle it with conflict resolution.

**If we want overlap detection:**
- Make it a separate ticket
- Design it properly (not just bounding boxes)
- Make it opt-in (flag: `--warn-overlaps`)

**For this ticket: APPEND. Period.**

Simple, predictable, works for all scenarios.

---

### Dr. Sarah "The Architect" Martinez → All

Let me formalize the **merge semantics** based on this discussion.

**Merge Strategy: Append with Deterministic Order**

**Rule 1: Array Append**
```
When composite module generates entries for existing array:
  result_array = user_entries ++ composite_entries
```

**Rule 2: Multiple Composite Modules**
```
When multiple composite modules generate same array:
  result_array = user_entries ++ module1_entries ++ module2_entries
  (module order = registration order)
```

**Rule 3: No Replace or Error**
```
Composite modules never replace or error on existing arrays.
Always append.
```

**Rule 4: Namespace Not Used**
```
Composite modules add to standard arrays (cutouts, pcbStands).
No separate namespaced arrays (lcdCutouts).
```

**Why these rules?**
1. **Predictable** - Always append, no surprises
2. **Composable** - Multiple modules can coexist
3. **User control** - User entries take precedence (rendered first)
4. **Simple** - One strategy, no edge cases

**Edge Cases Handled:**
- Empty user array + composite: Works (append to empty = [composite entries])
- No user array + composite: Works (create array = [composite entries])
- User array + no composite: Works (unchanged = [user entries])
- User array + multiple composites: Works (deterministic append order)

**My position: This is the complete merge strategy.**

---

## Moderator Summary

### Key Arguments

**PRO Append:**
1. **Natural semantics** - Feature arrays are additive
2. **User-friendly** - No forced choices, both features exist
3. **Composable** - Multiple modules can generate same array
4. **Simple** - One strategy, no edge cases
5. **Predictable** - Deterministic order (user first, then composites)

**PRO Overlap Detection:**
1. **Catch mistakes** - Help users avoid unintentional overlaps
2. **Better UX** - Warnings guide users
3. **Teachable moments** - Users learn about composite behavior

**CON Overlap Detection:**
1. **Complex** - Proper geometry checks require 100+ lines
2. **Incomplete** - Bounding box checks miss cases
3. **False positives** - Annoys users with intentional overlaps
4. **Scope creep** - Separate feature from conflict resolution

### Strong Consensus

- ✅ **APPEND is the right strategy**
- ✅ Replace is destructive (rejected)
- ✅ Error is user-hostile (rejected)
- ✅ Namespace adds complexity (rejected)
- ✅ Deterministic order matters for UX
- ✅ User entries first, composite entries appended

### Disagreement on Overlap Detection

- Jordan: PRO (warns users, better UX)
- Cutouts: CON (complex, incomplete, annoying)
- Alex: NEUTRAL (separate feature, not now)
- Sarah: NEUTRAL (good idea, but later)

**Resolution:** Overlap detection is **out of scope** for initial implementation. Can be added later as separate feature.

### Merge Strategy Formalized

```
For each feature array (cutouts, pcb_stands, etc.):
  1. Start with user-defined entries (if any)
  2. Append composite module 1 entries (if any)
  3. Append composite module 2 entries (if any)
  4. ... (in module registration order)
  
Result: Single array with all entries, deterministic order
```

### Implementation Requirements

**Merge Function:**
```go
func mergeDSLEntry(base map[string]any, entry DSLEntry) map[string]any {
    // Navigate to entry.Path (e.g., "features.cutouts")
    // If array exists: append entry.Value
    // If not: create array with entry.Value
    return merged
}
```

**Complexity:** ~30-50 lines for robust path navigation and array merging

### Open Questions Resolved

- ✅ **Append, not replace**
- ✅ **No errors on conflicts**
- ✅ **No namespacing**
- ✅ **User entries first, deterministic order**
- ❓ **Overlap detection: deferred to future work**

### Next Steps

1. **Implement merge function** with append semantics
2. **Document behavior** clearly for module authors
3. **Add tests** for multiple conflict scenarios
4. **Consider overlap detection** as future enhancement

---

## Round 4 Conclusion

**Strong consensus achieved: APPEND strategy with deterministic order.**

This resolves a major design question and provides clear implementation guidance.

**The merge strategy is simple, predictable, and user-friendly.**

Next round should explore: Validation timing and error handling for composite-generated DSL.

