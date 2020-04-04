package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func databaseDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("config")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("postgresProviderKey").Equals().Lit("postgres"),
			jen.ID("mariaDBProviderKey").Equals().Lit("mariadb"),
			jen.ID("sqliteProviderKey").Equals().Lit("sqlite"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideDatabase provides a database implementation dependent on the configuration"),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").PointerTo().ID("ServerConfig")).ID("ProvideDatabase").Params(utils.CtxParam(), jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual(proj.DatabaseV1Package(), "Database"), jen.Error()).Block(
			jen.Var().Defs(
				jen.ID("debug").Equals().ID("cfg").Dot("Database").Dot("Debug").Or().ID("cfg").Dot("Meta").Dot("Debug"),
				jen.ID("connectionDetails").Equals().ID("cfg").Dot("Database").Dot("ConnectionDetails"),
			),
			jen.Line(),
			jen.Switch(jen.ID("cfg").Dot("Database").Dot("Provider")).Block(
				jen.Case(jen.ID("postgresProviderKey")).Block(
					jen.List(jen.ID("rawDB"), jen.Err()).Assign().ID("postgres").Dot("ProvidePostgresDB").Call(jen.ID("logger"), jen.ID("connectionDetails")),
					jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("establish postgres database connection: %w"), jen.Err())),
					),
					jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RegisterAllViews").Call(),
					jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RecordStats").Call(jen.ID("rawDB"), jen.ID("cfg").Dot("Metrics").Dot("DBMetricsCollectionInterval")),
					jen.Line(),
					jen.ID("pgdb").Assign().Qual(proj.DatabaseV1Package("queriers", "postgres"), "ProvidePostgres").Call(jen.ID("debug"), jen.ID("rawDB"), jen.ID("logger")),
					jen.Line(),
					jen.Return().Qual(proj.DatabaseV1Package("client"), "ProvideDatabaseClient").Call(utils.CtxVar(), jen.ID("rawDB"), jen.ID("pgdb"), jen.ID("debug"), jen.ID("logger"))),
				jen.Case(jen.ID("mariaDBProviderKey")).Block(
					jen.List(jen.ID("rawDB"), jen.Err()).Assign().ID("mariadb").Dot("ProvideMariaDBConnection").Call(jen.ID("logger"), jen.ID("connectionDetails")),
					jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("establish mariadb database connection: %w"), jen.Err())),
					),
					jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RegisterAllViews").Call(),
					jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RecordStats").Call(jen.ID("rawDB"), jen.ID("cfg").Dot("Metrics").Dot("DBMetricsCollectionInterval")),
					jen.Line(),
					jen.ID("mdb").Assign().Qual(proj.DatabaseV1Package("queriers", "mariadb"), "ProvideMariaDB").Call(jen.ID("debug"), jen.ID("rawDB"), jen.ID("logger")),
					jen.Line(),
					jen.Return().Qual(proj.DatabaseV1Package("client"), "ProvideDatabaseClient").Call(utils.CtxVar(), jen.ID("rawDB"), jen.ID("mdb"), jen.ID("debug"), jen.ID("logger"))),
				jen.Case(jen.ID("sqliteProviderKey")).Block(
					jen.List(jen.ID("rawDB"), jen.Err()).Assign().ID("sqlite").Dot("ProvideSqliteDB").Call(jen.ID("logger"), jen.ID("connectionDetails")),
					jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("establish sqlite database connection: %w"), jen.Err())),
					),
					jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RegisterAllViews").Call(),
					jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RecordStats").Call(jen.ID("rawDB"), jen.ID("cfg").Dot("Metrics").Dot("DBMetricsCollectionInterval")),
					jen.Line(),
					jen.ID("sdb").Assign().Qual(proj.DatabaseV1Package("queriers", "sqlite"), "ProvideSqlite").Call(jen.ID("debug"), jen.ID("rawDB"), jen.ID("logger")),
					jen.Line(),
					jen.Return().Qual(proj.DatabaseV1Package("client"), "ProvideDatabaseClient").Call(utils.CtxVar(), jen.ID("rawDB"), jen.ID("sdb"), jen.ID("debug"), jen.ID("logger"))),
				jen.Default().Block(jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("invalid database type selected")))),
			),
		),
		jen.Line(),
	)
	return ret
}
