package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func actionsDotGo() *jen.File {
	ret := jen.NewFile("main")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("ErrUnavailableYet").Op("=").Qual("errors", "New").Call(jen.Lit("can't do this yet")),
	)
	ret.Add(jen.Null().Type().ID("actionFunc").Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Type().ID("Action").Struct(
		jen.ID("Action").ID("actionFunc"),
		jen.ID("Weight").ID("int"),
		jen.ID("Name").ID("string"),
	),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
