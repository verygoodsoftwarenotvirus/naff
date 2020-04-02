package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func helpersTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)

	utils.AddImports(proj, ret)

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
				jen.Err().Assign().ID("argIsNotPointerOrNil").Call(
					jen.VarPointer().ID("testingType").Values(),
				),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("error should not be returned when a pointer is provided"),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with non-pointer",
				jen.Err().Assign().ID("argIsNotPointerOrNil").Call(
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
				jen.Err().Assign().ID("argIsNotPointerOrNil").Call(
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
				).Assign().ID("argIsNotPointer").Call(
					jen.VarPointer().ID("testingType").Values(),
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
				).Assign().ID("argIsNotPointer").Call(
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
				).Assign().ID("argIsNotPointer").Call(
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
				).Assign().ID("argIsNotNil").Call(
					jen.VarPointer().ID("testingType").Values(),
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
				).Assign().ID("argIsNotNil").Call(
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
				).Assign().ID("argIsNotNil").Call(
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
			utils.BuildSubTest(
				"expected use",
				jen.ID("expected").Assign().Lit("whatever"),
				jen.ID("res").Assign().VarPointer().Qual("net/http", "Response").Valuesln(
					jen.ID("Body").MapAssign().Qual("io/ioutil", "NopCloser").Call(
						jen.Qual("strings", "NewReader").Call(
							jen.Qual("fmt", "Sprintf").Call(
								jen.Lit("{\"name\": %q}"),
								jen.ID("expected"),
							),
						),
					),
					jen.ID("StatusCode").MapAssign().Qual("net/http", "StatusOK"),
				),
				jen.Var().ID("out").ID("testingType"),
				jen.Line(),
				jen.Err().Assign().ID("unmarshalBody").Call(
					utils.CtxVar(),
					jen.ID("res"),
					jen.VarPointer().ID("out"),
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
			utils.BuildSubTest(
				"with good status but unmarshallable response",
				jen.ID("res").Assign().VarPointer().Qual("net/http", "Response").Valuesln(
					jen.ID("Body").MapAssign().Qual("io/ioutil", "NopCloser").Call(
						jen.Qual("strings", "NewReader").Call(jen.Lit("BLAH")),
					),
					jen.ID("StatusCode").MapAssign().Qual("net/http", "StatusOK"),
				),
				jen.Var().ID("out").ID("testingType"),
				jen.Line(),
				jen.Err().Assign().ID("unmarshalBody").Call(
					utils.CtxVar(),
					jen.ID("res"),
					jen.VarPointer().ID("out"),
				),
				utils.AssertError(
					jen.Err(),
					jen.Lit("error should be encountered unmarshaling invalid response into a valid struct"),
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with an erroneous error code",
				jen.ID("res").Assign().VarPointer().Qual("net/http", "Response").Valuesln(
					jen.ID("Body").MapAssign().Qual("io/ioutil", "NopCloser").Callln(
						jen.Qual("strings", "NewReader").Callln(
							jen.Func().Params().Params(
								jen.ID("string"),
							).Block(
								jen.List(
									jen.ID("bs"),
									jen.Err(),
								).Assign().Qual("encoding/json", "Marshal").Call(
									jen.VarPointer().Qual(proj.ModelsV1Package(), "ErrorResponse").Values(),
								),
								utils.RequireNoError(jen.Err(), nil),
								jen.Return().ID("string").Call(
									jen.ID("bs"),
								),
							).Call(),
						),
					),
					jen.ID("StatusCode").MapAssign().Qual("net/http", "StatusBadRequest"),
				),
				jen.Var().ID("out").Op("*").ID("testingType"),
				jen.Line(),
				jen.Err().Assign().ID("unmarshalBody").Call(
					utils.CtxVar(),
					jen.ID("res"),
					jen.VarPointer().ID("out"),
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
			utils.BuildSubTest(
				"with an erroneous error code and unmarshallable body",
				jen.ID("res").Assign().VarPointer().Qual("net/http", "Response").Valuesln(
					jen.ID("Body").MapAssign().Qual("io/ioutil", "NopCloser").Call(
						jen.Qual("strings", "NewReader").Call(jen.Lit("BLAH")),
					),
					jen.ID("StatusCode").MapAssign().Qual("net/http", "StatusBadRequest"),
				),
				jen.Var().ID("out").Op("*").ID("testingType"),
				jen.Line(),
				jen.Err().Assign().ID("unmarshalBody").Call(
					utils.CtxVar(),
					jen.ID("res"),
					jen.VarPointer().ID("out"),
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
			utils.BuildSubTest(
				"with nil target variable",
				jen.Err().Assign().ID("unmarshalBody").Call(
					utils.CtxVar(),
					jen.Nil(),
					jen.Nil(),
				),
				utils.AssertError(
					jen.Err(),
					jen.Lit("error should be encountered when passed nil"),
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with erroneous reader",
				jen.ID("expected").Assign().Qual("errors", "New").Call(
					jen.Lit("blah"),
				),
				jen.Line(),
				jen.ID("rc").Assign().ID("newMockReadCloser").Call(),
				jen.ID("rc").Dot("On").Call(
					jen.Lit("Read"),
					jen.Qual(utils.MockPkg, "AnythingOfType").Call(jen.Lit("[]uint8")),
				).Dot("Return").Call(
					jen.Lit(0),
					jen.ID("expected"),
				),
				jen.Line(),
				jen.ID("res").Assign().VarPointer().Qual("net/http", "Response").Valuesln(
					jen.ID("Body").MapAssign().ID("rc"),
					jen.ID("StatusCode").MapAssign().Qual("net/http", "StatusOK"),
				),
				jen.Var().ID("out").ID("testingType"),
				jen.Line(),
				jen.Err().Assign().ID("unmarshalBody").Call(
					utils.CtxVar(),
					jen.ID("res"),
					jen.VarPointer().ID("out"),
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
				jen.ID("name").Assign().Lit("whatever"),
				jen.ID("expected").Assign().Qual("fmt", "Sprintf").Call(jen.Lit(`{"name":%q}`), jen.ID("name")),
				jen.ID("x").Assign().VarPointer().ID("testingType").Values(
					jen.ID("Name").MapAssign().ID("name"),
				),
				jen.Line(),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("createBodyFromStruct").Call(
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
				).Assign().Qual("io/ioutil", "ReadAll").Call(
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
				jen.ID("x").Assign().VarPointer().ID("testBreakableStruct").Values(
					jen.ID("Thing").MapAssign().Lit("stuff"),
				),
				jen.List(
					jen.ID("_"),
					jen.Err(),
				).Assign().ID("createBodyFromStruct").Call(
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
