package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func initDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("frontend")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().ID("urlToUse").String(),
		jen.Line(),
	)

	ret.Add(
		jen.Const().Defs(
			jen.ID("seleniumHubAddr").Equals().Lit("http://selenium-hub:4444/wd/hub"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("init").Params().Block(
			jen.ID("urlToUse").Equals().ID("testutil").Dot("DetermineServiceURL").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog", "NewZeroLogger").Call(),
			jen.ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("url"), jen.ID("urlToUse")).Dot("Info").Call(jen.Lit("checking server")),
			jen.Qual(proj.TestutilV1Package(), "EnsureServerIsUp").Call(jen.ID("urlToUse")),
			jen.Line(),
			jen.Comment("NOTE: this is sad, but also the only thing that consistently works"),
			jen.Comment("see above for my vain attempts at a real solution to this problem."),
			jen.Qual("time", "Sleep").Call(jen.Lit(10).Times().Qual("time", "Second")),
			jen.Line(),
			jen.ID("fiftySpaces").Assign().Qual("strings", "Repeat").Call(jen.Lit("\n"), jen.Lit(50)),
			jen.Qual("fmt", "Printf").Call(jen.Lit("%s\tRunning tests%s"), jen.ID("fiftySpaces"), jen.ID("fiftySpaces")),
		),
		jen.Line(),
	)

	return ret
}
