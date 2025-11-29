---
Title: Module System Analysis - Multi-Array Module Support
Ticket: YAPP-MODULE-PIPELINE-001
Status: active
Topics:
  - dsl
  - modules
  - codegen
DocType: analysis
Intent: long-term
Owners: []
RelatedFiles:
  - Path: /home/manuel/workspaces/2025-11-28/yapp-film-developer/YAPP_Box/pkg/registry/schema.go
    Note: Defines FeatureModule interface and ArrayDecl type
  - Path: /home/manuel/workspaces/2025-11-28/yapp-film-developer/YAPP_Box/pkg/yappgen/features.go
    Note: Wires up modules using arrayFeatureModule and multiArrayFeatureModule helpers
  - Path: /home/manuel/workspaces/2025-11-28/yapp-film-developer/YAPP_Box/pkg/yappgen/modules_gen.go
    Note: Generated registry that registers all modules
  - Path: /home/manuel/workspaces/2025-11-28/yapp-film-developer/YAPP_Box/pkg/yappgen/modules/cutouts/registry.go
    Note: Example of multi-array module (returns multiple ArrayDecl for different faces)
  - Path: /home/manuel/workspaces/2025-11-28/yapp-film-developer/YAPP_Box/pkg/yappgen/modules/pcbstands/registry.go
    Note: Example of single-array module (returns one ArrayDecl)
  - Path: /home/manuel/workspaces/2025-11-28/yapp-film-developer/YAPP_Box/pkg/registry/registry.go
    Note: Global module registry implementation
Summary: Analysis of current module system architecture and identification of gaps for composite modules (e.g., LCD module that adds cutouts + mounting holes)
LastUpdated: 2025-11-28
---

# Module System Analysis - Multi-Array Module Support

## Overview

This document analyzes the current YAPP DSL module system to understand how modules work, identify limitations, and determine what's needed to support **composite modules** that can add entries to multiple different array types (e.g., an LCD module that adds a display cutout to `cutoutsFront` AND mounting holes to `pcbStands` or `boxMounts`).

## Current Module System Architecture

### Core Components

#### 1. Module Interface (`pkg/registry/schema.go`)

```60:66:pkg/registry/schema.go
// ArrayDecl represents a single OpenSCAD array declaration.
type ArrayDecl struct {
	// Name is the identifier written in the SCAD output (e.g., "pcbStands").
	Name string
	// Rows contains the positional parameter rows that make up the array.
	Rows [][]any
}
```

```51:58:pkg/registry/schema.go
// FeatureModule ties together a schema with the logic that converts validated
// DSL entries into OpenSCAD parameter arrays.
type FeatureModule interface {
	// Schema returns the module's schema metadata.
	Schema() ModuleSchema
	// Build consumes validated items and produces OpenSCAD parameter arrays.
	Build(items []map[string]any) ([]ArrayDecl, error)
}
```

**Key Points:**
- `Build()` returns `[]ArrayDecl` - can return multiple arrays
- Each `ArrayDecl` has a `Name` (SCAD array name) and `Rows` (parameter rows)
- Modules are registered by their schema `Path()` (e.g., `"features.cutouts"`)

#### 2. Module Registration (`pkg/yappgen/modules_gen.go`)

```16:24:pkg/yappgen/modules_gen.go
func init() {
	registry.Register(boxmounts.NewModule())
	registry.Register(connectors.NewModule())
	registry.Register(cutouts.NewModule())
	registry.Register(lighttubes.NewModule())
	registry.Register(pcbstands.NewModule())
	registry.Register(pushbuttons.NewModule())
	registry.Register(snapjoins.NewModule())
}
```

**Key Points:**
- Generated file (`modules_gen.go`) registers all modules
- Modules register themselves via `init()` function
- Registration happens at package import time

#### 3. Feature Module Wiring (`pkg/yappgen/features.go`)

```27:53:pkg/yappgen/features.go
var featureModules = []FeatureModule{
	newArrayFeatureModule("pcb_stands", "pcbStands",
		func(m *Model) *[]map[string]any { return &m.PcbStands },
		pcbstands.NewModule().Build, nil, nil),
	newArrayFeatureModule("connectors", "connectors",
		func(m *Model) *[]map[string]any { return &m.Connectors },
		connectors.NewModule().Build, nil, nil),
	newArrayFeatureModule("box_mounts", "boxMounts",
		func(m *Model) *[]map[string]any { return &m.BoxMounts },
		boxmounts.NewModule().Build, nil, nil),
	newArrayFeatureModule("push_buttons", "pushButtons",
		func(m *Model) *[]map[string]any { return &m.PushButtons },
		pushbuttons.NewModule().Build,
		func(m *Model, items []map[string]any) {
			m.PrintSwitchExtenders = len(items) > 0
		},
		pushButtonsHeader()),
	newArrayFeatureModule("snap_joins", "snapJoins",
		func(m *Model) *[]map[string]any { return &m.SnapJoins },
		snapjoins.NewModule().Build, nil, nil),
	newArrayFeatureModule("light_tubes", "lightTubes",
		func(m *Model) *[]map[string]any { return &m.LightTubes },
		lighttubes.NewModule().Build, nil, nil),
	newMultiArrayFeatureModule("cutouts",
		func(m *Model) *[]map[string]any { return &m.Cutouts },
		cutouts.NewModule().Build),
}
```

**Key Points:**
- Two helper types: `arrayFeatureModule` (single array) and `multiArrayFeatureModule` (multiple arrays)
- Each module is tied to ONE DSL key (e.g., `"pcb_stands"`)
- Each module reads from ONE field in the `Model` struct
- `arrayFeatureModule` enforces that exactly one `ArrayDecl` is returned
- `multiArrayFeatureModule` allows multiple `ArrayDecl` but they're typically related (like cutouts per face)

#### 4. Module Collection and Emission

```55:71:pkg/yappgen/features.go
func collectFeatureModules(resolved map[string]any, features map[string]any, model *Model) error {
	for _, module := range featureModules {
		if err := module.Collect(resolved, features, model); err != nil {
			return errors.Wrapf(err, "feature %s", module.Name())
		}
	}
	return nil
}

func emitFeatureModules(ctx context.Context, model *Model, b *strings.Builder) error {
	for _, module := range featureModules {
		if err := module.Emit(ctx, model, b); err != nil {
			return errors.Wrapf(err, "feature %s", module.Name())
		}
	}
	return nil
}
```

**Key Points:**
- `Collect()` reads DSL data and populates `Model` fields
- `Emit()` calls module's `Build()` and writes SCAD arrays
- Modules are processed sequentially in registration order

### Current Module Examples

#### Single-Array Module: `pcb_stands`

```27:42:pkg/yappgen/modules/pcbstands/registry.go
func (m *module) Build(items []map[string]any) ([]registry.ArrayDecl, error) {
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}

	rows, err := Build(typed)
	if err != nil {
		return nil, err
	}

	return []registry.ArrayDecl{{
		Name: "pcbStands",
		Rows: rows,
	}}, nil
}
```

**Characteristics:**
- Returns exactly one `ArrayDecl`
- Array name: `"pcbStands"`
- DSL key: `"features.pcb_stands"`
- Model field: `Model.PcbStands`

#### Multi-Array Module: `cutouts`

```32:59:pkg/yappgen/modules/cutouts/registry.go
func (m *module) Build(items []map[string]any) ([]registry.ArrayDecl, error) {
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}

	byFace, err := Build(typed)
	if err != nil {
		return nil, err
	}

	var decls []registry.ArrayDecl
	for name, rows := range byFace {
		if len(rows) == 0 {
			continue
		}
		decls = append(decls, registry.ArrayDecl{
			Name: name,
			Rows: rows,
		})
	}

	sort.Slice(decls, func(i, j int) bool {
		return decls[i].Name < decls[j].Name
	})

	return decls, nil
}
```

**Characteristics:**
- Returns multiple `ArrayDecl` (one per face)
- Array names: `"cutoutsFront"`, `"cutoutsBack"`, `"cutoutsLeft"`, `"cutoutsRight"`, `"cutoutsLid"`, `"cutoutsBase"`
- DSL key: `"features.cutouts"`
- Model field: `Model.Cutouts`
- **All arrays are the same type** (cutouts), just distributed by face

## Current Limitations

### 1. One DSL Key Per Module

**Current Behavior:**
- Each module is registered with ONE schema path (e.g., `"features.pcb_stands"`)
- Each module reads from ONE DSL key (e.g., `features.pcb_stands`)
- Each module writes to ONE `Model` field

**Implication:**
- Cannot have a module that reads from `features.lcd` and adds entries to multiple different array types
- Would need separate DSL entries for each component (e.g., `features.lcd_cutout` and `features.lcd_mounts`)

### 2. One Model Field Per Module

**Current Behavior:**
- `arrayFeatureModule` and `multiArrayFeatureModule` both take a field getter: `func(*Model) *[]map[string]any`
- This field stores the raw DSL data for that module
- The module's `Build()` function receives only this one field's data

**Implication:**
- A composite module cannot access data from multiple DSL keys
- Cannot combine data from `features.lcd` with `features.pcb` or other features

### 3. Array Type Restriction

**Current Behavior:**
- `cutouts` module returns multiple arrays, but they're all cutout arrays (just different faces)
- `arrayFeatureModule` enforces that exactly one array is returned
- `multiArrayFeatureModule` allows multiple arrays but they're typically related

**Implication:**
- Cannot have a module that adds to `cutoutsFront` AND `pcbStands` simultaneously
- Would need to manually add entries to both arrays in the DSL

### 4. No Cross-Module Dependencies

**Current Behavior:**
- Modules are processed independently
- No way for one module to reference or modify another module's output
- No way to coordinate between modules

**Implication:**
- Cannot have an LCD module that:
  - Adds a display cutout to `cutoutsFront`
  - Adds mounting holes to `pcbStands` (for mounting the LCD PCB)
  - Ensures cutout and mounting holes are aligned

## Use Case: LCD Module

### Desired Behavior

An LCD module should allow users to specify:

```yaml
features:
  lcd:
    - display:
        face: front
        position: [x, y]
        size: [width, height]
        cutout_shape: rectangle
      mounting:
        type: pcb_stands  # or box_mounts
        positions: [[x1, y1], [x2, y2], ...]
        standoff_height: 5
        diameter: 3
```

**Expected Output:**
- Add display cutout to `cutoutsFront` array
- Add mounting holes to `pcbStands` array (or `boxMounts`)
- Ensure cutout and mounting positions are aligned

### Current Workaround

Users must manually specify:

```yaml
features:
  cutouts:
    - face: front
      from_face_left: <calculated>
      from_face_bottom: <calculated>
      width: <lcd_width>
      length: <lcd_height>
      shape: rectangle
  pcb_stands:
    - x: <mount_x1>
      y: <mount_y1>
      height: 5
      diameter: 3
    - x: <mount_x2>
      y: <mount_y2>
      height: 5
      diameter: 3
    # ... more mounts
```

**Problems:**
- Manual calculation of positions
- No validation that cutout and mounts align
- Repetitive configuration
- Error-prone

## Available YAPP Arrays

From analysis of `YAPPgenerator_v3.scad` and existing modules:

1. **pcbStands** - PCB mounting standoffs
2. **connectors** - Connector mounting points
3. **boxMounts** - Box mounting points
4. **snapJoins** - Snap-fit joints
5. **lightTubes** - Light guide tubes
6. **pushButtons** - Push button extenders
7. **cutoutsFront/Back/Left/Right/Lid/Base** - Face cutouts

**Note:** A composite module might need to add entries to any combination of these arrays.

## Relevant Code Locations

### Core Module System

1. **`pkg/registry/schema.go`**
   - `FeatureModule` interface definition
   - `ArrayDecl` type definition
   - `ModuleSchema` interface

2. **`pkg/registry/registry.go`**
   - Global module registry
   - `Register()`, `Get()`, `All()` functions

3. **`pkg/yappgen/features.go`**
   - `FeatureModule` interface (local, wraps registry)
   - `arrayFeatureModule` helper
   - `multiArrayFeatureModule` helper
   - `collectFeatureModules()` function
   - `emitFeatureModules()` function
   - `featureModules` slice (wires up all modules)

4. **`pkg/yappgen/modules_gen.go`**
   - Generated module registration
   - Imports all module packages
   - Calls `registry.Register()` for each module

### Module Implementations

5. **`pkg/yappgen/modules/*/registry.go`**
   - Each module implements `registry.FeatureModule`
   - `Build()` function returns `[]ArrayDecl`
   - Schema implementation

6. **`pkg/yappgen/modules/*/module.go`**
   - Typed builder functions
   - Convert typed structs to SCAD parameter rows

7. **`pkg/yappgen/modules/*/schema.yaml`**
   - YAML schema definitions
   - Used by code generator to create Go types

### Code Generation

8. **`pkg/schemagen/discover.go`**
   - Discovers modules from `modules/*/schema.yaml`
   - Generates `modules_gen.go`

9. **`pkg/schemagen/templates/modules_gen.go.tmpl`**
   - Template for generated registry file

## Key Findings

### What Works Today

1. ✅ Modules can return multiple `ArrayDecl` (see `cutouts`)
2. ✅ Modules are registered automatically via code generation
3. ✅ Module system is extensible (add new module = add schema.yaml)
4. ✅ Modules can have complex schemas (nested objects, enums, validation)

### What Doesn't Work Today

1. ❌ Modules cannot add entries to **different array types** (e.g., cutouts + pcbStands)
2. ❌ Modules cannot read from **multiple DSL keys**
3. ❌ Modules cannot coordinate with **other modules**
4. ❌ No way to have **composite modules** that combine multiple features

### Architectural Constraints

1. **One-to-One Mapping:** One DSL key → One Model field → One module
2. **Sequential Processing:** Modules are processed independently, in order
3. **No Dependencies:** Modules cannot reference other modules' data or output
4. **No Merging:** Arrays from different modules are not merged (each module emits its own arrays)

## Next Steps

To support composite modules like LCD, we need to:

1. **Design a new module type** that can:
   - Read from one DSL key (e.g., `features.lcd`)
   - Output multiple `ArrayDecl` of **different types** (cutouts + pcbStands)
   - Optionally merge with existing arrays (append to existing `pcbStands` entries)

2. **Extend the module interface** to support:
   - Cross-array type output
   - Array merging/aggregation
   - Module dependencies

3. **Update the feature wiring** to:
   - Handle modules that output to multiple array types
   - Merge arrays from multiple modules
   - Validate array consistency

4. **Consider DSL design** for composite modules:
   - How to specify which arrays to populate
   - How to handle conflicts (multiple modules adding to same array)
   - How to ensure proper ordering

## Questions to Resolve

1. **Should composite modules merge with existing arrays?**
   - Option A: Append to existing arrays (e.g., LCD mounts added to existing `pcbStands`)
   - Option B: Create separate arrays (e.g., `lcdMounts` separate from `pcbStands`)

2. **How to handle array ordering?**
   - If LCD module adds to `pcbStands`, should it come before or after manual `pcb_stands` entries?
   - Should there be explicit ordering control?

3. **How to handle conflicts?**
   - What if LCD module and manual `pcb_stands` specify overlapping positions?
   - Should there be validation or warnings?

4. **Should composite modules be a separate type?**
   - New `compositeFeatureModule` type?
   - Or extend `multiArrayFeatureModule` to support different array types?

5. **How to handle module dependencies?**
   - Should LCD module be able to reference PCB dimensions?
   - Should it be able to reference other feature positions?

## Related Documentation

- Module authoring guide: `pkg/docs/tutorials/yapp-module-authoring-guide.md`
- Module system implementation: `ttmp/YAPP-PUSH-BUTTONS-001-implement-pushbuttons-dsl-parity/playbook/module-system-implementation-guide.md`
- ArrayDecl refactor: `ttmp/2025/11/17/YAPP-BUILDER-REFACTOR-001-implement-path-1-a-b-builder-contract-refactor-generate-decode-and-arraydecl-ir/`

