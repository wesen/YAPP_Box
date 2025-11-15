package yappgen

import (
	"context"
	"strings"
	"testing"
)

func TestBuildParams_InsertsUndefForOptional(t *testing.T) {
	item := map[string]any{
		"x":        10,
		"y":        20,
		"diameter": 6.0, // skipping height and pcb_gap
	}
	params, err := buildParams(item, pcbStandsSchema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(params) != 9 {
		t.Fatalf("expected 9 params, got %d", len(params))
	}
	// height and pcb_gap should be undef
	if _, ok := params[2].(ScadRaw); !ok || params[2] != Undef {
		t.Fatalf("expected params[2] to be undef, got %#v", params[2])
	}
	if _, ok := params[3].(ScadRaw); !ok || params[3] != Undef {
		t.Fatalf("expected params[3] to be undef, got %#v", params[3])
	}
	// diameter should be at index 4
	if params[4] != 6.0 {
		t.Fatalf("expected diameter at index 4 to be 6.0, got %#v", params[4])
	}
}

func TestBuildCutoutParams_ShapeSpecificZeros(t *testing.T) {
	// Circle uses radius; width and length should be 0
	item := map[string]any{
		"shape":     "circle",
		"from_back": 30,
		"from_left": 12,
		"radius":    4,
	}
	params, err := buildCutoutParams(item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(params) != 8 {
		t.Fatalf("expected 8 params incl. shape flag, got %d", len(params))
	}
	if params[2] != 0 || params[3] != 0 {
		t.Fatalf("expected width and length to be 0 for circle, got %v, %v", params[2], params[3])
	}
	if params[4] != 4 {
		t.Fatalf("expected radius 4 at index 4, got %#v", params[4])
	}
	if params[5] != ScadRaw("yappCircle") {
		t.Fatalf("expected yappCircle shape flag at index 5, got %#v", params[5])
	}
}

func TestDistributeCutouts_ByFace(t *testing.T) {
	cuts := []Cutout{
		{Face: "front", Item: map[string]any{"shape": "rectangle", "from_back": 10, "from_left": 5, "width": 8, "length": 3}},
		{Face: "left", Item: map[string]any{"shape": "circle", "from_back": 12, "from_left": 4, "radius": 2}},
	}
	m, err := distributeCutouts(cuts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m["cutoutsFront"]) != 1 || len(m["cutoutsLeft"]) != 1 {
		t.Fatalf("expected one cutout in front and left, got front=%d left=%d", len(m["cutoutsFront"]), len(m["cutoutsLeft"]))
	}
}

func TestEmitSCAD_ContainsUndefAndArrays(t *testing.T) {
	ctx := context.Background()
	model := &Model{
		ProjectName:        "test",
		PcbLength:          60,
		PcbWidth:           40,
		PcbThickness:       1.6,
		StandoffHeight:     2,
		WallThickness:      2,
		BasePlaneThickness: 1.5,
		LidPlaneThickness:  1.5,
		RoundRadius:        3,
		PaddingFront:       2, PaddingBack: 2, PaddingLeft: 2, PaddingRight: 2,
		PcbStands: []map[string]any{
			{"x": 10, "y": 10, "diameter": 6.0},
		},
	}
	out, err := EmitSCAD(ctx, model)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "pcbStands") {
		t.Fatalf("expected pcbStands array in output")
	}
	if !strings.Contains(s, "undef") {
		t.Fatalf("expected 'undef' placeholders in output")
	}
	if !strings.Contains(s, "printSwitchExtenders = false") {
		t.Fatalf("expected printSwitchExtenders default to be false")
	}
}

func TestBuildSnapJoins_SideFlag(t *testing.T) {
	items := []map[string]any{
		{"pos": 25, "width": 8, "side": "left"},
	}
	list, err := buildSnapJoins(items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected one row")
	}
	if list[0][2] != ScadRaw("yappLeft") {
		t.Fatalf("expected yappLeft flag at index 2, got %#v", list[0][2])
	}
}

func TestBuildPushButtons_PolygonPreset(t *testing.T) {
	items := []map[string]any{
		{
			"name":          "reset",
			"x":             12.0,
			"y":             8.0,
			"shape":         "polygon",
			"polygon":       "arrow",
			"angle":         45.0,
			"coordinate":    "box",
			"origin":        "left",
			"fillet_radius": 0.6,
			"lid": map[string]any{
				"protrusion":      1.5,
				"wall":            2.0,
				"plate_thickness": 2.4,
				"slack":           0.3,
				"snap_slack":      0.15,
				"no_fillet":       true,
			},
			"cap": map[string]any{
				"length": 8.0,
				"width":  6.0,
				"radius": 2.0,
			},
			"switch": map[string]any{
				"height":        5.0,
				"travel":        0.6,
				"pole_diameter": 3.0,
				"top_height":    4.1,
			},
		},
	}

	list, err := buildPushButtons(items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected one push button entry, got %d", len(list))
	}
	params := list[0]
	if len(params) < 21 {
		t.Fatalf("expected extended parameter list including flags, got %d", len(params))
	}
	if params[10] != ScadRaw("yappPolygon") {
		t.Fatalf("expected polygon flag at index 10, got %#v", params[10])
	}
	if params[17] != ScadRaw("shapeArrow") {
		t.Fatalf("expected shapeArrow token after base params, got %#v", params[17])
	}
	if params[18] != ScadRaw("yappCoordBox") {
		t.Fatalf("expected yappCoordBox token, got %#v", params[18])
	}
	if params[19] != ScadRaw("yappLeftOrigin") {
		t.Fatalf("expected yappLeftOrigin token, got %#v", params[19])
	}
	if params[20] != ScadRaw("yappNoFillet") {
		t.Fatalf("expected yappNoFillet token, got %#v", params[20])
	}
}

func TestEmitSCAD_SetsPushButtonsFlag(t *testing.T) {
	ctx := context.Background()
	model := &Model{
		PcbLength:      40,
		PcbWidth:       20,
		PcbThickness:   1.6,
		StandoffHeight: 2,
		PushButtons: []map[string]any{
			{
				"x": 10.0,
				"y": 12.0,
				"cap": map[string]any{
					"length": 6.0,
					"width":  6.0,
					"radius": 1.0,
				},
				"lid": map[string]any{
					"protrusion": 1.0,
				},
				"switch": map[string]any{
					"height":        4.0,
					"travel":        0.5,
					"pole_diameter": 2.5,
				},
			},
		},
		PrintSwitchExtenders: true,
	}

	out, err := EmitSCAD(ctx, model)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "pushButtons") {
		t.Fatalf("expected pushButtons array in output")
	}
	if !strings.Contains(s, "printSwitchExtenders = true") {
		t.Fatalf("expected printSwitchExtenders to be true when push buttons exist")
	}
}
