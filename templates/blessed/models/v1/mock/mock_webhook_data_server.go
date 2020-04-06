package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockWebhookDataServerDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().Qual(proj.ModelsV1Package(), "WebhookDataServer").Equals().Parens(jen.PointerTo().ID("WebhookDataServer")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("WebhookDataServer is a mocked models.WebhookDataServer for testing"),
		jen.Line(),
		jen.Type().ID("WebhookDataServer").Struct(jen.Qual(utils.MockPkg, "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreationInputMiddleware implements our interface requirements"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateInputMiddleware implements our interface requirements"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ListHandler implements our interface requirements"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateHandler implements our interface requirements"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ReadHandler implements our interface requirements"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateHandler implements our interface requirements"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("UpdateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveHandler implements our interface requirements"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataServer")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)
	return ret
}
