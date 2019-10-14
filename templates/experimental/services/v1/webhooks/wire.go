package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo() *jen.File {
	ret := jen.NewFile("webhooks")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideWebhooksService"), jen.ID("ProvideWebhookDataManager"), jen.ID("ProvideWebhookDataServer")),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideWebhookDataManager is an arbitrary function for dependency injection's sake").ID("ProvideWebhookDataManager").Params(jen.ID("db").ID("database").Dot(
		"Database",
	)).Params(jen.ID("models").Dot(
		"WebhookDataManager",
	)).Block(
		jen.Return().ID("db"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideWebhookDataServer is an arbitrary function for dependency injection's sake").ID("ProvideWebhookDataServer").Params(jen.ID("s").Op("*").ID("Service")).Params(jen.ID("models").Dot(
		"WebhookDataServer",
	)).Block(
		jen.Return().ID("s"),
	),

		jen.Line(),
	)
	return ret
}
