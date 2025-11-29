package errorx

// YAMLIngestContext contains context for YAML parsing errors.
type YAMLIngestContext struct {
	File    string
	Line    int
	Column  int
	Snippet string // Offending YAML snippet
}

func (c *YAMLIngestContext) Stage() StageCode {
	return StageIngestYAMLSyntax
}

// SchemaConstraintContext contains context for schema validation errors.
type SchemaConstraintContext struct {
	Module    string   // Module name (e.g., "connectors")
	FieldPath string   // Field path within module (e.g., "corner")
	Expected  string   // Expected type (e.g., "enum", "number")
	Allowed   []string // Allowed enum values (if applicable)
	Actual    any      // Actual value provided
	Min       *float64 // Minimum value (if applicable)
	Max       *float64 // Maximum value (if applicable)
	Line      int      // Line number in source YAML (0 if unknown)
	Column    int      // Column number in source YAML (0 if unknown)
}

func (c *SchemaConstraintContext) Stage() StageCode {
	return StageSchemaConstraints
}

// SchemaStructureContext contains context for Phase 1 structure validation errors.
type SchemaStructureContext struct {
	Module    string // Module name
	FieldPath string // Field path within module
	Expected  string // Expected type
	Actual    string // Actual type found
	Required  bool   // Whether field is required
	Line      int    // Line number in source YAML (0 if unknown)
	Column    int    // Column number in source YAML (0 if unknown)
}

func (c *SchemaStructureContext) Stage() StageCode {
	return StageSchemaStructure
}

// ExprDependencyContext contains context for missing dependency errors.
type ExprDependencyContext struct {
	Expression  string   // Expression text (e.g., "max(vars.height, 10)")
	MissingRefs []string // Missing variable references (e.g., ["vars.height"])
	Iterations  int      // Number of resolution iterations attempted
	Line        int      // Line number in source YAML (0 if unknown)
	Column      int      // Column number in source YAML (0 if unknown)
}

func (c *ExprDependencyContext) Stage() StageCode {
	return StageExprDependencyMissing
}

// ExprSyntaxContext contains context for expression syntax errors.
type ExprSyntaxContext struct {
	Expression string // Expression text
	Token      string // Offending token (if available)
	Position   int    // Character position in expression
	Line       int    // Line number in source YAML (0 if unknown)
	Column     int    // Column number in source YAML (0 if unknown)
}

func (c *ExprSyntaxContext) Stage() StageCode {
	return StageExprSyntax
}

// ExprRuntimeContext contains context for expression runtime errors.
type ExprRuntimeContext struct {
	Expression string // Expression text
	Function   string // Function name (if applicable)
	Message    string // Error message from evaluator
}

func (c *ExprRuntimeContext) Stage() StageCode {
	return StageExprRuntime
}

// StrictModeContext contains context for strict-mode violations.
type StrictModeContext struct {
	Offenders []string // List of offending keys/variables
	Kind      string   // "unknown-key" or "unused-var"
}

func (c *StrictModeContext) Stage() StageCode {
	if c.Kind == "unknown-key" {
		return StageStrictUnknownKey
	}
	return StageStrictUnusedVar
}

