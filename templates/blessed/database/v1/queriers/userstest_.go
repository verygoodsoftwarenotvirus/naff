package queriers

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/Masterminds/squirrel"
)

const (
	usersTableName = "users"
)

var (
	usersTableColumns = []string{
		fmt.Sprintf("%s.id", usersTableName),
		fmt.Sprintf("%s.username", usersTableName),
		fmt.Sprintf("%s.hashed_password", usersTableName),
		fmt.Sprintf("%s.password_last_changed_on", usersTableName),
		fmt.Sprintf("%s.two_factor_secret", usersTableName),
		fmt.Sprintf("%s.is_admin", usersTableName),
		fmt.Sprintf("%s.created_on", usersTableName),
		fmt.Sprintf("%s.updated_on", usersTableName),
		fmt.Sprintf("%s.archived_on", usersTableName),
	}
)

func usersTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, ret)

	ret.Add(buildBuildMockRowsFromUser(proj, dbvendor)...)
	ret.Add(buildBuildErroneousMockRowFromUser(proj, dbvendor)...)
	ret.Add(buildTestDB_buildGetUserQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_GetUser(proj, dbvendor)...)
	ret.Add(buildTestDB_buildGetUsersQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_GetUsers(proj, dbvendor)...)
	ret.Add(buildTestDB_buildGetUserByUsernameQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_GetUserByUsername(proj, dbvendor)...)
	ret.Add(buildTestDB_buildGetAllUserCountQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_GetAllUserCount(proj, dbvendor)...)
	ret.Add(buildTestDB_buildCreateUserQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_CreateUser(proj, dbvendor)...)
	ret.Add(buildTestDB_buildUpdateUserQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_UpdateUser(proj, dbvendor)...)
	ret.Add(buildTestDB_buildArchiveUserQuery(proj, dbvendor)...)
	ret.Add(buildTestDB_ArchiveUser(proj, dbvendor)...)

	return ret
}

func buildBuildMockRowsFromUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildMockRowsFromUser").Params(
			jen.ID("users").Spread().PointerTo().Qual(proj.ModelsV1Package(), "User"),
		).Params(
			jen.PointerTo().Qual("github.com/DATA-DOG/go-sqlmock", "Rows"),
		).Block(
			jen.ID("includeCount").Assign().Len(jen.ID("users")).GreaterThan().One(),
			jen.ID("columns").Assign().ID("usersTableColumns"),
			jen.Line(),
			jen.If(jen.ID("includeCount")).Block(
				utils.AppendItemsToList(jen.ID("columns"), jen.Lit("count")),
			),
			jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("columns")),
			jen.Line(),
			jen.For(
				jen.List(jen.Underscore(), jen.ID("user")).Assign().Range().ID("users"),
			).Block(
				jen.ID("rowValues").Assign().Index().Qual("database/sql/driver", "Value").Valuesln(
					jen.ID("user").Dot("ID"),
					jen.ID("user").Dot("Username"),
					jen.ID("user").Dot("HashedPassword"),
					jen.ID("user").Dot("PasswordLastChangedOn"),
					jen.ID("user").Dot("TwoFactorSecret"),
					jen.ID("user").Dot("IsAdmin"),
					jen.ID("user").Dot("CreatedOn"),
					jen.ID("user").Dot("UpdatedOn"),
					jen.ID("user").Dot("ArchivedOn"),
				),
				jen.Line(),
				jen.If(jen.ID("includeCount")).Block(
					utils.AppendItemsToList(jen.ID("rowValues"), jen.Len(jen.ID("users"))),
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

func buildBuildErroneousMockRowFromUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildErroneousMockRowFromUser").Params(jen.ID("user").PointerTo().Qual(proj.ModelsV1Package(), "User")).Params(jen.PointerTo().Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("usersTableColumns")).Dot("AddRow").Callln(
				jen.ID("user").Dot("ArchivedOn"),
				jen.ID("user").Dot("ID"),
				jen.ID("user").Dot("Username"),
				jen.ID("user").Dot("HashedPassword"),
				jen.ID("user").Dot("PasswordLastChangedOn"),
				jen.ID("user").Dot("TwoFactorSecret"),
				jen.ID("user").Dot("IsAdmin"),
				jen.ID("user").Dot("CreatedOn"),
				jen.ID("user").Dot("UpdatedOn"),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDB_buildGetUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", usersTableName): whateverValue,
		})

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
	}
	callArgs := []jen.Code{}

	return buildQueryTest(proj, dbvendor, models.DataType{Name: GetUserPalabra()}, "GetUser", qb, expectedArgs, callArgs, true, false, false, true, false, false, nil)
}

func buildTestDB_GetUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetUser", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE id = %s", getIncIndex(dbvendor, 0)),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "User").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Username").MapAssign().Lit("username"),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowFromUser").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUser").Call(utils.CtxVar(), jen.ID("expected").Dot("ID")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"surfaces sql.ErrNoRows",
				jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "User").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Username").MapAssign().Lit("username"),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUser").Call(utils.CtxVar(), jen.ID("expected").Dot("ID")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDB_buildGetUsersQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(append(usersTableColumns, fmt.Sprintf(countQuery, usersTableName))...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", usersTableName): nil,
		})

	expectedArgs := []jen.Code{}
	callArgs := []jen.Code{}

	qb = applyFleshedOutQueryFilter(qb, usersTableName)
	expectedArgs = appendFleshedOutQueryFilterArgs(expectedArgs)

	return buildQueryTest(proj, dbvendor, models.DataType{Name: GetUserPalabra()}, "GetUsers", qb, expectedArgs, callArgs, true, true, true, false, false, false, nil)
}

func buildTestDB_GetUsers(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetUsers", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedUsersQuery").Assign().Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE archived_on IS NULL LIMIT 20"),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expectedCountQuery").Assign().Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.ID("expectedCount").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserList").Valuesln(
					jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().Add(utils.FakeUint64Func()),
						jen.ID("Limit").MapAssign().Lit(20),
						jen.ID("TotalCount").MapAssign().ID("expectedCount")),
					jen.ID("Users").MapAssign().Index().Qual(proj.ModelsV1Package(), "User").Valuesln(
						jen.Valuesln(
							jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
							jen.ID("Username").MapAssign().Lit("username"),
						),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromUser").Call(jen.AddressOf().ID("expected").Dot("Users").Index(jen.Zero())),
					jen.ID("buildMockRowFromUser").Call(jen.AddressOf().ID("expected").Dot("Users").Index(jen.Zero())),
					jen.ID("buildMockRowFromUser").Call(jen.AddressOf().ID("expected").Dot("Users").Index(jen.Zero())),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUsers").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"surfaces sql.ErrNoRows",
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUsers").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call()),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error querying database",
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUsers").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call()),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with erroneous response from database",
				jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserList").Valuesln(
					jen.ID("Users").MapAssign().Index().Qual(proj.ModelsV1Package(), "User").Valuesln(
						jen.Valuesln(
							jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
							jen.ID("Username").MapAssign().Lit("username")),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromUser").Call(jen.AddressOf().ID("expected").Dot("Users").Index(jen.Zero()))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUsers").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call()),
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

func buildTestDB_buildGetUserByUsernameQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.username", usersTableName): whateverValue,
		})

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("Username"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("Username"),
	}

	return buildQueryTest(proj, dbvendor, models.DataType{Name: GetUserPalabra()}, "GetUserByUsername", qb, expectedArgs, callArgs, true, false, false, true, false, true, nil)
}

func buildTestDB_GetUserByUsername(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetUserByUsername", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE username = %s", getIncIndex(dbvendor, 0)),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "User").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Username").MapAssign().Lit("username"),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("Username")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowFromUser").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUserByUsername").Call(utils.CtxVar(), jen.ID("expected").Dot("Username")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"surfaces sql.ErrNoRows",
				jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "User").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Username").MapAssign().Lit("username"),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("Username")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUserByUsername").Call(utils.CtxVar(), jen.ID("expected").Dot("Username")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error querying database",
				jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "User").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Username").MapAssign().Lit("username"),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("Username")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUserByUsername").Call(utils.CtxVar(), jen.ID("expected").Dot("Username")),
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

func buildTestDB_buildGetAllUserCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(fmt.Sprintf(countQuery, usersTableName)).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", usersTableName): nil,
		})

	expectedArgs := []jen.Code{}
	callArgs := []jen.Code{}

	return buildQueryTest(proj, dbvendor, models.DataType{Name: GetUserPalabra()}, "GetAllUserCount", qb, expectedArgs, callArgs, false, false, false, false, false, false, nil)
}

func buildTestDB_GetAllUserCount(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetAllUserCount", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllUserCount").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error querying database",
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllUserCount").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call()),
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

func buildTestDB_buildCreateUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Insert(usersTableName).
		Columns(
			"username",
			"hashed_password",
			"two_factor_secret",
			"is_admin",
		).
		Values(
			whateverValue,
			whateverValue,
			whateverValue,
			false,
		)

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING id, created_on")
	}

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("Username"),
		jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
		jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecret"),
		jen.ID(utils.BuildFakeVarName("User")).Dot("IsAdmin"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("Input")),
	}

	return buildQueryTest(proj,
		dbvendor,
		models.DataType{Name: GetUserPalabra()},
		"CreateUser",
		qb,
		expectedArgs,
		callArgs,
		true,
		false,
		false,
		false,
		true,
		false,
		[]jen.Code{
			utils.BuildFakeVarWithCustomName(proj, "exampleInput", "BuildFakeUserDatabaseCreationInputFromUser", jen.ID(utils.BuildFakeVarName("User"))),
		},
	)
}

func buildTestDB_CreateUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	queryTail := ""
	if isPostgres(dbvendor) {
		queryTail = " RETURNING id, created_on"
	}

	var (
		createdOnCol, createdOnVal string
	)

	if isMariaDB(dbvendor) {
		createdOnCol = ",created_on"
		createdOnVal = ",UNIX_TIMESTAMP()"
	}

	var specialSnowflakePGTest jen.Code
	if isPostgres(dbvendor) {
		specialSnowflakePGTest = utils.BuildSubTest(
			"with postgres row exists error",
			jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "User").Valuesln(
				jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
				jen.ID("Username").MapAssign().Lit("username"),
				jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
			),
			jen.ID("expectedInput").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserInput").Valuesln(
				jen.ID("Username").MapAssign().ID("expected").Dot("Username"),
			),
			jen.ID("expectedQuery").Assign().Litf("INSERT INTO users (username,hashed_password,two_factor_secret,is_admin%s) VALUES (%s,%s,%s,%s%s)%s",
				createdOnCol,
				getIncIndex(dbvendor, 0),
				getIncIndex(dbvendor, 1),
				getIncIndex(dbvendor, 2),
				getIncIndex(dbvendor, 3),
				createdOnVal,
				queryTail,
			),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
				jen.ID("expected").Dot("Username"),
				jen.ID("expected").Dot("HashedPassword"),
				jen.ID("expected").Dot("TwoFactorSecret"),
				jen.ID("expected").Dot("IsAdmin"),
			).Dot("WillReturnError").Call(jen.AddressOf().Qual("github.com/lib/pq", "Error").Valuesln(
				jen.ID("Code").MapAssign().Qual("github.com/lib/pq", "ErrorCode").Call(jen.Lit("23505")),
			)),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateUser").Call(utils.CtxVar(), jen.ID("expectedInput")),
			utils.AssertError(jen.Err(), nil),
			utils.AssertNil(jen.ID("actual"), nil),
			utils.AssertEqual(jen.Err(), jen.Qual(proj.DatabaseV1Package("client"), "ErrUserExists"), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)
	} else {
		specialSnowflakePGTest = jen.Null()
	}
	queryTail = ""

	if isPostgres(dbvendor) {
		queryTail = " RETURNING id, created_on"
	}

	var badPathExpectFuncName string
	if isPostgres(dbvendor) {
		badPathExpectFuncName = "ExpectQuery"
	} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		badPathExpectFuncName = "ExpectExec"
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_CreateUser", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("INSERT INTO users (username,hashed_password,two_factor_secret,is_admin%s) VALUES (%s,%s,%s,%s%s)%s",
				createdOnCol,
				getIncIndex(dbvendor, 0),
				getIncIndex(dbvendor, 1),
				getIncIndex(dbvendor, 2),
				getIncIndex(dbvendor, 3),
				createdOnVal,
				queryTail),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				func() []jen.Code {
					out := []jen.Code{
						jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "User").Valuesln(
							jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
							jen.ID("Username").MapAssign().Lit("username"),
							jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
						),
						jen.ID("expectedInput").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserInput").Valuesln(
							jen.ID("Username").MapAssign().ID("expected").Dot("Username"),
						),
					}

					var expectMethodName, returnMethodName string
					if isPostgres(dbvendor) {
						expectMethodName = "ExpectQuery"
						returnMethodName = "WillReturnRows"
						out = append(out, jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("ID"), jen.Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call())))
					} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
						expectMethodName = "ExpectExec"
						returnMethodName = "WillReturnResult"
						out = append(
							out,
							jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.Add(utils.FakeUint64Func())),
							jen.Line(),
						)
					}

					out = append(out,
						jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
						jen.ID("mockDB").Dot(expectMethodName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
							jen.ID("expected").Dot("Username"),
							jen.ID("expected").Dot("HashedPassword"),
							jen.ID("expected").Dot("TwoFactorSecret"),
							jen.ID("expected").Dot("IsAdmin")).
							Dot(returnMethodName).Call(jen.ID("exampleRows")),
						jen.Line(),
					)

					if isSqlite(dbvendor) || isMariaDB(dbvendor) {
						out = append(out,
							jen.ID("expectedTimeQuery").Assign().Litf("SELECT created_on FROM users WHERE id = %s", getIncIndex(dbvendor, 0)),
							jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).
								Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("CreatedOn"))),
						)
					}

					out = append(out,
						jen.Line(),
						jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateUser").Call(utils.CtxVar(), jen.ID("expectedInput")),
						utils.AssertNoError(jen.Err(), nil),
						utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
						jen.Line(),
						utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
					)

					return out
				}()...,
			),
			jen.Line(),
			specialSnowflakePGTest,
			jen.Line(),
			utils.BuildSubTest(
				"with error querying database",
				jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "User").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Username").MapAssign().Lit("username"),
					jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedInput").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserInput").Valuesln(
					jen.ID("Username").MapAssign().ID("expected").Dot("Username"),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(badPathExpectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Username"),
					jen.ID("expected").Dot("HashedPassword"),
					jen.ID("expected").Dot("TwoFactorSecret"),
					jen.ID("expected").Dot("IsAdmin"),
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateUser").Call(utils.CtxVar(), jen.ID("expectedInput")),
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

func buildTestDB_buildUpdateUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Update(usersTableName).
		Set("username", whateverValue).
		Set("hashed_password", whateverValue).
		Set("two_factor_secret", whateverValue).
		Set("updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id": whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING updated_on")
	}

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("Username"),
		jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
		jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecret"),
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")),
	}

	return buildQueryTest(proj,
		dbvendor,
		models.DataType{Name: GetUserPalabra()},
		"UpdateUser",
		qb,
		expectedArgs,
		callArgs,
		true,
		false,
		false,
		true,
		false,
		true,
		nil,
	)
}

func buildTestDB_UpdateUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	queryTail := ""
	var updateUserExpectMethod, updateUserReturnMethod string
	if isPostgres(dbvendor) {
		queryTail = " RETURNING updated_on"
		updateUserExpectMethod = "ExpectQuery"
		updateUserReturnMethod = "WillReturnRows"
	} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		updateUserExpectMethod = "ExpectExec"
		updateUserReturnMethod = "WillReturnResult"
	}

	buildUpdateUserExampleRows := func() jen.Code {
		if isPostgres(dbvendor) {
			return jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("updated_on"))).Dot("AddRow").Call(jen.Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			return jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.Add(utils.FakeUint64Func()))
		}
		return jen.Null()
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_UpdateUser", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "User").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Username").MapAssign().Lit("username"),
					jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				buildUpdateUserExampleRows(),
				jen.ID("expectedQuery").Assign().Litf("UPDATE users SET username = %s, hashed_password = %s, two_factor_secret = %s, updated_on = %s WHERE id = %s%s",
					getIncIndex(dbvendor, 0),
					getIncIndex(dbvendor, 1),
					getIncIndex(dbvendor, 2),
					getTimeQuery(dbvendor),
					getIncIndex(dbvendor, 3),
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(updateUserExpectMethod).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Username"),
					jen.ID("expected").Dot("HashedPassword"),
					jen.ID("expected").Dot("TwoFactorSecret"),
					jen.ID("expected").Dot("ID"),
				).Dot(updateUserReturnMethod).Call(jen.ID("exampleRows")),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("UpdateUser").Call(utils.CtxVar(), jen.ID("expected")),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDB_buildArchiveUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Update(usersTableName).
		Set("updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Set("archived_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id": whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING archived_on")
	}

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
	}
	callArgs := []jen.Code{}

	return buildQueryTest(proj,
		dbvendor,
		models.DataType{Name: GetUserPalabra()},
		"ArchiveUser",
		qb,
		expectedArgs,
		callArgs,
		true,
		false,
		false,
		true,
		false,
		false,
		nil,
	)
}

func buildTestDB_ArchiveUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	queryTail := ""
	if isPostgres(dbvendor) {
		queryTail = " RETURNING archived_on"
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_ArchiveUser", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expected").Assign().AddressOf().Qual(proj.ModelsV1Package(), "User").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Username").MapAssign().Lit("username"),
					jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedQuery").Assign().Litf("UPDATE users SET updated_on = %s, archived_on = %s WHERE id = %s%s",
					getTimeQuery(dbvendor),
					getTimeQuery(dbvendor),
					getIncIndex(dbvendor, 0),
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("ID")).
					Dotln("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.Add(utils.FakeUint64Func()), jen.Add(utils.FakeUint64Func()))),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("ArchiveUser").Call(utils.CtxVar(), jen.ID("expected").Dot("ID")),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}
