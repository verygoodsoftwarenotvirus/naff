package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func initDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("debug").Op("=").ID("true"),
			jen.ID("nonexistentID").ID("uint64").Op("=").Qual("math", "MaxUint32"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("urlToUse").ID("string"),
			jen.ID("parsedURLToUse").Op("*").Qual("net/url", "URL"),
			jen.ID("premadeAdminUser").Op("=").Op("&").ID("types").Dot("User").Valuesln(jen.ID("ID").Op(":").Lit(1), jen.ID("TwoFactorSecret").Op(":").Lit("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="), jen.ID("Username").Op(":").Lit("exampleUser"), jen.ID("HashedPassword").Op(":").Lit("integration-tests-are-cool")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.Qual("context", "Background").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("parsedURLToUse").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "DetermineServiceURL").Call(),
			jen.ID("urlToUse").Op("=").ID("parsedURLToUse").Dot("String").Call(),
			jen.ID("logger").Op(":=").ID("logging").Dot("ProvideLogger").Call(jen.ID("logging").Dot("Config").Valuesln(jen.ID("Provider").Op(":").ID("logging").Dot("ProviderZerolog"))),
			jen.ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("URLKey"),
				jen.ID("urlToUse"),
			).Dot("Info").Call(jen.Lit("checking server")),
			jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "EnsureServerIsUp").Call(
				jen.ID("ctx"),
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
