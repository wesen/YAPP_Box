package cutouts

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

// Build converts DSL cutouts entries into the YAPP array format.
// Note: Cutouts are distributed by face, so this returns a map.
func Build(items []map[string]any) (map[string][][]any, error) {
	byFace := map[string][][]any{
		"cutoutsFront": {},
		"cutoutsBack":  {},
		"cutoutsLeft":  {},
		"cutoutsRight": {},
		"cutoutsLid":   {},
		"cutoutsBase":  {},
	}

	for idx, it := range items {
		label := fmt.Sprintf("cutouts[%d]", idx)

		// Unmarshal into typed struct
		data, err := yaml.Marshal(it)
		if err != nil {
			return nil, errors.Wrapf(err, "%s: marshal", label)
		}

		var item CutoutsItem
		if err := yaml.Unmarshal(data, &item); err != nil {
			return nil, errors.Wrapf(err, "%s: unmarshal", label)
		}

		// Apply defaults
		item.ApplyDefaults()

		// Custom validation
		if err := item.CustomValidate(); err != nil {
			return nil, errors.Wrapf(err, "%s", label)
		}

		// Get shape flag and dimension usage
		shapeFlag, usesWidth, usesLength, usesRadius, err := cutoutShapeFlag(item.Shape)
		if err != nil {
			return nil, errors.Wrapf(err, "%s", label)
		}

		// Set unused dimensions to 0
		width := item.Width
		if !usesWidth {
			width = 0
		}
		length := item.Length
		if !usesLength {
			length = 0
		}
		radius := item.Radius
		if !usesRadius {
			radius = 0
		}

		// Determine position values based on face
		// For side faces (front/back/left/right), pos_z can be used instead of from_left
		// for clarity (pos_z = vertical position from bottom)
		pos0 := item.FromBack
		pos1 := item.FromLeft
		
		// Check if this is a side face and pos_z is provided
		faceLower := strings.ToLower(strings.TrimSpace(item.Face))
		isSideFace := faceLower == "front" || faceLower == "back" || faceLower == "left" || faceLower == "right"
		if isSideFace && item.PosZ != nil {
			pos1 = *item.PosZ
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

		// Determine face array
		faceArray, err := faceArrayName(item.Face)
		if err != nil {
			return nil, errors.Wrapf(err, "%s", label)
		}

		byFace[faceArray] = append(byFace[faceArray], params)
	}

	return byFace, nil
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
