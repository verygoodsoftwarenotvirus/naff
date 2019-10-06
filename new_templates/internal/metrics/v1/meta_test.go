package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func metaTestDotGo() *jen.File {
	ret := jen.NewFile("metrics")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Func().ID("TestRegisterDefaultViews").Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
		jen.ID("t").Dot(
			"Parallel",
		).Call(),
		jen.ID("require").Dot(
			"NoError",
		).Call(jen.ID("t"), jen.ID("RegisterDefaultViews").Call()),
	),
	)
	return ret
}
