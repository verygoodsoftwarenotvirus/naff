package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func authServiceDotGo() *jen.File {
	ret := jen.NewFile("auth")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("serviceName").Op("=").Lit("auth_service"))
	ret.Add(jen.Null().Type().ID("OAuth2ClientValidator").Interface(jen.ID("ExtractOAuth2ClientFromRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("models").Dot(
		"OAuth2Client",
	), jen.ID("error"))).Type().ID("cookieEncoderDecoder").Interface(jen.ID("Encode").Params(jen.ID("name").ID("string"), jen.ID("value").Interface()).Params(jen.ID("string"), jen.ID("error")), jen.ID("Decode").Params(jen.List(jen.ID("name"), jen.ID("value")).ID("string"), jen.ID("dst").Interface()).Params(jen.ID("error"))).Type().ID("UserIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Type().ID("Service").Struct(
		jen.ID("config").ID("config").Dot(
			"AuthSettings",
		),
		jen.ID("logger").ID("logging").Dot(
			"Logger",
		),
		jen.ID("authenticator").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1", "Authenticator"),
		jen.ID("userIDFetcher").ID("UserIDFetcher"),
		jen.ID("userDB").ID("models").Dot(
			"UserDataManager",
		),
		jen.ID("oauth2ClientsService").ID("OAuth2ClientValidator"),
		jen.ID("encoderDecoder").ID("encoding").Dot(
			"EncoderDecoder",
		),
		jen.ID("cookieManager").ID("cookieEncoderDecoder"),
	),
	)
	ret.Add(jen.Func())
	return ret
}
