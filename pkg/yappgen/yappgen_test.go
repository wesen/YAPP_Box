package yappgen

import (
	"context"
	"strings"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/cutouts"
	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/pushbuttons"
	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/modules/snapjoins"
	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
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
	if _, ok := params[2].(scad.Raw); !ok || params[2] != scad.Undef {
		t.Fatalf("expected params[2] to be undef, got %#v", params[2])
	}
	if _, ok := params[3].(scad.Raw); !ok || params[3] != scad.Undef {
		t.Fatalf("expected params[3] to be undef, got %#v", params[3])
	}
	// diameter should be at index 4
	if params[4] != 6.0 {
		t.Fatalf("expected diameter at index 4 to be 6.0, got %#v", params[4])
	}
}

func TestBuildCutoutParams_ShapeSpecificZeros(t *testing.T) {
	// Circle uses radius; width and length should be 0
	fromBottom := float64(6)
	width := float64(0)
	length := float64(0)
	radius := float64(4)
	items := []cutouts.CutoutsItem{
		{
			Face:           "front",
			Shape:          "circle",
			FromFaceLeft:   12.0,
			FromFaceBottom: &fromBottom,
			Width:          &width,
			Length:         &length,
			Radius:         &radius,
		},
	}
	byFace, err := cutouts.Build(items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(byFace["cutoutsFront"]) != 1 {
		t.Fatalf("expected one cutout, got %d", len(byFace["cutoutsFront"]))
	}
	params := byFace["cutoutsFront"][0]
	if len(params) < 8 {
		t.Fatalf("expected at least 8 params incl. shape flag, got %d", len(params))
	}
	// Module Build function already sets width/length to 0 for circles
	if params[2].(float64) != 0 || params[3].(float64) != 0 {
		t.Fatalf("expected width and length to be 0 for circle, got %v, %v", params[2], params[3])
	}
	if params[4].(float64) != 4 {
		t.Fatalf("expected radius 4 at index 4, got %#v", params[4])
	}
	if params[5] != scad.Raw("yappCircle") {
		t.Fatalf("expected yappCircle shape flag at index 5, got %#v", params[5])
	}
}

func TestDistributeCutouts_ByFace(t *testing.T) {
	rectWidth := float64(8)
	rectLength := float64(3)
	rectRadius := float64(0)
	rectBottom := float64(8)
	circWidth := float64(0)
	circLength := float64(0)
	circRadius := float64(2)
	circBottom := float64(6)

	items := []cutouts.CutoutsItem{
		{
			Face:           "front",
			Shape:          "rectangle",
			FromFaceLeft:   5.0,
			FromFaceBottom: &rectBottom,
			Width:          &rectWidth,
			Length:         &rectLength,
			Radius:         &rectRadius,
		},
		{
			Face:           "left",
			Shape:          "circle",
			FromFaceLeft:   4.0,
			FromFaceBottom: &circBottom,
			Width:          &circWidth,
			Length:         &circLength,
			Radius:         &circRadius,
		},
	}
	m, err := cutouts.Build(items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m["cutoutsFront"]) != 1 || len(m["cutoutsLeft"]) != 1 {
		t.Fatalf("expected one cutout in front and left, got front=%d left=%d", len(m["cutoutsFront"]), len(m["cutoutsLeft"]))
	}
	_ = cutouts.NewModule() // ensure module compiles with interface
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
	items := []snapjoins.SnapJoinsItem{
		{Pos: 25, Width: 8, Side: "left"},
	}
	list, err := snapjoins.Build(items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected one row")
	}
	if list[0][2] != scad.Raw("yappLeft") {
		t.Fatalf("expected yappLeft flag at index 2, got %#v", list[0][2])
	}
}

func TestBuildPushButtons_PolygonPreset(t *testing.T) {
	name := "reset"
	shape := "polygon"
	polygon := "arrow"
	angle := 45.0
	coordinate := "box"
	origin := "left"
	filletRadius := 0.6
	lidWall := 2.0
	plateThickness := 2.4
	slack := 0.3
	snapSlack := 0.15
	lidNoFillet := true
	topHeight := 4.1

	items := []pushbuttons.PushButtonsItem{
		{
			Name:         &name,
			X:            12.0,
			Y:            8.0,
			Shape:        &shape,
			Polygon:      &polygon,
			Angle:        &angle,
			Coordinate:   &coordinate,
			Origin:       &origin,
			FilletRadius: &filletRadius,
			Lid: pushbuttons.PushButtonsItemLid{
				Protrusion:     1.5,
				Wall:           &lidWall,
				PlateThickness: &plateThickness,
				Slack:          &slack,
				SnapSlack:      &snapSlack,
				NoFillet:       &lidNoFillet,
			},
			Cap: pushbuttons.PushButtonsItemCap{
				Length: 8.0,
				Width:  6.0,
				Radius: 2.0,
			},
			Switch: pushbuttons.PushButtonsItemSwitch{
				Height:       5.0,
				Travel:       0.6,
				PoleDiameter: 3.0,
				TopHeight:    &topHeight,
			},
		},
	}

	list, err := pushbuttons.Build(items)
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
	if params[10] != scad.Raw("yappPolygon") {
		t.Fatalf("expected polygon flag at index 10, got %#v", params[10])
	}
	if params[17] != scad.Raw("shapeArrow") {
		t.Fatalf("expected shapeArrow token after base params, got %#v", params[17])
	}
	if params[18] != scad.Raw("yappCoordBox") {
		t.Fatalf("expected yappCoordBox token, got %#v", params[18])
	}
	if params[19] != scad.Raw("yappLeftOrigin") {
		t.Fatalf("expected yappLeftOrigin token, got %#v", params[19])
	}
	if params[20] != scad.Raw("yappNoFillet") {
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
