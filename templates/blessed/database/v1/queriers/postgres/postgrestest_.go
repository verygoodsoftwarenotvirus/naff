package postgres

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func postgresTestDotGo() *jen.File {
	ret := jen.NewFile("postgres")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("Postgres"), jen.Qual("github.com/DATA-DOG/go-sqlmock",
			"Sqlmock",
		)).Block(
			jen.List(jen.ID("db"), jen.ID("mock"), jen.ID("err")).Op(":=").ID("sqlmock").Dot("New").Call(),
			jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("p").Op(":=").ID("ProvidePostgres").Call(jen.ID("true"), jen.ID("db"), jen.ID("noop").Dot("ProvideNoopLogger").Call()),
			jen.Return().List(jen.ID("p").Assert(jen.Op("*").ID("Postgres")), jen.ID("mock")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("sqlMockReplacer").Op("=").Qual("strings", "NewReplacer").PairedCallln(
				jen.Lit("$"), jen.RawString(`\$`),
				jen.Lit("("), jen.RawString(`\(`),
				jen.Lit(")"), jen.RawString(`\)`),
				jen.Lit("="), jen.RawString(`\=`),
				jen.Lit("*"), jen.RawString(`\*`),
				jen.Lit("."), jen.RawString(`\.`),
				jen.Lit("+"), jen.RawString(`\+`),
				jen.Lit("?"), jen.RawString(`\?`),
				jen.Lit(","), jen.RawString(`\,`),
				jen.Lit("-"), jen.RawString(`\-`),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("formatQueryForSQLMock").Params(jen.ID("query").ID("string")).Params(jen.ID("string")).Block(
			jen.Return().ID("sqlMockReplacer").Dot("Replace").Call(jen.ID("query")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvidePostgres").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("buildTestService").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestPostgres_IsReady").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("assert").Dot("True").Call(jen.ID("t"), jen.ID("p").Dot("IsReady").Call(jen.Qual("context", "Background").Call())),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_logQueryBuildingError").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("p").Dot("logQueryBuildingError").Call(jen.Qual("errors", "New").Call(jen.Lit(""))),
			)),
		),
		jen.Line(),
	)
	return ret
}