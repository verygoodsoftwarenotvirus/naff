package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("SortAscending").ID("sortType").Op("=").Lit("asc"),
			jen.ID("SortDescending").ID("sortType").Op("=").Lit("desc"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("ContextKey").ID("string"),
			jen.ID("sortType").ID("string"),
			jen.ID("Pagination").Struct(
				jen.ID("Page").ID("uint64"),
				jen.ID("Limit").ID("uint8"),
				jen.ID("FilteredCount").ID("uint64"),
				jen.ID("TotalCount").ID("uint64"),
			),
			jen.ID("ErrorResponse").Struct(
				jen.ID("Message").ID("string"),
				jen.ID("Code").ID("int"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("error").Op("=").Parens(jen.Op("*").ID("ErrorResponse")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("er").Op("*").ID("ErrorResponse")).ID("Error").Params().Params(jen.ID("string")).Body(
			jen.Return().Qual("fmt", "Sprintf").Call(
				jen.Lit("%d: %s"),
				jen.ID("er").Dot("Code"),
				jen.ID("er").Dot("Message"),
			)),
		jen.Line(),
	)

	return code
}
