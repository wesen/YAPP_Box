---
Title: Debate Round 4 — Special Cases: How Do We Retire cutoutFeatureModule?
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
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/04-debate-round-3-developer-ergonomics-which-path-reduces-friction-fastest.md
      Note: Round 3 (Developer ergonomics - chose Path 2)
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/cutouts/module.go
      Note: Cutouts Build returns map[string][][]any
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/features.go
      Note: cutoutFeatureModule special-case
ExternalSources: []
Summary: How do we unify cutouts with other modules in Path 2 (A→C)?
LastUpdated: 2025-11-17
---

# Debate Round 4 — Special Cases: How Do We Retire cutoutFeatureModule?

## Question

**How do we handle cutouts' multi-array output in Path 2 (A→C)?**

**Context:** Cutouts is special because it distributes entries to 6 face arrays:
- `cutoutsFront`, `cutoutsBack`, `cutoutsLeft`, `cutoutsRight`, `cutoutsLid`, `cutoutsBase`

Current system uses `cutoutFeatureModule` (custom `FeatureModule` impl) that bypasses `registry.FeatureModule.Build()` and calls `cutouts.Build()` directly in `Emit()`.

**From Round 3:** We chose Path 2 (A→C) for better ergonomics. Now we need to decide how cutouts fits.

## Pre-Debate Research

### Current Cutouts Architecture

```go
// pkg/yappgen/modules/cutouts/module.go:15
func Build(items []map[string]any) (map[string][][]any, error) {
	byFace := map[string][][]any{
		"cutoutsFront": {},
		"cutoutsBack":  {},
		"cutoutsLeft":  {},
		"cutoutsRight": {},
		"cutoutsLid":   {},
		"cutoutsBase":  {},
	}
	
	for idx, it := range items {
		// ... decode item ...
		// Determine face array
		faceArray, err := faceArrayName(item.Face)
		if err != nil {
			return nil, errors.Wrapf(err, "%s", label)
		}
		byFace[faceArray] = append(byFace[faceArray], params)
	}
	
	return byFace, nil
}
```

**Finding:** `Build()` returns `map[string][][]any`, not `[][]any`. This is why it can't implement `registry.FeatureModule.Build()`.

### Current Registry Stub

```go
// pkg/yappgen/modules/cutouts/registry.go:31-42
func (m *module) Build(items []map[string]any) ([][]any, error) {
	// Note: cutouts return a map by face, not a simple array
	// This is a special case handled in features.go
	byFace, err := Build(items)
	if err != nil {
		return nil, err
	}
	// For now, return empty to satisfy interface
	// The actual cutout handling is in features.go cutoutFeatureModule
	_ = byFace
	return nil, nil
}
```

**Finding:** Registry Build is a stub that returns `nil, nil`. It computes `byFace` but discards it.

### Custom FeatureModule

```go
// pkg/yappgen/features.go:136-181
type cutoutFeatureModule struct{}

func (m *cutoutFeatureModule) Name() string {
	return "cutouts"
}

func (m *cutoutFeatureModule) Collect(resolved map[string]any, features map[string]any, model *Model) error {
	model.Cutouts = nil
	if features == nil {
		return nil
	}
	arr, ok := getArray(features, "cutouts")
	if !ok {
		return nil
	}
	items := normalizeArrayOfMaps(arr)
	model.Cutouts = items
	return nil
}

func (m *cutoutFeatureModule) Emit(ctx context.Context, model *Model, b *strings.Builder) error {
	if len(model.Cutouts) == 0 {
		return nil
	}
	byFace, err := cutouts.Build(model.Cutouts)  // Calls real Build
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

**Finding:** `cutoutFeatureModule` implements `FeatureModule` interface but doesn't use `registry.FeatureModule`. It has custom Collect and Emit logic.

### Lines of Special-Case Code

```bash
$ wc -l pkg/yappgen/features.go
182 pkg/yappgen/features.go

$ grep -A 50 "type cutoutFeatureModule" pkg/yappgen/features.go | wc -l
46
```

**Finding:** 46 lines of special-case code for cutouts in features.go.

---

## Opening Statements

### Six-Faced Friend (`cutouts/`)

*[Defensive but pragmatic]*

Look, I know I'm special. I've heard the complaints. But let me explain *why* I'm special.

**I distribute to 6 arrays because that's how YAPP works.** The OpenSCAD generator expects:

```scad
cutoutsFront = [[10, 20, 30, 40, 5, yappRectangle]];
cutoutsBack = [[15, 25, 35, 45, 5, yappCircle]];
// ... 4 more arrays
```

Each cutout specifies a `face` field, and I route it to the correct array. That's not a quirk—that's the *requirement*.

**Current architecture:**
1. User writes flat list: `cutouts: [{face: front, ...}, {face: back, ...}]`
2. I distribute by face: `{cutoutsFront: [...], cutoutsBack: [...]}`
3. Emission writes 6 separate arrays

**Path 2 (A→C) options:**

**Option 1: Keep custom FeatureModule**
```go
// Model stores typed slice
type Model struct {
	Cutouts []cutouts.CutoutsItem
}

// Custom Emit still calls cutouts.Build
func (m *cutoutFeatureModule) Emit(...) {
	byFace, err := cutouts.Build(model.Cutouts)  // Takes []CutoutsItem now
	// ... emit 6 arrays
}
```

**Option 2: Change Build signature to return map**
```go
// Build signature changes
func Build(items []CutoutsItem) (map[string][][]any, error) {
	// Same logic, just typed input
}

// Emit calls Build and handles map
```

**Option 3: Build returns single array, Emit distributes**
```go
// Build returns flat array with face metadata
func Build(items []CutoutsItem) ([]CutoutRow, error) {
	// CutoutRow has Face field
}

// Emit distributes by face
func (m *cutoutFeatureModule) Emit(...) {
	rows, err := cutouts.Build(model.Cutouts)
	byFace := make(map[string][][]any)
	for _, row := range rows {
		byFace[row.Face] = append(byFace[row.Face], row.Params)
	}
	// ... emit
}
```

**My position:** Option 1 or 2. Don't make me return a single array (Option 3) because face distribution is *my job*. It's domain logic that belongs in the cutouts module, not in features.go.

**Keep the custom FeatureModule.** It's 46 lines of code that handles a legitimate special case. Path 2 doesn't eliminate it—it just makes the input typed.

---

### The Orchestrator (`features.go`)

*[Tired of special-cases]*

Six-Faced Friend, I respect your domain logic argument. But let me show you what your special-case costs me.

**My code has two paths:**

```go
// Path 1: Normal modules (6 modules)
var featureModules = []FeatureModule{
	newArrayFeatureModule("pcb_stands", ...),
	newArrayFeatureModule("connectors", ...),
	// ... 4 more
}

// Path 2: Cutouts (1 module)
newCutoutFeatureModule(),
```

**Every time I add a feature to FeatureModule interface, I have to update:**
1. `arrayFeatureModule` (normal path)
2. `cutoutFeatureModule` (special path)

Example: When we added `ValidateStructure`, I had to add it to both. When we add metrics/logging/tracing, I'll have to add it to both.

**Option 1 (keep custom FeatureModule) perpetuates this.** You're still special, I still have two code paths.

**Option 2 (Build returns map) is interesting.** What if we change the interface?

```go
type FeatureModule interface {
	Schema() ModuleSchema
	Build(items []TypedItem) (BuildResult, error)
}

type BuildResult interface {
	// Single array or multi-array
}

type SingleArray [][]any
type MultiArray map[string][][]any
```

Now `arrayFeatureModule` returns `SingleArray`, cutouts returns `MultiArray`. Emission checks the type and handles accordingly.

**But** that's basically Option D from the analysis (capability interfaces). It's more complex than just accepting your special-case.

**My position:** I'm torn. Option 1 (keep custom FeatureModule) is simplest but perpetuates the split. Option 2 (Build returns map) requires interface changes but unifies the code path.

**Question for Six-Faced Friend:** If we do Option 2, can your Build signature change from `[]map[string]any` to `[]CutoutsItem`? That's the Path 2 requirement.

---

### The Contract (`registry/schema.go`)

*[Principled about interfaces]*

Let me talk about interface design.

**Current interface:**

```go
type FeatureModule interface {
	Schema() ModuleSchema
	Build(items []map[string]any) ([][]any, error)
}
```

This interface makes an assumption: **one module → one array**. That assumption is wrong for cutouts.

**Option 1 (keep custom FeatureModule):**
- Cutouts doesn't implement `registry.FeatureModule`
- Cutouts implements `yappgen.FeatureModule` (different interface)
- Two interfaces with the same name in different packages

That's confusing. When someone says "FeatureModule," which one do they mean?

**Option 2 (Build returns map):**
- Change `registry.FeatureModule.Build` to return `any` or a union type
- Emission checks type and handles accordingly

That's type erasure. We lose compile-time safety. The whole point of Path 2 is *more* type safety, not less.

**Option 3 (Build returns flat array, Emit distributes):**
- Build returns `[][]any` with face metadata
- Emission extracts face and distributes

This keeps the interface clean but moves domain logic (face distribution) out of the module.

**My position:** Option 3. The interface should be simple and uniform. Face distribution is *presentation logic* (how we emit SCAD), not *domain logic* (how we validate/transform cutouts).

**Analogy:** It's like a web API. The controller (Build) returns data. The serializer (Emit) formats it for output. Cutouts shouldn't care about SCAD array names—that's emission's job.

---

## Rebuttals

### Six-Faced Friend (`cutouts/`)

Contract, your analogy is wrong. Face distribution *is* domain logic.

Here's why: The `face` field isn't just metadata for emission. It affects *validation*. Look at the schema:

```yaml
# schema.yaml
fields:
  face:
    type: string
    required: true
    enum: [front, back, left, right, top, lid, bottom, base]
  
  from_face_bottom:
    type: number
    desc: Required for side faces (front, back, left, right)
  
  from_face_back:
    type: number
    desc: Required for horizontal faces (base, lid)
```

**Side faces require `from_face_bottom`.** Horizontal faces require `from_face_back`. That's domain logic that depends on the face value.

If I return a flat array and you distribute in Emit, you've separated the validation logic (in Build) from the distribution logic (in Emit). That's worse coupling, not better.

**Orchestrator's Option 2 is acceptable.** Change my Build signature to:

```go
func Build(items []CutoutsItem) (map[string][][]any, error) {
	// Typed input, map output
}
```

Then update `cutoutFeatureModule.Emit` to call the typed Build. The custom FeatureModule stays, but it's now using typed data.

**My revised position:** Option 2 with typed input. I keep my domain logic (face distribution), you keep your special-case handling (custom FeatureModule), but we get Path 2's type safety.

---

### The Orchestrator (`features.go`)

Six-Faced Friend, I can live with Option 2. Here's what it looks like:

**Before (current):**

```go
type cutoutFeatureModule struct{}

func (m *cutoutFeatureModule) Emit(ctx context.Context, model *Model, b *strings.Builder) error {
	if len(model.Cutouts) == 0 {
		return nil
	}
	byFace, err := cutouts.Build(model.Cutouts)  // []map[string]any → map
	// ... emit 6 arrays
}
```

**After (Path 2 + Option 2):**

```go
type cutoutFeatureModule struct{}

func (m *cutoutFeatureModule) Emit(ctx context.Context, model *Model, b *strings.Builder) error {
	if len(model.Cutouts) == 0 {
		return nil
	}
	byFace, err := cutouts.Build(model.Cutouts)  // []cutouts.CutoutsItem → map
	// ... emit 6 arrays
}
```

**Only change:** Input type. The special-case handling stays, but now it's type-safe.

**Lines of code changed:**
- `cutouts/module.go`: Change Build signature (1 line)
- `cutouts/module.go`: Remove marshal/unmarshal (8 lines deleted)
- `features.go`: No changes (model.Cutouts is already typed)

**My revised position:** Option 2. It's the minimal change that gets us Path 2's type safety without rewriting cutouts' domain logic.

---

### The Contract (`registry/schema.go`)

*[Sighs]*

Fine. I concede that face distribution is domain logic. Six-Faced Friend's validation argument is compelling—the face value affects which fields are required.

**Option 2 it is.** But let me document what we're accepting:

**Trade-off:**
- ✅ Cutouts gets typed input (`[]CutoutsItem`)
- ✅ Domain logic stays in cutouts module
- ❌ Cutouts still has custom FeatureModule (46 lines of special-case)
- ❌ `registry.FeatureModule.Build` interface doesn't apply to cutouts

**Future risk:** If we add another module that needs multi-array output (e.g., `ridgeExt` might have 4 arrays: `ridgeExtFront`, `ridgeExtBack`, `ridgeExtLeft`, `ridgeExtRight`), we'll have *two* custom FeatureModules.

**Mitigation:** Document the pattern. If ridgeExt needs multi-array output, it follows the same pattern as cutouts:
1. Typed Build returns `map[string][][]any`
2. Custom FeatureModule in features.go
3. Emit loops over map and writes arrays

**My revised position:** Option 2, with documentation that this is the pattern for multi-array modules.

---

## Moderator Summary

### Key Arguments

**For Option 1 (keep custom FeatureModule, typed input):**
- Minimal change (just input type)
- Domain logic stays in module
- Proven pattern (already works)

**For Option 2 (Build returns map, typed input):**
- Same as Option 1 (they're equivalent)
- Build signature changes but logic stays the same

**For Option 3 (Build returns flat array, Emit distributes):**
- Clean interface (all modules return `[][]any`)
- Rejected: Face distribution is domain logic, not presentation

### Tensions

1. **Interface purity vs. pragmatism:** Contract wanted uniform interface; Six-Faced Friend showed why it's not feasible.

2. **Special-case handling:** Orchestrator tired of custom code paths; accepted that cutouts legitimately needs it.

3. **Future modules:** What if ridgeExt also needs multi-array output?

### Emerging Consensus

**All three candidates accept Option 2:**
- Cutouts Build takes `[]CutoutsItem` (typed)
- Cutouts Build returns `map[string][][]any` (multi-array)
- Custom `cutoutFeatureModule` stays in features.go
- Pattern documented for future multi-array modules

**Key insight:** Face distribution is domain logic because it affects validation (side faces vs. horizontal faces have different required fields).

### Interesting Ideas

- **Six-Faced Friend's validation argument:** Face value determines required fields, so distribution must stay in module.
- **Orchestrator's minimal change:** Only input type changes, output stays the same.
- **Contract's pattern documentation:** If ridgeExt needs multi-array, follow cutouts' pattern.

### Open Questions

1. Do we need a helper function for multi-array emission? (Cutouts and potential ridgeExt would share it)
2. Should we rename `cutoutFeatureModule` to `multiArrayFeatureModule` to make the pattern clearer?
3. How do we document this pattern in the module authoring guide?

---

## Decision Point

**Consensus: Option 2 (typed input, map output, keep custom FeatureModule)**

**Implementation:**
1. Change `cutouts.Build` signature: `[]CutoutsItem` → `map[string][][]any`
2. Remove marshal/unmarshal boilerplate from cutouts/module.go
3. Keep `cutoutFeatureModule` in features.go (no changes needed)
4. Document multi-array pattern for future modules

**Next debate:** Codegen scope (Question 6) - What needs to be generated for Path 2?
