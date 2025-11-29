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
    - Path: pkg/cli/resolvercli/resolver.go
      Note: CLI integration passing PositionMap to resolver
    - Path: pkg/resolver/errorx/constructors.go
      Note: Taxonomy constructors with position support
    - Path: pkg/resolver/errorx/format.go
      Note: Format functions displaying positions
    - Path: pkg/resolver/errorx/taxonomy.go
      Note: Implements taxonomy schema design
    - Path: pkg/resolver/errorx/taxonomy_test.go
      Note: Unit tests for taxonomy constructors and AsTaxonomy unwrapping
    - Path: pkg/resolver/positions.go
      Note: PositionMap type for tracking line/column positions
    - Path: pkg/resolver/resolver.go
      Note: Resolver with PositionMap integration
    - Path: pkg/resolver/rules/dependency_graph.go
      Note: New DependencyGraphRule implementation showing dependency chains for missing variables
    - Path: pkg/resolver/rules/enum_suggest.go
      Note: Implements enum suggestion rule example
    - Path: pkg/resolver/rules/enum_suggest_test.go
      Note: Unit tests for enum suggestion rule
    - Path: pkg/resolver/rules/module_doc.go
      Note: |-
        ModuleDocEmbedRule that shows field tables for schema errors
        Rule that embeds module field tables in error messages
    - Path: pkg/resolver/rules/registry.go
      Note: Implements rule registry design
    - Path: pkg/resolver/rules/registry_test.go
      Note: Unit tests for registry matching and sorting
    - Path: pkg/resolver/rules/vars_scaffold.go
      Note: Implements vars scaffold rule example
    - Path: pkg/resolver/rules/vars_scaffold_test.go
      Note: Unit tests for vars scaffold rule
    - Path: pkg/resolver/rules/yaml_syntax.go
      Note: Implements YAML syntax rule example
    - Path: pkg/resolver/rules/yaml_syntax_test.go
      Note: Unit tests for YAML syntax rule
    - Path: pkg/resolver/validation.go
      Note: Validation functions with position tracking
    - Path: pkg/resolver/validation_extract.go
      Note: Extraction functions for enum values and constraints from schema errors
    - Path: ttmp/2025/11/29/YAPP-BUG-001-resolver-cannot-extract-expressions-from-paths-with-array-indices
      Note: Related bug ticket affecting DependencyGraphRule
    - Path: ttmp/2025/11/29/YAPP-BUG-001-resolver-cannot-extract-expressions-from-paths-with-array-indices/playbook/02-test-examples-cli.md
      Note: CLI test examples for DependencyGraphRule with verified output
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
