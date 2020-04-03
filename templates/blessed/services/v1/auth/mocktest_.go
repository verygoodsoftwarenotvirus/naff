package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().ID("_").ID("OAuth2ClientValidator").Equals().Parens(jen.PointerTo().ID("mockOAuth2ClientValidator")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("mockOAuth2ClientValidator").Struct(jen.Qual("github.com/stretchr/testify/mock",
			"Mock",
		)),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2ClientValidator")).ID("ExtractOAuth2ClientFromRequest").Params(
			utils.CtxParam(),
			jen.ID("req").ParamPointer().Qual("net/http", "Request"),
		).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("req")),
			jen.Return().List(
				jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")),
				jen.ID("args").Dot("Error").Call(jen.Lit(1)),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").ID("cookieEncoderDecoder").Equals().Parens(jen.PointerTo().ID("mockCookieEncoderDecoder")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("mockCookieEncoderDecoder").Struct(jen.Qual("github.com/stretchr/testify/mock",
			"Mock",
		)),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockCookieEncoderDecoder")).ID("Encode").Params(jen.ID("name").ID("string"), jen.ID("value").Interface()).Params(jen.ID("string"), jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot(
				"Called",
			).Call(jen.ID("name"), jen.ID("value")),
			jen.Return().List(
				jen.ID("args").Dot("String").Call(jen.Lit(0)),
				jen.ID("args").Dot("Error").Call(jen.Lit(1)),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockCookieEncoderDecoder")).ID("Decode").Params(jen.List(jen.ID("name"), jen.ID("value")).ID("string"), jen.ID("dst").Interface()).Params(jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot(
				"Called",
			).Call(jen.ID("name"), jen.ID("value"), jen.ID("dst")),
			jen.Return().ID("args").Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").Qual("net/http", "Handler").Equals().Parens(jen.PointerTo().ID("MockHTTPHandler")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("MockHTTPHandler").Struct(jen.Qual("github.com/stretchr/testify/mock",
			"Mock",
		)),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("MockHTTPHandler")).ID("ServeHTTP").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID("res"), jen.ID("req")),
		),
		jen.Line(),
	)
	return ret
}
