package requests

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func builderTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("testingType").Struct(jen.ID("Name").ID("string")),
			jen.ID("testBreakableStruct").Struct(jen.ID("Thing").Qual("encoding/json", "Number")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestNewBuilder").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("encoder").Op(":=").ID("encoding").Dot("ProvideClientEncoder").Call(
						jen.ID("logger"),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewBuilder").Call(
						jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
						jen.ID("logger"),
						jen.ID("encoder"),
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
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("encoder").Op(":=").ID("encoding").Dot("ProvideClientEncoder").Call(
						jen.ID("logger"),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewBuilder").Call(
						jen.ID("nil"),
						jen.ID("logger"),
						jen.ID("encoder"),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil encoder"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewBuilder").Call(
						jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
						jen.ID("logger"),
						jen.ID("nil"),
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
		jen.Func().ID("TestBuilder_URL").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("builder").Dot("URL").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuilder_SetURL").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("builder").Dot("SetURL").Call(jen.Op("&").Qual("net/url", "URL").Valuesln()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("builder").Dot("SetURL").Call(jen.ID("nil")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuilder_BuildURL").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("various urls"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("encoder").Op(":=").ID("encoding").Dot("ProvideClientEncoder").Call(
						jen.ID("logger"),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("NewBuilder").Call(
						jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
						jen.ID("logger"),
						jen.ID("encoder"),
					),
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
					jen.ID("c").Op(":=").ID("buildTestRequestBuilderWithInvalidURL").Call(),
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
		jen.Func().ID("TestBuilder_Must").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("helper").Dot("builder").Dot("Must").Call(
						jen.Op("&").Qual("net/http", "Request").Valuesln(),
						jen.ID("nil"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with panic"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleErr").Op(":=").Qual("errors", "New").Call(jen.Lit("blah")),
					jen.ID("mockPanicker").Op(":=").ID("panicking").Dot("NewMockPanicker").Call(),
					jen.ID("mockPanicker").Dot("On").Call(
						jen.Lit("Panic"),
						jen.ID("exampleErr"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("builder").Dot("panicker").Op("=").ID("mockPanicker"),
					jen.ID("helper").Dot("builder").Dot("Must").Call(
						jen.Op("&").Qual("net/http", "Request").Valuesln(),
						jen.ID("exampleErr"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockPanicker"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_buildRawURL").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("buildRawURL").Call(
						jen.ID("parsedExampleURL"),
						jen.Qual("net/url", "Values").Valuesln(),
						jen.ID("true"),
						jen.Lit("things"),
						jen.Lit("and"),
						jen.Lit("stuff"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuilder_buildAPIV1URL").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("builder").Dot("buildAPIV1URL").Call(
							jen.ID("helper").Dot("ctx"),
							jen.Qual("net/url", "Values").Valuesln(),
							jen.Lit("things"),
							jen.Lit("and"),
							jen.Lit("stuff"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuilder_buildUnversionedURL").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("various urls"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("encoder").Op(":=").ID("encoding").Dot("ProvideClientEncoder").Call(
						jen.ID("logger"),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.List(jen.ID("b"), jen.ID("err")).Op(":=").ID("NewBuilder").Call(
						jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
						jen.ID("logger"),
						jen.ID("encoder"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("testCases").Op(":=").Index().Struct(
						jen.ID("inputQuery").ID("valuer"),
						jen.ID("expectation").ID("string"),
						jen.ID("inputParts").Index().ID("string"),
					).Valuesln(jen.Valuesln(jen.ID("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/things"), jen.ID("inputParts").Op(":").Index().ID("string").Valuesln(jen.Lit("things"))), jen.Valuesln(jen.ID("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/stuff?key=value"), jen.ID("inputQuery").Op(":").Map(jen.ID("string")).Index().ID("string").Valuesln(jen.Lit("key").Op(":").Valuesln(jen.Lit("value"))), jen.ID("inputParts").Op(":").Index().ID("string").Valuesln(jen.Lit("stuff"))), jen.Valuesln(jen.ID("expectation").Op(":").Lit("https://todo.verygoodsoftwarenotvirus.ru/things/and/stuff?key=value1&key=value2&yek=eulav"), jen.ID("inputQuery").Op(":").Map(jen.ID("string")).Index().ID("string").Valuesln(jen.Lit("key").Op(":").Valuesln(jen.Lit("value1"), jen.Lit("value2")), jen.Lit("yek").Op(":").Valuesln(jen.Lit("eulav"))), jen.ID("inputParts").Op(":").Index().ID("string").Valuesln(jen.Lit("things"), jen.Lit("and"), jen.Lit("stuff")))),
					jen.For(jen.List(jen.ID("_"), jen.ID("tc")).Op(":=").Range().ID("testCases")).Body(
						jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
						jen.ID("actual").Op(":=").ID("b").Dot("buildUnversionedURL").Call(
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
					jen.ID("c").Op(":=").ID("buildTestRequestBuilderWithInvalidURL").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("actual").Op(":=").ID("c").Dot("buildUnversionedURL").Call(
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
		jen.Func().ID("TestBuilder_BuildWebsocketURL").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("encoder").Op(":=").ID("encoding").Dot("ProvideClientEncoder").Call(
						jen.ID("logger"),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewBuilder").Call(
						jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
						jen.ID("logger"),
						jen.ID("encoder"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("expected").Op(":=").Lit("ws://todo.verygoodsoftwarenotvirus.ru/api/v1/things/and/stuff"),
					jen.ID("actual").Op(":=").ID("c").Dot("BuildWebsocketURL").Call(
						jen.ID("ctx"),
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
		jen.Func().ID("TestBuilder_BuildHealthCheckRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expectedMethod").Op(":=").Qual("net/http", "MethodGet"),
					jen.ID("c").Op(":=").ID("buildTestRequestBuilder").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("BuildHealthCheckRequest").Call(jen.ID("ctx")),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("actual").Dot("Method"),
						jen.ID("expectedMethod"),
						jen.Lit("request should be a %s request"),
						jen.ID("expectedMethod"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuilder_buildDataRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("exampleData").Op(":=").Op("&").ID("testingType").Valuesln(jen.ID("Name").Op(":").Lit("whatever")),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("c").Op(":=").ID("buildTestRequestBuilder").Call(),
					jen.ID("expectedMethod").Op(":=").Qual("net/http", "MethodPost"),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("buildDataRequest").Call(
						jen.ID("ctx"),
						jen.ID("expectedMethod"),
						jen.ID("exampleURI"),
						jen.ID("exampleData"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedMethod"),
						jen.ID("req").Dot("Method"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleURI"),
						jen.ID("req").Dot("URL").Dot("String").Call(),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid structure"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("c").Op(":=").ID("buildTestRequestBuilder").Call(),
					jen.ID("x").Op(":=").Op("&").ID("testBreakableStruct").Valuesln(jen.ID("Thing").Op(":").Lit("stuff")),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("buildDataRequest").Call(
						jen.ID("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.ID("exampleURI"),
						jen.ID("x"),
					),
					jen.ID("require").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid client url"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("c").Op(":=").ID("buildTestRequestBuilderWithInvalidURL").Call(),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("buildDataRequest").Call(
						jen.ID("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.ID("c").Dot("url").Dot("String").Call(),
						jen.ID("exampleData"),
					),
					jen.ID("require").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("req"),
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
		jen.Func().ID("Test_mustParseURL").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("mustParseURL").Call(jen.ID("exampleURI")),
				),
			),
		),
		jen.Line(),
	)

	return code
}
