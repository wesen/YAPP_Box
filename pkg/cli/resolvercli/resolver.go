package resolvercli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver"
	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
)

// LoadOptions control how the CLI helper reads and resolves a DSL document.
type LoadOptions struct {
	MaxIterations int
	Strict        bool
}

// DecodeResult contains the result of parsing a YAML document.
type DecodeResult struct {
	Document  map[string]any      // Parsed document structure
	Comments  map[string][]string // Comments extracted from YAML nodes
	Positions PositionMap         // Path-to-position mapping for error reporting
}

// LoadResult contains the fully resolved document plus trace + comment metadata.
type LoadResult struct {
	Document map[string]any
	Trace    resolver.Trace
	Comments map[string][]string
	Raw      map[string]any
	Positions PositionMap // Path-to-position mapping for error reporting
}

// LoadAndResolveResult reads a YAML DSL file, resolves it, and returns metadata.
func LoadAndResolveResult(ctx context.Context, path string, opts LoadOptions) (*LoadResult, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		// Wrap file system errors with taxonomy
		return nil, errors.Wrapf(err, "read input %s", path)
	}

	decodeResult, err := decodeDocumentWithComments(raw, path)
	if err != nil {
		return nil, err
	}
	result, err := resolver.ResolveResult(ctx, decodeResult.Document, resolver.Options{
		MaxIterations: opts.MaxIterations,
		Strict:        opts.Strict,
	})
	if err != nil {
		return nil, err
	}
	rawCopy := deepCopy(decodeResult.Document).(map[string]any)

	return &LoadResult{
		Document:  result.Document,
		Trace:     result.Trace,
		Comments:  decodeResult.Comments,
		Raw:       rawCopy,
		Positions: decodeResult.Positions,
	}, nil
}

// LoadAndResolve preserves the legacy API that returns only the resolved document.
func LoadAndResolve(ctx context.Context, path string, opts LoadOptions) (map[string]any, error) {
	result, err := LoadAndResolveResult(ctx, path, opts)
	if err != nil {
		return nil, err
	}
	return result.Document, nil
}

var yamlLineColRe = regexp.MustCompile(`line (\d+):(\d+):`)

func decodeDocumentWithComments(raw []byte, filePath string) (*DecodeResult, error) {
	var root yaml.Node
	dec := yaml.NewDecoder(bytes.NewReader(raw))
	dec.KnownFields(true)
	if err := dec.Decode(&root); err != nil {
		// Extract line/column from yaml.v3 error message
		line, column := extractLineColumn(err)
		snippet := extractSnippet(raw, line)
		taxonomy := errorx.NewYAMLIngestTaxonomy(filePath, line, column, snippet)
		return nil, errors.Wrap(taxonomy, "parse yaml")
	}
	node := &root
	if root.Kind == yaml.DocumentNode {
		if len(root.Content) == 0 {
			return nil, errors.New("parse yaml: empty document")
		}
		node = root.Content[0]
	}
	comments := map[string][]string{}
	positions := make(PositionMap)
	buildPositionMap(node, "", positions)
	value, err := nodeToInterface(node, "", comments)
	if err != nil {
		return nil, err
	}
	doc, ok := value.(map[string]any)
	if !ok {
		return nil, errors.New("parse yaml: root document must be a mapping")
	}
	return &DecodeResult{
		Document:  doc,
		Comments:  comments,
		Positions: positions,
	}, nil
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

func nodeToInterface(node *yaml.Node, path string, comments map[string][]string) (any, error) {
	if c := extractComments(node); len(c) > 0 && path != "" {
		comments[path] = append(comments[path], c...)
	}
	switch node.Kind {
	case yaml.MappingNode:
		m := make(map[string]any, len(node.Content)/2)
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valNode := node.Content[i+1]
			key := strings.TrimSpace(keyNode.Value)
			if key == "" {
				return nil, errors.Errorf("parse yaml: empty key at line %d", keyNode.Line)
			}
			childPath := joinPath(path, key)
			val, err := nodeToInterface(valNode, childPath, comments)
			if err != nil {
				return nil, err
			}
			m[key] = val
		}
		return m, nil
	case yaml.SequenceNode:
		arr := make([]any, len(node.Content))
		for idx, child := range node.Content {
			childPath := joinPath(path, fmt.Sprintf("%d", idx))
			val, err := nodeToInterface(child, childPath, comments)
			if err != nil {
				return nil, err
			}
			arr[idx] = val
		}
		return arr, nil
	case yaml.ScalarNode:
		var out any
		if err := node.Decode(&out); err != nil {
			return nil, errors.Wrapf(err, "parse scalar at line %d", node.Line)
		}
		return out, nil
	default:
		return nil, errors.Errorf("parse yaml: unsupported node kind %d", node.Kind)
	}
}

func extractComments(node *yaml.Node) []string {
	if node == nil {
		return nil
	}
	raw := []string{node.HeadComment, node.LineComment, node.FootComment}
	var out []string
	for _, block := range raw {
		for _, line := range strings.Split(block, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "#") {
				line = strings.TrimSpace(strings.TrimPrefix(line, "#"))
			}
			if line == "" {
				continue
			}
			out = append(out, line)
		}
	}
	return out
}

func joinPath(base, next string) string {
	if base == "" {
		return next
	}
	if next == "" {
		return base
	}
	return base + "." + next
}

// extractLineColumn extracts line and column from yaml.v3 error message.
func extractLineColumn(err error) (line, column int) {
	msg := err.Error()
	matches := yamlLineColRe.FindStringSubmatch(msg)
	if len(matches) >= 3 {
		if l, err := strconv.Atoi(matches[1]); err == nil {
			line = l
		}
		if c, err := strconv.Atoi(matches[2]); err == nil {
			column = c
		}
	}
	return line, column
}

// extractSnippet extracts a snippet of YAML around the given line.
func extractSnippet(raw []byte, lineNum int) string {
	if lineNum <= 0 {
		return ""
	}
	lines := strings.Split(string(raw), "\n")
	if lineNum > len(lines) {
		return ""
	}
	// Show 3 lines of context (line before, error line, line after)
	start := lineNum - 2
	if start < 0 {
		start = 0
	}
	end := lineNum + 1
	if end > len(lines) {
		end = len(lines)
	}
	return strings.Join(lines[start:end], "\n")
}
