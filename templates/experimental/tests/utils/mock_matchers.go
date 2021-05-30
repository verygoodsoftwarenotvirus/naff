package utils

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockMatchersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("ContextMatcher").Interface().Op("=").ID("mock").Dot("MatchedBy").Call(jen.Func().Params(jen.Qual("context", "Context")).Params(jen.ID("bool")).Body(
				jen.Return().ID("true"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("HTTPRequestMatcher").Interface().Op("=").ID("mock").Dot("MatchedBy").Call(jen.Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("bool")).Body(
				jen.Return().ID("true"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("HTTPResponseWriterMatcher").Interface().Op("=").ID("mock").Dot("MatchedBy").Call(jen.Func().Params(jen.Qual("net/http", "ResponseWriter")).Params(jen.ID("bool")).Body(
				jen.Return().ID("true"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AuditLogEntryCreationInputMatcher is a matcher for use with testify/mock's MatchBy function."),
		jen.Line(),
		jen.Func().ID("AuditLogEntryCreationInputMatcher").Params(jen.ID("eventType").ID("string")).Params(jen.Params(jen.Op("*").ID("types").Dot("AuditLogEntryCreationInput")).Params(jen.ID("bool"))).Body(
			jen.Return().Func().Params(jen.ID("input").Op("*").ID("types").Dot("AuditLogEntryCreationInput")).Params(jen.ID("bool")).Body(
				jen.Return().ID("input").Dot("EventType").Op("==").ID("eventType"))),
		jen.Line(),
	)

	return code
}
