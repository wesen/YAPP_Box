package errorx

import (
	"fmt"
)

// StageCode represents the pipeline phase where an error occurred.
type StageCode string

const (
	// StageIngestYAMLSyntax indicates YAML parsing failed (syntax error).
	StageIngestYAMLSyntax StageCode = "ingest.yaml.syntax"
	// StageIngestFS indicates file system error (file not found, permission denied).
	StageIngestFS StageCode = "ingest.fs"
	// StageSchemaStructure indicates Phase 1 validation failed (missing required field, wrong type).
	StageSchemaStructure StageCode = "schema.structure"
	// StageSchemaConstraints indicates Phase 2 validation failed (min/max/enum constraint violation).
	StageSchemaConstraints StageCode = "schema.constraints"
	// StageExprSyntax indicates expression syntax error (invalid token, malformed expression).
	StageExprSyntax StageCode = "expr.syntax"
	// StageExprRuntime indicates expression runtime error (type mismatch, invalid function call).
	StageExprRuntime StageCode = "expr.runtime"
	// StageExprDependencyMissing indicates expression references undefined variable.
	StageExprDependencyMissing StageCode = "expr.dependency.missing"
	// StageStrictUnknownKey indicates strict mode: unknown top-level key.
	StageStrictUnknownKey StageCode = "strict.unknown-key"
	// StageStrictUnusedVar indicates strict mode: variable declared but never referenced.
	StageStrictUnusedVar StageCode = "strict.unused-var"
)

// SymptomCode represents the specific type of failure.
type SymptomCode string

const (
	// SymptomSyntax indicates syntax error (YAML or expression).
	SymptomSyntax SymptomCode = "syntax"
	// SymptomMissingRequired indicates required field is missing.
	SymptomMissingRequired SymptomCode = "missing_required"
	// SymptomTypeMismatch indicates field has wrong type (expected number, got string).
	SymptomTypeMismatch SymptomCode = "type_mismatch"
	// SymptomEnumMismatch indicates value not in allowed enum list.
	SymptomEnumMismatch SymptomCode = "enum_mismatch"
	// SymptomConstraintViolation indicates min/max constraint violated.
	SymptomConstraintViolation SymptomCode = "constraint_violation"
	// SymptomDependencyMissing indicates referenced variable/field doesn't exist.
	SymptomDependencyMissing SymptomCode = "dependency_missing"
	// SymptomUnknownKey indicates key not recognized (strict mode).
	SymptomUnknownKey SymptomCode = "unknown_key"
	// SymptomUnusedVar indicates variable declared but unused (strict mode).
	SymptomUnusedVar SymptomCode = "unused_var"
)

// Severity indicates error severity level.
type Severity string

const (
	// SeverityError indicates hard failure (must fix to proceed).
	SeverityError Severity = "error"
	// SeverityWarning indicates lint violation (can proceed but discouraged).
	SeverityWarning Severity = "warning"
	// SeverityInfo indicates informational message.
	SeverityInfo Severity = "info"
)

// TaxonomyContext is the interface for stage-specific context data.
type TaxonomyContext interface {
	// Stage returns the stage code this context belongs to.
	Stage() StageCode
}

// Taxonomy represents a structured error with metadata.
type Taxonomy struct {
	Stage    StageCode
	Symptom  SymptomCode
	Path     string      // Canonical DSL path (e.g., "features.connectors[0].corner")
	Severity Severity
	Context  TaxonomyContext
}

// Error implements error interface.
func (t *Taxonomy) Error() string {
	if t.Path != "" {
		return fmt.Sprintf("[%s:%s] %s: %v", t.Stage, t.Symptom, t.Path, t.Context)
	}
	return fmt.Sprintf("[%s:%s] %v", t.Stage, t.Symptom, t.Context)
}

// AsTaxonomy unwraps nested errors to find a Taxonomy.
// Similar to errors.As but specifically for Taxonomy.
func AsTaxonomy(err error) (*Taxonomy, bool) {
	if err == nil {
		return nil, false
	}
	if t, ok := err.(*Taxonomy); ok {
		return t, true
	}
	if t, ok := err.(interface{ Unwrap() error }); ok {
		return AsTaxonomy(t.Unwrap())
	}
	return nil, false
}

