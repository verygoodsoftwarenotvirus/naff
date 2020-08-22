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
	projectName,
	sourceModels,
	outputPackage string

	postgresEnabled,
	sqliteEnabled,
	mariaDBEnabled bool

	// generateCmd represents the generate command
	generateCmd = &cobra.Command{
		Use:   "gen",
		Short: "executes the templates to generate boilerplate",
		Long: `This command will prompt the user for a few things:
	1. The name of the project
	2. The directory where the files should end up
	3. The directory where the input models are kept	
	4. Which databases you care about supporting

INPUT:
	As input, NAFF accepts a package that, when combined with the location of
your GOPATH, will point to a valid directory with Go files in it. These files
can contain lots of things, but the only thing NAFF will actually look for are
struct definitions. Take this for an example:
		type Item struct{
			Name string
			Details string
		}

	NAFF can read that definition, and create all the code necessary to CRUD
that type. It will make sure there is a column in the database for every field,
and it will also automatically create the requisite ID, creation time, last
updated at time, and archived time fields, so you don't need to include those.
	NAFF only tolerates simple types and pointers to simple types, in fact
here's the whole list:
		bool
		string
		float32, float64
		int, int8, int16, int32, int64
		uint, uint8, uint16, uint32, uint64

	If there are any other types in your model definitions, you're going to
have a bad time. You might be looking at this list, and be concerned that there
aren't any time.Time fields, and you may be applauded for your astuteness. Time
is represented in NAFF codebases as unsigned 64-bit integers, or pointers where
relevant. If you're looking to replicate this behavior, I feel bad for you,
son, because I didn't implement that and I only thought of it while writing
these docs.
	If your types are pointers, the only real differences will be in the
database schema (where they won't have a NOT NULL clause), and in some other
areas (tests, update functions) there will be asterisk operators or nil checks.
	There are two currently supported struct field tags for NAFF:
		!creatable
		!editable

	If you wanted to have our above structure, only you wanted users to be
incapable of setting the Name, and incapable of editing the Details, you could
provide this:
		type Item struct{
			Name *string ` + "`" + `naff:"!creatable"` + "`" + `
			Details string ` + "`" + `naff:"!editable"` + "`" + `
		}

OUTPUT:
	NAFF will only ever create files that don't exist, it will never overwrite
files. If it's not creating a file you think should be there, make sure it
didn't barf any errors out, and make sure the file does not exist.
	Once NAFF does its thing, you will likely want to enter into the directory
and run:
		` + "`" + `make gamut` + "`" + `

	That will set your dependencies up, and then you should be good to code.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := models.CompleteSurvey(projectName, sourceModels, outputPackage)
			if err != nil {
				return err
			}

			if strings.TrimSpace(p.OutputPath) == "gitlab.com/verygoodsoftwarenotvirus/naff" {
				return errors.New("you want me to erase myself?")
			}

			if err := p.ParseModels(); err != nil {
				return err
			}

			project.RenderProject(p)
			return nil
		},
	}
)

func init() {
	generateCmd.Flags().StringVarP(&projectName, "name", "n", "", "project name")
	generateCmd.Flags().StringVarP(&sourceModels, "source-models", "m", "", "input models package")
	generateCmd.Flags().StringVarP(&outputPackage, "output-dir", "o", "", "package to generate")

	generateCmd.Flags().BoolVarP(&postgresEnabled, "enable-postgres", "", false, "enable postgres support")
	generateCmd.Flags().BoolVarP(&sqliteEnabled, "enable-sqlite", "", false, "enable sqlite support")
	generateCmd.Flags().BoolVarP(&mariaDBEnabled, "enable-mariadb", "", false, "enable mariadb support")

	rootCmd.AddCommand(generateCmd)
}
