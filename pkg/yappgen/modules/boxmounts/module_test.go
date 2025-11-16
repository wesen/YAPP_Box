package boxmounts

import (
	"reflect"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

func TestBuild_MinimalMount(t *testing.T) {
	input := []map[string]any{
		{
			"pos":        12.0,
			"screw_d":    3.0,
			"slot_width": 0.0,
			"height":     5.0,
			"faces": map[string]any{
				"left": true,
			},
		},
	}

	out, err := Build(input)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	want := []any{
		12.0,
		3.0,
		0.0,
		5.0,
		scad.Undef,
		scad.Raw("yappLeft"),
	}

	if got := out[0]; !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected build result\nwant: %#v\n got: %#v", want, got)
	}
}

func TestBuild_FlaggedMount(t *testing.T) {
	input := []map[string]any{
		{
			"pos":           20.0,
			"offset":        5.0,
			"screw_d":       3.5,
			"slot_width":    4.0,
			"height":        8.0,
			"fillet_radius": 1.0,
			"faces": map[string]any{
				"back":  true,
				"right": true,
			},
			"shell_part": "lid",
			"alignment":  "center",
			"origin":     "alt",
			"no_fillet":  true,
		},
	}

	out, err := Build(input)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	want := []any{
		[]any{20.0, 5.0},
		3.5,
		4.0,
		8.0,
		1.0,
		scad.Raw("yappRight"),
		scad.Raw("yappBack"),
		scad.Raw("yappNoFillet"),
		scad.Raw("yappLid"),
		scad.Raw("yappCenter"),
		scad.Raw("yappAltOrigin"),
	}

	if got := out[0]; !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected build result\nwant: %#v\n got: %#v", want, got)
	}
}
