package webhooks

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("webhooks")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is our collection of what we provide to other services"),
			jen.ID("Providers").Op("=").Qual("github.com/google/wire", "NewSet".Callln(
				jen.ID("ProvideWebhooksService"),
				jen.ID("ProvideWebhookDataManager"),
				jen.ID("ProvideWebhookDataServer"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebhookDataManager is an arbitrary function for dependency injection's sake"),
		jen.Line(),
		jen.Func().ID("ProvideWebhookDataManager").Params(jen.ID("db").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Database")).Params(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookDataManager")).Block(
			jen.Return().ID("db"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebhookDataServer is an arbitrary function for dependency injection's sake"),
		jen.Line(),
		jen.Func().ID("ProvideWebhookDataServer").Params(jen.ID("s").Op("*").ID("Service")).Params(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookDataServer")).Block(
			jen.Return().ID("s"),
		),
		jen.Line(),
	)
	return ret
}
