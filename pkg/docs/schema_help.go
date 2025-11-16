package docs

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"

	"github.com/go-go-golems/glazed/pkg/help"

	"github.com/wesen/yapp-encl-resolver/pkg/registry"
	"github.com/wesen/yapp-encl-resolver/pkg/schemagen"
)

// LoadModuleHelp generates and registers help pages for all module schemas.
func LoadModuleHelp(helpSystem *help.HelpSystem) error {
	modules := registry.All()
	for _, module := range modules {
		schema := module.Schema()
		
		// Load schema doc for field metadata
		schemaPath := fmt.Sprintf("pkg/yappgen/modules/%s/schema.yaml", strings.ReplaceAll(schema.Name(), "_", ""))
		doc, err := schemagen.LoadSchemaDoc(schemaPath)
		if err != nil {
			// Skip if schema file not found (legacy modules)
			continue
		}
		
		markdown := renderSchemaHelp(doc)
		
		// Parse markdown section (Glazed expects frontmatter + content)
		section, err := help.LoadSectionFromMarkdown([]byte(markdown))
		if err != nil {
			return fmt.Errorf("parse help markdown for %s: %w", schema.Name(), err)
		}
		
		helpSystem.AddSection(section)
	}
	return nil
}

func renderSchemaHelp(doc *schemagen.SchemaDoc) string {
	var buf bytes.Buffer

	buf.WriteString("---\n")
	buf.WriteString(fmt.Sprintf("Title: %s Module Reference\n", doc.Module))
	buf.WriteString(fmt.Sprintf("Slug: module-%s\n", doc.Module))
	buf.WriteString(fmt.Sprintf("Short: %s\n", doc.Description))
	buf.WriteString("Topics:\n  - yapp\n  - dsl\n  - modules\n")
	buf.WriteString("Commands:\n  - yappctl\n")
	buf.WriteString("IsTemplate: false\n")
	buf.WriteString("IsTopLevel: true\n")
	buf.WriteString("ShowPerDefault: true\n")
	buf.WriteString("SectionType: GeneralTopic\n")
	buf.WriteString(fmt.Sprintf("Order: %d\n", 100+doc.Order/10))
	buf.WriteString("---\n\n")

	buf.WriteString(fmt.Sprintf("# %s Module\n\n", doc.Module))
	buf.WriteString(doc.Description)
	buf.WriteString("\n\n")

	buf.WriteString("## YAML Location\n\n")
	buf.WriteString(fmt.Sprintf("```yaml\nfeatures:\n  %s:\n    - # array of items\n```\n\n", doc.Module))

	buf.WriteString("## Fields\n\n")
	buf.WriteString("| Field | Type | Required | Default | Description |\n")
	buf.WriteString("|-------|------|----------|---------|-------------|\n")

	renderFields(&buf, doc.Fields, "")

	if len(doc.Tests) > 0 {
		buf.WriteString("\n## Examples\n\n")
		for _, test := range doc.Tests {
			if test.ExpectValid {
				buf.WriteString(fmt.Sprintf("### %s\n\n", test.Name))
				if test.Description != "" {
					buf.WriteString(test.Description)
					buf.WriteString("\n\n")
				}
				buf.WriteString("```yaml\n")
				buf.WriteString(fmt.Sprintf("%s:\n", doc.Module))
				buf.WriteString("  - ")
				// TODO: render test input YAML
				buf.WriteString("# example\n")
				buf.WriteString("```\n\n")
			}
		}
	}

	buf.WriteString("## OpenSCAD Array\n\n")
	buf.WriteString(fmt.Sprintf("This module generates the `%s` array in OpenSCAD.\n", doc.ScadArray))

	return buf.String()
}

func renderFields(buf *bytes.Buffer, fields []*schemagen.SchemaField, prefix string) {
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
			renderFields(buf, field.Children, name+".")
		}
	}
}
