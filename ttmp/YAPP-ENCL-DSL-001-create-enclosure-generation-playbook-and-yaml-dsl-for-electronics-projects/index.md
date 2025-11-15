---
Title: Create enclosure-generation playbook and YAML DSL for electronics projects
Ticket: YAPP-ENCL-DSL-001
Status: active
Topics:
    - playbook
    - electronics
    - 3d-printing
    - yapp
    - openscad
    - automation
DocType: index
Intent: long-term
Owners:
    - manuel
RelatedFiles:
    - Path: /home/manuel/code/others/YAPP_Box/YAPP_Template_v3.scad
      Note: Canonical YAPP parameter orders and defaults
    - Path: /home/manuel/code/others/YAPP_Box/YAPPgenerator_v3.scad
      Note: Upstream OpenSCAD generator include
    - Path: /home/manuel/code/others/YAPP_Box/cmd/yapp-gen/main.go
      Note: Cobra CLI for YAML→YAPP generation (uses resolver)
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/emit.go
      Note: SCAD emitter for YAPP arrays and globals
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/map.go
      Note: Schema-based mapping to positional arrays with undef
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/model.go
      Note: Generator model and normalization
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/schema.go
      Note: YAPP parameter schemas and shape/flag mappings
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/YAPP-ENCL-DSL-001-create-enclosure-generation-playbook-and-yaml-dsl-for-electronics-projects/analysis/04-mvp-path-forward-semantic-fixes-required.md
      Note: MVP plan & semantic fixes
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/YAPP-ENCL-DSL-001-create-enclosure-generation-playbook-and-yaml-dsl-for-electronics-projects/debate/05-round-5-mvp-semantic-correctness-what-must-map-cleanly.md
      Note: Semantic rationale and alternatives
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/YAPP-ENCL-DSL-001-create-enclosure-generation-playbook-and-yaml-dsl-for-electronics-projects/reference/01-enclosure-dsl-language-reference.md
      Note: DSL spec to align (remove coordinates
    - Path: YAPP_Template_v3.scad
      Note: Canonical YAPP structure with parameter orders and defaults
    - Path: cmd/encl-resolve/main.go
      Note: Glazed v0.7.0 BareCommand CLI to resolve DSL
    - Path: go.mod
      Note: Dependencies (glazed v0.7.0
    - Path: pkg/resolver/resolver.go
      Note: Fixed-point evaluator with dotted-paths and functions
    - Path: pkg/resolver/resolver_test.go
      Note: Unit tests for happy/unhappy paths
    - Path: ttmp/PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/README.md
      Note: Real-world enclosure case study and measurements
    - Path: ttmp/PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/pico_temp_monitor_base_only.scad
      Note: Base-only variant for speed when validating geometry
    - Path: ttmp/PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/pico_temp_monitor_box.scad
      Note: OpenSCAD box model to test DSL mapping
    - Path: ttmp/PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/pico_temp_monitor_lid_only.scad
      Note: Lid-only variant for fast iteration on lid features
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/analysis/README.md
      Note: Automated doc analysis context for YAPPgenerator
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/code-examples-yappcircle-yappcenter-corrections.md
      Note: Known pitfalls in examples that impact DSL guidance
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/coordinate-systems-pcb-vs-box-vs-boxinside-comparison.md
      Note: Coordinate system comparison for feature placement semantics
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/default-values-for-optional-parameters.md
      Note: Reference defaults to seed DSL default values
    - Path: ttmp/YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/pcb-stands-parameter-order-verification.md
      Note: Ordering and naming caveats for stands mapping
    - Path: ttmp/YAPP-ENCL-DSL-001-*/analysis/03-debate-synthesis-dsl-to-yapp-translation-evaluation.md
      Note: Initial debate synthesis (4 rounds)
    - Path: ttmp/YAPP-ENCL-DSL-001-*/analysis/04-mvp-path-forward-semantic-fixes-required.md
      Note: Primary implementation guide with semantic fixes and checklist
    - Path: ttmp/YAPP-ENCL-DSL-001-*/debate/05-round-5-mvp-semantic-correctness-what-must-map-cleanly.md
      Note: Round 5 debate identifying 3 critical semantic issues
    - Path: ttmp/YAPP-ENCL-DSL-001-*/reference/01-enclosure-dsl-language-reference.md
      Note: 'DSL spec (needs update: remove coordinates section)'
    - Path: ttmp/YAPP-ENCL-DSL-001-*/reference/03-debate-format-and-candidates-dsl-to-yapp-translation.md
      Note: Debate candidate profiles
    - Path: ttmp/YAPP-ENCL-DSL-001-*/reference/04-debate-questions-dsl-to-yapp-translation.md
      Note: Debate questions structure
    - Path: ttmp/YAPP-ENCL-DSL-001-.../examples/01-minimal.resolved.yaml
      Note: Resolved output
    - Path: ttmp/YAPP-ENCL-DSL-001-.../examples/01-minimal.yaml
      Note: Minimal valid config
    - Path: ttmp/YAPP-ENCL-DSL-001-.../examples/02-vars-and-expr.resolved.yaml
      Note: Resolved output
    - Path: ttmp/YAPP-ENCL-DSL-001-.../examples/02-vars-and-expr.yaml
      Note: Vars and expressions example
    - Path: ttmp/YAPP-ENCL-DSL-001-.../examples/03-functions.resolved.yaml
      Note: Resolved output
    - Path: ttmp/YAPP-ENCL-DSL-001-.../examples/03-functions.yaml
      Note: Functions (min/max/round/ceil/floor/clamp)
    - Path: ttmp/YAPP-ENCL-DSL-001-.../examples/04-features.resolved.yaml
      Note: Resolved output
    - Path: ttmp/YAPP-ENCL-DSL-001-.../examples/04-features.yaml
      Note: 'Features: holes'
    - Path: ttmp/YAPP-ENCL-DSL-001-.../examples/05-chain-deps.resolved.yaml
      Note: Resolved output
    - Path: ttmp/YAPP-ENCL-DSL-001-.../examples/05-chain-deps.yaml
      Note: Chained vars resolution
    - Path: ttmp/YAPP-ENCL-DSL-001-.../examples/90-unresolved.yaml
      Note: 'Negative case: unresolved dependency'
ExternalSources: []
Summary: Create playbook and YAML DSL to generate electronics enclosures quickly
LastUpdated: 2025-11-08T19:40:05.070921836-05:00
---



















# Create enclosure-generation playbook and YAML DSL for electronics projects

## Overview

Goal: deliver an opinionated playbook and a YAML-based DSL that maps cleanly to YAPPgenerator so engineers can generate parametric OpenSCAD enclosures quickly and repeatably.

Deliverables:
- A practical playbook for the end-to-end flow (YAML → OpenSCAD → STL)
- A first-pass DSL schema with sensible defaults
- Examples and a mapping table from DSL keys to YAPPgenerator parameters
- References to prior findings and real-world example models

## Key Links

- Playbook: [playbook/01-generate-electronics-enclosures-end-to-end-playbook.md](./playbook/01-generate-electronics-enclosures-end-to-end-playbook.md)
- Analysis: [analysis/01-analysis-docs-and-requirements-for-enclosure-dsl.md](./analysis/01-analysis-docs-and-requirements-for-enclosure-dsl.md)
- References:
  - [YAPP: Coordinate Systems — PCB vs Box vs BoxInside](../YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/coordinate-systems-pcb-vs-box-vs-boxinside-comparison.md)
  - [YAPP: Default Values for Optional Parameters](../YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/default-values-for-optional-parameters.md)
  - [YAPP: PCB Stands — Parameter Order Verification](../YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/pcb-stands-parameter-order-verification.md)
  - [YAPP: Code Examples — yappCircle/yappCenter Corrections](../YAPP-DOCS-001-fix-yappgenerator-documentation-issues-and-update-gitbook/reference/code-examples-yappcircle-yappcenter-corrections.md)
  - Case Study: [PICO-TEMP-001 README](../PICO-TEMP-001-raspberry-pi-pico-temperature-monitor-enclosure/README.md)

## Status

Current status: **active**

## Next Steps

- Finalize DSL top-level schema and defaulting rules
- Map DSL keys to YAPPgenerator parameters and functions
- Create example YAML for PICO-TEMP-001 and generate STL
- Validate dimensions/tolerances; iterate defaults and docs
- Flesh out playbook with QA checks and troubleshooting

## Topics

- playbook
- electronics
- 3d-printing
- yapp
- openscad
- automation

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
