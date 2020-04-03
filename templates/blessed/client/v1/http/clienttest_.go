package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)
	utils.AddImports(proj, ret)

	// vars
	ret.Add(
		jen.Const().Defs(
			jen.ID("exampleURI").Equals().Lit("https://todo.verygoodsoftwarenotvirus.ru"),
			jen.ID("asciiControlChar").Equals().ID("string").Call(jen.ID("byte").Call(jen.Lit(0x7f))),
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
		jen.Line(),
		jen.Comment("begin helper funcs"),
		jen.Line(),
		jen.Line(),
	)

	ret.Add(buildMustParseURL()...)
	ret.Add(buildBuildTestClient()...)
	ret.Add(buildBuildTestClientWithInvalidURL()...)

	ret.Add(
		jen.Line(),
		jen.Line(),
		jen.Comment("end helper funcs"),
		jen.Line(),
	)

	ret.Add(buildTestV1Client_AuthenticatedClient()...)
	ret.Add(buildTestV1Client_PlainClient()...)
	ret.Add(buildTestV1Client_TokenSource()...)
	ret.Add(buildTestNewClient()...)
	ret.Add(buildTestNewSimpleClient()...)
	ret.Add(buildTestV1Client_CloseRequestBody()...)
	ret.Add(buildTestBuildURL()...)
	ret.Add(buildTestBuildVersionlessURL()...)
	ret.Add(buildTestV1Client_BuildWebsocketURL()...)
	ret.Add(buildTestV1Client_BuildHealthCheckRequest()...)
	ret.Add(buildTestV1Client_IsUp()...)
	ret.Add(buildTestV1Client_buildDataRequest()...)
	ret.Add(buildTestV1Client_executeRequest()...)
	ret.Add(buildTestV1Client_executeRawRequest()...)
	ret.Add(buildTestV1Client_checkExistence()...)
	ret.Add(buildTestV1Client_retrieve()...)
	ret.Add(buildTestV1Client_executeUnauthenticatedDataRequest()...)

	return ret
}

func buildMustParseURL() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("mustParseURL").Params(
			jen.ID("uri").ID("string"),
		).Params(
			jen.ParamPointer().Qual("net/url", "URL"),
		).Block(
			jen.List(
				jen.ID("u"),
				jen.Err(),
			).Assign().Qual("net/url", "Parse").Call(
				jen.ID("uri"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("panic").Call(
					jen.Err(),
				),
			),
			jen.Return().ID("u"),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildTestClient() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildTestClient").Params(
			jen.ID("t").ParamPointer().Qual("testing", "T"),
			jen.ID("ts").ParamPointer().Qual("net/http/httptest", "Server"),
		).Params(
			jen.PointerTo().ID(v1),
		).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("l").Assign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.ID("u").Assign().ID("mustParseURL").Call(
				jen.ID("ts").Dot("URL"),
			),
			jen.ID("c").Assign().ID("ts").Dot("Client").Call(),
			jen.Line(),
			jen.Return().VarPointer().ID(v1).Valuesln(
				jen.ID("URL").MapAssign().ID("u"),
				jen.ID("plainClient").MapAssign().ID("c"),
				jen.ID("logger").MapAssign().ID("l"),
				jen.ID("Debug").MapAssign().ID("true"),
				jen.ID("authedClient").MapAssign().ID("c"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildTestClientWithInvalidURL() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildTestClientWithInvalidURL").Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Params(jen.PointerTo().ID(v1)).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("l").Assign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.ID("u").Assign().ID("mustParseURL").Call(jen.Lit("https://verygoodsoftwarenotvirus.ru")),
			jen.ID("u").Dot("Scheme").Equals().Qual("fmt", "Sprintf").Call(
				jen.RawString(`%s://`),
				jen.ID("asciiControlChar"),
			),
			jen.Line(),
			jen.Return().VarPointer().ID(v1).Valuesln(
				jen.ID("URL").MapAssign().ID("u"),
				jen.ID("plainClient").MapAssign().Qual("net/http", "DefaultClient"),
				jen.ID("logger").MapAssign().ID("l"),
				jen.ID("Debug").MapAssign().ID("true"),
				jen.ID("authedClient").MapAssign().Qual("net/http", "DefaultClient"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_AuthenticatedClient() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_AuthenticatedClient").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("AuthenticatedClient").Call(),
				jen.Line(),
				utils.AssertEqual(
					jen.ID("ts").Dot("Client").Call(),
					jen.ID("actual"),
					jen.Lit("AuthenticatedClient should return the assigned authedClient"),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_PlainClient() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_PlainClient").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("PlainClient").Call(),
				jen.Line(),
				utils.AssertEqual(
					jen.ID("ts").Dot("Client").Call(),
					jen.ID("actual"),
					jen.Lit("PlainClient should return the assigned plainClient"),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_TokenSource() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_TokenSource").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"obligatory",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.List(
					jen.ID("c"),
					jen.Err(),
				).Assign().ID("NewClient").Callln(
					utils.CtxVar(),
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
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("TokenSource").Call(),
				jen.Line(),
				utils.AssertNotNil(jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestNewClient() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("NewClient").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.List(
					jen.ID("c"),
					jen.Err(),
				).Assign().ID("NewClient").Callln(
					utils.CtxVar(),
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
				utils.RequireNoError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with client but invalid timeout",
				jen.Line(),
				jen.List(
					jen.ID("c"),
					jen.Err(),
				).Assign().ID("NewClient").Callln(
					utils.CtxVar(),
					jen.Lit(""),
					jen.Lit(""),
					jen.ID("mustParseURL").Call(
						jen.ID("exampleURI"),
					),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.VarPointer().Qual("net/http", "Client").Valuesln(
						jen.ID("Timeout").MapAssign().Lit(0),
					),
					jen.Index().ID("string").Values(jen.Lit("*")),
					jen.ID("true"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("c"), nil),
				utils.RequireNoError(jen.Err(), nil),
				utils.AssertEqual(
					jen.ID("c").Dot("plainClient").Dot("Timeout"),
					jen.ID("defaultTimeout"),
					jen.Lit("NewClient should set the default timeout"),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestNewSimpleClient() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("NewSimpleClient").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"obligatory",
				jen.Line(),
				jen.List(
					jen.ID("c"),
					jen.Err(),
				).Assign().ID("NewSimpleClient").Callln(
					utils.CtxVar(),
					jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
					jen.ID("true"),
				),
				utils.AssertNotNil(jen.ID("c"), nil),
				utils.AssertNoError(jen.Err(), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_CloseRequestBody() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_CloseRequestBody").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"with error",
				jen.Line(),
				jen.ID("rc").Assign().ID("newMockReadCloser").Call(),
				jen.ID("rc").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.ID("res").Assign().VarPointer().Qual("net/http", "Response").Valuesln(
					jen.ID("Body").MapAssign().ID("rc"),
					jen.ID("StatusCode").MapAssign().Qual("net/http", "StatusOK"),
				),
				jen.Line(),
				jen.List(
					jen.ID("c"),
					jen.Err(),
				).Assign().ID("NewSimpleClient").Callln(
					utils.CtxVar(),
					jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
					jen.True(),
				),
				utils.AssertNotNil(jen.ID("c"), nil),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("c").Dot("closeResponseBody").Call(jen.ID("res")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestBuildURL() []jen.Code {
	testVar := func() *jen.Statement {
		return jen.ID("t")
	}

	lines := []jen.Code{
		utils.OuterTestFunc("BuildURL").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"various urls",
				utils.ParallelTest(testVar()),
				jen.Line(),
				jen.List(
					jen.ID("u"),
					jen.ID("_"),
				).Assign().Qual("net/url", "Parse").Call(
					jen.ID("exampleURI"),
				),
				jen.List(
					jen.ID("c"),
					jen.Err(),
				).Assign().ID("NewClient").Callln(
					utils.CtxVar(),
					jen.Lit(""),
					jen.Lit(""),
					jen.ID("u"),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Nil(),
					jen.Index().ID("string").Values(jen.Lit("*")),
					jen.ID("false"),
				),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("testCases").Assign().Index().Struct(
					jen.ID("expectation").ID("string"),
					jen.ID("inputParts").Index().ID("string"),
					jen.ID("inputQuery").ID("valuer"),
				).Valuesln(
					jen.Valuesln(
						jen.ID("expectation").MapAssign().Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/things"),
						jen.ID("inputParts").MapAssign().Index().ID("string").Values(jen.Lit("things")),
					),
					jen.Valuesln(
						jen.ID("expectation").MapAssign().Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/stuff?key=value"),
						jen.ID("inputQuery").MapAssign().Map(jen.ID("string")).Index().ID("string").Values(
							jen.Lit("key").MapAssign().Values(jen.Lit("value")),
						),
						jen.ID("inputParts").MapAssign().Index().ID("string").Values(jen.Lit("stuff")),
					),
					jen.Valuesln(
						jen.ID("expectation").MapAssign().Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/things/and/stuff?key=value1&key=value2&yek=eulav"),
						jen.ID("inputQuery").MapAssign().Map(jen.ID("string")).Index().ID("string").Valuesln(
							jen.Lit("key").MapAssign().Values(
								jen.Lit("value1"),
								jen.Lit("value2"),
							),
							jen.Lit("yek").MapAssign().Values(jen.Lit("eulav")),
						),
						jen.ID("inputParts").MapAssign().Index().ID("string").Values(jen.Lit("things"),
							jen.Lit("and"),
							jen.Lit("stuff"),
						),
					),
				),
				jen.Line(),
				jen.For(jen.List(
					jen.ID("_"),
					jen.ID("tc"),
				).Assign().Range().ID("testCases"),
				).Block(
					jen.ID("actual").Assign().ID("c").Dot("BuildURL").Call(
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
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with invalid URL parts",
				utils.ParallelTest(testVar()),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(testVar()),
				jen.Qual(utils.AssertPkg, "Empty").Call(
					testVar(),
					jen.ID("c").Dot("BuildURL").Call(
						jen.Nil(),
						jen.ID("asciiControlChar"),
					),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestBuildVersionlessURL() []jen.Code {
	testVar := func() *jen.Statement {
		return jen.ID("t")
	}

	lines := []jen.Code{
		utils.OuterTestFunc("BuildVersionlessURL").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"various urls",
				utils.ParallelTest(testVar()),
				jen.Line(),
				jen.List(
					jen.ID("u"),
					jen.ID("_"),
				).Assign().Qual("net/url", "Parse").Call(
					jen.ID("exampleURI"),
				),
				jen.List(
					jen.ID("c"),
					jen.Err(),
				).Assign().ID("NewClient").Callln(
					utils.CtxVar(),
					jen.Lit(""),
					jen.Lit(""),
					jen.ID("u"),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Nil(),
					jen.Index().ID("string").Values(jen.Lit("*")),
					jen.ID("false"),
				),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("testCases").Assign().Index().Struct(
					jen.ID("expectation").ID("string"),
					jen.ID("inputParts").Index().ID("string"),
					jen.ID("inputQuery").ID("valuer"),
				).Valuesln(
					jen.Valuesln(
						jen.ID("expectation").MapAssign().Lit("https://todo.verygoodsoftwarenotvirus.ru/things"),
						jen.ID("inputParts").MapAssign().Index().ID("string").Values(jen.Lit("things")),
					),
					jen.Valuesln(
						jen.ID("expectation").MapAssign().Lit("https://todo.verygoodsoftwarenotvirus.ru/stuff?key=value"),
						jen.ID("inputQuery").MapAssign().Map(jen.ID("string")).Index().ID("string").Values(
							jen.Lit("key").MapAssign().Values(jen.Lit("value")),
						),
						jen.ID("inputParts").MapAssign().Index().ID("string").Values(jen.Lit("stuff")),
					),
					jen.Valuesln(
						jen.ID("expectation").MapAssign().Lit("https://todo.verygoodsoftwarenotvirus.ru/things/and/stuff?key=value1&key=value2&yek=eulav"),
						jen.ID("inputQuery").MapAssign().Map(jen.ID("string")).Index().ID("string").Valuesln(
							jen.Lit("key").MapAssign().Values(
								jen.Lit("value1"),
								jen.Lit("value2"),
							),
							jen.Lit("yek").MapAssign().Values(jen.Lit("eulav")),
						),
						jen.ID("inputParts").MapAssign().Index().ID("string").Values(jen.Lit("things"),
							jen.Lit("and"),
							jen.Lit("stuff"),
						),
					),
				),
				jen.Line(),
				jen.For(jen.List(
					jen.ID("_"),
					jen.ID("tc"),
				).Assign().Range().ID("testCases"),
				).Block(
					jen.ID("actual").Assign().ID("c").Dot("buildVersionlessURL").Call(
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
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with invalid URL parts",
				utils.ParallelTest(testVar()),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(testVar()),
				jen.Qual(utils.AssertPkg, "Empty").Call(
					testVar(),
					jen.ID("c").Dot("buildVersionlessURL").Call(
						jen.Nil(),
						jen.ID("asciiControlChar"),
					),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_BuildWebsocketURL() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_BuildWebsocketURL").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.List(
					jen.ID("u"),
					jen.ID("_"),
				).Assign().Qual("net/url", "Parse").Call(
					jen.ID("exampleURI"),
				),
				jen.List(
					jen.ID("c"),
					jen.Err(),
				).Assign().ID("NewClient").Callln(
					utils.CtxVar(),
					jen.Lit(""),
					jen.Lit(""),
					jen.ID("u"),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Nil(),
					jen.Index().ID("string").Values(jen.Lit("*")),
					jen.ID("false"),
				),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("expected").Assign().Lit("ws://todo.verygoodsoftwarenotvirus.ru/api/v1/things/and/stuff"),
				jen.ID("actual").Assign().ID("c").Dot("BuildWebsocketURL").Call(
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
	}

	return lines
}

func buildTestV1Client_BuildHealthCheckRequest() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_BuildHealthCheckRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodGet"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildHealthCheckRequest").Call(utils.CtxVar()),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.Err(),
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
	}

	return lines
}

func buildTestV1Client_IsUp() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_IsUp").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
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
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.ID("actual").Assign().ID("c").Dot("IsUp").Call(utils.CtxVar()),
				utils.AssertTrue(
					jen.ID("actual"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"returns error with invalid URL",
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("IsUp").Call(utils.CtxVar()),
				utils.AssertFalse(jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with bad status code",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
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
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.ID("actual").Assign().ID("c").Dot("IsUp").Call(utils.CtxVar()),
				utils.AssertFalse(
					jen.ID("actual"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with timeout",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodGet"),
								nil,
							),
							jen.Qual("time", "Sleep").Call(
								jen.Lit(10).Times().Qual("time", "Hour"),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.ID("c").Dot("plainClient").Dot("Timeout").Equals().Lit(500).Times().Qual("time", "Millisecond"),
				jen.ID("actual").Assign().ID("c").Dot("IsUp").Call(utils.CtxVar()),
				utils.AssertFalse(
					jen.ID("actual"),
					nil,
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_buildDataRequest() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_buildDataRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			jen.ID("exampleData").Assign().VarPointer().ID("testingType").Values(jen.ID("Name").MapAssign().Lit("whatever")),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().ID("c").Dot("buildDataRequest").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("exampleData"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.AssertNoError(
					jen.Err(),
					nil,
				),
				utils.AssertEqual(
					jen.ID("expectedMethod"),
					jen.ID("req").Dot("Method"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid structure",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.ID("x").Assign().VarPointer().ID("testBreakableStruct").Values(jen.ID("Thing").MapAssign().Lit("stuff")),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().ID("c").Dot("buildDataRequest").Call(
					utils.CtxVar(),
					jen.Qual("net/http", "MethodPost"),
					jen.ID("ts").Dot("URL"),
					jen.ID("x"),
				),
				jen.Line(),
				utils.RequireNil(jen.ID("req"), nil),
				utils.AssertError(
					jen.Err(),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().ID("c").Dot("buildDataRequest").Call(
					utils.CtxVar(),
					jen.Qual("net/http", "MethodPost"),
					jen.ID("c").Dot("URL").Dot("String").Call(),
					jen.ID("exampleData"),
				),
				jen.Line(),
				utils.RequireNil(jen.ID("req"), nil),
				utils.AssertError(
					jen.Err(),
					nil,
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_checkExistence() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_checkExistence").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodHead"),
				jen.ID("ts").Assign().ID("httptest").Dot("NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK")),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("checkExistence").Call(
					utils.CtxVar(),
					jen.ID("req"),
				),
				utils.AssertTrue(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with timeout",
				utils.ExpectMethod("expectedMethod", "MethodHead"),
				jen.ID("ts").Assign().ID("httptest").Dot("NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							jen.Qual("time", "Sleep").Call(jen.Lit(10).Times().Qual("time", "Hour")),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("c").Dot("authedClient").Dot("Timeout").Equals().Lit(500).Times().Qual("time", "Millisecond"),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("checkExistence").Call(
					utils.CtxVar(),
					jen.ID("req"),
				),
				utils.AssertFalse(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
		),
	}

	return lines
}

func buildTestV1Client_retrieve() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_retrieve").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("exampleResponse").Assign().VarPointer().ID("argleBargle").Values(jen.ID("Name").MapAssign().Lit("whatever")),
				jen.Line(),
				jen.ID("ts").Assign().ID("httptest").Dot("NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.RequireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("exampleResponse")),
								nil,
							),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.Err().Equals().ID("c").Dot("retrieve").Call(
					utils.CtxVar(),
					jen.ID("req"),
					jen.VarPointer().ID("argleBargle").Values(),
				),
				utils.AssertNoError(
					jen.Err(),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with nil passed in",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.Qual("net/http", "MethodPost"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.Err().Equals().ID("c").Dot("retrieve").Call(
					utils.CtxVar(),
					jen.ID("req"),
					jen.Nil(),
				),
				utils.AssertError(
					jen.Err(),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with timeout",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							jen.Qual("time", "Sleep").Call(
								jen.Lit(10).Times().Qual("time", "Hour"),
							),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("c").Dot("authedClient").Dot("Timeout").Equals().Lit(500).Times().Qual("time", "Millisecond"),
				jen.Err().Equals().ID("c").Dot("retrieve").Call(
					utils.CtxVar(),
					jen.ID("req"),
					jen.VarPointer().ID("argleBargle").Values(),
				),
				utils.AssertError(
					jen.Err(),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with 404",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
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
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertEqual(
					jen.ID("ErrNotFound"),
					jen.ID("c").Dot("retrieve").Call(
						utils.CtxVar(),
						jen.ID("req"),
						jen.VarPointer().ID("argleBargle").Values(),
					),
					nil,
				),
			),
		),
	}

	return lines
}

func buildTestV1Client_executeRequest() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_executeRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			jen.ID("exampleResponse").Assign().VarPointer().ID("argleBargle").Values(jen.ID("Name").MapAssign().Lit("whatever")),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.RequireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(
									jen.ID("res"),
								).Dot("Encode").Call(jen.ID("exampleResponse")),
								nil,
							),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.Err().Equals().ID("c").Dot("executeRequest").Call(
					utils.CtxVar(),
					jen.ID("req"),
					jen.VarPointer().ID("argleBargle").Values(),
				),
				utils.AssertNoError(
					jen.Err(),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with timeout",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							jen.Qual("time", "Sleep").Call(jen.Lit(10).Times().Qual("time", "Hour")),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(
					jen.ID("req"),
					nil,
				),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("c").Dot("authedClient").Dot("Timeout").Equals().Lit(500).Times().Qual("time", "Millisecond"),
				jen.Err().Equals().ID("c").Dot("executeRequest").Call(
					utils.CtxVar(),
					jen.ID("req"),
					jen.VarPointer().ID("argleBargle").Values(),
				),
				utils.AssertError(jen.Err(), nil),
			),
			utils.BuildSubTest(
				"with 401",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.WriteHeader("StatusUnauthorized"),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertEqual(
					jen.ID("ErrUnauthorized"),
					jen.ID("c").Dot("executeRequest").Call(
						utils.CtxVar(),
						jen.ID("req"),
						jen.VarPointer().ID("argleBargle").Values(),
					),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with 404",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
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
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertEqual(
					jen.ID("ErrNotFound"),
					jen.ID("c").Dot("executeRequest").Call(
						utils.CtxVar(),
						jen.ID("req"),
						jen.VarPointer().ID("argleBargle").Values(),
					),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with unreadable response",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.RequireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(
									jen.ID("res"),
								).Dot("Encode").Call(jen.ID("exampleResponse")),
								nil,
							),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertError(
					jen.ID("c").Dot("executeRequest").Call(
						utils.CtxVar(),
						jen.ID("req"),
						jen.ID("argleBargle").Values(),
					),
					nil,
				),
			),
			jen.Line(),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_executeRawRequest() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_executeRawRequest").Block(
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
						jen.Lit(10).Times().Qual("time", "Hour"),
					),
				),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.List(
					jen.ID("res"),
					jen.Err(),
				).Assign().ID("c").Dot("executeRawRequest").Call(
					utils.CtxVar(),
					jen.VarPointer().Qual("net/http", "Client").Values(
						jen.ID("Timeout").MapAssign().Qual("time", "Second"),
					),
					jen.ID("req"),
				),
				utils.AssertNil(
					jen.ID("res"),
					nil,
				),
				utils.AssertError(
					jen.Err(),
					nil,
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_executeUnauthenticatedDataRequest() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_executeUnauthenticatedDataRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			jen.ID("exampleResponse").Assign().VarPointer().ID("argleBargle").Values(jen.ID("Name").MapAssign().Lit("whatever")),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.RequireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("exampleResponse")),
								nil,
							),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("in"),
					jen.ID("out"),
				).Assign().List(
					jen.VarPointer().ID("argleBargle").Values(),
					jen.VarPointer().ID("argleBargle").Values(),
				),
				jen.Line(),
				jen.List(
					jen.ID("body"),
					jen.Err(),
				).Assign().ID("createBodyFromStruct").Call(
					jen.ID("in"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("body"), nil),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.Err().Equals().ID("c").Dot("executeUnauthenticatedDataRequest").Call(
					utils.CtxVar(),
					jen.ID("req"),
					jen.ID("out"),
				),
				utils.AssertNoError(
					jen.Err(),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with 401",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.WriteHeader("StatusUnauthorized"),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("in"),
					jen.ID("out"),
				).Assign().List(
					jen.VarPointer().ID("argleBargle").Values(),
					jen.VarPointer().ID("argleBargle").Values(),
				),
				jen.Line(),
				jen.List(
					jen.ID("body"),
					jen.Err(),
				).Assign().ID("createBodyFromStruct").Call(
					jen.ID("in"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("body"), nil),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.Err().Equals().ID("c").Dot("executeUnauthenticatedDataRequest").Call(
					utils.CtxVar(),
					jen.ID("req"),
					jen.ID("out"),
				),
				utils.AssertError(
					jen.Err(),
					nil,
				),
				utils.AssertEqual(jen.ID("ErrUnauthorized"),
					jen.Err(), nil),
			),
			utils.BuildSubTest(
				"with 404",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
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
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("in"),
					jen.ID("out"),
				).Assign().List(
					jen.VarPointer().ID("argleBargle").Values(),
					jen.VarPointer().ID("argleBargle").Values(),
				),
				jen.Line(),
				jen.List(
					jen.ID("body"),
					jen.Err(),
				).Assign().ID("createBodyFromStruct").Call(
					jen.ID("in"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("body"), nil),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.Err().Equals().ID("c").Dot("executeUnauthenticatedDataRequest").Call(
					utils.CtxVar(),
					jen.ID("req"),
					jen.ID("out"),
				),
				utils.AssertError(
					jen.Err(),
					nil,
				),
				utils.AssertEqual(jen.ID("ErrNotFound"),
					jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with timeout",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							jen.Qual("time", "Sleep").Call(
								jen.Lit(10).Times().Qual("time", "Hour"),
							),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("in"),
					jen.ID("out"),
				).Assign().List(
					jen.VarPointer().ID("argleBargle").Values(),
					jen.VarPointer().ID("argleBargle").Values(),
				),
				jen.Line(),
				jen.List(
					jen.ID("body"),
					jen.Err(),
				).Assign().ID("createBodyFromStruct").Call(
					jen.ID("in"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("body"), nil),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.ID("c").Dot("plainClient").Dot("Timeout").Equals().Lit(500).Times().Qual("time", "Millisecond"),
				utils.AssertError(
					jen.ID("c").Dot("executeUnauthenticatedDataRequest").Call(
						utils.CtxVar(),
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
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.ID("in").Assign().VarPointer().ID("argleBargle").Values(),
				jen.Line(),
				jen.List(
					jen.ID("body"),
					jen.Err(),
				).Assign().ID("createBodyFromStruct").Call(
					jen.ID("in"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("body"), nil),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.Err().Equals().ID("c").Dot("executeUnauthenticatedDataRequest").Call(
					utils.CtxVar(),
					jen.ID("req"),
					jen.ID("testingType").Values(),
				),
				utils.AssertError(
					jen.Err(),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with unreadable response",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.RequireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("exampleResponse")),
								nil,
							),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					utils.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertError(
					jen.ID("c").Dot("executeUnauthenticatedDataRequest").Call(
						utils.CtxVar(),
						jen.ID("req"),
						jen.ID("argleBargle").Values(),
					),
					nil,
				),
			),
		),
		jen.Line(),
	}

	return lines
}
