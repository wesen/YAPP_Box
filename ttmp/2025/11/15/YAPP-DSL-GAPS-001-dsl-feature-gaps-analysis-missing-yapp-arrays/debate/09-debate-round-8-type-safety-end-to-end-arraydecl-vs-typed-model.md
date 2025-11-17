---
Title: Debate Round 8 — Type Safety End-to-End: ArrayDecl vs Typed Model
Ticket: YAPP-DSL-GAPS-001
Status: active
Topics:
    - yapp
    - architecture
    - codegen
    - debate
    - type-safety
DocType: debate
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/06-debate-round-5-codegen-scope-what-to-generate-for-path-1-a-b.md
      Note: Round 5 (Codegen scope)
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/model.go
      Note: Model struct with untyped fields
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/features.go
      Note: Collection and emission logic
ExternalSources: []
Summary: Should we keep []map[string]any in Model or is ArrayDecl IR enough for type safety in Path 1?
LastUpdated: 2025-11-17
---

# Debate Round 8 — Type Safety End-to-End: ArrayDecl vs Typed Model

## Question

**Is ArrayDecl IR (Path 1 Step B) sufficient for type safety, or do we need typed Model fields?**

**Context:** Path 1 (A→B) adds:
- Step A: Generated Decode() (eliminates marshal/unmarshal)
- Step B: ArrayDecl IR (unifies single/multi-array modules)

**But Model still stores `[]map[string]any`:**

```go
type Model struct {
	PcbStands   []map[string]any
	Connectors  []map[string]any
	// ... 5 more untyped fields
}
```

**Question:** Is this "good enough" or should we go further?

## Pre-Debate Research

### Current Type Flow

```
User YAML → Resolver (map[string]any)
         → features.Collect() extracts []any
         → normalizeArrayOfMaps() → []map[string]any
         → Model.PcbStands (stored as []map[string]any)
         → registry.Build() calls Decode()
         → Decode() → []PcbStandsItem (typed!)
         → Build() logic → [][]any
         → ArrayDecl wrapper → []ArrayDecl
         → Emit() → SCAD string
```

**Type boundaries:**
1. ❌ Resolver → Model.PcbStands (untyped)
2. ✅ Model.PcbStands → Decode() → typed struct (Path 1 Step A fixes this)
3. ✅ Build() → ArrayDecl (Path 1 Step B fixes this)
4. ✅ ArrayDecl → Emit() (Path 1 Step B fixes this)

**Remaining untyped boundary:** Resolver → Model storage

### Model Field Definitions

```go
// pkg/yappgen/model.go:37-43
// Features
PcbStands   []map[string]any
Connectors  []map[string]any
BoxMounts   []map[string]any
SnapJoins   []map[string]any
Cutouts     []map[string]any
PushButtons []map[string]any
LightTubes  []map[string]any
```

**Finding:** All 7 feature fields are `[]map[string]any`. No type information at storage level.

### Collection Logic

```go
// pkg/yappgen/features.go:97-119
func (m *arrayFeatureModule) Collect(resolved map[string]any, features map[string]any, model *Model) error {
	ptr := m.field(model)  // func(*Model) *[]map[string]any
	if features == nil {
		*ptr = nil
		return nil
	}
	arr, ok := getArray(features, m.key)
	if !ok {
		*ptr = nil
		return nil
	}
	items := normalizeArrayOfMaps(arr)  // []any → []map[string]any
	*ptr = items
	if m.afterCollect != nil {
		m.afterCollect(model, items)
	}
	return nil
}
```

**Finding:** `normalizeArrayOfMaps()` converts `[]any` to `[]map[string]any`. This is the untyped boundary.

### Field Accessor Pattern

```go
// pkg/yappgen/features.go:28-46
var featureModules = []FeatureModule{
	newArrayFeatureModule("pcb_stands", "pcbStands",
		func(m *Model) *[]map[string]any { return &m.PcbStands },
		pcbstands.Build, nil),
	newArrayFeatureModule("connectors", "connectors",
		func(m *Model) *[]map[string]any { return &m.Connectors },
		connectors.Build, nil),
	// ... 4 more
}
```

**Finding:** Each module has a field accessor that returns `*[]map[string]any`. This signature would need to change for typed Model.

---

## Opening Statements

### The Contract (`registry/schema.go`)

Let me analyze the **type safety gaps** in Path 1 (A→B).

**After Path 1, we have:**

```
✅ Decode() eliminates marshal/unmarshal (type-safe conversion)
✅ ArrayDecl unifies output (type-safe emission)
❌ Model stores []map[string]any (untyped storage)
```

**The remaining gap:**

```go
// In Collect()
arr, ok := getArray(features, "pcb_stands")  // []any
items := normalizeArrayOfMaps(arr)           // []map[string]any
model.PcbStands = items                       // Store untyped
```

**What could go wrong?**

1. **Wrong type in array:**
   ```go
   features["pcb_stands"] = []any{"not a map", 123, true}
   // normalizeArrayOfMaps silently skips non-maps
   // Model.PcbStands = [] (empty, no error)
   ```

2. **Wrong field types:**
   ```go
   features["pcb_stands"] = []any{
       map[string]any{"x": "not a number", "y": 20},
   }
   // Stored in Model, error only caught later in Decode()
   ```

3. **Missing required fields:**
   ```go
   features["pcb_stands"] = []any{
       map[string]any{"x": 10},  // missing required 'y'
   }
   // Stored in Model, error only caught later in Decode()
   ```

**My position:** The untyped Model storage is a **delayed error boundary**. Errors are caught in Decode(), not at storage time.

**Is this acceptable?** Maybe. Decode() happens immediately before Build(), so the delay is minimal. But it's still a gap.

**Alternative:** Call Decode() during Collect() and store typed data:

```go
func (m *arrayFeatureModule) Collect(...) error {
	arr, ok := getArray(features, m.key)
	if !ok {
		return nil
	}
	items := normalizeArrayOfMaps(arr)
	
	// Decode immediately
	typed, err := m.decoder(items)  // func([]map[string]any) ([]TypedItem, error)
	if err != nil {
		return errors.Wrapf(err, "decode %s", m.key)
	}
	
	// Store typed (but how? Model field is []map[string]any)
	// ... this requires changing Model field types
}
```

**Problem:** This requires typed Model fields, which is Path 2 (A→C), not Path 1.

---

### Alex Chen — "The Pragmatist"

Contract, you're overthinking this.

**Path 1 (A→B) gives us:**
- ✅ No more marshal/unmarshal (Decode() handles it)
- ✅ Unified output (ArrayDecl)
- ✅ Errors caught in Decode() (before Build())

**The "gap" you're worried about:**
- Model stores `[]map[string]any`
- Errors caught in Decode() instead of Collect()

**So what?** The error is still caught. It's caught 10 lines later in the call stack instead of immediately, but it's caught.

**Let me show you the actual flow:**

```go
// BuildModel() in model.go
func BuildModel(ctx context.Context, resolved map[string]any) (*Model, error) {
	m := &Model{}
	// ... populate globals ...
	
	features, _ := getMap(resolved, "features")
	if err := collectFeatureModules(resolved, features, m); err != nil {
		return nil, err  // Collect errors
	}
	
	return m, nil
}

// EmitSCAD() in emit.go
func EmitSCAD(ctx context.Context, m *Model) ([]byte, error) {
	var b strings.Builder
	// ... emit globals ...
	
	if err := emitFeatureModules(ctx, m, &b); err != nil {
		return nil, err  // Decode/Build errors
	}
	
	return []byte(b.String()), nil
}
```

**Errors are caught before SCAD emission.** The user never sees invalid SCAD because Decode() fails first.

**My position:** Path 1 (A→B) is sufficient. The "untyped Model" concern is theoretical, not practical.

**Cost of typed Model:**
- Change 7 Model field types
- Change field accessor signatures (6 call sites)
- Change Collect() to call Decode()
- Update all tests that construct Model

**Benefit:** Errors caught 10 lines earlier in the call stack.

**ROI:** Not worth it. Path 1 (A→B) solves the real problems (marshal/unmarshal boilerplate, ArrayDecl unification). Typed Model is perfectionism.

---

### Sam Park — "The Codegen Maintainer"

Let me talk about the **developer experience** of the untyped Model.

**Scenario:** I'm debugging a SCAD generation issue. I set a breakpoint in Build():

```go
func Build(items []map[string]any) ([]ArrayDecl, error) {
	typed, err := Decode(items)  // ← Breakpoint here
	if err != nil {
		return nil, err
	}
	// ...
}
```

**What do I see in the debugger?**

```
items = []map[string]any{
	{"x": 10.0, "y": 20.0, "height": 5.0},
	{"x": 15.0, "y": 25.0, "height": 7.0},
}
```

**Problem:** I have to mentally map `map[string]any` to the struct. I can't hover over `items[0]` and see "PcbStandsItem{X: 10.0, Y: 20.0, Height: 5.0}".

**With typed Model:**

```go
func Build(items []PcbStandsItem) ([]ArrayDecl, error) {
	// ← Breakpoint here
}
```

**Debugger shows:**

```
items = []PcbStandsItem{
	{X: 10.0, Y: 20.0, Height: 5.0},
	{X: 15.0, Y: 25.0, Height: 7.0},
}
```

**Much clearer.** I can see field names, types, and values.

**But** Alex is right that the cost is high. Changing Model field types touches a lot of code.

**My position:** Path 1 (A→B) is acceptable for now. Typed Model is a "nice to have" for debugging, not a "must have" for correctness.

**Compromise:** Add a debug helper:

```go
// pkg/yappgen/debug/helpers.go
func DumpModel(m *Model) {
	fmt.Printf("PcbStands: %d items\n", len(m.PcbStands))
	for i, item := range m.PcbStands {
		typed, _ := pcbstands.Decode([]map[string]any{item})
		fmt.Printf("  [%d]: %+v\n", i, typed[0])
	}
	// ... repeat for other features
}
```

**This gives us typed debugging without changing Model.**

---

## Rebuttals

### The Contract (`registry/schema.go`)

Alex and Sam, I concede that Path 1 (A→B) is *sufficient* for correctness. Errors are caught in Decode(), and that's good enough.

**But** let me make one more argument for typed Model: **interface clarity**.

**Current interface (Path 1):**

```go
type FeatureModule interface {
	Name() string
	Collect(resolved map[string]any, features map[string]any, model *Model) error
	Emit(ctx context.Context, model *Model, b *strings.Builder) error
}
```

**What does Collect() do?** You have to read the implementation to know it extracts `[]map[string]any` and stores it in Model.

**What does Emit() do?** You have to read the implementation to know it calls Decode() → Build() → emits ArrayDecls.

**With typed Model (hypothetical):**

```go
type FeatureModule[T any] interface {
	Name() string
	Collect(resolved map[string]any, features map[string]any) ([]T, error)
	Emit(ctx context.Context, items []T, b *strings.Builder) error
}
```

**Now the interface documents itself:**
- Collect() returns typed items
- Emit() takes typed items

**But** this requires generics, which adds complexity. And it's a bigger change than just typed Model fields.

**My revised position:** Path 1 (A→B) is acceptable. Typed Model would be better for clarity, but the cost/benefit doesn't justify it.

**Recommendation:** Document the type flow clearly so future maintainers understand where type boundaries are.

---

### Alex Chen — "The Pragmatist"

Contract, your generics proposal is interesting but way out of scope.

**Let me summarize the trade-offs:**

**Path 1 (A→B) - What we're doing:**
- ✅ Generate Decode() (eliminates marshal/unmarshal)
- ✅ Add ArrayDecl IR (unifies output)
- ✅ Model stores `[]map[string]any` (untyped)
- ✅ Errors caught in Decode() (before Build())

**Cost:** ~70 lines of code changes (7 modules + templates)
**Benefit:** Cleaner builders, unified output, no marshal/unmarshal

**Path 2 (A→C) - Typed Model (rejected):**
- ✅ Generate Decode()
- ✅ Model stores typed slices
- ✅ Errors caught in Collect() (earlier)

**Cost:** ~120 lines of code changes (7 modules + Model + features.go + tests)
**Benefit:** Type safety end-to-end, better debugging

**We chose Path 1 because:**
1. Lower cost (70 vs 120 lines)
2. Sufficient correctness (errors still caught)
3. Simpler migration (no Model field changes)

**The "untyped Model" concern is valid but not critical.** We're not sacrificing correctness, just moving the error boundary slightly later.

**My final position:** Path 1 (A→B) is the right choice. Typed Model is a future enhancement if we find the debugging experience is painful.

---

### Sam Park — "The Codegen Maintainer"

I agree with Alex. Let me add one more point: **testing**.

**With Path 1 (A→B), tests are simple:**

```go
func TestBuild_ValidInput(t *testing.T) {
	items := []map[string]any{
		{"x": 10.0, "y": 20.0},
	}
	
	decls, err := Build(items)
	// ... assertions
}
```

**With typed Model, tests need to construct typed structs:**

```go
func TestEmit_ValidInput(t *testing.T) {
	model := &Model{
		PcbStands: []pcbstands.PcbStandsItem{
			{X: 10.0, Y: 20.0},
		},
	}
	
	var b strings.Builder
	err := emitFeatureModules(ctx, model, &b)
	// ... assertions
}
```

**Typed structs are more verbose.** For simple tests, `map[string]any` is actually easier.

**But** for complex tests (nested objects, optional fields), typed structs are clearer.

**Trade-off:** Path 1 has easier test construction, Path 2 has clearer test intent.

**My final position:** Path 1 (A→B) is sufficient. We can always add typed Model later if we find the untyped storage is causing problems.

---

## Moderator Summary

### Key Arguments

**For Path 1 (A→B) - Keep untyped Model:**
- ✅ Errors still caught (in Decode(), before Build())
- ✅ Lower migration cost (70 vs 120 lines)
- ✅ Simpler test construction (map[string]any)
- ✅ Sufficient for correctness

**For Typed Model (Path 2 A→C) - Rejected:**
- ✅ Errors caught earlier (in Collect())
- ✅ Better debugging (typed values in debugger)
- ✅ Clearer interface (self-documenting types)
- ❌ Higher migration cost
- ❌ More complex tests

### Consensus Points

All three candidates agree:
- Path 1 (A→B) is sufficient for correctness
- Typed Model would be better for debugging/clarity
- The cost/benefit doesn't justify typed Model now
- Can revisit typed Model as future enhancement

### Type Safety Analysis

**Path 1 (A→B) type boundaries:**

```
User YAML → Resolver
         ↓ (untyped: map[string]any)
         → Model.PcbStands
         ↓ (untyped: []map[string]any)
         → Decode()
         ↓ (TYPED: []PcbStandsItem)
         → Build()
         ↓ (typed: [][]any)
         → ArrayDecl
         ↓ (typed: []ArrayDecl)
         → Emit()
         ↓ (typed: string)
         → SCAD output
```

**Untyped segment:** Resolver → Model → Decode()
**Typed segment:** Decode() → Emit()

**Error detection:** Decode() catches type errors before Build()

### Implementation Decision

**Consensus: Path 1 (A→B) with untyped Model is sufficient**

**Rationale:**
1. Errors are caught before SCAD emission (correctness preserved)
2. Lower migration cost (70 vs 120 lines)
3. Simpler implementation (no Model field changes)
4. Can add typed Model later if needed

**Trade-offs accepted:**
- ❌ Errors caught in Decode() instead of Collect() (10 lines later)
- ❌ Debugging shows map[string]any instead of typed structs
- ❌ Model fields don't document expected types

**Mitigations:**
- Document type flow clearly
- Add debug helpers (DumpModel) if needed
- Consider typed Model as future enhancement

### Open Questions

1. Should we add debug helpers for typed inspection?
2. Do we document the type boundaries in code comments?
3. Is there a middle ground (typed Model but keep untyped Collect)?

---

## Decision Point

**Consensus: Path 1 (A→B) with untyped Model**

**Key insight:** Type safety at the Build/Emit boundary (via Decode() + ArrayDecl) is sufficient. End-to-end type safety (typed Model) is a "nice to have," not a requirement.

**Implementation:**
1. Generate Decode() (Step A)
2. Add ArrayDecl IR (Step B)
3. Keep Model fields as `[]map[string]any`
4. Document type flow in code comments

**Future consideration:** If debugging untyped Model becomes painful, revisit typed Model fields.

**Next steps:** Synthesize all debate rounds into final implementation plan.
