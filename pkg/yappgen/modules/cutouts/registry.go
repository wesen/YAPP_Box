package cutouts

import (
	_ "embed"

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
	// TODO: implement schema-driven validation
	return nil
}

func (s *moduleSchema) ValidateConstraints(path string, data any) error {
	// TODO: implement constraint validation
	return nil
}

func (s *moduleSchema) Fields() []registry.FieldSpec {
	// TODO: return field specs from schema
	return nil
}

// Ensure moduleSchema implements ModuleSchema
var _ registry.ModuleSchema = &moduleSchema{}
