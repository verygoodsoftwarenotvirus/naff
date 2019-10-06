package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func configDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("defaultStartupDeadline").Op("=").Qual("time", "Minute").Var().ID("defaultCookieLifetime").Op("=").Lit(24).Op("*").Qual("time", "Hour").Var().ID("defaultMetricsCollectionInterval").Op("=").Lit(2).Op("*").Qual("time", "Second").Var().ID("defaultDatabaseMetricsCollectionInterval").Op("=").Lit(2).Op("*").Qual("time", "Second").Var().ID("randStringSize").Op("=").Lit(32))
	ret.Add(jen.Func().ID("init").Params().Block(
		jen.ID("b").Op(":=").ID("make").Call(jen.Index().ID("byte"), jen.Lit(64)),
		jen.If(
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").Qual("crypto/rand", "Read").Call(jen.ID("b")),
			jen.ID("err").Op("!=").ID("nil"),
		).Block(
			jen.ID("panic").Call(jen.ID("err")),
		),
	),
	)
	ret.Add(jen.Null().Type().ID("MetaSettings").Struct(
		jen.ID("Debug").ID("bool"),
		jen.ID("StartupDeadline").Qual("time", "Duration"),
	).Type().ID("ServerSettings").Struct(
		jen.ID("Debug").ID("bool"),
		jen.ID("HTTPPort").ID("uint16"),
	).Type().ID("FrontendSettings").Struct(
		jen.ID("StaticFilesDirectory").ID("string"),
		jen.ID("Debug").ID("bool"),
		jen.ID("CacheStaticFiles").ID("bool"),
	).Type().ID("AuthSettings").Struct(
		jen.ID("CookieDomain").ID("string"),
		jen.ID("CookieSecret").ID("string"),
		jen.ID("CookieLifetime").Qual("time", "Duration"),
		jen.ID("Debug").ID("bool"),
		jen.ID("SecureCookiesOnly").ID("bool"),
		jen.ID("EnableUserSignup").ID("bool"),
	).Type().ID("DatabaseSettings").Struct(
		jen.ID("Debug").ID("bool"),
		jen.ID("Provider").ID("string"),
		jen.ID("ConnectionDetails").ID("database").Dot(
			"ConnectionDetails",
		),
	).Type().ID("MetricsSettings").Struct(
		jen.ID("MetricsProvider").ID("metricsProvider"),
		jen.ID("TracingProvider").ID("tracingProvider"),
		jen.ID("DBMetricsCollectionInterval").Qual("time", "Duration"),
		jen.ID("RuntimeMetricsCollectionInterval").Qual("time", "Duration"),
	).Type().ID("ServerConfig").Struct(
		jen.ID("Meta").ID("MetaSettings"),
		jen.ID("Frontend").ID("FrontendSettings"),
		jen.ID("Auth").ID("AuthSettings"),
		jen.ID("Server").ID("ServerSettings"),
		jen.ID("Database").ID("DatabaseSettings"),
		jen.ID("Metrics").ID("MetricsSettings"),
	).Type().ID("MarshalFunc").Params(jen.ID("v").Interface()).Params(jen.Index().ID("byte"), jen.ID("error")),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
