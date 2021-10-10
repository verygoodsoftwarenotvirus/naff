package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func apiClientsDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.Underscore().Qual(proj.QuerybuildingPackage(), "APIClientSQLQueryBuilder").Equals().Parens(jen.PointerTo().ID(dbvendor.Singular())).Call(jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(buildBuildGetBatchOfAPIClientsQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAPIClientByClientIDQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAllAPIClientsCountQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAPIClientsQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAPIClientByDatabaseIDQuery(proj, dbvendor)...)
	code.Add(buildBuildCreateAPIClientQuery(proj, dbvendor)...)
	code.Add(buildBuildUpdateAPIClientQuery(proj, dbvendor)...)
	code.Add(buildBuildArchiveAPIClientQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAuditLogEntriesForAPIClientQuery(proj, dbvendor)...)

	return code
}

func buildBuildGetBatchOfAPIClientsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetBatchOfAPIClientsQuery returns a query that fetches every API client in the database within a bucketed range."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetBatchOfAPIClientsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName")).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Gt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).Op(":").ID("beginID"))).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Lt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).Op(":").ID("endID"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAPIClientByClientIDQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAPIClientByClientIDQuery returns a SQL query which requests a given API client by its database ID."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetAPIClientByClientIDQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("string")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachAPIClientClientIDToSpan").Call(
				jen.ID("span"),
				jen.ID("clientID"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName")).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableClientIDColumn"),
				).Op(":").ID("clientID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAllAPIClientsCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAllAPIClientsCountQuery returns a SQL query for the number of API clients"),
		jen.Newline(),
		jen.Comment("in the database, regardless of ownership."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetAllAPIClientsCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQueryOnly").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.ID("columnCountQueryTemplate"),
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName"),
				)).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName")).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAPIClientsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAPIClientsQuery returns a SQL query (and arguments) that will retrieve a list of API clients that"),
		jen.Newline(),
		jen.Comment("meet the given filter's criteria (if relevant) and belong to a given account."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetAPIClientsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").Uint64(), jen.ID("filter").PointerTo().Qual(proj.TypesPackage(), "QueryFilter")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
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
			jen.Return().ID("b").Dot("buildListQuery").Callln(
				jen.ID("ctx"),
				jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName"),
				jen.ID("nil"),
				jen.ID("nil"),
				jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableOwnershipColumn"),
				jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableColumns"),
				jen.ID("userID"),
				jen.ID("false"),
				jen.ID("filter"),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAPIClientByDatabaseIDQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAPIClientByDatabaseIDQuery returns a SQL query which requests a given API client by its database ID."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetAPIClientByDatabaseIDQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachAPIClientDatabaseIDToSpan").Call(
				jen.ID("span"),
				jen.ID("clientID"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName")).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableOwnershipColumn"),
				).Op(":").ID("userID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).Op(":").ID("clientID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildCreateAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildCreateAPIClientQuery returns a SQL query (and args) that will create the given APIClient in the database."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildCreateAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").PointerTo().Qual(proj.TypesPackage(), "APIClientCreationInput")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName")).
					Dotln("Columns").Callln(
					jen.Qual(proj.QuerybuildingPackage(), "ExternalIDColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableNameColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableClientIDColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableSecretKeyColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableOwnershipColumn"),
				).
					Dotln("Values").Callln(
					jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
					jen.ID("input").Dot("Name"),
					jen.ID("input").Dot("ClientID"),
					jen.ID("input").Dot("ClientSecret"),
					jen.ID("input").Dot("BelongsToUser"),
				).Add(utils.ConditionalCode(dbvendor.SingularPackageName() == "postgres", jen.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s"), jen.Qual(proj.QuerybuildingPackage(), "IDColumn"))))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildUpdateAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildUpdateAPIClientQuery returns a SQL query (and args) that will update a given API client in the database."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildUpdateAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").PointerTo().Qual(proj.TypesPackage(), "APIClient")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("BelongsToUser"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachAPIClientClientIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("ClientID"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachAPIClientDatabaseIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("ID"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName")).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableClientIDColumn"),
					jen.ID("input").Dot("ClientID"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "IDColumn").Op(":").ID("input").Dot("ID"), jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableOwnershipColumn").Op(":").ID("input").Dot("BelongsToUser"), jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildArchiveAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildArchiveAPIClientQuery returns a SQL query (and arguments) that will mark an API client as archived."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildArchiveAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachAPIClientDatabaseIDToSpan").Call(
				jen.ID("span"),
				jen.ID("clientID"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName")).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "IDColumn").Op(":").ID("clientID"), jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"), jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableOwnershipColumn").Op(":").ID("userID"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAuditLogEntriesForAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAuditLogEntriesForAPIClientQuery constructs a SQL query for fetching audit log entries belong to a user with a given ID."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetAuditLogEntriesForAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("uint64")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachAPIClientDatabaseIDToSpan").Call(
				jen.ID("span"),
				jen.ID("clientID"),
			),
			jen.Newline(),
			jen.ID("apiClientIDKey").Assign().Qual("fmt", "Sprintf").Callln(
				jen.ID("jsonPluckQuery"),
				jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
				jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableContextColumn"),
				utils.ConditionalCode(dbvendor.SingularPackageName() == "mysql", jen.ID("clientID")),
				jen.Qual(proj.InternalAuditPackage(), "APIClientAssignmentKey"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				func() jen.Code {
					if dbvendor.SingularPackageName() == "mysql" {
						return jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableColumns").Op("...")).
							Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")).
							Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Expr").Call(jen.ID("apiClientIDKey"))).
							Dotln("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%s.%s"),
							jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
							jen.Qual(proj.QuerybuildingPackage(), "CreatedOnColumn"),
						))
					}
					return jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableColumns").Op("...")).
						Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")).
						Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Values(jen.ID("apiClientIDKey").Op(":").ID("clientID"))).
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
