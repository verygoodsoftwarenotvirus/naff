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
	code := jen.NewFile("config")

	utils.AddImports(proj, code)
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
			constants.LoggerParam(),
		).Params(
			jen.PointerTo().Qual("database/sql", "DB"),
			jen.Error(),
		).Block(
			jen.Switch(jen.ID("cfg").Dot("Database").Dot("Provider")).Block(
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.Postgres) {
						return jen.Case(jen.ID("PostgresProviderKey")).Block(
							jen.Return(jen.Qual(proj.DatabaseV1Package("queriers", "postgres"), "ProvidePostgresDB").Call(
								jen.ID(constants.LoggerVarName),
								jen.ID("cfg").Dot("Database").Dot("ConnectionDetails"),
							)),
						)
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.MariaDB) {
						return jen.Case(jen.ID("MariaDBProviderKey")).Block(
							jen.Return(jen.Qual(proj.DatabaseV1Package("queriers", "mariadb"), "ProvideMariaDBConnection").Call(
								jen.ID(constants.LoggerVarName),
								jen.ID("cfg").Dot("Database").Dot("ConnectionDetails"),
							)),
						)
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.Sqlite) {
						return jen.Case(jen.ID("SqliteProviderKey")).Block(
							jen.Return(jen.Qual(proj.DatabaseV1Package("queriers", "sqlite"), "ProvideSqliteDB").Call(
								jen.ID(constants.LoggerVarName),
								jen.ID("cfg").Dot("Database").Dot("ConnectionDetails"),
							)),
						)
					}
					return jen.Null()
				}(),
				jen.Default().Block(
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
			constants.LoggerParam(),
			jen.ID("rawDB").PointerTo().Qual("database/sql", "DB"),
		).Params(jen.Qual(proj.DatabaseV1Package(), "DataManager"), jen.Error()).Block(
			jen.If(jen.ID("rawDB").IsEqualTo().Nil()).Block(
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
			jen.Var().ID("dbc").Qual(proj.DatabaseV1Package(), "DataManager"),
			jen.Switch(jen.ID("cfg").Dot("Database").Dot("Provider")).Block(
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.Postgres) {
						return jen.Case(jen.ID("PostgresProviderKey")).Block(
							jen.ID("dbc").Equals().Qual(proj.DatabaseV1Package("queriers", "postgres"), "ProvidePostgres").Call(
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
						return jen.Case(jen.ID("MariaDBProviderKey")).Block(
							jen.ID("dbc").Equals().Qual(proj.DatabaseV1Package("queriers", "mariadb"), "ProvideMariaDB").Call(
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
						return jen.Case(jen.ID("SqliteProviderKey")).Block(
							jen.ID("dbc").Equals().Qual(proj.DatabaseV1Package("queriers", "sqlite"), "ProvideSqlite").Call(
								jen.ID("debug"),
								jen.ID("rawDB"),
								jen.ID(constants.LoggerVarName),
							),
						)
					}
					return jen.Null()
				}(),
				jen.Default().Block(
					jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("invalid database type selected: %q"),
						jen.ID("cfg").Dot("Database").Dot("Provider"),
					)),
				),
			),
			jen.Line(),
			jen.Return(jen.Qual(proj.DatabaseV1Package("client"), "ProvideDatabaseClient").Call(
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
		).Block(
			jen.ID("sessionManager").Assign().Qual("github.com/alexedwards/scs/v2", "New").Call(),
			jen.Line(),
			jen.Switch(jen.ID("dbConf").Dot("Provider")).Block(
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.Postgres) {
						return jen.Case(jen.ID("PostgresProviderKey")).Block(
							jen.ID("sessionManager").Dot("Store").Equals().Qual(filepath.Join(SessionManagerLibrary, "postgresstore"), "New").Call(jen.ID("db")),
						)
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.MariaDB) {
						return jen.Case(jen.ID("MariaDBProviderKey")).Block(
							jen.ID("sessionManager").Dot("Store").Equals().Qual(filepath.Join(SessionManagerLibrary, "mysqlstore"), "New").Call(jen.ID("db")),
						)
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.Sqlite) {
						return jen.Case(jen.ID("SqliteProviderKey")).Block(
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
