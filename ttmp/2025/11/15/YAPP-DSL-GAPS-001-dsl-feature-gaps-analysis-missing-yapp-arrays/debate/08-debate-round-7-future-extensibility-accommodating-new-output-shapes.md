---
Title: Debate Round 7 — Future Extensibility: Accommodating New Output Shapes
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
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/07-debate-round-6-enforceability-keeping-modules-honest-post-refactor.md
      Note: Round 6 (Enforceability)
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/reference/01-yapp-feature-array-inventory.md
      Note: Complete YAPP feature inventory
ExternalSources: []
Summary: Does Path 1 (A→B) with ArrayDecl IR handle unknown future module needs?
LastUpdated: 2025-11-17
---

# Debate Round 7 — Future Extensibility: Accommodating New Output Shapes

## Question

**Does Path 1 (A→B) with ArrayDecl IR handle unknown future module needs?**

**Context:** We're implementing ArrayDecl to unify single-array and multi-array modules. But what about future modules with different output patterns?

**Known future modules:**
- `labelsPlane` - Single array (text labels)
- `ridgeExt` - Four arrays (ridgeExtFront/Back/Left/Right)
- `displayMounts` - Single array (LCD/OLED mounting)

## Pre-Debate Research

### YAPP Feature Array Inventory

From `YAPP_Template_v3.scad` and feature inventory:

**Currently implemented (7 modules):**
- ✅ pcbStands (single array)
- ✅ connectors (single array)
- ✅ snapJoins (single array)
- ✅ cutouts (6 arrays: Front/Back/Left/Right/Lid/Base)
- ✅ pushButtons (single array)
- ✅ boxMounts (single array)
- ✅ lightTubes (single array)

**Remaining Priority 2 (3 modules):**
- ❌ labelsPlane (single array)
- ❌ ridgeExt (4 arrays: ridgeExtFront/Back/Left/Right)
- ❌ Cutout masks (extension to cutouts, not a new module)

**Remaining Priority 3 (2 modules):**
- ❌ displayMounts (single array, 17 parameters)
- ❌ Multi-PCB support (not a feature array, structural change)

### RidgeExt Structure

```scad
// YAPP_Template_v3.scad:842-856
ridgeExtFront =
[
];

ridgeExtBack =
[
];

ridgeExtLeft =
[
];

ridgeExtRight =
[
];
```

**Finding:** RidgeExt is 4 separate arrays, similar to cutouts' 6 arrays. Each entry has:
- p(0) = pos (along edge)
- p(1) = width
- p(2) = height
- Flags: coordinate, origin, yappPCBName

**Pattern:** Face-based distribution like cutouts.

### LabelsPlane Structure

```scad
// YAPP_Template_v3.scad:797-799
labelsPlane =
[
];
```

**Finding:** Single array. Each entry has 12 parameters:
- Position, text, font, size, depth, face, rotation, alignment, direction, spacing

**Pattern:** Standard single-array module.

### DisplayMounts Structure

```scad
// YAPP_Template_v3.scad:894-896
displayMounts =
[
];
```

**Finding:** Single array. Each entry has 17 parameters (most complex module):
- Position, dimensions, window cutout, pin mounts, bevels, flags

**Pattern:** Standard single-array module (just more parameters).

### Current ArrayDecl Design

From Round 5 (Codegen scope):

```go
type ArrayDecl struct {
	Name string
	Rows [][]any
}

// Single-array modules return 1 ArrayDecl
func (m *module) Build(items []map[string]any) ([]ArrayDecl, error) {
	typed, err := Decode(items)
	rows, err := Build(typed)
	return []ArrayDecl{{Name: "pcbStands", Rows: rows}}, nil
}

// Multi-array modules return multiple ArrayDecls
func (m *module) Build(items []map[string]any) ([]ArrayDecl, error) {
	typed, err := Decode(items)
	byFace, err := Build(typed)
	var decls []ArrayDecl
	for name, rows := range byFace {
		decls = append(decls, ArrayDecl{Name: name, Rows: rows})
	}
	return decls, nil
}
```

---

## Opening Statements

### Jordan Rivera — "The Architect"

Let me analyze the future module landscape.

**Remaining modules by output pattern:**

1. **Single-array (8 total):**
   - Implemented: pcbStands, connectors, snapJoins, pushButtons, boxMounts, lightTubes (6)
   - Future: labelsPlane, displayMounts (2)

2. **Multi-array (2 total):**
   - Implemented: cutouts (6 arrays)
   - Future: ridgeExt (4 arrays)

**Pattern:** 80% single-array, 20% multi-array.

**ArrayDecl handles both:**

```go
// labelsPlane (single array)
return []ArrayDecl{{Name: "labelsPlane", Rows: rows}}, nil

// ridgeExt (4 arrays)
return []ArrayDecl{
	{Name: "ridgeExtFront", Rows: frontRows},
	{Name: "ridgeExtBack", Rows: backRows},
	{Name: "ridgeExtLeft", Rows: leftRows},
	{Name: "ridgeExtRight", Rows: rightRows},
}, nil
```

**My position:** ArrayDecl is sufficient for all known future modules. We have 10 total modules to implement, and ArrayDecl handles both patterns.

**But** let me play devil's advocate: **What about unknown unknowns?**

**Hypothetical future needs:**

1. **Conditional arrays:** Module returns different arrays based on config
   ```go
   if model.UseAdvancedFeature {
       return []ArrayDecl{{Name: "advanced", Rows: advRows}}, nil
   }
   return []ArrayDecl{{Name: "basic", Rows: basicRows}}, nil
   ```
   ✅ ArrayDecl handles this

2. **Dynamic array names:** Array name depends on data
   ```go
   for pcbName, rows := range byPCB {
       decls = append(decls, ArrayDecl{Name: pcbName + "Stands", Rows: rows})
   }
   ```
   ✅ ArrayDecl handles this

3. **Non-array output:** Module needs to emit something other than arrays
   ```go
   // Hypothetical: module sets a global variable instead of array
   return []ArrayDecl{}, nil  // Empty, but still valid
   ```
   ✅ ArrayDecl handles this (empty slice)

**Conclusion:** ArrayDecl is flexible enough for foreseeable needs.

---

### Six-Faced Friend (`cutouts/`)

*[Speaks from experience]*

Jordan's analysis is good, but let me talk about **ridgeExt specifically** because it's the next multi-array module.

**RidgeExt is like cutouts:**
- User writes flat list: `ridge_ext: [{face: front, ...}, {face: back, ...}]`
- Module distributes by face: `{ridgeExtFront: [...], ridgeExtBack: [...]}`
- Emission writes 4 separate arrays

**Key difference:** RidgeExt has 4 faces (sides), cutouts has 6 faces (sides + top/bottom).

**Implementation with ArrayDecl:**

```go
// pkg/yappgen/modules/ridgeext/module.go
func Build(items []map[string]any) (map[string][][]any, error) {
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}
	
	byFace := map[string][][]any{
		"ridgeExtFront": {},
		"ridgeExtBack":  {},
		"ridgeExtLeft":  {},
		"ridgeExtRight": {},
	}
	
	for _, item := range typed {
		face := faceArrayName(item.Face)  // front → ridgeExtFront
		params := []any{item.Pos, item.Width, item.Height, ...}
		byFace[face] = append(byFace[face], params)
	}
	
	return byFace, nil
}

// pkg/yappgen/modules/ridgeext/registry.go
func (m *module) Build(items []map[string]any) ([]ArrayDecl, error) {
	typed, err := Decode(items)
	byFace, err := Build(typed)
	
	var decls []ArrayDecl
	for name, rows := range byFace {
		if len(rows) > 0 {  // Only emit non-empty arrays
			decls = append(decls, ArrayDecl{Name: name, Rows: rows})
		}
	}
	return decls, nil
}
```

**My position:** ArrayDecl works perfectly for ridgeExt. The pattern is identical to cutouts, just with 4 faces instead of 6.

**Recommendation:** When we implement ridgeExt, copy the cutouts pattern:
1. `module.go` Build returns `map[string][][]any`
2. `registry.go` Build wraps in ArrayDecl
3. Custom `ridgeExtFeatureModule` in features.go (or reuse cutouts' pattern)

**Question:** Should we create a shared `multiArrayFeatureModule` helper that both cutouts and ridgeExt can use?

---

### The Generator (`schemagen/`)

Let me talk about **codegen implications for future modules**.

**Current generation (Path 1 A→B):**
- Generate Decode() in schema_gen.go
- Generate ArrayDecl wrapper in registry.go
- Module author writes Build() logic

**For labelsPlane (12 parameters):**

```go
// Generated Decode()
func Decode(items []map[string]any) ([]LabelsPlaneItem, error) {
	// ... 12 field extractions
}

// Generated ArrayDecl wrapper
func (m *module) Build(items []map[string]any) ([]ArrayDecl, error) {
	typed, err := Decode(items)
	rows, err := Build(typed)
	return []ArrayDecl{{Name: "labelsPlane", Rows: rows}}, nil
}

// Hand-written Build()
func Build(items []LabelsPlaneItem) ([][]any, error) {
	for _, item := range items {
		params := []any{
			item.PosX, item.PosY, item.Text, item.Font, item.Size,
			item.Depth, item.Face, item.Rotation, item.Alignment,
			item.Direction, item.Spacing, /* flags */
		}
		// ... append params
	}
}
```

**For ridgeExt (4 arrays):**

```go
// Generated Decode()
func Decode(items []map[string]any) ([]RidgeExtItem, error) {
	// ... field extractions
}

// Generated ArrayDecl wrapper (multi-array variant)
func (m *module) Build(items []map[string]any) ([]ArrayDecl, error) {
	typed, err := Decode(items)
	byFace, err := Build(typed)  // Returns map[string][][]any
	
	var decls []ArrayDecl
	for name, rows := range byFace {
		if len(rows) > 0 {
			decls = append(decls, ArrayDecl{Name: name, Rows: rows})
		}
	}
	return decls, nil
}

// Hand-written Build()
func Build(items []RidgeExtItem) (map[string][][]any, error) {
	// ... face distribution logic
}
```

**Challenge:** How does schemagen know whether to generate single-array or multi-array wrapper?

**Option 1: Schema metadata**

```yaml
# schema.yaml
module: ridge_ext
output: multi_array  # or "single_array"
array_names:
  - ridgeExtFront
  - ridgeExtBack
  - ridgeExtLeft
  - ridgeExtRight
```

**Option 2: Convention**

If `module.go` exports `Build() (map[string][][]any, error)`, generate multi-array wrapper.
If `module.go` exports `Build() ([][]any, error)`, generate single-array wrapper.

**Option 3: Don't generate wrapper**

Module author writes the wrapper manually in registry.go. We only generate Decode().

**My position:** Option 1 (schema metadata) is cleanest. It makes the output pattern explicit and allows generation to be deterministic.

**But** Option 3 (no wrapper generation) is simpler and gives module authors more control. We've already saved them from marshal/unmarshal with Decode().

---

## Rebuttals

### Jordan Rivera — "The Architect"

Generator, your Option 1 (schema metadata) is interesting, but I think it's overkill.

**Here's why:** We have 10 total modules. 8 are single-array, 2 are multi-array. The multi-array modules (cutouts, ridgeExt) have domain-specific distribution logic that's non-trivial.

**Generating the wrapper doesn't save much work:**

```go
// Single-array wrapper (10 lines)
func (m *module) Build(items []map[string]any) ([]ArrayDecl, error) {
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}
	rows, err := Build(typed)
	if err != nil {
		return nil, err
	}
	return []ArrayDecl{{Name: "pcbStands", Rows: rows}}, nil
}
```

That's boilerplate, yes, but it's **simple, obvious boilerplate**. A module author can write it in 30 seconds.

**Multi-array wrapper is more complex:**

```go
// Multi-array wrapper (15-20 lines)
func (m *module) Build(items []map[string]any) ([]ArrayDecl, error) {
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}
	byFace, err := Build(typed)  // Returns map[string][][]any
	if err != nil {
		return nil, err
	}
	
	var decls []ArrayDecl
	for name, rows := range byFace {
		if len(rows) > 0 {
			decls = append(decls, ArrayDecl{Name: name, Rows: rows})
		}
	}
	sort.Slice(decls, func(i, j int) bool {
		return decls[i].Name < decls[j].Name
	})
	return decls, nil
}
```

**But** this is still straightforward. And it only applies to 2 modules out of 10.

**My revised position:** Don't generate the ArrayDecl wrapper. Generate Decode() only. Module authors write the wrapper.

**Benefits:**
1. Simpler templates (no need for output metadata)
2. More control for module authors
3. Clear separation: generated (Decode) vs. hand-written (wrapper + Build logic)

**Cost:** 10-20 lines of boilerplate per module. But it's *readable* boilerplate.

---

### Six-Faced Friend (`cutouts/`)

Jordan, I agree with not generating the wrapper. But I want to address the **multiArrayFeatureModule** question.

**Should cutouts and ridgeExt share a helper?**

**Current state:** Cutouts has custom `cutoutFeatureModule` in features.go (46 lines).

**If ridgeExt follows the same pattern:** We'll have `ridgeExtFeatureModule` (similar 40-45 lines).

**That's duplication.** Let me propose a shared helper:

```go
// pkg/yappgen/features.go
type multiArrayFeatureModule struct {
	key      string
	field    func(*Model) *[]map[string]any
	builder  func([]map[string]any) ([]ArrayDecl, error)
}

func newMultiArrayFeatureModule(
	key string,
	field func(*Model) *[]map[string]any,
	builder func([]map[string]any) ([]ArrayDecl, error),
) FeatureModule {
	return &multiArrayFeatureModule{
		key:     key,
		field:   field,
		builder: builder,
	}
}

func (m *multiArrayFeatureModule) Collect(...) error {
	// Same as arrayFeatureModule.Collect
}

func (m *multiArrayFeatureModule) Emit(...) error {
	decls, err := m.builder(*m.field(model))
	if err != nil {
		return err
	}
	for _, decl := range decls {
		if len(decl.Rows) > 0 {
			writeArrayDecl(b, decl.Name, decl.Rows)
			b.WriteString("\n")
		}
	}
	return nil
}
```

**Usage:**

```go
var featureModules = []FeatureModule{
	// Single-array modules
	newArrayFeatureModule("pcb_stands", ...),
	newArrayFeatureModule("connectors", ...),
	// ...
	
	// Multi-array modules
	newMultiArrayFeatureModule("cutouts",
		func(m *Model) *[]map[string]any { return &m.Cutouts },
		cutouts.Build),  // Returns []ArrayDecl
	
	newMultiArrayFeatureModule("ridge_ext",
		func(m *Model) *[]map[string]any { return &m.RidgeExt },
		ridgeext.Build),  // Returns []ArrayDecl
}
```

**Benefits:**
1. No code duplication (cutouts and ridgeExt use same helper)
2. Consistent handling of multi-array modules
3. Easy to add future multi-array modules

**My revised position:** Create `multiArrayFeatureModule` helper. Cutouts and ridgeExt both use it.

---

### The Generator (`schemagen/`)

Jordan and Six-Faced Friend, I'm convinced.

**Revised plan:**

1. **Generate Decode() only** (not ArrayDecl wrapper)
2. **Module authors write wrapper** in registry.go
3. **Create multiArrayFeatureModule helper** for cutouts/ridgeExt

**Template changes:**
- Add Decode() generation to schema_gen.go.tmpl (~80 lines)
- Add Decode() tests to schema_gen_test.go.tmpl (~40 lines)
- No wrapper generation

**Total template complexity:** ~120 lines (down from ~195 lines in Round 5)

**Module authoring steps:**
1. Write schema.yaml
2. Run schemagen (generates Decode())
3. Write Build() logic (hand-written)
4. Write registry.go wrapper (10-20 lines boilerplate)
5. Register in features.go (newArrayFeatureModule or newMultiArrayFeatureModule)

**Future extensibility:**

- **labelsPlane:** Single-array, use newArrayFeatureModule
- **ridgeExt:** Multi-array, use newMultiArrayFeatureModule
- **displayMounts:** Single-array, use newArrayFeatureModule

**Unknown future needs:** If a module needs something other than single/multi-array, they can:
1. Implement custom FeatureModule (like cutouts did before)
2. Extend multiArrayFeatureModule with new behavior
3. Add a new helper (e.g., conditionalArrayFeatureModule)

**My final position:** Path 1 (A→B) with ArrayDecl + multiArrayFeatureModule helper is extensible enough for all foreseeable needs.

---

## Moderator Summary

### Key Decisions

**1. ArrayDecl is sufficient for all known future modules**
- 8 single-array modules: return 1 ArrayDecl
- 2 multi-array modules: return multiple ArrayDecls
- Handles conditional arrays, dynamic names, empty outputs

**2. Don't generate ArrayDecl wrapper**
- Generate Decode() only
- Module authors write wrapper (10-20 lines)
- Simpler templates, more control

**3. Create multiArrayFeatureModule helper**
- Shared helper for cutouts and ridgeExt
- Eliminates code duplication
- Makes multi-array pattern explicit

### Consensus Points

All three candidates agree:
- ArrayDecl handles all known future modules (labelsPlane, ridgeExt, displayMounts)
- Generating Decode() is valuable; generating wrapper is overkill
- multiArrayFeatureModule helper reduces duplication

### Implementation Plan

**Step A (Decode generation):**
1. Create pkg/yappgen/decode/helpers.go (shared helpers)
2. Add Decode() generation to schema_gen.go.tmpl
3. Add Decode() tests to schema_gen_test.go.tmpl
4. Update all 7 modules to use Decode()

**Step B (ArrayDecl IR):**
1. Define ArrayDecl type in pkg/registry/schema.go
2. Change FeatureModule.Build return type to []ArrayDecl
3. Create multiArrayFeatureModule helper in features.go
4. Update cutouts to use multiArrayFeatureModule
5. Update all 6 single-array modules to return []ArrayDecl
6. Update emission to loop over ArrayDecls

**Future modules:**
- labelsPlane: Use newArrayFeatureModule + Decode()
- ridgeExt: Use newMultiArrayFeatureModule + Decode()
- displayMounts: Use newArrayFeatureModule + Decode()

### Extensibility Analysis

**Can ArrayDecl handle:**
- ✅ Single-array modules (8 modules)
- ✅ Multi-array modules (2 modules)
- ✅ Conditional arrays (if/else logic)
- ✅ Dynamic array names (loop over map)
- ✅ Empty outputs (empty slice)
- ✅ Future unknown patterns (custom FeatureModule still possible)

**Limitations:**
- ❌ Non-array outputs (e.g., global variables) - but no known need
- ❌ Nested arrays (e.g., array of arrays) - but no known need

**Mitigation:** If future needs arise, we can:
1. Extend ArrayDecl (add fields)
2. Create new IR types (e.g., GlobalDecl)
3. Implement custom FeatureModule

### Open Questions

1. Should multiArrayFeatureModule sort ArrayDecls by name? (For deterministic output)
2. Do we document the multi-array pattern in module authoring guide?
3. Should we add a "module patterns" reference doc?

---

## Decision Point

**Consensus: Path 1 (A→B) with ArrayDecl + multiArrayFeatureModule is extensible**

**Key insights:**
- ArrayDecl handles 100% of known future modules
- multiArrayFeatureModule eliminates duplication for cutouts/ridgeExt
- Not generating wrapper keeps templates simple

**Final implementation:**
1. Generate Decode() only (not wrapper)
2. Add ArrayDecl IR with multiArrayFeatureModule helper
3. Module authors write 10-20 line wrapper (acceptable boilerplate)

**Next steps:** Synthesize all debate rounds into implementation plan and RFC.
