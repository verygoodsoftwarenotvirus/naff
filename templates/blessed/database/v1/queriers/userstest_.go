package queriers

import (
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, ret)
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))
	dbrn := dbvendor.RouteName()

	isPostgres := dbrn == "postgres"
	isSqlite := dbrn == "sqlite"
	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	ret.Add(
		jen.Func().ID("buildMockRowFromUser").Params(jen.ID("user").PointerTo().Qual(proj.ModelsV1Package(), "User")).Params(jen.ParamPointer().Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("usersTableColumns")).Dot("AddRow").Callln(
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
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildErroneousMockRowFromUser").Params(jen.ID("user").PointerTo().Qual(proj.ModelsV1Package(), "User")).Params(jen.ParamPointer().Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
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
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetUserQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedArgCount").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedQuery").Assign().Litf("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE id = %s", getIncIndex(dbvendor, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetUserQuery").Call(jen.ID("expectedUserID")),
				utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
				utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
				utils.AssertEqual(jen.ID("expectedUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64")), nil),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetUser", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE id = %s", getIncIndex(dbvendor, 0)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Valuesln(
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
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Valuesln(
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
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetUsersQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("expectedArgCount").Assign().Lit(0),
				jen.ID("expectedQuery").Assign().Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetUsersQuery").Call(jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call()),
				utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
				utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetUsers", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedUsersQuery").Assign().Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE archived_on IS NULL LIMIT 20"),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expectedCountQuery").Assign().Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.ID("expectedCount").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "UserList").Valuesln(
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
					jen.ID("buildMockRowFromUser").Call(jen.VarPointer().ID("expected").Dot("Users").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromUser").Call(jen.VarPointer().ID("expected").Dot("Users").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromUser").Call(jen.VarPointer().ID("expected").Dot("Users").Index(jen.Lit(0))),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUsers").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
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
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUsers").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call()),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous response from database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "UserList").Valuesln(
					jen.ID("Users").MapAssign().Index().Qual(proj.ModelsV1Package(), "User").Valuesln(
						jen.Valuesln(
							jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
							jen.ID("Username").MapAssign().Lit("username")),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromUser").Call(jen.VarPointer().ID("expected").Dot("Users").Index(jen.Lit(0)))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUsers").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call()),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching count"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expectedCountQuery").Assign().Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.ID("expectedCount").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "UserList").Valuesln(
					jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().Add(utils.FakeUint64Func()),
						jen.ID("Limit").MapAssign().Lit(20),
						jen.ID("TotalCount").MapAssign().ID("expectedCount"),
					),
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
					jen.ID("buildMockRowFromUser").Call(jen.VarPointer().ID("expected").Dot("Users").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromUser").Call(jen.VarPointer().ID("expected").Dot("Users").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromUser").Call(jen.VarPointer().ID("expected").Dot("Users").Index(jen.Lit(0))),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUsers").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call()),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetUserByUsernameQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("expectedUsername").Assign().Lit("username"),
				jen.ID("expectedArgCount").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedQuery").Assign().Litf("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE username = %s", getIncIndex(dbvendor, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetUserByUsernameQuery").Call(jen.ID("expectedUsername")),
				utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
				utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
				utils.AssertEqual(jen.ID("expectedUsername"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("string")), nil),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetUserByUsername", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Litf("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE username = %s", getIncIndex(dbvendor, 0)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Valuesln(
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
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Valuesln(
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
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Valuesln(
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
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetUserCountQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("expectedArgCount").Assign().Lit(0),
				jen.ID("expectedQuery").Assign().Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetUserCountQuery").Call(jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call()),
				utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
				utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetUserCount", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUserCount").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUserCount").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call()),
				utils.AssertError(jen.Err(), nil),
				utils.AssertZero(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	////////////

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
		jen.Func().IDf("Test%s_buildCreateUserQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUser").Assign().VarPointer().Qual(proj.ModelsV1Package(), "UserInput").Valuesln(
					jen.ID("Username").MapAssign().Lit("username"),
					jen.ID("Password").MapAssign().Lit("hashed password"),
					jen.ID("TwoFactorSecret").MapAssign().Lit("two factor secret"),
				),
				jen.ID("expectedArgCount").Assign().Lit(4),
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
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildCreateUserQuery").Call(jen.ID("exampleUser")),
				utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
				utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
			)),
		),
		jen.Line(),
	)
	queryTail = ""

	////////////

	if isPostgres {
		queryTail = " RETURNING id, created_on"
	}

	var specialSnowflakePGTest jen.Code
	if isPostgres {
		specialSnowflakePGTest = jen.ID("T").Dot("Run").Call(jen.Litf("with %s row exists error", dbvendor.SingularCommonName()), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Valuesln(
				jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
				jen.ID("Username").MapAssign().Lit("username"),
				jen.ID("CreatedOn").MapAssign().ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
			),
			jen.ID("expectedInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), "UserInput").Valuesln(
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
			).Dot("WillReturnError").Call(jen.VarPointer().Qual("github.com/lib/pq", "Error").Valuesln(
				jen.ID("Code").MapAssign().Qual("github.com/lib/pq", "ErrorCode").Call(jen.Lit("23505")),
			)),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateUser").Call(utils.CtxVar(), jen.ID("expectedInput")),
			utils.AssertError(jen.Err(), nil),
			utils.AssertNil(jen.ID("actual"), nil),
			utils.AssertEqual(jen.Err(), jen.Qual(proj.DatabaseV1Package("client"), "ErrUserExists"), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		))
	} else {
		specialSnowflakePGTest = jen.Null()
	}
	queryTail = ""

	////////////

	if isPostgres {
		queryTail = " RETURNING id, created_on"
	}

	var badPathExpectFuncName string
	if isPostgres {
		badPathExpectFuncName = "ExpectQuery"
	} else if isSqlite || isMariaDB {
		badPathExpectFuncName = "ExpectExec"
	}

	ret.Add(
		jen.Func().IDf("Test%s_CreateUser", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
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
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				func() []jen.Code {
					out := []jen.Code{
						jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Valuesln(
							jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
							jen.ID("Username").MapAssign().Lit("username"),
							jen.ID("CreatedOn").MapAssign().ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
						),
						jen.ID("expectedInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), "UserInput").Valuesln(
							jen.ID("Username").MapAssign().ID("expected").Dot("Username"),
						),
					}

					var expectMethodName, returnMethodName string
					if isPostgres {
						expectMethodName = "ExpectQuery"
						returnMethodName = "WillReturnRows"
						out = append(out, jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("ID"), jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call())))
					} else if isSqlite || isMariaDB {
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

					if isSqlite || isMariaDB {
						out = append(out,
							jen.ID("expectedTimeQuery").Assign().Litf("SELECT created_on FROM users WHERE id = %s", getIncIndex(dbvendor, 0)),
							jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).
								Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("CreatedOn"))),
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
			)),
			jen.Line(),
			specialSnowflakePGTest,
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Username").MapAssign().Lit("username"),
					jen.ID("CreatedOn").MapAssign().ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), "UserInput").Valuesln(
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
			)),
		),
		jen.Line(),
	)
	queryTail = ""

	////////////

	if isPostgres {
		queryTail = " RETURNING updated_on"
	}

	ret.Add(
		jen.Func().IDf("Test%s_buildUpdateUserQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUser").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Valuesln(
					jen.ID("ID").MapAssign().Lit(321),
					jen.ID("Username").MapAssign().Lit("username"),
					jen.ID("HashedPassword").MapAssign().Lit("hashed password"), jen.ID("TwoFactorSecret").MapAssign().Lit("two factor secret"),
				),
				jen.ID("expectedArgCount").Assign().Lit(4),
				jen.ID("expectedQuery").Assign().Litf("UPDATE users SET username = %s, hashed_password = %s, two_factor_secret = %s, updated_on = %s WHERE id = %s%s",
					getIncIndex(dbvendor, 0),
					getIncIndex(dbvendor, 1),
					getIncIndex(dbvendor, 2),
					getTimeQuery(dbvendor),
					getIncIndex(dbvendor, 3),
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildUpdateUserQuery").Call(jen.ID("exampleUser")),
				utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
				utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
			)),
		),
		jen.Line(),
	)
	queryTail = ""

	////////////

	var updateUserExpectMethod, updateUserReturnMethod string
	if isPostgres {
		queryTail = " RETURNING updated_on"
		updateUserExpectMethod = "ExpectQuery"
		updateUserReturnMethod = "WillReturnRows"
	} else if isSqlite || isMariaDB {
		updateUserExpectMethod = "ExpectExec"
		updateUserReturnMethod = "WillReturnResult"
	}

	buildUpdateUserExampleRows := func() jen.Code {
		if isPostgres {
			return jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("updated_on"))).Dot("AddRow").Call(jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
		} else if isSqlite || isMariaDB {
			return jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.Add(utils.FakeUint64Func()))
		}
		return jen.Null()
	}

	ret.Add(
		jen.Func().IDf("Test%s_UpdateUser", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Username").MapAssign().Lit("username"),
					jen.ID("CreatedOn").MapAssign().ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
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
			)),
		),
		jen.Line(),
	)
	queryTail = ""

	////////////

	if isPostgres {
		queryTail = " RETURNING archived_on"
	}

	ret.Add(
		jen.Func().IDf("Test%s_buildArchiveUserQuery", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedArgCount").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expectedQuery").Assign().Litf("UPDATE users SET updated_on = %s, archived_on = %s WHERE id = %s%s",
					getTimeQuery(dbvendor),
					getTimeQuery(dbvendor),
					getIncIndex(dbvendor, 0),
					queryTail,
				),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dot("buildArchiveUserQuery").Call(jen.ID("exampleUserID")),
				utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
				utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
				utils.AssertEqual(jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64")), nil),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_ArchiveUser", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("Username").MapAssign().Lit("username"),
					jen.ID("CreatedOn").MapAssign().ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
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
			)),
		),
		jen.Line(),
	)
	return ret
}
