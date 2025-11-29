package rules

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
)

// DependencyGraphRule matches missing variable errors and shows dependency chains.
type DependencyGraphRule struct{}

func (r *DependencyGraphRule) Match(t *errorx.Taxonomy) (bool, int) {
	return t.Stage == errorx.StageExprDependencyMissing, 95
}

func (r *DependencyGraphRule) Render(ctx context.Context, t *errorx.Taxonomy) (*RuleResult, error) {
	ec, ok := t.Context.(*errorx.ExprDependencyContext)
	if !ok {
		return nil, fmt.Errorf("expected ExprDependencyContext, got %T", t.Context)
	}

	if len(ec.MissingRefs) == 0 {
		// If we don't have missing refs but have a path, we can still show some info
		if t.Path == "" {
			return nil, fmt.Errorf("no missing references and no path in context")
		}
		// Fallback: show that dependencies couldn't be extracted
		var body strings.Builder
		body.WriteString(fmt.Sprintf("**Error location:** `%s`\n\n", t.Path))
		if ec.Expression != "" {
			body.WriteString(fmt.Sprintf("**Expression:** `%s`\n\n", ec.Expression))
		}
		body.WriteString("Unable to extract dependency information. Check the expression syntax.\n")
		
		return &RuleResult{
			Headline: fmt.Sprintf("Missing dependencies at `%s`", t.Path),
			Body:     body.String(),
			Severity: errorx.SeverityWarning,
		}, nil
	}

	// Build dependency chain visualization
	var body strings.Builder
	
	body.WriteString("Missing variables:\n")
	for _, ref := range ec.MissingRefs {
		body.WriteString(fmt.Sprintf("- `%s`\n", ref))
	}
	body.WriteString("\n")

	// Show where the error occurred
	body.WriteString(fmt.Sprintf("**Error location:** `%s`\n\n", t.Path))
	if ec.Expression != "" {
		body.WriteString(fmt.Sprintf("**Expression:** `%s`\n\n", ec.Expression))
	}

	// Build dependency chain
	body.WriteString("Dependency chain:\n")
	if len(ec.MissingRefs) == 1 {
		body.WriteString(fmt.Sprintf("  `%s` → `%s`\n\n", ec.MissingRefs[0], t.Path))
	} else {
		body.WriteString(fmt.Sprintf("  `%s` → `%s`\n\n", strings.Join(ec.MissingRefs, "`, `"), t.Path))
	}

	// Suggest resolution order
	body.WriteString("**Suggested resolution order:**\n")
	sortedRefs := make([]string, len(ec.MissingRefs))
	copy(sortedRefs, ec.MissingRefs)
	sort.Strings(sortedRefs)
	
	for i, ref := range sortedRefs {
		body.WriteString(fmt.Sprintf("  %d. Define `%s` first\n", i+1, ref))
	}

	return &RuleResult{
		Headline: fmt.Sprintf("Missing variables: %s", strings.Join(ec.MissingRefs, ", ")),
		Body:     body.String(),
		Severity: errorx.SeverityError,
		Actions: []Action{
			{Label: "See variable scaffold rule for YAML template", Command: ""},
		},
	}, nil
}

