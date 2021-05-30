package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhookDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("Webhook").Struct(
				jen.ID("LastUpdatedOn").Op("*").ID("uint64"),
				jen.ID("ArchivedOn").Op("*").ID("uint64"),
				jen.ID("ExternalID").ID("string"),
				jen.ID("Name").ID("string"),
				jen.ID("URL").ID("string"),
				jen.ID("Method").ID("string"),
				jen.ID("ContentType").ID("string"),
				jen.ID("Events").Index().ID("string"),
				jen.ID("DataTypes").Index().ID("string"),
				jen.ID("Topics").Index().ID("string"),
				jen.ID("ID").ID("uint64"),
				jen.ID("CreatedOn").ID("uint64"),
				jen.ID("BelongsToAccount").ID("uint64"),
			),
			jen.ID("WebhookCreationInput").Struct(
				jen.ID("Name").ID("string"),
				jen.ID("ContentType").ID("string"),
				jen.ID("URL").ID("string"),
				jen.ID("Method").ID("string"),
				jen.ID("Events").Index().ID("string"),
				jen.ID("DataTypes").Index().ID("string"),
				jen.ID("Topics").Index().ID("string"),
				jen.ID("BelongsToAccount").ID("uint64"),
			),
			jen.ID("WebhookUpdateInput").Struct(
				jen.ID("Name").ID("string"),
				jen.ID("ContentType").ID("string"),
				jen.ID("URL").ID("string"),
				jen.ID("Method").ID("string"),
				jen.ID("Events").Index().ID("string"),
				jen.ID("DataTypes").Index().ID("string"),
				jen.ID("Topics").Index().ID("string"),
				jen.ID("BelongsToAccount").ID("uint64"),
			),
			jen.ID("WebhookList").Struct(
				jen.ID("Webhooks").Index().Op("*").ID("Webhook"),
				jen.ID("Pagination"),
			),
			jen.ID("WebhookDataManager").Interface(
				jen.ID("GetWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("accountID")).ID("uint64")).Params(jen.Op("*").ID("Webhook"), jen.ID("error")),
				jen.ID("GetAllWebhooksCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")),
				jen.ID("GetAllWebhooks").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("resultChannel").Chan().Index().Op("*").ID("Webhook"), jen.ID("bucketSize").ID("uint16")).Params(jen.ID("error")),
				jen.ID("GetWebhooks").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("filter").Op("*").ID("QueryFilter")).Params(jen.Op("*").ID("WebhookList"), jen.ID("error")),
				jen.ID("CreateWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("WebhookCreationInput"), jen.ID("createdByUser").ID("uint64")).Params(jen.Op("*").ID("Webhook"), jen.ID("error")),
				jen.ID("UpdateWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("Webhook"), jen.ID("changedByUser").ID("uint64"), jen.ID("changes").Index().Op("*").ID("FieldChangeSummary")).Params(jen.ID("error")),
				jen.ID("ArchiveWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("accountID"), jen.ID("archivedByUserID")).ID("uint64")).Params(jen.ID("error")),
				jen.ID("GetAuditLogEntriesForWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("webhookID").ID("uint64")).Params(jen.Index().Op("*").ID("AuditLogEntry"), jen.ID("error")),
			),
			jen.ID("WebhookDataService").Interface(
				jen.ID("AuditEntryHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ListHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("CreateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ReadHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("UpdateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ArchiveHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Update merges an WebhookCreationInput with an Webhook."),
		jen.Line(),
		jen.Func().Params(jen.ID("w").Op("*").ID("Webhook")).ID("Update").Params(jen.ID("input").Op("*").ID("WebhookUpdateInput")).Params(jen.Index().Op("*").ID("FieldChangeSummary")).Body(
			jen.ID("changes").Op(":=").Index().Op("*").ID("FieldChangeSummary").Valuesln(),
			jen.If(jen.ID("input").Dot("Name").Op("!=").Lit("")).Body(
				jen.ID("changes").Op("=").ID("append").Call(
					jen.ID("changes"),
					jen.Op("&").ID("FieldChangeSummary").Valuesln(jen.ID("FieldName").Op(":").Lit("Name"), jen.ID("OldValue").Op(":").ID("w").Dot("Name"), jen.ID("NewValue").Op(":").ID("input").Dot("Name")),
				),
				jen.ID("w").Dot("Name").Op("=").ID("input").Dot("Name"),
			),
			jen.If(jen.ID("input").Dot("ContentType").Op("!=").Lit("")).Body(
				jen.ID("changes").Op("=").ID("append").Call(
					jen.ID("changes"),
					jen.Op("&").ID("FieldChangeSummary").Valuesln(jen.ID("FieldName").Op(":").Lit("ContentType"), jen.ID("OldValue").Op(":").ID("w").Dot("ContentType"), jen.ID("NewValue").Op(":").ID("input").Dot("ContentType")),
				),
				jen.ID("w").Dot("ContentType").Op("=").ID("input").Dot("ContentType"),
			),
			jen.If(jen.ID("input").Dot("URL").Op("!=").Lit("")).Body(
				jen.ID("changes").Op("=").ID("append").Call(
					jen.ID("changes"),
					jen.Op("&").ID("FieldChangeSummary").Valuesln(jen.ID("FieldName").Op(":").Lit("url"), jen.ID("OldValue").Op(":").ID("w").Dot("URL"), jen.ID("NewValue").Op(":").ID("input").Dot("URL")),
				),
				jen.ID("w").Dot("URL").Op("=").ID("input").Dot("URL"),
			),
			jen.If(jen.ID("input").Dot("Method").Op("!=").Lit("")).Body(
				jen.ID("changes").Op("=").ID("append").Call(
					jen.ID("changes"),
					jen.Op("&").ID("FieldChangeSummary").Valuesln(jen.ID("FieldName").Op(":").Lit("Method"), jen.ID("OldValue").Op(":").ID("w").Dot("Method"), jen.ID("NewValue").Op(":").ID("input").Dot("Method")),
				),
				jen.ID("w").Dot("Method").Op("=").ID("input").Dot("Method"),
			),
			jen.If(jen.ID("input").Dot("Events").Op("!=").ID("nil").Op("&&").ID("len").Call(jen.ID("input").Dot("Events")).Op(">").Lit(0)).Body(
				jen.ID("changes").Op("=").ID("append").Call(
					jen.ID("changes"),
					jen.Op("&").ID("FieldChangeSummary").Valuesln(jen.ID("FieldName").Op(":").Lit("Events"), jen.ID("OldValue").Op(":").ID("w").Dot("Events"), jen.ID("NewValue").Op(":").ID("input").Dot("Events")),
				),
				jen.ID("w").Dot("Events").Op("=").ID("input").Dot("Events"),
			),
			jen.If(jen.ID("input").Dot("DataTypes").Op("!=").ID("nil").Op("&&").ID("len").Call(jen.ID("input").Dot("DataTypes")).Op(">").Lit(0)).Body(
				jen.ID("changes").Op("=").ID("append").Call(
					jen.ID("changes"),
					jen.Op("&").ID("FieldChangeSummary").Valuesln(jen.ID("FieldName").Op(":").Lit("DataTypes"), jen.ID("OldValue").Op(":").ID("w").Dot("DataTypes"), jen.ID("NewValue").Op(":").ID("input").Dot("DataTypes")),
				),
				jen.ID("w").Dot("DataTypes").Op("=").ID("input").Dot("DataTypes"),
			),
			jen.If(jen.ID("input").Dot("Topics").Op("!=").ID("nil").Op("&&").ID("len").Call(jen.ID("input").Dot("Topics")).Op(">").Lit(0)).Body(
				jen.ID("changes").Op("=").ID("append").Call(
					jen.ID("changes"),
					jen.Op("&").ID("FieldChangeSummary").Valuesln(jen.ID("FieldName").Op(":").Lit("Topics"), jen.ID("OldValue").Op(":").ID("w").Dot("Topics"), jen.ID("NewValue").Op(":").ID("input").Dot("Topics")),
				),
				jen.ID("w").Dot("Topics").Op("=").ID("input").Dot("Topics"),
			),
			jen.Return().ID("changes"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("WebhookCreationInput")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates a WebhookUpdateInput."),
		jen.Line(),
		jen.Func().Params(jen.ID("w").Op("*").ID("WebhookUpdateInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("w"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("w").Dot("Name"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("w").Dot("URL"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					jen.Op("&").ID("urlValidator").Valuesln(),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("w").Dot("Method"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "In").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Qual("net/http", "MethodPut"),
						jen.Qual("net/http", "MethodPatch"),
						jen.Qual("net/http", "MethodPost"),
						jen.Qual("net/http", "MethodDelete"),
					),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("w").Dot("ContentType"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "In").Call(
						jen.Lit("application/json"),
						jen.Lit("application/xml"),
					),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("w").Dot("Events"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("w").Dot("DataTypes"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("WebhookUpdateInput")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates a WebhookUpdateInput."),
		jen.Line(),
		jen.Func().Params(jen.ID("w").Op("*").ID("WebhookUpdateInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("w"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("w").Dot("Name"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("w").Dot("URL"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					jen.Op("&").ID("urlValidator").Valuesln(),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("w").Dot("Method"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "In").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Qual("net/http", "MethodPut"),
						jen.Qual("net/http", "MethodPatch"),
						jen.Qual("net/http", "MethodPost"),
						jen.Qual("net/http", "MethodDelete"),
					),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("w").Dot("ContentType"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "In").Call(
						jen.Lit("application/json"),
						jen.Lit("application/xml"),
					),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("w").Dot("Events"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("w").Dot("DataTypes"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	return code
}
