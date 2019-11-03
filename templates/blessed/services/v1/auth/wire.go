package auth

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is our collection of what we provide to other services"),
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Callln(
				jen.ID("ProvideAuthService"),
				jen.ID("ProvideWebsocketAuthFunc"),
				jen.ID("ProvideOAuth2ClientValidator"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebsocketAuthFunc provides a WebsocketAuthFunc"),
		jen.Line(),
		jen.Func().ID("ProvideWebsocketAuthFunc").Params(jen.ID("svc").Op("*").ID("Service")).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "WebsocketAuthFunc")).Block(
			jen.Return().ID("svc").Dot("WebsocketAuthFunction"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideOAuth2ClientValidator converts an oauth2clients.Service to an OAuth2ClientValidator"),
		jen.Line(),
		jen.Func().ID("ProvideOAuth2ClientValidator").Params(jen.ID("s").Op("*").Qual(filepath.Join(pkgRoot, "services/v1/oauth2clients"), "Service")).Params(jen.ID("OAuth2ClientValidator")).Block(
			jen.Return().ID("s"),
		),
		jen.Line(),
	)
	return ret
}
