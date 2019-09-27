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
	"os/exec"
	"strings"

	tojen "gitlab.com/verygoodsoftwarenotvirus/naff/lib/tojen/gen"
	client "gitlab.com/verygoodsoftwarenotvirus/naff/new_templates/client/v1/http"
)

func runTojenForFile(filename, pkg string) (string, error) {
	cmd := exec.Command("tojen", "gen", "--formatted", fmt.Sprintf("--package=%s", pkg), filename)

	b, err := cmd.Output()
	if err != nil {
		return "", err
	}

	if f, err := format.Source(b); err == nil {
		b = f
	}

	return string(b), nil
}

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

func GetJenniferForFile(filename string) (string, error) {
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

func main() {
	for path, file := range client.Files {
		fp := fmt.Sprintf("/home/jeffrey/src/gitlab.com/verygoodsoftwarenotvirus/naff/templates/%s", path)
		_ = os.Remove(fp)

		ogfp := strings.Replace(fp, "naff/templates", "todo", 1)
		ogBytes, err := ioutil.ReadFile(ogfp)
		if err != nil {
			log.Fatal(err)
		}
		ogFile := string(ogBytes)
		_ = ogFile

		var b bytes.Buffer
		if err := file.Render(&b); err != nil {
			log.Fatal(err)
		}

		if err := ioutil.WriteFile(fp, b.Bytes(), os.ModePerm); err != nil {
			log.Fatal(err)
		}

		//diff, err := runDiffForFiles(ogfp, fp)
		//if err != nil {
		//	func(error) {}(err)
		//}
		//_ = diff
	}

	doTheThingForFile("/home/jeffrey/src/gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http/items_test.go", "client")

	//exprs := map[string]ast.Expr{}
	//sourcePath := filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com/verygoodsoftwarenotvirus/todo")
	//if err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
	//	outputFilepath := strings.Replace(path, "todo", "naff/template2/base_repository", 1)
	//	relativePath := strings.Replace(path, sourcePath+"/", "", 1)
	//
	//	if err != nil {
	//		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
	//		return err
	//	}
	//
	//	if info.IsDir() {
	//		if _, ok := skipDirectories[relativePath]; ok {
	//			return filepath.SkipDir
	//		}
	//
	//		if _, ok := iterableDirectories[relativePath]; ok {
	//			outputFilepath = strings.Replace(outputFilepath, "base_repository", "iterables", 1)
	//		}
	//
	//		return nil // os.MkdirAll(outputFilepath, info.Mode())
	//	}
	//	if _, ok := skipFiles[relativePath]; ok {
	//		return nil
	//	}
	//
	//	if strings.HasSuffix(path, ".go") {
	//		doTheThing(path)
	//	}
	//
	//	return nil
	//}); err != nil {
	//	log.Fatal(err)
	//}

	//kvs := []ast.Expr{}
	//for path, funcLit := range exprs {
	//	kvs = append(kvs, &ast.KeyValueExpr{
	//		Key: &ast.BasicLit{
	//			Value: fmt.Sprintf("%q", path),
	//			Kind:  token.STRING,
	//		},
	//		Value: funcLit,
	//	})
	//}
	//f.Decls[1].(*ast.GenDecl).Specs[0].(*ast.ValueSpec).Values[0].(*ast.CompositeLit).Elts = kvs

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
