package boxmounts

import (
	"reflect"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

func TestBuild_MinimalMount(t *testing.T) {
	left := true
	input := []BoxMountsItem{
		{
			Pos:       12.0,
			ScrewD:    3.0,
			SlotWidth: 0.0,
			Height:    5.0,
			Faces: BoxMountsItemFaces{
				Left: &left,
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
	back := true
	right := true
	offset := 5.0
	fillet := 1.0
	shellPart := "lid"
	alignment := "center"
	origin := "alt"
	noFillet := true
	input := []BoxMountsItem{
		{
			Pos:          20.0,
			Offset:       &offset,
			ScrewD:       3.5,
			SlotWidth:    4.0,
			Height:       8.0,
			FilletRadius: &fillet,
			Faces: BoxMountsItemFaces{
				Back:  &back,
				Right: &right,
			},
			ShellPart: &shellPart,
			Alignment: &alignment,
			Origin:    &origin,
			NoFillet:  &noFillet,
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
