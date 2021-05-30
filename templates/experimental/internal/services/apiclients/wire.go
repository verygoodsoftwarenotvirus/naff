package apiclients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Call(
				jen.ID("ProvideConfig"),
				jen.ID("ProvideAPIClientsService"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideConfig converts an auth config to a local config."),
		jen.Line(),
		jen.Func().ID("ProvideConfig").Params(jen.ID("cfg").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/authentication", "Config")).Params(jen.Op("*").ID("config")).Body(
			jen.Return().Op("&").ID("config").Valuesln(jen.ID("minimumUsernameLength").Op(":").ID("cfg").Dot("MinimumUsernameLength"), jen.ID("minimumPasswordLength").Op(":").ID("cfg").Dot("MinimumPasswordLength"))),
		jen.Line(),
	)

	return code
}
