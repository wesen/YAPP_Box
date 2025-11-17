---
Title: Debate Round 3 — Developer Ergonomics: Which Path Reduces Friction Fastest?
Ticket: YAPP-DSL-GAPS-001
Status: active
Topics:
    - yapp
    - architecture
    - codegen
    - debate
    - dx
DocType: debate
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/03-debate-round-2-path-choice-a-b-vs-a-c.md
      Note: Round 2 (Path choice)
    - Path: /home/manuel/code/others/YAPP_Box/pkg/docs/tutorials/yapp-module-authoring-guide.md
      Note: Current module authoring guide
ExternalSources: []
Summary: Which approach most improves the module authoring experience? Constraint - single-shot refactor, no phasing.
LastUpdated: 2025-11-17
---

# Debate Round 3 — Developer Ergonomics: Which Path Reduces Friction Fastest?

## Question

**Which path reduces developer friction fastest: Path 1 (A→B) or Path 2 (A→C)?**

**Context:** Current pain points:
- YAML marshal/unmarshal boilerplate in every builder
- Untyped map juggling (`map[string]any` → typed struct)
- Unclear patterns (pushbuttons uses custom decode, others use marshal/unmarshal)

**Constraint from Round 2:** Single-shot refactor, no phasing. Whichever path we choose, we implement it completely in one go.

## Pre-Debate Research

### Module Authoring Guide Complexity

```bash
$ wc -l pkg/docs/tutorials/yapp-module-authoring-guide.md
754 pkg/docs/tutorials/yapp-module-authoring-guide.md
```

**Finding:** 754-line guide. Let's see what it teaches:

```markdown
## Step 5: Implement Builder (lines 192-271)

Create `module.go` with your builder logic:

```go
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
    
    // Apply defaults
    item.ApplyDefaults()
    
    // Custom validation
    if err := item.CustomValidate(); err != nil {
        return nil, errors.Wrapf(err, "%s", label)
    }
    
    // Build positional array...
}
```
```

**Finding:** 80 lines of guide dedicated to explaining the marshal/unmarshal pattern. That's 10% of the entire guide.

### Boilerplate Error Wrapping

```bash
$ grep -r "errors.Wrapf.*marshal\|errors.Wrapf.*unmarshal" pkg/yappgen/modules/*/module.go | wc -l
12
```

**Finding:** 12 error wrapping sites for marshal/unmarshal across modules. Each one is 3-4 lines of boilerplate.

### Pushbuttons Exception

```go
// pkg/yappgen/modules/pushbuttons/module.go:14-18
func Build(items []map[string]any) ([][]any, error) {
	typedItems, err := decodePushButtonsItems(items)
	if err != nil {
		return nil, err
	}
	// ... build logic ...
}

// pkg/yappgen/modules/pushbuttons/module.go:108-150
func decodePushButtonsItems(items []map[string]any) ([]PushButtonsItem, error) {
	typed := make([]PushButtonsItem, len(items))
	for idx, raw := range items {
		label := fmt.Sprintf("push_buttons[%d]", idx)
		// Manual field extraction with type assertions
		// ... 40+ lines of field extraction ...
	}
	return typed, nil
}
```

**Finding:** Pushbuttons has a 40+ line custom decode function. It's more verbose than marshal/unmarshal but avoids the round-trip.

### Current Builder Structure

All non-pushbuttons modules follow this pattern:

```
1. Loop over items (1 line)
2. Format label (1 line)
3. Marshal to YAML (4 lines with error handling)
4. Unmarshal to struct (4 lines with error handling)
5. Apply defaults (1 line)
6. Custom validate (3 lines with error handling)
7. Build params array (10-30 lines)
8. Append to output (1 line)
```

**Total boilerplate per module:** ~15 lines (steps 1-6)
**Business logic per module:** ~10-30 lines (step 7)

---

## Opening Statements

### Sam Park — "The Codegen Maintainer"

Let me show you what each path looks like from a module author's perspective.

**Current state (no refactor):**

```go
func Build(items []map[string]any) ([][]any, error) {
	var out [][]any
	for idx, it := range items {
		label := fmt.Sprintf("pcb_stands[%d]", idx)
		
		// 8 lines of marshal/unmarshal boilerplate
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
		
		// 10-30 lines of actual business logic
		params := []any{item.X, item.Y, ptrOrUndef(item.Height), ...}
		out = append(out, params)
	}
	return out, nil
}
```

**Path 1 (A→B) with generated Decode:**

```go
func Build(items []map[string]any) ([][]any, error) {
	// Generated Decode handles marshal/unmarshal + defaults + validation
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}
	
	var out [][]any
	for _, item := range typed {
		// Just business logic
		params := []any{item.X, item.Y, ptrOrUndef(item.Height), ...}
		out = append(out, params)
	}
	return out, nil
}
```

**Path 2 (A→C) with typed Model:**

```go
func Build(items []PcbStandsItem) ([][]any, error) {
	// Items already typed and validated!
	var out [][]any
	for _, item := range items {
		// Just business logic
		params := []any{item.X, item.Y, ptrOrUndef(item.Height), ...}
		out = append(out, params)
	}
	return out, nil
}
```

**My analysis:**

- **Current:** 15 lines boilerplate + 10-30 lines logic = 25-45 lines total
- **Path 1:** 3 lines Decode call + 10-30 lines logic = 13-33 lines total (~50% reduction)
- **Path 2:** 0 lines boilerplate + 10-30 lines logic = 10-30 lines total (~60% reduction)

**But** Path 2 changes the signature from `[]map[string]any` to `[]PcbStandsItem`. That's a breaking change to the `Build` function signature.

**My position:** Path 1 has better ergonomics *for module authors* because they still write the same function signature. Path 2 is cleaner but requires changing every Build signature.

---

### Casey Thompson — "The New Hire"

*[Pulls up the module authoring guide]*

I tried to follow this guide last week. Here's where I got confused:

**Step 1:** Write schema.yaml ✅ (makes sense)

**Step 2:** Run schemagen ✅ (generates struct)

**Step 3:** Write Build function... and then I see this:

```go
// Unmarshal into typed struct
data, err := yaml.Marshal(it)
```

Wait, what? I just got `it` from the caller. Why am I marshaling it back to YAML? Isn't that... backwards?

I asked Sam, and he said "it's because `it` is `map[string]any` and we need a typed struct." But... why not just extract the fields directly?

Then I looked at pushbuttons and saw `decodePushButtonsItems`. So there *is* a way to avoid marshal/unmarshal! But the guide doesn't mention it. And pushbuttons' decode function is 40+ lines of manual field extraction.

**My confusion:**

1. **Two patterns:** Marshal/unmarshal vs. manual decode. Which should I use?
2. **Guide doesn't explain why:** Just says "do this" without explaining the tradeoff.
3. **Inconsistency:** 6 modules use marshal/unmarshal, 1 uses manual decode.

**What I want:**

- **One pattern** that works for all modules
- **Generated code** so I don't have to choose
- **Clear signature** so I know what types to expect

**My position:** I don't care about Path 1 vs. Path 2. I care about *consistency*. Right now, I have to read 6 different module.go files to figure out the "right" pattern. If either path gives me generated Decode, I'm happy.

But Path 2's typed signature (`[]PcbStandsItem`) is *way* clearer than Path 1's (`[]map[string]any`). When I see `[]PcbStandsItem`, I know exactly what I'm working with. When I see `[]map[string]any`, I have to read the schema to understand the structure.

---

### The Orchestrator (`features.go`)

*[Speaks up from the integration layer]*

Let me tell you about the pain I experience with the current system.

**My job:** Wire modules together. Here's what I do today:

```go
// features.go:97-120
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

I have to:
1. Extract `[]any` from features
2. Normalize to `[]map[string]any`
3. Store in Model
4. Later, pass to Build which does marshal/unmarshal

**That's three type conversions:**
- `[]any` → `[]map[string]any` (in Collect)
- `[]map[string]any` → YAML bytes (in Build)
- YAML bytes → typed struct (in Build)

**Path 1 (A→B) reduces it to two conversions:**
- `[]any` → `[]map[string]any` (in Collect)
- `[]map[string]any` → typed struct (generated Decode)

**Path 2 (A→C) reduces it to one conversion:**
- `[]any` → typed struct (generated Decode in Collect)

**My position:** Path 2 eliminates two unnecessary conversions. But it requires changing my `field` accessor signature from `func(*Model) *[]map[string]any` to `func(*Model) *[]PcbStandsItem`.

That's a breaking change to `newArrayFeatureModule`. But since we're doing single-shot refactor, that's fine. We change all 6 call sites at once.

**Ergonomics impact:** Path 2 makes my code simpler. No more `normalizeArrayOfMaps`. No more untyped storage. Just: extract from features, decode, store.

---

## Rebuttals

### Sam Park — "The Codegen Maintainer"

Casey, you've hit on the core issue: **inconsistency is the real enemy**.

I showed Path 1 as "better for module authors" because it doesn't change the Build signature. But you're right—that's actually *worse* for clarity. When you see:

```go
func Build(items []map[string]any) ([][]any, error)
```

You have no idea what's in those maps without reading the schema. But with Path 2:

```go
func Build(items []PcbStandsItem) ([][]any, error)
```

The signature *documents itself*. You know exactly what you're working with.

**I'm changing my position.** Path 2 has better ergonomics because:

1. **Self-documenting signatures:** Type names tell you what's expected
2. **Compile-time checks:** Pass wrong type → compiler error, not runtime panic
3. **IDE support:** Autocomplete works on `PcbStandsItem`, not on `map[string]any`

The Orchestrator's point about conversions is also compelling. Why do three conversions when we can do one?

**Revised position:** Path 2 (A→C). Yes, it changes Build signatures. But that's a *good* change. It makes the code clearer and safer.

---

### Casey Thompson — "The New Hire"

*[Nods enthusiastically]*

Sam, thank you! That's exactly what I was trying to say.

Let me show you what the module authoring guide would look like with each path.

**Path 1 (A→B) guide:**

```markdown
## Step 5: Implement Builder

Create `module.go`:

```go
func Build(items []map[string]any) ([][]any, error) {
	typed, err := Decode(items)  // Generated function
	if err != nil {
		return nil, err
	}
	
	var out [][]any
	for _, item := range typed {
		params := []any{item.X, item.Y, ...}
		out = append(out, params)
	}
	return out, nil
}
```
```

**Path 2 (A→C) guide:**

```markdown
## Step 5: Implement Builder

Create `module.go`:

```go
func Build(items []YourModuleItem) ([][]any, error) {
	var out [][]any
	for _, item := range items {
		params := []any{item.X, item.Y, ...}
		out = append(out, params)
	}
	return out, nil
}
```
```

**Path 2 is shorter and clearer.** No need to explain Decode because it's already been called. No need to explain error handling because validation happened before Build.

**Module authoring becomes:**
1. Write schema.yaml
2. Run schemagen (generates struct + Decode)
3. Write Build logic (just the params array)

That's it. No marshal/unmarshal explanation. No "why do we do this" confusion.

---

### The Orchestrator (`features.go`)

Let me address the "breaking change" concern.

**Path 1 changes:**
- Add `Decode()` generation (schemagen templates)
- Module authors update Build to call Decode
- No changes to my code

**Path 2 changes:**
- Add `Decode()` generation (schemagen templates)
- Change Model field types (7 fields)
- Change `newArrayFeatureModule` signature (6 call sites)
- Module authors update Build signature + remove Decode call

**Lines of code changed:**
- Path 1: ~70 lines (7 modules × 10 lines each)
- Path 2: ~120 lines (7 modules × 10 lines + Model + features.go)

**But** Path 2 eliminates more code than it adds:

- Remove `normalizeArrayOfMaps` calls (6 call sites)
- Remove marshal/unmarshal boilerplate (12 error wrapping sites)
- Remove untyped field accessors (6 functions)

**Net change:**
- Path 1: +70 lines, -12 lines = +58 lines
- Path 2: +120 lines, -50 lines = +70 lines

Path 2 is only 12 more lines of change, but it eliminates three type conversions and gives us compile-time safety.

**My revised position:** Path 2. The ergonomics win (self-documenting signatures, fewer conversions, compile-time checks) outweighs the slightly larger change scope.

---

## Moderator Summary

### Key Arguments

**For Path 1 (A→B):**
- Smaller change scope (~58 net lines)
- Doesn't change Build signatures (less breaking)
- Module authors still write `Build([]map[string]any)`

**For Path 2 (A→C):**
- Self-documenting signatures (`Build([]PcbStandsItem)`)
- Eliminates two type conversions (three → one)
- Better IDE support (autocomplete on typed structs)
- Shorter module authoring guide (no Decode explanation needed)
- Compile-time type safety

### Tensions

1. **Change scope vs. clarity:** Path 1 changes less code; Path 2 produces clearer code.

2. **Breaking changes:** Path 1 avoids signature changes; Path 2 embraces them for better types.

3. **Documentation burden:** Path 1 still requires explaining Decode; Path 2 makes it invisible.

### Position Changes

- **Sam:** Started favoring Path 1 (smaller change), switched to Path 2 (better clarity)
- **Casey:** Neutral on paths, just wants consistency; convinced by Path 2's self-documenting signatures
- **Orchestrator:** Favored Path 2 from start; showed it eliminates conversions

### Emerging Consensus

**All three candidates now favor Path 2 (A→C)** because:
- Self-documenting function signatures
- Fewer type conversions (better performance, less complexity)
- Better developer experience (IDE autocomplete, compile-time checks)
- Simpler module authoring guide

**Trade-off accepted:** Slightly larger change scope (70 vs. 58 lines) is worth the ergonomics win.

### Interesting Ideas

- **Casey's guide comparison:** Showed Path 2 guide is shorter and clearer
- **Orchestrator's conversion count:** Three conversions → one conversion
- **Sam's signature clarity argument:** `[]PcbStandsItem` is self-documenting

### Open Questions

1. How do we handle the single-shot refactor? (All 7 modules + Model + features.go in one commit)
2. What's the test strategy to ensure we don't break anything?
3. Do we need a migration checklist?

---

## Decision Point

**Strong consensus for Path 2 (A→C)** based on developer ergonomics:
- Clearer code (self-documenting signatures)
- Better tooling (IDE autocomplete)
- Simpler onboarding (shorter guide, fewer concepts)

**Next debate:** How do we handle cutouts' special-case in Path 2? (Question 5)
