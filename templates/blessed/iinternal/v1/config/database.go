package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func databaseDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("config")

	utils.AddImports(proj, code)

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

	code.Add(jen.Const().Defs(constDefs...), jen.Line())
	code.Add(buildProvideDatabaseConnection(proj)...)
	code.Add(buildProvideDatabaseClient(proj)...)
	code.Add(buildProvideSessionManager(proj)...)

	return code
}

func buildProvideDatabaseConnection(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideDatabaseConnection provides a database implementation dependent on the configuration."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").PointerTo().ID("ServerConfig")).ID("ProvideDatabaseConnection").Params(
			jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger"),
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
			jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger"),
			jen.ID("rawDB").PointerTo().Qual("database/sql", "DB"),
		).Params(jen.Qual(proj.DatabaseV1Package(), "DataManager"), jen.Error()).Block(
			jen.Var().Defs(
				jen.ID("debug").Equals().ID("cfg").Dot("Database").Dot("Debug").Or().ID("cfg").Dot("Meta").Dot("Debug"),
				jen.ID("connectionDetails").Equals().ID("cfg").Dot("Database").Dot("ConnectionDetails"),
			),
			jen.Line(),
			jen.Switch(jen.ID("cfg").Dot("Database").Dot("Provider")).Block(
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.Postgres) {
						return jen.Case(jen.ID("postgresProviderKey")).Block(
							jen.List(jen.ID("rawDB"), jen.Err()).Assign().Qual(proj.DatabaseV1Package("queriers", "postgres"), "ProvidePostgresDB").Call(
								jen.ID(constants.LoggerVarName),
								jen.ID("connectionDetails"),
							),
							jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
								jen.Return().List(
									jen.Nil(),
									jen.Qual("fmt", "Errorf").Call(
										jen.Lit("establish postgres database connection: %w"),
										jen.Err(),
									),
								),
							),
							jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RegisterAllViews").Call(),
							jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RecordStats").Call(
								jen.ID("rawDB"),
								jen.ID("cfg").Dot("Metrics").Dot("DBMetricsCollectionInterval"),
							),
							jen.Line(),
							jen.ID("pgdb").Assign().Qual(proj.DatabaseV1Package("queriers", "postgres"), "ProvidePostgres").Call(
								jen.ID("debug"),
								jen.ID("rawDB"),
								jen.ID(constants.LoggerVarName),
							),
							jen.Line(),
							jen.Return().Qual(proj.DatabaseV1Package("client"), "ProvideDatabaseClient").Call(
								constants.CtxVar(),
								jen.ID("rawDB"),
								jen.ID("pgdb"),
								jen.ID("debug"),
								jen.ID(constants.LoggerVarName),
							),
						)
					} else {
						return jen.Null()
					}

				}(),
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.MariaDB) {
						return jen.Case(jen.ID("mariaDBProviderKey")).Block(
							jen.List(jen.ID("rawDB"), jen.Err()).Assign().Qual(proj.DatabaseV1Package("queriers", "mariadb"), "ProvideMariaDBConnection").Call(
								jen.ID(constants.LoggerVarName),
								jen.ID("connectionDetails"),
							),
							jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
								jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(
									jen.Lit("establish mariadb database connection: %w"),
									jen.Err(),
								),
								),
							),
							jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RegisterAllViews").Call(),
							jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RecordStats").Call(
								jen.ID("rawDB"),
								jen.ID("cfg").Dot("Metrics").Dot("DBMetricsCollectionInterval"),
							),
							jen.Line(),
							jen.ID("mdb").Assign().Qual(proj.DatabaseV1Package("queriers", "mariadb"), "ProvideMariaDB").Call(
								jen.ID("debug"),
								jen.ID("rawDB"),
								jen.ID(constants.LoggerVarName),
							),
							jen.Line(),
							jen.Return().Qual(proj.DatabaseV1Package("client"), "ProvideDatabaseClient").Call(
								constants.CtxVar(),
								jen.ID("rawDB"),
								jen.ID("mdb"),
								jen.ID("debug"),
								jen.ID(constants.LoggerVarName),
							),
						)
					} else {
						return jen.Null()
					}

				}(),
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.Sqlite) {
						return jen.Case(jen.ID("sqliteProviderKey")).Block(
							jen.List(jen.ID("rawDB"), jen.Err()).Assign().Qual(proj.DatabaseV1Package("queriers", "sqlite"), "ProvideSqliteDB").Call(
								jen.ID(constants.LoggerVarName),
								jen.ID("connectionDetails"),
							),
							jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
								jen.Return().List(
									jen.Nil(),
									jen.Qual("fmt", "Errorf").Call(
										jen.Lit("establish sqlite database connection: %w"),
										jen.Err(),
									),
								),
							),
							jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RegisterAllViews").Call(),
							jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RecordStats").Call(
								jen.ID("rawDB"),
								jen.ID("cfg").Dot("Metrics").Dot("DBMetricsCollectionInterval"),
							),
							jen.Line(),
							jen.ID("sdb").Assign().Qual(proj.DatabaseV1Package("queriers", "sqlite"), "ProvideSqlite").Call(
								jen.ID("debug"), jen.ID("rawDB"), jen.ID(constants.LoggerVarName)),
							jen.Line(),
							jen.Return().Qual(proj.DatabaseV1Package("client"), "ProvideDatabaseClient").Call(
								constants.CtxVar(),
								jen.ID("rawDB"),
								jen.ID("sdb"),
								jen.ID("debug"),
								jen.ID(constants.LoggerVarName),
							),
						)
					} else {
						return jen.Null()
					}

				}(),
				jen.Default().Block(
					jen.Return().List(jen.Nil(), utils.Error("invalid database type selected")),
				),
			),
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
			jen.PointerTo().Qual("database/sql", "DB"),
		).Params(
			jen.PointerTo().Qual(utils.SessionManagerLibrary, "SessionManager"),
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
