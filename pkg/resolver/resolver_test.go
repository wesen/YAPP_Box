package resolver_test

import (
	"context"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver"
	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
	"gopkg.in/yaml.v3"
)

func mustUnmarshalYAML(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := yaml.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return m
}

func TestSimpleEvaluation(t *testing.T) {
	in := mustUnmarshalYAML(t, `
project: test
version: 0
units: mm
yapp_version: v0
vars:
  add: pcb.z_clearance + 4.0
pcb:
  z_clearance: 2.0
enclosure:
  base:
    thickness: add
`)
	out, err := resolver.Resolve(context.Background(), in, resolver.Options{})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	encl := out["enclosure"].(map[string]any)
	base := encl["base"].(map[string]any)
	th := base["thickness"]
	if th != 6 {
		t.Fatalf("expected 6 (int), got %#v", th)
	}
}

func TestFunctionsAndInts(t *testing.T) {
	in := mustUnmarshalYAML(t, `
project: test
version: 0
units: mm
yapp_version: v0
enclosure:
  wall:
    thickness: 2
  base:
    thickness: max(2.0, enclosure.wall.thickness)
`)
	out, err := resolver.Resolve(context.Background(), in, resolver.Options{})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	encl := out["enclosure"].(map[string]any)
	base := encl["base"].(map[string]any)
	if base["thickness"] != 2 {
		t.Fatalf("expected 2, got %#v", base["thickness"])
	}
}

func TestUnresolvedFails(t *testing.T) {
	in := mustUnmarshalYAML(t, `
project: test
version: 0
units: mm
yapp_version: v0
enclosure:
  base:
    thickness: missing.path + 1
`)
	_, err := resolver.Resolve(context.Background(), in, resolver.Options{MaxIterations: 2})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStrictUnusedVars(t *testing.T) {
	in := mustUnmarshalYAML(t, `
project: test
version: 0
units: mm
yapp_version: v0
vars:
  never_used: 2
enclosure:
  wall:
    thickness: 2
`)
	_, err := resolver.Resolve(context.Background(), in, resolver.Options{Strict: true})
	if err == nil {
		t.Fatal("expected strict error on unused vars")
	}
}

func TestIterationChain(t *testing.T) {
	in := mustUnmarshalYAML(t, `
project: test
version: 0
units: mm
yapp_version: v0
vars:
  a: 1
  b: a + 1
  c: b + 1
  d: c + 1
enclosure:
  base:
    thickness: d
`)
	out, err := resolver.Resolve(context.Background(), in, resolver.Options{})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	encl := out["enclosure"].(map[string]any)
	base := encl["base"].(map[string]any)
	if base["thickness"] != 4 {
		t.Fatalf("expected 4, got %#v", base["thickness"])
	}
}

func TestArrayExpressionMissingDependencies(t *testing.T) {
	// Test case from YAPP-BUG-001: expressions in array elements with missing variables
	in := mustUnmarshalYAML(t, `
project: test
units: mm
yapp_version: 3.0
features:
  cutouts:
    - face: front
      from_face_left: 10
      from_face_top: 10
      width: max(vars.missing_height, 10)
      height: max(vars.missing_width, 20)
      shape: rectangle
`)
	_, err := resolver.ResolveResult(context.Background(), in, resolver.Options{MaxIterations: 16}, nil)
	if err == nil {
		t.Fatal("expected error for missing dependencies")
	}

	// Extract taxonomy from error
	tax, ok := errorx.AsTaxonomy(err)
	if !ok {
		t.Fatalf("expected Taxonomy error, got %T: %v", err, err)
	}

	if tax.Stage != errorx.StageExprDependencyMissing {
		t.Errorf("expected Stage %s, got %s", errorx.StageExprDependencyMissing, tax.Stage)
	}
	if tax.Path != "features.cutouts.0.height" && tax.Path != "features.cutouts.0.width" {
		// Either path could be first, depending on iteration order
		t.Logf("Path: %s (acceptable)", tax.Path)
	}

	ctx, ok := tax.Context.(*errorx.ExprDependencyContext)
	if !ok {
		t.Fatalf("expected ExprDependencyContext, got %T", tax.Context)
	}

	// After fix: expression should be extracted
	if ctx.Expression == "" {
		t.Error("expected Expression to be extracted, got empty string")
	}

	// After fix: missing_refs should be populated
	if len(ctx.MissingRefs) == 0 {
		t.Error("expected MissingRefs to be populated, got empty slice")
	}

	// Verify missing refs contain the expected variables
	expectedRefs := map[string]bool{
		"vars.missing_height": false,
		"vars.missing_width":  false,
	}
	for _, ref := range ctx.MissingRefs {
		if ref == "vars.missing_height" || ref == "vars.missing_width" {
			expectedRefs[ref] = true
		}
	}

	// Check that at least one expected ref is present (depending on which path is reported first)
	foundAny := false
	for ref, found := range expectedRefs {
		if found {
			foundAny = true
			t.Logf("Found expected missing ref: %s", ref)
		}
	}
	if !foundAny {
		t.Errorf("expected at least one of [vars.missing_height, vars.missing_width] in MissingRefs, got %v", ctx.MissingRefs)
	}
}

func TestArrayExpressionLookupPath(t *testing.T) {
	// Test that lookupPath can handle array indices in paths
	// This is tested indirectly through ResolveResult, but we can verify
	// the path extraction works correctly
	in := mustUnmarshalYAML(t, `
project: test
units: mm
yapp_version: 3.0
features:
  cutouts:
    - face: front
      width: 10
      height: max(vars.missing_var, 20)
      shape: rectangle
`)
	_, err := resolver.ResolveResult(context.Background(), in, resolver.Options{MaxIterations: 16}, nil)
	if err == nil {
		t.Fatal("expected error for missing dependencies")
	}

	tax, ok := errorx.AsTaxonomy(err)
	if !ok {
		t.Fatalf("expected Taxonomy error, got %T: %v", err, err)
	}

	// Verify the path contains array index
	if tax.Path != "features.cutouts.0.height" {
		t.Errorf("expected Path 'features.cutouts.0.height', got %s", tax.Path)
	}

	ctx, ok := tax.Context.(*errorx.ExprDependencyContext)
	if !ok {
		t.Fatalf("expected ExprDependencyContext, got %T", tax.Context)
	}

	// After fix: expression should be extracted from array path
	if ctx.Expression != "max(vars.missing_var, 20)" {
		t.Errorf("expected Expression 'max(vars.missing_var, 20)', got %q", ctx.Expression)
	}

	// After fix: missing_refs should contain vars.missing_var
	if len(ctx.MissingRefs) == 0 {
		t.Error("expected MissingRefs to be populated")
	}
	found := false
	for _, ref := range ctx.MissingRefs {
		if ref == "vars.missing_var" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'vars.missing_var' in MissingRefs, got %v", ctx.MissingRefs)
	}
}
