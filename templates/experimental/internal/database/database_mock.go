package database

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func databaseMockDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("DataManager").Op("=").Parens(jen.Op("*").ID("MockDatabase")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildMockDatabase builds a mock database."),
		jen.Line(),
		jen.Func().ID("BuildMockDatabase").Params().Params(jen.Op("*").ID("MockDatabase")).Body(
			jen.Return().Op("&").ID("MockDatabase").Valuesln(jen.ID("AuditLogEntryDataManager").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Valuesln(), jen.ID("AccountDataManager").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountDataManager").Valuesln(), jen.ID("AccountUserMembershipDataManager").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(), jen.ID("ItemDataManager").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(), jen.ID("UserDataManager").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(), jen.ID("AdminUserDataManager").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AdminUserDataManager").Valuesln(), jen.ID("APIClientDataManager").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(), jen.ID("WebhookDataManager").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln())),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("MockDatabase").Struct(
				jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AdminUserDataManager"),
				jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager"),
				jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager"),
				jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager"),
				jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager"),
				jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager"),
				jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager"),
				jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountDataManager"),
				jen.ID("mock").Dot("Mock"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Migrate satisfies the DataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockDatabase")).ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("maxAttempts").ID("uint8"), jen.ID("ucc").Op("*").ID("types").Dot("TestUserCreationConfig")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("maxAttempts"),
				jen.ID("ucc"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("IsReady satisfies the DataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockDatabase")).ID("IsReady").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("maxAttempts").ID("uint8")).Params(jen.ID("ready").ID("bool")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("maxAttempts"),
			).Dot("Bool").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BeginTx satisfies the DataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockDatabase")).ID("BeginTx").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("options").Op("*").Qual("database/sql", "TxOptions")).Params(jen.Op("*").Qual("database/sql", "Tx"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("options"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual("database/sql", "Tx")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("querybuilding").Dot("SQLQueryBuilder").Op("=").Parens(jen.Op("*").ID("MockSQLQueryBuilder")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildMockSQLQueryBuilder builds a MockSQLQueryBuilder."),
		jen.Line(),
		jen.Func().ID("BuildMockSQLQueryBuilder").Params().Params(jen.Op("*").ID("MockSQLQueryBuilder")).Body(
			jen.Return().Op("&").ID("MockSQLQueryBuilder").Valuesln(jen.ID("AccountSQLQueryBuilder").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mock", "AccountSQLQueryBuilder").Valuesln(), jen.ID("AccountUserMembershipSQLQueryBuilder").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mock", "AccountUserMembershipSQLQueryBuilder").Valuesln(), jen.ID("AuditLogEntrySQLQueryBuilder").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mock", "AuditLogEntrySQLQueryBuilder").Valuesln(), jen.ID("ItemSQLQueryBuilder").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mock", "ItemSQLQueryBuilder").Valuesln(), jen.ID("APIClientSQLQueryBuilder").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mock", "APIClientSQLQueryBuilder").Valuesln(), jen.ID("UserSQLQueryBuilder").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mock", "UserSQLQueryBuilder").Valuesln(), jen.ID("WebhookSQLQueryBuilder").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mock", "WebhookSQLQueryBuilder").Valuesln())),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("MockSQLQueryBuilder").Struct(
				jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mock", "UserSQLQueryBuilder"),
				jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mock", "AccountSQLQueryBuilder"),
				jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mock", "AccountUserMembershipSQLQueryBuilder"),
				jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mock", "AuditLogEntrySQLQueryBuilder"),
				jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mock", "ItemSQLQueryBuilder"),
				jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mock", "APIClientSQLQueryBuilder"),
				jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mock", "WebhookSQLQueryBuilder"),
				jen.ID("mock").Dot("Mock"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildMigrationFunc implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockSQLQueryBuilder")).ID("BuildMigrationFunc").Params(jen.ID("db").Op("*").Qual("database/sql", "DB")).Params(jen.Params()).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("db")),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Params()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildTestUserCreationQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockSQLQueryBuilder")).ID("BuildTestUserCreationQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("testUserConfig").Op("*").ID("types").Dot("TestUserCreationConfig")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnValues").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("testUserConfig"),
			),
			jen.Return().List(jen.ID("returnValues").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("string")), jen.ID("returnValues").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("ResultIterator").Op("=").Parens(jen.Op("*").ID("MockResultIterator")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("MockResultIterator").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Scan satisfies the ResultIterator interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockResultIterator")).ID("Scan").Params(jen.ID("dest").Op("...").Interface()).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("dest").Op("...")).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Next satisfies the ResultIterator interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockResultIterator")).ID("Next").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("m").Dot("Called").Call().Dot("Bool").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Err satisfies the ResultIterator interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockResultIterator")).ID("Err").Params().Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call().Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Close satisfies the ResultIterator interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockResultIterator")).ID("Close").Params().Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call().Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("MockSQLResult").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LastInsertId implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockSQLResult")).ID("LastInsertId").Params().Params(jen.ID("int64"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("int64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("RowsAffected implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockSQLResult")).ID("RowsAffected").Params().Params(jen.ID("int64"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("int64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	return code
}
