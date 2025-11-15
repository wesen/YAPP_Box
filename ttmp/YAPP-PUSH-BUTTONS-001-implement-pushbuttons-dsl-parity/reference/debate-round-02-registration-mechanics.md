---
Title: Debate Round 2 — Registration Mechanics
Ticket: YAPP-PUSH-BUTTONS-001
Status: active
Topics:
    - yapp
    - dsl
    - schema
    - architecture
DocType: debate
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/yappgen/features.go
    - Path: pkg/yappgen/modules/pushbuttons/module.go
    - Path: debate-round-01-validation-parsing-flow.md
ExternalSources: []
Summary: Second debate round on how modules register with the central registry without touching core files
LastUpdated: 2025-11-15
---

# Debate Round 2 — Registration Mechanics

## Question

**What is the right API for modules to register their typed struct + schema with the central registry so new modules remain easy to plug in without touching core files?**

---

## Pre-Debate Research (Round 2)

### Candidate Research: Module Registration Patterns

**A (DSL Architect) researched:**

*[Examines current registration pattern]*

```bash
cat pkg/yappgen/features.go | head -40
```

**Current pattern:**
```go
var featureModules = []FeatureModule{
    newArrayFeatureModule("pcb_stands", "pcbStands", ...),
    newArrayFeatureModule("connectors", "connectors", ...),
    newArrayFeatureModule("push_buttons", "pushButtons", 
        ..., pushbuttons.Build, ...),
    // ...
}
```

**Finding:** **Centralized registration** — Every module must be explicitly added to the `featureModules` slice in `pkg/yappgen/features.go`. This means:
- ✅ Clear, explicit list of all modules in one place
- ✅ Deterministic ordering (modules process in array order)
- ❌ Core file edit required for each new module
- ❌ Merge conflicts when multiple people add modules

*[Researches Go registration patterns]*

**Pattern 1: Init-based registration** (database/sql model)
```go
// In module package
func init() {
    yappgen.RegisterModule(&PushButtonsModule{})
}

// Core imports module for side effects
import _ "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pushbuttons"
```

**Pattern 2: Explicit registration** (flag package model)
```go
// In main.go or init code
func initModules() {
    registry.Register(pushbuttons.New())
    registry.Register(pcbstands.New())
}
```

**Pattern 3: Discovery-based** (plugin system)
```go
// Scan filesystem for schema.yaml files
modules := yappgen.DiscoverModules("pkg/yappgen/modules/*/schema.yaml")
```

**Verdict:** For code-generation-based schemas, **discovery is most elegant**. The `schemagen` tool embeds module registration into generated code, and the core never needs manual editing.

---

**B (Module Author) researched:**

*[Tests module authoring workflows]*

**Current workflow (centralized registration):**
```bash
# 1. Create module
mkdir pkg/yappgen/modules/myfeature
vim pkg/yappgen/modules/myfeature/builder.go

# 2. Edit core file (error-prone!)
vim pkg/yappgen/features.go
# Add newArrayFeatureModule("myfeature", ...) to slice

# 3. Update imports
# Add import statement if builder is external

# 4. Test
go test ./pkg/yappgen
```

**Pain points:**
- Must remember to edit features.go (easy to forget)
- Import statement might need to be added
- Merge conflicts if teammates also adding modules

**Hypothetical: Init-based registration**
```bash
# 1. Create module with init function
mkdir pkg/yappgen/modules/myfeature
cat > pkg/yappgen/modules/myfeature/module.go <<EOF
package myfeature

import "github.com/wesen/yapp-encl-resolver/pkg/yappgen"

func init() {
    yappgen.RegisterModule(New())
}
EOF

# 2. Import for side effects somewhere
vim pkg/yappgen/modules.go
# Add: import _ "github.com/.../modules/myfeature"

# 3. Test
go test ./pkg/yappgen
```

**Better, but still requires core file edit** (the import statement). And init() side effects are often considered bad practice in Go (hard to test, ordering issues).

**Hypothetical: Schema-driven registration**
```bash
# 1. Create schema
cat > pkg/yappgen/modules/myfeature/schema.yaml <<EOF
module: my_feature
scad_array: myFeature
# ... rest of schema
EOF

# 2. Generate code (includes registration!)
go generate ./pkg/yappgen/modules/myfeature

# 3. Test
go test ./pkg/yappgen/modules/myfeature
```

**This is ideal!** No core file edits. The schema.yaml presence triggers registration. The `schemagen` tool embeds the registration boilerplate.

**Verdict:** Schema-driven discovery wins for developer ergonomics.

---

**C (Validation Engine) researched:**

*[Examines resolver integration needs]*

```bash
grep "RegisterModuleSchema\|register.*schema" pkg/resolver/*.go
# No results — registration not yet implemented
```

**Current state:** The resolver has no module registration mechanism. It would need to be added.

**Requirements for resolver:**
1. **Schemas available before first Resolve()** — Can't load lazily
2. **Path-based lookup** — Given "features.push_buttons", find the schema
3. **No circular dependencies** — Resolver can't import yappgen (creates cycle)

**Proposed API:**

```go
// In pkg/resolver/schema_registry.go
var moduleSchemas = map[string]ModuleSchema{}

func RegisterModuleSchema(path string, schema ModuleSchema) {
    moduleSchemas[path] = schema
}

func GetModuleSchema(path string) (ModuleSchema, bool) {
    s, ok := moduleSchemas[path]
    return s, ok
}
```

**Who calls RegisterModuleSchema?**

**Option 1: Module init() functions**
```go
// In pkg/yappgen/modules/pushbuttons/module.go
func init() {
    resolver.RegisterModuleSchema("features.push_buttons", GetSchema())
}
```

**Option 2: Central registration (current pattern)**
```go
// In pkg/yappgen/features.go
func init() {
    resolver.RegisterModuleSchema("features.push_buttons", pushbuttons.GetSchema())
    resolver.RegisterModuleSchema("features.pcb_stands", pcbstands.GetSchema())
}
```

**Option 3: Auto-discovery at startup**
```go
// In cmd/yappctl/main.go
func main() {
    yappgen.DiscoverAndRegisterModules()  // Scans embedded schemas
    // ... rest of CLI
}
```

**Verdict:** Option 3 (auto-discovery) avoids circular imports and core file edits. The yappgen package can import resolver, discover modules, and register them. CLI just calls one init function.

---

**D (Documentation Librarian) researched:**

*[Examines documentation discovery]*

```bash
cat pkg/docs/docs.go
```

**Current pattern:**
```go
//go:embed tutorials/*.md
var tutorialFS embed.FS

func Load(helpSystem *help.HelpSystem) error {
    return helpSystem.LoadSectionsFromFS(tutorialFS, "tutorials")
}
```

**Finding:** Documentation uses **filesystem embedding** (`go:embed`). Help topics are discovered by scanning the embedded FS, not hardcoded.

**For module documentation:**

```go
//go:embed modules/*/schema.yaml
var moduleSchemasFS embed.FS

func LoadModuleHelp(helpSystem *help.HelpSystem) error {
    entries, _ := fs.ReadDir(moduleSchemasFS, "modules")
    for _, entry := range entries {
        if !entry.IsDir() { continue }
        
        schemaPath := path.Join("modules", entry.Name(), "schema.yaml")
        schemaBytes, _ := fs.ReadFile(moduleSchemasFS, schemaPath)
        
        schema := parseSchema(schemaBytes)
        helpSystem.AddHelpPage(
            schema.Module,
            renderSchemaAsMarkdown(schema),
        )
    }
    return nil
}
```

**This is auto-discovery!** No manual registration. Just drop a `schema.yaml` file in `modules/myfeature/` and it appears in help automatically.

**Verdict:** Filesystem-based discovery is the right pattern for documentation. It should be the same for module registration—consistent patterns across the codebase.

---

**E (Tooling Engineer) researched:**

*[Examines code generation integration points]*

**Question:** Where does generated code hook into the registry?

**Option 1: Generated code includes registration**
```go
// schema_gen.go (AUTO-GENERATED)
package pushbuttons

import "github.com/wesen/yapp-encl-resolver/pkg/yappgen"

func init() {
    yappgen.RegisterModule(&module{
        schema: GetSchema(),
        builder: Build,
    })
}
```

**Pro:** Self-contained module packages
**Con:** Init() ordering issues, import cycles

**Option 2: Central discovery scans generated code**
```go
// In pkg/yappgen/discovery.go
func DiscoverModules() []FeatureModule {
    // Embed all schema.yaml files
    schemas := embeddedSchemas.ReadDir("modules")
    
    var modules []FeatureModule
    for _, entry := range schemas {
        schema := loadSchema(entry)
        module := newModuleFromSchema(schema)
        modules = append(modules, module)
    }
    return modules
}
```

**Pro:** No init() side effects, explicit initialization
**Con:** Need build-time or runtime schema parsing

**Option 3: Code generation creates registry file**
```go
// modules_gen.go (AUTO-GENERATED from all schema.yaml files)
package yappgen

import (
    "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pushbuttons"
    "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pcbstands"
)

var featureModules = []FeatureModule{
    pushbuttons.NewModule(),
    pcbstands.NewModule(),
    // ... auto-populated from discovered schemas
}
```

**Pro:** Explicit, no init(), no discovery overhead
**Con:** Still requires generating a central file (but it's automated!)

**My testing:**

```bash
# Prototype Option 3
cd pkg/yappgen
schemagen discover --output modules_gen.go

# Generated file includes all modules found in modules/*/schema.yaml
cat modules_gen.go | head -20
```

**Verdict:** **Option 3 is best** for code generation. The `schemagen discover` command runs as part of `go generate ./...` and produces a single registry file. This gives us:
- No manual edits to core files
- Explicit registration (no init magic)
- Deterministic ordering
- Easy to audit (generated file is committed to Git)

---

## Opening Statements (Round 2)

### A (DSL Architect) — Discovery-Based Registration

After researching Go patterns, I'm convinced that **filesystem-based discovery** is the right architecture for module registration.

*[Opens research notes]*

Look at how Go's own standard library handles this:
- `database/sql` — Drivers register via init(), but **discovery is standardized**
- `net/http` — Handlers register explicitly
- `plugin` package — Dynamic loading via filesystem

Our situation is closer to `database/sql` than `net/http`. We have a **fixed set of module packages** (not dynamic plugins), but we want **loose coupling** between modules and core.

**My proposal: Two-tier discovery**

**Tier 1: Build-time discovery** (schemagen)
```bash
go generate ./pkg/yappgen
# schemagen scans modules/*/schema.yaml
# Generates modules_gen.go with registry
```

**Generated registry:**
```go
// modules_gen.go - AUTO-GENERATED by schemagen discover
package yappgen

import (
    "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pushbuttons"
    "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pcbstands"
    // ... auto-imported based on discovered schemas
)

var featureModules = []FeatureModule{
    pushbuttons.NewModule(),
    pcbstands.NewModule(),
    // ... auto-registered
}

// For resolver integration
func init() {
    for _, mod := range featureModules {
        resolver.RegisterModuleSchema("features."+mod.Name(), mod.GetSchema())
    }
}
```

**Tier 2: Runtime discovery** (for documentation)
```go
//go:embed modules/*/schema.yaml
var moduleSchemasFS embed.FS

func LoadModuleHelp() {
    // Scan embedded FS, generate help pages
}
```

**Benefits:**

1. **No manual core edits** — Add schema.yaml, run `go generate`, done
2. **No init() in modules** — Central init() in generated code
3. **Explicit imports** — Generated code has proper imports (Go compiler validates)
4. **Auditable** — Commit `modules_gen.go` to Git, diff shows what changed
5. **Resolver integration** — Generated init() registers schemas automatically

**Trade-offs:**

- Need to run `go generate` after adding modules (but Makefile already does this)
- Generated file gets stale if you forget `go generate` (CI catches it)

This pattern follows **Go conventions** (code generation, explicit imports) while achieving **zero-touch module addition** (no manual registry edits).

### B (Module Author) — Schema-Driven Auto-Registration

A's proposal is close to what I want, but let me refine the **developer experience**.

*[Shows workflow]*

**Today (manual registration):**
```bash
# Create module
mkdir -p pkg/yappgen/modules/myfeature
vim pkg/yappgen/modules/myfeature/schema.yaml
vim pkg/yappgen/modules/myfeature/builder.go

# MANUAL STEP: Edit features.go to add to registry
vim pkg/yappgen/features.go
# (easy to forget!)

go test ./pkg/yappgen
```

**With A's discovery:**
```bash
# Create module
mkdir -p pkg/yappgen/modules/myfeature
vim pkg/yappgen/modules/myfeature/schema.yaml

# Generate (automatic via Makefile)
go generate ./pkg/yappgen

# Test
go test ./pkg/yappgen/modules/myfeature
```

**This is great!** No manual registry edit. But I have concerns:

**Concern 1: Import cycles**

A's proposal has the generated code in `pkg/yappgen/modules_gen.go` importing module packages:

```go
import "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pushbuttons"
```

This is fine if modules don't import `yappgen`. But what if a module needs access to shared types in `yappgen`? We get a cycle:
- `yappgen` imports `modules/pushbuttons`
- `modules/pushbuttons` imports `yappgen` (for shared types)

**Concern 2: Generated code location**

Where does `modules_gen.go` live? If it's in `pkg/yappgen/`, it's part of the yappgen package. But modules are in `pkg/yappgen/modules/*`. This feels awkward.

**Alternative proposal: Module manifest file**

Instead of generating Go code, generate a **manifest** that's read at startup:

```yaml
# modules_manifest.yaml (AUTO-GENERATED)
modules:
  - name: push_buttons
    package: github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pushbuttons
    schema: modules/pushbuttons/schema.yaml
  
  - name: pcb_stands
    package: github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pcbstands
    schema: modules/pcbstands/schema.yaml
```

At startup, the registry loads this manifest and uses `go:embed` to load schemas. No generated Go code, no import cycles.

**But wait...** this loses compile-time validation. If a module package moves, the manifest breaks at runtime.

**Revised proposal:** **Hybrid approach**

1. **Core modules** (push_buttons, pcb_stands, etc.) — Use A's generated registry (explicit imports)
2. **External/plugin modules** — Use manifest-based discovery (optional)

This gives us flexibility: tight integration for built-in modules, loose coupling for extensions.

### C (Validation Engine) — Resolver Registry Independence

Both A and B are designing registration for `yappgen`. Let me clarify what the **resolver** needs.

*[Opens pkg/resolver/resolver.go]*

The resolver should **not** depend on `yappgen`. That would create circular dependencies:
- `yappgen` imports `resolver` (to resolve expressions)
- `resolver` imports `yappgen` (for module schemas)

**Circular dependency = bad architecture.**

**My proposal: Shared registry package**

```
pkg/
  registry/
    registry.go       # ModuleSchema interface, registration API
  resolver/
    resolver.go       # Imports registry, uses ModuleSchema
  yappgen/
    features.go       # Imports registry, registers modules
    modules/
      pushbuttons/
        module.go     # Imports registry for types
```

**registry package:**
```go
package registry

type ModuleSchema interface {
    Name() string
    Path() string  // e.g., "features.push_buttons"
    Fields() []FieldSpec
    ValidateStructure(data any) error
    ValidateConstraints(data any) error
}

var modules = map[string]ModuleSchema{}

func RegisterModule(schema ModuleSchema) {
    modules[schema.Path()] = schema
}

func GetModule(path string) (ModuleSchema, bool) {
    m, ok := modules[path]
    return m, ok
}
```

**Resolver uses it:**
```go
func (r *Resolver) validateModules(doc map[string]any) error {
    features, _ := doc["features"].(map[string]any)
    
    for key, val := range features {
        path := "features." + key
        schema, ok := registry.GetModule(path)
        if !ok {
            continue  // Unknown module, skip
        }
        
        if err := schema.ValidateStructure(val); err != nil {
            return err
        }
    }
    return nil
}
```

**yappgen registers modules:**
```go
// In generated modules_gen.go
func init() {
    registry.RegisterModule(pushbuttons.GetSchema())
    registry.RegisterModule(pcbstands.GetSchema())
}
```

**Key point:** The `registry` package is **shared infrastructure**. It breaks the circular dependency and gives us a clean API boundary.

A's discovery approach can work, but the generated code should register with `pkg/registry`, not call resolver or yappgen directly.

### D (Documentation Librarian) — Consistency With Existing Patterns

I've been maintaining the docs package, and I want **consistency** across the codebase.

*[Opens pkg/docs/docs.go]*

We already have a discovery pattern:

```go
//go:embed tutorials/*.md
var tutorialFS embed.FS

func Load(helpSystem *help.HelpSystem) error {
    return helpSystem.LoadSectionsFromFS(tutorialFS, "tutorials")
}
```

This pattern works beautifully:
- Drop a new `.md` file in `tutorials/`
- It automatically appears in `glaze help`
- No registration, no manifest

**I want the same for module documentation:**

```go
//go:embed modules/*/schema.yaml
var moduleSchemasFS embed.FS

func LoadModuleHelp(helpSystem *help.HelpSystem) error {
    entries, _ := fs.ReadDir(moduleSchemasFS, "modules")
    for _, entry := range entries {
        schema := loadSchema(moduleSchemasFS, entry.Name())
        helpPage := renderSchemaAsMarkdown(schema)
        helpSystem.AddHelpPage("module-"+schema.Name(), helpPage)
    }
    return nil
}
```

**This is runtime discovery**—no code generation needed for help. But here's the thing: **we're duplicating discovery logic**.

A's proposal has `schemagen discover` for code generation. I need runtime discovery for docs. These should **share logic**:

**Shared discovery library:**
```go
// In pkg/yappgen/discovery/discovery.go
func DiscoverModules(fsys fs.FS) ([]ModuleInfo, error) {
    entries, _ := fs.ReadDir(fsys, "modules")
    var modules []ModuleInfo
    
    for _, entry := range entries {
        schemaPath := path.Join("modules", entry.Name(), "schema.yaml")
        schemaBytes, _ := fs.ReadFile(fsys, schemaPath)
        schema := parseSchema(schemaBytes)
        
        modules = append(modules, ModuleInfo{
            Name:   schema.Name,
            Path:   entry.Name(),
            Schema: schema,
        })
    }
    return modules, nil
}
```

**schemagen uses it:**
```go
func generateRegistry() error {
    modules := discovery.DiscoverModules(os.DirFS("pkg/yappgen"))
    // Generate modules_gen.go with imports
}
```

**docs package uses it:**
```go
//go:embed modules/*/schema.yaml
var schemasFS embed.FS

func LoadModuleHelp() {
    modules := discovery.DiscoverModules(schemasFS)
    // Render help pages
}
```

**One discovery algorithm, two consumers.** This ensures consistency—if a module is discovered for code generation, it's also discovered for docs.

### E (Tooling Engineer) — Making Discovery Bulletproof

Everyone's proposing discovery-based registration. Let me talk about **making it robust**.

*[Considers failure modes]*

**Problem 1: Partial generation**

What if `schemagen discover` fails halfway through? Do we get a partial `modules_gen.go`? How do we detect this?

**Solution: Atomic generation**
```bash
schemagen discover --output modules_gen.go.tmp
# Validate generated file
schemagen validate-generated modules_gen.go.tmp
# Atomic rename only if valid
mv modules_gen.go.tmp modules_gen.go
```

**Problem 2: Schema file errors**

What if one schema.yaml has syntax errors? Does the whole discovery fail?

**Solution: Fail fast with clear errors**
```bash
schemagen discover
# Output:
✓ Found modules/pushbuttons/schema.yaml (25 fields)
✗ modules/myfeature/schema.yaml:12 - invalid field type 'banana'
  
Error: Discovery failed due to invalid schemas.
Run 'schemagen validate modules/myfeature/schema.yaml' for details.
```

**Problem 3: Import path changes**

If a module package moves, the generated imports break. How do we detect this?

**Solution: Validate imports**
```go
// Generated code includes self-test
func init() {
    // Ensure all imported modules are accessible
    _ = pushbuttons.NewModule()
    _ = pcbstands.NewModule()
}
```

If an import is wrong, the code won't compile. Go's type system is our safety net.

**Problem 4: Forgetting go generate**

Developer adds schema.yaml but forgets to run `go generate`. Module isn't registered. Tests pass locally but fail in CI.

**Solution: Pre-commit hook + CI check**
```bash
# .git/hooks/pre-commit
#!/bin/bash
go generate ./pkg/yappgen
git diff --exit-code pkg/yappgen/modules_gen.go || {
    echo "Error: modules_gen.go is out of date. Commit the changes."
    exit 1
}

# In CI
go generate ./...
git diff --exit-code || {
    echo "Error: Generated files out of sync. Run 'go generate' locally."
    exit 1
}
```

**My implementation plan for schemagen discover:**

**Week 1 (after core schemagen is done):**
- Implement `schemagen discover` command
- Scan `modules/*/schema.yaml`
- Generate `modules_gen.go` with imports and registry

**Week 2:**
- Add validation (--check, --validate-generated flags)
- Error handling for malformed schemas
- CI integration documentation

**Week 3:**
- Pre-commit hook templates
- Testing with multiple modules
- Documentation for module authors

**Key principle:** Discovery should be **reliable and debuggable**. When it fails, errors must be actionable. No mysterious "module not found" at runtime—fail at generation or compile time.

---

## Rebuttals (Round 2)

### A (DSL Architect) responds to B

Your concern about import cycles is valid. Let me address it.

**Current architecture:**
```
pkg/yappgen/
  features.go           # FeatureModule interface
  model.go             # Model struct
  modules/
    pushbuttons/
      module.go         # Implements FeatureModule, imports yappgen
```

You're right that `modules/pushbuttons` might need to import `yappgen` for shared types. And if `modules_gen.go` (in yappgen package) imports pushbuttons, we get a cycle.

**Solution: Interface in separate package** (C's registry idea)

```
pkg/
  yappgen/
    registry/
      types.go         # FeatureModule interface, ModuleSchema
    features.go        # Uses registry.FeatureModule
    modules_gen.go     # Imports modules/*, registers with registry
    modules/
      pushbuttons/
        module.go      # Imports yappgen/registry (no cycle!)
```

The `yappgen/registry` subpackage holds interfaces. Modules import `yappgen/registry`, not `yappgen`. Generated code in `yappgen/modules_gen.go` can import both registry and modules without cycles.

This is a **standard Go pattern** for breaking circular dependencies (see `database/sql` and `database/sql/driver`).

### B (Module Author) responds to A and C

OK, the `registry` subpackage solves the import cycle. But I still worry about **developer confusion**.

New module authors need to understand:
- "Import `yappgen/registry` for interfaces"
- "Implement registry.FeatureModule"
- "Don't import yappgen directly"

That's three concepts instead of one. Can we simplify?

**Proposal: Single entrypoint package**

```go
// In pkg/yappgen/module/module.go
package module

// Public API for module authors
type Module interface {
    Name() string
    Build(items []any) ([][]any, error)
    GetSchema() Schema
}

type Schema struct {
    // ... schema definition
}
```

Module authors only import `yappgen/module`. The `yappgen` package internally uses the same types. No mental overhead about "which package do I import?"

This is how `net/http` works—there's `http.Handler`, not `http/handler.Handler`.

### C (Validation Engine) responds to A

I like the registry subpackage approach. But let's be precise about **where the registry lives**.

A suggested `pkg/yappgen/registry`. I say it should be **`pkg/registry`** (top-level).

**Why?**

1. **The resolver needs it** — If it's under yappgen, we're saying "this is a yappgen concept". But the resolver also uses module schemas. It's a **shared concept**.

2. **Future extensibility** — What if we add other kinds of registries? (e.g., custom functions for expressions, custom validators). A top-level `pkg/registry` can grow to accommodate these.

3. **Naming clarity** — `yappgen/registry` is awkward to import. You'd have `import ypgreg "github.com/.../yappgen/registry"` to avoid naming clashes. Just `import "github.com/.../registry"` is cleaner.

**Proposed structure:**
```
pkg/
  registry/
    registry.go        # RegisterModule, GetModule
    schema.go          # ModuleSchema interface
  resolver/
    resolver.go        # Imports pkg/registry
  yappgen/
    features.go        # Imports pkg/registry
    modules_gen.go     # Imports pkg/registry + modules
    modules/
      pushbuttons/
        module.go      # Imports pkg/registry
```

Everyone imports `pkg/registry`. Clean, simple, no ambiguity.

### D (Documentation Librarian) responds to all

I'm happy with the discovery approach, but I need **guarantees** about schema availability.

When I call `discovery.DiscoverModules(schemasFS)` to generate help, what if a schema is **embedded but invalid**? The help page renders broken YAML?

**Requirement: Schema validation at embed time**

```go
//go:embed modules/*/schema.yaml
var schemasFS embed.FS

func init() {
    // Validate all embedded schemas at startup
    modules, err := discovery.DiscoverModules(schemasFS)
    if err != nil {
        panic(fmt.Sprintf("Invalid embedded schema: %v", err))
    }
    registeredModules = modules
}
```

If a schema is malformed, the binary **won't start**. This forces schema correctness.

But wait—that's a **runtime check**. Can we move it to **compile time**?

**E's schemagen can help:** Add a generated file that validates schemas:

```go
// schemas_validated.go - AUTO-GENERATED
package yappgen

// This file is generated by 'schemagen discover --validate-embedded'.
// It ensures all embedded schemas are valid at compile time.

import (
    _ "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pushbuttons"
    // ^ Import forces compilation, which would fail if schema is bad
)

const schemasValidatedAt = "2025-11-15T10:30:00Z"
const schemasCount = 5
```

The imports force the Go compiler to parse the module packages. If their embedded schemas are invalid, `go build` fails.

This is **belt-and-suspenders** validation: schema validation during generation + compile-time import checks.

### E (Tooling Engineer) responds to all

Everyone's converging on discovery + generated registry. Let me summarize what `schemagen discover` needs to do:

**1. Scan for schemas**
```bash
find pkg/yappgen/modules -name schema.yaml
```

**2. Validate each schema**
```bash
for schema in *.yaml; do
    schemagen validate "$schema" || exit 1
done
```

**3. Generate registry code**
```go
// modules_gen.go
package yappgen

import (
    "github.com/wesen/yapp-encl-resolver/pkg/registry"
    "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pushbuttons"
    "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pcbstands"
)

func init() {
    registry.RegisterModule(pushbuttons.NewModule())
    registry.RegisterModule(pcbstands.NewModule())
}
```

**4. Self-validate**
```bash
# Ensure generated code compiles
go build -o /dev/null ./pkg/yappgen/modules_gen.go
```

**5. Atomic commit**
```bash
mv modules_gen.go.tmp modules_gen.go
```

**Implementation checklist:**
- [ ] `schemagen discover` command
- [ ] Schema scanning logic (filepath.Glob)
- [ ] Code generation template (text/template)
- [ ] Import path resolution
- [ ] Validation before generation
- [ ] CI integration docs
- [ ] Pre-commit hook template

This is **2-3 days of work** after the core schemagen tool is done. Very achievable.

---

## Moderator Summary (Round 2 — Registration)

### Key Arguments

**A (DSL Architect) — Discovery + Generated Registry:**
- **Strength:** Two-tier approach (build-time + runtime), follows Go conventions
- **Evolution:** Accepted registry subpackage to avoid circular imports
- **Trade-off:** Generated code must be committed, requires `go generate` discipline

**B (Module Author) — Developer Experience Focus:**
- **Strength:** Highlighted import cycle problem, pushed for simplified API
- **Evolution:** Accepted registry solution, proposed single entrypoint package
- **Trade-off:** Still two files per module (schema.yaml + builder.go)

**C (Validation Engine) — Clean Architecture:**
- **Strength:** Identified circular dependency risk, proposed shared registry package
- **Key insight:** Registry should be top-level (`pkg/registry`), not under yappgen
- **Trade-off:** One more package to understand, but cleaner dependencies

**D (Documentation Librarian) — Consistency + Validation:**
- **Strength:** Emphasized shared discovery logic, compile-time schema validation
- **Key contribution:** Proposed schemas_validated.go pattern for embed-time checking
- **Trade-off:** More generated files, but higher confidence

**E (Tooling Engineer) — Robust Implementation:**
- **Strength:** Detailed failure mode analysis, atomic generation, error handling
- **Key contribution:** Implementation checklist, CI integration plan
- **Trade-off:** Additional tooling complexity (pre-commit hooks, validation)

### Emerging Consensus

**Agreed architecture:**

1. **Top-level registry package** (`pkg/registry`)
   - ModuleSchema interface
   - RegisterModule/GetModule API
   - Imported by resolver, yappgen, and modules

2. **Discovery-based registration**
   - `schemagen discover` scans `modules/*/schema.yaml`
   - Generates `modules_gen.go` with imports and registration
   - No manual edits to core files

3. **Generated registry structure:**
```go
// pkg/yappgen/modules_gen.go - AUTO-GENERATED
package yappgen

import (
    "github.com/wesen/yapp-encl-resolver/pkg/registry"
    "github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pushbuttons"
)

func init() {
    registry.RegisterModule(pushbuttons.NewModule())
}
```

4. **Shared discovery library**
   - `pkg/yappgen/discovery` package
   - Used by both schemagen (code gen) and docs (runtime help)

5. **Safety guarantees:**
   - Schema validation before code generation
   - Compile-time import checks
   - CI verification (go generate + git diff)
   - Pre-commit hook (optional)

### Workflow Example

**Adding a new module:**

```bash
# 1. Create schema
mkdir -p pkg/yappgen/modules/myfeature
cat > pkg/yappgen/modules/myfeature/schema.yaml <<EOF
module: my_feature
scad_array: myFeature
fields:
  x: {type: number, required: true}
EOF

# 2. Generate code (automatic)
go generate ./pkg/yappgen
# Creates: modules/myfeature/schema_gen.go
# Updates: modules_gen.go (adds myfeature to registry)

# 3. Write builder
cat > pkg/yappgen/modules/myfeature/builder.go <<EOF
package myfeature
import "github.com/wesen/yapp-encl-resolver/pkg/registry"

func Build(items []MyFeatureDSL) ([][]any, error) {
    // ... implementation
}

func NewModule() registry.FeatureModule {
    return registry.NewTypedModule("my_feature", "myFeature", Build)
}
EOF

# 4. Test
go test ./pkg/yappgen/modules/myfeature

# 5. Done! No core file edits needed.
```

### Open Questions

1. **Module versioning** — How do we handle schema changes? (Planned: Question 4)
2. **External modules** — Can third parties create modules outside the main repo?
3. **Plugin architecture** — Should we support dynamic loading (.so files)?
4. **Registry initialization order** — Does it matter? (Probably not, but worth documenting)

### Next Debate Round

**Question 3 (suggested):** Given a straw-man implementation (e.g., refactoring `modules/pushbuttons` to use the new schema + registry pattern), what works well, what feels brittle, and what adjustments would we make before declaring the pattern ready?

This will be a **hands-on prototyping round** where candidates actually try the proposed architecture and report back on real issues.

---

## References

- [Debate Setup Document](./dsl-module-schema-debate-setup.md)
- [Debate Round 1: Validation + Parsing Flow](./debate-round-01-validation-parsing-flow.md)
- [Feature Module Registry Design](../design/feature-module-registry.md)
- [Go database/sql pattern](https://pkg.go.dev/database/sql) — Driver registration reference
- [Go plugin package](https://pkg.go.dev/plugin) — Dynamic loading reference

