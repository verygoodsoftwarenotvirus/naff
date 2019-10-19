package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func authServiceDotGo() *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("serviceName").Op("=").Lit("auth_service"),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("OAuth2ClientValidator").Interface(jen.ID("ExtractOAuth2ClientFromRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("models").Dot(
			"OAuth2Client",
	),
	jen.ID("error"))).Type().ID("cookieEncoderDecoder").Interface(jen.ID("Encode").Params(jen.ID("name").ID("string"), jen.ID("value").Interface()).Params(jen.ID("string"), jen.ID("error")), jen.ID("Decode").Params(jen.List(jen.ID("name"), jen.ID("value")).ID("string"), jen.ID("dst").Interface()).Params(jen.ID("error"))).Type().ID("UserIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Type().ID("Service").Struct(jen.ID("config").ID("config").Dot(
			"AuthSettings",
	),
	jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1",
			"Logger",
	),
	jen.ID("authenticator").ID("auth").Dot(
			"Authenticator",
	),
	jen.ID("userIDFetcher").ID("UserIDFetcher"), jen.ID("userDB").ID("models").Dot(
			"UserDataManager",
	),
	jen.ID("oauth2ClientsService").ID("OAuth2ClientValidator"), jen.ID("encoderDecoder").ID("encoding").Dot(
			"EncoderDecoder",
	),
	jen.ID("cookieManager").ID("cookieEncoderDecoder")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideAuthService builds a new AuthService"),
		jen.Line(),
		jen.Func().ID("ProvideAuthService").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1",
			"Logger",
	),
	jen.ID("cfg").Op("*").ID("config").Dot(
			"ServerConfig",
	),
	jen.ID("authenticator").ID("auth").Dot(
			"Authenticator",
	),
	jen.ID("database").ID("models").Dot(
			"UserDataManager",
	),
	jen.ID("oauth2ClientsService").ID("OAuth2ClientValidator"), jen.ID("userIDFetcher").ID("UserIDFetcher"), jen.ID("encoder").ID("encoding").Dot(
			"EncoderDecoder",
		)).Params(jen.Op("*").ID("Service")).Block(
			jen.ID("svc").Op(":=").Op("&").ID("Service").Valuesln(jen.ID("logger").Op(":").ID("logger").Dot(
				"WithName",
			).Call(jen.ID("serviceName")), jen.ID("encoderDecoder").Op(":").ID("encoder"), jen.ID("config").Op(":").ID("cfg").Dot(
				"Auth",
	),
	jen.ID("userDB").Op(":").ID("database"), jen.ID("oauth2ClientsService").Op(":").ID("oauth2ClientsService"), jen.ID("authenticator").Op(":").ID("authenticator"), jen.ID("userIDFetcher").Op(":").ID("userIDFetcher"), jen.ID("cookieManager").Op(":").ID("securecookie").Dot(
				"New",
			).Call(jen.ID("securecookie").Dot(
				"GenerateRandomKey",
			).Call(jen.Lit(64)), jen.Index().ID("byte").Call(jen.ID("cfg").Dot(
				"Auth",
			).Dot(
				"CookieSecret",
			)))),
			jen.Return().ID("svc"),
		),
		jen.Line(),
	)
	return ret
}
