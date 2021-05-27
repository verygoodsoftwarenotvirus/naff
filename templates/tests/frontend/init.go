package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func initDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("urlToUse").ID("string"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.ID("u").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "DetermineServiceURL").Call(),
			jen.ID("urlToUse").Op("=").ID("u").Dot("String").Call(),
			jen.ID("logger").Op(":=").ID("logging").Dot("ProvideLogger").Call(jen.ID("logging").Dot("Config").Valuesln(jen.ID("Provider").Op(":").ID("logging").Dot("ProviderZerolog"))),
			jen.ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("URLKey"),
				jen.ID("urlToUse"),
			).Dot("Info").Call(jen.Lit("checking server")),
			jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "EnsureServerIsUp").Call(
				jen.Qual("context", "Background").Call(),
				jen.ID("urlToUse"),
			),
			jen.ID("fiftySpaces").Op(":=").Qual("strings", "Repeat").Call(
				jen.Lit("\n"),
				jen.Lit(50),
			),
			jen.Qual("fmt", "Printf").Call(
				jen.Lit("%s\tRunning tests%s"),
				jen.ID("fiftySpaces"),
				jen.ID("fiftySpaces"),
			),
		),
		jen.Line(),
	)

	return code
}
