package queriers

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/Masterminds/squirrel"
)

const (
	oauth2ClientsTableName            = "oauth2_clients"
	oauth2ClientsTableOwnershipColumn = "belongs_to_user"
)

var (
	oauth2ClientsTableColumns = []string{
		fmt.Sprintf("%s.id", oauth2ClientsTableName),
		fmt.Sprintf("%s.name", oauth2ClientsTableName),
		fmt.Sprintf("%s.client_id", oauth2ClientsTableName),
		fmt.Sprintf("%s.scopes", oauth2ClientsTableName),
		fmt.Sprintf("%s.redirect_uri", oauth2ClientsTableName),
		fmt.Sprintf("%s.client_secret", oauth2ClientsTableName),
		fmt.Sprintf("%s.created_on", oauth2ClientsTableName),
		fmt.Sprintf("%s.last_updated_on", oauth2ClientsTableName),
		fmt.Sprintf("%s.archived_on", oauth2ClientsTableName),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn),
	}
)

func oauth2ClientsTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	code := jen.NewFilePathName(proj.DatabaseV1Package("queriers", "v1", spn), spn)

	utils.AddImports(proj, code)

	code.Add(buildBuildMockRowsFromOAuth2Client(proj, dbvendor)...)
	code.Add(buildBuildErroneousMockRowFromOAuth2Client(proj, dbvendor)...)
	code.Add(buildTestScanOAuth2Clients(proj, dbvendor)...)
	code.Add(buildTestDB_buildGetOAuth2ClientByClientIDQuery(proj, dbvendor)...)
	code.Add(buildTestDB_GetOAuth2ClientByClientID(proj, dbvendor)...)
	code.Add(buildTestDB_buildGetAllOAuth2ClientsQuery(proj, dbvendor)...)
	code.Add(buildTestDB_GetAllOAuth2Clients(proj, dbvendor)...)
	code.Add(buildTestDB_GetAllOAuth2ClientsForUser(proj, dbvendor)...)
	code.Add(buildTestDB_buildGetOAuth2ClientQuery(proj, dbvendor)...)
	code.Add(buildTestDB_GetOAuth2Client(proj, dbvendor)...)
	code.Add(buildTestDB_buildGetAllOAuth2ClientsCountQuery(proj, dbvendor)...)
	code.Add(buildTestDB_GetAllOAuth2ClientCount(dbvendor)...)
	code.Add(buildTestDB_buildGetOAuth2ClientsForUserQuery(proj, dbvendor)...)
	code.Add(buildTestDB_GetOAuth2ClientsForUser(proj, dbvendor)...)
	code.Add(buildTestDB_buildCreateOAuth2ClientQuery(proj, dbvendor)...)
	code.Add(buildTestDB_CreateOAuth2Client(proj, dbvendor)...)
	code.Add(buildTestDB_buildUpdateOAuth2ClientQuery(proj, dbvendor)...)
	code.Add(buildTestDB_UpdateOAuth2Client(proj, dbvendor)...)
	code.Add(buildTestDB_buildArchiveOAuth2ClientQuery(proj, dbvendor)...)
	code.Add(buildTestDB_ArchiveOAuth2Client(proj, dbvendor)...)

	return code
}

func buildBuildMockRowsFromOAuth2Client(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildMockRowsFromOAuth2Client").Params(
			jen.ID("clients").Spread().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
		).Params(
			jen.PointerTo().Qual("github.com/DATA-DOG/go-sqlmock", "Rows"),
		).Body(
			jen.ID("columns").Assign().ID("oauth2ClientsTableColumns"),
			jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("columns")),
			jen.Line(),
			jen.For(jen.List(jen.Underscore(), jen.ID("c")).Assign().Range().ID("clients")).Body(
				jen.ID("rowValues").Assign().Index().Qual("database/sql/driver", "Value").Valuesln(
					jen.ID("c").Dot("ID"),
					jen.ID("c").Dot("Name"),
					jen.ID("c").Dot("ClientID"),
					jen.Qual("strings", "Join").Call(jen.ID("c").Dot("Scopes"), jen.ID("scopesSeparator")),
					jen.ID("c").Dot("RedirectURI"),
					jen.ID("c").Dot("ClientSecret"),
					jen.ID("c").Dot("CreatedOn"),
					jen.ID("c").Dot("LastUpdatedOn"),
					jen.ID("c").Dot("ArchivedOn"),
					jen.ID("c").Dot(constants.UserOwnershipFieldName),
				),
				jen.ID(utils.BuildFakeVarName("Rows")).Dot("AddRow").Call(jen.ID("rowValues").Spread()),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildErroneousMockRowFromOAuth2Client(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildErroneousMockRowFromOAuth2Client").Params(jen.ID("c").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Params(jen.PointerTo().Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Body(
			jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("oauth2ClientsTableColumns")).Dot("AddRow").Callln(
				jen.ID("c").Dot("ArchivedOn"),
				jen.ID("c").Dot("Name"),
				jen.ID("c").Dot("ClientID"),
				jen.Qual("strings", "Join").Call(jen.ID("c").Dot("Scopes"), jen.ID("scopesSeparator")),
				jen.ID("c").Dot("RedirectURI"),
				jen.ID("c").Dot("ClientSecret"),
				jen.ID("c").Dot("CreatedOn"),
				jen.ID("c").Dot("LastUpdatedOn"),
				jen.ID("c").Dot(constants.UserOwnershipFieldName),
				jen.ID("c").Dot("ID"),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	}

	return lines
}

func buildTestScanOAuth2Clients(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	lines := []jen.Code{
		jen.Func().IDf("Test%s_ScanOAuth2Clients", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
				).Assign().ID(dbfl).Dot("scanOAuth2Clients").Call(jen.ID("mockRows")),
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
				).Assign().ID(dbfl).Dot("scanOAuth2Clients").Call(jen.ID("mockRows")),
				utils.AssertNoError(jen.Err(), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDB_buildGetOAuth2ClientByClientIDQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.client_id", oauth2ClientsTableName):   whateverValue,
			fmt.Sprintf("%s.archived_on", oauth2ClientsTableName): nil,
		})

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "OAuth2Client"),
	}

	return buildQueryTest(dbvendor, "GetOAuth2ClientByClientID", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_GetOAuth2ClientByClientID(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	stn := "OAuth2Client"

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.client_id", oauth2ClientsTableName):   whateverValue,
			fmt.Sprintf("%s.archived_on", oauth2ClientsTableName): nil,
		}).ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetOAuth2ClientByClientID", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, stn),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(stn)).Dot("ClientID")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowsFromOAuth2Client").Call(jen.ID(utils.BuildFakeVarName(stn)))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2ClientByClientID").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName(stn)).Dot("ClientID"),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName(stn)), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"surfaces sql.ErrNoRows",
				utils.BuildFakeVar(proj, stn),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(stn)).Dot("ClientID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2ClientByClientID").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName(stn)).Dot("ClientID"),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with erroneous row",
				utils.BuildFakeVar(proj, stn),
				jen.Line(),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(stn)).Dot("ClientID")).Dotln("WillReturnRows").Call(
					jen.ID("buildErroneousMockRowFromOAuth2Client").Call(
						jen.ID(utils.BuildFakeVarName(stn)),
					),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2ClientByClientID").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName(stn)).Dot("ClientID"),
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

func buildTestDB_buildGetAllOAuth2ClientsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", oauth2ClientsTableName): nil,
		})

	expectedArgs := []jen.Code{}
	callArgs := []jen.Code{}

	return buildQueryTest(dbvendor, "GetAllOAuth2Clients", qb, expectedArgs, callArgs, nil)
}

func buildTestDB_GetAllOAuth2Clients(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", oauth2ClientsTableName): nil,
		}).ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetAllOAuth2Clients", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID("expected").Assign().Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID(utils.BuildFakeVarName("OAuth2Client")),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowsFromOAuth2Client").Callln(
						jen.ID(utils.BuildFakeVarName("OAuth2Client")),
						jen.ID(utils.BuildFakeVarName("OAuth2Client")),
						jen.ID(utils.BuildFakeVarName("OAuth2Client")),
					),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2Clients").Call(constants.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"surfaces sql.ErrNoRows",
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2Clients").Call(constants.CtxVar()),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error executing query",
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2Clients").Call(constants.CtxVar()),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with erroneous response from database",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2Clients").Call(constants.CtxVar()),
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

func buildTestDB_GetAllOAuth2ClientsForUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): whateverValue,
			fmt.Sprintf("%s.archived_on", oauth2ClientsTableName):                           nil,
		}).
		OrderBy(fmt.Sprintf("%s.id", oauth2ClientsTableName)).
		ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetAllOAuth2ClientsForUser", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				utils.BuildFakeVar(proj, "OAuth2ClientList"),
				jen.Line(),
				jen.ID("expected").Assign().Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.AddressOf().ID(utils.BuildFakeVarName("OAuth2ClientList")).Dot("Clients").Index(jen.Zero()),
					jen.AddressOf().ID(utils.BuildFakeVarName("OAuth2ClientList")).Dot("Clients").Index(jen.One()),
					jen.AddressOf().ID(utils.BuildFakeVarName("OAuth2ClientList")).Dot("Clients").Index(jen.Lit(2)),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowsFromOAuth2Client").Call(jen.ID("expected").Spread())),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2ClientsForUser").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User")).Dot("ID")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"surfaces sql.ErrNoRows",
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2ClientsForUser").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User")).Dot("ID")),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with erroneous response from database",
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2ClientsForUser").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User")).Dot("ID")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with unscannable response",
				utils.BuildFakeVar(proj, "User"),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2ClientsForUser").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User")).Dot("ID")),
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

func buildTestDB_buildGetOAuth2ClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", oauth2ClientsTableName):                                    whateverValue,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): whateverValue,
			fmt.Sprintf("%s.archived_on", oauth2ClientsTableName):                           nil,
		})

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "OAuth2Client"),
	}

	return buildQueryTest(dbvendor, "GetOAuth2Client", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_GetOAuth2Client(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	fvn := utils.BuildFakeVarName("OAuth2Client")
	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", oauth2ClientsTableName):                                    whateverValue,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): whateverValue,
			fmt.Sprintf("%s.archived_on", oauth2ClientsTableName):                           nil,
		}).
		ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetOAuth2Client", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(fvn).Dot(constants.UserOwnershipFieldName), jen.ID(fvn).Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowsFromOAuth2Client").Call(jen.ID(fvn))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2Client").Call(constants.CtxVar(), jen.ID(fvn).Dot("ID"), jen.ID(fvn).Dot(constants.UserOwnershipFieldName)),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(fvn), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"surfaces sql.ErrNoRows",
				utils.BuildFakeVar(proj, "OAuth2Client"), jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(fvn).Dot(constants.UserOwnershipFieldName), jen.ID(fvn).Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2Client").Call(constants.CtxVar(), jen.ID(fvn).Dot("ID"), jen.ID(fvn).Dot(constants.UserOwnershipFieldName)),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with erroneous response from database",
				utils.BuildFakeVar(proj, "OAuth2Client"), jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(fvn).Dot(constants.UserOwnershipFieldName), jen.ID(fvn).Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.ID(fvn))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2Client").Call(constants.CtxVar(), jen.ID(fvn).Dot("ID"), jen.ID(fvn).Dot(constants.UserOwnershipFieldName)),
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

func buildTestDB_buildGetAllOAuth2ClientsCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(fmt.Sprintf(countQuery, oauth2ClientsTableName)).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", oauth2ClientsTableName): nil,
		})

	var (
		expectedArgs []jen.Code
		callArgs     []jen.Code
	)

	return buildQueryTest(dbvendor, "GetAllOAuth2ClientsCount", qb, expectedArgs, callArgs, nil)
}

func buildTestDB_GetAllOAuth2ClientCount(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(fmt.Sprintf(countQuery, oauth2ClientsTableName)).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", oauth2ClientsTableName): nil,
		}).ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetAllOAuth2ClientCount", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expectedQuery").Assign().Lit(expectedQuery),
				jen.ID("expectedCount").Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actualCount"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2ClientCount").Call(constants.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expectedCount"), jen.ID("actualCount"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDB_buildGetOAuth2ClientsForUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): whateverValue,
			fmt.Sprintf("%s.archived_on", oauth2ClientsTableName):                           nil,
		})

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
		jen.ID(constants.FilterVarName),
	}

	qb = applyFleshedOutQueryFilter(qb, oauth2ClientsTableName)
	expectedArgs = appendFleshedOutQueryFilterArgs(expectedArgs)
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "User"),
		jen.ID(constants.FilterVarName).Assign().Qual(proj.FakeModelsPackage(), "BuildFleshedOutQueryFilter").Call(),
	}

	return buildQueryTest(dbvendor, "GetOAuth2ClientsForUser", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_GetOAuth2ClientsForUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))
	tn := "OAuth2ClientList"

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): whateverValue,
			fmt.Sprintf("%s.archived_on", oauth2ClientsTableName):                           nil,
		}).
		OrderBy(fmt.Sprintf("%s.id", oauth2ClientsTableName)).
		Limit(20).
		ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetOAuth2ClientsForUser", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildFakeVar(proj, "User"),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				utils.BuildFakeVar(proj, tn),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowsFromOAuth2Client").Callln(
						jen.AddressOf().ID(utils.BuildFakeVarName(tn)).Dot("Clients").Index(jen.Zero()),
						jen.AddressOf().ID(utils.BuildFakeVarName(tn)).Dot("Clients").Index(jen.One()),
						jen.AddressOf().ID(utils.BuildFakeVarName(tn)).Dot("Clients").Index(jen.Lit(2)),
					),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2ClientsForUser").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(constants.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName(tn)), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with no rows returned from database",
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2ClientsForUser").Call(
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
				"with error reading from database",
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2ClientsForUser").Call(
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
				"with erroneous response",
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(
					jen.Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2Client").Call(),
				)),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2ClientsForUser").Call(
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

func buildTestDB_buildCreateOAuth2ClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Insert(oauth2ClientsTableName).
		Columns(
			"name",
			"client_id",
			"client_secret",
			"scopes",
			"redirect_uri",
			oauth2ClientsTableOwnershipColumn,
		).
		Values(
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
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("Name"),
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientSecret"),
		jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("Scopes"), jen.ID("scopesSeparator")),
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("RedirectURI"),
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "OAuth2Client"),
	}

	return buildQueryTest(dbvendor, "CreateOAuth2Client", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_CreateOAuth2Client(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	qb := queryBuilderForDatabase(dbvendor).
		Insert(oauth2ClientsTableName).
		Columns(
			"name",
			"client_id",
			"client_secret",
			"scopes",
			"redirect_uri",
			oauth2ClientsTableOwnershipColumn,
		).
		Values(
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
		)

	var (
		happyPathExpectMethodName string
		happyPathReturnMethodName string
		exampleRowsDecl           jen.Code
	)
	if isPostgres(dbvendor) {
		happyPathExpectMethodName = "ExpectQuery"
		happyPathReturnMethodName = "WillReturnRows"
		qb = qb.Suffix("RETURNING id, created_on")
		exampleRowsDecl = jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").
			Call(jen.Index().String().Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").
			Call(
				jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
				jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("CreatedOn"),
			)
	} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		happyPathExpectMethodName = "ExpectExec"
		happyPathReturnMethodName = "WillReturnResult"
		exampleRowsDecl = jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").
			Call(jen.Int64().Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID")), jen.One())
	}

	happyPathLines := []jen.Code{
		utils.BuildFakeVar(proj, "OAuth2Client"),
		utils.BuildFakeVarWithCustomName(proj, "expectedInput", "BuildFakeOAuth2ClientCreationInputFromClient", jen.ID(utils.BuildFakeVarName("OAuth2Client"))),
		exampleRowsDecl,
		jen.Line(),
		jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
		jen.ID("mockDB").Dot(happyPathExpectMethodName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
			jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("Name"),
			jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
			jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientSecret"),
			jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("Scopes"), jen.ID("scopesSeparator")),
			jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("RedirectURI"),
			jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
		).Dot(happyPathReturnMethodName).Call(jen.ID(utils.BuildFakeVarName("Rows"))),
		jen.Line(),
	}

	if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		happyPathLines = append(happyPathLines,
			jen.ID("mtt").Assign().AddressOf().ID("mockTimeTeller").Values(),
			jen.ID("mtt").Dot("On").Call(jen.Lit("Now")).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("CreatedOn")),
			jen.ID(dbfl).Dot("timeTeller").Equals().ID("mtt"),
			jen.Line(),
		)
	}

	happyPathLines = append(happyPathLines,
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateOAuth2Client").Call(constants.CtxVar(), jen.ID("expectedInput")),
		utils.AssertNoError(jen.Err(), nil),
		utils.AssertEqual(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.ID("actual"), nil),
		jen.Line(),
		func() jen.Code {
			if isSqlite(dbvendor) || isMariaDB(dbvendor) {
				return utils.AssertExpectationsFor("mtt")
			}
			return jen.Null()
		}(),
		utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
	)

	expectedQuery, _, _ := qb.ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_CreateOAuth2Client", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest("happy path", happyPathLines...),
			jen.Line(),
			utils.BuildSubTest(
				"with error writing to database",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				utils.BuildFakeVarWithCustomName(proj, "expectedInput", "BuildFakeOAuth2ClientCreationInputFromClient", jen.ID(utils.BuildFakeVarName("OAuth2Client"))),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(happyPathExpectMethodName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("Name"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientSecret"),
					jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("Scopes"), jen.ID("scopesSeparator")),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("RedirectURI"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateOAuth2Client").Call(constants.CtxVar(), jen.ID("expectedInput")),
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

func buildTestDB_buildUpdateOAuth2ClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Update(oauth2ClientsTableName).
		Set("client_id", whateverValue).
		Set("client_secret", whateverValue).
		Set("scopes", whateverValue).
		Set("redirect_uri", whateverValue).
		Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id":                              whateverValue,
			oauth2ClientsTableOwnershipColumn: whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING last_updated_on")
	}

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientSecret"),
		jen.Qual("strings", "Join").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("Scopes"), jen.ID("scopesSeparator")),
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("RedirectURI"),
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "OAuth2Client"),
	}

	return buildQueryTest(dbvendor, "UpdateOAuth2Client", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_UpdateOAuth2Client(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	qb := queryBuilderForDatabase(dbvendor).
		Update(oauth2ClientsTableName).
		Set("client_id", whateverValue).
		Set("client_secret", whateverValue).
		Set("scopes", whateverValue).
		Set("redirect_uri", whateverValue).
		Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id":                              whateverValue,
			oauth2ClientsTableOwnershipColumn: whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING last_updated_on")
	}

	expectedQuery, _, _ := qb.ToSql()

	var (
		mockDBExpect        jen.Code
		errFuncExpectMethod string
	)
	if isPostgres(dbvendor) {
		mockDBExpect = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
			Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(
			jen.Index().String().Values(jen.Lit("last_updated_on")),
		).Dot("AddRow").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
		errFuncExpectMethod = "ExpectQuery"
	} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		mockDBExpect = jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
			Dotln("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.One(), jen.One()))
		errFuncExpectMethod = "ExpectExec"
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_UpdateOAuth2Client", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				mockDBExpect,
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("UpdateOAuth2Client").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("OAuth2Client"))),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error writing to database",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(errFuncExpectMethod).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("UpdateOAuth2Client").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("OAuth2Client"))),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDB_buildArchiveOAuth2ClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Update(oauth2ClientsTableName).
		Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Set("archived_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id":                              whateverValue,
			oauth2ClientsTableOwnershipColumn: whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING archived_on")
	}

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "OAuth2Client"),
	}

	return buildQueryTest(dbvendor, "ArchiveOAuth2Client", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_ArchiveOAuth2Client(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	queryTail := ""
	if isPostgres(dbvendor) {
		queryTail = " RETURNING archived_on"
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_ArchiveOAuth2Client", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("UPDATE oauth2_clients SET last_updated_on = %s, archived_on = %s WHERE belongs_to_user = %s AND id = %s%s",
				getTimeQuery(dbvendor),
				getTimeQuery(dbvendor),
				getIncIndex(dbvendor, 0),
				getIncIndex(dbvendor, 1),
				queryTail,
			),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
				).
					Dotln("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.One(), jen.One())),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("ArchiveOAuth2Client").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error writing to database",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
				).
					Dotln("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("ArchiveOAuth2Client").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}
