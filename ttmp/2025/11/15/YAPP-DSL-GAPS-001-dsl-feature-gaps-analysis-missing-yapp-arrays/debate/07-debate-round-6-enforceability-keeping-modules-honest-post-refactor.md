---
Title: Debate Round 6 — Enforceability: Keeping Modules Honest Post-Refactor
Ticket: YAPP-DSL-GAPS-001
Status: active
Topics:
    - yapp
    - architecture
    - codegen
    - debate
    - testing
DocType: debate
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/06-debate-round-5-codegen-scope-what-to-generate-for-path-1-a-b.md
      Note: Round 5 (Codegen scope)
    - Path: /home/manuel/code/others/YAPP_Box/pkg/docs/tutorials/yapp-module-authoring-guide.md
      Note: Module authoring guide
ExternalSources: []
Summary: How do we prevent regression to old patterns after Path 1 (A→B) refactor?
LastUpdated: 2025-11-17
---

# Debate Round 6 — Enforceability: Keeping Modules Honest Post-Refactor

## Question

**How do we prevent regression to old patterns after implementing Path 1 (A→B)?**

**Context:** After refactor, we'll have:
- Generated `Decode()` functions that should be used
- Old marshal/unmarshal pattern that should be deprecated
- ArrayDecl IR that modules should return

**Risk:** Future module authors might:
- Copy old code that uses marshal/unmarshal
- Forget to use generated Decode()
- Return `[][]any` instead of `[]ArrayDecl`

## Pre-Debate Research

### Current Enforcement Mechanisms

```bash
$ find . -name ".golangci.yml" -o -name "golangci.yml"
# (no output)
```

**Finding:** No linter configuration. No automated enforcement of patterns.

### Current TODO Markers

```bash
$ grep -r "TODO:" pkg/yappgen/modules/*/registry.go | wc -l
15
```

**Finding:** 15 TODO markers in registry files. Most are:
- `TODO: implement constraint validation`
- `TODO: return field specs from schema`

These TODOs have been there since module creation. No enforcement to complete them.

### Test Coverage

```bash
$ find pkg/yappgen/modules -name "*_test.go" -exec grep -l "func Test.*Build" {} \;
pkg/yappgen/modules/boxmounts/module_test.go
pkg/yappgen/modules/snapjoins/module_test.go
pkg/yappgen/modules/connectors/module_test.go
pkg/yappgen/modules/pcbstands/module_test.go
```

**Finding:** 4 out of 7 modules have Build tests. Pushbuttons, lighttubes, cutouts don't have builder tests.

### Documentation Patterns

```bash
$ grep -n "marshal\|unmarshal" pkg/docs/tutorials/yapp-module-authoring-guide.md | head -7
215:        // Unmarshal into typed struct
218:            return nil, errors.Wrapf(err, "%s: marshal", label)
222:        if err := yaml.Unmarshal(data, &item); err != nil {
223:            return nil, errors.Wrapf(err, "%s: unmarshal", label)
642:            return nil, errors.Wrapf(err, "%s: marshal", label)
646:        if err := yaml.Unmarshal(data, &item); err != nil {
647:            return nil, errors.Wrapf(err, "%s: unmarshal", label)
```

**Finding:** Module authoring guide shows marshal/unmarshal pattern in 2 places (lines 215-223 and 642-647). After refactor, this will be outdated.

### Interface Definitions

```go
// pkg/registry/schema.go:53-58
type FeatureModule interface {
	Schema() ModuleSchema
	Build(items []map[string]any) ([][]any, error)
}

// pkg/yappgen/features.go:20-24
type FeatureModule interface {
	Name() string
	Collect(resolved map[string]any, features map[string]any, model *Model) error
	Emit(ctx context.Context, model *Model, b *strings.Builder) error
}
```

**Finding:** Two different `FeatureModule` interfaces in different packages. After ArrayDecl, `registry.FeatureModule.Build` will return `[]ArrayDecl`, not `[][]any`.

---

## Opening Statements

### The Contract (`registry/schema.go`)

Let me talk about **interface-level enforcement**.

**After Path 1 (A→B), the interface changes:**

```go
// Before
type FeatureModule interface {
	Schema() ModuleSchema
	Build(items []map[string]any) ([][]any, error)
}

// After
type FeatureModule interface {
	Schema() ModuleSchema
	Build(items []map[string]any) ([]ArrayDecl, error)
}
```

**This is compile-time enforcement.** If a module returns `[][]any`, it won't compile. Problem solved, right?

**Wrong.** Here's what a lazy module author might do:

```go
// "Compliant" but wrong
func (m *module) Build(items []map[string]any) ([]ArrayDecl, error) {
	// Still using marshal/unmarshal!
	data, _ := yaml.Marshal(items[0])
	var item MyItem
	yaml.Unmarshal(data, &item)
	
	// Wrap in ArrayDecl to satisfy interface
	rows := [][]any{{item.X, item.Y}}
	return []ArrayDecl{{Name: "myArray", Rows: rows}}, nil
}
```

They've satisfied the interface but **ignored the generated Decode()**.

**My position:** Interface changes enforce the return type, but they don't enforce *how* you build that return value. We need additional enforcement.

**Proposal 1: Deprecation markers**

```go
// pkg/yappgen/decode/deprecated.go
// Deprecated: Use generated Decode() instead of marshal/unmarshal.
// This function will be removed in v2.0.
func MarshalUnmarshalDecode(items []map[string]any) error {
	return errors.New("marshal/unmarshal pattern is deprecated; use generated Decode()")
}
```

If someone tries to use `yaml.Marshal`, they get a deprecation warning.

**Proposal 2: Linter rules**

Add golangci-lint rules:
- Ban `yaml.Marshal` in module.go files
- Require `Decode()` call in Build functions

---

### Sam Park — "The Codegen Maintainer"

Contract, your interface enforcement is good but incomplete. Let me show you the **testing gap**.

**Current state:** 4 out of 7 modules have Build tests. After refactor, we need to ensure:
1. Decode() is tested
2. Build() uses Decode()
3. ArrayDecl is returned correctly

**Proposal: Generate tests for Decode()**

When we generate Decode(), also generate tests:

```go
// Generated in schema_gen_test.go
func TestDecode_ValidInput(t *testing.T) {
	items := []map[string]any{
		{"x": 10.0, "y": 20.0, "height": 5.0},
	}
	
	typed, err := Decode(items)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	
	if len(typed) != 1 {
		t.Fatalf("expected 1 item, got %d", len(typed))
	}
	
	if typed[0].X != 10.0 {
		t.Errorf("expected X=10.0, got %v", typed[0].X)
	}
	// ... more assertions
}

func TestDecode_MissingRequired(t *testing.T) {
	items := []map[string]any{
		{"x": 10.0}, // missing required 'y'
	}
	
	_, err := Decode(items)
	if err == nil {
		t.Fatal("expected error for missing required field")
	}
	
	if !strings.Contains(err.Error(), "y: missing required field") {
		t.Errorf("wrong error: %v", err)
	}
}
```

**Benefits:**
1. Decode() is tested automatically
2. If someone breaks Decode(), tests fail
3. Error message validation ensures good UX

**Enforcement:** CI fails if generated tests don't pass. Can't merge broken Decode().

**My position:** Generate tests for Decode(). This ensures the generated code works and gives module authors confidence.

---

### Casey Thompson — "The New Hire"

*[Raises hand]*

I want to talk about **documentation enforcement**.

**Problem:** After refactor, the module authoring guide will be outdated. Lines 215-223 and 642-647 show the marshal/unmarshal pattern. If I follow the guide, I'll write deprecated code.

**Proposal: Update guide with new pattern**

```markdown
## Step 5: Implement Builder (UPDATED)

Create `module.go`:

```go
func Build(items []map[string]any) ([]ArrayDecl, error) {
	// Use generated Decode() - DO NOT use yaml.Marshal/Unmarshal
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}
	
	var rows [][]any
	for _, item := range typed {
		params := []any{item.X, item.Y, ptrOrUndef(item.Height), ...}
		rows = append(rows, params)
	}
	
	// Return ArrayDecl (not [][]any)
	return []ArrayDecl{{Name: "yourModule", Rows: rows}}, nil
}
```

**DO NOT:**
- ❌ Use `yaml.Marshal(it)` / `yaml.Unmarshal(data, &item)`
- ❌ Return `[][]any` directly
- ❌ Copy code from pre-refactor modules

**DO:**
- ✅ Use generated `Decode()` function
- ✅ Return `[]ArrayDecl`
- ✅ Follow this guide's pattern
```
```

**Enforcement idea:** Add a checklist to the guide:

```markdown
## Pre-Merge Checklist

Before submitting your module PR:
- [ ] Used generated Decode() (no marshal/unmarshal)
- [ ] Returned []ArrayDecl (not [][]any)
- [ ] Added tests for Build()
- [ ] Ran `go test ./...` successfully
- [ ] Ran `go run ./cmd/schemagen discover` to regenerate
```

**My position:** Documentation is the first line of defense. If the guide shows the right pattern, most people will follow it.

---

## Rebuttals

### The Contract (`registry/schema.go`)

Casey, your documentation proposal is excellent. But documentation alone isn't enough—people don't always read it.

**Combining enforcement layers:**

1. **Compile-time:** Interface change (ArrayDecl return type)
2. **Test-time:** Generated Decode tests (Sam's proposal)
3. **Lint-time:** Golangci-lint rules (my proposal)
4. **Review-time:** PR checklist (Casey's proposal)

**Let me detail the linter rules:**

```yaml
# .golangci.yml (new file)
linters:
  enable:
    - forbidigo  # Forbid specific function calls
    
linters-settings:
  forbidigo:
    forbid:
      - pattern: 'yaml\.Marshal'
        msg: 'Do not use yaml.Marshal in module builders; use generated Decode() instead'
        pkg: '^github\.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/.*/module\.go$'
      
      - pattern: 'yaml\.Unmarshal'
        msg: 'Do not use yaml.Unmarshal in module builders; use generated Decode() instead'
        pkg: '^github\.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/.*/module\.go$'
```

**This catches:**
- Direct yaml.Marshal/Unmarshal calls in module.go files
- Fails CI if someone tries to use the old pattern

**Exception:** Cutouts and pushbuttons might have legitimate uses. We can add `//nolint:forbidigo` with a comment explaining why.

**My revised position:** Four-layer enforcement (compile, test, lint, review).

---

### Sam Park — "The Codegen Maintainer"

Contract, I love the linter idea. Let me add **test enforcement details**.

**Problem:** Generating tests is good, but how do we ensure Build() *uses* Decode()?

**Proposal: Generate Build() wrapper test**

```go
// Generated in schema_gen_test.go
func TestBuild_UsesGeneratedDecode(t *testing.T) {
	// This test verifies Build() uses Decode() by checking error messages
	items := []map[string]any{
		{"x": "not a number"}, // Type error
	}
	
	_, err := Build(items)
	if err == nil {
		t.Fatal("expected error for invalid type")
	}
	
	// If Build() uses Decode(), error message will contain "expected number"
	if !strings.Contains(err.Error(), "expected number") {
		t.Errorf("Build() may not be using Decode(); error: %v", err)
	}
}
```

**This test fails if:**
- Build() doesn't call Decode()
- Build() uses marshal/unmarshal (different error messages)

**Additional test: ArrayDecl structure**

```go
func TestBuild_ReturnsArrayDecl(t *testing.T) {
	items := []map[string]any{
		{"x": 10.0, "y": 20.0},
	}
	
	decls, err := Build(items)
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

**My revised position:** Generate three test types:
1. Decode() validation tests
2. Build() uses Decode() test (error message check)
3. ArrayDecl structure tests

---

### Casey Thompson — "The New Hire"

*[Nods enthusiastically]*

This is great! But I have one more concern: **What about existing modules?**

After we refactor, we'll have:
- 7 modules updated to use Decode() + ArrayDecl
- Module authoring guide updated
- Linter rules in place

But what if someone looks at an old commit or an old branch? They'll see the marshal/unmarshal pattern and think it's still valid.

**Proposal: Migration guide**

Create a document: `docs/migrations/path-1-a-b-refactor.md`

```markdown
# Migration Guide: Path 1 (A→B) Refactor

## What Changed (2025-11-17)

We refactored the module builder system:
- ✅ Generated Decode() functions (no more marshal/unmarshal)
- ✅ ArrayDecl IR (unified single/multi-array handling)
- ✅ Linter rules (enforce new patterns)

## If You're Working on an Old Branch

**Before merging to main:**
1. Rebase on latest main
2. Update your module's Build() to use Decode()
3. Change return type from `[][]any` to `[]ArrayDecl`
4. Run `go run ./cmd/schemagen discover` to regenerate
5. Run `go test ./...` to verify

## If You're Adding a New Module

Follow the updated [Module Authoring Guide](../tutorials/yapp-module-authoring-guide.md).

**DO NOT:**
- Copy code from commits before 2025-11-17
- Use yaml.Marshal/Unmarshal
- Return [][]any directly

## If You See Linter Errors

```
module.go:25: Do not use yaml.Marshal in module builders
```

**Fix:** Replace marshal/unmarshal with generated Decode():

```go
// Before
data, _ := yaml.Marshal(it)
var item MyItem
yaml.Unmarshal(data, &item)

// After
typed, err := Decode(items)
if err != nil {
    return nil, err
}
item := typed[0]
```
```

**My position:** Document the migration explicitly so future contributors know the old pattern is deprecated.

---

## Moderator Summary

### Key Enforcement Mechanisms

**1. Compile-time (Interface)**
- Change `Build()` return type to `[]ArrayDecl`
- Won't compile if returning `[][]any`
- Enforces output format, not implementation

**2. Test-time (Generated Tests)**
- Generate Decode() validation tests
- Generate Build() uses Decode() test (error message check)
- Generate ArrayDecl structure tests
- CI fails if tests don't pass

**3. Lint-time (Golangci-lint)**
- Forbid `yaml.Marshal` in module.go files
- Forbid `yaml.Unmarshal` in module.go files
- CI fails if linter catches violations

**4. Review-time (Documentation + Checklist)**
- Update module authoring guide with new pattern
- Add pre-merge checklist
- Create migration guide for old branches

### Consensus Points

All three candidates agree on:
- **Multi-layer enforcement** (compile + test + lint + review)
- **Generated tests** for Decode() and Build()
- **Documentation updates** (guide + migration doc)
- **Linter rules** to catch deprecated patterns

### Implementation Plan

**Phase 1: Code changes**
1. Implement Path 1 (A→B) refactor
2. Generate Decode() + ArrayDecl wrapper
3. Update all 7 modules

**Phase 2: Enforcement infrastructure**
1. Add `.golangci.yml` with forbidigo rules
2. Generate Decode() tests in schema_gen_test.go
3. Generate Build() wrapper tests
4. Update CI to run linter

**Phase 3: Documentation**
1. Update module authoring guide (remove marshal/unmarshal)
2. Add pre-merge checklist
3. Create migration guide
4. Update README with refactor notes

**Phase 4: Validation**
1. Run linter on all modules (should pass)
2. Run all tests (should pass)
3. Try adding a test module following new guide
4. Verify enforcement catches violations

### Open Questions

1. Should we add a pre-commit hook to run the linter?
2. Do we need a "deprecated patterns" document listing old code to avoid?
3. Should generated tests be in a separate file (schema_gen_test.go) or module_test.go?

---

## Decision Point

**Consensus: Four-layer enforcement (compile + test + lint + review)**

**Implementation priority:**
1. Interface change (compile-time) - Part of refactor
2. Generated tests (test-time) - Add to schemagen templates
3. Documentation (review-time) - Update guide + add migration doc
4. Linter rules (lint-time) - Add .golangci.yml

**Next debate:** Future extensibility (Question 10) - Does Path 1 handle unknown future needs?
