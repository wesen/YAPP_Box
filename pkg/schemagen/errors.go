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
}

func (e ValidationError) Error() string {
	loc := e.File
	if e.Line > 0 {
		loc = fmt.Sprintf("%s:%d", loc, e.Line)
		if e.Column > 0 {
			loc = fmt.Sprintf("%s:%d:%d", e.File, e.Line, e.Column)
		}
	}
	path := e.Path
	if path != "" {
		path = " [" + path + "]"
	}
	return fmt.Sprintf("%s%s: %s", loc, path, e.Message)
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
		b.WriteString(err.Error())
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
}

func (ve *ValidationErrors) add(err ValidationError) {
	*ve = append(*ve, err)
}

func (ve ValidationErrors) Len() int { return len(ve) }
