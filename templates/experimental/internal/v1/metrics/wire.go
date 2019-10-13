package metrics

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func wireDotGo() *jen.File {
	ret := jen.NewFile("metrics")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideUnitCounter"), jen.ID("ProvideUnitCounterProvider")),

		jen.Line(),
	)
	return ret
}
