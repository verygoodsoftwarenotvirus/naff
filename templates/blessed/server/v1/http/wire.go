package httpserver

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(proj, ret)

	buildProviderSet := func() []jen.Code {
		lines := []jen.Code{
			jen.ID("paramFetcherProviders"),
			jen.ID("ProvideServer"),
			jen.ID("ProvideNamespace"),
		}

		// if proj.EnableNewsman {
		lines = append(lines, jen.ID("ProvideNewsmanTypeNameManipulationFunc"))
		// }

		return lines
	}

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is our wire superset of providers this package offers"),
			jen.ID("Providers").Equals().Qual("github.com/google/wire", "NewSet").Callln(
				buildProviderSet()...,
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideNamespace provides a namespace"),
		jen.Line(),
		jen.Func().ID("ProvideNamespace").Params().Params(jen.Qual(proj.InternalMetricsV1Package(), "Namespace")).Block(
			jen.Return().Lit("todo-service"),
		),
		jen.Line(),
	)

	// if proj.EnableNewsman {
	ret.Add(
		jen.Comment("ProvideNewsmanTypeNameManipulationFunc provides an WebhookIDFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideNewsmanTypeNameManipulationFunc").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "TypeNameManipulationFunc")).Block(
			jen.Return().Func().Params(jen.ID("s").ID("string")).Params(jen.ID("string")).Block(
				jen.ID("logger").Dot("WithName").Call(jen.Lit("events")).Dot("WithValue").Call(jen.Lit("type_name"), jen.ID("s")).Dot("Info").Call(jen.Lit("event occurred")),
				jen.Return().ID("s"),
			),
		),
		jen.Line(),
	)
	// }

	ret.Add(
		jen.Comment("provideHTTPServer provides an HTTP httpServer"),
		jen.Line(),
		jen.Func().ID("provideHTTPServer").Params().Params(jen.ParamPointer().Qual("net/http", "Server")).Block(
			jen.Comment("heavily inspired by https://blog.cloudflare.com/exposing-go-on-the-internet/"),
			jen.ID("srv").Assign().VarPointer().Qual("net/http", "Server").Valuesln(
				jen.ID("ReadTimeout").MapAssign().Lit(5).Times().Qual("time", "Second"),
				jen.ID("WriteTimeout").MapAssign().Lit(10).Times().Qual("time", "Second"),
				jen.ID("IdleTimeout").MapAssign().Lit(120).Times().Qual("time", "Second"),
				jen.ID("TLSConfig").MapAssign().VarPointer().Qual("crypto/tls", "Config").Valuesln(
					jen.ID("PreferServerCipherSuites").MapAssign().ID("true"),
					jen.Comment(`"Only use curves which have assembly implementations"`).Line().
						ID("CurvePreferences").MapAssign().Index().Qual("crypto/tls", "CurveID").Valuesln(
						jen.Qual("crypto/tls", "CurveP256"),
						jen.Qual("crypto/tls", "X25519"),
					),
					jen.ID("MinVersion").MapAssign().Qual("crypto/tls", "VersionTLS12"),
					jen.ID("CipherSuites").MapAssign().Index().ID("uint16").Valuesln(
						jen.Qual("crypto/tls", "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"),
						jen.Qual("crypto/tls", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"),
						jen.Qual("crypto/tls", "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305"),
						jen.Qual("crypto/tls", "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305"),
						jen.Qual("crypto/tls", "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"),
						jen.Qual("crypto/tls", "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"),
					),
				),
			),
			jen.Return().ID("srv"),
		),
		jen.Line(),
	)
	return ret
}
