package querier

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func querierTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("defaultLimit").Op("=").ID("uint8").Call(jen.Lit(20)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("newCountDBRowResponse").Params(jen.ID("count").ID("uint64")).Params(jen.Op("*").ID("sqlmock").Dot("Rows")).Body(
			jen.Return().ID("sqlmock").Dot("NewRows").Call(jen.Index().ID("string").Valuesln(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("count"))),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("newSuccessfulDatabaseResult").Params(jen.ID("returnID").ID("uint64")).Params(jen.ID("driver").Dot("Result")).Body(
			jen.Return().ID("sqlmock").Dot("NewResult").Call(
				jen.ID("int64").Call(jen.ID("returnID")),
				jen.Lit(1),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("formatQueryForSQLMock").Params(jen.ID("query").ID("string")).Params(jen.ID("string")).Body(
			jen.Return().Qual("strings", "NewReplacer").Call(
				jen.Lit("$"),
				jen.Lit(`\$`),
				jen.Lit("("),
				jen.Lit(`\(`),
				jen.Lit(")"),
				jen.Lit(`\)`),
				jen.Lit("="),
				jen.Lit(`\=`),
				jen.Lit("*"),
				jen.Lit(`\*`),
				jen.Lit("."),
				jen.Lit(`\.`),
				jen.Lit("+"),
				jen.Lit(`\+`),
				jen.Lit("?"),
				jen.Lit(`\?`),
				jen.Lit(","),
				jen.Lit(`\,`),
				jen.Lit("-"),
				jen.Lit(`\-`),
				jen.Lit("["),
				jen.Lit(`\[`),
				jen.Lit("]"),
				jen.Lit(`\]`),
			).Dot("Replace").Call(jen.ID("query"))),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("interfaceToDriverValue").Params(jen.ID("in").Index().Interface()).Params(jen.Index().ID("driver").Dot("Value")).Body(
			jen.ID("out").Op(":=").Index().ID("driver").Dot("Value").Values(),
			jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("in")).Body(
				jen.ID("out").Op("=").ID("append").Call(
					jen.ID("out"),
					jen.ID("driver").Dot("Value").Call(jen.ID("x")),
				)),
			jen.Return().ID("out"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("sqlmockExpecterWrapper").Struct(jen.ID("sqlmock").Dot("Sqlmock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("e").Op("*").ID("sqlmockExpecterWrapper")).ID("AssertExpectations").Params(jen.ID("t").ID("mock").Dot("TestingT")).Params(jen.ID("bool")).Body(
			jen.Return().ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("e").Dot("Sqlmock").Dot("ExpectationsWereMet").Call(),
				jen.Lit("not all database expectations were met"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestClient").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("SQLQuerier"), jen.Op("*").ID("sqlmockExpecterWrapper")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.List(jen.ID("db"), jen.ID("sqlMock"), jen.ID("err")).Op(":=").ID("sqlmock").Dot("New").Call(jen.ID("sqlmock").Dot("MonitorPingsOption").Call(jen.ID("true"))),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("c").Op(":=").Op("&").ID("SQLQuerier").Valuesln(jen.ID("db").Op(":").ID("db"), jen.ID("logger").Op(":").ID("logging").Dot("NewNoopLogger").Call(), jen.ID("timeFunc").Op(":").ID("defaultTimeFunc"), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.Lit("test")), jen.ID("sqlQueryBuilder").Op(":").ID("database").Dot("BuildMockSQLQueryBuilder").Call(), jen.ID("idStrategy").Op(":").ID("DefaultIDRetrievalStrategy")),
			jen.Return().List(jen.ID("c"), jen.Op("&").ID("sqlmockExpecterWrapper").Valuesln(jen.ID("Sqlmock").Op(":").ID("sqlMock"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildErroneousMockRow").Params().Params(jen.Op("*").ID("sqlmock").Dot("Rows")).Body(
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.Index().ID("string").Valuesln(jen.Lit("columns"), jen.Lit("don't"), jen.Lit("match"), jen.Lit("lol"))).Dot("AddRow").Call(
				jen.Lit("doesn't"),
				jen.Lit("matter"),
				jen.Lit("what"),
				jen.Lit("goes"),
			),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("expectAuditLogEntryInTransaction").Params(jen.ID("mockQueryBuilder").Op("*").ID("database").Dot("MockSQLQueryBuilder"), jen.ID("db").ID("sqlmock").Dot("Sqlmock"), jen.ID("returnErr").ID("error")).Body(
			jen.List(jen.ID("fakeAuditLogEntryQuery"), jen.ID("fakeAuditLogEntryArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
			jen.ID("mockQueryBuilder").Dot("AuditLogEntrySQLQueryBuilder").Dot("On").Call(
				jen.Lit("BuildCreateAuditLogEntryQuery"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
				jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("AuditLogEntryCreationInput").Values()),
			).Dot("Return").Call(
				jen.ID("fakeAuditLogEntryQuery"),
				jen.ID("fakeAuditLogEntryArgs"),
			),
			jen.ID("e").Op(":=").ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAuditLogEntryQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAuditLogEntryArgs")).Op("...")),
			jen.If(jen.ID("returnErr").Op("!=").ID("nil")).Body(
				jen.ID("e").Dot("WillReturnError").Call(jen.ID("returnErr"))).Else().Body(
				jen.ID("e").Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.Lit(123)))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_IsReady").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectPing").Call().Dot("WillDelayFor").Call(jen.Lit(0)),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("c").Dot("IsReady").Call(
							jen.ID("ctx"),
							jen.Lit(1),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error pinging database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectPing").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("c").Dot("IsReady").Call(
							jen.ID("ctx"),
							jen.Lit(1),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("exhausting all available queries"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("ctx"), jen.ID("cancel")).Op(":=").Qual("context", "WithTimeout").Call(
						jen.Qual("context", "Background").Call(),
						jen.Qual("time", "Second"),
					),
					jen.Defer().ID("cancel").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("c").Dot("IsReady").Call(
						jen.ID("ctx"),
						jen.Lit(1),
					),
					jen.ID("db").Dot("ExpectPing").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("c").Dot("IsReady").Call(
							jen.ID("ctx"),
							jen.Lit(1),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideDatabaseClient").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Var().Defs(
						jen.ID("migrationFunctionCalled").ID("bool"),
					),
					jen.ID("fakeMigrationFunc").Op(":=").Func().Params().Body(
						jen.ID("migrationFunctionCalled").Op("=").ID("true")),
					jen.List(jen.ID("db"), jen.ID("mockDB"), jen.ID("err")).Op(":=").ID("sqlmock").Dot("New").Call(jen.ID("sqlmock").Dot("MonitorPingsOption").Call(jen.ID("true"))),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("queryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("queryBuilder").Dot("On").Call(
						jen.Lit("BuildMigrationFunc"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").Qual("database/sql", "DB").Values()),
					).Dot("Return").Call(jen.ID("fakeMigrationFunc")),
					jen.ID("mockDB").Dot("ExpectPing").Call().Dot("WillDelayFor").Call(jen.Lit(0)),
					jen.ID("exampleConfig").Op(":=").Op("&").ID("config").Dot("Config").Valuesln(jen.ID("Debug").Op(":").ID("true"), jen.ID("RunMigrations").Op(":").ID("true"), jen.ID("MaxPingAttempts").Op(":").Lit(1)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideDatabaseClient").Call(
						jen.ID("ctx"),
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("db"),
						jen.ID("exampleConfig"),
						jen.ID("queryBuilder"),
						jen.ID("true"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("migrationFunctionCalled"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.Op("&").ID("sqlmockExpecterWrapper").Valuesln(jen.ID("Sqlmock").Op(":").ID("mockDB")),
						jen.ID("queryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with PostgresProvider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Var().Defs(
						jen.ID("migrationFunctionCalled").ID("bool"),
					),
					jen.ID("fakeMigrationFunc").Op(":=").Func().Params().Body(
						jen.ID("migrationFunctionCalled").Op("=").ID("true")),
					jen.List(jen.ID("db"), jen.ID("mockDB"), jen.ID("err")).Op(":=").ID("sqlmock").Dot("New").Call(jen.ID("sqlmock").Dot("MonitorPingsOption").Call(jen.ID("true"))),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("queryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("queryBuilder").Dot("On").Call(
						jen.Lit("BuildMigrationFunc"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").Qual("database/sql", "DB").Values()),
					).Dot("Return").Call(jen.ID("fakeMigrationFunc")),
					jen.ID("mockDB").Dot("ExpectPing").Call().Dot("WillDelayFor").Call(jen.Lit(0)),
					jen.ID("exampleConfig").Op(":=").Op("&").ID("config").Dot("Config").Valuesln(jen.ID("Provider").Op(":").ID("config").Dot("PostgresProvider"), jen.ID("Debug").Op(":").ID("true"), jen.ID("RunMigrations").Op(":").ID("true"), jen.ID("MaxPingAttempts").Op(":").Lit(1)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideDatabaseClient").Call(
						jen.ID("ctx"),
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("db"),
						jen.ID("exampleConfig"),
						jen.ID("queryBuilder"),
						jen.ID("true"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("migrationFunctionCalled"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.Op("&").ID("sqlmockExpecterWrapper").Valuesln(jen.ID("Sqlmock").Op(":").ID("mockDB")),
						jen.ID("queryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error initializing querier"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("db"), jen.ID("mockDB"), jen.ID("err")).Op(":=").ID("sqlmock").Dot("New").Call(jen.ID("sqlmock").Dot("MonitorPingsOption").Call(jen.ID("true"))),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("queryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockDB").Dot("ExpectPing").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("exampleConfig").Op(":=").Op("&").ID("config").Dot("Config").Valuesln(jen.ID("Debug").Op(":").ID("true"), jen.ID("RunMigrations").Op(":").ID("true"), jen.ID("MaxPingAttempts").Op(":").Lit(1)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideDatabaseClient").Call(
						jen.ID("ctx"),
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("db"),
						jen.ID("exampleConfig"),
						jen.ID("queryBuilder"),
						jen.ID("true"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.Op("&").ID("sqlmockExpecterWrapper").Valuesln(jen.ID("Sqlmock").Op(":").ID("mockDB")),
						jen.ID("queryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestDefaultTimeFunc").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotZero").Call(
						jen.ID("t"),
						jen.ID("defaultTimeFunc").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_currentTime").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("c").Dot("currentTime").Call(),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("handles nil"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Var().Defs(
						jen.ID("c").Op("*").ID("SQLQuerier"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("c").Dot("currentTime").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_rollbackTransaction").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.ID("db").Dot("ExpectRollback").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("tx"), jen.ID("err")).Op(":=").ID("c").Dot("db").Dot("BeginTx").Call(
						jen.ID("ctx"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("c").Dot("rollbackTransaction").Call(
						jen.ID("ctx"),
						jen.ID("tx"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_getIDFromResult").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("expected").Op(":=").ID("int64").Call(jen.Lit(123)),
					jen.ID("m").Op(":=").Op("&").ID("database").Dot("MockSQLResult").Values(),
					jen.ID("m").Dot("On").Call(jen.Lit("LastInsertId")).Dot("Return").Call(
						jen.ID("expected"),
						jen.ID("nil"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("actual").Op(":=").ID("c").Dot("getIDFromResult").Call(
						jen.ID("ctx"),
						jen.ID("m"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("uint64").Call(jen.ID("expected")),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("logs error"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("m").Op(":=").Op("&").ID("database").Dot("MockSQLResult").Values(),
					jen.ID("m").Dot("On").Call(jen.Lit("LastInsertId")).Dot("Return").Call(
						jen.ID("int64").Call(jen.Lit(0)),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("actual").Op(":=").ID("c").Dot("getIDFromResult").Call(
						jen.ID("ctx"),
						jen.ID("m"),
					),
					jen.ID("assert").Dot("Zero").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_handleRows").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("mockRows").Op(":=").Op("&").ID("database").Dot("MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.ID("nil")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(jen.ID("nil")),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("checkRowsForErrorAndClose").Call(
						jen.ID("ctx"),
						jen.ID("mockRows"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with row error"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expected").Op(":=").Qual("errors", "New").Call(jen.Lit("blah")),
					jen.ID("mockRows").Op(":=").Op("&").ID("database").Dot("MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.ID("expected")),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("checkRowsForErrorAndClose").Call(
						jen.ID("ctx"),
						jen.ID("mockRows"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.Qual("errors", "Is").Call(
							jen.ID("err"),
							jen.ID("expected"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with close error"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expected").Op(":=").Qual("errors", "New").Call(jen.Lit("blah")),
					jen.ID("mockRows").Op(":=").Op("&").ID("database").Dot("MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.ID("nil")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(jen.ID("expected")),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("checkRowsForErrorAndClose").Call(
						jen.ID("ctx"),
						jen.ID("mockRows"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.Qual("errors", "Is").Call(
							jen.ID("err"),
							jen.ID("expected"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_performCreateQueryIgnoringReturn").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.Lit(1))),
					jen.ID("err").Op(":=").ID("c").Dot("performWriteQueryIgnoringReturn").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("db"),
						jen.Lit("example"),
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_performCreateQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.Lit(1))),
					jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("c").Dot("performWriteQuery").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("db"),
						jen.ID("false"),
						jen.Lit("example"),
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("c").Dot("performWriteQuery").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("db"),
						jen.ID("false"),
						jen.Lit("example"),
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no rows returned"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("sqlmock").Dot("NewResult").Call(
						jen.ID("int64").Call(jen.Lit(1)),
						jen.Lit(0),
					)),
					jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("c").Dot("performWriteQuery").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("db"),
						jen.ID("false"),
						jen.Lit("example"),
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.Qual("errors", "Is").Call(
							jen.ID("err"),
							jen.Qual("database/sql", "ErrNoRows"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with ReturningStatementIDRetrievalStrategy"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("c").Dot("idStrategy").Op("=").ID("ReturningStatementIDRetrievalStrategy"),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("sqlmock").Dot("NewRows").Call(jen.Index().ID("string").Valuesln(jen.Lit("id"))).Dot("AddRow").Call(jen.ID("uint64").Call(jen.Lit(123)))),
					jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("c").Dot("performWriteQuery").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("db"),
						jen.ID("false"),
						jen.Lit("example"),
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with ReturningStatementIDRetrievalStrategy and error"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("c").Dot("idStrategy").Op("=").ID("ReturningStatementIDRetrievalStrategy"),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("id"), jen.ID("err")).Op(":=").ID("c").Dot("performWriteQuery").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("db"),
						jen.ID("false"),
						jen.Lit("example"),
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("assert").Dot("Zero").Call(
						jen.ID("t"),
						jen.ID("id"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("ignoring return with return statement id strategy"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("c").Dot("idStrategy").Op("=").ID("ReturningStatementIDRetrievalStrategy"),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.Lit(1))),
					jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("c").Dot("performWriteQuery").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("db"),
						jen.ID("true"),
						jen.Lit("example"),
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("ignoring return with return statement id strategy and error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("c").Dot("idStrategy").Op("=").ID("ReturningStatementIDRetrievalStrategy"),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("c").Dot("performWriteQuery").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("db"),
						jen.ID("true"),
						jen.Lit("example"),
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("ignoring return with return statement id strategy with no rows affected"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("c").Dot("idStrategy").Op("=").ID("ReturningStatementIDRetrievalStrategy"),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("sqlmock").Dot("NewResult").Call(
						jen.Lit(0),
						jen.Lit(0),
					)),
					jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("c").Dot("performWriteQuery").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("db"),
						jen.ID("true"),
						jen.Lit("example"),
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.Qual("errors", "Is").Call(
							jen.ID("err"),
							jen.Qual("database/sql", "ErrNoRows"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
