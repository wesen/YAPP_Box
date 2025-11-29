package rules

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
	"github.com/wesen/yapp-encl-resolver/pkg/schemagen"
)

// ModuleDocEmbedRule matches schema validation errors and embeds module field documentation.
type ModuleDocEmbedRule struct{}

func (r *ModuleDocEmbedRule) Match(t *errorx.Taxonomy) (bool, int) {
	// Match both structure and constraint validation errors
	return (t.Stage == errorx.StageSchemaStructure || t.Stage == errorx.StageSchemaConstraints), 80
}

func (r *ModuleDocEmbedRule) Render(ctx context.Context, t *errorx.Taxonomy) (*RuleResult, error) {
	var moduleName string
	
	// Extract module name from context
	switch ctx := t.Context.(type) {
	case *errorx.SchemaStructureContext:
		moduleName = ctx.Module
	case *errorx.SchemaConstraintContext:
		moduleName = ctx.Module
	default:
		return nil, fmt.Errorf("expected SchemaStructureContext or SchemaConstraintContext, got %T", t.Context)
	}
	
	if moduleName == "" {
		return nil, fmt.Errorf("module name is empty")
	}
	
	// Load schema doc
	schemaPath := fmt.Sprintf("pkg/yappgen/modules/%s/schema.yaml", strings.ReplaceAll(moduleName, "_", ""))
	doc, err := schemagen.LoadSchemaDoc(schemaPath)
	if err != nil {
		// Schema file not found - skip embedding (legacy modules)
		return nil, fmt.Errorf("schema doc not found for %s: %w", moduleName, err)
	}
	
	// Render fields table
	fieldsTable := renderFieldsTable(doc.Fields, "")
	if fieldsTable == "" {
		return nil, fmt.Errorf("no fields found for module %s", moduleName)
	}
	
	var body strings.Builder
	body.WriteString(fmt.Sprintf("## %s Module Fields\n\n", doc.Module))
	if doc.Description != "" {
		body.WriteString(doc.Description)
		body.WriteString("\n\n")
	}
	body.WriteString("| Field | Type | Required | Default | Description |\n")
	body.WriteString("|-------|------|----------|---------|-------------|\n")
	body.WriteString(fieldsTable)
	
	return &RuleResult{
		Headline: fmt.Sprintf("Module documentation: %s", doc.Module),
		Body:     body.String(),
		Severity: errorx.SeverityInfo, // Informational, not an error
		Actions: []Action{
			{Label: fmt.Sprintf("View full module docs: yappctl help module-%s", moduleName), Command: ""},
		},
	}, nil
}

// renderFieldsTable renders a markdown table of fields, similar to renderFields in schema_help.go
func renderFieldsTable(fields []*schemagen.SchemaField, prefix string) string {
	var buf bytes.Buffer
	
	for _, field := range fields {
		name := prefix + field.Name
		
		req := ""
		if field.Required {
			req = "✓"
		}
		
		defaultVal := "—"
		if field.Default != nil {
			var val any
			if err := field.Default.Decode(&val); err == nil {
				defaultVal = fmt.Sprintf("`%v`", val)
			}
		}
		
		typeStr := field.Type
		if len(field.Enum) > 0 {
			typeStr = fmt.Sprintf("enum (%s)", strings.Join(field.Enum, ", "))
		}
		
		desc := field.Description
		if desc == "" {
			desc = field.Name
		}
		
		buf.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s | %s |\n",
			name, typeStr, req, defaultVal, desc))
		
		// Recurse for nested objects
		if field.Type == "object" && len(field.Children) > 0 {
			buf.WriteString(renderFieldsTable(field.Children, name+"."))
		}
	}
	
	return buf.String()
}
