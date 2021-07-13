package config

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	serviceConfigs := []jen.Code{
		jen.ID("AuditLog").MapAssign().Qual(proj.AuditServicePackage(), "Config").Valuesln(
			jen.ID("Enabled").MapAssign().ID("true"),
		),
		jen.ID("Auth").MapAssign().Qual(proj.AuthServicePackage(), "Config").Valuesln(
			jen.ID("Cookies").MapAssign().Qual(proj.AuthServicePackage(), "CookieConfig").Valuesln(
				jen.ID("Name").MapAssign().Litf("%s_cookie", proj.Name.RouteName()),
				jen.ID("Domain").MapAssign().Lit("https://verygoodsoftwarenotvirus.ru"),
				jen.ID("Lifetime").MapAssign().Qual("time", "Second"),
			),
			jen.ID("MinimumUsernameLength").MapAssign().Lit(4),
			jen.ID("MinimumPasswordLength").MapAssign().Lit(8),
			jen.ID("EnableUserSignup").MapAssign().ID("true"),
		),
	}

	for _, typ := range proj.DataTypes {
		if typ.SearchEnabled {
			serviceConfigs = append(serviceConfigs,
				jen.ID(typ.Name.Plural()).MapAssign().Qual(proj.ServicePackage(typ.Name.PackageName()), "Config").Valuesln(
					jen.ID("SearchIndexPath").MapAssign().Litf("/%s_index_path", typ.Name.PluralRouteName()),
				),
			)
		}
	}

	code.Add(
		jen.Func().ID("TestServerConfig_EncodeToFile").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("cfg").Assign().AddressOf().ID("InstanceConfig").Valuesln(
						jen.ID("Server").MapAssign().Qual(proj.HTTPServerPackage(), "Config").Valuesln(
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
						jen.ID("Observability").MapAssign().Qual(proj.ObservabilityPackage(), "Config").Valuesln(
							jen.ID("Metrics").MapAssign().Qual(proj.MetricsPackage(), "Config").Valuesln(
								jen.ID("Provider").MapAssign().Lit(""),
								jen.ID("RouteToken").MapAssign().Lit(""),
								jen.ID("RuntimeMetricsCollectionInterval").MapAssign().Lit(2).PointerTo().Qual("time", "Second"),
							),
						),
						jen.ID("Services").MapAssign().ID("ServicesConfigurations").Valuesln(
							serviceConfigs...,
						),
						jen.ID("Database").MapAssign().Qual(proj.DatabasePackage("config"), "Config").Valuesln(
							jen.ID("Provider").MapAssign().Lit("postgres"),
							jen.ID("MetricsCollectionInterval").MapAssign().Lit(2).PointerTo().Qual("time", "Second"),
							jen.ID("Debug").MapAssign().ID("true"),
							jen.ID("RunMigrations").MapAssign().ID("true"),
							jen.ID("ConnectionDetails").MapAssign().Qual(proj.DatabasePackage(), "ConnectionDetails").Call(jen.Lit("postgres://username:passwords@host/table")),
						),
					),
					jen.Newline(),
					jen.List(jen.ID("f"),
						jen.ID("err")).Assign().Qual("io/ioutil", "TempFile").Call(
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("EncodeToFile").Call(
							jen.ID("f").Dot("Name").Call(),
							jen.Qual("encoding/json", "Marshal"),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error marshaling"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("cfg").Assign().AddressOf().ID("InstanceConfig").Values(),
					jen.Newline(),
					jen.List(jen.ID("f"),
						jen.ID("err")).Assign().Qual("io/ioutil", "TempFile").Call(
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
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
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestServerConfig_ProvideDatabaseClient").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("supported providers"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					constants.LoggerVar().Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.Newline(),
					jen.For(jen.List(jen.Underscore(), jen.ID("provider")).Assign().Range().Index().String().Values(
						utils.ConditionalCode(proj.DatabaseIsEnabled(models.Sqlite), jen.Lit("sqlite")),
						utils.ConditionalCode(proj.DatabaseIsEnabled(models.Postgres), jen.Lit("postgres")),
						utils.ConditionalCode(proj.DatabaseIsEnabled(models.MariaDB), jen.Lit("mariadb")),
					)).Body(
						jen.ID("cfg").Assign().AddressOf().ID("InstanceConfig").Valuesln(
							jen.ID("Database").MapAssign().Qual(proj.DatabasePackage("config"), "Config").Valuesln(
								jen.ID("Provider").MapAssign().ID("provider"),
							),
						),
						jen.Newline(),
						jen.List(jen.ID("x"),
							jen.ID("err")).Assign().ID("ProvideDatabaseClient").Call(
							jen.ID("ctx"),
							constants.LoggerVar(),
							jen.AddressOf().Qual("database/sql", "DB").Values(),
							jen.ID("cfg"),
						),
						jen.Qual(constants.AssertionLibrary, "NotNil").Call(
							jen.ID("t"),
							jen.ID("x"),
						),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil *sql.DB"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					constants.LoggerVar().Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("cfg").Assign().AddressOf().ID("InstanceConfig").Values(),
					jen.Newline(),
					jen.List(jen.ID("x"),
						jen.ID("err")).Assign().ID("ProvideDatabaseClient").Call(
						jen.ID("ctx"),
						constants.LoggerVar(),
						jen.ID("nil"),
						jen.ID("cfg"),
					),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(
						jen.ID("t"),
						jen.ID("x"),
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid provider"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					constants.LoggerVar().Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.Newline(),
					jen.ID("cfg").Assign().AddressOf().ID("InstanceConfig").Valuesln(
						jen.ID("Database").MapAssign().Qual(proj.DatabasePackage("config"), "Config").Valuesln(
							jen.ID("Provider").MapAssign().Lit("provider"),
						),
					),
					jen.Newline(),
					jen.List(jen.ID("x"),
						jen.ID("err")).Assign().ID("ProvideDatabaseClient").Call(
						jen.ID("ctx"),
						constants.LoggerVar(),
						jen.AddressOf().Qual("database/sql", "DB").Values(),
						jen.ID("cfg"),
					),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(
						jen.ID("t"),
						jen.ID("x"),
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Newline(),
	)

	return code
}
