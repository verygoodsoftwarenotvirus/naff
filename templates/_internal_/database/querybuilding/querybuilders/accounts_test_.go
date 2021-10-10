package querybuilders

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/Masterminds/squirrel"
)

func accountsTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	code.Add(buildTestSqlite_BuildGetAccountQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetAllAccountsCountQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetBatchOfAccountsQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetAccountsQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildCreateAccountQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildUpdateAccountQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildArchiveAccountQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetAuditLogEntriesForAccountQuery(proj, dbvendor)...)

	return code
}

func buildTestSqlite_BuildGetAccountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Select(
			"accounts.id",
			"accounts.external_id",
			"accounts.name",
			"accounts.billing_status",
			"accounts.contact_email",
			"accounts.contact_phone",
			"accounts.payment_processor_customer_id",
			"accounts.subscription_plan_id",
			"accounts.created_on",
			"accounts.last_updated_on",
			"accounts.archived_on",
			"accounts.belongs_to_user",
			"account_user_memberships.id",
			"account_user_memberships.belongs_to_user",
			"account_user_memberships.belongs_to_account",
			"account_user_memberships.account_roles",
			"account_user_memberships.default_account",
			"account_user_memberships.created_on",
			"account_user_memberships.last_updated_on",
			"account_user_memberships.archived_on",
		).
			From("accounts").
			Join("account_user_memberships ON account_user_memberships.belongs_to_account = accounts.id").
			Where(squirrel.Eq{
				"accounts.id":              whateverValue,
				"accounts.belongs_to_user": whateverValue,
				"accounts.archived_on":     nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAccountQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleAccount").Dot("BelongsToUser"), jen.ID("exampleAccount").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetAccountQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildGetAllAccountsCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(queryBuilderForDatabase(dbvendor).
		Select(fmt.Sprintf(columnCountQueryTemplate, "accounts")).
		From("accounts").
		Where(squirrel.Eq{
			"accounts.archived_on": nil,
		}))

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAllAccountsCountQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("actualQuery").Assign().ID("q").Dot("BuildGetAllAccountsCountQuery").Call(jen.ID("ctx")),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.Index().Interface().Values(),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildGetBatchOfAccountsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Select(
			"accounts.id",
			"accounts.external_id",
			"accounts.name",
			"accounts.billing_status",
			"accounts.contact_email",
			"accounts.contact_phone",
			"accounts.payment_processor_customer_id",
			"accounts.subscription_plan_id",
			"accounts.created_on",
			"accounts.last_updated_on",
			"accounts.archived_on",
			"accounts.belongs_to_user",
		).
			From("accounts").
			Where(squirrel.Gt{
				"accounts.id": whateverValue,
			}).
			Where(squirrel.Lt{
				"accounts.id": whateverValue,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetBatchOfAccountsQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.List(jen.ID("beginID"), jen.ID("endID")).Assign().List(jen.Uint64().Call(jen.Lit(1)), jen.Uint64().Call(jen.Lit(1000))),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("beginID"), jen.ID("endID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetBatchOfAccountsQuery").Call(
						jen.ID("ctx"),
						jen.ID("beginID"),
						jen.ID("endID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildGetAccountsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAccountsQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("filter").Assign().Qual(proj.FakeTypesPackage(), "BuildFleshedOutQueryFilter").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Litf("SELECT accounts.id, accounts.external_id, accounts.name, accounts.billing_status, accounts.contact_email, accounts.contact_phone, accounts.payment_processor_customer_id, accounts.subscription_plan_id, accounts.created_on, accounts.last_updated_on, accounts.archived_on, accounts.belongs_to_user, account_user_memberships.id, account_user_memberships.belongs_to_user, account_user_memberships.belongs_to_account, account_user_memberships.account_roles, account_user_memberships.default_account, account_user_memberships.created_on, account_user_memberships.last_updated_on, account_user_memberships.archived_on, (SELECT COUNT(accounts.id) FROM accounts WHERE accounts.archived_on IS NULL AND accounts.belongs_to_user = %s) as total_count, (SELECT COUNT(accounts.id) FROM accounts WHERE accounts.archived_on IS NULL AND accounts.belongs_to_user = %s AND accounts.created_on > %s AND accounts.created_on < %s AND accounts.last_updated_on > %s AND accounts.last_updated_on < %s) as filtered_count FROM accounts JOIN account_user_memberships ON account_user_memberships.belongs_to_account = accounts.id WHERE accounts.archived_on IS NULL AND accounts.belongs_to_user = %s AND accounts.created_on > %s AND accounts.created_on < %s AND accounts.last_updated_on > %s AND accounts.last_updated_on < %s GROUP BY accounts.id, account_user_memberships.id LIMIT 20 OFFSET 180", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), getIncIndex(dbvendor, 2), getIncIndex(dbvendor, 3), getIncIndex(dbvendor, 4), getIncIndex(dbvendor, 5), getIncIndex(dbvendor, 6), getIncIndex(dbvendor, 7), getIncIndex(dbvendor, 8), getIncIndex(dbvendor, 9), getIncIndex(dbvendor, 10)),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore"), jen.ID("exampleUser").Dot("ID"), jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetAccountsQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("false"),
						jen.ID("filter"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildCreateAccountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).Insert("accounts").
		Columns(
			"external_id",
			"name",
			"billing_status",
			"contact_email",
			"contact_phone",
			"belongs_to_user",
		).
		Values(
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
		)

	if dbvendor.SingularPackageName() == "postgres" {
		qb = qb.Suffix("RETURNING id")
	}

	expectedQuery, _ := buildQuery(qb)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildCreateAccountQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
					jen.Newline(),
					jen.ID("exIDGen").Assign().AddressOf().Qual(proj.QuerybuildingPackage(), "MockExternalIDGenerator").Values(),
					jen.ID("exIDGen").Dot("On").Call(jen.Lit("NewExternalID")).Dot("Return").Call(jen.ID("exampleAccount").Dot("ExternalID")),
					jen.ID("q").Dot("externalIDGenerator").Equals().ID("exIDGen"),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleAccount").Dot("ExternalID"), jen.ID("exampleAccount").Dot("Name"), jen.Qual(proj.TypesPackage(), "UnpaidAccountBillingStatus"), jen.ID("exampleAccount").Dot("ContactEmail"), jen.ID("exampleAccount").Dot("ContactPhone"), jen.ID("exampleAccount").Dot("BelongsToUser")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildAccountCreationQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("exIDGen"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildUpdateAccountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Update("accounts").
			Set("name", whateverValue).
			Set("contact_email", whateverValue).
			Set("contact_phone", whateverValue).
			Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
			Where(squirrel.Eq{
				"id":              whateverValue,
				"archived_on":     nil,
				"belongs_to_user": whateverValue,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildUpdateAccountQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleAccount").Dot("Name"), jen.ID("exampleAccount").Dot("ContactEmail"), jen.ID("exampleAccount").Dot("ContactPhone"), jen.ID("exampleAccount").Dot("BelongsToUser"), jen.ID("exampleAccount").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildUpdateAccountQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildArchiveAccountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Update("accounts").
			Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
			Set("archived_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
			Where(squirrel.Eq{
				"id":              whateverValue,
				"archived_on":     nil,
				"belongs_to_user": whateverValue,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildArchiveAccountQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("exampleAccount").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildArchiveAccountQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildGetAuditLogEntriesForAccountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	var accountIDKey string

	switch dbvendor.LowercaseAbbreviation() {
	case "m":
		accountIDKey = fmt.Sprintf(`JSON_CONTAINS(%s.%s, '%s', '$.%s')`, "audit_log", "context", "%d", "account_id")
	case "p":
		accountIDKey = fmt.Sprintf(`%s.%s->'%s'`, "audit_log", "context", "account_id")
	case "s":
		accountIDKey = fmt.Sprintf(`json_extract(%s.%s, '$.%s')`, "audit_log", "context", "account_id")
	}

	qb := queryBuilderForDatabase(dbvendor).Select(
		"audit_log.id",
		"audit_log.external_id",
		"audit_log.event_type",
		"audit_log.context",
		"audit_log.created_on",
	).
		From("audit_log")

	if dbvendor.SingularPackageName() == "mysql" {
		qb = qb.Where(squirrel.Expr(accountIDKey))
	} else {
		qb = qb.Where(squirrel.Eq{accountIDKey: whateverValue})
	}

	qb = qb.OrderBy("audit_log.created_on")

	expectedQuery, _ := buildQuery(qb)

	expectedQueryDecl := jen.ID("expectedQuery").Assign().Lit(expectedQuery)
	if dbvendor.SingularPackageName() == "mysql" {
		expectedQueryDecl = jen.ID("expectedQuery").Assign().Qual("fmt", "Sprintf").Call(jen.Lit(expectedQuery), jen.ID("exampleAccount").Dot("ID"))
	}

	expectedArgsDecl := jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleAccount").Dot("ID"))
	if dbvendor.SingularPackageName() == "mysql" {
		expectedArgsDecl = jen.ID("expectedArgs").Assign().Index().Interface().Call(jen.Nil())
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAuditLogEntriesForAccountQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Newline(),
					expectedQueryDecl,
					expectedArgsDecl,
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetAuditLogEntriesForAccountQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}
