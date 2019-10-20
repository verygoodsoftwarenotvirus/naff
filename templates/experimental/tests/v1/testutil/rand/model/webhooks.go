package model

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func webhooksDotGo() *jen.File {
	ret := jen.NewFile("randmodel")

	utils.AddImports(ret)

	ret.Add(
		jen.Comment("RandomWebhookInput creates a random WebhookCreationInput"),
		jen.Line(),
		jen.Func().ID("RandomWebhookInput").Params().Params(jen.Op("*").ID("models").Dot(
		"WebhookCreationInput",
	)).Block(
		jen.ID("x").Op(":=").Op("&").ID("models").Dot(
			"WebhookCreationInput",
		).Valuesln(
	jen.ID("Name").Op(":").ID("fake").Dot(
			"Word",
		).Call(), jen.ID("URL").Op(":").ID("fake").Dot(
			"DomainName",
		).Call(), jen.ID("ContentType").Op(":").Lit("application/json"), jen.ID("Method").Op(":").Lit("POST")),
		jen.Return().ID("x"),
	),
	jen.Line(),
	)
	return ret
}
