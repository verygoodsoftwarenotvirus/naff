package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2HandlerMockTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().ID("oauth2Handler").Equals().Parens(jen.PointerTo().ID("mockOAuth2Handler")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("mockOAuth2Handler").Struct(jen.Qual(utils.MockPkg,
			"Mock",
		)),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetAllowGetAccessRequest").Params(jen.ID("allowed").Bool()).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("allowed")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetClientAuthorizedHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ClientAuthorizedHandler")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetClientScopeHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ClientScopeHandler")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetClientInfoHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ClientInfoHandler")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetUserAuthorizationHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "UserAuthorizationHandler")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetAuthorizeScopeHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "AuthorizeScopeHandler")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetResponseErrorHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ResponseErrorHandler")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("SetInternalErrorHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "InternalErrorHandler")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("ValidationBearerToken").Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Qual("gopkg.in/oauth2.v3",
			"TokenInfo",
		),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot(
				"Called",
			).Call(jen.ID(constants.RequestVarName)),
			jen.Return().List(
				jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("gopkg.in/oauth2.v3", "TokenInfo")),
				jen.ID("args").Dot("Error").Call(jen.One()),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("HandleAuthorizeRequest").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot(
				"Called",
			).Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2Handler")).ID("HandleTokenRequest").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot(
				"Called",
			).Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)

	return ret
}
