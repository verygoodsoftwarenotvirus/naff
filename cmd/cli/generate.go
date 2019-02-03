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

	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

const (
	defaultFileExtension = ".tmpl"

	defaultTemplatePath = "test_template"
)

func fillSurvey() (*Project, error) {
	// the questions to ask
	questions := []*survey.Question{
		{
			Name:      "name",
			Prompt:    &survey.Input{Message: "Project name:"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name: "basePackage",
			Prompt: &survey.Input{
				Message: "Base import:",
				Help: `the package path that all the subrepositories will live in.
Something like gitlab.com/verygoodsoftwarenotvirus`,
			},
		},
	}

	// perform the questions
	var p Project
	return &p, survey.Ask(questions, &p)
}

var (
	outputPackage string
	// generateCmd represents the generate command
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputPackage == "" {
				return errors.New("no output package set, please run again with --output-package set")
			}

			if fileBox == nil {
				return errors.New("no filebox available")
			}

			p, err := fillSurvey()
			if err != nil {
				return err
			}
			if err := p.EnsureRootPath(); err != nil {
				return err
			}
			if err := p.RenderDirectory(fileBox); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	fileBox = packr.New("templates", defaultTemplatePath)

	generateCmd.Flags().StringVarP(&outputPackage, "output-package", "o", "", "Package to generate.")

	rootCmd.AddCommand(generateCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
