package pcbstands

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

// Build converts DSL pcb_stands entries into the YAPP array format.
func Build(items []map[string]any) ([][]any, error) {
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}

	var out [][]any
	for idx, item := range typed {
		label := fmt.Sprintf("pcb_stands[%d]", idx)

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

	// Prefer plural corners (multiple allowed); fall back to single corner
	if list := parseCorners(item.Corners); len(list) > 0 {
		cf, err := cornersListToFlags(list)
		if err != nil {
			return nil, err
		}
		flags = append(flags, cf...)
	} else {
		cornerFlags, err := cornerSelectors(strOrDefault(item.Corner, "single"))
		if err != nil {
			return nil, err
		}
		flags = append(flags, cornerFlags...)
	}

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

func parseCorners(raw any) []string {
	if raw == nil {
		return nil
	}
	// Dereference *any produced by schemagen
	if p, ok := raw.(*any); ok && p != nil {
		raw = *p
	}
	// Expect []any of strings
	switch t := raw.(type) {
	case []any:
		var out []string
		for _, v := range t {
			if s, ok := v.(string); ok {
				out = append(out, s)
			}
		}
		return out
	case []string:
		return t
	default:
		return nil
	}
}

func cornersListToFlags(list []string) ([]any, error) {
	var flags []any
	for _, v := range list {
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "all":
			flags = append(flags, scad.Raw("yappAllCorners"))
		case "front_left":
			flags = append(flags, scad.Raw("yappFrontLeft"))
		case "front_right":
			flags = append(flags, scad.Raw("yappFrontRight"))
		case "back_left":
			flags = append(flags, scad.Raw("yappBackLeft"))
		case "back_right":
			flags = append(flags, scad.Raw("yappBackRight"))
		default:
			return nil, errors.Errorf("invalid corners entry: %s", v)
		}
	}
	return flags, nil
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
