package connectors

import (
	_ "embed"

	"github.com/wesen/yapp-encl-resolver/pkg/registry"
)

//go:embed schema.yaml
var schemaYAML []byte

// NewModule returns a FeatureModule for connectors.
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
	return Build(items)
}

// Ensure module implements FeatureModule
var _ registry.FeatureModule = &module{}

type moduleSchema struct{}

func (s *moduleSchema) Name() string {
	return "connectors"
}

func (s *moduleSchema) Path() string {
	return "features.connectors"
}

func (s *moduleSchema) Description() string {
	return "Defines connector mounting points"
}

func (s *moduleSchema) ValidateStructure(path string, data any) error {
	return _svValidateStructureImpl(path, data)
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
