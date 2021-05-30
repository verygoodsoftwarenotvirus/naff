package httpclient

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func helpersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("testingType").Struct(jen.ID("Name").ID("string")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestErrorFromResponse").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("returns error for nil response"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("errorFromResponse").Call(jen.ID("nil")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestArgIsNotPointerOrNil").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("expected use"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("err").Op(":=").ID("argIsNotPointerOrNil").Call(jen.Op("&").ID("testingType").Valuesln()),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("error should not be returned when a pointer is provided"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with non-pointer"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("err").Op(":=").ID("argIsNotPointerOrNil").Call(jen.ID("testingType").Valuesln()),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("error should be returned when a non-pointer is provided"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("err").Op(":=").ID("argIsNotPointerOrNil").Call(jen.ID("nil")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("error should be returned when nil is provided"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestArgIsNotPointer").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("expected use"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("notAPointer"), jen.ID("err")).Op(":=").ID("argIsNotPointer").Call(jen.Op("&").ID("testingType").Valuesln()),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("notAPointer"),
						jen.Lit("expected `false` when a pointer is provided"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("error should not be returned when a pointer is provided"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with non-pointer"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("notAPointer"), jen.ID("err")).Op(":=").ID("argIsNotPointer").Call(jen.ID("testingType").Valuesln()),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("notAPointer"),
						jen.Lit("expected `true` when a non-pointer is provided"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("error should be returned when a non-pointer is provided"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("notAPointer"), jen.ID("err")).Op(":=").ID("argIsNotPointer").Call(jen.ID("nil")),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("notAPointer"),
						jen.Lit("expected `true` when nil is provided"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("error should be returned when nil is provided"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestArgIsNotNil").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without nil"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("isNil"), jen.ID("err")).Op(":=").ID("argIsNotNil").Call(jen.Op("&").ID("testingType").Valuesln()),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("isNil"),
						jen.Lit("expected `false` when a pointer is provided"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("error should not be returned when a pointer is provided"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with non-pointer"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("isNil"), jen.ID("err")).Op(":=").ID("argIsNotNil").Call(jen.ID("testingType").Valuesln()),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("isNil"),
						jen.Lit("expected `true` when a non-pointer is provided"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("error should not be returned when a non-pointer is provided"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("isNil"), jen.ID("err")).Op(":=").ID("argIsNotNil").Call(jen.ID("nil")),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("isNil"),
						jen.Lit("expected `true` when nil is provided"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("error should be returned when nil is provided"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestUnmarshalBody").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("expected").Op(":=").Lit("whatever"),
					jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").Valuesln(jen.ID("Body").Op(":").Qual("io", "NopCloser").Call(jen.Qual("strings", "NewReader").Call(jen.Qual("fmt", "Sprintf").Call(
						jen.Lit(`{"name": %q}`),
						jen.ID("expected"),
					))), jen.ID("StatusCode").Op(":").Qual("net/http", "StatusOK")),
					jen.Var().Defs(
						jen.ID("out").ID("testingType"),
					),
					jen.ID("err").Op(":=").ID("c").Dot("unmarshalBody").Call(
						jen.ID("ctx"),
						jen.ID("res"),
						jen.Op("&").ID("out"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("out").Dot("Name"),
						jen.ID("expected"),
						jen.Lit("expected marshaling to work"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("no error should be encountered unmarshalling into a valid struct"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with good status but unmarshallable response"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").Valuesln(jen.ID("Body").Op(":").Qual("io", "NopCloser").Call(jen.Qual("strings", "NewReader").Call(jen.Lit("BLAH"))), jen.ID("StatusCode").Op(":").Qual("net/http", "StatusOK")),
					jen.Var().Defs(
						jen.ID("out").ID("testingType"),
					),
					jen.ID("err").Op(":=").ID("c").Dot("unmarshalBody").Call(
						jen.ID("ctx"),
						jen.ID("res"),
						jen.Op("&").ID("out"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("error should be encountered unmarshalling invalid response into a valid struct"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with an erroneous error code"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").Valuesln(jen.ID("Body").Op(":").Qual("io", "NopCloser").Call(jen.Qual("strings", "NewReader").Call(jen.Func().Params().Params(jen.ID("string")).Body(
						jen.List(jen.ID("bs"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Op("&").ID("types").Dot("ErrorResponse").Valuesln()),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.Return().ID("string").Call(jen.ID("bs")),
					).Call())), jen.ID("StatusCode").Op(":").Qual("net/http", "StatusBadRequest")),
					jen.Var().Defs(
						jen.ID("out").Op("*").ID("testingType"),
					),
					jen.ID("err").Op(":=").ID("c").Dot("unmarshalBody").Call(
						jen.ID("ctx"),
						jen.ID("res"),
						jen.Op("&").ID("out"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("out"),
						jen.Lit("expected nil to be returned"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("error should be returned from the API"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with an erroneous error code and unmarshallable body"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").Valuesln(jen.ID("Body").Op(":").Qual("io", "NopCloser").Call(jen.Qual("strings", "NewReader").Call(jen.Lit("BLAH"))), jen.ID("StatusCode").Op(":").Qual("net/http", "StatusBadRequest")),
					jen.Var().Defs(
						jen.ID("out").Op("*").ID("testingType"),
					),
					jen.ID("err").Op(":=").ID("c").Dot("unmarshalBody").Call(
						jen.ID("ctx"),
						jen.ID("res"),
						jen.Op("&").ID("out"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("out"),
						jen.Lit("expected nil to be returned"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("error should be returned from the unmarshaller"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil target variable"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("unmarshalBody").Call(
						jen.ID("ctx"),
						jen.ID("nil"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("error should be encountered when passed nil"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with erroneous reader"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("expected").Op(":=").Qual("errors", "New").Call(jen.Lit("blah")),
					jen.ID("rc").Op(":=").ID("newMockReadCloser").Call(),
					jen.ID("rc").Dot("On").Call(
						jen.Lit("Read"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().ID("uint8").Valuesln()),
					).Dot("Return").Call(
						jen.Lit(0),
						jen.ID("expected"),
					),
					jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").Valuesln(jen.ID("Body").Op(":").ID("rc"), jen.ID("StatusCode").Op(":").Qual("net/http", "StatusOK")),
					jen.Var().Defs(
						jen.ID("out").ID("testingType"),
					),
					jen.ID("err").Op(":=").ID("c").Dot("unmarshalBody").Call(
						jen.ID("ctx"),
						jen.ID("res"),
						jen.Op("&").ID("out"),
					),
					jen.ID("assertErrorMatches").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.ID("expected"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("no error should be encountered unmarshalling into a valid struct"),
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

	return code
}
