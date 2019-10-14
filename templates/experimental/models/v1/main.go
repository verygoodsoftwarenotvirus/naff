package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func mainDotGo() *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("SortAscending").ID("sortType").Op("=").Lit("asc").Var().ID("SortDescending").ID("sortType").Op("=").Lit("desc"),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("ContextKey").ID("string").Type().ID("sortType").ID("string").Type().ID("Pagination").Struct(jen.ID("Page").ID("uint64"), jen.ID("Limit").ID("uint64"), jen.ID("TotalCount").ID("uint64")).Type().ID("CountResponse").Struct(jen.ID("Count").ID("uint64")),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("error").Op("=").Parens(jen.Op("*").ID("ErrorResponse")).Call(jen.ID("nil")),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("ErrorResponse").Struct(jen.ID("Message").ID("string"), jen.ID("Code").ID("uint")),

		jen.Line(),
	)
	ret.Add(jen.Func().Params(jen.ID("er").Op("*").ID("ErrorResponse")).ID("Error").Params().Params(jen.ID("string")).Block(
		jen.Return().Qual("fmt", "Sprintf").Call(jen.Lit("%d - %s"), jen.ID("er").Dot(
			"Code",
		), jen.ID("er").Dot(
			"Message",
		)),
	),

		jen.Line(),
	)
	return ret
}
