package data_scaffolder

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func exiterDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().ID("Quitter").Interface(
			jen.ID("Quit").Params(jen.ID("code").ID("int")),
			jen.ID("ComplainAndQuit").Params(jen.Op("...").Interface()),
			jen.ID("ComplainAndQuitf").Params(jen.ID("string"), jen.Op("...").Interface()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").ID("Quitter").Op("=").ID("fatalQuitter").Values(),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("fatalQuitter").Struct(),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("e").ID("fatalQuitter")).ID("Quit").Params(jen.ID("code").ID("int")).Body(
			jen.Qual("os", "Exit").Call(jen.ID("code"))),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("e").ID("fatalQuitter")).ID("ComplainAndQuit").Params(jen.ID("v").Op("...").Interface()).Body(
			jen.Qual("log", "Fatal").Call(jen.ID("v").Op("..."))),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("e").ID("fatalQuitter")).ID("ComplainAndQuitf").Params(jen.ID("format").ID("string"), jen.ID("args").Op("...").Interface()).Body(
			jen.Qual("log", "Fatalf").Call(
				jen.ID("format"),
				jen.ID("args").Op("..."),
			)),
		jen.Line(),
	)

	return code
}
