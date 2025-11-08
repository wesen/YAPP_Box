---
Title: 'Resolver CLI: Specification and Usage'
Ticket: YAPP-ENCL-DSL-001
Status: active
Topics:
    - playbook
    - electronics
    - 3d-printing
    - yapp
    - openscad
    - automation
DocType: reference
Intent: long-term
Owners:
    - manuel
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2025-11-08T16:18:37.633525415-05:00
---


# Resolver CLI: Specification and Usage

## Goal

Specify the command-line interface for a Go-based resolver that reads the Enclosure DSL YAML, evaluates variables and expressions to a fixed point, and outputs the resolved configuration.

## Context

The resolver implements deterministic evaluation: dotted-path references, numeric expressions, and function support (`min/max/round/floor/ceil/clamp`). It fails fast on invalid paths, cycles, division by zero, and non-numeric results.

## Quick Reference

Command synopsis:

```
encl-resolve --input INPUT.yaml [--output OUTPUT.yaml] [--format yaml|json] [--max-iterations 16]
```

Flags:
- --input, -i: Path to input DSL file (YAML)
- --output, -o: Path to write resolved file (defaults to stdout)
- --format, -f: Output format `yaml` (default) or `json`
- --max-iterations: Cap for fixed-point passes (default 16)
- --strict: Fail on unknown keys or unused vars (default false)
- --version: Print resolver version

Exit codes:
- 0: Success
- 2: Validation or resolution error (unresolved, cycle, invalid path)
- 3: I/O or parse error

Resolution rules:
- Numbers remain numbers; expressions are plain scalars with dotted paths.
- Evaluate until no value changes or `--max-iterations` reached.
- Detect and report cycles with a minimal dependency trace.

## Usage Examples

Resolve and write to stdout:

```bash
encl-resolve -i enclosure.yaml
```

Write JSON to file:

```bash
encl-resolve -i enclosure.yaml -o resolved.json -f json
```

Custom iteration cap:

```bash
encl-resolve -i enclosure.yaml --max-iterations 8
```

Sample output (YAML):

```yaml
pcb:
  length: 51
  width: 21
  thickness: 1.6
  z_clearance: 2
enclosure:
  wall:
    thickness: 2
    clearance: 0.5
  base:
    thickness: 2
```

## Related

- Language Reference: 01-enclosure-dsl-language-reference.md
- Playbook: ../playbook/01-generate-electronics-enclosures-end-to-end-playbook.md
