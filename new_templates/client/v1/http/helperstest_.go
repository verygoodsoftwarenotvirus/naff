package client

import jen "github.com/dave/jennifer/jen"

func helpersTestDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

	ret.Add(
		jen.Type().Id("testingType").Struct(
			jen.Id("Name").Id("string"),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("ArgIsNotPointerOrNil").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"expected use",
				jen.Id("err").Op(":=").Id("argIsNotPointerOrNil").Call(
					jen.Op("&").Id("testingType").Values(),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("error should not be returned when a pointer is provided"),
				),
			),
			buildSubTest(
				"with non-pointer",
				jen.Id("err").Op(":=").Id("argIsNotPointerOrNil").Call(
					jen.Id("testingType").Values(),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("error should be returned when a non-pointer is provided"),
				),
			),
			buildSubTest(
				"with nil",
				jen.Id("err").Op(":=").Id("argIsNotPointerOrNil").Call(
					jen.Id("nil"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("error should be returned when nil is provided"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("ArgIsNotPointer").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"expected use",
				jen.List(
					jen.Id("notAPointer"),
					jen.Id("err"),
				).Op(":=").Id("argIsNotPointer").Call(
					jen.Op("&").Id("testingType").Values(),
				),
				assertFalse(
					jen.Id("t"),
					jen.Id("notAPointer"),
					jen.Lit("expected `false` when a pointer is provided"),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("error should not be returned when a pointer is provided"),
				),
			),
			buildSubTest(
				"with non-pointer",
				jen.List(
					jen.Id("notAPointer"),
					jen.Id("err"),
				).Op(":=").Id("argIsNotPointer").Call(
					jen.Id("testingType").Values(),
				),
				assertTrue(
					jen.Id("t"),
					jen.Id("notAPointer"),
					jen.Lit("expected `true` when a non-pointer is provided"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("error should be returned when a non-pointer is provided"),
				),
			),
			buildSubTest(
				"with nil",
				jen.List(
					jen.Id("notAPointer"),
					jen.Id("err"),
				).Op(":=").Id("argIsNotPointer").Call(
					jen.Id("nil"),
				),
				assertTrue(
					jen.Id("t"),
					jen.Id("notAPointer"),
					jen.Lit("expected `true` when nil is provided"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("error should be returned when nil is provided"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("ArgIsNotNil").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"without nil",
				jen.List(
					jen.Id("isNil"),
					jen.Id("err"),
				).Op(":=").Id("argIsNotNil").Call(
					jen.Op("&").Id("testingType").Values(),
				),
				assertFalse(
					jen.Id("t"),
					jen.Id("isNil"),
					jen.Lit("expected `false` when a pointer is provided"),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("error should not be returned when a pointer is provided"),
				),
			),
			buildSubTest(
				"with non-pointer",
				jen.List(
					jen.Id("isNil"),
					jen.Id("err"),
				).Op(":=").Id("argIsNotNil").Call(
					jen.Id("testingType").Values(),
				),
				assertFalse(
					jen.Id("t"),
					jen.Id("isNil"),
					jen.Lit("expected `true` when a non-pointer is provided"),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("error should not be returned when a non-pointer is provided"),
				),
			),
			buildSubTest(
				"with nil",
				jen.List(
					jen.Id("isNil"),
					jen.Id("err"),
				).Op(":=").Id("argIsNotNil").Call(
					jen.Id("nil"),
				),
				assertTrue(
					jen.Id("t"),
					jen.Id("isNil"),
					jen.Lit("expected `true` when nil is provided"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("error should be returned when nil is provided"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("UnmarshalBody").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"expected use",
				jen.Id("expected").Op(":=").Lit("example"),
				jen.Id("res").Op(":=").Op("&").Qual("net/http", "Response").Values(
					jen.Id("Body").Op(":").Qual("io/ioutil", "NopCloser").Call(
						jen.Qual("strings", "NewReader").Call(
							jen.Qual("fmt", "Sprintf").Call(
								jen.Lit(`{"name": %q}`),
								jen.Id("expected"),
							),
						),
					),
					jen.Id("StatusCode").Op(":").Qual("net/http", "StatusOK"),
				),
				jen.Var().Id("out").Id("testingType"),
				jen.Id("err").Op(":=").Id("unmarshalBody").Call(
					jen.Id("res"),
					jen.Op("&").Id("out"),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("out").Dot("Name"),
					jen.Id("expected"),
					jen.Lit("expected marshaling to work"),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be encountered unmarshaling into a valid struct"),
				),
			),
			buildSubTest(
				"with good status but unmarshallable response",
				jen.Id("res").Op(":=").Op("&").Qual("net/http", "Response").Values(
					jen.Id("Body").Op(":").Qual("io/ioutil", "NopCloser").Call(
						jen.Qual("strings", "NewReader").Call(
							jen.Lit(`

					BLAH

							`),
						),
					),
					jen.Id("StatusCode").Op(":").Qual("net/http", "StatusOK"),
				),
				jen.Var().Id("out").Id("testingType"),
				jen.Id("err").Op(":=").Id("unmarshalBody").Call(
					jen.Id("res"),
					jen.Op("&").Id("out"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be encountered unmarshaling into a valid struct"),
				),
			),
			buildSubTest(
				"with an erroneous error code",
				jen.Id("res").Op(":=").Op("&").Qual("net/http", "Response").Values(
					jen.Id("Body").Op(":").Qual("io/ioutil", "NopCloser").Call(
						jen.Qual("strings", "NewReader").Call(
							jen.Func().Params().Params(
								jen.Id("string"),
							).Block(
								jen.Id("er").Op(":=").Op("&").Id("models").Dot("ErrorResponse").Values(),
								jen.List(
									jen.Id("bs"),
									jen.Id("err"),
								).Op(":=").Qual("encoding/json", "Marshal").Call(
									jen.Id("er"),
								),
								requireNoError(jen.Id("err"), nil),
								jen.Return().Id("string").Call(
									jen.Id("bs"),
								),
							).Call(),
						),
					),
					jen.Id("StatusCode").Op(":").Qual("net/http", "StatusBadRequest"),
				),
				jen.Var().Id("out").Op("*").Id("testingType"),
				jen.Id("err").Op(":=").Id("unmarshalBody").Call(
					jen.Id("res"),
					jen.Op("&").Id("out"),
				),
				assertNil(
					jen.Id("t"),
					jen.Id("out"),
					jen.Lit("expected nil to be returned"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("error should be returned from the API"),
				),
			),
			buildSubTest(
				"with an erroneous error code and unmarshallable body",
				jen.Id("res").Op(":=").Op("&").Qual("net/http", "Response").Values(
					jen.Id("Body").Op(":").Qual("io/ioutil", "NopCloser").Call(
						jen.Qual("strings", "NewReader").Call(
							jen.Lit(`

				BLAH

							`),
						),
					),
					jen.Id("StatusCode").Op(":").Qual("net/http", "StatusBadRequest"),
				),
				jen.Var().Id("out").Op("*").Id("testingType"),
				jen.Id("err").Op(":=").Id("unmarshalBody").Call(
					jen.Id("res"),
					jen.Op("&").Id("out"),
				),
				assertNil(
					jen.Id("t"),
					jen.Id("out"),
					jen.Lit("expected nil to be returned"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("error should be returned from the unmarshaller"),
				),
			),
			buildSubTest(
				"with nil target variable",
				jen.Id("err").Op(":=").Id("unmarshalBody").Call(
					jen.Id("nil"),
					jen.Id("nil"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("error should be encountered when passed nil"),
				),
			),
			buildSubTest(
				"with erroneous reader",
				jen.Id("expected").Op(":=").Qual("errors", "New").Call(
					jen.Lit("blah"),
				),
				jen.Id("rc").Op(":=").Qual(
					"gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/mock",
					"NewMockReadCloser",
				).Call(),
				jen.Id("rc").Dot("On").Call(
					jen.Lit("Read"),
					jen.Qual(mockPkg, "Anything"),
				).Dot("Return").Call(
					jen.Lit(0),
					jen.Id("expected"),
				),
				jen.Id("res").Op(":=").Op("&").Qual("net/http", "Response").Values(
					jen.Id("Body").Op(":").Id("rc"),
					jen.Id("StatusCode").Op(":").Qual("net/http", "StatusOK"),
				),
				jen.Var().Id("out").Id("testingType"),
				jen.Id("err").Op(":=").Id("unmarshalBody").Call(
					jen.Id("res"), jen.Op("&").Id("out"),
				),
				assertEqual(jen.Id("expected"), jen.Id("err"), nil),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be encountered unmarshaling into a valid struct"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Id("testBreakableStruct").Struct(
			jen.Id("Thing").Qual("encoding/json", "Number").Tag(map[string]string{
				"json": `"thing"`,
			}),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Id("TestCreateBodyFromStruct").Params(jen.Id("T").Op("*").Qual("testing", "T")).Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"expected use",
				jen.Id("expected").Op(":=").Lit(`{"name":"expected"}`),
				jen.Id("x").Op(":=").Op("&").Id("testingType").Values(
					jen.Id("Name").Op(":").Lit("expected"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("createBodyFromStruct").Call(
					jen.Id("x"),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("expected no error creating JSON from valid struct"),
				),
				jen.List(
					jen.Id("bs"),
					jen.Id("err"),
				).Op(":=").Qual("io/ioutil", "ReadAll").Call(
					jen.Id("actual"),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("expected no error reading JSON from valid struct"),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("expected"),
					jen.Id("string").Call(
						jen.Id("bs"),
					),
					jen.Lit("expected and actual JSON bodies don't match"),
				),
			),
			buildSubTest(
				"with unmarshallable struct",
				jen.Id("x").Op(":=").Op("&").Id("testBreakableStruct").Values(
					jen.Id("Thing").Op(":").Lit("stuff"),
				),
				jen.List(
					jen.Id("_"),
					jen.Id("err"),
				).Op(":=").Id("createBodyFromStruct").Call(
					jen.Id("x"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("expected no error creating JSON from valid struct"),
				),
			),
		),
		jen.Line(),
	)

	return ret
}
