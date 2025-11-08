package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/help"
	help_cmd "github.com/go-go-golems/glazed/pkg/help/cmd"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver"
)

// Ensure interface compliance
var _ cmds.BareCommand = &ResolveCommand{}

type ResolveSettings struct {
	Input         string `glazed.parameter:"input"`
	Output        string `glazed.parameter:"output"`
	Format        string `glazed.parameter:"format"`
	MaxIterations int    `glazed.parameter:"max-iterations"`
	Strict        bool   `glazed.parameter:"strict"`
}

type ResolveCommand struct {
	*cmds.CommandDescription
}

func NewResolveCommand() (*ResolveCommand, error) {
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}
	commandSettingsLayer, err := cli.NewCommandSettingsLayer()
	if err != nil {
		return nil, err
	}

	desc := cmds.NewCommandDescription(
		"resolve",
		cmds.WithShort("Resolve Enclosure DSL YAML (expressions, vars) to concrete values"),
		cmds.WithLong(`
Reads an input YAML describing an enclosure using the DSL, evaluates variables and expressions
to a fixed point, and outputs a fully-resolved configuration as YAML or JSON.

Examples:
  encl-resolve resolve --input enclosure.yaml
  encl-resolve resolve --input enclosure.yaml --output resolved.yaml
  encl-resolve resolve -i enclosure.yaml -o resolved.json --format json
  encl-resolve resolve -i enclosure.yaml --max-iterations 12 --strict
`),
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"input",
				parameters.ParameterTypeFile,
				parameters.WithHelp("Path to input DSL YAML"),
				parameters.WithRequired(true),
				parameters.WithShortFlag("i"),
			),
			parameters.NewParameterDefinition(
				"output",
				parameters.ParameterTypeString,
				parameters.WithDefault(""),
				parameters.WithHelp("Path to write resolved output (defaults to stdout)"),
				parameters.WithShortFlag("o"),
			),
			parameters.NewParameterDefinition(
				"format",
				parameters.ParameterTypeChoice,
				parameters.WithChoices("yaml", "json"),
				parameters.WithDefault("yaml"),
				parameters.WithHelp("Output format"),
				parameters.WithShortFlag("f"),
			),
			parameters.NewParameterDefinition(
				"max-iterations",
				parameters.ParameterTypeInteger,
				parameters.WithDefault(16),
				parameters.WithHelp("Maximum fixed-point evaluation passes"),
			),
			parameters.NewParameterDefinition(
				"strict",
				parameters.ParameterTypeBool,
				parameters.WithDefault(false),
				parameters.WithHelp("Enable strict validation (unknown keys, unused vars)"),
			),
		),
		cmds.WithLayersList(glazedLayer, commandSettingsLayer),
	)

	return &ResolveCommand{CommandDescription: desc}, nil
}

func (c *ResolveCommand) Run(ctx context.Context, parsed *layers.ParsedLayers) error {
	settings := &ResolveSettings{}
	if err := parsed.InitializeStruct(layers.DefaultSlug, settings); err != nil {
		return errors.Wrap(err, "parse parameters")
	}

	inBytes, err := os.ReadFile(settings.Input)
	if err != nil {
		return errors.Wrapf(err, "read input file %s", settings.Input)
	}
	var doc map[string]any
	if err := yaml.Unmarshal(inBytes, &doc); err != nil {
		return errors.Wrap(err, "parse YAML")
	}

	resolved, err := resolver.Resolve(ctx, doc, resolver.Options{
		MaxIterations: settings.MaxIterations,
		Strict:        settings.Strict,
	})
	if err != nil {
		return err
	}

	var out []byte
	switch settings.Format {
	case "yaml":
		out, err = yaml.Marshal(resolved)
		if err != nil {
			return errors.Wrap(err, "marshal YAML")
		}
	case "json":
		out, err = json.MarshalIndent(resolved, "", "  ")
		if err != nil {
			return errors.Wrap(err, "marshal JSON")
		}
	default:
		return errors.Errorf("unsupported format: %s", settings.Format)
	}

	if settings.Output == "" {
		fmt.Printf("%s", string(out))
		return nil
	}
	if err := os.WriteFile(settings.Output, out, 0o644); err != nil {
		return errors.Wrapf(err, "write output file %s", settings.Output)
	}
	return nil
}

func main() {
	root := &cobra.Command{
		Use:   "encl-resolve",
		Short: "Enclosure DSL resolver",
		Long:  "Resolve Enclosure DSL YAML (variables and expressions) into concrete numeric values",
	}

	cmd, err := NewResolveCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating command: %v\n", err)
		os.Exit(1)
	}
	cobraCmd, err := cli.BuildCobraCommand(cmd,
		cli.WithParserConfig(cli.CobraParserConfig{
			ShortHelpLayers: []string{layers.DefaultSlug},
			MiddlewaresFunc: cli.CobraCommandDefaultMiddlewares,
		}),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error building cobra command: %v\n", err)
		os.Exit(1)
	}
	root.AddCommand(cobraCmd)

	// Enhanced help
	helpSystem := help.NewHelpSystem()
	help_cmd.SetupCobraRootCommand(helpSystem, root)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}


