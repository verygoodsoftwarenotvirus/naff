package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func initDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildInitConstDefs()...)
	code.Add(buildInitVarDefs(proj)...)
	code.Add(buildInitInit(proj)...)
	code.Add(buildInitBuildHTTPClient()...)
	code.Add(buildInitInitializeClient(proj)...)

	return code
}

func buildInitConstDefs() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.ID("debug").Equals().True(),
			jen.ID("nonexistentID").Equals().Lit(999999999),
		),
		jen.Line(),
	}

	return lines
}

func buildInitVarDefs(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.ID("urlToUse").String(),
			jen.IDf("%sClient", proj.Name.UnexportedVarName()).PointerTo().Qual(proj.HTTPClientV1Package(), "V1Client"),
		),
		jen.Line(),
	}

	return lines
}

func buildInitInit(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("init").Params().Body(
			jen.ID("urlToUse").Equals().Qual(proj.TestUtilPackage(), "DetermineServiceURL").Call(),
			jen.ID(constants.LoggerVarName).Assign().Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog", "NewZeroLogger").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("url"), jen.ID("urlToUse")).Dot("Info").Call(jen.Lit("checking server")),
			jen.Qual(proj.TestUtilPackage(), "EnsureServerIsUp").Call(jen.ID("urlToUse")),
			jen.Line(),
			jen.Line(),
			jen.List(jen.ID("ogUser"), jen.Err()).Assign().Qual(proj.TestUtilPackage(), "CreateObligatoryUser").Call(jen.ID("urlToUse"), jen.ID("debug")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID(constants.LoggerVarName).Dot("Fatal").Call(jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("oa2Client"), jen.Err()).Assign().Qual(proj.TestUtilPackage(), "CreateObligatoryClient").Call(jen.ID("urlToUse"), jen.ID("ogUser")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID(constants.LoggerVarName).Dot("Fatal").Call(jen.Err()),
			),
			jen.Line(),
			jen.IDf("%sClient", proj.Name.UnexportedVarName()).Equals().ID("initializeClient").Call(jen.ID("oa2Client")),
			jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("Debug").Equals().ID("urlToUse").IsEqualTo().EmptyString().Comment("change this for debug logs"),
			jen.Line(),
			jen.ID("fiftySpaces").Assign().Qual("strings", "Repeat").Call(jen.Lit("\n"), jen.Lit(50)),
			jen.Qual("fmt", "Printf").Call(jen.Lit("%s\tRunning tests%s"), jen.ID("fiftySpaces"), jen.ID("fiftySpaces")),
		),
		jen.Line(),
	}

	return lines
}

func buildInitBuildHTTPClient() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildHTTPClient").Params().Params(jen.PointerTo().Qual("net/http", "Client")).Body(
			jen.Return(
				jen.AddressOf().Qual("net/http", "Client").Valuesln(
					jen.ID("Transport").MapAssign().Qual("net/http", "DefaultTransport"),
					jen.ID("Timeout").MapAssign().Lit(5).Times().Qual("time", "Second"),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildInitInitializeClient(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("initializeClient").Params(jen.ID("oa2Client").PointerTo().Qual(proj.TypesPackage(), "OAuth2Client")).Params(jen.PointerTo().Qual(proj.HTTPClientV1Package(), "V1Client")).Body(
			jen.List(jen.ID("uri"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.ID("urlToUse")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
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
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Return().ID("c"),
		),
		jen.Line(),
	}

	return lines
}
