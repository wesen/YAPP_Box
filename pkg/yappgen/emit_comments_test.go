package yappgen

import (
	"context"
	"strings"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver"
)

func TestEmitSCADIncludesProvenanceComments(t *testing.T) {
	ctx := context.Background()
	doc := map[string]any{
		"pcb": map[string]any{
			"length":      90.0,
			"width":       70.0,
			"thickness":   1.6,
			"z_clearance": 4.0,
		},
		"vars": map[string]any{
			"button_spacing": 18.0,
			"button_center":  "pcb.length / 2",
		},
		"features": map[string]any{
			"push_buttons": []any{
				map[string]any{
					"name":  "btn-left",
					"x":     "vars.button_center - vars.button_spacing",
					"y":     45.0,
					"shape": "circle",
					"cap": map[string]any{
						"length": 24.0,
						"width":  24.0,
						"radius": 12.0,
					},
					"lid": map[string]any{
						"protrusion":      0.0,
						"wall":            2.0,
						"plate_thickness": 2.4,
					},
					"switch": map[string]any{
						"height":        14.0,
						"travel":        0.8,
						"pole_diameter": 3.2,
					},
				},
			},
		},
	}

	result, err := resolver.ResolveResult(ctx, doc, resolver.Options{})
	if err != nil {
		t.Fatalf("resolver failed: %v", err)
	}

	model, err := BuildModel(ctx, result.Document, result.Trace)
	if err != nil {
		t.Fatalf("BuildModel failed: %v", err)
	}

	scadBytes, err := EmitSCAD(ctx, model)
	if err != nil {
		t.Fatalf("EmitSCAD failed: %v", err)
	}
	scad := string(scadBytes)

	if !strings.Contains(scad, `// pcbLength (source: pcb.length)`) {
		t.Fatalf("expected pcbLength comment, got:\n%s", scad)
	}
	if !strings.Contains(scad, `expr="vars.button_center - vars.button_spacing"`) {
		t.Fatalf("expected expression comment for push button x, got:\n%s", scad)
	}
	if !strings.Contains(scad, `deps: vars.button_center=`) {
		t.Fatalf("expected dependency list in push button comment, got:\n%s", scad)
	}
	if !strings.Contains(scad, `// pushButtons[0] ← features.push_buttons[0]`) {
		t.Fatalf("expected pushButtons row header comment, got:\n%s", scad)
	}
}
