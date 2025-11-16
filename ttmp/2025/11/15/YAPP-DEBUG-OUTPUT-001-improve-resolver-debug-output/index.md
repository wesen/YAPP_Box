---
Title: Improve resolver debug output
Ticket: YAPP-DEBUG-OUTPUT-001
Status: active
Topics:
    - tooling
    - yapp
    - dx
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: examples/yapp-demo-buttons.yaml
      Note: repro sample
    - Path: examples/yapp-demo-lighttubes-with-errors.yaml
      Note: test case demonstrating resolver error when snap_joins uses 'sides' array instead of 'side' string - useful for testing improved error messages
    - Path: pkg/resolver/resolver.go
      Note: unresolved expression diagnostics
ExternalSources: []
Summary: ""
LastUpdated: 2025-11-15T23:16:25.576380464-05:00
---




# Improve resolver debug output

## Overview

`yappctl resolve` currently stops with messages like `unresolved expressions after 16 passes: [features.box_mounts.0.alignment ...] (missing=map[..])`. The output neither echoes the raw YAML nor hints at which fields are treated as enums, making it hard to know whether the issue is a typo, a missing schema flag, or a non-numeric literal that the resolver is still trying to evaluate. This ticket tracks the work to:

- Capture the offending expression/value, its schema type (enum, string, expr), and the upstream dependency chain.
- Emit actionable suggestions (e.g., “treat as literal string because field is declared enum” or “value references vars.foo which is undefined”).
- Surface the enriched diagnostics through `yappctl` so end users aren’t forced to dig through resolver internals when schema validation fails.

## Problem Snapshot

- **Repro:** `go run ./cmd/yappctl resolve -i examples/yapp-demo-buttons.yaml`
- **Current output:** `unresolved expressions after 16 passes: [features.pcb_stands.0.corner ...]`
- **Expected:** Structured hints showing the literal input and why the resolver still thinks it is an expression.

See `pkg/resolver/resolver.go` (`collectUnresolved`, `isStringFieldPath`, `resolvePass`) for the current logic.

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- tooling
- yapp
- dx

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
