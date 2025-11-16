package schemagen

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// ValidateSchemaFile reads a schema YAML file and validates its structure.
func ValidateSchemaFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "read schema file")
	}
	errs, err := ValidateSchemaBytes(path, data)
	if err != nil {
		return err
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// ValidateSchemaBytes validates a schema document represented as bytes.
func ValidateSchemaBytes(filename string, data []byte) (ValidationErrors, error) {
	root, err := parseDocument(data)
	if err != nil {
		return nil, errors.Wrap(err, "parse schema YAML")
	}
	return validateSchemaNode(filename, root), nil
}

func parseDocument(data []byte) (*yaml.Node, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(false)

	var doc yaml.Node
	if err := decoder.Decode(&doc); err != nil {
		return nil, err
	}
	if len(doc.Content) == 0 {
		return nil, errors.New("empty schema document")
	}
	return doc.Content[0], nil
}

func validateSchemaNode(filename string, root *yaml.Node) ValidationErrors {
	var errs ValidationErrors

	if root.Kind != yaml.MappingNode {
		errs.add(ValidationError{
			File:    filename,
			Path:    "",
			Line:    root.Line,
			Column:  root.Column,
			Message: "schema root must be a mapping/object",
		})
		return errs
	}

	requiredStrings := []string{"module", "scad_array", "go_package", "description"}
	for _, key := range requiredStrings {
		if node := mapValue(root, key); node == nil {
			errs.add(newError(filename, root, key, fmt.Sprintf("missing required field %q", key)))
		} else if _, err := readString(node); err != nil {
			errs.add(newError(filename, node, key, fmt.Sprintf("%q must be a string", key)))
		} else if strings.TrimSpace(node.Value) == "" {
			errs.add(newError(filename, node, key, fmt.Sprintf("%q cannot be empty", key)))
		}
	}

	if orderNode := mapValue(root, "order"); orderNode != nil {
		if _, err := readInt(orderNode); err != nil {
			errs.add(newError(filename, orderNode, "order", "order must be an integer"))
		}
	}

	fieldsNode := mapValue(root, "fields")
	if fieldsNode == nil {
		errs.add(newError(filename, root, "fields", "missing required section 'fields'"))
	} else if fieldsNode.Kind != yaml.MappingNode {
		errs.add(newError(filename, fieldsNode, "fields", "'fields' must be a mapping of field definitions"))
	} else {
		for i := 0; i < len(fieldsNode.Content); i += 2 {
			keyNode := fieldsNode.Content[i]
			valueNode := fieldsNode.Content[i+1]
			fieldName := keyNode.Value
			errs = append(errs, validateField(filename, fieldName, valueNode, fmt.Sprintf("fields.%s", fieldName))...)
		}
	}

	if testsNode := mapValue(root, "tests"); testsNode != nil {
		errs = append(errs, validateTests(filename, testsNode)...)
	}

	return errs
}

func validateField(filename, fieldName string, node *yaml.Node, path string) ValidationErrors {
	var errs ValidationErrors

	if node.Kind != yaml.MappingNode {
		errs.add(newError(filename, node, path, fmt.Sprintf("field %q must be an object", fieldName)))
		return errs
	}

	typeNode := mapValue(node, "type")
	if typeNode == nil {
		errs.add(newError(filename, node, path, "missing required key 'type'"))
		return errs
	}
	typeStr, err := readString(typeNode)
	if err != nil {
		errs.add(newError(filename, typeNode, path+".type", "'type' must be a string"))
		return errs
	}

	switch typeStr {
	case "number", "string":
		// no extra validation
	case "array":
		// arrays can optionally specify element structure via fields/items
		if fieldsNode := mapValue(node, "fields"); fieldsNode != nil && fieldsNode.Kind != yaml.SequenceNode && fieldsNode.Kind != yaml.MappingNode {
			errs.add(newError(filename, fieldsNode, path+".fields", "'fields' inside array must be mapping or sequence"))
		}
	case "object":
		fieldsNode := mapValue(node, "fields")
		if fieldsNode == nil {
			errs.add(newError(filename, node, path, "object field must define nested 'fields'"))
		} else if fieldsNode.Kind != yaml.MappingNode {
			errs.add(newError(filename, fieldsNode, path+".fields", "'fields' must be a mapping"))
		} else {
			for i := 0; i < len(fieldsNode.Content); i += 2 {
				childKey := fieldsNode.Content[i]
				childVal := fieldsNode.Content[i+1]
				childName := childKey.Value
				childPath := fmt.Sprintf("%s.%s", path, childName)
				errs = append(errs, validateField(filename, childName, childVal, childPath)...)
			}
		}
	default:
		errs.add(newError(filename, typeNode, path+".type", fmt.Sprintf("unknown field type %q (valid: number, string, object, array)", typeStr)))
	}

	if requiredNode := mapValue(node, "required"); requiredNode != nil {
		if _, err := readBool(requiredNode); err != nil {
			errs.add(newError(filename, requiredNode, path+".required", "'required' must be true or false"))
		}
	}

	if minNode := mapValue(node, "min"); minNode != nil {
		if _, err := readNumber(minNode); err != nil {
			errs.add(newError(filename, minNode, path+".min", "'min' must be a number"))
		}
	}

	if maxNode := mapValue(node, "max"); maxNode != nil {
		if _, err := readNumber(maxNode); err != nil {
			errs.add(newError(filename, maxNode, path+".max", "'max' must be a number"))
		}
	}

	if enumNode := mapValue(node, "enum"); enumNode != nil {
		if enumNode.Kind != yaml.SequenceNode {
			errs.add(newError(filename, enumNode, path+".enum", "'enum' must be a list of strings"))
		} else {
			for idx, entry := range enumNode.Content {
				if _, err := readString(entry); err != nil {
					errs.add(newError(filename, entry, fmt.Sprintf("%s.enum[%d]", path, idx), "enum values must be strings"))
				}
			}
		}
	}

	return errs
}

func validateTests(filename string, node *yaml.Node) ValidationErrors {
	var errs ValidationErrors

	if node.Kind != yaml.SequenceNode {
		errs.add(newError(filename, node, "tests", "'tests' must be a list"))
		return errs
	}

	for idx, testNode := range node.Content {
		path := fmt.Sprintf("tests[%d]", idx)
		if testNode.Kind != yaml.MappingNode {
			errs.add(newError(filename, testNode, path, "test case must be an object"))
			continue
		}

		nameNode := mapValue(testNode, "name")
		if nameNode == nil {
			errs.add(newError(filename, testNode, path, "test case missing 'name'"))
		} else if _, err := readString(nameNode); err != nil {
			errs.add(newError(filename, nameNode, path+".name", "'name' must be a string"))
		}

		inputNode := mapValue(testNode, "input")
		if inputNode == nil {
			errs.add(newError(filename, testNode, path, "test case missing 'input'"))
		}

		expectValid := mapValue(testNode, "expect_valid")
		expectError := mapValue(testNode, "expect_error")
		if expectValid == nil && expectError == nil {
			errs.add(newError(filename, testNode, path, "test case must define 'expect_valid' or 'expect_error'"))
		}

		if expectValid != nil {
			if _, err := readBool(expectValid); err != nil {
				errs.add(newError(filename, expectValid, path+".expect_valid", "'expect_valid' must be true or false"))
			}
		}
		if expectError != nil {
			if _, err := readString(expectError); err != nil {
				errs.add(newError(filename, expectError, path+".expect_error", "'expect_error' must be a string"))
			}
		}
	}

	return errs
}

func mapValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(node.Content); i += 2 {
		k := node.Content[i]
		if k.Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

func readString(node *yaml.Node) (string, error) {
	if node == nil || node.Kind != yaml.ScalarNode {
		return "", errors.New("not a scalar string")
	}
	return node.Value, nil
}

func readBool(node *yaml.Node) (bool, error) {
	if node == nil || node.Kind != yaml.ScalarNode {
		return false, errors.New("not a scalar bool")
	}
	switch strings.ToLower(node.Value) {
	case "true", "yes", "on":
		return true, nil
	case "false", "no", "off":
		return false, nil
	default:
		return false, errors.New("invalid boolean value")
	}
}

func readNumber(node *yaml.Node) (float64, error) {
	if node == nil || node.Kind != yaml.ScalarNode {
		return 0, errors.New("not a scalar number")
	}
	return strconv.ParseFloat(node.Value, 64)
}

func readInt(node *yaml.Node) (int64, error) {
	if node == nil || node.Kind != yaml.ScalarNode {
		return 0, errors.New("not a scalar integer")
	}
	return strconv.ParseInt(node.Value, 10, 64)
}

func newError(filename string, node *yaml.Node, path string, msg string) ValidationError {
	line, column := 0, 0
	if node != nil {
		line = node.Line
		column = node.Column
	}
	return ValidationError{
		File:    filename,
		Path:    path,
		Line:    line,
		Column:  column,
		Message: msg,
	}
}
