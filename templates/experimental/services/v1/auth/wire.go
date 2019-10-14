package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo() *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideAuthService"), jen.ID("ProvideWebsocketAuthFunc"), jen.ID("ProvideOAuth2ClientValidator")),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideWebsocketAuthFunc provides a WebsocketAuthFunc").ID("ProvideWebsocketAuthFunc").Params(jen.ID("svc").Op("*").ID("Service")).Params(jen.ID("newsman").Dot(
		"WebsocketAuthFunc",
	)).Block(
		jen.Return().ID("svc").Dot(
			"WebsocketAuthFunction",
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideOAuth2ClientValidator converts an oauth2clients.Service to an OAuth2ClientValidator").ID("ProvideOAuth2ClientValidator").Params(jen.ID("s").Op("*").ID("oauth2clients").Dot(
		"Service",
	)).Params(jen.ID("OAuth2ClientValidator")).Block(
		jen.Return().ID("s"),
	),

		jen.Line(),
	)
	return ret
}
