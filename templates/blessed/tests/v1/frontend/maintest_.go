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
		jen.Func().ID("runTestOnAllSupportedBrowsers").Params(jen.ID("t").PointerTo().Qual("testing", "T"), jen.ID("tp").ID("testProvider")).Block(
			jen.For(jen.List(jen.Underscore(), jen.ID("bn")).Assign().Range().Index().String().Values(jen.Lit("firefox"), jen.Lit("chrome"))).Block(
				jen.ID("browserName").Assign().ID("bn"),
				jen.ID("caps").Assign().Qual("github.com/tebeka/selenium", "Capabilities").Values(jen.Lit("browserName").MapAssign().ID("browserName")),
				jen.List(jen.ID("wd"), jen.Err()).Assign().Qual("github.com/tebeka/selenium", "NewRemote").Call(jen.ID("caps"), jen.ID("seleniumHubAddr")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("panic").Call(jen.Err()),
				),
				jen.Line(),
				jen.ID("t").Dot("Run").Call(jen.ID("bn"), jen.ID("tp").Call(jen.ID("wd"))),
				utils.AssertNoError(jen.ID("wd").Dot("Quit").Call(), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("testProvider").Func().Params(jen.ID("driver").Qual("github.com/tebeka/selenium", "WebDriver")).Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestLoginPage").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("runTestOnAllSupportedBrowsers").Call(jen.ID("T"), jen.Func().Params(jen.ID("driver").Qual("github.com/tebeka/selenium", "WebDriver")).Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				jen.Return().Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
					jen.Comment("Navigate to the login page."),
					utils.RequireNoError(jen.ID("driver").Dot("Get").Call(jen.ID("urlToUse").Plus().Lit("/login")), nil),
					jen.Line(),
					jen.Comment("fetch the button."),
					jen.List(jen.ID("elem"), jen.Err()).Assign().ID("driver").Dot("FindElement").Call(jen.Qual("github.com/tebeka/selenium", "ByID"), jen.Lit("loginButton")),
					jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.ID("panic").Call(jen.Err()),
					),
					jen.Line(),
					jen.Comment("check that it is visible."),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("elem").Dot("IsDisplayed").Call(),
					utils.AssertNoError(jen.Err(), nil),
					utils.AssertTrue(jen.ID("actual"), nil),
				),
			)),
		),
		jen.Line(),
	)

	return ret
}
