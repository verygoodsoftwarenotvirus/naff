package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("invalidProvider").Op("=").Lit("blah"),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestConfig_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("PostgresProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ValidateWithContext").Call(jen.ID("ctx")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestConfig_ProvideDatabaseConnection").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for postgres"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("PostgresProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.List(jen.ID("db"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideDatabaseConnection").Call(jen.ID("logger")),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for sqlite"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("SqliteProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.List(jen.ID("db"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideDatabaseConnection").Call(jen.ID("logger")),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for mariadb"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("MariaDBProvider"), jen.ID("ConnectionDetails").Op(":").Lit("dbuser:hunter2@tcp(database:3306)/todo")),
					jen.List(jen.ID("db"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideDatabaseConnection").Call(jen.ID("logger")),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("invalidProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.List(jen.ID("db"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideDatabaseConnection").Call(jen.ID("logger")),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestConfig_ProvideDatabasePlaceholderFormat").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for postgres"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("PostgresProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.List(jen.ID("pf"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideDatabasePlaceholderFormat").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("pf"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for sqlite"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("SqliteProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.List(jen.ID("pf"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideDatabasePlaceholderFormat").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("pf"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for mariadb"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("MariaDBProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.List(jen.ID("pf"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideDatabasePlaceholderFormat").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("pf"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("invalidProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.List(jen.ID("pf"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideDatabasePlaceholderFormat").Call(),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("pf"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestConfig_ProvideJSONPluckQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for postgres"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("PostgresProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ProvideJSONPluckQuery").Call(),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for sqlite"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("SqliteProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ProvideJSONPluckQuery").Call(),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for mariadb"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("MariaDBProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ProvideJSONPluckQuery").Call(),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("invalidProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ProvideJSONPluckQuery").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestConfig_ProvideCurrentUnixTimestampQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for postgres"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("PostgresProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ProvideCurrentUnixTimestampQuery").Call(),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for sqlite"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("SqliteProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ProvideCurrentUnixTimestampQuery").Call(),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for mariadb"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("MariaDBProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ProvideCurrentUnixTimestampQuery").Call(),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("invalidProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ProvideCurrentUnixTimestampQuery").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideSessionManager").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cookieConfig").Op(":=").ID("authentication").Dot("CookieConfig").Values(),
					jen.ID("cfg").Op(":=").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("PostgresProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.List(jen.ID("sessionManager"), jen.ID("err")).Op(":=").ID("ProvideSessionManager").Call(
						jen.ID("cookieConfig"),
						jen.ID("cfg"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("sessionManager"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for postgres"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cookieConfig").Op(":=").ID("authentication").Dot("CookieConfig").Values(),
					jen.ID("cfg").Op(":=").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("PostgresProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.List(jen.ID("sessionManager"), jen.ID("err")).Op(":=").ID("ProvideSessionManager").Call(
						jen.ID("cookieConfig"),
						jen.ID("cfg"),
						jen.Op("&").Qual("database/sql", "DB").Values(),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("sessionManager"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for sqlite"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cookieConfig").Op(":=").ID("authentication").Dot("CookieConfig").Values(),
					jen.ID("cfg").Op(":=").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("SqliteProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.List(jen.ID("sessionManager"), jen.ID("err")).Op(":=").ID("ProvideSessionManager").Call(
						jen.ID("cookieConfig"),
						jen.ID("cfg"),
						jen.Op("&").Qual("database/sql", "DB").Values(),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("sessionManager"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for mariadb"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cookieConfig").Op(":=").ID("authentication").Dot("CookieConfig").Values(),
					jen.ID("cfg").Op(":=").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("MariaDBProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.List(jen.ID("db"), jen.ID("mockDB"), jen.ID("err")).Op(":=").ID("sqlmock").Dot("New").Call(),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
					jen.ID("mockDB").Dot("ExpectQuery").Call(jen.Lit("SELECT VERSION()")).Dot("WillReturnRows").Call(jen.ID("sqlmock").Dot("NewRows").Call(jen.Index().ID("string").Valuesln(
						jen.Lit("version"))).Dot("AddRow").Call(jen.Lit("1.2.3"))),
					jen.List(jen.ID("sessionManager"), jen.ID("err")).Op(":=").ID("ProvideSessionManager").Call(
						jen.ID("cookieConfig"),
						jen.ID("cfg"),
						jen.ID("db"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("sessionManager"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("mockDB").Dot("ExpectationsWereMet").Call(),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cookieConfig").Op(":=").ID("authentication").Dot("CookieConfig").Values(),
					jen.ID("cfg").Op(":=").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("invalidProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.List(jen.ID("sessionManager"), jen.ID("err")).Op(":=").ID("ProvideSessionManager").Call(
						jen.ID("cookieConfig"),
						jen.ID("cfg"),
						jen.Op("&").Qual("database/sql", "DB").Values(),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("sessionManager"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
