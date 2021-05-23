package querybuilding

import (
	"fmt"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func databaseTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	code := jen.NewFilePathName(proj.DatabasePackage("queriers", "v1", spn), spn)

	utils.AddImports(proj, code, false)

	code.Add(buildConstDecls()...)
	code.Add(buildBuildTestService(dbvendor)...)
	code.Add(buildDBVendorTestVarDecls(dbvendor)...)

	code.Add(
		jen.Func().ID("formatQueryForSQLMock").Params(jen.ID("query").String()).Params(jen.String()).Body(
			jen.Return().ID("sqlMockReplacer").Dot("Replace").Call(jen.ID("query")),
		),
		jen.Line(),
	)

	code.Add(buildEnsureArgCountMatchesQuery()...)

	if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		code.Add(buildInterfacesToDriverValues()...)
	}

	code.Add(buildDBVendorProviderTest(dbvendor)...)
	code.Add(buildTestDBVendor_IsReady(dbvendor)...)
	code.Add(buildTestDBVendor_logQueryBuildingError(dbvendor)...)

	if isMariaDB(dbvendor) || isSqlite(dbvendor) {
		code.Add(buildTestDBVendor_logIDRetrievalError(dbvendor)...)
	}

	if isPostgres(dbvendor) {
		code.Add(buildTest_joinUint64s())
	}

	code.Add(buildTestProviderFunc(dbvendor))

	return code
}

func buildConstDecls() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.ID("defaultLimit").Equals().Uint8().Call(jen.Lit(20)),
		),
	}

	return lines
}

func buildBuildTestService(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	lines := []jen.Code{
		jen.Func().ID("buildTestService").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Params(jen.PointerTo().ID(sn), jen.Qual("github.com/DATA-DOG/go-sqlmock", "Sqlmock")).Body(
			jen.List(jen.ID("db"), jen.ID("mock"), jen.Err()).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "New").Call(),
			utils.RequireNoError(jen.Err(), nil),
			jen.ID(dbfl).Assign().IDf("Provide%s", sn).Call(jen.True(), jen.ID("db"), jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Return().List(jen.ID(dbfl).Assert(jen.PointerTo().ID(sn)), jen.ID("mock")),
		),
		jen.Line(),
	}

	return lines
}

func buildDBVendorTestVarDecls(dbvendor wordsmith.SuperPalabra) []jen.Code {
	regexPattern := `\?+`
	if isPostgres(dbvendor) {
		regexPattern = `\$\d+`
	}

	lines := []jen.Code{
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
	}

	return lines
}

func buildFormatQueryForSQLMock(dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		//
	}

	return lines
}

func buildEnsureArgCountMatchesQuery() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("ensureArgCountMatchesQuery").Params(
			jen.ID("t").PointerTo().Qual("testing", "T"),
			jen.ID("query").String(),
			jen.ID("args").Index().Interface(),
		).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("queryArgCount").Assign().Len(jen.ID("queryArgRegexp").Dot("FindAllString").Call(jen.ID("query"), jen.Minus().One())),
			jen.Line(),
			jen.If(jen.Len(jen.ID("args")).GreaterThan().Zero()).Body(
				utils.AssertEqual(jen.ID("queryArgCount"), jen.Len(jen.ID("args")), nil),
			).Else().Body(
				utils.AssertZero(jen.ID("queryArgCount"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildInterfacesToDriverValues() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("interfacesToDriverValues").Params(jen.ID("in").Index().Interface()).Params(jen.ID("out").Index().Qual("database/sql/driver", "Value")).Body(
			jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().ID("in")).Body(
				jen.ID("out").Equals().Append(jen.ID("out"), jen.Qual("database/sql/driver", "Value").Call(jen.ID("x"))),
			),
			jen.Return(jen.ID("out")),
		),
	}

	return lines
}

func buildDBVendorProviderTest(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()

	lines := []jen.Code{
		jen.Func().IDf("TestProvide%s", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("buildTestService").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDBVendor_IsReady(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	lines := []jen.Code{
		jen.Func().IDf("Test%s_IsReady", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"obligatory",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.AssertTrue(jen.ID(dbfl).Dot("IsReady").Call(constants.CtxVar()), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDBVendor_logQueryBuildingError(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	lines := []jen.Code{
		jen.Func().IDf("Test%s_logQueryBuildingError", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(constants.ObligatoryError()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDBVendor_logIDRetrievalError(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	lines := []jen.Code{
		jen.Func().IDf("Test%s_logIDRetrievalError", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID(dbfl).Dot("logIDRetrievalError").Call(constants.ObligatoryError()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTest_joinUint64s() jen.Code {
	return jen.Func().ID("Test_joinUint64s").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
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
	)
}

func buildTestProviderFunc(dbvendor wordsmith.SuperPalabra) jen.Code {
	sn := dbvendor.Singular()

	var providerFuncName string
	if isPostgres(dbvendor) || isSqlite(dbvendor) {
		providerFuncName = fmt.Sprintf("Provide%sDB", sn)
	} else if isMariaDB(dbvendor) {
		providerFuncName = fmt.Sprintf("Provide%sConnection", sn)
	} else {
		panic(fmt.Sprintf("invalid dbvendor: %q", dbvendor.Singular()))
	}

	return jen.Func().IDf("Test%s", providerFuncName).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.List(jen.Underscore(), jen.Err()).Assign().ID(providerFuncName).Call(
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.EmptyString(),
				),
				utils.AssertNoError(jen.Err(), nil),
			),
		),
	)
}
