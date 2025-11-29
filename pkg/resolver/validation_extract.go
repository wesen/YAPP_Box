package resolver

import (
	"regexp"
	"strings"

	"github.com/wesen/yapp-encl-resolver/pkg/registry"
)

// extractConstraintInfo extracts enum values, min/max constraints, and actual value from error message and schema.
func extractConstraintInfo(err error, schema registry.ModuleSchema, path string, data any) (
	allowed []string,
	min, max *float64,
	actual any,
	fieldPath string,
) {
	errMsg := err.Error()
	
	// Extract field path from error message (e.g., "features.cutouts[0].face: invalid value...")
	fieldPath = extractFieldPathFromError(errMsg, path)
	
	// Extract allowed enum values from error message (pattern: "allowed: value1, value2, value3")
	allowed = extractAllowedValues(errMsg)
	
	// Extract actual value from error message (pattern: "invalid value %q" or similar)
	actual = extractActualValue(errMsg, data, fieldPath)
	
	// Query schema for min/max constraints
	if schema != nil && fieldPath != "" {
		schemaMin, schemaMax := findFieldConstraints(schema, fieldPath)
		if schemaMin != nil {
			min = schemaMin
		}
		if schemaMax != nil {
			max = schemaMax
		}
	}
	
	return allowed, min, max, actual, fieldPath
}

// extractAllowedValues parses "(allowed: value1, value2, value3)" from error message.
func extractAllowedValues(errMsg string) []string {
	// Pattern: (allowed: value1, value2, value3)
	re := regexp.MustCompile(`\(allowed:\s*([^)]+)\)`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) < 2 {
		return nil
	}
	
	allowedStr := strings.TrimSpace(matches[1])
	if allowedStr == "" {
		return nil
	}
	
	// Split by comma and trim each value
	parts := strings.Split(allowedStr, ",")
	allowed := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			allowed = append(allowed, part)
		}
	}
	
	return allowed
}

// extractActualValue tries to extract the actual value from error message or data.
func extractActualValue(errMsg string, data any, fieldPath string) any {
	// Try to extract from error message: "invalid value %q" or "invalid value %v"
	re := regexp.MustCompile(`invalid value\s+["']?([^"'\s,)]+)["']?`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) >= 2 {
		return matches[1]
	}
	
	// Fallback: try to extract from data structure if we have fieldPath
	if fieldPath != "" && data != nil {
		if value := getValueFromPath(data, fieldPath); value != nil {
			return value
		}
	}
	
	return nil
}

// extractFieldPathFromError extracts the field path from error message.
// Error format: "features.module[idx].field: message"
func extractFieldPathFromError(errMsg, modulePath string) string {
	// Split on colon to get path part
	parts := strings.SplitN(errMsg, ":", 2)
	if len(parts) == 0 {
		return ""
	}
	
	pathPart := strings.TrimSpace(parts[0])
	
	// Remove module prefix if present
	if strings.HasPrefix(pathPart, modulePath) {
		fieldPart := strings.TrimPrefix(pathPart, modulePath)
		// Remove leading dot
		fieldPart = strings.TrimPrefix(fieldPart, ".")
		return fieldPart
	}
	
	return pathPart
}

// findFieldConstraints searches schema Fields() for the given field path and returns min/max constraints.
func findFieldConstraints(schema registry.ModuleSchema, fieldPath string) (min, max *float64) {
	if schema == nil {
		return nil, nil
	}
	
	fields := schema.Fields()
	if len(fields) == 0 {
		return nil, nil
	}
	
	// fieldPath might be like "face" or "cutouts[0].face" or "connectors[0].corner"
	// We need to find the field by name, handling array indices
	fieldName := extractFieldName(fieldPath)
	
	for _, field := range fields {
		if field.Name == fieldName {
			return field.Min, field.Max
		}
		// Check children for nested fields
		if len(field.Children) > 0 {
			if min, max := findInChildren(field.Children, fieldPath); min != nil || max != nil {
				return min, max
			}
		}
	}
	
	return nil, nil
}

// extractFieldName extracts the field name from a path like "face" or "cutouts[0].face" -> "face"
func extractFieldName(fieldPath string) string {
	// Remove array indices: "cutouts[0].face" -> "cutouts.face"
	fieldPath = regexp.MustCompile(`\[\d+\]`).ReplaceAllString(fieldPath, "")
	
	// Get last component
	parts := strings.Split(fieldPath, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return fieldPath
}

// findInChildren recursively searches children for a field.
func findInChildren(children []registry.FieldSpec, fieldPath string) (min, max *float64) {
	fieldName := extractFieldName(fieldPath)
	
	for _, child := range children {
		if child.Name == fieldName {
			return child.Min, child.Max
		}
		if len(child.Children) > 0 {
			if min, max := findInChildren(child.Children, fieldPath); min != nil || max != nil {
				return min, max
			}
		}
	}
	
	return nil, nil
}

// getValueFromPath extracts a value from a nested data structure using a field path.
func getValueFromPath(data any, fieldPath string) any {
	// Simple implementation: handle basic cases like "face" or "connectors[0].corner"
	// This is a best-effort extraction
	
	// Remove array indices for now (we'd need to know the array structure)
	fieldName := extractFieldName(fieldPath)
	
	switch v := data.(type) {
	case map[string]any:
		if value, ok := v[fieldName]; ok {
			return value
		}
	case []any:
		// For arrays, try first element
		if len(v) > 0 {
			if item, ok := v[0].(map[string]any); ok {
				if value, ok := item[fieldName]; ok {
					return value
				}
			}
		}
	}
	
	return nil
}

