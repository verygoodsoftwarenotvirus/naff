package mock

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func mockOauth2ClientDataServerDotGo() *jen.File {
	ret := jen.NewFile("mock")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"OAuth2ClientDataServer",
	).Op("=").Parens(jen.Op("*").ID("OAuth2ClientDataServer")).Call(jen.ID("nil")),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("OAuth2ClientDataServer").Struct(jen.ID("mock").Dot(
		"Mock",
	)),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ListHandler is the obligatory implementation for our interface").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// CreateHandler is the obligatory implementation for our interface").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ReadHandler is the obligatory implementation for our interface").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ArchiveHandler is the obligatory implementation for our interface").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// CreationInputMiddleware is the obligatory implementation for our interface").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("next")),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// OAuth2ClientInfoMiddleware is the obligatory implementation for our interface").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("OAuth2ClientInfoMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("next")),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ExtractOAuth2ClientFromRequest is the obligatory implementation for our interface").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("ExtractOAuth2ClientFromRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("models").Dot(
		"OAuth2Client",
	), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("req")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Op("*").ID("models").Dot(
			"OAuth2Client",
		)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// HandleAuthorizeRequest is the obligatory implementation for our interface").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("HandleAuthorizeRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("res"), jen.ID("req")),
		jen.Return().ID("args").Dot(
			"Error",
		).Call(jen.Lit(0)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// HandleTokenRequest is the obligatory implementation for our interface").Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("HandleTokenRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("res"), jen.ID("req")),
		jen.Return().ID("args").Dot(
			"Error",
		).Call(jen.Lit(0)),
	),

		jen.Line(),
	)
	return ret
}
