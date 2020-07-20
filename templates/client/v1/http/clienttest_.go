package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)
	utils.AddImports(proj, code)

	// vars
	code.Add(
		jen.Const().Defs(
			jen.ID(utils.BuildFakeVarName("URI")).Equals().Lit("https://todo.verygoodsoftwarenotvirus.ru"),
			jen.ID("asciiControlChar").Equals().String().Call(jen.Byte().Call(jen.Lit(0x7f))),
		),
	)

	// types
	code.Add(
		jen.Type().Defs(
			jen.ID("argleBargle").Struct(
				jen.ID("Name").String(),
			),
			jen.Line(),
			jen.ID("valuer").Map(jen.String()).Index().String(),
		),
		jen.Line(),
	)

	code.Add(
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
	code.Add(
		jen.Line(),
		jen.Comment("begin helper funcs"),
		jen.Line(),
		jen.Line(),
	)

	code.Add(buildMustParseURL()...)
	code.Add(buildBuildTestClient()...)
	code.Add(buildBuildTestClientWithInvalidURL()...)

	code.Add(
		jen.Line(),
		jen.Line(),
		jen.Comment("end helper funcs"),
		jen.Line(),
	)

	code.Add(buildTestV1Client_AuthenticatedClient()...)
	code.Add(buildTestV1Client_PlainClient()...)
	code.Add(buildTestV1Client_TokenSource()...)
	code.Add(buildTestNewClient()...)
	code.Add(buildTestNewSimpleClient()...)
	code.Add(buildTestV1Client_CloseRequestBody()...)
	code.Add(buildTestBuildURL()...)
	code.Add(buildTestBuildVersionlessURL()...)
	code.Add(buildTestV1Client_BuildWebsocketURL()...)
	code.Add(buildTestV1Client_BuildHealthCheckRequest()...)
	code.Add(buildTestV1Client_IsUp()...)
	code.Add(buildTestV1Client_buildDataRequest()...)
	code.Add(buildTestV1Client_executeRequest()...)
	code.Add(buildTestV1Client_executeRawRequest()...)
	code.Add(buildTestV1Client_checkExistence()...)
	code.Add(buildTestV1Client_retrieve()...)
	code.Add(buildTestV1Client_executeUnauthenticatedDataRequest()...)

	return code
}

func buildMustParseURL() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("mustParseURL").Params(
			jen.ID("uri").String(),
		).Params(
			jen.PointerTo().Qual("net/url", "URL"),
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
			jen.ID("t").PointerTo().Qual("testing", "T"),
			jen.ID("ts").PointerTo().Qual("net/http/httptest", "Server"),
		).Params(
			jen.PointerTo().ID(v1),
		).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("l").Assign().Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.ID("u").Assign().ID("mustParseURL").Call(
				jen.ID("ts").Dot("URL"),
			),
			jen.ID("c").Assign().ID("ts").Dot("Client").Call(),
			jen.Line(),
			jen.Return().AddressOf().ID(v1).Valuesln(
				jen.ID("URL").MapAssign().ID("u"),
				jen.ID("plainClient").MapAssign().ID("c"),
				jen.ID(constants.LoggerVarName).MapAssign().ID("l"),
				jen.ID("Debug").MapAssign().True(),
				jen.ID("authedClient").MapAssign().ID("c"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildTestClientWithInvalidURL() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildTestClientWithInvalidURL").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Params(jen.PointerTo().ID(v1)).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("l").Assign().Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.ID("u").Assign().ID("mustParseURL").Call(jen.Lit("https://verygoodsoftwarenotvirus.ru")),
			jen.ID("u").Dot("Scheme").Equals().Qual("fmt", "Sprintf").Call(
				jen.RawString(`%s://`),
				jen.ID("asciiControlChar"),
			),
			jen.Line(),
			jen.Return().AddressOf().ID(v1).Valuesln(
				jen.ID("URL").MapAssign().ID("u"),
				jen.ID("plainClient").MapAssign().Qual("net/http", "DefaultClient"),
				jen.ID(constants.LoggerVarName).MapAssign().ID("l"),
				jen.ID("Debug").MapAssign().True(),
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
				jen.Line(),
				jen.List(jen.ID("c"), jen.Err()).Assign().ID("NewClient").Callln(
					constants.CtxVar(),
					jen.EmptyString(),
					jen.EmptyString(),
					jen.ID("mustParseURL").Call(
						jen.ID(utils.BuildFakeVarName("URI")),
					),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("ts").Dot("Client").Call(),
					jen.Index().String().Values(jen.Lit("*")),
					jen.False(),
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
				jen.Line(),
				jen.List(jen.ID("c"), jen.Err()).Assign().ID("NewClient").Callln(
					constants.CtxVar(),
					jen.EmptyString(),
					jen.EmptyString(),
					jen.ID("mustParseURL").Call(
						jen.ID(utils.BuildFakeVarName("URI")),
					),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("ts").Dot("Client").Call(),
					jen.Index().String().Values(jen.Lit("*")),
					jen.False(),
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
					constants.CtxVar(),
					jen.EmptyString(),
					jen.EmptyString(),
					jen.ID("mustParseURL").Call(
						jen.ID(utils.BuildFakeVarName("URI")),
					),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.AddressOf().Qual("net/http", "Client").Valuesln(
						jen.ID("Timeout").MapAssign().Zero(),
					),
					jen.Index().String().Values(jen.Lit("*")),
					jen.True(),
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
					constants.CtxVar(),
					jen.ID("mustParseURL").Call(jen.ID(utils.BuildFakeVarName("URI"))),
					jen.True(),
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
				jen.ID("rc").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().AddressOf().Qual("net/http", "Response").Valuesln(
					jen.ID("Body").MapAssign().ID("rc"),
					jen.ID("StatusCode").MapAssign().Qual("net/http", "StatusOK"),
				),
				jen.Line(),
				jen.List(
					jen.ID("c"),
					jen.Err(),
				).Assign().ID("NewSimpleClient").Callln(
					constants.CtxVar(),
					jen.ID("mustParseURL").Call(jen.ID(utils.BuildFakeVarName("URI"))),
					jen.True(),
				),
				utils.AssertNotNil(jen.ID("c"), nil),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("c").Dot("closeResponseBody").Call(jen.ID(constants.ResponseVarName)),
				jen.Line(),
				utils.AssertExpectationsFor("rc"),
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
				jen.Line(),
				jen.List(jen.ID("u"), jen.Underscore()).Assign().Qual("net/url", "Parse").Call(
					jen.ID(utils.BuildFakeVarName("URI")),
				),
				jen.Line(),
				jen.List(
					jen.ID("c"),
					jen.Err(),
				).Assign().ID("NewClient").Callln(
					constants.CtxVar(),
					jen.EmptyString(),
					jen.EmptyString(),
					jen.ID("u"),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Nil(),
					jen.Index().String().Values(jen.Lit("*")),
					jen.False(),
				),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("testCases").Assign().Index().Struct(
					jen.ID("expectation").String(),
					jen.ID("inputParts").Index().String(),
					jen.ID("inputQuery").ID("valuer"),
				).Valuesln(
					jen.Valuesln(
						jen.ID("expectation").MapAssign().Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/things"),
						jen.ID("inputParts").MapAssign().Index().String().Values(jen.Lit("things")),
					),
					jen.Valuesln(
						jen.ID("expectation").MapAssign().Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/stuff?key=value"),
						jen.ID("inputQuery").MapAssign().Map(jen.String()).Index().String().Values(
							jen.Lit("key").MapAssign().Values(jen.Lit("value")),
						),
						jen.ID("inputParts").MapAssign().Index().String().Values(jen.Lit("stuff")),
					),
					jen.Valuesln(
						jen.ID("expectation").MapAssign().Lit("https://todo.verygoodsoftwarenotvirus.ru/api/v1/things/and/stuff?key=value1&key=value2&yek=eulav"),
						jen.ID("inputQuery").MapAssign().Map(jen.String()).Index().String().Valuesln(
							jen.Lit("key").MapAssign().Values(
								jen.Lit("value1"),
								jen.Lit("value2"),
							),
							jen.Lit("yek").MapAssign().Values(jen.Lit("eulav")),
						),
						jen.ID("inputParts").MapAssign().Index().String().Values(jen.Lit("things"),
							jen.Lit("and"),
							jen.Lit("stuff"),
						),
					),
				),
				jen.Line(),
				jen.For(jen.List(
					jen.Underscore(),
					jen.ID("tc"),
				).Assign().Range().ID("testCases"),
				).Block(
					jen.ID("actual").Assign().ID("c").Dot("BuildURL").Call(
						jen.ID("tc").Dot("inputQuery").Dot("ToValues").Call(),
						jen.ID("tc").Dot("inputParts").Spread(),
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
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(testVar()),
				jen.Qual(constants.AssertPkg, "Empty").Call(
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
				jen.Line(),
				jen.List(jen.ID("u"), jen.Underscore()).Assign().Qual("net/url", "Parse").Call(
					jen.ID(utils.BuildFakeVarName("URI")),
				),
				jen.Line(),
				jen.List(
					jen.ID("c"),
					jen.Err(),
				).Assign().ID("NewClient").Callln(
					constants.CtxVar(),
					jen.EmptyString(),
					jen.EmptyString(),
					jen.ID("u"),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Nil(),
					jen.Index().String().Values(jen.Lit("*")),
					jen.False(),
				),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("testCases").Assign().Index().Struct(
					jen.ID("expectation").String(),
					jen.ID("inputParts").Index().String(),
					jen.ID("inputQuery").ID("valuer"),
				).Valuesln(
					jen.Valuesln(
						jen.ID("expectation").MapAssign().Lit("https://todo.verygoodsoftwarenotvirus.ru/things"),
						jen.ID("inputParts").MapAssign().Index().String().Values(jen.Lit("things")),
					),
					jen.Valuesln(
						jen.ID("expectation").MapAssign().Lit("https://todo.verygoodsoftwarenotvirus.ru/stuff?key=value"),
						jen.ID("inputQuery").MapAssign().Map(jen.String()).Index().String().Values(
							jen.Lit("key").MapAssign().Values(jen.Lit("value")),
						),
						jen.ID("inputParts").MapAssign().Index().String().Values(jen.Lit("stuff")),
					),
					jen.Valuesln(
						jen.ID("expectation").MapAssign().Lit("https://todo.verygoodsoftwarenotvirus.ru/things/and/stuff?key=value1&key=value2&yek=eulav"),
						jen.ID("inputQuery").MapAssign().Map(jen.String()).Index().String().Valuesln(
							jen.Lit("key").MapAssign().Values(
								jen.Lit("value1"),
								jen.Lit("value2"),
							),
							jen.Lit("yek").MapAssign().Values(jen.Lit("eulav")),
						),
						jen.ID("inputParts").MapAssign().Index().String().Values(jen.Lit("things"),
							jen.Lit("and"),
							jen.Lit("stuff"),
						),
					),
				),
				jen.Line(),
				jen.For(jen.List(
					jen.Underscore(),
					jen.ID("tc"),
				).Assign().Range().ID("testCases"),
				).Block(
					jen.ID("actual").Assign().ID("c").Dot("buildVersionlessURL").Call(
						jen.ID("tc").Dot("inputQuery").Dot("ToValues").Call(),
						jen.ID("tc").Dot("inputParts").Spread(),
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
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(testVar()),
				jen.Qual(constants.AssertPkg, "Empty").Call(
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
					jen.Underscore(),
				).Assign().Qual("net/url", "Parse").Call(
					jen.ID(utils.BuildFakeVarName("URI")),
				),
				jen.Line(),
				jen.List(
					jen.ID("c"),
					jen.Err(),
				).Assign().ID("NewClient").Callln(
					constants.CtxVar(),
					jen.EmptyString(),
					jen.EmptyString(),
					jen.ID("u"),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Nil(),
					jen.Index().String().Values(jen.Lit("*")),
					jen.False(),
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
				).Assign().ID("c").Dot("BuildHealthCheckRequest").Call(constants.CtxVar()),
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
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
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
				jen.ID("actual").Assign().ID("c").Dot("IsUp").Call(constants.CtxVar()),
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
				jen.ID("actual").Assign().ID("c").Dot("IsUp").Call(constants.CtxVar()),
				utils.AssertFalse(jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with bad status code",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
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
				jen.ID("actual").Assign().ID("c").Dot("IsUp").Call(constants.CtxVar()),
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
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
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
				jen.ID("actual").Assign().ID("c").Dot("IsUp").Call(constants.CtxVar()),
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
			jen.ID(utils.BuildFakeVarName("Data")).Assign().AddressOf().ID("testingType").Values(jen.ID("Name").MapAssign().Lit("whatever")),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().ID("c").Dot("buildDataRequest").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID(utils.BuildFakeVarName("Data")),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				utils.AssertNoError(
					jen.Err(),
					nil,
				),
				utils.AssertEqual(
					jen.ID("expectedMethod"),
					jen.ID(constants.RequestVarName).Dot("Method"),
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
				jen.ID("x").Assign().AddressOf().ID("testBreakableStruct").Values(jen.ID("Thing").MapAssign().Lit("stuff")),
				jen.List(
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().ID("c").Dot("buildDataRequest").Call(
					constants.CtxVar(),
					jen.Qual("net/http", "MethodPost"),
					jen.ID("ts").Dot("URL"),
					jen.ID("x"),
				),
				jen.Line(),
				utils.RequireNil(jen.ID(constants.RequestVarName), nil),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().ID("c").Dot("buildDataRequest").Call(
					constants.CtxVar(),
					jen.Qual("net/http", "MethodPost"),
					jen.ID("c").Dot("URL").Dot("String").Call(),
					jen.ID(utils.BuildFakeVarName("Data")),
				),
				jen.Line(),
				utils.RequireNil(jen.ID(constants.RequestVarName), nil),
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
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK")),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("checkExistence").Call(
					constants.CtxVar(),
					jen.ID(constants.RequestVarName),
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
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("c").Dot("authedClient").Dot("Timeout").Equals().Lit(500).Times().Qual("time", "Millisecond"),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("checkExistence").Call(
					constants.CtxVar(),
					jen.ID(constants.RequestVarName),
				),
				utils.AssertFalse(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
			),
		),
		jen.Line(),
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
				jen.ID(utils.BuildFakeVarName("Response")).Assign().AddressOf().ID("argleBargle").Values(jen.ID("Name").MapAssign().Lit("whatever")),
				jen.Line(),
				jen.ID("ts").Assign().ID("httptest").Dot("NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.RequireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName("Response"))),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.Err().Equals().ID("c").Dot("retrieve").Call(
					constants.CtxVar(),
					jen.ID(constants.RequestVarName),
					jen.AddressOf().ID("argleBargle").Values(),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.Qual("net/http", "MethodPost"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.Err().Equals().ID("c").Dot("retrieve").Call(
					constants.CtxVar(),
					jen.ID(constants.RequestVarName),
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
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("c").Dot("authedClient").Dot("Timeout").Equals().Lit(500).Times().Qual("time", "Millisecond"),
				jen.Err().Equals().ID("c").Dot("retrieve").Call(
					constants.CtxVar(),
					jen.ID(constants.RequestVarName),
					jen.AddressOf().ID("argleBargle").Values(),
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
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertEqual(
					jen.ID("ErrNotFound"),
					jen.ID("c").Dot("retrieve").Call(
						constants.CtxVar(),
						jen.ID(constants.RequestVarName),
						jen.AddressOf().ID("argleBargle").Values(),
					),
					nil,
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1Client_executeRequest() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_executeRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			jen.ID(utils.BuildFakeVarName("Response")).Assign().AddressOf().ID("argleBargle").Values(jen.ID("Name").MapAssign().Lit("whatever")),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.RequireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(
									jen.ID(constants.ResponseVarName),
								).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName("Response"))),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.Err().Equals().ID("c").Dot("executeRequest").Call(
					constants.CtxVar(),
					jen.ID(constants.RequestVarName),
					jen.AddressOf().ID("argleBargle").Values(),
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
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(
					jen.ID(constants.RequestVarName),
					nil,
				),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("c").Dot("authedClient").Dot("Timeout").Equals().Lit(500).Times().Qual("time", "Millisecond"),
				jen.Err().Equals().ID("c").Dot("executeRequest").Call(
					constants.CtxVar(),
					jen.ID(constants.RequestVarName),
					jen.AddressOf().ID("argleBargle").Values(),
				),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with 401",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertEqual(
					jen.ID("ErrUnauthorized"),
					jen.ID("c").Dot("executeRequest").Call(
						constants.CtxVar(),
						jen.ID(constants.RequestVarName),
						jen.AddressOf().ID("argleBargle").Values(),
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
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertEqual(
					jen.ID("ErrNotFound"),
					jen.ID("c").Dot("executeRequest").Call(
						constants.CtxVar(),
						jen.ID(constants.RequestVarName),
						jen.AddressOf().ID("argleBargle").Values(),
					),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with unreadable response",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.RequireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(
									jen.ID(constants.ResponseVarName),
								).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName("Response"))),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertError(
					jen.ID("c").Dot("executeRequest").Call(
						constants.CtxVar(),
						jen.ID(constants.RequestVarName),
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
						jen.ID(constants.RequestVarName).Dot("Method"),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.List(
					jen.ID(constants.ResponseVarName),
					jen.Err(),
				).Assign().ID("c").Dot("executeRawRequest").Call(
					constants.CtxVar(),
					jen.AddressOf().Qual("net/http", "Client").Values(
						jen.ID("Timeout").MapAssign().Qual("time", "Second"),
					),
					jen.ID(constants.RequestVarName),
				),
				utils.AssertNil(
					jen.ID(constants.ResponseVarName),
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
			jen.ID(utils.BuildFakeVarName("Response")).Assign().AddressOf().ID("argleBargle").Values(jen.ID("Name").MapAssign().Lit("whatever")),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.RequireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName("Response"))),
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
					jen.AddressOf().ID("argleBargle").Values(),
					jen.AddressOf().ID("argleBargle").Values(),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.Err().Equals().ID("c").Dot("executeUnauthenticatedDataRequest").Call(
					constants.CtxVar(),
					jen.ID(constants.RequestVarName),
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
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
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
					jen.AddressOf().ID("argleBargle").Values(),
					jen.AddressOf().ID("argleBargle").Values(),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.Err().Equals().ID("c").Dot("executeUnauthenticatedDataRequest").Call(
					constants.CtxVar(),
					jen.ID(constants.RequestVarName),
					jen.ID("out"),
				),
				utils.AssertError(
					jen.Err(),
					nil,
				),
				utils.AssertEqual(jen.ID("ErrUnauthorized"),
					jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with 404",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
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
					jen.AddressOf().ID("argleBargle").Values(),
					jen.AddressOf().ID("argleBargle").Values(),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.Err().Equals().ID("c").Dot("executeUnauthenticatedDataRequest").Call(
					constants.CtxVar(),
					jen.ID(constants.RequestVarName),
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
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
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
					jen.AddressOf().ID("argleBargle").Values(),
					jen.AddressOf().ID("argleBargle").Values(),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID("c").Dot("plainClient").Dot("Timeout").Equals().Lit(500).Times().Qual("time", "Millisecond"),
				utils.AssertError(
					jen.ID("c").Dot("executeUnauthenticatedDataRequest").Call(
						constants.CtxVar(),
						jen.ID(constants.RequestVarName),
						jen.ID("out"),
					),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with nil as output",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.ID("in").Assign().AddressOf().ID("argleBargle").Values(),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.ID("body"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.Err().Equals().ID("c").Dot("executeUnauthenticatedDataRequest").Call(
					constants.CtxVar(),
					jen.ID(constants.RequestVarName),
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
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
							jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID(constants.RequestVarName).Dot("Method"),
								jen.ID("expectedMethod"),
								nil,
							),
							utils.RequireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName("Response"))),
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
					jen.ID(constants.RequestVarName),
					jen.Err(),
				).Assign().Qual("net/http", "NewRequestWithContext").Call(
					constants.CtxVar(),
					jen.ID("expectedMethod"),
					jen.ID("ts").Dot("URL"),
					jen.Nil(),
				),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertError(
					jen.ID("c").Dot("executeUnauthenticatedDataRequest").Call(
						constants.CtxVar(),
						jen.ID(constants.RequestVarName),
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
