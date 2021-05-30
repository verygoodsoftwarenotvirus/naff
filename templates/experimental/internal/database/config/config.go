package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("PostgresProvider").Op("=").Lit("postgres"),
			jen.ID("MariaDBProvider").Op("=").Lit("mariadb"),
			jen.ID("SqliteProvider").Op("=").Lit("sqlite"),
			jen.ID("DefaultMetricsCollectionInterval").Op("=").Lit(2).Op("*").Qual("time", "Second"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("errInvalidDatabase").Op("=").Qual("errors", "New").Call(jen.Lit("invalid database")),
			jen.ID("errNilDBProvided").Op("=").Qual("errors", "New").Call(jen.Lit("invalid DB connection provided")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("Config").Struct(
				jen.ID("CreateTestUser").Op("*").ID("types").Dot("TestUserCreationConfig"),
				jen.ID("Provider").ID("string"),
				jen.ID("ConnectionDetails").ID("database").Dot("ConnectionDetails"),
				jen.ID("MetricsCollectionInterval").Qual("time", "Duration"),
				jen.ID("Debug").ID("bool"),
				jen.ID("RunMigrations").ID("bool"),
				jen.ID("MaxPingAttempts").ID("uint8"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("Config")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates an DatabaseSettings struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("Config")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("cfg"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("cfg").Dot("ConnectionDetails"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("cfg").Dot("Provider"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "In").Call(
						jen.ID("PostgresProvider"),
						jen.ID("MariaDBProvider"),
						jen.ID("SqliteProvider"),
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
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideDatabaseConnection provides a database implementation dependent on the configuration."),
		jen.Line(),
		jen.Func().ID("ProvideDatabaseConnection").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("cfg").Op("*").ID("Config")).Params(jen.Op("*").Qual("database/sql", "DB"), jen.ID("error")).Body(
			jen.Switch(jen.ID("cfg").Dot("Provider")).Body(
				jen.Case(jen.ID("PostgresProvider")).Body(
					jen.Return().ID("postgres").Dot("ProvidePostgresDB").Call(
						jen.ID("logger"),
						jen.ID("cfg").Dot("ConnectionDetails"),
					)),
				jen.Case(jen.ID("MariaDBProvider")).Body(
					jen.Return().ID("mariadb").Dot("ProvideMariaDBConnection").Call(
						jen.ID("logger"),
						jen.ID("cfg").Dot("ConnectionDetails"),
					)),
				jen.Case(jen.ID("SqliteProvider")).Body(
					jen.Return().ID("sqlite").Dot("ProvideSqliteDB").Call(
						jen.ID("logger"),
						jen.ID("cfg").Dot("ConnectionDetails"),
						jen.ID("cfg").Dot("MetricsCollectionInterval"),
					)),
				jen.Default().Body(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("%w: %q"),
						jen.ID("errInvalidDatabase"),
						jen.ID("cfg").Dot("Provider"),
					))),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideDatabasePlaceholderFormat provides ."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("Config")).ID("ProvideDatabasePlaceholderFormat").Params().Params(jen.ID("squirrel").Dot("PlaceholderFormat"), jen.ID("error")).Body(
			jen.Switch(jen.ID("cfg").Dot("Provider")).Body(
				jen.Case(jen.ID("PostgresProvider")).Body(
					jen.Return().List(jen.ID("squirrel").Dot("Dollar"), jen.ID("nil"))),
				jen.Case(jen.ID("MariaDBProvider"), jen.ID("SqliteProvider")).Body(
					jen.Return().List(jen.ID("squirrel").Dot("Question"), jen.ID("nil"))),
				jen.Default().Body(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("%w: %q"),
						jen.ID("errInvalidDatabase"),
						jen.ID("cfg").Dot("Provider"),
					))),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideJSONPluckQuery provides a query for extracting a value out of a JSON dictionary for a given database."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("Config")).ID("ProvideJSONPluckQuery").Params().Params(jen.ID("string")).Body(
			jen.Switch(jen.ID("cfg").Dot("Provider")).Body(
				jen.Case(jen.ID("PostgresProvider")).Body(
					jen.Return().Lit(`%s.%s->'%s'`)),
				jen.Case(jen.ID("MariaDBProvider")).Body(
					jen.Return().Lit(`JSON_CONTAINS(%s.%s, '%d', '$.%s')`)),
				jen.Case(jen.ID("SqliteProvider")).Body(
					jen.Return().Lit(`json_extract(%s.%s, '$.%s')`)),
				jen.Default().Body(
					jen.Return().Lit("")),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideCurrentUnixTimestampQuery provides a database implementation dependent on the configuration."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("Config")).ID("ProvideCurrentUnixTimestampQuery").Params().Params(jen.ID("string")).Body(
			jen.Switch(jen.ID("cfg").Dot("Provider")).Body(
				jen.Case(jen.ID("PostgresProvider")).Body(
					jen.Return().Lit(`extract(epoch FROM NOW())`)),
				jen.Case(jen.ID("MariaDBProvider")).Body(
					jen.Return().Lit(`UNIX_TIMESTAMP()`)),
				jen.Case(jen.ID("SqliteProvider")).Body(
					jen.Return().Lit(`(strftime('%s','now'))`)),
				jen.Default().Body(
					jen.Return().Lit("")),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideSessionManager provides a session manager based on some settings."),
		jen.Line(),
		jen.Func().Comment("There's not a great place to put this function. I don't think it belongs in Auth because it accepts a DB connection,").Comment("but it obviously doesn't belong in the database package, or maybe it does.").ID("ProvideSessionManager").Params(jen.ID("cookieConfig").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/authentication", "CookieConfig"), jen.ID("dbConf").ID("Config"), jen.ID("db").Op("*").Qual("database/sql", "DB")).Params(jen.Op("*").ID("scs").Dot("SessionManager"), jen.ID("error")).Body(
			jen.ID("sessionManager").Op(":=").ID("scs").Dot("New").Call(),
			jen.If(jen.ID("db").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("errNilDBProvided"))),
			jen.Switch(jen.ID("dbConf").Dot("Provider")).Body(
				jen.Case(jen.ID("PostgresProvider")).Body(
					jen.ID("sessionManager").Dot("Store").Op("=").ID("postgresstore").Dot("New").Call(jen.ID("db"))),
				jen.Case(jen.ID("MariaDBProvider")).Body(
					jen.ID("sessionManager").Dot("Store").Op("=").ID("mysqlstore").Dot("New").Call(jen.ID("db"))),
				jen.Case(jen.ID("SqliteProvider")).Body(
					jen.ID("sessionManager").Dot("Store").Op("=").ID("sqlite3store").Dot("New").Call(jen.ID("db"))),
				jen.Default().Body(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("%w: %q"),
						jen.ID("errInvalidDatabase"),
						jen.ID("dbConf").Dot("Provider"),
					))),
			),
			jen.ID("sessionManager").Dot("Lifetime").Op("=").ID("cookieConfig").Dot("Lifetime"),
			jen.ID("sessionManager").Dot("Cookie").Dot("Name").Op("=").ID("cookieConfig").Dot("Name"),
			jen.ID("sessionManager").Dot("Cookie").Dot("Domain").Op("=").ID("cookieConfig").Dot("Domain"),
			jen.ID("sessionManager").Dot("Cookie").Dot("HttpOnly").Op("=").ID("true"),
			jen.ID("sessionManager").Dot("Cookie").Dot("Path").Op("=").Lit("/"),
			jen.ID("sessionManager").Dot("Cookie").Dot("SameSite").Op("=").Qual("net/http", "SameSiteStrictMode"),
			jen.ID("sessionManager").Dot("Cookie").Dot("Secure").Op("=").ID("cookieConfig").Dot("SecureOnly"),
			jen.Return().List(jen.ID("sessionManager"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
