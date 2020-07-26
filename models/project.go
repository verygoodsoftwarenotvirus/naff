package models

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"hash/fnv"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"

	"github.com/gonum/graph/simple"
	"github.com/gonum/graph/topo"
	"gopkg.in/AlecAivazis/survey.v1"
)

const (
	metaFieldName = "_META_"

	belongsTo    = "belongs_to"
	restricted   = "restricted_to_user"
	notCreatable = "!creatable"
	notEditable  = "!editable"

	Postgres validDatabase = "postgres"
	MariaDB  validDatabase = "mariadb"
	Sqlite   validDatabase = "sqlite"
)

type depWrapper struct {
	dependency string
}

type validDatabase string

var (
	validDatabaseMap = map[validDatabase]struct{}{
		Postgres: {},
		Sqlite:   {},
		MariaDB:  {},
	}

	nameToValidDBMap = map[string]validDatabase{
		string(Postgres): Postgres,
		string(Sqlite):   Sqlite,
		string(MariaDB):  MariaDB,
	}
)

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

	enabledDatabases map[validDatabase]struct{}
}

func (p *Project) Validate() {
	if len(p.EnabledDatabases()) == 0 {
		log.Panic("no databases enabled!")
	}

	if len(p.DataTypes) == 0 {
		log.Panic("no types defined!")
	}

	for _, dt := range p.DataTypes {
		if len(dt.Fields) == 0 {
			log.Panicf("no fields defined for type %q!", dt.Name.Singular())
		}
	}

	if p.containsCyclicOwnerships() {
		log.Panic("error: cyclic ownership detected")
	}
}

func (p *Project) ParseModels() error {
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
		p.Validate() // trust but verify
		p.iterableServicesImports = append(p.iterableServicesImports, imps...)
	}

	return nil
}

// SearchEnabled returns true if any of the datatypes have SearchEnabled
func (p *Project) SearchEnabled() bool {
	for _, typ := range p.DataTypes {
		if typ.SearchEnabled {
			return true
		}
	}

	return false
}

// FindOwnerTypeChain returns the owner chain of a given object from highest ancestor to lowest
// so if C belongs to B belongs to A, then calling `FindOwnerTypeChain` for C would yield [A, B]
func (p *Project) FindOwnerTypeChain(typ DataType) []DataType {
	parentTypes := []DataType{}

	var parentType = &typ

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

// FindOwnerTypeChain returns the owner chain of a given object from highest ancestor to lowest
// so if C belongs to B belongs to A, then calling `FindOwnerTypeChain` for C would yield [A, B]
func (p *Project) FindOwnerTypeChainWithoutReversing(typ DataType) []DataType {
	parentTypes := []DataType{}

	var parentType = &typ

	for parentType != nil && parentType.BelongsToStruct != nil {
		newParent := p.FindType(parentType.BelongsToStruct.Singular())
		if newParent != nil {
			parentTypes = append(parentTypes, *newParent)
		}
		parentType = newParent
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

func (p *Project) ensureNoNilFields() {
	if p.enabledDatabases == nil {
		p.enabledDatabases = map[validDatabase]struct{}{}
	}
}

func (p *Project) EnableDatabase(database validDatabase) {
	p.ensureNoNilFields()

	if _, ok := validDatabaseMap[database]; !ok {
		log.Panicf("unknown database: %q", database)
	}

	if _, ok := p.enabledDatabases[database]; !ok {
		p.enabledDatabases[database] = struct{}{}
	}
}

func (p *Project) DisableDatabase(database validDatabase) {
	p.ensureNoNilFields()

	if _, ok := validDatabaseMap[database]; !ok {
		log.Panicf("unknown database: %q", database)
	}

	if _, ok := p.enabledDatabases[database]; ok {
		delete(p.enabledDatabases, database)
	}
}

func (p *Project) DatabaseIsEnabled(database validDatabase) bool {
	p.ensureNoNilFields()

	if _, ok := validDatabaseMap[database]; !ok {
		log.Panicf("unknown database: %q", database)
	}

	_, present := p.enabledDatabases[database]

	return present
}

func (p *Project) EnabledDatabases() []string {
	p.ensureNoNilFields()

	out := []string{}
	for k := range p.enabledDatabases {
		out = append(out, string(k))
	}
	return out
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
					fieldName := field.Names[0].Name
					df := DataField{
						Name:                  wordsmith.FromSingularPascalCase(fieldName),
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
						Pos:                   field.Pos(),
					}

					if x, identOK := field.Type.(*ast.Ident); identOK {
						df.Type = x.Name
					} else if y, starOK := field.Type.(*ast.StarExpr); starOK {
						df.Pointer = true
						df.Type = y.X.(*ast.Ident).Name
					}

					// this line doubles as a validation against invalid types
					df.UnderlyingType = GetTypeForTypeName(df.Type)

					if fieldName != metaFieldName && df.Type == "uintptr" {
						panic("invalid type!")
					}

					// check if this is the meta flag
					if fieldName == metaFieldName && df.Type == "uintptr" {
						// since found the meta flag, process its directives
						var tags string
						if field != nil && field.Tag != nil {
							tags = field.Tag.Value
						}

						for _, tag := range strings.Split(tags, " ") {
							// check belonging
							var alsoBelongsToUser bool
							tagWithoutBackticks := strings.ReplaceAll(tag, "`", "")

							if strings.Contains(tag, belongsTo) {
								tagWithoutBelongsTo := strings.ReplaceAll(tagWithoutBackticks, fmt.Sprintf("%s:", belongsTo), "")
								properOwner := strings.ReplaceAll(tagWithoutBelongsTo, `"`, ``)

								if strings.Contains(properOwner, ",") {
									properOwnerParts := strings.Split(properOwner, ",")

									if len(properOwnerParts) != 2 {
										panic("too many owners, a type may only be owned by another type and a user!")
									}

									if strings.ToLower(properOwnerParts[0]) != "user" && strings.ToLower(properOwnerParts[1]) != "user" {
										panic("too many owners, a type may only be owned by another type and a user!")
									}
									// we can't get here without this being true
									alsoBelongsToUser = true

									if strings.ToLower(properOwnerParts[0]) == "user" {
										properOwner = properOwnerParts[1]
									} else if strings.ToLower(properOwnerParts[1]) == "user" {
										properOwner = properOwnerParts[0]
									}
								}

								if properOwner == "__nobody__" {
									dt.BelongsToUser = false
									dt.BelongsToNobody = true
								} else if properOwner != "" {
									dt.BelongsToStruct = wordsmith.FromSingularPascalCase(properOwner)
									dt.BelongsToUser = alsoBelongsToUser
								}
							} else {
								if strings.Contains(tagWithoutBackticks, `restricted_to_user:"true"`) {
									dt.RestrictedToUser = true
								}
								if strings.Contains(tagWithoutBackticks, `search_enabled:"true"`) && dt.BelongsToUser && dt.BelongsToStruct == nil {
									dt.SearchEnabled = true
								}
								if dt.SearchEnabled {
									containsString := false
									for _, field := range dt.Fields {
										if field.Type == "string" {
											containsString = true
											break
										}
									}
									if !containsString {
										log.Panicf("no string fields present in type with search enabled: %q", dt.Name.Singular())
									}
								}
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

// GetTypeForTypeName blah
func GetTypeForTypeName(name string) types.Type {
	switch strings.ToLower(name) {
	case "bool":
		return types.Typ[types.Bool]
	case "int":
		return types.Typ[types.Int]
	case "int8":
		return types.Typ[types.Int8]
	case "int16":
		return types.Typ[types.Int16]
	case "int32":
		return types.Typ[types.Int32]
	case "int64":
		return types.Typ[types.Int64]
	case "uint":
		return types.Typ[types.Uint]
	case "uint8":
		return types.Typ[types.Uint8]
	case "uint16":
		return types.Typ[types.Uint16]
	case "uint32":
		return types.Typ[types.Uint32]
	case "uint64":
		return types.Typ[types.Uint64]
	case "uintptr":
		return types.Typ[types.Uintptr]
	case "float32":
		return types.Typ[types.Float32]
	case "float64":
		return types.Typ[types.Float64]
	case "string":
		return types.Typ[types.String]
	default:
		log.Panicf("invalid type: %q", name)
		return nil
	}
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

	supportedDBs := []string{}
	dbSupportPrompt := &survey.MultiSelect{
		Message: "Which databases would you like to generate code for?",
		Options: []string{string(Postgres), string(Sqlite), string(MariaDB)},
	}

	if dbPromptErr := survey.AskOne(dbSupportPrompt, &supportedDBs, nil); dbPromptErr != nil {
		return nil, dbPromptErr
	}

	targetDestination := filepath.Join(os.Getenv("GOPATH"), "src", p.OutputRepository)
	if strings.HasSuffix(targetDestination, "verygoodsoftwarenotvirus") {
		log.Fatal("I don't think you actually want to do that.")
	} else {
		os.RemoveAll(targetDestination)
	}

	proj := &Project{
		EnableNewsman: true,
		Name:          wordsmith.FromSingularPascalCase(p.Name),
		OutputPath:    p.OutputRepository,
		sourcePackage: p.ModelsPackage,
	}
	for _, db := range supportedDBs {
		proj.EnableDatabase(nameToValidDBMap[db])
	}

	return proj, nil
}
