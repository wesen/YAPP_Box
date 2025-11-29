package resolver

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/registry"
	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
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
			// Extract field path from error message if possible
			fieldPath := extractFieldPath(err.Error(), key)
			taxonomy := errorx.NewSchemaStructureTaxonomy(path, key, fieldPath, "", "", true)
			return errors.Wrapf(taxonomy, "validate %s", path)
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
			// For now, create basic taxonomy entry. TODO: Extract enum/allowed values from schema.
			taxonomy := errorx.NewSchemaConstraintTaxonomy(path, key, "", "constraint", nil, nil, nil, nil)
			return errors.Wrapf(taxonomy, "validate %s", path)
		}
	}

	return nil
}

// extractFieldPath tries to extract field path from error message.
// Schema validation errors often have format "path: message", so we extract the path part.
func extractFieldPath(errMsg, moduleName string) string {
	// Error messages from generated code often have format "features.module[idx].field: message"
	// Try to extract the field part after the last dot
	parts := strings.Split(errMsg, ":")
	if len(parts) > 0 {
		pathPart := strings.TrimSpace(parts[0])
		// Remove module prefix if present
		prefix := "features." + moduleName
		if strings.HasPrefix(pathPart, prefix) {
			fieldPart := strings.TrimPrefix(pathPart, prefix)
			// Remove array index if present (e.g., "[0]")
			fieldPart = strings.TrimPrefix(fieldPart, "[0]")
			fieldPart = strings.TrimPrefix(fieldPart, "[1]")
			fieldPart = strings.TrimPrefix(fieldPart, "[2]")
			fieldPart = strings.TrimPrefix(fieldPart, "[3]")
			fieldPart = strings.TrimPrefix(fieldPart, ".")
			return fieldPart
		}
	}
	return ""
}
