package httpserver

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	buildProviderSet := func() []jen.Code {
		lines := []jen.Code{

			jen.ID("paramFetcherProviders"),
			jen.ID("ProvideServer"),
			jen.ID("ProvideNamespace"),
		}

		// if pkg.EnableNewsman {
		lines = append(lines, jen.ID("ProvideNewsmanTypeNameManipulationFunc"))
		// }

		return lines
	}

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is our wire superset of providers this package offers"),
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Callln(
				buildProviderSet()...,
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideNamespace provides a namespace"),
		jen.Line(),
		jen.Func().ID("ProvideNamespace").Params().Params(jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "Namespace")).Block(
			jen.Return().Lit("todo-service"),
		),
		jen.Line(),
	)

	// if pkg.EnableNewsman {
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
		jen.Func().ID("provideHTTPServer").Params().Params(jen.Op("*").Qual("net/http", "Server")).Block(
			jen.Comment("heavily inspired by https://blog.cloudflare.com/exposing-go-on-the-internet/"),
			jen.ID("srv").Op(":=").Op("&").Qual("net/http", "Server").Valuesln(
				jen.ID("ReadTimeout").Op(":").Lit(5).Op("*").Qual("time", "Second"),
				jen.ID("WriteTimeout").Op(":").Lit(10).Op("*").Qual("time", "Second"),
				jen.ID("IdleTimeout").Op(":").Lit(120).Op("*").Qual("time", "Second"),
				jen.ID("TLSConfig").Op(":").Op("&").Qual("crypto/tls", "Config").Valuesln(
					jen.ID("PreferServerCipherSuites").Op(":").ID("true"),
					jen.Comment(`"Only use curves which have assembly implementations"`).Line().
						ID("CurvePreferences").Op(":").Index().Qual("crypto/tls", "CurveID").Valuesln(
						jen.Qual("crypto/tls", "CurveP256"),
						jen.Qual("crypto/tls", "X25519"),
					),
					jen.ID("MinVersion").Op(":").Qual("crypto/tls", "VersionTLS12"),
					jen.ID("CipherSuites").Op(":").Index().ID("uint16").Valuesln(
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
