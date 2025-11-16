package pcbstands

import (
	"fmt"
	"strings"

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

		flags, err := encodeFlags(item)
		if err != nil {
			return nil, errors.Wrapf(err, "%s: %v", label, err)
		}

		params = append(params, flags...)

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

func encodeFlags(item PcbStandsItem) ([]any, error) {
	var flags []any

	if flag, err := shellPartFlag(strOrDefault(item.ShellPart, "both")); err != nil {
		return nil, err
	} else if flag != "" {
		flags = append(flags, flag)
	}

	if flag, err := treatmentFlag(strOrDefault(item.Treatment, "pin")); err != nil {
		return nil, err
	} else if flag != "" {
		flags = append(flags, flag)
	}

	cornerFlags, err := cornerSelectors(strOrDefault(item.Corner, "single"))
	if err != nil {
		return nil, err
	}
	flags = append(flags, cornerFlags...)

	if flag, err := coordinateFlag(strOrDefault(item.Coordinate, "pcb")); err != nil {
		return nil, err
	} else if flag != "" {
		flags = append(flags, flag)
	}

	if boolVal(item.NoFillet) {
		flags = append(flags, scad.Raw("yappNoFillet"))
	}

	if name := strings.TrimSpace(strOrDefault(item.PcbName, "")); name != "" {
		flags = append(flags, []any{scad.Raw("yappPCBName"), name})
	}

	if boolVal(item.SelfThreading) {
		flags = append(flags, scad.Raw("yappSelfThreading"))
	}

	return flags, nil
}

func shellPartFlag(value string) (scad.Raw, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "both":
		return "", nil
	case "base_only":
		return scad.Raw("yappBaseOnly"), nil
	case "lid_only":
		return scad.Raw("yappLidOnly"), nil
	default:
		return "", errors.Errorf("invalid shell_part: %s", value)
	}
}

func treatmentFlag(value string) (scad.Raw, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "pin":
		return "", nil
	case "hole":
		return scad.Raw("yappHole"), nil
	case "top_pin":
		return scad.Raw("yappTopPin"), nil
	default:
		return "", errors.Errorf("invalid treatment: %s", value)
	}
}

func cornerSelectors(value string) ([]any, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "single":
		return nil, nil
	case "all":
		return []any{scad.Raw("yappAllCorners")}, nil
	case "front_left":
		return []any{scad.Raw("yappFrontLeft")}, nil
	case "front_right":
		return []any{scad.Raw("yappFrontRight")}, nil
	case "back_left":
		return []any{scad.Raw("yappBackLeft")}, nil
	case "back_right":
		return []any{scad.Raw("yappBackRight")}, nil
	default:
		return nil, errors.Errorf("invalid corner: %s", value)
	}
}

func coordinateFlag(value string) (scad.Raw, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "pcb":
		return "", nil
	case "box":
		return scad.Raw("yappCoordBox"), nil
	case "box_inside":
		return scad.Raw("yappCoordBoxInside"), nil
	default:
		return "", errors.Errorf("invalid coordinate: %s", value)
	}
}

func strOrDefault(v *string, def string) string {
	if v == nil {
		return def
	}
	return *v
}

func boolVal(v *bool) bool {
	return v != nil && *v
}
