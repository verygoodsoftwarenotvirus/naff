package apiclients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildWireVarDefs()...)
	code.Add(buildProvideOAuth2ClientDataServer(proj)...)

	return code
}

func buildWireVarDefs() []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.Comment("Providers are what we provide for dependency injection."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				jen.ID("ProvideOAuth2ClientsService"),
				jen.ID("ProvideOAuth2ClientDataServer"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideOAuth2ClientDataServer(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideOAuth2ClientDataServer is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideOAuth2ClientDataServer").Params(jen.ID("s").PointerTo().ID("Service")).Params(jen.Qual(proj.TypesPackage(), "OAuth2ClientDataServer")).Body(
			jen.Return().ID("s"),
		),
		jen.Line(),
	}

	return lines
}
