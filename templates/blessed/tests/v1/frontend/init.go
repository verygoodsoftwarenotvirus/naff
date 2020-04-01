package frontend

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func initDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("frontend")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Var().ID("urlToUse").ID("string"),
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
			utils.InlineFakeSeedFunc(),
			jen.ID("urlToUse").Equals().ID("testutil").Dot("DetermineServiceURL").Call(),
			jen.Line(),
			jen.ID("logger").Assign().Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog", "NewZeroLogger").Call(),
			jen.ID("logger").Dot("WithValue").Call(jen.Lit("url"), jen.ID("urlToUse")).Dot("Info").Call(jen.Lit("checking server")),
			jen.Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil"), "EnsureServerIsUp").Call(jen.ID("urlToUse")),
			jen.Line(),
			jen.Qual(utils.FakeLibrary, "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call()),
			jen.Line(),
			jen.Comment("NOTE: this is sad, but also the only thing that consistently works"),
			jen.Comment("see above for my vain attempts at a real solution to this problem"),
			jen.Qual("time", "Sleep").Call(jen.Lit(10).Times().Qual("time", "Second")),
			jen.Line(),
			jen.ID("fiftySpaces").Assign().Qual("strings", "Repeat").Call(jen.Lit("\n"), jen.Lit(50)),
			jen.Qual("fmt", "Printf").Call(jen.Lit("%s\tRunning tests%s"), jen.ID("fiftySpaces"), jen.ID("fiftySpaces")),
		),
		jen.Line(),
	)
	return ret
}
