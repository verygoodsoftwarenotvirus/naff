package fakes

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhookDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("BuildFakeWebhook builds a faked Webhook."),
		jen.Line(),
		jen.Func().ID("BuildFakeWebhook").Params().Params(jen.Op("*").ID("types").Dot("Webhook")).Body(
			jen.Return().Op("&").ID("types").Dot("Webhook").Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Uint32").Call()), jen.ID("ExternalID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "UUID").Call(), jen.ID("Name").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Word").Call(), jen.ID("ContentType").Op(":").Lit("application/json"), jen.ID("URL").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "URL").Call(), jen.ID("Method").Op(":").Qual("net/http", "MethodPost"), jen.ID("Events").Op(":").Index().ID("string").Valuesln(jen.Qual("github.com/brianvoe/gofakeit/v5", "Word").Call()), jen.ID("DataTypes").Op(":").Index().ID("string").Valuesln(jen.Qual("github.com/brianvoe/gofakeit/v5", "Word").Call()), jen.ID("Topics").Op(":").Index().ID("string").Valuesln(jen.Qual("github.com/brianvoe/gofakeit/v5", "Word").Call()), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.ID("uint32").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Date").Call().Dot("Unix").Call())), jen.ID("ArchivedOn").Op(":").ID("nil"), jen.ID("BelongsToAccount").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call())),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeWebhookList builds a faked WebhookList."),
		jen.Line(),
		jen.Func().ID("BuildFakeWebhookList").Params().Params(jen.Op("*").ID("types").Dot("WebhookList")).Body(
			jen.Var().Defs(
				jen.ID("examples").Index().Op("*").ID("types").Dot("Webhook"),
			),
			jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("exampleQuantity"), jen.ID("i").Op("++")).Body(
				jen.ID("examples").Op("=").ID("append").Call(
					jen.ID("examples"),
					jen.ID("BuildFakeWebhook").Call(),
				)),
			jen.Return().Op("&").ID("types").Dot("WebhookList").Valuesln(jen.ID("Pagination").Op(":").ID("types").Dot("Pagination").Valuesln(jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("FilteredCount").Op(":").ID("exampleQuantity").Op("/").Lit(2), jen.ID("TotalCount").Op(":").ID("exampleQuantity")), jen.ID("Webhooks").Op(":").ID("examples")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeWebhookUpdateInput builds a faked WebhookUpdateInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeWebhookUpdateInput").Params().Params(jen.Op("*").ID("types").Dot("WebhookUpdateInput")).Body(
			jen.ID("webhook").Op(":=").ID("BuildFakeWebhook").Call(),
			jen.Return().Op("&").ID("types").Dot("WebhookUpdateInput").Valuesln(jen.ID("Name").Op(":").ID("webhook").Dot("Name"), jen.ID("ContentType").Op(":").ID("webhook").Dot("ContentType"), jen.ID("URL").Op(":").ID("webhook").Dot("URL"), jen.ID("Method").Op(":").ID("webhook").Dot("Method"), jen.ID("Events").Op(":").ID("webhook").Dot("Events"), jen.ID("DataTypes").Op(":").ID("webhook").Dot("DataTypes"), jen.ID("Topics").Op(":").ID("webhook").Dot("Topics"), jen.ID("BelongsToAccount").Op(":").ID("webhook").Dot("BelongsToAccount")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeWebhookUpdateInputFromWebhook builds a faked WebhookUpdateInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeWebhookUpdateInputFromWebhook").Params(jen.ID("webhook").Op("*").ID("types").Dot("Webhook")).Params(jen.Op("*").ID("types").Dot("WebhookUpdateInput")).Body(
			jen.Return().Op("&").ID("types").Dot("WebhookUpdateInput").Valuesln(jen.ID("Name").Op(":").ID("webhook").Dot("Name"), jen.ID("ContentType").Op(":").ID("webhook").Dot("ContentType"), jen.ID("URL").Op(":").ID("webhook").Dot("URL"), jen.ID("Method").Op(":").ID("webhook").Dot("Method"), jen.ID("Events").Op(":").ID("webhook").Dot("Events"), jen.ID("DataTypes").Op(":").ID("webhook").Dot("DataTypes"), jen.ID("Topics").Op(":").ID("webhook").Dot("Topics"), jen.ID("BelongsToAccount").Op(":").ID("webhook").Dot("BelongsToAccount"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeWebhookCreationInput builds a faked WebhookCreationInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeWebhookCreationInput").Params().Params(jen.Op("*").ID("types").Dot("WebhookCreationInput")).Body(
			jen.ID("webhook").Op(":=").ID("BuildFakeWebhook").Call(),
			jen.Return().ID("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("webhook")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeWebhookCreationInputFromWebhook builds a faked WebhookCreationInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeWebhookCreationInputFromWebhook").Params(jen.ID("webhook").Op("*").ID("types").Dot("Webhook")).Params(jen.Op("*").ID("types").Dot("WebhookCreationInput")).Body(
			jen.Return().Op("&").ID("types").Dot("WebhookCreationInput").Valuesln(jen.ID("Name").Op(":").ID("webhook").Dot("Name"), jen.ID("ContentType").Op(":").ID("webhook").Dot("ContentType"), jen.ID("URL").Op(":").ID("webhook").Dot("URL"), jen.ID("Method").Op(":").ID("webhook").Dot("Method"), jen.ID("Events").Op(":").ID("webhook").Dot("Events"), jen.ID("DataTypes").Op(":").ID("webhook").Dot("DataTypes"), jen.ID("Topics").Op(":").ID("webhook").Dot("Topics"), jen.ID("BelongsToAccount").Op(":").ID("webhook").Dot("BelongsToAccount"))),
		jen.Line(),
	)

	return code
}
