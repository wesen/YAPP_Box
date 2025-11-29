package rules

import (
	"context"
	"sort"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
)

// Renderer is the interface that rules must implement.
type Renderer interface {
	// Match returns whether the rule applies and an optional score for ordering (higher wins).
	Match(t *errorx.Taxonomy) (ok bool, score int)
	// Render produces a help card for the taxonomy entry.
	Render(ctx context.Context, t *errorx.Taxonomy) (*RuleResult, error)
}

// RuleResult represents the output of a rule.
type RuleResult struct {
	Headline string        // Short summary (e.g., "YAML syntax error at line 5")
	Body     string        // Detailed help (can include code blocks, tables, links)
	Severity errorx.Severity // error, warning, or info
	Actions  []Action      // e.g., CLI commands, doc links (for future UI buttons)
}

// Action represents an actionable follow-up.
type Action struct {
	Label   string   // "View module docs"
	Command string   // "yappctl explain --module connectors"
	Args    []string // Command arguments
}

// Registry manages rule registration and rendering.
type Registry struct {
	rules []Renderer
}

// NewRegistry creates a new rule registry.
func NewRegistry() *Registry {
	return &Registry{
		rules: make([]Renderer, 0),
	}
}

// Register adds a rule to the registry.
func (r *Registry) Register(renderer Renderer) {
	r.rules = append(r.rules, renderer)
}

// RenderAll evaluates all registered rules and returns aggregated results.
// Results are sorted by score (descending) and severity (descending).
func (r *Registry) RenderAll(ctx context.Context, taxonomy *errorx.Taxonomy) ([]*RuleResult, error) {
	var matches []matchResult

	// Evaluate all rules
	for _, rule := range r.rules {
		ok, score := rule.Match(taxonomy)
		if !ok {
			continue
		}
		result, err := rule.Render(ctx, taxonomy)
		if err != nil {
			// Log error but continue with other rules
			continue
		}
		matches = append(matches, matchResult{
			result: result,
			score:  score,
		})
	}

	// Sort by score (descending), then by severity (error > warning > info)
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].score != matches[j].score {
			return matches[i].score > matches[j].score
		}
		return severityOrder(matches[i].result.Severity) > severityOrder(matches[j].result.Severity)
	})

	// Extract results and deduplicate by headline+body
	seen := make(map[string]bool)
	results := make([]*RuleResult, 0, len(matches))
	for _, m := range matches {
		// Create a key from headline and body to detect duplicates
		key := m.result.Headline + "|" + m.result.Body
		if seen[key] {
			continue
		}
		seen[key] = true
		results = append(results, m.result)
	}

	return results, nil
}

type matchResult struct {
	result *RuleResult
	score  int
}

// severityOrder returns numeric order for severity (higher = more severe).
func severityOrder(s errorx.Severity) int {
	switch s {
	case errorx.SeverityError:
		return 3
	case errorx.SeverityWarning:
		return 2
	case errorx.SeverityInfo:
		return 1
	default:
		return 0
	}
}

