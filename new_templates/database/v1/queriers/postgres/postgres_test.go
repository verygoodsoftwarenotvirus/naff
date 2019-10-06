package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func postgresTestDotGo() *jen.File {
	ret := jen.NewFile("postgres")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Func().ID("buildTestService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("Postgres"), jen.ID("sqlmock").Dot(
		"Sqlmock",
	)).Block(
		jen.List(jen.ID("db"), jen.ID("mock"), jen.ID("err")).Op(":=").ID("sqlmock").Dot(
			"New",
		).Call(),
		jen.ID("require").Dot(
			"NoError",
		).Call(jen.ID("t"), jen.ID("err")),
		jen.ID("p").Op(":=").ID("ProvidePostgres").Call(jen.ID("true"), jen.ID("db"), jen.ID("noop").Dot(
			"ProvideNoopLogger",
		).Call()),
		jen.Return().List(jen.ID("p").Assert(jen.Op("*").ID("Postgres")), jen.ID("mock")),
	),
	)
	ret.Add(jen.Null().Var().ID("sqlMockReplacer").Op("=").Qual("strings", "NewReplacer").Call(jen.Lit("$"), jen.Lit(`\$`), jen.Lit("("), jen.Lit(`\(`), jen.Lit(")"), jen.Lit(`\)`), jen.Lit("="), jen.Lit(`\=`), jen.Lit("*"), jen.Lit(`\*`), jen.Lit("."), jen.Lit(`\.`), jen.Lit("+"), jen.Lit(`\+`), jen.Lit("?"), jen.Lit(`\?`), jen.Lit(","), jen.Lit(`\,`), jen.Lit("-"), jen.Lit(`\-`)),
	)
	ret.Add(jen.Func().ID("formatQueryForSQLMock").Params(jen.ID("query").ID("string")).Params(jen.ID("string")).Block(
		jen.Return().ID("sqlMockReplacer").Dot(
			"Replace",
		).Call(jen.ID("query")),
	),
	)
	ret.Add(jen.Func().ID("TestProvidePostgres").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("buildTestService").Call(jen.ID("t")),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_IsReady").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("assert").Dot(
				"True",
			).Call(jen.ID("t"), jen.ID("p").Dot(
				"IsReady",
			).Call(jen.Qual("context", "Background").Call())),
		)),
	),
	)
	ret.Add(jen.Func().ID("Test_logQueryBuildingError").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("p").Dot(
				"logQueryBuildingError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit(""))),
		)),
	),
	)
	return ret
}
