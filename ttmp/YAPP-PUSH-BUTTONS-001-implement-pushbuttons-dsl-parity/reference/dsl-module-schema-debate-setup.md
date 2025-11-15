---
Title: DSL Module Schema Debate Setup
Ticket: YAPP-PUSH-BUTTONS-001
Status: active
Topics:
    - yapp
    - cli
    - dsl
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/docs/tutorials/yapp-dsl-module-howto.md
    - Path: pkg/yappgen/features.go
    - Path: pkg/yappgen/modules/pushbuttons/module.go
ExternalSources: []
Summary: Scaffolding for a debate about typed module structs and schema registration across the DSL
LastUpdated: 2025-11-15T17:14:43.753563032-05:00
---



# DSL Module Schema Debate Setup

## Goal

Capture the debate scaffolding for introducing typed module structs plus exported schemas across the DSL so we can evaluate competing approaches before implementing changes in `pkg/yappgen`.

## Context

We recently introduced standalone feature modules (starting with push buttons) and now want every module to define a strongly typed Go struct (consumed by the generator) together with a declarative schema artifact that powers:

- YAML validation + parsing (resolver)
- Auto-generated documentation (Glazed help system)
- Discovery/registry metadata for module authors

This debate will explore how to let modules contribute those artifacts with minimal friction while avoiding tight coupling or boilerplate.

## Candidate Profiles

| Candidate | Perspective | Key Concerns | Tooling / Evidence |
|-----------|-------------|--------------|--------------------|
| **A. The DSL Architect** — `pkg/yappgen/features.go` steward | Wants a principled registry that owns schema loading and struct instantiation | Avoiding ad-hoc hooks, ensuring compile-time safety, supporting schema references | Inspects `pkg/yappgen/features.go`, Glazed doc rendering, resolver interfaces |
| **B. Module Author** — representing future contributors | Wants a gentle authoring experience with minimal boilerplate | Ergonomics, IDE support, keeping module code in a single folder | Works from `pkg/yappgen/modules/pushbuttons`, runs `go test` |
| **C. Validation Engine** — personified `pkg/resolver` | Focuses on schema fidelity + evaluation pipeline | Needs canonical schema definitions, wants cross-module expression/var rules | Uses `pkg/resolver/resolver.go`, YAML fixtures, resolver tracing |
| **D. Documentation Librarian** — `pkg/docs/tutorials/*` maintainer | Needs rich metadata for docs/help | Ensuring module schemas include descriptions/examples, hooking into Glazed exporters | References `pkg/docs/tutorials`, `glaze help` content |
| **E. Runtime Pragmatist** — `cmd/yappctl generate` maintainer | Needs predictable structs at runtime and sane defaults | Ensuring parsed structs survive versioning, enabling feature detection in CLI | Uses `cmd/yappctl`, `pkg/cli/generatorcli`, integration playbooks |

## Debate Questions

1. **Validation + Parsing Flow:** Assuming schemas are authored in Go, how should the resolver + parser pipeline be structured so modules can define validation logic once and have it power YAML ingestion, expression evaluation, and typed struct creation without duplicating rules?
2. **Registration Mechanics:** What is the right API for modules to register their typed struct + schema with the central registry so new modules remain easy to plug in without touching core files?
3. **Prototype Round:** Given a straw-man implementation (e.g., refactoring `modules/pushbuttons` to use the new hooks), what works well, what feels brittle, and what adjustments would we make before declaring the pattern ready?
4. **Documentation + Tooling:** How should schemas expose narrative documentation (long-form descriptions, examples, resolver notes) so Glazed help pages and future tooling can surface module features automatically?

These primitives will drive the debate rounds and eventual design doc/RFC.

## Research Guide

Participants should gather concrete evidence before each round. Recommended starting points:

- **Resolver internals:** `pkg/resolver/resolver.go` and adjacent files show how expressions, strict mode flags, and schema-driven validation currently behave. Trace through `ResolveYAML` for end-to-end flows.
- **Generator modules:** `pkg/yappgen/modules/pushbuttons/` plus `pkg/yappgen/features.go` illustrate today’s module interface (config structs, builder helpers, registry wiring). Use this to anchor prototype discussions.
- **CLI entry points:** `cmd/yappctl` (and `pkg/cli/generatorcli`, `pkg/cli/resolvercli`) demonstrate how parsed structs feed SCAD/STL emission. Useful for understanding runtime expectations.
- **Documentation/export pipeline:** `pkg/docs/tutorials/yapp-dsl-module-howto.md`, `pkg/docs/tutorials/yapp-dsl-reference.md`, and `glaze help help-system` explain how schemas need to surface descriptions/examples for Glazed.
- **Playbooks & tickets:** Ticket index + design docs under `ttmp/YAPP-PUSH-BUTTONS-001-...` contain prior analysis, task breakdowns, and the feature-module registry write-up (`design/feature-module-registry.md`) for historical context.

Tip: capture commands (`rg`, `go test`, `yappctl resolve`, etc.) in each debate round so arguments remain reproducible.

## Quick Reference

- Debate framework: `/home/manuel/workspaces/.../playbook-using-debate-framework-for-technical-rfcs.md`
- Ticket index: `ttmp/YAPP-PUSH-BUTTONS-001-implement-pushbuttons-dsl-parity/index.md`
- Push button module reference implementation: `pkg/yappgen/modules/pushbuttons/`

## Usage Examples

Once rounds are drafted, link them here alongside resulting design and implementation notes.

## Related

- `pkg/docs/tutorials/yapp-dsl-module-howto.md`
- `pkg/yappgen/features.go`
