package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockWebhookDataServerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Underscore().Qual(proj.TypesPackage(), "WebhookDataServer").Equals().Parens(jen.PointerTo().ID("WebhookDataServer")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(buildMockWebhookDataServer()...)
	code.Add(buildMockWebhookCreationInputMiddleware()...)
	code.Add(buildMockWebhookUpdateInputMiddleware()...)
	code.Add(buildMockWebhookListHandler()...)
	code.Add(buildMockWebhookCreateHandler()...)
	code.Add(buildMockWebhookReadHandler()...)
	code.Add(buildMockWebhookUpdateHandler()...)
	code.Add(buildMockWebhookArchiveHandler()...)

	return code
}

func buildMockWebhookDataServer() []jen.Code {
	lines := []jen.Code{
		jen.Comment("WebhookDataServer is a mocked models.WebhookDataServer for testing"),
		jen.Line(),
		jen.Type().ID("WebhookDataServer").Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	}

	return lines
}

func buildMockWebhookCreationInputMiddleware() []jen.Code {
	lines := []jen.Code{
		jen.Comment("CreationInputMiddleware implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildMockWebhookUpdateInputMiddleware() []jen.Code {
	lines := []jen.Code{
		jen.Comment("UpdateInputMiddleware implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildMockWebhookListHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ListHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("ListHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildMockWebhookCreateHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("CreateHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("CreateHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildMockWebhookReadHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ReadHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("ReadHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildMockWebhookUpdateHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("UpdateHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("UpdateHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildMockWebhookArchiveHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ArchiveHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("ArchiveHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}
