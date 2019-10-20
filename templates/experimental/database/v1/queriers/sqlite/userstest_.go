package sqlite

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func usersTestDotGo() *jen.File {
	ret := jen.NewFile("sqlite")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("buildMockRowFromUser").Params(jen.ID("user").Op("*").ID("models").Dot(
			"User",
		)).Params(jen.Op("*").ID("sqlmock").Dot(
			"Rows",
		)).Block(
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.ID("usersTableColumns")).Dot(
				"AddRow",
			).Call(jen.ID("user").Dot("ID"),
				jen.ID("user").Dot(
					"Username",
				),
				jen.ID("user").Dot(
					"HashedPassword",
				),
				jen.ID("user").Dot(
					"PasswordLastChangedOn",
				),
				jen.ID("user").Dot(
					"TwoFactorSecret",
				),
				jen.ID("user").Dot(
					"IsAdmin",
				),
				jen.ID("user").Dot("CreatedOn"),
				jen.ID("user").Dot("UpdatedOn"),
				jen.ID("user").Dot("ArchivedOn")),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildErroneousMockRowFromUser").Params(jen.ID("user").Op("*").ID("models").Dot(
			"User",
		)).Params(jen.Op("*").ID("sqlmock").Dot(
			"Rows",
		)).Block(
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.ID("usersTableColumns")).Dot(
				"AddRow",
			).Call(jen.ID("user").Dot("ArchivedOn"),
				jen.ID("user").Dot("ID"),
				jen.ID("user").Dot(
					"Username",
				),
				jen.ID("user").Dot(
					"HashedPassword",
				),
				jen.ID("user").Dot(
					"PasswordLastChangedOn",
				),
				jen.ID("user").Dot(
					"TwoFactorSecret",
				),
				jen.ID("user").Dot(
					"IsAdmin",
				),
				jen.ID("user").Dot("CreatedOn"),
				jen.ID("user").Dot("UpdatedOn")),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestSqlite_buildGetUserQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("s"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE id = ?"),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("s").Dot(
					"buildGetUserQuery",
				).Call(jen.ID("expectedUserID")),
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
		jen.Func().ID("TestSqlite_GetUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE id = ?"),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"User",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(123), jen.ID("Username").Op(":").Lit("username")),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(
					"ExpectQuery",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
					"WithArgs",
				).Call(jen.ID("expected").Dot("ID")).Dot(
					"WillReturnRows",
				).Call(jen.ID("buildMockRowFromUser").Call(jen.ID("expected"))),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"GetUser",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE id = ?"),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"User",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(123), jen.ID("Username").Op(":").Lit("username")),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(
					"ExpectQuery",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
					"WithArgs",
				).Call(jen.ID("expected").Dot("ID")).Dot(
					"WillReturnError",
				).Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"GetUser",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestSqlite_buildGetUsersQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("s"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedArgCount").Op(":=").Lit(0),
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("s").Dot(
					"buildGetUsersQuery",
				).Call(jen.ID("models").Dot(
					"DefaultQueryFilter",
				).Call()),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot(
					"Len",
				).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestSqlite_GetUsers").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.ID("expectedUsersQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"UserList",
				).Valuesln(
					jen.ID("Pagination").Op(":").ID("models").Dot(
						"Pagination",
					).Valuesln(
						jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("TotalCount").Op(":").ID("expectedCount")), jen.ID("Users").Op(":").Index().ID("models").Dot(
						"User",
					).Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(123), jen.ID("Username").Op(":").Lit("username")))),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(
					"ExpectQuery",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).Dot(
					"WillReturnRows",
				).Call(jen.ID("buildMockRowFromUser").Call(jen.Op("&").ID("expected").Dot(
					"Users",
				).Index(jen.Lit(0))), jen.ID("buildMockRowFromUser").Call(jen.Op("&").ID("expected").Dot(
					"Users",
				).Index(jen.Lit(0))), jen.ID("buildMockRowFromUser").Call(jen.Op("&").ID("expected").Dot(
					"Users",
				).Index(jen.Lit(0)))),
				jen.ID("mockDB").Dot(
					"ExpectQuery",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).Dot(
					"WillReturnRows",
				).Call(jen.ID("sqlmock").Dot(
					"NewRows",
				).Call(jen.Index().ID("string").Valuesln(
					jen.Lit("count"))).Dot(
					"AddRow",
				).Call(jen.ID("expectedCount"))),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"GetUsers",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
					"DefaultQueryFilter",
				).Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUsersQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(
					"ExpectQuery",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).Dot(
					"WillReturnError",
				).Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"GetUsers",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
					"DefaultQueryFilter",
				).Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUsersQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(
					"ExpectQuery",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).Dot(
					"WillReturnError",
				).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"GetUsers",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
					"DefaultQueryFilter",
				).Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with erroneous response from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUsersQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"UserList",
				).Valuesln(
					jen.ID("Users").Op(":").Index().ID("models").Dot(
						"User",
					).Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(123), jen.ID("Username").Op(":").Lit("username")))),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(
					"ExpectQuery",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).Dot(
					"WillReturnRows",
				).Call(jen.ID("buildErroneousMockRowFromUser").Call(jen.Op("&").ID("expected").Dot(
					"Users",
				).Index(jen.Lit(0)))),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"GetUsers",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
					"DefaultQueryFilter",
				).Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.ID("expectedUsersQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"UserList",
				).Valuesln(
					jen.ID("Pagination").Op(":").ID("models").Dot(
						"Pagination",
					).Valuesln(
						jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("TotalCount").Op(":").ID("expectedCount")), jen.ID("Users").Op(":").Index().ID("models").Dot(
						"User",
					).Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(123), jen.ID("Username").Op(":").Lit("username")))),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(
					"ExpectQuery",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedUsersQuery"))).Dot(
					"WillReturnRows",
				).Call(jen.ID("buildMockRowFromUser").Call(jen.Op("&").ID("expected").Dot(
					"Users",
				).Index(jen.Lit(0))), jen.ID("buildMockRowFromUser").Call(jen.Op("&").ID("expected").Dot(
					"Users",
				).Index(jen.Lit(0))), jen.ID("buildMockRowFromUser").Call(jen.Op("&").ID("expected").Dot(
					"Users",
				).Index(jen.Lit(0)))),
				jen.ID("mockDB").Dot(
					"ExpectQuery",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).Dot(
					"WillReturnError",
				).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"GetUsers",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
					"DefaultQueryFilter",
				).Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestSqlite_buildGetUserByUsernameQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("s"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedUsername").Op(":=").Lit("username"),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE username = ?"),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("s").Dot(
					"buildGetUserByUsernameQuery",
				).Call(jen.ID("expectedUsername")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot(
					"Len",
				).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedUsername"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("string"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestSqlite_GetUserByUsername").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE username = ?"),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"User",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(123), jen.ID("Username").Op(":").Lit("username")),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(
					"ExpectQuery",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
					"WithArgs",
				).Call(jen.ID("expected").Dot(
					"Username",
				)).Dot(
					"WillReturnRows",
				).Call(jen.ID("buildMockRowFromUser").Call(jen.ID("expected"))),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"GetUserByUsername",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
					"Username",
				)),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE username = ?"),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"User",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(123), jen.ID("Username").Op(":").Lit("username")),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(
					"ExpectQuery",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
					"WithArgs",
				).Call(jen.ID("expected").Dot(
					"Username",
				)).Dot(
					"WillReturnError",
				).Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"GetUserByUsername",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
					"Username",
				)),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Lit("SELECT id, username, hashed_password, password_last_changed_on, two_factor_secret, is_admin, created_on, updated_on, archived_on FROM users WHERE username = ?"),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"User",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(123), jen.ID("Username").Op(":").Lit("username")),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(
					"ExpectQuery",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
					"WithArgs",
				).Call(jen.ID("expected").Dot(
					"Username",
				)).Dot(
					"WillReturnError",
				).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"GetUserByUsername",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
					"Username",
				)),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestSqlite_buildGetUserCountQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("s"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedArgCount").Op(":=").Lit(0),
				jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("s").Dot(
					"buildGetUserCountQuery",
				).Call(jen.ID("models").Dot(
					"DefaultQueryFilter",
				).Call()),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot(
					"Len",
				).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestSqlite_GetUserCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(
					"ExpectQuery",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
					"WillReturnRows",
				).Call(jen.ID("sqlmock").Dot(
					"NewRows",
				).Call(jen.Index().ID("string").Valuesln(
					jen.Lit("count"))).Dot(
					"AddRow",
				).Call(jen.ID("expected"))),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"GetUserCount",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
					"DefaultQueryFilter",
				).Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM users WHERE archived_on IS NULL LIMIT 20"),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(
					"ExpectQuery",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
					"WillReturnError",
				).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"GetUserCount",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
					"DefaultQueryFilter",
				).Call()),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot(
					"Zero",
				).Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestSqlite_buildCreateUserQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("s"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
					"UserInput",
				).Valuesln(
					jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("hashed password"), jen.ID("TwoFactorSecret").Op(":").Lit("two factor secret")),
				jen.ID("expectedArgCount").Op(":=").Lit(4),
				jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO users (username,hashed_password,two_factor_secret,is_admin) VALUES (?,?,?,?)"),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("s").Dot(
					"buildCreateUserQuery",
				).Call(jen.ID("exampleUser")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot(
					"Len",
				).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestSqlite_CreateUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"User",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(123), jen.ID("Username").Op(":").Lit("username"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
						"Unix",
					).Call())),
				jen.ID("expectedInput").Op(":=").Op("&").ID("models").Dot(
					"UserInput",
				).Valuesln(
					jen.ID("Username").Op(":").ID("expected").Dot(
						"Username",
					)),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO users (username,hashed_password,two_factor_secret,is_admin) VALUES (?,?,?,?)"),
				jen.ID("mockDB").Dot(
					"ExpectExec",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
					"WithArgs",
				).Call(jen.ID("expected").Dot(
					"Username",
				),
					jen.ID("expected").Dot(
						"HashedPassword",
					),
					jen.ID("expected").Dot(
						"TwoFactorSecret",
					),
					jen.ID("expected").Dot(
						"IsAdmin",
					)).Dot(
					"WillReturnResult",
				).Call(jen.ID("sqlmock").Dot(
					"NewResult",
				).Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.Lit(1))),
				jen.ID("expectedTimeQuery").Op(":=").Lit(`SELECT created_on FROM users WHERE id = ?`),
				jen.ID("mockDB").Dot(
					"ExpectQuery",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).Dot(
					"WillReturnRows",
				).Call(jen.ID("sqlmock").Dot(
					"NewRows",
				).Call(jen.Index().ID("string").Valuesln(
					jen.Lit("created_on"))).Dot(
					"AddRow",
				).Call(jen.ID("expected").Dot("CreatedOn"))),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"CreateUser",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"User",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(123), jen.ID("Username").Op(":").Lit("username"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
						"Unix",
					).Call())),
				jen.ID("expectedInput").Op(":=").Op("&").ID("models").Dot(
					"UserInput",
				).Valuesln(
					jen.ID("Username").Op(":").ID("expected").Dot(
						"Username",
					)),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO users (username,hashed_password,two_factor_secret,is_admin) VALUES (?,?,?,?)"),
				jen.ID("mockDB").Dot(
					"ExpectExec",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
					"WithArgs",
				).Call(jen.ID("expected").Dot(
					"Username",
				),
					jen.ID("expected").Dot(
						"HashedPassword",
					),
					jen.ID("expected").Dot(
						"TwoFactorSecret",
					),
					jen.ID("expected").Dot(
						"IsAdmin",
					)).Dot(
					"WillReturnError",
				).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"CreateUser",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestSqlite_buildUpdateUserQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("s"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
					"User",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(321), jen.ID("Username").Op(":").Lit("username"), jen.ID("HashedPassword").Op(":").Lit("hashed password"), jen.ID("TwoFactorSecret").Op(":").Lit("two factor secret")),
				jen.ID("expectedArgCount").Op(":=").Lit(4),
				jen.ID("expectedQuery").Op(":=").Lit("UPDATE users SET username = ?, hashed_password = ?, two_factor_secret = ?, updated_on = (strftime('%s','now')) WHERE id = ?"),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("s").Dot(
					"buildUpdateUserQuery",
				).Call(jen.ID("exampleUser")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot(
					"Len",
				).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestSqlite_UpdateUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"User",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(123), jen.ID("Username").Op(":").Lit("username"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
						"Unix",
					).Call())),
				jen.ID("expectedQuery").Op(":=").Lit("UPDATE users SET username = ?, hashed_password = ?, two_factor_secret = ?, updated_on = (strftime('%s','now')) WHERE id = ?"),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(
					"ExpectExec",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
					"WithArgs",
				).Call(jen.ID("expected").Dot(
					"Username",
				),
					jen.ID("expected").Dot(
						"HashedPassword",
					),
					jen.ID("expected").Dot(
						"TwoFactorSecret",
					),
					jen.ID("expected").Dot("ID")).Dot(
					"WillReturnResult",
				).Call(jen.ID("sqlmock").Dot(
					"NewResult",
				).Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.Lit(1))),
				jen.ID("err").Op(":=").ID("s").Dot(
					"UpdateUser",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("expected")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestSqlite_buildArchiveUserQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("s"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Lit("UPDATE users SET updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE id = ?"),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("s").Dot(
					"buildArchiveUserQuery",
				).Call(jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot(
					"Len",
				).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestSqlite_ArchiveUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"User",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(123), jen.ID("Username").Op(":").Lit("username"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
						"Unix",
					).Call())),
				jen.ID("expectedQuery").Op(":=").Lit("UPDATE users SET updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE id = ?"),
				jen.List(jen.ID("s"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot(
					"ExpectExec",
				).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
					"WithArgs",
				).Call(jen.ID("expected").Dot("ID")).Dot(
					"WillReturnResult",
				).Call(jen.ID("sqlmock").Dot(
					"NewResult",
				).Call(jen.Lit(1), jen.Lit(1))),
				jen.ID("err").Op(":=").ID("s").Dot(
					"ArchiveUser",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
					"ExpectationsWereMet",
				).Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)
	return ret
}
