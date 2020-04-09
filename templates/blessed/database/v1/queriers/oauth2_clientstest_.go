package queriers

import (
	"fmt"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/Masterminds/squirrel"
)

const (
	scopesSeparator                   = ","
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
		fmt.Sprintf("%s.updated_on", oauth2ClientsTableName),
		fmt.Sprintf("%s.archived_on", oauth2ClientsTableName),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn),
	}
)

func oauth2ClientsTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, ret)

	ret.Add(buildBuildMockRowsFromOAuth2Client(proj, dbvendor)...)
	ret.Add(buildBuildErroneousMockRowFromOAuth2Client(proj, dbvendor)...)
	ret.Add(buildTestDB_buildGetOAuth2ClientByClientIDQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_GetOAuth2ClientByClientID(proj, dbvendor)...)
	ret.Add(buildTestDB_buildGetAllOAuth2ClientsQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_GetAllOAuth2Clients(proj, dbvendor)...)
	ret.Add(buildTestDB_GetAllOAuth2ClientsForUser(proj, dbvendor)...)
	ret.Add(buildTestDB_buildGetOAuth2ClientQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_GetOAuth2Client(proj, dbvendor)...)
	ret.Add(buildTestDB_buildGetAllOAuth2ClientCountQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_GetAllOAuth2ClientCount(proj, dbvendor)...)
	ret.Add(buildTestDB_buildGetOAuth2ClientsQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_GetOAuth2Clients(proj, dbvendor)...)
	ret.Add(buildTestDB_buildCreateOAuth2ClientQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_CreateOAuth2Client(proj, dbvendor)...)
	ret.Add(buildTestDB_buildUpdateOAuth2ClientQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_UpdateOAuth2Client(proj, dbvendor)...)
	ret.Add(buildTestDB_buildArchiveOAuth2ClientQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_ArchiveOAuth2Client(proj, dbvendor)...)

	return ret
}

func buildBuildMockRowsFromOAuth2Client(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildMockRowsFromOAuth2Client").Params(
			jen.ID("clients").Spread().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
		).Params(
			jen.ParamPointer().Qual("github.com/DATA-DOG/go-sqlmock", "Rows"),
		).Block(
			jen.ID("includeCount").Assign().Len(jen.ID("clients")).GreaterThan().One(),
			jen.ID("columns").Assign().ID("oauth2ClientsTableColumns"),
			jen.Line(),
			jen.If(jen.ID("includeCount")).Block(
				utils.AppendItemsToList(jen.ID("columns"), jen.Lit("count")),
			),
			jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("columns")),
			jen.Line(),
			jen.For(jen.List(jen.Underscore(), jen.ID("c")).Assign().Range().ID("clients")).Block(
				jen.ID("rowValues").Assign().Index().Qual("database/sql/driver", "Value").Valuesln(
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
				jen.If(jen.ID("includeCount")).Block(
					utils.AppendItemsToList(jen.ID("rowValues"), jen.Len(jen.ID("clients"))),
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

func buildBuildErroneousMockRowFromOAuth2Client(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
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

	return buildQueryTest(proj, dbvendor, models.DataType{Name: GetOAuth2ClientPalabra()}, "GetOAuth2ClientByClientID", qb, expectedArgs, callArgs, true, false, false, false, true)
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
		jen.Func().IDf("Test%s_GetOAuth2ClientByClientID", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
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
					utils.CtxVar(),
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
					utils.CtxVar(),
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
					utils.CtxVar(),
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

	return buildQueryTest(proj, dbvendor, models.DataType{}, "GetAllOAuth2Clients", qb, expectedArgs, callArgs, false, true, false, false, false)
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
		jen.Func().IDf("Test%s_GetAllOAuth2Clients", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
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
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2Clients").Call(utils.CtxVar()),
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
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2Clients").Call(utils.CtxVar()),
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
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2Clients").Call(utils.CtxVar()),
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
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2Clients").Call(utils.CtxVar()),
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
		Select(append(oauth2ClientsTableColumns, fmt.Sprintf(countQuery, oauth2ClientsTableName))...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): whateverValue,
			fmt.Sprintf("%s.archived_on", oauth2ClientsTableName):                           nil,
		}).
		GroupBy(fmt.Sprintf("%s.id", oauth2ClientsTableName)).
		ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetAllOAuth2ClientsForUser", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleUser").Assign().AddressOf().Qual(proj.ModelsV1Package(), "User").Values(jen.ID("ID").MapAssign().One()),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowsFromOAuth2Client").Call(jen.ID("expected").Index(jen.Zero())),
					jen.ID("buildMockRowsFromOAuth2Client").Call(jen.ID("expected").Index(jen.One())),
					jen.ID("buildMockRowsFromOAuth2Client").Call(jen.ID("expected").Index(jen.Lit(2))),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllOAuth2ClientsForUser").Call(utils.CtxVar(), jen.ID("exampleUser").Dot("ID")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"surfaces sql.ErrNoRows",
				jen.ID("exampleUser").Assign().AddressOf().Qual(proj.ModelsV1Package(), "User").Values(jen.ID("ID").MapAssign().One()),
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
			utils.BuildSubTest(
				"with erroneous response from database",
				jen.ID("exampleUser").Assign().AddressOf().Qual(proj.ModelsV1Package(), "User").Values(jen.ID("ID").MapAssign().One()),
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
			utils.BuildSubTest(
				"with unscannable response",
				jen.ID("exampleUser").Assign().AddressOf().Qual(proj.ModelsV1Package(), "User").Values(jen.ID("ID").MapAssign().One()),
				utils.BuildFakeVar(proj, "OAuth2Client"),
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
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
	}

	return buildQueryTest(proj, dbvendor, models.DataType{Name: GetOAuth2ClientPalabra()}, "GetOAuth2Client", qb, expectedArgs, callArgs, false, false, false, false, false)
}

func buildTestDB_GetOAuth2Client(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(append(oauth2ClientsTableColumns, fmt.Sprintf(countQuery, oauth2ClientsTableName))...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): whateverValue,
			fmt.Sprintf("%s.archived_on", oauth2ClientsTableName):                           nil,
		}).
		GroupBy(fmt.Sprintf("%s.id", oauth2ClientsTableName)).
		ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetOAuth2Client", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
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
					Dotln("WithArgs").Call(jen.ID("expected").Dot("BelongsToUser"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowsFromOAuth2Client").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetOAuth2Client").Call(utils.CtxVar(), jen.ID("expected").Dot("ID"), jen.ID("expected").Dot("BelongsToUser")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"surfaces sql.ErrNoRows",
				utils.BuildFakeVar(proj, "OAuth2Client"), jen.Line(),
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
			utils.BuildSubTest(
				"with erroneous response from database",
				utils.BuildFakeVar(proj, "OAuth2Client"), jen.Line(),
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
	}

	return lines
}

func buildTestDB_buildGetAllOAuth2ClientCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(fmt.Sprintf(countQuery, oauth2ClientsTableName)).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", oauth2ClientsTableName): nil,
		})

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
	}

	return buildQueryTest(proj, dbvendor, models.DataType{Name: GetOAuth2ClientPalabra()}, "GetAllOAuth2ClientCount", qb, expectedArgs, callArgs, false, false, false, false, false)
}

func buildTestDB_GetAllOAuth2ClientCount(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	lines := []jen.Code{
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
	}

	return lines
}

func buildTestDB_buildGetOAuth2ClientsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(append(oauth2ClientsTableColumns, fmt.Sprintf(countQuery, oauth2ClientsTableName))...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): whateverValue,
			fmt.Sprintf("%s.archived_on", oauth2ClientsTableName):                           nil,
		}).
		GroupBy(fmt.Sprintf("%s.id", oauth2ClientsTableName))

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
	}

	return buildQueryTest(proj, dbvendor, models.DataType{Name: GetOAuth2ClientPalabra()}, "GetOAuth2Clients", qb, expectedArgs, callArgs, false, false, false, false, false)
}

func buildTestDB_GetOAuth2Clients(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetOAuth2Clients", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedListQuery").Assign().Lit("SELECT id, name, client_id, scopes, redirect_uri, client_secret, created_on, updated_on, archived_on, belongs_to_user FROM oauth2_clients WHERE archived_on IS NULL"),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",

				jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "OAuth2ClientList").Valuesln(
					jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().One(),
						jen.ID("Limit").MapAssign().Lit(20),
						jen.ID("TotalCount").MapAssign().Lit(111),
					),
					jen.ID("Clients").MapAssign().Index().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
						jen.Valuesln(
							jen.ID("ID").MapAssign().One(),
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
					jen.ID("buildMockRowsFromOAuth2Client").Call(jen.AddressOf().ID("expected").Dot("Clients").Index(jen.Zero())),
					jen.ID("buildMockRowsFromOAuth2Client").Call(jen.AddressOf().ID("expected").Dot("Clients").Index(jen.Zero())),
					jen.ID("buildMockRowsFromOAuth2Client").Call(jen.AddressOf().ID("expected").Dot("Clients").Index(jen.Zero())),
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
				jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "OAuth2ClientList").Valuesln(
					jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().One(),
						jen.ID("Limit").MapAssign().Lit(20),
						jen.ID("TotalCount").MapAssign().Lit(111),
					),
					jen.ID("Clients").MapAssign().Index().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
						jen.Valuesln(
							jen.ID("ID").MapAssign().One(),
							jen.ID("Name").MapAssign().Lit("name"),
							jen.ID("BelongsToUser").MapAssign().ID("expectedUserID"),
							jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
						),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromOAuth2Client").Call(jen.AddressOf().ID("expected").Dot("Clients").Index(jen.Zero()))),
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
				jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "OAuth2ClientList").Valuesln(
					jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().One(),
						jen.ID("Limit").MapAssign().Lit(20),
						jen.ID("TotalCount").MapAssign().Zero()),
					jen.ID("Clients").MapAssign().Index().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
						jen.Valuesln(
							jen.ID("ID").MapAssign().One(),
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
					jen.ID("buildMockRowsFromOAuth2Client").Call(jen.AddressOf().ID("expected").Dot("Clients").Index(jen.Zero())),
					jen.ID("buildMockRowsFromOAuth2Client").Call(jen.AddressOf().ID("expected").Dot("Clients").Index(jen.Zero())),
					jen.ID("buildMockRowsFromOAuth2Client").Call(jen.AddressOf().ID("expected").Dot("Clients").Index(jen.Zero())),
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
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
	}

	return buildQueryTest(proj, dbvendor, models.DataType{Name: GetOAuth2ClientPalabra()}, "CreateOAuth2Client", qb, expectedArgs, callArgs, false, false, false, false, false)
}

func buildTestDB_CreateOAuth2Client(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	queryTail := ""
	var (
		happyPathExpectMethodName     string
		happyPathReturnMethodName     string
		createOAuth2ClientExampleRows jen.Code
		sqliteTimeCreationAddendum    jen.Code
	)
	if isPostgres(dbvendor) {
		queryTail = " RETURNING id, created_on"
		happyPathExpectMethodName = "ExpectQuery"
		happyPathReturnMethodName = "WillReturnRows"
		createOAuth2ClientExampleRows = jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("ID"), jen.Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
	} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		happyPathExpectMethodName = "ExpectExec"
		happyPathReturnMethodName = "WillReturnResult"
		createOAuth2ClientExampleRows = jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.One())
		g := &jen.Group{}
		g.Add(
			jen.ID("expectedTimeQuery").Assign().Litf("SELECT created_on FROM oauth2_clients WHERE id = %s", getIncIndex(dbvendor, 0)), jen.Line(),
			jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).
				Dotln("WithArgs").Call(jen.ID("expected").Dot("ID")).
				Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("CreatedOn"))), jen.Line(),
		)
		sqliteTimeCreationAddendum = g
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_CreateOAuth2Client", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to_user) VALUES (%s,%s,%s,%s,%s,%s%s)%s",
				getIncIndex(dbvendor, 0),
				getIncIndex(dbvendor, 1),
				getIncIndex(dbvendor, 2),
				getIncIndex(dbvendor, 3),
				getIncIndex(dbvendor, 4),
				getIncIndex(dbvendor, 5),
				queryTail,
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID("expectedInput").Assign().AddressOf().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput").Valuesln(
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
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID("expectedInput").Assign().AddressOf().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput").Valuesln(
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
		Set("updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id":                              whateverValue,
			oauth2ClientsTableOwnershipColumn: whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING updated_on")
	}

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
	}

	return buildQueryTest(proj, dbvendor, models.DataType{Name: GetOAuth2ClientPalabra()}, "UpdateOAuth2Client", qb, expectedArgs, callArgs, false, false, false, false, false)
}

func buildTestDB_UpdateOAuth2Client(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	queryTail := ""
	var (
		mockDBExpect        jen.Code
		errFuncExpectMethod string
	)
	if isPostgres(dbvendor) {
		queryTail = " RETURNING updated_on"
		mockDBExpect = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WillReturnRows").Callln(
			jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("updated_on"))).Dot("AddRow").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
		)
		errFuncExpectMethod = "ExpectQuery"
	} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		mockDBExpect = jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
			Dotln("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.One()), jen.One())
		errFuncExpectMethod = "ExpectExec"
	}

	lines := []jen.Code{
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
				utils.BuildFakeVar(proj, "OAuth2Client"),
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
				utils.BuildFakeVar(proj, "OAuth2Client"),
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
	}

	return lines
}

func buildTestDB_buildArchiveOAuth2ClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Update(oauth2ClientsTableName).
		Set("updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Set("archived_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id":                              whateverValue,
			oauth2ClientsTableOwnershipColumn: whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING archived_on")
	}

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
	}

	return buildQueryTest(proj, dbvendor, models.DataType{Name: GetOAuth2ClientPalabra()}, "ArchiveOAuth2Client", qb, expectedArgs, callArgs, false, false, false, false, false)
}

func buildTestDB_ArchiveOAuth2Client(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	queryTail := ""
	if isPostgres(dbvendor) {
		queryTail = " RETURNING archived_on"
	}

	lines := []jen.Code{
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
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("exampleUserID"), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID")).
					Dotln("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.One()), jen.One()),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("ArchiveOAuth2Client").Call(utils.CtxVar(), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"), jen.ID("exampleUserID")),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error writing to database",
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("exampleUserID"), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("ArchiveOAuth2Client").Call(utils.CtxVar(), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"), jen.ID("exampleUserID")),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}
