---
Title: Debate Round 2 — Path Choice: A→B vs A→C
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
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/02-debate-round-1-go-no-go-builder-contract-refactor-urgency.md
      Note: Round 1 (Go/No-Go)
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/analysis/04-2025-11-17-builder-contract-and-codegen-options.md
      Note: Source analysis with Options A-E
ExternalSources: []
Summary: Which migration path delivers more value sooner: incremental (A→B) or comprehensive (A→C)?
LastUpdated: 2025-11-17
---

# Debate Round 2 — Path Choice: A→B vs A→C

## Question

**Which path should we take: Path 1 (A→B) or Path 2 (A→C)?**

**Path 1 (A→B):**
- **Step A:** Generate `Decode()` per module, keep `[][]any` output
- **Step B:** Add `ArrayDecl` IR to handle multi-array modules (cutouts)

**Path 2 (A→C [+D]):**
- **Step A:** Generate `Decode()` per module
- **Step C:** Store typed slices in `Model`, typed collectors
- **Step D (optional):** Capability interfaces for special outputs

## Pre-Debate Research

### Current Interface Contract

```go
// pkg/registry/schema.go:51-58
type FeatureModule interface {
	Schema() ModuleSchema
	Build(items []map[string]any) ([][]any, error)
}
```

**Finding:** Interface is simple, but `[][]any` return type can't express cutouts' multi-array output.

### Cutouts Special-Case

```go
// pkg/yappgen/features.go:136-181
type cutoutFeatureModule struct{}

func (m *cutoutFeatureModule) Emit(ctx context.Context, model *Model, b *strings.Builder) error {
	if len(model.Cutouts) == 0 {
		return nil
	}
	byFace, err := cutouts.Build(model.Cutouts)  // Returns map[string][][]any
	if err != nil {
		return errors.Wrap(err, "cutouts")
	}
	keys := make([]string, 0, len(byFace))
	for k := range byFace {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if len(byFace[k]) == 0 {
			continue
		}
		writeArrayDecl(b, k, byFace[k])
		b.WriteString("\n")
	}
	return nil
}
```

**Finding:** Cutouts bypasses `registry.FeatureModule.Build()` entirely. Has custom `FeatureModule` impl that calls `cutouts.Build()` directly in `Emit()`.

### Model Field Types

```bash
$ grep -A 1 "// Features" pkg/yappgen/model.go
	// Features
	PcbStands   []map[string]any
	Connectors  []map[string]any
	BoxMounts   []map[string]any
	SnapJoins   []map[string]any
	Cutouts     []map[string]any
	PushButtons []map[string]any
	LightTubes  []map[string]any
```

**Finding:** All features stored as `[]map[string]any`. No type information preserved from collection to emission.

### ArrayFeatureModule Usage

```bash
$ grep -A 5 "newArrayFeatureModule" pkg/yappgen/features.go | head -30
var featureModules = []FeatureModule{
	newArrayFeatureModule("pcb_stands", "pcbStands",
		func(m *Model) *[]map[string]any { return &m.PcbStands },
		pcbstands.Build, nil),
	newArrayFeatureModule("connectors", "connectors",
		func(m *Model) *[]map[string]any { return &m.Connectors },
		connectors.Build, nil),
	newArrayFeatureModule("box_mounts", "boxMounts",
		func(m *Model) *[]map[string]any { return &m.BoxMounts },
		boxmounts.Build, nil),
	newArrayFeatureModule("push_buttons", "pushButtons",
		func(m *Model) *[]map[string]any { return &m.PushButtons },
		pushbuttons.Build,
		func(m *Model, items []map[string]any) {
			m.PrintSwitchExtenders = len(items) > 0
		}),
	newArrayFeatureModule("snap_joins", "snapJoins",
		func(m *Model) *[]map[string]any { return &m.SnapJoins },
		snapjoins.Build, nil),
	newArrayFeatureModule("light_tubes", "lightTubes",
		func(m *Model) *[]map[string]any { return &m.LightTubes },
		lighttubes.Build, nil),
	newCutoutFeatureModule(),
}
```

**Finding:** 6 modules use `newArrayFeatureModule` helper. Cutouts uses `newCutoutFeatureModule`. All pass `func(m *Model) *[]map[string]any` accessor.

### Schemagen Template Complexity

```bash
$ wc -l pkg/schemagen/templates/*.tmpl
  23 pkg/schemagen/templates/schema_gen.go.tmpl
  20 pkg/schemagen/templates/schema_gen_test.go.tmpl
  14 pkg/schemagen/templates/modules_gen.go.tmpl
 165 pkg/schemagen/templates/schema_validate.go.tmpl
 222 total
```

**Finding:** Current templates are 222 lines total. `schema_validate.go.tmpl` is the most complex (165 lines).

---

## Opening Statements

### Jordan Rivera — "The Architect"

I'm going to make a controversial argument: **Path 1 (A→B) is a dead end.**

Look at the research. We have 6 modules using `newArrayFeatureModule` and 1 module (cutouts) using a custom `FeatureModule` impl. Path 1's Option B adds `ArrayDecl` IR to unify them. But here's the problem:

**Option B doesn't actually unify anything.** It just moves the special-case from `features.go` into the IR layer. You still have:

```go
// Hypothetical Option B
type ArrayDecl struct {
	Name string
	Rows [][]any
}

func (m *module) Build(items []map[string]any) ([]ArrayDecl, error) {
	// Cutouts returns 6 ArrayDecls
	// Everyone else returns 1 ArrayDecl
}
```

So now instead of `cutoutFeatureModule` being special, `cutouts.Build()` returns a different *number* of `ArrayDecl` objects. You've moved the special-case, not eliminated it.

**Path 2 (A→C) actually solves the problem.** Look at what happens when Model stores typed slices:

```go
// Path 2 Model
type Model struct {
	PcbStands   []pcbstands.PcbStandsItem
	Connectors  []connectors.ConnectorsItem
	Cutouts     []cutouts.CutoutsItem  // Still a single slice!
	// ...
}
```

Cutouts' special behavior (distributing to 6 faces) becomes an *implementation detail* of `cutouts.Build()`. The Model doesn't care. The registry doesn't care. Only `cutouts.Build()` and its emission logic care.

**My position:** Path 2 (A→C). It's more work upfront (change Model fields, update collectors), but it eliminates the special-case entirely instead of shuffling it around.

**Migration cost:** Yes, Path 2 touches more files (Model + 7 modules + features.go). But it's a one-time cost that pays off forever. Path 1 leaves us with `ArrayDecl` IR that's still accommodating special-cases.

---

### Alex Chen — "The Pragmatist"

Jordan, I respect the vision, but let's talk about *risk*.

**Path 1 (A→B) is two small steps:**
- Step A: Generate Decode → touches 7 module files + schemagen templates
- Step B: Add ArrayDecl IR → touches features.go + cutouts emission

Each step is independently testable. Each step can be reverted without cascading failures.

**Path 2 (A→C) is a big bang:**
- Change Model field types (7 fields)
- Change `newArrayFeatureModule` signature (6 call sites)
- Update collectors (7 modules)
- Update emission (features.go)
- Update all tests that construct Model

That's 20+ files changing in one commit. If something breaks, the revert is painful.

I ran a quick analysis of what Path 2 touches:

```bash
$ grep -r "Model struct" pkg/yappgen/
pkg/yappgen/model.go:type Model struct {

$ grep -r "newArrayFeatureModule" pkg/yappgen/
pkg/yappgen/features.go:	newArrayFeatureModule("pcb_stands", ...
pkg/yappgen/features.go:	newArrayFeatureModule("connectors", ...
pkg/yappgen/features.go:	newArrayFeatureModule("box_mounts", ...
pkg/yappgen/features.go:	newArrayFeatureModule("push_buttons", ...
pkg/yappgen/features.go:	newArrayFeatureModule("snap_joins", ...
pkg/yappgen/features.go:	newArrayFeatureModule("light_tubes", ...

$ grep -r "Model{" pkg/yappgen/ --count
pkg/yappgen/yappgen_test.go:3
```

**Finding:** 6 call sites to update, 3+ test constructions to fix. And that's just what I found in 5 minutes.

**My position:** Path 1 (A→B). Do Step A first (generate Decode), ship it, validate it works. *Then* evaluate Step B. If Step B turns out to be complex, we can pivot to Path 2. But we don't commit to Path 2's big bang without proof that Step A works.

**Incremental delivery:** Path 1 gives us value after Step A. Path 2 gives us value only after A+C are both done. That's 2x the risk exposure.

---

### Sam Park — "The Codegen Maintainer"

*[Opens laptop, pulls up templates]*

Let me talk about the codegen implications, because that's where the rubber meets the road.

**Option A (Decode generation) is straightforward:**

```go
// Generated Decode function (pseudocode)
func Decode(items []map[string]any) ([]PcbStandsItem, error) {
	out := make([]PcbStandsItem, 0, len(items))
	for idx, it := range items {
		// Field extraction with type assertions
		x, ok := it["x"].(float64)
		if !ok {
			return nil, fmt.Errorf("pcb_stands[%d].x: expected number", idx)
		}
		// ... more fields ...
		out = append(out, PcbStandsItem{X: x, Y: y, ...})
	}
	return out, nil
}
```

Template complexity: ~50 lines. I can generate this from the existing `SchemaField` metadata.

**Path 1 (A→B) adds ArrayDecl IR:**

```go
// Generated Build wrapper (pseudocode)
func (m *module) Build(items []map[string]any) ([]ArrayDecl, error) {
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}
	rows, err := Build(typed)  // Hand-written logic
	if err != nil {
		return nil, err
	}
	return []ArrayDecl{{Name: "pcbStands", Rows: rows}}, nil
}
```

Template complexity: ~30 lines. Doable.

**Path 2 (A→C) changes Model fields:**

Now I need to generate:
1. Decode function (same as Option A)
2. Typed collector wrapper
3. Model field type (already generated as `PcbStandsItem`)

But here's the catch: **I can't generate the Model struct itself** because it has non-feature fields (ProjectName, PcbLength, etc.). So Model becomes a hybrid: some fields hand-written, some fields typed by generated structs.

That's not a blocker, but it's *weird*. It means Model.go has imports like:

```go
import (
	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pcbstands"
	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/connectors"
	// ... 7 module imports
)
```

And if a module changes its struct name, Model.go breaks. That's a new dependency edge.

**My position:** Path 1 (A→B) is easier to generate. Path 2 (A→C) is doable but creates import coupling between Model and all modules.

**Codegen complexity:**
- Path 1: +80 lines of templates
- Path 2: +50 lines of templates + manual Model.go updates

---

## Rebuttals

### Jordan Rivera — "The Architect"

Alex, I hear your "big bang" concern, but let's be precise about what Path 2 actually changes.

**Path 2 does NOT require a single atomic commit.** We can stage it:

1. **Phase 1:** Generate Decode (Option A) — 7 modules, schemagen templates
2. **Phase 2:** Update Model fields one module at a time
   - Change `Model.PcbStands` to typed slice
   - Update `newArrayFeatureModule` call for pcbstands
   - Update tests
   - Ship, validate
3. **Phase 3:** Repeat for remaining 6 modules

Each phase is independently testable. If Phase 2 breaks pcbstands, we revert *one module*, not the whole system.

**Your "20+ files" number is misleading.** Yes, Path 2 eventually touches 20+ files, but so does Path 1 if you count Step A + Step B. The difference is *sequencing*, not total scope.

Sam, you're right about the import coupling. But that's actually a *feature*, not a bug. Right now, if I change `PcbStandsItem` struct, nothing breaks at compile time. With Path 2, if I change the struct, Model.go fails to compile. **That's type safety working as intended.**

**The real question:** Do we want compile-time checks between modules and Model? I say yes. You're saying it's "weird" to have those imports, but it's only weird because we've been living in untyped map land for so long.

**Path 1's ArrayDecl IR is a half-measure.** It unifies the emission interface but doesn't eliminate the special-case. Cutouts still returns multiple ArrayDecls. You've just moved the conditional logic from `features.go` to the emission loop.

---

### Alex Chen — "The Pragmatist"

Jordan, your "phase 2" plan sounds great in theory, but let's talk about *developer experience during migration*.

If we do Path 2 phased (one module at a time), we have a **mixed state** for weeks:

```go
type Model struct {
	PcbStands   []pcbstands.PcbStandsItem  // Migrated
	Connectors  []map[string]any            // Not migrated yet
	BoxMounts   []map[string]any            // Not migrated yet
	// ...
}
```

Now `newArrayFeatureModule` has to handle *both* typed slices and untyped maps. That's two code paths, not one. You've increased complexity during migration, not reduced it.

**Path 1 doesn't have this problem.** Step A (Decode) is additive—it doesn't change existing interfaces. Step B (ArrayDecl) is also additive—it extends the return type but doesn't break existing callers.

Sam's point about import coupling is critical. You're adding 7 new import edges to Model.go. That's 7 new ways for circular dependencies to creep in. Right now, Model.go imports nothing from modules/. That's a clean boundary.

**My revised position:** Do Step A (Decode) now. Evaluate Step B vs. Path 2 after we see how Decode works in practice. Don't commit to Path 2's import coupling until we have evidence it's worth it.

---

### Sam Park — "The Codegen Maintainer"

*[Looks at both proposals]*

Okay, I need to correct something I said earlier. I claimed Path 2 is harder to generate, but I was wrong.

I just sketched out the templates. **Path 2 is actually *simpler* for codegen** because I don't have to generate the Build wrapper. Here's why:

**Path 1 (A→B) requires:**
1. Generate Decode
2. Generate Build wrapper that calls Decode + hand-written Build + wraps in ArrayDecl
3. Generate emission loop that unwraps ArrayDecl

**Path 2 (A→C) requires:**
1. Generate Decode
2. That's it. Hand-written Build stays as-is.

The complexity I was worried about (Model.go imports) isn't codegen complexity—it's *architecture* complexity. And Jordan's right: that's type safety, not a bug.

**My revised position:** From a codegen perspective, Path 2 is cleaner. I generate Decode, module authors use it in Build, done. No IR layer, no wrappers, no special-case emission.

**But** I agree with Alex on migration risk. If we do Path 2, we need the phased rollout Jordan described. And we need to solve the "mixed state" problem Alex raised.

**Proposal:** Can we do Path 1 Step A (Decode) first, then decide between Step B and Path 2? That gives us data on how Decode works before committing to the bigger change.

---

## Moderator Summary

### Key Arguments

**For Path 1 (A→B):**
- Incremental: two small steps, independently testable
- Lower risk: each step can be reverted without cascading failures
- No import coupling: Model.go stays independent of modules
- Immediate value: Step A delivers benefits before Step B

**For Path 2 (A→C):**
- Eliminates special-case: cutouts becomes a normal module
- Type safety: compile-time checks between Model and modules
- Simpler codegen: no IR layer, no wrappers
- Long-term maintainability: proper type boundaries

### Tensions

1. **Incremental vs. Comprehensive:** Alex wants small steps; Jordan wants to solve the root problem.

2. **Import coupling:** Alex sees it as a new risk; Jordan sees it as type safety.

3. **Migration state:** Alex worries about mixed typed/untyped Model during Path 2 rollout.

4. **ArrayDecl IR value:** Jordan says it's a half-measure; Alex says it's a useful abstraction.

### Emerging Consensus

All three candidates agree:
- Option A (Decode) is the right first step
- Current `[]map[string]any` → `[][]any` flow is suboptimal
- Cutouts special-case should be addressed

**The real debate:** Do we commit to Path 2's type safety now, or keep options open with Path 1?

### Interesting Ideas

- **Jordan's phased Path 2:** Migrate Model fields one module at a time
- **Sam's codegen insight:** Path 2 is actually simpler to generate than Path 1+IR
- **Alex's "mixed state" concern:** Valid UX issue during migration

### Open Questions

1. Can we solve the "mixed state" problem for Path 2 phased rollout?
2. Is ArrayDecl IR (Path 1 Step B) actually useful, or just moving special-cases around?
3. Should we do Step A first and defer the Path 1 vs. Path 2 decision?

---

## Decision Point

**Emerging recommendation:** Do Option A (Decode) first, then evaluate.

**Why:**
- All candidates agree Decode is valuable
- Gives us real-world data on how generated code works
- Doesn't commit to Path 1 or Path 2 yet
- Can ship in 2-3 days

**Next debate:** If we do Option A first, how do we handle the "mixed state" during migration? (Relevant for both Path 1 and Path 2)
