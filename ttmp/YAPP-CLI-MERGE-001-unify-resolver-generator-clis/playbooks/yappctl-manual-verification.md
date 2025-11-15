---
Title: yappctl manual verification
Ticket: YAPP-CLI-MERGE-001
Status: active
Topics:
    - yapp
    - cli
    - tooling
DocType: playbook
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2025-11-15T15:57:29.181835352-05:00
---


# yappctl manual verification

## Purpose

Run the new `yappctl generate` verb end-to-end and confirm that both SCAD and STL outputs are produced with real geometry.

## Environment Assumptions

- Repository root checked out with `examples/yapp-demo-buttons.yaml` present.
- OpenSCAD available on `$PATH` (used for STL rendering).

## Commands

1. Generate SCAD and STL artifacts:

   ```bash
   go run ./cmd/yappctl generate \
     --input examples/yapp-demo-buttons.yaml \
     --scad-out examples/yapp-demo-buttons.scad \
     --stl-base examples/yapp-demo-buttons-base.stl \
     --stl-lid examples/yapp-demo-buttons-lid.stl
   ```

2. Verify STL files are non-empty (indicating real mesh data):

   ```bash
   stat -c '%n %s' examples/yapp-demo-buttons-*.stl
   ```

## Exit Criteria

- Command completes without errors.
- Both STL files report non-zero sizes (recent run: base ≈1.34 MB, lid ≈1.39 MB).

## Notes

- This is a temporary manual verification in lieu of full CLI integration tests; capture outputs if regressions appear.
