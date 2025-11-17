package decode

import (
	"strings"
	"testing"
)

func TestFormatLabel(t *testing.T) {
	if got := FormatLabel("pcb_stands", 3); got != "pcb_stands[3]" {
		t.Fatalf("unexpected label: %s", got)
	}
}

func TestGetFloat(t *testing.T) {
	label := "pcb_stands[0]"
	m := map[string]any{
		"x": 10,
	}
	val, err := GetFloat(m, "x", label)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 10 {
		t.Fatalf("unexpected value: %v", val)
	}

	if _, err := GetFloat(m, "y", label); err == nil || !strings.Contains(err.Error(), "missing required field") {
		t.Fatalf("expected missing error, got %v", err)
	}

	m["bad"] = "oops"
	if _, err := GetFloat(m, "bad", label); err == nil || !strings.Contains(err.Error(), "expected number") {
		t.Fatalf("expected type error, got %v", err)
	}
}

func TestGetOptionalFloat(t *testing.T) {
	label := "module[0]"
	m := map[string]any{}
	if val, err := GetOptionalFloat(m, "x", label); err != nil || val != nil {
		t.Fatalf("expected nil optional float, got %v err=%v", val, err)
	}
	m["x"] = 5.5
	val, err := GetOptionalFloat(m, "x", label)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val == nil || *val != 5.5 {
		t.Fatalf("unexpected value: %v", val)
	}
	m["x"] = "oops"
	if _, err := GetOptionalFloat(m, "x", label); err == nil || !strings.Contains(err.Error(), "expected number") {
		t.Fatalf("expected type error, got %v", err)
	}
}

func TestGetString(t *testing.T) {
	label := "module[1]"
	m := map[string]any{
		"name": "foo",
	}
	s, err := GetString(m, "name", label)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != "foo" {
		t.Fatalf("unexpected value: %s", s)
	}
	if _, err := GetString(m, "missing", label); err == nil || !strings.Contains(err.Error(), "missing required field") {
		t.Fatalf("expected missing error, got %v", err)
	}
	m["bad"] = 42
	if _, err := GetString(m, "bad", label); err == nil || !strings.Contains(err.Error(), "expected string") {
		t.Fatalf("expected type error, got %v", err)
	}
}

func TestGetBool(t *testing.T) {
	label := "module[2]"
	m := map[string]any{
		"flag": true,
	}
	val, err := GetBool(m, "flag", label)
	if err != nil || !val {
		t.Fatalf("unexpected result: %v err=%v", val, err)
	}
	if _, err := GetBool(m, "missing", label); err == nil || !strings.Contains(err.Error(), "missing required field") {
		t.Fatalf("expected missing error, got %v", err)
	}
	m["flag"] = "nope"
	if _, err := GetBool(m, "flag", label); err == nil || !strings.Contains(err.Error(), "expected bool") {
		t.Fatalf("expected type error, got %v", err)
	}
}

func TestGetObject(t *testing.T) {
	label := "module[3]"
	m := map[string]any{
		"child": map[string]any{"x": 1},
	}
	child, err := GetObject(m, "child", label)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if child["x"] != 1 {
		t.Fatalf("unexpected child: %v", child)
	}
	if _, err := GetObject(m, "missing", label); err == nil || !strings.Contains(err.Error(), "missing required field") {
		t.Fatalf("expected missing error, got %v", err)
	}
	m["child"] = "nope"
	if _, err := GetObject(m, "child", label); err == nil || !strings.Contains(err.Error(), "expected object") {
		t.Fatalf("expected type error, got %v", err)
	}
}

func TestGetArray(t *testing.T) {
	label := "module[4]"
	m := map[string]any{
		"items": []any{1, 2, 3},
	}
	arr, err := GetArray(m, "items", label)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(arr) != 3 {
		t.Fatalf("unexpected length: %d", len(arr))
	}
	if _, err := GetArray(m, "missing", label); err == nil || !strings.Contains(err.Error(), "missing required field") {
		t.Fatalf("expected missing error, got %v", err)
	}
	m["items"] = "nope"
	if _, err := GetArray(m, "items", label); err == nil || !strings.Contains(err.Error(), "expected array") {
		t.Fatalf("expected type error, got %v", err)
	}
}
