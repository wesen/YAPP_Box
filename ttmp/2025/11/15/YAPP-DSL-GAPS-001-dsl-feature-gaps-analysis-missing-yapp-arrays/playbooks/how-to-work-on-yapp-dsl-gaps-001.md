You’re taking over ticket `YAPP-DSL-GAPS-001` in repo `/home/manuel/code/others/YAPP_Box`.

- Read docmgr basics first
  - Run:
    ```bash
    cd /home/manuel/code/others/YAPP_Box

    docmgr help how-to-use
    docmgr list tickets --ticket YAPP-DSL-GAPS-001
    docmgr list docs --ticket YAPP-DSL-GAPS-001
    docmgr tasks list --ticket YAPP-DSL-GAPS-001
    ```

- Read up on this ticket (in order)
  1) `/home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/index.md`
  2) Diaries:
     - `/home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/log/01-2025-11-16-implementation-diary-lighttubes-module.md`
     - `/home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/log/02-2025-11-16-implementation-diary-polygon-cutout-shapes.md`
  3) Tasks and changelog:
     - `/home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/tasks.md`
     - `/home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/changelog.md`
  4) Background authoring:
     - `/home/manuel/code/others/YAPP_Box/pkg/docs/tutorials/yapp-module-authoring-guide.md`

- Start tasks
  - List and begin with the first unchecked task:
    ```bash
    docmgr tasks list --ticket YAPP-DSL-GAPS-001
    ```
  - As you complete a task:
    ```bash
    docmgr tasks check --ticket YAPP-DSL-GAPS-001 --id <ID>
    ```

- At EVERY step: relate files and keep a changelog
  - After creating/updating files:
    ```bash
    docmgr relate --ticket YAPP-DSL-GAPS-001 \
      --file-note "/ABSOLUTE/PATH/HERE:Short note on why this file matters"

    docmgr changelog update --ticket YAPP-DSL-GAPS-001 \
      --entry "What changed and why" \
      --file-note "/ABSOLUTE/PATH/HERE:Reason"
    ```

- Keep a running implementation diary
  - Append brief entries after meaningful steps (continue in the existing diaries or create a new `various/` note). Use:
    ```text
    Write up an implementation diary for what you just did, documenting every step, what you had to do, what worked, what didn't work, what you learned, what you should do in the future,
    ```

- Repo specifics you’ll need
  - Comparison artifacts in /tmp (already generated):
    - Cutouts (polygons): `/tmp/yapp_compare/cutouts_polygons/{dsl-base.stl,dsl-lid.stl,legacy-base.stl,legacy-lid.stl}`
    - LightTubes: `/tmp/yapp_compare/lighttubes/{dsl-base.stl,dsl-lid.stl,legacy-base.stl,legacy-lid.stl}`
    - All-faces cutouts: `/tmp/yapp_compare/cutouts_all_faces/{dsl-base.stl,dsl-lid.stl}`
  - Re-generate comparisons (cut/paste; absolute paths, includes fixed automatically via sed):
    ```bash
    ROOT="/home/manuel/code/others/YAPP_Box"

    # Cutouts (polygons)
    OUT="/tmp/yapp_compare/cutouts_polygons" && mkdir -p "$OUT" && \
    go run "$ROOT/cmd/yappctl" resolve -i "$ROOT/examples/compare/cutouts_polygons.yaml" -o "$OUT/resolved.yaml" && \
    go run "$ROOT/cmd/yappctl" generate -i "$OUT/resolved.yaml" -o "$OUT/dsl.scad" && \
    sed -i 's@include <.*YAPPgenerator_v3.scad>@include <'"$ROOT"'/YAPPgenerator_v3.scad>@' "$OUT/dsl.scad" && \
    openscad -o "$OUT/dsl-base.stl" -D 'printBaseShell=true;printLidShell=false;printSwitchExtenders=false' "$OUT/dsl.scad" && \
    openscad -o "$OUT/dsl-lid.stl"  -D 'printBaseShell=false;printLidShell=true;printSwitchExtenders=false'  "$OUT/dsl.scad" && \
    openscad -o "$OUT/legacy-base.stl" -D 'printBaseShell=true;printLidShell=false;printSwitchExtenders=false' \
      "$ROOT/examples/YAPP_Compare_cutouts_polygons_v3.scad" && \
    openscad -o "$OUT/legacy-lid.stl"  -D 'printBaseShell=false;printLidShell=true;printSwitchExtenders=false' \
      "$ROOT/examples/YAPP_Compare_cutouts_polygons_v3.scad"

    # LightTubes comparison
    OUT="/tmp/yapp_compare/lighttubes" && mkdir -p "$OUT" && \
    go run "$ROOT/cmd/yappctl" resolve -i "$ROOT/examples/yapp-demo-lighttubes.yaml" -o "$OUT/resolved.yaml" && \
    go run "$ROOT/cmd/yappctl" generate -i "$OUT/resolved.yaml" -o "$OUT/dsl.scad" && \
    sed -i 's@include <.*YAPPgenerator_v3.scad>@include <'"$ROOT"'/YAPPgenerator_v3.scad>@' "$OUT/dsl.scad" && \
    openscad -o "$OUT/dsl-base.stl" -D 'printBaseShell=true;printLidShell=false;printSwitchExtenders=false' "$OUT/dsl.scad" && \
    openscad -o "$OUT/dsl-lid.stl"  -D 'printBaseShell=false;printLidShell=true;printSwitchExtenders=false'  "$OUT/dsl.scad" && \
    openscad -o "$OUT/legacy-base.stl" -D 'printBaseShell=true;printLidShell=false;printSwitchExtenders=false' \
      "$ROOT/examples/YAPP_Demo_lightTubes_v30.scad" && \
    openscad -o "$OUT/legacy-lid.stl"  -D 'printBaseShell=false;printLidShell=true;printSwitchExtenders=false' \
      "$ROOT/examples/YAPP_Demo_lightTubes_v30.scad"

    # All-faces cutouts coverage
    OUT="/tmp/yapp_compare/cutouts_all_faces" && mkdir -p "$OUT" && \
    go run "$ROOT/cmd/yappctl" resolve -i "$ROOT/examples/compare/cutouts_all_faces.yaml" -o "$OUT/resolved.yaml" && \
    go run "$ROOT/cmd/yappctl" generate -i "$OUT/resolved.yaml" -o "$OUT/dsl.scad" && \
    sed -i 's@include <.*YAPPgenerator_v3.scad>@include <'"$ROOT"'/YAPPgenerator_v3.scad>@' "$OUT/dsl.scad" && \
    openscad -o "$OUT/dsl-base.stl" -D 'printBaseShell=true;printLidShell=false;printSwitchExtenders=false' "$OUT/dsl.scad" && \
    openscad -o "$OUT/dsl-lid.stl"  -D 'printBaseShell=false;printLidShell=true;printSwitchExtenders=false'  "$OUT/dsl.scad"
    ```

- Known gaps and immediate focus
  - Side-face cutouts mapping: legacy front/back use [posy,posz]; DSL currently maps `from_back→posy`, `from_left→posz`. This is ergonomics-only but confuses examples; we added follow-up tasks to make mapping face-aware and to align YAML examples to legacy.
  - Outstanding (see tasks list for full set): cutout masks, coordinate/origin flags, labelsPlane, ridgeExt, push_buttons missing flags, validation (ValidateStructure/ValidateConstraints), and possible removal of legacy ParamSpec helpers if unused.

- Useful docmgr helpers
  ```bash
  docmgr status --summary-only
  docmgr list tickets
  docmgr meta update --ticket YAPP-DSL-GAPS-001 --field Status --value active
  ```

- Where to start next (suggested)
  - Triage and implement “Face-aware cutout mapping” + “Adjust lighttubes example cutout heights” tasks; re-run the STL comparisons and record results in the diary.
  - Then pick up cutout masks (yappMaskDef) per ticket roadmap.