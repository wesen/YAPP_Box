package resolver

import (
	"context"
	"testing"
)

func TestResolveResultCapturesExpressionsAndDependencies(t *testing.T) {
	ctx := context.Background()
	doc := map[string]any{
		"pcb": map[string]any{
			"length":    90.0,
			"width":     70.0,
			"thickness": 1.6,
		},
		"vars": map[string]any{
			"button_spacing": 14.0,
			"button_center":  "pcb.length / 2",
		},
		"features": map[string]any{
			"push_buttons": []any{
				map[string]any{
					"x": "vars.button_center - vars.button_spacing",
					"y": 42.0,
					"cap": map[string]any{
						"length": 20.0,
						"width":  20.0,
						"radius": 10.0,
					},
					"lid": map[string]any{
						"protrusion":      0.0,
						"plate_thickness": 2.0,
					},
					"switch": map[string]any{
						"height":        10.0,
						"travel":        1.0,
						"pole_diameter": 3.0,
					},
				},
			},
		},
	}

	result, err := ResolveResult(ctx, doc, Options{}, nil)
	if err != nil {
		t.Fatalf("ResolveResult returned error: %v", err)
	}

	// Literal trace
	if entry, ok := result.Trace.Get("pcb.length"); !ok || entry.Kind != TraceKindLiteral || entry.Value != 90.0 {
		t.Fatalf("expected literal trace for pcb.length, got %#v", entry)
	}

	// Expression trace for vars.button_center
	centerEntry, ok := result.Trace.Get("vars.button_center")
	if !ok {
		t.Fatalf("missing trace entry for vars.button_center")
	}
	if centerEntry.Kind != TraceKindExpression {
		t.Fatalf("expected expression trace for vars.button_center, got %s", centerEntry.Kind)
	}
	if centerEntry.Expression == "" || len(centerEntry.Dependencies) != 1 || centerEntry.Dependencies[0] != "pcb.length" {
		t.Fatalf("unexpected dependency set for vars.button_center: %#v", centerEntry)
	}

	// Expression trace for push button x coordinate
	buttonEntry, ok := result.Trace.Get("features.push_buttons.0.x")
	if !ok {
		t.Fatalf("missing trace entry for push button x coordinate")
	}
	if buttonEntry.Kind != TraceKindExpression {
		t.Fatalf("expected expression trace for button x, got %s", buttonEntry.Kind)
	}
	if len(buttonEntry.Dependencies) != 2 {
		t.Fatalf("expected two dependencies for button x, got %#v", buttonEntry.Dependencies)
	}
	expectedDeps := map[string]struct{}{
		"vars.button_center":  {},
		"vars.button_spacing": {},
	}
	for _, dep := range buttonEntry.Dependencies {
		if _, ok := expectedDeps[dep]; !ok {
			t.Fatalf("unexpected dependency %s for button x", dep)
		}
	}
}
