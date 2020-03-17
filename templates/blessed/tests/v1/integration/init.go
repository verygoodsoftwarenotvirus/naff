package integration

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func initDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("integration")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("debug").Op("=").ID("true"),
			jen.ID("nonexistentID").Op("=").Lit(999999999),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("urlToUse").ID("string"),
			jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Op("*").Qual(filepath.Join(pkg.OutputPath, "client/v1/http"), "V1Client"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("init").Params().Block(
			jen.ID("urlToUse").Op("=").Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil"), "DetermineServiceURL").Call(),
			jen.ID("logger").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog", "NewZeroLogger").Call(),
			jen.Line(),
			jen.ID("logger").Dot("WithValue").Call(jen.Lit("url"), jen.ID("urlToUse")).Dot("Info").Call(jen.Lit("checking server")),
			jen.Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil"), "EnsureServerIsUp").Call(jen.ID("urlToUse")),
			jen.Line(),
			jen.Qual(utils.FakeLibrary, "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call()),
			jen.Line(),
			jen.List(jen.ID("ogUser"), jen.Err()).Op(":=").Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil"), "CreateObligatoryUser").Call(jen.ID("urlToUse"), jen.ID("debug")),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot("Fatal").Call(jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("oa2Client"), jen.Err()).Op(":=").Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil"), "CreateObligatoryClient").Call(jen.ID("urlToUse"), jen.ID("ogUser")),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot("Fatal").Call(jen.Err()),
			),
			jen.Line(),
			jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Op("=").ID("initializeClient").Call(jen.ID("oa2Client")), // VARME
			jen.Line(),
			jen.ID("fiftySpaces").Op(":=").Qual("strings", "Repeat").Call(jen.Lit("\n"), jen.Lit(50)),
			jen.Qual("fmt", "Printf").Call(jen.Lit("%s\tRunning tests%s"), jen.ID("fiftySpaces"), jen.ID("fiftySpaces")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildHTTPClient").Params().Params(jen.Op("*").Qual("net/http", "Client")).Block(
			jen.ID("httpc").Op(":=").Op("&").Qual("net/http", "Client").Valuesln(
				jen.ID("Transport").Op(":").Qual("net/http", "DefaultTransport"),
				jen.ID("Timeout").Op(":").Lit(5).Op("*").Qual("time", "Second"),
			),
			jen.Line(),
			jen.Return().ID("httpc"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("initializeClient").Params(jen.ID("oa2Client").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "client/v1/http"), "V1Client")).Block(
			jen.List(jen.ID("uri"), jen.Err()).Op(":=").Qual("net/url", "Parse").Call(jen.ID("urlToUse")),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("c"), jen.Err()).Op(":=").Qual(filepath.Join(pkg.OutputPath, "client/v1/http"), "NewClient").Callln(
				jen.Qual("context", "Background").Call(),
				jen.ID("oa2Client").Dot("ClientID"),
				jen.ID("oa2Client").Dot("ClientSecret"),
				jen.ID("uri"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog", "NewZeroLogger").Call(),
				jen.ID("buildHTTPClient").Call(),
				jen.ID("oa2Client").Dot("Scopes"),
				jen.ID("debug"),
			),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Return().ID("c"),
		),
		jen.Line(),
	)
	return ret
}
