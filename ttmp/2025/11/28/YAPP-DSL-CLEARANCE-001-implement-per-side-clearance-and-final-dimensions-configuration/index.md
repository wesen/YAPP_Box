---
Title: Implement per-side clearance and final dimensions configuration
Ticket: YAPP-DSL-CLEARANCE-001
Status: active
Topics:
    - yapp
    - dsl
    - enclosure
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: examples/yapp/enclosure-dimensions/error-*.yaml
      Note: Error examples demonstrating validation failures
    - Path: examples/yapp/enclosure-dimensions/final-dimensions-example.yaml
      Note: Example demonstrating final dimensions with computed padding
    - Path: examples/yapp/enclosure-dimensions/partial-final-dimensions-example.yaml
      Note: Example demonstrating partial final dimensions
    - Path: examples/yapp/enclosure-dimensions/per-side-clearance-example.yaml
      Note: Example demonstrating asymmetric per-side clearance
    - Path: examples/yapp/enclosure-dimensions/uniform-clearance-example.yaml
      Note: Example demonstrating backward-compatible uniform clearance
    - Path: examples/yapp/enclosure-dimensions/verify_scad.py
      Note: Verification script that parses SCAD files to validate padding and shell dimension computations
    - Path: pkg/docs/tutorials/yapp-dsl-reference.md
      Note: Updated with per-side clearance and final dimensions documentation
    - Path: pkg/resolver/validation.go
      Note: Validation framework where mutual exclusivity checks will be added
    - Path: pkg/yappgen/emit.go
      Note: SCAD emission logic (no changes needed
    - Path: pkg/yappgen/model.go
      Note: |-
        Model struct and BuildModel function that need updates
        Updated Model struct with FinalLength/FinalWidth fields
    - Path: pkg/yappgen/model_test.go
      Note: New comprehensive test file covering all validation and computation scenarios
    - Path: pkg/yappgen/validation.go
      Note: New validation file with validateEnclosureDimensions function for mutual exclusivity checks
    - Path: ttmp/2025/11/28/YAPP-DSL-CLEARANCE-001-implement-per-side-clearance-and-final-dimensions-configuration/design/01-per-side-clearance-and-final-dimensions.md
      Note: Complete design specification for per-side clearance and final dimensions feature
ExternalSources: []
Summary: ""
LastUpdated: 2025-11-28T19:30:12.840832306-05:00
---





# Implement per-side clearance and final dimensions configuration

## Overview

This ticket implements two new configuration modes for the YAML DSL:

1. **Per-side clearance**: Configure `paddingFront`, `paddingBack`, `paddingLeft`, `paddingRight` independently (asymmetric padding)
2. **Final dimensions**: Specify `shellLength` and `shellWidth` directly, with padding computed backwards

These modes are mutually exclusive—specifying both results in a validation error. The implementation maintains full backward compatibility with existing `enclosure.wall.clearance` (uniform padding).

See the [design document](./design/01-per-side-clearance-and-final-dimensions.md) for complete specifications.

## Key Links

- **[Design Document](./design/01-per-side-clearance-and-final-dimensions.md)** - Complete specification for per-side clearance and final dimensions feature
- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- yapp
- dsl
- enclosure

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
