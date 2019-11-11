package queriers

import (
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func databaseTestDotGo(pkg *models.Project, vendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(vendor.SingularPackageName())

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)
	sn := vendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	ret.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID(sn), jen.Qual("github.com/DATA-DOG/go-sqlmock", "Sqlmock")).Block(
			jen.List(jen.ID("db"), jen.ID("mock"), jen.ID("err")).Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "New").Call(),
			jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID(dbfl).Op(":=").IDf("Provide%s", sn).Call(jen.ID("true"), jen.ID("db"), jen.ID("noop").Dot("ProvideNoopLogger").Call()),
			jen.Return().List(jen.ID(dbfl).Assert(jen.Op("*").ID(sn)), jen.ID("mock")),
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
		jen.Func().IDf("TestProvide%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("buildTestService").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_IsReady", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("assert").Dot("True").Call(jen.ID("t"), jen.ID(dbfl).Dot("IsReady").Call(jen.Qual("context", "Background").Call())),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_logQueryBuildingError", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Qual("errors", "New").Call(jen.Lit(""))),
			)),
		),
		jen.Line(),
	)
	return ret
}
