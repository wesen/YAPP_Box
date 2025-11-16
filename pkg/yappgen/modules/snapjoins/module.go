package snapjoins

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

// Build converts DSL snap_joins entries into the YAPP array format.
func Build(items []map[string]any) ([][]any, error) {
	var out [][]any
	for idx, it := range items {
		label := fmt.Sprintf("snap_joins[%d]", idx)

		// Unmarshal into typed struct
		data, err := yaml.Marshal(it)
		if err != nil {
			return nil, errors.Wrapf(err, "%s: marshal", label)
		}

		var item SnapJoinsItem
		if err := yaml.Unmarshal(data, &item); err != nil {
			return nil, errors.Wrapf(err, "%s: unmarshal", label)
		}

		// Apply defaults
		item.ApplyDefaults()

		// Custom validation
		if err := item.CustomValidate(); err != nil {
			return nil, errors.Wrapf(err, "%s", label)
		}

		// Validate side enum
		sideFlag, err := snapSideFlag(item.Side)
		if err != nil {
			return nil, errors.Wrapf(err, "%s", label)
		}

		// Build positional array + side flag
		params := []any{
			item.Pos,
			item.Width,
			sideFlag,
		}

		out = append(out, params)
	}
	return out, nil
}

func snapSideFlag(side string) (scad.Raw, error) {
	switch strings.ToLower(strings.TrimSpace(side)) {
	case "left":
		return scad.Raw("yappLeft"), nil
	case "right":
		return scad.Raw("yappRight"), nil
	case "front":
		return scad.Raw("yappFront"), nil
	case "back":
		return scad.Raw("yappBack"), nil
	default:
		return "", errors.Errorf("invalid snap side: %s", side)
	}
}
