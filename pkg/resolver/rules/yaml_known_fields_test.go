package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
	"github.com/wesen/yapp-encl-resolver/pkg/schemagen"
)

func TestYamlKnownFieldsRule_Match(t *testing.T) {
	rule := &YamlKnownFieldsRule{}
	
	// Should match strict mode unknown key errors
	tax := errorx.NewStrictModeTaxonomy("unknown-key", []string{"unknown_key1"})
	ok, score := rule.Match(tax)
	if !ok {
		t.Error("expected Match to return true for unknown key error")
	}
	if score != 85 {
		t.Errorf("expected score 85, got %d", score)
	}
	
	// Should not match other error types
	tax2 := errorx.NewYAMLIngestTaxonomy("test.yaml", 1, 2, "snippet")
	ok2, _ := rule.Match(tax2)
	if ok2 {
		t.Error("expected Match to return false for non-unknown-key error")
	}
	
	// Should not match unused-var strict mode errors
	tax3 := errorx.NewStrictModeTaxonomy("unused-var", []string{"vars.unused"})
	ok3, _ := rule.Match(tax3)
	if ok3 {
		t.Error("expected Match to return false for unused-var error")
	}
}

func TestYamlKnownFieldsRule_Render(t *testing.T) {
	rule := &YamlKnownFieldsRule{}
	ctx := context.Background()
	
	offenders := []string{"unkown_key", "featurs"}
	tax := errorx.NewStrictModeTaxonomy("unknown-key", offenders)
	
	result, err := rule.Render(ctx, tax)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	
	if result.Headline == "" {
		t.Error("result Headline is empty")
	}
	if !strings.Contains(result.Headline, "Unknown keys") {
		t.Errorf("expected Headline about unknown keys, got %s", result.Headline)
	}
	
	if result.Body == "" {
		t.Error("result Body is empty")
	}
	if !strings.Contains(result.Body, "unkown_key") {
		t.Error("expected Body to contain unknown key")
	}
	if !strings.Contains(result.Body, "featurs") {
		t.Error("expected Body to contain second unknown key")
	}
	if !strings.Contains(result.Body, "Did you mean") {
		t.Error("expected Body to contain suggestions")
	}
	
	if result.Severity != errorx.SeverityWarning {
		t.Errorf("expected SeverityWarning, got %s", result.Severity)
	}
}

func TestYamlKnownFieldsRule_Render_NoOffenders(t *testing.T) {
	rule := &YamlKnownFieldsRule{}
	ctx := context.Background()
	
	tax := errorx.NewStrictModeTaxonomy("unknown-key", []string{})
	
	_, err := rule.Render(ctx, tax)
	if err == nil {
		t.Error("expected Render to return error when no offenders")
	}
	if !strings.Contains(err.Error(), "no unknown keys") {
		t.Errorf("expected error about no unknown keys, got %v", err)
	}
}

func TestNearestField(t *testing.T) {
	known := []string{"features", "project", "vars", "enclosure"}
	
	// Test exact match
	result := nearestField(known, "features")
	if result != "features" {
		t.Errorf("expected 'features' for exact match, got %s", result)
	}
	
	// Test typo
	result = nearestField(known, "featurs")
	if result != "features" {
		t.Errorf("expected 'features' for 'featurs', got %s", result)
	}
	
	// Test no close match
	result = nearestField(known, "xyzabc")
	if result != "" {
		t.Errorf("expected empty string for no close match, got %s", result)
	}
}

func TestExtractFieldNames(t *testing.T) {
	// This test would require creating SchemaField structures
	// For now, just test that function exists and doesn't panic
	fields := []*schemagen.SchemaField{}
	result := extractFieldNames(fields, "")
	if len(result) != 0 {
		t.Errorf("expected empty result for empty fields, got %v", result)
	}
}

