package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockOauth2ClientDataServerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Underscore().Qual(proj.TypesPackage(), "OAuth2ClientDataServer").Equals().Parens(jen.PointerTo().ID("OAuth2ClientDataServer")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(buildOAuth2ClientDataServer()...)
	code.Add(buildOAuth2ClientListHandler()...)
	code.Add(buildOAuth2ClientCreateHandler()...)
	code.Add(buildOAuth2ClientReadHandler()...)
	code.Add(buildOAuth2ClientArchiveHandler()...)
	code.Add(buildOAuth2ClientCreationInputMiddleware()...)
	code.Add(buildOAuth2ClientInfoMiddleware()...)
	code.Add(buildExtractOAuth2ClientFromRequest(proj)...)
	code.Add(buildOAuth2ClientHandleAuthorizeRequest()...)
	code.Add(buildOAuth2ClientHandleTokenRequest()...)

	return code
}

func buildOAuth2ClientDataServer() []jen.Code {
	lines := []jen.Code{
		jen.Comment("OAuth2ClientDataServer is a mocked models.OAuth2ClientDataServer for testing"),
		jen.Line(),
		jen.Type().ID("OAuth2ClientDataServer").Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	}

	return lines
}

func buildOAuth2ClientListHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ListHandler is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("ListHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2ClientCreateHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("CreateHandler is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("CreateHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2ClientReadHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ReadHandler is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("ReadHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2ClientArchiveHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ArchiveHandler is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("ArchiveHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2ClientCreationInputMiddleware() []jen.Code {
	lines := []jen.Code{
		jen.Comment("CreationInputMiddleware is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2ClientInfoMiddleware() []jen.Code {
	lines := []jen.Code{
		jen.Comment("OAuth2ClientInfoMiddleware is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("OAuth2ClientInfoMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildExtractOAuth2ClientFromRequest(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ExtractOAuth2ClientFromRequest is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("ExtractOAuth2ClientFromRequest").Params(constants.CtxParam(), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "OAuth2Client"),
			jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID(constants.RequestVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.TypesPackage(), "OAuth2Client")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2ClientHandleAuthorizeRequest() []jen.Code {
	lines := []jen.Code{
		jen.Comment("HandleAuthorizeRequest is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("HandleAuthorizeRequest").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
			jen.Return().ID("args").Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2ClientHandleTokenRequest() []jen.Code {
	lines := []jen.Code{
		jen.Comment("HandleTokenRequest is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("HandleTokenRequest").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
			jen.Return().ID("args").Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}
