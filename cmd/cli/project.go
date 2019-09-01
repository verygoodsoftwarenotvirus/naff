package main

import (
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

	"gitlab.com/verygoodsoftwarenotvirus/naff/embedded"

	"github.com/Masterminds/sprig"
	"github.com/codemodus/kace"
	"github.com/pkg/errors"
	"gopkg.in/AlecAivazis/survey.v1"
)

type dataType struct {
	Name   string
	Fields []dataField
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
	p := Project{}
	if surveyErr := survey.Ask(questions, &p); surveyErr != nil {
		return nil, surveyErr
	}
	os.RemoveAll(filepath.Join(os.Getenv("GOPATH"), "src", p.OutputRepository))

	return &p, nil
}

func (p *Project) parseModels() {
	fullModelsPath := filepath.Join(os.Getenv("GOPATH"), "src", p.ModelsPackage)

	packages, err := parser.ParseDir(token.NewFileSet(), fullModelsPath, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range packages {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				if dec, ok := n.(*ast.TypeSpec); ok {
					dt := dataType{
						Name:   dec.Name.Name,
						Fields: []dataField{},
					}

					if _, ok := dec.Type.(*ast.StructType); !ok {
						log.Println("ERROR: only structs allowed in model declarations")
						return false
					}

					for _, field := range dec.Type.(*ast.StructType).Fields.List {
						df := dataField{
							Name:                  field.Names[0].Name,
							ValidForCreationInput: true,
							ValidForUpdateInput:   true,
						}

						if x, ok := field.Type.(*ast.Ident); ok {
							df.Type = x.Name
						} else if y, ok2 := field.Type.(*ast.StarExpr); ok2 {
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
							_t := strings.ToLower(strings.TrimSpace(t))
							if _t == "!createable" {
								df.ValidForCreationInput = false
							} else if _t == "!editable" {
								df.ValidForUpdateInput = false
							}
						}
						dt.Fields = append(dt.Fields, df)
					}

					p.DataTypes = append(p.DataTypes, dt)
					p.IterableServicesImports = append(
						p.IterableServicesImports,
						filepath.Join(
							p.OutputRepository,
							"services",
							"v1",
							strings.ToLower(dt.Name)+"s",
						),
					)
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
		return errors.Wrap(err, "creating directory")
	}

	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "creating file")
	}

	err = t.Execute(f, data)
	if err != nil {
		return errors.Wrap(err, "executing template")
	}

	return f.Close()
}

func typeToPostgresType(t string) string {
	typeMap := map[string]string{
		"[]string": "CHARACTER VARYING",
		"string":   "CHARACTER VARYING",
		"*string":  "CHARACTER VARYING",
		"uint64":   "BIGINT",
		"*uint64":  "BIGINT",
		"bool":     "BOOLEAN",
		"*bool":    "BOOLEAN",
		"int":      "INTEGER",
		"*int":     "INTEGER",
		"uint":     "INTEGER",
		"*uint":    "INTEGER",
		"float64":  "NUMERIC",
	}

	if x, ok := typeMap[t]; ok {
		return x
	}

	log.Println("typeToPostgresType called for type: ", t)
	return t
}

func typeToSqliteType(t string) string {
	typeMap := map[string]string{
		"[]string": "CHARACTER VARYING",
		"string":   "CHARACTER VARYING",
		"*string":  "CHARACTER VARYING",
		"uint64":   "INTEGER",
		"*uint64":  "INTEGER",
		"bool":     "BOOLEAN",
		"*bool":    "BOOLEAN",
		"int":      "INTEGER",
		"*int":     "INTEGER",
		"uint":     "INTEGER",
		"*uint":    "INTEGER",
		"float64":  "REAL",
	}

	if x, ok := typeMap[t]; ok {
		return x
	}

	log.Println("typeToSqliteType called for type: ", t)
	return t
}

func typeExample(t string, pointer bool) interface{} {
	typeMap := map[string]interface{}{
		"[]string": `[]string{"example"}`,
		"string":   `"example"`,
		"*string":  `func(s string) *string { return &s }("example")`,
		"uint64":   "uint64(123)",
		"*uint64":  "func(u uint64) *uint64 { return &u }(123)",
		"bool":     false,
		"*bool":    "func(b bool) *bool { return &b }(false)",
		"int":      "int(456)",
		"*int":     "func(i int) *int { return &i }(123)",
		"uint":     "uint(456)",
		"*uint":    "func(i uint) *uint { return &i }(123)",
		"float64":  "float64(12.34)",
	}

	tn := t
	if pointer {
		tn = fmt.Sprintf("*%s", tn)
	}

	if x, ok := typeMap[tn]; ok {
		return x
	}

	return t
}

const (
	baseTemplateDirectory = "template/base_repository/"
	iterableDirectory     = "template/iterables/"
)

// RenderDirectory renders a directory full of templates with the project as the data
func (p *Project) RenderDirectory() error {
	//thisPackage := filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com/verygoodsoftwarenotvirus/naff")

	baseFiles, err := embedded.WalkDirs(baseTemplateDirectory, false)
	if err != nil {
		return errors.Wrap(err, "fetching directory from embedded files")
	}

	for _, path := range baseFiles {
		b, err := embedded.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "reading embedded file at path: %q", path)
		}
		renderPath := strings.Replace(
			path,
			baseTemplateDirectory,
			filepath.Join(os.Getenv("GOPATH"), "src", p.OutputRepository)+"/",
			1,
		)
		renderPath = strings.TrimSuffix(renderPath, ".tmpl")
		fmt.Printf("rendering file: %q\n", renderPath)

		if strings.HasSuffix(path, defaultFileExtension) {
			t := template.Must(template.New(path).Funcs(map[string]interface{}{
				"typeToPostgresType": typeToPostgresType,
				"typeToSqliteType":   typeToSqliteType,
				"typeExample":        typeExample,
				"camelCase":          kace.Camel,
				"pascal":             kace.Pascal,
			}).Funcs(sprig.TxtFuncMap()).Parse(string(b)))

			if renderErr := renderTemplateToPath(t, p, renderPath); renderErr != nil {
				return errors.Wrap(renderErr, "rendering template")
			}
		} else {
			if mkdirErr := os.MkdirAll(filepath.Dir(renderPath), os.ModePerm); mkdirErr != nil {
				return errors.Wrap(mkdirErr, "creating containing folder")
			}

			if renderErr := ioutil.WriteFile(renderPath, b, 0644); renderErr != nil {
				return errors.Wrap(renderErr, "rendering template")
			}
		}
	}

	dataFiles, err := embedded.WalkDirs(iterableDirectory, false)
	if err != nil {
		return errors.Wrap(err, "fetching directory from embedded files")
	}

	for _, dt := range p.DataTypes {
		for _, path := range dataFiles {
			fmt.Printf("iterating over file %q for data type %q\n", path, dt.Name)

			if strings.HasSuffix(path, defaultFileExtension) {
				b, err := embedded.ReadFile(path)
				if err != nil {
					return errors.Wrapf(err, "reading embedded file at path: %q", path)
				}

				t := template.Must(template.New(path).Funcs(map[string]interface{}{
					"typeToPostgresType": typeToPostgresType,
					"typeToSqliteType":   typeToSqliteType,
					"typeExample":        typeExample,
					"camelCase":          kace.Camel,
					"pascal":             kace.Pascal,
				}).Funcs(sprig.TxtFuncMap()).Parse(string(b)))

				renderPath := strings.Replace(
					path,
					iterableDirectory,
					filepath.Join(os.Getenv("GOPATH"), "src", p.OutputRepository)+"/",
					1,
				)

				renderPath = strings.TrimSuffix(renderPath, ".tmpl")
				renderPath = strings.Replace(renderPath, "services/v1/models/", fmt.Sprintf("services/v1/%ss/", strings.ToLower(dt.Name)), 1)
				renderPath = strings.Replace(renderPath, "frontend/v1/src/pages/models", fmt.Sprintf("frontend/v1/src/pages/%ss/", strings.ToLower(dt.Name)), 1)
				renderPath = strings.Replace(renderPath, "model.go", strings.ToLower(dt.Name)+".go", 1)
				renderPath = strings.Replace(renderPath, "mock_model_data_manager", fmt.Sprintf("mock_%s_data_manager", strings.ToLower(dt.Name)), 1)
				renderPath = strings.Replace(renderPath, "mock_model_data_server", fmt.Sprintf("mock_%s_data_server", strings.ToLower(dt.Name)), 1)
				renderPath = strings.Replace(renderPath, "model_test.go", strings.ToLower(dt.Name)+"_test.go", 1)
				renderPath = strings.Replace(renderPath, "models.go", strings.ToLower(dt.Name+"s.go"), 1)
				renderPath = strings.Replace(renderPath, "models_test.go", strings.ToLower(dt.Name+"s_test.go"), 1)

				type tt struct {
					Project
					dataType
					Name string
				}

				x := &tt{Project: *p, dataType: dt, Name: dt.Name}

				fmt.Printf("rendering file: %q\n", renderPath)
				if renderErr := renderTemplateToPath(t, x, renderPath); renderErr != nil {
					return renderErr
				}
			}
		}
	}

	return nil
}
