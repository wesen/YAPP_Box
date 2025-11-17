package assets

import _ "embed"

//go:embed YAPPgenerator_v3.scad
var YAPPGeneratorSCAD []byte

// GetGeneratorSCAD returns the embedded YAPP generator SCAD content.
func GetGeneratorSCAD() []byte {
	return YAPPGeneratorSCAD
}

