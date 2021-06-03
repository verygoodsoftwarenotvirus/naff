package mock

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockIterablesSQLQueryBuilderDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	pn := typ.Name.Plural()

	code.Add(
		jen.Var().ID("_").ID("querybuilding").Dotf("%sSQLQueryBuilder", sn).Op("=").Parens(jen.Op("*").IDf("%sSQLQueryBuilder", sn)).Call(jen.ID("nil")),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("%sSQLQueryBuilder is a mocked types.%sSQLQueryBuilder for testing.", sn, sn),
		jen.Newline(),
		jen.Type().IDf("%sSQLQueryBuilder", sn).Struct(jen.ID("mock").Dot("Mock")),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("Build%sExistsQuery implements our interface.", sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sSQLQueryBuilder", sn)).IDf("Build%sExistsQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.IDf("%sID", uvn), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildGet%sQuery implements our interface.", sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sSQLQueryBuilder", sn)).IDf("BuildGet%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.IDf("%sID", uvn), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildGetAll%sCountQuery implements our interface.", pn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sSQLQueryBuilder", sn)).IDf("BuildGetAll%sCountQuery", pn).Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Newline(),
			jen.Return().ID("returnArgs").Dot("String").Call(jen.Lit(0)),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildGetBatchOf%sQuery implements our interface.", pn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sSQLQueryBuilder", sn)).IDf("BuildGetBatchOf%sQuery", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("beginID"),
				jen.ID("endID"),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildGet%sQuery implements our interface.", pn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sSQLQueryBuilder", sn)).IDf("BuildGet%sQuery", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("forAdmin").ID("bool"), jen.ID("filter").Op("*").Qual(proj.TypesPackage(), "QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("forAdmin"),
				jen.ID("filter"),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildGet%sWithIDsQuery implements our interface.", pn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sSQLQueryBuilder", sn)).IDf("BuildGet%sWithIDsQuery", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("limit").ID("uint8"), jen.ID("ids").Index().ID("uint64"), jen.ID("forAdmin").ID("bool")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("limit"),
				jen.ID("ids"),
				jen.ID("forAdmin"),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildCreate%sQuery implements our interface.", sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sSQLQueryBuilder", sn)).IDf("BuildCreate%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn))).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildGetAuditLogEntriesFor%sQuery implements our interface.", sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sSQLQueryBuilder", sn)).IDf("BuildGetAuditLogEntriesFor%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.IDf("%sID", uvn).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildUpdate%sQuery implements our interface.", sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sSQLQueryBuilder", sn)).IDf("BuildUpdate%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual(proj.TypesPackage(), sn)).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildArchive%sQuery implements our interface.", sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sSQLQueryBuilder", sn)).IDf("BuildArchive%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.IDf("%sID", uvn), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Newline(),
	)

	return code
}
