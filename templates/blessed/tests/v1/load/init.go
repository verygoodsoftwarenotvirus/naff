package load

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func initDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("main")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Var().Defs(
			jen.ID("debug").ID("bool"),
			jen.ID("urlToUse").ID("string"),
			jen.ID("oa2Client").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("init").Params().Block(
			utils.InlineFakeSeedFunc(),
			jen.ID("urlToUse").Equals().Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil"), "DetermineServiceURL").Call(),
			jen.ID("logger").Assign().Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog", "NewZeroLogger").Call(),
			jen.Line(),
			jen.ID("logger").Dot("WithValue").Call(jen.Lit("url"), jen.ID("urlToUse")).Dot("Info").Call(jen.Lit("checking server")),
			jen.Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil"), "EnsureServerIsUp").Call(jen.ID("urlToUse")),
			jen.Line(),
			jen.Qual(utils.FakeLibrary, "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call()),
			jen.Line(),
			jen.List(jen.ID("u"), jen.Err()).Assign().Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil"), "CreateObligatoryUser").Call(jen.ID("urlToUse"), jen.ID("debug")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("logger").Dot("Fatal").Call(jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("oa2Client"), jen.Err()).Equals().Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil"), "CreateObligatoryClient").Call(jen.ID("urlToUse"), jen.ID("u")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("logger").Dot("Fatal").Call(jen.Err()),
			),
			jen.Line(),
			jen.ID("fiftySpaces").Assign().Qual("strings", "Repeat").Call(jen.Lit("\n"), jen.Lit(50)),
			jen.Qual("fmt", "Printf").Call(jen.Lit("%s\tRunning tests%s"), jen.ID("fiftySpaces"), jen.ID("fiftySpaces")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildHTTPClient").Params().Params(jen.ParamPointer().Qual("net/http", "Client")).Block(
			jen.ID("httpc").Assign().VarPointer().Qual("net/http", "Client").Valuesln(
				jen.ID("Transport").MapAssign().Qual("net/http", "DefaultTransport"),
				jen.ID("Timeout").MapAssign().Lit(5).Times().Qual("time", "Second"),
			),
			jen.Line(),
			jen.Return().ID("httpc"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("initializeClient").Params(jen.ID("oa2Client").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "client/v1/http"), "V1Client")).Block(
			jen.List(jen.ID("uri"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.ID("urlToUse")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("c"), jen.Err()).Assign().Qual(filepath.Join(pkg.OutputPath, "client/v1/http"), "NewClient").Callln(
				utils.InlineCtx(),
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
	return ret
}
