package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo() *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers represents what this library offers to external users in the form of dependencies"),
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Callln(
				jen.ID("ProvideUnitCounter"),
				jen.ID("ProvideUnitCounterProvider"),
			),
		),
		jen.Line(),
	)
	return ret
}