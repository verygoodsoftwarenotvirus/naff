package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const starterConfig = `project:
  name: "MyProject"
  module: "github.com/example/myproject"
  ios_bundle_id: "com.example.myproject"
  ios_module_name: "MyProject"

targets:
  backend: true
  ios: true

features:
  webhooks: true
  issuereports: true
  comments: true
  notifications: true
  payments: false
  waitlists: false
  uploadedmedia: true
  dataprivacy: true
  consumer_app: false
  admin_app: false

domains:
  - name: "example"
    entities:
      - name: "Widget"
        belongs_to_account: true
        created_by_user: true
        consumer_editable: true
        fields:
          - name: "Name"
            type: "string"
          - name: "Description"
            type: "string"
            required: false
`

func buildInitCmd() *cobra.Command {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Generate a starter project.yaml config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(outputPath); err == nil {
				return fmt.Errorf("%s already exists, refusing to overwrite", outputPath)
			}

			if err := os.WriteFile(outputPath, []byte(starterConfig), 0o644); err != nil {
				return fmt.Errorf("writing config: %w", err)
			}

			fmt.Printf("Wrote starter config to %s\n", outputPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "project.yaml", "output path for the config file")

	return cmd
}
