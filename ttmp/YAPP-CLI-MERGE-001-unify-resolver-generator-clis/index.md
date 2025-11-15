---
Title: Unify resolver/generator CLIs
Ticket: YAPP-CLI-MERGE-001
Status: active
Topics:
    - yapp
    - cli
    - tooling
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Plan to consolidate the expression resolver (encl-resolve) and generator (yapp-gen) CLIs into one Glazed-based tool with shared verbs/help."
LastUpdated: 2025-11-15T15:44:37-05:00
---


# Unify resolver/generator CLIs

## Overview

`encl-resolve` (Glazed) and `yapp-gen` (Cobra) now overlap heavily: both read DSL YAML, call `pkg/resolver`, and only differ in how they expose SCAD/STL generation vs. YAML outputs. This ticket evaluates how to merge them into a single modern CLI so every user gets the same help/playbooks and we stop duplicating plumbing. See `various/cli-merge-assessment.md` for the full analysis and proposed architecture.

Key ideas:
- adopt a unified Glazed root command with subcommands `resolve` and `generate`;
- factor generator logic into reusable helpers so both verbs share code;
- ship shims for existing binaries during migration.

## Key Links

- `various/cli-merge-assessment.md` — current analysis and plan
- `tasks.md` — implementation checklist
- `cmd/encl-resolve/main.go`, `cmd/yapp-gen/main.go` — existing CLIs to merge
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- yapp
- cli
- tooling

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
