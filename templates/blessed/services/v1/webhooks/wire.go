package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("webhooks")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is our collection of what we provide to other services."),
			jen.ID("Providers").Equals().Qual("github.com/google/wire", "NewSet").Callln(
				jen.ID("ProvideWebhooksService"),
				jen.ID("ProvideWebhookDataManager"),
				jen.ID("ProvideWebhookDataServer"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebhookDataManager is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideWebhookDataManager").Params(jen.ID("db").Qual(proj.DatabaseV1Package(), "Database")).Params(jen.Qual(proj.ModelsV1Package(), "WebhookDataManager")).Block(
			jen.Return().ID("db"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebhookDataServer is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideWebhookDataServer").Params(jen.ID("s").PointerTo().ID("Service")).Params(jen.Qual(proj.ModelsV1Package(), "WebhookDataServer")).Block(
			jen.Return().ID("s"),
		),
		jen.Line(),
	)

	return ret
}
