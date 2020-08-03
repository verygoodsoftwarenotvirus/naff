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

	code.Add(buildWireProviders()...)

	// if proj.EnableNewsman {
	code.Add(buildWireProvideWebsocketAuthFunc()...)
	// }

	code.Add(buildWireProvideOAuth2ClientValidator(proj)...)

	return code
}

func buildWireProviders() []jen.Code {
	providersLines := []jen.Code{
		jen.ID("ProvideAuthService"),
	}

	// if proj.EnableNewsman {
	providersLines = append(providersLines, jen.ID("ProvideWebsocketAuthFunc"))
	// }

	providersLines = append(providersLines, jen.ID("ProvideOAuth2ClientValidator"))

	lines := []jen.Code{
		jen.Var().Defs(
			jen.Comment("Providers is our collection of what we provide to other services."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				providersLines...,
			),
		),
		jen.Line(),
	}

	return lines
}

func buildWireProvideWebsocketAuthFunc() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideWebsocketAuthFunc provides a WebsocketAuthFunc."),
		jen.Line(),
		jen.Func().ID("ProvideWebsocketAuthFunc").Params(jen.ID("svc").PointerTo().ID("Service")).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "WebsocketAuthFunc")).Body(
			jen.Return().ID("svc").Dot("WebsocketAuthFunction"),
		),
		jen.Line(),
	}

	return lines
}

func buildWireProvideOAuth2ClientValidator(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideOAuth2ClientValidator converts an oauth2clients.Service to an OAuth2ClientValidator"),
		jen.Line(),
		jen.Func().ID("ProvideOAuth2ClientValidator").Params(jen.ID("s").PointerTo().Qual(proj.ServiceV1OAuth2ClientsPackage(), "Service")).Params(jen.ID("OAuth2ClientValidator")).Body(
			jen.Return().ID("s"),
		),
		jen.Line(),
	}

	return lines
}
