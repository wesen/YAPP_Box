package schemagen

import "testing"

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
}
