# Enclosure Dimension Configuration Examples

This directory contains examples demonstrating the three modes of enclosure dimension configuration:

1. **Uniform clearance** (backward compatible) - Single clearance value for all sides
2. **Per-side clearance** - Independent clearance values for front, back, left, right
3. **Final dimensions** - Specify target outer dimensions with automatic padding computation

## Valid Examples

### `uniform-clearance-example.yaml`
Demonstrates the traditional uniform clearance (backward compatible):
- Single clearance value: 1.5 mm applied to all sides
- All padding values are equal: `paddingFront = paddingBack = paddingLeft = paddingRight = 1.5`

**Generated SCAD:** `uniform-clearance-example.scad`
- Shows uniform padding: all sides set to 1.5 mm
- Resulting shell: 97.8 × 77.8 mm

### `per-side-clearance-example.yaml`
Demonstrates asymmetric padding with different clearance values for each side:
- Front: 2.0 mm (extra clearance for front-mounted components)
- Back: 1.5 mm
- Left: 1.0 mm
- Right: 1.0 mm

**Generated SCAD:** `per-side-clearance-example.scad`
- Shows individual padding variables: `paddingFront = 2`, `paddingBack = 1.5`, etc.

### `final-dimensions-example.yaml`
Demonstrates specifying target outer dimensions directly:
- Target length: 100.0 mm
- Target width: 80.0 mm
- Padding computed automatically: 2.6 mm per side (even distribution)

**Generated SCAD:** `final-dimensions-example.scad`
- Shows computed padding: `paddingFront = 2.6`, `paddingBack = 2.6`, etc.

### `partial-final-dimensions-example.yaml`
Demonstrates specifying only one dimension (length), with width using default clearance:
- Target length: 100.0 mm (padding computed: 2.6 mm per side)
- Width: Uses default clearance 1.0 mm

**Generated SCAD:** `partial-final-dimensions-example.scad`
- Shows mixed configuration: length uses computed padding, width uses default

## Error Examples

These files demonstrate validation errors and should **not** generate SCAD files.

### `error-conflicting-clearance-and-dimensions.yaml`
**Error:** Cannot specify both `wall.clearance` and `dimensions` simultaneously.

**Expected error message:**
```
enclosure configuration conflict: cannot specify both clearance 
(uniform or per-side) and final dimensions (enclosure.dimensions.length/width). 
Choose one mode.
```

### `error-incomplete-per-side-clearance.yaml`
**Error:** Per-side clearance requires all four sides (`front`, `back`, `left`, `right`).

**Expected error message:**
```
per-side clearance missing required keys: [right]
```

### `error-empty-dimensions.yaml`
**Error:** `dimensions` block exists but neither `length` nor `width` is specified.

**Expected error message:**
```
enclosure.dimensions specified but neither length nor width provided
```

### `error-negative-padding.yaml`
**Error:** Target final dimension is smaller than PCB + walls (would result in negative padding).

**Expected error message:**
```
final length 90.00 is too small for PCB (90.00) + walls (2.40 × 2)
```

## Usage

Generate SCAD from valid examples:
```bash
# Uniform clearance (backward compatible)
go run ./cmd/yappctl generate \
  --input examples/yapp/enclosure-dimensions/uniform-clearance-example.yaml \
  --scad-out examples/yapp/enclosure-dimensions/uniform-clearance-example.scad

# Per-side clearance
go run ./cmd/yappctl generate \
  --input examples/yapp/enclosure-dimensions/per-side-clearance-example.yaml \
  --scad-out examples/yapp/enclosure-dimensions/per-side-clearance-example.scad

# Final dimensions
go run ./cmd/yappctl generate \
  --input examples/yapp/enclosure-dimensions/final-dimensions-example.yaml \
  --scad-out examples/yapp/enclosure-dimensions/final-dimensions-example.scad
```

Test error examples (should fail with validation errors):
```bash
go run ./cmd/yappctl generate \
  --input examples/yapp/enclosure-dimensions/error-conflicting-clearance-and-dimensions.yaml \
  --scad-out /tmp/test.scad
```

## See Also

- [YAPP DSL Reference](../../../../pkg/docs/tutorials/yapp-dsl-reference.md) - Complete field reference
- [Design Document](../../../../../ttmp/2025/11/28/YAPP-DSL-CLEARANCE-001-implement-per-side-clearance-and-final-dimensions-configuration/design/01-per-side-clearance-and-final-dimensions.md) - Implementation details

