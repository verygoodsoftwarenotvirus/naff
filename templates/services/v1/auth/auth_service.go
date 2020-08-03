package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authServiceDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildAuthServiceConstantDefs()...)
	code.Add(buildAuthServiceTypeDefs(proj)...)
	code.Add(buildProvideAuthService(proj)...)

	return code
}

func buildAuthServiceConstantDefs() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.ID("serviceName").Equals().Lit("auth_service"),
		),
		jen.Line(),
	}

	return lines
}

func buildAuthServiceTypeDefs(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Type().Defs(
			jen.Comment("OAuth2ClientValidator is a stand-in interface, where we needed to abstract"),
			jen.Comment("a regular structure with an interface for testing purposes."),
			jen.ID("OAuth2ClientValidator").Interface(
				jen.ID("ExtractOAuth2ClientFromRequest").Params(constants.CtxParam(), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"), jen.Error()),
			),
			jen.Line(),
			jen.Comment("cookieEncoderDecoder is a stand-in interface for gorilla/securecookie"),
			jen.ID("cookieEncoderDecoder").Interface(
				jen.ID("Encode").Params(jen.ID("name").String(), jen.ID("value").Interface()).Params(jen.String(), jen.Error()),
				jen.ID("Decode").Params(jen.List(jen.ID("name"), jen.ID("value")).String(), jen.ID("dst").Interface()).Params(jen.Error()),
			),
			jen.Line(),
			jen.Comment("Service handles authentication service-wide"),
			jen.ID("Service").Struct(
				jen.ID("config").Qual(proj.InternalConfigV1Package(), "AuthSettings"),
				constants.LoggerParam(),
				jen.ID("authenticator").Qual(proj.InternalAuthV1Package(), "Authenticator"),
				jen.ID("userDB").Qual(proj.ModelsV1Package(), "UserDataManager"),
				jen.ID("oauth2ClientsService").ID("OAuth2ClientValidator"),
				jen.ID("encoderDecoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
				jen.ID("cookieManager").ID("cookieEncoderDecoder"),
				jen.ID("sessionManager").PointerTo().Qual(constants.SessionManagerLibrary, "SessionManager"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideAuthService(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideAuthService builds a new AuthService."),
		jen.Line(),
		jen.Func().ID("ProvideAuthService").Paramsln(
			constants.LoggerParam(),
			jen.ID("cfg").Qual(proj.InternalConfigV1Package(), "AuthSettings"),
			jen.ID("authenticator").Qual(proj.InternalAuthV1Package(), "Authenticator"),
			jen.ID("database").Qual(proj.ModelsV1Package(), "UserDataManager"),
			jen.ID("oauth2ClientsService").ID("OAuth2ClientValidator"),
			jen.ID("sessionManager").PointerTo().Qual(constants.SessionManagerLibrary, "SessionManager"),
			jen.ID("encoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
		).Params(jen.PointerTo().ID("Service"), jen.Error()).Body(
			jen.ID("svc").Assign().AddressOf().ID("Service").Valuesln(
				jen.ID(constants.LoggerVarName).MapAssign().ID(constants.LoggerVarName).Dot("WithName").Call(jen.ID("serviceName")),
				jen.ID("encoderDecoder").MapAssign().ID("encoder"),
				jen.ID("config").MapAssign().ID("cfg"),
				jen.ID("userDB").MapAssign().ID("database"),
				jen.ID("oauth2ClientsService").MapAssign().ID("oauth2ClientsService"),
				jen.ID("authenticator").MapAssign().ID("authenticator"),
				jen.ID("sessionManager").MapAssign().ID("sessionManager"),
				jen.ID("cookieManager").MapAssign().Qual("github.com/gorilla/securecookie", "New").Callln(
					jen.Qual("github.com/gorilla/securecookie", "GenerateRandomKey").Call(jen.Lit(64)),
					jen.Index().Byte().Call(jen.ID("cfg").Dot("CookieSecret")),
				),
			),
			jen.ID("svc").Dot("sessionManager").Dot("Lifetime").Equals().ID("cfg").Dot("CookieLifetime"),
			jen.Line(),
			jen.Return(jen.ID("svc"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}
