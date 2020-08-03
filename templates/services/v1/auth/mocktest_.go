package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Underscore().ID("OAuth2ClientValidator").Equals().Parens(jen.PointerTo().ID("mockOAuth2ClientValidator")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(buildMockOAuth2ClientValidator()...)
	code.Add(buildMockOAuth2ClientValidatorExtractOAuth2ClientFromRequest(proj)...)

	code.Add(
		jen.Var().Underscore().ID("cookieEncoderDecoder").Equals().Parens(jen.PointerTo().ID("mockCookieEncoderDecoder")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(buildMockCookieEncoderDecoder()...)
	code.Add(buildMockCookieEncoderDecoderEncode()...)
	code.Add(buildMockCookieEncoderDecoderDecode()...)

	code.Add(
		jen.Var().Underscore().Qual("net/http", "Handler").Equals().Parens(jen.PointerTo().ID("MockHTTPHandler")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(buildMockHTTPHandler()...)
	code.Add(buildMockHTTPHandlerServeHTTP()...)

	return code
}

func buildMockOAuth2ClientValidator() []jen.Code {
	lines := []jen.Code{
		jen.Type().ID("mockOAuth2ClientValidator").Struct(jen.Qual(constants.MockPkg,
			"Mock",
		)),
		jen.Line(),
	}

	return lines
}

func buildMockOAuth2ClientValidatorExtractOAuth2ClientFromRequest(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockOAuth2ClientValidator")).ID("ExtractOAuth2ClientFromRequest").Params(
			constants.CtxParam(),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID(constants.RequestVarName)),
			jen.Return().List(
				jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")),
				jen.ID("args").Dot("Error").Call(jen.One()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildMockCookieEncoderDecoder() []jen.Code {
	lines := []jen.Code{
		jen.Type().ID("mockCookieEncoderDecoder").Struct(jen.Qual(constants.MockPkg,
			"Mock",
		)),
		jen.Line(),
	}

	return lines
}

func buildMockCookieEncoderDecoderEncode() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockCookieEncoderDecoder")).ID("Encode").Params(jen.ID("name").String(), jen.ID("value").Interface()).Params(jen.String(), jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot(
				"Called",
			).Call(jen.ID("name"), jen.ID("value")),
			jen.Return().List(
				jen.ID("args").Dot("String").Call(jen.Zero()),
				jen.ID("args").Dot("Error").Call(jen.One()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildMockCookieEncoderDecoderDecode() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockCookieEncoderDecoder")).ID("Decode").Params(jen.List(jen.ID("name"), jen.ID("value")).String(), jen.ID("dst").Interface()).Params(jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot(
				"Called",
			).Call(jen.ID("name"), jen.ID("value"), jen.ID("dst")),
			jen.Return().ID("args").Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}

func buildMockHTTPHandler() []jen.Code {
	lines := []jen.Code{
		jen.Type().ID("MockHTTPHandler").Struct(jen.Qual(constants.MockPkg,
			"Mock",
		)),
		jen.Line(),
	}

	return lines
}

func buildMockHTTPHandlerServeHTTP() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("MockHTTPHandler")).ID("ServeHTTP").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Body(
			jen.ID("m").Dot(
				"Called",
			).Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}
