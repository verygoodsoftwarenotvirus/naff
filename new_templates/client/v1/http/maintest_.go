package client

import jen "github.com/dave/jennifer/jen"

func mainTestDotGo() *jen.File {
	ret := jen.NewFile("client")
	ret.Add(jen.Null())

	addImports(ret)

	// vars
	ret.Add(jen.Var().Id("exampleURI").Op("=").Lit("https://todo.verygoodsoftwarenotvirus.ru"))

	// types
	ret.Add(
		jen.Type().Id("argleBargle").Struct(
			jen.Id("Name").Id("string"),
		),
		jen.Line(),
	)
	ret.Add(jen.Type().Id("valuer").Map(jen.Id("string")).Index().Id("string"))

	ret.Add(
		jen.Func().Params(
			jen.Id("v").Id("valuer"),
		).Id("ToValues").Params().Params(
			jen.Qual("net/url", "Values"),
		).Block(
			jen.Return().Qual("net/url", "Values").Call(
				jen.Id("v"),
			),
		),
		jen.Line(),
	)

	// funcs
	ret.Add(
		jen.Func().Id("mustParseURL").Params(
			jen.Id("uri").Id("string"),
		).Params(
			jen.Op("*").Qual("net/url", "URL"),
		).Block(
			jen.List(
				jen.Id("u"),
				jen.Id("err"),
			).Op(":=").Qual("net/url", "Parse").Call(
				jen.Id("uri"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Id("panic").Call(
					jen.Id("err"),
				),
			),
			jen.Return().Id("u"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("buildTestClient").Params(
			jen.Id("t").Op("*").Qual("testing", "T"),
			jen.Id("ts").Op("*").Qual("net/http/httptest", "Server"),
		).Params(
			jen.Op("*").Id(v1),
		).Block(
			jen.Id("t").Dot("Helper").Call(),
			jen.Id("l").Op(":=").Qual(noopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.Id("u").Op(":=").Id("mustParseURL").Call(
				jen.Id("ts").Dot("URL"),
			),
			jen.Return().Op("&").Id(v1).Values(jen.Dict{
				jen.Id("URL"):          jen.Id("u"),
				jen.Id("plainClient"):  jen.Id("ts").Dot("Client").Call(),
				jen.Id("logger"):       jen.Id("l"),
				jen.Id("Debug"):        jen.Id("true"),
				jen.Id("authedClient"): jen.Id("ts").Dot("Client").Call(),
			}),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_AuthenticatedClient").Params(jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"obligatory",
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Id("actual").Op(":=").Id("c").Dot("AuthenticatedClient").Call(),
				assertEqual(
					jen.Id("t"),
					jen.Id("ts").Dot("Client").Call(),
					jen.Id("actual"),
					jen.Lit("AuthenticatedClient should return the assigned authedClient"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_PlainClient").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"obligatory",
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Id("actual").Op(":=").Id("c").Dot("PlainClient").Call(),
				assertEqual(
					jen.Id("t"),
					jen.Id("ts").Dot("Client").Call(),
					jen.Id("actual"),
					jen.Lit("PlainClient should return the assigned plainClient"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_TokenSource").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"obligatory",
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.List(
					jen.Id("c"),
					jen.Id("err"),
				).Op(":=").Id("NewClient").Call(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.Id("mustParseURL").Call(
						jen.Id("exampleURI"),
					),
					jen.Qual(noopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Id("ts").Dot("Client").Call(),
					jen.Index().Id("string").Values(jen.Lit("*")),
					jen.Id("false"),
				),
				requireNoError(jen.Id("err"), nil),
				jen.Id("actual").Op(":=").Id("c").Dot("TokenSource").Call(),
				assertNotNil(
					jen.Id("t"),
					jen.Id("actual"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestNewClient").Params(jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.List(
					jen.Id("c"),
					jen.Id("err"),
				).Op(":=").Id("NewClient").Call(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.Id("mustParseURL").Call(
						jen.Id("exampleURI"),
					),
					jen.Qual(noopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Id("ts").Dot("Client").Call(),
					jen.Index().Id("string").Values(jen.Lit("*")),
					jen.Id("false"),
				),
				requireNotNil(jen.Id("c"), nil),
				requireNoError(jen.Id("err"), nil),
			),

			buildSubTest(
				"with client but invalid timeout",
				jen.List(
					jen.Id("c"),
					jen.Id("err"),
				).Op(":=").Id("NewClient").Call(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.Id("mustParseURL").Call(
						jen.Id("exampleURI"),
					),
					jen.Qual(noopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Op("&").Qual("net/http", "Client").Values(
						jen.Id("Timeout").Op(":").Lit(0),
					),
					jen.Index().Id("string").Values(jen.Lit("*")),
					jen.Id("true"),
				),

				requireNotNil(jen.Id("c"), nil),
				requireNoError(jen.Id("err"), nil),
				assertEqual(
					jen.Id("t"),
					jen.Id("c").Dot("plainClient").Dot("Timeout"),
					jen.Id("defaultTimeout"),
					jen.Lit("NewClient should set the default timeout"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("NewSimpleClient").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"obligatory",
				jen.List(
					jen.Id("c"),
					jen.Id("err"),
				).Op(":=").Id("NewSimpleClient").Call(
					jen.Qual("context", "Background").Call(),
					jen.Id("mustParseURL").Call(
						jen.Id("exampleURI"),
					),
					jen.Id("true"),
				),
				assertNotNil(jen.Id("c"), nil),
				assertNoError(jen.Id("err"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_executeRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"with error",
				expectMethod("expectedMethod", "MethodPost"),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
						jen.Func().Params(jen.Id("res").Qual("net/http", "ResponseWriter"),
							jen.Id("req").Op("*").Qual("net/http", "Request")).Block(
							assertEqual(
								jen.Id("t"),
								jen.Id("req").Dot("Method"),
								jen.Id("expectedMethod"),
							),
							jen.Qual("time", "Sleep").Call(
								jen.Lit(10).Op("*").Qual("time", "Hour"),
							),
						),
					),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("req"),
					jen.Id("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Id("nil"),
				),
				requireNotNil(jen.Id("req"), nil),
				requireNoError(jen.Id("err"), nil),
				jen.List(
					jen.Id("res"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("executeRawRequest").Call(
					jen.Id("ctx"), jen.Op("&").Qual("net/http", "Client").Values(
						jen.Id("Timeout").Op(":").Qual("time", "Second"),
					),
					jen.Id("req"),
				),
				assertNil(
					jen.Id("t"),
					jen.Id("res"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("BuildURL").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"various urls",
				jen.Id("t").Dot("Parallel").Call(),
				jen.List(
					jen.Id("u"),
					jen.Id("_"),
				).Op(":=").Qual("net/url", "Parse").Call(
					jen.Id("exampleURI"),
				),
				jen.List(
					jen.Id("c"),
					jen.Id("err"),
				).Op(":=").Id("NewClient").Call(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""), jen.Lit(""),
					jen.Id("u"),
					jen.Qual(noopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Id("nil"),
					jen.Index().Id("string").Values(jen.Lit("*")),
					jen.Id("false"),
				),
				requireNoError(jen.Id("err"), nil),
				jen.Id("testCases").Op(":=").Index().Struct(
					jen.Id("expectation").Id("string"),
					jen.Id("inputParts").Index().Id("string"),
					jen.Id("inputQuery").Id("valuer"),
				).Values(
					jen.Values(jen.Dict{
						jen.Id("expectation"): jen.Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/things"),
						jen.Id("inputParts"):  jen.Index().Id("string").Values(jen.Lit("things")),
					}),
					jen.Values(jen.Dict{
						jen.Id("expectation"): jen.Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/stuff?key=value"),
						jen.Id("inputQuery"): jen.Map(jen.Id("string")).Index().Id("string").Values(jen.Dict{
							jen.Lit("key"): jen.Values(jen.Lit("value")),
						}),
						jen.Id("inputParts"): jen.Index().Id("string").Values(jen.Lit("stuff")),
					}),
					jen.Values(jen.Dict{
						jen.Id("expectation"): jen.Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/things/and/stuff?key=value1&key=value2&yek=eulav"),
						jen.Id("inputQuery"): jen.Map(jen.Id("string")).Index().Id("string").Values(jen.Dict{
							jen.Lit("key"): jen.Values(jen.Lit("value1"), jen.Lit("value2")),
							jen.Lit("yek"): jen.Values(jen.Lit("eulav")),
						}),
						jen.Id("inputParts"): jen.Index().Id("string").Values(jen.Lit("things"), jen.Lit("and"), jen.Lit("stuff")),
					}),
				),
				jen.For(jen.List(
					jen.Id("_"),
					jen.Id("tc"),
				).Op(":=").Range().Id("testCases"),
				).Block(
					jen.Id("actual").Op(":=").Id("c").Dot("BuildURL").Call(
						jen.Id("tc").Dot("inputQuery").Dot("ToValues").Call(),
						jen.Id("tc").Dot("inputParts").Op("..."),
					),
					assertEqual(
						jen.Id("t"),
						jen.Id("tc").Dot("expectation"),
						jen.Id("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildWebsocketURL").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.List(
					jen.Id("u"),
					jen.Id("_"),
				).Op(":=").Qual("net/url", "Parse").Call(
					jen.Id("exampleURI"),
				),
				jen.List(
					jen.Id("c"),
					jen.Id("err"),
				).Op(":=").Id("NewClient").Call(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.Id("u"),
					jen.Qual(noopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Id("nil"),
					jen.Index().Id("string").Values(jen.Lit("*")),
					jen.Id("false"),
				),
				requireNoError(jen.Id("err"), nil),
				jen.Id("expected").Op(":=").Lit("ws://todo.verygoodsoftwarenotvirus.ru/api/v1/things/and/stuff"),
				jen.Id("actual").Op(":=").Id("c").Dot("BuildWebsocketURL").Call(
					jen.Lit("things"),
					jen.Lit("and"),
					jen.Lit("stuff"),
				),
				assertEqual(jen.Id("expected"), jen.Id("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildHealthCheckRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildHealthCheckRequest").Call(),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("actual").Dot("Method"),
					jen.Id("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.Id("expectedMethod"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_IsUp").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet"),
						),
						writeHeader("StatusOK"),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Id("actual").Op(":=").Id("c").Dot("IsUp").Call(),
				assertTrue(
					jen.Id("t"),
					jen.Id("actual"),
				),
			),
			jen.Line(),
			buildSubTest(
				"with bad status code",
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet"),
						),
						writeHeader("StatusInternalServerError"),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Id("actual").Op(":=").Id("c").Dot("IsUp").Call(),
				assertFalse(
					jen.Id("t"),
					jen.Id("actual"),
				),
			),
			jen.Line(),
			buildSubTest(
				"with timeout",
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Qual("net/http", "MethodGet"),
						),
						jen.Qual("time", "Sleep").Call(
							jen.Lit(10).Op("*").Qual("time", "Hour"),
						),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Id("c").Dot("plainClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
				jen.Id("actual").Op(":=").Id("c").Dot("IsUp").Call(),
				assertFalse(
					jen.Id("t"),
					jen.Id("actual"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_buildDataRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				expectMethod("expectedMethod", "MethodPost"),
				jen.List(
					jen.Id("req"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("buildDataRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Op("&").Id("testingType").Values(jen.Id("Name").Op(":").Lit("name")),
				),
				requireNotNil(jen.Id("req"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("expectedMethod"),
					jen.Id("req").Dot("Method"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_makeRequest").Params(jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				createCtx(),
				expectMethod("expectedMethod", "MethodPost"),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Id("expectedMethod"),
						),
						jen.Id("require").Dot("NoError").Call(
							jen.Id("t"),
							jen.Qual("encoding/json", "NewEncoder").Call(
								jen.Id("res"),
							).Dot("Encode").Call(
								jen.Op("&").Id("argleBargle").Values(
									jen.Id("Name").Op(":").Lit("name"),
								),
							),
						),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("req"),
					jen.Id("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Id("nil"),
				),
				requireNotNil(jen.Id("req"), nil),
				requireNoError(jen.Id("err"), nil),
				jen.Line(),
				jen.Id("err").Op("=").Id("c").Dot("executeRequest").Call(
					jen.Id("ctx"),
					jen.Id("req"),
					jen.Op("&").Id("argleBargle").Values(),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
				),
			),
			jen.Line(),
			buildSubTest(
				"with 404",
				jen.Line(),
				createCtx(),
				expectMethod("expectedMethod", "MethodPost"),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Id("expectedMethod"),
						),
						writeHeader("StatusNotFound"),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("req"),
					jen.Id("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Id("nil"),
				),
				jen.Id("require").Dot("NotNil").Call(
					jen.Id("t"),
					jen.Id("req"),
				),
				requireNoError(jen.Id("err"), nil),
				jen.Line(),
				assertEqual(
					jen.Id("t"),
					jen.Id("ErrNotFound"),
					jen.Id("c").Dot("executeRequest").Call(
						jen.Id("ctx"),
						jen.Id("req"),
						jen.Op("&").Id("argleBargle").Values(),
					),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_makeUnauthedDataRequest").Params(jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				createCtx(),
				expectMethod("expectedMethod", "MethodPost"),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Id("expectedMethod"),
						),
						jen.Id("require").Dot("NoError").Call(
							jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(
								jen.Id("res"),
							).Dot("Encode").Call(
								jen.Op("&").Id("argleBargle").Values(
									jen.Id("Name").Op(":").Lit("name"),
								),
							),
						),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Line(),
				jen.List(
					jen.Id("in"),
					jen.Id("out"),
				).Op(":=").List(
					jen.Op("&").Id("argleBargle").Values(),
					jen.Op("&").Id("argleBargle").Values(),
				),
				jen.Line(),
				jen.List(
					jen.Id("body"),
					jen.Id("err"),
				).Op(":=").Id("createBodyFromStruct").Call(
					jen.Id("in"),
				),
				requireNoError(jen.Id("err"), nil),
				requireNotNil(jen.Id("body"), nil),
				jen.Line(),
				jen.List(
					jen.Id("req"),
					jen.Id("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Id("body"),
				),
				requireNoError(jen.Id("err"), nil),
				requireNotNil(jen.Id("req"), nil),
				jen.Line(),
				jen.Id("err").Op("=").Id("c").Dot("executeUnathenticatedDataRequest").Call(
					jen.Id("ctx"),
					jen.Id("req"),
					jen.Id("out"),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
				),
			),
			jen.Line(),
			buildSubTest(
				"with 404",
				createCtx(),
				expectMethod("expectedMethod", "MethodPost"),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Id("expectedMethod"),
						),
						writeHeader("StatusNotFound"),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Line(),
				jen.List(
					jen.Id("in"),
					jen.Id("out"),
				).Op(":=").List(
					jen.Op("&").Id("argleBargle").Values(),
					jen.Op("&").Id("argleBargle").Values(),
				),
				jen.Line(),
				jen.List(
					jen.Id("body"),
					jen.Id("err"),
				).Op(":=").Id("createBodyFromStruct").Call(
					jen.Id("in"),
				),
				requireNoError(jen.Id("err"), nil),
				requireNotNil(jen.Id("body"), nil),
				jen.Line(),
				jen.List(
					jen.Id("req"),
					jen.Id("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Id("body"),
				),
				requireNoError(jen.Id("err"), nil),
				requireNotNil(jen.Id("req"), nil),
				jen.Line(),
				jen.Id("err").Op("=").Id("c").Dot("executeUnathenticatedDataRequest").Call(
					jen.Id("ctx"),
					jen.Id("req"),
					jen.Id("out"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
				),
				assertEqual(jen.Id("ErrNotFound"), jen.Id("err"), nil),
			),
			jen.Line(),
			buildSubTest(
				"with timeout",
				createCtx(),
				expectMethod("expectedMethod", "MethodPost"),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Id("expectedMethod"),
						),
						jen.Qual("time", "Sleep").Call(
							jen.Lit(10).Op("*").Qual("time", "Hour"),
						),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Line(),
				jen.List(
					jen.Id("in"),
					jen.Id("out"),
				).Op(":=").List(
					jen.Op("&").Id("argleBargle").Values(),
					jen.Op("&").Id("argleBargle").Values(),
				),
				jen.Line(),
				jen.List(
					jen.Id("body"),
					jen.Id("err"),
				).Op(":=").Id("createBodyFromStruct").Call(
					jen.Id("in"),
				),
				requireNoError(jen.Id("err"), nil),
				requireNotNil(jen.Id("body"), nil),
				jen.Line(),
				jen.List(
					jen.Id("req"),
					jen.Id("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Id("body"),
				),
				requireNoError(jen.Id("err"), nil),
				requireNotNil(jen.Id("req"), nil),
				jen.Line(),
				jen.Id("c").Dot("plainClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
				assertError(
					jen.Id("t"),
					jen.Id("c").Dot("executeUnathenticatedDataRequest").Call(
						jen.Id("ctx"),
						jen.Id("req"),
						jen.Id("out"),
					),
				),
			),
			jen.Line(),
			buildSubTest(
				"with nil as output",
				createCtx(),
				expectMethod("expectedMethod", "MethodPost"),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Line(),
				jen.Id("in").Op(":=").Op("&").Id("argleBargle").Values(),
				jen.Line(),
				jen.List(
					jen.Id("body"),
					jen.Id("err"),
				).Op(":=").Id("createBodyFromStruct").Call(
					jen.Id("in"),
				),
				requireNoError(jen.Id("err"), nil),
				requireNotNil(jen.Id("body"), nil),
				jen.Line(),
				jen.List(
					jen.Id("req"),
					jen.Id("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Id("body"),
				),
				requireNoError(jen.Id("err"), nil),
				requireNotNil(jen.Id("req"), nil),
				jen.Line(),
				jen.Id("err").Op("=").Id("c").Dot("executeUnathenticatedDataRequest").Call(
					jen.Id("ctx"),
					jen.Id("req"),
					jen.Id("testingType").Values(),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestV1Client_retrieve").Params(jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				createCtx(),
				expectMethod("expectedMethod", "MethodPost"),
				jen.Id("ts").Op(":=").Id("httptest").Dot("NewTLSServer").
					Call(jen.Qual("net/http", "HandlerFunc").Call(
						jen.Func().Params(
							jen.Id("res").Qual("net/http", "ResponseWriter"),
							jen.Id("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.Id("t"),
								jen.Id("req").Dot("Method"),
								jen.Id("expectedMethod"),
							),
							jen.Id("require").Dot("NoError").Call(
								jen.Id("t"),
								jen.Qual("encoding/json", "NewEncoder").Call(
									jen.Id("res"),
								).Dot("Encode").Call(
									jen.Op("&").Id("argleBargle").Values(jen.Dict{
										jen.Id("Name"): jen.Lit("name"),
									}),
								),
							),
						),
					),
					),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Line(),
				jen.List(
					jen.Id("req"),
					jen.Id("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Id("nil"),
				),
				requireNotNil(jen.Id("req"), nil),
				requireNoError(jen.Id("err"), nil),
				jen.Line(),
				jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(
					jen.Id("ctx"),
					jen.Id("req"),
					jen.Op("&").Id("argleBargle").Values(),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
				),
			),
			jen.Line(),
			buildSubTest(
				"with nil passed in",
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Line(),
				jen.List(
					jen.Id("req"),
					jen.Id("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Id("ts").Dot("URL"),
					jen.Id("nil"),
				),
				requireNotNil(jen.Id("req"), nil),
				requireNoError(jen.Id("err"), nil),
				jen.Line(),
				jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(
					jen.Id("ctx"),
					jen.Id("req"),
					jen.Id("nil"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
				),
			),
			jen.Line(),
			buildSubTest(
				"with timeout",
				createCtx(),
				expectMethod("expectedMethod", "MethodPost"),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Id("expectedMethod"),
						),
						jen.Qual("time", "Sleep").Call(
							jen.Lit(10).Op("*").Qual("time", "Hour"),
						),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Line(),
				jen.List(
					jen.Id("req"),
					jen.Id("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Id("nil"),
				),
				requireNotNil(jen.Id("req"), nil),
				requireNoError(jen.Id("err"), nil),
				jen.Line(),
				jen.Id("c").Dot("authedClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
				jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(
					jen.Id("ctx"),
					jen.Id("req"),
					jen.Op("&").Id("argleBargle").Values(),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
				),
			),
			jen.Line(),
			buildSubTest(
				"with 404",
				createCtx(),
				expectMethod("expectedMethod", "MethodPost"),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Id("expectedMethod"),
						),
						writeHeader("StatusNotFound"),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Line(),
				jen.List(
					jen.Id("req"),
					jen.Id("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Id("expectedMethod"),
					jen.Id("ts").Dot("URL"),
					jen.Id("nil"),
				),
				requireNotNil(jen.Id("req"), nil),
				requireNoError(jen.Id("err"), nil),
				jen.Line(),
				assertEqual(
					jen.Id("t"),
					jen.Id("ErrNotFound"),
					jen.Id("c").Dot("retrieve").Call(
						jen.Id("ctx"),
						jen.Id("req"),
						jen.Op("&").Id("argleBargle").Values(),
					),
				),
			),
		),
	)

	return ret
}
