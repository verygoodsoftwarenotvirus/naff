package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/verygoodsoftwarenotvirus/naff/internal/config"
	"github.com/verygoodsoftwarenotvirus/naff/internal/pipeline"
	"github.com/verygoodsoftwarenotvirus/naff/templates"

	"github.com/spf13/cobra"
)

func buildGenerateCmd() *cobra.Command {
	var (
		configPath string
		outputDir  string
		clean      bool
		debug      bool
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a project from a YAML config",
		RunE: func(cmd *cobra.Command, args []string) error {
			// If --config was not explicitly provided, check for .naff.yaml in the output directory.
			if !cmd.Flags().Changed("config") {
				naffConfigPath := filepath.Join(outputDir, config.NaffConfigFileName)
				if _, err := os.Stat(naffConfigPath); err == nil {
					configPath = naffConfigPath
					fmt.Printf("Using %s as config source\n", naffConfigPath)
				}
			}

			cfg, err := config.LoadFromFile(configPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			if err = cfg.Validate(); err != nil {
				return fmt.Errorf("validating config: %w", err)
			}

			pipeline.TemplateFS = templates.FS
			p := pipeline.New(cfg, outputDir, clean, debug)
			return p.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "project.yaml", "path to project YAML config")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "./output", "output directory")
	cmd.Flags().BoolVar(&clean, "clean", false, "remove orphaned generated files")
	cmd.Flags().BoolVar(&debug, "debug", false, "stream post-generation make output live")

	return cmd
}
