package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func mainTestDotGo() *jen.File {
	ret := jen.NewFile("models")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Func().ID("TestErrorResponse_Error").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("_").Op("=").Parens(jen.Op("&").ID("ErrorResponse").Valuesln()).Dot(
				"Error",
			).Call(),
		)),
	),
	)
	return ret
}
