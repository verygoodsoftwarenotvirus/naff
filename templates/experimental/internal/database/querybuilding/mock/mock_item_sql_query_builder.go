package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockItemSQLQueryBuilderDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("querybuilding").Dot("ItemSQLQueryBuilder").Op("=").Parens(jen.Op("*").ID("ItemSQLQueryBuilder")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("ItemSQLQueryBuilder").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildItemExistsQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemSQLQueryBuilder")).ID("BuildItemExistsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
				jen.ID("accountID"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetItemQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemSQLQueryBuilder")).ID("BuildGetItemQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
				jen.ID("accountID"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAllItemsCountQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemSQLQueryBuilder")).ID("BuildGetAllItemsCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Return().ID("returnArgs").Dot("String").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetBatchOfItemsQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemSQLQueryBuilder")).ID("BuildGetBatchOfItemsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("beginID"),
				jen.ID("endID"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetItemsQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemSQLQueryBuilder")).ID("BuildGetItemsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("forAdmin").ID("bool"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("forAdmin"),
				jen.ID("filter"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetItemsWithIDsQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemSQLQueryBuilder")).ID("BuildGetItemsWithIDsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("limit").ID("uint8"), jen.ID("ids").Index().ID("uint64"), jen.ID("forAdmin").ID("bool")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("limit"),
				jen.ID("ids"),
				jen.ID("forAdmin"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildCreateItemQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemSQLQueryBuilder")).ID("BuildCreateItemQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("ItemCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAuditLogEntriesForItemQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemSQLQueryBuilder")).ID("BuildGetAuditLogEntriesForItemQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("itemID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUpdateItemQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemSQLQueryBuilder")).ID("BuildUpdateItemQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("Item")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildArchiveItemQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemSQLQueryBuilder")).ID("BuildArchiveItemQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
				jen.ID("accountID"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	return code
}
