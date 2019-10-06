package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func readerDotGo() *jen.File {
	ret := jen.NewFile("mock")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("_").Qual("io", "ReadCloser").Op("=").Parens(jen.Op("*").ID("ReadCloser")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Null().Type().ID("ReadCloser").Struct(
		jen.ID("mock").Dot(
			"Mock",
		),
	),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
