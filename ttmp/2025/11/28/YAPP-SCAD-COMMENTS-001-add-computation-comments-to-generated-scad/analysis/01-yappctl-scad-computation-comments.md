---
Title: YAPPCTL SCAD computation comments
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
Summary: "Explains how to capture resolver traces and surface them as readable SCAD comments per variable/feature row."
LastUpdated: 2025-11-28T16:21:14.660208436-05:00
---

# YAPPCTL SCAD computation comments

## 1. Goal and scope

- Make every emitted YAPP parameter (scalars and array rows) precede with a comment that explains **where the number came from**: original DSL path, expression string, and referenced variables/values.
- Keep the generated SCAD readable for hardware designers (short, single-line comments) while still comprehensive enough to debug calculations (list operands + resolved numbers).
- Leave YAPPgenerator untouched; all changes live inside `yappctl`, `resolver`, and `yappgen`.

## 2. Current generate flow (where to hook in)

1. `cmd/yappctl/generate_command.go` loads YAML, resolves expressions via `resolvercli.LoadAndResolve`, then delegates to `generatorcli.WriteSCAD`.
2. `pkg/cli/resolvercli/resolver.go` simply returns the resolved `map[string]any`; all expression strings are discarded once evaluated.
3. `pkg/yappgen/model.go` pulls numbers from that map into a `Model`.
4. `pkg/yappgen/emit.go` serializes the model into SCAD via helpers (`writeVarFloat`, `writeArrayDecl`), which currently have no context about original expressions.

Observation: once we leave `resolver.Resolve`, we've already lost the string/expression information we need for comments.

## 3. Where variables & expressions live today

- User-authored variables reside under `vars.*` inside the YAML (see `examples/control-panel-90x70.yaml`). During resolution they are hoisted into the evaluation environment (`buildEnv`), but the map retains the literal numbers only.
- Expressions are any string value that is **not** on a "string field" path. `resolvePass` treats them as expressions, evaluates them with `expr-lang`, then overwrites the map in-place.
- The resolver already tokenizes identifiers (see `extractVarRefs` and `findMissingDependencies` in `pkg/resolver/resolver.go`). That logic can be reused to capture dependencies without re-parsing.
- Feature data (e.g. `features.push_buttons[0].x`) remains as nested maps until `yappgen` walks them. We can deduce the original path for any emitted value because:
  - Globals map directly (e.g. `pcbLength` ← `pcb.length`).
  - Feature builders always know the index they are processing (see `pushbuttons.Build` label formatting).

## 4. Trace capture design

### 4.1 Resolver produces trace metadata

- Introduce `type TraceEntry struct { Path string; Expr string; Value any; Dependencies []string; Kind traceKind }`.
- When `resolvePass` successfully evaluates a string expression (or parses a numeric string), record a `TraceEntry` keyed by the dotted path that was just overwritten.
- Dependencies can be extracted with the existing `identRe` regexp, filtering out intrinsic function names. Additionally, dereference each dependency back to its **final numeric value** so comments can show both the symbolic reference and resolved number.
- `resolver.Resolve` should return `Result{ Doc map[string]any, Trace map[string]TraceEntry }` (or keep the old signature and add an option). `resolvercli.LoadAndResolve` would surface the trace alongside the resolved map so higher layers can opt into comments without breaking older consumers.

### 4.2 Preserve literal strings that were never expressions

- For string-typed paths (e.g. `face`, `shape`), we still want comments even though no expression existed. Emit entries with `Expr=""` and `Kind=Literal` so downstream code can still mention the source path.

## 5. Propagating traces to emission

1. Extend `yappgen.BuildModel` to accept a trace lookup (e.g. `type TraceLookup func(path string) (TraceEntry, bool)`).
2. When reading each scalar from the resolved map (`getFloat(resolved, "pcb.length")`), call the lookup to store a `ParamComment` struct on the model (e.g. `Model.Comments["pcbLength"] = traceFor("pcb.length")`).
3. For feature arrays, record the base path and index. Example: before normalizing the array in `arrayFeatureModule.Collect`, clone the items along with their `pathPrefix` (e.g. `features.push_buttons[<idx>]`). Alternatively, keep the path logic inside each module: `pushbuttons.Build` already iterates with `idx`, so we can reconstruct full paths like `features.push_buttons[2].x`.
4. Thread the trace lookup down to module builders so they can annotate each positional parameter they output with a `ParamComment`. This likely means extending `registry.ArrayDecl` to `registry.ArrayDecl{ Name string; Rows [][]any; RowComments [][]ParamComment }`.
5. Update `writeArrayDecl` to emit comment lines before each row when `RowComments` has entries. Example format:

```
// pushButtons[1] ← features.push_buttons[1] (button-right):
//   x = vars.button_row_center (62.0) + vars.button_spacing (28.0) => 90.0
//   shape = "circle" (literal)
[ 90, 45, ... ],
```

6. For globals, wrap `writeVarFloat` so it accepts `(name, value, comment)` and prints `// pcbLength (pcb.length) = 90.0 = pcb.length literal` just before the assignment.

## 6. Comment formatting guidelines

- Keep to two lines max per emitted value to avoid bloating SCAD:
  - Line 1: `// pcbLength ← pcb.length (resolved expression)`
  - Line 2 (optional): `//   deps: vars.xxx (12.0), pcb.width (70.0)`
- For array rows, reuse the existing label (e.g. `push_buttons[0] (button-left)`) so the comment remains aligned with YAML entries.
- Expression detail:
  - Show the original expression string exactly as authored.
  - Append `=` and the resolved numeric result when it differs from a literal.
  - For literals, say `literal 32.0`.

## 7. Example: push button location comments

Given the sample YAML (three push buttons) the trace table would capture:

| Path | Expr | Value | Dependencies |
|------|------|-------|--------------|
| `features.push_buttons[0].x` | `vars.button_row_center - vars.button_spacing` | 46.0 | `vars.button_row_center`, `vars.button_spacing` |
| `features.push_buttons[0].shape` | `circle` (literal) | `"circle"` | — |

Emission snippet (target):

```
// pushButtons[0] (button-left)
//   x = vars.button_row_center (62.0) - vars.button_spacing (16.0) => 46.0
//   y = vars.button_row_y (52.0)
//   shape = "circle" (literal)
[ 46, 52, 27, 27, 13.5, ... ],
```

This gives the exact context the user requested: both the variables involved and the expression math that produced the coordinates and shape flag.

## 8. Open questions / risks

- **Backwards compatibility:** `resolvercli.LoadAndResolve` currently returns only a map; need a plan to preserve existing API (maybe add `LoadResult.Doc`). `generatorcli.WriteSCAD` will also need an optional trace input.
- **Performance:** Capturing trace data adds bookkeeping per assignment, but the document sizes are small; still, keep trace structures lightweight and avoid copying large maps.
- **Custom modules:** Some feature modules may compute derived fields (e.g., `pushButtonShapeTokens`). Decide whether comments should stop at DSL fields or include module-level derivations (probably the former for now).
- **Array metadata:** Extending `registry.ArrayDecl` touches every module; ensure tests cover each builder (pcb stands, connectors, etc.).
- **Formatting overflow:** Complex expressions might be long; consider truncating or splitting long comments to keep SCAD legible.

## 9. Next steps

1. Implement resolver trace capture and expose it via CLI helpers.
2. Extend `yappgen` data structures to carry comment metadata down to emission helpers.
3. Prototype comment output for globals, then for one feature module (push buttons) before rolling out to others.
4. Add golden-file tests that assert the new comments for representative DSL documents.

Once these pieces are in place, we can iterate on wording/styling of the comments without touching the resolver again.
