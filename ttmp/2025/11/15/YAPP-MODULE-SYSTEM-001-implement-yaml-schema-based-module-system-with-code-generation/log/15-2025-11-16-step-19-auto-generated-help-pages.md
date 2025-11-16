---
Title: 2025-11-16 Step 19 - auto-generated help pages
Ticket: YAPP-MODULE-SYSTEM-001
Status: active
Topics:
    - yapp
    - dsl
    - codegen
    - architecture
DocType: log
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Schema-driven help pages now auto-load into yappctl
LastUpdated: 2025-11-15T21:40:16.914282997-05:00
---


# 2025-11-16 Step 19 - auto-generated help pages

<!-- Log entries in reverse chronological order (newest first) -->

## 2025-11-16 - Schema-driven help pages integrated

Implemented auto-generated documentation for all modules:

**Implementation:**
- Created `pkg/docs/schema_help.go` with `LoadModuleHelp()` and `renderSchemaHelp()`
- Iterates through `registry.All()` to find all registered modules
- Loads each module's `schema.yaml` via `schemagen.LoadSchemaDoc()`
- Renders markdown with frontmatter (Title, Slug, Topics, etc.)
- Generates field table with nested object support
- Parses markdown via `help.LoadSectionFromMarkdown()` and adds to help system
- Integrated into `pkg/docs/docs.go` Load() function

**Help Page Structure:**
- **YAML Location**: Shows where module appears in DSL
- **Fields Table**: Name, Type, Required, Default, Description (with nested fields indented)
- **Examples**: Renders test cases from schema (TODO: improve YAML rendering)
- **OpenSCAD Array**: Notes which SCAD array the module generates

**Testing:**
```bash
./yappctl help module-push_buttons
./yappctl help module-pcb_stands
./yappctl help module-connectors
./yappctl help module-snap_joins
./yappctl help module-cutouts
```

All 5 modules now have auto-generated help pages that update automatically when schemas change!

**Example output:**
- Field table shows nested objects (cap.length, lid.protrusion, switch.height)
- Enums display allowed values (shape: rectangle/circle/rounded_rect/...)
- Required fields marked with ✓
- Defaults shown when present
