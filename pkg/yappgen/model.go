package yappgen

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver"
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
	BaseWallHeight     float64
	LidWallHeight      float64
	RidgeHeight        float64
	RoundRadius        float64
	PaddingFront       float64
	PaddingBack        float64
	PaddingLeft        float64
	PaddingRight       float64

	// Features
	PcbStands   []map[string]any
	Connectors  []map[string]any
	BoxMounts   []map[string]any
	SnapJoins   []map[string]any
	Cutouts     []map[string]any
	PushButtons []map[string]any
	LightTubes  []map[string]any

	// Derived feature toggles
	PrintSwitchExtenders bool

	// Provenance metadata for comment emission.
	Provenance *Provenance
}

// BuildModel converts a resolved DSL document into a Model.
// The input is expected to be fully numeric where applicable (use pkg/resolver before calling).
func BuildModel(ctx context.Context, resolved map[string]any, trace resolver.Trace, comments map[string][]string) (*Model, error) {
	m := &Model{
		Provenance: NewProvenance(trace, resolved, comments),
	}

	// Project (optional)
	if v, ok := getString(resolved, "project"); ok {
		m.ProjectName = v
	}

	// PCB
	var ok bool
	if m.PcbLength, ok = getFloat(resolved, "pcb.length"); !ok {
		return nil, errors.Errorf("missing required pcb.length")
	}
	m.Provenance.AddScalar("pcbLength", "pcb.length")
	if m.PcbWidth, ok = getFloat(resolved, "pcb.width"); !ok {
		return nil, errors.Errorf("missing required pcb.width")
	}
	m.Provenance.AddScalar("pcbWidth", "pcb.width")
	if m.PcbThickness, ok = getFloat(resolved, "pcb.thickness"); !ok {
		return nil, errors.Errorf("missing required pcb.thickness")
	}
	m.Provenance.AddScalar("pcbThickness", "pcb.thickness")
	// standoff height maps from pcb.z_clearance when present; else default to 1.0 (YAPP default)
	if v, ok := getFloat(resolved, "pcb.z_clearance"); ok {
		m.StandoffHeight = v
		m.Provenance.AddScalar("standoffHeight", "pcb.z_clearance")
	} else {
		m.StandoffHeight = 1.0
	}
	// Optional standoff details
	if v, ok := getFloat(resolved, "pcb.standoffs.diameter"); ok {
		m.StandoffDiameter = v
		m.Provenance.AddScalar("standoffDiameter", "pcb.standoffs.diameter")
	}
	if v, ok := getFloat(resolved, "pcb.standoffs.screw_d"); ok {
		m.StandoffPinDia = v
		m.Provenance.AddScalar("standoffPinDiameter", "pcb.standoffs.screw_d")
	}
	if v, ok := getFloat(resolved, "tolerances.holes"); ok {
		m.StandoffHoleSlack = v
		m.Provenance.AddScalar("standoffHoleSlack", "tolerances.holes")
	}

	// Enclosure
	if v, ok := getFloat(resolved, "enclosure.wall.thickness"); ok {
		m.WallThickness = v
		m.Provenance.AddScalar("wallThickness", "enclosure.wall.thickness")
	}
	if v, ok := getFloat(resolved, "enclosure.base.thickness"); ok {
		m.BasePlaneThickness = v
		m.Provenance.AddScalar("basePlaneThickness", "enclosure.base.thickness")
		// Use base thickness for lid thickness if none provided (pragmatic default)
		m.LidPlaneThickness = v
		m.Provenance.AddScalar("lidPlaneThickness", "enclosure.base.thickness")
	}
	if v, ok := getFloat(resolved, "enclosure.lid.thickness"); ok {
		m.LidPlaneThickness = v
		m.Provenance.AddScalar("lidPlaneThickness", "enclosure.lid.thickness")
	}
	if v, ok := getFloat(resolved, "enclosure.base.wall_height"); ok {
		m.BaseWallHeight = v
		m.Provenance.AddScalar("baseWallHeight", "enclosure.base.wall_height")
	}
	if v, ok := getFloat(resolved, "enclosure.lid.wall_height"); ok {
		m.LidWallHeight = v
		m.Provenance.AddScalar("lidWallHeight", "enclosure.lid.wall_height")
	}
	if v, ok := getFloat(resolved, "enclosure.ridge.height"); ok {
		m.RidgeHeight = v
		m.Provenance.AddScalar("ridgeHeight", "enclosure.ridge.height")
	}
	if v, ok := getFloat(resolved, "enclosure.wall.fillet_radius"); ok {
		m.RoundRadius = v
		m.Provenance.AddScalar("roundRadius", "enclosure.wall.fillet_radius")
	}
	// Map a single clearance value to all paddings if present
	if v, ok := getFloat(resolved, "enclosure.wall.clearance"); ok {
		m.PaddingFront = v
		m.PaddingBack = v
		m.PaddingLeft = v
		m.PaddingRight = v
		m.Provenance.AddScalar("paddingFront", "enclosure.wall.clearance")
		m.Provenance.AddScalar("paddingBack", "enclosure.wall.clearance")
		m.Provenance.AddScalar("paddingLeft", "enclosure.wall.clearance")
		m.Provenance.AddScalar("paddingRight", "enclosure.wall.clearance")
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
		switch t := cur.(type) {
		case map[string]any:
			next, ok := t[part]
			if !ok {
				return nil, false
			}
			cur = next
		case []any:
			idx, err := strconv.Atoi(part)
			if err != nil || idx < 0 || idx >= len(t) {
				return nil, false
			}
			cur = t[idx]
		default:
			return nil, false
		}
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
