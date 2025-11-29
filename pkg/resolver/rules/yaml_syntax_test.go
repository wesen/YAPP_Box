package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
)

func TestYamlSyntaxPointerRule_Match(t *testing.T) {
	rule := &YamlSyntaxPointerRule{}
	
	// Should match YAML syntax errors
	tax := errorx.NewYAMLIngestTaxonomy("test.yaml", 1, 2, "snippet")
	ok, score := rule.Match(tax)
	if !ok {
		t.Error("expected Match to return true for YAML syntax error")
	}
	if score != 100 {
		t.Errorf("expected score 100, got %d", score)
	}
	
	// Should not match other error types
	tax2 := errorx.NewSchemaStructureTaxonomy("features.test", "test", "field", "string", "", true, 1, 2)
	ok2, _ := rule.Match(tax2)
	if ok2 {
		t.Error("expected Match to return false for non-YAML syntax error")
	}
}

func TestYamlSyntaxPointerRule_Render(t *testing.T) {
	rule := &YamlSyntaxPointerRule{}
	ctx := context.Background()
	
	tax := errorx.NewYAMLIngestTaxonomy("test.yaml", 5, 10, "key: value\ninvalid line")
	
	result, err := rule.Render(ctx, tax)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	
	if result.Headline == "" {
		t.Error("result Headline is empty")
	}
	if !strings.Contains(result.Headline, "test.yaml") {
		t.Errorf("expected Headline to contain filename, got %s", result.Headline)
	}
	if !strings.Contains(result.Headline, "5") {
		t.Errorf("expected Headline to contain line number, got %s", result.Headline)
	}
	
	if result.Body == "" {
		t.Error("result Body is empty")
	}
	if !strings.Contains(result.Body, "```yaml") {
		t.Error("expected Body to contain code block")
	}
	if !strings.Contains(result.Body, "key: value") {
		t.Error("expected Body to contain snippet")
	}
	
	if result.Severity != errorx.SeverityError {
		t.Errorf("expected SeverityError, got %s", result.Severity)
	}
	
	if len(result.Actions) == 0 {
		t.Error("expected at least one action")
	}
}

func TestYamlSyntaxPointerRule_Render_EmptySnippet(t *testing.T) {
	rule := &YamlSyntaxPointerRule{}
	ctx := context.Background()
	
	tax := errorx.NewYAMLIngestTaxonomy("test.yaml", 1, 2, "")
	
	result, err := rule.Render(ctx, tax)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	
	if !strings.Contains(result.Body, "Unable to extract snippet") {
		t.Error("expected Body to indicate snippet unavailable")
	}
}

func TestRenderSnippetWithPointer(t *testing.T) {
	snippet := "line1\nline2\nline3"
	body := renderSnippetWithPointer(snippet, 2, 5)
	
	if !strings.Contains(body, "```yaml") {
		t.Error("expected code block marker")
	}
	if !strings.Contains(body, "line2") {
		t.Error("expected snippet content")
	}
	if !strings.Contains(body, "^") {
		t.Error("expected pointer character")
	}
}

func TestRenderSnippetWithPointer_EmptySnippet(t *testing.T) {
	body := renderSnippetWithPointer("", 1, 1)
	if body != "Unable to extract snippet." {
		t.Errorf("expected 'Unable to extract snippet.', got %s", body)
	}
}

