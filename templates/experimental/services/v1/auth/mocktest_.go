package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func mockTestDotGo() *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("_").ID("OAuth2ClientValidator").Op("=").Parens(jen.Op("*").ID("mockOAuth2ClientValidator")).Call(jen.ID("nil")),
	jen.Line(),
	)

	ret.Add(
		jen.Type().ID("mockOAuth2ClientValidator").Struct(jen.ID("mock").Dot(
		"Mock",
	)),
	jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockOAuth2ClientValidator")).ID("ExtractOAuth2ClientFromRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("models").Dot(
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
	jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").ID("cookieEncoderDecoder").Op("=").Parens(jen.Op("*").ID("mockCookieEncoderDecoder")).Call(jen.ID("nil")),
	jen.Line(),
	)

	ret.Add(
		jen.Type().ID("mockCookieEncoderDecoder").Struct(jen.ID("mock").Dot(
		"Mock",
	)),
	jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockCookieEncoderDecoder")).ID("Encode").Params(jen.ID("name").ID("string"), jen.ID("value").Interface()).Params(jen.ID("string"), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("name"), jen.ID("value")),
		jen.Return().List(jen.ID("args").Dot(
			"String",
		).Call(jen.Lit(0)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockCookieEncoderDecoder")).ID("Decode").Params(jen.List(jen.ID("name"), jen.ID("value")).ID("string"), jen.ID("dst").Interface()).Params(jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("name"), jen.ID("value"), jen.ID("dst")),
		jen.Return().ID("args").Dot(
			"Error",
		).Call(jen.Lit(0)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").Qual("net/http", "Handler").Op("=").Parens(jen.Op("*").ID("MockHTTPHandler")).Call(jen.ID("nil")),
	jen.Line(),
	)

	ret.Add(
		jen.Type().ID("MockHTTPHandler").Struct(jen.ID("mock").Dot(
		"Mock",
	)),
	jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("MockHTTPHandler")).ID("ServeHTTP").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
		jen.ID("m").Dot(
			"Called",
		).Call(jen.ID("res"), jen.ID("req")),
	),
	jen.Line(),
	)
	return ret
}
