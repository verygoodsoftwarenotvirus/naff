package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func helpersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("chromeDisabled").ID("bool"),
			jen.ID("firefoxDisabled").ID("bool"),
			jen.ID("webkitDisabled").ID("bool"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.ID("chromeDisabled").Op("=").Qual("strings", "EqualFold").Call(
				jen.Qual("strings", "TrimSpace").Call(jen.Qual("os", "Getenv").Call(jen.Lit("CHROME_DISABLED"))),
				jen.Lit("y"),
			),
			jen.ID("firefoxDisabled").Op("=").Qual("strings", "EqualFold").Call(
				jen.Qual("strings", "TrimSpace").Call(jen.Qual("os", "Getenv").Call(jen.Lit("FIREFOX_DISABLED"))),
				jen.Lit("y"),
			),
			jen.ID("webkitDisabled").Op("=").Qual("strings", "EqualFold").Call(
				jen.Qual("strings", "TrimSpace").Call(jen.Qual("os", "Getenv").Call(jen.Lit("WEBKIT_DISABLED"))),
				jen.Lit("y"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("stringPointer").Params(jen.ID("s").ID("string")).Params(jen.Op("*").ID("string")).Body(
			jen.Return().Op("&").ID("s")),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("testHelper").Struct(
				jen.ID("pw").Op("*").ID("playwright").Dot("Playwright"),
				jen.List(jen.ID("Firefox"), jen.ID("Chrome"), jen.ID("Webkit")).ID("playwright").Dot("Browser"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("setupTestHelper").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("testHelper")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.List(jen.ID("pw"), jen.ID("err")).Op(":=").ID("playwright").Dot("Run").Call(),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
				jen.Lit("could not start playwright"),
			),
			jen.ID("th").Op(":=").Op("&").ID("testHelper").Valuesln(jen.ID("pw").Op(":").ID("pw")),
			jen.If(jen.Op("!").ID("chromeDisabled")).Body(
				jen.List(jen.ID("th").Dot("Chrome"), jen.ID("err")).Op("=").ID("pw").Dot("Chromium").Dot("Launch").Call(),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
					jen.Lit("could not launch browser"),
				),
				jen.ID("require").Dot("NotNil").Call(
					jen.ID("t"),
					jen.ID("th").Dot("Chrome"),
				),
			),
			jen.If(jen.Op("!").ID("firefoxDisabled")).Body(
				jen.List(jen.ID("th").Dot("Firefox"), jen.ID("err")).Op("=").ID("pw").Dot("Firefox").Dot("Launch").Call(),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
					jen.Lit("could not launch browser"),
				),
				jen.ID("require").Dot("NotNil").Call(
					jen.ID("t"),
					jen.ID("th").Dot("Firefox"),
				),
			),
			jen.If(jen.Op("!").ID("webkitDisabled")).Body(
				jen.List(jen.ID("th").Dot("Webkit"), jen.ID("err")).Op("=").ID("pw").Dot("WebKit").Dot("Launch").Call(),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
					jen.Lit("could not launch browser"),
				),
				jen.ID("require").Dot("NotNil").Call(
					jen.ID("t"),
					jen.ID("th").Dot("Webkit"),
				),
			),
			jen.Return().ID("th"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("h").Op("*").ID("testHelper")).ID("runForAllBrowsers").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("testName").ID("string"), jen.ID("testFunc").Func().Params(jen.ID("playwright").Dot("Browser")).Params(jen.Func().Params(jen.Op("*").Qual("testing", "T")))).Body(
			jen.If(jen.Op("!").ID("chromeDisabled")).Body(
				jen.ID("t").Dot("Run").Call(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s with chrome"),
						jen.ID("testName"),
					),
					jen.ID("testFunc").Call(jen.ID("h").Dot("Chrome")),
				)),
			jen.If(jen.Op("!").ID("firefoxDisabled")).Body(
				jen.ID("t").Dot("Run").Call(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s with firefox"),
						jen.ID("testName"),
					),
					jen.ID("testFunc").Call(jen.ID("h").Dot("Firefox")),
				)),
			jen.If(jen.Op("!").ID("webkitDisabled")).Body(
				jen.ID("t").Dot("Run").Call(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s with webkit"),
						jen.ID("testName"),
					),
					jen.ID("testFunc").Call(jen.ID("h").Dot("Webkit")),
				)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("boolPointer").Params(jen.ID("b").ID("bool")).Params(jen.Op("*").ID("bool")).Body(
			jen.Return().Op("&").ID("b")),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("artifactsDir").Op("=").Lit("/home/vgsnv/src/gitlab.com/verygoodsoftwarenotvirus/todo/artifacts"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("saveScreenshotTo").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("page").ID("playwright").Dot("Page"), jen.ID("path").ID("string")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("opts").Op(":=").ID("playwright").Dot("PageScreenshotOptions").Valuesln(jen.ID("FullPage").Op(":").ID("boolPointer").Call(jen.ID("true")), jen.ID("Path").Op(":").ID("stringPointer").Call(jen.Qual("path/filepath", "Join").Call(
				jen.ID("artifactsDir"),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.png"),
					jen.ID("path"),
				),
			)), jen.ID("Type").Op(":").ID("playwright").Dot("ScreenshotTypePng")),
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("page").Dot("Screenshot").Call(jen.ID("opts")),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("getScreenshotBytes").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("ss").ID("playwright").Dot("ElementHandle")).Params(jen.Index().ID("byte")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("opts").Op(":=").ID("playwright").Dot("ElementHandleScreenshotOptions").Valuesln(jen.ID("Type").Op(":").ID("playwright").Dot("ScreenshotTypePng")),
			jen.List(jen.ID("data"), jen.ID("err")).Op(":=").ID("ss").Dot("Screenshot").Call(jen.ID("opts")),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().ID("data"),
		),
		jen.Line(),
	)

	return code
}
