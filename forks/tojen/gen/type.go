package gen

import (
	"go/ast"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func funcType(s *ast.FuncType) jen.Code {
	var ret jen.Statement
	params := fieldList(s.Params)
	ret.Dot("Params").Call(params...)
	results := fieldList(s.Results)
	if len(results) > 0 {
		ret.Dot("Params").Call(results...)
	}
	return &ret
}

func arrayType(s *ast.ArrayType) jen.Code {
	return jen.Dot("Index").Call().Add(genExpr(s.Elt))
}

func structType(s *ast.StructType) jen.Code {
	args := fieldList(s.Fields)
	if len(args) <= 1 {
		return jen.Dot("Struct").Call(args...)
	}
	return jen.Dot("Struct").Callln(args...)
}

func interfaceType(s *ast.InterfaceType) jen.Code {
	args := fieldList(s.Methods)
	if len(args) <= 1 {
		return jen.Dot("Interface").Call(args...)
	}
	return jen.Dot("Interface").Callln(args...)
}
