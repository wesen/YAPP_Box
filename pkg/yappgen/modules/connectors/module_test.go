package connectors

import (
	"reflect"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

func TestBuild_Defaults(t *testing.T) {
	input := []map[string]any{
		{
			"x":            15.0,
			"y":            10.0,
			"stand_height": 5.0,
			"screw_d":      2.2,
			"screw_head_d": 4.2,
			"insert_d":     3.2,
			"outside_d":    7.5,
		},
	}

	out, err := Build(input)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 connector, got %d", len(out))
	}

	if got := len(out[0]); got != 10 {
		t.Fatalf("expected only positional params (10 entries), got %d", got)
	}
}

func TestBuild_FlagEncoding(t *testing.T) {
	input := []map[string]any{
		{
			"x":                  25.0,
			"y":                  30.0,
			"stand_height":       6.0,
			"screw_d":            3.0,
			"screw_head_d":       6.0,
			"insert_d":           4.0,
			"outside_d":          8.0,
			"corner":             "all",
			"coordinate":         "box_inside",
			"no_fillet":          true,
			"countersink":        true,
			"pcb_name":           "Sensor",
			"through_lid":        true,
			"self_threading":     true,
			"no_internal_fillet": true,
		},
	}

	out, err := Build(input)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	want := []any{
		25.0,
		30.0,
		6.0,
		3.0,
		6.0,
		4.0,
		8.0,
		scad.Undef,
		scad.Undef,
		scad.Undef,
		scad.Raw("yappAllCorners"),
		scad.Raw("yappCoordBoxInside"),
		scad.Raw("yappNoFillet"),
		scad.Raw("yappCountersink"),
		[]any{scad.Raw("yappPCBName"), "Sensor"},
		scad.Raw("yappThroughLid"),
		scad.Raw("yappSelfThreading"),
		scad.Raw("yappNoInternalFillet"),
	}

	if !reflect.DeepEqual(out[0], want) {
		t.Fatalf("unexpected build result\nwant: %#v\n got: %#v", want, out[0])
	}
}
