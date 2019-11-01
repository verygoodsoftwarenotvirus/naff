package models

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"

	"gopkg.in/AlecAivazis/survey.v1"
)

// DataType represents a data model
type DataType struct {
	Name   wordsmith.SuperPalabra
	Fields []DataField
}

// DataField represents a data model's field
type DataField struct {
	Name                  wordsmith.SuperPalabra
	Type                  string
	Pointer               bool
	DefaultAllowed        bool
	DefaultValue          string
	ValidForCreationInput bool
	ValidForUpdateInput   bool
}

type Project struct {
	sourcePackage           string
	OutputPath              string
	iterableServicesImports []string
	EnableWebhooks          bool

	Name      wordsmith.SuperPalabra
	DataTypes []DataType
}

func (p *Project) ParseModels(outputPath string) error {
	fullModelsPath := filepath.Join(os.Getenv("GOPATH"), "src", p.sourcePackage)

	packages, err := parser.ParseDir(token.NewFileSet(), fullModelsPath, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				if dec, ok := n.(*ast.TypeSpec); ok {
					dt := DataType{
						Name:   wordsmith.FromSingularPascalCase(dec.Name.Name),
						Fields: []DataField{},
					}

					if _, ok := dec.Type.(*ast.StructType); !ok {
						log.Println("ERROR: only structs allowed in model declarations")
						return false
					}

					for _, field := range dec.Type.(*ast.StructType).Fields.List {
						df := DataField{
							Name:                  wordsmith.FromSingularPascalCase(field.Names[0].Name),
							ValidForCreationInput: true,
							ValidForUpdateInput:   true,
						}

						if x, identOK := field.Type.(*ast.Ident); identOK {
							df.Type = x.Name
						} else if y, starOK := field.Type.(*ast.StarExpr); starOK {
							df.Pointer = true
							df.Type = y.X.(*ast.Ident).Name
						}

						var tag string
						if field != nil && field.Tag != nil {
							tag = strings.Replace(strings.Replace(
								strings.Replace(field.Tag.Value, `naff:`, "", 1),
								"`", "", -1), `"`, "", -1)
						}

						for _, t := range strings.Split(tag, ",") {
							switch strings.ToLower(strings.TrimSpace(t)) {
							case "!creatable":
								df.ValidForCreationInput = false
							case "!editable":
								df.ValidForUpdateInput = false
							}
						}
						dt.Fields = append(dt.Fields, df)
					}

					p.DataTypes = append(p.DataTypes, dt)
					p.iterableServicesImports = append(
						p.iterableServicesImports,
						filepath.Join(
							p.OutputPath,
							"services",
							"v1",
							strings.ToLower(dt.Name.Plural()),
						),
					)
				}
				return true
			})
		}
	}
	return nil
}

type projectSurvey struct {
	Name             string `survey:"name"`
	OutputRepository string `survey:"outputRepository"`
	ModelsPackage    string `survey:"modelsPackage"`
}

// CompleteSurvey asks the user questions to determine core project information
func CompleteSurvey() (*Project, error) {
	// the questions to ask
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "project name:"},
			Validate: survey.Required,
		},
		{
			Name: "outputRepository",
			Prompt: &survey.Input{
				Message: "output repository path:",
				Default: "gitlab.com/verygoodsoftwarenotvirus/whatever",
				Help:    `the package path that the generated project will live in`,
			},
		},
		{
			Name: "modelsPackage",
			Prompt: &survey.Input{
				Message: "models package:",
				Default: "gitlab.com/verygoodsoftwarenotvirus/naff/example_models/todo",
				Help:    `the input package that defines the base set of models`,
			},
		},
	}

	// perform the questions
	p := projectSurvey{}
	if surveyErr := survey.Ask(questions, &p); surveyErr != nil {
		return nil, surveyErr
	}
	os.RemoveAll(filepath.Join(os.Getenv("GOPATH"), "src", p.OutputRepository))

	return &Project{
		Name:          wordsmith.FromSingularPascalCase(p.Name),
		OutputPath:    p.OutputRepository,
		sourcePackage: p.ModelsPackage,
	}, nil
}
