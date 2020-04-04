package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func clientTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("buildTestClient").Params().Params(jen.PointerTo().ID("Client"), jen.PointerTo().Qual(proj.DatabaseV1Package(), "MockDatabase")).Block(
			jen.ID("db").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
			jen.ID("c").Assign().VarPointer().ID("Client").Valuesln(
				jen.ID("logger").MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				jen.ID("querier").MapAssign().ID("db"),
			),
			jen.Return(jen.List(jen.ID("c"), jen.ID("db"))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMigrate").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.Line(),
				jen.ID("c").Assign().VarPointer().ID("Client").Values(jen.ID("querier").MapAssign().ID("mockDB")),
				jen.ID("actual").Assign().ID("c").Dot("Migrate").Call(utils.CtxVar()),
				utils.AssertNoError(jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"bubbles up errors",
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.ID("c").Assign().VarPointer().ID("Client").Values(jen.ID("querier").MapAssign().ID("mockDB")),
				jen.ID("actual").Assign().ID("c").Dot("Migrate").Call(utils.CtxVar()),
				utils.AssertError(jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestIsReady").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"obligatory",
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("IsReady"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.True()),
				jen.Line(),
				jen.ID("c").Assign().VarPointer().ID("Client").Values(jen.ID("querier").MapAssign().ID("mockDB")),
				jen.ID("c").Dot("IsReady").Call(utils.CtxVar()),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideDatabaseClient").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideDatabaseClient").Call(
					utils.CtxVar(),
					jen.Nil(),
					jen.ID("mockDB"),
					jen.False(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				utils.AssertNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error migrating querier",
				jen.ID("expected").Assign().Qual("errors", "New").Call(jen.Lit("blah")),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.List(jen.ID("x"), jen.ID("actual")).Assign().ID("ProvideDatabaseClient").Call(
					utils.CtxVar(),
					jen.Nil(),
					jen.ID("mockDB"),
					jen.False(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				utils.AssertNil(jen.ID("x"), nil),
				utils.AssertError(jen.ID("actual"), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	return ret
}
