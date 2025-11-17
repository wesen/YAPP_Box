---
Title: Path 1 Step B Follow-up Playbook
Ticket: YAPP-BUILDER-REFACTOR-001
Status: active
Topics:
    - yapp
    - refactor
    - onboarding
DocType: playbook
Intent: short-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Concrete next steps for continuing the ArrayDecl refactor now that Step B scaffolding is in place.
LastUpdated: 2025-11-17
---

# Step B Continuation Playbook

## What’s already done
- Step A fully validated (`go test ./...`, SCAD diffs for `examples/04-features.yaml` and `examples/yapp-mvp-pcbstands-cutouts.yaml`).
- Step B scaffolding complete:
  - `pkg/registry/schema.go` now defines `ArrayDecl` and the new `FeatureModule` contract.
  - All module builders accept typed slices and their registries return `[]ArrayDecl`.
  - `pkg/yappgen/features.go` uses the new `arrayFeatureModule`/`multiArrayFeatureModule` helpers, so cutouts no longer need a bespoke emitter.
  - Unit/integration tests updated (`go test ./...` passes).

## Next concrete tasks
1. **Task 14 – Validate Step B outputs**
   - Re-run `go test ./...` after any subsequent changes.
   - Regenerate SCAD for `examples/04-features.yaml` and `examples/yapp-mvp-pcbstands-cutouts.yaml`, diff against `/tmp/refactor-baseline`.
2. **Task 15 – Enforce no yaml.Marshal/Unmarshal**
   - Add `.golangci.yml` with `forbidigo` rules from the architecture guide.
   - Wire `golangci-lint` into CI and verify `golangci-lint run` is clean locally.
3. **Task 16 – Build wrapper tests**
   - Extend `schema_gen_test.go.tmpl` to emit `TestBuild_UsesGeneratedDecode` and `TestBuild_ReturnsArrayDecl`.
   - Regenerate modules (`go run ./cmd/schemagen discover`) and ensure the new tests pass.
4. **Task 17 – Update module authoring guide**
   - Rewrite Step 5 to show typed Build functions + registry wrappers (see architecture guide section B7).
5. **Task 18 – Migration guide**
   - Document the Path 1 (A→B) changes under `pkg/docs/migrations/path-1-a-b-refactor.md`.
6. **Tasks 19–21 – Full validation & enforcement**
   - Run the complete test suite and regenerate *all* example SCADs for parity.
   - Dry-run creation of a new module using the updated authoring guide.
   - Confirm linters/tests fail if a module skips `Decode()` or returns the wrong shape.

## Testing checklist
- `go test ./...`
- `go run ./cmd/yappctl generate -i examples/04-features.yaml -o /tmp/after-step-b-04.scad`
- `go run ./cmd/yappctl generate -i examples/yapp-mvp-pcbstands-cutouts.yaml -o /tmp/after-step-b-mvp.scad`
- `diff /tmp/refactor-baseline/baseline-04-features.scad /tmp/after-step-b-04.scad`
- `diff /tmp/refactor-baseline/baseline-mvp.scad /tmp/after-step-b-mvp.scad`
- `golangci-lint run --config .golangci.yml` (once Task 15 lands)

## Helpful references
- `design/01-architecture-implementation-guide-path-1-a-b-builder-contract-refactor.md` sections B1–B7.
- `pkg/yappgen/features.go` for array vs multi-array helper usage.
- Existing module registries under `pkg/yappgen/modules/*/registry.go` for ArrayDecl wrapper examples.

## Tips
- Regenerate modules (`go run ./cmd/schemagen discover`) after editing templates; commit generated files with the change.
- When updating docs/playbooks, remember to run `docmgr relate`/`docmgr changelog` with file notes.
- Keep `/tmp/refactor-baseline` up to date—future interns rely on those SCAD snapshots for diffs.

