package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockWebhookSQLQueryBuilderDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("querybuilding").Dot("WebhookSQLQueryBuilder").Op("=").Parens(jen.Op("*").ID("WebhookSQLQueryBuilder")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("WebhookSQLQueryBuilder").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetWebhookQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookSQLQueryBuilder")).ID("BuildGetWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("webhookID"),
				jen.ID("accountID"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAllWebhooksCountQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookSQLQueryBuilder")).ID("BuildGetAllWebhooksCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("ctx")).Dot("String").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetBatchOfWebhooksQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookSQLQueryBuilder")).ID("BuildGetBatchOfWebhooksQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("beginID"),
				jen.ID("endID"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetWebhooksQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookSQLQueryBuilder")).ID("BuildGetWebhooksQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("filter"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildCreateWebhookQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookSQLQueryBuilder")).ID("BuildCreateWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("x").Op("*").ID("types").Dot("WebhookCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("x"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUpdateWebhookQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookSQLQueryBuilder")).ID("BuildUpdateWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("Webhook")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildArchiveWebhookQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookSQLQueryBuilder")).ID("BuildArchiveWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("webhookID"),
				jen.ID("accountID"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAuditLogEntriesForWebhookQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookSQLQueryBuilder")).ID("BuildGetAuditLogEntriesForWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("webhookID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("webhookID"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	return code
}
