package queriers

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
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
		fmt.Sprintf("%s.updated_on", webhooksTableName),
		fmt.Sprintf("%s.archived_on", webhooksTableName),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn),
	}
)

func webhooksTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, ret)

	ret.Add(buildBuildMockRowsFromWebhook(proj, dbvendor)...)
	ret.Add(buildBuildErroneousMockRowFromWebhook(proj, dbvendor)...)
	ret.Add(buildTestDB_buildGetWebhookQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_GetWebhook(proj, dbvendor)...)
	ret.Add(buildTestDB_buildGetAllWebhooksCountQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_GetAllWebhooksCount(proj, dbvendor)...)
	ret.Add(buildTestDB_buildGetAllWebhooksQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_GetAllWebhooks(proj, dbvendor)...)
	ret.Add(buildTestDB_buildGetWebhooksQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_GetWebhooks(proj, dbvendor)...)
	ret.Add(buildTestDB_buildWebhookCreationQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_CreateWebhook(proj, dbvendor)...)
	ret.Add(buildTestDB_buildUpdateWebhookQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_UpdateWebhook(proj, dbvendor)...)
	ret.Add(buildTestDB_buildArchiveWebhookQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_ArchiveWebhook(proj, dbvendor)...)

	return ret
}

func buildBuildMockRowsFromWebhook(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildMockRowsFromWebhook").Params(
			jen.ID("webhooks").Spread().PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
		).Params(
			jen.PointerTo().Qual("github.com/DATA-DOG/go-sqlmock", "Rows"),
		).Block(
			jen.ID("includeCount").Assign().Len(jen.ID("webhooks")).GreaterThan().One(),
			jen.ID("columns").Assign().ID("webhooksTableColumns"),
			jen.Line(),
			jen.If(jen.ID("includeCount")).Block(
				utils.AppendItemsToList(jen.ID("columns"), jen.Lit("count")),
			),
			jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("columns")),
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
					jen.ID("w").Dot("UpdatedOn"),
					jen.ID("w").Dot("ArchivedOn"),
					jen.ID("w").Dot("BelongsToUser"),
				),
				jen.Line(),
				jen.If(jen.ID("includeCount")).Block(
					utils.AppendItemsToList(jen.ID("rowValues"), jen.Len(jen.ID("webhooks"))),
				),
				jen.Line(),
				jen.ID("exampleRows").Dot("AddRow").Call(jen.ID("rowValues").Spread()),
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
			jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("webhooksTableColumns")).Dot("AddRow").Callln(
				jen.ID("w").Dot("ArchivedOn"),
				jen.ID("w").Dot("BelongsToUser"),
				jen.ID("w").Dot("Name"),
				jen.ID("w").Dot("ContentType"),
				jen.ID("w").Dot("URL"),
				jen.ID("w").Dot("Method"),
				jen.Qual("strings", "Join").Call(jen.ID("w").Dot("Events"), jen.ID("eventsSeparator")),
				jen.Qual("strings", "Join").Call(jen.ID("w").Dot("DataTypes"), jen.ID("typesSeparator")),
				jen.Qual("strings", "Join").Call(jen.ID("w").Dot("Topics"), jen.ID("topicsSeparator")),
				jen.ID("w").Dot("CreatedOn"),
				jen.ID("w").Dot("UpdatedOn"),
				jen.ID("w").Dot("ID"),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
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
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser"),
	}

	return buildQueryTest(proj,
		dbvendor,
		models.DataType{Name: GetWebhookPalabra()},
		"GetWebhook",
		qb,
		expectedArgs,
		callArgs,
		true,
		false,
		false,
		false,
		true,
		false,
		nil,
	)
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
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser"), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowsFromWebhook").Call(jen.ID(utils.BuildFakeVarName("Webhook")))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhook").Call(utils.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser")),
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
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser"), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhook").Call(utils.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser")),
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
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser"), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhook").Call(utils.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser")),
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
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser"), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.ID(utils.BuildFakeVarName("Webhook")))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhook").Call(utils.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser")),
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

	return buildQueryTest(proj,
		dbvendor,
		models.DataType{Name: GetWebhookPalabra()},
		"GetAllWebhooksCount",
		qb,
		expectedArgs,
		callArgs,
		false,
		false,
		false,
		false,
		false,
		false,
		nil,
	)
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
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooksCount").Call(utils.CtxVar()),
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
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooksCount").Call(utils.CtxVar()),
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

	return buildQueryTest(proj,
		dbvendor,
		models.DataType{Name: GetWebhookPalabra()},
		"GetAllWebhooks",
		qb,
		expectedArgs,
		callArgs,
		false,
		false,
		false,
		false,
		false,
		false,
		nil,
	)
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
			jen.ID("expectedListQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.Line(),
				utils.BuildFakeVar(proj, "WebhookList"),
				jen.ID(utils.BuildFakeVarName("WebhookList")).Dot("Limit").Equals().Zero(),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowsFromWebhook").Callln(
						jen.AddressOf().ID(utils.BuildFakeVarName("WebhookList")).Dot("Webhooks").Index(jen.Zero()),
						jen.AddressOf().ID(utils.BuildFakeVarName("WebhookList")).Dot("Webhooks").Index(jen.One()),
						jen.AddressOf().ID(utils.BuildFakeVarName("WebhookList")).Dot("Webhooks").Index(jen.Lit(2)),
					),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooks").Call(utils.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("WebhookList")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"surfaces sql.ErrNoRows",
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooks").Call(utils.CtxVar()),
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
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooks").Call(utils.CtxVar()),
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
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.ID(utils.BuildFakeVarName("Webhook")))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooks").Call(utils.CtxVar()),
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
		Select(append(webhooksTableColumns, fmt.Sprintf(countQuery, webhooksTableName))...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): whateverValue,
			fmt.Sprintf("%s.archived_on", webhooksTableName):                      nil,
		})

	qb = applyFleshedOutQueryFilter(qb, webhooksTableName)

	expectedArgs := appendFleshedOutQueryFilterArgs([]jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
	})
	callArgs := []jen.Code{}

	return buildQueryTest(proj,
		dbvendor,
		models.DataType{Name: GetWebhookPalabra()},
		"GetWebhooks",
		qb,
		expectedArgs,
		callArgs,
		true,
		false,
		true,
		true,
		false,
		false,
		nil,
	)
}

func buildTestDB_GetWebhooks(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(append(webhooksTableColumns, fmt.Sprintf(countQuery, webhooksTableName))...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): whateverValue,
			fmt.Sprintf("%s.archived_on", webhooksTableName):                      nil,
		}).
		GroupBy(fmt.Sprintf("%s.id", webhooksTableName)).
		Limit(20).
		ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetWebhooks", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildFakeVar(proj, "User"),
			jen.ID("expectedListQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.CreateDefaultQueryFilter(proj),
				utils.BuildFakeVar(proj, "WebhookList"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
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
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID(utils.FilterVarName),
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
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhooks").Call(
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID(utils.FilterVarName),
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
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhooks").Call(
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID(utils.FilterVarName),
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
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.ID(utils.BuildFakeVarName("Webhook")))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhooks").Call(
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID(utils.FilterVarName),
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
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("Webhook")),
	}

	return buildQueryTest(proj,
		dbvendor,
		models.DataType{Name: GetWebhookPalabra()},
		"WebhookCreation",
		qb,
		expectedArgs,
		callArgs,
		true,
		false,
		false,
		false,
		true,
		false,
		nil,
	)
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
			return jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(
				jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
				jen.ID(utils.BuildFakeVarName("Webhook")).Dot("CreatedOn"),
			)
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			createWebhookExpectFunc = "ExpectExec"
			createWebhookReturnFunc = "WillReturnResult"
			return jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID")), jen.One())
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
							jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser"),
						).Dot(createWebhookReturnFunc).Call(jen.ID("exampleRows")),
						jen.Line(),
					}

					if isSqlite(dbvendor) || isMariaDB(dbvendor) {
						out = append(out,
							jen.ID("expectedTimeQuery").Assign().Litf("SELECT created_on FROM webhooks WHERE id = %s", getIncIndex(dbvendor, 0)),
							jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).
								Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("CreatedOn"))),
							jen.Line(),
						)
					}

					out = append(out,
						jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateWebhook").Call(utils.CtxVar(), jen.ID("exampleInput")),
						utils.AssertNoError(jen.Err(), nil),
						utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Webhook")), jen.ID("actual"), nil),
						jen.Line(),
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
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser"),
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateWebhook").Call(utils.CtxVar(), jen.ID("exampleInput")),
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
		Set("updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id":                         whateverValue,
			webhooksTableOwnershipColumn: whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING updated_on")
	}

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Name"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ContentType"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("URL"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Method"),
		jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Events"), jen.ID("eventsSeparator")),
		jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("DataTypes"), jen.ID("typesSeparator")),
		jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("Topics"), jen.ID("topicsSeparator")),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("Webhook")),
	}

	return buildQueryTest(proj,
		dbvendor,
		models.DataType{Name: GetWebhookPalabra()},
		"UpdateWebhook",
		qb,
		expectedArgs,
		callArgs,
		true,
		false,
		false,
		false,
		true,
		false,
		nil,
	)
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
		Set("updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id":                         whateverValue,
			webhooksTableOwnershipColumn: whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING updated_on")
	}
	expectedQuery, _, _ := qb.ToSql()

	var updateWebhookExpectFunc, updateWebhookReturnFunc string
	buildUpdateWebhookExampleRows := func() jen.Code {
		if isPostgres(dbvendor) {
			updateWebhookExpectFunc = "ExpectQuery"
			updateWebhookReturnFunc = "WillReturnRows"
			return jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("updated_on"))).Dot("AddRow").Call(
				jen.ID(utils.BuildFakeVarName("Webhook")).Dot("UpdatedOn"),
			)
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			updateWebhookExpectFunc = "ExpectExec"
			updateWebhookReturnFunc = "WillReturnResult"
			return jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID")), jen.One())
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
						jen.ID("topicsSeparator")), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
				).Dot(updateWebhookReturnFunc).Call(jen.ID("exampleRows")),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("UpdateWebhook").Call(utils.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook"))),
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
						jen.ID("topicsSeparator")), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("UpdateWebhook").Call(utils.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook"))),
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
		Set("updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
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
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
		jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser"),
	}

	return buildQueryTest(proj,
		dbvendor,
		models.DataType{Name: GetWebhookPalabra()},
		"ArchiveWebhook",
		qb,
		expectedArgs,
		callArgs,
		true,
		false,
		false,
		false,
		true,
		false,
		nil,
	)
}

func buildTestDB_ArchiveWebhook(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	qb := queryBuilderForDatabase(dbvendor).
		Update(webhooksTableName).
		Set("updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
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
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
				).Dot("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.One(), jen.One())),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("ArchiveWebhook").Call(utils.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("BelongsToUser")),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}
