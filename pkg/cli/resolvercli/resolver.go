package resolvercli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver"
)

// LoadOptions control how the CLI helper reads and resolves a DSL document.
type LoadOptions struct {
	MaxIterations int
	Strict        bool
}

// LoadResult contains the fully resolved document plus trace + comment metadata.
type LoadResult struct {
	Document map[string]any
	Trace    resolver.Trace
	Comments map[string][]string
}

// LoadAndResolveResult reads a YAML DSL file, resolves it, and returns metadata.
func LoadAndResolveResult(ctx context.Context, path string, opts LoadOptions) (*LoadResult, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "read input %s", path)
	}

	doc, comments, err := decodeDocumentWithComments(raw)
	if err != nil {
		return nil, err
	}
	result, err := resolver.ResolveResult(ctx, doc, resolver.Options{
		MaxIterations: opts.MaxIterations,
		Strict:        opts.Strict,
	})
	if err != nil {
		return nil, err
	}
	return &LoadResult{
		Document: result.Document,
		Trace:    result.Trace,
		Comments: comments,
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

func decodeDocumentWithComments(raw []byte) (map[string]any, map[string][]string, error) {
	var root yaml.Node
	dec := yaml.NewDecoder(bytes.NewReader(raw))
	dec.KnownFields(true)
	if err := dec.Decode(&root); err != nil {
		return nil, nil, errors.Wrap(err, "parse yaml")
	}
	node := &root
	if root.Kind == yaml.DocumentNode {
		if len(root.Content) == 0 {
			return nil, nil, errors.New("parse yaml: empty document")
		}
		node = root.Content[0]
	}
	comments := map[string][]string{}
	value, err := nodeToInterface(node, "", comments)
	if err != nil {
		return nil, nil, err
	}
	doc, ok := value.(map[string]any)
	if !ok {
		return nil, nil, errors.New("parse yaml: root document must be a mapping")
	}
	return doc, comments, nil
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
