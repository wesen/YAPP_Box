package schemagen

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

//go:embed templates/schema_gen.go.tmpl
var schemaGenTemplate string

//go:embed templates/schema_gen_test.go.tmpl
var schemaTestTemplate string

//go:embed templates/modules_gen.go.tmpl
var modulesGenTemplate string

// GenerateOptions controls schemagen output.
type GenerateOptions struct {
	RootDir    string
	ModulePath string
}

// GenerateCode writes Go sources for all schema documents.
func GenerateCode(docs []*SchemaDoc, opts GenerateOptions) error {
	if opts.RootDir == "" {
		opts.RootDir = "."
	}
	if opts.ModulePath == "" {
		return fmt.Errorf("module path is required for code generation")
	}
	for _, doc := range docs {
		if err := generateModuleCode(doc, opts.RootDir); err != nil {
			return err
		}
		if err := generateModuleTests(doc, opts.RootDir); err != nil {
			return err
		}
	}
	if err := generateModulesRegistry(docs, opts.RootDir, opts.ModulePath); err != nil {
		return err
	}
	return nil
}

func generateModuleCode(doc *SchemaDoc, root string) error {
	structs := buildStructDefs(doc)
	tmpl, err := template.New("schema_gen").Parse(schemaGenTemplate)
	if err != nil {
		return fmt.Errorf("parse schema_gen template: %w", err)
	}

	data := map[string]any{
		"Package": doc.GoPackage,
		"Module":  doc.Module,
		"Structs": structs,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute schema_gen template: %w", err)
	}

	outPath := filepath.Join(resolveRoot(root, doc.ModuleDir()), "schema_gen.go")
	return writeGoFile(outPath, buf.Bytes())
}

func generateModuleTests(doc *SchemaDoc, root string) error {
	var positive []testData
	for _, test := range doc.Tests {
		if test.ExpectValid {
			var val any
			if test.InputNode != nil {
				if err := test.InputNode.Decode(&val); err != nil {
					continue
				}
			}
			data, _ := yaml.Marshal(val)
			positive = append(positive, testData{
				FuncName:  toCamel(test.Name),
				InputYAML: string(data),
			})
		}
	}
	if len(positive) == 0 {
		return nil
	}

	tmpl, err := template.New("schema_gen_test").Parse(schemaTestTemplate)
	if err != nil {
		return fmt.Errorf("parse schema_gen_test template: %w", err)
	}

	data := map[string]any{
		"Package":    doc.GoPackage,
		"RootStruct": doc.RootStructName(),
		"Tests":      positive,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute schema_gen_test template: %w", err)
	}

	outPath := filepath.Join(resolveRoot(root, doc.ModuleDir()), "schema_gen_test.go")
	return writeGoFile(outPath, buf.Bytes())
}

type testData struct {
	FuncName  string
	InputYAML string
}

func generateModulesRegistry(docs []*SchemaDoc, root, modulePath string) error {
	tmpl, err := template.New("modules_gen").Parse(modulesGenTemplate)
	if err != nil {
		return fmt.Errorf("parse modules_gen template: %w", err)
	}

	type moduleInfo struct {
		Module    string
		GoPackage string
	}
	var modules []moduleInfo
	for _, doc := range docs {
		modules = append(modules, moduleInfo{
			Module:    doc.Module,
			GoPackage: doc.GoPackage,
		})
	}

	data := map[string]any{
		"ModulePath": modulePath,
		"Modules":    modules,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute modules_gen template: %w", err)
	}

	outPath := filepath.Join(resolveRoot(root, filepath.Join("pkg", "yappgen")), "modules_gen.go")
	return writeGoFile(outPath, buf.Bytes())
}

type structFieldDef struct {
	GoName  string
	Type    string
	Tag     string
	Comment string
}

type structDef struct {
	Name               string
	Comment            string
	Fields             []structFieldDef
	DefaultAssignments []string
}

func buildStructDefs(doc *SchemaDoc) []structDef {
	var defs []structDef
	queue := []struct {
		name   string
		fields []*SchemaField
	}{
		{name: doc.RootStructName(), fields: doc.Fields},
	}
	seen := map[string]bool{}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if seen[current.name] {
			continue
		}
		seen[current.name] = true
		def := structDef{
			Name:    current.name,
			Comment: fmt.Sprintf("%s describes the %s schema.", current.name, doc.Module),
		}
		var assignments []string
		for _, f := range current.fields {
			goName := toCamel(f.Name)
			fieldType := scalarTypeForField(current.name, f)
			tag := fmt.Sprintf(`yaml:"%s`, f.Name)
			if !f.Required {
				tag += ",omitempty"
			}
			tag += `"`
			comment := f.Description
			if comment == "" {
				comment = fmt.Sprintf("%s field", f.Name)
			}
			def.Fields = append(def.Fields, structFieldDef{
				GoName:  goName,
				Type:    fieldType,
				Tag:     tag,
				Comment: comment,
			})
			if f.Default != nil && f.Type != "object" {
				assignments = append(assignments, buildDefaultAssignment("x."+goName, fieldType, f.Default))
			}
			if f.Type == "object" && len(f.Children) > 0 {
				childName := nestedStructName(current.name, f.Name)
				queue = append(queue, struct {
					name   string
					fields []*SchemaField
				}{name: childName, fields: f.Children})
			}
		}
		def.DefaultAssignments = assignments
		defs = append(defs, def)
	}
	return defs
}

func scalarTypeForField(parentName string, field *SchemaField) string {
	goType := ""
	switch field.Type {
	case "number":
		goType = "float64"
	case "string":
		goType = "string"
	case "bool":
		goType = "bool"
	case "object":
		goType = nestedStructName(parentName, field.Name)
	default:
		goType = "any"
	}
	if !field.Required {
		switch goType {
		case "float64":
			return "*float64"
		case "string":
			return "*string"
		case "bool":
			return "*bool"
		default:
			return "*" + goType
		}
	}
	return goType
}

func nestedStructName(parent, field string) string {
	return parent + toCamel(field)
}

func writeGoFile(path string, data []byte) error {
	formatted, err := format.Source(data)
	if err == nil {
		data = formatted
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func backtick(value string) string {
	if value == "" {
		return ""
	}
	return "`" + value + "`"
}

func escapeBackticks(s string) string {
	return strings.ReplaceAll(s, "`", "` + \"`\" + `")
}

func resolveRoot(root, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

func buildDefaultAssignment(field, fieldType string, defaultNode *yaml.Node) string {
	if defaultNode == nil {
		return ""
	}
	var val any
	if err := defaultNode.Decode(&val); err != nil {
		return ""
	}

	// For pointer types, check if field is nil before assigning
	isPointer := strings.HasPrefix(fieldType, "*")
	if isPointer {
		baseType := strings.TrimPrefix(fieldType, "*")
		switch baseType {
		case "string":
			if s, ok := val.(string); ok {
				return fmt.Sprintf("if %s == nil {\n\t\tv := %q\n\t\t%s = &v\n\t}", field, s, field)
			}
		case "float64":
			if n, ok := toNumber(val); ok {
				return fmt.Sprintf("if %s == nil {\n\t\tv := %v\n\t\t%s = &v\n\t}", field, n, field)
			}
		case "bool":
			if b, ok := val.(bool); ok {
				return fmt.Sprintf("if %s == nil {\n\t\tv := %v\n\t\t%s = &v\n\t}", field, b, field)
			}
		}
	}
	return ""
}

func toNumber(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case float64:
		return n, true
	case float32:
		return float64(n), true
	default:
		return 0, false
	}
}
