---
Title: Implement pushButtons DSL parity
Ticket: YAPP-PUSH-BUTTONS-001
Status: active
Topics:
    - yapp
    - cli
    - dsl
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: examples/YAPP_Demo_buttons_v31.scad
      Note: Baseline SCAD reference for two-button lid
    - Path: examples/YAPP_Demo_buttons2_v31.scad
      Note: Shape variations to match in DSL
    - Path: examples/yapp-demo-buttons.yaml
      Note: Current YAML missing push_buttons section
    - Path: pkg/yappgen/schema.go
      Note: Location for new push button schema
    - Path: pkg/yappgen/map.go
      Note: Builder utilities that need push button support
    - Path: pkg/docs/tutorials/yapp-dsl-getting-started.md
      Note: Tutorial covering DSL basics for future push button docs to reference
    - Path: pkg/docs/tutorials/yapp-dsl-reference.md
      Note: Field reference for all DSL sections
    - Path: pkg/docs/tutorials/yapp-dsl-module-howto.md
      Note: Developer tutorial explaining how to add new DSL modules via the registry
    - Path: pkg/yappgen/modules/pushbuttons/module.go
      Note: Sub-package implementing the push buttons FeatureModule builder
    - Path: pkg/yappgen/scad/scad.go
      Note: Shared SCAD raw token type consumed by modules
ExternalSources: []
Summary: "Close the gap between the YAML DSL and SCAD-only pushButtons arrays so tactile switch caps can be described in YAML, generated via yappctl, and documented in the help system."
LastUpdated: 2025-11-15T16:40:12-05:00
---


# Implement pushButtons DSL parity

## Overview

`pushButtons` are still SCAD-only: the example YAML (`examples/yapp-demo-buttons.yaml`) has a TODO because `pkg/yappgen` cannot build the positional array that `YAPPgenerator_v3.scad` expects. This ticket tracks the work to (1) define a DSL schema, (2) map it to the generator, (3) document it via the Glazed help system, and (4) prove parity with both `YAPP_Demo_buttons_v31.scad` and `YAPP_Demo_buttons2_v31.scad`.

The Glazed help system now carries a general DSL tutorial plus a deeper reference page, so once the push button node exists the guidance must flow into those two docs (no bespoke push-buttons page). Resolver behavior (variables + expressions) has also been documented, which means the new schema needs to harmonize with strict-mode validation and expression resolution from the start.

Downstream tooling (manual playbooks + `yappctl generate`) still relies on running the actual CLI and ensuring the resulting STL/SCAD artifacts contain push button geometry. We therefore need schema + builder code that plugs into the existing generator helpers without introducing regression risk for the rest of the DSL.

## Goals

- Capture the entire SCAD-level push button parameter set in a `features.push_buttons` schema, including enums for shape/origin/coordinate and optional presets.
- Transform the YAML nodes into the `pushButtons` array + `printSwitchExtenders` flag inside `pkg/yappgen` so `yappctl generate` emits working SCAD/STL files.
- Provide parity YAML for both SCAD examples (buttons + buttons2) and wire them into automated/manual checks.
- Extend the existing DSL tutorial/reference so end users can learn the feature from `yappctl help` (per Glazed documentation guidelines).

## Recent Progress

- Captured the parameter inventory and proposed YAML structure in `various/pushbuttons-dsl-analysis.md`.
- Replaced the old push-buttons-specific doc with two general DSL pages (`pkg/docs/tutorials/yapp-dsl-getting-started.md` and `pkg/docs/tutorials/yapp-dsl-reference.md`) that now explain resolver usage, variables, and expressions.
- Implemented the `push_buttons` schema, builder, and model wiring (including `printSwitchExtenders` toggles + OpenSCAD flags), added unit tests, and backfilled YAML examples (`examples/yapp-demo-buttons.yaml`, `examples/yapp-demo-buttons2.yaml`) that now round-trip to SCAD + STL via `yappctl generate`.
- Updated the CLI docs/playbooks so `yappctl help` surfaces the DSL material once the schema lands.
- Built a registry-driven feature plumbing layer (`pkg/yappgen/features.go`) so future DSL modules only need to register a `FeatureModule`; documented the approach in `design/feature-module-registry.md`.

## Module System Architecture Design (COMPLETED → SPUN OFF)

This ticket conducted the **design phase** for a complete module system redesign via three debate rounds. The architecture is now finalized, and **implementation has been spun off** to ticket **YAPP-MODULE-SYSTEM-001**.

### Debate Rounds (Design Documentation)
1. **Round 1: Validation + Parsing Flow** (`reference/debate-round-01-validation-parsing-flow.md`)
   - How modules define validation logic once for YAML ingestion, expression evaluation, and typed struct creation
   - Consensus: YAML schemas with code generation, two-phase validation
   
2. **Round 2: Registration Mechanics** (`reference/debate-round-02-registration-mechanics.md`)
   - How modules register with the central registry without touching core files
   - Consensus: Discovery-based registration via `schemagen discover`
   
3. **Round 3: Prototype Implementation** (`reference/debate-round-03-prototype-implementation.md`)
   - Hands-on prototyping of the architecture, identifying what works and what needs refinement
   - Confirmed: YAML-only approach (no hybrid), 6-week implementation timeline

### Architecture Summary (Final Consensus)
- **All modules use YAML schemas** (no hybrid with struct tags)
- **Code generation** via `schemagen` tool (generates structs, tests, registry)
- **Top-level `pkg/registry` package** breaks circular dependencies
- **Two-phase validation** (structure before resolution, constraints after)
- **Auto-generated documentation** from schemas
- **Discovery-based registration** (no core file edits to add modules)

### Implementation Status

**Implementation work moved to:** YAPP-MODULE-SYSTEM-001
- See that ticket for roadmap, tasks, and implementation guide
- This ticket retains design debates as documentation of architectural decisions

## Current Focus (Push Buttons Feature Work)

## Open Questions

- Do we allow custom polygon point lists, or restrict the MVP to existing SCAD presets (arrow, triangle, etc.)?
- How should the YAML expose shape presets (e.g., `shape: polygon` + `preset: arrow` vs a single enum entry)?
- Which defaults should we hard-code for lid wall thickness, plate thickness, and slack to mirror SCAD behavior without overwhelming the YAML?
- Should coordinate/origin flags follow the same naming as cutouts/connectors (`coordinate: pcb`, `origin: global`) or adopt clearer names?

See `various/pushbuttons-dsl-analysis.md` for the current field inventory and proposed YAML structure.

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- yapp
- cli
- dsl

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- design/ - Architecture and design documents
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
