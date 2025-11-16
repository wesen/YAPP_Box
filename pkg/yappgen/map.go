package yappgen

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

// buildParams builds a positional parameter slice according to the schema.
// - For missing optional fields, it inserts scad.Undef
// - For missing required fields, it returns an error
func buildParams(item map[string]any, schema []ParamSpec) ([]any, error) {
	out := make([]any, 0, len(schema))
	for _, p := range schema {
		v, has := item[p.Name]
		if has {
			out = append(out, v)
			continue
		}
		if p.Required {
			return nil, errors.Errorf("missing required parameter: %s", p.Name)
		}
		out = append(out, scad.Undef)
	}
	return out, nil
}

func formatAnyForDebug(v any) string {
	switch t := v.(type) {
	case scad.Raw:
		return string(t)
	case string:
		return fmt.Sprintf("%q", t)
	default:
		return fmt.Sprintf("%v", t)
	}
}
