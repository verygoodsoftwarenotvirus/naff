package client

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func helpersTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("client")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Type().ID("testingType").Struct(
			jen.ID("Name").ID("string").Tag(map[string]string{"json": "name"}),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("ArgIsNotPointerOrNil").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"expected use",
				jen.Err().Op(":=").ID("argIsNotPointerOrNil").Call(
					jen.Op("&").ID("testingType").Values(),
				),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("error should not be returned when a pointer is provided"),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with non-pointer",
				jen.Err().Op(":=").ID("argIsNotPointerOrNil").Call(
					jen.ID("testingType").Values(),
				),
				utils.AssertError(
					jen.Err(),
					jen.Lit("error should be returned when a non-pointer is provided"),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with nil",
				jen.Err().Op(":=").ID("argIsNotPointerOrNil").Call(
					jen.Nil(),
				),
				utils.AssertError(
					jen.Err(),
					jen.Lit("error should be returned when nil is provided"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("ArgIsNotPointer").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"expected use",
				jen.List(
					jen.ID("notAPointer"),
					jen.Err(),
				).Op(":=").ID("argIsNotPointer").Call(
					jen.Op("&").ID("testingType").Values(),
				),
				utils.AssertFalse(
					jen.ID("notAPointer"),
					jen.Lit("expected `false` when a pointer is provided"),
				),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("error should not be returned when a pointer is provided"),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with non-pointer",
				jen.List(
					jen.ID("notAPointer"),
					jen.Err(),
				).Op(":=").ID("argIsNotPointer").Call(
					jen.ID("testingType").Values(),
				),
				utils.AssertTrue(
					jen.ID("notAPointer"),
					jen.Lit("expected `true` when a non-pointer is provided"),
				),
				utils.AssertError(
					jen.Err(),
					jen.Lit("error should be returned when a non-pointer is provided"),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with nil",
				jen.List(
					jen.ID("notAPointer"),
					jen.Err(),
				).Op(":=").ID("argIsNotPointer").Call(
					jen.Nil(),
				),
				utils.AssertTrue(
					jen.ID("notAPointer"),
					jen.Lit("expected `true` when nil is provided"),
				),
				utils.AssertError(
					jen.Err(),
					jen.Lit("error should be returned when nil is provided"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("ArgIsNotNil").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without nil",
				jen.List(
					jen.ID("isNil"),
					jen.Err(),
				).Op(":=").ID("argIsNotNil").Call(
					jen.Op("&").ID("testingType").Values(),
				),
				utils.AssertFalse(
					jen.ID("isNil"),
					jen.Lit("expected `false` when a pointer is provided"),
				),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("error should not be returned when a pointer is provided"),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with non-pointer",
				jen.List(
					jen.ID("isNil"),
					jen.Err(),
				).Op(":=").ID("argIsNotNil").Call(
					jen.ID("testingType").Values(),
				),
				utils.AssertFalse(
					jen.ID("isNil"),
					jen.Lit("expected `true` when a non-pointer is provided"),
				),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("error should not be returned when a non-pointer is provided"),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with nil",
				jen.List(
					jen.ID("isNil"),
					jen.Err(),
				).Op(":=").ID("argIsNotNil").Call(
					jen.Nil(),
				),
				utils.AssertTrue(
					jen.ID("isNil"),
					jen.Lit("expected `true` when nil is provided"),
				),
				utils.AssertError(
					jen.Err(),
					jen.Lit("error should be returned when nil is provided"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("UnmarshalBody").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"expected use",
				jen.ID("expected").Op(":=").Add(utils.FakeStringFunc()),
				jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").Valuesln(
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
				jen.Err().Op(":=").ID("unmarshalBody").Call(
					jen.ID("res"),
					jen.Op("&").ID("out"),
				),
				utils.AssertEqual(
					jen.ID("out").Dot("Name"),
					jen.ID("expected"),
					jen.Lit("expected marshaling to work"),
				),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be encountered unmarshaling into a valid struct"),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with good status but unmarshallable response",
				jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").Valuesln(
					jen.ID("Body").Op(":").Qual("io/ioutil", "NopCloser").Call(
						jen.Qual("strings", "NewReader").Call(jen.Lit("BLAH")),
					),
					jen.ID("StatusCode").Op(":").Qual("net/http", "StatusOK"),
				),
				jen.Var().ID("out").ID("testingType"),
				jen.Line(),
				jen.Err().Op(":=").ID("unmarshalBody").Call(
					jen.ID("res"),
					jen.Op("&").ID("out"),
				),
				utils.AssertError(
					jen.Err(),
					jen.Lit("no error should be encountered unmarshaling into a valid struct"),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with an erroneous error code",
				jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").Valuesln(
					jen.ID("Body").Op(":").Qual("io/ioutil", "NopCloser").Callln(
						jen.Qual("strings", "NewReader").Callln(
							jen.Func().Params().Params(
								jen.ID("string"),
							).Block(
								jen.ID("er").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "ErrorResponse").Values(),
								jen.List(
									jen.ID("bs"),
									jen.Err(),
								).Op(":=").Qual("encoding/json", "Marshal").Call(
									jen.ID("er"),
								),
								utils.RequireNoError(jen.Err(), nil),
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
				jen.Err().Op(":=").ID("unmarshalBody").Call(
					jen.ID("res"),
					jen.Op("&").ID("out"),
				),
				utils.AssertNil(
					jen.ID("out"),
					jen.Lit("expected nil to be returned"),
				),
				utils.AssertError(
					jen.Err(),
					jen.Lit("error should be returned from the API"),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with an erroneous error code and unmarshallable body",
				jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").Valuesln(
					jen.ID("Body").Op(":").Qual("io/ioutil", "NopCloser").Call(
						jen.Qual("strings", "NewReader").Call(jen.Lit("BLAH")),
					),
					jen.ID("StatusCode").Op(":").Qual("net/http", "StatusBadRequest"),
				),
				jen.Var().ID("out").Op("*").ID("testingType"),
				jen.Line(),
				jen.Err().Op(":=").ID("unmarshalBody").Call(
					jen.ID("res"),
					jen.Op("&").ID("out"),
				),
				utils.AssertNil(
					jen.ID("out"),
					jen.Lit("expected nil to be returned"),
				),
				utils.AssertError(
					jen.Err(),
					jen.Lit("error should be returned from the unmarshaller"),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with nil target variable",
				jen.Err().Op(":=").ID("unmarshalBody").Call(
					jen.Nil(),
					jen.Nil(),
				),
				utils.AssertError(
					jen.Err(),
					jen.Lit("error should be encountered when passed nil"),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with erroneous reader",
				jen.ID("expected").Op(":=").Qual("errors", "New").Call(
					jen.Lit("blah"),
				),
				jen.Line(),
				jen.ID("rc").Op(":=").Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil/mock"), "NewMockReadCloser").Call(),
				jen.ID("rc").Dot("On").Call(
					jen.Lit("Read"),
					jen.Qual(utils.MockPkg, "AnythingOfType").Call(jen.Lit("[]uint8")),
				).Dot("Return").Call(
					jen.Lit(0),
					jen.ID("expected"),
				),
				jen.Line(),
				jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").Valuesln(
					jen.ID("Body").Op(":").ID("rc"),
					jen.ID("StatusCode").Op(":").Qual("net/http", "StatusOK"),
				),
				jen.Var().ID("out").ID("testingType"),
				jen.Line(),
				jen.Err().Op(":=").ID("unmarshalBody").Call(
					jen.ID("res"), jen.Op("&").ID("out"),
				),
				utils.AssertEqual(jen.ID("expected"), jen.Err(), nil),
				utils.AssertError(
					jen.Err(),
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
		utils.OuterTestFunc("CreateBodyFromStruct").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"expected use",
				jen.ID("name").Op(":=").Add(utils.FakeStringFunc()),
				jen.ID("expected").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit(`{"name":%q}`), jen.ID("name")),
				jen.ID("x").Op(":=").Op("&").ID("testingType").Values(
					jen.ID("Name").Op(":").ID("name"),
				),
				jen.Line(),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Op(":=").ID("createBodyFromStruct").Call(
					jen.ID("x"),
				),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("expected no error creating JSON from valid struct"),
				),
				jen.Line(),
				jen.List(
					jen.ID("bs"),
					jen.Err(),
				).Op(":=").Qual("io/ioutil", "ReadAll").Call(
					jen.ID("actual"),
				),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("expected no error reading JSON from valid struct"),
				),
				utils.AssertEqual(
					jen.ID("expected"),
					jen.ID("string").Call(
						jen.ID("bs"),
					),
					jen.Lit("expected and actual JSON bodies don't match"),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with unmarshallable struct",
				jen.ID("x").Op(":=").Op("&").ID("testBreakableStruct").Values(
					jen.ID("Thing").Op(":").Lit("stuff"),
				),
				jen.List(
					jen.ID("_"),
					jen.Err(),
				).Op(":=").ID("createBodyFromStruct").Call(
					jen.ID("x"),
				),
				jen.Line(),
				utils.AssertError(
					jen.Err(),
					jen.Lit("expected no error creating JSON from valid struct"),
				),
			),
		),
		jen.Line(),
	)

	return ret
}
