package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func roundtripperTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildTestBuildDefaultTransport()...)
	code.Add(buildTestDefaultRoundTripperRoundTrip()...)
	code.Add(buildTestNewDefaultRoundTripper()...)

	return code
}

func buildTestBuildDefaultTransport() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("Test_buildDefaultTransport").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().ID("buildDefaultTransport").Call(),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDefaultRoundTripperRoundTrip() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("Test_defaultRoundTripper_RoundTrip").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("transport").Assign().ID("newDefaultRoundTripper").Call(),
				jen.Line(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("https://verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				jen.List(jen.Underscore(), jen.Err()).Equals().ID("transport").Dot("RoundTrip").Call(jen.ID(constants.RequestVarName)),
				utils.AssertNoError(jen.Err(), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestNewDefaultRoundTripper() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("Test_newDefaultRoundTripper").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().ID("newDefaultRoundTripper").Call(),
			),
		),
		jen.Line(),
	}

	return lines
}
