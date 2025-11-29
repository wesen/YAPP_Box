package errorx

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatTaxonomy formats a taxonomy entry as a readable string.
func FormatTaxonomy(t *Taxonomy) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("Stage: %s\n", t.Stage))
	buf.WriteString(fmt.Sprintf("Symptom: %s\n", t.Symptom))
	buf.WriteString(fmt.Sprintf("Path: %s\n", t.Path))
	buf.WriteString(fmt.Sprintf("Severity: %s\n", t.Severity))
	buf.WriteString("\nContext:\n")

	switch ctx := t.Context.(type) {
	case *YAMLIngestContext:
		buf.WriteString(fmt.Sprintf("  File: %s\n", ctx.File))
		buf.WriteString(fmt.Sprintf("  Line: %d\n", ctx.Line))
		buf.WriteString(fmt.Sprintf("  Column: %d\n", ctx.Column))
		buf.WriteString(fmt.Sprintf("  Snippet: %s\n", ctx.Snippet))
	case *SchemaConstraintContext:
		buf.WriteString(fmt.Sprintf("  Module: %s\n", ctx.Module))
		buf.WriteString(fmt.Sprintf("  FieldPath: %s\n", ctx.FieldPath))
		buf.WriteString(fmt.Sprintf("  Expected: %s\n", ctx.Expected))
		buf.WriteString(fmt.Sprintf("  Actual: %v\n", ctx.Actual))
		if len(ctx.Allowed) > 0 {
			buf.WriteString(fmt.Sprintf("  Allowed: %v\n", ctx.Allowed))
		}
		if ctx.Line > 0 || ctx.Column > 0 {
			buf.WriteString(fmt.Sprintf("  Line: %d\n", ctx.Line))
			buf.WriteString(fmt.Sprintf("  Column: %d\n", ctx.Column))
		}
	case *SchemaStructureContext:
		buf.WriteString(fmt.Sprintf("  Module: %s\n", ctx.Module))
		buf.WriteString(fmt.Sprintf("  FieldPath: %s\n", ctx.FieldPath))
		buf.WriteString(fmt.Sprintf("  Expected: %s\n", ctx.Expected))
		buf.WriteString(fmt.Sprintf("  Actual: %s\n", ctx.Actual))
		buf.WriteString(fmt.Sprintf("  Required: %v\n", ctx.Required))
		if ctx.Line > 0 || ctx.Column > 0 {
			buf.WriteString(fmt.Sprintf("  Line: %d\n", ctx.Line))
			buf.WriteString(fmt.Sprintf("  Column: %d\n", ctx.Column))
		}
	case *ExprDependencyContext:
		buf.WriteString(fmt.Sprintf("  Expression: %s\n", ctx.Expression))
		buf.WriteString(fmt.Sprintf("  MissingRefs: %v\n", ctx.MissingRefs))
		buf.WriteString(fmt.Sprintf("  Iterations: %d\n", ctx.Iterations))
		if ctx.Line > 0 || ctx.Column > 0 {
			buf.WriteString(fmt.Sprintf("  Line: %d\n", ctx.Line))
			buf.WriteString(fmt.Sprintf("  Column: %d\n", ctx.Column))
		}
	case *ExprSyntaxContext:
		buf.WriteString(fmt.Sprintf("  Expression: %s\n", ctx.Expression))
		buf.WriteString(fmt.Sprintf("  Token: %s\n", ctx.Token))
		buf.WriteString(fmt.Sprintf("  Position: %d\n", ctx.Position))
		if ctx.Line > 0 || ctx.Column > 0 {
			buf.WriteString(fmt.Sprintf("  Line: %d\n", ctx.Line))
			buf.WriteString(fmt.Sprintf("  Column: %d\n", ctx.Column))
		}
	case *ExprRuntimeContext:
		buf.WriteString(fmt.Sprintf("  Expression: %s\n", ctx.Expression))
		buf.WriteString(fmt.Sprintf("  Function: %s\n", ctx.Function))
		buf.WriteString(fmt.Sprintf("  Message: %s\n", ctx.Message))
	case *StrictModeContext:
		buf.WriteString(fmt.Sprintf("  Kind: %s\n", ctx.Kind))
		buf.WriteString(fmt.Sprintf("  Offenders: %v\n", ctx.Offenders))
	default:
		buf.WriteString(fmt.Sprintf("  Unknown context type: %T\n", ctx))
	}

	return buf.String()
}

// FormatTaxonomyJSON formats a taxonomy entry as JSON.
func FormatTaxonomyJSON(t *Taxonomy) (string, error) {
	data := map[string]any{
		"stage":    string(t.Stage),
		"symptom":  string(t.Symptom),
		"path":     t.Path,
		"severity": string(t.Severity),
	}

	// Add context-specific fields
	switch ctx := t.Context.(type) {
	case *YAMLIngestContext:
		data["context"] = map[string]any{
			"file":    ctx.File,
			"line":    ctx.Line,
			"column":  ctx.Column,
			"snippet": ctx.Snippet,
		}
	case *SchemaConstraintContext:
		ctxData := map[string]any{
			"module":     ctx.Module,
			"field_path": ctx.FieldPath,
			"expected":   ctx.Expected,
			"actual":     ctx.Actual,
			"allowed":    ctx.Allowed,
		}
		if ctx.Line > 0 || ctx.Column > 0 {
			ctxData["line"] = ctx.Line
			ctxData["column"] = ctx.Column
		}
		data["context"] = ctxData
	case *SchemaStructureContext:
		ctxData := map[string]any{
			"module":     ctx.Module,
			"field_path": ctx.FieldPath,
			"expected":   ctx.Expected,
			"actual":     ctx.Actual,
			"required":   ctx.Required,
		}
		if ctx.Line > 0 || ctx.Column > 0 {
			ctxData["line"] = ctx.Line
			ctxData["column"] = ctx.Column
		}
		data["context"] = ctxData
	case *ExprDependencyContext:
		ctxData := map[string]any{
			"expression":   ctx.Expression,
			"missing_refs": ctx.MissingRefs,
			"iterations":   ctx.Iterations,
		}
		if ctx.Line > 0 || ctx.Column > 0 {
			ctxData["line"] = ctx.Line
			ctxData["column"] = ctx.Column
		}
		data["context"] = ctxData
	case *ExprSyntaxContext:
		ctxData := map[string]any{
			"expression": ctx.Expression,
			"token":      ctx.Token,
			"position":   ctx.Position,
		}
		if ctx.Line > 0 || ctx.Column > 0 {
			ctxData["line"] = ctx.Line
			ctxData["column"] = ctx.Column
		}
		data["context"] = ctxData
	case *ExprRuntimeContext:
		data["context"] = map[string]any{
			"expression": ctx.Expression,
			"function":   ctx.Function,
			"message":   ctx.Message,
		}
	case *StrictModeContext:
		data["context"] = map[string]any{
			"kind":      ctx.Kind,
			"offenders": ctx.Offenders,
		}
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

