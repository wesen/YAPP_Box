package resolvercli

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver"
)

// LoadOptions control how the CLI helper reads and resolves a DSL document.
type LoadOptions struct {
	MaxIterations int
	Strict        bool
}

// LoadResult contains the fully resolved document and the trace metadata.
type LoadResult struct {
	Document map[string]any
	Trace    resolver.Trace
}

// LoadAndResolveResult reads a YAML DSL file, resolves it, and returns trace metadata.
func LoadAndResolveResult(ctx context.Context, path string, opts LoadOptions) (*LoadResult, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "read input %s", path)
	}
	var doc map[string]any
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, errors.Wrapf(err, "parse yaml %s", path)
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
