package frontend

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func frontendServiceTestDotGo() *jen.File {
	ret := jen.NewFile("frontend")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestProvideFrontendService").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("ProvideFrontendService").Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("config").Dot(
				"FrontendSettings",
			).Valuesln()),
		)),
	),

		jen.Line(),
	)
	return ret
}
