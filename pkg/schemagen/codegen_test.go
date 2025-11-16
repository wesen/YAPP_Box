package schemagen

import (
	"os"
	"path/filepath"
	"testing"
)

const testSchema = `
module: demo_module
scad_array: demoArray
go_package: demo
description: Example
fields:
  x:
    type: number
    required: true
  inner:
    type: object
    required: true
    fields:
      name:
        type: string
tests:
  - name: valid_case
    desc: Should decode into struct
    input:
      x: 1
      inner:
        name: foo
    expect_valid: true
`

func TestGenerateCode_WritesFiles(t *testing.T) {
	tmpDir := t.TempDir()
	schemaPath := filepath.Join(tmpDir, "pkg", "yappgen", "modules", "demo", "schema.yaml")
	if err := os.MkdirAll(filepath.Dir(schemaPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(schemaPath, []byte(testSchema), 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}

	paths, err := DiscoverSchemaFiles(filepath.Join(tmpDir, "pkg", "yappgen", "modules"))
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	docs, err := LoadSchemaDocs(paths)
	if err != nil {
		t.Fatalf("load docs: %v", err)
	}
	if err := GenerateCode(docs, GenerateOptions{
		RootDir:    tmpDir,
		ModulePath: "github.com/example/project",
	}); err != nil {
		t.Fatalf("generate code: %v", err)
	}

	checkFiles := []string{
		filepath.Join(tmpDir, "pkg", "yappgen", "modules", "demo", "schema_gen.go"),
		filepath.Join(tmpDir, "pkg", "yappgen", "modules", "demo", "schema_gen_test.go"),
		filepath.Join(tmpDir, "pkg", "yappgen", "modules_gen.go"),
	}
	for _, f := range checkFiles {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("expected generated file %s, got error: %v", f, err)
		}
	}
}
