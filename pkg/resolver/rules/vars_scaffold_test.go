package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
)

func TestVarsScaffoldRule_Match(t *testing.T) {
	rule := &VarsScaffoldRule{}
	
	// Should match missing dependency errors
	tax := errorx.NewExprDependencyTaxonomy("features.test.0.x", "max(vars.height, 10)", []string{"vars.height"}, 5, 1, 2)
	ok, score := rule.Match(tax)
	if !ok {
		t.Error("expected Match to return true for missing dependency")
	}
	if score != 100 {
		t.Errorf("expected score 100, got %d", score)
	}
	
	// Should not match other error types
	tax2 := errorx.NewYAMLIngestTaxonomy("test.yaml", 1, 2, "snippet")
	ok2, _ := rule.Match(tax2)
	if ok2 {
		t.Error("expected Match to return false for non-dependency error")
	}
}

func TestVarsScaffoldRule_Render(t *testing.T) {
	rule := &VarsScaffoldRule{}
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
	if !strings.Contains(result.Headline, "missing variables") || !strings.Contains(result.Headline, "Declare") {
		t.Errorf("expected Headline about missing variables, got %s", result.Headline)
	}
	
	if result.Body == "" {
		t.Error("result Body is empty")
	}
	if !strings.Contains(result.Body, "```yaml") {
		t.Error("expected Body to contain YAML code block")
	}
	if !strings.Contains(result.Body, "vars:") {
		t.Error("expected Body to contain vars section")
	}
	if !strings.Contains(result.Body, "height") {
		t.Error("expected Body to contain height variable")
	}
	if !strings.Contains(result.Body, "width") {
		t.Error("expected Body to contain width variable")
	}
	
	if result.Severity != errorx.SeverityError {
		t.Errorf("expected SeverityError, got %s", result.Severity)
	}
	
	if len(result.Actions) == 0 {
		t.Error("expected at least one action")
	}
}

func TestVarsScaffoldRule_Render_NoMissingRefs(t *testing.T) {
	rule := &VarsScaffoldRule{}
	ctx := context.Background()
	
	tax := errorx.NewExprDependencyTaxonomy("features.test.0.x", "expression", []string{}, 5, 1, 2)
	
	_, err := rule.Render(ctx, tax)
	if err == nil {
		t.Error("expected Render to return error when no missing refs")
	}
	if !strings.Contains(err.Error(), "no missing references") {
		t.Errorf("expected error about no missing references, got %v", err)
	}
}

func TestScaffoldVars(t *testing.T) {
	missingRefs := []string{"vars.height", "vars.width"}
	expression := "max(vars.height, vars.width, 10)"
	
	yaml := scaffoldVars(missingRefs, expression)
	
	if !strings.Contains(yaml, "vars:") {
		t.Error("expected 'vars:' in scaffold")
	}
	if !strings.Contains(yaml, "height:") {
		t.Error("expected 'height:' in scaffold")
	}
	if !strings.Contains(yaml, "width:") {
		t.Error("expected 'width:' in scaffold")
	}
}

func TestScaffoldVars_WithoutVarsPrefix(t *testing.T) {
	missingRefs := []string{"height", "width"}
	expression := "max(height, width)"
	
	yaml := scaffoldVars(missingRefs, expression)
	
	if !strings.Contains(yaml, "height:") {
		t.Error("expected 'height:' in scaffold")
	}
	if !strings.Contains(yaml, "width:") {
		t.Error("expected 'width:' in scaffold")
	}
}

func TestInferDefaultFromExpression(t *testing.T) {
	// Test max() pattern
	result := inferDefaultFromExpression("height", "max(vars.height, 10)")
	if result != "10" {
		t.Errorf("expected default '10', got %s", result)
	}
	
	// Test min() pattern
	result = inferDefaultFromExpression("width", "min(vars.width, 5)")
	if result != "5" {
		t.Errorf("expected default '5', got %s", result)
	}
	
	// Test empty expression
	result = inferDefaultFromExpression("height", "")
	if result != "0" {
		t.Errorf("expected default '0' for empty expression, got %s", result)
	}
	
	// Test expression without variable
	result = inferDefaultFromExpression("height", "some other expression")
	if result != "0" {
		t.Errorf("expected default '0' for unrelated expression, got %s", result)
	}
}

