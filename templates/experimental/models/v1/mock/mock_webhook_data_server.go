package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func mockWebhookDataServerDotGo() *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"WebhookDataServer",
	).Op("=").Parens(jen.Op("*").ID("WebhookDataServer")).Call(jen.ID("nil")),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("WebhookDataServer").Struct(jen.ID("mock").Dot(
		"Mock",
	)),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// CreationInputMiddleware implements our interface requirements").Params(jen.ID("m").Op("*").ID("WebhookDataServer")).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("next")),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// UpdateInputMiddleware implements our interface requirements").Params(jen.ID("m").Op("*").ID("WebhookDataServer")).ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("next")),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ListHandler implements our interface requirements").Params(jen.ID("m").Op("*").ID("WebhookDataServer")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// CreateHandler implements our interface requirements").Params(jen.ID("m").Op("*").ID("WebhookDataServer")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ReadHandler implements our interface requirements").Params(jen.ID("m").Op("*").ID("WebhookDataServer")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// UpdateHandler implements our interface requirements").Params(jen.ID("m").Op("*").ID("WebhookDataServer")).ID("UpdateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ArchiveHandler implements our interface requirements").Params(jen.ID("m").Op("*").ID("WebhookDataServer")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
	),

		jen.Line(),
	)
	return ret
}
