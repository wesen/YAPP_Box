package resolvercli

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

// PositionMap maps DSL paths to their line/column positions in the source YAML.
type PositionMap map[string]Position

// Position represents a location in the source file.
type Position struct {
	Line   int
	Column int
}

// buildPositionMap traverses the yaml.Node tree and builds a map from paths to positions.
func buildPositionMap(node *yaml.Node, path string, positions PositionMap) {
	if node == nil {
		return
	}

	// Record position for this path if it's a scalar (leaf node)
	if node.Kind == yaml.ScalarNode && path != "" {
		positions[path] = Position{
			Line:   node.Line,
			Column: node.Column,
		}
	}

	switch node.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valNode := node.Content[i+1]
			key := keyNode.Value
			childPath := joinPath(path, key)
			// Record position for the key
			if path != "" {
				positions[childPath] = Position{
					Line:   keyNode.Line,
					Column: keyNode.Column,
				}
			}
			buildPositionMap(valNode, childPath, positions)
		}
	case yaml.SequenceNode:
		for idx, child := range node.Content {
			childPath := joinPath(path, fmt.Sprintf("%d", idx))
			buildPositionMap(child, childPath, positions)
		}
	}
}

// getPosition looks up the position for a given path, returning 0,0 if not found.
func (pm PositionMap) getPosition(path string) (line, column int) {
	if pos, ok := pm[path]; ok {
		return pos.Line, pos.Column
	}
	return 0, 0
}

