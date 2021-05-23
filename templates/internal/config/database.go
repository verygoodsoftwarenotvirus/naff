package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"path/filepath"
)

const (
	SessionManagerLibrary = "github.com/alexedwards/scs"
)

func databaseDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)
	code.ImportAlias("github.com/alexedwards/scs/v2", "scs")

	if proj.DatabaseIsEnabled(models.Postgres) {
		code.ImportName(filepath.Join(SessionManagerLibrary, "postgresstore"), "postgresstore")
	}
	if proj.DatabaseIsEnabled(models.Sqlite) {
		code.ImportName(filepath.Join(SessionManagerLibrary, "sqlite3store"), "sqlite3store")
	}
	if proj.DatabaseIsEnabled(models.MariaDB) {
		code.ImportName(filepath.Join(SessionManagerLibrary, "mysqlstore"), "mysqlstore")
	}

	code.Add(buildDatabaseConstantDeclarations(proj)...)
	code.Add(buildProvideDatabaseConnection(proj)...)
	code.Add(buildProvideDatabaseClient(proj)...)
	code.Add(buildProvideSessionManager(proj)...)

	return code
}

func buildDatabaseConstantDeclarations(proj *models.Project) []jen.Code {
	var constDefs []jen.Code
	if proj.DatabaseIsEnabled(models.Postgres) {
		constDefs = append(constDefs,
			jen.Comment("PostgresProviderKey is the string we use to refer to postgres"),
			jen.ID("PostgresProviderKey").Equals().Lit("postgres"),
		)
	}
	if proj.DatabaseIsEnabled(models.Sqlite) {
		constDefs = append(constDefs,
			jen.Comment("MariaDBProviderKey is the string we use to refer to mariaDB"),
			jen.ID("MariaDBProviderKey").Equals().Lit("mariadb"),
		)
	}
	if proj.DatabaseIsEnabled(models.MariaDB) {
		constDefs = append(constDefs,
			jen.Comment("SqliteProviderKey is the string we use to refer to sqlite"),
			jen.ID("SqliteProviderKey").Equals().Lit("sqlite"),
		)
	}

	lines := []jen.Code{
		jen.Const().Defs(constDefs...),
		jen.Line(),
	}

	return lines
}

func buildProvideDatabaseConnection(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideDatabaseConnection provides a database implementation dependent on the configuration."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").PointerTo().ID("ServerConfig")).ID("ProvideDatabaseConnection").Params(
			proj.LoggerParam(),
		).Params(
			jen.PointerTo().Qual("database/sql", "DB"),
			jen.Error(),
		).Body(
			jen.Switch(jen.ID("cfg").Dot("Database").Dot("Provider")).Body(
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.Postgres) {
						return jen.Case(jen.ID("PostgresProviderKey")).Body(
							jen.Return(jen.Qual(proj.DatabasePackage("queriers", "postgres"), "ProvidePostgresDB").Call(
								jen.ID(constants.LoggerVarName),
								jen.ID("cfg").Dot("Database").Dot("ConnectionDetails"),
							)),
						)
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.MariaDB) {
						return jen.Case(jen.ID("MariaDBProviderKey")).Body(
							jen.Return(jen.Qual(proj.DatabasePackage("queriers", "mariadb"), "ProvideMariaDBConnection").Call(
								jen.ID(constants.LoggerVarName),
								jen.ID("cfg").Dot("Database").Dot("ConnectionDetails"),
							)),
						)
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.Sqlite) {
						return jen.Case(jen.ID("SqliteProviderKey")).Body(
							jen.Return(jen.Qual(proj.DatabasePackage("queriers", "sqlite"), "ProvideSqliteDB").Call(
								jen.ID(constants.LoggerVarName),
								jen.ID("cfg").Dot("Database").Dot("ConnectionDetails"),
							)),
						)
					}
					return jen.Null()
				}(),
				jen.Default().Body(
					jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("invalid database type selected: %q"),
						jen.ID("cfg").Dot("Database").Dot("Provider"),
					)),
				),
			),
		),
	}

	return lines
}

func buildProvideDatabaseClient(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideDatabaseClient provides a database implementation dependent on the configuration."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").PointerTo().ID("ServerConfig")).ID("ProvideDatabaseClient").Params(
			constants.CtxParam(),
			proj.LoggerParam(),
			jen.ID("rawDB").PointerTo().Qual("database/sql", "DB"),
		).Params(jen.Qual(proj.DatabasePackage(), "DataManager"), jen.Error()).Body(
			jen.If(jen.ID("rawDB").IsEqualTo().Nil()).Body(
				jen.Return(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("nil DB connection provided"))),
			),
			jen.Line(),
			jen.ID("debug").Assign().ID("cfg").Dot("Database").Dot("Debug").Or().ID("cfg").Dot("Meta").Dot("Debug"),
			jen.Line(),
			jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RegisterAllViews").Call(),
			jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RecordStats").Call(
				jen.ID("rawDB"),
				jen.ID("cfg").Dot("Metrics").Dot("DBMetricsCollectionInterval"),
			),
			jen.Line(),
			jen.Var().ID("dbc").Qual(proj.DatabasePackage(), "DataManager"),
			jen.Switch(jen.ID("cfg").Dot("Database").Dot("Provider")).Body(
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.Postgres) {
						return jen.Case(jen.ID("PostgresProviderKey")).Body(
							jen.ID("dbc").Equals().Qual(proj.DatabasePackage("queriers", "postgres"), "ProvidePostgres").Call(
								jen.ID("debug"),
								jen.ID("rawDB"),
								jen.ID(constants.LoggerVarName),
							),
						)
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.MariaDB) {
						return jen.Case(jen.ID("MariaDBProviderKey")).Body(
							jen.ID("dbc").Equals().Qual(proj.DatabasePackage("queriers", "mariadb"), "ProvideMariaDB").Call(
								jen.ID("debug"),
								jen.ID("rawDB"),
								jen.ID(constants.LoggerVarName),
							),
						)
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.Sqlite) {
						return jen.Case(jen.ID("SqliteProviderKey")).Body(
							jen.ID("dbc").Equals().Qual(proj.DatabasePackage("queriers", "sqlite"), "ProvideSqlite").Call(
								jen.ID("debug"),
								jen.ID("rawDB"),
								jen.ID(constants.LoggerVarName),
							),
						)
					}
					return jen.Null()
				}(),
				jen.Default().Body(
					jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("invalid database type selected: %q"),
						jen.ID("cfg").Dot("Database").Dot("Provider"),
					)),
				),
			),
			jen.Line(),
			jen.Return(jen.Qual(proj.DatabasePackage("client"), "ProvideDatabaseClient").Call(
				constants.CtxVar(),
				jen.ID("rawDB"),
				jen.ID("dbc"),
				jen.ID("debug"),
				jen.ID(constants.LoggerVarName),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideSessionManager(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideSessionManager provides a session manager based on some settings."),
		jen.Line(),
		jen.Comment("There's not a great place to put this function. I don't think it belongs in Auth because it accepts a DB connection,"),
		jen.Line(),
		jen.Comment("but it obviously doesn't belong in the database package, or maybe it does"),
		jen.Line(),
		jen.Func().ID("ProvideSessionManager").Params(
			jen.ID("authConf").ID("AuthSettings"),
			jen.ID("dbConf").ID("DatabaseSettings"),
			jen.ID("db").PointerTo().Qual("database/sql", "DB"),
		).Params(
			jen.PointerTo().Qual("github.com/alexedwards/scs/v2", "SessionManager"),
		).Body(
			jen.ID("sessionManager").Assign().Qual("github.com/alexedwards/scs/v2", "New").Call(),
			jen.Line(),
			jen.Switch(jen.ID("dbConf").Dot("Provider")).Body(
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.Postgres) {
						return jen.Case(jen.ID("PostgresProviderKey")).Body(
							jen.ID("sessionManager").Dot("Store").Equals().Qual(filepath.Join(SessionManagerLibrary, "postgresstore"), "New").Call(jen.ID("db")),
						)
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.MariaDB) {
						return jen.Case(jen.ID("MariaDBProviderKey")).Body(
							jen.ID("sessionManager").Dot("Store").Equals().Qual(filepath.Join(SessionManagerLibrary, "mysqlstore"), "New").Call(jen.ID("db")),
						)
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.Sqlite) {
						return jen.Case(jen.ID("SqliteProviderKey")).Body(
							jen.ID("sessionManager").Dot("Store").Equals().Qual(filepath.Join(SessionManagerLibrary, "sqlite3store"), "New").Call(jen.ID("db")),
						)
					}
					return jen.Null()
				}(),
			),
			jen.Line(),
			jen.ID("sessionManager").Dot("Lifetime").Equals().ID("authConf").Dot("CookieLifetime"),
			jen.Comment("elaborate further here later if you so choose"),
			jen.Line(),
			jen.Return(jen.ID("sessionManager")),
		),
	}

	return lines
}
