package main

import (
	"bytes"
	"errors"
	tojen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/tojen/gen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	client "gitlab.com/verygoodsoftwarenotvirus/naff/new_templates/client/v1/http"
	cmdv1server "gitlab.com/verygoodsoftwarenotvirus/naff/new_templates/cmd/server/v1"
	twofactor "gitlab.com/verygoodsoftwarenotvirus/naff/new_templates/cmd/tools/two_factor"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os/exec"
)

func runDiffForFiles(file1, file2 string) (string, error) {
	cmd := exec.Command("diff", file1, file2)
	var e *exec.ExitError

	b, err := cmd.Output()
	if err != nil {
		if !errors.Is(err, e) {
			return "", err
		}
	}

	if f, err := format.Source(b); err == nil {
		b = f
	}

	return string(b), nil
}

func getJenniferForFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	retBytes, err := tojen.GenerateFileBytes(b, "main", true, true)
	if err != nil {
		return "", err
	}

	return string(retBytes), nil
}

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
		return "", err
	}

	if f, err := format.Source(b); err == nil {
		b = f
	}

	return string(b), nil
}

func main() {
	todoDataTypes := []models.DataType{
		{
			Name: models.Name{
				Singular:                "Item",
				Plural:                  "Items",
				RouteName:               "items",
				PluralRouteName:         "item",
				UnexportedVarName:       "item",
				PluralUnexportedVarName: "items",
			},
		},
	}

	client.RenderPackage(todoDataTypes)
	cmdv1server.RenderPackage(todoDataTypes)
	twofactor.RenderPackage(todoDataTypes)

	//sourcePath := filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com/verygoodsoftwarenotvirus/todo")
	//path := fmt.Sprintf("%s/cmd/toosls/two_factor/main.go", sourcePath)
	//thing := doTheThingForFile(path, "main")
	//if thing != nil {
	//	println()
	//}
}

func doTheThingForFile(path, pkg string) error {
	src, err := runTojenForFile(path, pkg)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	code, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		return err
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
	//if fileFunc != nil {
	//	exprs[path] = fileFunc
	//}

	var b bytes.Buffer
	if err := printer.Fprint(&b, fset, code); err != nil {
		log.Fatal(err)
	}
	final := b.String()
	print(final)

	return nil
}
