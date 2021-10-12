package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountsDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.Underscore().Qual(proj.QuerybuildingPackage(), "AccountSQLQueryBuilder").Equals().Parens(jen.PointerTo().ID(dbvendor.Singular())).Call(jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(buildBuildTestSqlite_BuildGetDefaultAccountIDForUserQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAllAccountsCountQuery(proj, dbvendor)...)
	code.Add(buildBuildGetBatchOfAccountsQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAccountsQuery(proj, dbvendor)...)
	code.Add(buildBuildAccountCreationQuery(proj, dbvendor)...)
	code.Add(buildBuildUpdateAccountQuery(proj, dbvendor)...)
	code.Add(buildBuildArchiveAccountQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAuditLogEntriesForAccountQuery(proj, dbvendor)...)

	return code
}

func buildBuildTestSqlite_BuildGetDefaultAccountIDForUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAccountQuery constructs a SQL query for fetching an account with a given ID belong to a user with a given ID."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.ID("columns").Assign().ID("append").Call(
				jen.Qual(proj.QuerybuildingPackage(), "AccountsTableColumns"),
				jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableColumns").Op("..."),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("columns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName")).
					Dotln("Join").Call(jen.Qual("fmt", "Sprintf").Callln(
					jen.Lit("%s ON %s.%s = %s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableAccountOwnershipColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				)).Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).Op(":").ID("accountID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableUserOwnershipColumn"),
				).Op(":").ID("userID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAllAccountsCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAllAccountsCountQuery returns a query that fetches the total number of accounts in the database."),
		jen.Newline(),
		jen.Comment("This query only gets generated once, and is otherwise returned from cache."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetAllAccountsCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.String()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQueryOnly").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.ID("columnCountQueryTemplate"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
				)).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName")).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetBatchOfAccountsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetBatchOfAccountsQuery returns a query that fetches every account in the database within a bucketed range."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetBatchOfAccountsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName")).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Gt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).Op(":").ID("beginID"))).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Lt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).Op(":").ID("endID"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAccountsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAccountsQuery builds a SQL query selecting accounts that adhere to a given QueryFilter and belong to a given account,"),
		jen.Newline(),
		jen.Comment("and returns both the query and the relevant args to pass to the query executor."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetAccountsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").String(), jen.ID("forAdmin").ID("bool"), jen.ID("filter").PointerTo().Qual(proj.TypesPackage(), "QueryFilter")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.Newline(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.Qual(proj.InternalTracingPackage(), "AttachFilterToSpan").Call(
					jen.ID("span"),
					jen.ID("filter").Dot("Page"),
					jen.ID("filter").Dot("Limit"),
					jen.String().Call(jen.ID("filter").Dot("SortBy")),
				),
			),
			jen.Newline(),
			jen.Var().ID("includeArchived").ID("bool"),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("includeArchived").Equals().ID("filter").Dot("IncludeArchived"),
			),
			jen.Newline(),
			jen.ID("columns").Assign().ID("append").Call(
				jen.Qual(proj.QuerybuildingPackage(), "AccountsTableColumns"),
				jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableColumns").Op("..."),
			),
			jen.List(jen.ID("filteredCountQuery"), jen.ID("filteredCountQueryArgs")).Assign().ID("b").Dot("buildFilteredCountQuery").Call(
				jen.ID("ctx"),
				jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
				jen.ID("nil"),
				jen.ID("nil"),
				jen.Qual(proj.QuerybuildingPackage(), "AccountsTableUserOwnershipColumn"),
				jen.ID("userID"),
				jen.ID("forAdmin"),
				jen.ID("includeArchived"),
				jen.ID("filter"),
			),
			jen.List(jen.ID("totalCountQuery"), jen.ID("totalCountQueryArgs")).Assign().ID("b").Dot("buildTotalCountQuery").Call(
				jen.ID("ctx"),
				jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
				jen.ID("nil"),
				jen.ID("nil"),
				jen.Qual(proj.QuerybuildingPackage(), "AccountsTableUserOwnershipColumn"),
				jen.ID("userID"),
				jen.ID("forAdmin"),
				jen.ID("includeArchived"),
			),
			jen.Newline(),
			jen.ID("builder").Assign().ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("append").Callln(
				jen.ID("columns"),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("(%s) as total_count"),
					jen.ID("totalCountQuery"),
				),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("(%s) as filtered_count"),
					jen.ID("filteredCountQuery"),
				),
			).Op("...")).
				Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName")).
				Dotln("Join").Call(jen.Qual("fmt", "Sprintf").Callln(
				jen.Lit("%s ON %s.%s = %s.%s"),
				jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
				jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
				jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableAccountOwnershipColumn"),
				jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
				jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
			)),
			jen.Newline(),
			jen.If(jen.Op("!").ID("forAdmin")).Body(
				jen.ID("builder").Equals().ID("builder").Dot("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
				).Op(":").ID("nil"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableUserOwnershipColumn"),
				).Op(":").ID("userID"))),
			),
			jen.Newline(),
			jen.ID("builder").Equals().ID("builder").Dot("GroupBy").Call(jen.Qual("fmt", "Sprintf").Callln(
				jen.Lit("%s.%s, %s.%s"),
				jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
				jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
				jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
			)),
			jen.Newline(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("builder").Equals().Qual(proj.QuerybuildingPackage(), "ApplyFilterToQueryBuilder").Call(
					jen.ID("filter"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
					jen.ID("builder"),
				),
			),
			jen.Newline(),
			jen.List(jen.ID("query"), jen.ID("selectArgs")).Assign().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("builder"),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("query"), jen.ID("append").Call(
				jen.ID("append").Call(
					jen.ID("filteredCountQueryArgs"),
					jen.ID("totalCountQueryArgs").Op("..."),
				),
				jen.ID("selectArgs").Op("..."),
			)),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildAccountCreationQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildAccountCreationQuery takes an account and returns a creation query for that account and the relevant arguments."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildAccountCreationQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").PointerTo().Qual(proj.TypesPackage(), "AccountCreationInput")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName")).
					Dotln("Columns").Callln(
					jen.Qual(proj.QuerybuildingPackage(), "ExternalIDColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableNameColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableBillingStatusColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableContactEmailColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableContactPhoneColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableUserOwnershipColumn"),
				).
					Dotln("Values").Callln(
					jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
					jen.ID("input").Dot("Name"),
					jen.Qual(proj.TypesPackage(), "UnpaidAccountBillingStatus"),
					jen.ID("input").Dot("ContactEmail"),
					jen.ID("input").Dot("ContactPhone"),
					jen.ID("input").Dot("BelongsToUser"),
				).Add(utils.ConditionalCode(dbvendor.SingularPackageName() == "postgres", jen.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s"), jen.Qual(proj.QuerybuildingPackage(), "IDColumn"))))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildUpdateAccountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildUpdateAccountQuery takes an account and returns an update SQL query, with the relevant query parameters."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildUpdateAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").PointerTo().Qual(proj.TypesPackage(), "Account")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("ID"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName")).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableNameColumn"),
					jen.ID("input").Dot("Name"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableContactEmailColumn"),
					jen.ID("input").Dot("ContactEmail"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableContactPhoneColumn"),
					jen.ID("input").Dot("ContactPhone"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "IDColumn").Op(":").ID("input").Dot("ID"), jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"), jen.Qual(proj.QuerybuildingPackage(), "AccountsTableUserOwnershipColumn").Op(":").ID("input").Dot("BelongsToUser"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildArchiveAccountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildArchiveAccountQuery returns a SQL query which marks a given account belonging to a given user as archived."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildArchiveAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName")).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "IDColumn").Op(":").ID("accountID"), jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"), jen.Qual(proj.QuerybuildingPackage(), "AccountsTableUserOwnershipColumn").Op(":").ID("userID"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAuditLogEntriesForAccountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAuditLogEntriesForAccountQuery constructs a SQL query for fetching audit log entries belong to an account with a given ID."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetAuditLogEntriesForAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.ID("accountIDKey").Assign().Qual("fmt", "Sprintf").Callln(
				jen.ID("jsonPluckQuery"),
				jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
				jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableContextColumn"),
				utils.ConditionalCode(dbvendor.SingularPackageName() == "mysql", jen.ID("accountID")),
				jen.Qual(proj.InternalAuditPackage(), "AccountAssignmentKey"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				func() jen.Code {
					if dbvendor.SingularPackageName() == "mysql" {
						return jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableColumns").Op("...")).
							Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")).
							Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Expr").Call(jen.ID("accountIDKey"))).
							Dotln("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%s.%s"),
							jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
							jen.Qual(proj.QuerybuildingPackage(), "CreatedOnColumn"),
						))
					}
					return jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableColumns").Op("...")).
						Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")).
						Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Values(jen.ID("accountIDKey").Op(":").ID("accountID"))).
						Dotln("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
						jen.Qual(proj.QuerybuildingPackage(), "CreatedOnColumn"),
					))
				}(),
			),
		),
		jen.Newline(),
	}

	return lines
}
