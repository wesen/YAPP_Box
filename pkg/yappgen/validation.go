package yappgen

import (
	"github.com/pkg/errors"
)

// validateEnclosureDimensions performs validation on enclosure dimension configuration.
// It checks for mutual exclusivity between clearance modes and final dimensions.
func validateEnclosureDimensions(resolved map[string]any) error {
	if _, ok := getMap(resolved, "enclosure"); !ok {
		return nil // No enclosure config, nothing to validate
	}

	wall, _ := getMap(resolved, "enclosure.wall")
	_, hasDims := getMap(resolved, "enclosure.dimensions")

	// Check for final dimensions
	hasFinalLength, _ := getFloat(resolved, "enclosure.dimensions.length")
	hasFinalWidth, _ := getFloat(resolved, "enclosure.dimensions.width")
	hasFinalDims := hasFinalLength > 0 || hasFinalWidth > 0

	// Check for per-side clearance
	hasPerSide := hasPerSideClearance(resolved)

	// Check for uniform clearance
	hasUniform := hasUniformClearance(resolved)

	// Mutual exclusivity check
	if (hasPerSide || hasUniform) && hasFinalDims {
		return errors.New("enclosure configuration conflict: cannot specify both clearance (uniform or per-side) and final dimensions (enclosure.dimensions.length/width). Choose one mode")
	}

	// Validate per-side completeness
	if hasPerSide {
		clearance, _ := wall["clearance"]
		clearanceMap, ok := clearance.(map[string]any)
		if !ok {
			return errors.New("invalid clearance configuration: expected map for per-side clearance")
		}
		required := []string{"front", "back", "left", "right"}
		missing := []string{}
		for _, key := range required {
			if _, ok := clearanceMap[key]; !ok {
				missing = append(missing, key)
			}
		}
		if len(missing) > 0 {
			return errors.Errorf("per-side clearance missing required keys: %v", missing)
		}
	}

	// Validate final dimensions: at least one must be specified
	if hasDims && !hasFinalDims {
		return errors.New("enclosure.dimensions specified but neither length nor width provided")
	}

	return nil
}

