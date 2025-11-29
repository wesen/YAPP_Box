---
Title: Module Pipeline System - Multi-Array Module Support
Ticket: YAPP-MODULE-PIPELINE-001
Status: active
Topics:
    - dsl
    - modules
    - codegen
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/registry/registry.go
      Note: Global module registry implementation - module storage and lookup
    - Path: pkg/registry/schema.go
      Note: Defines FeatureModule interface and ArrayDecl type - core module system contracts
    - Path: pkg/resolver/resolver.go
      Note: DSL resolution pipeline - potential insertion point for DSL transformation
    - Path: pkg/yappgen/features.go
      Note: |-
        Wires up modules using arrayFeatureModule and multiArrayFeatureModule helpers - current module integration point
        Current module integration point - would need modification or bypass for new approaches
    - Path: pkg/yappgen/modules/cutouts/registry.go
      Note: Example of multi-array module (returns multiple ArrayDecl for different faces)
    - Path: pkg/yappgen/modules/pcbstands/registry.go
      Note: Example of single-array module (returns one ArrayDecl)
    - Path: pkg/yappgen/modules_gen.go
      Note: Generated registry that registers all modules - module discovery and registration
    - Path: ttmp/2025/11/28/YAPP-COMPOSITE-001-implement-composite-module-system-approach-3/index.md
      Note: Implementation ticket for composite module system following Approach 3 design
    - Path: ttmp/2025/11/28/YAPP-MODULE-PIPELINE-001-module-pipeline-system-multi-array-module-support/debate/debate-round-01-foundation-should-we-do-this.md
      Note: First debate round exploring whether composite modules should be implemented
    - Path: ttmp/2025/11/28/YAPP-MODULE-PIPELINE-001-module-pipeline-system-multi-array-module-support/debate/debate-round-02-architecture-which-approach.md
      Note: Second debate round comparing 5 architectural approaches with pipeline analysis and complexity metrics
    - Path: ttmp/2025/11/28/YAPP-MODULE-PIPELINE-001-module-pipeline-system-multi-array-module-support/debate/debate-round-03-pipeline-integration.md
      Note: Pipeline integration debate exploring Hook A vs Hook C with performance analysis
    - Path: ttmp/2025/11/28/YAPP-MODULE-PIPELINE-001-module-pipeline-system-multi-array-module-support/debate/debate-round-04-conflict-resolution.md
      Note: Conflict resolution debate establishing APPEND strategy and merge semantics
    - Path: ttmp/2025/11/28/YAPP-MODULE-PIPELINE-001-module-pipeline-system-multi-array-module-support/debate/debate-round-05-validation-timing.md
      Note: Validation timing debate with provenance tracking solution for error attribution
    - Path: ttmp/2025/11/28/YAPP-MODULE-PIPELINE-001-module-pipeline-system-multi-array-module-support/debate/debate-round-06-developer-experience.md
      Note: Developer experience debate with tooling requirements analysis and auto-generated documentation proposal
    - Path: ttmp/2025/11/28/YAPP-MODULE-PIPELINE-001-module-pipeline-system-multi-array-module-support/design/02-composite-module-system-design.md
      Note: Complete design document for composite module system using Approach 3 with API sketches
    - Path: ttmp/2025/11/28/YAPP-MODULE-PIPELINE-001-module-pipeline-system-multi-array-module-support/reference/01-debate-format-and-candidates.md
      Note: Cast of 11 candidates (4 human personas
    - Path: ttmp/2025/11/28/YAPP-MODULE-PIPELINE-001-module-pipeline-system-multi-array-module-support/reference/02-debate-questions.md
      Note: 10 debate questions progressing from foundation to final decision
ExternalSources: []
Summary: Design composite module system using Approach 3 (Post-Processing Array Merge) to enable modules that add entries to multiple array types (e.g., LCD with cutouts + mounting holes)
LastUpdated: 2025-11-28T21:01:52.266838041-05:00
---












# Module Pipeline System - Multi-Array Module Support

## Overview

This ticket analyzes the current YAPP DSL module system to understand how modules work and identify what's needed to support **composite modules** that can add entries to multiple different array types.

**Example Use Case:** An LCD module that needs to:
- Add a display cutout to `cutoutsFront` array
- Add mounting holes to `pcbStands` array (for mounting the LCD PCB)
- Ensure cutout and mounting holes are aligned

**Current Limitation:** Modules can only:
- Read from one DSL key (e.g., `features.pcb_stands`)
- Write to one Model field
- Output arrays of the same type (e.g., `cutouts` module outputs multiple cutout arrays for different faces, but all are cutouts)

**Goal:** Design and implement support for modules that can output to **different array types** (e.g., cutouts + pcbStands) from a single DSL entry.

## Key Findings

✅ **What Works:**
- Modules can return multiple `ArrayDecl` (see `cutouts` module)
- Modules are registered automatically via code generation
- Module system is extensible

❌ **What Doesn't Work:**
- Modules cannot add entries to different array types (e.g., cutouts + pcbStands)
- Modules cannot read from multiple DSL keys
- No way to have composite modules that combine multiple features

See [analysis/01-module-system-analysis.md](./analysis/01-module-system-analysis.md) for detailed findings.

## Design Approaches

Multiple architectural approaches have been explored through 6 rounds of debate. See [design/01-composite-module-approaches-brainstorm.md](./design/01-composite-module-approaches-brainstorm.md) for comprehensive analysis of 5 approaches.

## Design Decision

**Selected: Approach 3 (Post-Processing Array Merge)** with future migration path to Approach 1 if needed.

See [design/02-composite-module-system-design.md](./design/02-composite-module-system-design.md) for complete design.

**Key Design Decisions:**
- **Integration:** Post-processing hook in BuildModel (after regular modules)
- **Conflict Strategy:** APPEND - composite entries appended to user entries
- **Validation:** Reuse existing validators, track provenance for error attribution
- **Developer Experience:** 1 day tooling investment (auto-generated docs + helpers)

**Rationale for Approach 3:**
- ✅ Simpler (1-2 weeks vs 3-4 weeks for Approach 1)
- ✅ Sufficient for 2-5 modules
- ✅ Good DX with minimal tooling
- ✅ Clear migration path if expression generation needed

**Implementation Phases:**
1. Core infrastructure (3 days) - interfaces, registry, processor, helpers
2. Tooling and docs (2 days) - auto-generated docs, array builders
3. LCD module (2-3 days) - first composite module
4. Documentation (2-3 days) - tutorials, examples

**Total Effort:** 2-3 weeks

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- dsl
- modules
- codegen

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
