package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Defs(
			jen.Comment("Providers are what we provide for dependency injection."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				jen.ID("ProvideOAuth2ClientsService"),
				jen.ID("ProvideOAuth2ClientDataServer"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideOAuth2ClientDataServer is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideOAuth2ClientDataServer").Params(jen.ID("s").PointerTo().ID("Service")).Params(jen.Qual(proj.ModelsV1Package(), "OAuth2ClientDataServer")).Block(
			jen.Return().ID("s"),
		),
		jen.Line(),
	)

	return code
}
