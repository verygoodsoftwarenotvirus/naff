package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/codemodus/kace"

	tojen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/tojen/gen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

var (
	errCanStillComplete = errors.New("something done goofed, but I'll allow it")
)

func makeFuncLit(x *ast.FuncDecl) *ast.FuncLit {
	return &ast.FuncLit{
		Type: x.Type,
		Body: x.Body,
	}
}

func runTojenForFile(filename, pkg string) (string, error) {
	//cmd := exec.Command("/home/jeffrey/bin/tojen", "gen", "--formatted", fmt.Sprintf("--package=%s", pkg), filename)
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	b, err := tojen.GenerateFileBytes(fileBytes, pkg, false, true)
	if err != nil {
		log.Printf("file %q futzed up\n", filename)

		ee := err.Error()
		splitE := strings.Split(ee, "\n")
		e := splitE[0]
		fuckedFile := strings.Join(splitE[1:], "\n")

		return fuckedFile, fmt.Errorf("%s: %w", e, errCanStillComplete)
	}

	var f []byte
	if f, err = format.Source(b); err == nil {
		b = f
	}

	return string(b), nil
}

func main() {
	allPackages := []string{
		"gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types",
	}

	for _, pkg := range allPackages {
		err := doTheThingForPackage(filepath.Base(pkg), pkg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func doTheThingForPackage(pkg, pkgPath string) error {
	sourcePath := filepath.Join(os.Getenv("GOPATH"), "src", pkgPath)
	outputPath := strings.Replace(strings.Replace(sourcePath, "verygoodsoftwarenotvirus/todo", "verygoodsoftwarenotvirus/naff/templates/experimental", 1), "_test.go", "test_.go", 1)

	files, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		return err
	}

	fileMap := map[string]string{}

	for _, file := range files {
		if !file.IsDir() {
			fp := filepath.Join(sourcePath, file.Name())
			op := strings.TrimPrefix(
				strings.Replace(filepath.Join(outputPath, file.Name()), filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com/verygoodsoftwarenotvirus/naff/"), "", 1),
				"/",
			)
			op = strings.Replace(op, "_test.go", "_test_.go", 1)

			fileMap[strings.ReplaceAll(op, "templates/experimental/", "")] = kace.Camel(strings.ReplaceAll(filepath.Base(fp), ".go", "DotGo"))

			var code string
			if code, err = doTheThingForFile(fp, pkg, op, true); err != nil && !errors.Is(err, errCanStillComplete) {
				log.Printf(`error rendering %q to %q:
		%v
	%s`, fp, op, err, code)
			}
		}
	}

	genDotGo := buildGenDotGo(pkg, fileMap)
	ggfp := fmt.Sprintf("%s/gen.go", outputPath)
	if err = ioutil.WriteFile(ggfp, []byte(genDotGo), 0644); err != nil {
		return err
	}

	return utils.RunGoimportsForFile(ggfp)
}

func doTheThingForFile(path, pkg, outputPath string, writeFile bool) (string, error) {
	src, tojenErr := runTojenForFile(path, pkg)
	if tojenErr != nil && !errors.Is(tojenErr, errCanStillComplete) {
		return "", tojenErr
	} else if errors.Is(tojenErr, errCanStillComplete) {
		if writeFile {
			if err := ioutil.WriteFile(outputPath, []byte(src), 0644); err != nil {
				return "", err
			}
		}
		return src, nil
	} else if tojenErr == nil {
		fset := token.NewFileSet()
		code, err := parser.ParseFile(fset, "", src, parser.AllErrors)
		if err != nil {
			return "", err
		}

		code.Decls = consolidateDeclarations(code)
		if len(code.Decls) >= 2 {
			switch code.Decls[1].(type) {
			case *ast.FuncDecl:
				code.Decls[1].(*ast.FuncDecl).Name.Name = kace.Camel(strings.ReplaceAll(filepath.Base(path), ".go", "DotGo"))
			default:
				println()
			}
		}

		var b bytes.Buffer
		if err = printer.Fprint(&b, fset, code); err != nil {
			return "", err
		}

		outputBytes := b.Bytes()

		replacers := []struct {
			replacer    string
			replacement string
		}{
			{
				replacer:    `\)\n\tcode\.Add\(`,
				replacement: ")\n\n\tcode.Add(",
			},
			{
				replacer:    `\,\n\n\t\tjen\.Newline\(\)\,`,
				replacement: ",\n\t\tjen.Newline(),",
			},
			{
				replacer:    `\.Comment\(\"\/\/\s`,
				replacement: `.Comment("`,
			},
			{
				replacer:    `jen\.Func\(\)\.Comment\(\"([0-9a-zA-Z\d\-\,\:\.\*\s\t\(\)\/\\\'` + "`" + `]+)\"\)`,
				replacement: "jen.Comment(\"$1\"),\n\t\tjen.Newline(),\n\t\tjen.Func()",
			},
			{
				replacer:    `code\.Add\(jen`,
				replacement: "code.Add(\n\t\tjen",
			},
			{
				replacer:    `code\.Add\(\n\t\tjen\.Null\(\)\,\n\t\tjen\.Newline\(\)\,\n\t\)`,
				replacement: "",
			},
			{
				replacer:    `\)\.Body\(jen`,
				replacement: ").Body(\n\t\tjen",
			},
			{
				replacer:    `\.Dot\(\n\s+\"(\w+)\"\,\n\s+\)`,
				replacement: `.Dot("$1")`,
			},
			{
				replacer:    `\)\n\treturn code\n\}`,
				replacement: ")\n\n\treturn code\n}",
			},
			{
				replacer:    `jen\.Null\(\)\.Var\(\)`,
				replacement: "jen.Var()",
			},
			{
				replacer:    `\n\tcode\.Add\(func\(\) jen\.Code \{\n\t\treturn nil\n\t\}\,\n\t\tjen\.Newline\(\)\,\n\t\)`,
				replacement: "",
			},
		}

		for _, r := range replacers {
			outputBytes = regexp.MustCompile(r.replacer).ReplaceAll(outputBytes, []byte(r.replacement))
		}

		outputPath = strings.Replace(outputPath, "_test.go", "_test_.go", 1)

		// i don't care
		_ = os.MkdirAll(filepath.Dir(outputPath), 0777)

		if writeFile {
			if err = ioutil.WriteFile(outputPath, outputBytes, 0644); err != nil {
				return "", err
			}

			if err = utils.RunGoimportsForFile(outputPath); err != nil {
				return "", err
			}
		}

		return b.String(), nil
	}
	return src, nil
}

func consolidateDeclarations(code *ast.File) []ast.Decl {
	decls := []ast.Decl{}
	funcLits := map[string]*ast.FuncLit{}
	callExprs := map[string]*ast.CallExpr{}

	for _, dec := range code.Decls {
		switch x := dec.(type) {
		case *ast.FuncDecl:
			if x.Name.Name == "genFile" {
				for i := 1; i < len(x.Body.List); i++ {
					statement := x.Body.List[i] // this is our ret.Adds
					switch y := statement.(type) {
					case *ast.ExprStmt:
						switch z := y.X.(type) {
						case *ast.CallExpr:
							if len(z.Args) > 0 {
								switch a := z.Args[0].(type) {
								case *ast.CallExpr:
									switch b := a.Fun.(type) {
									case *ast.Ident:
										ce, callExprFound := callExprs[b.Name]
										fl, funcLitFound := funcLits[b.Name]

										if callExprFound {
											z.Args[0] = ce
										} else if funcLitFound {
											// we've found a call to a function we have a literal for
											z.Args[0] = fl
										}
									}
								}
							}
						}
					}
				}
				//fileFunc = makeFuncLit(x)
				decls = append(decls, dec)
			} else {
				funcLits[x.Name.Name] = makeFuncLit(x)
				switch y := x.Body.List[0].(*ast.ReturnStmt).Results[0].(type) {
				case *ast.CallExpr:
					callExprs[x.Name.Name] = y
				default:
					println(y)
				}
			}
		default:
			decls = append(decls, dec)
		}
	}
	return decls
}

func buildGenDotGo(pkgName string, fileToFunctionMap map[string]string) string {
	var baseDir string
	for k := range fileToFunctionMap {
		baseDir = filepath.Dir(k)
	}

	start := fmt.Sprintf(`package %s

import (
	_ "embed"
	"path/filepath"
	"fmt"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)
	
const (
	packageName = "%s"
	
	basePackagePath = "%s"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
`, pkgName, pkgName, baseDir)

	var fileDecs string
	for file, fun := range fileToFunctionMap {
		fileDecs = fmt.Sprintf("%s\t%q:\t%s(proj),\n", fileDecs, strings.ReplaceAll(filepath.Base(file), "_test_.go", "_test.go"), fun)
	}

	end := `
}

	//for _, typ := range types {
	//	files[fmt.Sprintf("%s.go", typ.Name.PluralRouteName)] = itemsDotGo(typ)
	//	files[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName)] = itemsTestDotGo(typ)
	//}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}`

	return fmt.Sprintf("%s%s%s", start, fileDecs, end)
}
