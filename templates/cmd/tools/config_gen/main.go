package config_gen

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(determineConstants(proj)...)

	code.Add(
		jen.Type().ID("configFunc").Func().Params(jen.ID("filePath").String()).Params(jen.Error()),
		jen.Line(),
	)

	code.Add(renderFileMap(proj)...)
	code.Add(buildDevelopmentConfig(proj)...)
	code.Add(buildFrontendTestsConfig(proj)...)
	code.Add(buildCoverageConfig(proj)...)
	code.Add(buildBuildIntegrationTestForDBImplementation(proj)...)
	code.Add(buildMain()...)

	return code
}

func renderFileMap(proj *models.Project) []jen.Code {
	return []jen.Code{
		jen.Var().Defs(
			jen.ID("files").Equals().Map(jen.String()).ID("configFunc").Valuesln(
				jen.Lit("environments/local/config.toml").MapAssign().ID("developmentConfig"),
				jen.Lit("environments/testing/config_files/frontend-tests.toml").MapAssign().ID("frontendTestsConfig"),
				jen.Lit("environments/testing/config_files/coverage.toml").MapAssign().ID("coverageConfig"),
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.Postgres) {
						return jen.Lit("environments/testing/config_files/integration-tests-postgres.toml").MapAssign().ID("buildIntegrationTestForDBImplementation").Call(jen.ID("postgres"), jen.ID("postgresDBConnDetails"))
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.Sqlite) {
						return jen.Lit("environments/testing/config_files/integration-tests-sqlite.toml").MapAssign().ID("buildIntegrationTestForDBImplementation").Call(jen.ID("sqlite"), jen.Lit("/tmp/db"))
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if proj.DatabaseIsEnabled(models.MariaDB) {
						return jen.Lit("environments/testing/config_files/integration-tests-mariadb.toml").MapAssign().ID("buildIntegrationTestForDBImplementation").Call(jen.ID("mariadb"), jen.Litf("dbuser:hunter2@tcp(database:3306)/%s", proj.Name.RouteName()))
					}
					return jen.Null()
				}(),
			),
		),
		jen.Line(),
	}
}

func determineConstants(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.ID("defaultPort").Equals().Lit(8888),
		jen.ID("oneDay").Equals().Lit(24).Times().Qual("time", "Hour"),
		jen.ID("debugCookieSecret").Equals().Lit("HEREISA32CHARSECRETWHICHISMADEUP"),
		jen.ID("defaultFrontendFilepath").Equals().Lit("/frontend"),
		jen.ID("postgresDBConnDetails").Equals().Litf("postgres://dbuser:hunter2@database:5432/%s?sslmode=disable", proj.Name.RouteName()),
		jen.ID("metaDebug").Equals().Lit("meta.debug"),
		jen.ID("metaRunMode").Equals().Lit("meta.run_mode"),
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
	}

	for _, typ := range proj.DataTypes {
		if typ.SearchEnabled {
			lines = append(
				lines,
				jen.IDf("%sSearchIndexPath", typ.Name.PluralUnexportedVarName()).Equals().Litf("search.%s_index_path", typ.Name.PluralRouteName()),
			)
		}
	}

	lines = append(lines,
		jen.Line(),
		jen.Comment("run modes"),
		jen.ID("developmentEnv").Equals().Lit("development"),
		jen.ID("testingEnv").Equals().Lit("testing"),
		jen.Line(),
		jen.Comment("database providers"),
	)

	if proj.DatabaseIsEnabled(models.Postgres) {
		lines = append(lines, jen.ID(string(models.Postgres)).Equals().Lit(string(models.Postgres)))
	}

	if proj.DatabaseIsEnabled(models.Sqlite) {
		lines = append(lines, jen.ID(string(models.Sqlite)).Equals().Lit(string(models.Sqlite)))
	}

	if proj.DatabaseIsEnabled(models.MariaDB) {
		lines = append(lines, jen.ID(string(models.MariaDB)).Equals().Lit(string(models.MariaDB)))
	}

	lines = append(lines,
		jen.Line(),
		jen.Comment("search index paths"),
	)

	for _, typ := range proj.DataTypes {
		if typ.SearchEnabled {
			lines = append(
				lines,
				jen.IDf("default%sSearchIndexPath", typ.Name.Plural()).Equals().Litf("%s.bleve", typ.Name.PluralRouteName()),
			)
		}
	}

	return []jen.Code{
		jen.Const().Defs(lines...),
		jen.Line(),
	}
}

func buildDevelopmentConfig(proj *models.Project) []jen.Code {
	block := []jen.Code{
		jen.ID("cfg").Assign().ID("config").Dot("BuildConfig").Call(),
		jen.Line(),
		jen.ID("cfg").Dot("Set").Call(jen.ID("metaRunMode"), jen.ID("developmentEnv")),
		jen.ID("cfg").Dot("Set").Call(jen.ID("metaDebug"), jen.True()),
		jen.ID("cfg").Dot("Set").Call(jen.ID("metaStartupDeadline"), jen.Qual("time", "Minute")),
		jen.Line(),
		jen.ID("cfg").Dot("Set").Call(jen.ID("serverHTTPPort"), jen.ID("defaultPort")),
		jen.ID("cfg").Dot("Set").Call(jen.ID("serverDebug"), jen.True()),
		jen.Line(),
		jen.ID("cfg").Dot("Set").Call(jen.ID("frontendDebug"), jen.True()),
		jen.ID("cfg").Dot("Set").Call(jen.ID("frontendStaticFilesDir"), jen.ID("defaultFrontendFilepath")),
		jen.ID("cfg").Dot("Set").Call(jen.ID("frontendCacheStatics"), jen.False()),
		jen.Line(),
		jen.ID("cfg").Dot("Set").Call(jen.ID("authDebug"), jen.True()),
		jen.ID("cfg").Dot("Set").Call(jen.ID("authCookieDomain"), jen.Lit("localhost")),
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
		jen.ID("cfg").Dot("Set").Call(jen.ID("dbProvider"), jen.ID("postgres")),
		jen.ID("cfg").Dot("Set").Call(jen.ID("dbDeets"), jen.ID("postgresDBConnDetails")),
		jen.Line(),
	}

	for _, typ := range proj.DataTypes {
		if typ.SearchEnabled {
			block = append(block,
				jen.ID("cfg").Dot("Set").Call(jen.IDf("%sSearchIndexPath", typ.Name.PluralUnexportedVarName()), jen.IDf("default%sSearchIndexPath", typ.Name.Plural())),
			)
		}
	}

	block = append(block,
		jen.Line(),
		jen.If(jen.ID("writeErr").Assign().ID("cfg").Dot("WriteConfigAs").Call(jen.ID("filePath")), jen.ID("writeErr").DoesNotEqual().Nil()).Body(
			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("error writing developmentEnv config: %w"), jen.ID("writeErr"))),
		),
		jen.Line(),
		jen.Return().Nil(),
	)

	lines := []jen.Code{
		jen.Func().ID("developmentConfig").Params(jen.ID("filePath").String()).Params(jen.Error()).Body(
			block...,
		),
		jen.Line(),
		jen.Line(),
	}

	return lines
}

func buildFrontendTestsConfig(proj *models.Project) []jen.Code {
	block := []jen.Code{
		jen.ID("cfg").Assign().ID("config").Dot("BuildConfig").Call(),
		jen.Line(),
		jen.ID("cfg").Dot("Set").Call(jen.ID("metaRunMode"), jen.ID("developmentEnv")),
		jen.ID("cfg").Dot("Set").Call(jen.ID("metaStartupDeadline"), jen.Qual("time", "Minute")),
		jen.Line(),
		jen.ID("cfg").Dot("Set").Call(jen.ID("serverHTTPPort"), jen.ID("defaultPort")),
		jen.ID("cfg").Dot("Set").Call(jen.ID("serverDebug"), jen.True()),
		jen.Line(),
		jen.ID("cfg").Dot("Set").Call(jen.ID("frontendDebug"), jen.True()),
		jen.ID("cfg").Dot("Set").Call(jen.ID("frontendStaticFilesDir"), jen.ID("defaultFrontendFilepath")),
		jen.ID("cfg").Dot("Set").Call(jen.ID("frontendCacheStatics"), jen.False()),
		jen.Line(),
		jen.ID("cfg").Dot("Set").Call(jen.ID("authDebug"), jen.True()),
		jen.ID("cfg").Dot("Set").Call(jen.ID("authCookieDomain"), jen.Lit("localhost")),
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
		jen.ID("cfg").Dot("Set").Call(jen.ID("dbProvider"), jen.ID("postgres")),
		jen.ID("cfg").Dot("Set").Call(jen.ID("dbDeets"), jen.ID("postgresDBConnDetails")),
		jen.Line(),
	}

	for _, typ := range proj.DataTypes {
		if typ.SearchEnabled {
			block = append(block,
				jen.ID("cfg").Dot("Set").Call(jen.IDf("%sSearchIndexPath", typ.Name.PluralUnexportedVarName()), jen.IDf("default%sSearchIndexPath", typ.Name.Plural())),
			)
		}
	}

	block = append(block,
		jen.Line(),
		jen.If(jen.ID("writeErr").Assign().ID("cfg").Dot("WriteConfigAs").Call(jen.ID("filePath")), jen.ID("writeErr").DoesNotEqual().Nil()).Body(
			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("error writing developmentEnv config: %w"), jen.ID("writeErr"))),
		),
		jen.Line(),
		jen.Return().Nil(),
	)

	lines := []jen.Code{
		jen.Func().ID("frontendTestsConfig").Params(jen.ID("filePath").String()).Params(jen.Error()).Body(
			block...,
		),
		jen.Line(),
	}

	return lines
}

func buildCoverageConfig(proj *models.Project) []jen.Code {
	block := []jen.Code{
		jen.ID("cfg").Assign().ID("config").Dot("BuildConfig").Call(),
		jen.Line(),
		jen.ID("cfg").Dot("Set").Call(jen.ID("metaRunMode"), jen.ID("testingEnv")),
		jen.ID("cfg").Dot("Set").Call(jen.ID("metaDebug"), jen.True()),
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
		jen.ID("cfg").Dot("Set").Call(jen.ID("dbProvider"), jen.ID("postgres")),
		jen.ID("cfg").Dot("Set").Call(jen.ID("dbDeets"), jen.ID("postgresDBConnDetails")),
		jen.Line(),
	}

	for _, typ := range proj.DataTypes {
		if typ.SearchEnabled {
			block = append(block,
				jen.ID("cfg").Dot("Set").Call(jen.IDf("%sSearchIndexPath", typ.Name.PluralUnexportedVarName()), jen.IDf("default%sSearchIndexPath", typ.Name.Plural())),
				jen.Line(),
			)
		}
	}

	block = append(block,
		jen.If(jen.ID("writeErr").Assign().ID("cfg").Dot("WriteConfigAs").Call(jen.ID("filePath")), jen.ID("writeErr").DoesNotEqual().Nil()).Body(
			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("error writing coverage config: %w"), jen.ID("writeErr"))),
		),
		jen.Line(),
		jen.Return().Nil(),
	)

	lines := []jen.Code{
		jen.Func().ID("coverageConfig").Params(jen.ID("filePath").String()).Params(jen.Error()).Body(
			block...,
		),
		jen.Line(),
	}

	return lines
}

func buildBuildIntegrationTestForDBImplementation(proj *models.Project) []jen.Code {
	block := []jen.Code{
		jen.ID("cfg").Assign().Qual(proj.InternalConfigV1Package(), "BuildConfig").Call(),
		jen.Line(),
		jen.ID("cfg").Dot("Set").Call(jen.ID("metaRunMode"), jen.ID("testingEnv")),
		jen.ID("cfg").Dot("Set").Call(jen.ID("metaDebug"), jen.False()),
		jen.Line(),
		jen.ID("sd").Assign().Qual("time", "Minute"),
		func() jen.Code {
			if proj.DatabaseIsEnabled(models.MariaDB) {
				return jen.If(jen.ID("dbVendor").IsEqualTo().ID("mariadb")).Body(
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
		jen.ID("cfg").Dot("Set").Call(jen.ID("dbProvider"), jen.ID("dbVendor")),
		jen.ID("cfg").Dot("Set").Call(jen.ID("dbDeets"), jen.ID("dbDetails")),
		jen.Line(),
	}

	for _, typ := range proj.DataTypes {
		if typ.SearchEnabled {
			block = append(block,
				jen.ID("cfg").Dot("Set").Call(jen.IDf("%sSearchIndexPath", typ.Name.PluralUnexportedVarName()), jen.IDf("default%sSearchIndexPath", typ.Name.Plural())),
				jen.Line(),
			)
		}
	}

	block = append(block,
		jen.If(jen.ID("writeErr").Assign().ID("cfg").Dot("WriteConfigAs").Call(jen.ID("filePath")), jen.ID("writeErr").DoesNotEqual().Nil()).Body(
			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("error writing integration test config for %s: %w"), jen.ID("dbVendor"), jen.ID("writeErr"))),
		),
		jen.Line(),
		jen.Return().Nil(),
	)

	lines := []jen.Code{
		jen.Func().ID("buildIntegrationTestForDBImplementation").Params(jen.List(jen.ID("dbVendor"), jen.ID("dbDetails")).String()).Params(jen.ID("configFunc")).Body(
			jen.Return().Func().Params(jen.ID("filePath").String()).Params(jen.Error()).Body(
				block...,
			),
		),
		jen.Line(),
	}

	return lines
}

func buildMain() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("main").Params().Body(
			jen.For(jen.List(jen.ID("filePath"), jen.ID("fun")).Assign().Range().ID("files")).Body(
				jen.If(jen.Err().Assign().ID("fun").Call(jen.ID("filePath")), jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Qual("log", "Fatalf").Call(jen.Lit("error rendering %s: %v"), jen.ID("filePath"), jen.Err()),
				),
			),
		),
		jen.Line(),
	}

	return lines
}
