#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_ROOT="/tmp/yapp_compare"

case="${1:-all}"

run_cutouts_polygons() {
  out="${OUT_ROOT}/cutouts_polygons"
  mkdir -p "${out}"
  echo "[cutouts_polygons] Generating DSL → SCAD → STL"
  go run "${ROOT_DIR}/cmd/yappctl" resolve -i "${ROOT_DIR}/examples/compare/cutouts_polygons.yaml" -o "${out}/resolved.yaml" &&
  go run "${ROOT_DIR}/cmd/yappctl" generate -i "${out}/resolved.yaml" -o "${out}/dsl.scad" \
    --stl-base "${out}/dsl-base.stl" \
    --stl-lid "${out}/dsl-lid.stl" \
    --render-timeout 5m

  echo "[cutouts_polygons] Rendering legacy SCAD → STL (base/lid)"
  openscad -o "${out}/legacy-base.stl" \
    -D 'printBaseShell=true; printLidShell=false; printSwitchExtenders=false' \
    "${ROOT_DIR}/examples/YAPP_Compare_cutouts_polygons_v3.scad"

  openscad -o "${out}/legacy-lid.stl" \
    -D 'printBaseShell=false; printLidShell=true; printSwitchExtenders=false' \
    "${ROOT_DIR}/examples/YAPP_Compare_cutouts_polygons_v3.scad"

  echo "[cutouts_polygons] Outputs in ${out}"
}

run_lighttubes() {
  out="${OUT_ROOT}/lighttubes"
  mkdir -p "${out}"
  echo "[lighttubes] Generating DSL → SCAD → STL"
  go run "${ROOT_DIR}/cmd/yappctl" resolve -i "${ROOT_DIR}/examples/yapp-demo-lighttubes.yaml" -o "${out}/resolved.yaml" &&
  go run "${ROOT_DIR}/cmd/yappctl" generate -i "${out}/resolved.yaml" -o "${out}/dsl.scad" \
    --stl-base "${out}/dsl-base.stl" \
    --stl-lid "${out}/dsl-lid.stl" \
    --render-timeout 5m

  echo "[lighttubes] Rendering legacy SCAD → STL (base/lid)"
  openscad -o "${out}/legacy-base.stl" \
    -D 'printBaseShell=true; printLidShell=false; printSwitchExtenders=false' \
    "${ROOT_DIR}/examples/YAPP_Demo_lightTubes_v30.scad"

  openscad -o "${out}/legacy-lid.stl" \
    -D 'printBaseShell=false; printLidShell=true; printSwitchExtenders=false' \
    "${ROOT_DIR}/examples/YAPP_Demo_lightTubes_v30.scad"

  echo "[lighttubes] Outputs in ${out}"
}

if [[ "${case}" == "all" ]]; then
  run_cutouts_polygons && run_lighttubes
elif [[ "${case}" == "cutouts" || "${case}" == "cutouts_polygons" ]]; then
  run_cutouts_polygons
elif [[ "${case}" == "lighttubes" ]]; then
  run_lighttubes
else
  echo "Unknown case: ${case}"
  echo "Usage: $0 [all|cutouts|lighttubes]"
  exit 1
fi


