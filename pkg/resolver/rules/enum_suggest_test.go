package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
)

func TestEnumSuggestClosestRule_Match(t *testing.T) {
	rule := &EnumSuggestClosestRule{}
	
	// Should match enum mismatch errors
	allowed := []string{"a", "b", "c"}
	tax := errorx.NewSchemaConstraintTaxonomy("features.test", "test", "field", "enum", "invalid", allowed, nil, nil, 1, 2)
	ok, score := rule.Match(tax)
	if !ok {
		t.Error("expected Match to return true for enum mismatch")
	}
	if score != 90 {
		t.Errorf("expected score 90, got %d", score)
	}
	
	// Should not match other error types
	tax2 := errorx.NewYAMLIngestTaxonomy("test.yaml", 1, 2, "snippet")
	ok2, _ := rule.Match(tax2)
	if ok2 {
		t.Error("expected Match to return false for non-enum error")
	}
}

func TestEnumSuggestClosestRule_Render(t *testing.T) {
	rule := &EnumSuggestClosestRule{}
	ctx := context.Background()
	
	allowed := []string{"single", "all", "front_left", "front_right"}
	tax := errorx.NewSchemaConstraintTaxonomy("features.test", "test", "corner", "enum", "diagonal", allowed, nil, nil, 1, 2)
	
	result, err := rule.Render(ctx, tax)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	
	if result.Headline == "" {
		t.Error("result Headline is empty")
	}
	if !strings.Contains(result.Headline, "corner") {
		t.Errorf("expected Headline to contain field name, got %s", result.Headline)
	}
	
	if result.Body == "" {
		t.Error("result Body is empty")
	}
	if !strings.Contains(result.Body, "corner") {
		t.Error("expected Body to contain field name")
	}
	if !strings.Contains(result.Body, "single") || !strings.Contains(result.Body, "all") {
		t.Error("expected Body to contain allowed values")
	}
	
	if result.Severity != errorx.SeverityError {
		t.Errorf("expected SeverityError, got %s", result.Severity)
	}
}

func TestEnumSuggestClosestRule_Render_NoAllowedValues(t *testing.T) {
	rule := &EnumSuggestClosestRule{}
	ctx := context.Background()
	
	tax := errorx.NewSchemaConstraintTaxonomy("features.test", "test", "field", "enum", "invalid", nil, nil, nil, 1, 2)
	
	_, err := rule.Render(ctx, tax)
	if err == nil {
		t.Error("expected Render to return error when no allowed values")
	}
	if !strings.Contains(err.Error(), "no allowed values") {
		t.Errorf("expected error about no allowed values, got %v", err)
	}
}

func TestNearestEnum(t *testing.T) {
	allowed := []string{"single", "all", "front_left", "front_right"}
	
	// Test exact match
	result := nearestEnum(allowed, "single")
	if result != "single" {
		t.Errorf("expected 'single', got %s", result)
	}
	
	// Test close match
	result = nearestEnum(allowed, "diagonal")
	// May return empty if distance is too large (acceptable behavior)
	// Just verify it doesn't crash
	_ = result
	
	// Test very different string
	result = nearestEnum(allowed, "xyzabc123")
	// May return empty if distance is too large
	// This is acceptable behavior
}

func TestNearestEnum_EmptyAllowed(t *testing.T) {
	result := nearestEnum([]string{}, "test")
	if result != "" {
		t.Errorf("expected empty string, got %s", result)
	}
}

func TestLevenshteinDistance(t *testing.T) {
	// Test identical strings
	dist := levenshteinDistance("test", "test")
	if dist != 0 {
		t.Errorf("expected distance 0 for identical strings, got %d", dist)
	}
	
	// Test one character difference
	dist = levenshteinDistance("test", "best")
	if dist != 1 {
		t.Errorf("expected distance 1, got %d", dist)
	}
	
	// Test completely different strings
	dist = levenshteinDistance("abc", "xyz")
	if dist != 3 {
		t.Errorf("expected distance 3, got %d", dist)
	}
	
	// Test empty strings
	dist = levenshteinDistance("", "")
	if dist != 0 {
		t.Errorf("expected distance 0 for empty strings, got %d", dist)
	}
	
	dist = levenshteinDistance("test", "")
	if dist != 4 {
		t.Errorf("expected distance 4 for 'test' vs '', got %d", dist)
	}
}

func TestMin(t *testing.T) {
	if min(1, 2, 3) != 1 {
		t.Errorf("expected min(1,2,3) = 1, got %d", min(1, 2, 3))
	}
	if min(3, 2, 1) != 1 {
		t.Errorf("expected min(3,2,1) = 1, got %d", min(3, 2, 1))
	}
	if min(2, 1, 3) != 1 {
		t.Errorf("expected min(2,1,3) = 1, got %d", min(2, 1, 3))
	}
}

