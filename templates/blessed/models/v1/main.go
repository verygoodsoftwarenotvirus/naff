package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("SortAscending is the pre-determined Ascending sortType for external use"),
			jen.ID("SortAscending").ID("sortType").Op("=").Lit("asc"),
			jen.Comment("SortDescending is the pre-determined Descending sortType for external use"),
			jen.ID("SortDescending").ID("sortType").Op("=").Lit("desc"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("ContextKey represents strings to be used in Context objects. From the docs:"),
			jen.Comment(`		"The provided key must be comparable and should not be of type string or`),
			jen.Comment(`		any other built-in type to avoid collisions between packages using context."`),
			jen.ID("ContextKey").ID("string"),
			jen.ID("sortType").ID("string"),
			jen.Line(),
			jen.Comment("Pagination represents a pagination request"),
			jen.ID("Pagination").Struct(
				jen.ID("Page").ID("uint64").Tag(map[string]string{"json": "page"}),
				jen.ID("Limit").ID("uint64").Tag(map[string]string{"json": "limit"}),
				jen.ID("TotalCount").ID("uint64").Tag(map[string]string{"json": "total_count"}),
			),
			jen.Line(),
			jen.Comment("CountResponse is what we respond with when a user requests a count of data types"),
			jen.ID("CountResponse").Struct(
				jen.ID("Count").ID("uint64").Tag(map[string]string{"json": "count"}),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").ID("error").Op("=").Parens(jen.Op("*").ID("ErrorResponse")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ErrorResponse represents a response we might send to the user in the event of an error"),
		jen.Line(),
		jen.Type().ID("ErrorResponse").Struct(
			jen.ID("Message").ID("string").Tag(map[string]string{"json": "message"}),
			jen.ID("Code").ID("uint").Tag(map[string]string{"json": "code"}),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("er").Op("*").ID("ErrorResponse")).ID("Error").Params().Params(jen.ID("string")).Block(
			jen.Return().Qual("fmt", "Sprintf").Call(jen.Lit("%d: %s"), jen.ID("er").Dot("Code"), jen.ID("er").Dot("Message")),
		),
		jen.Line(),
	)
	return ret
}
