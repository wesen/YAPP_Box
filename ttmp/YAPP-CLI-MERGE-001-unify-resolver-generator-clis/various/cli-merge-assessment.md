---
Title: CLI merge assessment
Ticket: YAPP-CLI-MERGE-001
Status: active
Topics:
    - yapp
    - cli
    - tooling
DocType: analysis
Intent: long-term
Owners: []
RelatedFiles:
    - Path: cmd/encl-resolve/main.go
      Note: Glazed CLI that resolves DSL expressions only
    - Path: cmd/yapp-gen/main.go
      Note: Cobra CLI that emits SCAD/STLs
    - Path: pkg/resolver/resolver.go
      Note: Shared resolver logic used by both CLIs
    - Path: pkg/yappgen/emit.go
      Note: Generator/SCAD emission helpers
ExternalSources: []
Summary: "Compares encl-resolve and yapp-gen, outlines risks, and proposes a unified Glazed CLI with shared verbs."
LastUpdated: 2025-11-15T15:44:37-05:00
---

# CLI merge assessment

## 1. Current State

| CLI | File | Framework | Responsibilities | Key Behaviors |
|-----|------|-----------|------------------|---------------|
| `encl-resolve resolve` | `cmd/encl-resolve/main.go` | Glazed BareCommand under Cobra | Reads DSL YAML, runs `pkg/resolver.Resolve`, writes YAML/JSON | Supports `--strict`, `--format json`, Glazed help + layered params, no SCAD/STL awareness |
| `yapp-gen` | `cmd/yapp-gen/main.go` | Plain Cobra command | Reads DSL YAML, resolves expressions, builds `pkg/yappgen` model, writes SCAD and optional STL | Manual flag plumbing, custom OpenSCAD invocation, no Glazed help, pipelines resolver twice |

### Takeaways
- Both binaries duplicate resolver setup (YAML read, call `resolver.Resolve`, handle `strict`).
- Only `encl-resolve` loads the new help playbooks from `pkg/docs/tutorials`. Generator users never see them.
- `yapp-gen` owns STL/SCAD emission and has additional flags that the resolver CLI lacks.
- CLI frameworks differ (Glazed vs raw Cobra), so merging requires choosing a single abstraction.

## 2. Requirements & Constraints

1. **Preserve existing verbs**: users rely on `encl-resolve resolve` for automation pipelines and on `yapp-gen` for geometry. We need comparable subcommands with compatible flags.
2. **Expose documentation once**: loading `pkg/docs` in a unified CLI lets every verb surface the playbooks, removing duplication.
3. **Reduce duplicated plumbing**: parameter parsing, resolver invocation, and model building should live in shared helper packages instead of separate `main.go` files.
4. **Stage gradual migration**: provide shims/aliases so scripts that call the legacy binaries continue to work while we advertise the new CLI.

## 3. Proposed Architecture

Create a single Glazed-powered root command (working name `yappctl`) with subcommands:

| Subcommand | Replacement for | Description |
|------------|-----------------|-------------|
| `resolve`  | `encl-resolve resolve` | Full parity: same parameters (`--input`, `--out-file`, `--format`, `--max-iterations`, `--strict`). |
| `generate` | `yapp-gen` | Reads DSL YAML, optionally runs resolver (or consumes `--resolved-input`), calls `pkg/yappgen`, writes SCAD, and can emit STLs via `--stl-base`/`--stl-lid`. |
| `render` (optional alias) | `yapp-gen --stl-*` | Shortcut that enforces STL output; could simply be `generate` + validation. |

### Shared Libraries
- Move generator logic from `cmd/yapp-gen/run()` into a reusable package (e.g., `pkg/cli/generator`):
  - `func Generate(ctx context.Context, opts Options) (string, []byte, error)` → returns SCAD bytes and/or writes to disk.
  - `func RenderSTLs(ctx context.Context, scadPath string, opts RenderOptions) error` → wraps OpenSCAD invocation with consistent logging.
- Add a helper in `pkg/resolver` (or new `pkg/cli/resolver`) that wraps file I/O + options to avoid duplication.

### Framework Decisions
- Base everything on Glazed command descriptions so both verbs share consistent flag names, help output, and layered configuration.
- Load `pkg/docs` once in the root to surface the documentation study playbook and any future tutorials.
- Keep Cobra as the underlying runner via `cli.BuildCobraCommand`, same as current `encl-resolve`.

## 4. Migration Strategy

1. **Introduce new CLI (`cmd/yappctl`)** with Glazed subcommands `resolve` and `generate` implemented using the refactored helper packages.
2. **Retire legacy binaries**: once the unified CLI is ready, remove `encl-resolve` and `yapp-gen` (after a brief announcement) instead of maintaining compatibility shims.
3. **Update documentation**: add migration notes to `pkg/docs/tutorials`, `README.md`, and ticket index.md so users know which verbs replace which binaries.
4. **Expand tests**: add CLI golden tests (e.g., run `yappctl resolve` on fixture YAMLs, ensure output matches historical snapshots; run `yappctl generate` and diff SCAD/undef placeholders). Reuse existing unit tests in `pkg/resolver` / `pkg/yappgen` for lower-level coverage.
5. **Release plan**: communicate timelines clearly since binaries will disappear rather than proxy to the new CLI.

## 5. Risks / Open Questions

- **Naming**: do we rename the binary or repurpose `encl-resolve`? (Recommendation: ship `yappctl` and delete the old ones once stable.)
- **Backward-compatible output**: `yapp-gen` prints SCAD path silently; Glazed may add structured output. Need to match legacy stdout/stderr so scripts parsing logs do not break.
- **OpenSCAD availability**: ensure the unified CLI still allows custom `--openscad-bin` and graceful degradation when OpenSCAD is missing.
- **Large help store**: loading all tutorials slightly increases binary size, but also improves discoverability.

## 6. Recommended Next Steps

1. Approve this plan within ticket **YAPP-CLI-MERGE-001**.
2. Create implementation tasks (see `tasks.md`) for:
   - extracting generator helpers,
   - scaffolding the unified CLI,
   - wiring doc/help loading,
   - rewriting docs to reflect removal of old binaries.
3. Prototype `yappctl resolve` + `generate` behind a feature branch; verify parity with existing binaries.
4. Iterate on feedback, then announce timeline for deprecating standalone binaries.
