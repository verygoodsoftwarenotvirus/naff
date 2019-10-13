package http

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func wireDotGo() *jen.File {
	ret := jen.NewFile("httpserver")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("paramFetcherProviders"), jen.ID("ProvideServer"), jen.ID("ProvideNamespace"), jen.ID("ProvideNewsmanTypeNameManipulationFunc")),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideNamespace provides a namespace").ID("ProvideNamespace").Params().Params(jen.ID("metrics").Dot(
		"Namespace",
	)).Block(
		jen.Return().Lit("todo-service"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideNewsmanTypeNameManipulationFunc provides an WebhookIDFetcher").ID("ProvideNewsmanTypeNameManipulationFunc").Params(jen.ID("logger").ID("logging").Dot(
		"Logger",
	)).Params(jen.ID("newsman").Dot(
		"TypeNameManipulationFunc",
	)).Block(
		jen.Return().Func().Params(jen.ID("s").ID("string")).Params(jen.ID("string")).Block(
			jen.ID("logger").Dot(
				"WithName",
			).Call(jen.Lit("events")).Dot(
				"WithValue",
			).Call(jen.Lit("type_name"), jen.ID("s")).Dot(
				"Info",
			).Call(jen.Lit("event occurred")),
			jen.Return().ID("s"),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// provideHTTPServer provides an HTTP httpServer").ID("provideHTTPServer").Params().Params(jen.Op("*").Qual("net/http", "Server")).Block(
		jen.ID("srv").Op(":=").Op("&").Qual("net/http", "Server").Valuesln(jen.ID("ReadTimeout").Op(":").Lit(5).Op("*").Qual("time", "Second"), jen.ID("WriteTimeout").Op(":").Lit(10).Op("*").Qual("time", "Second"), jen.ID("IdleTimeout").Op(":").Lit(120).Op("*").Qual("time", "Second"), jen.ID("TLSConfig").Op(":").Op("&").Qual("crypto/tls", "Config").Valuesln(jen.ID("PreferServerCipherSuites").Op(":").ID("true"), jen.ID("CurvePreferences").Op(":").Index().Qual("crypto/tls", "CurveID").Valuesln(jen.Qual("crypto/tls", "CurveP256"), jen.Qual("crypto/tls", "X25519")), jen.ID("MinVersion").Op(":").Qual("crypto/tls", "VersionTLS12"), jen.ID("CipherSuites").Op(":").Index().ID("uint16").Valuesln(jen.Qual("crypto/tls", "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"), jen.Qual("crypto/tls", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"), jen.Qual("crypto/tls", "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305"), jen.Qual("crypto/tls", "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305"), jen.Qual("crypto/tls", "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"), jen.Qual("crypto/tls", "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256")))),
		jen.Return().ID("srv"),
	),

		jen.Line(),
	)
	return ret
}
