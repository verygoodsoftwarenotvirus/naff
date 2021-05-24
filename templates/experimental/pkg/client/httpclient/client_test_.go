package httpclient

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func clientTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestClient_AuthenticatedClient").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("ts").Dot("Client").Call(),
						jen.ID("c").Dot("AuthenticatedClient").Call(),
						jen.Lit("AuthenticatedClient should return the assigned authedClient"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_PlainClient").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("ts").Dot("Client").Call(),
						jen.ID("c").Dot("PlainClient").Call(),
						jen.Lit("PlainClient should return the assigned unauthenticatedClient"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestNewClient").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewClient").Call(
						jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
						jen.ID("UsingLogger").Call(jen.ID("logging").Dot("NewNonOperationalLogger").Call()),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("c"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil URL"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewClient").Call(
						jen.ID("nil"),
						jen.ID("UsingLogger").Call(jen.ID("logging").Dot("NewNonOperationalLogger").Call()),
					),
					jen.ID("require").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("c"),
					),
					jen.ID("require").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_RequestBuilder").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("c").Dot("RequestBuilder").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_BuildURL").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("various urls"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("NewClient").Call(jen.ID("mustParseURL").Call(jen.ID("exampleURI"))),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("testCases").Op(":=").Index().Struct(
						jen.ID("inputQuery").ID("valuer"),
						jen.ID("expectation").ID("string"),
						jen.ID("inputParts").Index().ID("string"),
					).Valuesln(jen.Valuesln(jen.ID("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/things"), jen.ID("inputParts").Op(":").Index().ID("string").Valuesln(jen.Lit("things"))), jen.Valuesln(jen.ID("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/stuff?key=value"), jen.ID("inputQuery").Op(":").Map(jen.ID("string")).Index().ID("string").Valuesln(jen.Lit("key").Op(":").Valuesln(jen.Lit("value"))), jen.ID("inputParts").Op(":").Index().ID("string").Valuesln(jen.Lit("stuff"))), jen.Valuesln(jen.ID("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/things/and/stuff?key=value1&key=value2&yek=eulav"), jen.ID("inputQuery").Op(":").Map(jen.ID("string")).Index().ID("string").Valuesln(jen.Lit("key").Op(":").Valuesln(jen.Lit("value1"), jen.Lit("value2")), jen.Lit("yek").Op(":").Valuesln(jen.Lit("eulav"))), jen.ID("inputParts").Op(":").Index().ID("string").Valuesln(jen.Lit("things"), jen.Lit("and"), jen.Lit("stuff")))),
					jen.For(jen.List(jen.ID("_"), jen.ID("tc")).Op(":=").Range().ID("testCases")).Body(
						jen.ID("actual").Op(":=").ID("c").Dot("BuildURL").Call(
							jen.ID("ctx"),
							jen.ID("tc").Dot("inputQuery").Dot("ToValues").Call(),
							jen.ID("tc").Dot("inputParts").Op("..."),
						),
						jen.ID("assert").Dot("Equal").Call(
							jen.ID("t"),
							jen.ID("tc").Dot("expectation"),
							jen.ID("actual"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid url parts"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("c").Dot("BuildURL").Call(
							jen.ID("ctx"),
							jen.ID("nil"),
							jen.ID("asciiControlChar"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_CloseRequestBody").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("rc").Op(":=").ID("newMockReadCloser").Call(),
					jen.ID("rc").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").Valuesln(jen.ID("Body").Op(":").ID("rc"), jen.ID("StatusCode").Op(":").Qual("net/http", "StatusOK")),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("NewClient").Call(jen.ID("mustParseURL").Call(jen.ID("exampleURI"))),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("c"),
					),
					jen.ID("c").Dot("closeResponseBody").Call(
						jen.ID("ctx"),
						jen.ID("res"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("rc"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildVersionlessURL").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("various urls"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("NewClient").Call(jen.ID("mustParseURL").Call(jen.ID("exampleURI"))),
					jen.ID("testCases").Op(":=").Index().Struct(
						jen.ID("inputQuery").ID("valuer"),
						jen.ID("expectation").ID("string"),
						jen.ID("inputParts").Index().ID("string"),
					).Valuesln(jen.Valuesln(jen.ID("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/things"), jen.ID("inputParts").Op(":").Index().ID("string").Valuesln(jen.Lit("things"))), jen.Valuesln(jen.ID("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/stuff?key=value"), jen.ID("inputQuery").Op(":").Map(jen.ID("string")).Index().ID("string").Valuesln(jen.Lit("key").Op(":").Valuesln(jen.Lit("value"))), jen.ID("inputParts").Op(":").Index().ID("string").Valuesln(jen.Lit("stuff"))), jen.Valuesln(jen.ID("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/things/and/stuff?key=value1&key=value2&yek=eulav"), jen.ID("inputQuery").Op(":").Map(jen.ID("string")).Index().ID("string").Valuesln(jen.Lit("key").Op(":").Valuesln(jen.Lit("value1"), jen.Lit("value2")), jen.Lit("yek").Op(":").Valuesln(jen.Lit("eulav"))), jen.ID("inputParts").Op(":").Index().ID("string").Valuesln(jen.Lit("things"), jen.Lit("and"), jen.Lit("stuff")))),
					jen.For(jen.List(jen.ID("_"), jen.ID("tc")).Op(":=").Range().ID("testCases")).Body(
						jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
						jen.ID("actual").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
							jen.ID("ctx"),
							jen.ID("tc").Dot("inputQuery").Dot("ToValues").Call(),
							jen.ID("tc").Dot("inputParts").Op("..."),
						),
						jen.ID("assert").Dot("Equal").Call(
							jen.ID("t"),
							jen.ID("tc").Dot("expectation"),
							jen.ID("actual"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid url parts"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("actual").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
						jen.ID("ctx"),
						jen.ID("nil"),
						jen.ID("asciiControlChar"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_BuildWebsocketURL").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("NewClient").Call(jen.ID("mustParseURL").Call(jen.ID("exampleURI"))),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expected").Op(":=").Lit("ws://todo.verygoodsoftwarenotvirus.ru/api/v1/things/and/stuff"),
					jen.ID("actual").Op(":=").ID("c").Dot("BuildWebsocketURL").Call(
						jen.ID("ctx"),
						jen.ID("nil"),
						jen.Lit("things"),
						jen.Lit("and"),
						jen.Lit("stuff"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_IsUp").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
				jen.ID("true"),
				jen.Qual("net/http", "MethodGet"),
				jen.Lit(""),
				jen.Lit("/_meta_/ready"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("actual").Op(":=").ID("c").Dot("IsUp").Call(jen.ID("ctx")),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("returns error with invalid url"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("actual").Op(":=").ID("c").Dot("IsUp").Call(jen.ID("ctx")),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with bad status code"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusInternalServerError"),
					),
					jen.ID("actual").Op(":=").ID("c").Dot("IsUp").Call(jen.ID("ctx")),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with timeout"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("actual").Op(":=").ID("c").Dot("IsUp").Call(jen.ID("ctx")),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_fetchAndUnmarshal").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("exampleResponse").Op(":=").Op("&").ID("argleBargle").Valuesln(jen.ID("Name").Op(":").Lit("whatever")),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.Lit("/"),
					),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("exampleResponse"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("spec").Dot("method"),
						jen.ID("ts").Dot("URL"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.Op("&").ID("argleBargle").Valuesln(),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with timeout"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.Lit("/"),
					),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("spec").Dot("method"),
						jen.ID("ts").Dot("URL"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.Op("&").ID("argleBargle").Valuesln(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with 401"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.Lit("/"),
					),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("spec").Dot("method"),
						jen.ID("ts").Dot("URL"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.Qual("errors", "Is").Call(
							jen.ID("c").Dot("fetchAndUnmarshal").Call(
								jen.ID("ctx"),
								jen.ID("req"),
								jen.Op("&").ID("argleBargle").Valuesln(),
							),
							jen.ID("ErrUnauthorized"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with 404"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.Lit("/"),
					),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusNotFound"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("spec").Dot("method"),
						jen.ID("ts").Dot("URL"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.Qual("errors", "Is").Call(
							jen.ID("c").Dot("fetchAndUnmarshal").Call(
								jen.ID("ctx"),
								jen.ID("req"),
								jen.Op("&").ID("argleBargle").Valuesln(),
							),
							jen.ID("ErrNotFound"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with unreadable response"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.Lit("/"),
					),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("exampleResponse"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("spec").Dot("method"),
						jen.ID("ts").Dot("URL"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("fetchAndUnmarshal").Call(
							jen.ID("ctx"),
							jen.ID("req"),
							jen.ID("argleBargle").Valuesln(),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_executeRawRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expectedMethod").Op(":=").Qual("net/http", "MethodPost"),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("expectedMethod"),
						jen.ID("ts").Dot("URL"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("c").Dot("fetchResponseToRequest").Call(
						jen.ID("ctx"),
						jen.Op("&").Qual("net/http", "Client").Valuesln(jen.ID("Timeout").Op(":").Qual("time", "Second")),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("res"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_checkExistence").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expectedMethod").Op(":=").Qual("net/http", "MethodHead"),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.ID("expectedMethod"),
						jen.Lit(""),
						jen.Lit("/"),
					),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("expectedMethod"),
						jen.ID("ts").Dot("URL"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("responseIsOK").Call(
						jen.ID("ctx"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with timeout"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expectedMethod").Op(":=").Qual("net/http", "MethodHead"),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("expectedMethod"),
						jen.ID("ts").Dot("URL"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("c").Dot("authedClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("responseIsOK").Call(
						jen.ID("ctx"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_retrieve").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expectedMethod").Op(":=").Qual("net/http", "MethodPost"),
					jen.ID("exampleResponse").Op(":=").Op("&").ID("argleBargle").Valuesln(jen.ID("Name").Op(":").Lit("whatever")),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.ID("expectedMethod"),
						jen.Lit(""),
						jen.Lit("/"),
					),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("exampleResponse"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("expectedMethod"),
						jen.ID("ts").Dot("URL"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.Op("&").ID("argleBargle").Valuesln(),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil passed in"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.ID("ts").Dot("URL"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with timeout"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expectedMethod").Op(":=").Qual("net/http", "MethodPost"),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("expectedMethod"),
						jen.ID("ts").Dot("URL"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("c").Dot("authedClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
					jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.Op("&").ID("argleBargle").Valuesln(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with 404"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expectedMethod").Op(":=").Qual("net/http", "MethodPost"),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.ID("expectedMethod"),
						jen.Lit(""),
						jen.Lit("/"),
					),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusNotFound"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("expectedMethod"),
						jen.ID("ts").Dot("URL"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.Qual("errors", "Is").Call(
							jen.ID("c").Dot("fetchAndUnmarshal").Call(
								jen.ID("ctx"),
								jen.ID("req"),
								jen.Op("&").ID("argleBargle").Valuesln(),
							),
							jen.ID("ErrNotFound"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_fetchAndUnmarshalWithoutAuthentication").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().ID("expectedMethod").Op("=").Qual("net/http", "MethodPost"),
			jen.ID("exampleResponse").Op(":=").Op("&").ID("argleBargle").Valuesln(jen.ID("Name").Op(":").Lit("whatever")),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.ID("expectedMethod"),
						jen.Lit(""),
						jen.Lit("/"),
					),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("exampleResponse"),
					),
					jen.List(jen.ID("in"), jen.ID("out")).Op(":=").List(jen.Op("&").ID("argleBargle").Valuesln(), jen.Op("&").ID("argleBargle").Valuesln()),
					jen.ID("body").Op(":=").ID("createBodyFromStruct").Call(
						jen.ID("t"),
						jen.ID("in"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("body"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("expectedMethod"),
						jen.ID("ts").Dot("URL"),
						jen.ID("body"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshalWithoutAuthentication").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("out"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with 401"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.ID("expectedMethod"),
						jen.Lit(""),
						jen.Lit("/"),
					),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.List(jen.ID("in"), jen.ID("out")).Op(":=").List(jen.Op("&").ID("argleBargle").Valuesln(), jen.Op("&").ID("argleBargle").Valuesln()),
					jen.ID("body").Op(":=").ID("createBodyFromStruct").Call(
						jen.ID("t"),
						jen.ID("in"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("body"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("expectedMethod"),
						jen.ID("ts").Dot("URL"),
						jen.ID("body"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshalWithoutAuthentication").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("out"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.Qual("errors", "Is").Call(
							jen.ID("err"),
							jen.ID("ErrUnauthorized"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with 404"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.ID("expectedMethod"),
						jen.Lit(""),
						jen.Lit("/"),
					),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusNotFound"),
					),
					jen.List(jen.ID("in"), jen.ID("out")).Op(":=").List(jen.Op("&").ID("argleBargle").Valuesln(), jen.Op("&").ID("argleBargle").Valuesln()),
					jen.ID("body").Op(":=").ID("createBodyFromStruct").Call(
						jen.ID("t"),
						jen.ID("in"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("body"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("expectedMethod"),
						jen.ID("ts").Dot("URL"),
						jen.ID("body"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshalWithoutAuthentication").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("out"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.Qual("errors", "Is").Call(
							jen.ID("err"),
							jen.ID("ErrNotFound"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with timeout"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("in"), jen.ID("out")).Op(":=").List(jen.Op("&").ID("argleBargle").Valuesln(), jen.Op("&").ID("argleBargle").Valuesln()),
					jen.ID("body").Op(":=").ID("createBodyFromStruct").Call(
						jen.ID("t"),
						jen.ID("in"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("body"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("expectedMethod"),
						jen.ID("ts").Dot("URL"),
						jen.ID("body"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("c").Dot("unauthenticatedClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("fetchAndUnmarshalWithoutAuthentication").Call(
							jen.ID("ctx"),
							jen.ID("req"),
							jen.ID("out"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil as output"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("in").Op(":=").Op("&").ID("argleBargle").Valuesln(),
					jen.ID("body").Op(":=").ID("createBodyFromStruct").Call(
						jen.ID("t"),
						jen.ID("in"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("body"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("expectedMethod"),
						jen.ID("ts").Dot("URL"),
						jen.ID("body"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshalWithoutAuthentication").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("testingType").Valuesln(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with unreadable response"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.ID("expectedMethod"),
						jen.Lit(""),
						jen.Lit("/"),
					),
					jen.List(jen.ID("c"), jen.ID("ts")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("exampleResponse"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.ID("expectedMethod"),
						jen.ID("ts").Dot("URL"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("fetchAndUnmarshalWithoutAuthentication").Call(
							jen.ID("ctx"),
							jen.ID("req"),
							jen.ID("argleBargle").Valuesln(),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
