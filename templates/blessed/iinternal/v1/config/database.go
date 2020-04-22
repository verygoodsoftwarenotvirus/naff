package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func databaseDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("config")

	utils.AddImports(proj, ret)

	var constDefs []jen.Code
	if proj.DatabasesIsEnabled(models.Postgres) {
		constDefs = append(constDefs, jen.ID("postgresProviderKey").Equals().Lit("postgres"))
	}
	if proj.DatabasesIsEnabled(models.Sqlite) {
		constDefs = append(constDefs, jen.ID("mariaDBProviderKey").Equals().Lit("mariadb"))
	}
	if proj.DatabasesIsEnabled(models.MariaDB) {
		constDefs = append(constDefs, jen.ID("sqliteProviderKey").Equals().Lit("sqlite"))
	}

	ret.Add(jen.Const().Defs(constDefs...), jen.Line())

	ret.Add(
		jen.Comment("ProvideDatabase provides a database implementation dependent on the configuration"),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").PointerTo().ID("ServerConfig")).ID("ProvideDatabase").Params(
			constants.CtxParam(),
			jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Qual(proj.DatabaseV1Package(), "Database"), jen.Error()).Block(
			jen.Var().Defs(
				jen.ID("debug").Equals().ID("cfg").Dot("Database").Dot("Debug").Or().ID("cfg").Dot("Meta").Dot("Debug"),
				jen.ID("connectionDetails").Equals().ID("cfg").Dot("Database").Dot("ConnectionDetails"),
			),
			jen.Line(),
			jen.Switch(jen.ID("cfg").Dot("Database").Dot("Provider")).Block(
				func() jen.Code {
					if proj.DatabasesIsEnabled(models.Postgres) {
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
					if proj.DatabasesIsEnabled(models.MariaDB) {
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
					if proj.DatabasesIsEnabled(models.Sqlite) {
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
	)
	return ret
}
