---
Title: Architecture & Implementation Guide: Path 1 (A→B) Builder Contract Refactor
Ticket: YAPP-DSL-GAPS-001
Status: active
Topics:
    - yapp
    - architecture
    - implementation
    - codegen
DocType: design
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/
      Note: All debate rounds with decision rationale
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/analysis/04-2025-11-17-builder-contract-and-codegen-options.md
      Note: Original options analysis
ExternalSources: []
Summary: Complete architecture and implementation guide for Path 1 (A→B) - Generate Decode() and add ArrayDecl IR
LastUpdated: 2025-11-17
---

# Architecture & Implementation Guide: Path 1 (A→B) Builder Contract Refactor

## Executive Summary

**Goal:** Eliminate marshal/unmarshal boilerplate and unify single/multi-array module handling.

**Approach:** Path 1 (A→B)
- **Step A:** Generate `Decode()` functions per module
- **Step B:** Add `ArrayDecl` IR for unified output

**Decision rationale:** [Debate Round 2](../debate/03-debate-round-2-path-choice-a-b-vs-a-c.md) - User chose Path 1 over Path 2 (typed Model)

**Migration:** Single-shot refactor (all 7 modules + infrastructure in one commit)

---

## Architecture Overview

### Current State (Before Refactor)

```
User YAML
  ↓
Resolver (map[string]any)
  ↓
features.Collect()
  ↓ normalizeArrayOfMaps()
Model.PcbStands ([]map[string]any)
  ↓
registry.Build([]map[string]any)
  ↓ yaml.Marshal() + yaml.Unmarshal()  ← BOILERPLATE
typed struct (PcbStandsItem)
  ↓
Build logic
  ↓
[][]any
  ↓ (cutouts special-case)
Emit → SCAD
```

**Problems:**
1. Marshal/unmarshal boilerplate in every module (12 error wrapping sites)
2. Cutouts returns `map[string][][]any`, can't implement `registry.FeatureModule.Build`
3. Two patterns: marshal/unmarshal (6 modules) vs custom decode (pushbuttons)

**Decision rationale:** [Debate Round 1](../debate/02-debate-round-1-go-no-go-builder-contract-refactor-urgency.md) - Compounding cost, go now not later

### Target State (After Refactor)

```
User YAML
  ↓
Resolver (map[string]any)
  ↓
features.Collect()
  ↓ normalizeArrayOfMaps()
Model.PcbStands ([]map[string]any)  ← Still untyped
  ↓
registry.Build([]map[string]any)
  ↓ GENERATED Decode()  ← NEW
typed struct (PcbStandsItem)
  ↓
Build logic (hand-written)
  ↓
[][]any
  ↓ ArrayDecl wrapper  ← NEW
[]ArrayDecl
  ↓
Emit → SCAD
```

**Improvements:**
1. ✅ No marshal/unmarshal (generated Decode())
2. ✅ Unified output (ArrayDecl handles single/multi-array)
3. ✅ One pattern (all modules use Decode())

**Decision rationale:** [Debate Round 8](../debate/09-debate-round-8-type-safety-end-to-end-arraydecl-vs-typed-model.md) - Untyped Model is sufficient

---

## Step A: Generate Decode()

**Goal:** Eliminate marshal/unmarshal boilerplate by generating type-safe decode functions.

**Decision rationale:** [Debate Round 5](../debate/06-debate-round-5-codegen-scope-what-to-generate-for-path-1-a-b.md) - Generate Decode() only, not wrapper

### A1: Create Shared Helpers Package

**File:** `pkg/yappgen/decode/helpers.go` (hand-written, ~100 lines)

**Purpose:** Reusable type extraction functions with rich error messages.

**Pseudocode:**

```go
package decode

import "fmt"

// GetFloat extracts a float64 from map with rich error messages
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
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("%s.%s: expected number, got %T", label, key, val)
	}
}

// GetOptionalFloat extracts optional float64 (returns nil if missing)
func GetOptionalFloat(m map[string]any, key, label string) (*float64, error) {
	val, exists := m[key]
	if !exists {
		return nil, nil
	}
	
	switch v := val.(type) {
	case float64:
		return &v, nil
	case int:
		f := float64(v)
		return &f, nil
	case int64:
		f := float64(v)
		return &f, nil
	default:
		return nil, fmt.Errorf("%s.%s: expected number, got %T", label, key, val)
	}
}

// GetString extracts a string from map
func GetString(m map[string]any, key, label string) (string, error) {
	val, exists := m[key]
	if !exists {
		return "", fmt.Errorf("%s.%s: missing required field", label, key)
	}
	
	s, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("%s.%s: expected string, got %T", label, key, val)
	}
	return s, nil
}

// GetOptionalString extracts optional string
func GetOptionalString(m map[string]any, key, label string) (*string, error) {
	val, exists := m[key]
	if !exists {
		return nil, nil
	}
	
	s, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("%s.%s: expected string, got %T", label, key, val)
	}
	return &s, nil
}

// GetBool extracts a bool from map
func GetBool(m map[string]any, key, label string) (bool, error) {
	val, exists := m[key]
	if !exists {
		return false, fmt.Errorf("%s.%s: missing required field", label, key)
	}
	
	b, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("%s.%s: expected bool, got %T", label, key, val)
	}
	return b, nil
}

// GetOptionalBool extracts optional bool
func GetOptionalBool(m map[string]any, key, label string) (*bool, error) {
	val, exists := m[key]
	if !exists {
		return nil, nil
	}
	
	b, ok := val.(bool)
	if !ok {
		return nil, fmt.Errorf("%s.%s: expected bool, got %T", label, key, val)
	}
	return &b, nil
}

// GetObject extracts a nested object (map[string]any)
func GetObject(m map[string]any, key, label string) (map[string]any, error) {
	val, exists := m[key]
	if !exists {
		return nil, fmt.Errorf("%s.%s: missing required field", label, key)
	}
	
	obj, ok := val.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s.%s: expected object, got %T", label, key, val)
	}
	return obj, nil
}

// GetOptionalObject extracts optional nested object
func GetOptionalObject(m map[string]any, key, label string) (map[string]any, error) {
	val, exists := m[key]
	if !exists {
		return nil, nil
	}
	
	obj, ok := val.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s.%s: expected object, got %T", label, key, val)
	}
	return obj, nil
}

// GetArray extracts an array ([]any)
func GetArray(m map[string]any, key, label string) ([]any, error) {
	val, exists := m[key]
	if !exists {
		return nil, fmt.Errorf("%s.%s: missing required field", label, key)
	}
	
	arr, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("%s.%s: expected array, got %T", label, key, val)
	}
	return arr, nil
}

// GetOptionalArray extracts optional array
func GetOptionalArray(m map[string]any, key, label string) ([]any, error) {
	val, exists := m[key]
	if !exists {
		return nil, nil
	}
	
	arr, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("%s.%s: expected array, got %T", label, key, val)
	}
	return arr, nil
}

// FormatLabel creates a consistent label for error messages
func FormatLabel(module string, idx int) string {
	return fmt.Sprintf("%s[%d]", module, idx)
}
```

**Why hand-written?** [Debate Round 5](../debate/06-debate-round-5-codegen-scope-what-to-generate-for-path-1-a-b.md#casey-thompson--the-new-hire) - Easier to debug, no duplication across modules

### A2: Update Schemagen Templates

**File:** `pkg/schemagen/templates/schema_gen.go.tmpl`

**Add:** Decode() function generation (~80 lines template code)

**Pseudocode:**

```go
// Template: schema_gen.go.tmpl
// Code generated by schemagen; DO NOT EDIT.

package {{.Package}}

import (
	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/decode"
)

{{range .Structs}}
// Decode converts []map[string]any to typed []{{.Name}}
func Decode(items []map[string]any) ([]{{.Name}}, error) {
	typed := make([]{{.Name}}, len(items))
	
	for idx, raw := range items {
		label := decode.FormatLabel("{{$.Module}}", idx)
		
		{{range .Fields}}
		{{if .Required}}
		// Required field: {{.YAMLName}}
		{{if eq .Type "number"}}
		{{.GoName}}, err := decode.GetFloat(raw, "{{.YAMLName}}", label)
		if err != nil {
			return nil, err
		}
		{{else if eq .Type "string"}}
		{{.GoName}}, err := decode.GetString(raw, "{{.YAMLName}}", label)
		if err != nil {
			return nil, err
		}
		{{else if eq .Type "bool"}}
		{{.GoName}}, err := decode.GetBool(raw, "{{.YAMLName}}", label)
		if err != nil {
			return nil, err
		}
		{{else if eq .Type "object"}}
		{{.GoName}}Raw, err := decode.GetObject(raw, "{{.YAMLName}}", label)
		if err != nil {
			return nil, err
		}
		{{.GoName}}, err := decode{{.GoName}}({{.GoName}}Raw, label+".{{.YAMLName}}")
		if err != nil {
			return nil, err
		}
		{{end}}
		{{else}}
		// Optional field: {{.YAMLName}}
		{{if eq .Type "number"}}
		{{.GoName}}, err := decode.GetOptionalFloat(raw, "{{.YAMLName}}", label)
		if err != nil {
			return nil, err
		}
		{{else if eq .Type "string"}}
		{{.GoName}}, err := decode.GetOptionalString(raw, "{{.YAMLName}}", label)
		if err != nil {
			return nil, err
		}
		{{else if eq .Type "bool"}}
		{{.GoName}}, err := decode.GetOptionalBool(raw, "{{.YAMLName}}", label)
		if err != nil {
			return nil, err
		}
		{{else if eq .Type "object"}}
		{{.GoName}}Raw, err := decode.GetOptionalObject(raw, "{{.YAMLName}}", label)
		if err != nil {
			return nil, err
		}
		var {{.GoName}} *{{.GoType}}
		if {{.GoName}}Raw != nil {
			decoded, err := decode{{.GoName}}({{.GoName}}Raw, label+".{{.YAMLName}}")
			if err != nil {
				return nil, err
			}
			{{.GoName}} = &decoded
		}
		{{end}}
		{{end}}
		{{end}}
		
		item := {{.Name}}{
			{{range .Fields}}{{.GoName}}: {{.GoName}},
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

{{if .HasNestedObjects}}
// Helper decoders for nested objects
{{range .NestedObjects}}
func decode{{.Name}}(m map[string]any, label string) ({{.Type}}, error) {
	{{range .Fields}}
	{{if .Required}}
	{{.GoName}}, err := decode.Get{{.TypeHelper}}(m, "{{.YAMLName}}", label)
	if err != nil {
		return {{$.Type}}{}, err
	}
	{{else}}
	{{.GoName}}, err := decode.GetOptional{{.TypeHelper}}(m, "{{.YAMLName}}", label)
	if err != nil {
		return {{$.Type}}{}, err
	}
	{{end}}
	{{end}}
	
	return {{.Type}}{
		{{range .Fields}}{{.GoName}}: {{.GoName}},
		{{end}}
	}, nil
}
{{end}}
{{end}}

{{end}}
```

**Array handling:**

```go
// Template addition for array fields
{{if eq .Type "array"}}
{{if .Required}}
{{.GoName}}Raw, err := decode.GetArray(raw, "{{.YAMLName}}", label)
if err != nil {
	return nil, err
}
{{.GoName}} := make([]{{.ItemType}}, len({{.GoName}}Raw))
for i, v := range {{.GoName}}Raw {
	{{if eq .ItemType "string"}}
	s, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("%s.{{.YAMLName}}[%d]: expected string, got %T", label, i, v)
	}
	{{.GoName}}[i] = s
	{{else if eq .ItemType "float64"}}
	// Similar for numbers
	{{end}}
}
{{else}}
// Optional array handling
{{end}}
{{end}}
```

**Why this approach?** [Debate Round 5](../debate/06-debate-round-5-codegen-scope-what-to-generate-for-path-1-a-b.md#the-generator-schemagen) - Template complexity ~80 lines, handles nested objects and arrays

### A3: Generate Decode() Tests

**File:** `pkg/schemagen/templates/schema_gen_test.go.tmpl`

**Add:** Test generation for Decode() (~40 lines template code)

**Pseudocode:**

```go
// Template: schema_gen_test.go.tmpl additions
package {{.Package}}

import (
	"strings"
	"testing"
)

func TestDecode_ValidInput(t *testing.T) {
	items := []map[string]any{
		{
			{{range .RequiredFields}}"{{.YAMLName}}": {{.TestValue}},
			{{end}}
		},
	}
	
	typed, err := Decode(items)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	
	if len(typed) != 1 {
		t.Fatalf("expected 1 item, got %d", len(typed))
	}
	
	{{range .RequiredFields}}
	if typed[0].{{.GoName}} != {{.TestValue}} {
		t.Errorf("expected {{.GoName}}={{.TestValue}}, got %v", typed[0].{{.GoName}})
	}
	{{end}}
}

func TestDecode_MissingRequired(t *testing.T) {
	{{range .RequiredFields}}
	t.Run("missing_{{.YAMLName}}", func(t *testing.T) {
		items := []map[string]any{
			{
				{{range $other := $.RequiredFields}}
				{{if ne $other.YAMLName .YAMLName}}"{{$other.YAMLName}}": {{$other.TestValue}},
				{{end}}
				{{end}}
			},
		}
		
		_, err := Decode(items)
		if err == nil {
			t.Fatal("expected error for missing required field")
		}
		
		if !strings.Contains(err.Error(), "{{.YAMLName}}: missing required field") {
			t.Errorf("wrong error: %v", err)
		}
	})
	{{end}}
}

func TestDecode_WrongType(t *testing.T) {
	{{range .RequiredFields}}
	t.Run("wrong_type_{{.YAMLName}}", func(t *testing.T) {
		items := []map[string]any{
			{
				{{range $other := $.RequiredFields}}
				{{if ne $other.YAMLName .YAMLName}}"{{$other.YAMLName}}": {{$other.TestValue}},
				{{else}}"{{.YAMLName}}": "not a {{.Type}}",
				{{end}}
				{{end}}
			},
		}
		
		_, err := Decode(items)
		if err == nil {
			t.Fatal("expected error for wrong type")
		}
		
		if !strings.Contains(err.Error(), "expected {{.Type}}") {
			t.Errorf("wrong error: %v", err)
		}
	})
	{{end}}
}
```

**Why generate tests?** [Debate Round 6](../debate/07-debate-round-6-enforceability-keeping-modules-honest-post-refactor.md#sam-park--the-codegen-maintainer) - Ensures generated code works, catches regressions

### A4: Update All 7 Modules

**For each module:** `pkg/yappgen/modules/*/module.go`

**Before:**

```go
func Build(items []map[string]any) ([][]any, error) {
	var out [][]any
	for idx, it := range items {
		label := fmt.Sprintf("pcb_stands[%d]", idx)
		
		// Marshal/unmarshal boilerplate (8 lines)
		data, err := yaml.Marshal(it)
		if err != nil {
			return nil, errors.Wrapf(err, "%s: marshal", label)
		}
		var item PcbStandsItem
		if err := yaml.Unmarshal(data, &item); err != nil {
			return nil, errors.Wrapf(err, "%s: unmarshal", label)
		}
		
		item.ApplyDefaults()
		if err := item.CustomValidate(); err != nil {
			return nil, errors.Wrapf(err, "%s", label)
		}
		
		// Build logic (10-30 lines)
		params := []any{item.X, item.Y, ptrOrUndef(item.Height), ...}
		out = append(out, params)
	}
	return out, nil
}
```

**After:**

```go
func Build(items []map[string]any) ([][]any, error) {
	// Use generated Decode()
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}
	
	var out [][]any
	for _, item := range typed {
		// Build logic (10-30 lines) - unchanged
		params := []any{item.X, item.Y, ptrOrUndef(item.Height), ...}
		out = append(out, params)
	}
	return out, nil
}
```

**Changes:**
- ❌ Remove: Marshal/unmarshal boilerplate (8 lines)
- ❌ Remove: ApplyDefaults() call (handled in Decode())
- ❌ Remove: CustomValidate() call (handled in Decode())
- ✅ Add: Decode() call (3 lines)

**Net change:** -5 lines per module, cleaner code

**Modules to update:**
1. `pkg/yappgen/modules/pcbstands/module.go`
2. `pkg/yappgen/modules/connectors/module.go`
3. `pkg/yappgen/modules/boxmounts/module.go`
4. `pkg/yappgen/modules/snapjoins/module.go`
5. `pkg/yappgen/modules/lighttubes/module.go`
6. `pkg/yappgen/modules/cutouts/module.go`
7. `pkg/yappgen/modules/pushbuttons/module.go` (already has custom decode, replace with generated)

---

## Step B: Add ArrayDecl IR

**Goal:** Unify single-array and multi-array module handling.

**Decision rationale:** [Debate Round 4](../debate/05-debate-round-4-special-cases-how-do-we-retire-cutoutfeaturemodule.md) - Keep custom FeatureModule for cutouts, use ArrayDecl

### B1: Define ArrayDecl Type

**File:** `pkg/registry/schema.go`

**Add:**

```go
// ArrayDecl represents a single OpenSCAD array declaration
type ArrayDecl struct {
	Name string   // Array name (e.g., "pcbStands", "cutoutsFront")
	Rows [][]any  // Array rows (positional parameters)
}
```

**Why this structure?** [Debate Round 7](../debate/08-debate-round-7-future-extensibility-accommodating-new-output-shapes.md#jordan-rivera--the-architect) - Handles both single and multi-array modules

### B2: Change FeatureModule Interface

**File:** `pkg/registry/schema.go`

**Before:**

```go
type FeatureModule interface {
	Schema() ModuleSchema
	Build(items []map[string]any) ([][]any, error)
}
```

**After:**

```go
type FeatureModule interface {
	Schema() ModuleSchema
	Build(items []map[string]any) ([]ArrayDecl, error)
}
```

**Impact:** All modules must return `[]ArrayDecl` instead of `[][]any` (compile-time enforcement)

**Why this change?** [Debate Round 6](../debate/07-debate-round-6-enforceability-keeping-modules-honest-post-refactor.md#the-contract-registryschema.go) - Interface-level enforcement

### B3: Add ArrayDecl Wrappers to Modules

**For single-array modules** (6 modules: pcbstands, connectors, boxmounts, snapjoins, lighttubes, pushbuttons):

**File:** `pkg/yappgen/modules/*/registry.go`

**Before:**

```go
func (m *module) Build(items []map[string]any) ([][]any, error) {
	return Build(items)
}
```

**After:**

```go
func (m *module) Build(items []map[string]any) ([]ArrayDecl, error) {
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}
	
	rows, err := Build(typed)
	if err != nil {
		return nil, err
	}
	
	return []registry.ArrayDecl{{
		Name: "pcbStands",  // Module-specific name
		Rows: rows,
	}}, nil
}
```

**Note:** Module author writes this wrapper (10-20 lines). Not generated.

**Why not generate?** [Debate Round 7](../debate/08-debate-round-7-future-extensibility-accommodating-new-output-shapes.md#jordan-rivera--the-architect-1) - Simple boilerplate, more control for authors

**For multi-array module** (cutouts):

**File:** `pkg/yappgen/modules/cutouts/registry.go`

**Before:**

```go
func (m *module) Build(items []map[string]any) ([][]any, error) {
	// Note: cutouts return a map by face, not a simple array
	// This is a special case handled in features.go
	byFace, err := Build(items)
	if err != nil {
		return nil, err
	}
	// For now, return empty to satisfy interface
	_ = byFace
	return nil, nil
}
```

**After:**

```go
func (m *module) Build(items []map[string]any) ([]registry.ArrayDecl, error) {
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}
	
	byFace, err := Build(typed)  // Returns map[string][][]any
	if err != nil {
		return nil, err
	}
	
	var decls []registry.ArrayDecl
	for name, rows := range byFace {
		if len(rows) > 0 {
			decls = append(decls, registry.ArrayDecl{
				Name: name,
				Rows: rows,
			})
		}
	}
	
	// Sort for deterministic output
	sort.Slice(decls, func(i, j int) bool {
		return decls[i].Name < decls[j].Name
	})
	
	return decls, nil
}
```

**Why this pattern?** [Debate Round 4](../debate/05-debate-round-4-special-cases-how-do-we-retire-cutoutfeaturemodule.md#six-faced-friend-cutouts-1) - Face distribution is domain logic

### B4: Update Module Build Signatures

**For single-array modules:**

**File:** `pkg/yappgen/modules/*/module.go`

**Before:**

```go
func Build(items []map[string]any) ([][]any, error) {
	typed, err := Decode(items)
	// ... build logic
	return out, nil
}
```

**After:**

```go
// Build converts typed items to YAPP array format
// Note: This returns [][]any, not []ArrayDecl
// The ArrayDecl wrapper is in registry.go
func Build(items []PcbStandsItem) ([][]any, error) {
	var out [][]any
	for _, item := range items {
		params := []any{item.X, item.Y, ptrOrUndef(item.Height), ...}
		out = append(out, params)
	}
	return out, nil
}
```

**Key change:** Signature now takes typed items (`[]PcbStandsItem`) instead of `[]map[string]any`

**For multi-array module (cutouts):**

**File:** `pkg/yappgen/modules/cutouts/module.go`

**Before:**

```go
func Build(items []map[string]any) (map[string][][]any, error) {
	// ... marshal/unmarshal ...
	// ... build logic
}
```

**After:**

```go
func Build(items []CutoutsItem) (map[string][][]any, error) {
	// No marshal/unmarshal, items already typed
	byFace := map[string][][]any{
		"cutoutsFront": {},
		"cutoutsBack":  {},
		"cutoutsLeft":  {},
		"cutoutsRight": {},
		"cutoutsLid":   {},
		"cutoutsBase":  {},
	}
	
	for _, item := range items {
		faceArray, err := faceArrayName(item.Face)
		if err != nil {
			return nil, err
		}
		params := []any{...}
		byFace[faceArray] = append(byFace[faceArray], params)
	}
	
	return byFace, nil
}
```

### B5: Create multiArrayFeatureModule Helper

**File:** `pkg/yappgen/features.go`

**Add:**

```go
type multiArrayFeatureModule struct {
	key     string
	field   func(*Model) *[]map[string]any
	builder func([]map[string]any) ([]registry.ArrayDecl, error)
}

func newMultiArrayFeatureModule(
	key string,
	field func(*Model) *[]map[string]any,
	builder func([]map[string]any) ([]registry.ArrayDecl, error),
) FeatureModule {
	return &multiArrayFeatureModule{
		key:     key,
		field:   field,
		builder: builder,
	}
}

func (m *multiArrayFeatureModule) Name() string {
	return m.key
}

func (m *multiArrayFeatureModule) Collect(resolved map[string]any, features map[string]any, model *Model) error {
	ptr := m.field(model)
	if features == nil {
		*ptr = nil
		return nil
	}
	arr, ok := getArray(features, m.key)
	if !ok {
		*ptr = nil
		return nil
	}
	items := normalizeArrayOfMaps(arr)
	*ptr = items
	return nil
}

func (m *multiArrayFeatureModule) Emit(ctx context.Context, model *Model, b *strings.Builder) error {
	ptr := m.field(model)
	if len(*ptr) == 0 {
		return nil
	}
	
	decls, err := m.builder(*ptr)
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

**Why this helper?** [Debate Round 7](../debate/08-debate-round-7-future-extensibility-accommodating-new-output-shapes.md#six-faced-friend-cutouts-1) - Eliminates duplication for cutouts/ridgeExt

### B6: Update arrayFeatureModule for ArrayDecl

**File:** `pkg/yappgen/features.go`

**Before:**

```go
type arrayFeatureModule struct {
	key          string
	scadName     string
	field        func(*Model) *[]map[string]any
	builder      func([]map[string]any) ([][]any, error)
	afterCollect func(*Model, []map[string]any)
}

func (m *arrayFeatureModule) Emit(ctx context.Context, model *Model, b *strings.Builder) error {
	ptr := m.field(model)
	if len(*ptr) == 0 {
		return nil
	}
	rows, err := m.builder(*ptr)
	if err != nil {
		return err
	}
	writeArrayDecl(b, m.scadName, rows)
	b.WriteString("\n")
	return nil
}
```

**After:**

```go
type arrayFeatureModule struct {
	key          string
	scadName     string
	field        func(*Model) *[]map[string]any
	builder      func([]map[string]any) ([]registry.ArrayDecl, error)  // ← Changed
	afterCollect func(*Model, []map[string]any)
}

func (m *arrayFeatureModule) Emit(ctx context.Context, model *Model, b *strings.Builder) error {
	ptr := m.field(model)
	if len(*ptr) == 0 {
		return nil
	}
	
	decls, err := m.builder(*ptr)  // ← Returns []ArrayDecl
	if err != nil {
		return err
	}
	
	// Single-array modules return exactly 1 ArrayDecl
	if len(decls) != 1 {
		return fmt.Errorf("arrayFeatureModule %s returned %d arrays, expected 1", m.key, len(decls))
	}
	
	writeArrayDecl(b, decls[0].Name, decls[0].Rows)
	b.WriteString("\n")
	return nil
}
```

### B7: Update Feature Module Registration

**File:** `pkg/yappgen/features.go`

**Before:**

```go
var featureModules = []FeatureModule{
	newArrayFeatureModule("pcb_stands", "pcbStands",
		func(m *Model) *[]map[string]any { return &m.PcbStands },
		pcbstands.Build, nil),
	// ... 5 more single-array modules
	newCutoutFeatureModule(),
}
```

**After:**

```go
var featureModules = []FeatureModule{
	newArrayFeatureModule("pcb_stands", "pcbStands",
		func(m *Model) *[]map[string]any { return &m.PcbStands },
		pcbstands.NewModule().Build, nil),  // ← Now returns []ArrayDecl
	newArrayFeatureModule("connectors", "connectors",
		func(m *Model) *[]map[string]any { return &m.Connectors },
		connectors.NewModule().Build, nil),
	newArrayFeatureModule("box_mounts", "boxMounts",
		func(m *Model) *[]map[string]any { return &m.BoxMounts },
		boxmounts.NewModule().Build, nil),
	newArrayFeatureModule("push_buttons", "pushButtons",
		func(m *Model) *[]map[string]any { return &m.PushButtons },
		pushbuttons.NewModule().Build,
		func(m *Model, items []map[string]any) {
			m.PrintSwitchExtenders = len(items) > 0
		}),
	newArrayFeatureModule("snap_joins", "snapJoins",
		func(m *Model) *[]map[string]any { return &m.SnapJoins },
		snapjoins.NewModule().Build, nil),
	newArrayFeatureModule("light_tubes", "lightTubes",
		func(m *Model) *[]map[string]any { return &m.LightTubes },
		lighttubes.NewModule().Build, nil),
	
	// Multi-array module
	newMultiArrayFeatureModule("cutouts",
		func(m *Model) *[]map[string]any { return &m.Cutouts },
		cutouts.NewModule().Build),
}
```

**Remove:** `newCutoutFeatureModule()` (replaced by `newMultiArrayFeatureModule`)

---

## Enforcement Infrastructure

**Decision rationale:** [Debate Round 6](../debate/07-debate-round-6-enforceability-keeping-modules-honest-post-refactor.md) - Four-layer enforcement

### E1: Linter Rules

**File:** `.golangci.yml` (new file)

```yaml
linters:
  enable:
    - forbidigo  # Forbid specific function calls
    - govet
    - errcheck
    - staticcheck
    
linters-settings:
  forbidigo:
    forbid:
      - pattern: 'yaml\.Marshal'
        msg: 'Do not use yaml.Marshal in module builders; use generated Decode() instead'
        pkg: '^github\.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/.*/module\.go$'
      
      - pattern: 'yaml\.Unmarshal'
        msg: 'Do not use yaml.Unmarshal in module builders; use generated Decode() instead'
        pkg: '^github\.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/.*/module\.go$'

run:
  timeout: 5m
  tests: true
```

**Add to CI:** `.github/workflows/lint.yml` (or equivalent)

```yaml
- name: Run golangci-lint
  run: golangci-lint run --config .golangci.yml
```

### E2: Generated Tests

**Already covered in A3** - Tests generated for Decode() validate:
1. Valid input decodes correctly
2. Missing required fields produce errors
3. Wrong types produce errors
4. Error messages are clear

**Add:** Build wrapper tests (generated in schema_gen_test.go.tmpl)

```go
func TestBuild_UsesGeneratedDecode(t *testing.T) {
	// Verify Build() uses Decode() by checking error messages
	items := []map[string]any{
		{"x": "not a number"},  // Type error
	}
	
	_, err := NewModule().Build(items)
	if err == nil {
		t.Fatal("expected error for invalid type")
	}
	
	// If Build() uses Decode(), error message will contain "expected number"
	if !strings.Contains(err.Error(), "expected number") {
		t.Errorf("Build() may not be using Decode(); error: %v", err)
	}
}

func TestBuild_ReturnsArrayDecl(t *testing.T) {
	items := []map[string]any{
		{"x": 10.0, "y": 20.0},
	}
	
	decls, err := NewModule().Build(items)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	if len(decls) != 1 {
		t.Fatalf("expected 1 ArrayDecl, got %d", len(decls))
	}
	
	if decls[0].Name != "pcbStands" {
		t.Errorf("wrong array name: %s", decls[0].Name)
	}
	
	if len(decls[0].Rows) != 1 {
		t.Errorf("expected 1 row, got %d", len(decls[0].Rows))
	}
}
```

### E3: Documentation Updates

**File:** `pkg/docs/tutorials/yapp-module-authoring-guide.md`

**Update Step 5 (Implement Builder):**

**Before:**

```markdown
## Step 5: Implement Builder

Create `module.go`:

```go
func Build(items []map[string]any) ([][]any, error) {
    var out [][]any
    for idx, it := range items {
        label := fmt.Sprintf("your_module[%d]", idx)
        
        // Unmarshal into typed struct
        data, err := yaml.Marshal(it)
        if err != nil {
            return nil, errors.Wrapf(err, "%s: marshal", label)
        }
        var item YourModuleItem
        if err := yaml.Unmarshal(data, &item); err != nil {
            return nil, errors.Wrapf(err, "%s: unmarshal", label)
        }
        
        item.ApplyDefaults()
        if err := item.CustomValidate(); err != nil {
            return nil, errors.Wrapf(err, "%s", label)
        }
        
        // Build positional array...
    }
}
```
```

**After:**

```markdown
## Step 5: Implement Builder

Create `module.go`:

```go
// Build converts typed items to YAPP array format
func Build(items []YourModuleItem) ([][]any, error) {
    var out [][]any
    for _, item := range items {
        // Build positional array matching YAPP parameter order
        params := []any{
            item.X,
            item.Y,
            ptrOrUndef(item.Diameter),
            // ... more parameters
        }
        
        // Add flags if needed
        if item.Shape != nil {
            flag := shapeToFlag(*item.Shape)
            params = append(params, flag)
        }
        
        out = append(out, params)
    }
    return out, nil
}
```

**Key changes:**
- ✅ Takes typed items (`[]YourModuleItem`) not `[]map[string]any`
- ✅ No marshal/unmarshal (handled by generated Decode())
- ✅ No ApplyDefaults() or CustomValidate() (handled by Decode())
- ✅ Just business logic (building params array)

## Step 6: Implement Registry Wrapper

Create `registry.go` (or update existing):

```go
func (m *module) Build(items []map[string]any) ([]registry.ArrayDecl, error) {
    // Use generated Decode()
    typed, err := Decode(items)
    if err != nil {
        return nil, err
    }
    
    // Call hand-written Build()
    rows, err := Build(typed)
    if err != nil {
        return nil, err
    }
    
    // Wrap in ArrayDecl
    return []registry.ArrayDecl{{
        Name: "yourModuleName",  // e.g., "pcbStands"
        Rows: rows,
    }}, nil
}
```

**For multi-array modules (like cutouts):**

```go
func (m *module) Build(items []map[string]any) ([]registry.ArrayDecl, error) {
    typed, err := Decode(items)
    if err != nil {
        return nil, err
    }
    
    byFace, err := Build(typed)  // Returns map[string][][]any
    if err != nil {
        return nil, err
    }
    
    var decls []registry.ArrayDecl
    for name, rows := range byFace {
        if len(rows) > 0 {
            decls = append(decls, registry.ArrayDecl{
                Name: name,
                Rows: rows,
            })
        }
    }
    
    sort.Slice(decls, func(i, j int) bool {
        return decls[i].Name < decls[j].Name
    })
    
    return decls, nil
}
```

## DO NOT:
- ❌ Use `yaml.Marshal` / `yaml.Unmarshal`
- ❌ Return `[][]any` from registry.Build() (must return `[]ArrayDecl`)
- ❌ Copy code from pre-refactor modules

## DO:
- ✅ Use generated `Decode()` function
- ✅ Return `[]registry.ArrayDecl` from registry.Build()
- ✅ Follow this guide's pattern
```

**Add:** Migration guide

**File:** `pkg/docs/migrations/path-1-a-b-refactor.md` (new file)

```markdown
# Migration Guide: Path 1 (A→B) Refactor

## What Changed (2025-11-17)

We refactored the module builder system:
- ✅ Generated Decode() functions (no more marshal/unmarshal)
- ✅ ArrayDecl IR (unified single/multi-array handling)
- ✅ Linter rules (enforce new patterns)

## If You're Adding a New Module

Follow the updated [Module Authoring Guide](../tutorials/yapp-module-authoring-guide.md).

**DO NOT:**
- Copy code from commits before 2025-11-17
- Use yaml.Marshal/Unmarshal
- Return [][]any from registry.Build()

## If You See Linter Errors

```
module.go:25: Do not use yaml.Marshal in module builders
```

**Fix:** The old marshal/unmarshal pattern is deprecated. Use generated Decode():

```go
// Before (DEPRECATED)
data, _ := yaml.Marshal(it)
var item MyItem
yaml.Unmarshal(data, &item)

// After (CORRECT)
typed, err := Decode(items)
if err != nil {
    return nil, err
}
item := typed[0]
```
```

### E4: Pre-Merge Checklist

**Add to:** `pkg/docs/tutorials/yapp-module-authoring-guide.md`

```markdown
## Pre-Merge Checklist

Before submitting your module PR:
- [ ] Used generated Decode() (no marshal/unmarshal)
- [ ] Returned []registry.ArrayDecl from registry.Build()
- [ ] Build() takes typed items ([]YourModuleItem)
- [ ] Added tests for Build()
- [ ] Ran `go test ./...` successfully
- [ ] Ran `golangci-lint run` successfully
- [ ] Ran `go run ./cmd/schemagen discover` to regenerate
```

---

## Implementation Checklist

### Phase 1: Infrastructure (Step A)

- [ ] **A1:** Create `pkg/yappgen/decode/helpers.go` with shared helpers
  - [ ] GetFloat, GetOptionalFloat
  - [ ] GetString, GetOptionalString
  - [ ] GetBool, GetOptionalBool
  - [ ] GetObject, GetOptionalObject
  - [ ] GetArray, GetOptionalArray
  - [ ] FormatLabel
  - [ ] Write tests for helpers

- [ ] **A2:** Update `pkg/schemagen/templates/schema_gen.go.tmpl`
  - [ ] Add Decode() function generation
  - [ ] Handle required vs optional fields
  - [ ] Handle nested objects
  - [ ] Handle arrays
  - [ ] Generate helper decoders for nested objects

- [ ] **A3:** Update `pkg/schemagen/templates/schema_gen_test.go.tmpl`
  - [ ] Generate TestDecode_ValidInput
  - [ ] Generate TestDecode_MissingRequired
  - [ ] Generate TestDecode_WrongType

- [ ] **A4:** Run schemagen to regenerate all modules
  ```bash
  go run ./cmd/schemagen discover
  ```

- [ ] **A5:** Update module Build() functions (7 modules)
  - [ ] pcbstands: Replace marshal/unmarshal with Decode()
  - [ ] connectors: Replace marshal/unmarshal with Decode()
  - [ ] boxmounts: Replace marshal/unmarshal with Decode()
  - [ ] snapjoins: Replace marshal/unmarshal with Decode()
  - [ ] lighttubes: Replace marshal/unmarshal with Decode()
  - [ ] cutouts: Replace marshal/unmarshal with Decode()
  - [ ] pushbuttons: Replace custom decode with generated Decode()

- [ ] **A6:** Run tests
  ```bash
  go test ./pkg/yappgen/modules/...
  ```

### Phase 2: ArrayDecl IR (Step B)

- [ ] **B1:** Define ArrayDecl in `pkg/registry/schema.go`

- [ ] **B2:** Change FeatureModule.Build interface
  - [ ] Update interface definition
  - [ ] Verify all implementations break (expected)

- [ ] **B3:** Add ArrayDecl wrappers to module registry.go (7 modules)
  - [ ] pcbstands: Single-array wrapper
  - [ ] connectors: Single-array wrapper
  - [ ] boxmounts: Single-array wrapper
  - [ ] snapjoins: Single-array wrapper
  - [ ] lighttubes: Single-array wrapper
  - [ ] pushbuttons: Single-array wrapper
  - [ ] cutouts: Multi-array wrapper

- [ ] **B4:** Update module Build() signatures (7 modules)
  - [ ] pcbstands: Take []PcbStandsItem
  - [ ] connectors: Take []ConnectorsItem
  - [ ] boxmounts: Take []BoxMountsItem
  - [ ] snapjoins: Take []SnapJoinsItem
  - [ ] lighttubes: Take []LightTubesItem
  - [ ] pushbuttons: Take []PushButtonsItem
  - [ ] cutouts: Take []CutoutsItem

- [ ] **B5:** Create multiArrayFeatureModule helper in `pkg/yappgen/features.go`

- [ ] **B6:** Update arrayFeatureModule for ArrayDecl in `pkg/yappgen/features.go`

- [ ] **B7:** Update feature module registration in `pkg/yappgen/features.go`
  - [ ] Update 6 single-array modules
  - [ ] Replace cutoutFeatureModule with multiArrayFeatureModule

- [ ] **B8:** Run tests
  ```bash
  go test ./pkg/yappgen/...
  go test ./pkg/registry/...
  ```

### Phase 3: Enforcement

- [ ] **E1:** Add linter rules
  - [ ] Create `.golangci.yml`
  - [ ] Add forbidigo rules for yaml.Marshal/Unmarshal
  - [ ] Update CI to run linter

- [ ] **E2:** Generate Build wrapper tests
  - [ ] Update schema_gen_test.go.tmpl
  - [ ] Add TestBuild_UsesGeneratedDecode
  - [ ] Add TestBuild_ReturnsArrayDecl
  - [ ] Regenerate tests

- [ ] **E3:** Update documentation
  - [ ] Update yapp-module-authoring-guide.md
  - [ ] Create path-1-a-b-refactor.md migration guide
  - [ ] Add pre-merge checklist

- [ ] **E4:** Run linter
  ```bash
  golangci-lint run
  ```

### Phase 4: Validation

- [ ] **V1:** Run all tests
  ```bash
  go test ./...
  ```

- [ ] **V2:** Test end-to-end with example YAML
  ```bash
  go run ./cmd/yappctl resolve -i examples/04-features.yaml -o /tmp/resolved.yaml
  go run ./cmd/yappctl generate -i /tmp/resolved.yaml -o /tmp/output.scad
  ```

- [ ] **V3:** Verify SCAD output is identical to pre-refactor
  ```bash
  # Generate SCAD before refactor
  git stash
  go run ./cmd/yappctl generate -i examples/04-features.yaml -o /tmp/before.scad
  
  # Generate SCAD after refactor
  git stash pop
  go run ./cmd/yappctl generate -i examples/04-features.yaml -o /tmp/after.scad
  
  # Compare
  diff /tmp/before.scad /tmp/after.scad
  # Should be identical (or only whitespace differences)
  ```

- [ ] **V4:** Try adding a test module following new guide
  - [ ] Create schema.yaml
  - [ ] Run schemagen
  - [ ] Write Build() logic
  - [ ] Write registry wrapper
  - [ ] Verify it works

---

## Testing Strategy

### Unit Tests

**For Decode() (generated):**
- Valid input decodes correctly
- Missing required fields produce errors
- Wrong types produce errors
- Optional fields handle nil
- Nested objects decode correctly
- Arrays decode correctly

**For Build() (hand-written):**
- Valid typed input produces correct [][]any
- Flags are appended correctly
- Optional fields use ptrOrUndef()

**For ArrayDecl wrapper (hand-written):**
- Single-array modules return 1 ArrayDecl
- Multi-array modules return N ArrayDecls
- Empty inputs return empty []ArrayDecl
- Array names are correct

### Integration Tests

**End-to-end:**
- Resolve YAML → Generate SCAD → Compare with expected output
- Test all 7 modules with example YAML
- Verify SCAD output is identical to pre-refactor

### Regression Tests

**Linter:**
- Verify linter catches yaml.Marshal usage
- Verify linter catches yaml.Unmarshal usage
- Verify linter allows generated Decode()

---

## Rollback Plan

**If something goes wrong:**

1. **Revert the commit:**
   ```bash
   git revert HEAD
   ```

2. **Or cherry-pick fixes:**
   ```bash
   # Fix specific module
   git checkout HEAD~1 -- pkg/yappgen/modules/pcbstands/
   ```

3. **Regenerate:**
   ```bash
   go run ./cmd/schemagen discover
   go test ./...
   ```

**Why single-shot is safe:**
- All changes in one commit
- Easy to revert if needed
- Tests validate correctness
- SCAD output comparison catches regressions

---

## Success Criteria

**After refactor:**

1. ✅ All tests pass (`go test ./...`)
2. ✅ Linter passes (`golangci-lint run`)
3. ✅ SCAD output identical to pre-refactor
4. ✅ No yaml.Marshal/Unmarshal in module.go files
5. ✅ All modules use generated Decode()
6. ✅ All modules return []ArrayDecl
7. ✅ Documentation updated
8. ✅ Can add new module following guide

**Metrics:**

- **Lines removed:** ~84 lines (12 marshal/unmarshal sites × 7 lines each)
- **Lines added:** ~100 lines (helpers) + ~80 lines (templates) + ~140 lines (wrappers) = ~320 lines
- **Net change:** +236 lines (but eliminates boilerplate for all future modules)
- **Modules updated:** 7
- **Files changed:** ~30 files

---

## References

### Debate Decisions

- [Round 1: Go/No-Go](../debate/02-debate-round-1-go-no-go-builder-contract-refactor-urgency.md) - Decision to refactor now
- [Round 2: Path Choice](../debate/03-debate-round-2-path-choice-a-b-vs-a-c.md) - User chose Path 1 (A→B)
- [Round 4: Special Cases](../debate/05-debate-round-4-special-cases-how-do-we-retire-cutoutfeaturemodule.md) - Cutouts pattern
- [Round 5: Codegen Scope](../debate/06-debate-round-5-codegen-scope-what-to-generate-for-path-1-a-b.md) - Generate Decode() only
- [Round 6: Enforceability](../debate/07-debate-round-6-enforceability-keeping-modules-honest-post-refactor.md) - Four-layer enforcement
- [Round 7: Future Extensibility](../debate/08-debate-round-7-future-extensibility-accommodating-new-output-shapes.md) - ArrayDecl handles future needs
- [Round 8: Type Safety](../debate/09-debate-round-8-type-safety-end-to-end-arraydecl-vs-typed-model.md) - Untyped Model is sufficient

### Code References

- `pkg/yappgen/decode/helpers.go` - Shared decode helpers
- `pkg/schemagen/templates/schema_gen.go.tmpl` - Decode() generation
- `pkg/registry/schema.go` - ArrayDecl definition
- `pkg/yappgen/features.go` - multiArrayFeatureModule helper
- `pkg/yappgen/modules/*/module.go` - Module Build() functions
- `pkg/yappgen/modules/*/registry.go` - ArrayDecl wrappers

---

**Ready to implement!** Follow the checklist, run tests frequently, and refer back to debate rounds for decision rationale.
