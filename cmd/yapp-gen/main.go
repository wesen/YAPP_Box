package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver"
	"github.com/wesen/yapp-encl-resolver/pkg/yappgen"
)

func main() {
	root := &cobra.Command{
		Use:   "yapp-gen",
		Short: "Generate a YAPP OpenSCAD file from the enclosure YAML DSL (MVP)",
		RunE: func(cmd *cobra.Command, args []string) error {
			inPath, _ := cmd.Flags().GetString("in")
			outPath, _ := cmd.Flags().GetString("out")
			strict, _ := cmd.Flags().GetBool("strict")
			if inPath == "" || outPath == "" {
				return errors.Errorf("--in and --out are required")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			return run(ctx, inPath, outPath, strict)
		},
	}
	root.Flags().String("in", "", "Input YAML DSL document")
	root.Flags().String("out", "", "Output SCAD file path")
	root.Flags().Bool("strict", false, "Enable strict resolver checks")

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, inPath, outPath string, strict bool) error {
	raw, err := os.ReadFile(inPath)
	if err != nil {
		return errors.Wrap(err, "read input")
	}
	var doc map[string]any
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return errors.Wrap(err, "parse yaml")
	}
	// Resolve expressions to numeric values
	resolved, err := resolver.Resolve(ctx, doc, resolver.Options{
		MaxIterations: 16,
		Strict:        strict,
	})
	if err != nil {
		return errors.Wrap(err, "resolve expressions")
	}
	// Build model and emit SCAD
	model, err := yappgen.BuildModel(ctx, resolved)
	if err != nil {
		return errors.Wrap(err, "build model")
	}
	out, err := yappgen.EmitSCAD(ctx, model)
	if err != nil {
		return errors.Wrap(err, "emit scad")
	}
	// Adjust include path to YAPPgenerator_v3.scad relative to the output directory.
	// Assumes this command is executed from the repo root or a cwd where the file exists.
	outDir := filepath.Dir(outPath)
	genAbs, err := filepath.Abs("YAPPgenerator_v3.scad")
	if err == nil {
		outDirAbs, _ := filepath.Abs(outDir)
		rel, rErr := filepath.Rel(outDirAbs, genAbs)
		if rErr == nil {
			rel = strings.ReplaceAll(rel, string(filepath.Separator), "/")
			updated := strings.Replace(string(out), "include <./YAPPgenerator_v3.scad>", fmt.Sprintf("include <%s>", rel), 1)
			out = []byte(updated)
		}
	}
	if err := os.WriteFile(outPath, out, 0o644); err != nil {
		return errors.Wrap(err, "write output")
	}
	return nil
}


