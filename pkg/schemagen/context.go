package schemagen

import (
	"fmt"
	"regexp"
	"strings"
)

type schemaContext struct {
	filename string
	lines    []string
}

func newSchemaContext(filename string, data []byte) *schemaContext {
	contents := string(data)
	contents = strings.ReplaceAll(contents, "\r\n", "\n")
	lines := strings.Split(contents, "\n")
	return &schemaContext{
		filename: filename,
		lines:    lines,
	}
}

func (ctx *schemaContext) snippet(line, column int) string {
	if ctx == nil || line <= 0 || line > len(ctx.lines) {
		return ""
	}
	start := max(1, line-1)
	end := min(len(ctx.lines), line+1)

	var b strings.Builder
	for l := start; l <= end; l++ {
		fmt.Fprintf(&b, "%4d | %s\n", l, ctx.lines[l-1])
		if l == line {
			caretCol := column
			if caretCol <= 0 {
				caretCol = len(ctx.lines[l-1])
				if caretCol == 0 {
					caretCol = 1
				}
			}
			if caretCol > len(ctx.lines[l-1]) {
				caretCol = len(ctx.lines[l-1])
			}
			if caretCol < 1 {
				caretCol = 1
			}
			b.WriteString("     | ")
			b.WriteString(strings.Repeat(" ", caretCol-1))
			b.WriteString("^\n")
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

var lineRegex = regexp.MustCompile(`line\s+(\d+)`)

func extractLine(msg string) int {
	m := lineRegex.FindStringSubmatch(msg)
	if len(m) < 2 {
		return 0
	}
	return atoiSafe(m[1])
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func atoiSafe(s string) int {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0
		}
		n = n*10 + int(r-'0')
	}
	return n
}
