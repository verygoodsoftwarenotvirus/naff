package gen

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strconv"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

const jenImp = "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
const utilsImp = "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"

func funcDecl(s *ast.FuncDecl) jen.Code {
	ret := jen.Qual(jenImp, "Func").Call()

	if s.Doc != nil {
		for _, com := range s.Doc.List {
			ret.Add(jen.ID(".Comment").Call(jen.Lit(com.Text)))
		}
	}

	if s.Recv != nil {
		ret.Dot("Params").Call(fieldList(s.Recv)...)
	}
	ret.Add(ident(s.Name))
	ret.Add(funcType(s.Type))
	ret.Add(blockStmt(s.Body))
	return ret
}

var paths = map[string]string{}
var formatting = false

// GenerateFileBytes takes an array of bytes and transforms it into jennifer
// code
func GenerateFileBytes(s []byte, packName string, main bool, formatting bool) ([]byte, error) {
	file := GenerateFile(s, packName, main)
	b := &bytes.Buffer{}
	err := file.Render(b)
	if err != nil {
		return s, err
	}
	ret := b.Bytes()
	if formatting {
		ret, err = goFormat(ret)
		if err != nil {
			return ret, err
		}
	}
	return ret, nil
}

func imports(imports []*ast.ImportSpec) (map[string]string, []jen.Code) {
	p := make(map[string]string)
	anonImports := []jen.Code{}
	for _, i := range imports {
		pathVal := i.Path.Value[1 : len(i.Path.Value)-1]
		name := pathVal

		if i.Name == nil {
			idx := strings.Index(pathVal, "/")
			if idx != -1 {
				name = pathVal[idx+1:]
			}
		} else {
			name = i.Name.String()
			if name == "." {
				panic(". imports not supported")
			}
			if name == "_" {
				anonImports = append(anonImports, jen.Lit(pathVal))
				continue
			}
		}
		p[name] = pathVal
	}
	return p, anonImports
}

// GenerateFile Generates a jennifer file given a series of bytes a package name
// and if you want a main function or not
func GenerateFile(s []byte, packName string, main bool) *jen.File {
	file := jen.NewFile(packName)
	astFile := parseFile(s)
	var anonImports []jen.Code
	// paths is a global variable to map the exported object to the import
	paths, anonImports = imports(astFile.Imports)

	// generate the generative code based on the file
	decls := []string{}
	for _, decl := range astFile.Decls {
		code, name := makeJenCode(decl)
		file.Add(code)
		decls = append(decls, name)
	}

	// generate the function that pieces together all the code
	var codes []jen.Code
	codes = append(codes, genNewJenFile(astFile.Name.String()))
	// add anon imports i.e. _ for side effects
	if len(anonImports) > 0 {
		codes = append(codes, jen.ID("code").Dot("Anon").Call(anonImports...))
	}
	// add the generated functions to the created jen file
	for _, name := range decls {
		codes = append(codes, jen.ID("code").Dot("Add").Callln(jen.ID(name).Call(), jen.Qual(jenImp, "Newline").Call()))
	}
	// return the created jen file
	codes = append(codes, jen.Return().ID("code"))
	// add the patch function to the output file
	file.Add(
		jen.Func().ID("genFile").Params(jen.ID("proj").PointerTo().Qual("gitlab.com/verygoodsoftwarenotvirus/naff/models", "Project")).Op("*").Qual(jenImp, "File").Body(codes...),
	)
	// if main then generate a main function that prints out the output of the
	// patch function
	if main {
		file.Add(genMainFunc())
	}
	return file
}

func genNewJenFile(name string) jen.Code {
	const varName = "code"
	return jen.ID(varName).Assign().Qual(jenImp, "NewFile").Call(jen.ID("packageName")).Newline().Newline().Qual(utilsImp, "AddImports").Call(jen.ID("proj"), jen.ID(varName), jen.False()).Newline().Newline()
}

func genMainFunc() jen.Code {
	return jen.Func().ID("main").Params().Body(
		jen.ID("code").Assign().ID("genFile").Call(),
		jen.Qual("fmt", "Printf").Call(
			jen.Lit("%#v"),
			jen.ID("code"),
		),
	)
}

func makeJenCode(s ast.Decl) (jen.Code, string) {
	inner := jen.Null()
	name := ""
	switch t := s.(type) {
	case *ast.GenDecl:
		name = "genDeclAt" + strconv.Itoa(int(t.TokPos))
		inner.Add(genDecl(t))
	case *ast.FuncDecl:
		name = "genFunc" + t.Name.String()
		inner.Add(funcDecl(t))
	default:
		log.Println("UNHANDLED TYPE")
	}
	return makeJenFileFunc(name, inner), name
}

func makeJenFileFunc(name string, block jen.Code) jen.Code {
	return jen.Func().ID(name).Params().Qual(jenImp, "Code").Body(
		jen.Return().Add(block),
	)
}

func parseFile(code []byte) *ast.File {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return f
}

func genDecl(g *ast.GenDecl) jen.Code {
	var wrapperRet *jen.Statement
	var decls []jen.Code

	for _, spec := range g.Specs {
		switch s := spec.(type) {
		case *ast.ValueSpec:
			wrapperRet = jen.Qual(jenImp, "Var").Call()
			decls = append(decls, valueSpec(s))
		case *ast.TypeSpec:
			wrapperRet = jen.Qual(jenImp, "Type").Call()
			decls = append(decls, typeSpec(s))
		}
	}

	if wrapperRet == nil {
		return jen.Nil()
	}

	return wrapperRet.Dot("Defs").Callln(decls...)
}

func typeSpec(s *ast.TypeSpec) jen.Code {
	return jen.ID("jen").Add(ident(s.Name)).Add(genExpr(s.Type))
}

func valueSpec(s *ast.ValueSpec) jen.Code {
	ret := jen.ID("jen")
	ret.Add(identsList(s.Names))
	ret.Add(genExpr(s.Type))
	if len(s.Values) > 0 {
		ret.Dot("Op").Call(jen.Lit("="))
		ret.Add(genExprs(s.Values))
	}
	return ret
}

func basicLit(b *ast.BasicLit) jen.Code {
	switch b.Kind {
	case token.INT:
		i, err := strconv.ParseInt(b.Value, 10, 32)
		if err != nil {
			return nil
		}
		return jen.Dot("Lit").Call(jen.Lit(int(i)))
	case token.FLOAT:
		return jen.Dot("Lit").Call(jen.ID(b.Value))
	case token.IMAG:
		panic("Cannot parse Imaginary Numbers")
	case token.CHAR:
		return jen.Dot("ID").Call(jen.ID("\"" + b.Value + "\""))
	case token.STRING:
		return jen.Dot("Lit").Call(jen.ID(b.Value))
	}
	return nil
}
