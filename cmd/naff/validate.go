package main

import (
	"fmt"

	"github.com/verygoodsoftwarenotvirus/naff/internal/config"

	"github.com/spf13/cobra"
)

func buildValidateCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a project YAML config without generating",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadFromFile(configPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			if err = cfg.Validate(); err != nil {
				return fmt.Errorf("validating config: %w", err)
			}

			fmt.Println("Config is valid.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "project.yaml", "path to project YAML config")

	return cmd
}
