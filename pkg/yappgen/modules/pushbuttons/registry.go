package pushbuttons

import (
	"github.com/wesen/yapp-encl-resolver/pkg/registry"
)

type module struct{}

var _ registry.FeatureModule = (*module)(nil)

// NewModule returns a registry.FeatureModule for push_buttons.
func NewModule() registry.FeatureModule {
	return &module{}
}

func (m *module) Schema() registry.ModuleSchema {
	return pushButtonsSchema{}
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
		Name: "pushButtons",
		Rows: rows,
	}}, nil
}

type pushButtonsSchema struct{}

var _ registry.ModuleSchema = (*pushButtonsSchema)(nil)

func (pushButtonsSchema) Name() string        { return "push_buttons" }
func (pushButtonsSchema) Path() string        { return "features.push_buttons" }
func (pushButtonsSchema) Description() string { return "Push button extenders (generated schema placeholder)" }
func (pushButtonsSchema) Fields() []registry.FieldSpec {
	return nil
}
func (pushButtonsSchema) ValidateStructure(path string, data any) error    { return _svValidateStructureImpl(path, data) }
func (pushButtonsSchema) ValidateConstraints(path string, data any) error  { return nil }

