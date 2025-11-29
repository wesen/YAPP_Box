---
Title: Implement Composite Module System - Approach 3
Ticket: YAPP-COMPOSITE-001
Status: active
Topics:
    - dsl
    - modules
    - implementation
    - codegen
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/yappgen/features.go
      Note: Pattern to follow for feature module coordination
    - Path: pkg/yappgen/model.go
      Note: Integration point - BuildModel function where composite hook will be added
    - Path: ttmp/2025/11/28/YAPP-COMPOSITE-001-implement-composite-module-system-approach-3/design/01-composite-module-system-design.md
      Note: Complete design specification with API sketches and implementation phases
    - Path: ttmp/2025/11/28/YAPP-COMPOSITE-001-implement-composite-module-system-approach-3/playbooks/01-getting-started-composite-modules.md
      Note: Getting started guide for new developers with context
    - Path: ttmp/2025/11/28/YAPP-MODULE-PIPELINE-001-module-pipeline-system-multi-array-module-support/index.md
      Note: Parent ticket with complete analysis
ExternalSources: []
Summary: Implement composite module system using Approach 3 (Post-Processing Array Merge) to enable LCD and other multi-array modules
LastUpdated: 2025-11-28T21:06:04.310718937-05:00
---




# Implement Composite Module System - Approach 3

## Overview

This ticket implements the composite module system using Approach 3 (Post-Processing Array Merge) to enable DSL modules that generate entries across multiple array types.

**Example Use Case:** LCD module that generates:
- Display cutout (added to `cutoutsFront` array)
- Mounting holes (added to `pcbStands` array)

**From:**
```yaml
# ~30 lines of manual cutouts + pcb_stands configuration
```

**To:**
```yaml
features:
  lcd:
    - display:
        face: front
        position: [50, 30]
        size: [80, 40]
      mounting:
        pattern: rectangle
        spacing: [34.0, 25.0]
```

**Implementation:** 2-3 weeks total

**See:** 
- [Design Document](./design/01-composite-module-system-design.md) - Complete specification
- [Getting Started Guide](./playbooks/01-getting-started-composite-modules.md) - New developer onboarding

**Parent Ticket:** YAPP-MODULE-PIPELINE-001 (analysis, debates, decision process)

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- dsl
- modules
- implementation
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
