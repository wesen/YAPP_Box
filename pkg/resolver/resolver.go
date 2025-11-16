package resolver

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/pkg/errors"
)

// Options configures the resolver behavior.
type Options struct {
	MaxIterations int
	Strict        bool
}

// Resolve evaluates expressions in the given document until reaching a fixed point
// or failing with an error (unresolved dependencies, cycles, invalid paths).
// It returns a deep copy with all expressions resolved to numeric values.
func Resolve(ctx context.Context, doc map[string]any, opts Options) (map[string]any, error) {
	if opts.MaxIterations <= 0 {
		opts.MaxIterations = 16
	}
	// Make a deep copy so we don't mutate input.
	state := deepCopy(doc).(map[string]any)

	// Optional strict top-level key validation
	if opts.Strict {
		if err := validateTopLevelKeys(state); err != nil {
			return nil, err
		}
	}

	// Phase 1: Validate structure before expression resolution
	if err := validateStructure(state); err != nil {
		return nil, errors.Wrap(err, "structure validation")
	}

	// Track used variables (vars.*) to optionally validate unused vars in strict mode.
	usedVars := map[string]struct{}{}

	var changed bool
	for iter := 0; iter < opts.MaxIterations; iter++ {
		select {
		case <-ctx.Done():
			return nil, errors.Wrap(ctx.Err(), "context")
		default:
		}
		changed = false
		resolved, iterChanged, err := resolvePass(state, usedVars)
		if err != nil {
			return nil, err
		}
		state = resolved
		changed = iterChanged
		if !changed {
			break
		}
	}

	unresolvedPaths := collectUnresolved(state)
	if len(unresolvedPaths) > 0 {
		// Try to report missing deps with a best-effort variable extraction.
		missing := map[string][]string{}
		for _, p := range unresolvedPaths {
			exprStr, _ := getPathString(state, p)
			missing[p] = findMissingDependencies(exprStr, state)
		}
		return nil, errors.Errorf("unresolved expressions after %d passes: %v (missing=%v)",
			opts.MaxIterations, unresolvedPaths, missing)
	}

	if opts.Strict {
		if err := validateUnusedVars(state, usedVars); err != nil {
			return nil, err
		}
	}

	// Phase 2: Validate constraints after expression resolution
	if err := validateConstraints(state); err != nil {
		return nil, errors.Wrap(err, "constraint validation")
	}

	return state, nil
}

// resolvePass performs one pass and tries to evaluate as many expressions as possible.
func resolvePass(state map[string]any, usedVars map[string]struct{}) (map[string]any, bool, error) {
	changed := false
	env := buildEnv(state)

	var walk func(path string, v any) (any, bool, error)
	walk = func(path string, v any) (any, bool, error) {
		switch t := v.(type) {
		case map[string]any:
			updated := false
			for k, vv := range t {
				newV, didChange, err := walk(joinPath(path, k), vv)
				if err != nil {
					return nil, false, err
				}
				if didChange {
					updated = true
					t[k] = newV
				}
			}
			return t, updated, nil
		case []any:
			updated := false
			for i, vv := range t {
				newV, didChange, err := walk(joinPath(path, fmt.Sprintf("%d", i)), vv)
				if err != nil {
					return nil, false, err
				}
				if didChange {
					updated = true
					t[i] = newV
				}
			}
			return t, updated, nil
		case string:
			// Treat some fields as literal strings, not expressions
			if isStringFieldPath(path) {
				return v, false, nil
			}
			// Try to parse as number if it looks like one
			if num, ok := parseNumericString(t); ok {
				return num, true, nil
			}
			// Treat as expression
			// Track vars.* usage (best effort)
			for _, ref := range extractVarRefs(t) {
				usedVars[ref] = struct{}{}
			}
			out, ok, err := tryEvalExpr(t, env)
			if err != nil {
				// Not resolvable yet is not an error here; only return error on fatal eval issues (syntax etc.)
				// We detect syntax by checking err type/message; be permissive and only consider "unexpected" tokens fatal.
				if isSyntaxError(err) {
					return nil, false, errors.Wrapf(err, "invalid expression at %s: %q", path, t)
				}
				return t, false, nil
			}
			if ok {
				changed = true
				return out, true, nil
			}
			return t, false, nil
		default:
			// numeric or other literal types remain as is
			return v, false, nil
		}
	}

	newState, iterChanged, err := walk("", state)
	if err != nil {
		return nil, false, err
	}
	return newState.(map[string]any), changed || iterChanged, nil
}

func buildEnv(state map[string]any) map[string]any {
	// Build an evaluation environment from the current state.
	// Flatten vars.* into the root so expressions can reference variables directly (e.g., "standoff_height").
	env := map[string]any{}
	// Copy state keys
	for k, v := range state {
		env[k] = v
	}
	// Hoist vars
	if vm, ok := state["vars"].(map[string]any); ok {
		for k, v := range vm {
			// Only add if not clashing with existing top-level keys
			if _, exists := env[k]; !exists {
				env[k] = v
			}
		}
	}
	return env
}

func toFloat(v any) (float64, error) {
	switch t := v.(type) {
	case int:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case float32:
		return float64(t), nil
	case float64:
		return t, nil
	default:
		return 0, errors.Errorf("expected number, got %T", v)
	}
}

func functionsEnv() map[string]func(params ...any) (any, error) {
	return map[string]func(params ...any) (any, error){
		"min": func(params ...any) (any, error) {
			if len(params) != 2 {
				return nil, errors.Errorf("min expects 2 params")
			}
			a, err := toFloat(params[0])
			if err != nil {
				return nil, err
			}
			b, err := toFloat(params[1])
			if err != nil {
				return nil, err
			}
			if a < b {
				return a, nil
			}
			return b, nil
		},
		"max": func(params ...any) (any, error) {
			if len(params) != 2 {
				return nil, errors.Errorf("max expects 2 params")
			}
			a, err := toFloat(params[0])
			if err != nil {
				return nil, err
			}
			b, err := toFloat(params[1])
			if err != nil {
				return nil, err
			}
			if a > b {
				return a, nil
			}
			return b, nil
		},
		"round": func(params ...any) (any, error) {
			if len(params) != 1 {
				return nil, errors.Errorf("round expects 1 param")
			}
			a, err := toFloat(params[0])
			if err != nil {
				return nil, err
			}
			return math.Round(a), nil
		},
		"floor": func(params ...any) (any, error) {
			if len(params) != 1 {
				return nil, errors.Errorf("floor expects 1 param")
			}
			a, err := toFloat(params[0])
			if err != nil {
				return nil, err
			}
			return math.Floor(a), nil
		},
		"ceil": func(params ...any) (any, error) {
			if len(params) != 1 {
				return nil, errors.Errorf("ceil expects 1 param")
			}
			a, err := toFloat(params[0])
			if err != nil {
				return nil, err
			}
			return math.Ceil(a), nil
		},
		"clamp": func(params ...any) (any, error) {
			if len(params) != 3 {
				return nil, errors.Errorf("clamp expects 3 params")
			}
			x, err := toFloat(params[0])
			if err != nil {
				return nil, err
			}
			lo, err := toFloat(params[1])
			if err != nil {
				return nil, err
			}
			hi, err := toFloat(params[2])
			if err != nil {
				return nil, err
			}
			if x < lo {
				return lo, nil
			}
			if x > hi {
				return hi, nil
			}
			return x, nil
		},
	}
}

func tryEvalExpr(expression string, env map[string]any) (any, bool, error) {
	opts := []expr.Option{
		expr.Env(env),
		expr.AllowUndefinedVariables(),
	}
	// Register functions explicitly
	for name, fn := range functionsEnv() {
		opts = append(opts, expr.Function(name, fn))
	}
	prog, err := expr.Compile(expression, opts...)
	if err != nil {
		return nil, false, err
	}
	v, err := expr.Run(prog, env)
	if err != nil {
		return nil, false, err
	}
	switch n := v.(type) {
	case int:
		return n, true, nil
	case int64:
		return int(n), true, nil
	case float32:
		return float64(n), true, nil
	case float64:
		// Cast to int if it's an integer value
		if isIntLike(n) {
			return int(n), true, nil
		}
		return n, true, nil
	default:
		// Non-numeric results are invalid
		return nil, false, errors.Errorf("non-numeric expression result: %T", v)
	}
}

func isIntLike(f float64) bool {
	_, frac := math.Modf(f)
	return math.Abs(frac) < 1e-9
}

func parseNumericString(s string) (any, bool) {
	ss := strings.TrimSpace(s)
	if i, err := strconv.Atoi(ss); err == nil {
		return i, true
	}
	if f, err := strconv.ParseFloat(ss, 64); err == nil {
		if isIntLike(f) {
			return int(f), true
		}
		return f, true
	}
	return nil, false
}

func deepCopy(v any) any {
	switch t := v.(type) {
	case map[string]any:
		m := make(map[string]any, len(t))
		for k, vv := range t {
			m[k] = deepCopy(vv)
		}
		return m
	case []any:
		a := make([]any, len(t))
		for i, vv := range t {
			a[i] = deepCopy(vv)
		}
		return a
	default:
		return t
	}
}

func joinPath(base, next string) string {
	if base == "" {
		return next
	}
	return base + "." + next
}

func collectUnresolved(state any) []string {
	var out []string
	ignore := map[string]struct{}{
		"project":                        {},
		"units":                          {},
		"yapp_version":                   {},
		"enclosure.lid.type":             {},
		"enclosure.lid.screws.positions": {},
		"coordinates.origin":             {},
		"coordinates.reference_plane":    {},
		"pcb.standoffs.type":             {},
	}
	var walk func(path string, v any)
	walk = func(path string, v any) {
		switch t := v.(type) {
		case map[string]any:
			for k, vv := range t {
				walk(joinPath(path, k), vv)
			}
		case []any:
			for i, vv := range t {
				walk(joinPath(path, fmt.Sprintf("%d", i)), vv)
			}
		case string:
			// Ignore known string-typed fields
			if _, ok := ignore[path]; ok || isStringFieldPath(path) {
				return
			}
			out = append(out, path)
		}
	}
	walk("", state)
	return out
}

var identRe = regexp.MustCompile(`\b[a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z0-9_]+)*\b`)

func isStringFieldPath(path string) bool {
	if path == "project" || path == "units" || path == "yapp_version" {
		return true
	}
	switch path {
	case "enclosure.lid.type", "enclosure.lid.screws.positions",
		"coordinates.origin", "coordinates.reference_plane",
		"pcb.standoffs.type":
		return true
	}
	// Any feature face selector
	if strings.HasSuffix(path, ".face") {
		return true
	}
	// Treat common enum-like fields as strings
	if strings.HasSuffix(path, ".shape") ||
		strings.HasSuffix(path, ".side") ||
		strings.HasSuffix(path, ".name") ||
		strings.HasSuffix(path, ".polygon") ||
		strings.HasSuffix(path, ".polygon_preset") ||
		strings.HasSuffix(path, ".shape_preset") ||
		strings.HasSuffix(path, ".coordinate") ||
		strings.HasSuffix(path, ".origin") {
		return true
	}
	return false
}

func extractVarRefs(exprStr string) []string {
	if exprStr == "" {
		return nil
	}
	var out []string
	matches := identRe.FindAllString(exprStr, -1)
	seen := map[string]struct{}{}
	for _, m := range matches {
		if !strings.HasPrefix(m, "vars.") {
			continue
		}
		if _, ok := seen[m]; ok {
			continue
		}
		seen[m] = struct{}{}
		out = append(out, m)
	}
	return out
}

func findMissingDependencies(exprStr string, state map[string]any) []string {
	if exprStr == "" {
		return nil
	}
	candidates := identRe.FindAllString(exprStr, -1)
	fns := map[string]struct{}{"min": {}, "max": {}, "round": {}, "floor": {}, "ceil": {}, "clamp": {}}
	var missing []string
	for _, c := range candidates {
		if _, isFn := fns[c]; isFn {
			continue
		}
		if _, ok := lookupPath(state, c); !ok {
			missing = append(missing, c)
			continue
		}
	}
	return missing
}

func getPathString(state map[string]any, path string) (string, bool) {
	v, ok := lookupPath(state, path)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func lookupPath(state any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	cur := state
	for _, p := range parts {
		switch t := cur.(type) {
		case map[string]any:
			v, ok := t[p]
			if !ok {
				return nil, false
			}
			cur = v
		case []any:
			// indexes are not expected in dotted paths here; bail
			return nil, false
		default:
			return nil, false
		}
	}
	return cur, true
}

func isSyntaxError(err error) bool {
	// Best-effort: expr returns descriptive errors, consider unexpected tokens fatal.
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unexpected") || strings.Contains(msg, "syntax")
}

func mergeMaps(dst, src map[string]any) {
	for k, v := range src {
		dst[k] = v
	}
}
