package pushbuttons

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

// Build converts DSL push button entries into the YAPP array format.
func Build(items []map[string]any) ([][]any, error) {
	typedItems, err := Decode(items)
	if err != nil {
		return nil, err
	}

	var out [][]any
	for idx, item := range typedItems {
		label := fmt.Sprintf("push_buttons[%d]", idx)
		if item.Name != nil && strings.TrimSpace(*item.Name) != "" {
			label = fmt.Sprintf("%s (%s)", label, *item.Name)
		}

		params, err := buildPushButtonParams(label, &item)
		if err != nil {
			return nil, err
		}
		out = append(out, params)
	}
	return out, nil
}

func buildPushButtonParams(label string, item *PushButtonsItem) ([]any, error) {
	shape := stringOr(item.Shape, "rectangle")
	preset := firstNonEmptyPointer(
		item.Polygon,
		item.PolygonPreset,
		item.ShapePreset,
	)
	shapeFlag, shapeExtras, err := pushButtonShapeTokens(shape, preset)
	if err != nil {
		return nil, errors.Wrapf(err, "%s", label)
	}

	angle, hasAngle := floatFromPtr(item.Angle)
	filletRadius, hasFillet := floatFromPtr(item.FilletRadius)
	buttonWall, hasWall := floatFromPtr(item.Lid.Wall)
	plateThickness, hasPlate := floatFromPtr(item.Lid.PlateThickness)
	buttonSlack, hasSlack := floatFromPtr(item.Lid.Slack)
	snapSlack, hasSnapSlack := floatFromPtr(item.Lid.SnapSlack)
	heightToPCB, hasHeightToPCB := floatFromPtr(item.Switch.TopHeight)

	params := []any{
		item.X,
		item.Y,
		item.Cap.Length,
		item.Cap.Width,
		item.Cap.Radius,
		item.Lid.Protrusion,
		item.Switch.Height,
		item.Switch.Travel,
		item.Switch.PoleDiameter,
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

	if coord := stringOr(item.Coordinate, ""); coord != "" {
		flag, err := pushButtonCoordinateFlag(coord)
		if err != nil {
			return nil, errors.Wrapf(err, "%s.coordinate", label)
		}
		params = append(params, flag)
	}

	if origin := stringOr(item.Origin, ""); origin != "" {
		flag, err := pushButtonOriginFlag(origin)
		if err != nil {
			return nil, errors.Wrapf(err, "%s.origin", label)
		}
		params = append(params, flag)
	}

	// PCB name flag for multi-board projects
	if pcbName := stringOr(item.PcbName, ""); pcbName != "" {
		params = append(params, []any{scad.Raw("yappPCBName"), pcbName})
	}

	if noFilletSet(item) {
		params = append(params, scad.Raw("yappNoFillet"))
	}

	return params, nil
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
	case "alt":
		return scad.Raw("yappAltOrigin"), nil
	default:
		return "", errors.Errorf("invalid origin: %s", origin)
	}
}

func floatFromPtr(ptr *float64) (float64, bool) {
	if ptr == nil {
		return 0, false
	}
	return *ptr, true
}

func stringOr(ptr *string, fallback string) string {
	if ptr == nil {
		return fallback
	}
	if strings.TrimSpace(*ptr) == "" {
		return fallback
	}
	return strings.TrimSpace(*ptr)
}

func firstNonEmptyPointer(values ...*string) string {
	for _, ptr := range values {
		if ptr != nil && strings.TrimSpace(*ptr) != "" {
			return strings.TrimSpace(*ptr)
		}
	}
	return ""
}

func noFilletSet(item *PushButtonsItem) bool {
	if item == nil {
		return false
	}
	if item.Lid.NoFillet != nil && *item.Lid.NoFillet {
		return true
	}
	if item.NoFillet != nil && *item.NoFillet {
		return true
	}
	return false
}
