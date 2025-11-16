package registry

import (
	"testing"
)

type mockSchema struct {
	name string
	path string
}

func (m mockSchema) Name() string                                    { return m.name }
func (m mockSchema) Path() string                                    { return m.path }
func (m mockSchema) Description() string                             { return "mock schema" }
func (m mockSchema) Fields() []FieldSpec                             { return nil }
func (m mockSchema) ValidateStructure(path string, data any) error   { return nil }
func (m mockSchema) ValidateConstraints(path string, data any) error { return nil }

type mockModule struct {
	schema ModuleSchema
}

var _ FeatureModule = (*mockModule)(nil)

func (m *mockModule) Schema() ModuleSchema {
	return m.schema
}

func (m *mockModule) Build(items []map[string]any) ([][]any, error) {
	return nil, nil
}

func TestRegisterPreservesOrder(t *testing.T) {
	SetTestRegistry(nil)

	modA := &mockModule{schema: mockSchema{name: "a", path: "features.a"}}
	modB := &mockModule{schema: mockSchema{name: "b", path: "features.b"}}

	Register(modA)
	Register(modB)

	all := All()
	if len(all) != 2 {
		t.Fatalf("expected 2 modules, got %d", len(all))
	}
	if all[0] != modA {
		t.Fatalf("expected first module to be modA")
	}
	if all[1] != modB {
		t.Fatalf("expected second module to be modB")
	}
}

func TestRegisterPanicsOnDuplicatePath(t *testing.T) {
	SetTestRegistry(nil)

	mod := &mockModule{schema: mockSchema{name: "dup", path: "features.dup"}}
	Register(mod)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on duplicate registration")
		}
	}()

	Register(mod)
}

func TestSetTestRegistryOverridesState(t *testing.T) {
	a := &mockModule{schema: mockSchema{name: "alpha", path: "features.alpha"}}
	b := &mockModule{schema: mockSchema{name: "beta", path: "features.beta"}}

	SetTestRegistry(map[string]FeatureModule{
		"features.beta":  b,
		"features.alpha": a,
	})

	if got, ok := Get("features.alpha"); !ok || got != a {
		t.Fatalf("expected alpha module to be present")
	}

	all := All()
	if len(all) != 2 {
		t.Fatalf("expected 2 modules in All, got %d", len(all))
	}
	if all[0] != a || all[1] != b {
		t.Fatalf("expected deterministic ordering [alpha, beta], got %+v", all)
	}
}
