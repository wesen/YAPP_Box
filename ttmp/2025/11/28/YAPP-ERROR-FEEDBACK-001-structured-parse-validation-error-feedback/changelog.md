# Changelog

## 2025-11-28

- Initial workspace created


## 2025-11-29

Implemented taxonomy schema (pkg/resolver/errorx) with stage/symptom enums and typed contexts. Updated error producers (resolvercli, resolver, validation, strict) to emit taxonomy entries. Created rule registry (pkg/resolver/rules) with three example rules (YAML syntax, enum suggest, vars scaffold). Wired into CLI with --show-taxonomy flag for raw output.

### Related Files

- cmd/yappctl/resolve_command.go — Added --show-taxonomy flag and rule rendering
- pkg/cli/resolvercli/resolver.go — Updated to emit YAML ingest taxonomy
- pkg/resolver/errorx/constructors.go — Factory functions for creating taxonomy entries
- pkg/resolver/errorx/contexts.go — Typed context structs for each stage
- pkg/resolver/errorx/format.go — Formatting functions for taxonomy display
- pkg/resolver/errorx/taxonomy.go — Core taxonomy types and enums
- pkg/resolver/resolver.go — Updated to emit expression dependency taxonomy
- pkg/resolver/rules/cli.go — CLI text renderer for rule results
- pkg/resolver/rules/default.go — Default registry with built-in rules
- pkg/resolver/rules/enum_suggest.go — Enum mismatch suggestion rule
- pkg/resolver/rules/registry.go — Rule registry with Renderer interface
- pkg/resolver/rules/vars_scaffold.go — Missing variable scaffold rule
- pkg/resolver/rules/yaml_syntax.go — YAML syntax error rule
- pkg/resolver/strict.go — Updated to emit strict-mode taxonomy
- pkg/resolver/validation.go — Updated to emit schema validation taxonomy

