package queriers

import (
	"path/filepath"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksTestDotGo(pkg *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(pkg, ret)
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))
	dbrn := dbvendor.RouteName()

	isPostgres := dbrn == "postgres"
	isSqlite := dbrn == "sqlite"
	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	/////////////

	ret.Add(
		jen.Func().ID("buildMockRowFromWebhook").Params(jen.ID("w").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook")).Params(jen.ParamPointer().Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("webhooksTableColumns")).Dot("AddRow").Callln(
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
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().ID("buildErroneousMockRowFromWebhook").Params(jen.ID("w").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook")).Params(jen.ParamPointer().Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
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
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_buildGetWebhookQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleWebhookID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.ID("expectedArgCount").Assign().Lit(2),
				jen.ID("expectedQuery").Assign().Litf("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to_user FROM webhooks WHERE belongs_to_user = %s AND id = %s", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetWebhookQuery").Call(jen.ID("exampleWebhookID"), jen.ID("exampleUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.Qual("github.com/stretchr/testify/assert", "Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("exampleWebhookID"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_GetWebhook", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to_user FROM webhooks WHERE belongs_to_user = %s AND id = %s", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Add(utils.FakeStringFunc()),
					jen.ID("Events").MapAssign().Index().ID("string").Values(jen.Lit("things")),
					jen.ID("DataTypes").MapAssign().Index().ID("string").Values(jen.Lit("things")),
					jen.ID("Topics").MapAssign().Index().ID("string").Values(jen.Lit("things")),
				),
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowFromWebhook").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhook").Call(utils.CtxVar(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Add(utils.FakeStringFunc()),
					jen.ID("Events").MapAssign().Index().ID("string").Values(jen.Lit("things")),
					jen.ID("DataTypes").MapAssign().Index().ID("string").Values(jen.Lit("things")),
					jen.ID("Topics").MapAssign().Index().ID("string").Values(jen.Lit("things")),
				),
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhook").Call(utils.CtxVar(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.Err()),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error from database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Add(utils.FakeStringFunc()),
				),
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhook").Call(utils.CtxVar(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with invalid response from database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Add(utils.FakeStringFunc()),
					jen.ID("Events").MapAssign().Index().ID("string").Values(jen.Lit("things")),
					jen.ID("DataTypes").MapAssign().Index().ID("string").Values(jen.Lit("things")),
					jen.ID("Topics").MapAssign().Index().ID("string").Values(jen.Lit("things")),
				),
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhook").Call(utils.CtxVar(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_buildGetWebhookCountQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.CreateDefaultQueryFilter(pkg),
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedArgCount").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedQuery").Assign().Litf("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL AND belongs_to_user = %s LIMIT 20", getIncIndex(dbvendor, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetWebhookCountQuery").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.Qual("github.com/stretchr/testify/assert", "Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_GetWebhookCount", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("SELECT COUNT(webhooks.id) FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_user = %s LIMIT 20", getIncIndex(dbvendor, 0)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhookCount").Call(utils.CtxVar(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error from database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhookCount").Call(utils.CtxVar(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Zero").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_buildGetAllWebhooksCountQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
				jen.Line(),
				jen.ID("actual").Assign().ID(dbfl).Dot("buildGetAllWebhooksCountQuery").Call(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_GetAllWebhooksCount", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooksCount").Call(utils.CtxVar()),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error from database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooksCount").Call(utils.CtxVar()),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Zero").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_buildGetAllWebhooksQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to_user FROM webhooks WHERE archived_on IS NULL"),
				jen.Line(),
				jen.ID("actual").Assign().ID(dbfl).Dot("buildGetAllWebhooksQuery").Call(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_GetAllWebhooks", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedListQuery").Assign().Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to_user FROM webhooks WHERE archived_on IS NULL"),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expectedCount").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedCountQuery").Assign().Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList").Valuesln(
					jen.ID("Pagination").MapAssign().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().Add(utils.FakeUint64Func()), jen.ID("TotalCount").MapAssign().ID("expectedCount")), jen.ID("Webhooks").MapAssign().Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
						jen.Valuesln(
							jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
							jen.ID("Name").MapAssign().Lit("name"),
						),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooks").Call(utils.CtxVar()),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooks").Call(utils.CtxVar()),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.Err()),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooks").Call(utils.CtxVar()),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error from database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("example").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Lit("name"),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.ID("example"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooks").Call(utils.CtxVar()),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching count"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expectedCount").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedCountQuery").Assign().Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList").Valuesln(
					jen.ID("Pagination").MapAssign().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Pagination").Valuesln(
						jen.ID("TotalCount").MapAssign().ID("expectedCount")),
					jen.ID("Webhooks").MapAssign().Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
						jen.Valuesln(
							jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
							jen.ID("Name").MapAssign().Lit("name"),
						),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooks").Call(utils.CtxVar()),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_GetAllWebhooksForUser", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedListQuery").Assign().Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to_user FROM webhooks WHERE archived_on IS NULL"),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleUser").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").MapAssign().Add(utils.FakeUint64Func())),
				jen.ID("expected").Assign().Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.Valuesln(
						jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
						jen.ID("Name").MapAssign().Lit("name"),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Index(jen.Lit(0))),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooksForUser").Call(utils.CtxVar(), jen.ID("exampleUser").Dot("ID")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleUser").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").MapAssign().Add(utils.FakeUint64Func())),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooksForUser").Call(utils.CtxVar(), jen.ID("exampleUser").Dot("ID")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.Err()),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleUser").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").MapAssign().Add(utils.FakeUint64Func())),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooksForUser").Call(utils.CtxVar(), jen.ID("exampleUser").Dot("ID")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous response from database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleUser").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").MapAssign().Add(utils.FakeUint64Func())),
				jen.ID("expected").Assign().Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.Valuesln(
						jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
						jen.ID("Name").MapAssign().Lit("name"),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Index(jen.Lit(0)))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllWebhooksForUser").Call(utils.CtxVar(), jen.ID("exampleUser").Dot("ID")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_buildGetWebhooksQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.List(jen.ID(dbfl), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("expectedArgCount").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedQuery").Assign().Litf("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to_user FROM webhooks WHERE archived_on IS NULL AND belongs_to_user = %s LIMIT 20", getIncIndex(dbvendor, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetWebhooksQuery").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.Qual("github.com/stretchr/testify/assert", "Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_GetWebhooks", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedListQuery").Assign().Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to_user FROM webhooks WHERE archived_on IS NULL"),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedCount").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedCountQuery").Assign().Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList").Valuesln(
					jen.ID("Pagination").MapAssign().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().Add(utils.FakeUint64Func()),
						jen.ID("Limit").MapAssign().Lit(20),
						jen.ID("TotalCount").MapAssign().ID("expectedCount"),
					),
					jen.ID("Webhooks").MapAssign().Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
						jen.Valuesln(
							jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
							jen.ID("Name").MapAssign().Lit("name"),
						),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhooks").Call(utils.CtxVar(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhooks").Call(utils.CtxVar(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.Err()),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhooks").Call(utils.CtxVar(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous response from database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Lit("name"),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhooks").Call(utils.CtxVar(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching count"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedCount").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedCountQuery").Assign().Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList").Valuesln(
					jen.ID("Pagination").MapAssign().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().Add(utils.FakeUint64Func()),
						jen.ID("Limit").MapAssign().Lit(20),
						jen.ID("TotalCount").MapAssign().ID("expectedCount"),
					),
					jen.ID("Webhooks").MapAssign().Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
						jen.Valuesln(
							jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
							jen.ID("Name").MapAssign().Lit("name"),
						),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.VarPointer().ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetWebhooks").Call(utils.CtxVar(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	/////////////

	var queryTail string
	if isPostgres {
		queryTail = " RETURNING id, created_on"
	}

	var (
		createdOnCol, createdOnVal string
	)

	if isMariaDB {
		createdOnCol = ",created_on"
		createdOnVal = ",UNIX_TIMESTAMP()"
	}

	ret.Add(
		jen.Func().IDf("Test%s_buildWebhookCreationQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleInput").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("Name").MapAssign().Lit("name"),
					jen.ID("ContentType").MapAssign().Lit("application/json"),
					jen.ID("URL").MapAssign().Lit("https://verygoodsoftwarenotvirus.ru"),
					jen.ID("Method").MapAssign().Qual("net/http", "MethodPatch"),
					jen.ID("Events").MapAssign().Index().ID("string").Values(),
					jen.ID("DataTypes").MapAssign().Index().ID("string").Values(),
					jen.ID("Topics").MapAssign().Index().ID("string").Values(),
					jen.ID("BelongsToUser").MapAssign().Add(utils.FakeUint64Func()),
				),
				jen.ID("expectedArgCount").Assign().Lit(8),
				jen.ID("expectedQuery").Assign().Litf("INSERT INTO webhooks (name,content_type,url,method,events,data_types,topics,belongs_to_user%s) VALUES (%s,%s,%s,%s,%s,%s,%s,%s%s)%s",
					createdOnCol,
					getIncIndex(dbvendor, 0),
					getIncIndex(dbvendor, 1),
					getIncIndex(dbvendor, 2),
					getIncIndex(dbvendor, 3),
					getIncIndex(dbvendor, 4),
					getIncIndex(dbvendor, 5),
					getIncIndex(dbvendor, 6),
					getIncIndex(dbvendor, 7),
					createdOnVal,
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildWebhookCreationQuery").Call(jen.ID("exampleInput")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.Qual("github.com/stretchr/testify/assert", "Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			)),
		),
		jen.Line(),
	)
	queryTail = ""

	/////////////

	if isPostgres {
		queryTail = " RETURNING id, created_on"
	}

	var createWebhookExpectFunc, createWebhookReturnFunc string

	buildCreateWebhookExampleRows := func() jen.Code {
		if isPostgres {
			createWebhookExpectFunc = "ExpectQuery"
			createWebhookReturnFunc = "WillReturnRows"
			return jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("ID"), jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
		} else if isSqlite || isMariaDB {
			createWebhookExpectFunc = "ExpectExec"
			createWebhookReturnFunc = "WillReturnResult"
			return jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.Add(utils.FakeUint64Func()))
		}
		return jen.Null()
	}

	ret.Add(
		jen.Func().IDf("Test%s_CreateWebhook", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("INSERT INTO webhooks (name,content_type,url,method,events,data_types,topics,belongs_to_user%s) VALUES (%s,%s,%s,%s,%s,%s,%s,%s%s)%s",
				createdOnCol,
				getIncIndex(dbvendor, 0),
				getIncIndex(dbvendor, 1),
				getIncIndex(dbvendor, 2),
				getIncIndex(dbvendor, 3),
				getIncIndex(dbvendor, 4),
				getIncIndex(dbvendor, 5),
				getIncIndex(dbvendor, 6),
				getIncIndex(dbvendor, 7),
				createdOnVal,
				queryTail,
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				func() []jen.Code {
					out := []jen.Code{
						jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
						jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
							jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
							jen.ID("Name").MapAssign().Lit("name"),
							jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
							jen.ID("CreatedOn").MapAssign().ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
						),
						jen.ID("expectedInput").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput").Valuesln(
							jen.ID("Name").MapAssign().ID("expected").Dot("Name"),
							jen.ID("BelongsToUser").MapAssign().ID("expected").Dot("BelongsToUser"),
						),
						buildCreateWebhookExampleRows(),
						jen.Line(),
						jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
						jen.ID("mockDB").Dot(createWebhookExpectFunc).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
							jen.ID("expected").Dot("Name"),
							jen.ID("expected").Dot("ContentType"),
							jen.ID("expected").Dot("URL"),
							jen.ID("expected").Dot("Method"),
							jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Events"), jen.ID("eventsSeparator")),
							jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("DataTypes"), jen.ID("typesSeparator")),
							jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Topics"), jen.ID("topicsSeparator")),
							jen.ID("expected").Dot("BelongsToUser"),
						).Dot(createWebhookReturnFunc).Call(jen.ID("exampleRows")),
						jen.Line(),
					}

					if isSqlite || isMariaDB {
						out = append(out,
							jen.ID("expectedTimeQuery").Assign().Litf("SELECT created_on FROM webhooks WHERE id = %s", getIncIndex(dbvendor, 0)),
							jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).
								Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("CreatedOn"))),
							jen.Line(),
						)
					}

					out = append(out,
						jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateWebhook").Call(utils.CtxVar(), jen.ID("expectedInput")),
						jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
						jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
						jen.Line(),
						jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
					)

					return out
				}()...,
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error interacting with database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Lit("name"),
					jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
					jen.ID("CreatedOn").MapAssign().ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedInput").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput").Valuesln(
					jen.ID("Name").MapAssign().ID("expected").Dot("Name"),
					jen.ID("BelongsToUser").MapAssign().ID("expected").Dot("BelongsToUser")),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(createWebhookExpectFunc).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Name"),
					jen.ID("expected").Dot("ContentType"),
					jen.ID("expected").Dot("URL"),
					jen.ID("expected").Dot("Method"),
					jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Events"), jen.ID("eventsSeparator")),
					jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("DataTypes"), jen.ID("typesSeparator")),
					jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Topics"), jen.ID("topicsSeparator")),
					jen.ID("expected").Dot("BelongsToUser"),
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateWebhook").Call(utils.CtxVar(), jen.ID("expectedInput")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)
	queryTail = ""

	/////////////

	if isPostgres {
		queryTail = " RETURNING updated_on"
	}

	ret.Add(
		jen.Func().IDf("Test%s_buildUpdateWebhookQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleInput").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("Name").MapAssign().Lit("name"),
					jen.ID("ContentType").MapAssign().Lit("application/json"),
					jen.ID("URL").MapAssign().Lit("https://verygoodsoftwarenotvirus.ru"),
					jen.ID("Method").MapAssign().Qual("net/http", "MethodPatch"),
					jen.ID("Events").MapAssign().Index().ID("string").Values(),
					jen.ID("DataTypes").MapAssign().Index().ID("string").Values(),
					jen.ID("Topics").MapAssign().Index().ID("string").Values(),
					jen.ID("BelongsToUser").MapAssign().Add(utils.FakeUint64Func()),
				),
				jen.ID("expectedArgCount").Assign().Lit(9),
				jen.ID("expectedQuery").Assign().Litf("UPDATE webhooks SET name = %s, content_type = %s, url = %s, method = %s, events = %s, data_types = %s, topics = %s, updated_on = %s WHERE belongs_to_user = %s AND id = %s%s",
					getIncIndex(dbvendor, 0),
					getIncIndex(dbvendor, 1),
					getIncIndex(dbvendor, 2),
					getIncIndex(dbvendor, 3),
					getIncIndex(dbvendor, 4),
					getIncIndex(dbvendor, 5),
					getIncIndex(dbvendor, 6),
					getTimeQuery(dbvendor),
					getIncIndex(dbvendor, 7),
					getIncIndex(dbvendor, 8),
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildUpdateWebhookQuery").Call(jen.ID("exampleInput")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.Qual("github.com/stretchr/testify/assert", "Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			)),
		),
		jen.Line(),
	)
	queryTail = ""

	/////////////

	if isPostgres {
		queryTail = " RETURNING updated_on"
	}

	var updateWebhookExpectFunc, updateWebhookReturnFunc string

	buildUpdateWebhookExampleRows := func() jen.Code {
		if isPostgres {
			updateWebhookExpectFunc = "ExpectQuery"
			updateWebhookReturnFunc = "WillReturnRows"
			return jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("updated_on"))).Dot("AddRow").Call(jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
		} else if isSqlite || isMariaDB {
			updateWebhookExpectFunc = "ExpectExec"
			updateWebhookReturnFunc = "WillReturnResult"
			return jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.Add(utils.FakeUint64Func()))
		}
		return jen.Null()
	}

	ret.Add(
		jen.Func().IDf("Test%s_UpdateWebhook", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("UPDATE webhooks SET name = %s, content_type = %s, url = %s, method = %s, events = %s, data_types = %s, topics = %s, updated_on = %s WHERE belongs_to_user = %s AND id = %s%s",
				getIncIndex(dbvendor, 0),
				getIncIndex(dbvendor, 1),
				getIncIndex(dbvendor, 2),
				getIncIndex(dbvendor, 3),
				getIncIndex(dbvendor, 4),
				getIncIndex(dbvendor, 5),
				getIncIndex(dbvendor, 6),
				getTimeQuery(dbvendor),
				getIncIndex(dbvendor, 7),
				getIncIndex(dbvendor, 8),
				queryTail,
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("Name").MapAssign().Lit("name"),
					jen.ID("ContentType").MapAssign().Lit("application/json"),
					jen.ID("URL").MapAssign().Lit("https://verygoodsoftwarenotvirus.ru"),
					jen.ID("Method").MapAssign().Qual("net/http", "MethodPatch"),
					jen.ID("Events").MapAssign().Index().ID("string").Values(),
					jen.ID("DataTypes").MapAssign().Index().ID("string").Values(),
					jen.ID("Topics").MapAssign().Index().ID("string").Values(),
					jen.ID("BelongsToUser").MapAssign().Add(utils.FakeUint64Func()),
				),
				buildUpdateWebhookExampleRows(),
				jen.Line(),
				jen.ID("mockDB").Dot(updateWebhookExpectFunc).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Name"),
					jen.ID("expected").Dot("ContentType"),
					jen.ID("expected").Dot("URL"),
					jen.ID("expected").Dot("Method"),
					jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Events"),
						jen.ID("eventsSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("DataTypes"),
						jen.ID("typesSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Topics"),
						jen.ID("topicsSeparator")), jen.ID("expected").Dot("BelongsToUser"),
					jen.ID("expected").Dot("ID"),
				).Dot(updateWebhookReturnFunc).Call(jen.ID("exampleRows")),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("UpdateWebhook").Call(utils.CtxVar(), jen.ID("expected")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error from database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("Name").MapAssign().Lit("name"), jen.ID("ContentType").MapAssign().Lit("application/json"), jen.ID("URL").MapAssign().Lit("https://verygoodsoftwarenotvirus.ru"), jen.ID("Method").MapAssign().Qual("net/http", "MethodPatch"), jen.ID("Events").MapAssign().Index().ID("string").Values(), jen.ID("DataTypes").MapAssign().Index().ID("string").Values(), jen.ID("Topics").MapAssign().Index().ID("string").Values(), jen.ID("BelongsToUser").MapAssign().Add(utils.FakeUint64Func())),
				jen.Line(),
				jen.ID("mockDB").Dot(updateWebhookExpectFunc).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Name"),
					jen.ID("expected").Dot("ContentType"),
					jen.ID("expected").Dot("URL"),
					jen.ID("expected").Dot("Method"),
					jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Events"),
						jen.ID("eventsSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("DataTypes"),
						jen.ID("typesSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Topics"),
						jen.ID("topicsSeparator")), jen.ID("expected").Dot("BelongsToUser"),
					jen.ID("expected").Dot("ID"),
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("UpdateWebhook").Call(utils.CtxVar(), jen.ID("expected")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)
	queryTail = ""

	/////////////

	if isPostgres {
		queryTail = " RETURNING archived_on"
	}

	ret.Add(
		jen.Func().IDf("Test%s_buildArchiveWebhookQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleWebhookID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedArgCount").Assign().Lit(2),
				jen.ID("expectedQuery").Assign().Litf("UPDATE webhooks SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND belongs_to_user = %s AND id = %s%s",
					getTimeQuery(dbvendor),
					getTimeQuery(dbvendor),
					getIncIndex(dbvendor, 0),
					getIncIndex(dbvendor, 1),
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildArchiveWebhookQuery").Call(jen.ID("exampleWebhookID"), jen.ID("exampleUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.Qual("github.com/stretchr/testify/assert", "Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("exampleWebhookID"), jen.ID("args").Index(jen.Add(utils.FakeUint64Func())).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)
	queryTail = ""

	/////////////

	if isPostgres {
		queryTail = " RETURNING archived_on"
	}

	ret.Add(
		jen.Func().IDf("Test%s_ArchiveWebhook", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Lit("name"),
					jen.ID("BelongsToUser").MapAssign().Lit(321),
					jen.ID("CreatedOn").MapAssign().ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedQuery").Assign().Litf("UPDATE webhooks SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND belongs_to_user = %s AND id = %s%s",
					getTimeQuery(dbvendor),
					getTimeQuery(dbvendor),
					getIncIndex(dbvendor, 0),
					getIncIndex(dbvendor, 1),
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("BelongsToUser"),
					jen.ID("expected").Dot("ID"),
				).Dot("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.Add(utils.FakeUint64Func()), jen.Add(utils.FakeUint64Func()))),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("ArchiveWebhook").Call(utils.CtxVar(), jen.ID("expected").Dot("ID"), jen.ID("expected").Dot("BelongsToUser")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)
	return ret
}
