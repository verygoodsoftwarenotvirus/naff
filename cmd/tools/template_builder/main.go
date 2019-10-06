package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/codemodus/kace"

	tojen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/tojen/gen"
)

func makeFuncLit(x *ast.FuncDecl) *ast.FuncLit {
	return &ast.FuncLit{
		Type: x.Type,
		Body: x.Body,
	}
}

func runGoimportsForFile(filename string) error {
	return exec.Command("/home/jeffrey/bin/goimports", "-w", filename).Run()
}

func runTojenForFile(filename, pkg string) (string, error) {
	//cmd := exec.Command("/home/jeffrey/bin/tojen", "gen", "--formatted", fmt.Sprintf("--package=%s", pkg), filename)
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		if e, ok := err.(*os.PathError); ok {
			if e.Err == syscall.EISDIR {
				return "", nil
			} else {
				return "", err
			}
		} else {
			return "", err
		}
	}

	b, err := tojen.GenerateFileBytes(fileBytes, pkg, false, true)
	if err != nil {
		log.Printf("file %q futzed up\n%s", filename, string(b))
		return "", nil
	}

	if f, err := format.Source(b); err == nil {
		b = f
	}

	return string(b), nil
}

func main() {
	allPackages := []string{
		// completed
		//		"gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http",
		//		"gitlab.com/verygoodsoftwarenotvirus/todo/cmd/server/v1",
		//		"gitlab.com/verygoodsoftwarenotvirus/todo/cmd/tools/two_factor",
		// to complete
		"gitlab.com/verygoodsoftwarenotvirus/todo/cmd/config_gen/v1",
		"gitlab.com/verygoodsoftwarenotvirus/todo/database/v1",
		"gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/client",
		"gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/queriers/mariadb",
		"gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/queriers/postgres",
		"gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/queriers/sqlite",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1/mock",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/config/v1",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock",
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

	//todoDataTypes := []models.DataType{
	//	{
	//		Name: models.Name{
	//			Singular:                "Item",
	//			Plural:                  "Items",
	//			RouteName:               "items",
	//			PluralRouteName:         "item",
	//			UnexportedVarName:       "item",
	//			PluralUnexportedVarName: "items",
	//		},
	//	},
	//}
	//
	//client.RenderPackage(todoDataTypes)
	//cmdv1server.RenderPackage(todoDataTypes)
	//twofactor.RenderPackage(todoDataTypes)

	//path := "gitlab.com/verygoodsoftwarenotvirus/todo/cmd/config_gen/v1"

	for _, pkg := range allPackages {
		err := doTheThingForPackage("main", pkg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func doTheThingForPackage(pkg, pkgPath string) error {
	sourcePath := filepath.Join(os.Getenv("GOPATH"), "src", pkgPath)
	outputPath := strings.Replace(sourcePath, "verygoodsoftwarenotvirus/todo", "verygoodsoftwarenotvirus/naff/new_templates", 1)

	files, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			fp := filepath.Join(sourcePath, file.Name())
			op := filepath.Join(outputPath, file.Name())
			if code, err := doTheThingForFile(fp, pkg, op, true); err != nil {
				log.Printf("error rendering file %q: %v\n\n%s", fp, err, code)
			}
		}
	}

	return nil
}

func doTheThingForFile(path, pkg, outputPath string, writeFile bool) (string, error) {
	src, err := runTojenForFile(path, pkg)
	if err != nil {
		return "", err
	}

	fset := token.NewFileSet()
	code, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		return "", err
	}

	funcLits := map[string]*ast.FuncLit{}
	callExprs := map[string]*ast.CallExpr{}
	//var fileFunc *ast.FuncLit
	decls := []ast.Decl{}

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
	code.Decls = decls

	code.Decls[1].(*ast.FuncDecl).Name.Name = kace.Camel(strings.ReplaceAll(filepath.Base(path), ".go", "DotGo"))
	//if fileFunc != nil {
	//	exprs[path] = fileFunc
	//}

	var b bytes.Buffer
	if err := printer.Fprint(&b, fset, code); err != nil {
		return "", err
	}

	_, existenceErr := os.Stat("/path/to/whatever")
	fileExists := existenceErr != nil && os.IsNotExist(err)

	if writeFile && !fileExists {
		if err := ioutil.WriteFile(outputPath, b.Bytes(), 0644); err != nil {
			return "", err
		}
		runGoimportsForFile(outputPath)
	}

	return b.String(), nil
}
