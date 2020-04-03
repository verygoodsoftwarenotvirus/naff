package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2HandlerMockTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().ID("oauth2Handler").Equals().Parens(jen.PointerTo().ID("mockOauth2Handler")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("mockOauth2Handler").Struct(jen.Qual("github.com/stretchr/testify/mock",
			"Mock",
		)),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOauth2Handler")).ID("SetAllowGetAccessRequest").Params(jen.ID("allowed").Bool()).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("allowed")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOauth2Handler")).ID("SetClientAuthorizedHandler").Params(jen.ID("handler").Qual("goproj.in/oauth2.v3/server", "ClientAuthorizedHandler")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOauth2Handler")).ID("SetClientScopeHandler").Params(jen.ID("handler").Qual("goproj.in/oauth2.v3/server", "ClientScopeHandler")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOauth2Handler")).ID("SetClientInfoHandler").Params(jen.ID("handler").Qual("goproj.in/oauth2.v3/server", "ClientInfoHandler")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOauth2Handler")).ID("SetUserAuthorizationHandler").Params(jen.ID("handler").Qual("goproj.in/oauth2.v3/server", "UserAuthorizationHandler")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOauth2Handler")).ID("SetAuthorizeScopeHandler").Params(jen.ID("handler").Qual("goproj.in/oauth2.v3/server", "AuthorizeScopeHandler")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOauth2Handler")).ID("SetResponseErrorHandler").Params(jen.ID("handler").Qual("goproj.in/oauth2.v3/server", "ResponseErrorHandler")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOauth2Handler")).ID("SetInternalErrorHandler").Params(jen.ID("handler").Qual("goproj.in/oauth2.v3/server", "InternalErrorHandler")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOauth2Handler")).ID("ValidationBearerToken").Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Qual("goproj.in/oauth2.v3",
			"TokenInfo",
		),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot(
				"Called",
			).Call(jen.ID("req")),
			jen.Return().List(
				jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("goproj.in/oauth2.v3", "TokenInfo")),
				jen.ID("args").Dot("Error").Call(jen.Lit(1)),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOauth2Handler")).ID("HandleAuthorizeRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot(
				"Called",
			).Call(jen.ID("res"), jen.ID("req")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOauth2Handler")).ID("HandleTokenRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot(
				"Called",
			).Call(jen.ID("res"), jen.ID("req")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)
	return ret
}
