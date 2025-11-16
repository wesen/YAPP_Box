---
Title: Resolver Debug Output Improvements
Ticket: YAPP-DEBUG-OUTPUT-001
Status: active
Topics:
    - tooling
    - yapp
    - dx
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2025-11-15T23:20:16.065977376-05:00
---


# Resolver Debug Output Improvements

## Executive Summary

`yappctl resolve` currently stops with terse errors such as:

```
unresolved expressions after 16 passes: [features.box_mounts.0.alignment ...]
```

There is no information about the literal value (e.g., `"center"`), whether the field is enum vs numeric, which upstream variables are missing, or what file/line triggered the issue. This doc captures the current debugger experience, reproduction steps, and a plan to add structured diagnostics so new contributors (and users) can identify typos or schema gaps without spelunking through resolver internals.

## Problem Statement

### Reproduction

1. Update the DSL sample `examples/yapp-demo-buttons.yaml` with new fields (`features.box_mounts`, `pcb_stands.corner`, etc.).
2. Run:
   ```
   go run ./cmd/yappctl resolve -i examples/yapp-demo-buttons.yaml -o /tmp/resolved.yaml
   ```
3. Observe failure after the resolver’s 16-pass loop; fields such as `features.box_mounts.0.alignment` are reported as unresolved even though they are literal enum strings.

### Pain Points

- **Ambiguous values:** The error only lists the JSON path; it does not echo the literal or expression.
- **No schema hinting:** Users cannot tell whether the resolver expected a number (expression) or should have treated the string as a literal enum (the real issue).
- **Missing dependency context:** When truly unresolved references occur (e.g., expression referencing `vars.foo`), the error message truncates after the path list and a `missing=map[]` stub.
- **Debug requires code dive:** The logic that decides when to treat strings as literals lives in `pkg/resolver/resolver.go` (`isStringFieldPath`, `collectUnresolved`), so interns must search source to understand why an enum wasn’t whitelisted.

## Proposed Solution

1. **Structured diagnostics object**
   - For each unresolved path, capture:
     - Path string (existing behavior)
     - Raw value snippet (first ~80 chars)
     - Field classification (`enum`, `string`, `number`, `expression`)
     - If expression: missing dependencies, evaluated error, iteration count
   - Return a JSON blob (and formatted text) so CLI/UI consumers can display richer info.

2. **Auto-classify literal enums**
   - Expand `isStringFieldPath` or, better, use schema metadata so any field declared as `enum:` or `type: string` stops the resolver from cycling.
   - Provide a warning when a field is not in the schema (typo) instead of silently attempting expression evaluation.

3. **CLI integration**
   - Add `--debug` flag to `yappctl resolve/generate` that prints the diagnostics table (path, value, reason, hints).
   - Default output should include a short summary and pointer to the debug log file in `/tmp/yappctl-debug-*.log`.

4. **Documentation**
   - Update the module authoring guide (`pkg/docs/tutorials/yapp-module-authoring-guide.md`) with a troubleshooting section referencing the new error format.
   - Add a quickstart paragraph to this ticket’s index so future contributors know where to look.

## Design Decisions

- **Schema-driven classification over hardcoded suffixes:** Instead of growing `isStringFieldPath`, reuse schemagen metadata (e.g., `ModuleSchema.Fields()`) to determine data types. As a stopgap we extended the suffix list (corner, shell_part, treatment, alignment), but this doc formalizes schema introspection as the long-term fix.
- **Single source of truth for diagnostics:** The resolver should return structured data; `yappctl` merely formats it. This avoids duplicating logic across CLI/UI tools.
- **Backwards-compatible CLI behavior:** Default error message remains brief to avoid breaking scripts, but includes “(run with --debug for details)” so advanced users can opt-in.

## Implementation Plan

1. **Diagnostics struct**
   - Define `resolver.Diagnostics` (path, valuePreview, fieldType, reason, missingDeps, suggestions).
   - Update `collectUnresolved` and `resolvePass` to populate it.

2. **Schema awareness**
   - Expose type metadata from `registry.ModuleSchema` (currently TODO) or load `schema.yaml` directly for the resolver.
   - Replace suffix heuristics with schema-derived classification; keep suffix fallback until all modules expose metadata.

3. **CLI output**
   - Teach `yappctl` to print a table when `--debug` is set; otherwise, append a short pointer:
     ```
     unresolved expressions... (see /tmp/yappctl-debug-1234.log for details)
     ```

4. **Docs & ticket updates**
   - Add reproduction examples and troubleshooting steps to `design-doc/01-resolver-debug-output-improvements.md` (this file).
   - Update ticket index/tasks as milestones are completed.

5. **Testing**
   - Add resolver unit tests covering:
     - Enum fields now correctly treated as literals.
     - Expressions referencing undefined vars include missing dependency hints.
     - CLI debug flag writes a log file with structured JSON.

## Open Questions

- Should diagnostics include full YAML snippets (risk of leaking secrets)? Proposed answer: include only single-line previews and file path references.
- Where should debug logs live? `/tmp` is fine for dev, but consider `~/.cache/yappctl/` for long-term.

## References

- `pkg/resolver/resolver.go` – `resolvePass`, `collectUnresolved`, `isStringFieldPath` (current limitations).
- `examples/yapp-demo-buttons.yaml` – real-world repro with new enum fields.
- Ticket index: `ttmp/2025/11/15/YAPP-DEBUG-OUTPUT-001-.../index.md`.
