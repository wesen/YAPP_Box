package resolver_test

import (
	"context"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver"
	"gopkg.in/yaml.v3"
)

func mustUnmarshalYAML(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := yaml.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return m
}

func TestSimpleEvaluation(t *testing.T) {
	in := mustUnmarshalYAML(t, `
project: test
version: 0
units: mm
yapp_version: v0
vars:
  add: pcb.z_clearance + 4.0
pcb:
  z_clearance: 2.0
enclosure:
  base:
    thickness: add
`)
	out, err := resolver.Resolve(context.Background(), in, resolver.Options{})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	encl := out["enclosure"].(map[string]any)
	base := encl["base"].(map[string]any)
	th := base["thickness"]
	if th != 6 {
		t.Fatalf("expected 6 (int), got %#v", th)
	}
}

func TestFunctionsAndInts(t *testing.T) {
	in := mustUnmarshalYAML(t, `
project: test
version: 0
units: mm
yapp_version: v0
enclosure:
  wall:
    thickness: 2
  base:
    thickness: max(2.0, enclosure.wall.thickness)
`)
	out, err := resolver.Resolve(context.Background(), in, resolver.Options{})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	encl := out["enclosure"].(map[string]any)
	base := encl["base"].(map[string]any)
	if base["thickness"] != 2 {
		t.Fatalf("expected 2, got %#v", base["thickness"])
	}
}

func TestUnresolvedFails(t *testing.T) {
	in := mustUnmarshalYAML(t, `
project: test
version: 0
units: mm
yapp_version: v0
enclosure:
  base:
    thickness: missing.path + 1
`)
	_, err := resolver.Resolve(context.Background(), in, resolver.Options{MaxIterations: 2})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStrictUnusedVars(t *testing.T) {
	in := mustUnmarshalYAML(t, `
project: test
version: 0
units: mm
yapp_version: v0
vars:
  never_used: 2
enclosure:
  wall:
    thickness: 2
`)
	_, err := resolver.Resolve(context.Background(), in, resolver.Options{Strict: true})
	if err == nil {
		t.Fatal("expected strict error on unused vars")
	}
}

func TestIterationChain(t *testing.T) {
	in := mustUnmarshalYAML(t, `
project: test
version: 0
units: mm
yapp_version: v0
vars:
  a: 1
  b: a + 1
  c: b + 1
  d: c + 1
enclosure:
  base:
    thickness: d
`)
	out, err := resolver.Resolve(context.Background(), in, resolver.Options{})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	encl := out["enclosure"].(map[string]any)
	base := encl["base"].(map[string]any)
	if base["thickness"] != 4 {
		t.Fatalf("expected 4, got %#v", base["thickness"])
	}
}


