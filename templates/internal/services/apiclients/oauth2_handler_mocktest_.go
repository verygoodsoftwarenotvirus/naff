package apiclients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2HandlerMockTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildOAuth2HandlerMockTestMockOAuth2Handler()...)
	code.Add(buildOAuth2HandlerMockTestSetAllowGetAccessRequest()...)
	code.Add(buildOAuth2HandlerMockTestSetClientAuthorizedHandler()...)
	code.Add(buildOAuth2HandlerMockTestSetClientScopeHandler()...)
	code.Add(buildOAuth2HandlerMockTestSetClientInfoHandler()...)
	code.Add(buildOAuth2HandlerMockTestSetUserAuthorizationHandler()...)
	code.Add(buildOAuth2HandlerMockTestSetAuthorizeScopeHandler()...)
	code.Add(buildOAuth2HandlerMockTestSetResponseErrorHandler()...)
	code.Add(buildOAuth2HandlerMockTestSetInternalErrorHandler()...)
	code.Add(buildOAuth2HandlerMockTestValidationBearerToken()...)
	code.Add(buildOAuth2HandlerMockTestHandleAuthorizeRequest()...)
	code.Add(buildOAuth2HandlerMockTestHandleTokenRequest()...)

	return code
}

func buildOAuth2HandlerMockTestMockOAuth2Handler() []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().ID("oauth2Handler").Equals().Parens(jen.PointerTo().ID("mockOAuth2Handler")).Call(jen.Nil()),
		jen.Line(),
		jen.Type().ID("mockOAuth2Handler").Struct(jen.Qual(constants.MockPkg,
			"Mock",
		)),
		jen.Line(),
	}

	return lines
}

func buildOAuth2HandlerMockTestSetAllowGetAccessRequest() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetAllowGetAccessRequest").Params(jen.ID("allowed").Bool()).Body(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("allowed")),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2HandlerMockTestSetClientAuthorizedHandler() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetClientAuthorizedHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ClientAuthorizedHandler")).Body(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2HandlerMockTestSetClientScopeHandler() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetClientScopeHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ClientScopeHandler")).Body(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2HandlerMockTestSetClientInfoHandler() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetClientInfoHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ClientInfoHandler")).Body(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2HandlerMockTestSetUserAuthorizationHandler() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetUserAuthorizationHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "UserAuthorizationHandler")).Body(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2HandlerMockTestSetAuthorizeScopeHandler() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetAuthorizeScopeHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "AuthorizeScopeHandler")).Body(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2HandlerMockTestSetResponseErrorHandler() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetResponseErrorHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ResponseErrorHandler")).Body(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2HandlerMockTestSetInternalErrorHandler() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetInternalErrorHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "InternalErrorHandler")).Body(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2HandlerMockTestValidationBearerToken() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("ValidationBearerToken").Params(
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params(
			jen.Qual("gopkg.in/oauth2.v3", "TokenInfo"),
			jen.Error(),
		).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID(constants.RequestVarName)),
			jen.Return().List(
				jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("gopkg.in/oauth2.v3", "TokenInfo")),
				jen.ID("args").Dot("Error").Call(jen.One()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2HandlerMockTestHandleAuthorizeRequest() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("HandleAuthorizeRequest").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Error()).Body(
			jen.Return().ID("m").Dot(
				"Called",
			).Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2HandlerMockTestHandleTokenRequest() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("HandleTokenRequest").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Error()).Body(
			jen.Return().ID("m").Dot(
				"Called",
			).Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}
