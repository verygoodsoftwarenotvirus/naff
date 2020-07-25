package queriers

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"strings"
)

const (
	webhooksTableName            = "webhooks"
	webhooksTableOwnershipColumn = "belongs_to_user"
)

var (
	webhooksTableColumns = []string{
		fmt.Sprintf("%s.id", webhooksTableName),
		fmt.Sprintf("%s.name", webhooksTableName),
		fmt.Sprintf("%s.content_type", webhooksTableName),
		fmt.Sprintf("%s.url", webhooksTableName),
		fmt.Sprintf("%s.method", webhooksTableName),
		fmt.Sprintf("%s.events", webhooksTableName),
		fmt.Sprintf("%s.data_types", webhooksTableName),
		fmt.Sprintf("%s.topics", webhooksTableName),
		fmt.Sprintf("%s.created_on", webhooksTableName),
		fmt.Sprintf("%s.last_updated_on", webhooksTableName),
		fmt.Sprintf("%s.archived_on", webhooksTableName),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn),
	}
)

func webhooksTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	code := jen.NewFilePathName(proj.DatabaseV1Package("queriers", "v1", spn), spn)

	utils.AddImports(proj, code)

	code.Add(buildBuildMockRowsFromWebhook(proj, dbvendor)...)
	code.Add(buildBuildErroneousMockRowFromWebhook(proj, dbvendor)...)
	code.Add(buildTestScanWebhooks(proj, dbvendor)...)
	code.Add(buildTestDB_buildGetWebhookQuery(proj, dbvendor)...)
	code.Add(buildTestDB_GetWebhook(proj, dbvendor)...)
	code.Add(buildTestDB_buildGetAllWebhooksCountQuery(proj, dbvendor)...)
	code.Add(buildTestDB_GetAllWebhooksCount(proj, dbvendor)...)
	code.Add(buildTestDB_buildGetAllWebhooksQuery(proj, dbvendor)...)
	code.Add(buildTestDB_GetAllWebhooks(proj, dbvendor)...)
	code.Add(buildTestDB_buildGetWebhooksQuery(proj, dbvendor)...)
	code.Add(buildTestDB_GetWebhooks(proj, dbvendor)...)
	code.Add(buildTestDB_buildWebhookCreationQuery(proj, dbvendor)...)
	code.Add(buildTestDB_CreateWebhook(proj, dbvendor)...)
	code.Add(buildTestDB_buildUpdateWebhookQuery(proj, dbvendor)...)
	code.Add(buildTestDB_UpdateWebhook(proj, dbvendor)...)
	code.Add(buildTestDB_buildArchiveWebhookQuery(proj, dbvendor)...)
	code.Add(buildTestDB_ArchiveWebhook(proj, dbvendor)...)

	return code
}

func buildBuildMockRowsFromWebhook(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildMockRowsFromWebhook").Params(
			jen.ID("webhooks").Spread().PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
		).Params(
			jen.PointerTo().Qual("github.com/DATA-DOG/go-sqlmock", "Rows"),
		).Block(
			jen.ID("columns").Assign().ID("webhooksTableColumns"),
			jen.Line(),
			jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("columns")),
			jen.Line(),
			jen.For(jen.List(jen.Underscore(), jen.ID("w")).Assign().Range().ID("webhooks")).Block(
				jen.ID("rowValues").Assign().Index().Qual("database/sql/driver", "Value").Valuesln(
					jen.ID("w").Dot("ID"),
					jen.ID("w").Dot("Name"),
					jen.ID("w").Dot("ContentType"),
					jen.ID("w").Dot("URL"),
					jen.ID("w").Dot("Method"),
					jen.Qual("strings", "Join").Call(jen.ID("w").Dot("Events"), jen.ID("eventsSeparator")),
					jen.Qual("strings", "Join").Call(jen.ID("w").Dot("DataTypes"), jen.ID("typesSeparator")),
					jen.Qual("strings", "Join").Call(jen.ID("w").Dot("Topics"), jen.ID("topicsSeparator")),
					jen.ID("w").Dot("CreatedOn"),
					jen.ID("w").Dot("LastUpdatedOn"),
					jen.ID("w").Dot("ArchivedOn"),
					jen.ID("w").Dot(constants.UserOwnershipFieldName),
				),
				jen.Line(),
				jen.ID(utils.BuildFakeVarName("Rows")).Dot("AddRow").Call(jen.ID("rowValues").Spread()),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildErroneousMockRowFromWebhook(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildErroneousMockRowFromWebhook").Params(jen.ID("w").PointerTo().Qual(proj.ModelsV1Package(), "Webhook")).Params(jen.PointerTo().Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("webhooksTableColumns")).Dot("AddRow").Callln(
				jen.ID("w").Dot("ArchivedOn"),
				jen.ID("w").Dot(constants.UserOwnershipFieldName),
				jen.ID("w").Dot("Name"),
				jen.ID("w").Dot("ContentType"),
				jen.ID("w").Dot("URL"),
				jen.ID("w").Dot("Method"),
				jen.Qual("strings", "Join").Call(jen.ID("w").Dot("Events"), jen.ID("eventsSeparator")),
				jen.Qual("strings", "Join").Call(jen.ID("w").Dot("DataTypes"), jen.ID("typesSeparator")),
				jen.Qual("strings", "Join").Call(jen.ID("w").Dot("Topics"), jen.ID("topicsSeparator")),
				jen.ID("w").Dot("CreatedOn"),
				jen.ID("w").Dot("LastUpdatedOn"),
				jen.ID("w").Dot("ID"),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	}

	return lines
}

func buildTestScanWebhooks(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	lines := []jen.Code{
		jen.Func().IDf("Test%s_ScanWebhooks", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"surfaces row errors",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockRows").Assign().AddressOf().Qual(proj.DatabaseV1Package(), "MockResultIterator").Values(),
				jen.Line(),
				jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.False()),
				jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(
					jen.Underscore(),
					jen.Err(),
				).Assign().ID(dbfl).Dot("scanWebhooks").Call(jen.ID("mockRows")),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"logs row closing errors",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockRows").Assign().AddressOf().Qual(proj.DatabaseV1Package(), "MockResultIterator").Values(),
				jen.Line(),
				jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.False()),
				jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.Nil()),
				jen.ID("mockRows").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(
					jen.Underscore(),
					jen.Err(),
				).Assign().ID(dbfl).Dot("scanWebhooks").Call(jen.ID("mockRows")),
				utils.AssertNoError(jen.Err(), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDB_buildGetWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", webhooksTableName):                               whateverValue,
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): whateverValue,
		})

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "Webhook"),
	}

	return buildQueryTest(proj, dbvendor, "GetWebhook", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_GetWebhook(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", webhooksTableName):                               whateverValue,
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): whateverValue,
		}).
		ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetWebhook", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowsFromWebhook").Call(jen.ID(utils.BuildFakeVarName("Webhook")))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhook").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"), jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName)),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Webhook")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"surfaces sql.ErrNoRows",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhook").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"), jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName)),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error from database",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID")).
					Dotln("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhook").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"), jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName)),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid response from database",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.ID(utils.BuildFakeVarName("Webhook")))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhook").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"), jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName)),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDB_buildGetAllWebhooksCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(fmt.Sprintf(countQuery, webhooksTableName)).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", webhooksTableName): nil,
		})

	expectedArgs := []jen.Code{}
	callArgs := []jen.Code{}

	return buildQueryTest(proj, dbvendor, "GetAllWebhooksCount", qb, expectedArgs, callArgs, nil)
}

func buildTestDB_GetAllWebhooksCount(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(fmt.Sprintf(countQuery, webhooksTableName)).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", webhooksTableName): nil,
		}).ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetAllWebhooksCount", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID(utils.BuildFakeVarName("Count")).Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").
					Call(jen.Index().String().Values(jen.Lit("count"))).Dot("AddRow").
					Call(jen.ID(utils.BuildFakeVarName("Count")))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooksCount").Call(constants.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Count")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error from database",
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooksCount").Call(constants.CtxVar()),
				utils.AssertError(jen.Err(), nil),
				utils.AssertZero(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDB_buildGetAllWebhooksQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", webhooksTableName): nil,
		})

	expectedArgs := []jen.Code{}
	callArgs := []jen.Code{}

	return buildQueryTest(proj, dbvendor, "GetAllWebhooks", qb, expectedArgs, callArgs, nil)
}

func buildTestDB_GetAllWebhooks(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", webhooksTableName): nil,
		}).ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetAllWebhooks", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.Line(),
				utils.BuildFakeVar(proj, "WebhookList"),
				jen.ID(utils.BuildFakeVarName("WebhookList")).Dot("Limit").Equals().Zero(),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowsFromWebhook").Callln(
						jen.AddressOf().ID(utils.BuildFakeVarName("WebhookList")).Dot("Webhooks").Index(jen.Zero()),
						jen.AddressOf().ID(utils.BuildFakeVarName("WebhookList")).Dot("Webhooks").Index(jen.One()),
						jen.AddressOf().ID(utils.BuildFakeVarName("WebhookList")).Dot("Webhooks").Index(jen.Lit(2)),
					),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooks").Call(constants.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("WebhookList")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"surfaces sql.ErrNoRows",
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooks").Call(constants.CtxVar()),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error querying database",
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooks").Call(constants.CtxVar()),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error from database",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.ID(utils.BuildFakeVarName("Webhook")))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooks").Call(constants.CtxVar()),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDB_buildGetWebhooksQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): whateverValue,
			fmt.Sprintf("%s.archived_on", webhooksTableName):                      nil,
		})

	qb = applyFleshedOutQueryFilter(qb, webhooksTableName)

	expectedArgs := appendFleshedOutQueryFilterArgs([]jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
	})
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
		jen.ID(constants.FilterVarName),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "User"),
		jen.ID(constants.FilterVarName).Assign().Qual(proj.FakeModelsPackage(), "BuildFleshedOutQueryFilter").Call(),
	}

	return buildQueryTest(proj, dbvendor, "GetWebhooks", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_GetWebhooks(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): whateverValue,
			fmt.Sprintf("%s.archived_on", webhooksTableName):                      nil,
		}).
		OrderBy(fmt.Sprintf("%s.id", webhooksTableName)).
		Limit(20).
		ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetWebhooks", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildFakeVar(proj, "User"),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.CreateDefaultQueryFilter(proj),
				utils.BuildFakeVar(proj, "WebhookList"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("User")).Dot("ID")).
					Dotln("WillReturnRows").Callln(
					jen.ID("buildMockRowsFromWebhook").Callln(
						jen.AddressOf().ID(utils.BuildFakeVarName("WebhookList")).Dot("Webhooks").Index(jen.Zero()),
						jen.AddressOf().ID(utils.BuildFakeVarName("WebhookList")).Dot("Webhooks").Index(jen.One()),
						jen.AddressOf().ID(utils.BuildFakeVarName("WebhookList")).Dot("Webhooks").Index(jen.Lit(2)),
					),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhooks").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(constants.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("WebhookList")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"surfaces sql.ErrNoRows",
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhooks").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(constants.FilterVarName),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error querying database",
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhooks").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(constants.FilterVarName),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with erroneous response from database",
				utils.CreateDefaultQueryFilter(proj),
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.ID(utils.BuildFakeVarName("Webhook")))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhooks").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(constants.FilterVarName),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDB_buildWebhookCreationQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Insert(webhooksTableName).
		Columns(
			"name",
			"content_type",
			"url",
			"method",
			"events",
			"data_types",
			"topics",
			webhooksTableOwnershipColumn,
		).
		Values(
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
		)

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING id, created_on")
	}

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Name"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ContentType"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("URL"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Method"),
		jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Events"), jen.ID("eventsSeparator")),
		jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("DataTypes"), jen.ID("typesSeparator")),
		jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Topics"), jen.ID("topicsSeparator")),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("Webhook")),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "Webhook"),
	}

	return buildQueryTest(proj, dbvendor, "WebhookCreation", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_CreateWebhook(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	qb := queryBuilderForDatabase(dbvendor).
		Insert(webhooksTableName).
		Columns(
			"name",
			"content_type",
			"url",
			"method",
			"events",
			"data_types",
			"topics",
			webhooksTableOwnershipColumn,
		).
		Values(
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
		)

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING id, created_on")
	}
	expectedQuery, _, _ := qb.ToSql()

	var createWebhookExpectFunc, createWebhookReturnFunc string
	buildCreateWebhookExampleRows := func() jen.Code {
		if isPostgres(dbvendor) {
			createWebhookExpectFunc = "ExpectQuery"
			createWebhookReturnFunc = "WillReturnRows"
			return jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(
				jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
				jen.ID(utils.BuildFakeVarName("Webhook")).Dot("CreatedOn"),
			)
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			createWebhookExpectFunc = "ExpectExec"
			createWebhookReturnFunc = "WillReturnResult"
			return jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID")), jen.One())
		}
		return jen.Null()
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_CreateWebhook", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				func() []jen.Code {
					out := []jen.Code{
						utils.BuildFakeVar(proj, "Webhook"),
						utils.BuildFakeVarWithCustomName(proj, "exampleInput", "BuildFakeWebhookCreationInputFromWebhook", jen.ID(utils.BuildFakeVarName("Webhook"))),
						buildCreateWebhookExampleRows(),
						jen.Line(),
						jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
						jen.ID("mockDB").Dot(createWebhookExpectFunc).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
							jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Name"),
							jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ContentType"),
							jen.ID(utils.BuildFakeVarName("Webhook")).Dot("URL"),
							jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Method"),
							jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Events"), jen.ID("eventsSeparator")),
							jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("DataTypes"), jen.ID("typesSeparator")),
							jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Topics"), jen.ID("topicsSeparator")),
							jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName),
						).Dot(createWebhookReturnFunc).Call(jen.ID(utils.BuildFakeVarName("Rows"))),
						jen.Line(),
					}

					if isSqlite(dbvendor) || isMariaDB(dbvendor) {
						out = append(out,
							jen.IDf("%stt", dbfl).Assign().AddressOf().ID("mockTimeTeller").Values(),
							jen.IDf("%stt", dbfl).Dot("On").Call(jen.Lit("Now")).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("CreatedOn")),
							jen.ID(dbfl).Dot("timeTeller").Equals().IDf("%stt", dbfl),
							jen.Line(),
						)
					}

					out = append(out,
						jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateWebhook").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Input"))),
						utils.AssertNoError(jen.Err(), nil),
						utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Webhook")), jen.ID("actual"), nil),
						jen.Line(),
						func() jen.Code {
							if isMariaDB(dbvendor) || isSqlite(dbvendor) {
								return utils.AssertExpectationsFor(fmt.Sprintf("%stt", dbfl))
							}
							return jen.Null()
						}(),
						utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
					)

					return out
				}()...,
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error interacting with database",
				utils.BuildFakeVar(proj, "Webhook"),
				utils.BuildFakeVarWithCustomName(proj, "exampleInput", "BuildFakeWebhookCreationInputFromWebhook", jen.ID(utils.BuildFakeVarName("Webhook"))),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(createWebhookExpectFunc).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Name"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ContentType"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("URL"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Method"),
					jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Events"), jen.ID("eventsSeparator")),
					jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("DataTypes"), jen.ID("typesSeparator")),
					jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Topics"), jen.ID("topicsSeparator")),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName),
				).Dot("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateWebhook").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Input"))),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDB_buildUpdateWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Update(webhooksTableName).
		Set("name", whateverValue).
		Set("content_type", whateverValue).
		Set("url", whateverValue).
		Set("method", whateverValue).
		Set("events", whateverValue).
		Set("data_types", whateverValue).
		Set("topics", whateverValue).
		Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id":                         whateverValue,
			webhooksTableOwnershipColumn: whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING last_updated_on")
	}

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Name"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ContentType"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("URL"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Method"),
		jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Events"), jen.ID("eventsSeparator")),
		jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("DataTypes"), jen.ID("typesSeparator")),
		jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Topics"), jen.ID("topicsSeparator")),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("Webhook")),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "Webhook"),
	}

	return buildQueryTest(proj, dbvendor, "UpdateWebhook", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_UpdateWebhook(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	qb := queryBuilderForDatabase(dbvendor).
		Update(webhooksTableName).
		Set("name", whateverValue).
		Set("content_type", whateverValue).
		Set("url", whateverValue).
		Set("method", whateverValue).
		Set("events", whateverValue).
		Set("data_types", whateverValue).
		Set("topics", whateverValue).
		Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id":                         whateverValue,
			webhooksTableOwnershipColumn: whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING last_updated_on")
	}
	expectedQuery, _, _ := qb.ToSql()

	var updateWebhookExpectFunc, updateWebhookReturnFunc string
	buildUpdateWebhookExampleRows := func() jen.Code {
		if isPostgres(dbvendor) {
			updateWebhookExpectFunc = "ExpectQuery"
			updateWebhookReturnFunc = "WillReturnRows"
			return jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("last_updated_on"))).Dot("AddRow").Call(
				jen.ID(utils.BuildFakeVarName("Webhook")).Dot("LastUpdatedOn"),
			)
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			updateWebhookExpectFunc = "ExpectExec"
			updateWebhookReturnFunc = "WillReturnResult"
			return jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID")), jen.One())
		}
		return jen.Null()
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_UpdateWebhook", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				buildUpdateWebhookExampleRows(),
				jen.ID("mockDB").Dot(updateWebhookExpectFunc).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Name"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ContentType"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("URL"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Method"),
					jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Events"),
						jen.ID("eventsSeparator")), jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("DataTypes"),
						jen.ID("typesSeparator")), jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Topics"),
						jen.ID("topicsSeparator")), jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
				).Dot(updateWebhookReturnFunc).Call(jen.ID(utils.BuildFakeVarName("Rows"))),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("UpdateWebhook").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook"))),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error from database",
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.ID("mockDB").Dot(updateWebhookExpectFunc).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Name"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ContentType"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("URL"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Method"),
					jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Events"),
						jen.ID("eventsSeparator")), jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("DataTypes"),
						jen.ID("typesSeparator")), jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Topics"),
						jen.ID("topicsSeparator")), jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
				).Dot("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("UpdateWebhook").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook"))),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDB_buildArchiveWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Update(webhooksTableName).
		Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Set("archived_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id":                         whateverValue,
			webhooksTableOwnershipColumn: whateverValue,
			"archived_on":                nil,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING archived_on")
	}

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "Webhook"),
	}

	return buildQueryTest(proj, dbvendor, "ArchiveWebhook", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_ArchiveWebhook(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	qb := queryBuilderForDatabase(dbvendor).
		Update(webhooksTableName).
		Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Set("archived_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id":                         whateverValue,
			webhooksTableOwnershipColumn: whateverValue,
			"archived_on":                nil,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING archived_on")
	}
	expectedQuery, _, _ := qb.ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_ArchiveWebhook", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.ID("expectedQuery").Assign().Lit(expectedQuery),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
				).Dot("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.One(), jen.One())),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("ArchiveWebhook").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"), jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName)),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}
