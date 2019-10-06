package httpserver

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("paramFetcherProviders"), jen.ID("ProvideServer"), jen.ID("ProvideNamespace"), jen.ID("ProvideNewsmanTypeNameManipulationFunc")),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
