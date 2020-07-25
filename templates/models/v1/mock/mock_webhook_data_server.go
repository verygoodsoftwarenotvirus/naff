package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockWebhookDataServerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("mock")

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Underscore().Qual(proj.ModelsV1Package(), "WebhookDataServer").Equals().Parens(jen.PointerTo().ID("WebhookDataServer")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(
		jen.Comment("WebhookDataServer is a mocked models.WebhookDataServer for testing"),
		jen.Line(),
		jen.Type().ID("WebhookDataServer").Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreationInputMiddleware implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateInputMiddleware implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ListHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("ListHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Block(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("CreateHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Block(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ReadHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("ReadHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Block(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("UpdateHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Block(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("ArchiveHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Block(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	)

	return code
}
