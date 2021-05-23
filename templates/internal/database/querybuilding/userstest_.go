package querybuilding

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"strings"

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
		fmt.Sprintf("%s.salt", usersTableName),
		fmt.Sprintf("%s.requires_password_change", usersTableName),
		fmt.Sprintf("%s.password_last_changed_on", usersTableName),
		fmt.Sprintf("%s.two_factor_secret", usersTableName),
		fmt.Sprintf("%s.two_factor_secret_verified_on", usersTableName),
		fmt.Sprintf("%s.is_admin", usersTableName),
		fmt.Sprintf("%s.created_on", usersTableName),
		fmt.Sprintf("%s.last_updated_on", usersTableName),
		fmt.Sprintf("%s.archived_on", usersTableName),
	}
)

func usersTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	code := jen.NewFilePathName(proj.DatabasePackage("queriers", "v1", spn), spn)

	utils.AddImports(proj, code, false)

	code.Add(buildBuildMockRowsFromUser(proj)...)
	code.Add(buildBuildErroneousMockRowFromUser(proj)...)
	code.Add(buildTestScanUsers(proj, dbvendor)...)
	code.Add(buildTestDB_buildGetUserQuery(proj, dbvendor)...)
	code.Add(buildTestDB_GetUser(proj, dbvendor)...)
	code.Add(buildTestDB_buildGetUserWithUnverifiedTwoFactorSecretQuery(proj, dbvendor)...)
	code.Add(buildTestDB_GetUserWithUnverifiedTwoFactorSecret(proj, dbvendor)...)
	code.Add(buildTestDB_buildGetUsersQuery(proj, dbvendor)...)
	code.Add(buildTestDB_GetUsers(proj, dbvendor)...)
	code.Add(buildTestDB_buildGetUserByUsernameQuery(proj, dbvendor)...)
	code.Add(buildTestDB_GetUserByUsername(proj, dbvendor)...)
	code.Add(buildTestDB_buildGetAllUsersCountQuery(dbvendor)...)
	code.Add(buildTestDB_GetAllUsersCount(dbvendor)...)
	code.Add(buildTestDB_buildCreateUserQuery(proj, dbvendor)...)
	code.Add(buildTestDB_CreateUser(proj, dbvendor)...)
	code.Add(buildTestDB_buildUpdateUserQuery(proj, dbvendor)...)
	code.Add(buildTestDB_UpdateUser(proj, dbvendor)...)
	code.Add(buildTestDB_buildUpdateUserPasswordQuery(proj, dbvendor)...)
	code.Add(buildTestDB_UpdateUserPassword(proj, dbvendor)...)
	code.Add(buildTestDB_buildVerifyUserTwoFactorSecretQuery(proj, dbvendor)...)
	code.Add(buildTestDB_VerifyUserTwoFactorSecret(proj, dbvendor)...)
	code.Add(buildTestDB_buildArchiveUserQuery(proj, dbvendor)...)
	code.Add(buildTestDB_ArchiveUser(proj, dbvendor)...)

	return code
}

func buildBuildMockRowsFromUser(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildMockRowsFromUser").Params(
			jen.ID("users").Spread().PointerTo().Qual(proj.TypesPackage(), "User"),
		).Params(
			jen.PointerTo().Qual("github.com/DATA-DOG/go-sqlmock", "Rows"),
		).Body(
			jen.ID("columns").Assign().ID("usersTableColumns"),
			jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("columns")),
			jen.Line(),
			jen.For(
				jen.List(jen.Underscore(), jen.ID("user")).Assign().Range().ID("users"),
			).Body(
				jen.ID("rowValues").Assign().Index().Qual("database/sql/driver", "Value").Valuesln(
					jen.ID("user").Dot("ID"),
					jen.ID("user").Dot("Username"),
					jen.ID("user").Dot("HashedPassword"),
					jen.ID("user").Dot("Salt"),
					jen.ID("user").Dot("RequiresPasswordChange"),
					jen.ID("user").Dot("PasswordLastChangedOn"),
					jen.ID("user").Dot("TwoFactorSecret"),
					jen.ID("user").Dot("TwoFactorSecretVerifiedOn"),
					jen.ID("user").Dot("IsAdmin"),
					jen.ID("user").Dot("CreatedOn"),
					jen.ID("user").Dot("LastUpdatedOn"),
					jen.ID("user").Dot("ArchivedOn"),
				),
				jen.Line(),
				jen.ID(utils.BuildFakeVarName("Rows")).Dot("AddRow").Call(jen.ID("rowValues").Spread()),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildErroneousMockRowFromUser(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildErroneousMockRowFromUser").Params(jen.ID("user").PointerTo().Qual(proj.TypesPackage(), "User")).Params(jen.PointerTo().Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Body(
			jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("usersTableColumns")).Dot("AddRow").Callln(
				jen.ID("user").Dot("ArchivedOn"),
				jen.ID("user").Dot("ID"),
				jen.ID("user").Dot("Username"),
				jen.ID("user").Dot("HashedPassword"),
				jen.ID("user").Dot("Salt"),
				jen.ID("user").Dot("RequiresPasswordChange"),
				jen.ID("user").Dot("PasswordLastChangedOn"),
				jen.ID("user").Dot("TwoFactorSecret"),
				jen.ID("user").Dot("TwoFactorSecretVerifiedOn"),
				jen.ID("user").Dot("IsAdmin"),
				jen.ID("user").Dot("CreatedOn"),
				jen.ID("user").Dot("LastUpdatedOn"),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	}

	return lines
}

func buildTestScanUsers(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	lines := []jen.Code{
		jen.Func().IDf("Test%s_ScanUsers", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"surfaces row errors",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockRows").Assign().AddressOf().Qual(proj.DatabasePackage(), "MockResultIterator").Values(),
				jen.Line(),
				jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.False()),
				jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(
					jen.Underscore(),
					jen.Err(),
				).Assign().ID(dbfl).Dot("scanUsers").Call(jen.ID("mockRows")),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"logs row closing errors",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockRows").Assign().AddressOf().Qual(proj.DatabasePackage(), "MockResultIterator").Values(),
				jen.Line(),
				jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.False()),
				jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.Nil()),
				jen.ID("mockRows").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(
					jen.Underscore(),
					jen.Err(),
				).Assign().ID(dbfl).Dot("scanUsers").Call(jen.ID("mockRows")),
				utils.AssertNoError(jen.Err(), nil),
			),
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
			fmt.Sprintf("%s.id", usersTableName):          whateverValue,
			fmt.Sprintf("%s.archived_on", usersTableName): nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.two_factor_secret_verified_on", usersTableName): nil,
		})

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "User"),
	}

	return buildQueryTest(dbvendor, "GetUser", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_GetUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", usersTableName):          whateverValue,
			fmt.Sprintf("%s.archived_on", usersTableName): nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.two_factor_secret_verified_on", usersTableName): nil,
		}).
		ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetUser", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("User")).Dot("Salt").Equals().Nil(),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("User")).Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowsFromUser").Call(jen.ID(utils.BuildFakeVarName("User")))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUser").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User")).Dot("ID")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("User")), jen.ID("actual"), nil),
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
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("User")).Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUser").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User")).Dot("ID")),
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

func buildTestDB_buildGetUserWithUnverifiedTwoFactorSecretQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", usersTableName):                            whateverValue,
			fmt.Sprintf("%s.archived_on", usersTableName):                   nil,
			fmt.Sprintf("%s.two_factor_secret_verified_on", usersTableName): nil,
		})

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "User"),
	}

	return buildQueryTest(dbvendor, "GetUserWithUnverifiedTwoFactorSecret", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_GetUserWithUnverifiedTwoFactorSecret(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", usersTableName):                            whateverValue,
			fmt.Sprintf("%s.two_factor_secret_verified_on", usersTableName): nil,
			fmt.Sprintf("%s.archived_on", usersTableName):                   nil,
		}).
		ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetUserWithUnverifiedTwoFactorSecret", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("User")).Dot("Salt").Equals().Nil(),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("User")).Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowsFromUser").Call(jen.ID(utils.BuildFakeVarName("User")))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUserWithUnverifiedTwoFactorSecret").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User")).Dot("ID")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("User")), jen.ID("actual"), nil),
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
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("User")).Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUserWithUnverifiedTwoFactorSecret").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User")).Dot("ID")),
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
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", usersTableName): nil,
		})

	expectedArgs := []jen.Code{}
	callArgs := []jen.Code{
		jen.ID(constants.FilterVarName),
	}

	qb = applyFleshedOutQueryFilter(qb, usersTableName)
	expectedArgs = appendFleshedOutQueryFilterArgs(expectedArgs)
	pql := []jen.Code{
		jen.ID(constants.FilterVarName).Assign().Qual(proj.FakeModelsPackage(), "BuildFleshedOutQueryFilter").Call(),
	}

	return buildQueryTest(dbvendor, "GetUsers", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_GetUsers(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", usersTableName): nil,
		}).
		OrderBy(fmt.Sprintf("%s.id", usersTableName)).
		Limit(20).
		ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetUsers", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedUsersQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				utils.BuildFakeVar(proj, "UserList"),
				jen.ID(utils.BuildFakeVarName("UserList")).Dot("Users").Index(jen.Zero()).Dot("Salt").Equals().Nil(),
				jen.ID(utils.BuildFakeVarName("UserList")).Dot("Users").Index(jen.One()).Dot("Salt").Equals().Nil(),
				jen.ID(utils.BuildFakeVarName("UserList")).Dot("Users").Index(jen.Lit(2)).Dot("Salt").Equals().Nil(),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowsFromUser").Callln(
						jen.AddressOf().ID(utils.BuildFakeVarName("UserList")).Dot("Users").Index(jen.Zero()),
						jen.AddressOf().ID(utils.BuildFakeVarName("UserList")).Dot("Users").Index(jen.One()),
						jen.AddressOf().ID(utils.BuildFakeVarName("UserList")).Dot("Users").Index(jen.Lit(2)),
					),
				), jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUsers").Call(constants.CtxVar(), jen.ID(constants.FilterVarName)),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("UserList")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"surfaces sql.ErrNoRows",
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUsers").Call(constants.CtxVar(), jen.ID(constants.FilterVarName)),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error querying database",
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).
					Dotln("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUsers").Call(constants.CtxVar(), jen.ID(constants.FilterVarName)),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with erroneous response from database",
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromUser").Call(
					jen.Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				)),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUsers").Call(constants.CtxVar(), jen.ID(constants.FilterVarName)),
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
			fmt.Sprintf("%s.username", usersTableName):    whateverValue,
			fmt.Sprintf("%s.archived_on", usersTableName): nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.two_factor_secret_verified_on", usersTableName): nil,
		})

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("Username"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("Username"),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "User"),
	}

	return buildQueryTest(dbvendor, "GetUserByUsername", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_GetUserByUsername(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.username", usersTableName):    whateverValue,
			fmt.Sprintf("%s.archived_on", usersTableName): nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.two_factor_secret_verified_on", usersTableName): nil,
		}).
		ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetUserByUsername", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("User")).Dot("Salt").Equals().Nil(),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("User")).Dot("Username")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowsFromUser").Call(jen.ID(utils.BuildFakeVarName("User")))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUserByUsername").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User")).Dot("Username")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("User")), jen.ID("actual"), nil),
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
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("User")).Dot("Username")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUserByUsername").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User")).Dot("Username")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error querying database",
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("User")).Dot("Username")).
					Dotln("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetUserByUsername").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User")).Dot("Username")),
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

func buildTestDB_buildGetAllUsersCountQuery(dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Select(fmt.Sprintf(countQuery, usersTableName)).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", usersTableName): nil,
		})

	expectedArgs := []jen.Code{}
	callArgs := []jen.Code{}

	return buildQueryTest(dbvendor, "GetAllUsersCount", qb, expectedArgs, callArgs, nil)
}

func buildTestDB_GetAllUsersCount(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).
		Select(fmt.Sprintf(countQuery, usersTableName)).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", usersTableName): nil,
		}).ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetAllUsersCount", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID(utils.BuildFakeVarName("Count")).Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID(utils.BuildFakeVarName("Count")))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllUsersCount").Call(constants.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Count")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error querying database",
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("GetAllUsersCount").Call(constants.CtxVar()),
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
			"salt",
			"two_factor_secret",
			"is_admin",
		).
		Values(
			whateverValue,
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
		jen.ID(utils.BuildFakeVarName("User")).Dot("Salt"),
		jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecret"),
		jen.ID(utils.BuildFakeVarName("User")).Dot("IsAdmin"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("Input")),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "User"),
		utils.BuildFakeVarWithCustomName(
			proj,
			utils.BuildFakeVarName("Input"),
			"BuildFakeUserDatabaseCreationInputFromUser",
			jen.ID(utils.BuildFakeVarName("User")),
		),
	}

	return buildQueryTest(dbvendor, "CreateUser", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_CreateUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	qb := queryBuilderForDatabase(dbvendor).
		Insert(usersTableName).
		Columns(
			"username",
			"hashed_password",
			"salt",
			"two_factor_secret",
			"is_admin",
		).
		Values(
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			false,
		)
	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING id, created_on")
	}
	expectedQuery, _, _ := qb.ToSql()

	var specialSnowflakePGTest jen.Code = jen.Null()
	if isPostgres(dbvendor) {
		specialSnowflakePGTest = utils.BuildSubTest(
			"with postgres row exists error",
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.Line(),
			utils.BuildFakeVar(proj, "User"),
			jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecretVerifiedOn").Equals().Nil(),
			utils.BuildFakeVarWithCustomName(proj, "expectedInput", "BuildFakeUserDatabaseCreationInputFromUser", jen.ID(utils.BuildFakeVarName("User"))),
			jen.Line(),
			jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
				jen.ID(utils.BuildFakeVarName("User")).Dot("Username"),
				jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
				jen.ID(utils.BuildFakeVarName("User")).Dot("Salt"),
				jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecret"),
				jen.False(),
			).Dot("WillReturnError").Call(jen.AddressOf().Qual("github.com/lib/pq", "Error").Values(
				jen.ID("Code").MapAssign().ID("postgresRowExistsErrorCode"),
			)),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateUser").Call(constants.CtxVar(), jen.ID("expectedInput")),
			utils.AssertError(jen.Err(), nil),
			utils.AssertNil(jen.ID("actual"), nil),
			utils.AssertEqual(jen.Err(), jen.Qual(proj.DatabasePackage("client"), "ErrUserExists"), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)
	}

	var badPathExpectFuncName string
	if isPostgres(dbvendor) {
		badPathExpectFuncName = "ExpectQuery"
	} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		badPathExpectFuncName = "ExpectExec"
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_CreateUser", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				func() []jen.Code {
					out := []jen.Code{
						jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
						jen.Line(),
						utils.BuildFakeVar(proj, "User"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecretVerifiedOn").Equals().Nil(),
						jen.ID(utils.BuildFakeVarName("User")).Dot("Salt").Equals().Nil(),
						utils.BuildFakeVarWithCustomName(proj, "expectedInput", "BuildFakeUserDatabaseCreationInputFromUser", jen.ID(utils.BuildFakeVarName("User"))),
						jen.Line(),
					}
					var expectMethodName, returnMethodName string
					if isPostgres(dbvendor) {
						expectMethodName = "ExpectQuery"
						returnMethodName = "WillReturnRows"
						out = append(out,
							jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").
								Call(jen.Index().String().Values(jen.Lit("id"), jen.Lit("created_on"))).
								Dot("AddRow").Call(
								jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
								jen.ID(utils.BuildFakeVarName("User")).Dot("CreatedOn"),
							),
						)
					} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
						expectMethodName = "ExpectExec"
						returnMethodName = "WillReturnResult"
						out = append(out,
							jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID(utils.BuildFakeVarName("User")).Dot("ID")), jen.One()),
						)
					}
					out = append(out,
						jen.ID("mockDB").Dot(expectMethodName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
							jen.ID(utils.BuildFakeVarName("User")).Dot("Username"),
							jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
							jen.ID(utils.BuildFakeVarName("User")).Dot("Salt"),
							jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecret"),
							jen.False(),
						).Dot(returnMethodName).Call(jen.ID(utils.BuildFakeVarName("Rows"))),
						jen.Line(),
					)

					if isSqlite(dbvendor) || isMariaDB(dbvendor) {
						out = append(out,
							jen.IDf("%stt", dbfl).Assign().AddressOf().ID("mockTimeTeller").Values(),
							jen.IDf("%stt", dbfl).Dot("On").Call(jen.Lit("Now")).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("User")).Dot("CreatedOn")),
							jen.ID(dbfl).Dot("timeTeller").Equals().IDf("%stt", dbfl),
							jen.Line(),
						)
					}

					out = append(out,
						jen.Line(),
						jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateUser").Call(constants.CtxVar(), jen.ID("expectedInput")),
						utils.AssertNoError(jen.Err(), nil),
						utils.AssertEqual(jen.ID(utils.BuildFakeVarName("User")), jen.ID("actual"), nil),
						jen.Line(),
						func() jen.Code {
							if isMariaDB(dbvendor) || isSqlite(dbvendor) {
								return utils.AssertExpectationsFor(fmt.Sprintf("%stt", dbfl))
							}
							return jen.Null()
						}(),
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
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecretVerifiedOn").Equals().Nil(),
				utils.BuildFakeVarWithCustomName(proj, "expectedInput", "BuildFakeUserDatabaseCreationInputFromUser", jen.ID(utils.BuildFakeVarName("User"))),
				jen.Line(),
				jen.ID("mockDB").Dot(badPathExpectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID(utils.BuildFakeVarName("User")).Dot("Username"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("Salt"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecret"),
					jen.False(),
				).Dot("WillReturnError").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot("CreateUser").Call(constants.CtxVar(), jen.ID("expectedInput")),
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
		Set("salt", whateverValue).
		Set("two_factor_secret", whateverValue).
		Set("two_factor_secret_verified_on", whateverValue).
		Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id": whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING last_updated_on")
	}

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("Username"),
		jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
		jen.ID(utils.BuildFakeVarName("User")).Dot("Salt"),
		jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecret"),
		jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecretVerifiedOn"),
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "User"),
	}

	return buildQueryTest(dbvendor, "UpdateUser", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_UpdateUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	qb := queryBuilderForDatabase(dbvendor).
		Update(usersTableName).
		Set("username", whateverValue).
		Set("hashed_password", whateverValue).
		Set("salt", whateverValue).
		Set("two_factor_secret", whateverValue).
		Set("two_factor_secret_verified_on", whateverValue).
		Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id": whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING last_updated_on")
	}
	expectedQuery, _, _ := qb.ToSql()

	var (
		updateUserExpectMethod string
		updateUserReturnMethod string
	)
	if isPostgres(dbvendor) {
		updateUserExpectMethod = "ExpectQuery"
		updateUserReturnMethod = "WillReturnRows"
	} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		updateUserExpectMethod = "ExpectExec"
		updateUserReturnMethod = "WillReturnResult"
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_UpdateUser", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.ID("expectedQuery").Assign().Lit(expectedQuery),
				func() jen.Code {
					if isPostgres(dbvendor) {
						return jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("last_updated_on"))).Dot("AddRow").Call(jen.Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
					} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
						return jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID(utils.BuildFakeVarName("User")).Dot("ID")), jen.One())
					}
					// this line can never be tested :(
					panic(fmt.Sprintf("invalid dbvendor: %q", dbvendor))
				}(),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(updateUserExpectMethod).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID(utils.BuildFakeVarName("User")).Dot("Username"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("Salt"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecret"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecretVerifiedOn"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				).Dot(updateUserReturnMethod).Call(jen.ID(utils.BuildFakeVarName("Rows"))),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("UpdateUser").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User"))),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDB_buildUpdateUserPasswordQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Update(usersTableName).
		Set("hashed_password", whateverValue).
		Set("requires_password_change", false).
		Set("password_last_changed_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id": whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING last_updated_on")
	}

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
		jen.False(),
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
		jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "User"),
	}

	return buildQueryTest(dbvendor, "UpdateUserPassword", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_UpdateUserPassword(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	qb := queryBuilderForDatabase(dbvendor).
		Update(usersTableName).
		Set("hashed_password", whateverValue).
		Set("requires_password_change", false).
		Set("password_last_changed_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id": whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING last_updated_on")
	}
	expectedQuery, _, _ := qb.ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_UpdateUserPassword", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.ID("expectedQuery").Assign().Lit(expectedQuery),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
					jen.False(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				).Dot("WillReturnResult").Call(jen.Qual("database/sql/driver", "ResultNoRows")),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("UpdateUserPassword").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
				),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDB_buildVerifyUserTwoFactorSecretQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).
		Update(usersTableName).
		Set("two_factor_secret_verified_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id": whateverValue,
		})

	expectedArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
	}
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "User"),
	}

	return buildQueryTest(dbvendor, "VerifyUserTwoFactorSecret", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_VerifyUserTwoFactorSecret(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	qb := queryBuilderForDatabase(dbvendor).
		Update(usersTableName).
		Set("two_factor_secret_verified_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id": whateverValue,
		})

	expectedQuery, _, _ := qb.ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_VerifyUserTwoFactorSecret", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.ID("expectedQuery").Assign().Lit(expectedQuery),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(
					jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				).
					Dotln("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.One(), jen.One())),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("VerifyUserTwoFactorSecret").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				),
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
	callArgs := []jen.Code{
		jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
	}
	pql := []jen.Code{
		utils.BuildFakeVar(proj, "User"),
	}

	return buildQueryTest(dbvendor, "ArchiveUser", qb, expectedArgs, callArgs, pql)
}

func buildTestDB_ArchiveUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	qb := queryBuilderForDatabase(dbvendor).
		Update(usersTableName).
		Set("archived_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(squirrel.Eq{
			"id": whateverValue,
		})

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING archived_on")
	}
	expectedQuery, _, _ := qb.ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_ArchiveUser", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.ID("expectedQuery").Assign().Lit(expectedQuery),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName("User")).Dot("ID")).
					Dotln("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.One(), jen.One())),
				jen.Line(),
				jen.Err().Assign().ID(dbfl).Dot("ArchiveUser").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User")).Dot("ID")),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}
