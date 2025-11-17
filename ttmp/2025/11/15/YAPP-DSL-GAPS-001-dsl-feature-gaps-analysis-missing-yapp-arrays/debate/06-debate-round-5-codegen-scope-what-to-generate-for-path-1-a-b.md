---
Title: Debate Round 5 — Codegen Scope: What to Generate for Path 1 (A→B)?
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
    - Path: /home/manuel/code/others/YAPP_Box/pkg/schemagen/
      Note: Code generator
    - Path: /home/manuel/code/others/YAPP_Box/pkg/schemagen/templates/
      Note: Generation templates
ExternalSources: []
Summary: What code generation is needed for Decode() and ArrayDecl IR in Path 1?
LastUpdated: 2025-11-17
---

# Debate Round 5 — Codegen Scope: What to Generate for Path 1 (A→B)?

## Question

**What code should schemagen generate for Path 1 (A→B)?**

**Path 1 consists of:**
- **Step A:** Generate `Decode()` per module
- **Step B:** Add `ArrayDecl` IR for multi-array modules

**Decision from previous rounds:** Path 1 (A→B), single-shot refactor.

## Pre-Debate Research

### Current Schemagen Output

```bash
$ ls -1 pkg/yappgen/modules/pcbstands/schema_*.go
pkg/yappgen/modules/pcbstands/schema_gen.go
pkg/yappgen/modules/pcbstands/schema_gen_test.go
pkg/yappgen/modules/pcbstands/schema_validate.go
```

**Finding:** Schemagen currently generates 3 files per module:
1. `schema_gen.go` - Typed structs, ApplyDefaults(), CustomValidate() stub
2. `schema_gen_test.go` - Tests from schema test cases
3. `schema_validate.go` - ValidateStructure() implementation

### Template Complexity

```bash
$ wc -l pkg/schemagen/templates/*.tmpl
  23 pkg/schemagen/templates/schema_gen.go.tmpl
  20 pkg/schemagen/templates/schema_gen_test.go.tmpl
  14 pkg/schemagen/templates/modules_gen.go.tmpl
 165 pkg/schemagen/templates/schema_validate.go.tmpl
 222 total
```

**Finding:** 222 lines total. Most complex is `schema_validate.go.tmpl` (165 lines).

### Pushbuttons Custom Decode

```go
// pkg/yappgen/modules/pushbuttons/module.go:108-150 (42 lines)
func decodePushButtonsItems(items []map[string]any) ([]PushButtonsItem, error) {
	typed := make([]PushButtonsItem, len(items))
	for idx, raw := range items {
		label := fmt.Sprintf("push_buttons[%d]", idx)
		
		// Manual field extraction with type assertions
		x, ok := getFloat(raw, "x")
		if !ok {
			return nil, fmt.Errorf("%s.x: expected number", label)
		}
		// ... 35 more lines of field extraction
	}
	return typed, nil
}
```

**Finding:** Manual decode is 42 lines for pushbuttons (complex nested structure). Simpler modules would be ~20 lines.

### ArrayDecl IR Concept

From analysis doc, Option B proposes:

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

### The Generator (`schemagen/`)

Let me break down what needs to be generated for Path 1.

**Step A: Generate Decode()**

Template pseudocode:

```go
// Generated in schema_gen.go
func Decode(items []map[string]any) ([]{{.RootStruct}}, error) {
	typed := make([]{{.RootStruct}}, len(items))
	for idx, raw := range items {
		label := fmt.Sprintf("{{.Module}}[%d]", idx)
		
		{{range .Fields}}
		{{if .Required}}
		{{.Name}}, ok := get{{.Type}}(raw, "{{.YAMLName}}")
		if !ok {
			return nil, fmt.Errorf("%s.{{.YAMLName}}: expected {{.Type}}", label)
		}
		{{else}}
		{{.Name}} := getOptional{{.Type}}(raw, "{{.YAMLName}}")
		{{end}}
		{{end}}
		
		item := {{.RootStruct}}{
			{{range .Fields}}{{.GoName}}: {{.Name}},
			{{end}}
		}
		item.ApplyDefaults()
		if err := item.CustomValidate(); err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		typed[idx] = item
	}
	return typed, nil
}
```

**Template complexity:** ~60 lines with field iteration and type handling.

**Helper functions needed:**

```go
func getFloat(m map[string]any, key string) (float64, bool)
func getString(m map[string]any, key string) (string, bool)
func getBool(m map[string]any, key string) (bool, bool)
func getOptionalFloat(m map[string]any, key string) *float64
// ... etc
```

These can be generated once in a shared file or copied into each module.

**Step B: Generate ArrayDecl wrapper**

Option 1: Generate wrapper in `registry.go`:

```go
// Generated in registry.go
func (m *module) Build(items []map[string]any) ([]ArrayDecl, error) {
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}
	rows, err := Build(typed)  // Hand-written
	if err != nil {
		return nil, err
	}
	return []ArrayDecl{{Name: "{{.ScadArray}}", Rows: rows}}, nil
}
```

Option 2: Don't generate wrapper, let module author write it.

**My recommendation:**

- **Generate Decode()** in `schema_gen.go` (~60 lines template)
- **Generate ArrayDecl wrapper** in `registry.go` (~15 lines template)
- **Generate helper functions** in `schema_gen.go` or shared file (~40 lines)

**Total new template code:** ~115 lines

**Complexity:** Moderate. Field iteration is already done in `schema_validate.go.tmpl`, so we can reuse that logic.

---

### Sam Park — "The Codegen Maintainer"

Generator, your template looks good, but let me talk about **debugging generated code**.

**Current problem:** When ValidateStructure fails, the error message is:

```
features.pcb_stands[2].height: expected number|string
```

That's from generated code in `schema_validate.go`. Users can't easily see *why* it failed—they have to read the generated code.

**With generated Decode(), errors will be:**

```
pcb_stands[2].height: expected number
```

Better, but still opaque. **We need good error messages.**

**Proposal: Generate with context**

```go
// Generated Decode with rich errors
func Decode(items []map[string]any) ([]PcbStandsItem, error) {
	for idx, raw := range items {
		label := fmt.Sprintf("pcb_stands[%d]", idx)
		
		x, ok := getFloat(raw, "x")
		if !ok {
			val, exists := raw["x"]
			if !exists {
				return nil, fmt.Errorf("%s.x: missing required field", label)
			}
			return nil, fmt.Errorf("%s.x: expected number, got %T", label, val)
		}
		// ...
	}
}
```

**Template complexity increases:** ~80 lines instead of 60 lines, but error messages are much better.

**My position:** Generate Decode with rich error messages. Users shouldn't have to read generated code to debug type mismatches.

**ArrayDecl wrapper:** Generate it. It's boilerplate that's identical for all single-array modules. Only cutouts needs custom handling.

---

### Casey Thompson — "The New Hire"

*[Looks at the generated code proposal]*

I have a practical question: **Where do I look when something goes wrong?**

**Scenario:** I add a new module, run schemagen, and get a compile error:

```
schema_gen.go:45:12: undefined: getNestedObject
```

Do I:
1. Fix the template?
2. Fix the schema.yaml?
3. File a bug?

**Current system:** When `schema_gen.go` has issues, I know it's generated, so I check the schema.yaml first. If that's correct, I know it's a schemagen bug.

**With Decode() generation:** Same principle should apply. If Decode() has issues, check schema.yaml first.

**My concern:** Helper functions. If `getFloat()` is generated per-module, that's 7 copies. If one has a bug, do I fix all 7? Or do I fix the template?

**Proposal:** Put helper functions in a shared package:

```go
// pkg/yappgen/decode/helpers.go (hand-written, not generated)
package decode

func GetFloat(m map[string]any, key string) (float64, bool) { ... }
func GetString(m map[string]any, key string) (string, bool) { ... }
// ...
```

Then generated Decode() imports and uses them:

```go
// Generated schema_gen.go
import "github.com/wesen/yapp-encl-resolver/pkg/yappgen/decode"

func Decode(items []map[string]any) ([]PcbStandsItem, error) {
	for idx, raw := range items {
		x, ok := decode.GetFloat(raw, "x")
		// ...
	}
}
```

**Benefits:**
1. Helpers are hand-written, easy to debug
2. No duplication across modules
3. Clear separation: generated code calls hand-written helpers

**My position:** Generate Decode() and ArrayDecl wrapper, but use shared hand-written helpers.

---

## Rebuttals

### The Generator (`schemagen/`)

Casey, your shared helpers proposal is excellent. Let me revise my plan.

**New architecture:**

```
pkg/yappgen/decode/
  helpers.go (hand-written)
    - GetFloat(m, key) (float64, bool)
    - GetOptionalFloat(m, key) *float64
    - GetString(m, key) (string, bool)
    - GetObject(m, key) (map[string]any, bool)
    - GetArray(m, key) ([]any, bool)
    - FormatLabel(module, idx) string

pkg/yappgen/modules/pcbstands/
  schema_gen.go (generated)
    - Decode() uses decode.GetFloat(), etc.
  registry.go (partially generated)
    - Build() wrapper uses Decode()
```

**Template changes:**

1. **schema_gen.go.tmpl** adds Decode() function (~50 lines, uses decode.* helpers)
2. **registry.go** gets generated Build() wrapper (~15 lines)
3. New **helpers.go** (hand-written, ~100 lines total)

**Total code:**
- Generated per module: ~65 lines (Decode + wrapper)
- Shared helpers: ~100 lines (one-time)
- Template complexity: ~80 lines

Sam's rich error messages fit perfectly with shared helpers:

```go
// helpers.go
func GetFloat(m map[string]any, key, label string) (float64, error) {
	val, exists := m[key]
	if !exists {
		return 0, fmt.Errorf("%s.%s: missing required field", label, key)
	}
	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("%s.%s: expected number, got %T", label, key, val)
	}
}
```

**My revised position:** Generate Decode() + ArrayDecl wrapper, use shared hand-written helpers with rich error messages.

---

### Sam Park — "The Codegen Maintainer"

Perfect. Let me add one more thing: **nested objects**.

Pushbuttons has nested objects (`cap`, `lid`, `switch`). How do we generate Decode for those?

**Option 1: Recursive decode**

```go
// Generated for pushbuttons
func Decode(items []map[string]any) ([]PushButtonsItem, error) {
	for idx, raw := range items {
		label := fmt.Sprintf("push_buttons[%d]", idx)
		
		capRaw, err := decode.GetObject(raw, "cap", label)
		if err != nil {
			return nil, err
		}
		cap, err := decodeCap(capRaw, label+".cap")
		if err != nil {
			return nil, err
		}
		// ...
	}
}

func decodeCap(m map[string]any, label string) (PushButtonsItemCap, error) {
	length, err := decode.GetFloat(m, "length", label)
	if err != nil {
		return PushButtonsItemCap{}, err
	}
	// ... more fields
	return PushButtonsItemCap{Length: length, ...}, nil
}
```

**Option 2: Inline decode**

Just extract fields inline without helper functions. Gets messy for deep nesting.

**My position:** Option 1 (recursive decode). Generate helper functions for nested objects. Template complexity increases (~120 lines for complex modules) but it's still maintainable.

---

### Casey Thompson — "The New Hire"

Sam, recursive decode makes sense. One more question: **What about arrays of objects?**

Snap joins has a `corners` array:

```yaml
snap_joins:
  - side: left
    corners: [front_left, back_left]
```

How do we decode that?

**Current schema:**

```yaml
fields:
  corners:
    type: array
    items:
      type: string
```

**Generated Decode:**

```go
cornersRaw, err := decode.GetArray(raw, "corners", label)
if err != nil {
	return nil, err
}
corners := make([]string, len(cornersRaw))
for i, v := range cornersRaw {
	s, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("%s.corners[%d]: expected string, got %T", label, i, v)
	}
	corners[i] = s
}
```

**Template complexity:** Moderate. Need to handle `items.type` from schema.

**My position:** Generate array decode logic. It's boilerplate that varies by item type (string, number, object).

---

## Moderator Summary

### Key Decisions

**What to generate:**
1. **Decode() function** in `schema_gen.go` per module
2. **ArrayDecl wrapper** in `registry.go` per module
3. **Shared helpers** in `pkg/yappgen/decode/helpers.go` (hand-written)

**Template complexity:**
- Decode() generation: ~80 lines (handles fields, nested objects, arrays)
- ArrayDecl wrapper: ~15 lines
- Shared helpers: ~100 lines (one-time, hand-written)

**Total:** ~195 lines of new code (95 lines template + 100 lines helpers)

### Consensus Points

1. **Shared helpers:** All candidates agree helpers should be hand-written in a shared package
2. **Rich error messages:** Error messages should include field path and type information
3. **Nested objects:** Generate recursive decode functions for nested structures
4. **Arrays:** Generate array decode logic based on `items.type` from schema

### Implementation Plan

**Step A (Decode generation):**
1. Create `pkg/yappgen/decode/helpers.go` with GetFloat, GetString, GetObject, GetArray, etc.
2. Add Decode() generation to `schema_gen.go.tmpl`
3. Handle nested objects (recursive decode functions)
4. Handle arrays (item type-specific decode)
5. Test with all 7 modules

**Step B (ArrayDecl wrapper):**
1. Define `ArrayDecl` type in `pkg/registry/schema.go`
2. Generate Build() wrapper in `registry.go` for single-array modules
3. Update `cutoutFeatureModule` to handle ArrayDecl
4. Update emission to loop over ArrayDecls

### Open Questions

1. Should ArrayDecl be in `registry` or `yappgen` package?
2. Do we generate ArrayDecl wrapper for cutouts, or leave it custom?
3. How do we test generated Decode() functions? (Generate tests too?)

---

## Decision Point

**Consensus: Generate Decode() + ArrayDecl wrapper with shared helpers**

**Next steps:**
1. Implement shared helpers package
2. Update templates for Decode() generation
3. Test with existing modules before adding ArrayDecl IR

**Next debate:** Enforceability (Question 9) - How do we prevent regression to old patterns?
