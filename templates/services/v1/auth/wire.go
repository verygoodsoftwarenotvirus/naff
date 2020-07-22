package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	buildProviderSet := func() []jen.Code {
		lines := []jen.Code{
			jen.ID("ProvideAuthService"),
		}

		// if proj.EnableNewsman {
		lines = append(lines, jen.ID("ProvideWebsocketAuthFunc"))
		// }

		lines = append(lines, jen.ID("ProvideOAuth2ClientValidator"))

		return lines
	}

	code.Add(
		jen.Var().Defs(
			jen.Comment("Providers is our collection of what we provide to other services."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				buildProviderSet()...,
			),
		),
		jen.Line(),
	)

	// if proj.EnableNewsman {
	code.Add(
		jen.Comment("ProvideWebsocketAuthFunc provides a WebsocketAuthFunc."),
		jen.Line(),
		jen.Func().ID("ProvideWebsocketAuthFunc").Params(jen.ID("svc").PointerTo().ID("Service")).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "WebsocketAuthFunc")).Block(
			jen.Return().ID("svc").Dot("WebsocketAuthFunction"),
		),
		jen.Line(),
	)
	// }

	code.Add(
		jen.Comment("ProvideOAuth2ClientValidator converts an oauth2clients.Service to an OAuth2ClientValidator"),
		jen.Line(),
		jen.Func().ID("ProvideOAuth2ClientValidator").Params(jen.ID("s").PointerTo().Qual(proj.ServiceV1OAuth2ClientsPackage(), "Service")).Params(jen.ID("OAuth2ClientValidator")).Block(
			jen.Return().ID("s"),
		),
		jen.Line(),
	)

	return code
}
