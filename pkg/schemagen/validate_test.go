package schemagen

import (
	"errors"
	"testing"
)

const minimalSchema = `
module: test_feature
scad_array: testFeature
go_package: testfeature
description: Test feature schema
fields:
  x:
    type: number
tests:
  - name: valid
    input:
      x: 1
    expect_valid: true
`

func TestValidateSchemaBytes_Success(t *testing.T) {
	errs, err := ValidateSchemaBytes("test.yaml", []byte(minimalSchema))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %v", errs)
	}
}

func TestValidateSchemaBytes_MissingModule(t *testing.T) {
	doc := `
scad_array: testFeature
go_package: testfeature
description: Missing module
fields:
  x:
    type: number
`
	errs, err := ValidateSchemaBytes("broken.yaml", []byte(doc))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(errs) == 0 {
		t.Fatalf("expected validation errors")
	}

	if errs[0].Snippet == "" {
		t.Fatalf("expected snippet for missing module error")
	}

	found := false
	for _, e := range errs {
		if e.Path == "module" || e.Message == `missing required field "module"` {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected missing module error, got %v", errs)
	}
}

func TestValidateSchemaBytes_InvalidFieldType(t *testing.T) {
	doc := `
module: test_feature
scad_array: testFeature
go_package: testfeature
description: Invalid field type
fields:
  bad:
    type: banana
tests:
  - name: invalid
    input: {}
    expect_valid: false
`
	errs, err := ValidateSchemaBytes("broken.yaml", []byte(doc))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(errs) == 0 {
		t.Fatalf("expected validation errors")
	}
	found := false
	for _, e := range errs {
		if e.Path == "fields.bad.type" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected invalid field type error, got %v", errs)
	}
	if errs[0].Snippet == "" {
		t.Fatalf("expected snippet for invalid field type error")
	}
}

func TestValidateSchemaBytes_SyntaxError(t *testing.T) {
	doc := `
module broken
scad_array: broken
go_package: broken
description: Missing colon
fields:
  x:
    type number
`
	_, err := ValidateSchemaBytes("broken.yaml", []byte(doc))
	if err == nil {
		t.Fatalf("expected syntax error")
	}
	ve, ok := err.(ValidationErrors)
	if !ok {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
	if len(ve) == 0 {
		t.Fatalf("expected at least one error, got none")
	}
	if ve[0].Snippet == "" {
		t.Fatalf("expected snippet for syntax error")
	}
}

func TestWrapYAMLErrorAddsHint(t *testing.T) {
	ctx := newSchemaContext("fake.yaml", []byte("foo\nbar\n"))
	err := wrapYAMLError(ctx, errors.New("yaml: line 2: could not find expected ':'"))
	ve, ok := err.(ValidationErrors)
	if !ok || len(ve) == 0 {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
	if ve[0].Hint == "" {
		t.Fatalf("expected hint to be populated, got %+v", ve[0])
	}
}

func TestValidateSchemaBytes_BoolField(t *testing.T) {
	doc := `
module: demo
scad_array: demoArr
go_package: demo
description: Bool field test
fields:
  flag:
    type: bool
tests:
  - name: ok
    input:
      flag: true
    expect_valid: true
`
	errs, err := ValidateSchemaBytes("bool.yaml", []byte(doc))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(errs) != 0 {
		t.Fatalf("expected bool field to validate, got %v", errs)
	}
}
