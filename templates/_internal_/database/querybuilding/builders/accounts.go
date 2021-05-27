package builders

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("querybuilding").Dot("AccountSQLQueryBuilder").Op("=").Parens(jen.Op("*").ID("Postgres")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAccountQuery constructs a SQL query for fetching an account with a given ID belong to a user with a given ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildGetAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("columns").Op(":=").ID("append").Call(
				jen.ID("querybuilding").Dot("AccountsTableColumns"),
				jen.ID("querybuilding").Dot("AccountsUserMembershipTableColumns").Op("..."),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("columns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("AccountsTableName")).Dot("Join").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s ON %s.%s = %s.%s"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableAccountOwnershipColumn"),
					jen.ID("querybuilding").Dot("AccountsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				)).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("accountID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsTableName"),
					jen.ID("querybuilding").Dot("AccountsTableUserOwnershipColumn"),
				).Op(":").ID("userID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAllAccountsCountQuery returns a query that fetches the total number of accounts in the database."),
		jen.Line(),
		jen.Comment("This query only gets generated once, and is otherwise returned from cache."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildGetAllAccountsCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.List(jen.ID("allAccountsCountQuery"), jen.ID("_"), jen.ID("err")).Op(":=").ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
				jen.ID("columnCountQueryTemplate"),
				jen.ID("querybuilding").Dot("AccountsTableName"),
			)).Dot("From").Call(jen.ID("querybuilding").Dot("AccountsTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s"),
				jen.ID("querybuilding").Dot("AccountsTableName"),
				jen.ID("querybuilding").Dot("ArchivedOnColumn"),
			).Op(":").ID("nil"))).Dot("ToSql").Call(),
			jen.ID("b").Dot("logQueryBuildingError").Call(
				jen.ID("span"),
				jen.ID("err"),
			),
			jen.Return().ID("allAccountsCountQuery"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetBatchOfAccountsQuery returns a query that fetches every account in the database within a bucketed range."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildGetBatchOfAccountsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("AccountsTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("AccountsTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Gt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("beginID"))).Dot("Where").Call(jen.ID("squirrel").Dot("Lt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("endID"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAccountsQuery builds a SQL query selecting accounts that adhere to a given QueryFilter and belong to a given account,"),
		jen.Line(),
		jen.Comment("and returns both the query and the relevant args to pass to the query executor."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildGetAccountsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("forAdmin").ID("bool"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("tracing").Dot("AttachFilterToSpan").Call(
					jen.ID("span"),
					jen.ID("filter").Dot("Page"),
					jen.ID("filter").Dot("Limit"),
					jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
				)),
			jen.Var().Defs(
				jen.ID("includeArchived").ID("bool"),
			),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("includeArchived").Op("=").ID("filter").Dot("IncludeArchived")),
			jen.ID("columns").Op(":=").ID("append").Call(
				jen.ID("querybuilding").Dot("AccountsTableColumns"),
				jen.ID("querybuilding").Dot("AccountsUserMembershipTableColumns").Op("..."),
			),
			jen.List(jen.ID("filteredCountQuery"), jen.ID("filteredCountQueryArgs")).Op(":=").ID("b").Dot("buildFilteredCountQuery").Call(
				jen.ID("ctx"),
				jen.ID("querybuilding").Dot("AccountsTableName"),
				jen.ID("querybuilding").Dot("AccountsTableUserOwnershipColumn"),
				jen.ID("userID"),
				jen.ID("forAdmin"),
				jen.ID("includeArchived"),
				jen.ID("filter"),
			),
			jen.List(jen.ID("totalCountQuery"), jen.ID("totalCountQueryArgs")).Op(":=").ID("b").Dot("buildTotalCountQuery").Call(
				jen.ID("ctx"),
				jen.ID("querybuilding").Dot("AccountsTableName"),
				jen.ID("querybuilding").Dot("AccountsTableUserOwnershipColumn"),
				jen.ID("userID"),
				jen.ID("forAdmin"),
				jen.ID("includeArchived"),
			),
			jen.ID("builder").Op(":=").ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("append").Call(
				jen.ID("columns"),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("(%s) as total_count"),
					jen.ID("totalCountQuery"),
				),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("(%s) as filtered_count"),
					jen.ID("filteredCountQuery"),
				),
			).Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("AccountsTableName")).Dot("Join").Call(jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s ON %s.%s = %s.%s"),
				jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
				jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
				jen.ID("querybuilding").Dot("AccountsUserMembershipTableAccountOwnershipColumn"),
				jen.ID("querybuilding").Dot("AccountsTableName"),
				jen.ID("querybuilding").Dot("IDColumn"),
			)),
			jen.If(jen.Op("!").ID("forAdmin")).Body(
				jen.ID("builder").Op("=").ID("builder").Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsTableName"),
					jen.ID("querybuilding").Dot("AccountsTableUserOwnershipColumn"),
				).Op(":").ID("userID")))),
			jen.ID("builder").Op("=").ID("builder").Dot("GroupBy").Call(jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("(%s.%s, %s.%s)"),
				jen.ID("querybuilding").Dot("AccountsTableName"),
				jen.ID("querybuilding").Dot("IDColumn"),
				jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
				jen.ID("querybuilding").Dot("IDColumn"),
			)),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("builder").Op("=").ID("querybuilding").Dot("ApplyFilterToQueryBuilder").Call(
					jen.ID("filter"),
					jen.ID("querybuilding").Dot("AccountsTableName"),
					jen.ID("builder"),
				)),
			jen.List(jen.ID("query"), jen.ID("selectArgs")).Op(":=").ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("builder"),
			),
			jen.Return().List(jen.ID("query"), jen.ID("append").Call(
				jen.ID("append").Call(
					jen.ID("filteredCountQueryArgs"),
					jen.ID("totalCountQueryArgs").Op("..."),
				),
				jen.ID("selectArgs").Op("..."),
			)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildAccountCreationQuery takes an account and returns a creation query for that account and the relevant arguments."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildAccountCreationQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("AccountCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.ID("querybuilding").Dot("AccountsTableName")).Dot("Columns").Call(
					jen.ID("querybuilding").Dot("ExternalIDColumn"),
					jen.ID("querybuilding").Dot("AccountsTableNameColumn"),
					jen.ID("querybuilding").Dot("AccountsTableBillingStatusColumn"),
					jen.ID("querybuilding").Dot("AccountsTableContactEmailColumn"),
					jen.ID("querybuilding").Dot("AccountsTableContactPhoneColumn"),
					jen.ID("querybuilding").Dot("AccountsTableUserOwnershipColumn"),
				).Dot("Values").Call(
					jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
					jen.ID("input").Dot("Name"),
					jen.ID("types").Dot("UnpaidAccountBillingStatus"),
					jen.ID("input").Dot("ContactEmail"),
					jen.ID("input").Dot("ContactPhone"),
					jen.ID("input").Dot("BelongsToUser"),
				).Dot("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("RETURNING %s"),
					jen.ID("querybuilding").Dot("IDColumn"),
				)),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUpdateAccountQuery takes an account and returns an update SQL query, with the relevant query parameters."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildUpdateAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("Account")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("ID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("AccountsTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("AccountsTableNameColumn"),
					jen.ID("input").Dot("Name"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("AccountsTableContactEmailColumn"),
					jen.ID("input").Dot("ContactEmail"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("AccountsTableContactPhoneColumn"),
					jen.ID("input").Dot("ContactPhone"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("input").Dot("ID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"), jen.ID("querybuilding").Dot("AccountsTableUserOwnershipColumn").Op(":").ID("input").Dot("BelongsToUser"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildArchiveAccountQuery returns a SQL query which marks a given account belonging to a given user as archived."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildArchiveAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("AccountsTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("accountID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"), jen.ID("querybuilding").Dot("AccountsTableUserOwnershipColumn").Op(":").ID("userID"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAuditLogEntriesForAccountQuery constructs a SQL query for fetching audit log entries"),
		jen.Line(),
		jen.Comment("associated with a given account."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildGetAuditLogEntriesForAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("accountIDKey").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.ID("jsonPluckQuery"),
				jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
				jen.ID("querybuilding").Dot("AuditLogEntriesTableContextColumn"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "AccountAssignmentKey"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("accountIDKey").Op(":").ID("accountID"))).Dot("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
					jen.ID("querybuilding").Dot("CreatedOnColumn"),
				)),
			),
		),
		jen.Line(),
	)

	return code
}
