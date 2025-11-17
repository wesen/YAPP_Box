package cutouts

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

// Build converts typed cutouts entries into the YAPP array format.
// Note: Cutouts are distributed by face, so this returns a map.
func Build(items []CutoutsItem) (map[string][][]any, error) {
	byFace := map[string][][]any{
		"cutoutsFront": {},
		"cutoutsBack":  {},
		"cutoutsLeft":  {},
		"cutoutsRight": {},
		"cutoutsLid":   {},
		"cutoutsBase":  {},
	}

	for idx, item := range items {
		label := fmt.Sprintf("cutouts[%d]", idx)

		// Get shape flag and dimension usage
		shapeFlag, usesWidth, usesLength, usesRadius, err := cutoutShapeFlag(item.Shape)
		if err != nil {
			return nil, errors.Wrapf(err, "%s", label)
		}

		// Set unused dimensions to 0; allow missing dims to default to 0
		width := floatOrZero(item.Width)
		if !usesWidth {
			width = 0
		}
		length := floatOrZero(item.Length)
		if !usesLength {
			length = 0
		}
		radius := floatOrZero(item.Radius)
		if !usesRadius {
			radius = 0
		}

		// Determine position values based on face and new field names
		// All faces use from_face_left for horizontal position
		// Side faces use from_face_bottom for vertical position
		// Base/lid use from_face_back for depth position
		faceLower := strings.ToLower(strings.TrimSpace(item.Face))
		isSideFace := faceLower == "front" || faceLower == "back" || faceLower == "left" || faceLower == "right"
		isHorizFace := faceLower == "base" || faceLower == "lid" || faceLower == "top" || faceLower == "bottom"

		var pos0, pos1 float64

		if isSideFace {
			// Side faces: from_face_left → horizontal (pos0), from_face_bottom → vertical (pos1)
			pos0 = item.FromFaceLeft
			if item.FromFaceBottom != nil {
				pos1 = *item.FromFaceBottom
			}
		} else if isHorizFace {
			// Horizontal faces: from_face_left → Y, from_face_back → X
			// But YAPP arrays expect [X, Y] order, so we swap
			if item.FromFaceBack != nil {
				pos0 = *item.FromFaceBack // from_face_back becomes pos0 (X/back-to-front)
			}
			pos1 = item.FromFaceLeft // from_face_left becomes pos1 (Y/left-to-right)
		}

		// Build positional array with shape flag at position 5
		params := []any{
			pos0,
			pos1,
			width,
			length,
			radius,
			shapeFlag,
			ptrOrUndef(item.Depth),
			ptrOrUndef(item.Angle),
		}

		// For polygon shapes, add the preset shape after depth/angle
		if strings.ToLower(strings.TrimSpace(item.Shape)) == "polygon" {
			if item.Polygon == nil || *item.Polygon == "" {
				return nil, errors.Errorf("%s: polygon preset is required when shape=polygon", label)
			}
			presetFlag, err := polygonPresetFlag(*item.Polygon)
			if err != nil {
				return nil, errors.Wrapf(err, "%s", label)
			}
			params = append(params, presetFlag)
		}

		// Add mask preset and optional offsets/rotation (if provided)
		if maskElems, err := encodeMask(item); err != nil {
			return nil, errors.Wrapf(err, "%s", label)
		} else if len(maskElems) > 0 {
			params = append(params, maskElems...)
		}

		// Add coordinate and origin flags
		flags, err := encodeFlags(item)
		if err != nil {
			return nil, errors.Wrapf(err, "%s", label)
		}
		params = append(params, flags...)

		// Determine face array
		faceArray, err := faceArrayName(item.Face)
		if err != nil {
			return nil, errors.Wrapf(err, "%s", label)
		}

		byFace[faceArray] = append(byFace[faceArray], params)
	}

	return byFace, nil
}

func floatOrZero(ptr *float64) float64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func ptrOrUndef(ptr *float64) any {
	if ptr == nil {
		return scad.Undef
	}
	return *ptr
}

func cutoutShapeFlag(shape string) (scad.Raw, bool, bool, bool, error) {
	normalized := strings.ToLower(strings.TrimSpace(shape))
	switch normalized {
	case "rectangle":
		return scad.Raw("yappRectangle"), true, true, false, nil
	case "circle":
		return scad.Raw("yappCircle"), false, false, true, nil
	case "rounded_rect", "rounded-rect":
		return scad.Raw("yappRoundedRect"), true, true, true, nil
	case "circle_with_flats", "circle-with-flats":
		return scad.Raw("yappCircleWithFlats"), true, true, true, nil
	case "circle_with_key", "circle-with-key":
		return scad.Raw("yappCircleWithKey"), true, true, true, nil
	case "polygon":
		return scad.Raw("yappPolygon"), true, true, false, nil
	default:
		return "", false, false, false, errors.Errorf("unsupported cutout shape: %s", shape)
	}
}

func polygonPresetFlag(preset string) (scad.Raw, error) {
	switch strings.ToLower(strings.TrimSpace(preset)) {
	case "hexagon":
		return scad.Raw("shapeHexagon"), nil
	case "arrow":
		return scad.Raw("shapeArrow"), nil
	case "6pt_star", "6pt-star":
		return scad.Raw("shape6ptStar"), nil
	case "iso_triangle", "iso-triangle":
		return scad.Raw("shapeIsoTriangle"), nil
	case "iso_triangle2", "iso-triangle2":
		return scad.Raw("shapeIsoTriangle2"), nil
	case "triangle":
		return scad.Raw("shapeTriangle"), nil
	case "triangle2":
		return scad.Raw("shapeTriangle2"), nil
	default:
		return "", errors.Errorf("unsupported polygon preset: %s", preset)
	}
}

func maskPresetSymbol(preset string) (scad.Raw, error) {
	switch strings.ToLower(strings.TrimSpace(preset)) {
	case "honeycomb":
		return scad.Raw("maskHoneycomb"), nil
	case "hex_circles", "hex-circles":
		return scad.Raw("maskHexCircles"), nil
	case "circles":
		return scad.Raw("maskCircles"), nil
	case "squares":
		return scad.Raw("maskSquares"), nil
	case "bars":
		return scad.Raw("maskBars"), nil
	case "offset_bars", "offset-bars":
		return scad.Raw("maskOffsetBars"), nil
	default:
		return "", errors.Errorf("unsupported mask preset: %s", preset)
	}
}

// encodeMask returns zero or more elements to append to the cutout parameter list
// - First: the preset mask object (e.g., maskHoneycomb)
// - Then: optional offsets vector [[yappMaskDef, hOffset, vOffset, rotation]] if any offsets/rotation provided
func encodeMask(item CutoutsItem) ([]any, error) {
	if item.Mask == nil {
		return nil, nil
	}
	// Require preset when mask is specified
	if item.Mask.Preset == nil || strings.TrimSpace(*item.Mask.Preset) == "" {
		return nil, errors.Errorf("mask.preset is required when mask is specified")
	}
	var out []any
	presetSym, err := maskPresetSymbol(*item.Mask.Preset)
	if err != nil {
		return nil, err
	}
	out = append(out, presetSym)

	// Only emit offsets vector if any value is provided (defaults are 0)
	if item.Mask.OffsetX != nil || item.Mask.OffsetY != nil || item.Mask.Rotation != nil {
		offX := floatOrZero(item.Mask.OffsetX)
		offY := floatOrZero(item.Mask.OffsetY)
		rot := floatOrZero(item.Mask.Rotation)
		out = append(out, []any{scad.Raw("yappMaskDef"), offX, offY, rot})
	}
	return out, nil
}

func faceArrayName(face string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(face)) {
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

func encodeFlags(item CutoutsItem) ([]any, error) {
	var flags []any

	// Coordinate flag (default is pcb, only emit if not default)
	if flag, err := coordinateFlag(strOrDefault(item.Coordinate, "pcb")); err != nil {
		return nil, err
	} else if flag != "" {
		flags = append(flags, flag)
	}

	// Origin flag (default is global, only emit if not default)
	if flag, err := originFlag(strOrDefault(item.Origin, "global")); err != nil {
		return nil, err
	} else if flag != "" {
		flags = append(flags, flag)
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
	case "center":
		return scad.Raw("yappCenter"), nil
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
