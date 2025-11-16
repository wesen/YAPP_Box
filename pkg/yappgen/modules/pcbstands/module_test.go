package pcbstands

import (
	"reflect"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

func TestBuild_Defaults(t *testing.T) {
	input := []map[string]any{
		{
			"x": 5.0,
			"y": 7.5,
		},
	}

	out, err := Build(input)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 stand, got %d", len(out))
	}

	if got := len(out[0]); got != 9 {
		t.Fatalf("expected only positional params (9 entries), got %d", got)
	}
}

func TestBuild_FlagEncoding(t *testing.T) {
	input := []map[string]any{
		{
			"x":              10.0,
			"y":              12.0,
			"shell_part":     "lid_only",
			"treatment":      "hole",
			"corner":         "all",
			"coordinate":     "box_inside",
			"no_fillet":      true,
			"pcb_name":       "Sensor",
			"self_threading": true,
		},
	}

	out, err := Build(input)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	want := []any{
		10.0,
		12.0,
		scad.Undef,
		scad.Undef,
		scad.Undef,
		scad.Undef,
		scad.Undef,
		scad.Undef,
		scad.Undef,
		scad.Raw("yappLidOnly"),
		scad.Raw("yappHole"),
		scad.Raw("yappAllCorners"),
		scad.Raw("yappCoordBoxInside"),
		scad.Raw("yappNoFillet"),
		[]any{scad.Raw("yappPCBName"), "Sensor"},
		scad.Raw("yappSelfThreading"),
	}

	if !reflect.DeepEqual(out[0], want) {
		t.Fatalf("unexpected build result\nwant: %#v\n got: %#v", want, out[0])
	}
}
