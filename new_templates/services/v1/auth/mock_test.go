package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func mockTestDotGo() *jen.File {
	ret := jen.NewFile("auth")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("OAuth2ClientValidator").Op("=").Parens(jen.Op("*").ID("mockOAuth2ClientValidator")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Null().Type().ID("mockOAuth2ClientValidator").Struct(
		jen.ID("mock").Dot(
			"Mock",
		),
	),
	)
	ret.Add(jen.Func().Params(jen.ID("m").Op("*").ID("mockOAuth2ClientValidator")).ID("ExtractOAuth2ClientFromRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("models").Dot(
		"OAuth2Client",
	), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("req")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Op("*").ID("models").Dot(
			"OAuth2Client",
		)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),
	)
	ret.Add(jen.Null().Var().ID("_").ID("cookieEncoderDecoder").Op("=").Parens(jen.Op("*").ID("mockCookieEncoderDecoder")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Null().Type().ID("mockCookieEncoderDecoder").Struct(
		jen.ID("mock").Dot(
			"Mock",
		),
	),
	)
	ret.Add(jen.Func().Params(jen.ID("m").Op("*").ID("mockCookieEncoderDecoder")).ID("Encode").Params(jen.ID("name").ID("string"), jen.ID("value").Interface()).Params(jen.ID("string"), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("name"), jen.ID("value")),
		jen.Return().List(jen.ID("args").Dot(
			"String",
		).Call(jen.Lit(0)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),
	)
	ret.Add(jen.Func().Params(jen.ID("m").Op("*").ID("mockCookieEncoderDecoder")).ID("Decode").Params(jen.List(jen.ID("name"), jen.ID("value")).ID("string"), jen.ID("dst").Interface()).Params(jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("name"), jen.ID("value"), jen.ID("dst")),
		jen.Return().ID("args").Dot(
			"Error",
		).Call(jen.Lit(0)),
	),
	)
	ret.Add(jen.Null().Var().ID("_").Qual("net/http", "Handler").Op("=").Parens(jen.Op("*").ID("MockHTTPHandler")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Null().Type().ID("MockHTTPHandler").Struct(
		jen.ID("mock").Dot(
			"Mock",
		),
	),
	)
	ret.Add(jen.Func().Params(jen.ID("m").Op("*").ID("MockHTTPHandler")).ID("ServeHTTP").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
		jen.ID("m").Dot(
			"Called",
		).Call(jen.ID("res"), jen.ID("req")),
	),
	)
	return ret
}
