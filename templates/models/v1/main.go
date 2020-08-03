package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("models")

	utils.AddImports(proj, code)

	code.Add(buildMainConstantDefs()...)
	code.Add(buildMainTypeDefs()...)
	code.Add(buildMainErrorInterfaceImplementation()...)
	code.Add(buildMainErrorResponseDotError()...)

	return code
}

func buildMainConstantDefs() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.Comment("SortAscending is the pre-determined Ascending sortType for external use"),
			jen.ID("SortAscending").ID("sortType").Equals().Lit("asc"),
			jen.Comment("SortDescending is the pre-determined Descending sortType for external use"),
			jen.ID("SortDescending").ID("sortType").Equals().Lit("desc"),
		),
		jen.Line(),
	}

	return lines
}

func buildMainTypeDefs() []jen.Code {
	lines := []jen.Code{
		jen.Type().Defs(
			jen.Comment("ContextKey represents strings to be used in Context objects. From the docs:"),
			jen.Comment(`		"The provided key must be comparable and should not be of type string or`),
			jen.Comment(`		any other built-in type to avoid collisions between packages using context."`),
			jen.ID("ContextKey").String(),
			jen.ID("sortType").String(),
			jen.Line(),
			jen.Comment("Pagination represents a pagination request."),
			jen.ID("Pagination").Struct(
				jen.ID("Page").Uint64().Tag(jsonTag("page")),
				jen.ID("Limit").Uint8().Tag(jsonTag("limit")),
			),
			jen.Line(),
			jen.Comment("CountResponse is what we respond with when a user requests a count of data types."),
			jen.ID("CountResponse").Struct(
				jen.ID("Count").Uint64().Tag(jsonTag("count")),
			),
			jen.Line(),
			jen.Comment("ErrorResponse represents a response we might send to the user in the event of an error."),
			jen.ID("ErrorResponse").Struct(
				jen.ID("Message").String().Tag(jsonTag("message")),
				jen.ID("Code").ID("uint").Tag(jsonTag("code")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildMainErrorInterfaceImplementation() []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().Error().Equals().Parens(jen.PointerTo().ID("ErrorResponse")).Call(jen.Nil()),
		jen.Line(),
	}

	return lines
}

func buildMainErrorResponseDotError() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("er").PointerTo().ID("ErrorResponse")).ID("Error").Params().Params(jen.String()).Body(
			jen.Return().Qual("fmt", "Sprintf").Call(jen.Lit("%d: %s"), jen.ID("er").Dot("Code"), jen.ID("er").Dot("Message")),
		),
		jen.Line(),
	}

	return lines
}
