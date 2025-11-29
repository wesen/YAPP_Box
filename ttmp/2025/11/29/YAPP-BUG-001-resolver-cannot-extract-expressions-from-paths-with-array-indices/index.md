---
Title: Resolver cannot extract expressions from paths with array indices
Ticket: YAPP-BUG-001
Status: complete
Topics:
    - yapp
    - bug
    - resolver
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/resolver/resolver.go
      Note: |-
        Contains lookupPath and getPathString functions that fail on array indices
        Fixed lookupPath() to handle numeric path segments (array indices) - core bug fix
    - Path: pkg/resolver/resolver_test.go
      Note: Added integration tests for array expression extraction with missing dependencies
    - Path: pkg/resolver/rules/dependency_graph.go
      Note: DependencyGraphRule affected by this bug
ExternalSources: []
Summary: Resolver cannot extract expression strings from paths containing array indices, causing missing dependency information in error reporting
LastUpdated: 2025-11-29T12:51:01.151213047-05:00
---





# Resolver cannot extract expressions from paths with array indices

## Overview

<!-- Provide a brief overview of the ticket, its goals, and current status -->

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- yapp
- bug
- resolver

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
