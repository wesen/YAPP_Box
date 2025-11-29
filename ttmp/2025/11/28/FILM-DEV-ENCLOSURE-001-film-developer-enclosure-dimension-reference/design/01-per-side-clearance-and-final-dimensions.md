---
Title: Per-Side Clearance and Final Dimensions Configuration
Ticket: FILM-DEV-ENCLOSURE-001
Status: active
Topics:
    - yapp
    - enclosure
    - dsl
    - design
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/yappgen/model.go
      Note: Model struct and BuildModel function that will need updates
    - Path: pkg/yappgen/emit.go
      Note: SCAD emission logic that writes padding variables
    - Path: pkg/resolver/validation.go
      Note: Validation framework for mutual exclusivity checks
    - Path: pkg/yappgen/assets/YAPPgenerator_v3.scad
      Note: SCAD generator that consumes padding variables
ExternalSources: []
Summary: Design for adding per-side clearance and final dimension configuration to YAML DSL with mutual exclusivity validation
LastUpdated: 2025-11-28T20:00:00-05:00
---

# Per-Side Clearance and Final Dimensions Configuration

## Overview

This design adds two new configuration modes to the YAML DSL for controlling enclosure dimensions:

1. **Per-side clearance**: Set `paddingFront`, `paddingBack`, `paddingLeft`, `paddingRight` independently
2. **Final dimensions**: Set `shellLength` and `shellWidth` directly, with padding computed backwards

These modes are **mutually exclusive**—specifying both results in a validation error. This prevents ambiguity and ensures predictable behavior.

## Goals

- Enable asymmetric padding (different clearance on each side)
- Enable direct specification of final outer dimensions
- Maintain backward compatibility with existing `enclosure.wall.clearance`
- Provide clear validation errors when conflicting options are used
- Support expression evaluation in all new fields

## Non-Goals

- Supporting both modes simultaneously (would be ambiguous)
- Changing the underlying SCAD generator (only DSL changes)
- Adding per-side wall thickness (out of scope)

---

## YAML Schema Changes

### Option 1: Per-Side Clearance

Replace the single `enclosure.wall.clearance` with per-side options:

```yaml
enclosure:
  wall:
    thickness: 2.4
    # Option A: Uniform clearance (backward compatible)
    clearance: 1.5
    
    # Option B: Per-side clearance (new)
    clearance:
      front: 2.0
      back: 1.5
      left: 1.0
      right: 1.0
```

**Schema Structure:**
- `enclosure.wall.clearance` can be:
  - A number (uniform, backward compatible)
  - An object with `front`, `back`, `left`, `right` keys (new)

### Option 2: Final Dimensions

Add new top-level `enclosure.dimensions` block:

```yaml
enclosure:
  wall:
    thickness: 2.4
  dimensions:
    # Specify final outer dimensions
    length: 100.0  # shellLength
    width: 80.0    # shellWidth
```

**Schema Structure:**
- `enclosure.dimensions.length` - Final outer length (shellLength)
- `enclosure.dimensions.width` - Final outer width (shellWidth)
- Both are optional individually, but if one is specified, the other should typically be specified too (validation can warn but not error)

---

## Validation Rules

### Mutual Exclusivity

**Rule:** Cannot specify both per-side clearance AND final dimensions.

**Validation Logic:**
```go
hasPerSideClearance := hasPerSideClearanceConfig(doc)
hasFinalDimensions := hasFinalDimensionsConfig(doc)
hasUniformClearance := hasUniformClearance(doc)

if hasPerSideClearance && hasFinalDimensions {
    return error("cannot specify both per-side clearance and final dimensions")
}

if hasUniformClearance && hasFinalDimensions {
    return error("cannot specify both uniform clearance and final dimensions")
}
```

**Error Message:**
```
enclosure configuration conflict: cannot specify both per-side clearance 
(enclosure.wall.clearance.front/back/left/right) and final dimensions 
(enclosure.dimensions.length/width). Choose one mode.
```

### Partial Specification

**Per-side clearance:** All four sides must be specified if using object form:
```yaml
# Valid
clearance:
  front: 2.0
  back: 1.5
  left: 1.0
  right: 1.0

# Invalid (missing 'right')
clearance:
  front: 2.0
  back: 1.5
  left: 1.0
```

**Final dimensions:** Both length and width should be specified, but we allow partial:
```yaml
# Valid (both specified)
dimensions:
  length: 100.0
  width: 80.0

# Valid but unusual (only length specified - width computed from PCB + default padding)
dimensions:
  length: 100.0

# Invalid (neither specified - use clearance instead)
dimensions: {}
```

---

## Computation Logic

### Forward: Clearance → Final Dimensions

**Current behavior (uniform clearance):**
```
shellLength = (maxLength(pcb) + clearance + clearance) + (wallThickness × 2)
shellWidth  = (maxWidth(pcb) + clearance + clearance) + (wallThickness × 2)
```

**New behavior (per-side clearance):**
```
shellLength = (maxLength(pcb) + clearance.front + clearance.back) + (wallThickness × 2)
shellWidth  = (maxWidth(pcb) + clearance.left + clearance.right) + (wallThickness × 2)
```

### Reverse: Final Dimensions → Clearance

When final dimensions are specified, compute padding backwards:

**For length:**
```
targetShellLength = enclosure.dimensions.length
currentShellLength = (maxLength(pcb) + paddingFront + paddingBack) + (wallThickness × 2)

# Solve for paddingFront + paddingBack:
totalPaddingNeeded = targetShellLength - (maxLength(pcb) + wallThickness × 2)
```

**Distribution strategy:** Distribute evenly or allow specification:
- **Option A (simple):** Distribute evenly: `paddingFront = paddingBack = totalPaddingNeeded / 2`
- **Option B (flexible):** Allow ratio specification:
  ```yaml
  dimensions:
    length: 100.0
    padding_ratio: 0.6  # 60% front, 40% back
  ```

**Recommendation:** Start with Option A (even distribution) for simplicity. Option B can be added later if needed.

**For width:**
```
targetShellWidth = enclosure.dimensions.width
totalPaddingNeeded = targetShellWidth - (maxWidth(pcb) + wallThickness × 2)
paddingLeft = paddingRight = totalPaddingNeeded / 2
```

**Edge Cases:**
- If computed padding is negative: Error (final dimension too small for PCB + walls)
- If only one dimension specified: Compute padding for that dimension, use default clearance for the other

---

## Implementation Plan

### Phase 1: Schema Updates

**File:** `pkg/yappgen/model.go`

**Changes:**
1. Add new fields to `Model` struct:
   ```go
   type Model struct {
       // ... existing fields ...
       
       // New: Per-side padding (mutually exclusive with FinalDimensions)
       PaddingFront  float64  // >0 if per-side configured
       PaddingBack   float64  // >0 if per-side configured
       PaddingLeft   float64  // >0 if per-side configured
       PaddingRight  float64  // >0 if per-side configured
       
       // New: Final dimensions (mutually exclusive with per-side padding)
       FinalLength   float64  // >0 if final dimensions configured
       FinalWidth    float64  // >0 if final dimensions configured
   }
   ```

2. Update `BuildModel` function:
   ```go
   func BuildModel(ctx context.Context, resolved map[string]any, ...) (*Model, error) {
       m := &Model{...}
       
       // Check for mutual exclusivity
       hasPerSide := hasPerSideClearance(resolved)
       hasFinalDims := hasFinalDimensions(resolved)
       hasUniform := hasUniformClearance(resolved)
       
       if (hasPerSide || hasUniform) && hasFinalDims {
           return nil, errors.New("cannot specify both clearance and final dimensions")
       }
       
       // Process based on mode
       if hasFinalDims {
           // Compute padding from final dimensions
           m.FinalLength = getFloat(resolved, "enclosure.dimensions.length")
           m.FinalWidth = getFloat(resolved, "enclosure.dimensions.width")
           m.computePaddingFromDimensions()
       } else if hasPerSide {
           // Use per-side clearance
           m.PaddingFront = getFloat(resolved, "enclosure.wall.clearance.front")
           m.PaddingBack = getFloat(resolved, "enclosure.wall.clearance.back")
           m.PaddingLeft = getFloat(resolved, "enclosure.wall.clearance.left")
           m.PaddingRight = getFloat(resolved, "enclosure.wall.clearance.right")
       } else if hasUniform {
           // Backward compatible uniform clearance
           clearance := getFloat(resolved, "enclosure.wall.clearance")
           m.PaddingFront = clearance
           m.PaddingBack = clearance
           m.PaddingLeft = clearance
           m.PaddingRight = clearance
       }
       
       return m, nil
   }
   ```

### Phase 2: Validation

**File:** `pkg/resolver/validation.go` (or new `pkg/yappgen/validation.go`)

**Add validation function:**
```go
func validateEnclosureDimensions(doc map[string]any) error {
    enclosure, ok := doc["enclosure"].(map[string]any)
    if !ok {
        return nil
    }
    
    wall, _ := enclosure["wall"].(map[string]any)
    dimensions, hasDims := enclosure["dimensions"].(map[string]any)
    
    // Check for final dimensions
    hasFinalLength := hasDims && hasKey(dimensions, "length")
    hasFinalWidth := hasDims && hasKey(dimensions, "width")
    hasFinalDims := hasFinalLength || hasFinalWidth
    
    // Check for per-side clearance
    clearance, hasClearance := wall["clearance"]
    hasPerSide := false
    if hasClearance {
        if clearanceMap, ok := clearance.(map[string]any); ok {
            hasPerSide = hasKey(clearanceMap, "front") || 
                        hasKey(clearanceMap, "back") ||
                        hasKey(clearanceMap, "left") ||
                        hasKey(clearanceMap, "right")
        }
    }
    
    // Check for uniform clearance
    hasUniform := hasClearance && !hasPerSide
    
    // Mutual exclusivity check
    if (hasPerSide || hasUniform) && hasFinalDims {
        return errors.New("enclosure configuration conflict: cannot specify both " +
            "clearance (uniform or per-side) and final dimensions. Choose one mode.")
    }
    
    // Validate per-side completeness
    if hasPerSide {
        clearanceMap := clearance.(map[string]any)
        required := []string{"front", "back", "left", "right"}
        missing := []string{}
        for _, key := range required {
            if !hasKey(clearanceMap, key) {
                missing = append(missing, key)
            }
        }
        if len(missing) > 0 {
            return errors.Errorf("per-side clearance missing required keys: %v", missing)
        }
    }
    
    return nil
}
```

**Integration:** Call this validation in `BuildModel` before processing, or add to resolver's `validateStructure` phase.

### Phase 3: Padding Computation

**File:** `pkg/yappgen/model.go`

**Add method to compute padding from final dimensions:**
```go
func (m *Model) computePaddingFromDimensions() error {
    if m.FinalLength <= 0 && m.FinalWidth <= 0 {
        return errors.New("at least one final dimension must be specified")
    }
    
    // Compute padding for length dimension
    if m.FinalLength > 0 {
        totalPaddingNeeded := m.FinalLength - (m.PcbLength + m.WallThickness*2)
        if totalPaddingNeeded < 0 {
            return errors.Errorf("final length %.2f is too small for PCB (%.2f) + walls (%.2f × 2)",
                m.FinalLength, m.PcbLength, m.WallThickness)
        }
        // Distribute evenly
        m.PaddingFront = totalPaddingNeeded / 2
        m.PaddingBack = totalPaddingNeeded / 2
    }
    
    // Compute padding for width dimension
    if m.FinalWidth > 0 {
        totalPaddingNeeded := m.FinalWidth - (m.PcbWidth + m.WallThickness*2)
        if totalPaddingNeeded < 0 {
            return errors.Errorf("final width %.2f is too small for PCB (%.2f) + walls (%.2f × 2)",
                m.FinalWidth, m.PcbWidth, m.WallThickness)
        }
        // Distribute evenly
        m.PaddingLeft = totalPaddingNeeded / 2
        m.PaddingRight = totalPaddingNeeded / 2
    }
    
    // If only one dimension specified, use default clearance for the other
    if m.FinalLength <= 0 {
        defaultClearance := 1.0 // or from config
        m.PaddingFront = defaultClearance
        m.PaddingBack = defaultClearance
    }
    if m.FinalWidth <= 0 {
        defaultClearance := 1.0
        m.PaddingLeft = defaultClearance
        m.PaddingRight = defaultClearance
    }
    
    return nil
}
```

### Phase 4: SCAD Emission

**File:** `pkg/yappgen/emit.go`

**No changes needed** - the existing emission logic already writes `paddingFront`, `paddingBack`, `paddingLeft`, `paddingRight` individually. The new code path will just populate these values differently.

---

## Examples

### Example 1: Per-Side Clearance

```yaml
pcb:
  length: 90
  width: 70
  thickness: 1.6

enclosure:
  wall:
    thickness: 2.4
    clearance:
      front: 2.0   # Extra clearance for front-mounted components
      back: 1.5
      left: 1.0
      right: 1.0
```

**Result:**
- `paddingFront = 2.0`
- `paddingBack = 1.5`
- `paddingLeft = 1.0`
- `paddingRight = 1.0`
- `shellLength = 90 + 2.0 + 1.5 + (2.4 × 2) = 98.3`
- `shellWidth = 70 + 1.0 + 1.0 + (2.4 × 2) = 76.8`

### Example 2: Final Dimensions

```yaml
pcb:
  length: 90
  width: 70
  thickness: 1.6

enclosure:
  wall:
    thickness: 2.4
  dimensions:
    length: 100.0  # Target final length
    width: 80.0    # Target final width
```

**Computation:**
- Length padding needed: `100.0 - (90 + 2.4 × 2) = 5.2`
  - `paddingFront = 2.6`
  - `paddingBack = 2.6`
- Width padding needed: `80.0 - (70 + 2.4 × 2) = 5.2`
  - `paddingLeft = 2.6`
  - `paddingRight = 2.6`

**Result:**
- `shellLength = 100.0` ✓
- `shellWidth = 80.0` ✓

### Example 3: Partial Final Dimensions

```yaml
pcb:
  length: 90
  width: 70

enclosure:
  wall:
    thickness: 2.4
    clearance: 1.5  # Used for width (no final width specified)
  dimensions:
    length: 100.0  # Only length specified
```

**Computation:**
- Length: Compute from final dimension → `paddingFront = paddingBack = 2.6`
- Width: Use uniform clearance → `paddingLeft = paddingRight = 1.5`

**Result:**
- `shellLength = 100.0` ✓
- `shellWidth = 70 + 1.5 + 1.5 + (2.4 × 2) = 76.8`

### Example 4: Error - Conflicting Configuration

```yaml
enclosure:
  wall:
    clearance:
      front: 2.0
      back: 1.5
      left: 1.0
      right: 1.0
  dimensions:
    length: 100.0  # ERROR: Cannot specify both
    width: 80.0
```

**Error:**
```
enclosure configuration conflict: cannot specify both per-side clearance 
(enclosure.wall.clearance.front/back/left/right) and final dimensions 
(enclosure.dimensions.length/width). Choose one mode.
```

---

## Backward Compatibility

### Existing YAML Files

**Current syntax (still valid):**
```yaml
enclosure:
  wall:
    clearance: 1.5  # Uniform clearance
```

**Behavior:** Unchanged - maps to all four padding variables as before.

### Migration Path

No migration required. Existing files continue to work. Users can opt into new features incrementally:

1. **Step 1:** Keep using `clearance: 1.5` (uniform)
2. **Step 2:** Switch to per-side when needed:
   ```yaml
   clearance:
     front: 2.0
     back: 1.5
     left: 1.5
     right: 1.5
   ```
3. **Step 3:** Or switch to final dimensions:
   ```yaml
   dimensions:
     length: 100.0
     width: 80.0
   ```

---

## Testing Strategy

### Unit Tests

1. **Validation tests:**
   - ✅ Per-side clearance alone (valid)
   - ✅ Final dimensions alone (valid)
   - ✅ Uniform clearance alone (valid, backward compat)
   - ❌ Per-side + final dimensions (error)
   - ❌ Uniform + final dimensions (error)
   - ❌ Incomplete per-side clearance (error)

2. **Computation tests:**
   - Forward: Per-side clearance → final dimensions
   - Reverse: Final dimensions → padding (even distribution)
   - Edge case: Negative padding (error)
   - Edge case: Partial final dimensions

3. **Emission tests:**
   - Verify SCAD output has correct padding variables
   - Verify provenance tracking

### Integration Tests

1. **End-to-end:**
   - Generate SCAD from per-side clearance YAML
   - Generate SCAD from final dimensions YAML
   - Verify OpenSCAD renders correctly

2. **Backward compatibility:**
   - Existing YAML files still work
   - Generated SCAD matches previous output

---

## Open Questions

1. **Padding distribution:** Should we support asymmetric distribution when computing from final dimensions? (e.g., 60% front, 40% back)
   - **Decision:** Start with even distribution, add ratio later if needed

2. **Partial final dimensions:** Should we allow only `length` or only `width`?
   - **Decision:** Yes, allow partial. Use default clearance for unspecified dimension.

3. **Validation phase:** Should this be in resolver (before resolution) or in BuildModel (after resolution)?
   - **Decision:** Both - structure validation in resolver, constraint validation in BuildModel

4. **Default clearance:** What should be the default when computing from partial final dimensions?
   - **Decision:** Use `1.0` as default, or allow `enclosure.wall.clearance` to provide default for unspecified dimension

---

## Implementation Checklist

- [ ] Update `Model` struct with new fields
- [ ] Add validation function for mutual exclusivity
- [ ] Update `BuildModel` to handle per-side clearance
- [ ] Update `BuildModel` to handle final dimensions
- [ ] Add `computePaddingFromDimensions` method
- [ ] Add unit tests for validation
- [ ] Add unit tests for computation
- [ ] Add integration tests
- [ ] Update DSL reference documentation
- [ ] Add examples to documentation
- [ ] Update changelog

---

## References

- **Current implementation:** `pkg/yappgen/model.go` lines 138-148
- **SCAD emission:** `pkg/yappgen/emit.go` lines 68-78
- **Validation framework:** `pkg/resolver/validation.go`
- **Analysis document:** `analysis/01-enclosure-dimension-computation-analysis.md`

