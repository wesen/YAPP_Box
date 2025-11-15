---
Title: CLI verb merge plan
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
      Note: Glazed CLI verb wrapping the resolver only
    - Path: cmd/yapp-gen/main.go
      Note: Cobra CLI handling SCAD/STL generation
    - Path: pkg/resolver/resolver.go
      Note: Expression resolver shared by both binaries
    - Path: pkg/yappgen
      Note: Generator utilities for SCAD emission
ExternalSources: []
Summary: "Deep-dive on encl-resolve vs yapp-gen command stacks, the required refactors, and the shape of the unified Glazed CLI without legacy shims."
LastUpdated: 2025-11-15T15:45:41.882702532-05:00
---

# CLI verb merge plan

## 1. Current verbs in `cmd/`

| Binary / Verb | Framework | Responsibilities | Docs/Help | References |
|---------------|-----------|------------------|-----------|------------|
| `cmd/encl-resolve` → `resolve` | Glazed BareCommand wrapped by Cobra | Parses DSL YAML, runs `resolver.Resolve()` and emits YAML/JSON, supports Glazed layers/options | loads `pkg/docs/tutorials` via `pkg/docs.Load()` | `cmd/encl-resolve/main.go`, `pkg/resolver/resolver.go` |
| `cmd/yapp-gen` (root command) | Plain Cobra | Parses DSL YAML, re-runs resolver, builds `pkg/yappgen` model, writes SCAD, optionally renders STL via OpenSCAD | **none** (no Glazed), manual flag plumbing | `cmd/yapp-gen/main.go`, `pkg/yappgen/*` |

Key differences:

- Framework + UX: `encl-resolve` already uses Glazed layers, structured help, tutorials, and the resolver verb is a `cmds.BareCommand`; `yapp-gen` is a standalone Cobra command lacking help integration (see `cmd/encl-resolve/main.go:13-157` vs `cmd/yapp-gen/main.go:15-157`).
- Output responsibilities: resolver only cares about YAML/JSON; generator duplicates resolver setup and layers OpenSCAD/STL emission on top (SCAD writing, relative include rewriting, STL rendering).
- Flag surface: overlap on `--strict`, but generator has `--stl-base`, `--stl-lid`, `--openscad-bin`; resolver has `--format`, `--out-file`, `--max-iterations`.
- Documentation loading: only resolver wires the Glazed help system (`pkg/docs.Load` + `help_cmd.SetupCobraRootCommand`), so generator users do not see tutorials/playbooks.

## 2. Can we merge them?

Yes—both verbs sit on the same resolver+generator graph:

1. Read YAML (`os.ReadFile`) and decode.
2. Call `resolver.Resolve` (`pkg/resolver/resolver.go:26-170`).
3. Either emit the resolved doc (`encl-resolve`) or continue with `yappgen.BuildModel` + `yappgen.EmitSCAD` (`pkg/yappgen/*.go`).
4. Optionally run OpenSCAD for STLs (generator only).

The only divergence is the post-resolver pipeline. Extracting those operations into reusable helpers lets us compose a single Glazed-based CLI (working name `yappctl`) with multiple verbs:

- `yappctl resolve`: wraps the existing `ResolveCommand` logic (same flags, same JSON/YAML support).
- `yappctl generate`: wraps the generator pipeline, reusing resolver helpers for parsing and strict handling, and exposing options for SCAD output + STL rendering.
- (optional future) `yappctl render`: guard-rail verb that forces STL rendering and ensures OpenSCAD options are set.

Because the user explicitly prohibited legacy shims, the current binaries (`encl-resolve`, `yapp-gen`) should be removed once `yappctl` lands and reaches feature parity. No compatibility wrappers are allowed.

## 3. Refactor roadmap

1. **Shared resolver I/O helper**  
   - Create `pkg/cli/resolver` (or similar) exposing `func LoadAndResolve(ctx context.Context, inPath string, opts resolver.Options) (map[string]any, error)` that encapsulates file I/O, YAML decoding, strict handling, and formatting.  
   - Both verbs call this helper to avoid duplicated file handling (currently duplicated in `cmd/encl-resolve/main.go:77-119` and `cmd/yapp-gen/main.go:54-76`).

2. **Generator service package**  
   - Extract logic from `cmd/yapp-gen/run`/`maybeRenderSTLs`/`renderSTL` into `pkg/cli/generate` (or `pkg/yappgen/cli`).  
   - Provide composable pieces:
     - `GenerateSCAD(ctx, in resolver.Document, opts GenerateOptions) (*GenerateResult, error)` returning SCAD bytes, include rewrite metadata, etc.  
     - `RenderSTL(ctx context.Context, scadPath, output string, variant STLVariant, opts RenderOptions)` to encapsulate OpenSCAD invocation and flag wiring.
   - Keep CLI-specific messaging (fmt.Printf) at the command layer, not the helper.

3. **Unified Glazed CLI (`cmd/yappctl`)**  
   - Scaffold a Cobra root using `cli.NewGlazedCommand` to load tutorials/playbooks from `pkg/docs/tutorials` exactly once.  
   - Port the existing `ResolveCommand` struct unchanged except for package paths; mount under root.  
   - Implement a new `GenerateCommand` as a Glazed `cmds.BareCommand` with flags:
     - `--input`, `--strict`, `--max-iterations` (shared with resolver).  
     - `--scad-out`, `--stl-base`, `--stl-lid`, `--openscad-bin`, `--render-timeout`, `--render-only`.  
   - Ensure command descriptions reference the tutorials and migration guidance.

4. **Retire legacy binaries**  
   - Remove `cmd/encl-resolve` and `cmd/yapp-gen` once `yappctl` has parity + tests; no shims or stub binaries per constraint.

5. **Testing + fixtures**  
   - Introduce CLI integration tests that run `yappctl resolve` and `yappctl generate` against fixtures in `examples/` and compare outputs to golden files (SCAD, resolved YAML, optional STL presence).  
   - Provide fake OpenSCAD binary or `--dry-run` flag for CI to avoid heavy renders.

6. **Docs + help updates**  
   - Update `pkg/docs/tutorials` (see `glaze help help-system`) with a playbook describing the new CLI and pointing to `ttmp/YAPP-ENCL-DSL-001.../playbook/02-yapp-documentation-study-playbook.md`.  
   - Refresh README/playbooks to reference `yappctl`.  
   - Log migration notes in `ttmp/YAPP-CLI-MERGE-001.../changelog.md`.

## 4. Work breakdown

| Step | Output | Blocking files |
|------|--------|----------------|
| Extract resolver helper | `pkg/cli/resolver/io.go` (new) | `cmd/encl-resolve/main.go`, `cmd/yapp-gen/main.go` |
| Extract generator helper | `pkg/cli/generator/service.go` (new) | `pkg/yappgen/*`, CLI runner |
| Scaffold `cmd/yappctl` | new root command w/ docs loaded | `pkg/docs`, new helper pkgs |
| Implement verbs | `resolve` + `generate` Glazed commands | helper pkgs ready |
| Delete old binaries | remove `cmd/encl-resolve`, `cmd/yapp-gen` | unified CLI stable |
| Update docs/tests | README, tutorials, integration tests | new CLI established |

## 5. Open considerations

- **Timeouts**: `yapp-gen` hardcodes 30s contexts for STL rendering; need configurable timeouts and better error messages in shared helpers.
- **SCAD include rewrites**: currently performed inline (string replace). Extract to helper to keep behavior identical.
- **Help system load order**: ensure `yappctl` root loads `pkg/docs` before adding commands so all verbs share the same tutorial set.
- **Binary name**: recommended `yappctl` to signal multi-verb tool; confirm final name with stakeholders before removing old commands.
