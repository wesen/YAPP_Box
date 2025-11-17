---
Title: Intern Handoff: Getting Started with Path 1 (A→B) Refactor
Ticket: YAPP-BUILDER-REFACTOR-001
Status: active
Topics:
    - yapp
    - refactor
    - onboarding
DocType: playbook
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Step-by-step guide for intern to understand context and begin implementation
LastUpdated: 2025-11-17
---

# Intern Handoff: Getting Started with Path 1 (A→B) Refactor

## Welcome!

You're taking over **YAPP-BUILDER-REFACTOR-001**, a refactor to improve the YAPP DSL module builder system. This document will get you up to speed quickly.

## What This Ticket Is About

**Problem:** Module builders currently use a clunky pattern:
- Every builder does `yaml.Marshal()` then `yaml.Unmarshal()` to convert `map[string]any` to typed structs
- 12 error wrapping sites across 7 modules doing the same thing
- Cutouts module has special-case handling because it outputs 6 arrays instead of 1

**Solution:** Path 1 (A→B)
- **Step A:** Generate `Decode()` functions to eliminate marshal/unmarshal boilerplate
- **Step B:** Add `ArrayDecl` IR to unify single-array and multi-array modules

**Outcome:** Cleaner code, consistent patterns, easier to add new modules.

## Reading Order (Essential Documents)

### 1. Architecture Guide (START HERE - 30 min read)

**File:** [design/01-architecture-implementation-guide-path-1-a-b-builder-contract-refactor.md](../../../../2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/design/01-architecture-implementation-guide-path-1-a-b-builder-contract-refactor.md)

**What it covers:**
- Before/after architecture diagrams
- Complete pseudocode for all changes
- Implementation checklist (phases 1-4)
- Testing strategy
- Success criteria

**Key sections:**
- "Step A: Generate Decode()" - Shows helper functions and template changes
- "Step B: Add ArrayDecl IR" - Shows interface changes and wrapper patterns
- "Implementation Checklist" - Your roadmap

### 2. Tasks List (5 min read)

**File:** [tasks.md](../tasks.md)

**What it covers:**
- 21 high-level tasks organized by phase
- Estimated time: 3-4 days total
- Success criteria at the end

### 3. Options Analysis (Optional - 20 min read)

**File:** [analysis/04-builder-contract-and-codegen-options.md](../../../../2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/analysis/04-2025-11-17-builder-contract-and-codegen-options.md)

**What it covers:**
- Original analysis of Options A-E
- Why we chose Path 1 (A→B) over Path 2 (A→C)
- Current state pain points with code references

**Read this if:** You want to understand why we're doing this refactor.

### 4. Debate Rounds (Optional - 1-2 hours)

**Files:** `debate/` directory in parent ticket

**What they cover:**
- Round 1: Why refactor now (not defer)
- Round 2: Why Path 1 (A→B) not Path 2 (A→C)
- Round 5: What to generate (Decode() only, not wrapper)
- Round 6: How to enforce (linter + tests + docs)

**Read these if:** You want to understand the decision rationale and trade-offs.

## How to Start

### Step 1: Set Up Your Environment

```bash
cd /home/manuel/code/others/YAPP_Box

# Run tests to establish baseline
go test ./...

# Generate SCAD from examples to establish baseline
go run ./cmd/yappctl generate -i examples/04-features.yaml -o /tmp/baseline-04-features.scad
go run ./cmd/yappctl generate -i examples/yapp-mvp-pcbstands-cutouts.yaml -o /tmp/baseline-mvp.scad

# Save baselines for comparison later
mkdir -p /tmp/refactor-baseline
cp /tmp/baseline-*.scad /tmp/refactor-baseline/
```

### Step 2: Read the Architecture Guide

Open [design/01-architecture-implementation-guide-path-1-a-b-builder-contract-refactor.md](../../../../2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/design/01-architecture-implementation-guide-path-1-a-b-builder-contract-refactor.md) and read:
- "Architecture Overview" section
- "Step A: Generate Decode()" section
- "Step B: Add ArrayDecl IR" section
- "Implementation Checklist"

**Time:** 30 minutes

### Step 3: Start with Phase 1 (Step A)

Follow the tasks in [tasks.md](../tasks.md) starting with:
1. Create shared helpers package
2. Update schemagen templates
3. Regenerate modules
4. Update module Build() functions

**After each task:**
```bash
# Run tests
go test ./pkg/yappgen/modules/...

# Check if code compiles
go build ./...
```

### Step 4: Validate Step A Before Moving to Step B

**Critical:** Don't start Step B until Step A is complete and validated.

```bash
# Run all tests
go test ./...

# Generate SCAD and compare with baseline
go run ./cmd/yappctl generate -i examples/04-features.yaml -o /tmp/after-step-a.scad
diff /tmp/refactor-baseline/baseline-04-features.scad /tmp/after-step-a.scad
# Should be identical (or only comment differences)
```

### Step 5: Continue with Phase 2 (Step B)

Follow remaining tasks for ArrayDecl implementation.

## Key Concepts to Understand

### What is Decode()?

**Before (current):**
```go
// Every module does this
data, err := yaml.Marshal(it)
var item PcbStandsItem
yaml.Unmarshal(data, &item)
```

**After (generated):**
```go
// Generated function in schema_gen.go
func Decode(items []map[string]any) ([]PcbStandsItem, error) {
    // Type-safe extraction with rich error messages
    for idx, raw := range items {
        x, err := decode.GetFloat(raw, "x", label)
        // ... extract all fields
        item := PcbStandsItem{X: x, Y: y, ...}
        typed = append(typed, item)
    }
    return typed, nil
}
```

### What is ArrayDecl?

**Problem:** Some modules return 1 array (pcbstands), others return 6 arrays (cutouts).

**Solution:** Unified IR:
```go
type ArrayDecl struct {
    Name string   // "pcbStands" or "cutoutsFront"
    Rows [][]any  // The actual array data
}

// Single-array module
return []ArrayDecl{{Name: "pcbStands", Rows: rows}}

// Multi-array module
return []ArrayDecl{
    {Name: "cutoutsFront", Rows: frontRows},
    {Name: "cutoutsBack", Rows: backRows},
    // ... 4 more
}
```

### What is multiArrayFeatureModule?

**Problem:** Cutouts has 46 lines of custom code in features.go for multi-array handling.

**Solution:** Reusable helper that both cutouts and future ridgeExt can use:
```go
newMultiArrayFeatureModule("cutouts",
    func(m *Model) *[]map[string]any { return &m.Cutouts },
    cutouts.NewModule().Build)  // Returns []ArrayDecl
```

## Common Questions

**Q: Why not generate the ArrayDecl wrapper too?**

A: It's only 10-20 lines of simple boilerplate. Generating it adds template complexity for minimal benefit. See [Debate Round 7](../../../../2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/08-debate-round-7-future-extensibility-accommodating-new-output-shapes.md#jordan-rivera--the-architect-1).

**Q: Why not change Model to store typed slices?**

A: That's Path 2 (A→C), which was rejected. Higher cost (~120 lines vs ~70 lines), and the benefit (errors caught 10 lines earlier) doesn't justify it. See [Debate Round 8](../../../../2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/09-debate-round-8-type-safety-end-to-end-arraydecl-vs-typed-model.md).

**Q: What if I break something?**

A: Run tests frequently! After each file change, run `go test ./...`. If tests pass, you're probably okay. If SCAD output changes, you've introduced a regression.

**Q: Can I do this in phases?**

A: No. This is a single-shot refactor. All 7 modules + infrastructure in one commit. This makes it easier to revert if needed.

**Q: What if the templates are confusing?**

A: Look at existing templates in `pkg/schemagen/templates/`. The `schema_validate.go.tmpl` already does field iteration - you can copy that pattern for Decode().

## Testing Strategy

**After each phase:**

1. **Unit tests:** `go test ./pkg/yappgen/modules/...`
2. **Integration tests:** `go test ./...`
3. **SCAD comparison:** Generate SCAD and diff with baseline
4. **Linter:** `golangci-lint run` (after Phase 3)

**Final validation:**

1. All tests pass
2. SCAD output identical to baseline
3. Can add a test module following new guide
4. Linter catches violations

## If You Get Stuck

**Code questions:**
- Check the architecture guide pseudocode
- Look at existing module code for patterns
- Read the debate round linked in the architecture guide

**Decision questions:**
- Read the relevant debate round
- Check the options analysis document

**Build/test questions:**
- Run `go test -v` for verbose output
- Check error messages carefully
- Compare with baseline SCAD output

**Template questions:**
- Look at existing templates in `pkg/schemagen/templates/`
- Test template changes on one module first
- Regenerate and check generated code

## Success Checklist

When you're done, verify:

- [ ] All 7 modules use generated Decode()
- [ ] No yaml.Marshal/Unmarshal in module.go files
- [ ] All modules return []ArrayDecl
- [ ] Tests pass (`go test ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] SCAD output identical to baseline
- [ ] Documentation updated
- [ ] Can add new module following guide

## Estimated Timeline

- **Phase 1 (Step A):** 1-2 days
  - Helpers: 3-4 hours
  - Templates: 4-6 hours
  - Module updates: 2-3 hours

- **Phase 2 (Step B):** 1 day
  - ArrayDecl + interface: 2-3 hours
  - Wrappers: 2-3 hours
  - multiArrayFeatureModule: 2-3 hours

- **Phase 3 (Enforcement):** 0.5 days
  - Linter: 1 hour
  - Docs: 2-3 hours

- **Phase 4 (Validation):** 0.5 days
  - Testing: 2-3 hours
  - Verification: 1 hour

**Total:** 3-4 days

## Next Steps

1. Read the architecture guide (30 min)
2. Set up baseline (15 min)
3. Start with Task 1: Create helpers package
4. Work through tasks sequentially
5. Test after each phase
6. Ask questions if stuck!

Good luck! 🚀
