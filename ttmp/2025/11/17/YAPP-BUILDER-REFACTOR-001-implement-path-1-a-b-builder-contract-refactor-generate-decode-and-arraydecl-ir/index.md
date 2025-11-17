---
Title: Implement Path 1 (A→B) Builder Contract Refactor - Generate Decode() and ArrayDecl IR
Ticket: YAPP-BUILDER-REFACTOR-001
Status: active
Topics:
    - yapp
    - refactor
    - codegen
    - architecture
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/code/others/YAPP_Box/pkg/docs/tutorials/yapp-module-authoring-guide.md
      Note: Module authoring guide to update
    - Path: /home/manuel/code/others/YAPP_Box/pkg/registry/schema.go
      Note: FeatureModule interface and ArrayDecl definition
    - Path: /home/manuel/code/others/YAPP_Box/pkg/schemagen/codegen.go
      Note: Extended schemagen data for Decode generation
    - Path: /home/manuel/code/others/YAPP_Box/pkg/schemagen/templates/schema_gen.go.tmpl
      Note: Adds Decode() function generation
    - Path: /home/manuel/code/others/YAPP_Box/pkg/schemagen/templates/schema_gen_test.go.tmpl
      Note: Adds decode tests generation
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/decode/helpers.go
      Note: New decode helpers for generated modules
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/decode/helpers_test.go
      Note: Unit tests for decode helpers
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/features.go
      Note: arrayFeatureModule and multiArrayFeatureModule helpers
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/boxmounts/module.go
      Note: Build now decodes typed items
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/connectors/module.go
      Note: Build rewired to use Decode()
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/cutouts/module.go
      Note: Build translates typed cutouts into face arrays
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/lighttubes/module.go
      Note: Build driven by Decode() values
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/pcbstands/module.go
      Note: Build uses generated Decode() output
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/pushbuttons/module.go
      Note: Dropped manual decode in favor of generated Decode()
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/snapjoins/module.go
      Note: Build switched to generated decoders
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/analysis/04-2025-11-17-builder-contract-and-codegen-options.md
      Note: Original options analysis
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/02-debate-round-1-go-no-go-builder-contract-refactor-urgency.md
      Note: Debate Round 1 - Go/No-Go decision
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/03-debate-round-2-path-choice-a-b-vs-a-c.md
      Note: Debate Round 2 - Path choice
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/design/01-architecture-implementation-guide-path-1-a-b-builder-contract-refactor.md
      Note: Complete architecture and implementation guide
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/17/YAPP-BUILDER-REFACTOR-001-implement-path-1-a-b-builder-contract-refactor-generate-decode-and-arraydecl-ir/playbook/01-intern-handoff-getting-started-with-path-1-a-b-refactor.md
      Note: Intern onboarding guide
ExternalSources: []
Summary: 'Refactor module builder contract: generate Decode() to eliminate marshal/unmarshal boilerplate, add ArrayDecl IR to unify single/multi-array handling. Single-shot refactor of 7 modules.'
LastUpdated: 2025-11-17T11:52:41.254711306-05:00
---









# Implement Path 1 (A→B) Builder Contract Refactor - Generate Decode() and ArrayDecl IR

## Overview

This ticket implements **Path 1 (A→B)** to refactor the YAPP DSL module builder contract:

- **Step A:** Generate `Decode()` functions per module (eliminate marshal/unmarshal boilerplate)
- **Step B:** Add `ArrayDecl` IR to unify single-array and multi-array module handling

**Goal:** Cleaner builder code, consistent patterns, better developer experience for module authors.

**Approach:** Single-shot refactor (all 7 modules + infrastructure in one commit).

**Parent Ticket:** [YAPP-DSL-GAPS-001](../../2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/index.md)

## For the Intern: Quick Start

**Read these documents in order:**

1. **Architecture Guide** (START HERE): [design/01-architecture-implementation-guide-path-1-a-b-builder-contract-refactor.md](../../../2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/design/01-architecture-implementation-guide-path-1-a-b-builder-contract-refactor.md)
   - Complete implementation guide with pseudocode
   - Detailed checklists for each phase
   - Links to debate decisions for rationale

2. **Options Analysis**: [analysis/04-builder-contract-and-codegen-options.md](../../../2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/analysis/04-2025-11-17-builder-contract-and-codegen-options.md)
   - Original analysis of Options A-E
   - Current state pain points
   - Why we chose Path 1 (A→B)

3. **Debate Rounds** (for decision rationale):
   - [Round 1: Go/No-Go](../../../2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/02-debate-round-1-go-no-go-builder-contract-refactor-urgency.md) - Why now
   - [Round 2: Path Choice](../../../2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/03-debate-round-2-path-choice-a-b-vs-a-c.md) - Why Path 1
   - [Round 5: Codegen Scope](../../../2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/06-debate-round-5-codegen-scope-what-to-generate-for-path-1-a-b.md) - What to generate
   - [Round 6: Enforceability](../../../2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/07-debate-round-6-enforceability-keeping-modules-honest-post-refactor.md) - How to enforce

**Then:** Follow the tasks in [tasks.md](./tasks.md) sequentially.

## What You're Implementing

**Current problem:**
```go
// Every module does this (12 error wrapping sites across 7 modules)
data, err := yaml.Marshal(it)
if err != nil {
    return nil, errors.Wrapf(err, "%s: marshal", label)
}
var item PcbStandsItem
if err := yaml.Unmarshal(data, &item); err != nil {
    return nil, errors.Wrapf(err, "%s: unmarshal", label)
}
```

**After refactor:**
```go
// Generated Decode() handles this
typed, err := Decode(items)
if err != nil {
    return nil, err
}
```

**Plus:** Unified output via ArrayDecl IR (handles both single-array and multi-array modules)

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- yapp
- refactor
- codegen
- architecture

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
