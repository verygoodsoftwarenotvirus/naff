package httpclient

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func optionsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestClient_SetOptions").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("expectedURL"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.Lit("https://notarealplace.lol")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("NotEqual").Call(
						jen.ID("t"),
						jen.ID("expectedURL"),
						jen.ID("c").Dot("URL").Call(),
						jen.Lit("expected and actual URLs match somehow"),
					),
					jen.ID("exampleOption").Op(":=").Func().Params(jen.ID("client").Op("*").ID("Client")).Params(jen.ID("error")).Body(
						jen.ID("client").Dot("url").Op("=").ID("expectedURL"),
						jen.Return().ID("nil"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("SetOptions").Call(jen.ID("exampleOption")),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedURL"),
						jen.ID("c").Dot("URL").Call(),
						jen.Lit("expected and actual URLs do not match"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("exampleOption").Op(":=").Func().Params(jen.ID("client").Op("*").ID("Client")).Params(jen.ID("error")).Body(
						jen.Return().Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("SetOptions").Call(jen.ID("exampleOption")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestUsingJSON").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewClient").Call(
						jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
						jen.ID("UsingJSON").Call(),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("c"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Lit("application/json"),
						jen.ID("c").Dot("encoder").Dot("ContentType").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestUsingXML").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewClient").Call(
						jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
						jen.ID("UsingXML").Call(),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("c"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Lit("application/xml"),
						jen.ID("c").Dot("encoder").Dot("ContentType").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestUsingLogger").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("expectedURL"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.Lit("https://todo.verygoodsoftwarenotvirus.ru")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewClient").Call(
						jen.ID("expectedURL"),
						jen.ID("UsingLogger").Call(jen.ID("logging").Dot("NewNoopLogger").Call()),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("c"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestUsingDebug").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewClient").Call(
						jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
						jen.ID("UsingDebug").Call(jen.ID("true")),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("c"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("true"),
						jen.ID("c").Dot("debug"),
						jen.Lit("REPLACE ME"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestUsingTimeout").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("expected").Op(":=").Qual("time", "Minute"),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewClient").Call(
						jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
						jen.ID("UsingTimeout").Call(jen.ID("expected")),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("c"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("c").Dot("authedClient").Dot("Timeout"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with fallback to default timeout"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewClient").Call(
						jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
						jen.ID("UsingTimeout").Call(jen.Lit(0)),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("c"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("defaultTimeout"),
						jen.ID("c").Dot("authedClient").Dot("Timeout"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestUsingCookie").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").Op("&").Qual("net/http", "Cookie").Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewClient").Call(
						jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
						jen.ID("UsingCookie").Call(jen.ID("exampleInput")),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("c"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("c").Dot("authMethod").Op("==").ID("cookieAuthMethod"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil cooki9e"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewClient").Call(
						jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
						jen.ID("UsingCookie").Call(jen.ID("nil")),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("c"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestUsingPASETO").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewClient").Call(
						jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
						jen.ID("UsingPASETO").Call(
							jen.ID("t").Dot("Name").Call(),
							jen.Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()),
						),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("c"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("c").Dot("authMethod").Op("==").ID("pasetoAuthMethod"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
