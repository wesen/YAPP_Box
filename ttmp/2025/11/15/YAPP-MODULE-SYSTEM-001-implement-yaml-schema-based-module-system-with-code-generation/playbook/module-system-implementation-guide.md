---
Title: Module System Implementation Guide - MVP
Ticket: YAPP-MODULE-SYSTEM-001
Status: active
Topics:
    - yapp
    - dsl
    - implementation
    - guide
DocType: playbook
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../reference/debate-round-01-validation-parsing-flow.md
    - Path: ../reference/debate-round-02-registration-mechanics.md
    - Path: ../reference/debate-round-03-prototype-implementation.md
    - Path: ../design/feature-module-registry.md
ExternalSources: []
Summary: Complete implementation guide for building the YAPP DSL module system - MVP scope, written for newcomers
LastUpdated: 2025-11-15
---

# YAPP DSL Module System - Implementation Guide (MVP)

## For the Complete Newcomer

This guide assumes you're new to this project, possibly new to Go, and definitely new to DSL (Domain-Specific Language) design. We'll explain everything from scratch.

## What is YAPP?

**YAPP** (Yet Another Parametric Projectbox) is an OpenSCAD library for generating custom 3D-printable enclosures for electronics projects. Think of it like a smart template that generates box designs based on parameters.

**Example:** "I have a 65mm × 45mm circuit board with 4 mounting holes and need cutouts for a USB port and power jack" → YAPP generates the 3D model.

## What Problem Are We Solving?

**Current situation:**
- Users write OpenSCAD code directly (programming language for 3D models)
- Example: `pcbStands = [[3, 3, 6, 2.5], [62, 3, 6, 2.5], ...];` (confusing positional arrays)
- Hard to remember what each number means
- Easy to make mistakes

**What we're building:**
- A YAML-based DSL (human-friendly text format)
- Users write readable configuration files
- Example:
  ```yaml
  pcb_stands:
    - x: 3
      y: 3
      diameter: 6
      pin_diameter: 2.5
  ```
- Much clearer!

**The Module System:** Each type of feature (pcb_stands, push_buttons, cutouts) is a "module" that knows how to:
1. Validate the YAML input (check for errors)
2. Convert YAML to OpenSCAD arrays (translate to the old format)

## Status Snapshot (2025-11-16)

- `pkg/registry` and `cmd/schemagen` scaffolding are live; `schemagen validate`/`discover` run against the repo.
- The first production schema (`pkg/yappgen/modules/pushbuttons/schema.yaml`) now exists, including nested cap/lid/switch definitions and CLI-validated tests.
- `schemagen` emits contextual error snippets with hints, so schema authors get actionable feedback before code generation.
- `schemagen discover` currently lists/validates schemas and stops at the “code generation coming soon” placeholder—implement code emission next.

## Core Concepts You Need to Understand

### 1. Schema

A **schema** is a definition of what data is valid. Like a form with rules.

**Example schema concept:**
```
Field "x" must be:
  - A number
  - Required (can't be missing)
  - Minimum value: 0
  - Description: "X coordinate on PCB in mm"
```

We store schemas as **YAML files** (one per module).

### 2. Module

A **module** is a self-contained feature handler. Each module has:
- **Schema** (what fields are valid)
- **Builder** (converts YAML → OpenSCAD arrays)

**Example modules:**
- `pcb_stands` - mounting points for circuit boards
- `push_buttons` - tactile switch extenders
- `cutouts` - holes in the enclosure walls

### 3. Registry

The **registry** is a central list of all modules. Think of it like a phone book:
- Key: module path (e.g., "features.push_buttons")
- Value: module object (schema + builder)

### 4. Code Generation

**Code generation** means writing a program that writes code. We'll build a tool called `schemagen` that:
- Reads YAML schema files
- Generates Go code (structs, tests, registration)

This automates boring work and prevents mistakes.

### 5. Two-Phase Validation

Validation happens in **two phases** because of expressions:

**Phase 1: Structure validation** (before expression resolution)
- Check: Does field exist? Is it the right type?
- Example: "field 'x' must be present and must be a number"

**Phase 2: Constraint validation** (after expression resolution)
- Check: Is the number in the right range?
- Example: "field 'x' must be ≥ 0"

**Why two phases?**
```yaml
vars:
  base_x: 5
push_buttons:
  - x: base_x + 10  # This is an expression, not yet a number
```

We can't check `x ≥ 0` until we resolve `base_x + 10` to `15`.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     User writes YAML                        │
│  features:                                                  │
│    push_buttons:                                            │
│      - x: 10                                                │
│        y: 10                                                │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│              Resolver (pkg/resolver)                        │
│  1. Load module schemas from registry                       │
│  2. Validate structure (Phase 1)                            │
│  3. Resolve expressions (vars.foo → actual values)          │
│  4. Validate constraints (Phase 2)                          │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│              Builder (pkg/yappgen)                          │
│  1. Get module from registry                                │
│  2. Call module's Build function                            │
│  3. Convert YAML maps → OpenSCAD arrays                     │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│              OpenSCAD Code                                  │
│  pushButtons = [                                            │
│    [10, 10, 6, 6, 1, ...],  // Positional array             │
│  ];                                                         │
└─────────────────────────────────────────────────────────────┘
```

## What You Need to Build

### Part 1: Registry Package (`pkg/registry`)

**Purpose:** Central registration point for all modules. Breaks circular dependencies between packages.

**What to create:**

**File: `pkg/registry/schema.go`**
```
Define interfaces (contracts):

interface ModuleSchema:
  - Name() returns string (e.g., "push_buttons")
  - Path() returns string (e.g., "features.push_buttons")
  - ValidateStructure(path, data) returns error
  - ValidateConstraints(path, data) returns error
  - Description() returns string
  - Fields() returns list of FieldSpec

struct FieldSpec:
  - Name (string)
  - Type (NumberField, StringField, ObjectField, ArrayField)
  - Required (boolean)
  - Description (string)
  - Example (any type)
  - Default (any type)
  - Min, Max (pointer to float64, can be nil)
  - Enum (list of strings)
  - Children (list of FieldSpec, for nested objects)

enum FieldType:
  - NumberField
  - StringField
  - ObjectField
  - ArrayField
```

**File: `pkg/registry/registry.go`**
```
interface FeatureModule:
  - Schema() returns ModuleSchema
  - Build(items) returns (SCAD arrays, error)

Global variables:
  - modules: map[string]FeatureModule
  - moduleOrder: list of strings

function Register(module):
  - Get path from module.Schema().Path()
  - If path already in modules, panic (duplicate registration)
  - Add module to modules map
  - Append path to moduleOrder list

function Get(path) returns (module, found):
  - Look up path in modules map
  - Return module and whether it was found

function All() returns list of modules:
  - Iterate through moduleOrder (NOT the map!)
  - For each path in order, get module from map
  - Return ordered list

function SetTestRegistry(testModules):
  - Replace global modules with test modules
  - Used only in unit tests
```

**Why separate package?**
- `pkg/resolver` needs module schemas (for validation)
- `pkg/yappgen` needs module schemas (for building)
- If registry was in yappgen, resolver would import yappgen
- If yappgen imports resolver (for expression evaluation), we get circular dependency
- Solution: Both import `pkg/registry` (no cycle!)

### Part 2: YAML Schema Format

**Purpose:** Define what fields each module accepts, in a human-readable format.

**File location:** `pkg/yappgen/modules/{module_name}/schema.yaml`

**Schema structure:**
```yaml
# Top-level metadata
module: push_buttons               # Module name (snake_case)
order: 100                         # Processing order (higher = later)
scad_array: pushButtons            # OpenSCAD array name (camelCase)
go_package: pushbuttons            # Go package name
description: |                     # Multi-line description
  What this module does...
  How it works...

# Field definitions
fields:
  field_name:
    type: number                   # number, string, object, array
    required: true                 # Must be present in YAML
    min: 0                        # For numbers: minimum value
    max: 100                      # For numbers: maximum value
    enum: [option1, option2]      # For strings: allowed values
    default: value                # Default if not provided
    desc: "What this field is"    # Human-readable description
    example: 10.0                 # Example value for docs
  
  nested_object:
    type: object
    required: true
    desc: "A nested object"
    fields:                       # Nested fields (recursive structure)
      sub_field:
        type: number
        required: true
        desc: "Field inside the object"

# Optional: Custom validation notes
custom_validation: |
  - If field A is set, field B is ignored
  - If field C is "circle", field D must be > 0

# Test cases
tests:
  - name: valid_example
    desc: "Description of what this tests"
    input:
      field_name: 10
      nested_object:
        sub_field: 5
    expect_valid: true
  
  - name: missing_required
    desc: "Should fail when required field is missing"
    input:
      nested_object:
        sub_field: 5
      # missing field_name
    expect_error: "missing required field 'field_name'"
```

**Example: Push Buttons Schema**

See the prototype in debate round 3 for a complete example. Key points:
- 158 lines for a complex module (push_buttons)
- Nested objects: `cap`, `lid`, `switch` (3 levels deep)
- Test cases catch regressions
- Self-documenting (descriptions + examples)

### Part 3: Code Generation Tool (`cmd/schemagen`)

**Purpose:** Read schema YAML files and generate Go code. Automates boring work.

**Command structure:**
```
schemagen validate <schema.yaml>    # Check if schema is valid
schemagen discover                  # Find all schemas and generate registry
```

**What it does:**

**Command: `schemagen validate`**
```
1. Read schema.yaml file
2. Parse YAML (check syntax)
3. Validate schema structure:
   - Required fields present (module, scad_array, fields)
   - Field types are valid (number, string, object, array)
   - Nested fields are properly structured
   - Test cases are well-formed
4. If errors, show line numbers and helpful messages:
   "Line 42: field 'unknown_type' has invalid type 'banana'"
   "Valid types: number, string, object, array"
5. Exit with code 0 (success) or 1 (failure)
```

**Command: `schemagen discover`**
```
1. Scan for schema files:
   - Find all files matching "pkg/yappgen/modules/*/schema.yaml"
   - Store list of schema paths

2. For each schema:
   - Validate schema (reuse validate logic)
   - Parse to extract: module name, package name, fields
   - If any schema invalid, abort with error

3. Generate Go code files:
   
   A. For each module, generate `schema_gen.go`:
      - Define struct types from fields
      - Example:
        type PushButtonDSL struct {
            X float64 `yaml:"x"`
            Y float64 `yaml:"y"`
            Cap CapDSL `yaml:"cap"`
        }
        
        type CapDSL struct {
            Length float64 `yaml:"length"`
            Width float64 `yaml:"width"`
        }
      
      - Generate ApplyDefaults() method:
        func (p *PushButtonDSL) ApplyDefaults() {
            if p.Shape == "" {
                p.Shape = "rounded_rect"  // from schema default
            }
        }
      
      - Generate CustomValidate() stub:
        func (p *PushButtonDSL) CustomValidate() error {
            // Module author fills this in for complex validation
            return nil
        }
   
   B. For each module, generate `schema_gen_test.go`:
      - For each test case in schema:
        func TestBasicButton(t *testing.T) {
            input := /* YAML from test case */
            _, err := Build(input)
            // Check expected result
        }
   
   C. Generate central registry file `pkg/yappgen/modules_gen.go`:
      // AUTO-GENERATED - DO NOT EDIT
      package yappgen
      
      import (
          "pkg/registry"
          "pkg/yappgen/modules/pushbuttons"
          "pkg/yappgen/modules/pcbstands"
          // ... imports for all discovered modules
      )
      
      func init() {
          registry.Register(pushbuttons.NewModule())
          registry.Register(pcbstands.NewModule())
          // ... register all modules
      }

4. Validate generated code:
   - Run go/format to check it's valid Go
   - Run go/parser to check it parses
   - If invalid, abort and show error

5. Write generated files atomically:
   - Write to temp files first (.tmp extension)
   - Only rename to real names if all succeeded
   - Prevents partial generation on failure
```

**Error handling (critical!):**

Bad error:
```
Error: yaml: line 42: could not find expected ':'
```

Good error:
```
Error in pkg/yappgen/modules/pushbuttons/schema.yaml:
  Line 42: YAML syntax error - missing colon after key
  
  41 |   protrusion:
  42 |     type number
         ^^^^^^^^^^ Did you mean 'type: number'?
  43 |     required: true

Run 'schemagen validate pushbuttons/schema.yaml' for details.
```

Use `gopkg.in/yaml.v3` for line number tracking.

### Part 4: Module Implementation (Example: Push Buttons)

Each module needs three files:

**File: `pkg/yappgen/modules/pushbuttons/schema.yaml`**
- Copy from prototype in debate round 3
- 158 lines defining all fields, nested objects, test cases

**File: `pkg/yappgen/modules/pushbuttons/schema_gen.go`**
- AUTO-GENERATED by schemagen
- Contains struct types (PushButtonDSL, CapDSL, LidDSL, SwitchDSL)
- Contains ApplyDefaults() method
- Contains CustomValidate() stub

**File: `pkg/yappgen/modules/pushbuttons/builder.go`**
- HAND-WRITTEN by module author
- Contains Build() function
- Contains NewModule() function

**Pseudocode for `builder.go`:**
```go
package pushbuttons

import "pkg/registry"
import "pkg/yappgen/scad"

// Build converts validated YAML to OpenSCAD arrays
function Build(items list of PushButtonDSL) returns (arrays, error):
    result = empty list
    
    for each button in items:
        // Call CustomValidate for business logic
        if error from button.CustomValidate():
            return error
        
        // Build positional array (order matters!)
        params = [
            button.X,
            button.Y,
            button.Cap.Length,
            button.Cap.Width,
            button.Cap.Radius,
            button.Lid.Protrusion,
            button.Switch.Height,
            button.Switch.Travel,
            button.Switch.PoleDiameter,
            // ... more parameters
        ]
        
        // Add flags (shape, coordinate system, etc.)
        if button.Shape == "circle":
            params.append(scad.Raw("yappCircle"))
        else if button.Shape == "rectangle":
            params.append(scad.Raw("yappRectangle"))
        
        result.append(params)
    
    return result, nil

// NewModule creates a FeatureModule for registration
function NewModule() returns FeatureModule:
    return &pushButtonModule{
        schema: GetSchema(),
        builder: Build,
    }

// GetSchema returns the parsed schema
function GetSchema() returns ModuleSchema:
    // Load and parse schema.yaml (embedded with go:embed)
    // Return as ModuleSchema implementation
```

**Where OpenSCAD array order comes from:**

Check `YAPPgenerator_v3.scad` and `YAPP_Template_v3.scad` in project root. These files document the positional parameter order for each YAPP feature.

Example from YAPP docs:
```
pushButtons = [
  [0] = x position
  [1] = y position
  [2] = cap length
  [3] = cap width
  [4] = cap radius
  [5] = lid protrusion
  ...
]
```

Your builder must output arrays in exactly this order.

### Part 5: Resolver Integration

**Purpose:** Validate user YAML before building OpenSCAD.

**File: `pkg/resolver/validation.go`**
```
function validateAgainstSchemas(document) returns error:
    features = document["features"]
    if features is nil:
        return nil  // No features to validate
    
    for each key, value in features:
        path = "features." + key
        module, found = registry.Get(path)
        if not found:
            continue  // Unknown module, skip
        
        schema = module.Schema()
        
        // Phase 1: Structure validation
        if error from schema.ValidateStructure(path, value):
            return error
    
    return nil

function validateConstraints(document) returns error:
    features = document["features"]
    if features is nil:
        return nil
    
    for each key, value in features:
        path = "features." + key
        module, found = registry.Get(path)
        if not found:
            continue
        
        schema = module.Schema()
        
        // Phase 2: Constraint validation
        if error from schema.ValidateConstraints(path, value):
            return error
    
    return nil
```

**Modify existing `Resolve()` function:**
```
function Resolve(document, options) returns (resolved, error):
    // NEW: Validate before expression resolution
    if error from validateAgainstSchemas(document):
        return error
    
    // EXISTING: Expression resolution loop
    for iteration = 0 to options.MaxIterations:
        resolve expressions in document
        if no changes:
            break
    
    // NEW: Validate after resolution
    if error from validateConstraints(document):
        return error
    
    return document, nil
```

### Part 6: Documentation Auto-Generation

**Purpose:** Generate help pages from schemas automatically.

**File: `pkg/docs/schema_help.go`**
```
function RenderSchemaHelp(schema) returns markdown:
    output = "# " + schema.Name() + "\n\n"
    output += schema.Description() + "\n\n"
    
    output += "## Fields\n\n"
    output += "| Field | Type | Required | Default | Description |\n"
    output += "|-------|------|----------|---------|-------------|\n"
    
    fields = schema.Fields()
    renderFields(fields, "", output)  // Recursive for nested
    
    return output

function renderFields(fields, prefix, output):
    for each field in fields:
        fullName = prefix + field.Name
        
        req = ""
        if field.Required:
            req = "✓"
        
        defaultVal = field.Default or "—"
        
        output += "| `" + fullName + "` | " 
        output += fieldTypeName(field.Type) + " | "
        output += req + " | "
        output += defaultVal + " | "
        output += field.Description + " |\n"
        
        // Recurse for nested objects
        if field.Type == ObjectField and field.Children not empty:
            renderFields(field.Children, fullName + ".", output)

function LoadModuleHelp():
    // Embedded schemas (go:embed)
    schemaFiles = embedded "modules/*/schema.yaml"
    
    for each schemaFile in schemaFiles:
        schema = parseSchema(schemaFile)
        markdown = RenderSchemaHelp(schema)
        helpSystem.AddHelpPage("module-" + schema.Name(), markdown)
```

**Integration:**
```
function main():
    // Load module help during CLI initialization
    LoadModuleHelp()
    
    // User can now run:
    // go run ./cmd/yappctl help push_buttons
```

### Part 7: Testing Strategy

**Unit tests for schemagen:**
```
TestValidateSchema_ValidSchema:
    - Create a valid schema YAML
    - Call validate
    - Expect no errors

TestValidateSchema_InvalidType:
    - Create schema with type: "banana"
    - Call validate
    - Expect error mentioning line number and valid types

TestDiscover_FindsSchemas:
    - Create temp directory with test schemas
    - Call discover
    - Expect all schemas found

TestGenerateRegistry_ValidOutput:
    - Call discover with test schemas
    - Check generated Go code compiles
    - Check imports are correct
```

**Unit tests for modules:**
```
TestPushButtonBuild_ValidInput:
    - Create valid PushButtonDSL struct
    - Call Build()
    - Expect SCAD array with correct parameter order

TestPushButtonValidate_MissingRequired:
    - Create PushButtonDSL missing required field
    - Call Validate()
    - Expect error

TestSchemaValidation_AllTestCases:
    - Load schema.yaml
    - For each test case:
        - Run validation
        - Check result matches expected
```

**Integration tests:**
```
TestFullPipeline:
    1. Write test YAML file
    2. Call Resolve()
    3. Call BuildModel()
    4. Call EmitSCAD()
    5. Check OpenSCAD output is valid
```

### Part 8: CI Integration

**Add to `.github/workflows/` or CI config:**
```yaml
name: Schema Validation

on: [push, pull_request]

jobs:
  validate:
    steps:
      - name: Check generated files are up-to-date
        run: |
          go generate ./...
          git diff --exit-code
        # Fails if someone forgot to run go generate
      
      - name: Validate all schemas
        run: |
          for schema in pkg/yappgen/modules/*/schema.yaml; do
            go run ./cmd/schemagen validate "$schema"
          done
      
      - name: Run tests
        run: go test ./...
```

**Add pre-commit hook (`.git/hooks/pre-commit`):**
```bash
#!/bin/bash
# Auto-run go generate before committing

go generate ./pkg/yappgen
git add pkg/yappgen/*_gen.go

# Fail if generated code changed but wasn't committed
git diff --cached --exit-code pkg/yappgen/*_gen.go || {
    echo "Generated files changed. Re-committing..."
    exit 1
}
```

## Implementation Order

Follow this sequence to minimize risk:

### Week 1: Core Infrastructure
1. Create `pkg/registry` package (interfaces + registration)
2. Write tests for registry (registration, lookup, ordering)
3. Create `cmd/schemagen` skeleton (CLI structure, command parsing)

### Week 2: Code Generation
1. Implement `schemagen validate` command
   - YAML parsing with line numbers
   - Schema structure validation
   - Error message formatting
2. Implement `schemagen discover` command
   - Filesystem scanning
   - Code generation templates
   - Generated code validation
3. Test with a minimal example schema

### Week 3: Push Buttons Conversion (Pilot)
1. Create `pkg/yappgen/modules/pushbuttons/schema.yaml`
   - Copy structure from debate round 3 prototype
2. Run `go generate` to create schema_gen.go
3. Write `builder.go` (Build function + NewModule)
4. Write tests
5. Fix any issues with generated code
6. Update schemagen if needed

### Week 4: Resolver Integration
1. Add `validateAgainstSchemas()` to resolver
2. Add `validateConstraints()` to resolver
3. Modify `Resolve()` to call both validation phases
4. Test with push_buttons YAML examples
5. Ensure error messages are helpful

### Week 5: Remaining Modules
1. Convert `pcb_stands` to YAML schema (simpler than push_buttons)
2. Convert `connectors` to YAML schema
3. Convert `snap_joins` to YAML schema
4. Convert `cutouts` to YAML schema
5. Run full test suite

### Week 6: Documentation & Polish
1. Implement schema → markdown rendering
2. Integrate with Glazed help system
3. Add error message → help link suggestions
4. Write module authoring guide
5. Update CI configuration
6. Write pre-commit hook template

## Key Files and Their Purposes

```
pkg/
  registry/
    schema.go           # ModuleSchema interface, FieldSpec
    registry.go         # Register(), Get(), All()
    registry_test.go    # Unit tests

  resolver/
    resolver.go         # Main Resolve() function (modify)
    validation.go       # NEW: validateAgainstSchemas(), validateConstraints()
    validation_test.go  # NEW: Validation tests
  
  yappgen/
    modules_gen.go      # GENERATED: Central registry, auto-registers modules
    
    modules/
      pushbuttons/
        schema.yaml           # HAND-WRITTEN: Field definitions, tests
        schema_gen.go         # GENERATED: Struct types, ApplyDefaults()
        schema_gen_test.go    # GENERATED: Tests from schema
        builder.go            # HAND-WRITTEN: Build(), NewModule()
        builder_test.go       # HAND-WRITTEN: Additional tests
      
      pcbstands/
        schema.yaml
        schema_gen.go
        builder.go
      
      # ... more modules
  
  docs/
    schema_help.go      # NEW: RenderSchemaHelp()
    docs.go             # MODIFY: Add LoadModuleHelp()

cmd/
  schemagen/
    main.go             # CLI entry point
    validate.go         # validate command
    discover.go         # discover command
    generate.go         # Code generation logic
    templates/          # Go code templates
  
  yappctl/
    main.go             # MODIFY: Call LoadModuleHelp()
```

## Common Pitfalls to Avoid

### 1. Circular Dependencies
**Problem:** Package A imports B, B imports A → compile error

**Solution:** Use the registry package as a shared intermediary
- resolver imports registry (not yappgen)
- yappgen imports registry (not resolver)
- modules import registry (not yappgen)

### 2. Zero Value Confusion
**Problem:** How do you know if user set field to 0 or didn't set it?

**Solution:** Use pointers for optional fields with non-zero defaults
```
// Wrong:
type Config struct {
    Angle float64  // Default should be 45, but user might set 0
}

// Right:
type Config struct {
    Angle *float64  // nil = not set, &0.0 = set to zero
}

func (c *Config) ApplyDefaults() {
    if c.Angle == nil {
        angle := 45.0
        c.Angle = &angle
    }
}
```

### 3. Module Ordering
**Problem:** Modules process in random order (map iteration)

**Solution:** Track registration order in a slice, iterate that
```
// Wrong:
for name, module := range modules {
    process(module)  // Random order!
}

// Right:
for _, name := range moduleOrder {
    module := modules[name]
    process(module)  // Deterministic order!
}
```

### 4. Partial Code Generation
**Problem:** Generator crashes halfway, leaves invalid .go files

**Solution:** Write to temp files, only rename on success
```
// Wrong:
write to schema_gen.go
crash!  // Now schema_gen.go is broken

// Right:
write to schema_gen.go.tmp
validate generated code
if valid:
    rename schema_gen.go.tmp to schema_gen.go
```

### 5. Bad Error Messages
**Problem:** "yaml: line 42: syntax error" (not helpful)

**Solution:** Show context and suggest fixes
```
Error in schema.yaml line 42:
  41 |   protrusion:
  42 |     type number
         ^^^^^^^^^^^
  43 |     required: true

Missing colon after 'type'. Did you mean 'type: number'?
```

## Testing Your Implementation

### Manual Testing Checklist

1. **Create a test schema:**
   ```bash
   mkdir -p pkg/yappgen/modules/testmod
   cat > pkg/yappgen/modules/testmod/schema.yaml << 'EOF'
   module: test_feature
   scad_array: testFeature
   go_package: testmod
   description: Test module
   fields:
     x:
       type: number
       required: true
   EOF
   ```

2. **Validate schema:**
   ```bash
   go run ./cmd/schemagen validate pkg/yappgen/modules/testmod/schema.yaml
   # Should pass with no errors
   ```

3. **Test invalid schema:**
   ```bash
   # Add invalid field type
   echo "    banana: true" >> pkg/yappgen/modules/testmod/schema.yaml
   go run ./cmd/schemagen validate pkg/yappgen/modules/testmod/schema.yaml
   # Should show helpful error
   ```

4. **Run discovery:**
   ```bash
   go run ./cmd/schemagen discover
   # Should find test module and generate modules_gen.go
   ```

5. **Check generated code compiles:**
   ```bash
   go build ./pkg/yappgen
   # Should compile without errors
   ```

6. **Test resolver integration:**
   ```yaml
   # Create test.yaml
   features:
     test_feature:
       - x: 10
   ```
   ```bash
   go run ./cmd/yappctl resolve --input test.yaml
   # Should validate successfully
   ```

7. **Test with invalid input:**
   ```yaml
   features:
     test_feature:
       - y: 10  # Wrong: should be 'x'
   ```
   ```bash
   go run ./cmd/yappctl resolve --input test.yaml
   # Should show: "missing required field 'x'"
   ```

## Success Criteria

You've successfully implemented the MVP when:

1. ✅ Push buttons schema exists (schema.yaml)
2. ✅ `schemagen validate` catches schema errors with helpful messages
3. ✅ `schemagen discover` generates valid Go code (modules_gen.go)
4. ✅ Generated structs compile and have ApplyDefaults()
5. ✅ Tests auto-generate from schema test cases
6. ✅ Resolver validates YAML against schemas (two phases)
7. ✅ Build process works (YAML → validated data → OpenSCAD arrays)
8. ✅ `go run ./cmd/yappctl help push_buttons` shows auto-generated docs
9. ✅ CI fails if someone forgets to run `go generate`
10. ✅ Error messages include line numbers and suggestions

## Resources for Learning

### Go Language
- Official tour: https://go.dev/tour/
- Effective Go: https://go.dev/doc/effective_go
- Go by Example: https://gobyexample.com/

### YAML
- YAML spec: https://yaml.org/spec/
- Go YAML library: https://pkg.go.dev/gopkg.in/yaml.v3

### Code Generation in Go
- `text/template`: https://pkg.go.dev/text/template
- `go/format`: https://pkg.go.dev/go/format
- `go/parser`: https://pkg.go.dev/go/parser

### OpenSCAD
- Read `YAPPgenerator_v3.scad` in project root
- YAPP documentation: Comments in `YAPP_Template_v3.scad`

### This Project
- Debate documents (in ../reference/) explain design decisions
- Feature registry design (in ../design/) explains architecture
- Existing modules (in pkg/yappgen/) show current patterns

## Getting Help

If you get stuck:

1. **Read the debates** - They explain WHY decisions were made
2. **Look at existing code** - `pkg/yappgen/features.go` shows current registry
3. **Check test files** - They show expected behavior
4. **Ask specific questions** - "Why two-phase validation?" vs. "How do I code?"

## Glossary

**DSL** - Domain-Specific Language. A mini-language designed for a specific task (our YAML format).

**OpenSCAD** - Programming language for 3D models. YAPP generates this.

**Schema** - Definition of valid data structure. Like a form with rules.

**Registry** - Central list mapping module names to module implementations.

**Code Generation** - Writing a program that writes code (schemagen).

**Module** - Self-contained feature handler (schema + builder).

**Builder** - Function that converts YAML → OpenSCAD arrays.

**Resolution** - Process of evaluating expressions (vars.x + 10 → actual number).

**Two-Phase Validation** - Validate structure before resolution, constraints after.

**Circular Dependency** - Package A imports B, B imports A (causes compile error).

**Registry Pattern** - Central registration point that breaks circular dependencies.

---

**Remember:** This is an MVP. Don't add features not mentioned here. Build the simplest thing that works, test it thoroughly, and iterate based on real usage.

Good luck! 🚀

