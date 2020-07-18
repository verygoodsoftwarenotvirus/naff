package httpserver

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serverDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("maxTimeout").Equals().Lit(120).Times().Qual("time", "Second"),
			jen.ID("serverNamespace").Equals().Lit("todo-service"),
		),
		jen.Line(),
	)

	makeServerStructDeclLines := func() []jen.Code {

		lines := []jen.Code{
			jen.ID("DebugMode").Bool(),
			jen.Line(),
			jen.Comment("Services."),
			jen.ID("authService").PointerTo().Qual(proj.ServiceV1AuthPackage(), "Service"),
			jen.ID("frontendService").PointerTo().Qual(proj.ServiceV1FrontendPackage(), "Service"),
			jen.ID("usersService").Qual(proj.ModelsV1Package(), "UserDataServer"),
			jen.ID("oauth2ClientsService").Qual(proj.ModelsV1Package(), "OAuth2ClientDataServer"),
			jen.ID("webhooksService").Qual(proj.ModelsV1Package(), "WebhookDataServer"),
		}

		for _, typ := range proj.DataTypes {
			tsn := typ.Name.Singular()
			tpuvn := typ.Name.PluralUnexportedVarName()
			lines = append(lines,
				jen.IDf("%sService", tpuvn).Qual(proj.ModelsV1Package(), fmt.Sprintf("%sDataServer", tsn)))
		}

		lines = append(lines,
			jen.Line(),
			jen.Comment("infra things."),
			jen.ID("db").Qual(proj.DatabaseV1Package(), "DataManager"),
			jen.ID("config").PointerTo().Qual(proj.InternalConfigV1Package(), "ServerConfig"),
			jen.ID("router").PointerTo().Qual("github.com/go-chi/chi", "Mux"),
			jen.ID("httpServer").PointerTo().Qual("net/http", "Server"),
			jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger"),
			jen.ID("encoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
		)

		// if proj.EnableNewsman {
		lines = append(lines,
			jen.ID("newsManager").PointerTo().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Newsman"),
		)
		// }

		return lines
	}

	ret.Add(
		jen.Type().Defs(
			jen.Comment("Server is our API httpServer."),
			jen.ID("Server").Struct(
				makeServerStructDeclLines()...,
			),
		),
		jen.Line(),
	)

	buildProvideServerParams := func() []jen.Code {
		lines := []jen.Code{
			constants.CtxParam(),
			jen.ID("cfg").PointerTo().Qual(proj.InternalConfigV1Package(), "ServerConfig"),
			jen.ID("authService").PointerTo().Qual(proj.ServiceV1AuthPackage(), "Service"),
			jen.ID("frontendService").PointerTo().Qual(proj.ServiceV1FrontendPackage(), "Service"),
		}

		for _, typ := range proj.DataTypes {
			lines = append(lines, jen.IDf("%sService", typ.Name.PluralUnexportedVarName()).Qual(proj.ModelsV1Package(), fmt.Sprintf("%sDataServer", typ.Name.Singular())))
		}

		lines = append(lines,
			jen.ID("usersService").Qual(proj.ModelsV1Package(), "UserDataServer"),
			jen.ID("oauth2Service").Qual(proj.ModelsV1Package(), "OAuth2ClientDataServer"),
			jen.ID("webhooksService").Qual(proj.ModelsV1Package(), "WebhookDataServer"),
			jen.ID("db").Qual(proj.DatabaseV1Package(), "DataManager"),
			jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger"),
			jen.ID("encoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
		)

		// if proj.EnableNewsman {
		lines = append(lines,
			jen.ID("newsManager").PointerTo().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Newsman"),
		)
		// }

		return lines
	}

	buildServerDecLines := func() []jen.Code {
		lines := []jen.Code{
			jen.ID("DebugMode").MapAssign().ID("cfg").Dot("Server").Dot("Debug"),
			jen.Comment("infra things"),
			jen.ID("db").MapAssign().ID("db"),
			jen.ID("config").MapAssign().ID("cfg"),
			jen.ID("encoder").MapAssign().ID("encoder"),
			jen.ID("httpServer").MapAssign().ID("provideHTTPServer").Call(),
			jen.ID(constants.LoggerVarName).MapAssign().ID(constants.LoggerVarName).Dot("WithName").Call(jen.Lit("api_server")),
		}

		// if proj.EnableNewsman {
		lines = append(lines,
			jen.ID("newsManager").MapAssign().ID("newsManager"),
		)
		// }

		lines = append(lines,
			jen.Comment("services"),
			jen.ID("webhooksService").MapAssign().ID("webhooksService"),
			jen.ID("frontendService").MapAssign().ID("frontendService"),
			jen.ID("usersService").MapAssign().ID("usersService"),
			jen.ID("authService").MapAssign().ID("authService"),
		)

		for _, typ := range proj.DataTypes {
			tpuvn := typ.Name.PluralUnexportedVarName()
			lines = append(lines, jen.IDf("%sService", tpuvn).MapAssign().IDf("%sService", tpuvn))
		}

		lines = append(lines, jen.ID("oauth2ClientsService").MapAssign().ID("oauth2Service"))

		return lines
	}

	buildProvideServerLines := func() []jen.Code {
		lines := []jen.Code{
			jen.If(jen.Len(jen.ID("cfg").Dot("Auth").Dot("CookieSecret")).LessThan().Lit(32)).Block(
				jen.Err().Assign().Qual("errors", "New").Call(jen.Lit("cookie secret is too short, must be at least 32 characters in length")),
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("cookie secret failure")),
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.ID("srv").Assign().AddressOf().ID("Server").Valuesln(
				buildServerDecLines()...,
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("cfg").Dot("ProvideTracing").Call(jen.ID(constants.LoggerVarName)), jen.Err().DoesNotEqual().ID("nil").And().Err().DoesNotEqual().Qual(proj.InternalConfigV1Package(), "ErrInvalidTracingProvider")).Block(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("ih"), jen.Err()).Assign().ID("cfg").Dot("ProvideInstrumentationHandler").Call(jen.ID(constants.LoggerVarName)),
			jen.If(jen.Err().DoesNotEqual().ID("nil").And().Err().DoesNotEqual().Qual(proj.InternalConfigV1Package(), "ErrInvalidMetricsProvider")).Block(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.If(jen.ID("ih").DoesNotEqual().ID("nil")).Block(
				jen.ID("srv").Dot("setupRouter").Call(jen.ID("cfg").Dot("Frontend"), jen.ID("ih")),
			),
			jen.Line(),
			jen.ID("srv").Dot("httpServer").Dot("Handler").Equals().AddressOf().Qual("go.opencensus.io/plugin/ochttp", "Handler").Valuesln(
				jen.ID("Handler").MapAssign().ID("srv").Dot("router"),
				jen.ID("FormatSpanName").MapAssign().ID("formatSpanNameForRequest"),
			),
			jen.Line(),
		}

		// if proj.EnableNewsman {
		lines = append(lines,
			jen.List(jen.ID("allWebhooks"), jen.Err()).Assign().ID("db").Dot("GetAllWebhooks").Call(constants.CtxVar()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("initializing webhooks: %w"), jen.Err())),
			),
			jen.Line(),
			jen.For(jen.ID("i").Assign().Zero(), jen.ID("i").LessThan().ID("len").Call(jen.ID("allWebhooks").Dot("Webhooks")), jen.ID("i").Op("++")).Block(
				jen.ID("wh").Assign().ID("allWebhooks").Dot("Webhooks").Index(jen.ID("i")),
				jen.Comment("NOTE: we must guarantee that whatever is stored in the database is valid, otherwise"),
				jen.Comment("newsman will try (and fail) to execute requests constantly"),
				jen.ID("l").Assign().ID("wh").Dot("ToListener").Call(jen.ID("srv").Dot(constants.LoggerVarName)),
				jen.ID("srv").Dot("newsManager").Dot("TuneIn").Call(jen.ID("l")),
			),
			jen.Line(),
		)
		// }

		lines = append(lines,
			jen.Return().List(jen.ID("srv"), jen.Nil()),
		)

		return lines
	}

	ret.Add(
		jen.Comment("ProvideServer builds a new Server instance."),
		jen.Line(),
		jen.Func().ID("ProvideServer").Paramsln(
			buildProvideServerParams()...,
		).Params(jen.PointerTo().ID("Server"), jen.Error()).Block(
			buildProvideServerLines()...,
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(`func (s *Server) logRoutes() {
	if err := chi.Walk(s.router, func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		s.logger.WithValues(map[string]interface{}{
			"method": method,
			"route":  route,
		}).Debug("route found")

		return nil
	}); err != nil {
		s.logger.Error(err, "logging routes")
	}
}
`),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Serve serves HTTP traffic."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Server")).ID("Serve").Params().Block(
			jen.ID("s").Dot("httpServer").Dot("Addr").Equals().Qual("fmt", "Sprintf").Call(jen.Lit(":%d"), jen.ID("s").Dot("config").Dot("Server").Dot("HTTPPort")),
			jen.ID("s").Dot(constants.LoggerVarName).Dot("Debug").Call(utils.FormatString("Listening for HTTP requests on %q", jen.ID("s").Dot("httpServer").Dot("Addr"))),
			jen.Line(),
			jen.Comment("returns ErrServerClosed on graceful close."),
			jen.If(jen.Err().Assign().ID("s").Dot("httpServer").Dot("ListenAndServe").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("server shutting down")),
				jen.If(jen.Err().IsEqualTo().Qual("net/http", "ErrServerClosed")).Block(
					jen.Comment("NOTE: there is a chance that next line won't have time to run,"),
					jen.Comment("as main() doesn't wait for this goroutine to stop."),
					jen.Qual("os", "Exit").Call(jen.Zero()),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("provideHTTPServer provides an HTTP httpServer."),
		jen.Line(),
		jen.Func().ID("provideHTTPServer").Params().Params(jen.PointerTo().Qual("net/http", "Server")).Block(
			jen.Comment("heavily inspired by https://blog.cloudflare.com/exposing-go-on-the-internet/"),
			jen.ID("srv").Assign().AddressOf().Qual("net/http", "Server").Valuesln(
				jen.ID("ReadTimeout").MapAssign().Lit(5).Times().Qual("time", "Second"),
				jen.ID("WriteTimeout").MapAssign().Lit(10).Times().Qual("time", "Second"),
				jen.ID("IdleTimeout").MapAssign().Lit(120).Times().Qual("time", "Second"),
				jen.ID("TLSConfig").MapAssign().AddressOf().Qual("crypto/tls", "Config").Valuesln(
					jen.ID("PreferServerCipherSuites").MapAssign().True(),
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
