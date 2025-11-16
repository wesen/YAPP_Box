package schemagen

import (
	"fmt"
	"strings"
)

// ValidationError represents a single schema validation issue.
type ValidationError struct {
	File    string
	Path    string
	Line    int
	Column  int
	Message string
	Snippet string
	Hint    string
}

func (e ValidationError) Error() string {
	loc := e.File
	switch {
	case e.Line > 0 && e.Column > 0:
		loc = fmt.Sprintf("%s:%d:%d", loc, e.Line, e.Column)
	case e.Line > 0:
		loc = fmt.Sprintf("%s:%d", loc, e.Line)
	}

	path := ""
	if e.Path != "" {
		path = " [" + e.Path + "]"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%s%s: %s", loc, path, e.Message)
	if e.Snippet != "" {
		fmt.Fprintf(&b, "\n%s", e.Snippet)
	}
	if e.Hint != "" {
		fmt.Fprintf(&b, "\nHint: %s", e.Hint)
	}
	return b.String()
}

// ValidationErrors aggregates multiple validation errors.
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	if len(ve) == 1 {
		return ve[0].Error()
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d schema validation errors:\n", len(ve))
	for _, err := range ve {
		b.WriteString(" - ")
		b.WriteString(strings.ReplaceAll(err.Error(), "\n", "\n   "))
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
}

func (ve *ValidationErrors) add(err ValidationError) {
	*ve = append(*ve, err)
}

func (ve ValidationErrors) Len() int { return len(ve) }
