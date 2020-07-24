package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func clientTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("dbclient")

	utils.AddImports(proj, code)

	code.Add(
		jen.Const().Defs(
			jen.ID("defaultLimit").Equals().Uint8().Call(jen.Lit(20)),
		),
	)

	code.Add(
		jen.Func().ID("buildTestClient").Params().Params(jen.PointerTo().ID("Client"), jen.PointerTo().Qual(proj.DatabaseV1Package(), "MockDatabase")).Block(
			jen.ID("db").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
			jen.ID("c").Assign().AddressOf().ID("Client").Valuesln(
				jen.ID(constants.LoggerVarName).MapAssign().Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				jen.ID("querier").MapAssign().ID("db"),
			),
			jen.Return(jen.List(jen.ID("c"), jen.ID("db"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestMigrate").Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.Line(),
				jen.ID("c").Assign().AddressOf().ID("Client").Values(jen.ID("querier").MapAssign().ID("mockDB")),
				jen.ID("actual").Assign().ID("c").Dot("Migrate").Call(constants.CtxVar()),
				utils.AssertNoError(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"bubbles up errors",
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.ID("c").Assign().AddressOf().ID("Client").Values(jen.ID("querier").MapAssign().ID("mockDB")),
				jen.ID("actual").Assign().ID("c").Dot("Migrate").Call(constants.CtxVar()),
				utils.AssertError(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestIsReady").Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"obligatory",
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("IsReady"), jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(jen.True()),
				jen.Line(),
				jen.ID("c").Assign().AddressOf().ID("Client").Values(jen.ID("querier").MapAssign().ID("mockDB")),
				jen.ID("c").Dot("IsReady").Call(constants.CtxVar()),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideDatabaseClient").Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideDatabaseClient").Call(
					constants.CtxVar(),
					jen.Nil(),
					jen.ID("mockDB"),
					jen.True(),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				utils.AssertNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error migrating querier",
				jen.ID("expected").Assign().Qual("errors", "New").Call(jen.Lit("blah")),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.List(jen.ID("x"), jen.ID("actual")).Assign().ID("ProvideDatabaseClient").Call(
					constants.CtxVar(),
					jen.Nil(),
					jen.ID("mockDB"),
					jen.True(),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				utils.AssertNil(jen.ID("x"), nil),
				utils.AssertError(jen.ID("actual"), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)

	return code
}
