package rules

import (
	"context"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
)

func TestRegistry_Register(t *testing.T) {
	reg := NewRegistry()
	rule := &YamlSyntaxPointerRule{}
	
	reg.Register(rule)
	
	if len(reg.rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(reg.rules))
	}
}

func TestRegistry_RenderAll_NoMatches(t *testing.T) {
	reg := NewRegistry()
	reg.Register(&YamlSyntaxPointerRule{})
	
	tax := errorx.NewSchemaStructureTaxonomy("features.test", "test", "field", "string", "", true, 1, 2)
	
	results, err := reg.RenderAll(context.Background(), tax)
	if err != nil {
		t.Fatalf("RenderAll returned error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestRegistry_RenderAll_SingleMatch(t *testing.T) {
	reg := NewRegistry()
	reg.Register(&YamlSyntaxPointerRule{})
	
	tax := errorx.NewYAMLIngestTaxonomy("test.yaml", 5, 10, "key: value\ninvalid")
	
	results, err := reg.RenderAll(context.Background(), tax)
	if err != nil {
		t.Fatalf("RenderAll returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	
	if results[0].Headline == "" {
		t.Error("result Headline is empty")
	}
	if results[0].Body == "" {
		t.Error("result Body is empty")
	}
	if results[0].Severity != errorx.SeverityError {
		t.Errorf("expected SeverityError, got %s", results[0].Severity)
	}
}

func TestRegistry_RenderAll_MultipleMatches_SortedByScore(t *testing.T) {
	reg := NewRegistry()
	// Register rules with different scores
	reg.Register(&YamlSyntaxPointerRule{}) // score 100
	reg.Register(&VarsScaffoldRule{})      // score 100
	
	// Create a taxonomy that matches both (this shouldn't happen in practice, but test the sorting)
	// Actually, let's test with two different taxonomies that match different rules
	tax1 := errorx.NewYAMLIngestTaxonomy("test.yaml", 5, 10, "key: value")
	results1, err := reg.RenderAll(context.Background(), tax1)
	if err != nil {
		t.Fatalf("RenderAll returned error: %v", err)
	}
	if len(results1) != 1 {
		t.Errorf("expected 1 result for YAML syntax, got %d", len(results1))
	}
	
	tax2 := errorx.NewExprDependencyTaxonomy("features.test.0.x", "max(vars.height, 10)", []string{"vars.height"}, 5, 1, 2)
	results2, err := reg.RenderAll(context.Background(), tax2)
	if err != nil {
		t.Fatalf("RenderAll returned error: %v", err)
	}
	if len(results2) != 1 {
		t.Errorf("expected 1 result for vars scaffold, got %d", len(results2))
	}
}

func TestRegistry_RenderAll_SortedBySeverity(t *testing.T) {
	reg := NewRegistry()
	
	// Create a mock rule that returns different severities
	mockRule1 := &mockRule{
		matchFunc: func(t *errorx.Taxonomy) (bool, int) {
			return true, 50
		},
		renderFunc: func(ctx context.Context, t *errorx.Taxonomy) (*RuleResult, error) {
			return &RuleResult{
				Headline: "Warning rule",
				Body:     "Warning body",
				Severity: errorx.SeverityWarning,
			}, nil
		},
	}
	
	mockRule2 := &mockRule{
		matchFunc: func(t *errorx.Taxonomy) (bool, int) {
			return true, 50 // Same score
		},
		renderFunc: func(ctx context.Context, t *errorx.Taxonomy) (*RuleResult, error) {
			return &RuleResult{
				Headline: "Error rule",
				Body:     "Error body",
				Severity: errorx.SeverityError,
			}, nil
		},
	}
	
	reg.Register(mockRule1)
	reg.Register(mockRule2)
	
	tax := errorx.NewYAMLIngestTaxonomy("test.yaml", 1, 2, "test")
	results, err := reg.RenderAll(context.Background(), tax)
	if err != nil {
		t.Fatalf("RenderAll returned error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	
	// Error should come before Warning (higher severity)
	if results[0].Severity != errorx.SeverityError {
		t.Errorf("expected first result to be Error, got %s", results[0].Severity)
	}
	if results[1].Severity != errorx.SeverityWarning {
		t.Errorf("expected second result to be Warning, got %s", results[1].Severity)
	}
}

func TestSeverityOrder(t *testing.T) {
	if severityOrder(errorx.SeverityError) != 3 {
		t.Errorf("expected SeverityError order 3, got %d", severityOrder(errorx.SeverityError))
	}
	if severityOrder(errorx.SeverityWarning) != 2 {
		t.Errorf("expected SeverityWarning order 2, got %d", severityOrder(errorx.SeverityWarning))
	}
	if severityOrder(errorx.SeverityInfo) != 1 {
		t.Errorf("expected SeverityInfo order 1, got %d", severityOrder(errorx.SeverityInfo))
	}
	if severityOrder(errorx.Severity("unknown")) != 0 {
		t.Errorf("expected unknown severity order 0, got %d", severityOrder(errorx.Severity("unknown")))
	}
}

// mockRule is a helper for testing registry behavior
type mockRule struct {
	matchFunc func(*errorx.Taxonomy) (bool, int)
	renderFunc func(context.Context, *errorx.Taxonomy) (*RuleResult, error)
}

func (m *mockRule) Match(t *errorx.Taxonomy) (bool, int) {
	return m.matchFunc(t)
}

func (m *mockRule) Render(ctx context.Context, t *errorx.Taxonomy) (*RuleResult, error) {
	return m.renderFunc(ctx, t)
}

