package generatorcli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/yappgen"
	"github.com/wesen/yapp-encl-resolver/pkg/yappgen/assets"
)

const defaultGeneratorInclude = "YAPPgenerator_v3.scad"

// SCADOptions controls how the SCAD file is written.
type SCADOptions struct {
	// OutputPath is the location where the SCAD file should be written.
	OutputPath string
	// GeneratorInclude is an optional override for the include path. When empty,
	// defaults to YAPPgenerator_v3.scad in the repo root.
	GeneratorInclude string
	// CopyGenerator, when true, extracts the embedded generator to the output directory.
	CopyGenerator bool
	// GeneratorPath is an optional path to a custom YAPPgenerator_v3.scad file.
	// If empty, uses the embedded version.
	GeneratorPath string
}

// WriteSCAD builds the YAPP model, emits the SCAD, adjusts include paths, and writes the file.
func WriteSCAD(ctx context.Context, resolved map[string]any, opts SCADOptions) (string, *yappgen.Model, error) {
	if opts.OutputPath == "" {
		return "", nil, errors.New("output path is required")
	}
	model, err := yappgen.BuildModel(ctx, resolved)
	if err != nil {
		return "", nil, errors.Wrap(err, "build model")
	}
	scad, err := yappgen.EmitSCAD(ctx, model)
	if err != nil {
		return "", nil, errors.Wrap(err, "emit scad")
	}

	// Handle generator copying
	// Copy generator to output directory for standalone SCAD files or when explicitly requested
	var generatorInclude string
	if opts.CopyGenerator {
		generatorInclude, err = copyGeneratorToOutput(opts.OutputPath, opts.GeneratorPath)
		if err != nil {
			return "", nil, errors.Wrap(err, "copy generator")
		}
	} else {
		generatorInclude = opts.GeneratorInclude
	}

	scad = rewriteInclude(scad, opts.OutputPath, generatorInclude)

	if err := os.WriteFile(opts.OutputPath, scad, 0o644); err != nil {
		return "", nil, errors.Wrapf(err, "write SCAD %s", opts.OutputPath)
	}
	return opts.OutputPath, model, nil
}

// STLOptions controls how OpenSCAD rendering runs.
type STLOptions struct {
	BasePath string
	LidPath  string
	OpenSCAD string
	// PrintSwitchExtenders toggles tactile button extender geometry.
	PrintSwitchExtenders bool
}

// RenderSTLs renders any requested STL files using OpenSCAD.
func RenderSTLs(ctx context.Context, scadPath string, opts STLOptions) error {
	if opts.BasePath == "" && opts.LidPath == "" {
		return nil
	}
	bin := opts.OpenSCAD
	if bin == "" {
		bin = "openscad"
	}
	if opts.BasePath != "" {
		if err := renderSingleSTL(ctx, bin, scadPath, opts.BasePath, true, false, opts.PrintSwitchExtenders); err != nil {
			return errors.Wrap(err, "base STL")
		}
	}
	if opts.LidPath != "" {
		if err := renderSingleSTL(ctx, bin, scadPath, opts.LidPath, false, true, opts.PrintSwitchExtenders); err != nil {
			return errors.Wrap(err, "lid STL")
		}
	}
	return nil
}

func renderSingleSTL(ctx context.Context, bin, scadPath, outPath string, printBase, printLid bool, printExtenders bool) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return errors.Wrap(err, "create STL directory")
	}
	fmt.Printf("Rendering STL %s from %s\n", outPath, scadPath)
	args := []string{
		"-o", outPath,
		"-D", fmt.Sprintf("printBaseShell=%t", printBase),
		"-D", fmt.Sprintf("printLidShell=%t", printLid),
		"-D", fmt.Sprintf("printSwitchExtenders=%t", printExtenders),
		"-D", "printDisplayClips=false",
		scadPath,
	}
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "openscad")
	}
	return nil
}

func rewriteInclude(scad []byte, outputPath, generatorInclude string) []byte {
	if generatorInclude == "" {
		generatorInclude = defaultGeneratorInclude
	}
	
	// If generatorInclude is just a filename (no path), it's in the same directory
	if filepath.Base(generatorInclude) == generatorInclude {
		// Use relative path in same directory
		updated := strings.Replace(string(scad), "include <./YAPPgenerator_v3.scad>", fmt.Sprintf("include <%s>", generatorInclude), 1)
		return []byte(updated)
	}
	
	// Otherwise compute relative path from output to generator
	genAbs, err := filepath.Abs(generatorInclude)
	if err != nil {
		return scad
	}
	outDir := filepath.Dir(outputPath)
	outDirAbs, err := filepath.Abs(outDir)
	if err != nil {
		return scad
	}
	rel, err := filepath.Rel(outDirAbs, genAbs)
	if err != nil {
		return scad
	}
	rel = filepath.ToSlash(rel)
	updated := strings.Replace(string(scad), "include <./YAPPgenerator_v3.scad>", fmt.Sprintf("include <%s>", rel), 1)
	return []byte(updated)
}

// copyGeneratorToOutput extracts the generator to the output directory.
// Returns the path to use in the include directive (relative to output SCAD).
func copyGeneratorToOutput(outputPath, customGeneratorPath string) (string, error) {
	outDir := filepath.Dir(outputPath)
	generatorDest := filepath.Join(outDir, "YAPPgenerator_v3.scad")

	var generatorContent []byte
	if customGeneratorPath != "" {
		// Use custom generator path
		content, err := os.ReadFile(customGeneratorPath)
		if err != nil {
			return "", errors.Wrapf(err, "read custom generator %s", customGeneratorPath)
		}
		generatorContent = content
	} else {
		// Use embedded generator
		generatorContent = assets.GetGeneratorSCAD()
	}

	if err := os.WriteFile(generatorDest, generatorContent, 0o644); err != nil {
		return "", errors.Wrapf(err, "write generator to %s", generatorDest)
	}

	fmt.Printf("Copied YAPPgenerator_v3.scad to %s\n", generatorDest)
	return "YAPPgenerator_v3.scad", nil // Relative path in same directory
}
