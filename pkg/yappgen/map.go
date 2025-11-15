package yappgen

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// buildParams builds a positional parameter slice according to the schema.
// - For missing optional fields, it inserts Undef
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
		out = append(out, Undef)
	}
	return out, nil
}

// buildPcbStands converts items to YAPP array entries.
func buildPcbStands(items []map[string]any) ([][]any, error) {
	var out [][]any
	for _, it := range items {
		params, err := buildParams(it, pcbStandsSchema)
		if err != nil {
			return nil, errors.Wrap(err, "pcb_stands")
		}
		out = append(out, params)
	}
	return out, nil
}

// buildConnectors converts items to YAPP array entries.
func buildConnectors(items []map[string]any) ([][]any, error) {
	var out [][]any
	for _, it := range items {
		params, err := buildParams(it, connectorsSchema)
		if err != nil {
			return nil, errors.Wrap(err, "connectors")
		}
		out = append(out, params)
	}
	return out, nil
}

// buildSnapJoins converts items to YAPP array entries.
// Appends the side flag as a raw token after positional params.
func buildSnapJoins(items []map[string]any) ([][]any, error) {
	var out [][]any
	for _, it := range items {
		params, err := buildParams(it, snapJoinsSchema)
		if err != nil {
			return nil, errors.Wrap(err, "snap_joins")
		}
		// side is required in DSL for snap joins
		sideStr, _ := it["side"].(string)
		if sideStr == "" {
			return nil, errors.Errorf("snap_joins item missing 'side'")
		}
		flag, ok := SnapSideFlag(sideStr)
		if !ok {
			return nil, errors.Errorf("invalid snap side: %s", sideStr)
		}
		params = append(params, flag)
		out = append(out, params)
	}
	return out, nil
}

// buildPushButtons converts DSL push button entries into the YAPP array format.
func buildPushButtons(items []map[string]any) ([][]any, error) {
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
			params = append(params, ScadRaw("yappNoFillet"))
		}

		out = append(out, params)
	}
	return out, nil
}

// buildCutoutParams converts a single cutout item to positional params + shape flag.
func buildCutoutParams(item map[string]any) ([]any, error) {
	shape, _ := item["shape"].(string)
	if shape == "" {
		return nil, errors.Errorf("cutout missing required 'shape'")
	}
	flag, useW, useL, useR, ok := ShapeFlag(shape)
	if !ok {
		return nil, errors.Errorf("unsupported cutout shape: %s", shape)
	}
	// Normalize coordinate keys: prefer explicit names; fall back if older keys present
	if _, ok := item["from_back"]; !ok {
		if v, ok2 := item["x"]; ok2 {
			item["from_back"] = v
		}
	}
	if _, ok := item["from_left"]; !ok {
		// Try y first; then z (older docs suggested z for side faces)
		if v, ok2 := item["y"]; ok2 {
			item["from_left"] = v
		} else if v, ok3 := item["z"]; ok3 {
			item["from_left"] = v
		}
	}
	// Initialize width/length/radius defaults based on shape usage
	if !useW {
		item["width"] = 0
	}
	if !useL {
		item["length"] = 0
	}
	if !useR {
		item["radius"] = 0
	}
	// Build base params
	params, err := buildParams(item, cutoutsSchema)
	if err != nil {
		return nil, err
	}
	// Insert shape flag after the first 5 numeric params
	params = append(params[:5], append([]any{flag}, params[5:]...)...)
	return params, nil
}

// distributeCutouts splits cutouts by face array name and builds params accordingly.
func distributeCutouts(cuts []Cutout) (byFace map[string][][]any, err error) {
	byFace = map[string][][]any{
		"cutoutsFront": {},
		"cutoutsBack":  {},
		"cutoutsLeft":  {},
		"cutoutsRight": {},
		"cutoutsLid":   {},
		"cutoutsBase":  {},
	}
	for _, c := range cuts {
		params, e := buildCutoutParams(c.Item)
		if e != nil {
			return nil, errors.Wrap(e, "cutouts")
		}
		key, e := faceArrayName(c.Face)
		if e != nil {
			return nil, e
		}
		byFace[key] = append(byFace[key], params)
	}
	return byFace, nil
}

func faceArrayName(face string) (string, error) {
	switch face {
	case "front":
		return "cutoutsFront", nil
	case "back":
		return "cutoutsBack", nil
	case "left":
		return "cutoutsLeft", nil
	case "right":
		return "cutoutsRight", nil
	case "top", "lid":
		return "cutoutsLid", nil
	case "bottom", "base":
		return "cutoutsBase", nil
	default:
		return "", errors.Errorf("invalid cutout face: %s", face)
	}
}

func formatAnyForDebug(v any) string {
	switch t := v.(type) {
	case ScadRaw:
		return string(t)
	case string:
		return fmt.Sprintf("%q", t)
	default:
		return fmt.Sprintf("%v", t)
	}
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
	return Undef
}

func pushButtonShapeTokens(shape, preset string) (ScadRaw, []any, error) {
	s := strings.ToLower(strings.TrimSpace(shape))
	if s == "" {
		s = "rectangle"
	}
	switch s {
	case "rectangle":
		return ScadRaw("yappRectangle"), nil, nil
	case "circle":
		return ScadRaw("yappCircle"), nil, nil
	case "rounded_rect", "rounded-rect", "roundedrect", "rounded":
		return ScadRaw("yappRoundedRect"), nil, nil
	case "circle_with_flats", "circle-with-flats", "circleflats":
		return ScadRaw("yappCircleWithFlats"), nil, nil
	case "circle_with_key", "circle-with-key", "circlekey":
		return ScadRaw("yappCircleWithKey"), nil, nil
	case "polygon":
		token, err := pushButtonPolygonPreset(preset)
		if err != nil {
			return "", nil, err
		}
		if token == "" {
			return ScadRaw("yappPolygon"), nil, nil
		}
		return ScadRaw("yappPolygon"), []any{token}, nil
	case "arrow", "triangle", "triangle2", "hexagon", "iso_triangle", "iso-triangle", "iso_triangle2", "iso-triangle2", "6pt_star", "6pt-star", "star6":
		token, err := pushButtonPolygonPreset(s)
		if err != nil {
			return "", nil, err
		}
		return ScadRaw("yappPolygon"), []any{token}, nil
	default:
		// Allow direct yapp token names (e.g., yappPolygon) or shape* tokens
		if strings.HasPrefix(s, "yapp") {
			return ScadRaw(shape), nil, nil
		}
		if strings.HasPrefix(shape, "shape") {
			return ScadRaw("yappPolygon"), []any{ScadRaw(shape)}, nil
		}
		return "", nil, errors.Errorf("unsupported push button shape: %s", shape)
	}
}

func pushButtonPolygonPreset(name string) (ScadRaw, error) {
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
		return ScadRaw(n), nil
	}
	return "", errors.Errorf("unknown polygon preset: %s", name)
}

var polygonPresetTokens = map[string]ScadRaw{
	"shapeArrow":        ScadRaw("shapeArrow"),
	"shapeTriangle":     ScadRaw("shapeTriangle"),
	"shapeTriangle2":    ScadRaw("shapeTriangle2"),
	"shapeIsoTriangle":  ScadRaw("shapeIsoTriangle"),
	"shapeIsoTriangle2": ScadRaw("shapeIsoTriangle2"),
	"shapeHexagon":      ScadRaw("shapeHexagon"),
	"shape6ptStar":      ScadRaw("shape6ptStar"),
}

var polygonPresetAliases = map[string]ScadRaw{
	"arrow":         ScadRaw("shapeArrow"),
	"triangle":      ScadRaw("shapeTriangle"),
	"triangle2":     ScadRaw("shapeTriangle2"),
	"iso_triangle":  ScadRaw("shapeIsoTriangle"),
	"iso-triangle":  ScadRaw("shapeIsoTriangle"),
	"iso_triangle2": ScadRaw("shapeIsoTriangle2"),
	"iso-triangle2": ScadRaw("shapeIsoTriangle2"),
	"hexagon":       ScadRaw("shapeHexagon"),
	"6pt_star":      ScadRaw("shape6ptStar"),
	"6pt-star":      ScadRaw("shape6ptStar"),
	"star6":         ScadRaw("shape6ptStar"),
}

func pushButtonCoordinateFlag(coord string) (ScadRaw, error) {
	switch strings.ToLower(strings.TrimSpace(coord)) {
	case "pcb":
		return ScadRaw("yappCoordPCB"), nil
	case "box":
		return ScadRaw("yappCoordBox"), nil
	case "box_inside", "box-inside", "inside":
		return ScadRaw("yappCoordBoxInside"), nil
	default:
		return "", errors.Errorf("invalid coordinate: %s", coord)
	}
}

func pushButtonOriginFlag(origin string) (ScadRaw, error) {
	switch strings.ToLower(strings.TrimSpace(origin)) {
	case "global":
		return ScadRaw("yappGlobalOrigin"), nil
	case "left":
		return ScadRaw("yappLeftOrigin"), nil
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
