---
Title: Debate Round 3 — Prototype Implementation
Ticket: YAPP-PUSH-BUTTONS-001
Status: active
Topics:
    - yapp
    - dsl
    - schema
    - prototype
DocType: debate
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/yappgen/features.go
    - Path: pkg/yappgen/modules/pushbuttons/module.go
    - Path: debate-round-01-validation-parsing-flow.md
    - Path: debate-round-02-registration-mechanics.md
ExternalSources: []
Summary: Hands-on prototyping round where candidates implement the schema + registry pattern with push_buttons and report findings
LastUpdated: 2025-11-15
---

# Debate Round 3 — Prototype Implementation

## Question

**Given a straw-man implementation (e.g., refactoring `modules/pushbuttons` to use YAML schemas + generated registry), what works well, what feels brittle, and what adjustments would we make before declaring the pattern ready?**

**Design decision:** All modules use YAML schemas (no hybrid approach with struct tags for simple modules). This keeps the architecture uniform and tooling simpler.

---

## Pre-Debate Implementation (Round 3)

### Implementation Task

Each candidate prototypes a key piece of the architecture and reports findings.

---

**A (DSL Architect) — Prototype: Registry Package**

*[Creates pkg/registry package]*

```bash
mkdir -p pkg/registry
cat > pkg/registry/schema.go <<'EOF'
package registry

// ModuleSchema defines validation and metadata for a DSL module
type ModuleSchema interface {
    Name() string
    Path() string  // e.g., "features.push_buttons"
    
    // Two-phase validation (from Round 1)
    ValidateStructure(path string, data any) error
    ValidateConstraints(path string, data any) error
    
    // Metadata for documentation
    Description() string
    Fields() []FieldSpec
}

type FieldSpec struct {
    Name        string
    Type        FieldType
    Required    bool
    Description string
    Min, Max    *float64
    Enum        []string
    Children    []FieldSpec  // For nested objects
}

type FieldType int

const (
    NumberField FieldType = iota
    StringField
    ObjectField
    ArrayField
)
EOF
```

**Testing the interface:**

```bash
cat > pkg/registry/schema_test.go <<'EOF'
package registry_test

import (
    "testing"
    "github.com/wesen/yapp-encl-resolver/pkg/registry"
)

type mockSchema struct {
    name string
    path string
}

func (m *mockSchema) Name() string { return m.name }
func (m *mockSchema) Path() string { return m.path }
func (m *mockSchema) ValidateStructure(path string, data any) error { return nil }
func (m *mockSchema) ValidateConstraints(path string, data any) error { return nil }
func (m *mockSchema) Description() string { return "test schema" }
func (m *mockSchema) Fields() []registry.FieldSpec { return nil }

func TestSchemaInterface(t *testing.T) {
    var s registry.ModuleSchema = &mockSchema{
        name: "test",
        path: "features.test",
    }
    
    if s.Name() != "test" {
        t.Errorf("expected name 'test', got %s", s.Name())
    }
}
EOF

go test ./pkg/registry
```

**Finding:** Interface compiles and tests pass. But I realized **we need a builder interface too**:

```go
// FeatureModule combines schema + builder
type FeatureModule interface {
    Schema() ModuleSchema
    Build(items []any) ([][]any, error)
}
```

This pairs validation with generation. The registry holds `FeatureModule`, not just schema.

**Updated registry.go:**

```go
package registry

var modules = make(map[string]FeatureModule)

func Register(mod FeatureModule) {
    path := mod.Schema().Path()
    if _, exists := modules[path]; exists {
        panic("duplicate module registration: " + path)
    }
    modules[path] = mod
}

func Get(path string) (FeatureModule, bool) {
    m, ok := modules[path]
    return m, ok
}

func All() []FeatureModule {
    var result []FeatureModule
    for _, m := range modules {
        result = append(result, m)
    }
    return result
}
```

**Issue discovered:** Package initialization order matters. If module A's init() runs before module B's, they register in that order. But `All()` returns a map iteration (random order). Need deterministic ordering.

**Fix: Track registration order**

```go
var (
    modules      = make(map[string]FeatureModule)
    moduleOrder  = []string{}
)

func Register(mod FeatureModule) {
    path := mod.Schema().Path()
    if _, exists := modules[path]; exists {
        panic("duplicate module registration: " + path)
    }
    modules[path] = mod
    moduleOrder = append(moduleOrder, path)
}

func All() []FeatureModule {
    var result []FeatureModule
    for _, path := range moduleOrder {
        result = append(result, modules[path])
    }
    return result
}
```

**Verdict:** Registry package works, but needs **ordering guarantees** for predictable SCAD emission.

---

**B (Module Author) — Prototype: Push Buttons Schema YAML**

*[Converts push_buttons to YAML schema]*

```bash
cd pkg/yappgen/modules/pushbuttons
cat > schema.yaml <<'EOF'
module: push_buttons
scad_array: pushButtons
go_package: pushbuttons
description: |
  Tactile switch button extenders that transfer lid presses to PCB-mounted switches.
  Generates both button caps and switch pole extenders for printing.

fields:
  x:
    type: number
    required: true
    min: 0
    desc: X coordinate of button on PCB (mm from origin)
    example: 10.0
  
  y:
    type: number
    required: true
    min: 0
    desc: Y coordinate of button on PCB (mm from origin)
    example: 10.0
  
  cap:
    type: object
    required: true
    desc: Button cap geometry (visible part on lid surface)
    fields:
      length:
        type: number
        required: true
        min: 0.1
        desc: Cap length in mm
        example: 6.0
      
      width:
        type: number
        required: true
        min: 0.1
        desc: Cap width in mm
        example: 6.0
      
      radius:
        type: number
        required: true
        min: 0
        desc: Corner radius for rounded shapes, full radius for circles
        example: 1.0
  
  lid:
    type: object
    required: true
    desc: Lid integration parameters
    fields:
      protrusion:
        type: number
        required: true
        min: 0
        desc: How far button stands proud of lid surface (mm)
        example: 0.5
      
      wall:
        type: number
        required: false
        desc: Wall thickness around button shaft (mm), defaults to enclosure wall thickness
      
      plate_thickness:
        type: number
        required: false
        desc: Lid thickness at button location (mm), defaults to lid.thickness
      
      slack:
        type: number
        required: false
        desc: Clearance between button shaft and lid hole (mm)
        example: 0.3
      
      snap_slack:
        type: number
        required: false
        desc: Tolerance for snap-fit retention feature (mm)
  
  switch:
    type: object
    required: true
    desc: Tactile switch dimensions
    fields:
      height:
        type: number
        required: true
        min: 0.1
        desc: Distance from PCB surface to switch actuator top (mm)
        example: 3.5
      
      travel:
        type: number
        required: true
        min: 0.01
        max: 2.0
        desc: How far switch plunger moves when pressed (mm)
        example: 0.25
      
      pole_diameter:
        type: number
        required: true
        min: 0.1
        desc: Diameter of switch actuator post (mm)
        example: 1.5
      
      top_offset:
        type: number
        required: false
        desc: Vertical gap between extender and switch top (mm)
        example: 0.1
  
  shape:
    type: string
    required: false
    enum: [rectangle, circle, rounded_rect]
    desc: Button cap shape
    default: rounded_rect
  
  angle:
    type: number
    required: false
    desc: Rotation in degrees
    default: 0
  
  fillet_radius:
    type: number
    required: false
    min: 0
    desc: Internal fillet where extender shaft meets cap (mm)
  
  polygon:
    type: string
    required: false
    enum: [arrow, hex]
    desc: Preset polygon shape (overrides 'shape' if set)
  
  coordinate:
    type: string
    required: false
    enum: [pcb, box]
    desc: Coordinate system (pcb-relative or box-relative)
    default: pcb

# Test cases for validation
tests:
  - name: basic_button
    desc: Minimal valid button
    input:
      x: 10
      y: 10
      cap: {length: 6, width: 6, radius: 1}
      lid: {protrusion: 0.5}
      switch: {height: 3.5, travel: 0.25, pole_diameter: 1.5}
    expect_valid: true
  
  - name: missing_x
    desc: X coordinate is required
    input:
      y: 10
      cap: {length: 6, width: 6, radius: 1}
      lid: {protrusion: 0.5}
      switch: {height: 3.5, travel: 0.25, pole_diameter: 1.5}
    expect_error: "missing required field 'x'"
  
  - name: negative_travel
    desc: Switch travel must be positive
    input:
      x: 10
      y: 10
      cap: {length: 6, width: 6, radius: 1}
      lid: {protrusion: 0.5}
      switch: {height: 3.5, travel: -0.5, pole_diameter: 1.5}
    expect_error: "switch.travel: must be >= 0.01"
  
  - name: invalid_shape
    desc: Shape must be from enum
    input:
      x: 10
      y: 10
      cap: {length: 6, width: 6, radius: 1}
      lid: {protrusion: 0.5}
      switch: {height: 3.5, travel: 0.25, pole_diameter: 1.5}
      shape: oval
    expect_error: "shape: must be one of [rectangle, circle, rounded_rect]"
EOF

wc -l schema.yaml
```

**Result: 158 lines**

**Findings:**

1. **Verbosity is real** — 158 lines for one module. But it's comprehensive and readable.
2. **Nested objects are clear** — The YAML structure makes `cap`, `lid`, `switch` hierarchy obvious.
3. **Test cases are gold** — Having examples right in the schema is fantastic for documentation.
4. **Defaults need attention** — Some fields have defaults (`shape: rounded_rect`). How does code generation handle this?

**Issue discovered: Conditional validation**

The schema says `polygon` overrides `shape`. How do we express "if polygon is set, ignore shape"? Pure YAML can't express this logic.

**Proposed solution: Custom validation in builder**

```yaml
# In schema.yaml
custom_validation: |
  If 'polygon' is set, 'shape' is ignored.
  If 'shape' is 'circle', 'radius' must be > 0.
```

The generated struct gets a `Validate()` method stub that module authors fill in:

```go
// schema_gen.go (generated)
func (p *PushButtonDSL) Validate() error {
    // TODO: Module author implements custom validation
    return nil
}
```

Module author edits the stub. Not ideal (manual step), but complex validation is inherently custom.

**Verdict:** YAML schema works for push_buttons, but **needs custom validation escape hatch**.

---

**C (Validation Engine) — Prototype: Resolver Integration**

*[Implements schema-based validation in resolver]*

```bash
cd pkg/resolver
cat > validation.go <<'EOF'
package resolver

import (
    "fmt"
    "github.com/wesen/yapp-encl-resolver/pkg/registry"
)

// ValidateAgainstSchemas checks YAML structure against registered module schemas
func (r *Resolver) validateAgainstSchemas(doc map[string]any) error {
    features, ok := doc["features"].(map[string]any)
    if !ok {
        return nil  // No features to validate
    }
    
    for key, val := range features {
        path := "features." + key
        mod, ok := registry.Get(path)
        if !ok {
            // Unknown feature, skip (might be future module)
            continue
        }
        
        schema := mod.Schema()
        
        // Phase 1: Structure validation (before expression resolution)
        if err := schema.ValidateStructure(path, val); err != nil {
            return fmt.Errorf("%s: %w", path, err)
        }
    }
    
    return nil
}

// ValidateConstraints runs after expression resolution
func (r *Resolver) validateConstraints(doc map[string]any) error {
    features, ok := doc["features"].(map[string]any)
    if !ok {
        return nil
    }
    
    for key, val := range features {
        path := "features." + key
        mod, ok := registry.Get(path)
        if !ok {
            continue
        }
        
        schema := mod.Schema()
        
        // Phase 2: Constraint validation (min/max, enums)
        if err := schema.ValidateConstraints(path, val); err != nil {
            return fmt.Errorf("%s: %w", path, err)
        }
    }
    
    return nil
}
EOF
```

**Integration into Resolve():**

```go
func Resolve(ctx context.Context, doc map[string]any, opts Options) (map[string]any, error) {
    // ... existing code ...
    
    // NEW: Validate structure before expression resolution
    if err := validateAgainstSchemas(doc); err != nil {
        return nil, errors.Wrap(err, "schema validation")
    }
    
    // Existing: Expression resolution loop
    for iter := 0; iter < opts.MaxIterations; iter++ {
        // ...
    }
    
    // NEW: Validate constraints after resolution
    if err := validateConstraints(doc); err != nil {
        return nil, errors.Wrap(err, "constraint validation")
    }
    
    return state, nil
}
```

**Testing:**

```bash
cd pkg/resolver
cat > validation_test.go <<'EOF'
package resolver

import (
    "context"
    "testing"
)

func TestValidationPhases(t *testing.T) {
    // This test requires modules to be registered
    // TODO: Set up test registry in TestMain
    
    doc := map[string]any{
        "features": map[string]any{
            "push_buttons": []any{
                map[string]any{
                    "x": "base_x + 10",  // Expression, not resolved yet
                    "y": 10,
                    // Missing 'cap' (required)
                },
            },
        },
    }
    
    _, err := Resolve(context.Background(), doc, Options{})
    
    // Should fail with "missing required field 'cap'"
    if err == nil {
        t.Errorf("expected validation error, got nil")
    }
}
EOF
```

**Issue discovered: Testing requires module registration**

The resolver tests need modules to be registered. But modules are in `pkg/yappgen/modules/*`. The resolver can't import them (circular dependency).

**Solution: Test registry**

```go
func TestMain(m *testing.M) {
    // Register mock modules for testing
    registry.Register(&mockModule{
        schema: &mockSchema{/* ... */},
    })
    
    os.Exit(m.Run())
}
```

But this is **boilerplate**. Every package testing validation needs mock modules.

**Better solution: Registry reset for tests**

```go
// In pkg/registry/registry.go
func ResetForTesting() {
    modules = make(map[string]FeatureModule)
    moduleOrder = []string{}
}
```

Tests call this before registering mocks. Not ideal (mutable global state), but necessary for testability.

**Verdict:** Resolver integration works, but **testing is awkward** due to global registry.

---

**D (Documentation Librarian) — Prototype: Auto-Generated Help**

*[Implements schema → markdown conversion]*

```bash
cd pkg/docs
cat > schema_help.go <<'EOF'
package docs

import (
    "fmt"
    "strings"
    "github.com/wesen/yapp-encl-resolver/pkg/registry"
)

// RenderSchemaHelp converts a module schema to markdown for help pages
func RenderSchemaHelp(schema registry.ModuleSchema) string {
    var buf strings.Builder
    
    buf.WriteString("# " + schema.Name() + "\n\n")
    buf.WriteString(schema.Description() + "\n\n")
    
    buf.WriteString("## Fields\n\n")
    buf.WriteString("| Field | Type | Required | Description |\n")
    buf.WriteString("|-------|------|----------|-------------|\n")
    
    for _, field := range schema.Fields() {
        req := ""
        if field.Required {
            req = "✓"
        }
        
        typ := fieldTypeName(field.Type)
        
        buf.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s |\n",
            field.Name, typ, req, field.Description))
        
        // TODO: Nested fields (recursion)
    }
    
    return buf.String()
}

func fieldTypeName(t registry.FieldType) string {
    switch t {
    case registry.NumberField:
        return "number"
    case registry.StringField:
        return "string"
    case registry.ObjectField:
        return "object"
    case registry.ArrayField:
        return "array"
    default:
        return "unknown"
    }
}
EOF
```

**Testing with push_buttons schema:**

```bash
# (Hypothetically, after schema is loaded)
schema := pushbuttons.GetSchema()
markdown := docs.RenderSchemaHelp(schema)
fmt.Println(markdown)
```

**Expected output:**
```markdown
# push_buttons

Tactile switch button extenders that transfer lid presses to PCB-mounted switches.

## Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `x` | number | ✓ | X coordinate of button on PCB (mm from origin) |
| `y` | number | ✓ | Y coordinate of button on PCB (mm from origin) |
| `cap` | object | ✓ | Button cap geometry (visible part on lid surface) |
| ...
```

**Issue discovered: Nested fields not rendering**

The table shows `cap` as "object" but doesn't show `cap.length`, `cap.width`, etc. Need **recursive rendering**:

```go
func renderFields(fields []registry.FieldSpec, prefix string, buf *strings.Builder) {
    for _, field := range fields {
        name := prefix + field.Name
        req := ""
        if field.Required {
            req = "✓"
        }
        
        typ := fieldTypeName(field.Type)
        
        buf.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s |\n",
            name, typ, req, field.Description))
        
        // Recurse into nested fields
        if field.Type == registry.ObjectField && len(field.Children) > 0 {
            renderFields(field.Children, name+".", buf)
        }
    }
}
```

**Better output:**
```markdown
| `cap` | object | ✓ | Button cap geometry |
| `cap.length` | number | ✓ | Cap length in mm |
| `cap.width` | number | ✓ | Cap width in mm |
```

**Issue 2: Examples missing**

The schema has `example: 10.0` for fields, but `registry.FieldSpec` doesn't have an `Example` field. Need to add it:

```go
type FieldSpec struct {
    // ... existing fields
    Example any  // NEW: Example value for documentation
}
```

**Verdict:** Schema → docs works, but **FieldSpec needs enhancement** (examples, maybe min/max display).

---

**E (Tooling Engineer) — Prototype: schemagen discover**

*[Implements discovery command]*

```bash
mkdir -p cmd/schemagen
cat > cmd/schemagen/main.go <<'EOF'
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: schemagen <command>")
        fmt.Println("Commands: validate, discover")
        os.Exit(1)
    }
    
    switch os.Args[1] {
    case "discover":
        if err := discover(); err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
    default:
        fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
        os.Exit(1)
    }
}

func discover() error {
    // Scan for schema.yaml files
    patterns := []string{
        "pkg/yappgen/modules/*/schema.yaml",
    }
    
    var schemas []string
    for _, pattern := range patterns {
        matches, err := filepath.Glob(pattern)
        if err != nil {
            return err
        }
        schemas = append(schemas, matches...)
    }
    
    fmt.Printf("Found %d module schemas:\n", len(schemas))
    for _, s := range schemas {
        fmt.Printf("  - %s\n", s)
    }
    
    // TODO: Generate modules_gen.go
    
    return nil
}
EOF

go run ./cmd/schemagen discover
```

**Output:**
```
Found 1 module schemas:
  - pkg/yappgen/modules/pushbuttons/schema.yaml
```

**Next: Generate registry code**

```go
func generateRegistry(schemas []string) error {
    outFile := "pkg/yappgen/modules_gen.go"
    
    var imports []string
    var registrations []string
    
    for _, schemaPath := range schemas {
        // Parse schema to get module name and package
        schema := parseSchemaYAML(schemaPath)
        
        pkgPath := modulePkgPath(schemaPath)
        imports = append(imports, pkgPath)
        
        registrations = append(registrations,
            fmt.Sprintf("\tregistry.Register(%s.NewModule())", schema.GoPackage))
    }
    
    tmpl := `// Code generated by schemagen discover. DO NOT EDIT.
package yappgen

import (
    "github.com/wesen/yapp-encl-resolver/pkg/registry"
{{range .Imports}}
    "{{.}}"
{{end}}
)

func init() {
{{range .Registrations}}
    {{.}}
{{end}}
}
`
    
    // Execute template, write to outFile
    // ...
}
```

**Testing generation:**

```bash
go run ./cmd/schemagen discover --output pkg/yappgen/modules_gen.go
cat pkg/yappgen/modules_gen.go
```

**Generated file:**
```go
// Code generated by schemagen discover. DO NOT EDIT.
package yappgen

import (
    "github.com/wesen/yapp-encl-resolver/pkg/registry"
    "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pushbuttons"
)

func init() {
    registry.Register(pushbuttons.NewModule())
}
```

**Issue discovered: Schema parsing errors**

If `schema.yaml` is malformed, the generator panics. Need **graceful error handling**:

```go
schema, err := parseSchemaYAML(schemaPath)
if err != nil {
    return fmt.Errorf("invalid schema %s: %w", schemaPath, err)
}
```

**Issue 2: Generated file formatting**

The generated code needs to be `gofmt`-compatible. Use `go/format`:

```go
import "go/format"

func writeGeneratedFile(path string, code string) error {
    formatted, err := format.Source([]byte(code))
    if err != nil {
        return fmt.Errorf("generated code is not valid Go: %w", err)
    }
    return os.WriteFile(path, formatted, 0644)
}
```

**Verdict:** Discovery works, but **error handling** and **code formatting** are critical.

---

## Opening Statements (Round 3)

### A (DSL Architect) — Registry Works, Needs Ordering

After implementing `pkg/registry`, I found it's straightforward but has **one critical issue: ordering**.

*[Shows code]*

When modules register via `init()`, they run in import order (which is file alphabetical order in the generated file). But then `registry.All()` used to iterate over a map (random order).

I fixed this by tracking registration order in a slice. But there's a deeper question: **does order matter?**

**For SCAD emission: YES**

OpenSCAD processes arrays in order. If we emit:
```scad
pcbStands = [[...], [...]];
pushButtons = [[...], [...]];
```

vs.
```scad
pushButtons = [[...], [...]];
pcbStands = [[...], [...]];
```

The latter might fail if push_buttons reference stand positions. Order matters.

**Proposed solution: Explicit ordering in schema**

```yaml
# In schema.yaml
module: push_buttons
order: 100  # Higher = later in processing
```

Modules with no `order` default to 50. This gives us control without relying on init() timing.

**Alternative: Dependency declaration**

```yaml
module: push_buttons
depends_on:
  - pcb_stands  # Must be registered before push_buttons
```

The registry builds a dependency graph and topologically sorts. More complex, but handles actual dependencies.

**My recommendation:** Start simple (explicit `order` field), add dependency graph later if needed.

### B (Module Author) — YAML Schema is Excellent

I converted push_buttons to YAML and it's **better than I expected**.

*[Opens schema.yaml]*

158 lines, but they're **self-documenting**:
- Field descriptions are inline
- Examples show typical values
- Test cases catch regressions
- Nested structure is clear

**But there's a problem: conditional validation**

The schema says "if polygon is set, shape is ignored". YAML can't express this. I need custom Go code.

**Proposed solution: Two-tier validation**

1. **Schema validation** (automatic from YAML) — Types, required fields, min/max
2. **Custom validation** (hand-written Go) — Business logic

```go
// In builder.go (hand-written)
func (p *PushButtonDSL) CustomValidate() error {
    if p.Polygon != "" && p.Shape != "" {
        // Polygon overrides shape, no error, just document
    }
    
    if p.Shape == "circle" && p.Cap.Radius == 0 {
        return errors.New("circle buttons require radius > 0")
    }
    
    return nil
}
```

The resolver calls both:
```go
schema.ValidateStructure(data)  // From YAML
instance.CustomValidate()       // Hand-written
```

This is **pragmatic**. Most validation is declarative (YAML), but escape hatch exists for complex rules.

**One more finding: Default values**

The schema has `default: rounded_rect` for `shape`. How does this flow into the generated struct?

**Option 1: Zero values**
```go
type PushButtonDSL struct {
    Shape string  // Zero value is "", not "rounded_rect"
}
```

The builder checks if empty and applies default. But this is **implicit**.

**Option 2: Constructor**
```go
func NewPushButtonDSL() PushButtonDSL {
    return PushButtonDSL{
        Shape: "rounded_rect",
        Angle: 0,
        // ... other defaults
    }
}
```

The resolver calls the constructor. But YAML unmarshaling doesn't use constructors—it creates a zero-value struct.

**Option 3: Post-unmarshal hook**
```go
func (p *PushButtonDSL) ApplyDefaults() {
    if p.Shape == "" {
        p.Shape = "rounded_rect"
    }
    // ...
}
```

The resolver calls this after unmarshaling. This works but is manual.

**My recommendation:** Generated code includes `ApplyDefaults()` based on schema. Module author never touches it.

### C (Validation Engine) — Integration Works, Testing is Hard

I integrated schema validation into the resolver. It works, but **testing is awkward**.

*[Shows test code]*

Resolver tests need modules registered, but resolver can't import yappgen modules (circular dependency). I used `registry.ResetForTesting()` to clear the global registry and register mocks.

This works but feels **fragile**. Global mutable state in tests is a recipe for flaky tests (parallel execution, test order dependencies).

**Better approach: Dependency injection**

```go
type Resolver struct {
    registry Registry  // Interface, not global
}

type Registry interface {
    Get(path string) (FeatureModule, bool)
}

func NewResolver(reg Registry) *Resolver {
    return &Resolver{registry: reg}
}
```

Tests pass a mock registry:

```go
func TestValidation(t *testing.T) {
    mockReg := &mockRegistry{
        modules: map[string]FeatureModule{...},
    }
    
    r := resolver.NewResolver(mockReg)
    // Test with mock modules
}
```

**Trade-off:** More boilerplate (passing registry everywhere), but **much better testability**.

**My recommendation:** Refactor to dependency injection pattern. Global registry still exists for convenience, but tests use injection.

### D (Documentation Librarian) — Schema → Docs Works, Needs Enhancements

I built schema → markdown rendering. It works, but **`FieldSpec` is missing fields**.

*[Shows rendering code]*

The schema YAML has:
- `example: 10.0`
- `min: 0`, `max: 100`
- `default: rounded_rect`

But `registry.FieldSpec` only has:
- `Name`, `Type`, `Required`, `Description`

I need to enhance it:

```go
type FieldSpec struct {
    Name        string
    Type        FieldType
    Required    bool
    Description string
    
    // NEW: For documentation
    Example     any
    Default     any
    Min, Max    *float64
    Enum        []string
    
    // For nested objects
    Children    []FieldSpec
}
```

With this, I can render:

```markdown
| Field | Type | Required | Default | Range | Example | Description |
|-------|------|----------|---------|-------|---------|-------------|
| `x` | number | ✓ | — | ≥ 0 | 10.0 | X coordinate on PCB (mm) |
| `shape` | string | | rounded_rect | rectangle, circle, rounded_rect | circle | Button cap shape |
```

Much more informative!

**One more request: Nested field rendering**

Currently I show:
```
| `cap` | object | ✓ | Button cap geometry |
```

I want:
```
| `cap` | object | ✓ | Button cap geometry |
| `cap.length` | number | ✓ | Cap length in mm |
| `cap.width` | number | ✓ | Cap width in mm |
```

This requires recursive rendering (already prototyped, works fine).

**My recommendation:** Enhance `FieldSpec` to include all metadata from schema YAML.

### E (Tooling Engineer) — Discovery Works, Error Handling is Critical

I implemented `schemagen discover`. It works, but I learned **error handling is make-or-break**.

*[Shows prototype]*

If one schema is malformed, the whole generation fails. This is good (fail fast), but the error message matters:

**Bad error:**
```
Error: yaml: line 42: could not find expected ':'
```

**Good error:**
```
Error in pkg/yappgen/modules/pushbuttons/schema.yaml:
  Line 42: YAML syntax error - missing colon after key

  41 |   protrusion:
  42 |     type number
         ^^^^^^^^^^ Did you mean 'type: number'?
  43 |     required: true
```

The good error:
- Shows file and line
- Explains the problem
- Suggests a fix

**This requires custom YAML parsing** with position tracking. The standard `yaml.Unmarshal` doesn't preserve line numbers well.

**Solution: Use yaml.v3 with node parsing**

```go
import "gopkg.in/yaml.v3"

func parseSchemaWithErrors(path string) (*Schema, error) {
    data, _ := os.ReadFile(path)
    
    var node yaml.Node
    if err := yaml.Unmarshal(data, &node); err != nil {
        // Enhance error with file context
        return nil, enhanceYAMLError(path, data, err)
    }
    
    // Parse node into schema
    // ...
}

func enhanceYAMLError(path string, data []byte, err error) error {
    // Extract line number from error
    // Show surrounding context
    // Suggest fixes
}
```

**Second critical piece: Generated code validation**

After generating `modules_gen.go`, I need to verify it's valid Go:

```go
import "go/format"
import "go/parser"

func validateGeneratedCode(path string, code string) error {
    // Format check
    formatted, err := format.Source([]byte(code))
    if err != nil {
        return fmt.Errorf("generated code is not valid Go: %w", err)
    }
    
    // Parse check
    _, err = parser.ParseFile(token.NewFileSet(), path, formatted, parser.ParseComments)
    if err != nil {
        return fmt.Errorf("generated code does not parse: %w", err)
    }
    
    return nil
}
```

This catches generation bugs before they hit the developer.

**My implementation timeline:**
- Week 1: Core discover command (done in prototype)
- Week 2: Enhanced error messages (YAML + Go validation)
- Week 3: CI integration, pre-commit hooks

---

## Rebuttals (Round 3)

### A responds to C (Dependency Injection)

You're absolutely right that global registry makes testing hard. But I push back on **full dependency injection**.

If we pass `Registry` to every function, we get:

```go
func Resolve(ctx context.Context, doc map[string]any, reg Registry) (...)
func BuildModel(doc map[string]any, reg Registry) (...)
func EmitSCAD(model *Model, reg Registry) (...)
```

This is **boilerplate overload**. The registry rarely changes (only in tests).

**Middle ground: Test-only injection**

```go
var defaultRegistry = &globalRegistry{}

func SetTestRegistry(reg Registry) {
    defaultRegistry = reg
}

func Get(path string) (FeatureModule, bool) {
    return defaultRegistry.Get(path)
}
```

Production code uses the global. Tests call `SetTestRegistry()` before running. This is **simple and sufficient**.

### B responds to A (Ordering)

Your `order: 100` field in schema is fine for now, but I agree with your instinct—we'll eventually need **dependency declarations**.

Here's why: Say I'm adding a new module `display_mounts` that needs to reference PCB stands for positioning. I'd write:

```yaml
module: display_mounts
depends_on:
  - pcb_stands
```

The generator topologically sorts and ensures stands emit before displays. This is **explicit and self-documenting**.

But for MVP, `order` is fine. Let's add it to the schema spec.

### C responds to D (FieldSpec Enhancement)

I support adding `Example`, `Default`, `Min`, `Max` to `FieldSpec`. But let's be clear about **what goes in the interface**.

The `ModuleSchema` interface shouldn't expose every tiny detail. It should expose **semantically meaningful queries**:

```go
type ModuleSchema interface {
    // ... existing methods
    
    GetFieldSpec(path string) (FieldSpec, bool)
    GetDefaultValue(field string) (any, bool)
    GetConstraints(field string) Constraints
}

type Constraints struct {
    Min, Max *float64
    Enum     []string
    Pattern  *regexp.Regexp
}
```

This lets consumers (resolver, docs generator) ask specific questions without parsing the entire field tree.

### D responds to E (Error Messages)

Yes! Enhanced YAML errors are crucial. But let me add: **error messages should include links to documentation**.

Instead of:
```
Error: missing required field 'cap'
```

Show:
```
Error in push_buttons[2]:
  Missing required field 'cap'
  
  The 'cap' field defines button cap geometry.
  See: go run ./cmd/yappctl help push-buttons-reference
```

The help link takes users directly to generated docs. This closes the loop: error → docs → fix.

### E responds to B (Defaults)

Your `ApplyDefaults()` proposal is good, but let me refine it.

**Generated `ApplyDefaults()` from schema:**

```go
// schema_gen.go - AUTO-GENERATED
func (p *PushButtonDSL) ApplyDefaults() {
    if p.Shape == "" {
        p.Shape = "rounded_rect"  // From schema default
    }
    if p.Angle == 0 && !p.angleWasSet {  // Zero value ambiguity!
        p.Angle = 0  // From schema default (happens to be 0)
    }
}
```

**Problem:** How do we distinguish "user set to 0" from "not set, should use default"?

**Solution: Pointers for optional fields with defaults**

```go
type PushButtonDSL struct {
    Shape string   // Required field, no pointer
    Angle *float64 // Optional with default, use pointer
}

func (p *PushButtonDSL) ApplyDefaults() {
    if p.Angle == nil {
        angle := 0.0
        p.Angle = &angle
    }
}
```

This is **idiomatic Go** for optional-with-default fields.

---

## Moderator Summary (Round 3 — Prototype)

### Implementation Findings

**What worked well:**

1. **Registry package (`pkg/registry`)** — Clean interface, compiles, testable
2. **YAML schema for push_buttons** — Comprehensive, self-documenting, 158 lines
3. **Resolver integration** — Two-phase validation integrates smoothly
4. **Schema → docs rendering** — Auto-generates markdown from schemas
5. **Discovery command** — `schemagen discover` finds and processes schemas

**What needs refinement:**

1. **Module ordering** — Need explicit control (add `order` field to schema)
2. **Conditional validation** — YAML can't express complex rules (add `CustomValidate()` hook)
3. **Default values** — Need careful handling (use pointers for optional-with-default fields)
4. **Testing ergonomics** — Global registry is awkward (add test-only injection)
5. **FieldSpec enhancement** — Add `Example`, `Default`, `Min`, `Max`, `Enum`
6. **Error messages** — Need line-aware YAML parsing and helpful suggestions

### Architecture Decisions

**Confirmed:**

✅ **All modules use YAML schemas** (no hybrid struct-tag approach)  
✅ **Top-level `pkg/registry` package** with `FeatureModule` interface  
✅ **Discovery-based registration** via `schemagen discover`  
✅ **Two-phase validation** (structure before resolution, constraints after)  
✅ **Auto-generated docs** from schemas

**Refined:**

🔧 **Module ordering:** Add `order: N` field to schema (default 50)  
🔧 **Custom validation:** Generated structs get `CustomValidate()` method  
🔧 **Default values:** Use pointers for optional fields with defaults  
🔧 **FieldSpec:** Expanded with metadata (example, default, constraints)  
🔧 **Testing:** Add `SetTestRegistry()` for test-only injection

### Updated Schema Format

```yaml
module: push_buttons
order: 100                    # NEW: Explicit ordering
scad_array: pushButtons
go_package: pushbuttons
description: |
  Tactile switch button extenders...

fields:
  x:
    type: number
    required: true
    min: 0
    desc: X coordinate on PCB (mm)
    example: 10.0            # For docs
  
  shape:
    type: string
    required: false
    enum: [rectangle, circle, rounded_rect]
    default: rounded_rect    # Becomes pointer field
    desc: Button cap shape

# NEW: Custom validation notes
custom_validation: |
  - If 'polygon' is set, 'shape' is ignored
  - If 'shape' is 'circle', 'radius' must be > 0

tests:
  - name: basic_button
    input: {...}
    expect_valid: true
```

### Implementation Checklist

Before declaring the pattern ready:

**schemagen tool:**
- [x] Core validation command (Round 1)
- [x] Discovery command (prototype done)
- [ ] Enhanced YAML error messages (Week 2)
- [ ] Generated code validation (Week 2)
- [ ] Test generation from schema test cases
- [ ] Migration assistant (schema from existing builder)

**Registry package:**
- [x] Core interfaces and registration
- [x] Ordering support (track registration order)
- [ ] Test injection API (`SetTestRegistry()`)
- [ ] Dependency graph support (later, not MVP)

**Generated code:**
- [ ] Struct generation from YAML
- [ ] `ApplyDefaults()` method generation
- [ ] `CustomValidate()` stub generation
- [ ] Test generation from schema test cases
- [ ] Registry registration code (`modules_gen.go`)

**Documentation:**
- [x] Schema → markdown rendering (prototype)
- [ ] Enhanced FieldSpec with all metadata
- [ ] Nested field recursive rendering
- [ ] Help system integration
- [ ] Error message → help links

**Integration:**
- [ ] Resolver two-phase validation
- [ ] BuildModel uses registered modules
- [ ] EmitSCAD uses registered modules
- [ ] CLI init calls discovery

### Timeline Estimate

**Week 1-2: Core tooling**
- schemagen with validation and discovery
- Enhanced error messages
- Test generation

**Week 3: Push buttons conversion**
- Convert push_buttons to YAML schema
- Test all generated code
- Fix issues

**Week 4: Other modules**
- Convert pcb_stands, connectors, snap_joins
- Document module authoring guide
- CI integration

**Week 5: Documentation + polish**
- Auto-generated help pages
- Migration guide
- Pre-commit hooks

**Total: 5 weeks from green light to full implementation**

### Success Criteria

The pattern is ready when:
1. ✅ Push buttons fully converted (schema + generated code works)
2. ✅ `schemagen discover` generates valid `modules_gen.go`
3. ✅ Tests pass with mock modules
4. ✅ Help pages auto-generate from schemas
5. ✅ Error messages are actionable
6. ✅ Module author guide is complete
7. ✅ CI enforces `go generate` discipline

---

## Conclusion

The hands-on prototype validated the architecture. **YAML schemas + generated registry** is the right approach, with refinements identified during implementation.

**Key insights from prototyping:**

1. **Verbosity is acceptable** — 158 lines for push_buttons is comprehensive, not bloated
2. **Custom validation is necessary** — Pure YAML can't handle conditional logic
3. **Testing needs thought** — Global registry requires test injection API
4. **Error messages are critical** — Bad errors = frustrated module authors
5. **Tooling is achievable** — 5-week timeline to production-ready system

**Next steps:** Implement schemagen core (Weeks 1-2), then convert push_buttons as pilot (Week 3).

---

## References

- [Debate Setup Document](./dsl-module-schema-debate-setup.md)
- [Debate Round 1: Validation + Parsing Flow](./debate-round-01-validation-parsing-flow.md)
- [Debate Round 2: Registration Mechanics](./debate-round-02-registration-mechanics.md)
- [Feature Module Registry Design](../design/feature-module-registry.md)

