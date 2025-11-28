---
Title: YAML comment provenance
Ticket: YAPP-SCAD-COMMENTS-001
Status: active
Topics:
    - yapp
    - scad
    - codegen
    - dx
DocType: analysis
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Plan for capturing YAML AST comments and threading them into SCAD provenance output."
LastUpdated: 2025-11-28T17:05:00-05:00
---

# YAML comment provenance

## 1. Goal

Carry the author’s YAML comments (inline, head, foot) all the way through `yappctl generate` so SCAD files show both the original notes and the computed provenance. We no longer need backwards compatibility with the old pure-`map[string]any` flow, so we can redesign around `yaml.Node`.

## 2. Where code will change

| Area | Files | Why |
| --- | --- | --- |
| YAML loader | `pkg/cli/resolvercli/resolver.go` | Currently uses `yaml.Unmarshal` → loses comments. Needs a `yaml.Node` walker that produces both the numeric map and a `map[path][]string` of comments. |
| Resolver core & trace | `pkg/resolver/resolver.go`, `pkg/resolver/trace.go` | `ResolveResult` must accept the comment map, keep it aligned with the `TraceEntry` paths, and return it (new `Comments` field). |
| CLI plumbing | `cmd/yappctl/generate_command.go`, `pkg/cli/generatorcli/generator.go` | Need to pass the comment map through to the generator alongside the trace. |
| Provenance/model | `pkg/yappgen/provenance.go`, `pkg/yappgen/model.go`, `pkg/yappgen/features.go`, `pkg/yappgen/emit.go` | Extend `Provenance` to accept comment lookups, inject YAML comments before the generated provenance lines for scalars and feature rows. |
| Tests/fixtures | `pkg/resolver/trace_test.go`, `pkg/yappgen/emit_comments_test.go`, new sample YAML | Verify comments persist and appear in SCAD output. |

## 3. Clean implementation approach

### 3.1 Parse YAML into `yaml.Node`

1. Replace `yaml.Unmarshal` with `yaml.Decoder` into a root `yaml.Node`.
2. Walk the AST once to:
   - Build the `map[string]any` structure (`mapFromNode` helper).
   - Collect comment strings (`HeadComment`, `LineComment`, `FootComment`) into `map[string][]string`, keyed by `features.push_buttons.0.x` style paths.
   - Optionally capture `node.Line` for future use (line numbers in SCAD comments).

### 3.2 Resolver + trace changes

1. Update `resolver.Result` to include `Comments map[string][]string`.
2. Pass the comment map into `traceRecorder` so we can easily surface YAML comments next to computed expressions later (either by merging into `TraceEntry` or by storing separately and letting provenance merge them).
3. Ensure all recursive walkers (e.g., `resolvePass`, `collectUnresolved`) continue to work even though the original AST is no longer available—the numeric map remains identical.

### 3.3 CLI + generator plumbing

1. `resolvercli.LoadResult` should now return `(Document, Trace, Comments)`.
2. `cmd/yappctl` already consumes trace metadata; extend those call sites to forward comments to `generatorcli.WriteSCAD`.
3. Update `generatorcli.WriteSCAD` signature to accept the new comment map and pass it into `yappgen.BuildModel`.

### 3.4 Provenance upgrades

1. Extend `NewProvenance` to accept the comment map.
2. Update `ScalarComment` and `DescribeFeatureRow` to emit YAML comments first, e.g.:

```
// # front button nearest the display
// pushButtons[0] ← features.push_buttons[0]
//   x ...
```

3. Deduplicate so a YAML comment tied to a scalar only prints once even if consumed multiple times.

### 3.5 Testing strategy

1. Add a YAML fixture with inline comments around globals and feature entries.
2. Extend `pkg/resolver/trace_test.go` to assert `Result.Comments["pcb.length"]` etc.
3. Update `pkg/yappgen/emit_comments_test.go` to expect the YAML comment lines above the provenance block.
4. End-to-end: run `yappctl generate` on `examples/control-panel-90x70.yaml` (already full of comments once we add them) and ensure the SCAD output matches snapshots.

## 4. Next steps

1. Implement the `yaml.Node` loader + comment map builder in `resolvercli`.
2. Modify `resolver.ResolveResult` / `traceRecorder` to propagate comments.
3. Thread the comment map through CLI → generator → provenance.
4. Update emitter tests and fixtures.
5. Verify by regenerating SCAD for `control-panel-90x70`.
