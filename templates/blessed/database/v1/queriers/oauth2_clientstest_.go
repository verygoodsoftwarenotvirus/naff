package queriers

import (
	"path/filepath"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsTestDotGo(pkg *models.Project, vendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(vendor.SingularPackageName())

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)
	sn := vendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))
	dbrn := vendor.RouteName()

	isPostgres := dbrn == "postgres"
	isSqlite := dbrn == "sqlite"
	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	ret.Add(
		jen.Func().ID("buildMockRowFromOAuth2Client").Params(jen.ID("c").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")).Params(jen.Op("*").Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("oauth2ClientsTableColumns")).Dot("AddRow").Callln(
				jen.ID("c").Dot("ID"),
				jen.ID("c").Dot("Name"),
				jen.ID("c").Dot("ClientID"),
				jen.Qual("strings", "Join").Call(jen.ID("c").Dot("Scopes"), jen.ID("scopesSeparator")),
				jen.ID("c").Dot("RedirectURI"),
				jen.ID("c").Dot("ClientSecret"),
				jen.ID("c").Dot("CreatedOn"),
				jen.ID("c").Dot("UpdatedOn"),
				jen.ID("c").Dot("ArchivedOn"),
				jen.ID("c").Dot("BelongsTo"),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildErroneousMockRowFromOAuth2Client").Params(jen.ID("c").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")).Params(jen.Op("*").Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("oauth2ClientsTableColumns")).Dot("AddRow").Callln(
				jen.ID("c").Dot("ArchivedOn"),
				jen.ID("c").Dot("Name"),
				jen.ID("c").Dot("ClientID"),
				jen.Qual("strings", "Join").Call(jen.ID("c").Dot("Scopes"), jen.ID("scopesSeparator")),
				jen.ID("c").Dot("RedirectURI"),
				jen.ID("c").Dot("ClientSecret"),
				jen.ID("c").Dot("CreatedOn"),
				jen.ID("c").Dot("UpdatedOn"),
				jen.ID("c").Dot("BelongsTo"),
				jen.ID("c").Dot("ID"),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetOAuth2ClientByClientIDQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedClientID").Op(":=").Lit("ClientID"),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND client_id = %s", getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetOAuth2ClientByClientIDQuery").Call(jen.ID("expectedClientID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedClientID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("string"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetOAuth2ClientByClientID", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleClientID").Op(":=").Lit("EXAMPLE"),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.Line(),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND client_id = %s", getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("exampleClientID")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetOAuth2ClientByClientID").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleClientID").Op(":=").Lit("EXAMPLE"),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND client_id = %s", getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("exampleClientID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetOAuth2ClientByClientID").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous row"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleClientID").Op(":=").Lit("EXAMPLE"),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.Line(),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND client_id = %s", getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("exampleClientID")).Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetOAuth2ClientByClientID").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetAllOAuth2ClientsQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.Line(),
				jen.ID("actualQuery").Op(":=").ID(dbfl).Dot("buildGetAllOAuth2ClientsQuery").Call(),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetAllOAuth2Clients", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Index().Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.Valuesln(
						jen.ID("ID").Op(":").Lit(123),
						jen.ID("Name").Op(":").Lit("name"),
						jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
						jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					),
				),
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Lit(0))),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllOAuth2Clients").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllOAuth2Clients").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error executing query"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllOAuth2Clients").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous response from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Index().Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.Valuesln(
						jen.ID("ID").Op(":").Lit(123),
						jen.ID("Name").Op(":").Lit("name"),
						jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
						jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					),
				),
				jen.Line(),
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Lit(0)))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllOAuth2Clients").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetAllOAuth2ClientsForUser", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Lit(123)),
				jen.ID("expected").Op(":=").Index().Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.Valuesln(
						jen.ID("ID").Op(":").Lit(123),
						jen.ID("Name").Op(":").Lit("name"),
						jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
						jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					),
				),
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Lit(0))),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllOAuth2ClientsForUser").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot("ID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Lit(123)),
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllOAuth2ClientsForUser").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot("ID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous response from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Lit(123)),
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllOAuth2ClientsForUser").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot("ID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with unscannable response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(jen.ID("ID").Op(":").Lit(123)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllOAuth2ClientsForUser").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUser").Dot("ID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetOAuth2ClientQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedClientID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedArgCount").Op(":=").Lit(2),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = %s AND id = %s", getIncIndex(dbrn, 0), getIncIndex(dbrn, 1)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetOAuth2ClientQuery").Call(jen.ID("expectedClientID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedClientID"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetOAuth2Client", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					jen.ID("Scopes").Op(":").Index().ID("string").Values(jen.Lit("things")),
				),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = %s AND id = %s", getIncIndex(dbrn, 0), getIncIndex(dbrn, 1)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("BelongsTo"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expected").Dot("BelongsTo")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					jen.ID("Scopes").Op(":").Index().ID("string").Values(jen.Lit("things")),
				),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = %s AND id = %s", getIncIndex(dbrn, 0), getIncIndex(dbrn, 1)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("BelongsTo"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expected").Dot("BelongsTo")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous response from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = %s AND id = %s", getIncIndex(dbrn, 0), getIncIndex(dbrn, 1)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("BelongsTo"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expected").Dot("BelongsTo")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetOAuth2ClientCountQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetOAuth2ClientCountQuery").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetOAuth2ClientCount", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", getIncIndex(dbrn, 0)),
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actualCount"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetOAuth2ClientCount").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetAllOAuth2ClientCountQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Lit("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.Line(),
				jen.ID("actual").Op(":=").ID(dbfl).Dot("buildGetAllOAuth2ClientCountQuery").Call(),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetAllOAuth2ClientCount", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actualCount"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllOAuth2ClientCount").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetOAuth2ClientsQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetOAuth2ClientsQuery").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetOAuth2Clients", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientList").Valuesln(
					jen.ID("Pagination").Op(":").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Pagination").Valuesln(
						jen.ID("Page").Op(":").Lit(1),
						jen.ID("Limit").Op(":").Lit(20),
						jen.ID("TotalCount").Op(":").Lit(111),
					),
					jen.ID("Clients").Op(":").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(123),
							jen.ID("Name").Op(":").Lit("name"),
							jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
							jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
						),
					),
				),
				jen.Line(),
				jen.ID("filter").Op(":=").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.ID("expectedCountQuery").Op(":=").Litf("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.Op("&").ID("expected").Dot("Clients").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.Op("&").ID("expected").Dot("Clients").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.Op("&").ID("expected").Dot("Clients").Index(jen.Lit(0))),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expected").Dot("TotalCount"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetOAuth2Clients").Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with no rows returned from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetOAuth2Clients").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error reading from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetOAuth2Clients").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientList").Valuesln(
					jen.ID("Pagination").Op(":").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Pagination").Valuesln(
						jen.ID("Page").Op(":").Lit(1),
						jen.ID("Limit").Op(":").Lit(20),
						jen.ID("TotalCount").Op(":").Lit(111),
					),
					jen.ID("Clients").Op(":").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(123),
							jen.ID("Name").Op(":").Lit("name"),
							jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
							jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
						),
					),
				),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.Op("&").ID("expected").Dot("Clients").Index(jen.Lit(0)))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetOAuth2Clients").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientList").Valuesln(
					jen.ID("Pagination").Op(":").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Pagination").Valuesln(
						jen.ID("Page").Op(":").Lit(1),
						jen.ID("Limit").Op(":").Lit(20),
						jen.ID("TotalCount").Op(":").Lit(0)),
					jen.ID("Clients").Op(":").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(123),
							jen.ID("Name").Op(":").Lit("name"),
							jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
							jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
						),
					),
				),
				jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.ID("expectedCountQuery").Op(":=").Litf("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.Op("&").ID("expected").Dot("Clients").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.Op("&").ID("expected").Dot("Clients").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.Op("&").ID("expected").Dot("Clients").Index(jen.Lit(0))),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetOAuth2Clients").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	var (
		queryTail string

		createdOnCol, createdOnVal string
	)
	if isPostgres {
		queryTail = " RETURNING id, created_on"
	}

	if isMariaDB {
		createdOnCol = ",created_on"
		createdOnVal = ",UNIX_TIMESTAMP()"
	}

	ret.Add(
		jen.Func().IDf("Test%s_buildCreateOAuth2ClientQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ClientID").Op(":").Lit("ClientID"),
					jen.ID("ClientSecret").Op(":").Lit("ClientSecret"),
					jen.ID("Scopes").Op(":").Index().ID("string").Values(jen.Lit("blah")),
					jen.ID("RedirectURI").Op(":").Lit("RedirectURI"),
					jen.ID("BelongsTo").Op(":").Lit(123),
				),
				jen.ID("expectedArgCount").Op(":=").Lit(6),
				jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to%s) VALUES (%s,%s,%s,%s,%s,%s%s)%s",
					createdOnCol,
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					getIncIndex(dbrn, 2),
					getIncIndex(dbrn, 3),
					getIncIndex(dbrn, 4),
					getIncIndex(dbrn, 5),
					createdOnVal,
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildCreateOAuth2ClientQuery").Call(jen.ID("exampleInput")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleInput").Dot("Name"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("string"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleInput").Dot("ClientID"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("string"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleInput").Dot("ClientSecret"), jen.ID("args").Index(jen.Lit(2)).Assert(jen.ID("string"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleInput").Dot("Scopes").Index(jen.Lit(0)), jen.ID("args").Index(jen.Lit(3)).Assert(jen.ID("string"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleInput").Dot("RedirectURI"), jen.ID("args").Index(jen.Lit(4)).Assert(jen.ID("string"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleInput").Dot("BelongsTo"), jen.ID("args").Index(jen.Lit(5)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	queryTail = ""
	var (
		happyPathExpectMethodName     string
		happyPathReturnMethodName     string
		createOAuth2ClientExampleRows jen.Code
		sqliteTimeCreationAddendum    jen.Code
	)
	if isPostgres {
		queryTail = " RETURNING id, created_on"
		happyPathExpectMethodName = "ExpectQuery"
		happyPathReturnMethodName = "WillReturnRows"
		createOAuth2ClientExampleRows = jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("ID"), jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
	} else if isSqlite || isMariaDB {
		happyPathExpectMethodName = "ExpectExec"
		happyPathReturnMethodName = "WillReturnResult"
		createOAuth2ClientExampleRows = jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.Lit(1))
		g := &jen.Group{}
		g.Add(
			jen.ID("expectedTimeQuery").Op(":=").Litf("SELECT created_on FROM oauth2_clients WHERE id = %s", getIncIndex(dbrn, 0)), jen.Line(),
			jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).
				Dotln("WithArgs").Call(jen.ID("expected").Dot("ID")).
				Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("CreatedOn"))), jen.Line(),
		)
		sqliteTimeCreationAddendum = g
	}

	ret.Add(
		jen.Func().IDf("Test%s_CreateOAuth2Client", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientCreationInput").Valuesln(
					jen.ID("Name").Op(":").ID("expected").Dot("Name"),
					jen.ID("BelongsTo").Op(":").ID("expected").Dot("BelongsTo"),
				),
				createOAuth2ClientExampleRows,
				jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to%s) VALUES (%s,%s,%s,%s,%s,%s%s)%s",
					createdOnCol,
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					getIncIndex(dbrn, 2),
					getIncIndex(dbrn, 3),
					getIncIndex(dbrn, 4),
					getIncIndex(dbrn, 5),
					createdOnVal,
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(happyPathExpectMethodName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Name"),
					jen.ID("expected").Dot("ClientID"),
					jen.ID("expected").Dot("ClientSecret"),
					jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Scopes"), jen.ID("scopesSeparator")),
					jen.ID("expected").Dot("RedirectURI"),
					jen.ID("expected").Dot("BelongsTo"),
				).Dot(happyPathReturnMethodName).Call(jen.ID("exampleRows")),
				jen.Line(),
				sqliteTimeCreationAddendum,
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("CreateOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientCreationInput").Valuesln(
					jen.ID("Name").Op(":").ID("expected").Dot("Name"),
					jen.ID("BelongsTo").Op(":").ID("expected").Dot("BelongsTo")),
				jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to%s) VALUES (%s,%s,%s,%s,%s,%s%s)%s",
					createdOnCol,
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					getIncIndex(dbrn, 2),
					getIncIndex(dbrn, 3),
					getIncIndex(dbrn, 4),
					getIncIndex(dbrn, 5),
					createdOnVal,
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(happyPathExpectMethodName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Name"),
					jen.ID("expected").Dot("ClientID"),
					jen.ID("expected").Dot("ClientSecret"),
					jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Scopes"), jen.ID("scopesSeparator")),
					jen.ID("expected").Dot("RedirectURI"),
					jen.ID("expected").Dot("BelongsTo"),
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("CreateOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	queryTail = ""
	if isPostgres {
		queryTail = " RETURNING updated_on"
	}

	ret.Add(
		jen.Func().IDf("Test%s_buildUpdateOAuth2ClientQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ClientID").Op(":").Lit("ClientID"),
					jen.ID("ClientSecret").Op(":").Lit("ClientSecret"),
					jen.ID("Scopes").Op(":").Index().ID("string").Values(jen.Lit("blah")),
					jen.ID("RedirectURI").Op(":").Lit("RedirectURI"),
					jen.ID("BelongsTo").Op(":").Lit(123),
				),
				jen.ID("expectedArgCount").Op(":=").Lit(6),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE oauth2_clients SET client_id = %s, client_secret = %s, scopes = %s, redirect_uri = %s, updated_on = %s WHERE belongs_to = %s AND id = %s%s",
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					getIncIndex(dbrn, 2),
					getIncIndex(dbrn, 3),
					getTimeQuery(dbrn),
					getIncIndex(dbrn, 4),
					getIncIndex(dbrn, 5),
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildUpdateOAuth2ClientQuery").Call(jen.ID("expected")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("ClientID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("string"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("ClientSecret"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("string"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("Scopes").Index(jen.Lit(0)), jen.ID("args").Index(jen.Lit(2)).Assert(jen.ID("string"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("RedirectURI"), jen.ID("args").Index(jen.Lit(3)).Assert(jen.ID("string"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("BelongsTo"), jen.ID("args").Index(jen.Lit(4)).Assert(jen.ID("uint64"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("ID"), jen.ID("args").Index(jen.Lit(5)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	queryTail = ""
	var (
		mockDBExpect        jen.Code
		errFuncExpectMethod string
	)
	if isPostgres {
		queryTail = " RETURNING updated_on"
		mockDBExpect = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WillReturnRows").Callln(
			jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("updated_on"))).Dot("AddRow").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
		)
		errFuncExpectMethod = "ExpectQuery"
	} else if isSqlite || isMariaDB {
		mockDBExpect = jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
			Dotln("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.Lit(1), jen.Lit(1)))
		errFuncExpectMethod = "ExpectExec"
	}

	ret.Add(
		jen.Func().IDf("Test%s_UpdateOAuth2Client", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE oauth2_clients SET client_id = %s, client_secret = %s, scopes = %s, redirect_uri = %s, updated_on = %s WHERE belongs_to = %s AND id = %s%s",
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					getIncIndex(dbrn, 2),
					getIncIndex(dbrn, 3),
					getTimeQuery(dbrn),
					getIncIndex(dbrn, 4),
					getIncIndex(dbrn, 5),
					queryTail,
				),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Values(),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				mockDBExpect,
				jen.Line(),
				jen.ID("err").Op(":=").ID(dbfl).Dot("UpdateOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleInput")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE oauth2_clients SET client_id = %s, client_secret = %s, scopes = %s, redirect_uri = %s, updated_on = %s WHERE belongs_to = %s AND id = %s%s",
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					getIncIndex(dbrn, 2),
					getIncIndex(dbrn, 3),
					getTimeQuery(dbrn),
					getIncIndex(dbrn, 4),
					getIncIndex(dbrn, 5),
					queryTail,
				),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client").Values(),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(errFuncExpectMethod).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.ID("err").Op(":=").ID(dbfl).Dot("UpdateOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleInput")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	queryTail = ""
	if isPostgres {
		queryTail = " RETURNING archived_on"
	}

	ret.Add(
		jen.Func().IDf("Test%s_buildArchiveOAuth2ClientQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedClientID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedArgCount").Op(":=").Lit(2),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE oauth2_clients SET updated_on = %s, archived_on = %s WHERE belongs_to = %s AND id = %s%s",
					getTimeQuery(dbrn),
					getTimeQuery(dbrn),
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildArchiveOAuth2ClientQuery").Call(jen.ID("expectedClientID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedClientID"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_ArchiveOAuth2Client", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE oauth2_clients SET updated_on = %s, archived_on = %s WHERE belongs_to = %s AND id = %s%s",
					getTimeQuery(dbrn),
					getTimeQuery(dbrn),
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					queryTail,
				),
				jen.ID("exampleClientID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("exampleUserID"), jen.ID("exampleClientID")).
					Dotln("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.Lit(1), jen.Lit(1))),
				jen.Line(),
				jen.ID("err").Op(":=").ID(dbfl).Dot("ArchiveOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE oauth2_clients SET updated_on = %s, archived_on = %s WHERE belongs_to = %s AND id = %s%s",
					getTimeQuery(dbrn),
					getTimeQuery(dbrn),
					getIncIndex(dbrn, 0),
					getIncIndex(dbrn, 1),
					queryTail,
				),
				jen.ID("exampleClientID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("exampleUserID"), jen.ID("exampleClientID")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.ID("err").Op(":=").ID(dbfl).Dot("ArchiveOAuth2Client").Call(jen.Qual("context", "Background").Call(), jen.ID("exampleClientID"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)
	return ret
}
