#!/usr/bin/env bash
# Validation script for Task #16: pushbuttons flags (yappAltOrigin, yappPCBName)
#
# Usage: ./scripts/validate-task16-pushbuttons-flags.sh
#
# Validates that:
# 1. YAML with pcb_name resolves correctly (no resolver errors)
# 2. SCAD generation includes yappAltOrigin flag when origin: alt
# 3. SCAD generation includes yappPCBName flag when pcb_name is set
# 4. Unit tests pass

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Script is in: ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/scripts/
# Find repo root by looking for cmd/yappctl or go.mod
ROOT="$SCRIPT_DIR"
while [ ! -f "$ROOT/go.mod" ] && [ ! -d "$ROOT/cmd/yappctl" ] && [ "$ROOT" != "/" ]; do
    ROOT="$(dirname "$ROOT")"
done
if [ "$ROOT" = "/" ]; then
    echo "Error: Could not find repo root"
    exit 1
fi
YAPPCTL="$ROOT/yappctl"

echo "=== Task #16 Validation: pushbuttons flags ==="
echo ""

# 1. Build yappctl
echo "1. Building yappctl..."
(cd "$ROOT" && go build ./cmd/yappctl)
echo "✓ Built"

# 2. Test YAML resolution
echo ""
echo "2. Testing YAML resolution with pcb_name..."
cat > /tmp/test-pushbuttons.yaml << 'YAML'
project: Test Pushbuttons Flags
yapp_version: v3

pcb:
  length: 30
  width: 40
  thickness: 1.6

enclosure:
  wall:
    thickness: 2.0
    clearance: 1
  base:
    thickness: 1.0
    wall_height: 8
  lid:
    thickness: 1.0
    wall_height: 13

features:
  push_buttons:
    - name: button_with_alt_origin
      x: 10
      y: 12
      origin: alt
      cap: { length: 8, width: 6, radius: 0.5 }
      lid: { protrusion: 2.5 }
      switch: { height: 5.5, travel: 1.0, pole_diameter: 3.5 }
    
    - name: button_with_pcb_name
      x: 20
      y: 15
      pcb_name: Sensor
      cap: { length: 6, width: 6, radius: 3 }
      lid: { protrusion: 2.0 }
      switch: { height: 4.0, travel: 0.8, pole_diameter: 3.0 }
    
    - name: button_with_both_flags
      x: 30
      y: 18
      origin: alt
      pcb_name: Aux
      cap: { length: 10, width: 8, radius: 1 }
      lid: { protrusion: 3.0 }
      switch: { height: 6.0, travel: 1.2, pole_diameter: 4.0 }
YAML

"$YAPPCTL" resolve -i /tmp/test-pushbuttons.yaml -o /tmp/test-pushbuttons-resolved.yaml
echo "✓ YAML resolves successfully"

# 3. Generate SCAD and check flags
echo ""
echo "3. Generating SCAD and checking flags..."
"$YAPPCTL" generate -i /tmp/test-pushbuttons-resolved.yaml -o /tmp/test-pushbuttons.scad --copy-generator

echo ""
echo "4. Checking for yappAltOrigin flag..."
if grep -q "yappAltOrigin" /tmp/test-pushbuttons.scad; then
    echo "✓ yappAltOrigin found"
    grep "yappAltOrigin" /tmp/test-pushbuttons.scad -B 2 -A 2
else
    echo "✗ yappAltOrigin NOT FOUND"
    exit 1
fi

echo ""
echo "5. Checking for yappPCBName flag..."
if grep -q "yappPCBName" /tmp/test-pushbuttons.scad; then
    echo "✓ yappPCBName found"
    grep "yappPCBName" /tmp/test-pushbuttons.scad -B 2 -A 2
else
    echo "✗ yappPCBName NOT FOUND"
    exit 1
fi

# 6. Show full pushButtons array
echo ""
echo "6. Full pushButtons array:"
grep -A 20 "^pushButtons" /tmp/test-pushbuttons.scad | head -25

# 7. Generate STLs to verify geometry
echo ""
echo "7. Generating STLs to verify geometry..."
OUT_DIR="/tmp/validate-task16-stls"
mkdir -p "$OUT_DIR"

"$YAPPCTL" generate \
  -i /tmp/test-pushbuttons-resolved.yaml \
  -o "$OUT_DIR/test-pushbuttons.scad" \
  --stl-base "$OUT_DIR/base.stl" \
  --stl-lid "$OUT_DIR/lid.stl" \
  --copy-generator

if [ -f "$OUT_DIR/base.stl" ] && [ -f "$OUT_DIR/lid.stl" ]; then
    echo "✓ STLs generated successfully"
    ls -lh "$OUT_DIR"/*.stl
    echo ""
    echo "   Check the lid STL for button cutouts with alternate origin positioning"
    echo "   (buttons with origin: alt should be positioned differently than default)"
else
    echo "✗ STL generation failed"
    exit 1
fi

# 8. Run unit tests
echo ""
echo "8. Running unit tests..."
(cd "$ROOT" && go test ./pkg/yappgen/modules/pushbuttons/... -v)

echo ""
echo "=== Validation Complete ==="
echo "✓ All checks passed"
echo ""
echo "What to look for:"
echo "  1. SCAD file should contain 'yappAltOrigin' flag for buttons with origin: alt"
echo "  2. SCAD file should contain '[yappPCBName, \"Sensor\"]' for buttons with pcb_name"
echo "  3. STL files should generate without errors"
echo "  4. Lid STL should show button cutouts (check visually in OpenSCAD or mesh viewer)"
echo ""
echo "Generated files:"
echo "  - SCAD: $OUT_DIR/test-pushbuttons.scad"
echo "  - Base STL: $OUT_DIR/base.stl"
echo "  - Lid STL: $OUT_DIR/lid.stl"
