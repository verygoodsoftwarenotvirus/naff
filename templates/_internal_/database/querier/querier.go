package querier

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func querierDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("idRetrievalStrategy").ID("int"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Const().Defs(
			jen.ID("name").Op("=").Lit("db_client"),
			jen.ID("loggerName").Op("=").ID("name"),
			jen.ID("tracingName").Op("=").ID("name"),
			jen.ID("defaultBatchSize").Op("=").Lit(1000),
			jen.ID("DefaultIDRetrievalStrategy").ID("idRetrievalStrategy").Op("=").ID("iota"),
			jen.ID("ReturningStatementIDRetrievalStrategy"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("database").Dot("DataManager").Op("=").Parens(jen.Op("*").ID("SQLQuerier")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("SQLQuerier").Struct(
				jen.ID("config").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/config", "Config"),
				jen.ID("db").Op("*").Qual("database/sql", "DB"),
				jen.ID("sqlQueryBuilder").ID("querybuilding").Dot("SQLQueryBuilder"),
				jen.ID("timeFunc").Func().Params().Params(jen.ID("uint64")),
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("migrateOnce").Qual("sync", "Once"),
				jen.ID("idStrategy").ID("idRetrievalStrategy"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideDatabaseClient provides a new DataManager client."),
		jen.Line(),
		jen.Func().ID("ProvideDatabaseClient").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("db").Op("*").Qual("database/sql", "DB"), jen.ID("cfg").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/config", "Config"), jen.ID("sqlQueryBuilder").ID("querybuilding").Dot("SQLQueryBuilder"), jen.ID("shouldCreateTestUser").ID("bool")).Params(jen.ID("database").Dot("DataManager"), jen.ID("error")).Body(
			jen.ID("tracer").Op(":=").ID("tracing").Dot("NewTracer").Call(jen.ID("tracingName")),
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("c").Op(":=").Op("&").ID("SQLQuerier").Valuesln(jen.ID("db").Op(":").ID("db"), jen.ID("config").Op(":").ID("cfg"), jen.ID("tracer").Op(":").ID("tracer"), jen.ID("timeFunc").Op(":").ID("defaultTimeFunc"), jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("loggerName")), jen.ID("sqlQueryBuilder").Op(":").ID("sqlQueryBuilder"), jen.ID("idStrategy").Op(":").ID("DefaultIDRetrievalStrategy")),
			jen.If(jen.ID("cfg").Dot("Provider").Op("==").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/config", "PostgresProvider")).Body(
				jen.ID("c").Dot("idStrategy").Op("=").ID("ReturningStatementIDRetrievalStrategy")),
			jen.If(jen.ID("cfg").Dot("Debug")).Body(
				jen.ID("c").Dot("logger").Dot("SetLevel").Call(jen.ID("logging").Dot("DebugLevel"))),
			jen.If(jen.ID("cfg").Dot("RunMigrations")).Body(
				jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("migrating querier")),
				jen.Var().Defs(
					jen.ID("testUser").Op("*").ID("types").Dot("TestUserCreationConfig"),
				),
				jen.If(jen.ID("shouldCreateTestUser")).Body(
					jen.ID("testUser").Op("=").ID("cfg").Dot("CreateTestUser")),
				jen.If(jen.ID("err").Op(":=").ID("c").Dot("Migrate").Call(
					jen.ID("ctx"),
					jen.ID("cfg").Dot("MaxPingAttempts"),
					jen.ID("testUser"),
				), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("migrating database"),
					))),
				jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("querier migrated!")),
			),
			jen.Return().List(jen.ID("c"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("IsReady is a simple wrapper around the core querier IsReady call."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("IsReady").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("maxAttempts").ID("uint8")).Params(jen.ID("ready").ID("bool")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("attemptCount").Op(":=").Lit(0),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("interval").Op(":").Qual("time", "Second").Dot("String").Call(), jen.Lit("max_attempts").Op(":").ID("maxAttempts"))),
			jen.For(jen.Op("!").ID("ready")).Body(
				jen.ID("err").Op(":=").ID("q").Dot("db").Dot("PingContext").Call(jen.ID("ctx")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("logger").Dot("WithValue").Call(
						jen.Lit("attempt_count"),
						jen.ID("attemptCount"),
					).Dot("Debug").Call(jen.Lit("ping failed, waiting for db")),
					jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
					jen.ID("attemptCount").Op("++"),
					jen.If(jen.ID("attemptCount").Op(">=").ID("int").Call(jen.ID("maxAttempts"))).Body(
						jen.Break()),
				).Else().Body(
					jen.ID("ready").Op("=").ID("true"),
					jen.Return().ID("ready"),
				),
			),
			jen.Return().ID("false"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("defaultTimeFunc").Params().Params(jen.ID("uint64")).Body(
			jen.Return().ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call())),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("currentTime").Params().Params(jen.ID("uint64")).Body(
			jen.If(jen.ID("q").Op("==").ID("nil").Op("||").ID("q").Dot("timeFunc").Op("==").ID("nil")).Body(
				jen.Return().ID("defaultTimeFunc").Call()),
			jen.Return().ID("q").Dot("timeFunc").Call(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("checkRowsForErrorAndClose").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("rows").ID("database").Dot("ResultIterator")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger"),
			jen.If(jen.ID("err").Op(":=").ID("rows").Dot("Err").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("row error"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("row error"),
				),
			),
			jen.If(jen.ID("err").Op(":=").ID("rows").Dot("Close").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("closing database rows"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("closing database rows"),
				),
			),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("rollbackTransaction").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("tx").Op("*").Qual("database/sql", "Tx")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("err").Op(":=").ID("tx").Dot("Rollback").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("q").Dot("logger"),
					jen.ID("span"),
					jen.Lit("rolling back transaction"),
				)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("getIDFromResult").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("database/sql", "Result")).Params(jen.ID("uint64")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.List(jen.ID("id"), jen.ID("err")).Op(":=").ID("res").Dot("LastInsertId").Call(),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("q").Dot("logger").Dot("WithValue").Call(
						jen.ID("keys").Dot("RowIDErrorKey"),
						jen.ID("true"),
					),
					jen.ID("span"),
					jen.Lit("fetching row ID"),
				)),
			jen.Return().ID("uint64").Call(jen.ID("id")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("getOneRow").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("querier").ID("database").Dot("Querier"), jen.List(jen.ID("queryDescription"), jen.ID("query")).ID("string"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Row")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachDatabaseQueryToSpan").Call(
				jen.ID("span"),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s single row fetch query"),
					jen.ID("queryDescription"),
				),
				jen.ID("query"),
				jen.ID("args"),
			),
			jen.Return().ID("querier").Dot("QueryRowContext").Call(
				jen.ID("ctx"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("performReadQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("querier").ID("database").Dot("Querier"), jen.List(jen.ID("queryDescription"), jen.ID("query")).ID("string"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Rows"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("query"),
				jen.ID("query"),
			),
			jen.ID("tracing").Dot("AttachDatabaseQueryToSpan").Call(
				jen.ID("span"),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s fetch query"),
					jen.ID("queryDescription"),
				),
				jen.ID("query"),
				jen.ID("args"),
			),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("querier").Dot("QueryContext").Call(
				jen.ID("ctx"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning user"),
				))),
			jen.If(jen.ID("rowsErr").Op(":=").ID("rows").Dot("Err").Call(), jen.ID("rowsErr").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("rowsErr"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning user"),
				))),
			jen.Return().List(jen.ID("rows"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("performCountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("querier").ID("database").Dot("Querier"), jen.List(jen.ID("query"), jen.ID("queryDesc")).ID("string")).Params(jen.ID("uint64"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachDatabaseQueryToSpan").Call(
				jen.ID("span"),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s count query"),
					jen.ID("queryDesc"),
				),
				jen.ID("query"),
				jen.ID("nil"),
			),
			jen.Var().Defs(
				jen.ID("count").ID("uint64"),
			),
			jen.If(jen.ID("err").Op(":=").ID("q").Dot("getOneRow").Call(
				jen.ID("ctx"),
				jen.ID("querier"),
				jen.ID("queryDesc"),
				jen.ID("query"),
			).Dot("Scan").Call(jen.Op("&").ID("count")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("q").Dot("logger"),
					jen.ID("span"),
					jen.Lit("executing count query"),
				))),
			jen.Return().List(jen.ID("count"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("performBooleanQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("querier").ID("database").Dot("Querier"), jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Params(jen.ID("bool"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Var().Defs(
				jen.ID("exists").ID("bool"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("DatabaseQueryKey"),
				jen.ID("query"),
			),
			jen.ID("tracing").Dot("AttachDatabaseQueryToSpan").Call(
				jen.ID("span"),
				jen.Lit("boolean query"),
				jen.ID("query"),
				jen.ID("args"),
			),
			jen.ID("err").Op(":=").ID("querier").Dot("QueryRowContext").Call(
				jen.ID("ctx"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			).Dot("Scan").Call(jen.Op("&").ID("exists")),
			jen.If(jen.Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.Return().List(jen.ID("false"), jen.ID("nil"))),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("false"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing boolean query"),
				))),
			jen.Return().List(jen.ID("exists"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("performWriteQueryIgnoringReturn").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("querier").ID("database").Dot("Querier"), jen.List(jen.ID("queryDescription"), jen.ID("query")).ID("string"), jen.ID("args").Index().Interface()).Params(jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("performWriteQuery").Call(
				jen.ID("ctx"),
				jen.ID("querier"),
				jen.ID("true"),
				jen.ID("queryDescription"),
				jen.ID("query"),
				jen.ID("args"),
			),
			jen.Return().ID("err"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("performWriteQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("querier").ID("database").Dot("Querier"), jen.ID("ignoreReturn").ID("bool"), jen.List(jen.ID("queryDescription"), jen.ID("query")).ID("string"), jen.ID("args").Index().Interface()).Params(jen.ID("uint64"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("query"),
				jen.ID("query"),
			).Dot("WithValue").Call(
				jen.Lit("description"),
				jen.ID("queryDescription"),
			).Dot("WithValue").Call(
				jen.Lit("args"),
				jen.ID("args"),
			),
			jen.ID("tracing").Dot("AttachDatabaseQueryToSpan").Call(
				jen.ID("span"),
				jen.ID("queryDescription"),
				jen.ID("query"),
				jen.ID("args"),
			),
			jen.If(jen.ID("q").Dot("idStrategy").Op("==").ID("ReturningStatementIDRetrievalStrategy").Op("&&").Op("!").ID("ignoreReturn")).Body(
				jen.Var().Defs(
					jen.ID("id").ID("uint64"),
				),
				jen.If(jen.ID("err").Op(":=").ID("querier").Dot("QueryRowContext").Call(
					jen.ID("ctx"),
					jen.ID("query"),
					jen.ID("args").Op("..."),
				).Dot("Scan").Call(jen.Op("&").ID("id")), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("executing %s query"),
						jen.ID("queryDescription"),
					))),
				jen.ID("logger").Dot("Debug").Call(jen.Lit("query executed successfully")),
				jen.Return().List(jen.ID("id"), jen.ID("nil")),
			).Else().If(jen.ID("q").Dot("idStrategy").Op("==").ID("ReturningStatementIDRetrievalStrategy")).Body(
				jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("querier").Dot("ExecContext").Call(
					jen.ID("ctx"),
					jen.ID("query"),
					jen.ID("args").Op("..."),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("executing %s query"),
						jen.ID("queryDescription"),
					))),
				jen.Var().Defs(
					jen.ID("affectedRowCount").ID("int64"),
				),
				jen.If(jen.List(jen.ID("affectedRowCount"), jen.ID("err")).Op("=").ID("res").Dot("RowsAffected").Call(), jen.ID("affectedRowCount").Op("==").Lit(0).Op("||").ID("err").Op("!=").ID("nil")).Body(
					jen.ID("logger").Dot("Debug").Call(jen.Lit("no rows modified by query")),
					jen.ID("span").Dot("AddEvent").Call(jen.Lit("no_rows_modified")),
					jen.Return().List(jen.Lit(0), jen.Qual("database/sql", "ErrNoRows")),
				),
				jen.ID("logger").Dot("Debug").Call(jen.Lit("query executed successfully")),
				jen.Return().List(jen.Lit(0), jen.ID("nil")),
			),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("querier").Dot("ExecContext").Call(
				jen.ID("ctx"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing query"),
				))),
			jen.If(jen.ID("res").Op("!=").ID("nil")).Body(
				jen.If(jen.List(jen.ID("rowCount"), jen.ID("err")).Op(":=").ID("res").Dot("RowsAffected").Call(), jen.ID("err").Op("==").ID("nil").Op("&&").ID("rowCount").Op("==").Lit(0)).Body(
					jen.ID("logger").Dot("Debug").Call(jen.Lit("no rows modified by query")),
					jen.ID("span").Dot("AddEvent").Call(jen.Lit("no_rows_modified")),
					jen.Return().List(jen.Lit(0), jen.Qual("database/sql", "ErrNoRows")),
				)),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("query executed successfully")),
			jen.Return().List(jen.ID("q").Dot("getIDFromResult").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
