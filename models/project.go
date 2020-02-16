package models

import (
	"fmt"
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

const (
	metaFlag = "_META_"

	belongsTo    = "belongs_to"
	notCreatable = "!creatable"
	notEditable  = "!editable"
)

// DataType represents a data model
type DataType struct {
	Name            wordsmith.SuperPalabra
	BelongsToUser   bool
	BelongsToNobody bool
	BelongsToStruct wordsmith.SuperPalabra
	Fields          []DataField
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
	EnableNewsman           bool

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
		dts, imps, err := parseModels(p.OutputPath, pkg.Files)
		if err != nil {
			return fmt.Errorf("attempting to read package %s: %w", pkg.Name, err)
		}

		p.DataTypes = append(p.DataTypes, dts...)
		p.iterableServicesImports = append(p.iterableServicesImports, imps...)
	}

	return nil
}

func parseModels(outputPath string, pkgFiles map[string]*ast.File) (dataTypes []DataType, imports []string, returnErr error) {
	for _, file := range pkgFiles {
		ast.Inspect(file, func(n ast.Node) bool {
			if dec, ok := n.(*ast.TypeSpec); ok {
				dt := DataType{
					Name:          wordsmith.FromSingularPascalCase(dec.Name.Name),
					Fields:        []DataField{},
					BelongsToUser: true,
				}

				if _, ok := dec.Type.(*ast.StructType); !ok {
					// log.Println("ERROR: only structs allowed in model declarations")
					return true
				}

				for _, field := range dec.Type.(*ast.StructType).Fields.List {
					fn := field.Names[0].Name
					df := DataField{
						Name:                  wordsmith.FromSingularPascalCase(fn),
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					}

					if x, identOK := field.Type.(*ast.Ident); identOK {
						df.Type = x.Name
					} else if y, starOK := field.Type.(*ast.StarExpr); starOK {
						df.Pointer = true
						df.Type = y.X.(*ast.Ident).Name
					}

					// check if this is the meta flag
					if fn == metaFlag && df.Type == "uintptr" {
						// since found the meta flag, process its directives
						var tag string
						if field != nil && field.Tag != nil {
							tag = field.Tag.Value
						}

						// check belonging
						if strings.Contains(tag, "belongs_to") {
							tagWithoutBackticks := strings.ReplaceAll(tag, "`", "")
							tagWithoutBelongsTo := strings.ReplaceAll(tagWithoutBackticks, fmt.Sprintf("%s:", belongsTo), "")
							ownerWithoutQuotes := strings.ReplaceAll(tagWithoutBelongsTo, `"`, ``)

							if ownerWithoutQuotes == "__nobody__" {
								dt.BelongsToUser = false
								dt.BelongsToNobody = true
							} else if ownerWithoutQuotes != "" {
								dt.BelongsToStruct = wordsmith.FromSingularPascalCase(ownerWithoutQuotes)
							}
						}
					} else {
						var tag string
						if field != nil && field.Tag != nil {
							tag = strings.Replace(
								strings.Replace(
									strings.Replace(field.Tag.Value,
										`naff:`, "", 1,
									),
									"`", "", -1,
								),
								`"`, "", -1,
							)
						}

						for _, t := range strings.Split(tag, ",") {
							switch strings.ToLower(strings.TrimSpace(t)) {
							case notCreatable:
								df.ValidForCreationInput = false
							case notEditable:
								df.ValidForUpdateInput = false
							}
						}

						dt.Fields = append(dt.Fields, df)
					}
				}

				dataTypes = append(dataTypes, dt)

				// BEGIN check for invalid ownership arrangements
				var owners = map[string]string{}
				var names = []string{}
				for _, typ := range dataTypes {
					names = append(names, typ.Name.Singular())
					if typ.BelongsToStruct != nil && typ.BelongsToStruct.Plural() != "" {
						owners[typ.Name.Singular()] = typ.BelongsToStruct.Singular()
					}
				}

				for dt, owner := range owners {
					var ownerFound bool
					for _, name := range names {
						if name == owner {
							ownerFound = true
						}
					}

					if !ownerFound {
						returnErr = fmt.Errorf("invalid ownership arrangement: %s belongs to %q, which does not exist", dt, owner)
						return false
					}
				}
				// END check for invalid ownership arrangements

				imports = append(
					imports,
					filepath.Join(
						outputPath,
						"services",
						"v1",
						strings.ToLower(dt.Name.Plural()),
					),
				)
			}
			return true
		})
	}

	return
}

type projectSurvey struct {
	Name             string `survey:"name"`
	OutputRepository string `survey:"outputRepository"`
	ModelsPackage    string `survey:"modelsPackage"`
	EnableNewsman    bool   `survey:"enableNewsman"`
}

// CompleteSurvey asks the user questions to determine core project information
func CompleteSurvey(projectName, sourceModels, outputPackage string) (*Project, error) {
	// the questions to ask
	questions := []*survey.Question{}
	p := projectSurvey{}

	if projectName == "" {
		questions = append(questions, &survey.Question{
			Name:     "name",
			Prompt:   &survey.Input{Message: "project name:"},
			Validate: survey.Required,
		})
	} else {
		p.Name = projectName
	}

	if sourceModels == "" {
		questions = append(questions, &survey.Question{
			Name: "modelsPackage",
			Prompt: &survey.Input{
				Message: "models package:",
				Help:    `the input package that defines the base set of models (i.e. gitlab.com/verygoodsoftwarenotvirus/naff/example_models/todo)`,
			},
			Validate: survey.Required,
		})
	} else {
		p.ModelsPackage = sourceModels
	}

	if outputPackage == "" {
		questions = append(questions, &survey.Question{
			Name: "outputRepository",
			Prompt: &survey.Input{
				Message: "output repository path:",
				Help:    `the package path that the generated project will live in (i.e. gitlab.com/verygoodsoftwarenotvirus/whatever)`,
			},
			Validate: survey.Required,
		})
	} else {
		p.OutputRepository = outputPackage
	}

	// {
	// 	Name: "enableNewsman",
	// 	Prompt: &survey.Confirm{
	// 		Message: "enable newsman?",
	// 		Default: true,
	// 		Help:    "generates newsman code",
	// 	},
	// 	Validate: survey.Required,
	// },

	// perform the questions
	if surveyErr := survey.Ask(questions, &p); surveyErr != nil {
		return nil, surveyErr
	}

	targetDestination := filepath.Join(os.Getenv("GOPATH"), "src", p.OutputRepository)
	if strings.HasSuffix(targetDestination, "verygoodsoftwarenotvirus") {
		log.Fatal("I don't think you actually want to do that.")
	} else {
		os.RemoveAll(targetDestination)
	}

	return &Project{
		EnableNewsman: true,
		Name:          wordsmith.FromSingularPascalCase(p.Name),
		OutputPath:    p.OutputRepository,
		sourcePackage: p.ModelsPackage,
	}, nil
}
