package main

import (
	"context"
	"time"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/pkg/errors"

	"github.com/wesen/yapp-encl-resolver/pkg/cli/generatorcli"
	"github.com/wesen/yapp-encl-resolver/pkg/cli/resolvercli"
)

// Ensure interface compliance
var _ cmds.BareCommand = &GenerateCommand{}

type GenerateSettings struct {
	Input         string `glazed.parameter:"input"`
	SCADOut       string `glazed.parameter:"scad-out"`
	MaxIterations int    `glazed.parameter:"max-iterations"`
	Strict        bool   `glazed.parameter:"strict"`
	BaseSTL       string `glazed.parameter:"stl-base"`
	LidSTL        string `glazed.parameter:"stl-lid"`
	AllSTL        string `glazed.parameter:"stl-all"`
	OpenSCADBin   string `glazed.parameter:"openscad-bin"`
	RenderTimeout string `glazed.parameter:"render-timeout"`
	CopyGenerator bool   `glazed.parameter:"copy-generator"`
	GeneratorPath string `glazed.parameter:"generator-path"`
	QualityValue  int    `glazed.parameter:"quality-value"`
}

type GenerateCommand struct {
	*cmds.CommandDescription
}

func NewGenerateCommand() (*GenerateCommand, error) {
	commandSettingsLayer, err := cli.NewCommandSettingsLayer()
	if err != nil {
		return nil, err
	}

	desc := cmds.NewCommandDescription(
		"generate",
		cmds.WithShort("Generate SCAD (and optional STLs) from the enclosure DSL"),
		cmds.WithLong(`
Reads the enclosure DSL, resolves expressions, emits a SCAD file compatible with the YAPP generator,
and optionally invokes OpenSCAD to render base/lid STL files.

Examples:
  yappctl generate --input enclosure.yaml --scad-out output.scad
  yappctl generate -i enclosure.yaml -o output.scad --stl-base base.stl --stl-lid lid.stl
`),
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"input",
				parameters.ParameterTypeString,
				parameters.WithRequired(true),
				parameters.WithShortFlag("i"),
				parameters.WithHelp("Path to input DSL YAML"),
			),
			parameters.NewParameterDefinition(
				"scad-out",
				parameters.ParameterTypeString,
				parameters.WithRequired(true),
				parameters.WithShortFlag("o"),
				parameters.WithHelp("Path to write the generated SCAD file"),
			),
			parameters.NewParameterDefinition(
				"max-iterations",
				parameters.ParameterTypeInteger,
				parameters.WithDefault(16),
				parameters.WithHelp("Maximum resolver passes"),
			),
			parameters.NewParameterDefinition(
				"strict",
				parameters.ParameterTypeBool,
				parameters.WithDefault(false),
				parameters.WithHelp("Enable strict resolver validation"),
			),
			parameters.NewParameterDefinition(
				"stl-base",
				parameters.ParameterTypeString,
				parameters.WithDefault(""),
				parameters.WithHelp("Optional path for the base STL output"),
			),
			parameters.NewParameterDefinition(
				"stl-lid",
				parameters.ParameterTypeString,
				parameters.WithDefault(""),
				parameters.WithHelp("Optional path for the lid STL output"),
			),
			parameters.NewParameterDefinition(
				"stl-all",
				parameters.ParameterTypeString,
				parameters.WithDefault(""),
				parameters.WithHelp("Optional path for a single STL containing base+lid (+extenders if present)"),
			),
			parameters.NewParameterDefinition(
				"openscad-bin",
				parameters.ParameterTypeString,
				parameters.WithDefault("openscad"),
				parameters.WithHelp("OpenSCAD binary to use for STL rendering"),
			),
			parameters.NewParameterDefinition(
				"render-timeout",
				parameters.ParameterTypeString,
				parameters.WithDefault(""),
				parameters.WithHelp("Duration for STL rendering context (e.g., 30s, 2m). Disabled by default."),
			),
			parameters.NewParameterDefinition(
				"copy-generator",
				parameters.ParameterTypeBool,
				parameters.WithDefault(false),
				parameters.WithHelp("Copy YAPPgenerator_v3.scad to output directory for standalone SCAD files"),
			),
			parameters.NewParameterDefinition(
				"generator-path",
				parameters.ParameterTypeString,
				parameters.WithDefault(""),
				parameters.WithHelp("Optional path to YAPPgenerator_v3.scad (uses embedded version if not specified)"),
			),
			parameters.NewParameterDefinition(
				"quality-value",
				parameters.ParameterTypeInteger,
				parameters.WithDefault(0),
				parameters.WithHelp("Override the YAPP generator renderQuality (1-32). Zero keeps generator defaults"),
			),
		),
		cmds.WithLayersList(commandSettingsLayer),
	)

	return &GenerateCommand{CommandDescription: desc}, nil
}

func (c *GenerateCommand) Run(ctx context.Context, parsed *layers.ParsedLayers) error {
	settings := &GenerateSettings{}
	if err := parsed.InitializeStruct(layers.DefaultSlug, settings); err != nil {
		return errors.Wrap(err, "parse parameters")
	}

	resolved, err := resolvercli.LoadAndResolve(ctx, settings.Input, resolvercli.LoadOptions{
		MaxIterations: settings.MaxIterations,
		Strict:        settings.Strict,
	})
	if err != nil {
		return err
	}

	// Auto-enable generator copying if rendering STLs (OpenSCAD needs the generator file)
	copyGenerator := settings.CopyGenerator || settings.BaseSTL != "" || settings.LidSTL != "" || settings.AllSTL != ""

	scadPath, model, err := generatorcli.WriteSCAD(ctx, resolved, generatorcli.SCADOptions{
		OutputPath:    settings.SCADOut,
		CopyGenerator: copyGenerator,
		GeneratorPath: settings.GeneratorPath,
	})
	if err != nil {
		return err
	}

	if settings.BaseSTL == "" && settings.LidSTL == "" && settings.AllSTL == "" {
		return nil
	}

	renderCtx := ctx
	var cancel context.CancelFunc
	if settings.RenderTimeout != "" {
		duration, err := time.ParseDuration(settings.RenderTimeout)
		if err != nil {
			return errors.Wrap(err, "parse render-timeout")
		}
		if duration > 0 {
			renderCtx, cancel = context.WithTimeout(ctx, duration)
			defer cancel()
		}
	}

	return generatorcli.RenderSTLs(renderCtx, scadPath, generatorcli.STLOptions{
		BasePath:             settings.BaseSTL,
		LidPath:              settings.LidSTL,
		AllPath:              settings.AllSTL,
		OpenSCAD:             settings.OpenSCADBin,
		PrintSwitchExtenders: model != nil && model.PrintSwitchExtenders,
		QualityValue:         settings.QualityValue,
	})
}
