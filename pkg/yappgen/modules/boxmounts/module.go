package boxmounts

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

// Build converts typed box_mounts entries into the YAPP array format.
func Build(items []BoxMountsItem) ([][]any, error) {
	var out [][]any
	for idx, item := range items {
		label := fmt.Sprintf("box_mounts[%d]", idx)

		if err := validateFaces(item.Faces); err != nil {
			return nil, errors.Wrapf(err, "%s: %v", label, err)
		}

		// Note: negative slot_width is valid in YAPP (indicates vertical orientation)
		posVal := positionValue(item.Pos, item.Offset)

		params := []any{
			posVal,
			item.ScrewD,
			item.SlotWidth,
			item.Height,
			ptrOrUndef(item.FilletRadius),
		}

		faceFlags := facesToFlags(item.Faces)
		params = append(params, faceFlags...)

		if boolVal(item.NoFillet) {
			params = append(params, scad.Raw("yappNoFillet"))
		}

		if shellFlag := shellPartFlag(item.ShellPart); shellFlag != "" {
			params = append(params, shellFlag)
		}

		if alignmentFlag := alignmentFlag(item.Alignment); alignmentFlag != "" {
			params = append(params, alignmentFlag)
		}

		if originFlag := originFlag(item.Origin); originFlag != "" {
			params = append(params, originFlag)
		}

		out = append(out, params)
	}
	return out, nil
}

// Removed: negative slot_width is valid (indicates vertical slot orientation in YAPP)

func ptrOrUndef(ptr *float64) any {
	if ptr == nil {
		return scad.Undef
	}
	return *ptr
}

func positionValue(pos float64, offset *float64) any {
	if offset == nil {
		return pos
	}
	return []any{pos, *offset}
}

func validateFaces(faces BoxMountsItemFaces) error {
	if boolVal(faces.Left) || boolVal(faces.Right) || boolVal(faces.Front) || boolVal(faces.Back) {
		return nil
	}
	return errors.New("faces must include at least one wall")
}

func facesToFlags(faces BoxMountsItemFaces) []any {
	var flags []any
	if boolVal(faces.Left) {
		flags = append(flags, scad.Raw("yappLeft"))
	}
	if boolVal(faces.Right) {
		flags = append(flags, scad.Raw("yappRight"))
	}
	if boolVal(faces.Front) {
		flags = append(flags, scad.Raw("yappFront"))
	}
	if boolVal(faces.Back) {
		flags = append(flags, scad.Raw("yappBack"))
	}
	return flags
}

func shellPartFlag(value *string) scad.Raw {
	if value == nil {
		return ""
	}
	switch strings.ToLower(strings.TrimSpace(*value)) {
	case "lid":
		return scad.Raw("yappLid")
	default:
		return ""
	}
}

func alignmentFlag(value *string) scad.Raw {
	if value == nil {
		return ""
	}
	if strings.EqualFold(*value, "center") {
		return scad.Raw("yappCenter")
	}
	return ""
}

func originFlag(value *string) scad.Raw {
	if value == nil {
		return ""
	}
	if strings.EqualFold(*value, "alt") {
		return scad.Raw("yappAltOrigin")
	}
	return ""
}

func boolVal(v *bool) bool {
	return v != nil && *v
}
