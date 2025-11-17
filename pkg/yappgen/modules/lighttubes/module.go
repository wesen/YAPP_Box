package lighttubes

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

// Build converts DSL light_tubes entries into the YAPP array format.
func Build(items []map[string]any) ([][]any, error) {
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}

	var out [][]any
	for idx, item := range typed {
		label := fmt.Sprintf("light_tubes[%d]", idx)

		// Build positional parameters matching YAPP order:
		// p(0) = posx
		// p(1) = posy
		// p(2) = tubeLength
		// p(3) = tubeWidth
		// p(4) = tubeWall
		// p(5) = gapAbovePcb
		// p(6) = tubeType (yappCircle|yappRectangle)
		// p(7) = lensThickness (optional, default 0)
		// p(8) = height (optional, default standoffHeight+pcbThickness)
		// p(9) = filletRadius (optional, default 0)
		params := []any{
			item.X,                // [0] posx
			item.Y,                // [1] posy
			item.TubeLength,       // [2] tubeLength
			item.TubeWidth,        // [3] tubeWidth
			item.TubeWall,         // [4] tubeWall
			item.GapAbovePcb,      // [5] gapAbovePcb
			shapeFlag(item.Shape), // [6] tubeType
		}

		// Add optional positional parameters
		if item.LensThickness != nil && *item.LensThickness != 0 {
			params = append(params, *item.LensThickness) // [7] lensThickness
		}
		if item.Height != nil {
			params = append(params, *item.Height) // [8] height
		}
		if item.FilletRadius != nil {
			params = append(params, *item.FilletRadius) // [9] filletRadius
		}

		// Add flags
		flags, err := encodeFlags(item)
		if err != nil {
			return nil, errors.Wrapf(err, "%s", label)
		}
		params = append(params, flags...)

		out = append(out, params)
	}
	return out, nil
}

func shapeFlag(shape string) scad.Raw {
	switch strings.ToLower(strings.TrimSpace(shape)) {
	case "rectangle":
		return scad.Raw("yappRectangle")
	case "circle":
		return scad.Raw("yappCircle")
	default:
		return scad.Raw("yappCircle")
	}
}

func encodeFlags(item LightTubesItem) ([]any, error) {
	var flags []any

	// Coordinate flag (default is pcb, so only emit if not default)
	if flag, err := coordinateFlag(strOrDefault(item.Coordinate, "pcb")); err != nil {
		return nil, err
	} else if flag != "" {
		flags = append(flags, flag)
	}

	// Origin flag (default is global, so only emit if not default)
	if flag, err := originFlag(strOrDefault(item.Origin, "global")); err != nil {
		return nil, err
	} else if flag != "" {
		flags = append(flags, flag)
	}

	// No fillet flag
	if boolVal(item.NoFillet) {
		flags = append(flags, scad.Raw("yappNoFillet"))
	}

	// PCB name flag
	if name := strings.TrimSpace(strOrDefault(item.PcbName, "")); name != "" {
		flags = append(flags, []any{scad.Raw("yappPCBName"), name})
	}

	return flags, nil
}

func coordinateFlag(value string) (scad.Raw, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "pcb":
		return "", nil // Default, don't emit
	case "box":
		return scad.Raw("yappCoordBox"), nil
	case "box_inside":
		return scad.Raw("yappCoordBoxInside"), nil
	default:
		return "", errors.Errorf("invalid coordinate: %s", value)
	}
}

func originFlag(value string) (scad.Raw, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "global":
		return "", nil // Default, don't emit
	case "alt":
		return scad.Raw("yappAltOrigin"), nil
	default:
		return "", errors.Errorf("invalid origin: %s", value)
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
