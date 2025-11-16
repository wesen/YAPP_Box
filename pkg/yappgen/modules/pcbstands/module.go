package pcbstands

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

// Build converts DSL pcb_stands entries into the YAPP array format.
func Build(items []map[string]any) ([][]any, error) {
	var out [][]any
	for idx, it := range items {
		label := fmt.Sprintf("pcb_stands[%d]", idx)

		// Unmarshal into typed struct
		data, err := yaml.Marshal(it)
		if err != nil {
			return nil, errors.Wrapf(err, "%s: marshal", label)
		}

		var item PcbStandsItem
		if err := yaml.Unmarshal(data, &item); err != nil {
			return nil, errors.Wrapf(err, "%s: unmarshal", label)
		}

		// Apply defaults
		item.ApplyDefaults()

		// Custom validation
		if err := item.CustomValidate(); err != nil {
			return nil, errors.Wrapf(err, "%s", label)
		}

		// Build positional array
		params := []any{
			item.X,
			item.Y,
			ptrOrUndef(item.Height),
			ptrOrUndef(item.PcbGap),
			ptrOrUndef(item.Diameter),
			ptrOrUndef(item.PinDiameter),
			ptrOrUndef(item.HoleSlack),
			ptrOrUndef(item.FilletRadius),
			ptrOrUndef(item.PinLength),
		}

		out = append(out, params)
	}
	return out, nil
}

func ptrOrUndef(ptr *float64) any {
	if ptr == nil {
		return scad.Undef
	}
	return *ptr
}
