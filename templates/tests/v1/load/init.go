package load

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func initDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("main")

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Defs(
			jen.ID("debug").Bool(),
			jen.ID("urlToUse").String(),
			jen.ID("oa2Client").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("init").Params().Block(
			jen.ID("urlToUse").Equals().Qual(proj.TestutilV1Package(), "DetermineServiceURL").Call(),
			jen.ID(constants.LoggerVarName).Assign().Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog", "NewZeroLogger").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("url"), jen.ID("urlToUse")).Dot("Info").Call(jen.Lit("checking server")),
			jen.Qual(proj.TestutilV1Package(), "EnsureServerIsUp").Call(jen.ID("urlToUse")),
			jen.Line(),
			jen.List(jen.ID("u"), jen.Err()).Assign().Qual(proj.TestutilV1Package(), "CreateObligatoryUser").Call(jen.ID("urlToUse"), jen.ID("debug")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Fatal").Call(jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("oa2Client"), jen.Err()).Equals().Qual(proj.TestutilV1Package(), "CreateObligatoryClient").Call(jen.ID("urlToUse"), jen.ID("u")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Fatal").Call(jen.Err()),
			),
			jen.Line(),
			jen.ID("fiftySpaces").Assign().Qual("strings", "Repeat").Call(jen.Lit("\n"), jen.Lit(50)),
			jen.Qual("fmt", "Printf").Call(jen.Lit("%s\tRunning tests%s"), jen.ID("fiftySpaces"), jen.ID("fiftySpaces")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildHTTPClient").Params().Params(jen.PointerTo().Qual("net/http", "Client")).Block(
			jen.ID("httpc").Assign().AddressOf().Qual("net/http", "Client").Valuesln(
				jen.ID("Transport").MapAssign().Qual("net/http", "DefaultTransport"),
				jen.ID("Timeout").MapAssign().Lit(5).Times().Qual("time", "Second"),
			),
			jen.Line(),
			jen.Return().ID("httpc"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("initializeClient").Params(jen.ID("oa2Client").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Params(jen.PointerTo().Qual(proj.HTTPClientV1Package(), "V1Client")).Block(
			jen.List(jen.ID("uri"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.ID("urlToUse")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("c"), jen.Err()).Assign().Qual(proj.HTTPClientV1Package(), "NewClient").Callln(
				constants.InlineCtx(),
				jen.ID("oa2Client").Dot("ClientID"),
				jen.ID("oa2Client").Dot("ClientSecret"),
				jen.ID("uri"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog", "NewZeroLogger").Call(),
				jen.ID("buildHTTPClient").Call(),
				jen.ID("oa2Client").Dot("Scopes"),
				jen.ID("debug"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Return().ID("c"),
		),
		jen.Line(),
	)

	return code
}
