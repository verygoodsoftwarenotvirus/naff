package models

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"hash/fnv"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/gonum/graph/simple"
	"github.com/gonum/graph/topo"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"

	"gopkg.in/AlecAivazis/survey.v1"
)

const (
	metaFlag = "meta"

	belongsTo    = "belongs_to"
	notCreatable = "!creatable"
	notEditable  = "!editable"
)

type depWrapper struct {
	dependency string
}

func (dw depWrapper) ID() int {
	x := fnv.New32a()

	if _, err := x.Write([]byte(dw.dependency)); err != nil {
		panic(err)
	}

	return int(x.Sum32())
}

type Project struct {
	sourcePackage string

	OutputPath                string
	OutputPathStringToReplace string
	OutputPathSubstitution    string

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
		if p.containsCyclicOwnerships() {
			return errors.New("error: cyclic ownership detected")
		}

		p.iterableServicesImports = append(p.iterableServicesImports, imps...)
	}

	return nil
}

// FindOwnerTypeChain returns the owner chain of a given object from highest ancestor to lowest
// so if C belongs to B belongs to A, then calling `FindOwnerTypeChain` for C would yield [A, B]
func (p *Project) FindOwnerTypeChain(typ DataType) []DataType {
	parentTypes := []DataType{}

	var parentType *DataType = &typ

	for parentType != nil && parentType.BelongsToStruct != nil {
		newParent := p.FindType(parentType.BelongsToStruct.Singular())
		if newParent != nil {
			parentTypes = append(parentTypes, *newParent)
		}
		parentType = newParent
	}

	// reverse it
	for i, j := 0, len(parentTypes)-1; i < j; i, j = i+1, j-1 {
		parentTypes[i], parentTypes[j] = parentTypes[j], parentTypes[i]
	}

	return parentTypes
}

func (p *Project) FindType(name string) *DataType {
	for _, typ := range p.DataTypes {
		if typ.Name.Singular() == name {
			return &typ
		}
	}

	return nil
}

func (p *Project) FindDependentsOfType(parentType DataType) []DataType {
	dependents := []DataType{}
	for _, typ := range p.DataTypes {
		if typ.BelongsToStruct != nil && typ.BelongsToStruct.Singular() == parentType.Name.Singular() {
			dependents = append(dependents, typ)
		}
	}

	return dependents
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

func (p *Project) containsCyclicOwnerships() bool {
	g := simple.NewDirectedGraph(0, math.Inf(1))

	for _, typ := range p.DataTypes {
		w := depWrapper{dependency: typ.Name.Singular()}
		if !g.Has(w) {
			g.AddNode(w)
		}
	}

	for _, typ := range p.DataTypes {
		w := depWrapper{dependency: typ.Name.Singular()}

		if typ.BelongsToStruct != nil {
			g.SetEdge(simple.Edge{F: w, T: depWrapper{typ.BelongsToStruct.Singular()}})
		}
	}

	cycles := topo.CyclesIn(g)
	return len(cycles) != 0
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
