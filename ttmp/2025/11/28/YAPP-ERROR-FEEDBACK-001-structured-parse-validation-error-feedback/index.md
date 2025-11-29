---
Title: Structured parse/validation error feedback
Ticket: YAPP-ERROR-FEEDBACK-001
Status: active
Topics:
    - yapp
    - dx
    - errors
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: cmd/yappctl/resolve_command.go
      Note: CLI integration with --show-taxonomy flag
    - Path: pkg/resolver/errorx/taxonomy.go
      Note: Implements taxonomy schema design
    - Path: pkg/resolver/rules/enum_suggest.go
      Note: Implements enum suggestion rule example
    - Path: pkg/resolver/rules/registry.go
      Note: Implements rule registry design
    - Path: pkg/resolver/rules/vars_scaffold.go
      Note: Implements vars scaffold rule example
    - Path: pkg/resolver/rules/yaml_syntax.go
      Note: Implements YAML syntax rule example
ExternalSources: []
Summary: ""
LastUpdated: 2025-11-28T16:52:31.520526183-05:00
---


# Structured parse/validation error feedback

## Overview

<!-- Provide a brief overview of the ticket, its goals, and current status -->

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- yapp
- dx
- errors

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
