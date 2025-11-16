package schemagen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SchemaDoc represents a parsed module schema file.
type SchemaDoc struct {
	Path        string
	Module      string
	Order       int
	ScadArray   string
	GoPackage   string
	Description string
	Fields      []*SchemaField
	Tests       []SchemaTest
}

// SchemaField captures metadata for a single schema field.
type SchemaField struct {
	Name        string
	Type        string
	Required    bool
	Default     *yaml.Node
	Description string
	Enum        []string
	Children    []*SchemaField
}

// SchemaTest represents an entry in the schema tests section.
type SchemaTest struct {
	Name        string
	Description string
	ExpectValid bool
	InputNode   *yaml.Node
}

// LoadSchemaDoc parses a schema.yaml file into a SchemaDoc structure.
func LoadSchemaDoc(path string) (*SchemaDoc, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, err
	}
	if len(root.Content) == 0 {
		return nil, fmt.Errorf("schema %s is empty", path)
	}
	docNode := root.Content[0]

	doc := &SchemaDoc{
		Path: path,
	}
	doc.Module = stringField(docNode, "module")
	doc.Order = intField(docNode, "order")
	doc.ScadArray = stringField(docNode, "scad_array")
	doc.GoPackage = stringField(docNode, "go_package")
	doc.Description = stringField(docNode, "description")

	if fieldsNode := mapValueNode(docNode, "fields"); fieldsNode != nil {
		doc.Fields = parseFields(fieldsNode)
	}
	if testsNode := mapValueNode(docNode, "tests"); testsNode != nil {
		doc.Tests = parseTests(testsNode)
	}
	return doc, nil
}

func stringField(node *yaml.Node, key string) string {
	if v := mapValueNode(node, key); v != nil && v.Kind == yaml.ScalarNode {
		return v.Value
	}
	return ""
}

func intField(node *yaml.Node, key string) int {
	if v := mapValueNode(node, key); v != nil && v.Kind == yaml.ScalarNode {
		var out int
		fmt.Sscanf(v.Value, "%d", &out)
		return out
	}
	return 0
}

func parseFields(node *yaml.Node) []*SchemaField {
	var fields []*SchemaField
	if node == nil || node.Kind != yaml.MappingNode {
		return fields
	}
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]
		field := &SchemaField{
			Name:        keyNode.Value,
			Type:        stringField(valueNode, "type"),
			Required:    boolField(valueNode, "required"),
			Description: stringField(valueNode, "desc"),
		}
		if enumNode := mapValueNode(valueNode, "enum"); enumNode != nil && enumNode.Kind == yaml.SequenceNode {
			for _, entry := range enumNode.Content {
				if entry.Kind == yaml.ScalarNode {
					field.Enum = append(field.Enum, entry.Value)
				}
			}
		}
		if defNode := mapValueNode(valueNode, "default"); defNode != nil {
			field.Default = defNode
		}
		if defNode := mapValueNode(valueNode, "default"); defNode != nil {
			field.Default = defNode
		}
		if field.Type == "object" {
			field.Children = parseFields(mapValueNode(valueNode, "fields"))
		}
		fields = append(fields, field)
	}
	return fields
}

func parseTests(node *yaml.Node) []SchemaTest {
	var tests []SchemaTest
	if node.Kind != yaml.SequenceNode {
		return tests
	}
	for _, entry := range node.Content {
		if entry.Kind != yaml.MappingNode {
			continue
		}
		test := SchemaTest{
			Name:        stringField(entry, "name"),
			Description: stringField(entry, "desc"),
			ExpectValid: boolField(entry, "expect_valid"),
		}
		test.InputNode = mapValueNode(entry, "input")
		tests = append(tests, test)
	}
	return tests
}

func mapValueNode(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

func boolField(node *yaml.Node, key string) bool {
	if v := mapValueNode(node, key); v != nil && v.Kind == yaml.ScalarNode {
		switch strings.ToLower(v.Value) {
		case "true", "yes", "on":
			return true
		}
	}
	return false
}

// RootStructName returns the exported struct name for the schema items.
func (d *SchemaDoc) RootStructName() string {
	return toCamel(d.Module) + "Item"
}

// ModuleDir returns the directory containing the schema file.
func (d *SchemaDoc) ModuleDir() string {
	return filepath.Dir(d.Path)
}

func toCamel(s string) string {
	if s == "" {
		return ""
	}
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == '.'
	})
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
	}
	return strings.Join(parts, "")
}
