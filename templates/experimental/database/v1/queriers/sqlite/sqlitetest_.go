package sqlite

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func sqliteTestDotGo() *jen.File {
	ret := jen.NewFile("sqlite")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("Sqlite"), jen.ID("sqlmock").Dot(
			"Sqlmock",
		)).Block(
			jen.List(jen.ID("db"), jen.ID("mock"), jen.ID("err")).Op(":=").ID("sqlmock").Dot("New").Call(),
			jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("s").Op(":=").ID("ProvideSqlite").Call(jen.ID("true"), jen.ID("db"), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.Return().List(jen.ID("s").Assert(jen.Op("*").ID("Sqlite")), jen.ID("mock")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("sqlMockReplacer").Op("=").Qual("strings", "NewReplacer").Call(jen.Lit("$"), jen.Lit(`\$`), jen.Lit("("), jen.Lit(`\(`), jen.Lit(")"), jen.Lit(`\)`), jen.Lit("="), jen.Lit(`\=`), jen.Lit("*"), jen.Lit(`\*`), jen.Lit("."), jen.Lit(`\.`), jen.Lit("+"), jen.Lit(`\+`), jen.Lit("?"), jen.Lit(`\?`), jen.Lit(","), jen.Lit(`\,`), jen.Lit("-"), jen.Lit(`\-`)),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("formatQueryForSQLMock").Params(jen.ID("query").ID("string")).Params(jen.ID("string")).Block(
			jen.Return().ID("sqlMockReplacer").Dot(
				"Replace",
			).Call(jen.ID("query")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideSqlite").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("buildTestService").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestSqlite_IsReady").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("s"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("assert").Dot(
					"True",
				).Call(jen.ID("t"), jen.ID("s").Dot(
					"IsReady",
				).Call(jen.Qual("context", "Background").Call())),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestSqlite_logQueryBuildingError").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("s"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("s").Dot(
					"logQueryBuildingError",
				).Call(jen.Qual("errors", "New").Call(jen.Lit(""))),
			)),
		),
		jen.Line(),
	)
	return ret
}
