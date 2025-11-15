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

// LoadAndResolve reads a YAML DSL file from disk, decodes it, and runs the
// expression resolver with the provided options.
func LoadAndResolve(ctx context.Context, path string, opts LoadOptions) (map[string]any, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "read input %s", path)
	}
	var doc map[string]any
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, errors.Wrapf(err, "parse yaml %s", path)
	}
	resolved, err := resolver.Resolve(ctx, doc, resolver.Options{
		MaxIterations: opts.MaxIterations,
		Strict:        opts.Strict,
	})
	if err != nil {
		return nil, err
	}
	return resolved, nil
}
