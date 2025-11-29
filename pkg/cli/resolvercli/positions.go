package resolvercli

import (
	"fmt"
	"gopkg.in/yaml.v3"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver"
)

// buildPositionMap traverses the yaml.Node tree and builds a map from paths to positions.
func buildPositionMap(node *yaml.Node, path string, positions resolver.PositionMap) {
	if node == nil {
		return
	}

	// Record position for this path if it's a scalar (leaf node)
	if node.Kind == yaml.ScalarNode && path != "" {
		positions[path] = resolver.Position{
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
				positions[childPath] = resolver.Position{
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


