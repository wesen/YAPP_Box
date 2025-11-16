package snapjoins

import (
	"reflect"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

func TestBuild_Basic(t *testing.T) {
	input := []map[string]any{
		{
			"pos":   20.0,
			"width": 10.0,
			"side":  "left",
		},
	}

	out, err := Build(input)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	want := []any{20.0, 10.0, scad.Raw("yappLeft")}
	if !reflect.DeepEqual(out[0], want) {
		t.Fatalf("unexpected build result\nwant: %#v\n got: %#v", want, out[0])
	}
}

func TestBuild_WithFlags(t *testing.T) {
	input := []map[string]any{
		{
			"pos":       25.0,
			"width":     8.0,
			"side":      "left",
			"alignment": "center",
			"symmetric": true,
		},
	}

	out, err := Build(input)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	want := []any{
		25.0,
		8.0,
		scad.Raw("yappLeft"),
		scad.Raw("yappCenter"),
		scad.Raw("yappSymmetric"),
	}

	if !reflect.DeepEqual(out[0], want) {
		t.Fatalf("unexpected build result\nwant: %#v\n got: %#v", want, out[0])
	}
}

func TestBuild_DiamondSymmetric(t *testing.T) {
	input := []map[string]any{
		{
			"pos":       30.0,
			"width":     4.0,
			"side":      "right",
			"alignment": "center",
			"diamond":   true,
			"symmetric": true,
		},
	}

	out, err := Build(input)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	want := []any{
		30.0,
		4.0,
		scad.Raw("yappRight"),
		scad.Raw("yappCenter"),
		scad.Raw("yappSymmetric"),
		scad.Raw("yappRectangle"),
	}

	if !reflect.DeepEqual(out[0], want) {
		t.Fatalf("unexpected build result\nwant: %#v\n got: %#v", want, out[0])
	}
}

