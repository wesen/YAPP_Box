package yappgen

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// Model holds the normalized, resolved configuration needed to emit a YAPP SCAD file.
type Model struct {
	ProjectName string

	// Global/box parameters (subset for MVP)
	PcbLength         float64
	PcbWidth          float64
	PcbThickness      float64
	StandoffHeight    float64
	StandoffDiameter  float64 // optional (<=0 means use YAPP default)
	StandoffPinDia    float64 // optional
	StandoffHoleSlack float64 // optional

	// Enclosure parameters (optional; <=0 means use YAPP defaults)
	WallThickness      float64
	BasePlaneThickness float64
	LidPlaneThickness  float64
	RoundRadius        float64
	PaddingFront       float64
	PaddingBack        float64
	PaddingLeft        float64
	PaddingRight       float64

	// Features
	PcbStands   []map[string]any
	Connectors  []map[string]any
	SnapJoins   []map[string]any
	Cutouts     []Cutout
	PushButtons []map[string]any

	// Derived feature toggles
	PrintSwitchExtenders bool
}

// Cutout is a normalized representation that includes the target face.
type Cutout struct {
	Face string         // front|back|left|right|lid|base
	Item map[string]any // original item fields (already resolved)
	Raw  map[string]any // optional raw for future use
}

// BuildModel converts a resolved DSL document into a Model.
// The input is expected to be fully numeric where applicable (use pkg/resolver before calling).
func BuildModel(ctx context.Context, resolved map[string]any) (*Model, error) {
	m := &Model{}

	// Project (optional)
	if v, ok := getString(resolved, "project"); ok {
		m.ProjectName = v
	}

	// PCB
	var ok bool
	if m.PcbLength, ok = getFloat(resolved, "pcb.length"); !ok {
		return nil, errors.Errorf("missing required pcb.length")
	}
	if m.PcbWidth, ok = getFloat(resolved, "pcb.width"); !ok {
		return nil, errors.Errorf("missing required pcb.width")
	}
	if m.PcbThickness, ok = getFloat(resolved, "pcb.thickness"); !ok {
		return nil, errors.Errorf("missing required pcb.thickness")
	}
	// standoff height maps from pcb.z_clearance when present; else default to 1.0 (YAPP default)
	if v, ok := getFloat(resolved, "pcb.z_clearance"); ok {
		m.StandoffHeight = v
	} else {
		m.StandoffHeight = 1.0
	}
	// Optional standoff details
	if v, ok := getFloat(resolved, "pcb.standoffs.diameter"); ok {
		m.StandoffDiameter = v
	}
	if v, ok := getFloat(resolved, "pcb.standoffs.screw_d"); ok {
		m.StandoffPinDia = v
	}
	if v, ok := getFloat(resolved, "tolerances.holes"); ok {
		m.StandoffHoleSlack = v
	}

	// Enclosure
	if v, ok := getFloat(resolved, "enclosure.wall.thickness"); ok {
		m.WallThickness = v
	}
	if v, ok := getFloat(resolved, "enclosure.base.thickness"); ok {
		m.BasePlaneThickness = v
		// Use base thickness for lid thickness if none provided (pragmatic default)
		m.LidPlaneThickness = v
	}
	if v, ok := getFloat(resolved, "enclosure.lid.thickness"); ok {
		m.LidPlaneThickness = v
	}
	if v, ok := getFloat(resolved, "enclosure.wall.fillet_radius"); ok {
		m.RoundRadius = v
	}
	// Map a single clearance value to all paddings if present
	if v, ok := getFloat(resolved, "enclosure.wall.clearance"); ok {
		m.PaddingFront = v
		m.PaddingBack = v
		m.PaddingLeft = v
		m.PaddingRight = v
	}

	features, _ := getMap(resolved, "features")
	if err := collectFeatureModules(resolved, features, m); err != nil {
		return nil, err
	}

	return m, nil
}

func getFloat(root map[string]any, path string) (float64, bool) {
	v, ok := lookupPath(root, path)
	if !ok {
		return 0, false
	}
	switch t := v.(type) {
	case float64:
		return t, true
	case int:
		return float64(t), true
	default:
		return 0, false
	}
}

func getString(root map[string]any, path string) (string, bool) {
	v, ok := lookupPath(root, path)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func getMap(root map[string]any, path string) (map[string]any, bool) {
	v, ok := lookupPath(root, path)
	if !ok {
		return nil, false
	}
	m, ok := v.(map[string]any)
	return m, ok
}

func getArray(root map[string]any, path string) ([]any, bool) {
	v, ok := lookupPath(root, path)
	if !ok {
		return nil, false
	}
	a, ok := v.([]any)
	return a, ok
}

func lookupPath(root any, path string) (any, bool) {
	cur := root
	for _, part := range splitPath(path) {
		asMap, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		next, ok := asMap[part]
		if !ok {
			return nil, false
		}
		cur = next
	}
	return cur, true
}

func splitPath(p string) []string {
	out := make([]string, 0, 8)
	start := 0
	for i := 0; i < len(p); i++ {
		if p[i] == '.' {
			out = append(out, p[start:i])
			start = i + 1
		}
	}
	out = append(out, p[start:])
	return out
}

func normalizeArrayOfMaps(in []any) []map[string]any {
	out := make([]map[string]any, 0, len(in))
	for _, v := range in {
		if m, ok := v.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}
	return v
}
