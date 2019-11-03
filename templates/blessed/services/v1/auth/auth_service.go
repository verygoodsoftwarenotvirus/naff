package auth

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authServiceDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("serviceName").Op("=").Lit("auth_service"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("OAuth2ClientValidator is a stand-in interface, where we needed to abstract"),
			jen.Comment("a regular structure with an interface for testing purposes"),
			jen.ID("OAuth2ClientValidator").Interface(
				jen.ID("ExtractOAuth2ClientFromRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"), jen.ID("error")),
			),
			jen.Line(),
			jen.Comment("cookieEncoderDecoder is a stand-in interface for gorilla/securecookie"),
			jen.ID("cookieEncoderDecoder").Interface(
				jen.ID("Encode").Params(jen.ID("name").ID("string"), jen.ID("value").Interface()).Params(jen.ID("string"), jen.ID("error")),
				jen.ID("Decode").Params(jen.List(jen.ID("name"), jen.ID("value")).ID("string"), jen.ID("dst").Interface()).Params(jen.ID("error")),
			),
			jen.Line(),
			jen.Comment("UserIDFetcher is a function that fetches user IDs"),
			jen.ID("UserIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
			jen.Line(),
			jen.Comment("Service handles authentication service-wide"),
			jen.ID("Service").Struct(
				jen.ID("config").Qual(filepath.Join(pkgRoot, "internal/v1/config"), "AuthSettings"),
				jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
				jen.ID("authenticator").Qual(filepath.Join(pkgRoot, "internal/v1/auth"), "Authenticator"),
				jen.ID("userIDFetcher").ID("UserIDFetcher"),
				jen.ID("userDB").Qual(filepath.Join(pkgRoot, "models/v1"), "UserDataManager"),
				jen.ID("oauth2ClientsService").ID("OAuth2ClientValidator"),
				jen.ID("encoderDecoder").Qual(filepath.Join(pkgRoot, "internal/v1/encoding"), "EncoderDecoder"),
				jen.ID("cookieManager").ID("cookieEncoderDecoder"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideAuthService builds a new AuthService"),
		jen.Line(),
		jen.Func().ID("ProvideAuthService").Paramsln(
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
			jen.ID("cfg").Op("*").Qual(filepath.Join(pkgRoot, "internal/v1/config"), "ServerConfig"),
			jen.ID("authenticator").Qual(filepath.Join(pkgRoot, "internal/v1/auth"), "Authenticator"),
			jen.ID("database").Qual(filepath.Join(pkgRoot, "models/v1"), "UserDataManager"),
			jen.ID("oauth2ClientsService").ID("OAuth2ClientValidator"),
			jen.ID("userIDFetcher").ID("UserIDFetcher"),
			jen.ID("encoder").Qual(filepath.Join(pkgRoot, "internal/v1/encoding"), "EncoderDecoder"),
		).Params(jen.Op("*").ID("Service")).Block(
			jen.ID("svc").Op(":=").Op("&").ID("Service").Valuesln(
				jen.ID("logger").Op(":").ID("logger").Dot("WithName").Call(jen.ID("serviceName")),
				jen.ID("encoderDecoder").Op(":").ID("encoder"),
				jen.ID("config").Op(":").ID("cfg").Dot("Auth"),
				jen.ID("userDB").Op(":").ID("database"),
				jen.ID("oauth2ClientsService").Op(":").ID("oauth2ClientsService"),
				jen.ID("authenticator").Op(":").ID("authenticator"),
				jen.ID("userIDFetcher").Op(":").ID("userIDFetcher"),
				jen.ID("cookieManager").Op(":").ID("securecookie").Dot("New").Callln(
					jen.ID("securecookie").Dot("GenerateRandomKey").Call(jen.Lit(64)),
					jen.Index().ID("byte").Call(jen.ID("cfg").Dot("Auth").Dot("CookieSecret")),
				),
			),
			jen.Line(),
			jen.Return().ID("svc"),
		),
		jen.Line(),
	)
	return ret
}
