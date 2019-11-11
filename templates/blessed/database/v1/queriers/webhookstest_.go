package queriers

import (
	"path/filepath"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksTestDotGo(pkg *models.Project, vendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(vendor.SingularPackageName())

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)
	sn := vendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))
	dbrn := vendor.RouteName()

	isPostgres := dbrn == "postgres"
	isSqlite := dbrn == "sqlite"
	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	/////////////

	ret.Add(
		jen.Func().ID("buildMockRowFromWebhook").Params(jen.ID("w").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook")).Params(jen.Op("*").Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("webhooksTableColumns")).Dot("AddRow").Callln(
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
				jen.ID("w").Dot("BelongsTo"),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().ID("buildErroneousMockRowFromWebhook").Params(jen.ID("w").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook")).Params(jen.Op("*").Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("webhooksTableColumns")).Dot("AddRow").Callln(
				jen.ID("w").Dot("ArchivedOn"),
				jen.ID("w").Dot("BelongsTo"),
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
		jen.Func().IDf("Test%s_buildGetWebhookQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleWebhookID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(2),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE belongs_to = %s AND id = %s", getIncIndex(dbrn, 0), getIncIndex(dbrn, 1)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetWebhookQuery").Call(jen.ID("exampleWebhookID"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleWebhookID"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_GetWebhook", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE belongs_to = %s AND id = %s", getIncIndex(dbrn, 0), getIncIndex(dbrn, 1)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("Events").Op(":").Index().ID("string").Values(jen.Lit("things")),
					jen.ID("DataTypes").Op(":").Index().ID("string").Values(jen.Lit("things")),
					jen.ID("Topics").Op(":").Index().ID("string").Values(jen.Lit("things")),
				),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowFromWebhook").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetWebhook").Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE belongs_to = %s AND id = %s", getIncIndex(dbrn, 0), getIncIndex(dbrn, 1)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"),
					jen.ID("Events").Op(":").Index().ID("string").Values(jen.Lit("things")),
					jen.ID("DataTypes").Op(":").Index().ID("string").Values(jen.Lit("things")),
					jen.ID("Topics").Op(":").Index().ID("string").Values(jen.Lit("things")),
				),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetWebhook").Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE belongs_to = %s AND id = %s", getIncIndex(dbrn, 0), getIncIndex(dbrn, 1)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
				),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetWebhook").Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with invalid response from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE belongs_to = %s AND id = %s", getIncIndex(dbrn, 0), getIncIndex(dbrn, 1)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("Events").Op(":").Index().ID("string").Values(jen.Lit("things")),
					jen.ID("DataTypes").Op(":").Index().ID("string").Values(jen.Lit("things")),
					jen.ID("Topics").Op(":").Index().ID("string").Values(jen.Lit("things")),
				),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetWebhook").Call(jen.ID("ctx"), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_buildGetWebhookCountQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetWebhookCountQuery").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_GetWebhookCount", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", getIncIndex(dbrn, 0)),
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetWebhookCount").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", getIncIndex(dbrn, 0)),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetWebhookCount").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Zero").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_buildGetAllWebhooksCountQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
				jen.Line(),
				jen.ID("actual").Op(":=").ID(dbfl).Dot("buildGetAllWebhooksCountQuery").Call(),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_GetAllWebhooksCount", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllWebhooksCount").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllWebhooksCount").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Zero").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_buildGetAllWebhooksQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
				jen.Line(),
				jen.ID("actual").Op(":=").ID(dbfl).Dot("buildGetAllWebhooksQuery").Call(),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_GetAllWebhooks", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList").Valuesln(
					jen.ID("Pagination").Op(":").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Pagination").Valuesln(
						jen.ID("Page").Op(":").Lit(1), jen.ID("TotalCount").Op(":").ID("expectedCount")), jen.ID("Webhooks").Op(":").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(123),
							jen.ID("Name").Op(":").Lit("name"),
						),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllWebhooks").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllWebhooks").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllWebhooks").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("example").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.ID("example"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllWebhooks").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList").Valuesln(
					jen.ID("Pagination").Op(":").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Pagination").Valuesln(
						jen.ID("TotalCount").Op(":").ID("expectedCount")),
					jen.ID("Webhooks").Op(":").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(123),
							jen.ID("Name").Op(":").Lit("name"),
						),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllWebhooks").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_GetAllWebhooksForUser", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expected").Op(":=").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.Valuesln(
						jen.ID("ID").Op(":").Lit(123),
						jen.ID("Name").Op(":").Lit("name"),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Index(jen.Lit(0))),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllWebhooksForUser").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot("ID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllWebhooksForUser").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot("ID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllWebhooksForUser").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot("ID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous response from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expected").Op(":=").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.Valuesln(
						jen.ID("ID").Op(":").Lit(123),
						jen.ID("Name").Op(":").Lit("name"),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.Op("&").ID("expected").Index(jen.Lit(0)))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllWebhooksForUser").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot("ID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_buildGetWebhooksQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetWebhooksQuery").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	/////////////

	ret.Add(
		jen.Func().IDf("Test%s_GetWebhooks", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList").Valuesln(
					jen.ID("Pagination").Op(":").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Pagination").Valuesln(
						jen.ID("Page").Op(":").Lit(1),
						jen.ID("Limit").Op(":").Lit(20),
						jen.ID("TotalCount").Op(":").ID("expectedCount"),
					),
					jen.ID("Webhooks").Op(":").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(123),
							jen.ID("Name").Op(":").Lit("name"),
						),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetWebhooks").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetWebhooks").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetWebhooks").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous response from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetWebhooks").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList").Valuesln(
					jen.ID("Pagination").Op(":").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Pagination").Valuesln(
						jen.ID("Page").Op(":").Lit(1),
						jen.ID("Limit").Op(":").Lit(20),
						jen.ID("TotalCount").Op(":").ID("expectedCount"),
					),
					jen.ID("Webhooks").Op(":").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(123),
							jen.ID("Name").Op(":").Lit("name"),
						),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot("Webhooks").Index(jen.Lit(0))),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetWebhooks").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
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
		jen.Func().IDf("Test%s_buildWebhookCreationQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("ContentType").Op(":").Lit("application/json"),
					jen.ID("URL").Op(":").Lit("https://verygoodsoftwarenotvirus.ru"),
					jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"),
					jen.ID("Events").Op(":").Index().ID("string").Values(),
					jen.ID("DataTypes").Op(":").Index().ID("string").Values(),
					jen.ID("Topics").Op(":").Index().ID("string").Values(),
					jen.ID("BelongsTo").Op(":").Lit(1),
				),
				jen.ID("expectedArgCount").Op(":=").Lit(8),
				jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO webhooks (name,content_type,url,method,events,data_types,topics,belongs_to%s) VALUES (%s,%s,%s,%s,%s,%s,%s,%s%s)%s",
					createdOnCol,
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					getIncIndex(dbrn, 2),
					getIncIndex(dbrn, 3),
					getIncIndex(dbrn, 4),
					getIncIndex(dbrn, 5),
					getIncIndex(dbrn, 6),
					getIncIndex(dbrn, 7),
					createdOnVal,
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildWebhookCreationQuery").Call(jen.ID("exampleInput")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
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
			return jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("ID"), jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
		} else if isSqlite || isMariaDB {
			createWebhookExpectFunc = "ExpectExec"
			createWebhookReturnFunc = "WillReturnResult"
			return jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.Lit(1))
		}
		return jen.Null()
	}

	buildCreateWebhookHappyPathBody := func() []jen.Code {
		out := []jen.Code{
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
				jen.ID("ID").Op(":").Lit(123),
				jen.ID("Name").Op(":").Lit("name"),
				jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
				jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
			),
			jen.ID("expectedInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput").Valuesln(
				jen.ID("Name").Op(":").ID("expected").Dot("Name"),
				jen.ID("BelongsTo").Op(":").ID("expected").Dot("BelongsTo"),
			),
			buildCreateWebhookExampleRows(),
			jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO webhooks (name,content_type,url,method,events,data_types,topics,belongs_to%s) VALUES (%s,%s,%s,%s,%s,%s,%s,%s%s)%s",
				createdOnCol,
				getIncIndex(dbrn, 0),
				getIncIndex(dbrn, 1),
				getIncIndex(dbrn, 2),
				getIncIndex(dbrn, 3),
				getIncIndex(dbrn, 4),
				getIncIndex(dbrn, 5),
				getIncIndex(dbrn, 6),
				getIncIndex(dbrn, 7),
				createdOnVal,
				queryTail,
			),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(createWebhookExpectFunc).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
				jen.ID("expected").Dot("Name"),
				jen.ID("expected").Dot("ContentType"),
				jen.ID("expected").Dot("URL"),
				jen.ID("expected").Dot("Method"),
				jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Events"), jen.ID("eventsSeparator")),
				jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("DataTypes"), jen.ID("typesSeparator")),
				jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Topics"), jen.ID("topicsSeparator")),
				jen.ID("expected").Dot("BelongsTo"),
			).Dot(createWebhookReturnFunc).Call(jen.ID("exampleRows")),
			jen.Line(),
		}

		if isSqlite || isMariaDB {
			out = append(out,
				jen.ID("expectedTimeQuery").Op(":=").Litf("SELECT created_on FROM webhooks WHERE id = %s", getIncIndex(dbrn, 0)),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("CreatedOn"))),
				jen.Line(),
			)
		}

		out = append(out,
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("CreateWebhook").Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.Line(),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return out
	}

	ret.Add(
		jen.Func().IDf("Test%s_CreateWebhook", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				buildCreateWebhookHappyPathBody()...,
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error interacting with database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput").Valuesln(
					jen.ID("Name").Op(":").ID("expected").Dot("Name"),
					jen.ID("BelongsTo").Op(":").ID("expected").Dot("BelongsTo")),
				jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO webhooks (name,content_type,url,method,events,data_types,topics,belongs_to%s) VALUES (%s,%s,%s,%s,%s,%s,%s,%s%s)%s",
					createdOnCol,
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					getIncIndex(dbrn, 2),
					getIncIndex(dbrn, 3),
					getIncIndex(dbrn, 4),
					getIncIndex(dbrn, 5),
					getIncIndex(dbrn, 6),
					getIncIndex(dbrn, 7),
					createdOnVal,
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(createWebhookExpectFunc).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Name"),
					jen.ID("expected").Dot("ContentType"),
					jen.ID("expected").Dot("URL"),
					jen.ID("expected").Dot("Method"),
					jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Events"), jen.ID("eventsSeparator")),
					jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("DataTypes"), jen.ID("typesSeparator")),
					jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Topics"), jen.ID("topicsSeparator")),
					jen.ID("expected").Dot("BelongsTo"),
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("CreateWebhook").Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
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
		jen.Func().IDf("Test%s_buildUpdateWebhookQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("ContentType").Op(":").Lit("application/json"),
					jen.ID("URL").Op(":").Lit("https://verygoodsoftwarenotvirus.ru"),
					jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"),
					jen.ID("Events").Op(":").Index().ID("string").Values(),
					jen.ID("DataTypes").Op(":").Index().ID("string").Values(),
					jen.ID("Topics").Op(":").Index().ID("string").Values(),
					jen.ID("BelongsTo").Op(":").Lit(1),
				),
				jen.ID("expectedArgCount").Op(":=").Lit(9),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE webhooks SET name = %s, content_type = %s, url = %s, method = %s, events = %s, data_types = %s, topics = %s, updated_on = %s WHERE belongs_to = %s AND id = %s%s",
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					getIncIndex(dbrn, 2),
					getIncIndex(dbrn, 3),
					getIncIndex(dbrn, 4),
					getIncIndex(dbrn, 5),
					getIncIndex(dbrn, 6),
					getTimeQuery(dbrn),
					getIncIndex(dbrn, 7),
					getIncIndex(dbrn, 8),
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildUpdateWebhookQuery").Call(jen.ID("exampleInput")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
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
			return jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("updated_on"))).Dot("AddRow").Call(jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
		} else if isSqlite || isMariaDB {
			updateWebhookExpectFunc = "ExpectExec"
			updateWebhookReturnFunc = "WillReturnResult"
			return jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.Lit(1))
		}
		return jen.Null()
	}

	ret.Add(
		jen.Func().IDf("Test%s_UpdateWebhook", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("ContentType").Op(":").Lit("application/json"),
					jen.ID("URL").Op(":").Lit("https://verygoodsoftwarenotvirus.ru"),
					jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"),
					jen.ID("Events").Op(":").Index().ID("string").Values(),
					jen.ID("DataTypes").Op(":").Index().ID("string").Values(),
					jen.ID("Topics").Op(":").Index().ID("string").Values(),
					jen.ID("BelongsTo").Op(":").Lit(1),
				),
				buildUpdateWebhookExampleRows(),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE webhooks SET name = %s, content_type = %s, url = %s, method = %s, events = %s, data_types = %s, topics = %s, updated_on = %s WHERE belongs_to = %s AND id = %s%s",
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					getIncIndex(dbrn, 2),
					getIncIndex(dbrn, 3),
					getIncIndex(dbrn, 4),
					getIncIndex(dbrn, 5),
					getIncIndex(dbrn, 6),
					getTimeQuery(dbrn),
					getIncIndex(dbrn, 7),
					getIncIndex(dbrn, 8),
					queryTail,
				),
				jen.Line(),
				jen.ID("mockDB").Dot(updateWebhookExpectFunc).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Name"),
					jen.ID("expected").Dot("ContentType"),
					jen.ID("expected").Dot("URL"),
					jen.ID("expected").Dot("Method"),
					jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Events"),
						jen.ID("eventsSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("DataTypes"),
						jen.ID("typesSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Topics"),
						jen.ID("topicsSeparator")), jen.ID("expected").Dot("BelongsTo"),
					jen.ID("expected").Dot("ID"),
				).Dot(updateWebhookReturnFunc).Call(jen.ID("exampleRows")),
				jen.Line(),
				jen.ID("err").Op(":=").ID(dbfl).Dot("UpdateWebhook").Call(jen.Qual("context", "Background").Call(), jen.ID("expected")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("Name").Op(":").Lit("name"), jen.ID("ContentType").Op(":").Lit("application/json"), jen.ID("URL").Op(":").Lit("https://verygoodsoftwarenotvirus.ru"), jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"), jen.ID("Events").Op(":").Index().ID("string").Values(), jen.ID("DataTypes").Op(":").Index().ID("string").Values(), jen.ID("Topics").Op(":").Index().ID("string").Values(), jen.ID("BelongsTo").Op(":").Lit(1)),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE webhooks SET name = %s, content_type = %s, url = %s, method = %s, events = %s, data_types = %s, topics = %s, updated_on = %s WHERE belongs_to = %s AND id = %s%s",
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					getIncIndex(dbrn, 2),
					getIncIndex(dbrn, 3),
					getIncIndex(dbrn, 4),
					getIncIndex(dbrn, 5),
					getIncIndex(dbrn, 6),
					getTimeQuery(dbrn),
					getIncIndex(dbrn, 7),
					getIncIndex(dbrn, 8),
					queryTail,
				),
				jen.Line(),
				jen.ID("mockDB").Dot(updateWebhookExpectFunc).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Name"),
					jen.ID("expected").Dot("ContentType"),
					jen.ID("expected").Dot("URL"),
					jen.ID("expected").Dot("Method"),
					jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Events"),
						jen.ID("eventsSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("DataTypes"),
						jen.ID("typesSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Topics"),
						jen.ID("topicsSeparator")), jen.ID("expected").Dot("BelongsTo"),
					jen.ID("expected").Dot("ID"),
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.ID("err").Op(":=").ID(dbfl).Dot("UpdateWebhook").Call(jen.Qual("context", "Background").Call(), jen.ID("expected")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
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
		jen.Func().IDf("Test%s_buildArchiveWebhookQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleWebhookID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedArgCount").Op(":=").Lit(2),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE webhooks SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND belongs_to = %s AND id = %s%s",
					getTimeQuery(dbrn),
					getTimeQuery(dbrn),
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildArchiveWebhookQuery").Call(jen.ID("exampleWebhookID"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleWebhookID"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
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
		jen.Func().IDf("Test%s_ArchiveWebhook", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").Lit(321),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE webhooks SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND belongs_to = %s AND id = %s%s",
					getTimeQuery(dbrn),
					getTimeQuery(dbrn),
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("BelongsTo"),
					jen.ID("expected").Dot("ID"),
				).Dot("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.Lit(1), jen.Lit(1))),
				jen.Line(),
				jen.ID("err").Op(":=").ID(dbfl).Dot("ArchiveWebhook").Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expected").Dot("BelongsTo")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)
	return ret
}
