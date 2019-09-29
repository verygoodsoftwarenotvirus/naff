package client

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func helpersTestDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

	ret.Add(
		jen.Type().ID("testingType").Struct(
			jen.ID("Name").ID("string").Tag(map[string]string{"json": "name"}),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("ArgIsNotPointerOrNil").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTestWithoutContext(
				"expected use",
				jen.ID("err").Op(":=").ID("argIsNotPointerOrNil").Call(
					jen.Op("&").ID("testingType").Values(),
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("error should not be returned when a pointer is provided"),
				),
			),
			jen.Line(),
			buildSubTestWithoutContext(
				"with non-pointer",
				jen.ID("err").Op(":=").ID("argIsNotPointerOrNil").Call(
					jen.ID("testingType").Values(),
				),
				assertError(
					jen.ID("err"),
					jen.Lit("error should be returned when a non-pointer is provided"),
				),
			),
			jen.Line(),
			buildSubTestWithoutContext(
				"with nil",
				jen.ID("err").Op(":=").ID("argIsNotPointerOrNil").Call(
					jen.ID("nil"),
				),
				assertError(
					jen.ID("err"),
					jen.Lit("error should be returned when nil is provided"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("ArgIsNotPointer").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTestWithoutContext(
				"expected use",
				jen.List(
					jen.ID("notAPointer"),
					jen.ID("err"),
				).Op(":=").ID("argIsNotPointer").Call(
					jen.Op("&").ID("testingType").Values(),
				),
				jen.Line(),
				assertFalse(
					jen.ID("notAPointer"),
					jen.Lit("expected `false` when a pointer is provided"),
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("error should not be returned when a pointer is provided"),
				),
			),
			jen.Line(),
			buildSubTestWithoutContext(
				"with non-pointer",
				jen.List(
					jen.ID("notAPointer"),
					jen.ID("err"),
				).Op(":=").ID("argIsNotPointer").Call(
					jen.ID("testingType").Values(),
				),
				jen.Line(),
				assertTrue(
					jen.ID("notAPointer"),
					jen.Lit("expected `true` when a non-pointer is provided"),
				),
				assertError(
					jen.ID("err"),
					jen.Lit("error should be returned when a non-pointer is provided"),
				),
			),
			jen.Line(),
			buildSubTestWithoutContext(
				"with nil",
				jen.List(
					jen.ID("notAPointer"),
					jen.ID("err"),
				).Op(":=").ID("argIsNotPointer").Call(
					jen.ID("nil"),
				),
				jen.Line(),
				assertTrue(
					jen.ID("notAPointer"),
					jen.Lit("expected `true` when nil is provided"),
				),
				assertError(
					jen.ID("err"),
					jen.Lit("error should be returned when nil is provided"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("ArgIsNotNil").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTestWithoutContext(
				"without nil",
				jen.List(
					jen.ID("isNil"),
					jen.ID("err"),
				).Op(":=").ID("argIsNotNil").Call(
					jen.Op("&").ID("testingType").Values(),
				),
				jen.Line(),
				assertFalse(
					jen.ID("isNil"),
					jen.Lit("expected `false` when a pointer is provided"),
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("error should not be returned when a pointer is provided"),
				),
			),
			jen.Line(),
			buildSubTestWithoutContext(
				"with non-pointer",
				jen.List(
					jen.ID("isNil"),
					jen.ID("err"),
				).Op(":=").ID("argIsNotNil").Call(
					jen.ID("testingType").Values(),
				),
				jen.Line(),
				assertFalse(
					jen.ID("isNil"),
					jen.Lit("expected `true` when a non-pointer is provided"),
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("error should not be returned when a non-pointer is provided"),
				),
			),
			jen.Line(),
			buildSubTestWithoutContext(
				"with nil",
				jen.List(
					jen.ID("isNil"),
					jen.ID("err"),
				).Op(":=").ID("argIsNotNil").Call(
					jen.ID("nil"),
				),
				jen.Line(),
				assertTrue(
					jen.ID("isNil"),
					jen.Lit("expected `true` when nil is provided"),
				),
				assertError(
					jen.ID("err"),
					jen.Lit("error should be returned when nil is provided"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("UnmarshalBody").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTestWithoutContext(
				"expected use",
				jen.ID("expected").Op(":=").Lit("example"),
				jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").ValuesLn(
					jen.ID("Body").Op(":").Qual("io/ioutil", "NopCloser").Call(
						jen.Qual("strings", "NewReader").Call(
							jen.Qual("fmt", "Sprintf").Call(
								jen.Lit("{\"name\": %q}"),
								jen.ID("expected"),
							),
						),
					),
					jen.ID("StatusCode").Op(":").Qual("net/http", "StatusOK"),
				),
				jen.Var().ID("out").ID("testingType"),
				jen.Line(),
				jen.ID("err").Op(":=").ID("unmarshalBody").Call(
					jen.ID("res"),
					jen.Op("&").ID("out"),
				),
				assertEqual(
					jen.ID("out").Dot("Name"),
					jen.ID("expected"),
					jen.Lit("expected marshaling to work"),
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be encountered unmarshaling into a valid struct"),
				),
			),
			jen.Line(),
			buildSubTestWithoutContext(
				"with good status but unmarshallable response",
				jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").ValuesLn(
					jen.ID("Body").Op(":").Qual("io/ioutil", "NopCloser").Call(
						jen.Qual("strings", "NewReader").Call(jen.Lit("BLAH")),
					),
					jen.ID("StatusCode").Op(":").Qual("net/http", "StatusOK"),
				),
				jen.Var().ID("out").ID("testingType"),
				jen.Line(),
				jen.ID("err").Op(":=").ID("unmarshalBody").Call(
					jen.ID("res"),
					jen.Op("&").ID("out"),
				),
				assertError(
					jen.ID("err"),
					jen.Lit("no error should be encountered unmarshaling into a valid struct"),
				),
			),
			jen.Line(),
			buildSubTestWithoutContext(
				"with an erroneous error code",
				jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").ValuesLn(
					jen.ID("Body").Op(":").Qual("io/ioutil", "NopCloser").CallLn(
						jen.Qual("strings", "NewReader").CallLn(
							jen.Func().Params().Params(
								jen.ID("string"),
							).Block(
								jen.ID("er").Op(":=").Op("&").Qual(modelsPkg, "ErrorResponse").Values(),
								jen.List(
									jen.ID("bs"),
									jen.ID("err"),
								).Op(":=").Qual("encoding/json", "Marshal").Call(
									jen.ID("er"),
								),
								requireNoError(jen.ID("err"), nil),
								jen.Return().ID("string").Call(
									jen.ID("bs"),
								),
							).Call(),
						),
					),
					jen.ID("StatusCode").Op(":").Qual("net/http", "StatusBadRequest"),
				),
				jen.Var().ID("out").Op("*").ID("testingType"),
				jen.Line(),
				jen.ID("err").Op(":=").ID("unmarshalBody").Call(
					jen.ID("res"),
					jen.Op("&").ID("out"),
				),
				assertNil(
					jen.ID("out"),
					jen.Lit("expected nil to be returned"),
				),
				assertError(
					jen.ID("err"),
					jen.Lit("error should be returned from the API"),
				),
			),
			jen.Line(),
			buildSubTestWithoutContext(
				"with an erroneous error code and unmarshallable body",
				jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").ValuesLn(
					jen.ID("Body").Op(":").Qual("io/ioutil", "NopCloser").Call(
						jen.Qual("strings", "NewReader").Call(jen.Lit("BLAH")),
					),
					jen.ID("StatusCode").Op(":").Qual("net/http", "StatusBadRequest"),
				),
				jen.Var().ID("out").Op("*").ID("testingType"),
				jen.Line(),
				jen.ID("err").Op(":=").ID("unmarshalBody").Call(
					jen.ID("res"),
					jen.Op("&").ID("out"),
				),
				assertNil(
					jen.ID("out"),
					jen.Lit("expected nil to be returned"),
				),
				assertError(
					jen.ID("err"),
					jen.Lit("error should be returned from the unmarshaller"),
				),
			),
			jen.Line(),
			buildSubTestWithoutContext(
				"with nil target variable",
				jen.ID("err").Op(":=").ID("unmarshalBody").Call(
					jen.ID("nil"),
					jen.ID("nil"),
				),
				assertError(
					jen.ID("err"),
					jen.Lit("error should be encountered when passed nil"),
				),
			),
			jen.Line(),
			buildSubTestWithoutContext(
				"with erroneous reader",
				jen.ID("expected").Op(":=").Qual("errors", "New").Call(
					jen.Lit("blah"),
				),
				jen.Line(),
				jen.ID("rc").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/mock", "NewMockReadCloser").Call(),
				jen.ID("rc").Dot("On").Call(
					jen.Lit("Read"),
					jen.Qual(mockPkg, "Anything"),
				).Dot("Return").Call(
					jen.Lit(0),
					jen.ID("expected"),
				),
				jen.Line(),
				jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").ValuesLn(
					jen.ID("Body").Op(":").ID("rc"),
					jen.ID("StatusCode").Op(":").Qual("net/http", "StatusOK"),
				),
				jen.Var().ID("out").ID("testingType"),
				jen.Line(),
				jen.ID("err").Op(":=").ID("unmarshalBody").Call(
					jen.ID("res"), jen.Op("&").ID("out"),
				),
				assertEqual(jen.ID("expected"), jen.ID("err"), nil),
				assertError(
					jen.ID("err"),
					jen.Lit("no error should be encountered unmarshaling into a valid struct"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("testBreakableStruct").Struct(
			jen.ID("Thing").Qual("encoding/json", "Number").Tag(map[string]string{
				"json": "thing",
			}),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("CreateBodyFromStruct").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTestWithoutContext(
				"expected use",
				jen.ID("expected").Op(":=").Lit(`{"name":"expected"}`),
				jen.ID("x").Op(":=").Op("&").ID("testingType").Values(
					jen.ID("Name").Op(":").Lit("expected"),
				),
				jen.Line(),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("createBodyFromStruct").Call(
					jen.ID("x"),
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("expected no error creating JSON from valid struct"),
				),
				jen.Line(),
				jen.List(
					jen.ID("bs"),
					jen.ID("err"),
				).Op(":=").Qual("io/ioutil", "ReadAll").Call(
					jen.ID("actual"),
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("expected no error reading JSON from valid struct"),
				),
				assertEqual(
					jen.ID("expected"),
					jen.ID("string").Call(
						jen.ID("bs"),
					),
					jen.Lit("expected and actual JSON bodies don't match"),
				),
			),
			jen.Line(),
			buildSubTestWithoutContext(
				"with unmarshallable struct",
				jen.ID("x").Op(":=").Op("&").ID("testBreakableStruct").Values(
					jen.ID("Thing").Op(":").Lit("stuff"),
				),
				jen.List(
					jen.ID("_"),
					jen.ID("err"),
				).Op(":=").ID("createBodyFromStruct").Call(
					jen.ID("x"),
				),
				jen.Line(),
				assertError(
					jen.ID("err"),
					jen.Lit("expected no error creating JSON from valid struct"),
				),
			),
		),
		jen.Line(),
	)

	return ret
}
