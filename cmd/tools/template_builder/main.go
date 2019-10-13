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

	if f, err := format.Source(b); err == nil {
		b = f
	}

	return string(b), nil
}

func main() {
	allPackages := []string{
		//// completed
		// "gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http",
		//// to be completed
		"gitlab.com/verygoodsoftwarenotvirus/todo/cmd/server/v1",
		"gitlab.com/verygoodsoftwarenotvirus/todo/cmd/tools/two_factor",
		"gitlab.com/verygoodsoftwarenotvirus/todo/cmd/config_gen/v1",
		"gitlab.com/verygoodsoftwarenotvirus/todo/database/v1",
		"gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/client",
		"gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/queriers/mariadb",
		"gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/queriers/postgres",
		"gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/queriers/sqlite",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/config",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock",
		"gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
		"gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock",
		"gitlab.com/verygoodsoftwarenotvirus/todo/server/v1",
		"gitlab.com/verygoodsoftwarenotvirus/todo/server/v1/http",
		"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/auth",
		"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/frontend",
		"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/items",
		"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/oauth2clients",
		"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/users",
		"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/webhooks",
		"gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/frontend",
		"gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/integration",
		"gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/load",
		"gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil",
		"gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/mock",
		"gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/rand/model",
	}

	//path := "gitlab.com/verygoodsoftwarenotvirus/todo/cmd/config_gen/v1"

	for _, pkg := range allPackages {
		err := doTheThingForPackage(filepath.Base(pkg), pkg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func doTheThingForPackage(pkg, pkgPath string) error {
	sourcePath := filepath.Join(os.Getenv("GOPATH"), "src", pkgPath)
	outputPath := strings.Replace(sourcePath, "verygoodsoftwarenotvirus/todo", "verygoodsoftwarenotvirus/naff/templates/experimental", 1)

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

			fileMap[strings.ReplaceAll(op, "templates/experimental/", "")] = kace.Camel(strings.ReplaceAll(filepath.Base(fp), ".go", "DotGo"))

			if code, err := doTheThingForFile(fp, pkg, op, true); err != nil && !errors.Is(err, errCanStillComplete) {
				log.Printf(`error rendering
	%q to %q:
		%v
	%s`, fp, op, err, code)
			}
		}
	}

	genDotGo := buildGenDotGo(pkg, fileMap)
	ggfp := fmt.Sprintf("%s/gen.go", outputPath)

	if err := ioutil.WriteFile(ggfp, []byte(genDotGo), 0644); err != nil {
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
		if err := printer.Fprint(&b, fset, code); err != nil {
			return "", err
		}

		if writeFile {
			if err := ioutil.WriteFile(outputPath, b.Bytes(), 0644); err != nil {
				return "", err
			}

			if err := utils.RunGoimportsForFile(outputPath); err != nil {
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
	start := fmt.Sprintf(`package %s

import (
	"fmt"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
`, pkgName)

	var fileDecs string
	for file, fun := range fileToFunctionMap {
		fileDecs = fmt.Sprintf("%s\t%q:\t%s(),\n", fileDecs, file, fun)
	}

	end := `
}

	//for _, typ := range types {
	//	files[fmt.Sprintf("client/v1/http/%s.go", typ.Name.PluralRouteName)] = itemsDotGo(typ)
	//	files[fmt.Sprintf("client/v1/http/%s_test.go", typ.Name.PluralRouteName)] = itemsTestDotGo(typ)
	//}

	for path, file := range files {
		if err := utils.RenderFile(path, file); err != nil {
			return err
		}
	}

	return nil
}`

	return fmt.Sprintf("%s%s%s", start, fileDecs, end)
}
