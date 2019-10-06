package models

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func mainDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("SortAscending").ID("sortType").Op("=").Lit("asc").Var().ID("SortDescending").ID("sortType").Op("=").Lit("desc"))
	ret.Add(jen.Null().Type().ID("ContextKey").ID("string").Type().ID("sortType").ID("string").Type().ID("Pagination").Struct(
		jen.ID("Page").ID("uint64"),
		jen.ID("Limit").ID("uint64"),
		jen.ID("TotalCount").ID("uint64"),
	).Type().ID("CountResponse").Struct(
		jen.ID("Count").ID("uint64"),
	),
	)
	ret.Add(jen.Null().Var().ID("_").ID("error").Op("=").Parens(jen.Op("*").ID("ErrorResponse")).Call(jen.ID("nil")))
	ret.Add(jen.Null().Type().ID("ErrorResponse").Struct(
		jen.ID("Message").ID("string"),
		jen.ID("Code").ID("uint"),
	),
	)
	ret.Add(jen.Func().Params(jen.ID("er").Op("*").ID("ErrorResponse")).ID("Error").Params().Params(jen.ID("string")).Block(
		jen.Return().Qual("fmt", "Sprintf").Call(jen.Lit("%d - %s"), jen.ID("er").Dot(
			"Code",
		), jen.ID("er").Dot(
			"Message",
		)),
	),
	)
	return ret
}
