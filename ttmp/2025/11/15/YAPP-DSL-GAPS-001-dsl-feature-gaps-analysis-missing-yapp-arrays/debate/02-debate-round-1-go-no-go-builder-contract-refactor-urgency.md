---
Title: Debate Round 1 — Go/No-Go: Builder Contract Refactor Urgency
Ticket: YAPP-DSL-GAPS-001
Status: active
Topics:
    - yapp
    - architecture
    - codegen
    - debate
DocType: debate
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/01-debate-builder-contract-migration-format-and-candidates.md
      Note: Candidate profiles
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/reference/02-reference-debate-questions-builder-contract-migration.md
      Note: Question set
ExternalSources: []
Summary: Should we invest in refactoring the builder contract now or defer?
LastUpdated: 2025-11-17
---

# Debate Round 1 — Go/No-Go: Builder Contract Refactor Urgency

## Question

**Is refactoring the builder contract urgent, or should we defer?**

**Context:** Current system works but has pain points: YAML re-marshal in every builder, untyped `[]map[string]any` → `[][]any` flow, special-case handling for cutouts.

## Pre-Debate Research

### Module Count and Builder Complexity

```bash
$ ls -1d pkg/yappgen/modules/*/ | wc -l
7

$ find pkg/yappgen/modules -name "module.go" -exec wc -l {} + | tail -1
1455 total
```

**Finding:** 7 modules, ~1455 lines of builder code total (~208 lines/module average).

### YAML Marshal/Unmarshal Pattern

```bash
$ grep -r "yaml\.Marshal(it)" pkg/yappgen/modules --count
pkg/yappgen/modules/cutouts/module.go:1
pkg/yappgen/modules/lighttubes/module.go:1
pkg/yappgen/modules/boxmounts/module.go:1
pkg/yappgen/modules/snapjoins/module.go:1
pkg/yappgen/modules/pcbstands/module.go:1
pkg/yappgen/modules/connectors/module.go:1

$ grep -r "yaml\.Unmarshal(data, &item)" pkg/yappgen/modules --count
pkg/yappgen/modules/cutouts/module.go:1
pkg/yappgen/modules/lighttubes/module.go:1
pkg/yappgen/modules/boxmounts/module.go:1
pkg/yappgen/modules/snapjoins/module.go:1
pkg/yappgen/modules/pcbstands/module.go:1
pkg/yappgen/modules/connectors/module.go:1
```

**Finding:** All 6 non-pushbuttons modules use the marshal→unmarshal pattern. Pushbuttons has custom `decodePushButtonsItems()` helper.

### Recent Module Churn

```bash
$ git log --oneline --since="2 months ago" -- pkg/yappgen/modules/*/module.go | head -14
94cc4b3 Update cutouts schema
f88c302 Add cutout masks
a7abdfd Completed yappAltOrigin and yappPCBName
f7fd532 Add cutout origin and pos_z functionality
3b49176 Add polygon cutout shapes
b9b2011 Add lighttubes support
a6eb179 Adjust some more menial stuff
a352a27 Add hinges and all the other things
6a748d3 Continue working towards the button box example
9d64f50 Continue working towards the button box
361b9d1 Add more modules and integrate them
f51f1a8 Use templates for the code gen
e618d80 Continue working on schemagen
8b69fba Add push buttons module
```

**Finding:** 14 commits to module builders in 2 months. Active development, frequent schema/builder changes.

### Model Storage

```go
// pkg/yappgen/model.go:11-47
type Model struct {
	ProjectName string
	// ... globals ...
	
	// Features
	PcbStands   []map[string]any
	Connectors  []map[string]any
	BoxMounts   []map[string]any
	SnapJoins   []map[string]any
	Cutouts     []map[string]any
	PushButtons []map[string]any
	LightTubes  []map[string]any
	
	PrintSwitchExtenders bool
}
```

**Finding:** Model stores untyped `[]map[string]any` for all features. No compile-time type safety between collection and emission.

---

## Opening Statements

### Alex Chen — "The Pragmatist"

Look, I ran the numbers. We have 7 modules, 1455 lines of builder code, and 14 commits touching module builders in the last 2 months. That's active churn.

But here's the thing: **it works**. Every one of those 14 commits shipped successfully. The YAML marshal→unmarshal pattern is boilerplate, sure, but it's *consistent* boilerplate. Copy-paste from an existing module and you're done in 5 minutes.

I grep'd for error handling in builders:

```bash
$ grep -r "return nil, errors" pkg/yappgen/modules/*/module.go | wc -l
47
```

47 error return sites across 7 modules. That's error handling we've already written and tested. A refactor means re-validating all of that.

**My position:** Defer. We're adding 2-3 modules per month based on git history. Let's finish the feature parity work (labelsPlane, ridgeExt, displayMounts) *first*, then refactor when the dust settles. Don't refactor a moving target.

**ROI question:** What's the cost? If Option A (generate Decode) touches 7 module files + schemagen templates + tests, that's a 3-4 day investment minimum. What do we get? Slightly cleaner builders. Not worth it *right now*.

---

### Jordan Rivera — "The Architect"

Alex, I hear you on "it works," but let's talk about what "works" actually means here.

I traced the data flow from YAML to SCAD:

```
User YAML → Resolver (map[string]any) 
         → Model.PcbStands ([]map[string]any)
         → Build() marshal→unmarshal → typed struct
         → Build() logic → [][]any
         → writeArrayDecl() → SCAD string
```

We have **three type boundaries** with no compile-time checks:
1. Resolver output → Model storage
2. Model storage → Builder input
3. Builder output → Emission

Every one of those is a runtime failure waiting to happen. And we *have* had failures—remember the cutouts `from_face_bottom` bug where we passed a string instead of a number? That should have been caught at compile time.

Alex says "consistent boilerplate," but I see **technical debt compounding**. Every new module copies the marshal→unmarshal pattern. That's not consistency, that's proliferation.

**My position:** Go now. The longer we wait, the more modules we have to migrate. We're at 7 modules today. If we wait until we have 12-15 modules (post-labelsPlane, ridgeExt, displayMounts), the migration cost triples.

**Investment vs. payoff:** Yes, Option A is 3-4 days. But it pays dividends *every time* someone adds a module. Casey (our hypothetical new hire) shouldn't have to understand YAML marshaling internals just to write a builder.

---

### Casey Thompson — "The New Hire"

*[Raises hand nervously]*

I... I actually tried to add a test module last week following the module authoring guide. I got stuck for 2 hours on this part:

```go
// From the guide:
data, err := yaml.Marshal(it)
if err != nil {
    return nil, errors.Wrapf(err, "%s: marshal", label)
}

var item YourModuleItem
if err := yaml.Unmarshal(data, &item); err != nil {
    return nil, errors.Wrapf(err, "%s: unmarshal", label)
}
```

I kept thinking: "Why are we marshaling and then immediately unmarshaling? Isn't `it` already the data I need?"

Then I read the comment in `pushbuttons/module.go`:

```go
// decodePushButtonsItems converts []map[string]any to typed structs
// without going through YAML marshal/unmarshal round-trip
```

So... pushbuttons *doesn't* use the pattern? Why not? Is that the better way? Should I copy pushbuttons or copy pcbstands?

**My position:** I don't know if it's urgent, but it's *confusing*. The fact that there are two patterns (marshal/unmarshal vs custom decode) means I have to make a choice, and I don't know which is "right."

If we're going to defer, can we at least document *why* we do the marshal/unmarshal dance? Because right now it feels like magic.

---

## Rebuttals

### Alex Chen — "The Pragmatist"

Jordan, you're right about the type boundaries, but let's be real: **we've shipped 7 modules with this pattern and how many runtime type errors have we actually seen in production?**

I checked our test suite:

```bash
$ grep -r "func Test.*Build" pkg/yappgen/modules --count
pkg/yappgen/modules/boxmounts/module_test.go:2
pkg/yappgen/modules/connectors/module_test.go:2
pkg/yappgen/modules/pcbstands/module_test.go:2
pkg/yappgen/modules/snapjoins/module_test.go:1
```

We have builder tests. They catch the type errors. The `from_face_bottom` bug you mentioned? That was caught *during development*, not in production, because we have `ValidateStructure` and `ValidateConstraints`.

Casey, I feel your pain. The marshal/unmarshal thing *is* weird. But here's the pragmatic truth: **it works, it's testable, and it's copy-pasteable**. 

Pushbuttons is different because it has nested objects (`cap`, `lid`, `switch`). That's why it needs custom decode logic. For flat structures (pcbstands, connectors, boxmounts), marshal/unmarshal is actually simpler than writing manual field extraction.

**My counter-proposal:** Let's finish the remaining Priority 2 modules (labelsPlane, ridgeExt) using the current pattern. *Then* we refactor. That way we're not changing the foundation while building on top of it.

Timeline: 2-3 weeks to finish P2 modules, then 1 week for Option A refactor. Total: 4 weeks. If we refactor now, we still have to add those modules anyway, so we're not saving time.

---

### Jordan Rivera — "The Architect"

Alex, your "it works" argument is the classic engineer's trap. Yes, it works *today*. But let's project forward.

You said we're adding 2-3 modules per month. That's 6-9 more modules in the next quarter. By the time we "finish P2 modules," we'll have 10 modules, not 7. The refactor cost goes from 3-4 days to 5-6 days.

**Compounding cost curve:**
- Refactor at 7 modules: 4 days
- Refactor at 10 modules: 6 days
- Refactor at 15 modules: 9 days

Every month we defer, we add 1-2 days to the migration cost. That's not pragmatism, that's procrastination.

Casey's confusion is a *signal*, not noise. If a new developer can't understand why we marshal/unmarshal, that's a code smell. And your solution—"let's document it"—is treating the symptom, not the disease.

**The real question:** Do we want to build 6-9 more modules on a foundation we *know* we're going to replace? That's throwaway work.

**My revised position:** Go now, but do Option A (generate Decode) *only*. Don't touch Option B/C/D yet. That's a 2-3 day investment that:
1. Eliminates marshal/unmarshal boilerplate
2. Gives Casey a clear pattern to follow
3. Doesn't block P2 module work (we can add labelsPlane with generated Decode)

We're not refactoring the whole system. We're removing one pain point that affects every future module.

---

### Casey Thompson — "The New Hire"

*[Looks between Alex and Jordan]*

Can I ask a dumb question? If we do Option A and generate `Decode()`, does that mean I don't have to write the marshal/unmarshal code anymore?

*[Jordan nods]*

Okay, then I vote for that. Because honestly, the hardest part of adding a module isn't the business logic (building the `[][]any` array). It's all the ceremony around it—the marshal/unmarshal, the error wrapping, the label formatting.

If schemagen can generate that, then adding a module becomes:
1. Write `schema.yaml`
2. Write `Build()` logic (just the array construction)
3. Done

That sounds way better than:
1. Write `schema.yaml`
2. Copy-paste marshal/unmarshal boilerplate
3. Write `Build()` logic
4. Hope I didn't mess up the error wrapping

Alex, you said "copy-paste is fine," but copy-paste is how bugs happen. I've already seen it—in `lighttubes/module.go`, the error label says `"light_tubes[%d]"` but in `boxmounts/module.go` it says `"box_mounts[%d]"`. That's inconsistent. Generated code would be *consistent*.

---

## Moderator Summary

### Key Arguments

**For "Go Now" (Jordan, Casey):**
- Compounding cost: 7 modules today, 10+ modules in 3 months
- Type safety: eliminate runtime failures at type boundaries
- Developer experience: reduce cognitive load for new module authors
- Consistency: generated code is uniform, copy-paste introduces drift

**For "Defer" (Alex):**
- Working system: 7 modules shipped successfully with current pattern
- Test coverage: builder tests catch type errors during development
- Stable foundation: finish P2 modules before refactoring
- ROI timing: refactor cost is similar whether we do it now or in 3 weeks

### Tensions

1. **Pragmatism vs. Principle:** Alex prioritizes shipping features; Jordan prioritizes long-term maintainability.

2. **Copy-paste vs. Generation:** Alex sees boilerplate as acceptable; Jordan/Casey see it as technical debt.

3. **Timing:** Alex wants to defer until P2 modules are done; Jordan wants to refactor before adding more modules.

### Emerging Consensus

Both sides agree:
- The marshal/unmarshal pattern is suboptimal
- Option A (generate Decode) is the right first step
- Full type safety (Option C) is overkill for now

**The real debate:** Timing. Do we refactor at 7 modules or 10 modules?

### Interesting Ideas

- **Casey's "ceremony vs. logic" framing:** Highlights that the pain point isn't the business logic, it's the boilerplate.
- **Jordan's compounding cost curve:** Quantifies the "defer" cost in concrete terms (4 days → 6 days → 9 days).
- **Alex's "stable foundation" argument:** Valid concern about refactoring during active feature work.

### Open Questions

1. Can we do Option A *and* continue P2 module work in parallel?
2. What's the actual implementation time for Option A? (Alex says 3-4 days, Jordan says 2-3 days)
3. If we defer, what's the decision gate? (After P2 modules? After all Priority 2 features?)

---

## Decision Point

**The question for Round 2:** If we agree to do Option A (generate Decode), should we do it *now* or after P2 modules?

This sets up the next debate: **Path 1 (A→B) vs. Path 2 (A→C)** and the sequencing of changes.
