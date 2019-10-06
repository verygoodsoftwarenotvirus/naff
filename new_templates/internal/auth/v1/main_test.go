package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func mainTestDotGo() *jen.File {
	ret := jen.NewFile("auth_test")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Func().ID("TestProvideBcryptHashCost").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("auth").Dot(
				"ProvideBcryptHashCost",
			).Call(),
		)),
	),
	)
	return ret
}
