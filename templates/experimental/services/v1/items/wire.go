package items

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo() *jen.File {
	ret := jen.NewFile("items")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideItemsService"), jen.ID("ProvideItemDataManager"), jen.ID("ProvideItemDataServer")),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideItemDataManager turns a database into an ItemDataManager").ID("ProvideItemDataManager").Params(jen.ID("db").ID("database").Dot(
		"Database",
	)).Params(jen.ID("models").Dot(
		"ItemDataManager",
	)).Block(
		jen.Return().ID("db"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideItemDataServer is an arbitrary function for dependency injection's sake").ID("ProvideItemDataServer").Params(jen.ID("s").Op("*").ID("Service")).Params(jen.ID("models").Dot(
		"ItemDataServer",
	)).Block(
		jen.Return().ID("s"),
	),

		jen.Line(),
	)
	return ret
}
