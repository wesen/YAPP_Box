package connectors

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

// Build converts DSL connectors entries into the YAPP array format.
func Build(items []map[string]any) ([][]any, error) {
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}

	var out [][]any
	for idx, item := range typed {
		label := fmt.Sprintf("connectors[%d]", idx)

		// Build positional array
		params := []any{
			item.X,
			item.Y,
			item.StandHeight,
			item.ScrewD,
			item.ScrewHeadD,
			item.InsertD,
			item.OutsideD,
			ptrOrUndef(item.InsertDepth),
			ptrOrUndef(item.PcbGap),
			ptrOrUndef(item.FilletRadius),
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

func encodeFlags(item ConnectorsItem) ([]any, error) {
	var flags []any

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

	if boolVal(item.Countersink) {
		flags = append(flags, scad.Raw("yappCountersink"))
	}

	if name := strings.TrimSpace(strOrDefault(item.PcbName, "")); name != "" {
		flags = append(flags, []any{scad.Raw("yappPCBName"), name})
	}

	if boolVal(item.ThroughLid) {
		flags = append(flags, scad.Raw("yappThroughLid"))
	}

	if boolVal(item.SelfThreading) {
		flags = append(flags, scad.Raw("yappSelfThreading"))
	}

	if boolVal(item.NoInternalFillet) {
		flags = append(flags, scad.Raw("yappNoInternalFillet"))
	}

	return flags, nil
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
