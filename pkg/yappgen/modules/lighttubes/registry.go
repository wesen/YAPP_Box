package lighttubes

import (
	_ "embed"

	"github.com/wesen/yapp-encl-resolver/pkg/registry"
)

//go:embed schema.yaml
var schemaYAML []byte

// NewModule returns a FeatureModule for light_tubes.
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

func (m *module) Build(items []map[string]any) ([]registry.ArrayDecl, error) {
	typed, err := Decode(items)
	if err != nil {
		return nil, err
	}

	rows, err := Build(typed)
	if err != nil {
		return nil, err
	}

	return []registry.ArrayDecl{{
		Name: "lightTubes",
		Rows: rows,
	}}, nil
}

var _ registry.FeatureModule = &module{}

type moduleSchema struct{}

func (s *moduleSchema) Name() string {
	return "light_tubes"
}

func (s *moduleSchema) Path() string {
	return "features.light_tubes"
}

func (s *moduleSchema) Description() string {
	return "Defines LED light pipes that extend from the PCB through the lid."
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

var _ registry.ModuleSchema = &moduleSchema{}
