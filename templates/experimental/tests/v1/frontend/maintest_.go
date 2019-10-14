package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func mainTestDotGo() *jen.File {
	ret := jen.NewFile("frontend")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("runTestOnAllSupportedBrowsers").Params(jen.ID("T").Op("*").Qual("testing", "T"), jen.ID("tp").ID("testProvider")).Block(
		jen.For(jen.List(jen.ID("_"), jen.ID("bn")).Op(":=").Range().Index().ID("string").Valuesln(jen.Lit("firefox"), jen.Lit("chrome"))).Block(
			jen.ID("browserName").Op(":=").ID("bn"),
			jen.ID("caps").Op(":=").ID("selenium").Dot(
				"Capabilities",
			).Valuesln(jen.Lit("browserName").Op(":").ID("browserName")),
			jen.List(jen.ID("wd"), jen.ID("err")).Op(":=").ID("selenium").Dot(
				"NewRemote",
			).Call(jen.ID("caps"), jen.ID("seleniumHubAddr")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Defer().ID("wd").Dot(
				"Quit",
			).Call(),
			jen.ID("T").Dot(
				"Run",
			).Call(jen.ID("bn"), jen.ID("tp").Call(jen.ID("wd"))),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("testProvider").Params(jen.ID("driver").ID("selenium").Dot(
		"WebDriver",
	)).Params(jen.Params(jen.ID("t").Op("*").Qual("testing", "T"))),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestLoginPage").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("runTestOnAllSupportedBrowsers").Call(jen.ID("T"), jen.Func().Params(jen.ID("driver").ID("selenium").Dot(
			"WebDriver",
		)).Params(jen.Params(jen.ID("t").Op("*").Qual("testing", "T"))).Block(
			jen.Return().Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("require").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("driver").Dot(
					"Get",
				).Call(jen.ID("urlToUse").Op("+").Lit("/login"))),
				jen.List(jen.ID("elem"), jen.ID("err")).Op(":=").ID("driver").Dot(
					"FindElement",
				).Call(jen.ID("selenium").Dot(
					"ByID",
				), jen.Lit("loginButton")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("panic").Call(jen.ID("err")),
				),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("elem").Dot(
					"IsDisplayed",
				).Call(),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot(
					"True",
				).Call(jen.ID("t"), jen.ID("actual")),
			),
		)),
	),

		jen.Line(),
	)
	return ret
}
