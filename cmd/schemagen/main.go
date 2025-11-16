package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/wesen/yapp-encl-resolver/pkg/schemagen"
)

func main() {
	rootCmd := newRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "schemagen error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schemagen",
		Short: "Generate Go code from YAPP module schemas",
		Long: `schemagen validates YAML schemas for YAPP modules and generates Go code.

Use 'schemagen validate <schema.yaml>' to lint a single schema, or
'schemagen discover' to scan the modules directory and emit generated files.`,
	}

	cmd.AddCommand(newValidateCommand())
	cmd.AddCommand(newDiscoverCommand())

	return cmd
}

func newValidateCommand() *cobra.Command {
	var (
		schemaPath string
	)

	c := &cobra.Command{
		Use:   "validate <schema.yaml>",
		Short: "Validate a single module schema file",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 && schemaPath == "" {
				return errors.New("provide schema path as argument or --schema flag")
			}
			if len(args) == 1 {
				schemaPath = args[0]
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return runValidate(ctx, schemaPath)
		},
	}

	c.Flags().StringVar(&schemaPath, "schema", "", "Path to schema YAML file")

	return c
}

type discoverOptions struct {
	ModulesDir string
	OutputDir  string
}

func newDiscoverCommand() *cobra.Command {
	opts := discoverOptions{
		ModulesDir: "pkg/yappgen/modules",
	}

	c := &cobra.Command{
		Use:   "discover",
		Short: "Discover schema files and generate code",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return runDiscover(ctx, opts)
		},
	}

	c.Flags().StringVar(&opts.ModulesDir, "modules-dir", opts.ModulesDir, "Directory containing module schema.yaml files")
	c.Flags().StringVar(&opts.OutputDir, "output-dir", "", "Optional directory to write generated files (defaults to in-place)")

	return c
}

func runValidate(ctx context.Context, schemaPath string) error {
	if schemaPath == "" {
		return errors.New("schema path is required")
	}
	if err := schemagen.ValidateSchemaFile(schemaPath); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Schema %s is valid\n", schemaPath)
	return nil
}

func runDiscover(ctx context.Context, opts discoverOptions) error {
	if opts.ModulesDir == "" {
		return errors.New("modules directory is required")
	}
	schemaFiles, err := schemagen.DiscoverSchemaFiles(opts.ModulesDir)
	if err != nil {
		return err
	}
	if len(schemaFiles) == 0 {
		fmt.Fprintf(os.Stdout, "No schema.yaml files found under %s\n", opts.ModulesDir)
		return nil
	}

	fmt.Fprintf(os.Stdout, "Found %d schema file(s):\n", len(schemaFiles))
	for _, path := range schemaFiles {
		fmt.Fprintf(os.Stdout, " - %s\n", path)
	}

	if err := schemagen.ValidateSchemaFiles(schemaFiles); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Validation succeeded for all schemas. Code generation coming soon.\n")
	return schemagen.ErrCodeGenerationNotImplemented
}
