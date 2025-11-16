package registry

// FieldType enumerates supported schema field categories.
type FieldType string

const (
	// NumberField represents scalar numeric values.
	NumberField FieldType = "number"
	// StringField represents scalar string values.
	StringField FieldType = "string"
	// ObjectField represents nested objects (maps).
	ObjectField FieldType = "object"
	// ArrayField represents repeated entries.
	ArrayField FieldType = "array"
)

// FieldSpec captures metadata about a single schema field. It is intended to
// be consumed by documentation generators, validators, and tooling.
type FieldSpec struct {
	Name        string      `json:"name" yaml:"name"`
	Type        FieldType   `json:"type" yaml:"type"`
	Required    bool        `json:"required" yaml:"required"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Example     any         `json:"example,omitempty" yaml:"example,omitempty"`
	Default     any         `json:"default,omitempty" yaml:"default,omitempty"`
	Min         *float64    `json:"min,omitempty" yaml:"min,omitempty"`
	Max         *float64    `json:"max,omitempty" yaml:"max,omitempty"`
	Enum        []string    `json:"enum,omitempty" yaml:"enum,omitempty"`
	Children    []FieldSpec `json:"children,omitempty" yaml:"children,omitempty"`
}

// ModuleSchema defines the behavior all module schemas must provide. Concrete
// implementations are typically generated from YAML schema definitions.
type ModuleSchema interface {
	// Name returns the short module identifier (e.g. "push_buttons").
	Name() string
	// Path returns the fully-qualified registry path (e.g. "features.push_buttons").
	Path() string
	// Description returns a human-readable summary of the module.
	Description() string
	// Fields returns a flattened view of the schema fields for documentation.
	Fields() []FieldSpec
	// ValidateStructure performs phase-1 validation (shape/type checking).
	ValidateStructure(path string, data any) error
	// ValidateConstraints performs phase-2 validation (post-resolution checks).
	ValidateConstraints(path string, data any) error
}

// FeatureModule ties together a schema with the logic that converts validated
// DSL entries into OpenSCAD parameter arrays.
type FeatureModule interface {
	// Schema returns the module's schema metadata.
	Schema() ModuleSchema
	// Build consumes validated items and produces OpenSCAD parameter arrays.
	Build(items []map[string]any) ([][]any, error)
}
