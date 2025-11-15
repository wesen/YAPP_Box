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
)

const defaultGeneratorInclude = "YAPPgenerator_v3.scad"

// SCADOptions controls how the SCAD file is written.
type SCADOptions struct {
	// OutputPath is the location where the SCAD file should be written.
	OutputPath string
	// GeneratorInclude is an optional override for the include path. When empty,
	// defaults to YAPPgenerator_v3.scad in the repo root.
	GeneratorInclude string
}

// WriteSCAD builds the YAPP model, emits the SCAD, adjusts include paths, and writes the file.
func WriteSCAD(ctx context.Context, resolved map[string]any, opts SCADOptions) (string, error) {
	if opts.OutputPath == "" {
		return "", errors.New("output path is required")
	}
	model, err := yappgen.BuildModel(ctx, resolved)
	if err != nil {
		return "", errors.Wrap(err, "build model")
	}
	scad, err := yappgen.EmitSCAD(ctx, model)
	if err != nil {
		return "", errors.Wrap(err, "emit scad")
	}

	scad = rewriteInclude(scad, opts.OutputPath, opts.GeneratorInclude)

	if err := os.WriteFile(opts.OutputPath, scad, 0o644); err != nil {
		return "", errors.Wrapf(err, "write SCAD %s", opts.OutputPath)
	}
	return opts.OutputPath, nil
}

// STLOptions controls how OpenSCAD rendering runs.
type STLOptions struct {
	BasePath string
	LidPath  string
	OpenSCAD string
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
		if err := renderSingleSTL(ctx, bin, scadPath, opts.BasePath, true, false); err != nil {
			return errors.Wrap(err, "base STL")
		}
	}
	if opts.LidPath != "" {
		if err := renderSingleSTL(ctx, bin, scadPath, opts.LidPath, false, true); err != nil {
			return errors.Wrap(err, "lid STL")
		}
	}
	return nil
}

func renderSingleSTL(ctx context.Context, bin, scadPath, outPath string, printBase, printLid bool) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return errors.Wrap(err, "create STL directory")
	}
	fmt.Printf("Rendering STL %s from %s\n", outPath, scadPath)
	args := []string{
		"-o", outPath,
		"-D", fmt.Sprintf("printBaseShell=%t", printBase),
		"-D", fmt.Sprintf("printLidShell=%t", printLid),
		"-D", "printSwitchExtenders=false",
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
