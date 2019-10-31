package queriers

import (
	"path/filepath"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
)

func usersTestDotGo(pkgRoot string, vendor *wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(vendor.SingularPackageName())

	utils.AddImports(ret)
	sn := vendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))
	dbrn := vendor.RouteName()

	ret.Add(
		jen.Func().ID("buildMockRowFromUser").Params(jen.ID("user").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "User")).Params(jen.Op("*").Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("usersTableColumns")).Dot("AddRow").Callln(
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
		jen.Func().ID("buildErroneousMockRowFromUser").Params(jen.ID("user").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "User")).Params(jen.Op("*").Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("usersTableColumns")).Dot("AddRow").Callln(
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
		jen.Func().IDf("Test%s_buildGetUserQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE id = %s", getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetUserQuery").Call(jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetUser", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE id = %s", getIncIndex(dbrn, 0)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "User").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Username").Op(":").Lit("username"),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("ID")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowFromUser").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetUser").Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE id = %s", getIncIndex(dbrn, 0)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "User").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Username").Op(":").Lit("username"),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetUser").Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetUsersQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(0),
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetUsersQuery").Call(jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call()),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetUsers", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.ID("expectedUsersQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "UserList").Valuesln(
					jen.ID("Pagination").Op(":").Qual(filepath.Join(pkgRoot, "models/v1"), "Pagination").Valuesln(
						jen.ID("Page").Op(":").Lit(1),
						jen.ID("Limit").Op(":").Lit(20),
						jen.ID("TotalCount").Op(":").ID("expectedCount")),
					jen.ID("Users").Op(":").Index().Qual(filepath.Join(pkgRoot, "models/v1"), "User").Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(123),
							jen.ID("Username").Op(":").Lit("username"),
						),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromUser").Call(jen.Op("&").ID("expected").Dot("Users").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromUser").Call(jen.Op("&").ID("expected").Dot("Users").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromUser").Call(jen.Op("&").ID("expected").Dot("Users").Index(jen.Lit(0))),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).Dot("WillReturnRows").Callln(
					jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetUsers").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUsersQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetUsers").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUsersQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetUsers").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous response from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUsersQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "UserList").Valuesln(
					jen.ID("Users").Op(":").Index().Qual(filepath.Join(pkgRoot, "models/v1"), "User").Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(123),
							jen.ID("Username").Op(":").Lit("username")),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRowFromUser").Call(jen.Op("&").ID("expected").Dot("Users").Index(jen.Lit(0)))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetUsers").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.ID("expectedUsersQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "UserList").Valuesln(
					jen.ID("Pagination").Op(":").Qual(filepath.Join(pkgRoot, "models/v1"), "Pagination").Valuesln(
						jen.ID("Page").Op(":").Lit(1),
						jen.ID("Limit").Op(":").Lit(20),
						jen.ID("TotalCount").Op(":").ID("expectedCount"),
					),
					jen.ID("Users").Op(":").Index().Qual(filepath.Join(pkgRoot, "models/v1"), "User").Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(123),
							jen.ID("Username").Op(":").Lit("username"),
						),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).Dot("WillReturnRows").Callln(
					jen.ID("buildMockRowFromUser").Call(jen.Op("&").ID("expected").Dot("Users").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromUser").Call(jen.Op("&").ID("expected").Dot("Users").Index(jen.Lit(0))),
					jen.ID("buildMockRowFromUser").Call(jen.Op("&").ID("expected").Dot("Users").Index(jen.Lit(0))),
				),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetUsers").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetUserByUsernameQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("expectedUsername").Op(":=").Lit("username"),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE username = %s", getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetUserByUsernameQuery").Call(jen.ID("expectedUsername")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedUsername"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("string"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetUserByUsername", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE username = %s", getIncIndex(dbrn, 0)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "User").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Username").Op(":").Lit("username"),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("Username")).
					Dotln("WillReturnRows").Call(jen.ID("buildMockRowFromUser").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetUserByUsername").Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("Username")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE username = %s", getIncIndex(dbrn, 0)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "User").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Username").Op(":").Lit("username"),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("Username")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetUserByUsername").Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("Username")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE username = %s", getIncIndex(dbrn, 0)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "User").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Username").Op(":").Lit("username"),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("Username")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetUserByUsername").Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("Username")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetUserCountQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(0),
				jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetUserCountQuery").Call(jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call()),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetUserCount", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetUserCount").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetUserCount").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Zero").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildCreateUserQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "UserInput").Valuesln(
					jen.ID("Username").Op(":").Lit("username"),
					jen.ID("Password").Op(":").Lit("hashed password"),
					jen.ID("TwoFactorSecret").Op(":").Lit("two factor secret"),
				),
				jen.ID("expectedArgCount").Op(":=").Lit(4),
				jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO users (username,hashed_password,two_factor_secret,is_admin) VALUES (%s,%s,%s,%s) RETURNING id, created_on", getIncIndex(dbrn, 0), getIncIndex(dbrn, 1), getIncIndex(dbrn, 2), getIncIndex(dbrn, 3)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildCreateUserQuery").Call(jen.ID("exampleUser")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			)),
		),
		jen.Line(),
	)

	var specialSnowflakePGTest jen.Code
	if dbrn == "postgres" {
		specialSnowflakePGTest = jen.ID("T").Dot("Run").Call(jen.Litf("with %s row exists error", vendor.SingularCommonName()), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "User").Valuesln(
				jen.ID("ID").Op(":").Lit(123),
				jen.ID("Username").Op(":").Lit("username"),
				jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
			),
			jen.ID("expectedInput").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "UserInput").Valuesln(
				jen.ID("Username").Op(":").ID("expected").Dot("Username"),
			),
			jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO users (username,hashed_password,two_factor_secret,is_admin) VALUES (%s,%s,%s,%s) RETURNING id, created_on", getIncIndex(dbrn, 0), getIncIndex(dbrn, 1), getIncIndex(dbrn, 2), getIncIndex(dbrn, 3)),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
				jen.ID("expected").Dot("Username"),
				jen.ID("expected").Dot("HashedPassword"),
				jen.ID("expected").Dot("TwoFactorSecret"),
				jen.ID("expected").Dot("IsAdmin"),
			).Dot("WillReturnError").Call(jen.Op("&").Qual("github.com/lib/pq", "Error").Valuesln(
				jen.ID("Code").Op(":").Qual("github.com/lib/pq", "ErrorCode").Call(jen.Lit("23505")),
			)),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("CreateUser").Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("err"), jen.Qual(filepath.Join(pkgRoot, "database/v1/client"), "ErrUserExists")),
			jen.Line(),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		))
	} else {
		specialSnowflakePGTest = jen.Null()
	}

	ret.Add(
		jen.Func().IDf("Test%s_CreateUser", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "User").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Username").Op(":").Lit("username"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedInput").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "UserInput").Valuesln(
					jen.ID("Username").Op(":").ID("expected").Dot("Username"),
				),
				jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("ID"), jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call())),
				jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO users (username,hashed_password,two_factor_secret,is_admin) VALUES (%s,%s,%s,%s) RETURNING id, created_on", getIncIndex(dbrn, 0), getIncIndex(dbrn, 1), getIncIndex(dbrn, 2), getIncIndex(dbrn, 3)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Username"),
					jen.ID("expected").Dot("HashedPassword"),
					jen.ID("expected").Dot("TwoFactorSecret"),
					jen.ID("expected").Dot("IsAdmin")).
					Dot("WillReturnRows").Call(jen.ID("exampleRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("CreateUser").Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			specialSnowflakePGTest,
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "User").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Username").Op(":").Lit("username"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedInput").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "UserInput").Valuesln(
					jen.ID("Username").Op(":").ID("expected").Dot("Username"),
				),
				jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO users (username,hashed_password,two_factor_secret,is_admin) VALUES (%s,%s,%s,%s) RETURNING id, created_on", getIncIndex(dbrn, 0), getIncIndex(dbrn, 1), getIncIndex(dbrn, 2), getIncIndex(dbrn, 3)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Username"),
					jen.ID("expected").Dot("HashedPassword"),
					jen.ID("expected").Dot("TwoFactorSecret"),
					jen.ID("expected").Dot("IsAdmin"),
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot("CreateUser").Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildUpdateUserQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUser").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "User").Valuesln(
					jen.ID("ID").Op(":").Lit(321),
					jen.ID("Username").Op(":").Lit("username"),
					jen.ID("HashedPassword").Op(":").Lit("hashed password"), jen.ID("TwoFactorSecret").Op(":").Lit("two factor secret"),
				),
				jen.ID("expectedArgCount").Op(":=").Lit(4),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE users SET username = %s, hashed_password = %s, two_factor_secret = %s, updated_on = %s WHERE id = %s RETURNING updated_on", getIncIndex(dbrn, 0), getIncIndex(dbrn, 1), getIncIndex(dbrn, 2), getTimeQuery(dbrn), getIncIndex(dbrn, 3)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildUpdateUserQuery").Call(jen.ID("exampleUser")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_UpdateUser", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "User").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Username").Op(":").Lit("username"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("updated_on"))).Dot("AddRow").Call(jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call())),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE users SET username = %s, hashed_password = %s, two_factor_secret = %s, updated_on = %s WHERE id = %s RETURNING updated_on", getIncIndex(dbrn, 0), getIncIndex(dbrn, 1), getIncIndex(dbrn, 2), getTimeQuery(dbrn), getIncIndex(dbrn, 3)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Username"),
					jen.ID("expected").Dot("HashedPassword"),
					jen.ID("expected").Dot("TwoFactorSecret"),
					jen.ID("expected").Dot("ID"),
				).Dot("WillReturnRows").Call(jen.ID("exampleRows")),
				jen.Line(),
				jen.ID("err").Op(":=").ID(dbfl).Dot("UpdateUser").Call(jen.Qual("context", "Background").Call(), jen.ID("expected")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildArchiveUserQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE users SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE id = %s RETURNING archived_on", getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildArchiveUserQuery").Call(jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_ArchiveUser", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "User").Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Username").Op(":").Lit("username"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE users SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE id = %s RETURNING archived_on", getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("ID")).
					Dotln("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.Lit(1), jen.Lit(1))),
				jen.Line(),
				jen.ID("err").Op(":=").ID(dbfl).Dot("ArchiveUser").Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)
	return ret
}
