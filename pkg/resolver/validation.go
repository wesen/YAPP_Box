package resolver

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/registry"
	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
)

// validateStructure performs Phase 1 validation (before expression resolution).
// Checks field types and required fields against module schemas.
// positions is optional and can be nil if position tracking is not needed.
func validateStructure(doc map[string]any, positions PositionMap) error {
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
			line, column := 0, 0
			if positions != nil {
				line, column = positions.GetPosition(path)
			}
			taxonomy := errorx.NewSchemaStructureTaxonomy(path, key, fieldPath, "", "", true, line, column)
			return errors.Wrapf(taxonomy, "validate %s", path)
		}
	}

	return nil
}

// validateConstraints performs Phase 2 validation (after expression resolution).
// Checks min/max/enum constraints against resolved values.
// positions is optional and can be nil if position tracking is not needed.
func validateConstraints(doc map[string]any, positions PositionMap) error {
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
			// Extract enum values, min/max constraints, and actual value from error and schema
			allowed, min, max, actual, fieldPath := extractConstraintInfo(err, schema, path, value)
			
			line, column := 0, 0
			if positions != nil {
				// Try to get position for the specific field if we have fieldPath
				if fieldPath != "" {
					fieldFullPath := path
					if !strings.HasPrefix(fieldPath, "[") {
						fieldFullPath = path + "." + fieldPath
					} else {
						fieldFullPath = path + fieldPath
					}
					line, column = positions.GetPosition(fieldFullPath)
				}
				// Fallback to module path if field position not found
				if line == 0 && column == 0 {
					line, column = positions.GetPosition(path)
				}
			}
			
			// Determine expected type based on what we found
			expected := "constraint"
			if len(allowed) > 0 {
				expected = "enum"
			} else if min != nil || max != nil {
				expected = "number"
			}
			
			taxonomy := errorx.NewSchemaConstraintTaxonomy(path, key, fieldPath, expected, actual, allowed, min, max, line, column)
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
