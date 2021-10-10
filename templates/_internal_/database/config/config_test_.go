package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			jen.ID("invalidProvider").Op("=").Lit("blah"),
		),
		jen.Newline(),
	)

	firstSupportedDatabase := ""
	for _, db := range proj.EnabledDatabases() {
		firstSupportedDatabase = db
	}

	code.Add(
		jen.Func().ID("TestConfig_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").IDf("%sProvider", firstSupportedDatabase), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ValidateWithContext").Call(jen.ID("ctx")),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestConfig_ProvideDatabaseConnection").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			utils.ConditionalCode(proj.DatabaseIsEnabled(models.Postgres),
				jen.ID("T").Dot("Run").Call(
					jen.Lit("standard for postgres"),
					jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
						jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("PostgresProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
						jen.Newline(),
						jen.List(jen.ID("db"), jen.ID("err")).Op(":=").ID("ProvideDatabaseConnection").Call(
							jen.ID("logger"),
							jen.ID("cfg"),
						),
						jen.Qual(constants.AssertionLibrary, "NotNil").Call(
							jen.ID("t"),
							jen.ID("db"),
						),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
					),
				),
			),
			jen.Newline(),
			utils.ConditionalCode(proj.DatabaseIsEnabled(models.MySQL),
				jen.ID("T").Dot("Run").Call(
					jen.Lit("standard for mysql"),
					jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
						jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("MariaDBProvider"), jen.ID("ConnectionDetails").Op(":").Lit("dbuser:hunter2@tcp(database:3306)/todo")),

						jen.Newline(),
						jen.List(jen.ID("db"), jen.ID("err")).Op(":=").ID("ProvideDatabaseConnection").Call(
							jen.ID("logger"),
							jen.ID("cfg"),
						),
						jen.Qual(constants.AssertionLibrary, "NotNil").Call(
							jen.ID("t"),
							jen.ID("db"),
						),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("invalidProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),

					jen.Newline(),
					jen.List(jen.ID("db"), jen.ID("err")).Op(":=").ID("ProvideDatabaseConnection").Call(
						jen.ID("logger"),
						jen.ID("cfg"),
					),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestConfig_ProvideDatabasePlaceholderFormat").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for postgres"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("PostgresProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.Newline(),
					jen.List(jen.ID("pf"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideDatabasePlaceholderFormat").Call(),
					jen.Qual(constants.AssertionLibrary, "NotNil").Call(
						jen.ID("t"),
						jen.ID("pf"),
					),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Newline(),
			utils.ConditionalCode(proj.DatabaseIsEnabled(models.MySQL),
				jen.ID("T").Dot("Run").Call(
					jen.Lit("standard for mysql"),
					jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
						jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("MariaDBProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
						jen.Newline(),
						jen.List(jen.ID("pf"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideDatabasePlaceholderFormat").Call(),
						jen.Qual(constants.AssertionLibrary, "NotNil").Call(
							jen.ID("t"),
							jen.ID("pf"),
						),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("invalidProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.Newline(),
					jen.List(jen.ID("pf"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideDatabasePlaceholderFormat").Call(),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(
						jen.ID("t"),
						jen.ID("pf"),
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestConfig_ProvideJSONPluckQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for postgres"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("PostgresProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "NotEmpty").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ProvideJSONPluckQuery").Call(),
					),
				),
			),
			jen.Newline(),
			utils.ConditionalCode(proj.DatabaseIsEnabled(models.MySQL),
				jen.ID("T").Dot("Run").Call(
					jen.Lit("standard for mysql"),
					jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
						jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("MariaDBProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
						jen.Newline(),
						jen.Qual(constants.AssertionLibrary, "NotEmpty").Call(
							jen.ID("t"),
							jen.ID("cfg").Dot("ProvideJSONPluckQuery").Call(),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("invalidProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Empty").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ProvideJSONPluckQuery").Call(),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestConfig_ProvideCurrentUnixTimestampQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for postgres"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("PostgresProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "NotEmpty").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ProvideCurrentUnixTimestampQuery").Call(),
					),
				),
			),
			jen.Newline(),
			utils.ConditionalCode(proj.DatabaseIsEnabled(models.MySQL),
				jen.ID("T").Dot("Run").Call(
					jen.Lit("standard for mysql"),
					jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
						jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("MariaDBProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
						jen.Newline(),
						jen.Qual(constants.AssertionLibrary, "NotEmpty").Call(
							jen.ID("t"),
							jen.ID("cfg").Dot("ProvideCurrentUnixTimestampQuery").Call(),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("invalidProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Empty").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ProvideCurrentUnixTimestampQuery").Call(),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestProvideSessionManager").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("cookieConfig").Op(":=").Qual(proj.AuthServicePackage(), "CookieConfig").Values(),
					jen.ID("cfg").Op(":=").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("PostgresProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.Newline(),
					jen.List(jen.ID("sessionManager"), jen.ID("err")).Op(":=").ID("ProvideSessionManager").Call(
						jen.ID("cookieConfig"),
						jen.ID("cfg"),
						jen.ID("nil"),
					),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(
						jen.ID("t"),
						jen.ID("sessionManager"),
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard for postgres"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("cookieConfig").Op(":=").Qual(proj.AuthServicePackage(), "CookieConfig").Values(),
					jen.ID("cfg").Op(":=").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("PostgresProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.Newline(),
					jen.List(jen.ID("sessionManager"), jen.ID("err")).Op(":=").ID("ProvideSessionManager").Call(
						jen.ID("cookieConfig"),
						jen.ID("cfg"),
						jen.Op("&").Qual("database/sql", "DB").Values(),
					),
					jen.Qual(constants.AssertionLibrary, "NotNil").Call(
						jen.ID("t"),
						jen.ID("sessionManager"),
					),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Newline(),
			utils.ConditionalCode(proj.DatabaseIsEnabled(models.MySQL),
				jen.ID("T").Dot("Run").Call(
					jen.Lit("standard for mysql"),
					jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
						jen.ID("cookieConfig").Op(":=").Qual(proj.AuthServicePackage(), "CookieConfig").Values(),
						jen.ID("cfg").Op(":=").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("MariaDBProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
						jen.Newline(),
						jen.List(jen.ID("db"), jen.ID("mockDB"), jen.ID("err")).Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "New").Call(),
						jen.Qual(constants.MustAssertPkg, "NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.Qual(constants.MustAssertPkg, "NotNil").Call(
							jen.ID("t"),
							jen.ID("mockDB"),
						),
						jen.Newline(),
						jen.ID("mockDB").Dot("ExpectQuery").Call(jen.Lit("SELECT VERSION()")).Dot("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("version"))).Dot("AddRow").Call(jen.Lit("1.2.3"))),
						jen.Newline(),
						jen.List(jen.ID("sessionManager"), jen.ID("err")).Op(":=").ID("ProvideSessionManager").Call(
							jen.ID("cookieConfig"),
							jen.ID("cfg"),
							jen.ID("db"),
						),
						jen.Qual(constants.AssertionLibrary, "NotNil").Call(
							jen.ID("t"),
							jen.ID("sessionManager"),
						),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("mockDB").Dot("ExpectationsWereMet").Call(),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("cookieConfig").Op(":=").Qual(proj.AuthServicePackage(), "CookieConfig").Values(),
					jen.ID("cfg").Op(":=").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("invalidProvider"), jen.ID("ConnectionDetails").Op(":").Lit("example_connection_string")),
					jen.Newline(),
					jen.List(jen.ID("sessionManager"), jen.ID("err")).Op(":=").ID("ProvideSessionManager").Call(
						jen.ID("cookieConfig"),
						jen.ID("cfg"),
						jen.Op("&").Qual("database/sql", "DB").Values(),
					),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(
						jen.ID("t"),
						jen.ID("sessionManager"),
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Newline(),
	)

	return code
}
