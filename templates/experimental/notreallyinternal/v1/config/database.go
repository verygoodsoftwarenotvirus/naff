package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func databaseDotGo() *jen.File {
	ret := jen.NewFile("config")

	utils.AddImports(ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("postgresProviderKey").Op("=").Lit("postgres"),
			jen.ID("mariaDBProviderKey").Op("=").Lit("mariadb"),
			jen.ID("sqliteProviderKey").Op("=").Lit("sqlite"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideDatabase provides a database implementation dependent on the configuration"),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("ServerConfig")).ID("ProvideDatabase").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1", "Database"), jen.ID("error")).Block(
			jen.Var().Defs(
				jen.ID("debug").Op("=").ID("cfg").Dot("Database").Dot("Debug").Op("||").ID("cfg").Dot("Meta").Dot("Debug"),
				jen.ID("connectionDetails").Op("=").ID("cfg").Dot("Database").Dot("ConnectionDetails"),
			),
			jen.Line(),
			jen.Switch(jen.ID("cfg").Dot("Database").Dot("Provider")).Block(
				jen.Case(jen.ID("postgresProviderKey")).Block(
					jen.List(jen.ID("rawDB"), jen.ID("err")).Op(":=").ID("postgres").Dot("ProvidePostgresDB").Call(jen.ID("logger"), jen.ID("connectionDetails")),
					jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("establish postgres database connection: %w"), jen.ID("err"))),
					),
					jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RegisterAllViews").Call(),
					jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RecordStats").Call(jen.ID("rawDB"), jen.ID("cfg").Dot("Metrics").Dot("DBMetricsCollectionInterval")),
					jen.Line(),
					jen.ID("pgdb").Op(":=").ID("postgres").Dot("ProvidePostgres").Call(jen.ID("debug"), jen.ID("rawDB"), jen.ID("logger")),
					jen.Line(),
					jen.Return().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/client", "ProvideDatabaseClient").Call(jen.ID("ctx"), jen.ID("rawDB"), jen.ID("pgdb"), jen.ID("debug"), jen.ID("logger"))),
				jen.Case(jen.ID("mariaDBProviderKey")).Block(
					jen.List(jen.ID("rawDB"), jen.ID("err")).Op(":=").ID("mariadb").Dot("ProvideMariaDBConnection").Call(jen.ID("logger"), jen.ID("connectionDetails")),
					jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("establish mariadb database connection: %w"), jen.ID("err"))),
					),
					jen.Line(),
					jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RegisterAllViews").Call(),
					jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RecordStats").Call(jen.ID("rawDB"), jen.ID("cfg").Dot("Metrics").Dot("DBMetricsCollectionInterval")),
					jen.Line(),
					jen.ID("mdb").Op(":=").ID("mariadb").Dot("ProvideMariaDB").Call(jen.ID("debug"), jen.ID("rawDB"), jen.ID("logger")),
					jen.Line(),
					jen.Return().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/client", "ProvideDatabaseClient").Call(jen.ID("ctx"), jen.ID("rawDB"), jen.ID("mdb"), jen.ID("debug"), jen.ID("logger"))),
				jen.Case(jen.ID("sqliteProviderKey")).Block(
					jen.List(jen.ID("rawDB"), jen.ID("err")).Op(":=").ID("sqlite").Dot("ProvideSqliteDB").Call(jen.ID("logger"), jen.ID("connectionDetails")),
					jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("establish sqlite database connection: %w"), jen.ID("err"))),
					),
					jen.Line(),
					jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RegisterAllViews").Call(),
					jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "RecordStats").Call(jen.ID("rawDB"), jen.ID("cfg").Dot("Metrics").Dot("DBMetricsCollectionInterval")),
					jen.Line(),
					jen.ID("sdb").Op(":=").ID("sqlite").Dot("ProvideSqlite").Call(jen.ID("debug"), jen.ID("rawDB"), jen.ID("logger")),
					jen.Line(),
					jen.Return().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/client", "ProvideDatabaseClient").Call(jen.ID("ctx"), jen.ID("rawDB"), jen.ID("sdb"), jen.ID("debug"), jen.ID("logger"))),
				jen.Default().Block(
					jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("invalid database type selected")))),
			),
		),
		jen.Line(),
	)
	return ret
}
