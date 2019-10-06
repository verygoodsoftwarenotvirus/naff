package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func usersDotGo() *jen.File {
	ret := jen.NewFile("dbclient")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"UserDataManager",
	).Op("=").Parens(jen.Op("*").ID("Client")).Call(jen.ID("nil")).Var().ID("ErrUserExists").Op("=").Qual("errors", "New").Call(jen.Lit("error: username already exists")),
	)
	ret.Add(jen.Func().ID("attachUsernameToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("username").ID("string")).Block(
		jen.If(
			jen.ID("span").Op("!=").ID("nil"),
		).Block(
			jen.ID("span").Dot(
				"AddAttributes",
			).Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("username"), jen.ID("username"))),
		),
	),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
