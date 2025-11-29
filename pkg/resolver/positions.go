package resolver

// PositionMap maps DSL paths to their line/column positions in the source YAML.
type PositionMap map[string]Position

// Position represents a location in the source file.
type Position struct {
	Line   int
	Column int
}

// GetPosition looks up the position for a given path, returning 0,0 if not found.
func (pm PositionMap) GetPosition(path string) (line, column int) {
	if pos, ok := pm[path]; ok {
		return pos.Line, pos.Column
	}
	return 0, 0
}

