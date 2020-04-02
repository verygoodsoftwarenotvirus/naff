package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("frontend")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("runTestOnAllSupportedBrowsers").Params(jen.ID("T").ParamPointer().Qual("testing", "T"), jen.ID("tp").ID("testProvider")).Block(
			jen.For(jen.List(jen.ID("_"), jen.ID("bn")).Assign().Range().Index().ID("string").Values(jen.Lit("firefox"), jen.Lit("chrome"))).Block(
				jen.ID("browserName").Assign().ID("bn"),
				jen.ID("caps").Assign().Qual("github.com/tebeka/selenium", "Capabilities").Values(jen.Lit("browserName").MapAssign().ID("browserName")),
				jen.List(jen.ID("wd"), jen.Err()).Assign().Qual("github.com/tebeka/selenium", "NewRemote").Call(jen.ID("caps"), jen.ID("seleniumHubAddr")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("panic").Call(jen.Err()),
				),
				jen.Defer().ID("wd").Dot("Quit").Call(),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.ID("bn"), jen.ID("tp").Call(jen.ID("wd"))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("testProvider").Func().Params(jen.ID("driver").Qual("github.com/tebeka/selenium", "WebDriver")).Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestLoginPage").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("runTestOnAllSupportedBrowsers").Call(jen.ID("T"), jen.Func().Params(jen.ID("driver").Qual("github.com/tebeka/selenium", "WebDriver")).Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.Return().Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
					jen.Comment("Navigate to the login page"),
					jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.ID("driver").Dot("Get").Call(jen.ID("urlToUse").Op("+").Lit("/login"))),
					jen.Line(),
					jen.Comment("fetch the button"),
					jen.List(jen.ID("elem"), jen.Err()).Assign().ID("driver").Dot("FindElement").Call(jen.Qual("github.com/tebeka/selenium", "ByID"), jen.Lit("loginButton")),
					jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.ID("panic").Call(jen.Err()),
					),
					jen.Line(),
					jen.Comment("check that it is visible"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("elem").Dot("IsDisplayed").Call(),
					jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
					jen.Qual("github.com/stretchr/testify/assert", "True").Call(jen.ID("t"), jen.ID("actual")),
				),
			)),
		),
		jen.Line(),
	)
	return ret
}
