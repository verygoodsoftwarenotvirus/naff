package client

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func clientTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(pkg, ret)

	ret.Add(
		utils.FakeSeedFunc(),
	)

	ret.Add(
		jen.Func().ID("buildTestClient").Params().Params(jen.Op("*").ID("Client"), jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "MockDatabase")).Block(
			jen.ID("db").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
			jen.ID("c").Op(":=").Op("&").ID("Client").Valuesln(
				jen.ID("logger").Op(":").Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				jen.ID("querier").Op(":").ID("db"),
			),
			jen.Return(jen.List(jen.ID("c"), jen.ID("db"))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMigrate").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Nil()),
				jen.Line(),
				jen.ID("c").Op(":=").Op("&").ID("Client").Values(jen.ID("querier").Op(":").ID("mockDB")),
				jen.ID("actual").Op(":=").ID("c").Dot("Migrate").Call(utils.CtxVar()),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("bubbles up errors"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.ID("c").Op(":=").Op("&").ID("Client").Values(jen.ID("querier").Op(":").ID("mockDB")),
				jen.ID("actual").Op(":=").ID("c").Dot("Migrate").Call(utils.CtxVar()),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestIsReady").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("IsReady"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.True()),
				jen.Line(),
				jen.ID("c").Op(":=").Op("&").ID("Client").Values(jen.ID("querier").Op(":").ID("mockDB")),
				jen.ID("c").Dot("IsReady").Call(utils.CtxVar()),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideDatabaseClient").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("ProvideDatabaseClient").Call(
					utils.CtxVar(),
					jen.Nil(),
					jen.ID("mockDB"),
					jen.False(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error migrating querier"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expected").Op(":=").Qual("errors", "New").Call(jen.Lit("blah")),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.List(jen.ID("x"), jen.ID("actual")).Op(":=").ID("ProvideDatabaseClient").Call(
					utils.CtxVar(),
					jen.Nil(),
					jen.ID("mockDB"),
					jen.False(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("x")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.ID("actual")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)
	return ret
}
