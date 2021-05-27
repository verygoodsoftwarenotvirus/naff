package mock

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountDataManagerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dot("AccountDataManager").Op("=").Parens(jen.Op("*").ID("AccountDataManager")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("AccountDataManager").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AccountExists is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountDataManager")).ID("AccountExists").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("bool"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("userID"),
			),
			jen.Return().List(jen.ID("args").Dot("Bool").Call(jen.Lit(0)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAccount is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountDataManager")).ID("GetAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("types").Dot("Account"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("userID"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("Account")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllAccountsCount is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountDataManager")).ID("GetAllAccountsCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllAccounts is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountDataManager")).ID("GetAllAccounts").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("results").Chan().Index().Op("*").ID("types").Dot("Account"), jen.ID("bucketSize").ID("uint16")).Params(jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("results"),
				jen.ID("bucketSize"),
			),
			jen.Return().ID("args").Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAccounts is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountDataManager")).ID("GetAccounts").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.Op("*").ID("types").Dot("AccountList"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("filter"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("AccountList")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAccountsForAdmin is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountDataManager")).ID("GetAccountsForAdmin").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.Op("*").ID("types").Dot("AccountList"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("filter"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("AccountList")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateAccount is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountDataManager")).ID("CreateAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("AccountCreationInput"), jen.ID("createdByUser").ID("uint64")).Params(jen.Op("*").ID("types").Dot("Account"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("input"),
				jen.ID("createdByUser"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("Account")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateAccount is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountDataManager")).ID("UpdateAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("types").Dot("Account"), jen.ID("changedByUser").ID("uint64"), jen.ID("changes").Index().Op("*").ID("types").Dot("FieldChangeSummary")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("updated"),
				jen.ID("changedByUser"),
				jen.ID("changes"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveAccount is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountDataManager")).ID("ArchiveAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID"), jen.ID("archivedByUser")).ID("uint64")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("userID"),
				jen.ID("archivedByUser"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAuditLogEntriesForAccount is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountDataManager")).ID("GetAuditLogEntriesForAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64")).Params(jen.Index().Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().Op("*").ID("types").Dot("AuditLogEntry")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	return code
}
