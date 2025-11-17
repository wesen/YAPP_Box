package boxmounts

import (
	_ "embed"

	"github.com/wesen/yapp-encl-resolver/pkg/registry"
)

//go:embed schema.yaml
var schemaYAML []byte

// NewModule returns a FeatureModule for box_mounts.
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
		Name: "boxMounts",
		Rows: rows,
	}}, nil
}

// Ensure module implements FeatureModule
var _ registry.FeatureModule = &module{}

type moduleSchema struct{}

func (s *moduleSchema) Name() string {
	return "box_mounts"
}

func (s *moduleSchema) Path() string {
	return "features.box_mounts"
}

func (s *moduleSchema) Description() string {
	return "Defines external mounting tabs"
}

func (s *moduleSchema) ValidateStructure(path string, data any) error {
	return _svValidateStructureImpl(path, data)
}

func (s *moduleSchema) ValidateConstraints(path string, data any) error {
	return nil
}

func (s *moduleSchema) Fields() []registry.FieldSpec {
	return nil
}

// Ensure moduleSchema implements ModuleSchema
var _ registry.ModuleSchema = &moduleSchema{}
