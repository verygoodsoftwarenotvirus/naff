package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(proj, ret)

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

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is our collection of what we provide to other services"),
			jen.ID("Providers").Equals().Qual("github.com/google/wire", "NewSet").Callln(
				buildProviderSet()...,
			),
		),
		jen.Line(),
	)

	// if proj.EnableNewsman {
	ret.Add(
		jen.Comment("ProvideWebsocketAuthFunc provides a WebsocketAuthFunc"),
		jen.Line(),
		jen.Func().ID("ProvideWebsocketAuthFunc").Params(jen.ID("svc").PointerTo().ID("Service")).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "WebsocketAuthFunc")).Block(
			jen.Return().ID("svc").Dot("WebsocketAuthFunction"),
		),
		jen.Line(),
	)
	// }

	ret.Add(
		jen.Comment("ProvideOAuth2ClientValidator converts an oauth2clients.Service to an OAuth2ClientValidator"),
		jen.Line(),
		jen.Func().ID("ProvideOAuth2ClientValidator").Params(jen.ID("s").PointerTo().Qual(proj.ServiceV1OAuth2ClientsPackage(), "Service")).Params(jen.ID("OAuth2ClientValidator")).Block(
			jen.Return().ID("s"),
		),
		jen.Line(),
	)
	return ret
}
