package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("SortAscending is the pre-determined Ascending sortType for external use"),
			jen.ID("SortAscending").ID("sortType").Equals().Lit("asc"),
			jen.Comment("SortDescending is the pre-determined Descending sortType for external use"),
			jen.ID("SortDescending").ID("sortType").Equals().Lit("desc"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("ContextKey represents strings to be used in Context objects. From the docs:"),
			jen.Comment(`		"The provided key must be comparable and should not be of type string or`),
			jen.Comment(`		any other built-in type to avoid collisions between packages using context."`),
			jen.ID("ContextKey").String(),
			jen.ID("sortType").String(),
			jen.Line(),
			jen.Comment("Pagination represents a pagination request"),
			jen.ID("Pagination").Struct(
				jen.ID("Page").Uint64().Tag(map[string]string{"json": "page"}),
				jen.ID("Limit").Uint64().Tag(map[string]string{"json": "limit"}),
				jen.ID("TotalCount").Uint64().Tag(map[string]string{"json": "total_count"}),
			),
			jen.Line(),
			jen.Comment("CountResponse is what we respond with when a user requests a count of data types"),
			jen.ID("CountResponse").Struct(
				jen.ID("Count").Uint64().Tag(map[string]string{"json": "count"}),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Underscore().Error().Equals().Parens(jen.PointerTo().ID("ErrorResponse")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ErrorResponse represents a response we might send to the user in the event of an error"),
		jen.Line(),
		jen.Type().ID("ErrorResponse").Struct(
			jen.ID("Message").String().Tag(map[string]string{"json": "message"}),
			jen.ID("Code").ID("uint").Tag(map[string]string{"json": "code"}),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("er").PointerTo().ID("ErrorResponse")).ID("Error").Params().Params(jen.String()).Block(
			jen.Return().Qual("fmt", "Sprintf").Call(jen.Lit("%d: %s"), jen.ID("er").Dot("Code"), jen.ID("er").Dot("Message")),
		),
		jen.Line(),
	)

	return ret
}
