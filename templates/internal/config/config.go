package config

import (
	"fmt"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildConfigConstantVariableDeclarations()...)
	code.Add(buildConfigVariableDeclarations()...)
	code.Add(buildInit()...)
	code.Add(buildTypeDefinitions(proj)...)
	code.Add(buildEncodeToFile()...)
	code.Add(buildBuildConfig()...)
	code.Add(buildParseConfigFile()...)
	code.Add(buildRandString()...)

	return code
}

func buildConfigConstantVariableDeclarations() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.Comment("DevelopmentRunMode is the run mode for a development environment"),
			jen.ID("DevelopmentRunMode").ID("runMode").Equals().Lit("development"),
			jen.Comment("TestingRunMode is the run mode for a testing environment"),
			jen.ID("TestingRunMode").ID("runMode").Equals().Lit("testing"),
			jen.Comment("ProductionRunMode is the run mode for a production environment"),
			jen.ID("ProductionRunMode").ID("runMode").Equals().Lit("production"),
			jen.Line(),
			jen.ID("defaultStartupDeadline").Equals().Qual("time", "Minute"),
			jen.ID("defaultRunMode").Equals().ID("DevelopmentRunMode"),
			jen.ID("defaultCookieLifetime").Equals().Lit(24).Times().Qual("time", "Hour"),
			jen.ID("defaultMetricsCollectionInterval").Equals().Lit(2).Times().Qual("time", "Second"),
			jen.ID("defaultDatabaseMetricsCollectionInterval").Equals().Lit(2).Times().Qual("time", "Second"),
			jen.ID("randStringSize").Equals().Lit(32),
		),
		jen.Line(),
	}

	return lines
}

func buildConfigVariableDeclarations() []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.ID("validModes").Equals().Map(jen.ID("runMode")).Struct().Valuesln(
				jen.ID("DevelopmentRunMode").MapAssign().Values(),
				jen.ID("TestingRunMode").MapAssign().Values(),
				jen.ID("ProductionRunMode").MapAssign().Values(),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildInit() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("init").Params().Body(
			jen.ID("b").Assign().ID("make").Call(jen.Index().Byte(), jen.Lit(64)),
			jen.If(jen.List(jen.Underscore(), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("panic").Call(jen.Err()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTypeDefinitions(proj *models.Project) []jen.Code {
	searchSettingsFields := []jen.Code{}
	for _, typ := range proj.DataTypes {
		if typ.SearchEnabled {
			searchSettingsFields = append(searchSettingsFields,
				jen.Commentf("%sIndexPath indicates where our %s search index files should go.", typ.Name.Plural(), typ.Name.PluralCommonName()),
				jen.IDf("%sIndexPath", typ.Name.Plural()).Qual(proj.InternalSearchV1Package(), "IndexPath").Tag(map[string]string{
					"mapstructure": fmt.Sprintf("%s_index_path", typ.Name.PluralRouteName()),
					"json":         fmt.Sprintf("%s_index_path", typ.Name.PluralRouteName()),
					"toml":         fmt.Sprintf("%s_index_path,omitempty", typ.Name.PluralRouteName()),
				}),
			)
		}
	}

	lines := []jen.Code{
		jen.ID("runMode").String(),
		jen.Line(),
		jen.Comment("MetaSettings is primarily used for development."),
		jen.ID("MetaSettings").Struct(
			jen.Comment("Debug enables debug mode service-wide"),
			jen.Comment("NOTE: this debug should override all other debugs, which is to say, if this is enabled, all of them are enabled."),
			jen.ID("Debug").Bool().Tag(map[string]string{
				"mapstructure": "debug",
				"json":         "debug",
				"toml":         "debug,omitempty",
			}),
			jen.Comment("StartupDeadline indicates how long the service can take to spin up. This includes database migrations, configuring services, etc."),
			jen.ID("StartupDeadline").Qual("time", "Duration").Tag(map[string]string{
				"mapstructure": "startup_deadline",
				"json":         "startup_deadline",
				"toml":         "startup_deadline,omitempty",
			}),
			jen.Comment("RunMode indicates the current run mode"),
			jen.ID("RunMode").ID("runMode").Tag(map[string]string{
				"mapstructure": "run_mode",
				"json":         "run_mode",
				"toml":         "run_mode,omitempty",
			}),
		),
		jen.Line(),
		jen.Comment("ServerSettings describes the settings pertinent to the HTTP serving portion of the service."),
		jen.ID("ServerSettings").Struct(
			jen.Comment("Debug determines if debug logging or other development conditions are active."),
			jen.ID("Debug").Bool().Tag(map[string]string{
				"mapstructure": "debug",
				"json":         "debug",
				"toml":         "debug,omitempty",
			}),
			jen.Comment("HTTPPort indicates which port to serve HTTP traffic on."),
			jen.ID("HTTPPort").ID("uint16").Tag(map[string]string{
				"mapstructure": "http_port",
				"json":         "http_port",
				"toml":         "http_port,omitempty",
			}),
		),
		jen.Line(),
		jen.Comment("FrontendSettings describes the settings pertinent to the frontend."),
		jen.ID("FrontendSettings").Struct(
			jen.Comment("StaticFilesDirectory indicates which directory contains our static files for the frontend (i.e. CSS/JS/HTML files)"),
			jen.ID("StaticFilesDirectory").String().Tag(map[string]string{
				"mapstructure": "static_files_directory",
				"json":         "static_files_directory",
				"toml":         "static_files_directory,omitempty",
			}),
			jen.Comment("Debug determines if debug logging or other development conditions are active."),
			jen.ID("Debug").Bool().Tag(map[string]string{
				"mapstructure": "debug",
				"json":         "debug",
				"toml":         "debug,omitempty",
			}),
			jen.Comment("CacheStaticFiles indicates whether or not to load the static files directory into memory via afero's MemMapFs."),
			jen.ID("CacheStaticFiles").Bool().Tag(map[string]string{
				"mapstructure": "cache_static_files",
				"json":         "cache_static_files",
				"toml":         "cache_static_files,omitempty",
			}),
		),
		jen.Line(),
		jen.Comment("AuthSettings represents our authentication configuration."),
		jen.ID("AuthSettings").Struct(
			jen.Comment("CookieDomain indicates what domain the cookies will have set for them."),
			jen.ID("CookieDomain").String().Tag(map[string]string{
				"mapstructure": "cookie_domain",
				"json":         "cookie_domain",
				"toml":         "cookie_domain,omitempty",
			}),
			jen.Comment("CookieSecret indicates the secret the cookie builder should use."),
			jen.ID("CookieSecret").String().Tag(map[string]string{
				"mapstructure": "cookie_secret",
				"json":         "cookie_secret",
				"toml":         "cookie_secret,omitempty",
			}),
			jen.Comment("CookieLifetime indicates how long the cookies built should last."),
			jen.ID("CookieLifetime").Qual("time", "Duration").Tag(map[string]string{
				"mapstructure": "cookie_lifetime",
				"json":         "cookie_lifetime",
				"toml":         "cookie_lifetime,omitempty",
			}),
			jen.Comment("Debug determines if debug logging or other development conditions are active."),
			jen.ID("Debug").Bool().Tag(map[string]string{
				"mapstructure": "debug",
				"json":         "debug",
				"toml":         "debug,omitempty",
			}),
			jen.Comment("SecureCookiesOnly indicates if the cookies built should be marked as HTTPS only."),
			jen.ID("SecureCookiesOnly").Bool().Tag(map[string]string{
				"mapstructure": "secure_cookies_only",
				"json":         "secure_cookies_only",
				"toml":         "secure_cookies_only,omitempty",
			}),
			jen.Comment("EnableUserSignup enables user signups."),
			jen.ID("EnableUserSignup").Bool().Tag(map[string]string{
				"mapstructure": "enable_user_signup",
				"json":         "enable_user_signup",
				"toml":         "enable_user_signup,omitempty",
			}),
		),
		jen.Line(),
		jen.Comment("DatabaseSettings represents our database configuration."),
		jen.ID("DatabaseSettings").Struct(
			jen.Comment("Debug determines if debug logging or other development conditions are active."),
			jen.ID("Debug").Bool().Tag(map[string]string{
				"mapstructure": "debug",
				"json":         "debug",
				"toml":         "debug,omitempty",
			}),
			jen.Comment("Provider indicates what database we'll connect to (postgres, mysql, etc.)"),
			jen.ID("Provider").String().Tag(map[string]string{
				"mapstructure": "provider",
				"json":         "provider",
				"toml":         "provider,omitempty",
			}),
			jen.Comment("ConnectionDetails indicates how our database driver should connect to the instance."),
			jen.ID("ConnectionDetails").Qual(proj.DatabaseV1Package(), "ConnectionDetails").Tag(map[string]string{
				"mapstructure": "connection_details",
				"json":         "connection_details",
				"toml":         "connection_details,omitempty",
			}),
		),
		jen.Line(),
		jen.Comment("MetricsSettings contains settings about how we report our metrics."),
		jen.ID("MetricsSettings").Struct(
			jen.Comment("MetricsProvider indicates where our metrics should go."),
			jen.ID("MetricsProvider").ID("metricsProvider").Tag(map[string]string{
				"mapstructure": "metrics_provider",
				"json":         "metrics_provider",
				"toml":         "metrics_provider,omitempty",
			}),
			jen.Comment("TracingProvider indicates where our traces should go."),
			jen.ID("TracingProvider").ID("tracingProvider").Tag(map[string]string{
				"mapstructure": "tracing_provider",
				"json":         "tracing_provider",
				"toml":         "tracing_provider,omitempty",
			}),
			jen.Comment("DBMetricsCollectionInterval is the interval we collect database statistics at."),
			jen.ID("DBMetricsCollectionInterval").Qual("time", "Duration").Tag(map[string]string{
				"mapstructure": "database_metrics_collection_interval",
				"json":         "database_metrics_collection_interval",
				"toml":         "database_metrics_collection_interval,omitempty",
			}),
			jen.Comment("RuntimeMetricsCollectionInterval  is the interval we collect runtime statistics at."),
			jen.ID("RuntimeMetricsCollectionInterval").Qual("time", "Duration").Tag(map[string]string{
				"mapstructure": "runtime_metrics_collection_interval",
				"json":         "runtime_metrics_collection_interval",
				"toml":         "runtime_metrics_collection_interval,omitempty",
			}),
		),
		jen.Line(),
	}

	if proj.SearchEnabled() {
		lines = append(lines,
			jen.Comment("SearchSettings contains settings regarding search indices."),
			jen.ID("SearchSettings").Struct(
				searchSettingsFields...,
			),
			jen.Line(),
		)
	}

	lines = append(lines,
		jen.Comment("ServerConfig is our server configuration struct. It is comprised of all the other setting structs"),
		jen.Comment("For information on this structs fields, refer to their definitions."),
		jen.ID("ServerConfig").Struct(
			jen.ID("Meta").ID("MetaSettings").Tag(map[string]string{
				"mapstructure": "meta",
				"json":         "meta",
				"toml":         "meta,omitempty",
			}),
			jen.ID("Frontend").ID("FrontendSettings").Tag(map[string]string{
				"mapstructure": "frontend",
				"json":         "frontend",
				"toml":         "frontend,omitempty",
			}),
			jen.ID("Auth").ID("AuthSettings").Tag(map[string]string{
				"mapstructure": "auth",
				"json":         "auth",
				"toml":         "auth,omitempty",
			}),
			jen.ID("Server").ID("ServerSettings").Tag(map[string]string{
				"mapstructure": "server",
				"json":         "server",
				"toml":         "server,omitempty",
			}),
			jen.ID("Database").ID("DatabaseSettings").Tag(map[string]string{
				"mapstructure": "database",
				"json":         "database",
				"toml":         "database,omitempty",
			}),
			jen.ID("Metrics").ID("MetricsSettings").Tag(map[string]string{
				"mapstructure": "metrics",
				"json":         "metrics",
				"toml":         "metrics,omitempty",
			}),
			func() jen.Code {
				if proj.SearchEnabled() {
					return jen.ID("Search").ID("SearchSettings").Tag(map[string]string{
						"mapstructure": "search",
						"json":         "search",
						"toml":         "search,omitempty",
					})
				} else {
					return jen.Null()
				}
			}(),
		),
		jen.Line(),
	)

	return []jen.Code{jen.Type().Defs(lines...)}
}

func buildEncodeToFile() []jen.Code {
	lines := []jen.Code{
		jen.Comment("EncodeToFile renders your config to a file given your favorite encoder."),
		jen.Line(),
		jen.Func().Params(
			jen.ID("cfg").PointerTo().ID("ServerConfig")).ID("EncodeToFile").
			Params(
				jen.ID("path").String(),
				jen.ID("marshaler").Func().Params(jen.ID("v").Interface()).Params(jen.Index().Byte(), jen.Error()),
			).Params(jen.Error()).Body(
			jen.List(jen.ID("byteSlice"), jen.Err()).Assign().ID("marshaler").Call(jen.PointerTo().ID("cfg")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().Err(),
			),
			jen.Line(),
			jen.Return().Qual("io/ioutil", "WriteFile").Call(jen.ID("path"), jen.ID("byteSlice"), jen.Op("0600")),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildConfig() []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildConfig is a constructor function that initializes a viper config."),
		jen.Line(),
		jen.Func().ID("BuildConfig").Params().Params(jen.PointerTo().Qual("github.com/spf13/viper", "Viper")).Body(
			jen.ID("cfg").Assign().ID("viper").Dot("New").Call(),
			jen.Line(),
			jen.Comment("meta stuff."),
			jen.ID("cfg").Dot("SetDefault").Call(jen.Lit("meta.run_mode"), jen.ID("defaultRunMode")),
			jen.ID("cfg").Dot("SetDefault").Call(jen.Lit("meta.startup_deadline"), jen.ID("defaultStartupDeadline")),
			jen.Line(),
			jen.Comment("auth stuff."),
			jen.ID("cfg").Dot("SetDefault").Call(jen.Lit("auth.cookie_lifetime"), jen.ID("defaultCookieLifetime")),
			jen.ID("cfg").Dot("SetDefault").Call(jen.Lit("auth.enable_user_signup"), jen.True()),
			jen.Line(),
			jen.Comment("metrics stuff."),
			jen.ID("cfg").Dot("SetDefault").Call(jen.Lit("metrics.database_metrics_collection_interval"), jen.ID("defaultMetricsCollectionInterval")),
			jen.ID("cfg").Dot("SetDefault").Call(jen.Lit("metrics.runtime_metrics_collection_interval"), jen.ID("defaultDatabaseMetricsCollectionInterval")),
			jen.Line(),
			jen.Comment("server stuff."),
			jen.ID("cfg").Dot("SetDefault").Call(jen.Lit("server.http_port"), jen.Lit(80)),
			jen.Line(),
			jen.Return().ID("cfg"),
		),
		jen.Line(),
	}

	return lines
}

func buildParseConfigFile() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ParseConfigFile parses a configuration file."),
		jen.Line(),
		jen.Func().ID("ParseConfigFile").Params(jen.ID("filename").String()).Params(jen.PointerTo().ID("ServerConfig"), jen.Error()).Body(
			jen.ID("cfg").Assign().ID("BuildConfig").Call(),
			jen.ID("cfg").Dot("SetConfigFile").Call(jen.ID("filename")),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("cfg").Dot("ReadInConfig").Call(), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("trying to read the config file: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Var().ID("serverConfig").PointerTo().ID("ServerConfig"),
			jen.If(jen.Err().Assign().ID("cfg").Dot("Unmarshal").Call(jen.AddressOf().ID("serverConfig")), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("trying to unmarshal the config: %w"), jen.Err())),
			),
			jen.Line(),
			jen.If(
				jen.List(jen.Underscore(), jen.ID("ok")).Assign().ID("validModes").Index(
					jen.ID("serverConfig").Dot("Meta").Dot("RunMode"),
				),
				jen.Not().ID("ok"),
			).Body(
				jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("invalid run mode: %q"), jen.ID("serverConfig").Dot("Meta").Dot("RunMode"))),
			),
			jen.Line(),
			jen.Comment("set the cookie secret to something (relatively) secure if not provided"),
			jen.If(jen.ID("serverConfig").Dot("Auth").Dot("CookieSecret").IsEqualTo().EmptyString()).Body(
				jen.ID("serverConfig").Dot("Auth").Dot("CookieSecret").Equals().ID("randString").Call(jen.ID("randStringSize")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("serverConfig"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildRandString() []jen.Code {
	lines := []jen.Code{
		jen.Comment("randString produces a random string."),
		jen.Line(),
		jen.Comment("https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/"),
		jen.Line(),
		jen.Func().ID("randString").Params(jen.ID("size").Uint()).Params(jen.String()).Body(
			jen.ID("b").Assign().ID("make").Call(jen.Index().Byte(), jen.ID("size")),
			jen.If(jen.List(jen.Underscore(), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Return().Qual("encoding/base32", "StdEncoding").Dot("WithPadding").Call(jen.Qual("encoding/base32", "NoPadding")).Dot("EncodeToString").Call(jen.ID("b")),
		),
		jen.Line(),
	}

	return lines
}
