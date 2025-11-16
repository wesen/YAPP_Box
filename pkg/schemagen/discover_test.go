package schemagen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverSchemaFiles_FindsSchemas(t *testing.T) {
	dir := t.TempDir()
	pushDir := filepath.Join(dir, "pushbuttons")
	if err := os.MkdirAll(pushDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pushDir, "schema.yaml"), []byte(minimalSchema), 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}

	files, err := DiscoverSchemaFiles(dir)
	if err != nil {
		t.Fatalf("DiscoverSchemaFiles: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 schema, got %d", len(files))
	}
	if files[0] != filepath.Join(pushDir, "schema.yaml") {
		t.Fatalf("unexpected schema path %s", files[0])
	}
}

func TestValidateSchemaFiles_AggregatesErrors(t *testing.T) {
	dir := t.TempDir()
	badDir := filepath.Join(dir, "bad")
	if err := os.MkdirAll(badDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	badSchema := `
scad_array: bad
go_package: bad
description: missing module
fields:
  x:
    type: number
`
	path := filepath.Join(badDir, "schema.yaml")
	if err := os.WriteFile(path, []byte(badSchema), 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}

	err := ValidateSchemaFiles([]string{path})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if _, ok := err.(ValidationErrors); !ok {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
}
