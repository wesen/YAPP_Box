package connectors

import (
	"reflect"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/scad"
)

func TestBuild_Defaults(t *testing.T) {
	input := []ConnectorsItem{
		{
			X:           15.0,
			Y:           10.0,
			StandHeight: 5.0,
			ScrewD:      2.2,
			ScrewHeadD:  4.2,
			InsertD:     3.2,
			OutsideD:    7.5,
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
	corner := "all"
	coordinate := "box_inside"
	noFillet := true
	countersink := true
	pcbName := "Sensor"
	throughLid := true
	selfThreading := true
	noInternalFillet := true

	input := []ConnectorsItem{
		{
			X:                25.0,
			Y:                30.0,
			StandHeight:      6.0,
			ScrewD:           3.0,
			ScrewHeadD:       6.0,
			InsertD:          4.0,
			OutsideD:         8.0,
			Corner:           &corner,
			Coordinate:       &coordinate,
			NoFillet:         &noFillet,
			Countersink:      &countersink,
			PcbName:          &pcbName,
			ThroughLid:       &throughLid,
			SelfThreading:    &selfThreading,
			NoInternalFillet: &noInternalFillet,
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
