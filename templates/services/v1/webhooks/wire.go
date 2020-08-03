package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildProviders()...)
	code.Add(buildProvideWebhookDataManager(proj)...)
	code.Add(buildProvideWebhookDataServer(proj)...)

	return code
}

func buildProviders() []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.Comment("Providers is our collection of what we provide to other services."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				jen.ID("ProvideWebhooksService"),
				jen.ID("ProvideWebhookDataManager"),
				jen.ID("ProvideWebhookDataServer"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideWebhookDataManager(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideWebhookDataManager is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideWebhookDataManager").Params(jen.ID("db").Qual(proj.DatabaseV1Package(), "DataManager")).Params(jen.Qual(proj.ModelsV1Package(), "WebhookDataManager")).Body(
			jen.Return().ID("db"),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideWebhookDataServer(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideWebhookDataServer is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideWebhookDataServer").Params(jen.ID("s").PointerTo().ID("Service")).Params(jen.Qual(proj.ModelsV1Package(), "WebhookDataServer")).Body(
			jen.Return().ID("s"),
		),
		jen.Line(),
	}

	return lines
}
