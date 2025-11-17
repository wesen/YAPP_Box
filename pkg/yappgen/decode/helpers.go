package decode

import (
	"fmt"

	"github.com/pkg/errors"
)

// FormatLabel formats a consistent label for module entries (e.g. pcb_stands[0]).
func FormatLabel(module string, idx int) string {
	return fmt.Sprintf("%s[%d]", module, idx)
}

// GetFloat extracts a numeric value from the map.
func GetFloat(m map[string]any, key, label string) (float64, error) {
	val, ok := m[key]
	if !ok {
		return 0, missingFieldErr(label, key)
	}
	return coerceNumber(val, label, key)
}

// GetOptionalFloat extracts an optional numeric value from the map.
func GetOptionalFloat(m map[string]any, key, label string) (*float64, error) {
	val, ok := m[key]
	if !ok || val == nil {
		return nil, nil
	}
	num, err := coerceNumber(val, label, key)
	if err != nil {
		return nil, err
	}
	return &num, nil
}

// GetString extracts a string value from the map.
func GetString(m map[string]any, key, label string) (string, error) {
	val, ok := m[key]
	if !ok {
		return "", missingFieldErr(label, key)
	}
	str, err := coerceString(val, label, key)
	if err != nil {
		return "", err
	}
	return str, nil
}

// GetOptionalString extracts an optional string value from the map.
func GetOptionalString(m map[string]any, key, label string) (*string, error) {
	val, ok := m[key]
	if !ok || val == nil {
		return nil, nil
	}
	str, err := coerceString(val, label, key)
	if err != nil {
		return nil, err
	}
	return &str, nil
}

// GetBool extracts a bool value from the map.
func GetBool(m map[string]any, key, label string) (bool, error) {
	val, ok := m[key]
	if !ok {
		return false, missingFieldErr(label, key)
	}
	boolean, err := coerceBool(val, label, key)
	if err != nil {
		return false, err
	}
	return boolean, nil
}

// GetOptionalBool extracts an optional bool value from the map.
func GetOptionalBool(m map[string]any, key, label string) (*bool, error) {
	val, ok := m[key]
	if !ok || val == nil {
		return nil, nil
	}
	boolean, err := coerceBool(val, label, key)
	if err != nil {
		return nil, err
	}
	return &boolean, nil
}

// GetObject extracts a nested map from the map.
func GetObject(m map[string]any, key, label string) (map[string]any, error) {
	val, ok := m[key]
	if !ok {
		return nil, missingFieldErr(label, key)
	}
	obj, err := coerceObject(val, label, key)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// GetOptionalObject extracts an optional nested map from the map.
func GetOptionalObject(m map[string]any, key, label string) (map[string]any, error) {
	val, ok := m[key]
	if !ok || val == nil {
		return nil, nil
	}
	obj, err := coerceObject(val, label, key)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// GetArray extracts an array from the map.
func GetArray(m map[string]any, key, label string) ([]any, error) {
	val, ok := m[key]
	if !ok {
		return nil, missingFieldErr(label, key)
	}
	arr, err := coerceArray(val, label, key)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

// GetOptionalArray extracts an optional array from the map.
func GetOptionalArray(m map[string]any, key, label string) ([]any, error) {
	val, ok := m[key]
	if !ok || val == nil {
		return nil, nil
	}
	arr, err := coerceArray(val, label, key)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

func missingFieldErr(label, key string) error {
	return errors.Errorf("%s.%s: missing required field", label, key)
}

func typeErr(label, key, expected string, actual any) error {
	return errors.Errorf("%s.%s: expected %s, got %T", label, key, expected, actual)
}

func coerceNumber(value any, label, key string) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case *float64:
		if v == nil {
			return 0, typeErr(label, key, "number", nil)
		}
		return *v, nil
	case *int:
		if v == nil {
			return 0, typeErr(label, key, "number", nil)
		}
		return float64(*v), nil
	case *int64:
		if v == nil {
			return 0, typeErr(label, key, "number", nil)
		}
		return float64(*v), nil
	default:
		return 0, typeErr(label, key, "number", value)
	}
}

func coerceString(value any, label, key string) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case *string:
		if v == nil {
			return "", typeErr(label, key, "string", nil)
		}
		return *v, nil
	default:
		return "", typeErr(label, key, "string", value)
	}
}

func coerceBool(value any, label, key string) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case *bool:
		if v == nil {
			return false, typeErr(label, key, "bool", nil)
		}
		return *v, nil
	default:
		return false, typeErr(label, key, "bool", value)
	}
}

func coerceObject(value any, label, key string) (map[string]any, error) {
	obj, ok := value.(map[string]any)
	if !ok {
		return nil, typeErr(label, key, "object", value)
	}
	return obj, nil
}

func coerceArray(value any, label, key string) ([]any, error) {
	arr, ok := value.([]any)
	if !ok {
		return nil, typeErr(label, key, "array", value)
	}
	return arr, nil
}
