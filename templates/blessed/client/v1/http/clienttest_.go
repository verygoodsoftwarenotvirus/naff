package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func mainTestDotGo() *jen.File {
	ret := jen.NewFile("client")
	utils.AddImports(ret)

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
			jen.ID("t").Op("*").Qual("testing", "T"),
			jen.ID("ts").Op("*").Qual("net/http/httptest", "Server"),
		).Params(
			jen.Op("*").ID(v1),
		).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("l").Op(":=").Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.ID("u").Op(":=").ID("mustParseURL").Call(
				jen.ID("ts").Dot("URL"),
			),
			jen.ID("c").Op(":=").ID("ts").Dot("Client").Call(),
			jen.Line(),
			jen.Return().Op("&").ID(v1).Valuesln(
				jen.ID("URL").Op(":").ID("u"),
				jen.ID("plainClient").Op(":").ID("c"),
				jen.ID("logger").Op(":").ID("l"),
				jen.ID("Debug").Op(":").ID("true"),
				jen.ID("authedClient").Op(":").ID("c"),
			),
		),
		jen.Line(),
		jen.Line(),
		jen.Comment("end helper funcs"),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_AuthenticatedClient").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("c").Dot("AuthenticatedClient").Call(),
				jen.Line(),
				utils.AssertEqual(
					jen.ID("ts").Dot("Client").Call(),
					jen.ID("actual"),
					jen.Lit("AuthenticatedClient should return the assigned authedClient"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_PlainClient").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("c").Dot("PlainClient").Call(),
				jen.Line(),
				utils.AssertEqual(
					jen.ID("ts").Dot("Client").Call(),
					jen.ID("actual"),
					jen.Lit("PlainClient should return the assigned plainClient"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_TokenSource").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.List(
					jen.ID("c"),
					jen.ID("err"),
				).Op(":=").ID("NewClient").Callln(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.ID("mustParseURL").Call(
						jen.ID("exampleURI"),
					),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("ts").Dot("Client").Call(),
					jen.Index().ID("string").Values(jen.Lit("*")),
					jen.ID("false"),
				),
				utils.RequireNoError(jen.ID("err"), nil),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("c").Dot("TokenSource").Call(),
				jen.Line(),
				utils.AssertNotNil(jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("NewClient").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.List(
					jen.ID("c"),
					jen.ID("err"),
				).Op(":=").ID("NewClient").Callln(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.ID("mustParseURL").Call(
						jen.ID("exampleURI"),
					),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("ts").Dot("Client").Call(),
					jen.Index().ID("string").Values(jen.Lit("*")),
					jen.ID("false"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("c"), nil),
				utils.RequireNoError(jen.ID("err"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with client but invalid timeout",
				jen.List(
					jen.ID("c"),
					jen.ID("err"),
				).Op(":=").ID("NewClient").Callln(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.ID("mustParseURL").Call(
						jen.ID("exampleURI"),
					),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Op("&").Qual("net/http", "Client").Valuesln(
						jen.ID("Timeout").Op(":").Lit(0),
					),
					jen.Index().ID("string").Values(jen.Lit("*")),
					jen.ID("true"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("c"), nil),
				utils.RequireNoError(jen.ID("err"), nil),
				utils.AssertEqual(
					jen.ID("c").Dot("plainClient").Dot("Timeout"),
					jen.ID("defaultTimeout"),
					jen.Lit("NewClient should set the default timeout"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("NewSimpleClient").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.List(
					jen.ID("c"),
					jen.ID("err"),
				).Op(":=").ID("NewSimpleClient").Callln(
					jen.Qual("context", "Background").Call(),
					jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
					jen.ID("true"),
				),
				utils.AssertNotNil(jen.ID("c"), nil),
				utils.AssertNoError(jen.ID("err"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_executeRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"with error",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
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
					jen.ID("t"),
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
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.ID("err"), nil),
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
				utils.AssertNil(
					jen.ID("res"),
					nil,
				),
				utils.AssertError(
					jen.ID("err"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("BuildURL").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"various urls",
				utils.ParallelTest(jen.ID("t")),
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
				).Op(":=").ID("NewClient").Callln(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.ID("u"),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("nil"),
					jen.Index().ID("string").Values(jen.Lit("*")),
					jen.ID("false"),
				),
				utils.RequireNoError(jen.ID("err"), nil),
				jen.Line(),
				jen.ID("testCases").Op(":=").Index().Struct(
					jen.ID("expectation").ID("string"),
					jen.ID("inputParts").Index().ID("string"),
					jen.ID("inputQuery").ID("valuer"),
				).Valuesln(
					jen.Valuesln(
						jen.ID("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/things"),
						jen.ID("inputParts").Op(":").Index().ID("string").Values(jen.Lit("things")),
					),
					jen.Valuesln(
						jen.ID("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/stuff?key=value"),
						jen.ID("inputQuery").Op(":").Map(jen.ID("string")).Index().ID("string").Values(
							jen.Lit("key").Op(":").Values(jen.Lit("value")),
						),
						jen.ID("inputParts").Op(":").Index().ID("string").Values(jen.Lit("stuff")),
					),
					jen.Valuesln(
						jen.ID("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/things/and/stuff?key=value1&key=value2&yek=eulav"),
						jen.ID("inputQuery").Op(":").Map(jen.ID("string")).Index().ID("string").Valuesln(
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
					utils.AssertEqual(
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
		utils.OuterTestFunc("V1Client_BuildWebsocketURL").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
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
				).Op(":=").ID("NewClient").Callln(
					jen.Qual("context", "Background").Call(),
					jen.Lit(""),
					jen.Lit(""),
					jen.ID("u"),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("nil"),
					jen.Index().ID("string").Values(jen.Lit("*")),
					jen.ID("false"),
				),
				utils.RequireNoError(jen.ID("err"), nil),
				jen.Line(),
				jen.ID("expected").Op(":=").Lit("ws://todo.verygoodsoftwarenotvirus.ru/api/v1/things/and/stuff"),
				jen.ID("actual").Op(":=").ID("c").Dot("BuildWebsocketURL").Call(
					jen.Lit("things"),
					jen.Lit("and"),
					jen.Lit("stuff"),
				),
				jen.Line(),
				utils.AssertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_BuildHealthCheckRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodGet"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildHealthCheckRequest").Call(),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(
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
		utils.OuterTestFunc("V1Client_IsUp").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodGet"),
								nil,
							),
							utils.WriteHeader("StatusOK"),
						),
					),
				),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.ID("actual").Op(":=").ID("c").Dot("IsUp").Call(),
				utils.AssertTrue(
					jen.ID("actual"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with bad status code",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodGet"),
								nil,
							),
							utils.WriteHeader("StatusInternalServerError"),
						),
					),
				),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.ID("actual").Op(":=").ID("c").Dot("IsUp").Call(),
				utils.AssertFalse(
					jen.ID("actual"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with timeout",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
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
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.ID("c").Dot("plainClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
				jen.ID("actual").Op(":=").ID("c").Dot("IsUp").Call(),
				utils.AssertFalse(
					jen.ID("actual"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_buildDataRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				utils.ExpectMethod("expectedMethod", "MethodPost"),
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
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.AssertNoError(
					jen.ID("err"),
					nil,
				),
				utils.AssertEqual(
					jen.ID("expectedMethod"),
					jen.ID("req").Dot("Method"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_makeRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.RequireNoError(
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
					jen.ID("t"),
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
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.ID("err"), nil),
				jen.Line(),
				jen.ID("err").Op("=").ID("c").Dot("executeRequest").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.Op("&").ID("argleBargle").Values(),
				),
				utils.AssertNoError(
					jen.ID("err"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with 404",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.WriteHeader("StatusNotFound"),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
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
				utils.RequireNotNil(
					jen.ID("req"),
					nil,
				),
				utils.RequireNoError(jen.ID("err"), nil),
				jen.Line(),
				utils.AssertEqual(
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
		utils.OuterTestFunc("V1Client_makeUnauthedDataRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.RequireNoError(
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
					jen.ID("t"),
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
				utils.RequireNoError(jen.ID("err"), nil),
				utils.RequireNotNil(jen.ID("body"), nil),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				utils.RequireNoError(jen.ID("err"), nil),
				utils.RequireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.ID("err").Op("=").ID("c").Dot("executeUnathenticatedDataRequest").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.ID("out"),
				),
				utils.AssertNoError(
					jen.ID("err"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with 404",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.WriteHeader("StatusNotFound"),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
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
				utils.RequireNoError(jen.ID("err"), nil),
				utils.RequireNotNil(jen.ID("body"), nil),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				utils.RequireNoError(jen.ID("err"), nil),
				utils.RequireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.ID("err").Op("=").ID("c").Dot("executeUnathenticatedDataRequest").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.ID("out"),
				),
				utils.AssertError(
					jen.ID("err"),
					nil,
				),
				utils.AssertEqual(jen.ID("ErrNotFound"),
					jen.ID("err"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with timeout",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
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
					jen.ID("t"),
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
				utils.RequireNoError(jen.ID("err"), nil),
				utils.RequireNotNil(jen.ID("body"), nil),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				utils.RequireNoError(jen.ID("err"), nil),
				utils.RequireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.ID("c").Dot("plainClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
				utils.AssertError(
					jen.ID("c").Dot("executeUnathenticatedDataRequest").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("out"),
					),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with nil as output",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
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
				utils.RequireNoError(jen.ID("err"), nil),
				utils.RequireNotNil(jen.ID("body"), nil),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").Qual("net/http", "NewRequest").Call(
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				utils.RequireNoError(jen.ID("err"), nil),
				utils.RequireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.ID("err").Op("=").ID("c").Dot("executeUnathenticatedDataRequest").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.ID("testingType").Values(),
				),
				utils.AssertError(
					jen.ID("err"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_retrieve").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").ID("httptest").Dot("NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.RequireNoError(
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
					jen.ID("t"),
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
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.ID("err"), nil),
				jen.Line(),
				jen.ID("err").Op("=").ID("c").Dot("retrieve").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.Op("&").ID("argleBargle").Values(),
				),
				utils.AssertNoError(
					jen.ID("err"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with nil passed in",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
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
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.ID("err"), nil),
				jen.Line(),
				jen.ID("err").Op("=").ID("c").Dot("retrieve").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.ID("nil"),
				),
				utils.AssertError(
					jen.ID("err"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with timeout",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
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
					jen.ID("t"),
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
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.ID("err"), nil),
				jen.Line(),
				jen.ID("c").Dot("authedClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
				jen.ID("err").Op("=").ID("c").Dot("retrieve").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.Op("&").ID("argleBargle").Values(),
				),
				utils.AssertError(
					jen.ID("err"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with 404",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.WriteHeader("StatusNotFound"),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
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
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.ID("err"), nil),
				jen.Line(),
				utils.AssertEqual(
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
