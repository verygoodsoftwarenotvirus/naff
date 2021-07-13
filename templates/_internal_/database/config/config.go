package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			utils.ConditionalCode(proj.DatabaseIsEnabled(models.Postgres), jen.Comment("PostgresProvider is the string used to refer to postgres.").Newline().ID("PostgresProvider").Op("=").Lit("postgres")),
			utils.ConditionalCode(proj.DatabaseIsEnabled(models.MariaDB), jen.Comment("MariaDBProvider is the string used to refer to mariaDB.").Newline().ID("MariaDBProvider").Op("=").Lit("mariadb")),
			utils.ConditionalCode(proj.DatabaseIsEnabled(models.Sqlite), jen.Comment("SqliteProvider is the string used to refer to sqlite.").Newline().ID("SqliteProvider").Op("=").Lit("sqlite")),
			jen.Newline(),
			jen.Comment("DefaultMetricsCollectionInterval is the default amount of time we wait between database metrics queries."),
			jen.ID("DefaultMetricsCollectionInterval").Op("=").Lit(2).Op("*").Qual("time", "Second"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("errInvalidDatabase").Op("=").Qual("errors", "New").Call(jen.Lit("invalid database")),
			jen.ID("errNilDBProvided").Op("=").Qual("errors", "New").Call(jen.Lit("invalid DB connection provided")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Type().Defs(
			jen.Comment("Config represents our database configuration."),
			jen.ID("Config").Struct(
				jen.ID("CreateTestUser").Op("*").Qual(proj.TypesPackage(), "TestUserCreationConfig").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("CreateTestUser"), false)),
				jen.ID("Provider").ID("string").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Provider"), false)),
				jen.ID("ConnectionDetails").Qual(proj.DatabasePackage(), "ConnectionDetails").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("ConnectionDetails"), false)),
				jen.ID("MetricsCollectionInterval").Qual("time", "Duration").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("MetricsCollectionInterval"), false)),
				jen.ID("Debug").ID("bool").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Debug"), false)),
				jen.ID("RunMigrations").ID("bool").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("RunMigrations"), false)),
				jen.ID("MaxPingAttempts").ID("uint8").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("MaxPingAttempts"), false)),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("Config")).Call(jen.ID("nil")),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates an DatabaseSettings struct."),
		jen.Newline(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("Config")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Callln(
				jen.ID("ctx"),
				jen.ID("cfg"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("cfg").Dot("ConnectionDetails"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("cfg").Dot("Provider"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "In").Call(
						utils.ConditionalCode(proj.DatabaseIsEnabled(models.Postgres), jen.ID("PostgresProvider")),
						utils.ConditionalCode(proj.DatabaseIsEnabled(models.MariaDB), jen.ID("MariaDBProvider")),
						utils.ConditionalCode(proj.DatabaseIsEnabled(models.Sqlite), jen.ID("SqliteProvider")),
					),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("cfg").Dot("CreateTestUser"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "When").Call(
						jen.ID("cfg").Dot("CreateTestUser").Op("!=").ID("nil"),
						jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					).Dot("Else").Call(jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Nil")),
				),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvideDatabaseConnection provides a database implementation dependent on the configuration."),
		jen.Newline(),
		jen.Func().ID("ProvideDatabaseConnection").Params(jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"), jen.ID("cfg").Op("*").ID("Config")).Params(jen.Op("*").Qual("database/sql", "DB"), jen.ID("error")).Body(
			jen.Switch(jen.ID("cfg").Dot("Provider")).Body(
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.Postgres),
					jen.Case(jen.ID("PostgresProvider")).Body(
						jen.Return().Qual(proj.QuerybuildingPackage("postgres"), "ProvidePostgresDB").Call(
							jen.ID("logger"),
							jen.ID("cfg").Dot("ConnectionDetails"),
						)),
				),
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.MariaDB),
					jen.Case(jen.ID("MariaDBProvider")).Body(
						jen.Return().Qual(proj.QuerybuildingPackage("mariadb"), "ProvideMariaDBConnection").Call(
							jen.ID("logger"),
							jen.ID("cfg").Dot("ConnectionDetails"),
						)),
				),
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.Sqlite),
					jen.Case(jen.ID("SqliteProvider")).Body(
						jen.Return().Qual(proj.QuerybuildingPackage("sqlite"), "ProvideSqliteDB").Call(
							jen.ID("logger"),
							jen.ID("cfg").Dot("ConnectionDetails"),
							jen.ID("cfg").Dot("MetricsCollectionInterval"),
						)),
				),
				jen.Default().Body(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("%w: %q"),
						jen.ID("errInvalidDatabase"),
						jen.ID("cfg").Dot("Provider"),
					))),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvideDatabasePlaceholderFormat provides ."),
		jen.Newline(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("Config")).ID("ProvideDatabasePlaceholderFormat").Params().Params(jen.Qual(constants.SQLGenerationLibrary, "PlaceholderFormat"), jen.ID("error")).Body(
			jen.Switch(jen.ID("cfg").Dot("Provider")).Body(
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.Postgres),
					jen.Case(jen.ID("PostgresProvider")).Body(
						jen.Return().List(jen.Qual(constants.SQLGenerationLibrary, "Dollar"), jen.ID("nil"))),
				),
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.Sqlite) || proj.DatabaseIsEnabled(models.MariaDB),
					jen.Case(
						utils.ConditionalCode(proj.DatabaseIsEnabled(models.MariaDB), jen.ID("MariaDBProvider")),
						utils.ConditionalCode(proj.DatabaseIsEnabled(models.Sqlite), jen.ID("SqliteProvider")),
					).Body(
						jen.Return().List(jen.Qual(constants.SQLGenerationLibrary, "Question"), jen.ID("nil"))),
				),
				jen.Default().Body(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("%w: %q"),
						jen.ID("errInvalidDatabase"),
						jen.ID("cfg").Dot("Provider"),
					))),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvideJSONPluckQuery provides a query for extracting a value out of a JSON dictionary for a given database."),
		jen.Newline(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("Config")).ID("ProvideJSONPluckQuery").Params().Params(jen.ID("string")).Body(
			jen.Switch(jen.ID("cfg").Dot("Provider")).Body(
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.Postgres),
					jen.Case(jen.ID("PostgresProvider")).Body(
						jen.Return().RawString(`%s.%s->'%s'`)),
				),
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.MariaDB),
					jen.Case(jen.ID("MariaDBProvider")).Body(
						jen.Return().RawString(`JSON_CONTAINS(%s.%s, '%d', '$.%s')`)),
				),
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.Sqlite),
					jen.Case(jen.ID("SqliteProvider")).Body(
						jen.Return().RawString(`json_extract(%s.%s, '$.%s')`)),
				),
				jen.Default().Body(
					jen.Return().Lit("")),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvideCurrentUnixTimestampQuery provides a database implementation dependent on the configuration."),
		jen.Newline(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("Config")).ID("ProvideCurrentUnixTimestampQuery").Params().Params(jen.ID("string")).Body(
			jen.Switch(jen.ID("cfg").Dot("Provider")).Body(
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.Postgres),
					jen.Case(jen.ID("PostgresProvider")).Body(
						jen.Return().RawString(`extract(epoch FROM NOW())`)),
				),
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.MariaDB),
					jen.Case(jen.ID("MariaDBProvider")).Body(
						jen.Return().RawString(`UNIX_TIMESTAMP()`)),
				),
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.Sqlite),
					jen.Case(jen.ID("SqliteProvider")).Body(
						jen.Return().RawString(`(strftime('%s','now'))`)),
				),
				jen.Default().Body(
					jen.Return().Lit("")),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvideSessionManager provides a session manager based on some settings."),
		jen.Newline(),
		jen.Comment("There's not a great place to put this function. I don't think it belongs in Auth because it accepts a DB connection,"),
		jen.Newline(),
		jen.Comment("but it obviously doesn't belong in the database package, or maybe it does."),
		jen.Newline(),
		jen.Func().ID("ProvideSessionManager").Params(jen.ID("cookieConfig").Qual(proj.AuthServicePackage(), "CookieConfig"), jen.ID("dbConf").ID("Config"), jen.ID("db").Op("*").Qual("database/sql", "DB")).Params(jen.Op("*").Qual("github.com/alexedwards/scs/v2", "SessionManager"), jen.ID("error")).Body(
			jen.ID("sessionManager").Op(":=").Qual("github.com/alexedwards/scs/v2", "New").Call(),
			jen.Newline(),
			jen.If(jen.ID("db").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("errNilDBProvided")),
			),
			jen.Newline(),
			jen.Switch(jen.ID("dbConf").Dot("Provider")).Body(
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.Postgres),
					jen.Case(jen.ID("PostgresProvider")).Body(
						jen.ID("sessionManager").Dot("Store").Op("=").Qual("github.com/alexedwards/scs/postgresstore", "New").Call(jen.ID("db"))),
				),
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.MariaDB),
					jen.Case(jen.ID("MariaDBProvider")).Body(
						jen.ID("sessionManager").Dot("Store").Op("=").Qual("github.com/alexedwards/scs/mysqlstore", "New").Call(jen.ID("db"))),
				),
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.Sqlite),
					jen.Case(jen.ID("SqliteProvider")).Body(
						jen.ID("sessionManager").Dot("Store").Op("=").Qual("github.com/alexedwards/scs/sqlite3store", "New").Call(jen.ID("db"))),
				),
				jen.Default().Body(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("%w: %q"),
						jen.ID("errInvalidDatabase"),
						jen.ID("dbConf").Dot("Provider"),
					))),
			),
			jen.Newline(),
			jen.ID("sessionManager").Dot("Lifetime").Op("=").ID("cookieConfig").Dot("Lifetime"),
			jen.ID("sessionManager").Dot("Cookie").Dot("Name").Op("=").ID("cookieConfig").Dot("Name"),
			jen.ID("sessionManager").Dot("Cookie").Dot("Domain").Op("=").ID("cookieConfig").Dot("Domain"),
			jen.ID("sessionManager").Dot("Cookie").Dot("HttpOnly").Op("=").ID("true"),
			jen.ID("sessionManager").Dot("Cookie").Dot("Path").Op("=").Lit("/"),
			jen.ID("sessionManager").Dot("Cookie").Dot("SameSite").Op("=").Qual("net/http", "SameSiteStrictMode"),
			jen.ID("sessionManager").Dot("Cookie").Dot("Secure").Op("=").ID("cookieConfig").Dot("SecureOnly"),
			jen.Newline(),
			jen.Return().List(jen.ID("sessionManager"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	return code
}
