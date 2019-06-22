package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/codemodus/kace"
	"gopkg.in/AlecAivazis/survey.v1"
)

type dataType struct {
	Name              string
	PluralLowercase   string
	SingularLowercase string
	PluralTitlecase   string
	SingularTitlecase string
	Fields            []dataField
}

type dataField struct {
	Name                  string
	Type                  string
	Pointer               bool
	ValidForCreationInput bool
	ValidForUpdateInput   bool
}

// Project is our project type
type Project struct {
	Name                    string `survey:"name"`
	DataTypes               []dataType
	IterableServicesImports []string
	OutputRepository        string `survey:"outputRepository"`
	ModelsPackage           string `survey:"modelsPackage"`
}

func fillSurvey() (*Project, error) {
	// the questions to ask
	questions := []*survey.Question{
		{
			Name:      "name",
			Prompt:    &survey.Input{Message: "project name:"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name: "outputRepository",
			Prompt: &survey.Input{
				Message: "output repository path:",
				Default: "gitlab.com/verygoodsoftwarenotvirus/whateverfarts",
				Help: `the package path that all the subrepositories will live in.
Something like gitlab.com/verygoodsoftwarenotvirus`,
			},
		},
		{
			Name: "modelsPackage",
			Prompt: &survey.Input{
				Message: "models package:",
				Default: "gitlab.com/verygoodsoftwarenotvirus/naff/example_models/a",
				Help: `the package path that all the subrepositories will live in.
Something like gitlab.com/verygoodsoftwarenotvirus`,
			},
		},
	}
	_ = questions

	// perform the questions
	p := Project{}
	if surveyErr := survey.Ask(questions, &p); surveyErr != nil {
		return nil, surveyErr
	}
	p.parseModels()

	return &p, nil
}

func (p *Project) parseModels() {
	fullModelsPath := filepath.Join(os.Getenv("GOPATH"), "src", p.ModelsPackage)

	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, fullModelsPath, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range packages {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				if dec, ok := n.(*ast.TypeSpec); ok {
					dt := dataType{
						Name:              dec.Name.Name,
						SingularLowercase: strings.ToLower(dec.Name.Name),
						SingularTitlecase: kace.Pascal(dec.Name.Name),
						PluralLowercase:   strings.ToLower(dec.Name.Name) + "s",
						PluralTitlecase:   kace.Pascal(dec.Name.Name) + "s",
						Fields:            []dataField{},
					}

					if _, ok := dec.Type.(*ast.StructType); !ok {
						log.Println("ERROR: only structs allowed in model declarations")
						return false
					}

					for _, field := range dec.Type.(*ast.StructType).Fields.List {
						df := dataField{
							Name: field.Names[0].Name,
						}

						if x, ok := field.Type.(*ast.Ident); ok {
							df.Type = x.Name
						} else if y, ok2 := field.Type.(*ast.StarExpr); ok2 {
							df.Pointer = true
							df.Type = y.X.(*ast.Ident).Name
						}

						tag := strings.Replace(strings.Replace(
							strings.Replace(field.Tag.Value, `naff:`, "", 1),
							"`", "", -1), `"`, "", -1)

						for _, t := range strings.Split(tag, ",") {
							_t := strings.ToLower(strings.TrimSpace(t))
							if _t == "createable" {
								df.ValidForCreationInput = true
							} else if _t == "editable" {
								df.ValidForUpdateInput = true
							}
						}
						dt.Fields = append(dt.Fields, df)
					}

					p.DataTypes = append(p.DataTypes, dt)
					p.IterableServicesImports = append(p.IterableServicesImports, filepath.Join(p.OutputRepository, "services", "v1", dt.PluralLowercase))
				}
				return true
			})
		}
	}

}

// ValidateName validates a project name
func (p Project) ValidateName(input string) error {
	if len(input) == 0 {
		return errors.New("empty string")
	}
	return nil
}

// ValidatePath validates an output path
func (p Project) ValidatePath(ans interface{}) error {
	return survey.Required(ans)
}

// RootPath returns what it calculates as the root path for the project.
func (p *Project) RootPath() string {
	return fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), p.OutputRepository)
}

// EnsureOutputDir ensures that our base path
func (p *Project) EnsureOutputDir() error {
	gp := os.Getenv("GOPATH")
	if gp == "" {
		return errors.New("GOPATH is not set")
	}

	folder := fmt.Sprintf("%s/src/%s", gp, p.OutputRepository)
	fi, err := os.Stat(folder)
	if err != nil {
		if os.IsNotExist(err) {
			if mkdirErr := os.MkdirAll(folder, 0777); mkdirErr != nil {
				return mkdirErr
			}
		}
	} else if !fi.IsDir() {
		return errors.New("destination is not a directory")
	}

	return nil
}

// RenderTemplateToPath renders a template file to a specific path
func renderTemplateToPath(t *template.Template, data interface{}, path string) error {
	p := filepath.Dir(path)
	if err := os.MkdirAll(p, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	err = t.Execute(f, data)
	if err != nil {
		return err
	}

	return f.Close()
}

const (
	baseTemplateDirectory = "template/base_repository/"
	iterableDirectory     = "template/iterables/"
)

func typeToPostgresType(t string) string {
	typeMap := map[string]string{
		"string":  "TEXT",
		"*string": "TEXT",
		"uint64":  "BIGINT",
		"*uint64": "BIGINT",
		"bool":    "BOOLEAN",
		"*bool":   "BOOLEAN",
	}

	if x, ok := typeMap[t]; ok {
		return x
	}

	log.Println("typeToPostgresType called for type: ", t)
	return t
}

// RenderDirectory renders a directory full of templates with the project as the data
func (p *Project) RenderDirectory() error {
	thisPackage := filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com/verygoodsoftwarenotvirus/naff")

	var fp string
	fp, _ = filepath.Abs(baseTemplateDirectory)
	if walkErr := filepath.Walk(fp,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}

			if !info.IsDir() {

				b, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				renderPath := strings.Replace(
					path,
					filepath.Join(thisPackage, baseTemplateDirectory),
					filepath.Join(os.Getenv("GOPATH"), "src", p.OutputRepository),
					1,
				)
				renderPath = strings.ReplaceAll(renderPath, ".tmpl", "")

				if strings.HasSuffix(path, defaultFileExtension) {
					t := template.Must(template.New(path).Funcs(map[string]interface{}{
						"typeToPostgresType": typeToPostgresType,
					}).Funcs(sprig.TxtFuncMap()).Parse(string(b)))

					if renderErr := renderTemplateToPath(t, p, renderPath); renderErr != nil {
						return renderErr
					}
				} else {
					if mkdirErr := os.MkdirAll(filepath.Dir(renderPath), os.ModePerm); mkdirErr != nil {
						return mkdirErr
					}

					return ioutil.WriteFile(renderPath, b, info.Mode())
				}
			}
			return nil
		},
	); walkErr != nil {
		return walkErr
	}

	for _, dt := range p.DataTypes {
		fp, _ = filepath.Abs(iterableDirectory)
		if walkErr := filepath.Walk(fp, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}

			if strings.HasSuffix(path, defaultFileExtension) {
				b, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				t := template.Must(template.New(path).Funcs(map[string]interface{}{
					"typeToPostgresType": typeToPostgresType,
				}).Funcs(sprig.TxtFuncMap()).Parse(string(b)))

				renderPath := strings.Replace(
					path,
					filepath.Join(thisPackage, iterableDirectory),
					filepath.Join(os.Getenv("GOPATH"), "src", p.OutputRepository),
					1,
				)
				renderPath = strings.ReplaceAll(renderPath, ".tmpl", "")
				renderPath = strings.Replace(renderPath, "services/v1/models/", fmt.Sprintf("services/v1/%s/", dt.PluralLowercase), 1)
				renderPath = strings.Replace(renderPath, "model.go", dt.SingularLowercase+".go", 1)
				renderPath = strings.Replace(renderPath, "mock_model_data_manager", fmt.Sprintf("mock_%s_data_manager", dt.SingularLowercase), 1)
				renderPath = strings.Replace(renderPath, "mock_model_data_server", fmt.Sprintf("mock_%s_data_server", dt.SingularLowercase), 1)
				renderPath = strings.Replace(renderPath, "model_test.go", dt.SingularLowercase+"_test.go", 1)
				renderPath = strings.Replace(renderPath, "models.go", dt.PluralLowercase+".go", 1)
				renderPath = strings.Replace(renderPath, "models_test.go", dt.PluralLowercase+"_test.go", 1)

				type tt struct {
					Project
					dataType
				}

				x := &tt{Project: *p, dataType: dt}

				if renderErr := renderTemplateToPath(t, x, renderPath); renderErr != nil {
					return renderErr
				}
			}
			return nil
		}); walkErr != nil {
			return walkErr
		}
	}

	return nil
}
