package mariadb

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func oauth2ClientsTestDotGo() *jen.File {
	ret := jen.NewFile("mariadb")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("buildMockRowFromOAuth2Client").Params(jen.ID("c").Op("*").ID("models").Dot(
		"OAuth2Client",
	)).Params(jen.Op("*").ID("sqlmock").Dot(
		"Rows",
	)).Block(
		jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot(
			"NewRows",
		).Call(jen.ID("oauth2ClientsTableColumns")).Dot(
			"AddRow",
		).Call(jen.ID("c").Dot(
			"ID",
	),
	jen.ID("c").Dot(
			"Name",
	),
	jen.ID("c").Dot(
			"ClientID",
	),
	jen.Qual("strings", "Join").Call(jen.ID("c").Dot(
			"Scopes",
	),
	jen.ID("scopesSeparator")), jen.ID("c").Dot(
			"RedirectURI",
	),
	jen.ID("c").Dot(
			"ClientSecret",
	),
	jen.ID("c").Dot(
			"CreatedOn",
	),
	jen.ID("c").Dot(
			"UpdatedOn",
	),
	jen.ID("c").Dot(
			"ArchivedOn",
	),
	jen.ID("c").Dot(
			"BelongsTo",
		)),
		jen.Return().ID("exampleRows"),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildErroneousMockRowFromOAuth2Client").Params(jen.ID("c").Op("*").ID("models").Dot(
		"OAuth2Client",
	)).Params(jen.Op("*").ID("sqlmock").Dot(
		"Rows",
	)).Block(
		jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot(
			"NewRows",
		).Call(jen.ID("oauth2ClientsTableColumns")).Dot(
			"AddRow",
		).Call(jen.ID("c").Dot(
			"ArchivedOn",
	),
	jen.ID("c").Dot(
			"Name",
	),
	jen.ID("c").Dot(
			"ClientID",
	),
	jen.Qual("strings", "Join").Call(jen.ID("c").Dot(
			"Scopes",
	),
	jen.ID("scopesSeparator")), jen.ID("c").Dot(
			"RedirectURI",
	),
	jen.ID("c").Dot(
			"ClientSecret",
	),
	jen.ID("c").Dot(
			"CreatedOn",
	),
	jen.ID("c").Dot(
			"UpdatedOn",
	),
	jen.ID("c").Dot(
			"BelongsTo",
	),
	jen.ID("c").Dot(
			"ID",
		)),
		jen.Return().ID("exampleRows"),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildGetOAuth2ClientByClientIDQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedClientID").Op(":=").Lit("ClientID"),
			jen.ID("expectedArgCount").Op(":=").Lit(1),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND client_id = ?"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildGetOAuth2ClientByClientIDQuery",
			).Call(jen.ID("expectedClientID")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedClientID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("string"))),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_GetOAuth2ClientByClientID").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleClientID").Op(":=").Lit("EXAMPLE"),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND client_id = ?"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("exampleClientID")).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetOAuth2ClientByClientID",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleClientID").Op(":=").Lit("EXAMPLE"),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND client_id = ?"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("exampleClientID")).Dot(
				"WillReturnError",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetOAuth2ClientByClientID",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous row"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleClientID").Op(":=").Lit("EXAMPLE"),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND client_id = ?"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("exampleClientID")).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.ID("expected"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetOAuth2ClientByClientID",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildGetAllOAuth2ClientsQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.ID("actualQuery").Op(":=").ID("m").Dot(
				"buildGetAllOAuth2ClientsQuery",
			).Call(),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_GetAllOAuth2Clients").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Index().Op("*").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call()))),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Lit(0))), jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Lit(0))), jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Lit(0)))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllOAuth2Clients",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllOAuth2Clients",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error executing query"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllOAuth2Clients",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous response from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Index().Op("*").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call()))),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Lit(0)))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllOAuth2Clients",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_GetAllOAuth2ClientsForUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123)),
			jen.ID("expected").Op(":=").Index().Op("*").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call()))),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Lit(0))), jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Lit(0))), jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Lit(0)))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllOAuth2ClientsForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot(
				"ID",
			)),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123)),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllOAuth2ClientsForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot(
				"ID",
			)),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous response from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123)),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllOAuth2ClientsForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot(
				"ID",
			)),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with unscannable response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.ID("expected"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllOAuth2ClientsForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot(
				"ID",
			)),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildGetOAuth2ClientQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedClientID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedArgCount").Op(":=").Lit(2),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildGetOAuth2ClientQuery",
			).Call(jen.ID("expectedClientID"), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedClientID"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_GetOAuth2Client").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call()), jen.ID("Scopes").Op(":").Index().ID("string").Valuesln(
	jen.Lit("things"))),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"BelongsTo",
	),
	jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetOAuth2Client",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
	),
	jen.ID("expected").Dot(
				"BelongsTo",
			)),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call()), jen.ID("Scopes").Op(":").Index().ID("string").Valuesln(
	jen.Lit("things"))),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"BelongsTo",
	),
	jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnError",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetOAuth2Client",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
	),
	jen.ID("expected").Dot(
				"BelongsTo",
			)),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous response from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"BelongsTo",
	),
	jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.ID("expected"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetOAuth2Client",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
	),
	jen.ID("expected").Dot(
				"BelongsTo",
			)),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildGetOAuth2ClientCountQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedArgCount").Op(":=").Lit(1),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildGetOAuth2ClientCountQuery",
			).Call(jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_GetOAuth2ClientCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnRows",
			).Call(jen.ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(
	jen.Lit("count"))).Dot(
				"AddRow",
			).Call(jen.ID("expectedCount"))),
			jen.List(jen.ID("actualCount"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetOAuth2ClientCount",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildGetAllOAuth2ClientCountQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.ID("actualQuery").Op(":=").ID("m").Dot(
				"buildGetAllOAuth2ClientCountQuery",
			).Call(),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_GetAllOAuth2ClientCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(
	jen.Lit("count"))).Dot(
				"AddRow",
			).Call(jen.ID("expectedCount"))),
			jen.List(jen.ID("actualCount"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllOAuth2ClientCount",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildGetOAuth2ClientsQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedArgCount").Op(":=").Lit(1),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildGetOAuth2ClientsQuery",
			).Call(jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_GetOAuth2Clients").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2ClientList",
			).Valuesln(
	jen.ID("Pagination").Op(":").ID("models").Dot(
				"Pagination",
			).Valuesln(
	jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("TotalCount").Op(":").Lit(111)), jen.ID("Clients").Op(":").Index().ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())))),
			jen.ID("filter").Op(":=").ID("models").Dot(
				"DefaultQueryFilter",
			).Call(),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromOAuth2Client").Call(jen.Op("&").ID("expected").Dot(
				"Clients",
			).Index(jen.Lit(0))), jen.ID("buildMockRowFromOAuth2Client").Call(jen.Op("&").ID("expected").Dot(
				"Clients",
			).Index(jen.Lit(0))), jen.ID("buildMockRowFromOAuth2Client").Call(jen.Op("&").ID("expected").Dot(
				"Clients",
			).Index(jen.Lit(0)))),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnRows",
			).Call(jen.ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(
	jen.Lit("count"))).Dot(
				"AddRow",
			).Call(jen.ID("expected").Dot(
				"TotalCount",
			))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetOAuth2Clients",
			).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with no rows returned from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetOAuth2Clients",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error reading from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetOAuth2Clients",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2ClientList",
			).Valuesln(
	jen.ID("Pagination").Op(":").ID("models").Dot(
				"Pagination",
			).Valuesln(
	jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("TotalCount").Op(":").Lit(111)), jen.ID("Clients").Op(":").Index().ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())))),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.Op("&").ID("expected").Dot(
				"Clients",
			).Index(jen.Lit(0)))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetOAuth2Clients",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2ClientList",
			).Valuesln(
	jen.ID("Pagination").Op(":").ID("models").Dot(
				"Pagination",
			).Valuesln(
	jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("TotalCount").Op(":").Lit(0)), jen.ID("Clients").Op(":").Index().ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())))),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromOAuth2Client").Call(jen.Op("&").ID("expected").Dot(
				"Clients",
			).Index(jen.Lit(0))), jen.ID("buildMockRowFromOAuth2Client").Call(jen.Op("&").ID("expected").Dot(
				"Clients",
			).Index(jen.Lit(0))), jen.ID("buildMockRowFromOAuth2Client").Call(jen.Op("&").ID("expected").Dot(
				"Clients",
			).Index(jen.Lit(0)))),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetOAuth2Clients",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildCreateOAuth2ClientQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ClientID").Op(":").Lit("ClientID"), jen.ID("ClientSecret").Op(":").Lit("ClientSecret"), jen.ID("Scopes").Op(":").Index().ID("string").Valuesln(
	jen.Lit("blah")), jen.ID("RedirectURI").Op(":").Lit("RedirectURI"), jen.ID("BelongsTo").Op(":").Lit(123)),
			jen.ID("expectedArgCount").Op(":=").Lit(6),
			jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to,created_on) VALUES (?,?,?,?,?,?,UNIX_TIMESTAMP())"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildCreateOAuth2ClientQuery",
			).Call(jen.ID("exampleInput")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleInput").Dot(
				"Name",
	),
	jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleInput").Dot(
				"ClientID",
	),
	jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleInput").Dot(
				"ClientSecret",
	),
	jen.ID("args").Index(jen.Lit(2)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleInput").Dot(
				"Scopes",
			).Index(jen.Lit(0)), jen.ID("args").Index(jen.Lit(3)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleInput").Dot(
				"RedirectURI",
	),
	jen.ID("args").Index(jen.Lit(4)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleInput").Dot(
				"BelongsTo",
	),
	jen.ID("args").Index(jen.Lit(5)).Assert(jen.ID("uint64"))),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_CreateOAuth2Client").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedInput").Op(":=").Op("&").ID("models").Dot(
				"OAuth2ClientCreationInput",
			).Valuesln(
	jen.ID("Name").Op(":").ID("expected").Dot(
				"Name",
	),
	jen.ID("BelongsTo").Op(":").ID("expected").Dot(
				"BelongsTo",
			)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedCreationQuery").Op(":=").Lit("INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to,created_on) VALUES (?,?,?,?,?,?,UNIX_TIMESTAMP())"),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCreationQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"Name",
	),
	jen.ID("expected").Dot(
				"ClientID",
	),
	jen.ID("expected").Dot(
				"ClientSecret",
	),
	jen.Qual("strings", "Join").Call(jen.ID("expected").Dot(
				"Scopes",
	),
	jen.ID("scopesSeparator")), jen.ID("expected").Dot(
				"RedirectURI",
	),
	jen.ID("expected").Dot(
				"BelongsTo",
			)).Dot(
				"WillReturnResult",
			).Call(jen.ID("sqlmock").Dot(
				"NewResult",
			).Call(jen.ID("int64").Call(jen.ID("expected").Dot(
				"ID",
			)), jen.Lit(1))),
			jen.ID("expectedTimeQuery").Op(":=").Lit("SELECT created_on FROM oauth2_clients WHERE id = ?"),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnRows",
			).Call(jen.ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(
	jen.Lit("created_on"))).Dot(
				"AddRow",
			).Call(jen.ID("expected").Dot(
				"CreatedOn",
			))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"CreateOAuth2Client",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedInput").Op(":=").Op("&").ID("models").Dot(
				"OAuth2ClientCreationInput",
			).Valuesln(
	jen.ID("Name").Op(":").ID("expected").Dot(
				"Name",
	),
	jen.ID("BelongsTo").Op(":").ID("expected").Dot(
				"BelongsTo",
			)),
			jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to,created_on) VALUES (?,?,?,?,?,?,UNIX_TIMESTAMP())"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"Name",
	),
	jen.ID("expected").Dot(
				"ClientID",
	),
	jen.ID("expected").Dot(
				"ClientSecret",
	),
	jen.Qual("strings", "Join").Call(jen.ID("expected").Dot(
				"Scopes",
	),
	jen.ID("scopesSeparator")), jen.ID("expected").Dot(
				"RedirectURI",
	),
	jen.ID("expected").Dot(
				"BelongsTo",
			)).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"CreateOAuth2Client",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildUpdateOAuth2ClientQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ClientID").Op(":").Lit("ClientID"), jen.ID("ClientSecret").Op(":").Lit("ClientSecret"), jen.ID("Scopes").Op(":").Index().ID("string").Valuesln(
	jen.Lit("blah")), jen.ID("RedirectURI").Op(":").Lit("RedirectURI"), jen.ID("BelongsTo").Op(":").Lit(123)),
			jen.ID("expectedArgCount").Op(":=").Lit(6),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE oauth2_clients SET client_id = ?, client_secret = ?, scopes = ?, redirect_uri = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildUpdateOAuth2ClientQuery",
			).Call(jen.ID("expected")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"ClientID",
	),
	jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"ClientSecret",
	),
	jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"Scopes",
			).Index(jen.Lit(0)), jen.ID("args").Index(jen.Lit(2)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"RedirectURI",
	),
	jen.ID("args").Index(jen.Lit(3)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"BelongsTo",
	),
	jen.ID("args").Index(jen.Lit(4)).Assert(jen.ID("uint64"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"ID",
	),
	jen.ID("args").Index(jen.Lit(5)).Assert(jen.ID("uint64"))),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_UpdateOAuth2Client").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE oauth2_clients SET client_id = ?, client_secret = ?, scopes = ?, redirect_uri = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Values(),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnResult",
			).Call(jen.ID("sqlmock").Dot(
				"NewResult",
			).Call(jen.Lit(1), jen.Lit(1))),
			jen.ID("err").Op(":=").ID("m").Dot(
				"UpdateOAuth2Client",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleInput")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE oauth2_clients SET client_id = ?, client_secret = ?, scopes = ?, redirect_uri = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Values(),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("err").Op(":=").ID("m").Dot(
				"UpdateOAuth2Client",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleInput")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildArchiveOAuth2ClientQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedClientID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedArgCount").Op(":=").Lit(2),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE oauth2_clients SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildArchiveOAuth2ClientQuery",
			).Call(jen.ID("expectedClientID"), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedClientID"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_ArchiveOAuth2Client").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE oauth2_clients SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"),
			jen.ID("exampleClientID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("exampleUserID"), jen.ID("exampleClientID")).Dot(
				"WillReturnResult",
			).Call(jen.ID("sqlmock").Dot(
				"NewResult",
			).Call(jen.Lit(1), jen.Lit(1))),
			jen.ID("err").Op(":=").ID("m").Dot(
				"ArchiveOAuth2Client",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID"), jen.ID("exampleUserID")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE oauth2_clients SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"),
			jen.ID("exampleClientID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("exampleUserID"), jen.ID("exampleClientID")).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("err").Op(":=").ID("m").Dot(
				"ArchiveOAuth2Client",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID"), jen.ID("exampleUserID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)
	return ret
}
