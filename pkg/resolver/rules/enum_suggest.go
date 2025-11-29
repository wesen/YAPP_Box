package rules

import (
	"context"
	"fmt"
	"strings"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
)

// EnumSuggestClosestRule matches enum mismatches and suggests the closest valid value.
type EnumSuggestClosestRule struct{}

func (r *EnumSuggestClosestRule) Match(t *errorx.Taxonomy) (bool, int) {
	if t.Stage != errorx.StageSchemaConstraints || t.Symptom != errorx.SymptomEnumMismatch {
		return false, 0
	}
	_, ok := t.Context.(*errorx.SchemaConstraintContext)
	return ok, 90
}

func (r *EnumSuggestClosestRule) Render(ctx context.Context, t *errorx.Taxonomy) (*RuleResult, error) {
	sc, ok := t.Context.(*errorx.SchemaConstraintContext)
	if !ok {
		return nil, fmt.Errorf("expected SchemaConstraintContext, got %T", t.Context)
	}

	if len(sc.Allowed) == 0 {
		return nil, fmt.Errorf("no allowed values in context")
	}

	actualStr := fmt.Sprint(sc.Actual)
	suggestion := nearestEnum(sc.Allowed, actualStr)

	var body strings.Builder
	body.WriteString(fmt.Sprintf("Field `%s` must be one of: %s\n\n", sc.FieldPath, strings.Join(sc.Allowed, ", ")))
	
	if suggestion != "" && suggestion != actualStr {
		body.WriteString(fmt.Sprintf("💡 Did you mean `%s`?\n\n", suggestion))
	}

	body.WriteString("Valid values:\n")
	for _, val := range sc.Allowed {
		body.WriteString(fmt.Sprintf("- `%s`\n", val))
	}

	return &RuleResult{
		Headline: fmt.Sprintf("`%s` must be one of: %s", sc.FieldPath, strings.Join(sc.Allowed, ", ")),
		Body:     body.String(),
		Severity: errorx.SeverityError,
	}, nil
}

// nearestEnum finds the enum value closest to the given string using Levenshtein distance.
func nearestEnum(allowed []string, actual string) string {
	if len(allowed) == 0 {
		return ""
	}

	best := allowed[0]
	bestDist := levenshteinDistance(actual, best)

	for _, candidate := range allowed[1:] {
		dist := levenshteinDistance(actual, candidate)
		if dist < bestDist {
			best = candidate
			bestDist = dist
		}
	}

	// Only suggest if distance is reasonable (not too far)
	if bestDist > len(actual)/2 && bestDist > len(best)/2 {
		return ""
	}

	return best
}

// levenshteinDistance computes the Levenshtein distance between two strings.
func levenshteinDistance(s1, s2 string) int {
	r1, r2 := []rune(s1), []rune(s2)
	column := make([]int, len(r1)+1)

	for y := 1; y <= len(r1); y++ {
		column[y] = y
	}

	for x := 1; x <= len(r2); x++ {
		column[0] = x
		lastDiag := x - 1
		for y := 1; y <= len(r1); y++ {
			oldDiag := column[y]
			cost := 0
			if r1[y-1] != r2[x-1] {
				cost = 1
			}
			column[y] = min(
				column[y]+1,
				column[y-1]+1,
				lastDiag+cost,
			)
			lastDiag = oldDiag
		}
	}

	return column[len(r1)]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

