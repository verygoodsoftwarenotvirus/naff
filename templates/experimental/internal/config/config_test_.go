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
		jen.Func().ID("TestServerConfig_EncodeToFile").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("ServerConfig").Valuesln(jen.ID("Server").Op(":").ID("server").Dot("Config").Valuesln(jen.ID("HTTPPort").Op(":").Lit(1234), jen.ID("Debug").Op(":").ID("false"), jen.ID("StartupDeadline").Op(":").Qual("time", "Minute")), jen.ID("AuditLog").Op(":").ID("audit").Dot("Config").Valuesln(jen.ID("Enabled").Op(":").ID("true")), jen.ID("Meta").Op(":").ID("MetaSettings").Valuesln(jen.ID("RunMode").Op(":").ID("DevelopmentRunMode")), jen.ID("Encoding").Op(":").ID("encoding").Dot("Config").Valuesln(jen.ID("ContentType").Op(":").Lit("application/json")), jen.ID("Auth").Op(":").ID("authentication").Dot("Config").Valuesln(jen.ID("Cookies").Op(":").ID("authentication").Dot("CookieConfig").Valuesln(jen.ID("Name").Op(":").Lit("todocookie"), jen.ID("Domain").Op(":").Lit("https://verygoodsoftwarenotvirus.ru"), jen.ID("Lifetime").Op(":").Qual("time", "Second")), jen.ID("MinimumUsernameLength").Op(":").Lit(4), jen.ID("MinimumPasswordLength").Op(":").Lit(8), jen.ID("EnableUserSignup").Op(":").ID("true")), jen.ID("Observability").Op(":").ID("observability").Dot("Config").Valuesln(jen.ID("Metrics").Op(":").ID("metrics").Dot("Config").Valuesln(jen.ID("Provider").Op(":").Lit(""), jen.ID("RouteToken").Op(":").Lit(""), jen.ID("RuntimeMetricsCollectionInterval").Op(":").Lit(2).Op("*").Qual("time", "Second"))), jen.ID("Frontend").Op(":").ID("frontend").Dot("Config").Valuesln(), jen.ID("Search").Op(":").ID("search").Dot("Config").Valuesln(jen.ID("ItemsIndexPath").Op(":").Lit("/items_index_path")), jen.ID("Database").Op(":").ID("config").Dot("Config").Valuesln(jen.ID("Provider").Op(":").Lit("postgres"), jen.ID("MetricsCollectionInterval").Op(":").Lit(2).Op("*").Qual("time", "Second"), jen.ID("Debug").Op(":").ID("true"), jen.ID("RunMigrations").Op(":").ID("true"), jen.ID("ConnectionDetails").Op(":").ID("database").Dot("ConnectionDetails").Call(jen.Lit("postgres://username:passwords@host/table")))),
					jen.List(jen.ID("f"), jen.ID("err")).Op(":=").Qual("io/ioutil", "TempFile").Call(
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("EncodeToFile").Call(
							jen.ID("f").Dot("Name").Call(),
							jen.Qual("encoding/json", "Marshal"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error marshaling"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("ServerConfig").Valuesln(),
					jen.List(jen.ID("f"), jen.ID("err")).Op(":=").Qual("io/ioutil", "TempFile").Call(
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("EncodeToFile").Call(
							jen.ID("f").Dot("Name").Call(),
							jen.Func().Params(jen.Interface()).Params(jen.Index().ID("byte"), jen.ID("error")).Body(
								jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("blah")))),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestServerConfig_ProvideDatabaseClient").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("supported providers"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.For(jen.List(jen.ID("_"), jen.ID("provider")).Op(":=").Range().Index().ID("string").Valuesln(jen.Lit("sqlite"), jen.Lit("postgres"), jen.Lit("mariadb"))).Body(
						jen.ID("cfg").Op(":=").Op("&").ID("ServerConfig").Valuesln(jen.ID("Database").Op(":").ID("config").Dot("Config").Valuesln(jen.ID("Provider").Op(":").ID("provider"))),
						jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideDatabaseClient").Call(
							jen.ID("ctx"),
							jen.ID("logger"),
							jen.Op("&").Qual("database/sql", "DB").Valuesln(),
						),
						jen.ID("assert").Dot("NotNil").Call(
							jen.ID("t"),
							jen.ID("x"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil *sql.DB"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("ServerConfig").Valuesln(),
					jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideDatabaseClient").Call(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("x"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("ServerConfig").Valuesln(jen.ID("Database").Op(":").ID("config").Dot("Config").Valuesln(jen.ID("Provider").Op(":").Lit("provider"))),
					jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideDatabaseClient").Call(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.Op("&").Qual("database/sql", "DB").Valuesln(),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("x"),
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
