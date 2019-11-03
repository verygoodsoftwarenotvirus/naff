// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	project "gitlab.com/verygoodsoftwarenotvirus/naff/templates"

	"github.com/spf13/cobra"
)

var (
	outputPackage string

	// generateCmd represents the generate command
	generateCmd = &cobra.Command{
		Use:   "gen",
		Short: "executes the templates to generate boilerplate",
		Long: `This command will prompt the user for a few things:
	1. The name of the project
	2. The directory where the files should end up
	3. The directory where the input models are kept

Input models are probably not necessary, but they may as well be, if you try to use this tool without any, you're going to have a bad time.
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := models.CompleteSurvey()
			if err != nil {
				return err
			}

			if strings.HasPrefix(strings.TrimSpace(p.OutputPath), "gitlab.com/verygoodsoftwarenotvirus/naff") {
				return errors.New("you want me to erase myself?")
			}

			if err := p.ParseModels(p.OutputPath); err != nil {
				return err
			}

			if err := project.RenderProject(p); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	generateCmd.Flags().StringVarP(&outputPackage, "output-package", "o", "", "Package to generate.")

	rootCmd.AddCommand(generateCmd)
}
