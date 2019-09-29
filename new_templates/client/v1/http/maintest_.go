package client

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func mainTestDotGo() *jen.File {
	ret := jen.NewFile("client")
	ret.Add(jen.Null())

	addImports(ret)

	// vars
	ret.Add(
		jen.Const().Defs(
			jen.ID("exampleURI").Op("=").Lit("https://todo.verygoodsoftwarenotvirus.ru"),
		),
	)

	// types
	ret.Add(
		jen.Type().Defs(
			jen.ID("argleBargle").Struct(
				jen.ID("Name").ID("string"),
			),
			jen.Line(),
			jen.ID("valuer").Map(jen.ID("string")).Index().ID("string"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(
			jen.ID("v").ID("valuer"),
		).ID("ToValues").Params().Params(
			jen.Qual("net/url", "Values"),
		).Block(
			jen.Return().Qual("net/url", "Values").Call(
				jen.ID("v"),
			),
		),
		jen.Line(),
	)

	// funcs
	ret.Add(
		jen.Comment("begin helper funcs"),
		jen.Line(),
		jen.Line(),
		jen.Func().ID("mustParseURL").Params(
			jen.ID("uri").ID("string"),
		).Params(
			jen.Op("*").Qual("net/url", "URL"),
		).Block(
			jen.List(
				jen.ID("u"),
				jen.ID("err"),
			).Op(":=").Qual("net/url", "Parse").Call(
				jen.ID("uri"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("panic").Call(
					jen.ID("err"),
				),
			),
			jen.Return().ID("u"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildTestClient").Params(
			jen.ID(t).Op("*").Qual("testing", T),
			jen.ID("ts").Op("*").Qual("net/http/httptest", "Server"),
		).Params(
			jen.Op("*").ID(v1),
		).Block(
			jen.ID(t).Dot("Helper").Call(),
			jen.Line(),
			jen.ID("l").Op(":=").Qual(noopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.ID("u").Op(":=").ID("mustParseURL").Call(
				jen.ID("ts").Dot("URL"),
			),
			jen.Line(),
			jen.Return().Op("&").ID(v1).ValuesLn(
				jen.ID("URL").Op(":").ID("u"),
				jen.ID("plainClient").Op(":").ID("ts").Dot("Client").Call(),
				jen.ID("logger").Op(":").ID("l"),
				jen.ID("Debug").Op(":").ID("true"),
				jen.ID("authedClient").Op(":").ID("ts").Dot("Client").Call(),
			),
		),
		jen.Line(),
		jen.Line(),
		jen.Comment("end helper funcs"),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_AuthenticatedClient").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTestWithoutContext(
				"obligatory",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("c").Dot("AuthenticatedClient").Call(),
				jen.Line(),
				assertEqual(
					jen.ID("ts").Dot("Client").Call(),
					jen.ID("actual"),
					jen.Lit("AuthenticatedClient should return the assigned authedClient"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_PlainClient").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTestWithoutContext(
				"obligatory",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("c").Dot("PlainClient").Call(),
				jen.Line(),
				assertEqual(
					jen.ID("ts").Dot("Client").Call(),
					jen.ID("actual"),
					jen.Lit("PlainClient should return the assigned plainClient"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_TokenSource").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTestWithoutContext(
				"obligatory",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.List(
					jen.ID("c"),
					jen.ID("err"),
				).Op(":=").ID("NewClient").CallLn(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.ID("mustParseURL").Call(
						jen.ID("exampleURI"),
					),
					jen.Qual(noopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("ts").Dot("Client").Call(),
					jen.Index().ID("string").Values(jen.Lit("*")),
					jen.ID("false"),
				),
				requireNoError(jen.ID("err"), nil),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("c").Dot("TokenSource").Call(),
				jen.Line(),
				assertNotNil(
					jen.ID("actual"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("NewClient").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTestWithoutContext(
				"happy path",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.List(
					jen.ID("c"),
					jen.ID("err"),
				).Op(":=").ID("NewClient").CallLn(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.ID("mustParseURL").Call(
						jen.ID("exampleURI"),
					),
					jen.Qual(noopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("ts").Dot("Client").Call(),
					jen.Index().ID("string").Values(jen.Lit("*")),
					jen.ID("false"),
				),
				jen.Line(),
				requireNotNil(jen.ID("c"), nil),
				requireNoError(jen.ID("err"), nil),
			),
			jen.Line(),
			buildSubTestWithoutContext(
				"with client but invalid timeout",
				jen.List(
					jen.ID("c"),
					jen.ID("err"),
				).Op(":=").ID("NewClient").CallLn(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.ID("mustParseURL").Call(
						jen.ID("exampleURI"),
					),
					jen.Qual(noopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Op("&").Qual("net/http", "Client").ValuesLn(
						jen.ID("Timeout").Op(":").Lit(0),
					),
					jen.Index().ID("string").Values(jen.Lit("*")),
					jen.ID("true"),
				),
				jen.Line(),
				requireNotNil(jen.ID("c"), nil),
				requireNoError(jen.ID("err"), nil),
				assertEqual(
					jen.ID("c").Dot("plainClient").Dot("Timeout"),
					jen.ID("defaultTimeout"),
					jen.Lit("NewClient should set the default timeout"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("NewSimpleClient").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTestWithoutContext(
				"obligatory",
				jen.List(
					jen.ID("c"),
					jen.ID("err"),
				).Op(":=").ID("NewSimpleClient").CallLn(
					jen.Qual("context", "Background").Call(),
					jen.ID("mustParseURL").Call(
						jen.ID("exampleURI"),
					),
					jen.ID("true"),
				),
				assertNotNil(jen.ID("c"), nil),
				assertNoError(jen.ID("err"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_executeRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"with error",
				expectMethod("expectedMethod", "MethodPost"),
				jen.Line(),
				buildTestServer(
					"ts",
					assertEqual(
						jen.ID("req").Dot("Method"),
						jen.ID("expectedMethod"),
						nil,
					),
					jen.Qual("time", "Sleep").Call(
						jen.Lit(10).Op("*").Qual("time", "Hour"),
					),
				),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("nil"),
				),
				requireNotNil(jen.ID("req"), nil),
				requireNoError(jen.ID("err"), nil),
				jen.Line(),
				jen.List(
					jen.ID("res"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("executeRawRequest").Call(
					jen.ID("ctx"),
					jen.Op("&").Qual("net/http", "Client").Values(
						jen.ID("Timeout").Op(":").Qual("time", "Second"),
					),
					jen.ID("req"),
				),
				assertNil(
					jen.ID("res"),
					nil,
				),
				assertError(
					jen.ID("err"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("BuildURL").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTestWithoutContext(
				"various urls",
				parallelTest(jen.ID("t")),
				jen.Line(),
				jen.List(
					jen.ID("u"),
					jen.ID("_"),
				).Op(":=").Qual("net/url", "Parse").Call(
					jen.ID("exampleURI"),
				),
				jen.List(
					jen.ID("c"),
					jen.ID("err"),
				).Op(":=").ID("NewClient").CallLn(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.ID("u"),
					jen.Qual(noopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("nil"),
					jen.Index().ID("string").Values(jen.Lit("*")),
					jen.ID("false"),
				),
				requireNoError(jen.ID("err"), nil),
				jen.Line(),
				jen.ID("testCases").Op(":=").Index().Struct(
					jen.ID("expectation").ID("string"),
					jen.ID("inputParts").Index().ID("string"),
					jen.ID("inputQuery").ID("valuer"),
				).ValuesLn(
					jen.ValuesLn(
						jen.ID("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/things"),
						jen.ID("inputParts").Op(":").Index().ID("string").Values(jen.Lit("things")),
					),
					jen.ValuesLn(
						jen.ID("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/stuff?key=value"),
						jen.ID("inputQuery").Op(":").Map(jen.ID("string")).Index().ID("string").Values(
							jen.Lit("key").Op(":").Values(jen.Lit("value")),
						),
						jen.ID("inputParts").Op(":").Index().ID("string").Values(jen.Lit("stuff")),
					),
					jen.ValuesLn(
						jen.ID("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/things/and/stuff?key=value1&key=value2&yek=eulav"),
						jen.ID("inputQuery").Op(":").Map(jen.ID("string")).Index().ID("string").ValuesLn(
							jen.Lit("key").Op(":").Values(
								jen.Lit("value1"),
								jen.Lit("value2"),
							),
							jen.Lit("yek").Op(":").Values(jen.Lit("eulav")),
						),
						jen.ID("inputParts").Op(":").Index().ID("string").Values(jen.Lit("things"),
							jen.Lit("and"),
							jen.Lit("stuff"),
						),
					),
				),
				jen.Line(),
				jen.For(jen.List(
					jen.ID("_"),
					jen.ID("tc"),
				).Op(":=").Range().ID("testCases"),
				).Block(
					jen.ID("actual").Op(":=").ID("c").Dot("BuildURL").Call(
						jen.ID("tc").Dot("inputQuery").Dot("ToValues").Call(),
						jen.ID("tc").Dot("inputParts").Op("..."),
					),
					assertEqual(
						jen.ID("tc").Dot("expectation"),
						jen.ID("actual"),
						nil,
					),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildWebsocketURL").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTestWithoutContext(
				"happy path",
				jen.List(
					jen.ID("u"),
					jen.ID("_"),
				).Op(":=").Qual("net/url", "Parse").Call(
					jen.ID("exampleURI"),
				),
				jen.List(
					jen.ID("c"),
					jen.ID("err"),
				).Op(":=").ID("NewClient").CallLn(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.ID("u"),
					jen.Qual(noopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("nil"),
					jen.Index().ID("string").Values(jen.Lit("*")),
					jen.ID("false"),
				),
				requireNoError(jen.ID("err"), nil),
				jen.Line(),
				jen.ID("expected").Op(":=").Lit("ws://todo.verygoodsoftwarenotvirus.ru/api/v1/things/and/stuff"),
				jen.ID("actual").Op(":=").ID("c").Dot("BuildWebsocketURL").Call(
					jen.Lit("things"),
					jen.Lit("and"),
					jen.Lit("stuff"),
				),
				jen.Line(),
				assertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildHealthCheckRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTestWithoutContext(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildHealthCheckRequest").Call(),
				jen.Line(),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.ID("actual").Dot("Method"),
					jen.ID("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.ID("expectedMethod"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_IsUp").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTestWithoutContext(
				"happy path",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodGet"),
								nil,
							),
							writeHeader("StatusOK"),
						),
					),
				),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.ID("actual").Op(":=").ID("c").Dot("IsUp").Call(),
				assertTrue(
					jen.ID("actual"),
					nil,
				),
			),
			jen.Line(),
			buildSubTestWithoutContext(
				"with bad status code",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodGet"),
								nil,
							),
							writeHeader("StatusInternalServerError"),
						),
					),
				),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.ID("actual").Op(":=").ID("c").Dot("IsUp").Call(),
				assertFalse(
					jen.ID("actual"),
					nil,
				),
			),
			jen.Line(),
			buildSubTestWithoutContext(
				"with timeout",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodGet"),
								nil,
							),
							jen.Qual("time", "Sleep").Call(
								jen.Lit(10).Op("*").Qual("time", "Hour"),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.ID("c").Dot("plainClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
				jen.ID("actual").Op(":=").ID("c").Dot("IsUp").Call(),
				assertFalse(
					jen.ID("actual"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_buildDataRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTestWithoutContext(
				"happy path",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.Line(),
				expectMethod("expectedMethod", "MethodPost"),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("buildDataRequest").Call(
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Op("&").ID("testingType").Values(
						jen.ID("Name").Op(":").Lit("name"),
					),
				),
				jen.Line(),
				requireNotNil(jen.ID("req"), nil),
				assertNoError(
					jen.ID("err"),
					nil,
				),
				assertEqual(
					jen.ID("expectedMethod"),
					jen.ID("req").Dot("Method"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_makeRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							requireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(
									jen.ID("res"),
								).Dot("Encode").Call(
									jen.Op("&").ID("argleBargle").Values(
										jen.ID("Name").Op(":").Lit("name"),
									),
								),
								nil,
							),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("nil"),
				),
				requireNotNil(jen.ID("req"), nil),
				requireNoError(jen.ID("err"), nil),
				jen.Line(),
				jen.ID("err").Op("=").ID("c").Dot("executeRequest").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.Op("&").ID("argleBargle").Values(),
				),
				assertNoError(
					jen.ID("err"),
					nil,
				),
			),
			jen.Line(),
			buildSubTest(
				"with 404",
				expectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							writeHeader("StatusNotFound"),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("nil"),
				),
				requireNotNil(
					jen.ID("req"),
					nil,
				),
				requireNoError(jen.ID("err"), nil),
				jen.Line(),
				assertEqual(
					jen.ID("ErrNotFound"),
					jen.ID("c").Dot("executeRequest").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.Op("&").ID("argleBargle").Values(),
					),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_makeUnauthedDataRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							requireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(
									jen.ID("res"),
								).Dot("Encode").Call(
									jen.Op("&").ID("argleBargle").Values(
										jen.ID("Name").Op(":").Lit("name"),
									),
								),
								nil,
							),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("in"),
					jen.ID("out"),
				).Op(":=").List(
					jen.Op("&").ID("argleBargle").Values(),
					jen.Op("&").ID("argleBargle").Values(),
				),
				jen.Line(),
				jen.List(
					jen.ID("body"),
					jen.ID("err"),
				).Op(":=").ID("createBodyFromStruct").Call(
					jen.ID("in"),
				),
				requireNoError(jen.ID("err"), nil),
				requireNotNil(jen.ID("body"), nil),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				requireNoError(jen.ID("err"), nil),
				requireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.ID("err").Op("=").ID("c").Dot("executeUnathenticatedDataRequest").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.ID("out"),
				),
				assertNoError(
					jen.ID("err"),
					nil,
				),
			),
			jen.Line(),
			buildSubTest(
				"with 404",
				expectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							writeHeader("StatusNotFound"),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("in"),
					jen.ID("out"),
				).Op(":=").List(
					jen.Op("&").ID("argleBargle").Values(),
					jen.Op("&").ID("argleBargle").Values(),
				),
				jen.Line(),
				jen.List(
					jen.ID("body"),
					jen.ID("err"),
				).Op(":=").ID("createBodyFromStruct").Call(
					jen.ID("in"),
				),
				requireNoError(jen.ID("err"), nil),
				requireNotNil(jen.ID("body"), nil),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				requireNoError(jen.ID("err"), nil),
				requireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.ID("err").Op("=").ID("c").Dot("executeUnathenticatedDataRequest").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.ID("out"),
				),
				assertError(
					jen.ID("err"),
					nil,
				),
				assertEqual(jen.ID("ErrNotFound"),
					jen.ID("err"), nil),
			),
			jen.Line(),
			buildSubTest(
				"with timeout",
				expectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							jen.Qual("time", "Sleep").Call(
								jen.Lit(10).Op("*").Qual("time", "Hour"),
							),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("in"),
					jen.ID("out"),
				).Op(":=").List(
					jen.Op("&").ID("argleBargle").Values(),
					jen.Op("&").ID("argleBargle").Values(),
				),
				jen.Line(),
				jen.List(
					jen.ID("body"),
					jen.ID("err"),
				).Op(":=").ID("createBodyFromStruct").Call(
					jen.ID("in"),
				),
				requireNoError(jen.ID("err"), nil),
				requireNotNil(jen.ID("body"), nil),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				requireNoError(jen.ID("err"), nil),
				requireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.ID("c").Dot("plainClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
				assertError(
					jen.ID("c").Dot("executeUnathenticatedDataRequest").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("out"),
					),
					nil,
				),
			),
			jen.Line(),
			buildSubTest(
				"with nil as output",
				expectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.ID("in").Op(":=").Op("&").ID("argleBargle").Values(),
				jen.Line(),
				jen.List(
					jen.ID("body"),
					jen.ID("err"),
				).Op(":=").ID("createBodyFromStruct").Call(
					jen.ID("in"),
				),
				requireNoError(jen.ID("err"), nil),
				requireNotNil(jen.ID("body"), nil),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				requireNoError(jen.ID("err"), nil),
				requireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.ID("err").Op("=").ID("c").Dot("executeUnathenticatedDataRequest").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.ID("testingType").Values(),
				),
				assertError(
					jen.ID("err"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_retrieve").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").ID("httptest").Dot("NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							requireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(
									jen.ID("res"),
								).Dot("Encode").Call(
									jen.Op("&").ID("argleBargle").Values(
										jen.ID("Name").Op(":").Lit("name"),
									),
								),
								nil,
							),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("nil"),
				),
				requireNotNil(jen.ID("req"), nil),
				requireNoError(jen.ID("err"), nil),
				jen.Line(),
				jen.ID("err").Op("=").ID("c").Dot("retrieve").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.Op("&").ID("argleBargle").Values(),
				),
				assertNoError(
					jen.ID("err"),
					nil,
				),
			),
			jen.Line(),
			buildSubTest(
				"with nil passed in",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.ID("ts").Dot("URL"),
					jen.ID("nil"),
				),
				requireNotNil(jen.ID("req"), nil),
				requireNoError(jen.ID("err"), nil),
				jen.Line(),
				jen.ID("err").Op("=").ID("c").Dot("retrieve").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.ID("nil"),
				),
				assertError(
					jen.ID("err"),
					nil,
				),
			),
			jen.Line(),
			buildSubTest(
				"with timeout",
				expectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							jen.Qual("time", "Sleep").Call(
								jen.Lit(10).Op("*").Qual("time", "Hour"),
							),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("nil"),
				),
				requireNotNil(jen.ID("req"), nil),
				requireNoError(jen.ID("err"), nil),
				jen.Line(),
				jen.ID("c").Dot("authedClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
				jen.ID("err").Op("=").ID("c").Dot("retrieve").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.Op("&").ID("argleBargle").Values(),
				),
				assertError(
					jen.ID("err"),
					nil,
				),
			),
			jen.Line(),
			buildSubTest(
				"with 404",
				expectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							writeHeader("StatusNotFound"),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("nil"),
				),
				requireNotNil(jen.ID("req"), nil),
				requireNoError(jen.ID("err"), nil),
				jen.Line(),
				assertEqual(
					jen.ID("ErrNotFound"),
					jen.ID("c").Dot("retrieve").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.Op("&").ID("argleBargle").Values(),
					),
					nil,
				),
			),
		),
	)

	return ret
}
