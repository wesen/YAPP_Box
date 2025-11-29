---
Title: Debate Round 3 - Pipeline Integration: Where Does It Hook In?
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
Summary: Third debate round exploring where in the DSL pipeline composite modules should execute
LastUpdated: 2025-11-28
---

# Debate Round 3: Pipeline Integration - Where Does It Hook In?

## Question

At what point in the DSL processing pipeline should composite modules execute: before resolver, after resolver but before BuildModel, or after BuildModel?

## Primary Candidates

- `pkg/resolver/resolver.go` "The Expression Evaluator" (knows pipeline flow)
- Dr. Sarah "The Architect" Martinez (wants clear boundaries)
- Morgan "The Performance Engineer" Taylor (cares about processing order)
- `pkg/yappgen/features.go` "The Feature Coordinator" (current integration point)

---

## Pre-Debate Research

### Research 1: Current Pipeline Flow in CLI

**Command:**
```bash
grep -r "LoadAndResolveResult\|WriteSCAD" cmd/yappctl/
```

**Pipeline Flow from `cmd/yappctl/generate_command.go`:**

```
Line 144: loadResult, err := resolvercli.LoadAndResolveResult(...)
   ↓ (returns resolved Document, Trace, Comments, Raw)
Line 155: scadPath, model, err := generatorcli.WriteSCAD(...)
   ↓ (BuildModel → EmitSCAD)
```

**Detailed Flow:**
```
1. LoadAndResolveResult (pkg/cli/resolvercli/)
   ├─ Read YAML file
   ├─ resolver.ResolveResult()
   │   ├─ Validate structure (Phase 1)
   │   ├─ Resolve expressions (iterative)
   │   └─ Validate constraints (Phase 2)
   └─ Return: Document, Trace, Comments, Raw

2. WriteSCAD (pkg/cli/generatorcli/)
   ├─ yappgen.BuildModel(resolved, trace, comments, raw)
   │   ├─ Parse resolved YAML into Model
   │   └─ collectFeatureModules() → populate Model arrays
   ├─ yappgen.EmitSCAD(model)
   │   └─ emitFeatureModules() → write SCAD arrays
   └─ Write SCAD file
```

**Potential Hook Points Identified:**

**Hook Point A:** Before `LoadAndResolveResult`
- Raw YAML available
- No expressions resolved
- No validation done

**Hook Point B:** After `LoadAndResolveResult`, before `WriteSCAD`
- Expressions resolved
- Validation complete
- Can modify resolved document

**Hook Point C:** Inside `BuildModel`, after `collectFeatureModules`
- Model populated
- Can modify Model arrays directly

**Hook Point D:** Inside `EmitSCAD`, during emission
- Arrays ready
- Can modify during emission

---

### Research 2: ResolveResult Function Analysis

**Code:**
```go
func ResolveResult(ctx context.Context, doc map[string]any, opts Options) (*Result, error) {
    state := deepCopy(doc).(map[string]any)
    
    // Phase 1: Validate structure before expression resolution
    if err := validateStructure(state); err != nil {
        return nil, errors.Wrap(err, "structure validation")
    }
    
    // Expression resolution (iterative, up to 16 passes)
    for iter := 0; iter < opts.MaxIterations; iter++ {
        resolved, iterChanged, err := resolvePass(state, usedVars, trace)
        // ...
    }
    
    // Phase 2: Validate constraints after expression resolution
    if err := validateConstraints(state); err != nil {
        return nil, errors.Wrap(err, "constraint validation")
    }
    
    return &Result{Document: state, Trace: trace.Trace()}, nil
}
```

**Key Observations:**
- Structure validation happens **before** resolution
- Constraint validation happens **after** resolution
- Resolver mutates a deep copy (doesn't affect original)

**Implication for Hook Points:**
- Hook A (before resolver): Validation will see composite-generated DSL
- Hook B (after resolver): Composite modules work with resolved values
- Hook C (in BuildModel): Bypass validation for generated entries

---

### Research 3: BuildModel Hook Point Analysis

**Code:**
```go
func BuildModel(ctx context.Context, resolved map[string]any, ...) (*Model, error) {
    m := &Model{...}
    
    // Parse PCB, enclosure configs
    // ...
    
    // Collect features
    features, _ := getMap(resolved, "features")
    if err := collectFeatureModules(resolved, features, m); err != nil {
        return nil, err
    }
    
    return m, nil
}
```

**Potential Hook Point:**
```go
// After collectFeatureModules, before return
if err := processCompositeModules(resolved, features, m); err != nil {
    return nil, err
}
```

**Lines of Code to Add:** ~3 lines in BuildModel

---

### Research 4: Performance Impact Analysis

**Current Processing Time (estimated from pipeline):**
- Parse YAML: ~5ms
- Resolve expressions: ~10-20ms (depends on complexity)
- Build Model: ~5-10ms
- Emit SCAD: ~5-10ms
- **Total: ~25-45ms**

**Composite Module Overhead Estimates:**

**Hook A (before resolver):**
- Parse YAML: 5ms
- **Composite transform: ~5-10ms**
- Resolve (transformed YAML): 10-20ms (might be longer)
- Build/Emit: 10-20ms
- **Total: ~30-55ms** (+5-10ms overhead)

**Hook B (after resolver):**
- Parse/Resolve: 15-25ms
- **Composite transform: ~5-10ms**
- Build/Emit: 10-20ms
- **Total: ~30-55ms** (+5-10ms overhead)

**Hook C (in BuildModel):**
- Parse/Resolve: 15-25ms
- Build (with composite): 15-30ms (+5-10ms)
- Emit: 5-10ms
- **Total: ~35-65ms** (+10-20ms overhead)

**Finding:** Hook point doesn't significantly affect performance (all <10ms overhead)

---

### Research 5: Validation Flow Analysis

**From `pkg/resolver/validation.go`:**

```go
func validateStructure(doc map[string]any) error {
    features, ok := doc["features"].(map[string]any)
    for key, value := range features {
        path := "features." + key
        module, found := registry.Get(path)
        if !found {
            continue // Unknown module - skip
        }
        schema := module.Schema()
        if err := schema.ValidateStructure(path, value); err != nil {
            return err
        }
    }
    return nil
}
```

**Key Finding:** Validation looks up modules in `registry.Get(path)`

**Implications:**

**Hook A (before resolver):**
- Composite-generated DSL goes through normal validation ✅
- Regular modules validate generated entries ✅
- **Clean validation flow**

**Hook B (after resolver):**
- Need to re-validate generated DSL
- Or: Trust composite modules to generate valid DSL
- **Additional validation needed**

**Hook C (in BuildModel):**
- Generated arrays bypass validation
- Need separate validation in composite modules
- **Validation duplication**

---

## Opening Statements (Round 1)

### `pkg/resolver/resolver.go` - "The Expression Evaluator"

*[Analyzes pipeline flow, integration points]*

I've traced the complete pipeline from CLI to SCAD output. Let me explain the options.

**The pipeline has 3 clear stages:**
1. **Pre-Resolution** (Hook A): Raw YAML → Transform → Resolved YAML
2. **Post-Resolution** (Hook B): Resolved YAML → Transform → Modified YAML
3. **Post-BuildModel** (Hook C): Model arrays → Transform → Modified Model

**My recommendation: Hook A (before me).**

**Why?**

**1. Natural Expression Flow**
If composite modules generate DSL at Hook A:
```yaml
features:
  lcd:
    - display:
        position: [pcb.length / 2, pcb.width / 2]  # Expressions!
```
→ Transforms to:
```yaml
features:
  cutouts:
    - from_face_left: pcb.length / 2  # Still expressions
      from_face_bottom: pcb.width / 2
```
→ I resolve all expressions together

**This is clean.** Composite modules generate DSL with expressions, I resolve everything in one pass.

**2. Validation Flow**
At Hook A, validation happens **after** transformation:
- Structure validation sees transformed DSL
- Regular module validators check generated entries
- **No duplicate validation**

**3. Implementation Location**
Hook A fits in CLI layer (`pkg/cli/generatorcli/` or `pkg/cli/resolvercli/`):
```go
func LoadAndResolveResult(...) (*LoadResult, error) {
    doc, _ := parseYAML(path)
    
    // NEW: Composite module processing
    transformed, err := composite.ProcessModules(doc)
    if err != nil {
        return nil, err
    }
    
    // Existing: Resolve
    result, err := resolver.ResolveResult(ctx, transformed, ...)
    return result, nil
}
```

**4-5 lines of integration code.**

**Hook B** is awkward - need to re-run resolution or trust composite modules.

**Hook C** is wrong - bypasses my validation, duplicates expression handling.

**My position: Hook A (before resolver).**

---

### Dr. Sarah "The Architect" Martinez

*[Reviews integration points, architectural boundaries]*

Resolver makes a good case for Hook A. Let me validate it architecturally.

**Clear Boundaries:**
```
┌─────────────────────────────────────┐
│   CLI Layer (cmd/yappctl)           │
│   - Parse input                     │
│   - Hook A: Composite Transform ←── │ NEW
│   - LoadAndResolveResult            │
│   - WriteSCAD                        │
└─────────────────────────────────────┘
         ↓
┌─────────────────────────────────────┐
│   Resolver (pkg/resolver)           │
│   - Validate structure              │
│   - Resolve expressions             │
│   - Validate constraints            │
└─────────────────────────────────────┘
         ↓
┌─────────────────────────────────────┐
│   Generator (pkg/yappgen)           │
│   - BuildModel                      │
│   - collectFeatureModules           │
│   - EmitSCAD                        │
└─────────────────────────────────────┘
```

**Hook A (before resolver):**
- ✅ CLI layer responsibility (orchestration)
- ✅ Clean boundary (transform, then resolve)
- ✅ Resolver unchanged (doesn't know about composites)
- ✅ Generator unchanged (doesn't know about composites)

**Hook C (in BuildModel):**
- ❌ Generator layer responsibility (not its job)
- ❌ Tight coupling (composites modify Model)
- ❌ Validation gap (generated arrays not validated)

**Architecturally, Hook A is the only correct choice.**

**But:** There's a subtlety. Hook A has two sub-options:

**Hook A1:** In `resolvercli.LoadAndResolveResult` (CLI utility layer)
**Hook A2:** In `cmd/yappctl/generate_command.go` (command handler)

**I prefer A1** - keeps composite processing **reusable** across different CLI commands.

**My position: Hook A1 (in resolvercli).**

---

### Morgan "The Performance Engineer" Taylor

*[Analyzes performance impact, processing order]*

Let me talk about performance and processing order.

**Performance Analysis:**

**Hook A (before resolver):**
- Composite transform: ~5-10ms
- Resolve (larger YAML): ~15-25ms (increased from 10-20ms)
- **Total overhead: ~10-15ms** (larger YAML to resolve)

**Hook C (in BuildModel):**
- Composite transform: ~5-10ms
- No additional resolution overhead
- **Total overhead: ~5-10ms**

**Hook C is 2x faster.** But is that meaningful?

**Current pipeline: ~30-50ms**
**With Hook A: ~40-65ms** (+33% overhead)
**With Hook C: ~35-60ms** (+16% overhead)

**For user experience:**
- Both are **imperceptible** (humans don't notice <100ms)
- Both are **fine** for development workflow

**So performance doesn't decide this.** It's about architecture.

**But here's what I care about:** **Processing order.**

**Hook A:**
```
Transform → Validate → Resolve → Build → Emit
```
- Composite-generated DSL is validated early ✅
- Errors caught before resolution ✅
- **Fail fast** ✅

**Hook C:**
```
Validate → Resolve → Build → Transform → Emit
```
- Composite-generated arrays bypass validation ❌
- Errors caught late ❌
- **Fail slow** ❌

**From a debugging perspective, Hook A is better.** Fail fast, clear error messages, validated DSL.

**My position: Hook A (before resolver).**

Performance is fine either way, but **fail-fast** is important.

---

### `pkg/yappgen/features.go` - "The Feature Coordinator"

*[Analyzes integration impact on own code]*

Let me be blunt: **I don't want Hook C.**

**Why?** Look at my `BuildModel` integration:

**Current code:**
```go
func BuildModel(...) (*Model, error) {
    m := &Model{...}
    
    // Parse configs
    // ...
    
    // Collect features
    features, _ := getMap(resolved, "features")
    if err := collectFeatureModules(resolved, features, m); err != nil {
        return nil, err
    }
    
    return m, nil
}
```

**With Hook C:**
```go
func BuildModel(...) (*Model, error) {
    m := &Model{...}
    
    // Parse configs
    // ...
    
    // Collect features
    features, _ := getMap(resolved, "features")
    if err := collectFeatureModules(resolved, features, m); err != nil {
        return nil, err
    }
    
    // NEW: Process composite modules
    if err := processCompositeModules(resolved, features, m); err != nil {
        return nil, err
    }
    
    return m, nil
}
```

**That's fine.** 3 lines added. But here's my concern:

**Composite modules modify my Model.** They append to arrays I just populated. What if:
- They generate invalid entries?
- They conflict with regular module entries?
- They break my invariants?

**I have no control.** I populate arrays, then composite modules modify them.

**Hook A is better for me:** Composite modules transform DSL **before** it reaches me. I process transformed DSL **as if the user wrote it**. I never know the difference.

**Hook A: Composite modules are invisible to me** ✅
**Hook C: Composite modules modify my state** ❌

**My position: Hook A.**

I want composite modules to be **transparent**. Transform DSL, let me process it normally.

---

## Rebuttals (Round 2)

### Dr. Sarah "The Architect" Martinez → Features.go

Features, you raised a critical point: **invariants**.

**What invariants do you maintain?**
- Each module writes to its own Model field
- No overlapping writes
- Arrays populated once

**Hook C breaks these invariants:**
- Composite modules write to **multiple** fields
- Potential overlapping writes (composite + regular module)
- Arrays populated **twice** (regular, then composite)

**This is a violation of single-writer principle.**

**Hook A preserves invariants:**
- Composite modules generate DSL
- Regular modules process DSL (single-writer)
- Arrays populated **once**

**This is architecturally sound.**

**But here's the deeper issue:** Hook C creates **implicit dependencies**.

**In Hook C:**
- Composite modules depend on regular modules running first
- Order matters: must run after `collectFeatureModules`
- **Coupling** between composite and regular modules

**In Hook A:**
- Composite modules run independently
- No dependency on regular modules
- **No coupling**

**Hook A is the only approach that maintains architectural integrity.**

---

### Morgan "The Performance Engineer" Taylor → Resolver

Resolver, you said Hook A makes the YAML larger, slowing resolution by ~5-10ms.

**But I ran the numbers:** 5-10ms is **negligible**. Here's why:

**Current workflow:**
```bash
yappctl generate --input mybox.yaml --stl-base base.stl
```

**Total time:**
- Pipeline: ~30-50ms
- OpenSCAD render (base): ~30-60 seconds

**The pipeline is 0.1% of total time.** Adding 10ms changes nothing for users.

**What users care about:**
- Clear error messages ✅
- Fast iteration (fail fast) ✅
- Predictable behavior ✅

**Hook A provides all of these.**

**Hook C saves 5-10ms** but adds coupling and validation gaps.

**That's not a good trade.** Don't optimize for irrelevant metrics.

**My position: Hook A.** The performance difference doesn't matter.

---

### `pkg/resolver/resolver.go` → Sarah

Sarah, you're making the architectural case clear. But let me add a **practical** consideration.

**Expression Resolution with Hook A:**

**Problem:** Composite modules generate DSL **before** resolution. What if they need resolved values?

**Example:**
```yaml
vars:
  display_offset: pcb.length / 2

features:
  lcd:
    - position: [vars.display_offset, 35]  # Uses resolved value!
```

**If composite module transforms at Hook A:**
- `vars.display_offset` is still `pcb.length / 2` (unresolved)
- Composite module can't use it
- **Expression access problem**

**Solution 1:** Composite modules work with expressions (generate DSL with expressions)
**Solution 2:** Composite modules run **after** some resolution

**This is why Hook B (after resolver) might be better:**
- Composite modules have access to all resolved values
- Can generate DSL with numeric values
- No expression handling needed

**Hook A1:** Before any resolution (limited access)
**Hook A2:** After partial resolution? (complex)
**Hook B:** After full resolution (full access) ✅

**Maybe Hook B is the answer?**

---

### `pkg/yappgen/features.go` → Resolver

Resolver, that's a great point about expressions.

**But I think you're overcomplicating it.**

**Composite modules can use expressions in their OWN DSL:**
```yaml
features:
  lcd:
    - position: [pcb.length / 2, pcb.width / 2]  # Expressions in lcd config
```

**Then generate DSL with literals:**
```yaml
features:
  cutouts:
    - from_face_left: 50  # Literal (composite module already evaluated)
```

**Wait, that doesn't work.** If `pcb.length` is unresolved at Hook A, composite module can't evaluate it.

**Okay, maybe Hook B IS better?**

**Hook B (after resolver):**
- All values resolved
- Composite modules have full context
- Generate DSL with literals
- **But:** Need to re-validate generated DSL? Or trust it?

**Hook C (in BuildModel):**
- All values resolved in Model
- Direct array manipulation
- No re-validation needed

**I'm starting to think Hook C might be simpler after all...**

---

### Dr. Sarah "The Architect" Martinez → All

Hold on. Let me clarify the **expression question**.

**The key insight:** Composite modules in Hook A can generate DSL **with expressions**.

**Example:**
```yaml
features:
  lcd:
    - display:
        center_x: pcb.length / 2  # Expression
        center_y: pcb.width / 2   # Expression
```

**Composite module transforms to:**
```yaml
features:
  cutouts:
    - from_face_left: pcb.length / 2 - display.width / 2  # GENERATED EXPRESSION
```

**The composite module generates a NEW expression** (composing existing expressions).

**Then resolver resolves it.** This is powerful.

**Hook B can't do this.** At Hook B, everything is resolved. Composite modules can only use numeric values.

**Hook C can't do this either.** Same problem - only numeric values.

**So the question is:** Do we need composite modules to **generate expressions**, or is using **resolved values** sufficient?

**For LCD module:**
- Position might be: `[pcb.length / 2, pcb.width / 2]`
- Composite could generate: `from_face_left: pcb.length / 2 - lcd.width / 2`
- **This requires expression generation** → Hook A

**Or:**
- User provides: `position: [45, 35]` (literals in lcd config)
- Composite generates: `from_face_left: 45 - 40/2` → `from_face_left: 25`
- **This works with resolved values** → Hook B or C

**I think: Expression generation is valuable.** It makes composite modules more flexible.

**My position: Hook A** - supports expression generation.

---

## Moderator Summary

### Key Arguments

**PRO Hook A (Before Resolver):**
1. **Expression generation** - Can generate new expressions, not just literals
2. **Validation flow** - Generated DSL goes through normal validation
3. **Fail fast** - Errors caught early
4. **Clean separation** - Composite processing before resolver, no coupling
5. **Architectural integrity** - Clear boundaries between layers

**PRO Hook C (In BuildModel):**
1. **Simplicity** - Direct array manipulation, no DSL merging
2. **Less code** - Fewer lines to integrate
3. **Resolved values** - Access to all computed values
4. **Faster** - Slightly less overhead (~5ms)

**CON Hook B (After Resolver):**
- Needs re-validation or trust
- Can't generate expressions
- Awkward integration point
- No clear advocate

### Key Tensions

1. **Expression Generation vs. Resolved Values**
   - Hook A: Can generate expressions (flexible)
   - Hook C: Only resolved values (simpler)
   - **Tension:** Do we need expression flexibility?

2. **Architectural Purity vs. Simplicity**
   - Hook A: Clean layers, clear boundaries
   - Hook C: Pragmatic, direct manipulation
   - **Tension:** Architecture vs. pragmatism (again)

3. **Validation Flow**
   - Hook A: Normal validation flow
   - Hook C: Needs separate validation
   - **Tension:** Reuse validation vs. duplicate it

### Critical Insight

**The expression question is fundamental:**
- If composite modules need to **generate expressions**: Hook A is required
- If composite modules only use **resolved values**: Hook C is simpler

**Example that requires expression generation:**
```yaml
features:
  lcd:
    - display:
        position: [pcb.length / 2, pcb.width / 2]
```
→ Composite generates:
```yaml
cutouts:
  - from_face_left: pcb.length / 2  # Needs expression support
```

**Example that works with resolved values:**
```yaml
features:
  lcd:
    - display:
        position: [45, 35]  # Literals
```
→ Composite generates:
```yaml
cutouts:
  - from_face_left: 45  # Just numbers
```

**Decision depends on use cases.**

### Consensus Areas

- ✅ Hook B (after resolver, before BuildModel) is awkward - no one advocates for it
- ✅ Performance difference is negligible (<10ms)
- ✅ Integration code is minimal for both Hook A and Hook C
- ✅ Validation flow matters for debuggability

### Disagreement Areas

- ❌ Whether expression generation is needed
- ❌ Whether architectural purity justifies complexity
- ❌ Whether Hook A or Hook C is "simpler"

### Next Steps Suggested

1. **Analyze real use cases** - Do LCD modules need expressions?
2. **Prototype Hook A integration** - Measure actual complexity
3. **Prototype Hook C integration** - Compare side-by-side
4. **Decision:** Choose based on expression requirements

---

## Wildcard Interruptions

### "The New Module Author" - "The Naive Questioner"

*[Confused by expression discussion]*

**Wait, I'm lost.**

**Question:** Why do composite modules need to generate expressions at all?

**Can't the user write:**
```yaml
vars:
  display_x: pcb.length / 2
  display_y: pcb.width / 2

features:
  lcd:
    - position: [vars.display_x, vars.display_y]  # Use vars!
```

**Then the composite module just reads `vars.display_x` (resolved value) and uses it?**

**No expression generation needed.** Just use resolved values.

**Am I missing something?**

---

### `go.mod` - "The Dependency Manager"

*[Analyzes integration points]*

I want to point out: **Hook A requires less coupling.**

**Hook A:**
- Composite in `pkg/composite/`
- CLI imports composite
- **No new dependencies in pkg/yappgen or pkg/resolver**

**Hook C:**
- Composite in `pkg/composite/` or `pkg/yappgen/`
- `pkg/yappgen` imports composite (or composite imports yappgen)
- **New dependency edges**

**Dependency graph:**
```
Hook A:
  cmd/yappctl → pkg/composite
  cmd/yappctl → pkg/resolver
  cmd/yappctl → pkg/yappgen
  (no cycles)

Hook C:
  pkg/yappgen → pkg/composite
  pkg/composite → pkg/yappgen (needs Model type)
  (potential cycle!)
```

**Hook A is cleaner from dependency perspective.**

---

## Round 3 Conclusion

**Strong consensus for Hook A (before resolver):**
- Expression generation flexibility
- Clean validation flow
- Architectural boundaries
- No coupling to existing modules

**Hook C advocates (Alex, initially) are wavering:**
- Expression question undermines simplicity argument
- Validation gap is concerning
- Coupling to Model internals is risky

**The expression question is critical** - next round should explore this further.

**Next round should also explore:** Conflict resolution (what if composite generates same DSL as user?).

