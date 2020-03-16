package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("frontend")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Func().ID("runTestOnAllSupportedBrowsers").Params(jen.ID("T").Op("*").Qual("testing", "T"), jen.ID("tp").ID("testProvider")).Block(
			jen.For(jen.List(jen.ID("_"), jen.ID("bn")).Op(":=").Range().Index().ID("string").Values(jen.Lit("firefox"), jen.Lit("chrome"))).Block(
				jen.ID("browserName").Op(":=").ID("bn"),
				jen.ID("caps").Op(":=").Qual("github.com/tebeka/selenium", "Capabilities").Values(jen.Lit("browserName").Op(":").ID("browserName")),
				jen.List(jen.ID("wd"), jen.ID("err")).Op(":=").Qual("github.com/tebeka/selenium", "NewRemote").Call(jen.ID("caps"), jen.ID("seleniumHubAddr")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("panic").Call(jen.ID("err")),
				),
				jen.Defer().ID("wd").Dot("Quit").Call(),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.ID("bn"), jen.ID("tp").Call(jen.ID("wd"))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("testProvider").Func().Params(jen.ID("driver").Qual("github.com/tebeka/selenium", "WebDriver")).Func().Params(jen.ID("t").Op("*").Qual("testing", "T")),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestLoginPage").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("runTestOnAllSupportedBrowsers").Call(jen.ID("T"), jen.Func().Params(jen.ID("driver").Qual("github.com/tebeka/selenium", "WebDriver")).Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.Return().Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.Comment("Navigate to the login page"),
					jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.ID("driver").Dot("Get").Call(jen.ID("urlToUse").Op("+").Lit("/login"))),
					jen.Line(),
					jen.Comment("fetch the button"),
					jen.List(jen.ID("elem"), jen.ID("err")).Op(":=").ID("driver").Dot("FindElement").Call(jen.Qual("github.com/tebeka/selenium", "ByID"), jen.Lit("loginButton")),
					jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
						jen.ID("panic").Call(jen.ID("err")),
					),
					jen.Line(),
					jen.Comment("check that it is visible"),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("elem").Dot("IsDisplayed").Call(),
					jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
					jen.Qual("github.com/stretchr/testify/assert", "True").Call(jen.ID("t"), jen.ID("actual")),
				),
			)),
		),
		jen.Line(),
	)
	return ret
}
