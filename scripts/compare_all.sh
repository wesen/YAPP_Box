#!/usr/bin/env bash
#
# compare_all.sh - Automated comparison of DSL YAML examples against legacy SCAD
#
# Usage: ./scripts/compare_all.sh [clean]
#
# Generates DSL SCAD and STLs, compares with legacy SCAD STLs, and reports results.
# Pass "clean" argument to delete /tmp/yapp_compare before starting.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
YAPPCTL="$ROOT/yappctl"
OUT_ROOT="/tmp/yapp_compare"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if yappctl binary exists
if [ ! -f "$YAPPCTL" ]; then
    echo "Building yappctl..."
    (cd "$ROOT" && go build ./cmd/yappctl)
fi

# Clean if requested
if [ "${1:-}" = "clean" ]; then
    echo "Cleaning $OUT_ROOT..."
    rm -rf "$OUT_ROOT"
fi

echo "=========================================="
echo "YAPP DSL vs Legacy SCAD Comparison"
echo "=========================================="
echo ""

#
# Helper function to generate DSL STLs
#
generate_dsl() {
    local name=$1
    local yaml_path=$2
    local out_dir="$OUT_ROOT/$name"
    
    echo -e "${YELLOW}Generating DSL: $name${NC}"
    mkdir -p "$out_dir"
    
    # Generate SCAD and STLs (yappctl auto-copies generator when rendering STLs)
    "$YAPPCTL" generate \
        -i "$yaml_path" \
        -o "$out_dir/dsl.scad" \
        --stl-base "$out_dir/dsl-base.stl" \
        --stl-lid "$out_dir/dsl-lid.stl" \
        > "$out_dir/generate.log" 2>&1
    
    echo "  ✓ DSL SCAD: $out_dir/dsl.scad"
    echo "  ✓ DSL Base: $(du -h "$out_dir/dsl-base.stl" | cut -f1)"
    echo "  ✓ DSL Lid:  $(du -h "$out_dir/dsl-lid.stl" | cut -f1)"
}

#
# Helper function to generate legacy STLs
#
generate_legacy() {
    local name=$1
    local scad_path=$2
    local out_dir="$OUT_ROOT/$name"
    
    echo -e "${YELLOW}Generating Legacy: $name${NC}"
    
    openscad -o "$out_dir/legacy-base.stl" \
        -D 'printBaseShell=true;printLidShell=false;printSwitchExtenders=false' \
        "$scad_path" \
        > "$out_dir/legacy-base.log" 2>&1
    
    openscad -o "$out_dir/legacy-lid.stl" \
        -D 'printBaseShell=false;printLidShell=true;printSwitchExtenders=false' \
        "$scad_path" \
        > "$out_dir/legacy-lid.log" 2>&1
    
    echo "  ✓ Legacy Base: $(du -h "$out_dir/legacy-base.stl" | cut -f1)"
    echo "  ✓ Legacy Lid:  $(du -h "$out_dir/legacy-lid.stl" | cut -f1)"
    echo ""
}

#
# Helper function to compare and report
#
compare_pair() {
    local name=$1
    local out_dir="$OUT_ROOT/$name"
    
    local dsl_base_size=$(stat -c%s "$out_dir/dsl-base.stl" 2>/dev/null || echo 0)
    local dsl_lid_size=$(stat -c%s "$out_dir/dsl-lid.stl" 2>/dev/null || echo 0)
    local legacy_base_size=$(stat -c%s "$out_dir/legacy-base.stl" 2>/dev/null || echo 0)
    local legacy_lid_size=$(stat -c%s "$out_dir/legacy-lid.stl" 2>/dev/null || echo 0)
    
    echo "  Comparison for $name:"
    echo "    Base:  DSL=$(numfmt --to=iec-i --suffix=B $dsl_base_size) vs Legacy=$(numfmt --to=iec-i --suffix=B $legacy_base_size)"
    echo "    Lid:   DSL=$(numfmt --to=iec-i --suffix=B $dsl_lid_size) vs Legacy=$(numfmt --to=iec-i --suffix=B $legacy_lid_size)"
    
    # Check if sizes are close (within 10%)
    if [ $legacy_base_size -gt 0 ]; then
        local base_diff=$(( (dsl_base_size - legacy_base_size) * 100 / legacy_base_size ))
        if [ ${base_diff#-} -lt 10 ]; then
            echo -e "    ${GREEN}✓ Base sizes match (${base_diff}% diff)${NC}"
        else
            echo -e "    ${RED}⚠ Base size difference: ${base_diff}%${NC}"
        fi
    fi
    
    if [ $legacy_lid_size -gt 0 ]; then
        local lid_diff=$(( (dsl_lid_size - legacy_lid_size) * 100 / legacy_lid_size ))
        if [ ${lid_diff#-} -lt 10 ]; then
            echo -e "    ${GREEN}✓ Lid sizes match (${lid_diff}% diff)${NC}"
        else
            echo -e "    ${RED}⚠ Lid size difference: ${lid_diff}%${NC}"
        fi
    fi
    echo ""
}

#
# Comparison 1: Polygon cutouts
#
echo "=========================================="
echo "1. Polygon Cutouts (hexagon, 6pt star, arrow)"
echo "=========================================="
generate_dsl "cutouts_polygons" "$ROOT/examples/compare/cutouts_polygons.yaml"
generate_legacy "cutouts_polygons" "$ROOT/examples/YAPP_Compare_cutouts_polygons_v3.scad"
compare_pair "cutouts_polygons"

#
# Comparison 2: Light tubes
#
echo "=========================================="
echo "2. Light Tubes (LED indicators)"
echo "=========================================="
generate_dsl "lighttubes" "$ROOT/examples/yapp-demo-lighttubes.yaml"
generate_legacy "lighttubes" "$ROOT/examples/YAPP_Demo_lightTubes_v30.scad"
compare_pair "lighttubes"

#
# Comparison 3: All faces cutouts
#
echo "=========================================="
echo "3. Cutouts All Faces (coverage test)"
echo "=========================================="
generate_dsl "cutouts_all_faces" "$ROOT/examples/compare/cutouts_all_faces.yaml"
echo "  (No legacy comparison for all-faces test)"
echo ""

#
# Comparison 4: Buttons v30 demo
#
echo "=========================================="
echo "4. Buttons v30 Demo (full feature set)"
echo "=========================================="
generate_dsl "buttons_v30" "$ROOT/examples/yapp-demo-buttons-v30.yaml"
echo "  (No legacy comparison, DSL-specific example)"
echo ""

echo "=========================================="
echo "Summary"
echo "=========================================="
echo "All comparison artifacts in: $OUT_ROOT"
echo ""
echo "To inspect visually:"
echo "  - Open STL files in your 3D viewer"
echo "  - Compare dsl-*.stl vs legacy-*.stl"
echo ""
echo "Generated files:"
find "$OUT_ROOT" -name "*.stl" -o -name "*.scad" | sort

