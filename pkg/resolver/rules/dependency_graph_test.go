package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
)

func TestDependencyGraphRule_Match(t *testing.T) {
	rule := &DependencyGraphRule{}
	
	// Should match missing dependency errors
	tax := errorx.NewExprDependencyTaxonomy("features.test.0.x", "max(vars.height, 10)", []string{"vars.height"}, 5, 1, 2)
	ok, score := rule.Match(tax)
	if !ok {
		t.Error("expected Match to return true for missing dependency")
	}
	if score != 95 {
		t.Errorf("expected score 95, got %d", score)
	}
	
	// Should not match other error types
	tax2 := errorx.NewYAMLIngestTaxonomy("test.yaml", 1, 2, "snippet")
	ok2, _ := rule.Match(tax2)
	if ok2 {
		t.Error("expected Match to return false for non-dependency error")
	}
}

func TestDependencyGraphRule_Render(t *testing.T) {
	rule := &DependencyGraphRule{}
	ctx := context.Background()
	
	missingRefs := []string{"vars.height", "vars.width"}
	tax := errorx.NewExprDependencyTaxonomy("features.test.0.x", "max(vars.height, vars.width)", missingRefs, 5, 1, 2)
	
	result, err := rule.Render(ctx, tax)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	
	if result.Headline == "" {
		t.Error("result Headline is empty")
	}
	if !strings.Contains(result.Headline, "Missing variables") {
		t.Errorf("expected Headline about missing variables, got %s", result.Headline)
	}
	if !strings.Contains(result.Headline, "vars.height") {
		t.Errorf("expected Headline to contain vars.height, got %s", result.Headline)
	}
	
	if result.Body == "" {
		t.Error("result Body is empty")
	}
	if !strings.Contains(result.Body, "Missing variables") {
		t.Error("expected Body to contain 'Missing variables'")
	}
	if !strings.Contains(result.Body, "vars.height") {
		t.Error("expected Body to contain vars.height")
	}
	if !strings.Contains(result.Body, "vars.width") {
		t.Error("expected Body to contain vars.width")
	}
	if !strings.Contains(result.Body, "features.test.0.x") {
		t.Error("expected Body to contain error path")
	}
	if !strings.Contains(result.Body, "Dependency chain") {
		t.Error("expected Body to contain dependency chain")
	}
	if !strings.Contains(result.Body, "Suggested resolution order") {
		t.Error("expected Body to contain resolution order")
	}
	
	if result.Severity != errorx.SeverityError {
		t.Errorf("expected SeverityError, got %s", result.Severity)
	}
}

func TestDependencyGraphRule_Render_NoMissingRefs(t *testing.T) {
	rule := &DependencyGraphRule{}
	ctx := context.Background()
	
	// Test with empty missing refs but valid path
	tax := errorx.NewExprDependencyTaxonomy("features.test.0.x", "expression", []string{}, 5, 1, 2)
	
	result, err := rule.Render(ctx, tax)
	if err != nil {
		t.Fatalf("Render should handle empty missing refs gracefully, got error: %v", err)
	}
	
	if result == nil {
		t.Fatal("expected result even with empty missing refs")
	}
	
	if result.Severity != errorx.SeverityWarning {
		t.Errorf("expected SeverityWarning for empty refs, got %s", result.Severity)
	}
	
	if !strings.Contains(result.Body, "features.test.0.x") {
		t.Error("expected Body to contain error path even with empty refs")
	}
}

func TestDependencyGraphRule_Render_SingleMissingRef(t *testing.T) {
	rule := &DependencyGraphRule{}
	ctx := context.Background()
	
	missingRefs := []string{"vars.height"}
	tax := errorx.NewExprDependencyTaxonomy("features.test.0.x", "max(vars.height, 10)", missingRefs, 5, 1, 2)
	
	result, err := rule.Render(ctx, tax)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	
	if !strings.Contains(result.Body, "vars.height") {
		t.Error("expected Body to contain vars.height")
	}
	if !strings.Contains(result.Body, "→") {
		t.Error("expected Body to contain dependency chain arrow")
	}
}

func TestDependencyGraphRule_Render_NoExpression(t *testing.T) {
	rule := &DependencyGraphRule{}
	ctx := context.Background()
	
	missingRefs := []string{"vars.height"}
	tax := errorx.NewExprDependencyTaxonomy("features.test.0.x", "", missingRefs, 5, 1, 2)
	
	result, err := rule.Render(ctx, tax)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	
	// Should still work without expression
	if result.Body == "" {
		t.Error("result Body is empty")
	}
	if !strings.Contains(result.Body, "vars.height") {
		t.Error("expected Body to contain vars.height")
	}
}

