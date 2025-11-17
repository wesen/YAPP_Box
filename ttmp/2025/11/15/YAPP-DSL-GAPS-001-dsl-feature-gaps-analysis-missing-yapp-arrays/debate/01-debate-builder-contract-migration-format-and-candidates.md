---
Title: Debate — Builder Contract Migration: Format and Candidates
Ticket: YAPP-DSL-GAPS-001
Status: active
Topics:
    - yapp
    - architecture
    - codegen
    - debate
DocType: debate
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/analysis/04-2025-11-17-builder-contract-and-codegen-options.md
      Note: Source analysis with Options A-E
    - Path: /home/manuel/code/others/YAPP_Box/pkg/registry/schema.go
      Note: FeatureModule interface contract
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/features.go
      Note: Orchestrator with special-case handling
    - Path: /home/manuel/code/others/YAPP_Box/pkg/schemagen/
      Note: Code generator
    - Path: /home/manuel/code/others/YAPP_Box/pkg/yappgen/modules/cutouts/
      Note: Multi-array output module
ExternalSources: []
Summary: Participants, personas, and rules for Path 1 (A→B) vs Path 2 (A→C [+D]) debate on builder contract migration.
LastUpdated: 2025-11-17
---

# Debate — Builder Contract Migration: Format and Candidates

## Overview

This debate explores two migration paths for improving the YAPP DSL builder contract:

- **Path 1 (A→B):** Generate `Decode` + keep `[][]any` output → Add `ArrayDecl` IR
- **Path 2 (A→C [+D]):** Generate `Decode` + typed `Model` slices → Optional capability interfaces

**Goal:** Surface trade-offs, data-driven arguments, and pragmatic recommendations.

## Debate Format

- **Presidential debate style** with opening statements and rebuttals
- **Data-driven arguments** using actual code references, grep results, and analysis
- **Position changes allowed** when evidence contradicts assumptions
- **Moderator summaries** extract key tensions and consensus

## Candidates

### Human Developer Personas

#### Alex Chen — "The Pragmatist"
**Role:** Senior engineer focused on velocity and risk mitigation

**Philosophy:** Ship incrementally, reduce migration risk, favor simple changes over perfect architecture.

**Core concerns:**
- How many files change per step?
- Can we revert easily if something breaks?
- What's the ROI of each phase?

**Tools:** grep for usage patterns, file counts, git blame for ownership

**Personality:** Skeptical of big refactors, wants proof of value before investing time

---

#### Jordan Rivera — "The Architect"
**Role:** Staff engineer focused on long-term maintainability

**Philosophy:** Type safety enables scale; invest in structure now to prevent future pain.

**Core concerns:**
- End-to-end type safety from YAML → SCAD
- Clear contracts between layers
- Extensibility for future module types

**Tools:** Interface analysis, type flow tracing, dependency graphs

**Personality:** Patient, principled, willing to pay upfront cost for long-term gains

---

#### Sam Park — "The Codegen Maintainer"
**Role:** Owner of `pkg/schemagen/` and code generation pipeline

**Philosophy:** Balance generator complexity with output quality; templates should be maintainable.

**Core concerns:**
- Template complexity and testability
- Code generation performance
- Debugging generated code

**Tools:** Template analysis, generation benchmarks, output inspection

**Personality:** Pragmatic about codegen limits, knows when to stop generating and write by hand

---

### Code Entity Personas

#### `pkg/registry/schema.go` — "The Contract"
**Stats:** 58 lines, defines `FeatureModule` and `ModuleSchema` interfaces

**Perspective:** "I define the contract between modules and the system. Changes to my interfaces ripple everywhere."

**Wants:** Clear, stable interfaces; minimal special-cases

**Fears:** Interface proliferation, breaking changes to 7 modules

**Personality:** Principled, conservative, wants explicit over implicit

**Tools:** Can analyze interface implementations and call sites

---

#### `pkg/yappgen/features.go` — "The Orchestrator"
**Stats:** 182 lines, manages feature collection and emission

**Perspective:** "I wire everything together. Special-cases like `cutoutFeatureModule` are my burden."

**Wants:** Uniform handling for all modules; fewer custom paths

**Fears:** More special-case logic, harder debugging

**Personality:** Tired of special-cases, wants simplification

**Tools:** Can trace module registration and emission flow

---

#### `pkg/yappgen/modules/cutouts/` — "Six-Faced Friend"
**Stats:** 315 lines in `module.go`, returns `map[string][][]any` not `[][]any`

**Perspective:** "I'm special because I distribute to 6 face arrays. Don't break my output format."

**Wants:** Keep multi-array output capability

**Fears:** Being forced into single-array mold, losing face distribution

**Personality:** Defensive but pragmatic; willing to adapt if output preserved

**Tools:** Can show actual Build() signature and emission logic

---

#### `pkg/schemagen/` — "The Generator"
**Stats:** 14 files, 1200+ lines, generates structs/tests/validators

**Perspective:** "I already generate structs and validators. More generation is feasible but adds complexity."

**Wants:** Clear generation rules, testable templates

**Fears:** Template spaghetti, hard-to-debug generated code

**Personality:** Capable but cautious; prefers simple templates

**Tools:** Can analyze template complexity and generation time

---

### Wildcards

#### "The New Hire" — Casey Thompson
**Role:** Junior developer onboarding to YAPP DSL

**Perspective:** "I just want to understand how to add a module without reading 10 files."

**Wants:** Clear patterns, good examples, helpful errors

**Fears:** Cognitive overload, hidden magic, unclear contracts

**Personality:** Naive questions, exposes documentation gaps

**Tools:** Reads docs, follows examples, gets confused by inconsistency

---

## Debate Rules

1. **Lead with data:** Reference actual code, grep results, file counts
2. **Show your work:** Mention commands/queries used
3. **Adjust positions:** Change stance when evidence contradicts
4. **No hand-waving:** Specific examples over abstract principles
5. **Respect rebuttals:** Address counterarguments directly

## Question Mapping

See [Reference — Debate Questions: Builder Contract Migration](../reference/02-reference-debate-questions-builder-contract-migration.md) for the full question set and candidate assignments.

## Debate Schedule

- **Round 1:** Question 1 (Go/No-Go)
- **Round 2:** Question 2 (Which path first)
- **Round 3:** Question 3 (Developer ergonomics)
- **Round 4:** Question 5 (Special-cases)
- **Round 5:** Question 6 (Codegen scope)
- **Round 6:** Question 7 (Type-safety end-to-end)
- **Round 7:** Question 9 (Enforceability)
- **Round 8:** Question 10 (Future extensibility)

## References

- [Analysis: Builder Contract and Codegen Options](../analysis/04-2025-11-17-builder-contract-and-codegen-options.md)
- [Debate Questions Reference](../reference/02-reference-debate-questions-builder-contract-migration.md)
