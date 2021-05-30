package mariadb

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mariadbTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("defaultLimit").Op("=").ID("uint8").Call(jen.Lit(20)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("MariaDB"), jen.ID("sqlmock").Dot("Sqlmock")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.List(jen.ID("_"), jen.ID("mock"), jen.ID("err")).Op(":=").ID("sqlmock").Dot("New").Call(),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().List(jen.ID("ProvideMariaDB").Call(jen.ID("logging").Dot("NewNoopLogger").Call()), jen.ID("mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("assertArgCountMatchesQuery").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("queryArgCount").Op(":=").ID("len").Call(jen.Qual("regexp", "MustCompile").Call(jen.Lit(`\?+`)).Dot("FindAllString").Call(
				jen.ID("query"),
				jen.Op("-").Lit(1),
			)),
			jen.If(jen.ID("len").Call(jen.ID("args")).Op(">").Lit(0)).Body(
				jen.ID("assert").Dot("Equal").Call(
					jen.ID("t"),
					jen.ID("queryArgCount"),
					jen.ID("len").Call(jen.ID("args")),
				)).Else().Body(
				jen.ID("assert").Dot("Zero").Call(
					jen.ID("t"),
					jen.ID("queryArgCount"),
				)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideMariaDB").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("buildTestService").Call(jen.ID("t")),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestMariaDB_logQueryBuildingError").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.ID("ctx")),
					jen.ID("q").Dot("logQueryBuildingError").Call(
						jen.ID("span"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideMariaDBConnection").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("ProvideMariaDBConnection").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.Lit(""),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
