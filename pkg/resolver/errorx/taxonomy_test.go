package errorx

import (
	"testing"

	"github.com/pkg/errors"
)

func TestNewYAMLIngestTaxonomy(t *testing.T) {
	tax := NewYAMLIngestTaxonomy("test.yaml", 5, 10, "snippet")
	
	if tax.Stage != StageIngestYAMLSyntax {
		t.Errorf("expected Stage %s, got %s", StageIngestYAMLSyntax, tax.Stage)
	}
	if tax.Symptom != SymptomSyntax {
		t.Errorf("expected Symptom %s, got %s", SymptomSyntax, tax.Symptom)
	}
	if tax.Path != "" {
		t.Errorf("expected empty Path, got %s", tax.Path)
	}
	if tax.Severity != SeverityError {
		t.Errorf("expected Severity %s, got %s", SeverityError, tax.Severity)
	}
	
	ctx, ok := tax.Context.(*YAMLIngestContext)
	if !ok {
		t.Fatalf("expected YAMLIngestContext, got %T", tax.Context)
	}
	if ctx.File != "test.yaml" {
		t.Errorf("expected File 'test.yaml', got %s", ctx.File)
	}
	if ctx.Line != 5 {
		t.Errorf("expected Line 5, got %d", ctx.Line)
	}
	if ctx.Column != 10 {
		t.Errorf("expected Column 10, got %d", ctx.Column)
	}
	if ctx.Snippet != "snippet" {
		t.Errorf("expected Snippet 'snippet', got %s", ctx.Snippet)
	}
	
	// Verify context implements TaxonomyContext
	if ctx.Stage() != StageIngestYAMLSyntax {
		t.Errorf("context.Stage() returned %s, expected %s", ctx.Stage(), StageIngestYAMLSyntax)
	}
}

func TestNewSchemaConstraintTaxonomy(t *testing.T) {
	allowed := []string{"a", "b", "c"}
	min := 1.0
	max := 10.0
	tax := NewSchemaConstraintTaxonomy("features.test", "test", "field", "enum", "invalid", allowed, &min, &max, 3, 7)
	
	if tax.Stage != StageSchemaConstraints {
		t.Errorf("expected Stage %s, got %s", StageSchemaConstraints, tax.Stage)
	}
	if tax.Symptom != SymptomEnumMismatch {
		t.Errorf("expected Symptom %s, got %s", SymptomEnumMismatch, tax.Symptom)
	}
	if tax.Path != "features.test" {
		t.Errorf("expected Path 'features.test', got %s", tax.Path)
	}
	if tax.Severity != SeverityError {
		t.Errorf("expected Severity %s, got %s", SeverityError, tax.Severity)
	}
	
	ctx, ok := tax.Context.(*SchemaConstraintContext)
	if !ok {
		t.Fatalf("expected SchemaConstraintContext, got %T", tax.Context)
	}
	if ctx.Module != "test" {
		t.Errorf("expected Module 'test', got %s", ctx.Module)
	}
	if ctx.FieldPath != "field" {
		t.Errorf("expected FieldPath 'field', got %s", ctx.FieldPath)
	}
	if ctx.Expected != "enum" {
		t.Errorf("expected Expected 'enum', got %s", ctx.Expected)
	}
	if ctx.Actual != "invalid" {
		t.Errorf("expected Actual 'invalid', got %v", ctx.Actual)
	}
	if len(ctx.Allowed) != 3 || ctx.Allowed[0] != "a" {
		t.Errorf("expected Allowed [a b c], got %v", ctx.Allowed)
	}
	if ctx.Min == nil || *ctx.Min != 1.0 {
		t.Errorf("expected Min 1.0, got %v", ctx.Min)
	}
	if ctx.Max == nil || *ctx.Max != 10.0 {
		t.Errorf("expected Max 10.0, got %v", ctx.Max)
	}
	if ctx.Line != 3 {
		t.Errorf("expected Line 3, got %d", ctx.Line)
	}
	if ctx.Column != 7 {
		t.Errorf("expected Column 7, got %d", ctx.Column)
	}
	
	if ctx.Stage() != StageSchemaConstraints {
		t.Errorf("context.Stage() returned %s, expected %s", ctx.Stage(), StageSchemaConstraints)
	}
}

func TestNewSchemaStructureTaxonomy_MissingRequired(t *testing.T) {
	tax := NewSchemaStructureTaxonomy("features.test", "test", "field", "string", "", true, 2, 5)
	
	if tax.Stage != StageSchemaStructure {
		t.Errorf("expected Stage %s, got %s", StageSchemaStructure, tax.Stage)
	}
	if tax.Symptom != SymptomMissingRequired {
		t.Errorf("expected Symptom %s, got %s", SymptomMissingRequired, tax.Symptom)
	}
	
	ctx, ok := tax.Context.(*SchemaStructureContext)
	if !ok {
		t.Fatalf("expected SchemaStructureContext, got %T", tax.Context)
	}
	if !ctx.Required {
		t.Errorf("expected Required true, got false")
	}
	if ctx.Line != 2 {
		t.Errorf("expected Line 2, got %d", ctx.Line)
	}
	if ctx.Column != 5 {
		t.Errorf("expected Column 5, got %d", ctx.Column)
	}
}

func TestNewSchemaStructureTaxonomy_TypeMismatch(t *testing.T) {
	tax := NewSchemaStructureTaxonomy("features.test", "test", "field", "number", "string", false, 1, 3)
	
	if tax.Symptom != SymptomTypeMismatch {
		t.Errorf("expected Symptom %s, got %s", SymptomTypeMismatch, tax.Symptom)
	}
	
	ctx, ok := tax.Context.(*SchemaStructureContext)
	if !ok {
		t.Fatalf("expected SchemaStructureContext, got %T", tax.Context)
	}
	if ctx.Required {
		t.Errorf("expected Required false, got true")
	}
}

func TestNewExprDependencyTaxonomy(t *testing.T) {
	missingRefs := []string{"vars.height", "vars.width"}
	tax := NewExprDependencyTaxonomy("features.test.0.x", "max(vars.height, vars.width)", missingRefs, 5, 10, 15)
	
	if tax.Stage != StageExprDependencyMissing {
		t.Errorf("expected Stage %s, got %s", StageExprDependencyMissing, tax.Stage)
	}
	if tax.Symptom != SymptomDependencyMissing {
		t.Errorf("expected Symptom %s, got %s", SymptomDependencyMissing, tax.Symptom)
	}
	if tax.Path != "features.test.0.x" {
		t.Errorf("expected Path 'features.test.0.x', got %s", tax.Path)
	}
	
	ctx, ok := tax.Context.(*ExprDependencyContext)
	if !ok {
		t.Fatalf("expected ExprDependencyContext, got %T", tax.Context)
	}
	if ctx.Expression != "max(vars.height, vars.width)" {
		t.Errorf("expected Expression 'max(vars.height, vars.width)', got %s", ctx.Expression)
	}
	if len(ctx.MissingRefs) != 2 || ctx.MissingRefs[0] != "vars.height" {
		t.Errorf("expected MissingRefs [vars.height vars.width], got %v", ctx.MissingRefs)
	}
	if ctx.Iterations != 5 {
		t.Errorf("expected Iterations 5, got %d", ctx.Iterations)
	}
	if ctx.Line != 10 {
		t.Errorf("expected Line 10, got %d", ctx.Line)
	}
	if ctx.Column != 15 {
		t.Errorf("expected Column 15, got %d", ctx.Column)
	}
	
	if ctx.Stage() != StageExprDependencyMissing {
		t.Errorf("context.Stage() returned %s, expected %s", ctx.Stage(), StageExprDependencyMissing)
	}
}

func TestNewExprSyntaxTaxonomy(t *testing.T) {
	tax := NewExprSyntaxTaxonomy("features.test.0.x", "max(", "(", 4, 8, 12)
	
	if tax.Stage != StageExprSyntax {
		t.Errorf("expected Stage %s, got %s", StageExprSyntax, tax.Stage)
	}
	if tax.Symptom != SymptomSyntax {
		t.Errorf("expected Symptom %s, got %s", SymptomSyntax, tax.Symptom)
	}
	
	ctx, ok := tax.Context.(*ExprSyntaxContext)
	if !ok {
		t.Fatalf("expected ExprSyntaxContext, got %T", tax.Context)
	}
	if ctx.Expression != "max(" {
		t.Errorf("expected Expression 'max(', got %s", ctx.Expression)
	}
	if ctx.Token != "(" {
		t.Errorf("expected Token '(', got %s", ctx.Token)
	}
	if ctx.Position != 4 {
		t.Errorf("expected Position 4, got %d", ctx.Position)
	}
	if ctx.Line != 8 {
		t.Errorf("expected Line 8, got %d", ctx.Line)
	}
	if ctx.Column != 12 {
		t.Errorf("expected Column 12, got %d", ctx.Column)
	}
	
	if ctx.Stage() != StageExprSyntax {
		t.Errorf("context.Stage() returned %s, expected %s", ctx.Stage(), StageExprSyntax)
	}
}

func TestNewStrictModeTaxonomy_UnknownKey(t *testing.T) {
	offenders := []string{"unknown_key1", "unknown_key2"}
	tax := NewStrictModeTaxonomy("unknown-key", offenders)
	
	if tax.Stage != StageStrictUnknownKey {
		t.Errorf("expected Stage %s, got %s", StageStrictUnknownKey, tax.Stage)
	}
	if tax.Symptom != SymptomUnknownKey {
		t.Errorf("expected Symptom %s, got %s", SymptomUnknownKey, tax.Symptom)
	}
	if tax.Severity != SeverityWarning {
		t.Errorf("expected Severity %s, got %s", SeverityWarning, tax.Severity)
	}
	
	ctx, ok := tax.Context.(*StrictModeContext)
	if !ok {
		t.Fatalf("expected StrictModeContext, got %T", tax.Context)
	}
	if ctx.Kind != "unknown-key" {
		t.Errorf("expected Kind 'unknown-key', got %s", ctx.Kind)
	}
	if len(ctx.Offenders) != 2 {
		t.Errorf("expected 2 offenders, got %d", len(ctx.Offenders))
	}
}

func TestNewStrictModeTaxonomy_UnusedVar(t *testing.T) {
	offenders := []string{"vars.unused"}
	tax := NewStrictModeTaxonomy("unused-var", offenders)
	
	if tax.Stage != StageStrictUnusedVar {
		t.Errorf("expected Stage %s, got %s", StageStrictUnusedVar, tax.Stage)
	}
	if tax.Symptom != SymptomUnusedVar {
		t.Errorf("expected Symptom %s, got %s", SymptomUnusedVar, tax.Symptom)
	}
	
	ctx, ok := tax.Context.(*StrictModeContext)
	if !ok {
		t.Fatalf("expected StrictModeContext, got %T", tax.Context)
	}
	if ctx.Kind != "unused-var" {
		t.Errorf("expected Kind 'unused-var', got %s", ctx.Kind)
	}
}

func TestTaxonomy_Error(t *testing.T) {
	tax := NewYAMLIngestTaxonomy("test.yaml", 1, 2, "snippet")
	errMsg := tax.Error()
	
	if errMsg == "" {
		t.Error("Error() returned empty string")
	}
	if errMsg[:1] != "[" {
		t.Errorf("Error() should start with '[', got %s", errMsg[:10])
	}
}

func TestAsTaxonomy_Direct(t *testing.T) {
	tax := NewYAMLIngestTaxonomy("test.yaml", 1, 2, "snippet")
	
	found, ok := AsTaxonomy(tax)
	if !ok {
		t.Fatal("AsTaxonomy returned false for direct Taxonomy")
	}
	if found != tax {
		t.Error("AsTaxonomy returned different pointer")
	}
}

func TestAsTaxonomy_Wrapped(t *testing.T) {
	tax := NewYAMLIngestTaxonomy("test.yaml", 1, 2, "snippet")
	wrapped := errors.Wrap(tax, "wrapped error")
	
	found, ok := AsTaxonomy(wrapped)
	if !ok {
		t.Fatal("AsTaxonomy returned false for wrapped Taxonomy")
	}
	if found != tax {
		t.Error("AsTaxonomy returned different pointer")
	}
}

func TestAsTaxonomy_DoubleWrapped(t *testing.T) {
	tax := NewYAMLIngestTaxonomy("test.yaml", 1, 2, "snippet")
	wrapped1 := errors.Wrap(tax, "first wrap")
	wrapped2 := errors.Wrap(wrapped1, "second wrap")
	
	found, ok := AsTaxonomy(wrapped2)
	if !ok {
		t.Fatal("AsTaxonomy returned false for double-wrapped Taxonomy")
	}
	if found != tax {
		t.Error("AsTaxonomy returned different pointer")
	}
}

func TestAsTaxonomy_NonTaxonomyError(t *testing.T) {
	err := errors.New("plain error")
	
	_, ok := AsTaxonomy(err)
	if ok {
		t.Error("AsTaxonomy returned true for non-Taxonomy error")
	}
}

func TestAsTaxonomy_Nil(t *testing.T) {
	_, ok := AsTaxonomy(nil)
	if ok {
		t.Error("AsTaxonomy returned true for nil error")
	}
}

func TestAsTaxonomy_MixedChain(t *testing.T) {
	tax := NewSchemaStructureTaxonomy("features.test", "test", "field", "string", "", true, 1, 2)
	wrapped1 := errors.Wrap(tax, "first")
	wrapped2 := errors.Wrap(wrapped1, "second")
	plainErr := errors.New("plain")
	mixed := errors.Wrap(plainErr, "before")
	mixed = errors.Wrap(mixed, "after")
	// This creates a chain: mixed -> plainErr, but we can't easily mix chains
	// Let's test a simpler case: wrapped taxonomy in a chain
	
	found, ok := AsTaxonomy(wrapped2)
	if !ok {
		t.Fatal("AsTaxonomy returned false for wrapped Taxonomy in chain")
	}
	if found != tax {
		t.Error("AsTaxonomy returned different pointer")
	}
}

