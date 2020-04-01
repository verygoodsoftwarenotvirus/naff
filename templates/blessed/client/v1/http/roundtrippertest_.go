package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func roundtripperTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile(packageName)

	utils.AddImports(pkg, ret)

	ret.Add(jen.Line())

	ret.Add(
		jen.Func().ID("Test_buildDefaultTransport").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				jen.ID("_").Equals().ID("buildDefaultTransport").Call(),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_defaultRoundTripper_RoundTrip").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				jen.ID("transport").Assign().ID("newDefaultRoundTripper").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("https://verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				jen.List(jen.ID("_"), jen.Err()).Equals().ID("transport").Dot("RoundTrip").Call(jen.ID("req")),
				utils.AssertNoError(jen.Err(), nil),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_newDefaultRoundTripper").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				jen.ID("_").Equals().ID("newDefaultRoundTripper").Call(),
			)),
		),
		jen.Line(),
	)

	return ret
}
