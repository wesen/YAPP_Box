package rules

import (
	"context"
	"fmt"
	"strings"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
)

// VarsScaffoldRule matches missing variable errors and generates YAML scaffold.
type VarsScaffoldRule struct{}

func (r *VarsScaffoldRule) Match(t *errorx.Taxonomy) (bool, int) {
	return t.Stage == errorx.StageExprDependencyMissing, 100
}

func (r *VarsScaffoldRule) Render(ctx context.Context, t *errorx.Taxonomy) (*RuleResult, error) {
	ec, ok := t.Context.(*errorx.ExprDependencyContext)
	if !ok {
		return nil, fmt.Errorf("expected ExprDependencyContext, got %T", t.Context)
	}

	if len(ec.MissingRefs) == 0 {
		return nil, fmt.Errorf("no missing references in context")
	}

	yaml := scaffoldVars(ec.MissingRefs, ec.Expression)

	var body strings.Builder
	body.WriteString("Add these variables to your YAML file:\n\n")
	body.WriteString("```yaml\n")
	body.WriteString(yaml)
	body.WriteString("\n```\n")

	return &RuleResult{
		Headline: "Declare missing variables",
		Body:     body.String(),
		Severity: errorx.SeverityError,
		Actions: []Action{
			{Label: "Copy the YAML above and add it to your vars section", Command: ""},
		},
	}, nil
}

// scaffoldVars generates YAML scaffold for missing variables.
// It tries to infer default values from expressions when possible.
func scaffoldVars(missingRefs []string, expression string) string {
	var buf strings.Builder
	buf.WriteString("vars:\n")

	for _, ref := range missingRefs {
		// Extract variable name (remove "vars." prefix if present)
		varName := ref
		if strings.HasPrefix(ref, "vars.") {
			varName = strings.TrimPrefix(ref, "vars.")
		}

		// Try to infer default from expression
		defaultVal := inferDefaultFromExpression(varName, expression)

		buf.WriteString(fmt.Sprintf("  %s: %s  # TODO: set appropriate value\n", varName, defaultVal))
	}

	return buf.String()
}

// inferDefaultFromExpression tries to extract a default value from an expression.
// For example, "max(vars.height, 10)" suggests default 10.
func inferDefaultFromExpression(varName, expression string) string {
	// Simple heuristic: look for numeric literals in the expression
	// This is a basic implementation; could be enhanced with proper expression parsing
	if expression == "" {
		return "0"
	}

	// Look for patterns like "max(vars.X, N)" or "min(vars.X, N)"
	varNameInExpr := varName
	if !strings.Contains(expression, varNameInExpr) {
		// Try with vars. prefix
		varNameInExpr = "vars." + varName
	}

	// Try to find numeric literals near the variable reference
	// This is a simple heuristic - could be improved
	if strings.Contains(expression, "max(") || strings.Contains(expression, "min(") {
		// Look for second argument which is often a default
		parts := strings.Split(expression, ",")
		if len(parts) >= 2 {
			// Try to extract number from second part
			for _, part := range parts[1:] {
				part = strings.TrimSpace(part)
				part = strings.TrimSuffix(part, ")")
				part = strings.TrimSpace(part)
				// Check if it's a number
				if len(part) > 0 && (part[0] >= '0' && part[0] <= '9') {
					return part
				}
			}
		}
	}

	// Default fallback
	return "0"
}

