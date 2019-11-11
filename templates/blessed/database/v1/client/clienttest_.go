package client

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func clientTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Func().ID("buildTestClient").Params().Params(jen.Op("*").ID("Client"), jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "MockDatabase")).Block(
			jen.ID("db").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
			jen.ID("c").Op(":=").Op("&").ID("Client").Valuesln(
				jen.ID("logger").Op(":").ID("noop").Dot("ProvideNoopLogger").Call(),
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
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("nil")),
				jen.Line(),
				jen.ID("c").Op(":=").Op("&").ID("Client").Values(jen.ID("querier").Op(":").ID("mockDB")),
				jen.ID("actual").Op(":=").ID("c").Dot("Migrate").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("bubbles up errors"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.ID("c").Op(":=").Op("&").ID("Client").Values(jen.ID("querier").Op(":").ID("mockDB")),
				jen.ID("actual").Op(":=").ID("c").Dot("Migrate").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestIsReady").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("IsReady"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("true")),
				jen.Line(),
				jen.ID("c").Op(":=").Op("&").ID("Client").Values(jen.ID("querier").Op(":").ID("mockDB")),
				jen.ID("c").Dot(
					"IsReady",
				).Call(jen.Qual("context", "Background").Call()),
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
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideDatabaseClient").Call(
					jen.Qual("context", "Background").Call(),
					jen.ID("nil"),
					jen.ID("mockDB"),
					jen.ID("false"),
					jen.ID("noop").Dot("ProvideNoopLogger").Call(),
				),
				jen.ID("assert").Dot("NotNil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error migrating querier"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").Qual("errors", "New").Call(jen.Lit("blah")),
				jen.ID("mockDB").Op(":=").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("On").Call(jen.Lit("Migrate"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.List(jen.ID("x"), jen.ID("actual")).Op(":=").ID("ProvideDatabaseClient").Call(
					jen.Qual("context", "Background").Call(),
					jen.ID("nil"),
					jen.ID("mockDB"),
					jen.ID("false"),
					jen.ID("noop").Dot("ProvideNoopLogger").Call(),
				),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("x")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)
	return ret
}
