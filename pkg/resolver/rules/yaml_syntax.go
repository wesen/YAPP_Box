package rules

import (
	"context"
	"fmt"
	"strings"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
)

// YamlSyntaxPointerRule matches YAML syntax errors and shows the offending line with a pointer.
type YamlSyntaxPointerRule struct{}

func (r *YamlSyntaxPointerRule) Match(t *errorx.Taxonomy) (bool, int) {
	return t.Stage == errorx.StageIngestYAMLSyntax && t.Symptom == errorx.SymptomSyntax, 100
}

func (r *YamlSyntaxPointerRule) Render(ctx context.Context, t *errorx.Taxonomy) (*RuleResult, error) {
	yc, ok := t.Context.(*errorx.YAMLIngestContext)
	if !ok {
		return nil, fmt.Errorf("expected YAMLIngestContext, got %T", t.Context)
	}

	body := renderSnippetWithPointer(yc.Snippet, yc.Line, yc.Column)

	return &RuleResult{
		Headline: fmt.Sprintf("YAML syntax error at %s:%d:%d", yc.File, yc.Line, yc.Column),
		Body:     body,
		Severity: errorx.SeverityError,
		Actions: []Action{
			{Label: "Check indentation and trailing commas", Command: ""},
		},
	}, nil
}

// renderSnippetWithPointer formats a YAML snippet with a visual pointer at the error location.
func renderSnippetWithPointer(snippet string, lineNum, column int) string {
	if snippet == "" {
		return "Unable to extract snippet."
	}

	lines := strings.Split(snippet, "\n")
	if len(lines) == 0 {
		return snippet
	}

	var buf strings.Builder
	buf.WriteString("```yaml\n")

	// Find the error line within the snippet
	// Snippet typically includes context lines, so we need to find which line has the error
	errorLineIdx := -1
	for i := range lines {
		// Heuristic: error is usually on the middle line of a 3-line snippet
		// or we can check if this line is around lineNum
		if i == len(lines)/2 || (lineNum > 0 && i < len(lines)) {
			errorLineIdx = i
			break
		}
	}
	if errorLineIdx < 0 {
		errorLineIdx = len(lines) / 2
	}

	// Print all lines
	for i, line := range lines {
		buf.WriteString(line)
		buf.WriteString("\n")
		// Add pointer on the error line
		if i == errorLineIdx && column > 0 {
			// Calculate pointer position (account for leading spaces)
			pointerPos := column - 1
			if pointerPos < 0 {
				pointerPos = 0
			}
			// Don't exceed line length
			if pointerPos > len(line) {
				pointerPos = len(line)
			}
			pointer := strings.Repeat(" ", pointerPos) + "^"
			buf.WriteString(pointer)
			buf.WriteString("\n")
		}
	}

	buf.WriteString("```\n")
	buf.WriteString("\nCommon YAML syntax issues:\n")
	buf.WriteString("- Missing colons after keys\n")
	buf.WriteString("- Incorrect indentation (use spaces, not tabs)\n")
	buf.WriteString("- Trailing commas in mappings\n")
	buf.WriteString("- Unquoted strings with special characters")

	return buf.String()
}

