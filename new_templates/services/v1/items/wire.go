package items

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo() *jen.File {
	ret := jen.NewFile("items")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideItemsService"), jen.ID("ProvideItemDataManager"), jen.ID("ProvideItemDataServer")),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
