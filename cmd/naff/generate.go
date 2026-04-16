package main

import (
	"fmt"

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
