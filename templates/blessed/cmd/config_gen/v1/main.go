package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("main")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("defaultPort").Equals().Lit(8888),
			jen.ID("oneDay").Equals().Lit(24).Times().Qual("time", "Hour"),
			jen.ID("debugCookieSecret").Equals().Lit("HEREISA32CHARSECRETWHICHISMADEUP"),
			jen.ID("defaultFrontendFilepath").Equals().Lit("/frontend"),
			jen.ID("postgresDBConnDetails").Equals().Lit("postgres://dbuser:hunter2@database:5432/todo?sslmode=disable"),
			jen.ID("metaDebug").Equals().Lit("meta.debug"),
			jen.ID("metaStartupDeadline").Equals().Lit("meta.startup_deadline"),
			jen.ID("serverHTTPPort").Equals().Lit("server.http_port"),
			jen.ID("serverDebug").Equals().Lit("server.debug"),
			jen.ID("frontendDebug").Equals().Lit("frontend.debug"),
			jen.ID("frontendStaticFilesDir").Equals().Lit("frontend.static_files_directory"),
			jen.ID("frontendCacheStatics").Equals().Lit("frontend.cache_static_files"),
			jen.ID("authDebug").Equals().Lit("auth.debug"),
			jen.ID("authCookieDomain").Equals().Lit("auth.cookie_domain"),
			jen.ID("authCookieSecret").Equals().Lit("auth.cookie_secret"),
			jen.ID("authCookieLifetime").Equals().Lit("auth.cookie_lifetime"),
			jen.ID("authSecureCookiesOnly").Equals().Lit("auth.secure_cookies_only"),
			jen.ID("authEnableUserSignup").Equals().Lit("auth.enable_user_signup"),
			jen.ID("metricsProvider").Equals().Lit("metrics.metrics_provider"),
			jen.ID("metricsTracer").Equals().Lit("metrics.tracing_provider"),
			jen.ID("metricsDBCollectionInterval").Equals().Lit("metrics.database_metrics_collection_interval"),
			jen.ID("metricsRuntimeCollectionInterval").Equals().Lit("metrics.runtime_metrics_collection_interval"),
			jen.ID("dbDebug").Equals().Lit("database.debug"),
			jen.ID("dbProvider").Equals().Lit("database.provider"),
			jen.ID("dbDeets").Equals().Lit("database.connection_details"),
			func() jen.Code {
				if proj.DatabasesIsEnabled(models.Postgres) {
					return jen.ID("postgres").Equals().Lit("postgres")
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if proj.DatabasesIsEnabled(models.Sqlite) {
					return jen.ID("sqlite").Equals().Lit("sqlite")
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if proj.DatabasesIsEnabled(models.MariaDB) {
					return jen.ID("mariadb").Equals().Lit("mariadb")
				}
				return jen.Null()
			}(),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("configFunc").Func().Params(jen.ID("filepath").String()).Params(jen.Error()),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("files").Equals().Map(jen.String()).ID("configFunc").Valuesln(
				jen.Lit("config_files/coverage.toml").MapAssign().ID("coverageConfig"),
				jen.Lit("config_files/development.toml").MapAssign().ID("developmentConfig"),
				func() jen.Code {
					if proj.DatabasesIsEnabled(models.Postgres) {
						return jen.Lit("config_files/integration-tests-postgres.toml").MapAssign().ID("buildIntegrationTestForDBImplementation").Call(jen.ID("postgres"), jen.ID("postgresDBConnDetails"))
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if proj.DatabasesIsEnabled(models.Sqlite) {
						return jen.Lit("config_files/integration-tests-sqlite.toml").MapAssign().ID("buildIntegrationTestForDBImplementation").Call(jen.ID("sqlite"), jen.Lit("/tmp/db"))
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if proj.DatabasesIsEnabled(models.MariaDB) {
						return jen.Lit("config_files/integration-tests-mariadb.toml").MapAssign().ID("buildIntegrationTestForDBImplementation").Call(jen.ID("mariadb"), jen.Lit("dbuser:hunter2@tcp(database:3306)/todo"))
					}
					return jen.Null()
				}(),
				jen.Lit("config_files/production.toml").MapAssign().ID("productionConfig")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("developmentConfig").Params(jen.ID("filepath").String()).Params(jen.Error()).Block(
			jen.ID("cfg").Assign().ID("config").Dot("BuildConfig").Call(),
			jen.Line(),
			jen.ID("cfg").Dot("Set").Call(jen.ID("metaStartupDeadline"), jen.Qual("time", "Minute")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("serverHTTPPort"), jen.ID("defaultPort")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("serverDebug"), jen.True()),
			jen.Line(),
			jen.ID("cfg").Dot("Set").Call(jen.ID("frontendDebug"), jen.True()),
			jen.ID("cfg").Dot("Set").Call(jen.ID("frontendStaticFilesDir"), jen.ID("defaultFrontendFilepath")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("frontendCacheStatics"), jen.False()),
			jen.Line(),
			jen.ID("cfg").Dot("Set").Call(jen.ID("authDebug"), jen.True()),
			jen.ID("cfg").Dot("Set").Call(jen.ID("authCookieDomain"), jen.EmptyString()),
			jen.ID("cfg").Dot("Set").Call(jen.ID("authCookieSecret"), jen.ID("debugCookieSecret")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("authCookieLifetime"), jen.ID("oneDay")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("authSecureCookiesOnly"), jen.False()),
			jen.ID("cfg").Dot("Set").Call(jen.ID("authEnableUserSignup"), jen.True()),
			jen.Line(),
			jen.ID("cfg").Dot("Set").Call(jen.ID("metricsProvider"), jen.Lit("prometheus")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("metricsTracer"), jen.Lit("jaeger")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("metricsDBCollectionInterval"), jen.Qual("time", "Second")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("metricsRuntimeCollectionInterval"), jen.Qual("time", "Second")),
			jen.Line(),
			jen.ID("cfg").Dot("Set").Call(jen.ID("dbDebug"), jen.True()),
			jen.ID("cfg").Dot("Set").Call(jen.ID("dbProvider"), jen.Lit("postgres")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("dbDeets"), jen.ID("postgresDBConnDetails")),
			jen.Line(),
			jen.Return().ID("cfg").Dot("WriteConfigAs").Call(jen.ID("filepath")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("coverageConfig").Params(jen.ID("filepath").String()).Params(jen.Error()).Block(
			jen.ID("cfg").Assign().ID("config").Dot("BuildConfig").Call(),
			jen.Line(),
			jen.ID("cfg").Dot("Set").Call(jen.ID("serverHTTPPort"), jen.ID("defaultPort")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("serverDebug"), jen.True()),
			jen.Line(),
			jen.ID("cfg").Dot("Set").Call(jen.ID("frontendDebug"), jen.True()),
			jen.ID("cfg").Dot("Set").Call(jen.ID("frontendStaticFilesDir"), jen.ID("defaultFrontendFilepath")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("frontendCacheStatics"), jen.False()),
			jen.Line(),
			jen.ID("cfg").Dot("Set").Call(jen.ID("authDebug"), jen.False()),
			jen.ID("cfg").Dot("Set").Call(jen.ID("authCookieSecret"), jen.ID("debugCookieSecret")),
			jen.Line(),
			jen.ID("cfg").Dot("Set").Call(jen.ID("dbDebug"), jen.False()),
			jen.ID("cfg").Dot("Set").Call(jen.ID("dbProvider"), jen.Lit("postgres")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("dbDeets"), jen.ID("postgresDBConnDetails")),
			jen.Line(),
			jen.Return().ID("cfg").Dot("WriteConfigAs").Call(jen.ID("filepath")),
		),
		jen.Line(),
	)
	ret.Add(
		jen.Func().ID("productionConfig").Params(jen.ID("filepath").String()).Params(jen.Error()).Block(
			jen.ID("cfg").Assign().ID("config").Dot("BuildConfig").Call(),
			jen.Line(),
			jen.ID("cfg").Dot("Set").Call(jen.ID("metaDebug"), jen.False()),
			jen.ID("cfg").Dot("Set").Call(jen.ID("metaStartupDeadline"), jen.Qual("time", "Minute")),
			jen.Line(),
			jen.ID("cfg").Dot("Set").Call(jen.ID("serverHTTPPort"), jen.ID("defaultPort")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("serverDebug"), jen.False()),
			jen.Line(),
			jen.ID("cfg").Dot("Set").Call(jen.ID("frontendDebug"), jen.False()),
			jen.ID("cfg").Dot("Set").Call(jen.ID("frontendStaticFilesDir"), jen.ID("defaultFrontendFilepath")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("frontendCacheStatics"), jen.False()),
			jen.Line(),
			jen.ID("cfg").Dot("Set").Call(jen.ID("authDebug"), jen.False()),
			jen.ID("cfg").Dot("Set").Call(jen.ID("authCookieDomain"), jen.EmptyString()),
			jen.ID("cfg").Dot("Set").Call(jen.ID("authCookieSecret"), jen.ID("debugCookieSecret")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("authCookieLifetime"), jen.ID("oneDay")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("authSecureCookiesOnly"), jen.False()),
			jen.ID("cfg").Dot("Set").Call(jen.ID("authEnableUserSignup"), jen.True()),
			jen.Line(),
			jen.ID("cfg").Dot("Set").Call(jen.ID("metricsProvider"), jen.Lit("prometheus")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("metricsTracer"), jen.Lit("jaeger")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("metricsDBCollectionInterval"), jen.Qual("time", "Second")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("metricsRuntimeCollectionInterval"), jen.Qual("time", "Second")),
			jen.Line(),
			jen.ID("cfg").Dot("Set").Call(jen.ID("dbDebug"), jen.False()),
			jen.ID("cfg").Dot("Set").Call(jen.ID("dbProvider"), jen.Lit("postgres")),
			jen.ID("cfg").Dot("Set").Call(jen.ID("dbDeets"), jen.ID("postgresDBConnDetails")),
			jen.Line(),
			jen.Return().ID("cfg").Dot("WriteConfigAs").Call(jen.ID("filepath")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildIntegrationTestForDBImplementation").Params(jen.List(jen.ID("dbprov"), jen.ID("dbDeet")).String()).Params(jen.ID("configFunc")).Block(
			jen.Return().Func().Params(jen.ID("filepath").String()).Params(jen.Error()).Block(
				jen.ID("cfg").Assign().Qual(proj.InternalConfigV1Package(), "BuildConfig").Call(),
				jen.Line(),
				jen.ID("cfg").Dot("Set").Call(jen.ID("metaDebug"), jen.False()),
				jen.Line(),
				jen.ID("sd").Assign().Qual("time", "Minute"),
				func() jen.Code {
					if proj.DatabasesIsEnabled(models.MariaDB) {
						jen.If(jen.ID("dbprov").IsEqualTo().ID("mariadb")).Block(
							jen.ID("sd").Equals().Lit(5).Times().Qual("time", "Minute"),
						)
					}
					return jen.Null()
				}(),
				jen.ID("cfg").Dot("Set").Call(jen.ID("metaStartupDeadline"), jen.ID("sd")),
				jen.Line(),
				jen.ID("cfg").Dot("Set").Call(jen.ID("serverHTTPPort"), jen.ID("defaultPort")),
				jen.ID("cfg").Dot("Set").Call(jen.ID("serverDebug"), jen.True()),
				jen.Line(),
				jen.ID("cfg").Dot("Set").Call(jen.ID("frontendStaticFilesDir"), jen.ID("defaultFrontendFilepath")),
				jen.ID("cfg").Dot("Set").Call(jen.ID("authCookieSecret"), jen.ID("debugCookieSecret")),
				jen.Line(),
				jen.ID("cfg").Dot("Set").Call(jen.ID("metricsProvider"), jen.Lit("prometheus")),
				jen.ID("cfg").Dot("Set").Call(jen.ID("metricsTracer"), jen.Lit("jaeger")),
				jen.Line(),
				jen.ID("cfg").Dot("Set").Call(jen.ID("dbDebug"), jen.False()),
				jen.ID("cfg").Dot("Set").Call(jen.ID("dbProvider"), jen.ID("dbprov")),
				jen.ID("cfg").Dot("Set").Call(jen.ID("dbDeets"), jen.ID("dbDeet")),
				jen.Line(),
				jen.Return().ID("cfg").Dot("WriteConfigAs").Call(jen.ID("filepath")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("main").Params().Block(
			jen.For(jen.List(jen.ID("filepath"), jen.ID("fun")).Assign().Range().ID("files")).Block(
				jen.If(jen.Err().Assign().ID("fun").Call(jen.ID("filepath")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Qual("log", "Fatal").Call(jen.Err()),
				),
			),
		),
		jen.Line(),
	)

	return ret
}
