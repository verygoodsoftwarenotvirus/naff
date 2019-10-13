package config

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func configDotGo() *jen.File {
	ret := jen.NewFile("config")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("defaultStartupDeadline").Op("=").Qual("time", "Minute").Var().ID("defaultCookieLifetime").Op("=").Lit(24).Op("*").Qual("time", "Hour").Var().ID("defaultMetricsCollectionInterval").Op("=").Lit(2).Op("*").Qual("time", "Second").Var().ID("defaultDatabaseMetricsCollectionInterval").Op("=").Lit(2).Op("*").Qual("time", "Second").Var().ID("randStringSize").Op("=").Lit(32),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("init").Params().Block(
		jen.ID("b").Op(":=").ID("make").Call(jen.Index().ID("byte"), jen.Lit(64)),
		jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("panic").Call(jen.ID("err")),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("MetaSettings").Struct(jen.ID("Debug").ID("bool"), jen.ID("StartupDeadline").Qual("time", "Duration")).Type().ID("ServerSettings").Struct(jen.ID("Debug").ID("bool"), jen.ID("HTTPPort").ID("uint16")).Type().ID("FrontendSettings").Struct(jen.ID("StaticFilesDirectory").ID("string"), jen.ID("Debug").ID("bool"), jen.ID("CacheStaticFiles").ID("bool")).Type().ID("AuthSettings").Struct(jen.ID("CookieDomain").ID("string"), jen.ID("CookieSecret").ID("string"), jen.ID("CookieLifetime").Qual("time", "Duration"), jen.ID("Debug").ID("bool"), jen.ID("SecureCookiesOnly").ID("bool"), jen.ID("EnableUserSignup").ID("bool")).Type().ID("DatabaseSettings").Struct(jen.ID("Debug").ID("bool"), jen.ID("Provider").ID("string"), jen.ID("ConnectionDetails").ID("database").Dot(
		"ConnectionDetails",
	)).Type().ID("MetricsSettings").Struct(jen.ID("MetricsProvider").ID("metricsProvider"), jen.ID("TracingProvider").ID("tracingProvider"), jen.ID("DBMetricsCollectionInterval").Qual("time", "Duration"), jen.ID("RuntimeMetricsCollectionInterval").Qual("time", "Duration")).Type().ID("ServerConfig").Struct(jen.ID("Meta").ID("MetaSettings"), jen.ID("Frontend").ID("FrontendSettings"), jen.ID("Auth").ID("AuthSettings"), jen.ID("Server").ID("ServerSettings"), jen.ID("Database").ID("DatabaseSettings"), jen.ID("Metrics").ID("MetricsSettings")).Type().ID("MarshalFunc").Params(jen.ID("v").Interface()).Params(jen.Index().ID("byte"), jen.ID("error")),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// EncodeToFile renders your config to a file given your favorite encoder").Params(jen.ID("cfg").Op("*").ID("ServerConfig")).ID("EncodeToFile").Params(jen.ID("path").ID("string"), jen.ID("marshaler").ID("MarshalFunc")).Params(jen.ID("error")).Block(
		jen.List(jen.ID("byteSlice"), jen.ID("err")).Op(":=").ID("marshaler").Call(jen.Op("*").ID("cfg")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().ID("err"),
		),
		jen.Return().Qual("io/ioutil", "WriteFile").Call(jen.ID("path"), jen.ID("byteSlice"), jen.Lit(644)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// BuildConfig is a constructor function that initializes a viper config.").ID("BuildConfig").Params().Params(jen.Op("*").ID("viper").Dot(
		"Viper",
	)).Block(
		jen.ID("cfg").Op(":=").ID("viper").Dot(
			"New",
		).Call(),
		jen.ID("cfg").Dot(
			"SetDefault",
		).Call(jen.Lit("meta.startup_deadline"), jen.ID("defaultStartupDeadline")),
		jen.ID("cfg").Dot(
			"SetDefault",
		).Call(jen.Lit("auth.cookie_secret"), jen.ID("randString").Call()),
		jen.ID("cfg").Dot(
			"SetDefault",
		).Call(jen.Lit("auth.cookie_lifetime"), jen.ID("defaultCookieLifetime")),
		jen.ID("cfg").Dot(
			"SetDefault",
		).Call(jen.Lit("auth.enable_user_signup"), jen.ID("true")),
		jen.ID("cfg").Dot(
			"SetDefault",
		).Call(jen.Lit("metrics.database_metrics_collection_interval"), jen.ID("defaultMetricsCollectionInterval")),
		jen.ID("cfg").Dot(
			"SetDefault",
		).Call(jen.Lit("metrics.runtime_metrics_collection_interval"), jen.ID("defaultDatabaseMetricsCollectionInterval")),
		jen.ID("cfg").Dot(
			"SetDefault",
		).Call(jen.Lit("server.http_port"), jen.Lit(80)),
		jen.Return().ID("cfg"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ParseConfigFile parses a configuration file").ID("ParseConfigFile").Params(jen.ID("filename").ID("string")).Params(jen.Op("*").ID("ServerConfig"), jen.ID("error")).Block(
		jen.ID("cfg").Op(":=").ID("BuildConfig").Call(),
		jen.ID("cfg").Dot(
			"SetConfigFile",
		).Call(jen.ID("filename")),
		jen.If(jen.ID("err").Op(":=").ID("cfg").Dot(
			"ReadInConfig",
		).Call(), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("trying to read the config file: %w"), jen.ID("err"))),
		),
		jen.Null().Var().ID("serverConfig").Op("*").ID("ServerConfig"),
		jen.If(jen.ID("err").Op(":=").ID("cfg").Dot(
			"Unmarshal",
		).Call(jen.Op("&").ID("serverConfig")), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("trying to unmarshal the config: %w"), jen.ID("err"))),
		),
		jen.Return().List(jen.ID("serverConfig"), jen.ID("nil")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// randString produces a random string").Comment("// https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/").ID("randString").Params().Params(jen.ID("string")).Block(
		jen.ID("b").Op(":=").ID("make").Call(jen.Index().ID("byte"), jen.ID("randStringSize")),
		jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("panic").Call(jen.ID("err")),
		),
		jen.Return().Qual("encoding/base32", "StdEncoding").Dot(
			"WithPadding",
		).Call(jen.Qual("encoding/base32", "NoPadding")).Dot(
			"EncodeToString",
		).Call(jen.ID("b")),
	),

		jen.Line(),
	)
	return ret
}
