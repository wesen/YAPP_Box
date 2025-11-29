---
Title: Getting Started - Composite Module System Implementation
Ticket: YAPP-COMPOSITE-001
Status: active
Topics:
  - dsl
  - modules
  - implementation
  - onboarding
DocType: playbook
Intent: long-term
Owners: []
RelatedFiles:
  - Path: /home/manuel/workspaces/2025-11-28/yapp-film-developer/YAPP_Box/ttmp/2025/11/28/YAPP-MODULE-PIPELINE-001-module-pipeline-system-multi-array-module-support/index.md
    Note: Parent ticket with complete analysis and design decision process
  - Path: /home/manuel/workspaces/2025-11-28/yapp-film-developer/YAPP_Box/ttmp/2025/11/28/YAPP-COMPOSITE-001-implement-composite-module-system-approach-3/design/01-composite-module-system-design.md
    Note: Complete design document with API sketches and implementation plan
Summary: Onboarding guide for new developers implementing the composite module system
LastUpdated: 2025-11-28
---

# Getting Started - Composite Module System Implementation

## Welcome!

This document helps you get started implementing the composite module system for YAPP. It provides all the context, background, and guidance you need to understand what this is about and begin implementation.

**Estimated time to read:** 15-20 minutes
**Estimated time to first code:** 30 minutes after reading

---

## What Is This Project?

### The Problem

YAPP users need to add LCD displays, sensors, and other composite components to their enclosures. These components require **multiple types of features**:
- A display cutout (hole in the enclosure)
- Mounting holes (to secure the display PCB)
- Sometimes additional features (connectors, light pipes, etc.)

**Currently, users must manually configure each piece:**
```yaml
features:
  cutouts:
    - face: front
      from_face_left: 50
      from_face_bottom: 30
      width: 80
      length: 40
      shape: rectangle
  
  pcb_stands:
    - x: 45
      y: 25
      # ... 16 more lines for 4 mounting holes
```

**This is:**
- Repetitive (~30 lines for one display)
- Error-prone (manual coordinate calculations)
- Hard to validate (no way to ensure cutout and mounts align)

### The Solution

**Composite modules** let users specify complete components in one place:
```yaml
features:
  lcd:
    - display:
        face: front
        position: [50, 30]
        size: [80, 40]
      mounting:
        pattern: rectangle
        spacing: [34.0, 25.0]
        hole_diameter: 2.6
```

**The composite module:**
- Generates the cutout entry automatically
- Generates mounting hole entries automatically
- Ensures they're properly aligned
- Reduces configuration from ~30 lines to ~10 lines

---

## What You're Building

### System Overview

You're implementing a **post-processing system** that runs after regular modules:

```
Input YAML
    ↓
[Resolver] - Resolves expressions, validates
    ↓
[BuildModel]
    ├─ collectFeatureModules() - Regular modules (pcb_stands, cutouts, etc.)
    └─ processCompositeModules() - NEW! Composite modules (lcd, sensor, etc.)
    ↓
Model (with arrays from both regular and composite modules)
    ↓
[EmitSCAD] - Generates OpenSCAD code
    ↓
SCAD Output
```

### Key Components

**1. pkg/composite/ Package**
- Core interfaces and utilities for composite modules
- Registry for managing composite modules
- Helper functions for type conversion and array construction
- Provenance tracking for error attribution

**2. Integration Point**
- 3-line hook in `pkg/yappgen/model.go` BuildModel function
- Calls `composite.ProcessCompositeModules()` after regular modules

**3. First Module: LCD**
- Example implementation in `pkg/composite/modules/lcd/`
- Generates cutouts + mounting holes
- ~100 lines of code

**4. Tooling and Documentation**
- Auto-generate array format docs from existing schemas
- Auto-generate array construction helpers
- Composite module authoring guide

---

## Background Context

### Where This Came From

**Parent Ticket:** YAPP-MODULE-PIPELINE-001

**Design Process:**
1. **Analysis Phase** - Analyzed current module system, identified gaps
2. **Brainstorming Phase** - Explored 5 different architectural approaches
3. **Debate Phase** - 6 rounds of technical debate evaluating approaches
4. **Decision Phase** - Selected Approach 3 (Post-Processing Array Merge)

**Key Documents to Read:**
1. **Design Doc** (this ticket): `design/01-composite-module-system-design.md` - Complete spec
2. **Module System Analysis** (parent ticket): Understanding current limitations
3. **Debate Round 2** (parent ticket): Detailed API walk-throughs for Approach 1 vs Approach 3

### Why Approach 3?

**5 approaches were considered:**
1. Pre-Processing DSL Transform
2. DSL Generation and Merging
3. **Post-Processing Array Merge** ← Selected
4. Hybrid Approach
5. DSL Macros

**Approach 3 was selected because:**
- ✅ Simpler implementation (2-3 weeks vs 3-4 weeks for Approach 1)
- ✅ Sufficient for 2-5 composite modules
- ✅ Minimal disruption (4 lines changed in existing code)
- ✅ Good DX with 1 day tooling investment
- ✅ Clear migration path to Approach 1 if needed

**Trade-off:** No expression generation support (Approach 1 feature), but can migrate later if needed.

---

## How the System Works

### The Module Lifecycle

**1. User writes YAML:**
```yaml
features:
  lcd:
    - display:
        face: front
        position: [50, 30]
        size: [80, 40]
      mounting:
        positions: [[45, 25], [105, 55]]
        standoff_height: 5
```

**2. Resolver processes DSL:**
- Validates structure
- Resolves expressions (e.g., `pcb.length / 2` → `45.0`)
- Validates constraints

**3. BuildModel collects regular features:**
```go
model.PcbStands = []  // Empty (user didn't define any)
model.Cutouts = []    // Empty (user didn't define any)
```

**4. processCompositeModules() runs:**
- Detects `features.lcd`
- Looks up LCD module in registry
- Calls `lcdModule.PostProcess(model, data, resolved)`

**5. LCD module appends to arrays:**
```go
// Generate cutout
cutout := composite.CreateCutoutEntry("front", 50, "rectangle", opts)
model.Cutouts = append(model.Cutouts, cutout)

// Generate mounting holes
stand1 := composite.CreatePcbStandEntry(45, 25, opts)
model.PcbStands = append(model.PcbStands, stand1)
// ... more stands
```

**6. EmitSCAD generates OpenSCAD:**
```openscad
cutoutsFront = [
  [50, 30, 80, 40, 0, yappRectangle]
];

pcbStands = [
  [45, 25, 5, undef, 3, ...]
];
```

**Result:** User's concise YAML becomes complete SCAD with cutouts and mounting holes.

---

## Core Concepts

### 1. CompositeModule Interface

```go
type CompositeModule interface {
    Name() string  // "lcd"
    Path() string  // "features.lcd"
    
    // PostProcess appends entries to Model arrays
    PostProcess(ctx, model, moduleData, resolved) error
}
```

**What PostProcess does:**
- Receives Model with arrays from regular modules
- Receives moduleData (the lcd configuration from YAML)
- Receives resolved (full document for context)
- **Appends** entries to model.Cutouts, model.PcbStands, etc.
- Tracks provenance for error messages

### 2. Registry Pattern

```go
var globalRegistry *Registry

composite.Register(lcdModule)  // Register on init

module, found := composite.Get("features.lcd")  // Lookup
if found {
    module.PostProcess(...)
}
```

**Same pattern as regular modules** - familiar to developers.

### 3. APPEND Merge Strategy

**When composite module generates entries:**
```
result = user_entries ++ composite_module_entries
```

**Order is deterministic:**
- User entries come first
- Composite entries appended
- Multiple composite modules: registration order

**Example:**
```
User defines: cutouts = [cutout_A]
LCD generates: cutout_B
Result: cutouts = [cutout_A, cutout_B]
```

### 4. Provenance Tracking

**Purpose:** Clear error messages showing which module generated which entry.

**Without provenance:**
```
Error: features.cutouts[2].shape: invalid value
```

**With provenance:**
```
Error: features.cutouts[2].shape: invalid value
       (generated by composite module 'lcd')
```

**How it works:**
```go
startIdx := len(model.Cutouts)
model.Cutouts = append(model.Cutouts, cutout)
model.CompositeProvenance.RecordSource("Cutouts", startIdx, "lcd")
```

---

## Architecture Deep Dive

### Package Structure

```
pkg/composite/
├── module.go          # CompositeModule interface definition
├── registry.go        # Registry for composite modules
├── processor.go       # ProcessCompositeModules() - called from BuildModel
├── helpers.go         # Type conversion: ToFloat, GetString, etc.
├── builders.go        # Array construction: CreateCutoutEntry, etc. (auto-generated)
├── provenance.go      # Provenance tracking for error attribution
└── modules/
    └── lcd/
        ├── module.go      # LCD composite module implementation
        └── module_test.go # Unit tests
```

### Integration Points

**1. In pkg/yappgen/model.go:**
```go
type Model struct {
    // ... existing fields ...
    CompositeProvenance *composite.Provenance  // NEW - 1 line
}

func BuildModel(...) (*Model, error) {
    // ... existing code ...
    
    // NEW - 3 lines
    if err := composite.ProcessCompositeModules(ctx, m, resolved, features); err != nil {
        return nil, errors.Wrap(err, "composite modules")
    }
    
    return m, nil
}
```

**Total changes to existing code: 4 lines**

**2. In pkg/composite/processor.go:**
```go
func ProcessCompositeModules(ctx, model, resolved, features) error {
    // Iterate features
    // Lookup composite modules in registry
    // Call PostProcess for each
    // Track errors
}
```

**New code: ~50-70 lines**

---

## What Makes This Different from Regular Modules

| Aspect | Regular Modules | Composite Modules |
|--------|----------------|-------------------|
| **Package** | `pkg/yappgen/modules/` | `pkg/composite/modules/` |
| **Interface** | `registry.FeatureModule` | `composite.CompositeModule` |
| **Processing** | During `collectFeatureModules` | After `collectFeatureModules` |
| **Output** | One array type | Multiple array types |
| **Input** | `features.pcb_stands` | `features.lcd` |
| **Function** | `Build() → []ArrayDecl` | `PostProcess() → error` |
| **Schema** | `schema.yaml` required | Optional (DSL structure only) |
| **Code Gen** | Extensive (structs, tests) | Minimal (helpers only) |

**Key difference:** Composite modules **append to multiple Model arrays** directly.

---

## Implementation Roadmap

### Phase 1: Core Infrastructure (3 days)

**Goal:** Build the foundation - interfaces, registry, processor, helpers.

**Files to create:**
- `pkg/composite/module.go` (~30 lines)
- `pkg/composite/registry.go` (~80 lines)
- `pkg/composite/processor.go` (~60 lines)
- `pkg/composite/helpers.go` (~100 lines)
- `pkg/composite/provenance.go` (~80 lines)
- Tests (~200 lines)

**Integration:**
- Modify `pkg/yappgen/model.go` (4 lines)

**Success criteria:**
- Registry can register and lookup modules
- Helpers convert types correctly
- Provenance tracks sources
- ProcessCompositeModules iterates correctly
- All tests pass

### Phase 2: Tooling and Documentation (2 days)

**Goal:** Auto-generate docs and helpers from existing schemas.

**Files to create:**
- `cmd/schemagen/docs.go` (~150 lines)
- `cmd/schemagen/builders.go` (~200 lines - template)
- `pkg/composite/builders.go` (auto-generated, ~300 lines)
- `pkg/docs/composite-array-formats.md` (auto-generated)

**Success criteria:**
- `schemagen docs` generates complete array format reference
- `schemagen builders` generates CreateXXXEntry helpers
- Documentation is clear and comprehensive

### Phase 3: LCD Module (2-3 days)

**Goal:** Implement first composite module as proof of concept.

**Files to create:**
- `pkg/composite/modules/lcd/module.go` (~150 lines)
- `pkg/composite/modules/lcd/module_test.go` (~150 lines)
- `examples/composite/test-lcd-basic.yaml`
- `examples/composite/test-lcd-complex.yaml`

**Success criteria:**
- LCD module generates cutout + mounting holes
- Unit tests pass
- Integration test: YAML → SCAD works end-to-end
- Generated SCAD can be rendered in OpenSCAD

### Phase 4: Documentation and Polish (2-3 days)

**Goal:** Complete developer documentation and examples.

**Files to create:**
- `pkg/docs/tutorials/composite-module-authoring-guide.md`
- Multiple example YAML files
- Troubleshooting guide

**Success criteria:**
- Guide explains how to write composite modules
- Examples cover common patterns
- Troubleshooting addresses anticipated issues

**Total: 2-3 weeks to complete system**

---

## Prerequisites

### Knowledge Required

**1. Go Programming**
- Interfaces and structs
- Error handling patterns
- Testing with `testing` package
- Basic generics (for type conversion helpers)

**2. YAPP Module System**
- Read: `pkg/docs/tutorials/yapp-module-authoring-guide.md`
- Understand: How regular modules work (Build functions, schemas, registry)
- Study: `pkg/yappgen/modules/pcbstands/` (simple example)
- Study: `pkg/yappgen/modules/cutouts/` (multi-array example)

**3. YAPP DSL Pipeline**
- Understand: Resolver → BuildModel → EmitSCAD flow
- Read: `pkg/docs/tutorials/yapp-dsl-reference.md` (pipeline section)
- Trace: `pkg/cli/generatorcli/generator.go` WriteSCAD function

**4. Testing**
- Unit testing with table-driven tests
- Integration testing with YAML files
- Understanding of test fixtures and assertions

### Environment Setup

```bash
# Clone and build
cd /home/manuel/workspaces/2025-11-28/yapp-film-developer/YAPP_Box

# Verify build works
go build ./...

# Run tests
go test ./...

# Test CLI
go run ./cmd/yappctl help

# Generate an example to understand pipeline
go run ./cmd/yappctl generate --input examples/yapp-demo-buttons.yaml --scad-out /tmp/test.scad
```

---

## Key Files to Understand

### Before You Start Coding

**Read these files to understand the context:**

**1. Design Document** (30 min read)
```bash
cat ttmp/2025/11/28/YAPP-COMPOSITE-001-.../design/01-composite-module-system-design.md
```
- Complete API specifications
- Implementation phases
- Testing strategy
- Success metrics

**2. Current Module System** (20 min)
```bash
# Simple module example
cat pkg/yappgen/modules/pcbstands/module.go
cat pkg/yappgen/modules/pcbstands/registry.go

# Complex module example
cat pkg/yappgen/modules/cutouts/module.go
cat pkg/yappgen/modules/cutouts/registry.go
```

**3. Pipeline Integration** (15 min)
```bash
# See where BuildModel is called
cat pkg/cli/generatorcli/generator.go

# See what BuildModel does
cat pkg/yappgen/model.go
```

**4. Debate Rounds for Deep Understanding** (optional, 1-2 hours)

If you want to understand WHY decisions were made:
- Round 1: Foundation - Why build this?
- Round 2: Architecture - Why Approach 3?
- Round 3: Pipeline - Where to integrate?
- Round 4: Conflicts - Why APPEND strategy?
- Round 5: Validation - When to validate?
- Round 6: Developer Experience - What tooling is needed?

**Location:** `ttmp/2025/11/28/YAPP-MODULE-PIPELINE-001-.../debate/`

---

## Implementation Guide

### Step 1: Understand the Model Structure

**Read `pkg/yappgen/model.go`:**

```go
type Model struct {
    // Arrays that composite modules can append to:
    PcbStands   []map[string]any  // PCB mounting standoffs
    Connectors  []map[string]any  // Connectors
    BoxMounts   []map[string]any  // Box mounting points
    SnapJoins   []map[string]any  // Snap-fit joints
    Cutouts     []map[string]any  // Face cutouts
    PushButtons []map[string]any  // Push button extenders
    LightTubes  []map[string]any  // Light guide tubes
}
```

**Key insight:** All arrays are `[]map[string]any` - composite modules append maps to these slices.

### Step 2: Understand Array Formats

**Each array has a different format:**

**Cutouts:**
```go
{
    "face": "front",              // Required: front, back, left, right, lid, base
    "from_face_left": 50.0,       // Required: horizontal position
    "from_face_bottom": 30.0,     // Optional: vertical (side faces)
    "width": 80.0,                // Optional: width dimension
    "length": 40.0,               // Optional: length dimension
    "shape": "rectangle",         // Required: rectangle, circle, etc.
}
```

**PcbStands:**
```go
{
    "x": 45.0,         // Required: X coordinate
    "y": 25.0,         // Required: Y coordinate
    "height": 5.0,     // Optional: standoff height
    "diameter": 3.0,   // Optional: standoff diameter
}
```

**Where to find this info:**
- Look at module schemas: `pkg/yappgen/modules/*/schema.yaml`
- Read module builders: `pkg/yappgen/modules/*/module.go`
- **Phase 2 will auto-generate** complete reference doc

### Step 3: Start with Phase 1

**Begin with the core infrastructure:**

```bash
mkdir -p pkg/composite
cd pkg/composite
```

**Create files in order:**
1. `module.go` - CompositeModule interface
2. `registry.go` - Registry implementation
3. `helpers.go` - Type conversion utilities
4. `provenance.go` - Source tracking
5. `processor.go` - ProcessCompositeModules implementation

**Write tests alongside each file.**

**Pattern to follow:** Look at `pkg/registry/` for registry patterns.

### Step 4: Add Integration Hook

**In `pkg/yappgen/model.go`:**

Find the `BuildModel` function and add the hook after `collectFeatureModules`:

```go
// Collect features (existing)
if err := collectFeatureModules(resolved, features, m); err != nil {
    return nil, err
}

// Process composite modules (NEW)
if err := composite.ProcessCompositeModules(ctx, m, resolved, features); err != nil {
    return nil, errors.Wrap(err, "composite modules")
}
```

**Test:** Verify existing examples still work (no composite modules yet).

### Step 5: Build Tooling (Phase 2)

**Auto-generate helpers and docs:**

```bash
# Add to cmd/schemagen
go run ./cmd/schemagen docs --output pkg/docs/composite-array-formats.md
go run ./cmd/schemagen builders --output pkg/composite/builders.go
```

**Verify:**
- Array format docs are comprehensive
- Builder helpers cover all module types
- Documentation is clear

### Step 6: Implement LCD Module (Phase 3)

**Create LCD module:**

```bash
mkdir -p pkg/composite/modules/lcd
```

**Follow this structure:**
```go
// module.go
type Module struct{}

func (m *Module) PostProcess(ctx, model, moduleData, resolved) error {
    // 1. Parse moduleData (extract display, mounting configs)
    // 2. Generate cutout using CreateCutoutEntry helper
    // 3. Append to model.Cutouts with provenance
    // 4. Generate mounting holes using CreatePcbStandEntry helper
    // 5. Append to model.PcbStands with provenance
    return nil
}
```

**Register:**
```go
func init() {
    composite.Register(&Module{})
}
```

### Step 7: Test End-to-End

**Create test YAML:**
```yaml
# examples/composite/test-lcd.yaml
features:
  lcd:
    - display:
        face: front
        position: [50, 30]
        size: [80, 40]
      mounting:
        positions: [[45, 25], [105, 55]]
        standoff_height: 5
        diameter: 3
```

**Generate and verify:**
```bash
go run ./cmd/yappctl generate --input examples/composite/test-lcd.yaml --scad-out /tmp/test.scad
cat /tmp/test.scad | grep -A5 "cutoutsFront\|pcbStands"
```

**Expected output:**
- cutoutsFront array with 1 entry (LCD display cutout)
- pcbStands array with 2 entries (mounting holes)

---

## Common Patterns and Examples

### Pattern 1: Parse Module Data

```go
func (m *Module) PostProcess(ctx, model, moduleData, resolved) error {
    // Extract array
    items, err := composite.ToArray(moduleData)
    if err != nil {
        return errors.Wrap(err, "lcd: expected array")
    }
    
    for i, item := range items {
        itemMap, err := composite.ToMap(item)
        if err != nil {
            continue
        }
        
        // Process item
    }
}
```

### Pattern 2: Extract Configuration

```go
// Get required field
display, found := composite.GetMap(itemMap, "display")
if !found {
    return errors.Errorf("lcd[%d]: missing display configuration", i)
}

// Get field with default
face := composite.GetStringOrDefault(display, "face", "front")

// Get required array
position, found := composite.GetArray(display, "position")
if !found {
    return errors.Errorf("lcd[%d].display: missing position", i)
}
```

### Pattern 3: Generate and Append Entry

```go
// Generate cutout entry
cutout := composite.CreateCutoutEntry(
    face,
    posX,
    "rectangle",
    &composite.CutoutOptions{
        FromFaceBottom: &posY,
        Width:          &width,
        Length:         &height,
    },
)

// Append with provenance
startIdx := len(model.Cutouts)
model.Cutouts = append(model.Cutouts, cutout)
model.CompositeProvenance.RecordSource("Cutouts", startIdx, "lcd")
```

### Pattern 4: Loop and Generate Multiple Entries

```go
// Generate mounting holes
for j, pos := range mountingPositions {
    posArray, _ := composite.ToArray(pos)
    x, _ := composite.ToFloat(posArray[0])
    y, _ := composite.ToFloat(posArray[1])
    
    stand := composite.CreatePcbStandEntry(x, y, opts)
    
    idx := len(model.PcbStands)
    model.PcbStands = append(model.PcbStands, stand)
    model.CompositeProvenance.RecordSource("PcbStands", idx, "lcd")
}
```

---

## Testing Strategy

### Unit Tests

**Test each component in isolation:**

```go
// Test registry
func TestRegistry_RegisterAndGet(t *testing.T)
func TestRegistry_DuplicatePanic(t *testing.T)

// Test helpers
func TestToFloat_ValidInputs(t *testing.T)
func TestToFloat_InvalidInputs(t *testing.T)

// Test provenance
func TestProvenance_RecordAndGet(t *testing.T)

// Test LCD module
func TestLCD_PostProcess_BasicDisplay(t *testing.T)
func TestLCD_PostProcess_WithExistingEntries(t *testing.T)
func TestLCD_PostProcess_InvalidConfig(t *testing.T)
```

### Integration Tests

**Test full pipeline with YAML files:**

```bash
# Test: Basic LCD
go run ./cmd/yappctl generate --input examples/composite/test-lcd-basic.yaml --scad-out /tmp/lcd-basic.scad

# Test: LCD + user cutouts (APPEND behavior)
go run ./cmd/yappctl generate --input examples/composite/test-lcd-with-user-cutouts.yaml --scad-out /tmp/lcd-append.scad

# Test: Multiple composite modules
go run ./cmd/yappctl generate --input examples/composite/test-multiple-composites.yaml --scad-out /tmp/multi.scad
```

**Verify SCAD output:**
- Arrays contain expected entries
- Order is correct (user first, then composite)
- No duplicate or missing entries

---

## Troubleshooting Guide

### Common Issues

**Issue 1: "unknown composite module"**
```
Error: features.lcd not recognized
```
**Fix:** Module not registered. Add `composite.Register(&lcdModule{})` in `init()`.

**Issue 2: "type assertion failed"**
```
panic: interface conversion: interface {} is float64, not string
```
**Fix:** Use helper functions: `composite.ToFloat(v)` instead of `v.(float64)`.

**Issue 3: "index out of range"**
```
panic: runtime error: index out of range [1] with length 1
```
**Fix:** Validate array length before accessing: `if len(arr) < 2 { return error }`.

**Issue 4: "cutouts not rendered"**
```
SCAD generated but no cutouts visible
```
**Fix:** Check provenance - was entry actually appended? Check array field names match schema.

---

## Quick Start Checklist

Before you write any code:

- [ ] Read this document completely
- [ ] Read design document (`design/01-composite-module-system-design.md`)
- [ ] Understand current module system (read pcbstands and cutouts modules)
- [ ] Trace pipeline flow (Resolver → BuildModel → EmitSCAD)
- [ ] Set up development environment (build and test)

Phase 1 (Core Infrastructure):

- [ ] Create `pkg/composite/module.go` with CompositeModule interface
- [ ] Create `pkg/composite/registry.go` with Registry
- [ ] Create `pkg/composite/helpers.go` with type conversion functions
- [ ] Create `pkg/composite/provenance.go` with source tracking
- [ ] Create `pkg/composite/processor.go` with ProcessCompositeModules
- [ ] Add 4-line integration in `pkg/yappgen/model.go`
- [ ] Write unit tests for all components
- [ ] Verify: `go test ./pkg/composite/` passes

Phase 2 (Tooling):

- [ ] Implement `cmd/schemagen docs` command
- [ ] Implement `cmd/schemagen builders` command  
- [ ] Generate `pkg/composite/builders.go`
- [ ] Generate `pkg/docs/composite-array-formats.md`
- [ ] Verify: Documentation is clear and comprehensive

Phase 3 (LCD Module):

- [ ] Create `pkg/composite/modules/lcd/module.go`
- [ ] Implement PostProcess with cutout + mounting hole generation
- [ ] Write unit tests for LCD module
- [ ] Create example YAML files
- [ ] Verify: End-to-end test generates correct SCAD

Phase 4 (Documentation):

- [ ] Write composite module authoring guide
- [ ] Create multiple example YAML files
- [ ] Write troubleshooting guide
- [ ] Update main documentation

---

## Next Steps

1. **Read the design doc thoroughly** - Understand all interfaces and patterns
2. **Study existing modules** - See how regular modules work
3. **Set up development environment** - Build, test, run examples
4. **Start with Phase 1** - Build core infrastructure
5. **Ask questions early** - Don't guess, clarify design decisions
6. **Test continuously** - Write tests alongside implementation
7. **Document as you go** - Capture decisions and patterns

---

## Getting Help

### Documentation Resources

**Design and Architecture:**
- Design doc: `design/01-composite-module-system-design.md`
- Parent ticket analysis: `../YAPP-MODULE-PIPELINE-001-.../analysis/01-module-system-analysis.md`
- Debate rounds: `../YAPP-MODULE-PIPELINE-001-.../debate/` (6 rounds)

**Code Examples:**
- Regular modules: `pkg/yappgen/modules/*/module.go`
- Registry pattern: `pkg/registry/registry.go`
- Pipeline flow: `pkg/cli/generatorcli/generator.go`

**Testing Examples:**
- Module tests: `pkg/yappgen/modules/*/module_test.go`
- Integration: Example YAML files in `examples/`

### Questions to Consider

**Before starting Phase 1:**
- Do I understand the difference between regular and composite modules?
- Do I know what the Model structure looks like?
- Can I trace the pipeline from YAML to SCAD?

**During Phase 1:**
- Are my interfaces following Go conventions?
- Are errors properly wrapped with context?
- Are tests covering edge cases?

**During Phase 3:**
- Is the LCD module following the patterns from design doc?
- Are provenance entries tracked for all generated arrays?
- Does end-to-end test produce correct SCAD?

---

## Success Criteria

**You'll know you're done when:**

Phase 1:
- ✅ All core interfaces implemented and tested
- ✅ Registry can manage modules
- ✅ Helpers handle type conversions
- ✅ Provenance tracks sources
- ✅ Integration hook compiles and tests pass

Phase 2:
- ✅ Array format docs auto-generated
- ✅ Builder helpers auto-generated
- ✅ Documentation is comprehensive

Phase 3:
- ✅ LCD module generates cutouts + mounting holes
- ✅ Example YAML generates correct SCAD
- ✅ Can render STL in OpenSCAD

Phase 4:
- ✅ Authoring guide is clear for future developers
- ✅ Multiple examples demonstrate patterns
- ✅ System is ready for second composite module

---

## Migration Path to Approach 1 (Future)

**If you later need expression generation:**

The design includes a migration path from Approach 3 → Approach 1.

**When to migrate:**
- Need to generate expressions (not just use resolved values)
- Have 5+ composite modules (complexity justifies better abstraction)
- Coupling to Model becomes maintenance burden

**How to migrate:**
- Refactor PostProcess → Transform
- Move processing from BuildModel to CLI layer (before resolver)
- Estimated effort: 1-2 weeks

**Interfaces are designed for easy migration** - parameters are similar, just different return types.

---

## Additional Resources

### Parent Ticket (YAPP-MODULE-PIPELINE-001)

**Location:** `ttmp/2025/11/28/YAPP-MODULE-PIPELINE-001-module-pipeline-system-multi-array-module-support/`

**Contains:**
- Complete analysis of current module system
- 5 architectural approaches explored
- 6 rounds of technical debate
- Comparison matrices and decision rationale

**Read if you want to understand:** Why this design, what alternatives were considered, what trade-offs were made.

### Related Tickets

**YAPP-PUSH-BUTTONS-001** - Module system design
- How the current module system was designed
- Feature module registry architecture
- Schema-based validation approach

**YAPP-BUILDER-REFACTOR-001** - ArrayDecl refactor
- How modules return arrays
- Registry interface evolution

---

## Final Notes

**Philosophy:** Start simple, iterate based on real usage.

**Approach 3 is pragmatic:**
- Solves the immediate need (2-3 composite modules)
- Minimal disruption to existing code
- Fast implementation
- Clear migration path if needs evolve

**Your mission:**
- Build the system as designed
- Document clearly for future developers
- Test thoroughly
- Ship the LCD module

**You've got this!** The design is solid, the patterns are clear, and the path is well-defined. 🚀

