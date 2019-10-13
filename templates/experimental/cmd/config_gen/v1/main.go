package v1

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func mainDotGo() *jen.File {
	ret := jen.NewFile("main")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("defaultPort").Op("=").Lit(8888).Var().ID("oneDay").Op("=").Lit(24).Op("*").Qual("time", "Hour").Var().ID("debugCookieSecret").Op("=").Lit("HEREISA32CHARSECRETWHICHISMADEUP").Var().ID("defaultFrontendFilepath").Op("=").Lit("/frontend").Var().ID("postgresDBConnDetails").Op("=").Lit("postgres://dbuser:hunter2@database:5432/todo?sslmode=disable").Var().ID("metaDebug").Op("=").Lit("meta.debug").Var().ID("metaStartupDeadline").Op("=").Lit("meta.startup_deadline").Var().ID("serverHTTPPort").Op("=").Lit("server.http_port").Var().ID("serverDebug").Op("=").Lit("server.debug").Var().ID("frontendDebug").Op("=").Lit("frontend.debug").Var().ID("frontendStaticFilesDir").Op("=").Lit("frontend.static_files_directory").Var().ID("frontendCacheStatics").Op("=").Lit("frontend.cache_static_files").Var().ID("authDebug").Op("=").Lit("auth.debug").Var().ID("authCookieDomain").Op("=").Lit("auth.cookie_domain").Var().ID("authCookieSecret").Op("=").Lit("auth.cookie_secret").Var().ID("authCookieLifetime").Op("=").Lit("auth.cookie_lifetime").Var().ID("authSecureCookiesOnly").Op("=").Lit("auth.secure_cookies_only").Var().ID("authEnableUserSignup").Op("=").Lit("auth.enable_user_signup").Var().ID("metricsProvider").Op("=").Lit("metrics.metrics_provider").Var().ID("metricsTracer").Op("=").Lit("metrics.tracing_provider").Var().ID("metricsDBCollectionInterval").Op("=").Lit("metrics.database_metrics_collection_interval").Var().ID("metricsRuntimeCollectionInterval").Op("=").Lit("metrics.runtime_metrics_collection_interval").Var().ID("dbDebug").Op("=").Lit("database.debug").Var().ID("dbProvider").Op("=").Lit("database.provider").Var().ID("dbDeets").Op("=").Lit("database.connection_details"),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("configFunc").Params(jen.ID("filepath").ID("string")).Params(jen.ID("error")),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("files").Op("=").Map(jen.ID("string")).ID("configFunc").Valuesln(jen.Lit("config_files/coverage.toml").Op(":").ID("coverageConfig"), jen.Lit("config_files/development.toml").Op(":").ID("developmentConfig"), jen.Lit("config_files/integration-tests-postgres.toml").Op(":").ID("buildIntegrationTestForDBImplementation").Call(jen.Lit("postgres"), jen.ID("postgresDBConnDetails")), jen.Lit("config_files/integration-tests-sqlite.toml").Op(":").ID("buildIntegrationTestForDBImplementation").Call(jen.Lit("sqlite"), jen.Lit("/tmp/db")), jen.Lit("config_files/integration-tests-mariadb.toml").Op(":").ID("buildIntegrationTestForDBImplementation").Call(jen.Lit("mariadb"), jen.Lit("dbuser:hunter2@tcp(database:3306)/todo")), jen.Lit("config_files/production.toml").Op(":").ID("productionConfig")),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("developmentConfig").Params(jen.ID("filepath").ID("string")).Params(jen.ID("error")).Block(
		jen.ID("cfg").Op(":=").ID("config").Dot(
			"BuildConfig",
		).Call(),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("metaStartupDeadline"), jen.Qual("time", "Minute")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("serverHTTPPort"), jen.ID("defaultPort")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("serverDebug"), jen.ID("true")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("frontendDebug"), jen.ID("true")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("frontendStaticFilesDir"), jen.ID("defaultFrontendFilepath")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("frontendCacheStatics"), jen.ID("false")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("authDebug"), jen.ID("true")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("authCookieDomain"), jen.Lit("")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("authCookieSecret"), jen.ID("debugCookieSecret")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("authCookieLifetime"), jen.ID("oneDay")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("authSecureCookiesOnly"), jen.ID("false")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("authEnableUserSignup"), jen.ID("true")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("metricsProvider"), jen.Lit("prometheus")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("metricsTracer"), jen.Lit("jaeger")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("metricsDBCollectionInterval"), jen.Qual("time", "Second")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("metricsRuntimeCollectionInterval"), jen.Qual("time", "Second")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("dbDebug"), jen.ID("true")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("dbProvider"), jen.Lit("postgres")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("dbDeets"), jen.ID("postgresDBConnDetails")),
		jen.Return().ID("cfg").Dot(
			"WriteConfigAs",
		).Call(jen.ID("filepath")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("coverageConfig").Params(jen.ID("filepath").ID("string")).Params(jen.ID("error")).Block(
		jen.ID("cfg").Op(":=").ID("config").Dot(
			"BuildConfig",
		).Call(),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("serverHTTPPort"), jen.ID("defaultPort")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("serverDebug"), jen.ID("true")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("frontendDebug"), jen.ID("true")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("frontendStaticFilesDir"), jen.ID("defaultFrontendFilepath")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("frontendCacheStatics"), jen.ID("false")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("authDebug"), jen.ID("false")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("authCookieSecret"), jen.ID("debugCookieSecret")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("dbDebug"), jen.ID("false")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("dbProvider"), jen.Lit("postgres")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("dbDeets"), jen.ID("postgresDBConnDetails")),
		jen.Return().ID("cfg").Dot(
			"WriteConfigAs",
		).Call(jen.ID("filepath")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("productionConfig").Params(jen.ID("filepath").ID("string")).Params(jen.ID("error")).Block(
		jen.ID("cfg").Op(":=").ID("config").Dot(
			"BuildConfig",
		).Call(),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("metaDebug"), jen.ID("false")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("metaStartupDeadline"), jen.Qual("time", "Minute")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("serverHTTPPort"), jen.ID("defaultPort")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("serverDebug"), jen.ID("false")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("frontendDebug"), jen.ID("false")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("frontendStaticFilesDir"), jen.ID("defaultFrontendFilepath")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("frontendCacheStatics"), jen.ID("false")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("authDebug"), jen.ID("false")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("authCookieDomain"), jen.Lit("")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("authCookieSecret"), jen.ID("debugCookieSecret")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("authCookieLifetime"), jen.ID("oneDay")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("authSecureCookiesOnly"), jen.ID("false")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("authEnableUserSignup"), jen.ID("true")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("metricsProvider"), jen.Lit("prometheus")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("metricsTracer"), jen.Lit("jaeger")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("metricsDBCollectionInterval"), jen.Qual("time", "Second")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("metricsRuntimeCollectionInterval"), jen.Qual("time", "Second")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("dbDebug"), jen.ID("false")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("dbProvider"), jen.Lit("postgres")),
		jen.ID("cfg").Dot(
			"Set",
		).Call(jen.ID("dbDeets"), jen.ID("postgresDBConnDetails")),
		jen.Return().ID("cfg").Dot(
			"WriteConfigAs",
		).Call(jen.ID("filepath")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("buildIntegrationTestForDBImplementation").Params(jen.List(jen.ID("dbprov"), jen.ID("dbDeet")).ID("string")).Params(jen.ID("configFunc")).Block(
		jen.Return().Func().Params(jen.ID("filepath").ID("string")).Params(jen.ID("error")).Block(
			jen.ID("cfg").Op(":=").ID("config").Dot(
				"BuildConfig",
			).Call(),
			jen.ID("cfg").Dot(
				"Set",
			).Call(jen.ID("metaDebug"), jen.ID("false")),
			jen.ID("cfg").Dot(
				"Set",
			).Call(jen.ID("metaStartupDeadline"), jen.Qual("time", "Minute")),
			jen.ID("cfg").Dot(
				"Set",
			).Call(jen.ID("serverHTTPPort"), jen.ID("defaultPort")),
			jen.ID("cfg").Dot(
				"Set",
			).Call(jen.ID("serverDebug"), jen.ID("true")),
			jen.ID("cfg").Dot(
				"Set",
			).Call(jen.ID("frontendStaticFilesDir"), jen.ID("defaultFrontendFilepath")),
			jen.ID("cfg").Dot(
				"Set",
			).Call(jen.ID("authCookieSecret"), jen.ID("debugCookieSecret")),
			jen.ID("cfg").Dot(
				"Set",
			).Call(jen.ID("metricsProvider"), jen.Lit("prometheus")),
			jen.ID("cfg").Dot(
				"Set",
			).Call(jen.ID("metricsTracer"), jen.Lit("jaeger")),
			jen.ID("cfg").Dot(
				"Set",
			).Call(jen.ID("dbDebug"), jen.ID("false")),
			jen.ID("cfg").Dot(
				"Set",
			).Call(jen.ID("dbProvider"), jen.ID("dbprov")),
			jen.ID("cfg").Dot(
				"Set",
			).Call(jen.ID("dbDeets"), jen.ID("dbDeet")),
			jen.Return().ID("cfg").Dot(
				"WriteConfigAs",
			).Call(jen.ID("filepath")),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("main").Params().Block(
		jen.For(jen.List(jen.ID("filepath"), jen.ID("fun")).Op(":=").Range().ID("files")).Block(
			jen.If(jen.ID("err").Op(":=").ID("fun").Call(jen.ID("filepath")), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Qual("log", "Fatal").Call(jen.ID("err")),
			),
		),
	),

		jen.Line(),
	)
	return ret
}
