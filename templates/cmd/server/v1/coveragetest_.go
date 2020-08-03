package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func coverageTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("main")

	utils.AddImports(proj, code)

	code.Add(buildTestMain()...)

	return code
}

func buildTestMain() []jen.Code {
	return []jen.Code{
		jen.Line(),
		jen.Func().ID("TestRunMain").Params(jen.Underscore().PointerTo().Qual("testing", "T")).Body(
			jen.Comment("This test is built specifically to capture the coverage that the integration"),
			jen.Comment("tests exhibit. We run the main function (i.e. a production server)"),
			jen.Comment("on an independent goroutine and sleep for long enough that the integration"),
			jen.Comment("tests can run, then we quit."),
			jen.List(jen.ID("d"), jen.Err()).Assign().Qual("time", "ParseDuration").Call(jen.Qual("os", "Getenv").Call(jen.Lit("RUNTIME_DURATION"))),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Qual("log", "Fatal").Call(jen.Err()),
			),
			jen.Line(),
			jen.Go().ID("main").Call(),
			jen.Line(),
			jen.Qual("time", "Sleep").Call(jen.ID("d")),
		),
		jen.Line(),
	}
}
