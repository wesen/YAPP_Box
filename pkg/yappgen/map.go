package yappgen

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

// buildParams builds a positional parameter slice according to the schema.
// - For missing optional fields, it inserts scad.Undef
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
		out = append(out, scad.Undef)
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
	case scad.Raw:
		return string(t)
	case string:
		return fmt.Sprintf("%q", t)
	default:
		return fmt.Sprintf("%v", t)
	}
}
