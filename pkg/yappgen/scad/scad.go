package scad

// Raw represents an identifier or token that must be emitted as-is (no quotes).
type Raw string

// Undef is the OpenSCAD 'undef' literal.
const Undef Raw = "undef"
