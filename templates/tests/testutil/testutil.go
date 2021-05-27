package testutil

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func testutilDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.Qual("github.com/brianvoe/gofakeit/v5", "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call())),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("errArbitrary").Op("=").Qual("errors", "New").Call(jen.Lit("blah")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BrokenSessionContextDataFetcher is a deliberately broken sessionContextDataFetcher."),
		jen.Line(),
		jen.Func().ID("BrokenSessionContextDataFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
			jen.Return().List(jen.ID("nil"), jen.ID("errArbitrary"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateBodyFromStruct takes any value in and returns an io.ReadCloser for an http.Request's body."),
		jen.Line(),
		jen.Func().ID("CreateBodyFromStruct").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("in").Interface()).Params(jen.Qual("io", "ReadCloser")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.List(jen.ID("out"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("in")),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().Qual("io", "NopCloser").Call(jen.Qual("bytes", "NewReader").Call(jen.ID("out"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildArbitraryImage builds an image with a bunch of colors in it."),
		jen.Line(),
		jen.Func().ID("BuildArbitraryImage").Params(jen.ID("widthAndHeight").ID("int")).Params(jen.Qual("image", "Image")).Body(
			jen.ID("img").Op(":=").Qual("image", "NewRGBA").Call(jen.Qual("image", "Rectangle").Valuesln(jen.ID("Min").Op(":").Qual("image", "Point").Values(), jen.ID("Max").Op(":").Qual("image", "Point").Valuesln(jen.ID("X").Op(":").ID("widthAndHeight"), jen.ID("Y").Op(":").ID("widthAndHeight")))),
			jen.For(jen.ID("x").Op(":=").Lit(0), jen.ID("x").Op("<").ID("widthAndHeight"), jen.ID("x").Op("++")).Body(
				jen.For(jen.ID("y").Op(":=").Lit(0), jen.ID("y").Op("<").ID("widthAndHeight"), jen.ID("y").Op("++")).Body(
					jen.ID("img").Dot("Set").Call(
						jen.ID("x"),
						jen.ID("y"),
						jen.Qual("image/color", "RGBA").Valuesln(jen.ID("R").Op(":").ID("uint8").Call(jen.ID("x").Op("%").Qual("math", "MaxUint8")), jen.ID("G").Op(":").ID("uint8").Call(jen.ID("y").Op("%").Qual("math", "MaxUint8")), jen.ID("B").Op(":").ID("uint8").Call(jen.ID("x").Op("+").ID("y").Op("%").Qual("math", "MaxUint8")), jen.ID("A").Op(":").Qual("math", "MaxUint8")),
					))),
			jen.Return().ID("img"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildArbitraryImagePNGBytes builds an image with a bunch of colors in it."),
		jen.Line(),
		jen.Func().ID("BuildArbitraryImagePNGBytes").Params(jen.ID("widthAndHeight").ID("int")).Params(jen.Index().ID("byte")).Body(
			jen.Var().Defs(
				jen.ID("b").Qual("bytes", "Buffer"),
			),
			jen.If(jen.ID("err").Op(":=").Qual("image/png", "Encode").Call(
				jen.Op("&").ID("b"),
				jen.ID("BuildArbitraryImage").Call(jen.ID("widthAndHeight")),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
			jen.Return().ID("b").Dot("Bytes").Call(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AssertAppropriateNumberOfTestsRan ensures the expected number of tests are run in a given suite."),
		jen.Line(),
		jen.Func().ID("AssertAppropriateNumberOfTestsRan").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("totalExpectedTestCount").ID("uint"), jen.ID("stats").Op("*").ID("suite").Dot("SuiteInformation")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.If(jen.ID("stats").Dot("Passed").Call()).Body(
				jen.ID("require").Dot("Equal").Call(
					jen.ID("t"),
					jen.ID("int").Call(jen.ID("totalExpectedTestCount")),
					jen.ID("len").Call(jen.ID("stats").Dot("TestStats")),
					jen.Lit("expected total number of tests run to equal %d, but it was %d"),
					jen.ID("totalExpectedTestCount"),
					jen.ID("len").Call(jen.ID("stats").Dot("TestStats")),
				)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildTestRequest builds an arbitrary *http.Request."),
		jen.Line(),
		jen.Func().ID("BuildTestRequest").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").Qual("net/http", "Request")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodOptions"),
				jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
				jen.ID("nil"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("req"),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().ID("req"),
		),
		jen.Line(),
	)

	return code
}
