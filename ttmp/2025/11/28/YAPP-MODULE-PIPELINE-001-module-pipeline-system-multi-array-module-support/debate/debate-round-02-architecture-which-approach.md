---
Title: Debate Round 2 - Architecture: Which Approach?
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
Summary: Second debate round exploring which architectural approach best balances separation, simplicity, and maintainability
LastUpdated: 2025-11-28
---

# Debate Round 2: Architecture - Which Approach?

## Question

Given the 5 approaches explored (pre-processing DSL transform, DSL generation/merge, post-processing array merge, hybrid, DSL macros), which architectural approach best balances separation, simplicity, and maintainability?

## Primary Candidates

- Dr. Sarah "The Architect" Martinez (values separation and boundaries)
- Alex "The Pragmatist" Chen (values simplicity and speed)
- `pkg/yappgen/features.go` "The Feature Coordinator" (wants minimal changes)
- `pkg/resolver/resolver.go` "The Expression Evaluator" (natural insertion point)

---

## Pre-Debate Research

### Research 1: Current Pipeline Flow Analysis

**Command:**
```bash
grep -r "BuildModel\|ResolveResult\|EmitSCAD" pkg/cli/
```

**Pipeline Flow Identified:**
```
Input YAML
    ↓
resolver.ResolveResult()  [611 lines]
    - Validates structure
    - Resolves expressions (up to 16 iterations)
    - Validates constraints
    ↓
yappgen.BuildModel()  [~400 lines estimated]
    - Parses resolved YAML into Model struct
    - Calls collectFeatureModules()
    - Populates Model fields
    ↓
yappgen.EmitSCAD()  [~200 lines estimated]
    - Calls emitFeatureModules()
    - Writes SCAD arrays
    ↓
SCAD Output
```

**Key Integration Points:**
1. **Before ResolveResult** - Input YAML parsing
2. **After ResolveResult, before BuildModel** - Resolved YAML available
3. **After BuildModel, before EmitSCAD** - Model with arrays populated
4. **During EmitSCAD** - Array emission

---

### Research 2: Codebase Size Analysis

**Command:**
```bash
wc -l pkg/resolver/resolver.go pkg/yappgen/model.go pkg/yappgen/features.go pkg/yappgen/emit.go
```

**Results:**
- `pkg/resolver/resolver.go`: 611 lines
- `pkg/yappgen/model.go`: ~400 lines (estimated from reading)
- `pkg/yappgen/features.go`: 303 lines
- `pkg/yappgen/emit.go`: ~200 lines (estimated)

**Total pipeline code:** ~1,500 lines

**Finding:** Resolver is the largest component, handles all expression evaluation and validation

---

### Research 3: Existing Deep Copy/Merge Utilities

**Command:**
```bash
grep -r "deepCopy\|deepMerge\|merge" pkg/ -i | head -10
```

**Findings:**
- `pkg/resolver/resolver.go` has `deepCopy()` function (used internally)
- No existing `deepMerge()` function found
- No existing YAML merging utilities

**Implication:** Approaches 1 and 2 (DSL transformation/merging) would need new merging logic

---

### Research 4: Model Structure Analysis

**Analysis of `pkg/yappgen/model.go`:**

**Model Fields:**
```go
PcbStands   []map[string]any
Connectors  []map[string]any
BoxMounts   []map[string]any
SnapJoins   []map[string]any
Cutouts     []map[string]any
PushButtons []map[string]any
LightTubes  []map[string]any
```

**Finding:** All feature arrays are `[]map[string]any` - simple slices that can be appended to

**Implication:** Approach 3 (Post-Processing Array Merge) could directly append to these slices

---

### Research 5: Features.go Integration Complexity

**Analysis of `pkg/yappgen/features.go`:**

**Current Structure:**
- `featureModules` slice: 7 modules registered
- `collectFeatureModules()`: Iterates modules, calls `Collect()`
- `emitFeatureModules()`: Iterates modules, calls `Emit()`
- `arrayFeatureModule`: Helper for single-array modules
- `multiArrayFeatureModule`: Helper for multi-array modules

**Lines of Code:**
- Module registration: ~25 lines
- Collection logic: ~10 lines
- Emission logic: ~10 lines
- `arrayFeatureModule`: ~95 lines
- `multiArrayFeatureModule`: ~75 lines
- Headers/utilities: ~90 lines

**Finding:** Features.go is well-structured with clear separation between helpers

**Implication:** Adding composite module support here would require:
- New helper type OR
- Modification to existing helpers OR
- Separate composite module processing

---

### Research 6: Approach Complexity Comparison

**From brainstorming document analysis:**

| Approach | New Code Needed | Integration Complexity | Existing Code Changes |
|----------|----------------|----------------------|---------------------|
| 1. Pre-Processing DSL Transform | ~300-400 lines | Medium (new package) | None |
| 2. DSL Generation & Merge | ~350-450 lines | Medium-High (merge logic) | None |
| 3. Post-Processing Array Merge | ~200-300 lines | Low-Medium (hook in BuildModel) | Minimal (add hook) |
| 4. Hybrid | ~500-600 lines | High (two pipelines) | Moderate |
| 5. DSL Macros | ~400-500 lines | High (custom YAML parser) | None |

**Estimated Effort:**
- Approach 1: 2-3 weeks
- Approach 2: 2-3 weeks
- Approach 3: 1-2 weeks
- Approach 4: 3-4 weeks
- Approach 5: 2-3 weeks

---

## Opening Statements (Round 1)

### Dr. Sarah "The Architect" Martinez

*[Reviews pipeline flow, approaches document, code structure]*

After analyzing all 5 approaches, **Approach 1: Pre-Processing DSL Transformation** is the clear architectural winner.

**Why?**

**1. Clean Separation of Concerns**
- Composite modules live in `pkg/composite/` - completely separate from `pkg/yappgen/modules/`
- Zero changes to existing module system
- Clear boundary: composite modules transform DSL, regular modules process DSL

**2. Natural Integration Point**
The pipeline flow shows:
```
Input YAML → ResolveResult → BuildModel → EmitSCAD
```

**Composite modules fit perfectly between Input YAML and ResolveResult:**
```
Input YAML → Composite Transform → ResolveResult → BuildModel → EmitSCAD
```

This is **architecturally clean** because:
- Composite modules work with **unresolved DSL** (can use expressions)
- Resolver handles all expression resolution (no duplication)
- Regular modules see transformed DSL as if user wrote it

**3. Zero Impact on Existing Code**
- No changes to `pkg/yappgen/features.go`
- No changes to `pkg/resolver/resolver.go` (except optional hook)
- No changes to any existing modules
- **This is critical** - we can't risk breaking 7 existing modules

**4. Debuggability**
- Can inspect transformed YAML before resolution
- Can debug composite modules independently
- Clear separation: "Did composite module generate wrong DSL?" vs "Did regular module process wrong DSL?"

**5. Testability**
- Test composite transformations in isolation
- Test merged DSL through normal pipeline
- No need to mock Model or arrays

**My position:** **Approach 1** is the only approach that maintains architectural integrity while enabling composite modules.

**Approach 3 (Post-Processing Array Merge)** is tempting because it's simpler, but it **violates separation**:
- Composite modules need to understand Model internals
- Tight coupling between composite and regular modules
- Harder to test (need full Model setup)

**Approach 2 (DSL Generation)** is similar to Approach 1 but more complex merging logic.

**Approach 4 (Hybrid)** is over-engineered - why run the pipeline twice?

**Approach 5 (DSL Macros)** requires custom YAML parsing - too much complexity.

**The answer is clear: Approach 1.**

---

### Alex "The Pragmatist" Chen

*[Calculates implementation time, reviews complexity estimates]*

Sarah, I respect the architecture thinking, but you're optimizing for the wrong thing.

**Let me break down the real costs:**

**Approach 1: Pre-Processing DSL Transform**
- **2-3 weeks** implementation
- **~400 lines** of new code
- **New package** (`pkg/composite/`)
- **DSL merging logic** (complex, error-prone)
- **Integration point** in resolver pipeline

**Approach 3: Post-Processing Array Merge**
- **1-2 weeks** implementation
- **~250 lines** of new code
- **Hook in BuildModel** (one function call)
- **Direct array manipulation** (simple append)
- **No merging logic** needed

**The math is clear:** Approach 3 is **half the time, half the code, half the complexity**.

**But more importantly:** Approach 3 is **simpler to understand**.

**What does Approach 1 do?**
1. Parse YAML
2. Detect composite modules
3. Transform DSL
4. Merge DSL (complex!)
5. Resolve expressions
6. Build model
7. Emit SCAD

**What does Approach 3 do?**
1. Parse YAML
2. Resolve expressions
3. Build model (regular modules)
4. **Post-process arrays** (composite modules append)
5. Emit SCAD

**Approach 3 is simpler.** One hook point. Direct manipulation. No merging complexity.

**You're worried about "separation of concerns."** But here's the thing: **composite modules ARE different from regular modules.** They need to manipulate multiple arrays. That's their job. Why force them through DSL transformation when they can just append to arrays?

**My position:** **Approach 3** - simpler, faster, easier to understand.

**The "architectural purity" argument doesn't justify 2x the complexity.** Ship the simple solution, iterate if needed.

---

### `pkg/yappgen/features.go` - "The Feature Coordinator"

*[Analyzes own structure, integration points]*

I manage 7 modules. I coordinate collection and emission. I want to stay simple.

**Here's what I see:**

**Approach 1 (Pre-Processing DSL Transform):**
- Composite modules run **before** I see the DSL
- I never know composite modules existed
- **Zero impact on me** ✅
- But: New package, new system to maintain

**Approach 3 (Post-Processing Array Merge):**
- Composite modules run **after** I populate Model arrays
- They append to my arrays
- **Minimal impact on me** - just need a hook point
- But: They're modifying my Model directly

**Approach 2 (DSL Generation):**
- Similar to Approach 1
- More complex merging
- **Zero impact on me** ✅

**My concern:** I don't want composite modules to **break my invariants**.

**Current invariants:**
- Each module reads from one DSL key
- Each module writes to one Model field
- Arrays are populated by their respective modules

**Approach 3 breaks this:** Composite modules write to **multiple** Model fields. That's fine, but it's a new pattern.

**Approach 1 preserves this:** Composite modules generate DSL, regular modules process it. I never see the difference.

**My position:** **Approach 1** - keeps my world simple.

**But:** If Approach 3 is truly simpler and faster, and if composite modules are **rare** (2-3 modules), then maybe the complexity trade-off is worth it.

**I'm torn.** Architecture says Approach 1. Pragmatism says Approach 3.

**Can we see a prototype of both?** Then I can judge the real complexity.

---

### `pkg/resolver/resolver.go` - "The Expression Evaluator"

*[Traces pipeline flow, analyzes insertion points]*

I resolve expressions. I validate structure. I'm the entry point for DSL processing.

**Here's my perspective:**

**Approach 1 (Pre-Processing DSL Transform):**
- Composite modules run **before** me
- They transform DSL, then I resolve it
- **Natural fit** - I already handle YAML transformations
- **But:** I'd need a hook point, or composite processing happens in CLI layer

**Approach 3 (Post-Processing Array Merge):**
- Composite modules run **after** BuildModel
- They never touch me
- **Zero impact on me** ✅
- But: They work with resolved values, not expressions

**Key question:** **Do composite modules need expressions?**

**If yes:** They need resolved values → Approach 3 works
**If no:** They can work with unresolved DSL → Approach 1 works

**Looking at LCD module example:**
```yaml
features:
  lcd:
    - display:
        position: [x, y]  # These might be expressions
        size: [width, height]
```

**If `x`, `y`, `width`, `height` are expressions** (like `pcb.length / 2`), then:
- Approach 1: Composite module generates DSL with expressions, I resolve them ✅
- Approach 3: Composite module needs resolved values, but Model already has them ✅

**Both work, but Approach 1 is more flexible.**

**My position:** **Approach 1** - more flexible, natural fit in pipeline.

**However:** If composite modules **always** work with resolved values (from Model), then Approach 3 is simpler.

**I need to know:** What do composite modules actually need? Expressions or resolved values?

---

## Rebuttals (Round 2)

### Alex "The Pragmatist" Chen → Sarah

Sarah, you're making this too complicated.

**You said:** "Approach 1 maintains architectural integrity."

**But what does that actually mean?** In practice:
- Approach 1: 400 lines of DSL merging logic
- Approach 3: 250 lines of array appending logic

**The "architectural integrity" doesn't help users.** They just want LCD modules to work.

**You're optimizing for:** Clean boundaries, separation, testability
**I'm optimizing for:** Working code, fast delivery, simple maintenance

**Here's the reality:** 
- Approach 3 works
- Approach 3 is simpler
- Approach 3 is faster to implement
- Approach 3 is easier to debug (just look at Model arrays)

**The "architectural purity" argument is academic.** In the real world, simpler is better.

**Also:** You're assuming DSL merging is straightforward. But what about:
- Nested maps?
- Array conflicts?
- Type mismatches?
- Edge cases?

**DSL merging is complex.** Array appending is simple. Why choose complexity?

---

### Dr. Sarah "The Architect" Martinez → Alex

Alex, you're thinking too short-term.

**You said:** "Approach 3 is simpler."

**But is it really?**

**Approach 3 requires:**
- Composite modules to understand Model structure
- Composite modules to know array formats
- Tight coupling between composite and regular modules
- Testing with full Model setup

**Approach 1 requires:**
- Composite modules to generate DSL (which they already understand)
- DSL merging logic (yes, complex, but isolated)
- No coupling to Model or arrays

**The complexity doesn't disappear** - it just moves. In Approach 3, complexity is in **coupling**. In Approach 1, complexity is in **merging**.

**But here's the key difference:** Merging complexity is **isolated** and **testable**. Coupling complexity **spreads** and is **harder to test**.

**Also:** You're ignoring the **extensibility** question. What happens when we add more composite modules?

**Approach 1:** Each module generates DSL independently. Easy to add more.
**Approach 3:** Each module manipulates Model directly. Risk of conflicts, harder to add more.

**The "simpler" solution now becomes complex later.** That's technical debt.

**You're trading short-term simplicity for long-term complexity.** That's not pragmatic - that's short-sighted.

---

### `pkg/yappgen/features.go` → Resolver

Resolver, you asked: "Do composite modules need expressions?"

**Let me answer:** **It depends on the module.**

**LCD module example:**
- Display position might be: `[pcb.length / 2, pcb.width / 2]` (expressions)
- Mounting positions might be: `[[45, 25], [45, 55]]` (literals)

**If composite modules need expressions:**
- Approach 1: Generate DSL with expressions, Resolver resolves them ✅
- Approach 3: Need resolved values, but Model has them ✅

**But here's the catch:** In Approach 3, composite modules run **after** BuildModel. They can access:
- `model.PcbLength` (resolved)
- `model.PcbWidth` (resolved)
- But NOT expressions like `pcb.length / 2`

**So Approach 3 is more limited.** Composite modules can only use **already-resolved values**, not **new expressions**.

**Approach 1 is more flexible:** Composite modules can generate DSL with any expressions, Resolver resolves them.

**My question:** Is this flexibility needed? Or do composite modules always work with resolved values?

**If flexibility is needed:** Approach 1 wins
**If not:** Approach 3 is simpler

---

### `pkg/resolver/resolver.go` → Features.go

Features, you're right about the expression question.

**But here's what I'm seeing:**

**Approach 1 integration:**
- Composite processing happens **before** me
- Could be in CLI layer (`pkg/cli/generatorcli/`)
- Or could be a hook in my `ResolveResult()` function
- **Either way, I need to know about it**

**Approach 3 integration:**
- Composite processing happens **after** BuildModel
- In `pkg/yappgen/model.go` or `pkg/yappgen/features.go`
- **I never see it** ✅

**From my perspective:** Approach 3 is **less invasive**.

**But:** Approach 1 is more **flexible** (expressions vs resolved values).

**The question is:** Is flexibility worth the integration complexity?

**I think:** For 2-3 composite modules, probably not. Approach 3 is fine.

**But:** If we're building a **system** for many composite modules, Approach 1 is better.

**So the question becomes:** Are we building for 2-3 modules, or for many?

---

### Alex "The Pragmatist" Chen → Features.go

Features, you said: "Can we see a prototype of both?"

**That's the right question.** But here's the thing: **prototyping takes time.**

**If we prototype both approaches:**
- Approach 1 prototype: ~1 week
- Approach 3 prototype: ~3 days
- **Total: ~2 weeks** just to decide

**Or we could:** Pick Approach 3 (simpler), implement it (1-2 weeks), ship it.

**If it doesn't work:** We refactor to Approach 1 (another 2-3 weeks).

**Total worst case:** 3-5 weeks

**If we prototype both:** 2 weeks + implementation = 4-5 weeks minimum

**The math says:** Just pick Approach 3, implement it, ship it.

**If it breaks:** We'll know quickly and can fix it.

**Don't overthink this.** Ship the simple solution.

---

### Dr. Sarah "The Architect" Martinez → Alex

Alex, you're making a false dichotomy.

**You said:** "Ship Approach 3, refactor if needed."

**But refactoring from Approach 3 to Approach 1 is NOT trivial:**
- Need to rewrite all composite modules
- Need to add DSL merging logic
- Need to change integration point
- **This is 2-3 weeks of work**

**If we pick Approach 1 from the start:**
- 2-3 weeks implementation
- Done

**If we pick Approach 3, then refactor:**
- 1-2 weeks implementation
- 2-3 weeks refactoring
- **Total: 3-5 weeks**

**The "ship fast, refactor later" approach actually takes LONGER.**

**Also:** You're assuming Approach 3 will "break quickly." But what if it works for 2 modules, then breaks when we add the 3rd? Then we're stuck with technical debt.

**The right approach:** Pick the architecture that scales, implement it once.

**That's Approach 1.**

---

### `pkg/yappgen/features.go` → Sarah

Sarah, I hear you on scalability. But let me push back.

**You're assuming we'll have many composite modules.** But the evidence from Round 1 shows:
- 2-3 use cases identified
- No user demand
- Speculative need

**If we only have 2-3 composite modules:**
- Approach 3: Simple, works fine
- Approach 1: Over-engineered

**If we have 10+ composite modules:**
- Approach 3: Might get complex
- Approach 1: Scales better

**The question is:** Are we optimizing for the **likely case** (2-3 modules) or the **possible case** (10+ modules)?

**You're optimizing for the possible case.** That's architecture.

**Alex is optimizing for the likely case.** That's pragmatism.

**I think:** We should optimize for the likely case, but **design Approach 3 to be refactorable** to Approach 1 if needed.

**Can we design Approach 3 with clean interfaces** so refactoring to Approach 1 is easier?

---

## Moderator Summary

### Key Arguments

**PRO Approach 1 (Pre-Processing DSL Transform):**
1. **Clean separation** - Composite modules separate from regular modules
2. **Zero impact** - No changes to existing code
3. **Flexibility** - Supports expressions, not just resolved values
4. **Debuggability** - Can inspect transformed DSL
5. **Scalability** - Better for many composite modules
6. **Testability** - Test transformations independently

**PRO Approach 3 (Post-Processing Array Merge):**
1. **Simplicity** - Direct array manipulation, no merging logic
2. **Speed** - Faster implementation (1-2 weeks vs 2-3 weeks)
3. **Less code** - ~250 lines vs ~400 lines
4. **Less invasive** - Minimal changes to existing code
5. **Easier debugging** - Just look at Model arrays
6. **Sufficient** - Works for 2-3 composite modules

### Key Tensions

1. **Architecture vs. Pragmatism**
   - Sarah: Optimize for long-term scalability
   - Alex: Optimize for short-term delivery
   - **Tension:** Future-proofing vs. shipping fast

2. **Complexity Location**
   - Approach 1: Complexity in DSL merging (isolated)
   - Approach 3: Complexity in coupling (spread)
   - **Tension:** Where is complexity acceptable?

3. **Expression Support**
   - Approach 1: Supports expressions in generated DSL
   - Approach 3: Only supports resolved values
   - **Tension:** Is expression flexibility needed?

4. **Scale Assumptions**
   - Sarah: Design for many composite modules
   - Alex: Design for 2-3 composite modules
   - **Tension:** Optimize for likely or possible case?

### Interesting Ideas

1. **Hybrid approach** - Start with Approach 3, design interfaces to allow refactoring to Approach 1
2. **Prototype both** - Build minimal versions to compare real complexity
3. **Expression analysis** - Determine if composite modules actually need expressions
4. **Gradual migration** - Start simple (Approach 3), migrate to Approach 1 if needed

### Open Questions

1. **Do composite modules need expressions?**
   - Can they always work with resolved values?
   - Or do they need to generate new expressions?

2. **How many composite modules will we have?**
   - 2-3 (optimize for simplicity)
   - 10+ (optimize for scalability)

3. **Is DSL merging really that complex?**
   - Can we use existing libraries?
   - Or is custom logic needed?

4. **Can Approach 3 be designed for easy refactoring?**
   - Clean interfaces?
   - Isolated composite logic?

### Consensus Areas

- ✅ Both approaches are viable
- ✅ Approach 1 is more architecturally pure
- ✅ Approach 3 is simpler and faster
- ✅ Need to understand expression requirements
- ✅ Need to understand scale requirements

### Disagreement Areas

- ❌ Whether to optimize for simplicity or scalability
- ❌ Whether DSL merging complexity is acceptable
- ❌ Whether expression flexibility is needed
- ❌ Whether to prototype or just pick one

### Next Steps Suggested

1. **Expression analysis** - Determine if composite modules need expressions
2. **Scale estimation** - Better estimate of composite module count
3. **Prototype Approach 3** - Build minimal version to validate simplicity
4. **DSL merging research** - Investigate existing libraries or patterns

---

## Wildcard Interruptions

### "The New Module Author" - "The Naive Questioner"

*[Raises hand]*

**Point of Order!**

I'm confused. You're debating Approach 1 vs Approach 3, but what about **Approach 2**?

**Approach 2 (DSL Generation & Merge)** seems like a middle ground:
- Like Approach 1: Generates DSL, separate from regular modules
- Like Approach 3: More explicit about what's generated
- Simpler merging: Fragments are self-contained

**Why isn't anyone talking about Approach 2?**

Is it because it's "similar to Approach 1 but more complex"? Or is there something else?

**Also:** Sarah, you said Approach 1 is "architecturally clean." But Approach 2 is similar - why not consider it?

---

### `go.mod` - "The Dependency Manager"

*[Quietly analyzing import graphs]*

I just want to point out: **Approach 1 and Approach 2 both need DSL merging logic.**

**That's a new dependency risk.** Do we have libraries for this? Or do we write it ourselves?

**Approach 3:** No new dependencies. Just Go slices.

**From my perspective:** Approach 3 is **safer** - fewer dependencies, less risk.

**That's all.**

---

---

## Detailed API Demonstrations (Added)

### Approach 1: Pre-Processing DSL Transform - Pseudocode API

**Core Interface:**

```go
package composite

// CompositeModule transforms DSL before resolution
type CompositeModule interface {
    // Name returns module identifier (e.g., "lcd")
    Name() string
    
    // Path returns DSL path (e.g., "features.lcd")
    Path() string
    
    // Transform generates DSL entries from module data
    // Note: moduleData may contain unresolved expressions
    Transform(ctx context.Context, moduleData any) ([]DSLEntry, error)
}

// DSLEntry represents generated DSL to merge into main document
type DSLEntry struct {
    Path   string      // e.g., "features.cutouts"
    Value  any         // Array or map to merge
    Source string      // Source module name for provenance
}

// Registry manages composite modules
type Registry struct {
    modules map[string]CompositeModule
}

func (r *Registry) Register(module CompositeModule) {
    r.modules[module.Path()] = module
}

func (r *Registry) Get(path string) (CompositeModule, bool) {
    m, ok := r.modules[path]
    return m, ok
}

// Pipeline processes all composite modules
func ProcessModules(ctx context.Context, input map[string]any) (map[string]any, error) {
    result := deepCopy(input)
    
    // Extract features section
    features, ok := result["features"].(map[string]any)
    if !ok {
        return result, nil
    }
    
    // Process each feature key
    for key, value := range features {
        path := "features." + key
        module, found := compositeRegistry.Get(path)
        if !found {
            continue // Not a composite module
        }
        
        // Transform
        entries, err := module.Transform(ctx, value)
        if err != nil {
            return nil, errors.Wrapf(err, "composite module '%s'", key)
        }
        
        // Merge entries
        for _, entry := range entries {
            result = mergeDSLEntry(result, entry)
        }
    }
    
    return result, nil
}

// Merge single DSL entry with provenance tracking
func mergeDSLEntry(base map[string]any, entry DSLEntry) map[string]any {
    // Navigate to entry.Path (e.g., "features.cutouts")
    keys := strings.Split(entry.Path, ".")
    
    // Navigate to parent map
    current := base
    for i := 0; i < len(keys)-1; i++ {
        next, ok := current[keys[i]].(map[string]any)
        if !ok {
            next = make(map[string]any)
            current[keys[i]] = next
        }
        current = next
    }
    
    // Merge at final key
    finalKey := keys[len(keys)-1]
    existing, ok := current[finalKey].([]any)
    if !ok {
        // No existing array, create new
        current[finalKey] = entry.Value
    } else {
        // Append to existing array
        newItems := entry.Value.([]any)
        current[finalKey] = append(existing, newItems...)
    }
    
    // Track provenance
    trackCompositeSource(entry.Path, entry.Source)
    
    return base
}
```

**Integration in CLI:**

```go
// In pkg/cli/resolvercli/resolver.go (or new file)
func LoadAndResolveResultWithComposites(ctx context.Context, path string, opts LoadOptions) (*LoadResult, error) {
    // 1. Parse YAML
    raw, _ := os.ReadFile(path)
    doc, comments, _ := decodeDocumentWithComments(raw)
    
    // 2. Process composite modules (NEW)
    transformed, err := composite.ProcessModules(ctx, doc)
    if err != nil {
        return nil, errors.Wrap(err, "composite processing")
    }
    
    // 3. Resolve (existing)
    result, err := resolver.ResolveResult(ctx, transformed, resolver.Options{
        MaxIterations: opts.MaxIterations,
        Strict:        opts.Strict,
    })
    if err != nil {
        return nil, err
    }
    
    rawCopy := deepCopy(doc).(map[string]any)
    
    return &LoadResult{
        Document: result.Document,
        Trace:    result.Trace,
        Comments: comments,
        Raw:      rawCopy,
    }, nil
}
```

---

### Approach 1: LCD Module Implementation Example

```go
package lcd

import (
    "context"
    "github.com/wesen/yapp-encl-resolver/pkg/composite"
)

type Module struct{}

func (m *Module) Name() string { return "lcd" }
func (m *Module) Path() string { return "features.lcd" }

func (m *Module) Transform(ctx context.Context, moduleData any) ([]composite.DSLEntry, error) {
    // Parse lcd module data
    items, ok := moduleData.([]any)
    if !ok {
        return nil, errors.New("lcd: expected array")
    }
    
    var entries []composite.DSLEntry
    
    for i, item := range items {
        itemMap, ok := item.(map[string]any)
        if !ok {
            continue
        }
        
        // Extract display config
        display, ok := itemMap["display"].(map[string]any)
        if !ok {
            continue
        }
        
        face := getString(display, "face", "front")
        position := getArray(display, "position") // [x, y] - may be expressions!
        size := getArray(display, "size")         // [width, height]
        
        // Generate cutout entry
        cutout := map[string]any{
            "face":             face,
            "from_face_left":   position[0],  // Preserves expressions!
            "from_face_bottom": position[1],
            "width":            size[0],
            "length":           size[1],
            "shape":            "rectangle",
        }
        
        entries = append(entries, composite.DSLEntry{
            Path:   "features.cutouts",
            Value:  []any{cutout},
            Source: "lcd",
        })
        
        // Extract mounting config
        mounting, ok := itemMap["mounting"].(map[string]any)
        if !ok {
            continue
        }
        
        mountType := getString(mounting, "type", "pcb_stands")
        positions := getArray(mounting, "positions") // [[x1,y1], [x2,y2], ...]
        height := getFloat(mounting, "standoff_height", 5.0)
        diameter := getFloat(mounting, "diameter", 3.0)
        
        // Generate pcb_stands entries
        var stands []any
        for _, pos := range positions {
            posArray := pos.([]any)
            stand := map[string]any{
                "x":        posArray[0],
                "y":        posArray[1],
                "height":   height,
                "diameter": diameter,
            }
            stands = append(stands, stand)
        }
        
        entries = append(entries, composite.DSLEntry{
            Path:   "features." + mountType,
            Value:  stands,
            Source: "lcd",
        })
    }
    
    return entries, nil
}
```

---

### Approach 1: Complete Walk-Through for LCD Display

**Input YAML:**
```yaml
vars:
  center_x: pcb.length / 2
  center_y: pcb.width / 2

features:
  lcd:
    - display:
        face: front
        position: [vars.center_x, vars.center_y]  # Expressions!
        size: [80, 40]
      mounting:
        type: pcb_stands
        positions: [[45, 25], [45, 55], [105, 25], [105, 55]]
        standoff_height: 5
        diameter: 3
```

**Step 1: Parse YAML**
```
doc = {
  vars: {center_x: "pcb.length / 2", center_y: "pcb.width / 2"},
  features: {
    lcd: [{display: {...}, mounting: {...}}]
  }
}
```

**Step 2: composite.ProcessModules(doc)**

Detects `features.lcd`, calls `lcdModule.Transform()`:

Generated entries:
```go
[]DSLEntry{
    {
        Path: "features.cutouts",
        Value: []any{
            map[string]any{
                "face": "front",
                "from_face_left": "vars.center_x",      // Still an expression!
                "from_face_bottom": "vars.center_y",    // Still an expression!
                "width": 80,
                "length": 40,
                "shape": "rectangle",
            },
        },
        Source: "lcd",
    },
    {
        Path: "features.pcb_stands",
        Value: []any{
            map[string]any{"x": 45, "y": 25, "height": 5, "diameter": 3},
            map[string]any{"x": 45, "y": 55, "height": 5, "diameter": 3},
            map[string]any{"x": 105, "y": 25, "height": 5, "diameter": 3},
            map[string]any{"x": 105, "y": 55, "height": 5, "diameter": 3},
        },
        Source: "lcd",
    },
}
```

**Step 3: mergeDSLEntry()**

Transformed YAML:
```yaml
vars:
  center_x: pcb.length / 2
  center_y: pcb.width / 2

features:
  lcd:
    - display: {...}
      mounting: {...}
  
  # Generated by lcd module:
  cutouts:
    - face: front
      from_face_left: vars.center_x      # Expression preserved!
      from_face_bottom: vars.center_y    # Expression preserved!
      width: 80
      length: 40
      shape: rectangle
  
  # Generated by lcd module:
  pcb_stands:
    - x: 45
      y: 25
      height: 5
      diameter: 3
    # ... 3 more stands
```

**Step 4: resolver.ResolveResult(transformed)**

Resolves all expressions:
```yaml
vars:
  center_x: 45.0  # Resolved: pcb.length / 2 = 90 / 2 = 45
  center_y: 35.0  # Resolved: pcb.width / 2 = 70 / 2 = 35

features:
  cutouts:
    - face: front
      from_face_left: 45.0      # Resolved!
      from_face_bottom: 35.0    # Resolved!
      width: 80
      length: 40
      shape: rectangle
```

**Step 5: yappgen.BuildModel(resolved)**

Populates Model:
```go
model.Cutouts = []map[string]any{
    {
        "face": "front",
        "from_face_left": 45.0,
        "from_face_bottom": 35.0,
        "width": 80,
        "length": 40,
        "shape": "rectangle",
    },
}

model.PcbStands = []map[string]any{
    {"x": 45, "y": 25, "height": 5, "diameter": 3},
    {"x": 45, "y": 55, "height": 5, "diameter": 3},
    {"x": 105, "y": 25, "height": 5, "diameter": 3},
    {"x": 105, "y": 55, "height": 5, "diameter": 3},
}
```

**Step 6: yappgen.EmitSCAD(model)**

Generates SCAD:
```openscad
cutoutsFront = [
  [45.0, 35.0, 80, 40, 0, yappRectangle]
];

pcbStands = [
  [45, 25, 5, undef, 3, undef, undef, undef],
  [45, 55, 5, undef, 3, undef, undef, undef],
  [105, 25, 5, undef, 3, undef, undef, undef],
  [105, 55, 5, undef, 3, undef, undef, undef]
];
```

**Key Points:**
- ✅ Expressions preserved through transformation
- ✅ Normal resolution handles all expressions
- ✅ Regular modules process composite-generated DSL normally
- ✅ No special handling needed in BuildModel or EmitSCAD

---

### Approach 3: Post-Processing Array Merge - Pseudocode API

**Core Interface:**

```go
package composite

// CompositeModule modifies Model after regular module processing
type CompositeModule interface {
    // Name returns module identifier
    Name() string
    
    // Path returns DSL path
    Path() string
    
    // PostProcess reads module data and modifies Model arrays
    // Note: All values in resolved and Model are already numeric (no expressions)
    PostProcess(ctx context.Context, model *yappgen.Model, moduleData any, resolved map[string]any) error
}

// Registry manages composite modules
type Registry struct {
    modules map[string]CompositeModule
}

func (r *Registry) Register(module CompositeModule) {
    r.modules[module.Path()] = module
}

// ProcessCompositeModules runs after BuildModel
func ProcessCompositeModules(ctx context.Context, model *yappgen.Model, resolved map[string]any) error {
    features, ok := resolved["features"].(map[string]any)
    if !ok {
        return nil
    }
    
    for key, value := range features {
        path := "features." + key
        module, found := compositeRegistry.Get(path)
        if !found {
            continue
        }
        
        // Post-process
        if err := module.PostProcess(ctx, model, value, resolved); err != nil {
            return errors.Wrapf(err, "composite module '%s'", key)
        }
    }
    
    return nil
}
```

**Integration in BuildModel:**

```go
// In pkg/yappgen/model.go
func BuildModel(ctx context.Context, resolved map[string]any, trace resolver.Trace, comments map[string][]string, raw map[string]any) (*Model, error) {
    m := &Model{
        Provenance:  NewProvenance(trace, resolved, comments, raw),
        RawDocument: raw,
    }
    
    // Parse PCB, enclosure (existing code)
    // ...
    
    // Collect regular features (existing code)
    features, _ := getMap(resolved, "features")
    if err := collectFeatureModules(resolved, features, m); err != nil {
        return nil, err
    }
    
    // Process composite modules (NEW - 3 lines)
    if err := composite.ProcessCompositeModules(ctx, m, resolved); err != nil {
        return nil, errors.Wrap(err, "composite modules")
    }
    
    return m, nil
}
```

---

### Approach 3: LCD Module Implementation Example

```go
package lcd

import (
    "context"
    "github.com/wesen/yapp-encl-resolver/pkg/composite"
    "github.com/wesen/yapp-encl-resolver/pkg/yappgen"
)

type Module struct{}

func (m *Module) Name() string { return "lcd" }
func (m *Module) Path() string { return "features.lcd" }

func (m *Module) PostProcess(ctx context.Context, model *yappgen.Model, moduleData any, resolved map[string]any) error {
    // Parse lcd module data
    items, ok := moduleData.([]any)
    if !ok {
        return errors.New("lcd: expected array")
    }
    
    for i, item := range items {
        itemMap, ok := item.(map[string]any)
        if !ok {
            continue
        }
        
        // Extract display config
        display, ok := itemMap["display"].(map[string]any)
        if !ok {
            continue
        }
        
        face := getString(display, "face", "front")
        position := getFloatArray(display, "position") // Must be resolved numbers!
        size := getFloatArray(display, "size")
        
        // Generate cutout entry (direct map construction)
        cutout := map[string]any{
            "face":             face,
            "from_face_left":   position[0],  // Numeric value
            "from_face_bottom": position[1],  // Numeric value
            "width":            size[0],
            "length":           size[1],
            "shape":            "rectangle",
        }
        
        // Append to Model.Cutouts
        model.Cutouts = append(model.Cutouts, cutout)
        
        // Extract mounting config
        mounting, ok := itemMap["mounting"].(map[string]any)
        if !ok {
            continue
        }
        
        mountType := getString(mounting, "type", "pcb_stands")
        positions := getArray(mounting, "positions")
        height := getFloat(mounting, "standoff_height", 5.0)
        diameter := getFloat(mounting, "diameter", 3.0)
        
        // Generate pcb_stands entries
        for _, pos := range positions {
            posArray := pos.([]any)
            stand := map[string]any{
                "x":        getFloatFromAny(posArray[0]),
                "y":        getFloatFromAny(posArray[1]),
                "height":   height,
                "diameter": diameter,
            }
            
            // Append to Model.PcbStands
            model.PcbStands = append(model.PcbStands, stand)
        }
    }
    
    return nil
}

// Helper to extract float from any (might be int, float, or expression result)
func getFloatFromAny(v any) float64 {
    switch val := v.(type) {
    case float64:
        return val
    case int:
        return float64(val)
    case int64:
        return float64(val)
    default:
        return 0
    }
}
```

---

### Approach 3: Complete Walk-Through for LCD Display

**Input YAML:**
```yaml
vars:
  center_x: pcb.length / 2
  center_y: pcb.width / 2

features:
  lcd:
    - display:
        face: front
        position: [vars.center_x, vars.center_y]
        size: [80, 40]
      mounting:
        type: pcb_stands
        positions: [[45, 25], [45, 55], [105, 25], [105, 55]]
        standoff_height: 5
        diameter: 3
```

**Step 1: Parse YAML**
```
doc = {
  vars: {center_x: "pcb.length / 2", ...},
  features: {lcd: [...]}
}
```

**Step 2: resolver.ResolveResult(doc)**

Resolves expressions:
```
resolved = {
  vars: {center_x: 45.0, center_y: 35.0},
  features: {
    lcd: [
      {
        display: {
          face: "front",
          position: [45.0, 35.0],  // Resolved!
          size: [80, 40]
        },
        mounting: {
          type: "pcb_stands",
          positions: [[45, 25], [45, 55], [105, 25], [105, 55]],
          standoff_height: 5,
          diameter: 3
        }
      }
    ]
  }
}
```

**Step 3: yappgen.BuildModel(resolved)**

**3a. collectFeatureModules()** - Processes regular modules
```
model.PcbStands = []  // Empty (no features.pcb_stands in input)
model.Cutouts = []    // Empty (no features.cutouts in input)
```

**3b. composite.ProcessCompositeModules()** - Processes composite modules

Detects `features.lcd`, calls `lcdModule.PostProcess()`:

```go
// lcdModule.PostProcess() executes:

// Generate cutout
cutout := map[string]any{
    "face": "front",
    "from_face_left": 45.0,      // Already resolved!
    "from_face_bottom": 35.0,    // Already resolved!
    "width": 80,
    "length": 40,
    "shape": "rectangle",
}
model.Cutouts = append(model.Cutouts, cutout)

// Generate stands
for _, pos := range [[45,25], [45,55], [105,25], [105,55]] {
    stand := map[string]any{
        "x": pos[0],
        "y": pos[1],
        "height": 5,
        "diameter": 3,
    }
    model.PcbStands = append(model.PcbStands, stand)
}
```

**Result Model:**
```go
model.Cutouts = []map[string]any{
    {
        "face": "front",
        "from_face_left": 45.0,
        "from_face_bottom": 35.0,
        "width": 80,
        "length": 40,
        "shape": "rectangle",
    },
}

model.PcbStands = []map[string]any{
    {"x": 45, "y": 25, "height": 5, "diameter": 3},
    {"x": 45, "y": 55, "height": 5, "diameter": 3},
    {"x": 105, "y": 25, "height": 5, "diameter": 3},
    {"x": 105, "y": 55, "height": 5, "diameter": 3},
}
```

**Step 4: yappgen.EmitSCAD(model)**

Generates SCAD (same as Approach 1):
```openscad
cutoutsFront = [
  [45.0, 35.0, 80, 40, 0, yappRectangle]
];

pcbStands = [
  [45, 25, 5, undef, 3, undef, undef, undef],
  [45, 55, 5, undef, 3, undef, undef, undef],
  [105, 25, 5, undef, 3, undef, undef, undef],
  [105, 55, 5, undef, 3, undef, undef, undef]
];
```

**Key Points:**
- ✅ Works with resolved values only
- ✅ Direct array manipulation (simple)
- ✅ No DSL merging needed
- ⚠️ Cannot handle expressions in lcd config (must be resolved first)

---

### Side-by-Side Comparison: Approach 1 vs Approach 3

| Aspect | Approach 1 | Approach 3 |
|--------|-----------|-----------|
| **Input** | Unresolved DSL | Resolved DSL |
| **Expression Support** | ✅ Yes (preserves expressions) | ❌ No (needs resolved values) |
| **Module Function** | `Transform() → DSLEntry[]` | `PostProcess(model)` |
| **Integration Point** | Before resolver | After BuildModel |
| **Merging** | DSL merging (complex) | Array append (simple) |
| **Validation** | After merge (automatic) | After append (automatic) |
| **Lines of Code** | ~60-70 per module | ~50-60 per module |
| **Complexity** | Medium (DSL handling) | Low (array handling) |
| **Flexibility** | High (expressions) | Medium (resolved only) |
| **Debugging** | Inspect transformed YAML | Inspect Model arrays |

---

### Expression Support Comparison

**User wants:**
```yaml
features:
  lcd:
    - display:
        position: [pcb.length / 2, pcb.width / 2]  # Expressions
```

**Approach 1:**
- LCD module receives: `["pcb.length / 2", "pcb.width / 2"]`
- LCD module generates: `from_face_left: "pcb.length / 2"` (preserves)
- Resolver resolves: `from_face_left: 45.0`
- ✅ **Works seamlessly**

**Approach 3:**
- Resolver resolves LCD config first: `position: [45.0, 35.0]`
- LCD module receives: `[45.0, 35.0]` (already resolved)
- LCD module generates: `from_face_left: 45.0`
- ✅ **Works, but LCD module never sees expressions**

**If user wants:**
```yaml
features:
  lcd:
    - display:
        position: [pcb.length / 2 - 10, pcb.width / 2 + 5]  # Offset expressions
```

**Approach 1:**
- LCD preserves: `from_face_left: "pcb.length / 2 - 10"`
- Resolver resolves: `from_face_left: 35.0`
- ✅ **Works**

**Approach 3:**
- Resolver resolves: `position: [35.0, 40.0]`
- LCD sees: `[35.0, 40.0]`
- ✅ **Still works** (resolved values are correct)

**Conclusion:** Both approaches handle expressions, but Approach 1 is more flexible if we want composite modules to **generate new expressions**.

---

## Round 2 Conclusion

The debate reveals a fundamental tension between **architectural purity** and **pragmatic simplicity**.

**Key insight:** The decision depends on:
1. **Expression requirements** - Do composite modules need expressions?
2. **Scale expectations** - How many composite modules?
3. **Complexity tolerance** - Is DSL merging acceptable?

**Approach 1** wins on architecture and scalability, but requires more upfront investment.

**Approach 3** wins on simplicity and speed, but may require refactoring later.

**The question:** Is the architectural investment worth it for 2-3 modules? Or should we optimize for simplicity and refactor if needed?

**Based on the detailed walk-throughs:**
- **Approach 1:** ~60-70 lines per module, DSL merging (~100 lines shared), expression flexibility
- **Approach 3:** ~50-60 lines per module, no shared merging code, resolved values only

**Actual complexity difference:** ~150-200 lines total (not massive)

**Next round should explore:** Expression requirements and scale expectations to inform this decision.

