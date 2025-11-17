package cutouts

import (
	_ "embed"

	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/wesen/yapp-encl-resolver/pkg/registry"
)

//go:embed schema.yaml
var schemaYAML []byte

// NewModule returns a FeatureModule for cutouts.
func NewModule() registry.FeatureModule {
	return &module{
		schema: &moduleSchema{},
	}
}

type module struct {
	schema registry.ModuleSchema
}

func (m *module) Schema() registry.ModuleSchema {
	return m.schema
}

func (m *module) Build(items []map[string]any) ([][]any, error) {
	// Note: cutouts return a map by face, not a simple array
	// This is a special case handled in features.go
	byFace, err := Build(items)
	if err != nil {
		return nil, err
	}
	// For now, return empty to satisfy interface
	// The actual cutout handling is in features.go cutoutFeatureModule
	_ = byFace
	return nil, nil
}

// Ensure module implements FeatureModule
var _ registry.FeatureModule = &module{}

type moduleSchema struct{}

func (s *moduleSchema) Name() string {
	return "cutouts"
}

func (s *moduleSchema) Path() string {
	return "features.cutouts"
}

func (s *moduleSchema) Description() string {
	return "Defines cutouts in enclosure faces"
}

func (s *moduleSchema) ValidateStructure(path string, data any) error {
	return _svValidateStructureImpl(path, data)
}

func (s *moduleSchema) ValidateConstraints(path string, data any) error {
	// Phase-2 validation: verify enums and basic cross-field constraints
	// Expect data to be a slice of cutout entries (maps)
	rv := reflect.ValueOf(data)
	if !rv.IsValid() || rv.Kind() != reflect.Slice {
		// Nothing to validate (feature missing or wrong type handled elsewhere)
		return nil
	}

	for i := 0; i < rv.Len(); i++ {
		raw := rv.Index(i).Interface()
		item, ok := raw.(map[string]any)
		if !ok {
			// Skip non-map entries; structure validation should catch this
			continue
		}

		// Helper: build field path like features.cutouts[0].field
		fieldPath := func(field string) string {
			return path + "[" + toIndex(i) + "]." + field
		}

		// Validate face enum (if provided)
		if v, ok := getString(item, "face"); ok && v != "" {
			n := normalizeEnum(v)
			if !inSet(n, allowedFaces) {
				return errors.Errorf("%s: invalid value %q (allowed: front, back, left, right, top, lid, bottom, base)", fieldPath("face"), v)
			}
		}

		// Validate shape enum (if provided)
		shapeVal, hasShape := getString(item, "shape")
		if hasShape && shapeVal != "" {
			shape := normalizeEnum(shapeVal)
			if !inSet(shape, allowedShapes) {
				return errors.Errorf("%s: invalid value %q (allowed: rectangle, circle, rounded_rect, circle_with_flats, circle_with_key, polygon)", fieldPath("shape"), shapeVal)
			}
			// Polygon preset required and must be valid when shape=polygon
			if shape == "polygon" {
				presetVal, hasPreset := getString(item, "polygon")
				if !hasPreset || strings.TrimSpace(presetVal) == "" {
					return errors.Errorf("%s: polygon preset is required when shape=polygon", fieldPath("polygon"))
				}
				preset := normalizeEnum(presetVal)
				if !inSet(preset, allowedPolygonPresets) {
					return errors.Errorf("%s: invalid value %q (allowed: hexagon, arrow, 6pt_star, iso_triangle, iso_triangle2, triangle, triangle2)", fieldPath("polygon"), presetVal)
				}
			}
		}

		// Validate mask.preset enum (if mask provided)
		if maskRaw, ok := item["mask"]; ok && maskRaw != nil {
			if maskMap, ok := maskRaw.(map[string]any); ok {
				if v, ok := getString(maskMap, "preset"); ok && v != "" {
					n := normalizeEnum(v)
					if !inSet(n, allowedMaskPresets) {
						return errors.Errorf("%s: invalid value %q (allowed: honeycomb, hex_circles, circles, squares, bars, offset_bars)", fieldPath("mask.preset"), v)
					}
				} else {
					// If mask is specified, require preset
					return errors.Errorf("%s: required field", fieldPath("mask.preset"))
				}
			}
		}

		// Validate coordinate enum (if provided)
		if v, ok := getString(item, "coordinate"); ok && v != "" {
			n := normalizeEnum(v)
			if !inSet(n, allowedCoordinate) {
				return errors.Errorf("%s: invalid value %q (allowed: pcb, box, box_inside)", fieldPath("coordinate"), v)
			}
		}

		// Validate origin enum (if provided)
		if v, ok := getString(item, "origin"); ok && v != "" {
			n := normalizeEnum(v)
			if !inSet(n, allowedOrigin) {
				return errors.Errorf("%s: invalid value %q (allowed: global, center, alt)", fieldPath("origin"), v)
			}
		}
	}

	return nil
}

func (s *moduleSchema) Fields() []registry.FieldSpec {
	// TODO: return field specs from schema
	return nil
}

// Ensure moduleSchema implements ModuleSchema
var _ registry.ModuleSchema = &moduleSchema{}

// --- helpers ---

// normalizeEnum makes enum comparison forgiving by lowercasing and replacing dashes with underscores.
func normalizeEnum(v string) string {
	out := strings.ToLower(strings.TrimSpace(v))
	out = strings.ReplaceAll(out, "-", "_")
	return out
}

func getString(m map[string]any, key string) (string, bool) {
	raw, ok := m[key]
	if !ok || raw == nil {
		return "", false
	}
	if s, ok := raw.(string); ok {
		return s, true
	}
	return "", false
}

func inSet(v string, set map[string]struct{}) bool {
	_, ok := set[v]
	return ok
}

// toIndex converts an int index to string without pulling in fmt just for Sprintf
func toIndex(i int) string {
	// Simple itoa without import cycles
	const digits = "0123456789"
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	n := i
	for n > 0 {
		pos--
		buf[pos] = digits[n%10]
		n /= 10
	}
	return string(buf[pos:])
}

var allowedShapes = map[string]struct{}{
	"rectangle":         {},
	"circle":            {},
	"rounded_rect":      {},
	"circle_with_flats": {},
	"circle_with_key":   {},
	"polygon":           {},
}

var allowedFaces = map[string]struct{}{
	"front":  {},
	"back":   {},
	"left":   {},
	"right":  {},
	"top":    {},
	"lid":    {},
	"bottom": {},
	"base":   {},
}

var allowedPolygonPresets = map[string]struct{}{
	"hexagon":       {},
	"arrow":         {},
	"6pt_star":      {},
	"iso_triangle":  {},
	"iso_triangle2": {},
	"triangle":      {},
	"triangle2":     {},
}

var allowedMaskPresets = map[string]struct{}{
	"honeycomb":    {},
	"hex_circles":  {},
	"circles":      {},
	"squares":      {},
	"bars":         {},
	"offset_bars":  {},
}

var allowedCoordinate = map[string]struct{}{
	"pcb":        {},
	"box":        {},
	"box_inside": {},
}

var allowedOrigin = map[string]struct{}{
	"global": {},
	"center": {},
	"alt":    {},
}
