package queriers

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"log"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func databaseTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	code := jen.NewFilePathName(proj.DatabaseV1Package("queriers", "v1", spn), spn)

	utils.AddImports(proj, code)

	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	code.Add(
		jen.Const().Defs(
			jen.ID("defaultLimit").Equals().Uint8().Call(jen.Lit(20)),
		),
	)

	code.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").PointerTo().Qual("testprojects", "T")).Params(jen.PointerTo().ID(sn), jen.Qual("github.com/DATA-DOG/go-sqlmock", "Sqlmock")).Block(
			jen.List(jen.ID("db"), jen.ID("mock"), jen.Err()).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "New").Call(),
			utils.RequireNoError(jen.Err(), nil),
			jen.ID(dbfl).Assign().IDf("Provide%s", sn).Call(jen.True(), jen.ID("db"), jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Return().List(jen.ID(dbfl).Assert(jen.PointerTo().ID(sn)), jen.ID("mock")),
		),
		jen.Line(),
	)

	regexPattern := `\?+`
	if isPostgres(dbvendor) {
		regexPattern = `\$\d+`
	}

	code.Add(
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
				jen.Lit("["), jen.RawString(`\[`),
				jen.Lit("]"), jen.RawString(`\]`),
			),
			jen.ID("queryArgRegexp").Equals().Qual("regexp", "MustCompile").Call(jen.RawString(regexPattern)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("formatQueryForSQLMock").Params(jen.ID("query").String()).Params(jen.String()).Block(
			jen.Return().ID("sqlMockReplacer").Dot("Replace").Call(jen.ID("query")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("ensureArgCountMatchesQuery").Params(
			jen.ID("t").PointerTo().Qual("testprojects", "T"),
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

	if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		code.Add(
			jen.Func().ID("interfacesToDriverValues").Params(jen.ID("in").Index().Interface()).Params(jen.ID("out").Index().Qual("database/sql/driver", "Value")).Block(
				jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().ID("in")).Block(
					jen.ID("out").Equals().Append(jen.ID("out"), jen.Qual("database/sql/driver", "Value").Call(jen.ID("x"))),
				),
				jen.Return(jen.ID("out")),
			),
		)
	}

	code.Add(
		jen.Func().IDf("TestProvide%s", sn).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("buildTestService").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().IDf("Test%s_IsReady", sn).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"obligatory",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.AssertTrue(jen.ID(dbfl).Dot("IsReady").Call(constants.CtxVar()), nil),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().IDf("Test%s_logQueryBuildingError", sn).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(constants.ObligatoryError()),
			),
		),
		jen.Line(),
	)

	if isMariaDB(dbvendor) || isSqlite(dbvendor) {
		code.Add(
			jen.Func().IDf("Test%s_logIDRetrievalError", sn).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
				jen.ID("T").Dot("Parallel").Call(),
				jen.Line(),
				utils.BuildSubTestWithoutContext(
					"obligatory",
					jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID(dbfl).Dot("logIDRetrievalError").Call(constants.ObligatoryError()),
				),
			),
			jen.Line(),
		)
	}

	if isPostgres(dbvendor) {
		code.Add(
			jen.Func().ID("Test_joinUint64s").Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
				jen.ID("T").Dot("Parallel").Call(),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(
					jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testprojects", "T")).Block(
						jen.ID(utils.BuildFakeVarName("Input")).Assign().Index().Uint64().Values(
							jen.Lit(123),
							jen.Lit(456),
							jen.Lit(789),
						),
						jen.Line(),
						jen.ID("expected").Assign().Lit("123,456,789"),
						jen.ID("actual").Assign().ID("joinUint64s").Call(jen.ID(utils.BuildFakeVarName("Input"))),
						jen.Line(),
						utils.AssertEqual(
							jen.ID("expected"),
							jen.ID("actual"),
							jen.Lit("expected %s to equal %s"),
							jen.ID("expected"),
							jen.ID("actual"),
						),
					),
				),
			),
		)
	}

	var (
		providerFuncName string
	)
	if isPostgres(dbvendor) || isSqlite(dbvendor) {
		providerFuncName = fmt.Sprintf("Provide%sDB", sn)
	} else if isMariaDB(dbvendor) {
		providerFuncName = fmt.Sprintf("Provide%sConnection", sn)
	} else {
		log.Panicf("invalid dbvendor: %q", dbvendor.Singular())
	}

	code.Add(
		jen.Func().IDf("Test%s", providerFuncName).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testprojects", "T")).Block(
					jen.List(jen.Underscore(), jen.Err()).Assign().ID(providerFuncName).Call(
						jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
						jen.EmptyString(),
					),
					utils.AssertNoError(jen.Err(), nil),
				),
			),
		),
	)

	return code
}
