package rules

import (
	"fmt"
	"strings"
)

// RenderToText formats rule results as text for CLI output.
func RenderToText(results []*RuleResult) string {
	if len(results) == 0 {
		return ""
	}

	var buf strings.Builder
	for i, result := range results {
		if i > 0 {
			buf.WriteString("\n\n")
		}

		// Headline
		buf.WriteString(fmt.Sprintf("❌ %s\n", result.Headline))

		// Body
		if result.Body != "" {
			buf.WriteString("\n")
			buf.WriteString(result.Body)
			buf.WriteString("\n")
		}

		// Actions
		if len(result.Actions) > 0 {
			buf.WriteString("\n")
			for _, action := range result.Actions {
				if action.Command != "" {
					cmd := action.Command
					if len(action.Args) > 0 {
						cmd += " " + strings.Join(action.Args, " ")
					}
					buf.WriteString(fmt.Sprintf("💡 %s: %s\n", action.Label, cmd))
				} else {
					buf.WriteString(fmt.Sprintf("💡 %s\n", action.Label))
				}
			}
		}
	}

	return buf.String()
}

