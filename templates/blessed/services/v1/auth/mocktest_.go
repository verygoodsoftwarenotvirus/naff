package auth

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Var().ID("_").ID("OAuth2ClientValidator").Op("=").Parens(jen.Op("*").ID("mockOAuth2ClientValidator")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("mockOAuth2ClientValidator").Struct(jen.Qual("github.com/stretchr/testify/mock",
			"Mock",
		)),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockOAuth2ClientValidator")).ID("ExtractOAuth2ClientFromRequest").Params(utils.CtxParam(), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"),
			"OAuth2Client",
		),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot(
				"Called",
			).Call(jen.ID("req")),
			jen.Return().List(jen.ID("args").Dot(
				"Get",
			).Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"),
				"OAuth2Client",
			)), jen.ID("args").Dot("Error").Call(jen.Add(utils.FakeUint64Func()))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").ID("cookieEncoderDecoder").Op("=").Parens(jen.Op("*").ID("mockCookieEncoderDecoder")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("mockCookieEncoderDecoder").Struct(jen.Qual("github.com/stretchr/testify/mock",
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
			).Call(jen.Lit(0)), jen.ID("args").Dot("Error").Call(jen.Add(utils.FakeUint64Func()))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockCookieEncoderDecoder")).ID("Decode").Params(jen.List(jen.ID("name"), jen.ID("value")).ID("string"), jen.ID("dst").Interface()).Params(jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot(
				"Called",
			).Call(jen.ID("name"), jen.ID("value"), jen.ID("dst")),
			jen.Return().ID("args").Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").Qual("net/http", "Handler").Op("=").Parens(jen.Op("*").ID("MockHTTPHandler")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("MockHTTPHandler").Struct(jen.Qual("github.com/stretchr/testify/mock",
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
