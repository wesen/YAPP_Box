package rules

import (
	"context"
	"fmt"
	"strings"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
	"github.com/wesen/yapp-encl-resolver/pkg/schemagen"
)

// YamlKnownFieldsRule matches unknown key errors and suggests known field names.
// Currently works with strict mode unknown keys. Future enhancement: detect unknown fields
// in module structures when yaml.v3 KnownFields is used with struct decoding.
type YamlKnownFieldsRule struct{}

func (r *YamlKnownFieldsRule) Match(t *errorx.Taxonomy) (bool, int) {
	// Match strict mode unknown key errors
	return t.Stage == errorx.StageStrictUnknownKey && t.Symptom == errorx.SymptomUnknownKey, 85
}

func (r *YamlKnownFieldsRule) Render(ctx context.Context, t *errorx.Taxonomy) (*RuleResult, error) {
	sc, ok := t.Context.(*errorx.StrictModeContext)
	if !ok {
		return nil, fmt.Errorf("expected StrictModeContext, got %T", t.Context)
	}

	if sc.Kind != "unknown-key" {
		return nil, fmt.Errorf("expected unknown-key kind, got %s", sc.Kind)
	}

	if len(sc.Offenders) == 0 {
		return nil, fmt.Errorf("no unknown keys in context")
	}

	// For top-level unknown keys, suggest known top-level keys
	knownTopLevel := []string{"project", "version", "units", "yapp_version", "vars", "pcb", "enclosure", "features", "coordinates", "tolerances", "show"}
	
	var body strings.Builder
	body.WriteString("Unknown keys detected:\n")
	for _, key := range sc.Offenders {
		body.WriteString(fmt.Sprintf("- `%s`\n", key))
	}
	body.WriteString("\n")

	// Suggest closest known fields for each unknown key
	body.WriteString("**Did you mean:**\n")
	for _, unknownKey := range sc.Offenders {
		suggestion := nearestField(knownTopLevel, unknownKey)
		if suggestion != "" && suggestion != unknownKey {
			body.WriteString(fmt.Sprintf("- `%s` → `%s`\n", unknownKey, suggestion))
		}
	}

	// For module-level unknown fields (future enhancement)
	// This would require detecting unknown fields in module structures
	// For now, provide general guidance
	if strings.HasPrefix(t.Path, "features.") {
		moduleName := strings.TrimPrefix(t.Path, "features.")
		// Try to load schema and suggest known fields
		schemaPath := fmt.Sprintf("pkg/yappgen/modules/%s/schema.yaml", strings.ReplaceAll(moduleName, "_", ""))
		if doc, err := schemagen.LoadSchemaDoc(schemaPath); err == nil {
			knownFields := extractFieldNames(doc.Fields, "")
			if len(knownFields) > 0 {
				body.WriteString("\n**Known fields for this module:**\n")
				for _, field := range knownFields {
					body.WriteString(fmt.Sprintf("- `%s`\n", field))
				}
			}
		}
	}

	return &RuleResult{
		Headline: fmt.Sprintf("Unknown keys: %s", strings.Join(sc.Offenders, ", ")),
		Body:     body.String(),
		Severity: errorx.SeverityWarning,
		Actions: []Action{
			{Label: "Check module documentation for valid fields", Command: ""},
		},
	}, nil
}

// nearestField finds the field name closest to the given string using Levenshtein distance.
func nearestField(allowed []string, actual string) string {
	if len(allowed) == 0 {
		return ""
	}

	best := allowed[0]
	bestDist := levenshteinDistance(actual, best)

	for _, candidate := range allowed[1:] {
		dist := levenshteinDistance(actual, candidate)
		if dist < bestDist {
			best = candidate
			bestDist = dist
		}
	}

	// Only suggest if distance is reasonable (not too far)
	if bestDist > len(actual)/2 && bestDist > len(best)/2 {
		return ""
	}

	return best
}

// extractFieldNames extracts all field names from a schema document recursively.
func extractFieldNames(fields []*schemagen.SchemaField, prefix string) []string {
	var names []string
	for _, field := range fields {
		name := prefix + field.Name
		names = append(names, name)
		// Recurse for nested objects
		if field.Type == "object" && len(field.Children) > 0 {
			names = append(names, extractFieldNames(field.Children, name+".")...)
		}
	}
	return names
}

