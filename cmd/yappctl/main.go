package main

import (
	"fmt"
	"os"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/help"
	help_cmd "github.com/go-go-golems/glazed/pkg/help/cmd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/wesen/yapp-encl-resolver/pkg/docs"
)

func main() {
	root := &cobra.Command{
		Use:   "yappctl",
		Short: "Unified CLI for the YAPP enclosure DSL",
		Long: `yappctl bundles resolver and generator workflows for the YAPP enclosure DSL.

Use 'yappctl resolve' to evaluate expressions/vars, and 'yappctl generate' to produce SCAD and STLs.`,
	}

	helpSystem := help.NewHelpSystem()
	if err := docs.Load(helpSystem); err != nil {
		fmt.Fprintf(os.Stderr, "error loading CLI help docs: %v\n", err)
	}
	help_cmd.SetupCobraRootCommand(helpSystem, root)

	if err := registerCommands(root); err != nil {
		fmt.Fprintf(os.Stderr, "error registering commands: %v\n", err)
		os.Exit(1)
	}

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func registerCommands(root *cobra.Command) error {
	resolveCmd, err := NewResolveCommand()
	if err != nil {
		return errors.Wrap(err, "resolve command")
	}
	if err := addBareCommand(root, resolveCmd); err != nil {
		return err
	}

	generateCmd, err := NewGenerateCommand()
	if err != nil {
		return errors.Wrap(err, "generate command")
	}
	if err := addBareCommand(root, generateCmd); err != nil {
		return err
	}
	return nil
}

func addBareCommand(root *cobra.Command, command cmds.BareCommand) error {
	cobraCmd, err := cli.BuildCobraCommand(command,
		cli.WithParserConfig(cli.CobraParserConfig{
			ShortHelpLayers: []string{layers.DefaultSlug},
			MiddlewaresFunc: cli.CobraCommandDefaultMiddlewares,
		}),
	)
	if err != nil {
		return err
	}
	root.AddCommand(cobraCmd)
	return nil
}
