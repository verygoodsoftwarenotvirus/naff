package mariadb

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func webhooksTestDotGo() *jen.File {
	ret := jen.NewFile("mariadb")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("buildMockRowFromWebhook").Params(jen.ID("w").Op("*").ID("models").Dot(
		"Webhook",
	)).Params(jen.Op("*").ID("sqlmock").Dot(
		"Rows",
	)).Block(
		jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot(
			"NewRows",
		).Call(jen.ID("webhooksTableColumns")).Dot(
			"AddRow",
		).Call(jen.ID("w").Dot(
			"ID",
		), jen.ID("w").Dot(
			"Name",
		), jen.ID("w").Dot(
			"ContentType",
		), jen.ID("w").Dot(
			"URL",
		), jen.ID("w").Dot(
			"Method",
		), jen.Qual("strings", "Join").Call(jen.ID("w").Dot(
			"Events",
		), jen.ID("eventsSeparator")), jen.Qual("strings", "Join").Call(jen.ID("w").Dot(
			"DataTypes",
		), jen.ID("typesSeparator")), jen.Qual("strings", "Join").Call(jen.ID("w").Dot(
			"Topics",
		), jen.ID("topicsSeparator")), jen.ID("w").Dot(
			"CreatedOn",
		), jen.ID("w").Dot(
			"UpdatedOn",
		), jen.ID("w").Dot(
			"ArchivedOn",
		), jen.ID("w").Dot(
			"BelongsTo",
		)),
		jen.Return().ID("exampleRows"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("buildErroneousMockRowFromWebhook").Params(jen.ID("w").Op("*").ID("models").Dot(
		"Webhook",
	)).Params(jen.Op("*").ID("sqlmock").Dot(
		"Rows",
	)).Block(
		jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot(
			"NewRows",
		).Call(jen.ID("webhooksTableColumns")).Dot(
			"AddRow",
		).Call(jen.ID("w").Dot(
			"ArchivedOn",
		), jen.ID("w").Dot(
			"BelongsTo",
		), jen.ID("w").Dot(
			"Name",
		), jen.ID("w").Dot(
			"ContentType",
		), jen.ID("w").Dot(
			"URL",
		), jen.ID("w").Dot(
			"Method",
		), jen.Qual("strings", "Join").Call(jen.ID("w").Dot(
			"Events",
		), jen.ID("eventsSeparator")), jen.Qual("strings", "Join").Call(jen.ID("w").Dot(
			"DataTypes",
		), jen.ID("typesSeparator")), jen.Qual("strings", "Join").Call(jen.ID("w").Dot(
			"Topics",
		), jen.ID("topicsSeparator")), jen.ID("w").Dot(
			"CreatedOn",
		), jen.ID("w").Dot(
			"UpdatedOn",
		), jen.ID("w").Dot(
			"ID",
		)),
		jen.Return().ID("exampleRows"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_buildGetWebhookQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleWebhookID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedArgCount").Op(":=").Lit(2),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE belongs_to = ? AND id = ?"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildGetWebhookQuery",
			).Call(jen.ID("exampleWebhookID"), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("exampleWebhookID"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_GetWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE belongs_to = ? AND id = ?"),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("Events").Op(":").Index().ID("string").Valuesln(jen.Lit("things")), jen.ID("DataTypes").Op(":").Index().ID("string").Valuesln(jen.Lit("things")), jen.ID("Topics").Op(":").Index().ID("string").Valuesln(jen.Lit("things"))),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID"), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromWebhook").Call(jen.ID("expected"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetWebhook",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
			), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE belongs_to = ? AND id = ?"),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("Events").Op(":").Index().ID("string").Valuesln(jen.Lit("things")), jen.ID("DataTypes").Op(":").Index().ID("string").Valuesln(jen.Lit("things")), jen.ID("Topics").Op(":").Index().ID("string").Valuesln(jen.Lit("things"))),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID"), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnError",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetWebhook",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
			), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE belongs_to = ? AND id = ?"),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID"), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetWebhook",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
			), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with invalid response from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE belongs_to = ? AND id = ?"),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("Events").Op(":").Index().ID("string").Valuesln(jen.Lit("things")), jen.ID("DataTypes").Op(":").Index().ID("string").Valuesln(jen.Lit("things")), jen.ID("Topics").Op(":").Index().ID("string").Valuesln(jen.Lit("things"))),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID"), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.ID("expected"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetWebhook",
			).Call(jen.ID("ctx"), jen.ID("expected").Dot(
				"ID",
			), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_buildGetWebhookCountQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedArgCount").Op(":=").Lit(1),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildGetWebhookCountQuery",
			).Call(jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_GetWebhookCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnRows",
			).Call(jen.ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(jen.Lit("count"))).Dot(
				"AddRow",
			).Call(jen.ID("expected"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetWebhookCount",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetWebhookCount",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Zero",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_buildGetAllWebhooksCountQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
			jen.ID("actualQuery").Op(":=").ID("m").Dot(
				"buildGetAllWebhooksCountQuery",
			).Call(),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_GetAllWebhooksCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(jen.Lit("count"))).Dot(
				"AddRow",
			).Call(jen.ID("expected"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllWebhooksCount",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllWebhooksCount",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Zero",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_buildGetAllWebhooksQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
			jen.ID("actual").Op(":=").ID("m").Dot(
				"buildGetAllWebhooksQuery",
			).Call(),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_GetAllWebhooks").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
			jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"WebhookList",
			).Valuesln(jen.ID("Pagination").Op(":").ID("models").Dot(
				"Pagination",
			).Valuesln(jen.ID("Page").Op(":").Lit(1), jen.ID("TotalCount").Op(":").ID("expectedCount")), jen.ID("Webhooks").Op(":").Index().ID("models").Dot(
				"Webhook",
			).Valuesln(jen.Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")))),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot(
				"Webhooks",
			).Index(jen.Lit(0))), jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot(
				"Webhooks",
			).Index(jen.Lit(0))), jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot(
				"Webhooks",
			).Index(jen.Lit(0)))),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(jen.Lit("count"))).Dot(
				"AddRow",
			).Call(jen.ID("expectedCount"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllWebhooks",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllWebhooks",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllWebhooks",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
			jen.ID("example").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.ID("example"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllWebhooks",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error fetching count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
			jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"WebhookList",
			).Valuesln(jen.ID("Pagination").Op(":").ID("models").Dot(
				"Pagination",
			).Valuesln(jen.ID("TotalCount").Op(":").ID("expectedCount")), jen.ID("Webhooks").Op(":").Index().ID("models").Dot(
				"Webhook",
			).Valuesln(jen.Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")))),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot(
				"Webhooks",
			).Index(jen.Lit(0))), jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot(
				"Webhooks",
			).Index(jen.Lit(0))), jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot(
				"Webhooks",
			).Index(jen.Lit(0)))),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllWebhooks",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_GetAllWebhooksForUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(123)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
			jen.ID("expected").Op(":=").Index().ID("models").Dot(
				"Webhook",
			).Valuesln(jen.Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"))),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Index(jen.Lit(0))), jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Index(jen.Lit(0))), jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Index(jen.Lit(0)))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllWebhooksForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot(
				"ID",
			)),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(123)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllWebhooksForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot(
				"ID",
			)),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(123)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllWebhooksForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot(
				"ID",
			)),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with erroneous response from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(123)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
			jen.ID("expected").Op(":=").Index().ID("models").Dot(
				"Webhook",
			).Valuesln(jen.Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"))),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.Op("&").ID("expected").Index(jen.Lit(0)))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllWebhooksForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot(
				"ID",
			)),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_buildGetWebhooksQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedArgCount").Op(":=").Lit(1),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildGetWebhooksQuery",
			).Call(jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_GetWebhooks").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
			jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"WebhookList",
			).Valuesln(jen.ID("Pagination").Op(":").ID("models").Dot(
				"Pagination",
			).Valuesln(jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("TotalCount").Op(":").ID("expectedCount")), jen.ID("Webhooks").Op(":").Index().ID("models").Dot(
				"Webhook",
			).Valuesln(jen.Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")))),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot(
				"Webhooks",
			).Index(jen.Lit(0))), jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot(
				"Webhooks",
			).Index(jen.Lit(0))), jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot(
				"Webhooks",
			).Index(jen.Lit(0)))),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(jen.Lit("count"))).Dot(
				"AddRow",
			).Call(jen.ID("expectedCount"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetWebhooks",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetWebhooks",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetWebhooks",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with erroneous response from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildErroneousMockRowFromWebhook").Call(jen.ID("expected"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetWebhooks",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error fetching count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, content_type, url, method, events, data_types, topics, created_on, updated_on, archived_on, belongs_to FROM webhooks WHERE archived_on IS NULL"),
			jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM webhooks WHERE archived_on IS NULL"),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"WebhookList",
			).Valuesln(jen.ID("Pagination").Op(":").ID("models").Dot(
				"Pagination",
			).Valuesln(jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("TotalCount").Op(":").ID("expectedCount")), jen.ID("Webhooks").Op(":").Index().ID("models").Dot(
				"Webhook",
			).Valuesln(jen.Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name")))),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot(
				"Webhooks",
			).Index(jen.Lit(0))), jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot(
				"Webhooks",
			).Index(jen.Lit(0))), jen.ID("buildMockRowFromWebhook").Call(jen.Op("&").ID("expected").Dot(
				"Webhooks",
			).Index(jen.Lit(0)))),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetWebhooks",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_buildWebhookCreationQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("Name").Op(":").Lit("name"), jen.ID("ContentType").Op(":").Lit("application/json"), jen.ID("URL").Op(":").Lit("https://verygoodsoftwarenotvirus.ru"), jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"), jen.ID("Events").Op(":").Index().ID("string").Valuesln(), jen.ID("DataTypes").Op(":").Index().ID("string").Valuesln(), jen.ID("Topics").Op(":").Index().ID("string").Valuesln(), jen.ID("BelongsTo").Op(":").Lit(1)),
			jen.ID("expectedArgCount").Op(":=").Lit(8),
			jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO webhooks (name,content_type,url,method,events,data_types,topics,belongs_to,created_on) VALUES (?,?,?,?,?,?,?,?,UNIX_TIMESTAMP())"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildWebhookCreationQuery",
			).Call(jen.ID("exampleInput")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_CreateWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedInput").Op(":=").Op("&").ID("models").Dot(
				"WebhookCreationInput",
			).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
				"Name",
			), jen.ID("BelongsTo").Op(":").ID("expected").Dot(
				"BelongsTo",
			)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedCreationQuery").Op(":=").Lit("INSERT INTO webhooks (name,content_type,url,method,events,data_types,topics,belongs_to,created_on) VALUES (?,?,?,?,?,?,?,?,UNIX_TIMESTAMP())"),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCreationQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"Name",
			), jen.ID("expected").Dot(
				"ContentType",
			), jen.ID("expected").Dot(
				"URL",
			), jen.ID("expected").Dot(
				"Method",
			), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot(
				"Events",
			), jen.ID("eventsSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot(
				"DataTypes",
			), jen.ID("typesSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot(
				"Topics",
			), jen.ID("topicsSeparator")), jen.ID("expected").Dot(
				"BelongsTo",
			)).Dot(
				"WillReturnResult",
			).Call(jen.ID("sqlmock").Dot(
				"NewResult",
			).Call(jen.ID("int64").Call(jen.ID("expected").Dot(
				"ID",
			)), jen.Lit(1))),
			jen.ID("expectedTimeQuery").Op(":=").Lit("SELECT created_on FROM webhooks WHERE id = ?"),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(jen.Lit("created_on"))).Dot(
				"AddRow",
			).Call(jen.ID("expected").Dot(
				"CreatedOn",
			))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"CreateWebhook",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error interacting with database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedInput").Op(":=").Op("&").ID("models").Dot(
				"WebhookCreationInput",
			).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
				"Name",
			), jen.ID("BelongsTo").Op(":").ID("expected").Dot(
				"BelongsTo",
			)),
			jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO webhooks (name,content_type,url,method,events,data_types,topics,belongs_to,created_on) VALUES (?,?,?,?,?,?,?,?,UNIX_TIMESTAMP())"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"Name",
			), jen.ID("expected").Dot(
				"ContentType",
			), jen.ID("expected").Dot(
				"URL",
			), jen.ID("expected").Dot(
				"Method",
			), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot(
				"Events",
			), jen.ID("eventsSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot(
				"DataTypes",
			), jen.ID("typesSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot(
				"Topics",
			), jen.ID("topicsSeparator")), jen.ID("expected").Dot(
				"BelongsTo",
			)).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"CreateWebhook",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_buildUpdateWebhookQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("Name").Op(":").Lit("name"), jen.ID("ContentType").Op(":").Lit("application/json"), jen.ID("URL").Op(":").Lit("https://verygoodsoftwarenotvirus.ru"), jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"), jen.ID("Events").Op(":").Index().ID("string").Valuesln(), jen.ID("DataTypes").Op(":").Index().ID("string").Valuesln(), jen.ID("Topics").Op(":").Index().ID("string").Valuesln(), jen.ID("BelongsTo").Op(":").Lit(1)),
			jen.ID("expectedArgCount").Op(":=").Lit(9),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE webhooks SET name = ?, content_type = ?, url = ?, method = ?, events = ?, data_types = ?, topics = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildUpdateWebhookQuery",
			).Call(jen.ID("exampleInput")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_UpdateWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("Name").Op(":").Lit("name"), jen.ID("ContentType").Op(":").Lit("application/json"), jen.ID("URL").Op(":").Lit("https://verygoodsoftwarenotvirus.ru"), jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"), jen.ID("Events").Op(":").Index().ID("string").Valuesln(), jen.ID("DataTypes").Op(":").Index().ID("string").Valuesln(), jen.ID("Topics").Op(":").Index().ID("string").Valuesln(), jen.ID("BelongsTo").Op(":").Lit(1)),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE webhooks SET name = ?, content_type = ?, url = ?, method = ?, events = ?, data_types = ?, topics = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"Name",
			), jen.ID("expected").Dot(
				"ContentType",
			), jen.ID("expected").Dot(
				"URL",
			), jen.ID("expected").Dot(
				"Method",
			), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot(
				"Events",
			), jen.ID("eventsSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot(
				"DataTypes",
			), jen.ID("typesSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot(
				"Topics",
			), jen.ID("topicsSeparator")), jen.ID("expected").Dot(
				"BelongsTo",
			), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnResult",
			).Call(jen.ID("sqlmock").Dot(
				"NewResult",
			).Call(jen.Lit(1), jen.Lit(1))),
			jen.ID("err").Op(":=").ID("m").Dot(
				"UpdateWebhook",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("Name").Op(":").Lit("name"), jen.ID("ContentType").Op(":").Lit("application/json"), jen.ID("URL").Op(":").Lit("https://verygoodsoftwarenotvirus.ru"), jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"), jen.ID("Events").Op(":").Index().ID("string").Valuesln(), jen.ID("DataTypes").Op(":").Index().ID("string").Valuesln(), jen.ID("Topics").Op(":").Index().ID("string").Valuesln(), jen.ID("BelongsTo").Op(":").Lit(1)),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE webhooks SET name = ?, content_type = ?, url = ?, method = ?, events = ?, data_types = ?, topics = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"Name",
			), jen.ID("expected").Dot(
				"ContentType",
			), jen.ID("expected").Dot(
				"URL",
			), jen.ID("expected").Dot(
				"Method",
			), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot(
				"Events",
			), jen.ID("eventsSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot(
				"DataTypes",
			), jen.ID("typesSeparator")), jen.Qual("strings", "Join").Call(jen.ID("expected").Dot(
				"Topics",
			), jen.ID("topicsSeparator")), jen.ID("expected").Dot(
				"BelongsTo",
			), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("err").Op(":=").ID("m").Dot(
				"UpdateWebhook",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_buildArchiveWebhookQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleWebhookID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedArgCount").Op(":=").Lit(2),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE webhooks SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildArchiveWebhookQuery",
			).Call(jen.ID("exampleWebhookID"), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("exampleWebhookID"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestMariaDB_ArchiveWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Webhook",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").Lit(321), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE webhooks SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"BelongsTo",
			), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnResult",
			).Call(jen.ID("sqlmock").Dot(
				"NewResult",
			).Call(jen.Lit(1), jen.Lit(1))),
			jen.ID("err").Op(":=").ID("m").Dot(
				"ArchiveWebhook",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
			), jen.ID("expected").Dot(
				"BelongsTo",
			)),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),

		jen.Line(),
	)
	return ret
}
