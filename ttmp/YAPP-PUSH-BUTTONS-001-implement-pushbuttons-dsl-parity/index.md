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
ExternalSources: []
Summary: "Close the gap between the YAML DSL and SCAD-only pushButtons arrays so tactile switch caps can be described in YAML, generated via yappctl, and documented in the help system."
LastUpdated: 2025-11-15T16:06:38.651474377-05:00
---


# Implement pushButtons DSL parity

## Overview

`pushButtons` are still SCAD-only: the example YAML (`examples/yapp-demo-buttons.yaml`) has a TODO because `pkg/yappgen` cannot build the positional array that `YAPPgenerator_v3.scad` expects. This ticket tracks the work to (1) define a DSL schema, (2) map it to the generator, (3) document it via the Glazed help system, and (4) prove parity with both `YAPP_Demo_buttons_v31.scad` and `YAPP_Demo_buttons2_v31.scad`.

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
