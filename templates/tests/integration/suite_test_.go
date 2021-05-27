package integration

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func suiteTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("cookieAuthType").Op("=").Lit("cookie"),
			jen.ID("pasetoAuthType").Op("=").Lit("PASETO"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("globalClientExceptions").Index().ID("string"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("testClientWrapper").Struct(
				jen.ID("main").Op("*").ID("httpclient").Dot("Client"),
				jen.ID("admin").Op("*").ID("httpclient").Dot("Client"),
				jen.ID("authType").ID("string"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestIntegration").Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("suite").Dot("Run").Call(
				jen.ID("t"),
				jen.ID("new").Call(jen.ID("TestSuite")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("TestSuite").Struct(
				jen.ID("suite").Dot("Suite"),
				jen.ID("ctx").Qual("context", "Context"),
				jen.ID("user").Op("*").ID("types").Dot("User"),
				jen.ID("cookie").Op("*").Qual("net/http", "Cookie"),
				jen.List(jen.ID("cookieClient"), jen.ID("pasetoClient"), jen.ID("adminCookieClient"), jen.ID("adminPASETOClient")).Op("*").ID("httpclient").Dot("Client"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("suite").Dot("SetupTestSuite").Op("=").Parens(jen.Op("*").ID("TestSuite")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("SetupTest").Params().Body(
			jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
			jen.ID("testName").Op(":=").ID("t").Dot("Name").Call(),
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
				jen.Qual("context", "Background").Call(),
				jen.ID("testName"),
			),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.List(jen.ID("s").Dot("ctx"), jen.ID("_")).Op("=").ID("tracing").Dot("StartCustomSpan").Call(
				jen.ID("ctx"),
				jen.ID("testName"),
			),
			jen.List(jen.ID("s").Dot("user"), jen.ID("s").Dot("cookie"), jen.ID("s").Dot("cookieClient"), jen.ID("s").Dot("pasetoClient")).Op("=").ID("createUserAndClientForTest").Call(
				jen.ID("s").Dot("ctx"),
				jen.ID("t"),
			),
			jen.List(jen.ID("s").Dot("adminCookieClient"), jen.ID("s").Dot("adminPASETOClient")).Op("=").ID("buildAdminCookieAndPASETOClients").Call(
				jen.ID("s").Dot("ctx"),
				jen.ID("t"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("runForEachClientExcept").Params(jen.ID("name").ID("string"), jen.ID("subtestBuilder").Func().Params(jen.Op("*").ID("testClientWrapper")).Params(jen.Func().Params()), jen.ID("exceptions").Op("...").ID("string")).Body(
			jen.For(jen.List(jen.ID("a"), jen.ID("c")).Op(":=").Range().ID("s").Dot("eachClientExcept").Call(jen.ID("exceptions").Op("..."))).Body(
				jen.List(jen.ID("authType"), jen.ID("testClients")).Op(":=").List(jen.ID("a"), jen.ID("c")),
				jen.ID("s").Dot("Run").Call(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s via %s"),
						jen.ID("name"),
						jen.ID("authType"),
					),
					jen.ID("subtestBuilder").Call(jen.ID("testClients")),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("eachClientExcept").Params(jen.ID("exceptions").Op("...").ID("string")).Params(jen.Map(jen.ID("string")).Op("*").ID("testClientWrapper")).Body(
			jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
			jen.ID("clients").Op(":=").Map(jen.ID("string")).Op("*").ID("testClientWrapper").Valuesln(jen.ID("cookieAuthType").Op(":").Valuesln(jen.ID("authType").Op(":").ID("cookieAuthType"), jen.ID("main").Op(":").ID("s").Dot("cookieClient"), jen.ID("admin").Op(":").ID("s").Dot("adminCookieClient")), jen.ID("pasetoAuthType").Op(":").Valuesln(jen.ID("authType").Op(":").ID("pasetoAuthType"), jen.ID("main").Op(":").ID("s").Dot("pasetoClient"), jen.ID("admin").Op(":").ID("s").Dot("adminPASETOClient"))),
			jen.For(jen.List(jen.ID("_"), jen.ID("name")).Op(":=").Range().ID("exceptions")).Body(
				jen.ID("delete").Call(
					jen.ID("clients"),
					jen.ID("name"),
				)),
			jen.For(jen.List(jen.ID("_"), jen.ID("name")).Op(":=").Range().ID("globalClientExceptions")).Body(
				jen.ID("delete").Call(
					jen.ID("clients"),
					jen.ID("name"),
				)),
			jen.ID("require").Dot("NotEmpty").Call(
				jen.ID("t"),
				jen.ID("clients"),
			),
			jen.Return().ID("clients"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("suite").Dot("WithStats").Op("=").Parens(jen.Op("*").ID("TestSuite")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("checkTestRunsForPositiveResultsThatOccurredTooQuickly").Params(jen.ID("stats").Op("*").ID("suite").Dot("SuiteInformation")).Body(
			jen.Var().Defs(
				jen.ID("minimumTestThreshold").Op("=").Lit(1).Op("*").Qual("time", "Millisecond"),
			),
			jen.If(jen.ID("stats").Dot("Passed").Call()).Body(
				jen.For(jen.List(jen.ID("testName"), jen.ID("stat")).Op(":=").Range().ID("stats").Dot("TestStats")).Body(
					jen.If(jen.ID("stat").Dot("End").Dot("Sub").Call(jen.ID("stat").Dot("Start")).Op("<").ID("minimumTestThreshold").Op("&&").ID("stat").Dot("Passed")).Body(
						jen.ID("s").Dot("T").Call().Dot("Fatalf").Call(
							jen.Lit("suspiciously quick test execution time: %q"),
							jen.ID("testName"),
						)))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("HandleStats").Params(jen.ID("_").ID("string"), jen.ID("stats").Op("*").ID("suite").Dot("SuiteInformation")).Body(
			jen.Var().Defs(
				jen.ID("totalExpectedTestCount").Op("=").Lit(69),
			),
			jen.ID("s").Dot("checkTestRunsForPositiveResultsThatOccurredTooQuickly").Call(jen.ID("stats")),
			jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "AssertAppropriateNumberOfTestsRan").Call(
				jen.ID("s").Dot("T").Call(),
				jen.ID("totalExpectedTestCount"),
				jen.ID("stats"),
			),
		),
		jen.Line(),
	)

	return code
}
