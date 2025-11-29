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
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/wesen/yapp-encl-resolver/pkg/cli/resolvercli"
	"github.com/wesen/yapp-encl-resolver/pkg/resolver/errorx"
	"github.com/wesen/yapp-encl-resolver/pkg/resolver/rules"
)

// Ensure interface compliance
var _ cmds.BareCommand = &ResolveCommand{}

type ResolveSettings struct {
	Input         string `glazed.parameter:"input"`
	OutFile       string `glazed.parameter:"out-file"`
	Format        string `glazed.parameter:"format"`
	MaxIterations int    `glazed.parameter:"max-iterations"`
	Strict        bool   `glazed.parameter:"strict"`
	ShowTaxonomy  bool   `glazed.parameter:"show-taxonomy"`
}

type ResolveCommand struct {
	*cmds.CommandDescription
}

func NewResolveCommand() (*ResolveCommand, error) {
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
  yappctl resolve --input enclosure.yaml
  yappctl resolve --input enclosure.yaml --out-file resolved.yaml
  yappctl resolve -i enclosure.yaml -o resolved.json --format json
  yappctl resolve -i enclosure.yaml --max-iterations 12 --strict
`),
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"input",
				parameters.ParameterTypeString,
				parameters.WithHelp("Path to input DSL YAML (not read into memory by the framework)"),
				parameters.WithRequired(true),
				parameters.WithShortFlag("i"),
			),
			parameters.NewParameterDefinition(
				"out-file",
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
			parameters.NewParameterDefinition(
				"show-taxonomy",
				parameters.ParameterTypeBool,
				parameters.WithDefault(false),
				parameters.WithHelp("Show raw taxonomy structure instead of executing rules"),
			),
		),
		cmds.WithLayersList(commandSettingsLayer),
	)

	return &ResolveCommand{CommandDescription: desc}, nil
}

func (c *ResolveCommand) Run(ctx context.Context, parsed *layers.ParsedLayers) error {
	settings := &ResolveSettings{}
	if err := parsed.InitializeStruct(layers.DefaultSlug, settings); err != nil {
		return errors.Wrap(err, "parse parameters")
	}

	resolved, err := resolvercli.LoadAndResolve(ctx, settings.Input, resolvercli.LoadOptions{
		MaxIterations: settings.MaxIterations,
		Strict:        settings.Strict,
	})
	if err != nil {
		// Try to extract taxonomy from error
		if taxonomy, ok := errorx.AsTaxonomy(err); ok {
			if settings.ShowTaxonomy {
				// Show raw taxonomy structure
				if settings.Format == "json" {
					jsonStr, jsonErr := errorx.FormatTaxonomyJSON(taxonomy)
					if jsonErr != nil {
						return errors.Wrap(jsonErr, "format taxonomy as JSON")
					}
					fmt.Println(jsonStr)
				} else {
					fmt.Println(errorx.FormatTaxonomy(taxonomy))
				}
				return nil
			}
			// Execute rules and show help
			reg := rules.DefaultRegistry()
			results, ruleErr := reg.RenderAll(ctx, taxonomy)
			if ruleErr == nil && len(results) > 0 {
				fmt.Fprintf(os.Stderr, "\n%s\n", rules.RenderToText(results))
			}
		}
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

	if settings.OutFile == "" {
		fmt.Printf("%s", string(out))
		return nil
	}
	if err := os.WriteFile(settings.OutFile, out, 0o644); err != nil {
		return errors.Wrapf(err, "write output file %s", settings.OutFile)
	}
	return nil
}
