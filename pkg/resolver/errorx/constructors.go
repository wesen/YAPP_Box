package errorx

// NewYAMLIngestTaxonomy creates a taxonomy entry for YAML syntax errors.
func NewYAMLIngestTaxonomy(file string, line, column int, snippet string) *Taxonomy {
	return &Taxonomy{
		Stage:    StageIngestYAMLSyntax,
		Symptom:  SymptomSyntax,
		Path:     "",
		Severity: SeverityError,
		Context: &YAMLIngestContext{
			File:    file,
			Line:    line,
			Column:  column,
			Snippet: snippet,
		},
	}
}

// NewSchemaConstraintTaxonomy creates a taxonomy entry for schema constraint violations.
func NewSchemaConstraintTaxonomy(path, module, fieldPath string, expected string, actual any, allowed []string, min, max *float64, line, column int) *Taxonomy {
	return &Taxonomy{
		Stage:    StageSchemaConstraints,
		Symptom:  SymptomEnumMismatch, // Default to enum mismatch if allowed values provided
		Path:     path,
		Severity: SeverityError,
		Context: &SchemaConstraintContext{
			Module:    module,
			FieldPath: fieldPath,
			Expected:  expected,
			Allowed:   allowed,
			Actual:    actual,
			Min:       min,
			Max:       max,
			Line:      line,
			Column:    column,
		},
	}
}

// NewSchemaStructureTaxonomy creates a taxonomy entry for Phase 1 structure validation errors.
func NewSchemaStructureTaxonomy(path, module, fieldPath, expected, actual string, required bool, line, column int) *Taxonomy {
	symptom := SymptomTypeMismatch
	if required {
		symptom = SymptomMissingRequired
	}
	return &Taxonomy{
		Stage:    StageSchemaStructure,
		Symptom:  symptom,
		Path:     path,
		Severity: SeverityError,
		Context: &SchemaStructureContext{
			Module:    module,
			FieldPath: fieldPath,
			Expected:  expected,
			Actual:    actual,
			Required:  required,
			Line:      line,
			Column:    column,
		},
	}
}

// NewExprDependencyTaxonomy creates a taxonomy entry for missing dependency errors.
func NewExprDependencyTaxonomy(path, expression string, missingRefs []string, iterations int, line, column int) *Taxonomy {
	return &Taxonomy{
		Stage:    StageExprDependencyMissing,
		Symptom:  SymptomDependencyMissing,
		Path:     path,
		Severity: SeverityError,
		Context: &ExprDependencyContext{
			Expression:  expression,
			MissingRefs: missingRefs,
			Iterations:  iterations,
			Line:        line,
			Column:      column,
		},
	}
}

// NewExprSyntaxTaxonomy creates a taxonomy entry for expression syntax errors.
func NewExprSyntaxTaxonomy(path, expression, token string, position int, line, column int) *Taxonomy {
	return &Taxonomy{
		Stage:    StageExprSyntax,
		Symptom:  SymptomSyntax,
		Path:     path,
		Severity: SeverityError,
		Context: &ExprSyntaxContext{
			Expression: expression,
			Token:      token,
			Position:   position,
			Line:       line,
			Column:     column,
		},
	}
}

// NewStrictModeTaxonomy creates a taxonomy entry for strict-mode violations.
func NewStrictModeTaxonomy(kind string, offenders []string) *Taxonomy {
	var stage StageCode
	var symptom SymptomCode
	if kind == "unknown-key" {
		stage = StageStrictUnknownKey
		symptom = SymptomUnknownKey
	} else {
		stage = StageStrictUnusedVar
		symptom = SymptomUnusedVar
	}
	return &Taxonomy{
		Stage:    stage,
		Symptom:  symptom,
		Path:     "",
		Severity: SeverityWarning, // Strict mode violations are warnings
		Context: &StrictModeContext{
			Offenders: offenders,
			Kind:      kind,
		},
	}
}

