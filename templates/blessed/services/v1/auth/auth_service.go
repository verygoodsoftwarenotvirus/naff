package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authServiceDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("serviceName").Equals().Lit("auth_service"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("OAuth2ClientValidator is a stand-in interface, where we needed to abstract"),
			jen.Comment("a regular structure with an interface for testing purposes"),
			jen.ID("OAuth2ClientValidator").Interface(
				jen.ID("ExtractOAuth2ClientFromRequest").Params(utils.CtxParam(), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.ModelsV1Package(), "OAuth2Client"), jen.ID("error")),
			),
			jen.Line(),
			jen.Comment("cookieEncoderDecoder is a stand-in interface for gorilla/securecookie"),
			jen.ID("cookieEncoderDecoder").Interface(
				jen.ID("Encode").Params(jen.ID("name").ID("string"), jen.ID("value").Interface()).Params(jen.ID("string"), jen.ID("error")),
				jen.ID("Decode").Params(jen.List(jen.ID("name"), jen.ID("value")).ID("string"), jen.ID("dst").Interface()).Params(jen.ID("error")),
			),
			jen.Line(),
			jen.Comment("UserIDFetcher is a function that fetches user IDs"),
			jen.ID("UserIDFetcher").Func().Params(jen.ParamPointer().Qual("net/http", "Request")).Params(jen.ID("uint64")),
			jen.Line(),
			jen.Comment("Service handles authentication service-wide"),
			jen.ID("Service").Struct(
				jen.ID("config").Qual(proj.InternalConfigV1Package(), "AuthSettings"),
				jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
				jen.ID("authenticator").Qual(proj.InternalAuthV1Package(), "Authenticator"),
				jen.ID("userIDFetcher").ID("UserIDFetcher"),
				jen.ID("userDB").Qual(proj.ModelsV1Package(), "UserDataManager"),
				jen.ID("oauth2ClientsService").ID("OAuth2ClientValidator"),
				jen.ID("encoderDecoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
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
			jen.ID("cfg").Op("*").Qual(proj.InternalConfigV1Package(), "ServerConfig"),
			jen.ID("authenticator").Qual(proj.InternalAuthV1Package(), "Authenticator"),
			jen.ID("database").Qual(proj.ModelsV1Package(), "UserDataManager"),
			jen.ID("oauth2ClientsService").ID("OAuth2ClientValidator"),
			jen.ID("userIDFetcher").ID("UserIDFetcher"),
			jen.ID("encoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
		).Params(jen.Op("*").ID("Service")).Block(
			jen.ID("svc").Assign().VarPointer().ID("Service").Valuesln(
				jen.ID("logger").MapAssign().ID("logger").Dot("WithName").Call(jen.ID("serviceName")),
				jen.ID("encoderDecoder").MapAssign().ID("encoder"),
				jen.ID("config").MapAssign().ID("cfg").Dot("Auth"),
				jen.ID("userDB").MapAssign().ID("database"),
				jen.ID("oauth2ClientsService").MapAssign().ID("oauth2ClientsService"),
				jen.ID("authenticator").MapAssign().ID("authenticator"),
				jen.ID("userIDFetcher").MapAssign().ID("userIDFetcher"),
				jen.ID("cookieManager").MapAssign().Qual("github.com/gorilla/securecookie", "New").Callln(
					jen.Qual("github.com/gorilla/securecookie", "GenerateRandomKey").Call(jen.Lit(64)),
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
