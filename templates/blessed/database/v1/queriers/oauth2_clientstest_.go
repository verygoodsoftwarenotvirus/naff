package queriers

import (
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, ret)
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))
	dbrn := dbvendor.RouteName()

	isPostgres := dbrn == "postgres"
	isSqlite := dbrn == "sqlite"
	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	ret.Add(
		jen.Func().ID("buildMockRowFromOAuth2Client").Params(jen.ID("c").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Params(jen.ParamPointer().Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("oauth2ClientsTableColumns")).Dot("AddRow").Callln(
				jen.ID("c").Dot("ID"),
				jen.ID("c").Dot("Name"),
				jen.ID("c").Dot("ClientID"),
				jen.Qual("strings", "Join").Call(jen.ID("c").Dot("Scopes"), jen.ID("scopesSeparator")),
				jen.ID("c").Dot("RedirectURI"),
				jen.ID("c").Dot("ClientSecret"),
				jen.ID("c").Dot("CreatedOn"),
				jen.ID("c").Dot("UpdatedOn"),
				jen.ID("c").Dot("ArchivedOn"),
				jen.ID("c").Dot("BelongsToUser"),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildErroneousMockRowFromOAuth2Client").Params(jen.ID("c").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Params(jen.ParamPointer().Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("oauth2ClientsTableColumns")).Dot("AddRow").Callln(
				jen.ID("c").Dot("ArchivedOn"),
				jen.ID("c").Dot("Name"),
				jen.ID("c").Dot("ClientID"),
				jen.Qual("strings", "Join").Call(jen.ID("c").Dot("Scopes"), jen.ID("scopesSeparator")),
				jen.ID("c").Dot("RedirectURI"),
				jen.ID("c").Dot("ClientSecret"),
				jen.ID("c").Dot("CreatedOn"),
				jen.ID("c").Dot("UpdatedOn"),
				jen.ID("c").Dot("BelongsToUser"),
				jen.ID("c").Dot("ID"),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetOAuth2ClientByClientIDQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedClientID").Assign().Lit("ClientID"),
				jen.ID("expectedArgCount").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedQuery").Assign().Litf("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to_user FROM oauth2_clients WHERE archived_on IS NULL AND client_id = %s", getIncIndex(dbvendor, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetOAuth2ClientByClientIDQuery").Call(jen.ID("expectedClientID")),
				utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
				utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
				utils.AssertEqual(jen.ID("expectedClientID"), jen.ID("args").Index(jen.Zero()).Assert(jen.String()), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetOAuth2ClientByClientID", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to_user FROM oauth2_clients WHERE archived_on IS NULL AND client_id = %s", getIncIndex(dbvendor, 0)),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("exampleClientID").Assign().Lit("EXAMPLE"),
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Lit("name"),
					jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
					jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.Line(),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("exampleClientID")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2ClientByClientID").Call(utils.CtxVar(), jen.ID("exampleClientID")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"surfaces sql.ErrNoRows",
				jen.ID("exampleClientID").Assign().Lit("EXAMPLE"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("exampleClientID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2ClientByClientID").Call(utils.CtxVar(), jen.ID("exampleClientID")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with erroneous row",
				jen.ID("exampleClientID").Assign().Lit("EXAMPLE"),
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Lit("name"),
					jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
					jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.Line(),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("exampleClientID")).Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2ClientByClientID").Call(utils.CtxVar(), jen.ID("exampleClientID")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetAllOAuth2ClientsQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedQuery").Assign().Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to_user FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.Line(),
				jen.ID("actualQuery").Assign().ID(dbfl).Dot("buildGetAllOAuth2ClientsQuery").Call(),
				utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetAllOAuth2Clients", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to_user FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.Valuesln(
						jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
						jen.ID("Name").MapAssign().Lit("name"),
						jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
						jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Zero())),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Zero())),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Zero())),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2Clients").Call(utils.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"surfaces sql.ErrNoRows",
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2Clients").Call(utils.CtxVar()),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error executing query",
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2Clients").Call(utils.CtxVar()),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with erroneous response from database",
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.Valuesln(
						jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
						jen.ID("Name").MapAssign().Lit("name"),
						jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
						jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Zero()))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2Clients").Call(utils.CtxVar()),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetAllOAuth2ClientsForUser", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to_user FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("exampleUser").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Values(jen.ID("ID").MapAssign().Add(utils.FakeUint64Func())),
				jen.ID("expected").Assign().Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.Valuesln(
						jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
						jen.ID("Name").MapAssign().Lit("name"),
						jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
						jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Zero())),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Zero())),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected").Index(jen.Zero())),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2ClientsForUser").Call(utils.CtxVar(), jen.ID("exampleUser").Dot("ID")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"surfaces sql.ErrNoRows",
				jen.ID("exampleUser").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Values(jen.ID("ID").MapAssign().Add(utils.FakeUint64Func())),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2ClientsForUser").Call(utils.CtxVar(), jen.ID("exampleUser").Dot("ID")),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with erroneous response from database",
				jen.ID("exampleUser").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Values(jen.ID("ID").MapAssign().Add(utils.FakeUint64Func())),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2ClientsForUser").Call(utils.CtxVar(), jen.ID("exampleUser").Dot("ID")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with unscannable response",
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("exampleUser").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Values(jen.ID("ID").MapAssign().Add(utils.FakeUint64Func())),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Lit("name"),
					jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
					jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2ClientsForUser").Call(utils.CtxVar(), jen.ID("exampleUser").Dot("ID")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetOAuth2ClientQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedClientID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedArgCount").Assign().Lit(2),
				jen.ID("expectedQuery").Assign().Litf("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to_user FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to_user = %s AND id = %s", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetOAuth2ClientQuery").Call(jen.ID("expectedClientID"), jen.ID("expectedUserID")),
				utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
				utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
				utils.AssertEqual(jen.ID("expectedUserID"), jen.ID("args").Index(jen.Zero()).Assert(jen.Uint64()), nil),
				utils.AssertEqual(jen.ID("expectedClientID"), jen.ID("args").Index(jen.Add(utils.FakeUint64Func())).Assert(jen.Uint64()), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetOAuth2Client", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to_user FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to_user = %s AND id = %s", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1)),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Lit("name"),
					jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
					jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					jen.ID("Scopes").MapAssign().Index().String().Values(jen.Lit("things")),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("BelongsToUser"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowFromOAuth2Client").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2Client").Call(utils.CtxVar(), jen.ID("expected").Dot("ID"), jen.ID("expected").Dot("BelongsToUser")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"surfaces sql.ErrNoRows",
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Lit("name"),
					jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
					jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					jen.ID("Scopes").MapAssign().Index().String().Values(jen.Lit("things")),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("BelongsToUser"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2Client").Call(utils.CtxVar(), jen.ID("expected").Dot("ID"), jen.ID("expected").Dot("BelongsToUser")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with erroneous response from database",
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Lit("name"),
					jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
					jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("BelongsToUser"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2Client").Call(utils.CtxVar(), jen.ID("expected").Dot("ID"), jen.ID("expected").Dot("BelongsToUser")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetOAuth2ClientCountQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedArgCount").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedQuery").Assign().Litf("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to_user = %s LIMIT 20", getIncIndex(dbvendor, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetOAuth2ClientCountQuery").Call(jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
				utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
				utils.AssertEqual(jen.ID("expectedUserID"), jen.ID("args").Index(jen.Zero()).Assert(jen.Uint64()), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetOAuth2ClientCount", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedQuery").Assign().Litf("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to_user = %s LIMIT 20", getIncIndex(dbvendor, 0)),
				jen.ID("expectedCount").Assign().Uint64().Call(jen.Lit(666)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actualCount"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2ClientCount").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expectedCount"), jen.ID("actualCount"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetAllOAuth2ClientCountQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().Lit("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.Line(),
				jen.ID("actual").Assign().ID(dbfl).Dot("buildGetAllOAuth2ClientCountQuery").Call(),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetAllOAuth2ClientCount", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expectedQuery").Assign().Lit("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL"),
				jen.ID("expectedCount").Assign().Uint64().Call(jen.Lit(666)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actualCount"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2ClientCount").Call(utils.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expectedCount"), jen.ID("actualCount"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetOAuth2ClientsQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedArgCount").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedQuery").Assign().Litf("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to_user FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to_user = %s LIMIT 20", getIncIndex(dbvendor, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetOAuth2ClientsQuery").Call(jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
				utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
				utils.AssertEqual(jen.ID("expectedUserID"), jen.ID("args").Index(jen.Zero()).Assert(jen.Uint64()), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetOAuth2Clients", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedListQuery").Assign().Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to_user FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",

				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2ClientList").Valuesln(
					jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().Add(utils.FakeUint64Func()),
						jen.ID("Limit").MapAssign().Lit(20),
						jen.ID("TotalCount").MapAssign().Lit(111),
					),
					jen.ID("Clients").MapAssign().Index().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
						jen.Valuesln(
							jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
							jen.ID("Name").MapAssign().Lit("name"),
							jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
							jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
						),
					),
				),
				jen.Line(),
				jen.ID(utils.FilterVarName).Assign().Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call(),
				jen.ID("expectedCountQuery").Assign().Litf("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to_user = %s LIMIT 20", getIncIndex(dbvendor, 0)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.VarPointer().ID("expected").Dot("Clients").Index(jen.Zero())),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.VarPointer().ID("expected").Dot("Clients").Index(jen.Zero())),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.VarPointer().ID("expected").Dot("Clients").Index(jen.Zero())),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expected").Dot("TotalCount"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2Clients").Call(utils.CtxVar(), jen.ID(utils.FilterVarName), jen.ID("expectedUserID")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with no rows returned from database",
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2Clients").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error reading from database",
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2Clients").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with erroneous response",
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2ClientList").Valuesln(
					jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().Add(utils.FakeUint64Func()),
						jen.ID("Limit").MapAssign().Lit(20),
						jen.ID("TotalCount").MapAssign().Lit(111),
					),
					jen.ID("Clients").MapAssign().Index().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
						jen.Valuesln(
							jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
							jen.ID("Name").MapAssign().Lit("name"),
							jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
							jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
						),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.VarPointer().ID("expected").Dot("Clients").Index(jen.Zero()))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2Clients").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error fetching count",
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2ClientList").Valuesln(
					jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().Add(utils.FakeUint64Func()),
						jen.ID("Limit").MapAssign().Lit(20),
						jen.ID("TotalCount").MapAssign().Zero()),
					jen.ID("Clients").MapAssign().Index().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
						jen.Valuesln(
							jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
							jen.ID("Name").MapAssign().Lit("name"),
							jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
							jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
						),
					),
				),
				jen.ID("expectedCountQuery").Assign().Litf("SELECT COUNT(id) FROM oauth2_clients WHERE archived_on IS NULL AND belongs_to_user = %s LIMIT 20", getIncIndex(dbvendor, 0)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.VarPointer().ID("expected").Dot("Clients").Index(jen.Zero())),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.VarPointer().ID("expected").Dot("Clients").Index(jen.Zero())),
					jen.ID("buildMockRowFromOAuth2Client").Call(jen.VarPointer().ID("expected").Dot("Clients").Index(jen.Zero())),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2Clients").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
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
		jen.Func().IDf("Test%s_buildCreateOAuth2ClientQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ClientID").MapAssign().Lit("ClientID"),
					jen.ID("ClientSecret").MapAssign().Lit("ClientSecret"),
					jen.ID("Scopes").MapAssign().Index().String().Values(jen.Lit("blah")),
					jen.ID("RedirectURI").MapAssign().Lit("RedirectURI"),
					jen.ID("BelongsToUser").MapAssign().Add(utils.FakeUint64Func()),
				),
				jen.ID("expectedArgCount").Assign().Lit(6),
				jen.ID("expectedQuery").Assign().Litf("INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to_user%s) VALUES (%s,%s,%s,%s,%s,%s%s)%s",
					createdOnCol,
					getIncIndex(dbvendor, 0),
					getIncIndex(dbvendor, 1),
					getIncIndex(dbvendor, 2),
					getIncIndex(dbvendor, 3),
					getIncIndex(dbvendor, 4),
					getIncIndex(dbvendor, 5),
					createdOnVal,
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildCreateOAuth2ClientQuery").Call(jen.ID("exampleInput")),
				utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
				utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
				utils.AssertEqual(jen.ID("exampleInput").Dot("Name"), jen.ID("args").Index(jen.Zero()).Assert(jen.String()), nil),
				utils.AssertEqual(jen.ID("exampleInput").Dot("ClientID"), jen.ID("args").Index(jen.Add(utils.FakeUint64Func())).Assert(jen.String()), nil),
				utils.AssertEqual(jen.ID("exampleInput").Dot("ClientSecret"), jen.ID("args").Index(jen.Lit(2)).Assert(jen.String()), nil),
				utils.AssertEqual(jen.ID("exampleInput").Dot("Scopes").Index(jen.Zero()), jen.ID("args").Index(jen.Lit(3)).Assert(jen.String()), nil),
				utils.AssertEqual(jen.ID("exampleInput").Dot("RedirectURI"), jen.ID("args").Index(jen.Lit(4)).Assert(jen.String()), nil),
				utils.AssertEqual(jen.ID("exampleInput").Dot("BelongsToUser"), jen.ID("args").Index(jen.Lit(5)).Assert(jen.Uint64()), nil),
			),
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
		createOAuth2ClientExampleRows = jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("ID"), jen.Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
	} else if isSqlite || isMariaDB {
		happyPathExpectMethodName = "ExpectExec"
		happyPathReturnMethodName = "WillReturnResult"
		createOAuth2ClientExampleRows = jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.Add(utils.FakeUint64Func()))
		g := &jen.Group{}
		g.Add(
			jen.ID("expectedTimeQuery").Assign().Litf("SELECT created_on FROM oauth2_clients WHERE id = %s", getIncIndex(dbvendor, 0)), jen.Line(),
			jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).
				Dotln("WithArgs").Call(jen.ID("expected").Dot("ID")).
				Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("CreatedOn"))), jen.Line(),
		)
		sqliteTimeCreationAddendum = g
	}

	ret.Add(
		jen.Func().IDf("Test%s_CreateOAuth2Client", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to_user%s) VALUES (%s,%s,%s,%s,%s,%s%s)%s",
				createdOnCol,
				getIncIndex(dbvendor, 0),
				getIncIndex(dbvendor, 1),
				getIncIndex(dbvendor, 2),
				getIncIndex(dbvendor, 3),
				getIncIndex(dbvendor, 4),
				getIncIndex(dbvendor, 5),
				createdOnVal,
				queryTail,
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Lit("name"),
					jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
					jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput").Valuesln(
					jen.ID("Name").MapAssign().ID("expected").Dot("Name"),
					jen.ID("BelongsToUser").MapAssign().ID("expected").Dot("BelongsToUser"),
				),
				createOAuth2ClientExampleRows,
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(happyPathExpectMethodName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Name"),
					jen.ID("expected").Dot("ClientID"),
					jen.ID("expected").Dot("ClientSecret"),
					jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Scopes"), jen.ID("scopesSeparator")),
					jen.ID("expected").Dot("RedirectURI"),
					jen.ID("expected").Dot("BelongsToUser"),
				).Dot(happyPathReturnMethodName).Call(jen.ID("exampleRows")),
				jen.Line(),
				sqliteTimeCreationAddendum,
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateOAuth2Client").Call(utils.CtxVar(), jen.ID("expectedInput")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error writing to database",
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Name").MapAssign().Lit("name"),
					jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
					jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput").Valuesln(
					jen.ID("Name").MapAssign().ID("expected").Dot("Name"),
					jen.ID("BelongsToUser").MapAssign().ID("expected").Dot("BelongsToUser")),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(happyPathExpectMethodName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Name"),
					jen.ID("expected").Dot("ClientID"),
					jen.ID("expected").Dot("ClientSecret"),
					jen.Qual("strings", "Join").Call(jen.ID("expected").Dot("Scopes"), jen.ID("scopesSeparator")),
					jen.ID("expected").Dot("RedirectURI"),
					jen.ID("expected").Dot("BelongsToUser"),
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateOAuth2Client").Call(utils.CtxVar(), jen.ID("expectedInput")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	)

	queryTail = ""
	if isPostgres {
		queryTail = " RETURNING updated_on"
	}

	ret.Add(
		jen.Func().IDf("Test%s_buildUpdateOAuth2ClientQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ClientID").MapAssign().Lit("ClientID"),
					jen.ID("ClientSecret").MapAssign().Lit("ClientSecret"),
					jen.ID("Scopes").MapAssign().Index().String().Values(jen.Lit("blah")),
					jen.ID("RedirectURI").MapAssign().Lit("RedirectURI"),
					jen.ID("BelongsToUser").MapAssign().Add(utils.FakeUint64Func()),
				),
				jen.ID("expectedArgCount").Assign().Lit(6),
				jen.ID("expectedQuery").Assign().Litf("UPDATE oauth2_clients SET client_id = %s, client_secret = %s, scopes = %s, redirect_uri = %s, updated_on = %s WHERE belongs_to_user = %s AND id = %s%s",
					getIncIndex(dbvendor, 0),
					getIncIndex(dbvendor, 1),
					getIncIndex(dbvendor, 2),
					getIncIndex(dbvendor, 3),
					getTimeQuery(dbvendor),
					getIncIndex(dbvendor, 4),
					getIncIndex(dbvendor, 5),
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildUpdateOAuth2ClientQuery").Call(jen.ID("expected")),
				utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
				utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
				utils.AssertEqual(jen.ID("expected").Dot("ClientID"), jen.ID("args").Index(jen.Zero()).Assert(jen.String()), nil),
				utils.AssertEqual(jen.ID("expected").Dot("ClientSecret"), jen.ID("args").Index(jen.Add(utils.FakeUint64Func())).Assert(jen.String()), nil),
				utils.AssertEqual(jen.ID("expected").Dot("Scopes").Index(jen.Zero()), jen.ID("args").Index(jen.Lit(2)).Assert(jen.String()), nil),
				utils.AssertEqual(jen.ID("expected").Dot("RedirectURI"), jen.ID("args").Index(jen.Lit(3)).Assert(jen.String()), nil),
				utils.AssertEqual(jen.ID("expected").Dot("BelongsToUser"), jen.ID("args").Index(jen.Lit(4)).Assert(jen.Uint64()), nil),
				utils.AssertEqual(jen.ID("expected").Dot("ID"), jen.ID("args").Index(jen.Lit(5)).Assert(jen.Uint64()), nil),
			),
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
			jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("updated_on"))).Dot("AddRow").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
		)
		errFuncExpectMethod = "ExpectQuery"
	} else if isSqlite || isMariaDB {
		mockDBExpect = jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
			Dotln("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.Add(utils.FakeUint64Func()), jen.Add(utils.FakeUint64Func())))
		errFuncExpectMethod = "ExpectExec"
	}

	ret.Add(
		jen.Func().IDf("Test%s_UpdateOAuth2Client", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("UPDATE oauth2_clients SET client_id = %s, client_secret = %s, scopes = %s, redirect_uri = %s, updated_on = %s WHERE belongs_to_user = %s AND id = %s%s",
				getIncIndex(dbvendor, 0),
				getIncIndex(dbvendor, 1),
				getIncIndex(dbvendor, 2),
				getIncIndex(dbvendor, 3),
				getTimeQuery(dbvendor),
				getIncIndex(dbvendor, 4),
				getIncIndex(dbvendor, 5),
				queryTail,
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("exampleInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				mockDBExpect,
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("UpdateOAuth2Client").Call(utils.CtxVar(), jen.ID("exampleInput")),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error writing to database",
				jen.ID("exampleInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(errFuncExpectMethod).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("UpdateOAuth2Client").Call(utils.CtxVar(), jen.ID("exampleInput")),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	)

	queryTail = ""
	if isPostgres {
		queryTail = " RETURNING archived_on"
	}

	ret.Add(
		jen.Func().IDf("Test%s_buildArchiveOAuth2ClientQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedClientID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedArgCount").Assign().Lit(2),
				jen.ID("expectedQuery").Assign().Litf("UPDATE oauth2_clients SET updated_on = %s, archived_on = %s WHERE belongs_to_user = %s AND id = %s%s",
					getTimeQuery(dbvendor),
					getTimeQuery(dbvendor),
					getIncIndex(dbvendor, 0),
					getIncIndex(dbvendor, 1),
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildArchiveOAuth2ClientQuery").Call(jen.ID("expectedClientID"), jen.ID("expectedUserID")),
				utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
				utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
				utils.AssertEqual(jen.ID("expectedUserID"), jen.ID("args").Index(jen.Zero()).Assert(jen.Uint64()), nil),
				utils.AssertEqual(jen.ID("expectedClientID"), jen.ID("args").Index(jen.Add(utils.FakeUint64Func())).Assert(jen.Uint64()), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_ArchiveOAuth2Client", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("UPDATE oauth2_clients SET updated_on = %s, archived_on = %s WHERE belongs_to_user = %s AND id = %s%s",
				getTimeQuery(dbvendor),
				getTimeQuery(dbvendor),
				getIncIndex(dbvendor, 0),
				getIncIndex(dbvendor, 1),
				queryTail,
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("exampleClientID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("exampleUserID"), jen.ID("exampleClientID")).
					Dotln("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.Add(utils.FakeUint64Func()), jen.Add(utils.FakeUint64Func()))),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("ArchiveOAuth2Client").Call(utils.CtxVar(), jen.ID("exampleClientID"), jen.ID("exampleUserID")),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error writing to database",
				jen.ID("exampleClientID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("exampleUserID"), jen.ID("exampleClientID")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("ArchiveOAuth2Client").Call(utils.CtxVar(), jen.ID("exampleClientID"), jen.ID("exampleUserID")),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	)
	return ret
}
