package pcbstands

import (
	_ "embed"

	"github.com/wesen/yapp-encl-resolver/pkg/registry"
)

//go:embed schema.yaml
var schemaYAML []byte

// NewModule returns a FeatureModule for pcb_stands.
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
		Name: "pcbStands",
		Rows: rows,
	}}, nil
}

// Ensure module implements FeatureModule
var _ registry.FeatureModule = &module{}

type moduleSchema struct{}

func (s *moduleSchema) Name() string {
	return "pcb_stands"
}

func (s *moduleSchema) Path() string {
	return "features.pcb_stands"
}

func (s *moduleSchema) Description() string {
	return "Defines PCB mounting standoffs"
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
