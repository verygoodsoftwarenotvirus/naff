package config

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	serviceConfigs := []jen.Code{
		jen.ID("AuditLog").MapAssign().ID("audit").Dot("Config").Valuesln(
			jen.ID("Enabled").MapAssign().ID("true"),
		),
		jen.ID("Auth").MapAssign().Qual(proj.AuthServicePackage(), "Config").Valuesln(
			jen.ID("Cookies").MapAssign().Qual(proj.AuthServicePackage(), "CookieConfig").Valuesln(
				jen.ID("Name").MapAssign().Lit("todocookie"),
				jen.ID("Domain").MapAssign().Lit("https://verygoodsoftwarenotvirus.ru"),
				jen.ID("Lifetime").MapAssign().Qual("time", "Second"),
			),
			jen.ID("MinimumUsernameLength").MapAssign().Lit(4),
			jen.ID("MinimumPasswordLength").MapAssign().Lit(8),
			jen.ID("EnableUserSignup").MapAssign().ID("true"),
		),
	}

	for _, typ := range proj.DataTypes {
		serviceConfigs = append(serviceConfigs,
			jen.ID(typ.Name.Plural()).MapAssign().Qual(proj.ServicePackage(typ.Name.PackageName()), "Config").Valuesln(
				jen.ID("SearchIndexPath").MapAssign().Litf("/%s_index_path", typ.Name.PluralRouteName()),
			),
		)
	}

	code.Add(
		jen.Func().ID("TestServerConfig_EncodeToFile").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("cfg").Op(":=").Op("&").ID("InstanceConfig").Valuesln(
						jen.ID("Server").MapAssign().ID("server").Dot("Config").Valuesln(
							jen.ID("HTTPPort").MapAssign().Lit(1234),
							jen.ID("Debug").MapAssign().ID("false"),
							jen.ID("StartupDeadline").MapAssign().Qual("time", "Minute"),
						),
						jen.ID("Meta").MapAssign().ID("MetaSettings").Valuesln(
							jen.ID("RunMode").MapAssign().ID("DevelopmentRunMode"),
						),
						jen.ID("Encoding").MapAssign().Qual(proj.EncodingPackage(), "Config").Valuesln(
							jen.ID("ContentType").MapAssign().Lit("application/json"),
						),
						jen.ID("Observability").MapAssign().ID("observability").Dot("Config").Valuesln(
							jen.ID("Metrics").MapAssign().Qual(proj.MetricsPackage(), "Config").Valuesln(
								jen.ID("Provider").MapAssign().Lit(""),
								jen.ID("RouteToken").MapAssign().Lit(""),
								jen.ID("RuntimeMetricsCollectionInterval").MapAssign().Lit(2).Op("*").Qual("time", "Second"),
							),
						),
						jen.ID("Services").MapAssign().ID("ServicesConfigurations").Valuesln(
							serviceConfigs...,
						),
						jen.ID("Database").MapAssign().ID("config").Dot("Config").Valuesln(
							jen.ID("Provider").MapAssign().Lit("postgres"),
							jen.ID("MetricsCollectionInterval").MapAssign().Lit(2).Op("*").Qual("time", "Second"),
							jen.ID("Debug").MapAssign().ID("true"),
							jen.ID("RunMigrations").MapAssign().ID("true"),
							jen.ID("ConnectionDetails").MapAssign().ID("database").Dot("ConnectionDetails").Call(jen.Lit("postgres://username:passwords@host/table")),
						),
					),
					jen.Line(),
					jen.List(jen.ID("f"),
						jen.ID("err")).Op(":=").Qual("io/ioutil", "TempFile").Call(
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Line(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("EncodeToFile").Call(
							jen.ID("f").Dot("Name").Call(),
							jen.Qual("encoding/json", "Marshal"),
						),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error marshaling"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("cfg").Op(":=").Op("&").ID("InstanceConfig").Values(),
					jen.Line(),
					jen.List(jen.ID("f"),
						jen.ID("err")).Op(":=").Qual("io/ioutil", "TempFile").Call(
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Line(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("EncodeToFile").Call(
							jen.ID("f").Dot("Name").Call(),
							jen.Func().Params(jen.Interface()).Params(jen.Index().ID("byte"),
								jen.ID("error")).Body(
								jen.Return().List(jen.ID("nil"),
									jen.Qual("errors", "New").Call(jen.Lit("blah")),
								),
							),
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
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("supported providers"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.Line(),
					jen.For(jen.List(jen.ID("_"), jen.ID("provider")).Op(":=").Range().Index().ID("string").Values(
						jen.Lit("sqlite"),
						jen.Lit("postgres"),
						jen.Lit("mariadb"),
					)).Body(
						jen.ID("cfg").Op(":=").Op("&").ID("InstanceConfig").Valuesln(
							jen.ID("Database").MapAssign().ID("config").Dot("Config").Valuesln(
								jen.ID("Provider").MapAssign().ID("provider"),
							),
						),
						jen.Line(),
						jen.List(jen.ID("x"),
							jen.ID("err")).Op(":=").ID("ProvideDatabaseClient").Call(
							jen.ID("ctx"),
							jen.ID("logger"),
							jen.Op("&").Qual("database/sql", "DB").Values(),
							jen.ID("cfg"),
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
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil *sql.DB"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("InstanceConfig").Values(),
					jen.Line(),
					jen.List(jen.ID("x"),
						jen.ID("err")).Op(":=").ID("ProvideDatabaseClient").Call(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("nil"),
						jen.ID("cfg"),
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
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.Line(),
					jen.ID("cfg").Op(":=").Op("&").ID("InstanceConfig").Valuesln(
						jen.ID("Database").MapAssign().ID("config").Dot("Config").Valuesln(
							jen.ID("Provider").MapAssign().Lit("provider"),
						),
					),
					jen.Line(),
					jen.List(jen.ID("x"),
						jen.ID("err")).Op(":=").ID("ProvideDatabaseClient").Call(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.Op("&").Qual("database/sql", "DB").Values(),
						jen.ID("cfg"),
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
