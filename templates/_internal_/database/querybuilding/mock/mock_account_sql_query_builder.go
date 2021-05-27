package mock

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockAccountSQLQueryBuilderDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("querybuilding").Dot("AccountSQLQueryBuilder").Op("=").Parens(jen.Op("*").ID("AccountSQLQueryBuilder")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("AccountSQLQueryBuilder").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildTransferAccountOwnershipQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountSQLQueryBuilder")).ID("BuildTransferAccountOwnershipQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("currentOwnerID"), jen.ID("newOwnerID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("currentOwnerID"),
				jen.ID("newOwnerID"),
				jen.ID("accountID"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAccountQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountSQLQueryBuilder")).ID("BuildGetAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("userID"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAllAccountsCountQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountSQLQueryBuilder")).ID("BuildGetAllAccountsCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Return().ID("returnArgs").Dot("String").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetBatchOfAccountsQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountSQLQueryBuilder")).ID("BuildGetBatchOfAccountsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
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
		jen.Comment("BuildGetAccountsQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountSQLQueryBuilder")).ID("BuildGetAccountsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("forAdmin").ID("bool"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("forAdmin"),
				jen.ID("filter"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildAccountCreationQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountSQLQueryBuilder")).ID("BuildAccountCreationQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("AccountCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUpdateAccountQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountSQLQueryBuilder")).ID("BuildUpdateAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("Account")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildArchiveAccountQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountSQLQueryBuilder")).ID("BuildArchiveAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("userID"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAuditLogEntriesForAccountQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountSQLQueryBuilder")).ID("BuildGetAuditLogEntriesForAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	return code
}
