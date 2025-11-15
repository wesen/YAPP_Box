package pushbuttons

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

// Build converts DSL push button entries into the YAPP array format.
func Build(items []map[string]any) ([][]any, error) {
	var out [][]any
	for idx, it := range items {
		label := fmt.Sprintf("push_buttons[%d]", idx)
		if name, _ := it["name"].(string); name != "" {
			label = fmt.Sprintf("%s (%s)", label, name)
		}
		capMap, err := getMapField(it, "cap", label, true)
		if err != nil {
			return nil, err
		}
		lidMap, err := getMapField(it, "lid", label, true)
		if err != nil {
			return nil, err
		}
		switchMap, err := getMapField(it, "switch", label, true)
		if err != nil {
			return nil, err
		}

		x, err := requireNumber(it, "x", label)
		if err != nil {
			return nil, err
		}
		y, err := requireNumber(it, "y", label)
		if err != nil {
			return nil, err
		}
		capLength, err := requireNumber(capMap, "length", label+".cap")
		if err != nil {
			return nil, err
		}
		capWidth, err := requireNumber(capMap, "width", label+".cap")
		if err != nil {
			return nil, err
		}
		capRadius, err := requireNumber(capMap, "radius", label+".cap")
		if err != nil {
			return nil, err
		}
		protrusion, err := requireNumber(lidMap, "protrusion", label+".lid")
		if err != nil {
			return nil, err
		}
		switchHeight, err := requireNumber(switchMap, "height", label+".switch")
		if err != nil {
			return nil, err
		}
		switchTravel, err := requireNumber(switchMap, "travel", label+".switch")
		if err != nil {
			return nil, err
		}
		poleDiameter, err := requireNumber(switchMap, "pole_diameter", label+".switch")
		if err != nil {
			return nil, err
		}
		heightToPCB, hasHeightToPCB, err := optionalNumber(switchMap, "top_height", label+".switch")
		if err != nil {
			return nil, err
		}

		shapeStr, _ := optionalString(it, "shape")
		presetStr := firstNonEmptyString(
			getOptionalString(it, "polygon"),
			getOptionalString(it, "polygon_preset"),
			getOptionalString(it, "shape_preset"),
		)
		shapeFlag, shapeExtras, err := pushButtonShapeTokens(shapeStr, presetStr)
		if err != nil {
			return nil, errors.Wrapf(err, "%s", label)
		}

		angle, hasAngle, err := optionalNumber(it, "angle", label)
		if err != nil {
			return nil, err
		}
		filletRadius, hasFillet, err := optionalNumber(it, "fillet_radius", label)
		if err != nil {
			return nil, err
		}
		buttonWall, hasWall, err := optionalNumber(lidMap, "wall", label+".lid")
		if err != nil {
			return nil, err
		}
		plateThickness, hasPlate, err := optionalNumber(lidMap, "plate_thickness", label+".lid")
		if err != nil {
			return nil, err
		}
		buttonSlack, hasSlack, err := optionalNumber(lidMap, "slack", label+".lid")
		if err != nil {
			return nil, err
		}
		snapSlack, hasSnapSlack, err := optionalNumber(lidMap, "snap_slack", label+".lid")
		if err != nil {
			return nil, err
		}

		coordStr, _ := optionalString(it, "coordinate")
		originStr, _ := optionalString(it, "origin")
		noFillet, _ := optionalBool(lidMap, "no_fillet")
		if !noFillet {
			// allow root override
			noFillet, _ = optionalBool(it, "no_fillet")
		}

		params := []any{
			x,
			y,
			capLength,
			capWidth,
			capRadius,
			protrusion,
			switchHeight,
			switchTravel,
			poleDiameter,
			valueOrUndef(heightToPCB, hasHeightToPCB),
			shapeFlag,
			valueOrUndef(angle, hasAngle),
			valueOrUndef(filletRadius, hasFillet),
			valueOrUndef(buttonWall, hasWall),
			valueOrUndef(plateThickness, hasPlate),
			valueOrUndef(buttonSlack, hasSlack),
			valueOrUndef(snapSlack, hasSnapSlack),
		}

		if len(shapeExtras) > 0 {
			params = append(params, shapeExtras...)
		}
		if coordStr != "" {
			flag, err := pushButtonCoordinateFlag(coordStr)
			if err != nil {
				return nil, errors.Wrapf(err, "%s.coordinate", label)
			}
			params = append(params, flag)
		}
		if originStr != "" {
			flag, err := pushButtonOriginFlag(originStr)
			if err != nil {
				return nil, errors.Wrapf(err, "%s.origin", label)
			}
			params = append(params, flag)
		}
		if noFillet {
			params = append(params, scad.Raw("yappNoFillet"))
		}

		out = append(out, params)
	}
	return out, nil
}

func requireNumber(m map[string]any, key, ctx string) (float64, error) {
	v, ok := m[key]
	if !ok {
		return 0, errors.Errorf("%s missing required field '%s'", ctx, key)
	}
	if v == nil {
		return 0, errors.Errorf("%s field '%s' cannot be null", ctx, key)
	}
	f, err := toFloat64(v)
	if err != nil {
		return 0, errors.Wrapf(err, "%s.%s", ctx, key)
	}
	return f, nil
}

func optionalNumber(m map[string]any, key, ctx string) (float64, bool, error) {
	v, ok := m[key]
	if !ok || v == nil {
		return 0, false, nil
	}
	f, err := toFloat64(v)
	if err != nil {
		return 0, false, errors.Wrapf(err, "%s.%s", ctx, key)
	}
	return f, true, nil
}

func optionalString(m map[string]any, key string) (string, bool) {
	v, ok := m[key]
	if !ok || v == nil {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	return strings.TrimSpace(s), true
}

func getOptionalString(m map[string]any, key string) string {
	if s, ok := optionalString(m, key); ok {
		return s
	}
	return ""
}

func optionalBool(m map[string]any, key string) (bool, bool) {
	v, ok := m[key]
	if !ok || v == nil {
		return false, false
	}
	b, ok := v.(bool)
	return b, ok
}

func getMapField(root map[string]any, key, ctx string, required bool) (map[string]any, error) {
	v, ok := root[key]
	if !ok || v == nil {
		if required {
			return nil, errors.Errorf("%s missing required object '%s'", ctx, key)
		}
		return nil, nil
	}
	m, ok := v.(map[string]any)
	if !ok {
		return nil, errors.Errorf("%s.%s must be an object", ctx, key)
	}
	return m, nil
}

func toFloat64(v any) (float64, error) {
	switch t := v.(type) {
	case float64:
		return t, nil
	case float32:
		return float64(t), nil
	case int:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case uint:
		return float64(t), nil
	case uint32:
		return float64(t), nil
	case uint64:
		return float64(t), nil
	default:
		return 0, errors.Errorf("expected number, got %T", v)
	}
}

func valueOrUndef(val float64, ok bool) any {
	if ok {
		return val
	}
	return scad.Undef
}

func pushButtonShapeTokens(shape, preset string) (scad.Raw, []any, error) {
	s := strings.ToLower(strings.TrimSpace(shape))
	if s == "" {
		s = "rectangle"
	}
	switch s {
	case "rectangle":
		return scad.Raw("yappRectangle"), nil, nil
	case "circle":
		return scad.Raw("yappCircle"), nil, nil
	case "rounded_rect", "rounded-rect", "roundedrect", "rounded":
		return scad.Raw("yappRoundedRect"), nil, nil
	case "circle_with_flats", "circle-with-flats", "circleflats":
		return scad.Raw("yappCircleWithFlats"), nil, nil
	case "circle_with_key", "circle-with-key", "circlekey":
		return scad.Raw("yappCircleWithKey"), nil, nil
	case "polygon":
		token, err := pushButtonPolygonPreset(preset)
		if err != nil {
			return "", nil, err
		}
		if token == "" {
			return scad.Raw("yappPolygon"), nil, nil
		}
		return scad.Raw("yappPolygon"), []any{token}, nil
	case "arrow", "triangle", "triangle2", "hexagon", "iso_triangle", "iso-triangle", "iso_triangle2", "iso-triangle2", "6pt_star", "6pt-star", "star6":
		token, err := pushButtonPolygonPreset(s)
		if err != nil {
			return "", nil, err
		}
		return scad.Raw("yappPolygon"), []any{token}, nil
	default:
		if strings.HasPrefix(s, "yapp") {
			return scad.Raw(shape), nil, nil
		}
		if strings.HasPrefix(shape, "shape") {
			return scad.Raw("yappPolygon"), []any{scad.Raw(shape)}, nil
		}
		return "", nil, errors.Errorf("unsupported push button shape: %s", shape)
	}
}

func pushButtonPolygonPreset(name string) (scad.Raw, error) {
	n := strings.TrimSpace(name)
	if n == "" {
		return "", errors.New("polygon preset is required when shape=polygon")
	}
	if token, ok := polygonPresetTokens[n]; ok {
		return token, nil
	}
	lower := strings.ToLower(n)
	if token, ok := polygonPresetAliases[lower]; ok {
		return token, nil
	}
	if strings.HasPrefix(n, "shape") {
		return scad.Raw(n), nil
	}
	return "", errors.Errorf("unknown polygon preset: %s", name)
}

var polygonPresetTokens = map[string]scad.Raw{
	"shapeArrow":        scad.Raw("shapeArrow"),
	"shapeTriangle":     scad.Raw("shapeTriangle"),
	"shapeTriangle2":    scad.Raw("shapeTriangle2"),
	"shapeIsoTriangle":  scad.Raw("shapeIsoTriangle"),
	"shapeIsoTriangle2": scad.Raw("shapeIsoTriangle2"),
	"shapeHexagon":      scad.Raw("shapeHexagon"),
	"shape6ptStar":      scad.Raw("shape6ptStar"),
}

var polygonPresetAliases = map[string]scad.Raw{
	"arrow":         scad.Raw("shapeArrow"),
	"triangle":      scad.Raw("shapeTriangle"),
	"triangle2":     scad.Raw("shapeTriangle2"),
	"iso_triangle":  scad.Raw("shapeIsoTriangle"),
	"iso-triangle":  scad.Raw("shapeIsoTriangle"),
	"iso_triangle2": scad.Raw("shapeIsoTriangle2"),
	"iso-triangle2": scad.Raw("shapeIsoTriangle2"),
	"hexagon":       scad.Raw("shapeHexagon"),
	"6pt_star":      scad.Raw("shape6ptStar"),
	"6pt-star":      scad.Raw("shape6ptStar"),
	"star6":         scad.Raw("shape6ptStar"),
}

func pushButtonCoordinateFlag(coord string) (scad.Raw, error) {
	switch strings.ToLower(strings.TrimSpace(coord)) {
	case "pcb":
		return scad.Raw("yappCoordPCB"), nil
	case "box":
		return scad.Raw("yappCoordBox"), nil
	case "box_inside", "box-inside", "inside":
		return scad.Raw("yappCoordBoxInside"), nil
	default:
		return "", errors.Errorf("invalid coordinate: %s", coord)
	}
}

func pushButtonOriginFlag(origin string) (scad.Raw, error) {
	switch strings.ToLower(strings.TrimSpace(origin)) {
	case "global":
		return scad.Raw("yappGlobalOrigin"), nil
	case "left":
		return scad.Raw("yappLeftOrigin"), nil
	default:
		return "", errors.Errorf("invalid origin: %s", origin)
	}
}

func firstNonEmptyString(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
