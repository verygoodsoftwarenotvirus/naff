package queriers

import (
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func databaseTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	ret := jen.NewFilePathName(proj.DatabaseV1Package("queriers", "v1", spn), spn)

	utils.AddImports(proj, ret)

	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	ret.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Params(jen.PointerTo().ID(sn), jen.Qual("github.com/DATA-DOG/go-sqlmock", "Sqlmock")).Block(
			jen.List(jen.ID("db"), jen.ID("mock"), jen.Err()).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "New").Call(),
			utils.RequireNoError(jen.Err(), nil),
			jen.ID(dbfl).Assign().IDf("Provide%s", sn).Call(jen.True(), jen.ID("db"), jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Return().List(jen.ID(dbfl).Assert(jen.PointerTo().ID(sn)), jen.ID("mock")),
		),
		jen.Line(),
	)

	regexPattern := `\?+`
	if isPostgres(dbvendor) {
		regexPattern = `\$\d+`
	}

	ret.Add(
		jen.Var().Defs(
			jen.ID("sqlMockReplacer").Equals().Qual("strings", "NewReplacer").PairedCallln(
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
			jen.ID("queryArgRegexp").Equals().Qual("regexp", "MustCompile").Call(jen.RawString(regexPattern)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("formatQueryForSQLMock").Params(jen.ID("query").String()).Params(jen.String()).Block(
			jen.Return().ID("sqlMockReplacer").Dot("Replace").Call(jen.ID("query")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("ensureArgCountMatchesQuery").Params(
			jen.ID("t").PointerTo().Qual("testing", "T"),
			jen.ID("query").String(),
			jen.ID("args").Index().Interface(),
		).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("queryArgCount").Assign().Len(jen.ID("queryArgRegexp").Dot("FindAllString").Call(jen.ID("query"), jen.Minus().One())),
			jen.Line(),
			jen.If(jen.Len(jen.ID("args")).GreaterThan().Zero()).Block(
				utils.AssertEqual(jen.ID("queryArgCount"), jen.Len(jen.ID("args")), nil),
			).Else().Block(
				utils.AssertZero(jen.ID("queryArgCount"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestProvide%s", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("buildTestService").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_IsReady", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"obligatory",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.AssertTrue(jen.ID(dbfl).Dot("IsReady").Call(utils.CtxVar()), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_logQueryBuildingError", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Qual("errors", "New").Call(jen.EmptyString())),
			),
		),
		jen.Line(),
	)

	if isMariaDB(dbvendor) || isSqlite(dbvendor) {
		ret.Add(
			jen.Func().IDf("Test%s_logIDRetrievalError", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				jen.ID("T").Dot("Parallel").Call(),
				jen.Line(),
				utils.BuildSubTestWithoutContext(
					"obligatory",
					jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID(dbfl).Dot("logIDRetrievalError").Call(jen.Qual("errors", "New").Call(jen.EmptyString())),
				),
			),
			jen.Line(),
		)
	}

	return ret
}
