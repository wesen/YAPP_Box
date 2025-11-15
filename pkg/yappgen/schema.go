package yappgen

// ParamSpec describes a positional parameter for a YAPP feature.
type ParamSpec struct {
	Name     string
	Required bool
}

var pcbStandsSchema = []ParamSpec{
	{Name: "x", Required: true},
	{Name: "y", Required: true},
	{Name: "height", Required: false},
	{Name: "pcb_gap", Required: false},
	{Name: "diameter", Required: false},
	{Name: "pin_diameter", Required: false},
	{Name: "hole_slack", Required: false},
	{Name: "fillet_radius", Required: false},
	{Name: "pin_length", Required: false},
}

var connectorsSchema = []ParamSpec{
	{Name: "x", Required: true},
	{Name: "y", Required: true},
	{Name: "stand_height", Required: true},
	{Name: "screw_d", Required: true},
	{Name: "screw_head_d", Required: true},
	{Name: "insert_d", Required: true},
	{Name: "outside_d", Required: true},
	{Name: "insert_depth", Required: false},
	{Name: "pcb_gap", Required: false},
	{Name: "fillet_radius", Required: false},
}

// snapJoins pos params; note: side is a flag appended after positional params.
var snapJoinsSchema = []ParamSpec{
	{Name: "pos", Required: true},
	{Name: "width", Required: true},
}

// cutouts params; note: shape_flag is not a field in the item; it's derived from 'shape'.
var cutoutsSchema = []ParamSpec{
	{Name: "from_back", Required: true},
	{Name: "from_left", Required: true},
	{Name: "width", Required: true},  // may be set to 0 for some shapes
	{Name: "length", Required: true}, // may be set to 0 for some shapes
	{Name: "radius", Required: true}, // may be set to 0 for some shapes
	// p(5) is shape flag (yappRectangle, etc.) — handled separately
	{Name: "depth", Required: false},
	{Name: "angle", Required: false},
}

// ShapeFlag returns the YAPP flag for a cutout shape and indicates which dimension fields are used.
// Unused dimensions should be emitted as 0.
func ShapeFlag(shape string) (flag ScadRaw, usesWidth bool, usesLength bool, usesRadius bool, ok bool) {
	switch shape {
	case "rectangle":
		return ScadRaw("yappRectangle"), true, true, false, true
	case "circle":
		return ScadRaw("yappCircle"), false, false, true, true
	case "rounded_rect":
		return ScadRaw("yappRoundedRect"), true, true, true, true
	case "circle_with_flats":
		return ScadRaw("yappCircleWithFlats"), true, true, true, true // width, length (distance between flats), radius
	case "circle_with_key":
		return ScadRaw("yappCircleWithKey"), true, true, true, true // width=key width, length=key depth, radius
	default:
		return "", false, false, false, false
	}
}

// SnapSideFlag maps DSL side enum to the YAPP side flag.
func SnapSideFlag(side string) (ScadRaw, bool) {
	switch side {
	case "left":
		return ScadRaw("yappLeft"), true
	case "right":
		return ScadRaw("yappRight"), true
	case "front":
		return ScadRaw("yappFront"), true
	case "back":
		return ScadRaw("yappBack"), true
	default:
		return "", false
	}
}


