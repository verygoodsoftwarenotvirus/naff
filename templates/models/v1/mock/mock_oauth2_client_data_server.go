package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockOauth2ClientDataServerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("mock")

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Underscore().Qual(proj.ModelsV1Package(), "OAuth2ClientDataServer").Equals().Parens(jen.PointerTo().ID("OAuth2ClientDataServer")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(
		jen.Comment("OAuth2ClientDataServer is a mocked models.OAuth2ClientDataServer for testprojects"),
		jen.Line(),
		jen.Type().ID("OAuth2ClientDataServer").Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ListHandler is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateHandler is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ReadHandler is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveHandler is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreationInputMiddleware is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("OAuth2ClientInfoMiddleware is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("OAuth2ClientInfoMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ExtractOAuth2ClientFromRequest is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("ExtractOAuth2ClientFromRequest").Params(constants.CtxParam(), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID(constants.RequestVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("HandleAuthorizeRequest is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("HandleAuthorizeRequest").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
			jen.Return().ID("args").Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("HandleTokenRequest is the obligatory implementation for our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataServer")).ID("HandleTokenRequest").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
			jen.Return().ID("args").Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)

	return code
}
