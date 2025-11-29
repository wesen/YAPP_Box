---
Title: Composite Module System - Design Document
Ticket: YAPP-MODULE-PIPELINE-001
Status: active
Topics:
  - dsl
  - modules
  - codegen
  - architecture
  - design
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
  - Path: /home/manuel/workspaces/2025-11-28/yapp-film-developer/YAPP_Box/pkg/yappgen/model.go
    Note: Model structure that composite modules modify - integration point
  - Path: /home/manuel/workspaces/2025-11-28/yapp-film-developer/YAPP_Box/pkg/yappgen/features.go
    Note: Feature module coordination - where composite hook will be added
  - Path: /home/manuel/workspaces/2025-11-28/yapp-film-developer/YAPP_Box/pkg/registry/registry.go
    Note: Existing registry pattern to follow for composite modules
Summary: Complete design for composite module system using Approach 3 (Post-Processing Array Merge) with future migration path to Approach 1
LastUpdated: 2025-11-28
---

# Composite Module System - Design Document

## Executive Summary

This design document specifies a **composite module system** that enables DSL modules to generate entries across multiple array types (e.g., an LCD module that adds both cutouts and mounting holes). The design uses **Approach 3: Post-Processing Array Merge** as the initial implementation, with a clear migration path to **Approach 1: Pre-Processing DSL Transform** if needed for expression generation support.

**Key Design Decisions:**
- **Integration Point:** Post-processing hook in BuildModel (after regular modules)
- **Conflict Strategy:** APPEND - composite entries appended to existing arrays
- **Validation:** Reuse existing validators, track provenance for error attribution
- **Developer Experience:** Invest in tooling (auto-generated docs, helpers, examples)

**Estimated Effort:**
- Core implementation: 1-2 weeks
- Tooling and documentation: 1 day
- First module (LCD): 2-3 days
- **Total: ~2-3 weeks**

**Migration Path:** If expression generation becomes necessary, refactor to Approach 1 (estimated 1-2 weeks additional).

---

## Problem Statement

### Current Limitation

The existing module system allows modules to output multiple arrays of the **same type** (e.g., cutouts module outputs `cutoutsFront`, `cutoutsBack`, etc.), but **cannot** output arrays of **different types** (e.g., cutouts + pcbStands).

### User Pain Points

**Example: LCD Display Module**

Users currently must manually specify:
```yaml
vars:
  display_window_width: 36.0
  display_window_height: 17.0
  display_center_y: (row_center_y - 4.0) + button_radius + 6.0 + (display_window_height / 2)
  display_mount_dx: 17.0
  display_mount_dy: 12.5

features:
  cutouts:
    - face: lid
      from_face_back: vars.row_center_x
      from_face_left: vars.display_center_y
      width: vars.display_window_width
      length: vars.display_window_height
      shape: rectangle
  
  # Display mounting holes (4 manual entries)
  # ... 16 lines of repetitive YAML
```

**Problems:**
- ~30 lines of YAML + 6 variables for one display
- Manual coordinate calculations
- No validation that cutout and mounts align
- Error-prone and repetitive

**Desired Solution:**
```yaml
features:
  lcd:
    - display:
        face: lid
        position: [row_center_x, display_center_y]
        size: [36.0, 17.0]
      mounting:
        pattern: rectangle
        spacing: [34.0, 25.0]
        hole_diameter: 2.6
```

**Benefits:**
- ~10 lines total (70% reduction)
- No manual calculations
- Module ensures alignment
- Clear, declarative intent

---

## Goals and Non-Goals

### Goals

1. ✅ Enable composite modules that output to multiple array types
2. ✅ Maintain clean separation from existing module system
3. ✅ Provide good developer experience for module authors
4. ✅ Support 2-5 composite modules initially
5. ✅ Design for future migration to Approach 1 if needed

### Non-Goals

1. ❌ Expression generation support (Approach 1 feature, deferred)
2. ❌ Overlap detection for geometry validation (future work)
3. ❌ Cross-module dependencies (keep modules independent)
4. ❌ Support for 10+ composite modules initially (optimize for 2-5)
5. ❌ Backward compatibility with legacy displayMounts (separate ticket)

---

## Proposed Architecture

### System Overview

```
Input YAML
    ↓
[Resolver]
    - Validates structure
    - Resolves expressions
    - Validates constraints
    ↓
Resolved YAML
    ↓
[BuildModel]
    - collectFeatureModules() → populates Model arrays
    - processCompositeModules() → appends to Model arrays ← NEW
    ↓
Model (with merged arrays)
    ↓
[EmitSCAD]
    ↓
SCAD Output
```

### Integration Point

**Location:** `pkg/yappgen/model.go` in `BuildModel()` function

**Before (current):**
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

**After (with composite modules):**
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
    
    // Process composite modules (NEW - 3 lines)
    if err := processCompositeModules(ctx, m, resolved, features); err != nil {
        return nil, errors.Wrap(err, "composite modules")
    }
    
    return m, nil
}
```

---

## Core API Design

### Package Structure

```
pkg/composite/
├── module.go          # CompositeModule interface
├── registry.go        # Module registration
├── processor.go       # processCompositeModules() implementation
├── helpers.go         # Type conversion helpers
├── builders.go        # Array construction helpers (auto-generated)
├── provenance.go      # Source tracking for error attribution
└── modules/
    ├── lcd/
    │   ├── module.go      # LCD composite module
    │   └── module_test.go # Tests
    └── (future composite modules)
```

### Interface Definitions

**Core Module Interface:**

```go
package composite

// CompositeModule processes DSL and appends to Model arrays.
// Executes after regular feature modules have populated the Model.
type CompositeModule interface {
    // Name returns the module identifier (e.g., "lcd")
    Name() string
    
    // Path returns the DSL path (e.g., "features.lcd")
    Path() string
    
    // PostProcess reads module DSL data and appends entries to Model arrays.
    // 
    // Parameters:
    //   - model: The Model with arrays already populated by regular modules
    //   - moduleData: Raw data from features.lcd (from resolved document)
    //   - resolved: Complete resolved document (for accessing other values)
    //
    // Returns error if processing fails.
    //
    // Note: All values in moduleData and resolved are numeric (expressions already resolved).
    PostProcess(ctx context.Context, model *yappgen.Model, moduleData any, resolved map[string]any) error
}
```

**Registry Interface:**

```go
package composite

// Registry manages composite module registration and lookup.
type Registry struct {
    // modules maps "features.modulename" → CompositeModule
    // Private field, not exported
}

// Register adds a composite module to the registry.
// Panics if module is nil or path is already registered.
func (r *Registry) Register(module CompositeModule)

// Get retrieves a composite module by DSL path.
// Returns (module, true) if found, (nil, false) if not found.
func (r *Registry) Get(path string) (CompositeModule, bool)

// All returns all registered modules in registration order.
func (r *Registry) All() []CompositeModule

// Global registry instance
var globalRegistry *Registry

// Register adds module to global registry (convenience function)
func Register(module CompositeModule)
```

**Processor Interface:**

```go
package composite

// processCompositeModules iterates through features and executes registered composite modules.
// Called from BuildModel after collectFeatureModules.
//
// Parameters:
//   - model: Model with arrays populated by regular modules
//   - resolved: Complete resolved document
//   - features: features section from resolved document
//
// Returns error if any composite module fails.
func ProcessCompositeModules(ctx context.Context, model *yappgen.Model, resolved map[string]any, features map[string]any) error
```

---

## Helper Functions Design

### Type Conversion Helpers

```go
package composite

// ToFloat converts any to float64.
// Handles: float64, int, int64, float32
// Returns error if conversion fails.
func ToFloat(v any) (float64, error)

// ToString converts any to string.
// Returns error if v is not a string.
func ToString(v any) (string, error)

// ToArray converts any to []any.
// Returns error if v is not an array.
func ToArray(v any) ([]any, error)

// ToBool converts any to bool.
// Returns error if v is not a bool.
func ToBool(v any) (bool, error)
```

### Map Extraction Helpers

```go
package composite

// GetFloat extracts float64 from map.
// Returns (value, true) if key exists and is numeric.
// Returns (0, false) if key missing or not numeric.
func GetFloat(m map[string]any, key string) (float64, bool)

// GetString extracts string from map.
func GetString(m map[string]any, key string) (string, bool)

// GetArray extracts array from map.
func GetArray(m map[string]any, key string) ([]any, bool)

// GetMap extracts nested map from map.
func GetMap(m map[string]any, key string) (map[string]any, bool)

// GetFloatOrDefault extracts float with default value if missing.
func GetFloatOrDefault(m map[string]any, key string, defaultValue float64) float64

// GetStringOrDefault extracts string with default value if missing.
func GetStringOrDefault(m map[string]any, key string, defaultValue string) string
```

### Array Construction Helpers (Auto-Generated from Schemas)

```go
package composite

// CreateCutoutEntry creates a cutout entry for model.Cutouts.
// Auto-generated from pkg/yappgen/modules/cutouts/schema.yaml
//
// Required parameters:
//   - face: Which face (front, back, left, right, lid, base)
//   - fromFaceLeft: Horizontal position (mm)
//   - shape: Cutout shape (rectangle, circle, rounded_rect, etc.)
//
// Optional parameters (use nil for defaults):
//   - fromFaceBottom: Vertical position for side faces
//   - fromFaceBack: Depth position for lid/base
//   - width, length, radius: Dimensions
//   - depth, angle: Advanced options
//
// Returns map ready to append to model.Cutouts.
func CreateCutoutEntry(face string, fromFaceLeft float64, shape string, opts *CutoutOptions) map[string]any

// CutoutOptions holds optional parameters for CreateCutoutEntry.
type CutoutOptions struct {
    FromFaceBottom *float64
    FromFaceBack   *float64
    Width          *float64
    Length         *float64
    Radius         *float64
    Depth          *float64
    Angle          *float64
    Coordinate     string // "pcb", "box", "box_inside"
    Origin         string // "global", "center", "alt"
}

// CreatePcbStandEntry creates a PCB stand entry for model.PcbStands.
// Auto-generated from pkg/yappgen/modules/pcbstands/schema.yaml
func CreatePcbStandEntry(x, y float64, opts *PcbStandOptions) map[string]any

// PcbStandOptions holds optional parameters.
type PcbStandOptions struct {
    Height       *float64
    Diameter     *float64
    PinDiameter  *float64
    FilletRadius *float64
    Corner       string
    Coordinate   string
    // ... etc
}

// Similar helpers for:
// - CreateConnectorEntry
// - CreateBoxMountEntry
// - CreateSnapJoinEntry
// - CreateLightTubeEntry
// - CreatePushButtonEntry
```

**Note:** These are **auto-generated** from existing module schemas using schemagen.

---

## Provenance Tracking Design

### Purpose

Track which composite module generated which array entries for clear error messages.

### Interface

```go
package composite

// Provenance tracks source of composite-generated entries.
type Provenance struct {
    // sources maps array path → source module name
    // e.g., "model.Cutouts[2]" → "lcd"
}

// RecordSource tracks that an entry was generated by a module.
//
// Parameters:
//   - arrayName: "Cutouts", "PcbStands", etc.
//   - index: Entry index in array
//   - moduleName: Source module ("lcd", "sensor", etc.)
func (p *Provenance) RecordSource(arrayName string, index int, moduleName string)

// GetSource retrieves the source module for an entry.
// Returns empty string if entry was user-defined (not composite-generated).
func (p *Provenance) GetSource(arrayName string, index int) string

// FormatErrorWithSource enhances error messages with provenance.
//
// Example:
//   Input:  "features.cutouts[2].shape: invalid"
//   Output: "features.cutouts[2].shape: invalid (generated by composite module 'lcd')"
func (p *Provenance) FormatErrorWithSource(err error, path string) error
```

### Integration

```go
// In pkg/yappgen/model.go
type Model struct {
    // ... existing fields ...
    
    // CompositeProvenance tracks composite module sources (NEW)
    CompositeProvenance *composite.Provenance
}
```

---

## Module Implementation Pattern

### LCD Module Example (High-Level Sketch)

```go
package lcd

import (
    "context"
    "github.com/wesen/yapp-encl-resolver/pkg/composite"
    "github.com/wesen/yapp-encl-resolver/pkg/yappgen"
)

type Module struct{}

var _ composite.CompositeModule = (*Module)(nil)

func (m *Module) Name() string { return "lcd" }
func (m *Module) Path() string { return "features.lcd" }

// PostProcess generates cutout and mounting hole entries from LCD config.
func (m *Module) PostProcess(ctx context.Context, model *yappgen.Model, moduleData any, resolved map[string]any) error {
    // 1. Parse and validate moduleData
    items := parseItems(moduleData)  // Helper: extract []any, validate structure
    
    for i, item := range items {
        // 2. Extract display configuration
        display := extractDisplay(item)  // Helper: parse display section
        mounting := extractMounting(item) // Helper: parse mounting section
        
        // 3. Generate cutout entry using builder helper
        cutout := composite.CreateCutoutEntry(
            display.Face,
            display.Position[0],
            "rectangle",
            &composite.CutoutOptions{
                FromFaceBottom: &display.Position[1],
                Width:          &display.Size[0],
                Length:         &display.Size[1],
            },
        )
        
        // 4. Append to model with provenance tracking
        startIdx := len(model.Cutouts)
        model.Cutouts = append(model.Cutouts, cutout)
        model.CompositeProvenance.RecordSource("Cutouts", startIdx, "lcd")
        
        // 5. Generate mounting hole entries
        for j, pos := range mounting.Positions {
            stand := composite.CreatePcbStandEntry(
                pos[0], pos[1],
                &composite.PcbStandOptions{
                    Height:   &mounting.StandoffHeight,
                    Diameter: &mounting.Diameter,
                },
            )
            
            // 6. Append to model with provenance
            standIdx := len(model.PcbStands)
            model.PcbStands = append(model.PcbStands, stand)
            model.CompositeProvenance.RecordSource("PcbStands", standIdx, "lcd")
        }
    }
    
    return nil
}

// Helper types for clarity
type displayConfig struct {
    Face     string
    Position [2]float64
    Size     [2]float64
}

type mountingConfig struct {
    Positions       [][2]float64
    StandoffHeight  float64
    Diameter        float64
}

// parseItems extracts and validates array structure
func parseItems(data any) ([]map[string]any, error)

// extractDisplay parses display section from item
func extractDisplay(item map[string]any) (*displayConfig, error)

// extractMounting parses mounting section from item
func extractMounting(item map[string]any) (*mountingConfig, error)
```

**Lines of Code Estimate:**
- PostProcess function: ~40 lines
- Helper functions: ~60 lines
- Total: ~100 lines per module

---

## Documentation Auto-Generation

### Array Format Reference Generator

**Tool:** `schemagen docs` command

**Purpose:** Generate comprehensive documentation for composite module authors showing what fields each array type accepts.

**Implementation Sketch:**

```go
package main

// Command: go run ./cmd/schemagen docs --output pkg/docs/composite-array-formats.md
func generateArrayFormatDocs() {
    // 1. Iterate all registered modules
    modules := registry.All()
    
    for _, module := range modules {
        schema := module.Schema()
        
        // 2. Extract schema metadata
        // name := schema.Name()            // "pcb_stands"
        // path := schema.Path()            // "features.pcb_stands"
        // fields := schema.Fields()        // []FieldSpec
        // description := schema.Description()
        
        // 3. Generate markdown section
        // ## features.pcb_stands → model.PcbStands
        // Description: Defines PCB mounting standoffs
        //
        // Required fields:
        // - x: number - X coordinate
        // - y: number - Y coordinate
        //
        // Optional fields:
        // - height: number - Standoff height (default: pcb.z_clearance)
        // ...
        //
        // Example:
        // ```go
        // stand := map[string]any{
        //     "x": 45.0,
        //     "y": 25.0,
        //     "height": 5.0,
        // }
        // model.PcbStands = append(model.PcbStands, stand)
        // ```
    }
    
    // 4. Write to output file
}
```

**Output:** `pkg/docs/composite-array-formats.md` with complete array format reference.

---

## Testing Strategy

### Unit Tests for Composite Modules

```go
package lcd

// TestPostProcess_BasicDisplay tests basic LCD display configuration
func TestPostProcess_BasicDisplay(t *testing.T) {
    // Setup: Create empty Model
    model := &yappgen.Model{
        CompositeProvenance: composite.NewProvenance(),
    }
    
    // Input: LCD module data
    moduleData := []any{
        map[string]any{
            "display": map[string]any{
                "face":     "front",
                "position": []any{50.0, 30.0},
                "size":     []any{80.0, 40.0},
            },
            "mounting": map[string]any{
                "positions":        []any{[]any{45.0, 25.0}},
                "standoff_height":  5.0,
                "diameter":         3.0,
            },
        },
    }
    
    resolved := map[string]any{}  // Minimal resolved doc
    
    // Execute
    err := module.PostProcess(context.Background(), model, moduleData, resolved)
    
    // Assert
    assert.NoError(t, err)
    assert.Len(t, model.Cutouts, 1)
    assert.Len(t, model.PcbStands, 1)
    
    // Verify cutout
    cutout := model.Cutouts[0]
    assert.Equal(t, "front", cutout["face"])
    assert.Equal(t, 50.0, cutout["from_face_left"])
    
    // Verify provenance
    source := model.CompositeProvenance.GetSource("Cutouts", 0)
    assert.Equal(t, "lcd", source)
}

// TestPostProcess_WithExistingEntries tests APPEND behavior
func TestPostProcess_WithExistingEntries(t *testing.T) {
    // Setup: Model with existing user-defined entries
    model := &yappgen.Model{
        Cutouts: []map[string]any{
            {"face": "back", "shape": "circle"},  // User entry
        },
        CompositeProvenance: composite.NewProvenance(),
    }
    
    // Execute LCD module
    // ...
    
    // Assert: User entry preserved, LCD entry appended
    assert.Len(t, model.Cutouts, 2)
    assert.Equal(t, "back", model.Cutouts[0]["face"])   // User entry first
    assert.Equal(t, "front", model.Cutouts[1]["face"])  // LCD entry appended
}

// TestPostProcess_InvalidInput tests error handling
func TestPostProcess_InvalidInput(t *testing.T) {
    // Test cases:
    // - Missing display section
    // - Invalid position array
    // - Missing required mounting fields
    // ... etc
}
```

### Integration Tests

```yaml
# examples/composite/test-lcd-basic.yaml
features:
  lcd:
    - display:
        face: front
        position: [50, 30]
        size: [80, 40]
      mounting:
        type: pcb_stands
        positions: [[45, 25], [45, 55]]
        standoff_height: 5
        diameter: 3
```

```bash
# Test: Generate SCAD and verify arrays
go run ./cmd/yappctl generate --input examples/composite/test-lcd-basic.yaml --scad-out /tmp/test.scad

# Verify output contains:
# - cutoutsFront with LCD display cutout
# - pcbStands with 2 mounting holes
```

---

## Migration Path to Approach 1

### When to Migrate

**Triggers:**
1. Need for expression generation in composite modules
2. Complexity of direct Model manipulation becomes unwieldy
3. Number of composite modules exceeds ~5
4. Coupling to Model internals causes maintenance issues

### Migration Strategy

**Approach 3 → Approach 1 Refactor:**

**Step 1: Create DSL Transform Layer**
- Add `pkg/composite/transform.go` with Transform interface
- Composite modules implement Transform() instead of PostProcess()

**Step 2: Add Pre-Processing Hook**
- Move composite processing from BuildModel to CLI layer
- Process before resolver.ResolveResult()

**Step 3: Migrate Modules One-by-One**
- Refactor LCD module: PostProcess → Transform
- Keep old PostProcess temporarily for compatibility
- Test, verify, remove old code

**Step 4: Remove Post-Processing Hook**
- Remove processCompositeModules() from BuildModel
- Clean up Model.CompositeProvenance (move to DSL provenance)

**Estimated Migration Effort:** 1-2 weeks

**Design Principle:** Keep Approach 3 interfaces **close to Approach 1** to ease migration.

**Interface Similarity:**

**Approach 3:**
```go
PostProcess(ctx, model, moduleData, resolved) error
```

**Approach 1:**
```go
Transform(ctx, moduleData, resolved) ([]DSLEntry, error)
```

**Similar parameters, different return type.** Migration is straightforward.

---

## Implementation Phases

### Phase 1: Core Infrastructure (Week 1, Days 1-3)

**Deliverables:**
- `pkg/composite/module.go` - CompositeModule interface
- `pkg/composite/registry.go` - Registry implementation
- `pkg/composite/processor.go` - processCompositeModules()
- `pkg/composite/helpers.go` - Type conversion helpers
- `pkg/composite/provenance.go` - Provenance tracking
- Integration hook in `pkg/yappgen/model.go` (3 lines)

**Tests:**
- Registry tests (register, lookup, duplicate detection)
- Helper tests (type conversions, map extraction)
- Provenance tests (tracking, error formatting)

**Estimated Lines of Code:** ~200-300 lines

---

### Phase 2: Tooling and Documentation (Week 1, Days 4-5)

**Deliverables:**
- `cmd/schemagen docs` command - Auto-generate array format docs
- `pkg/composite/builders.go` - Auto-generated array construction helpers
- `pkg/docs/composite-array-formats.md` - Generated reference docs
- `pkg/docs/tutorials/composite-module-authoring-guide.md` - Author guide

**Code Generation Template:**
```go
// For each module schema, generate:
func Create{{.ModuleName}}Entry(required params, opts *{{.ModuleName}}Options) map[string]any {
    entry := map[string]any{
        // Required fields from schema
        {{range .RequiredFields}}
        "{{.Name}}": {{.Name}},
        {{end}}
    }
    
    // Optional fields from opts
    if opts != nil {
        {{range .OptionalFields}}
        if opts.{{.GoName}} != nil {
            entry["{{.Name}}"] = *opts.{{.GoName}}
        }
        {{end}}
    }
    
    return entry
}
```

**Estimated Lines of Code:** ~150-200 lines (template + generator)

---

### Phase 3: First Composite Module - LCD (Week 2, Days 1-2)

**Deliverables:**
- `pkg/composite/modules/lcd/module.go` - Complete LCD implementation
- `pkg/composite/modules/lcd/module_test.go` - Comprehensive tests
- `examples/composite/test-lcd-basic.yaml` - Basic example
- `examples/composite/test-lcd-complex.yaml` - Complex example
- `pkg/composite/modules_gen.go` - Generated registration (if using schemagen)

**LCD Module Features:**
- Display cutout generation (front, back, or lid face)
- Mounting hole patterns (rectangle, corners, custom)
- Configuration validation
- Provenance tracking
- Error messages

**Estimated Lines of Code:** ~150-200 lines (module + tests)

---

### Phase 4: Documentation and Examples (Week 2, Days 3-5)

**Deliverables:**
- Complete composite module authoring guide
- Multiple example YAML files
- Integration with existing yappctl help system
- Tutorial walkthrough
- Troubleshooting guide

**Examples:**
- `examples/composite/lcd-simple.yaml` - Minimal LCD config
- `examples/composite/lcd-with-buttons.yaml` - LCD + buttons
- `examples/composite/multiple-composites.yaml` - LCD + sensor
- `examples/composite/user-and-composite.yaml` - User cutouts + LCD

---

## API Surface Summary

### Public Interfaces

**Module Authors Use:**
```go
composite.CompositeModule         // Interface to implement
composite.Register()              // Register module
composite.ToFloat(), ToString()   // Type helpers
composite.GetFloat(), GetString() // Map helpers
composite.CreateCutoutEntry()     // Array builders (auto-generated)
composite.CreatePcbStandEntry()   // Array builders (auto-generated)
```

**Internal (Framework Uses):**
```go
composite.ProcessCompositeModules() // Called from BuildModel
composite.Provenance                // Error attribution
composite.Registry                  // Module management
```

**Generated Artifacts:**
```
pkg/composite/builders.go           # Array construction helpers
pkg/docs/composite-array-formats.md # Array format reference
pkg/composite/modules_gen.go        # Module registration (optional)
```

---

## Error Handling Design

### Error Attribution

**Without provenance:**
```
Error: features.cutouts[2].shape: invalid value "oval"
```

**With provenance:**
```
Error: features.cutouts[2].shape: invalid value "oval"
       (generated by composite module 'lcd')
```

### Error Categories

**1. User Configuration Errors:**
```
Error: composite module 'lcd': invalid display configuration
       lcd[0].display.position: expected 2 elements, got 1
```

**2. Module Implementation Errors:**
```
Error: composite module 'lcd': failed to generate cutout
       caused by: unknown face "top" (allowed: front, back, left, right, lid, base)
```

**3. Validation Errors (Post-Processing):**
```
Error: features.cutouts[2].shape: invalid value "oval"
       (generated by composite module 'lcd')
       Hint: lcd module may have a bug
```

---

## Decision: Approach 3 as Starting Point

### Rationale

**Based on 6 rounds of debate, Approach 3 is selected because:**

1. **Simpler Implementation**
   - ~200-300 lines core infrastructure
   - 1-2 weeks total implementation
   - Less architectural complexity

2. **Faster Time to Value**
   - Can ship first LCD module in 2-3 weeks
   - Validate approach with real usage
   - Iterate based on feedback

3. **Good DX with Minimal Investment**
   - 1 day tooling investment
   - Auto-generated docs and helpers
   - Clear example modules

4. **Lower Risk**
   - Minimal changes to existing code (3 lines in BuildModel)
   - No impact on resolver or regular modules
   - Easy to test and debug

5. **Clear Migration Path**
   - If expression generation needed → migrate to Approach 1
   - Interfaces designed for easy migration
   - Estimated migration: 1-2 weeks

### Approach 1 Migration Criteria

**Migrate when:**
- ✅ Need to generate expressions in composite modules
- ✅ Have 5+ composite modules (complexity justifies abstraction)
- ✅ Tight coupling to Model becomes maintenance burden
- ✅ Users request DSL-level composition features

**Until then:** Approach 3 is sufficient.

---

## Design Principles

### 1. Separation of Concerns

**Composite modules are separate from regular modules:**
- Different package: `pkg/composite/` vs `pkg/yappgen/modules/`
- Different interface: `CompositeModule` vs `registry.FeatureModule`
- Different processing: Post-processing vs collection

### 2. Minimal Disruption

**Changes to existing code:**
- `pkg/yappgen/model.go`: +3 lines (processCompositeModules call)
- `pkg/yappgen/model.go`: +1 field (CompositeProvenance)
- **Total: 4 lines changed**

**No changes to:**
- Resolver
- Existing modules
- Emission logic
- Validation logic

### 3. Explicit Over Implicit

**Composite module behavior is explicit:**
- Clear PostProcess function shows what arrays are modified
- Provenance tracking shows what was generated
- Error messages clearly attribute source

### 4. Pragmatic Over Perfect

**Start simple, add complexity only when needed:**
- No expression generation (Approach 1 feature, deferred)
- No overlap detection (future work)
- No cross-module dependencies (keep independent)
- **Solve the 80% case**, leave 20% for future iterations

### 5. Design for Migration

**Interfaces similar to Approach 1:**
- Easy to refactor PostProcess → Transform
- Similar parameters and structure
- Migration path is clear and low-risk

---

## Success Metrics

### Initial Success (Week 3)

- ✅ LCD module implemented and working
- ✅ Documentation complete
- ✅ Example YAML files generate correct SCAD
- ✅ 2-3 users test and provide feedback

### Migration Success (if needed)

- ✅ Clear pain points identified (expression generation needed)
- ✅ Migration completed in 1-2 weeks
- ✅ No breaking changes to existing composite modules
- ✅ Approach 1 benefits realized (expression support)

### Long-term Success (6 months)

- ✅ 3-5 composite modules in production
- ✅ Users report improved authoring experience
- ✅ No major bugs or architectural issues
- ✅ Clear decision on Approach 3 vs Approach 1 based on real usage

---

## Risks and Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|-----------|
| Tight coupling to Model internals | Medium | Medium | Design for migration to Approach 1, keep interfaces clean |
| Expression generation needed | Medium | High | Clear migration path to Approach 1 (~1-2 weeks) |
| Documentation falls out of sync | Low | Medium | Auto-generate docs from schemas |
| Only 2 composite modules | High | Low | That's fine - Approach 3 is optimized for 2-5 modules |
| Users request features Approach 3 can't support | Medium | Medium | Migrate to Approach 1 when justified |
| Helper functions incomplete | Low | Medium | Generate from schemas, extend as needed |
| Validation gaps | Low | High | Reuse existing validators, add provenance tracking |

---

## Alternatives Considered

### Why Not Approach 1 (Pre-Processing DSL Transform)?

**Pros:**
- Cleaner architecture
- Expression generation support
- Better separation

**Cons:**
- More complex (DSL merging logic)
- Longer implementation (3-4 weeks)
- Over-engineered for 2-3 modules

**Decision:** Start with Approach 3, migrate to Approach 1 if expression generation becomes necessary.

### Why Not Approach 2 (DSL Generation & Merge)?

Similar to Approach 1 but more complex merging. No clear advantage over Approach 1.

### Why Not Approach 4 (Hybrid)?

Too complex - runs pipeline twice. Over-engineered.

### Why Not Approach 5 (DSL Macros)?

Requires custom YAML parser. Too much infrastructure for limited benefit.

---

## Open Questions

### Resolved by Debates

- ✅ Which approach? **Approach 3 with migration path to Approach 1**
- ✅ Integration point? **Post-processing hook in BuildModel**
- ✅ Conflict strategy? **APPEND with deterministic order**
- ✅ Validation timing? **After merge, with provenance tracking**
- ✅ Developer experience? **Good with 1 day tooling investment**

### Remaining Questions

**1. Should we auto-generate composite module registration?**
- Option A: Manual registration (simple, 3 lines per module)
- Option B: Auto-generate like regular modules (consistent)
- **Lean toward:** Manual initially, auto-generate if we get 5+ modules

**2. Should helpers validate or trust callers?**
- Option A: Helpers validate (e.g., CreateCutoutEntry checks face enum)
- Option B: Helpers trust (validation happens in resolver anyway)
- **Lean toward:** Helpers validate for immediate feedback

**3. Should we support "disable" flags?**
```yaml
lcd:
  - display: {...}
    mounting: {...}
    disable_cutout: true  # Don't generate cutout
```
- **Lean toward:** Not initially, add if requested

---

## Next Steps

### Immediate (This Week)

1. Review and approve this design doc
2. Create implementation ticket/tasks
3. Set up pkg/composite/ package structure

### Implementation (Weeks 2-3)

1. Build core infrastructure (Phase 1)
2. Build tooling and docs (Phase 2)
3. Implement LCD module (Phase 3)
4. Write tutorials and examples (Phase 4)

### Follow-up (Week 4+)

1. User testing and feedback
2. Additional composite modules (if needed)
3. Evaluate migration to Approach 1 (if needed)

---

## References

### Debate Rounds

- [Round 1: Foundation](../debate/debate-round-01-foundation-should-we-do-this.md) - Should we do this?
- [Round 2: Architecture](../debate/debate-round-02-architecture-which-approach.md) - Which approach?
- [Round 3: Pipeline](../debate/debate-round-03-pipeline-integration.md) - Where to hook in?
- [Round 4: Conflicts](../debate/debate-round-04-conflict-resolution.md) - Conflict resolution
- [Round 5: Validation](../debate/debate-round-05-validation-timing.md) - Validation strategy
- [Round 6: Developer Experience](../debate/debate-round-06-developer-experience.md) - Authoring UX

### Analysis Documents

- [Module System Analysis](../analysis/01-module-system-analysis.md) - Current system analysis
- [Approaches Brainstorm](./01-composite-module-approaches-brainstorm.md) - 5 approaches explored

### Related Documentation

- Module Authoring Guide: `pkg/docs/tutorials/yapp-module-authoring-guide.md`
- Module System Guide: `ttmp/YAPP-PUSH-BUTTONS-001-.../playbook/module-system-implementation-guide.md`

