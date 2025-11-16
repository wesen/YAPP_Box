package resolver

import (
	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/registry"
)

// validateStructure performs Phase 1 validation (before expression resolution).
// Checks field types and required fields against module schemas.
func validateStructure(doc map[string]any) error {
	features, ok := doc["features"].(map[string]any)
	if !ok || features == nil {
		return nil // No features to validate
	}

	for key, value := range features {
		path := "features." + key
		module, found := registry.Get(path)
		if !found {
			// Unknown module - skip (allows forward compatibility)
			continue
		}

		schema := module.Schema()
		if err := schema.ValidateStructure(path, value); err != nil {
			return errors.Wrapf(err, "validate %s", path)
		}
	}

	return nil
}

// validateConstraints performs Phase 2 validation (after expression resolution).
// Checks min/max/enum constraints against resolved values.
func validateConstraints(doc map[string]any) error {
	features, ok := doc["features"].(map[string]any)
	if !ok || features == nil {
		return nil
	}

	for key, value := range features {
		path := "features." + key
		module, found := registry.Get(path)
		if !found {
			continue
		}

		schema := module.Schema()
		if err := schema.ValidateConstraints(path, value); err != nil {
			return errors.Wrapf(err, "validate %s", path)
		}
	}

	return nil
}
