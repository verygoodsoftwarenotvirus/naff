package data_scaffolder

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func exiterTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("_").ID("Quitter").Op("=").Op("&").ID("mockQuitter").Values(),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("mockQuitter").Struct(jen.ID("mock").Dot("Mock")),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockQuitter")).ID("Quit").Params(jen.ID("code").ID("int")).Body(
			jen.ID("m").Dot("Called").Call(jen.ID("code"))),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockQuitter")).ID("ComplainAndQuit").Params(jen.ID("v").Op("...").Interface()).Body(
			jen.ID("m").Dot("Called").Call(jen.ID("v").Op("..."))),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockQuitter")).ID("ComplainAndQuitf").Params(jen.ID("s").ID("string"), jen.ID("args").Op("...").Interface()).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("s"),
				jen.ID("args"),
			)),
		jen.Line(),
	)

	return code
}
