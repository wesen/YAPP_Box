package yappgen

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver"
)

// Provenance tracks how SCAD outputs map back to DSL paths and expressions.
type Provenance struct {
	trace        resolver.Trace
	resolved     map[string]any
	comments     map[string][]string
	raw          map[string]any
	scalarPaths  map[string]string
	featurePaths map[string][]string
}

// NewProvenance initializes a provenance tracker.
func NewProvenance(trace resolver.Trace, resolved map[string]any, comments map[string][]string, raw map[string]any) *Provenance {
	return &Provenance{
		trace:        trace,
		resolved:     resolved,
		comments:     comments,
		raw:          raw,
		scalarPaths:  map[string]string{},
		featurePaths: map[string][]string{},
	}
}

// AddScalar records the DSL path that feeds a SCAD global variable.
func (p *Provenance) AddScalar(scadName, path string) {
	if p == nil || scadName == "" || path == "" {
		return
	}
	p.scalarPaths[scadName] = path
}

// RegisterFeature records base paths for feature array items.
func (p *Provenance) RegisterFeature(featureKey string, count int) {
	if p == nil || featureKey == "" || count <= 0 {
		return
	}
	bases := make([]string, count)
	for i := 0; i < count; i++ {
		bases[i] = fmt.Sprintf("features.%s.%d", featureKey, i)
	}
	p.featurePaths[featureKey] = bases
}

// ScalarComment returns the comment lines for a SCAD scalar variable.
func (p *Provenance) ScalarComment(scadName string) []string {
	if p == nil {
		return nil
	}
	path, ok := p.scalarPaths[scadName]
	if !ok {
		return nil
	}
	lines := append([]string{}, p.emitCommentLines(path, 0)...)
	if line := p.formatFieldLine(path, scadName, p.lookupValue(path), 0); line != "" {
		lines = append(lines, line)
	}
	return lines
}

// DescribeFeatureRow produces comment lines for a single feature array row.
func (p *Provenance) DescribeFeatureRow(featureKey, scadName string, idx int, item map[string]any) []string {
	if p == nil || item == nil {
		return nil
	}
	bases, ok := p.featurePaths[featureKey]
	if !ok || idx < 0 || idx >= len(bases) {
		return nil
	}
	basePath := bases[idx]
	header := fmt.Sprintf("// %s[%d] ← %s", scadName, idx, humanizePath(basePath))
	lines := append([]string{}, p.emitCommentLines(basePath, 0)...)
	lines = append(lines, p.yamlCommentLines(basePath, 0)...)
	lines = append(lines, header)
	lines = append(lines, p.describeMap(basePath, "", item, 2)...)
	return lines
}

func (p *Provenance) describeMap(path, labelPrefix string, data map[string]any, indent int) []string {
	if data == nil {
		return nil
	}
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var lines []string
	for _, key := range keys {
		childLabel := joinLabel(labelPrefix, key)
		childPath := joinPathParts(path, key)
		lines = append(lines, p.describeValue(childPath, childLabel, data[key], indent)...)
	}
	return lines
}

func (p *Provenance) describeArray(path, labelPrefix string, data []any, indent int) []string {
	var lines []string
	for idx, val := range data {
		childLabel := fmt.Sprintf("%s[%d]", labelPrefix, idx)
		childPath := fmt.Sprintf("%s.%d", path, idx)
		lines = append(lines, p.describeValue(childPath, childLabel, val, indent)...)
	}
	return lines
}

func (p *Provenance) describeValue(path, label string, value any, indent int) []string {
	switch t := value.(type) {
	case map[string]any:
		lines := append([]string{}, p.emitCommentLines(path, indent)...)
		lines = append(lines, p.describeMap(path, label, t, indent)...)
		return lines
	case []any:
		lines := append([]string{}, p.emitCommentLines(path, indent)...)
		lines = append(lines, p.describeArray(path, label, t, indent)...)
		return lines
	default:
		lines := append([]string{}, p.emitCommentLines(path, indent)...)
		if line := p.formatFieldLine(path, label, t, indent); line != "" {
			lines = append(lines, line)
		}
		return lines
	}
}

func (p *Provenance) formatFieldLine(path, label string, value any, indent int) string {
	if path == "" || p == nil {
		return ""
	}
	humanPath := humanizePath(path)
	if label == "" {
		label = humanPath
	}
	prefix := strings.Repeat(" ", indent)
	if entry, ok := p.trace.Get(path); ok {
		switch entry.Kind {
		case resolver.TraceKindExpression:
			expr := entry.Expression
			if expr == "" {
				expr = "expression"
			}
			line := fmt.Sprintf("%s// %s (source: %s) expr=\"%s\" => %s",
				prefix, label, humanPath, expr, formatDisplayValue(entry.Value))
			if depStr := p.formatDependencies(entry.Dependencies); depStr != "" {
				line += " " + depStr
			}
			return line
		default:
			note := "literal"
			if entry.Expression != "" {
				note = fmt.Sprintf("literal %q", entry.Expression)
			}
			return fmt.Sprintf("%s// %s (source: %s) = %s (%s)",
				prefix, label, humanPath, formatDisplayValue(entry.Value), note)
		}
	}
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%s// %s (source: %s) = %s",
		prefix, label, humanPath, formatDisplayValue(value))
}

func (p *Provenance) formatDependencies(refs []string) string {
	if len(refs) == 0 || p.resolved == nil {
		return ""
	}
	var parts []string
	for _, ref := range refs {
		val, ok := lookupPath(p.resolved, ref)
		if !ok {
			continue
		}
		segment := fmt.Sprintf("%s=%s", humanizePath(ref), formatDisplayValue(val))
		if comment := p.firstComment(ref); comment != "" {
			segment = fmt.Sprintf("%s (%s)", segment, comment)
		}
		parts = append(parts, segment)
	}
	if len(parts) == 0 {
		return ""
	}
	return "deps: " + strings.Join(parts, ", ")
}

func (p *Provenance) lookupValue(path string) any {
	if p == nil || path == "" {
		return nil
	}
	val, _ := lookupPath(p.resolved, path)
	return val
}

func (p *Provenance) emitCommentLines(path string, indent int) []string {
	if p == nil || path == "" {
		return nil
	}
	raw := p.comments[path]
	if len(raw) == 0 {
		return nil
	}
	prefix := strings.Repeat(" ", indent)
	lines := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s// %s", prefix, line))
	}
	return lines
}

func (p *Provenance) firstComment(path string) string {
	if p == nil || path == "" {
		return ""
	}
	raw := p.comments[path]
	if len(raw) == 0 {
		return ""
	}
	return raw[0]
}

func (p *Provenance) yamlCommentLines(path string, indent int) []string {
	if p == nil || path == "" || p.raw == nil {
		return nil
	}
	val, ok := lookupPath(p.raw, path)
	if !ok {
		return nil
	}
	data, err := yaml.Marshal(val)
	if err != nil {
		return nil
	}
	content := strings.TrimRight(string(data), "\n")
	if content == "" {
		return nil
	}
	lines := strings.Split(content, "\n")
	prefix := strings.Repeat(" ", indent)
	out := make([]string, 0, len(lines)+1)
	out = append(out, fmt.Sprintf("%s// YAML %s:", prefix, humanizePath(path)))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, fmt.Sprintf("%s//   %s", prefix, line))
	}
	return out
}

func humanizePath(path string) string {
	if path == "" {
		return ""
	}
	parts := strings.Split(path, ".")
	out := strings.Builder{}
	for _, part := range parts {
		if part == "" {
			continue
		}
		if _, err := strconv.Atoi(part); err == nil {
			out.WriteString("[")
			out.WriteString(part)
			out.WriteString("]")
			continue
		}
		if out.Len() > 0 {
			out.WriteString(".")
		}
		out.WriteString(part)
	}
	return out.String()
}

func joinLabel(prefix, key string) string {
	if prefix == "" {
		return key
	}
	if key == "" {
		return prefix
	}
	return prefix + "." + key
}

func joinPathParts(base, key string) string {
	if base == "" {
		return key
	}
	if key == "" {
		return base
	}
	return base + "." + key
}

func formatDisplayValue(v any) string {
	switch t := v.(type) {
	case float64:
		return formatFloat(t)
	case float32:
		return formatFloat(float64(t))
	case int:
		return fmt.Sprintf("%d", t)
	case int64:
		return fmt.Sprintf("%d", t)
	case string:
		return fmt.Sprintf("%q", t)
	case bool:
		return fmt.Sprintf("%t", t)
	default:
		return fmt.Sprintf("%v", t)
	}
}
