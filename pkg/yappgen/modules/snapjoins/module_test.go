package snapjoins

import (
	"reflect"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

func TestBuild_Basic(t *testing.T) {
	input := []SnapJoinsItem{
		{
			Pos:   20.0,
			Width: 10.0,
			Side:  "left",
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
	alignment := "center"
	symmetric := true
	input := []SnapJoinsItem{
		{
			Pos:       25.0,
			Width:     8.0,
			Side:      "left",
			Alignment: &alignment,
			Symmetric: &symmetric,
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
	alignment := "center"
	diamond := true
	symmetric := true
	input := []SnapJoinsItem{
		{
			Pos:       30.0,
			Width:     4.0,
			Side:      "right",
			Alignment: &alignment,
			Diamond:   &diamond,
			Symmetric: &symmetric,
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
