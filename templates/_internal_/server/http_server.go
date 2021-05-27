package server

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpServerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("serverNamespace").Op("=").Lit("todo_service"),
			jen.ID("loggerName").Op("=").Lit("api_server"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("HTTPServer").Struct(
				jen.ID("authService").ID("types").Dot("AuthService"),
				jen.ID("accountsService").ID("types").Dot("AccountDataService"),
				jen.ID("frontendService").ID("frontend").Dot("Service"),
				jen.ID("auditService").ID("types").Dot("AuditLogEntryDataService"),
				jen.ID("usersService").ID("types").Dot("UserDataService"),
				jen.ID("adminService").ID("types").Dot("AdminService"),
				jen.ID("apiClientsService").ID("types").Dot("APIClientDataService"),
				jen.ID("webhooksService").ID("types").Dot("WebhookDataService"),
				jen.ID("itemsService").ID("types").Dot("ItemDataService"),
				jen.ID("encoder").ID("encoding").Dot("ServerEncoderDecoder"),
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("router").ID("routing").Dot("Router"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("httpServer").Op("*").Qual("net/http", "Server"),
				jen.ID("panicker").ID("panicking").Dot("Panicker"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideHTTPServer builds a new HTTPServer instance."),
		jen.Line(),
		jen.Func().ID("ProvideHTTPServer").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("serverSettings").ID("Config"), jen.ID("metricsHandler").ID("metrics").Dot("InstrumentationHandler"), jen.ID("authService").ID("types").Dot("AuthService"), jen.ID("auditService").ID("types").Dot("AuditLogEntryDataService"), jen.ID("usersService").ID("types").Dot("UserDataService"), jen.ID("accountsService").ID("types").Dot("AccountDataService"), jen.ID("apiClientsService").ID("types").Dot("APIClientDataService"), jen.ID("itemsService").ID("types").Dot("ItemDataService"), jen.ID("webhooksService").ID("types").Dot("WebhookDataService"), jen.ID("adminService").ID("types").Dot("AdminService"), jen.ID("frontendService").ID("frontend").Dot("Service"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("encoder").ID("encoding").Dot("ServerEncoderDecoder"), jen.ID("router").ID("routing").Dot("Router")).Params(jen.Op("*").ID("HTTPServer"), jen.ID("error")).Body(
			jen.ID("srv").Op(":=").Op("&").ID("HTTPServer").Valuesln(jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("loggerName")), jen.ID("encoder").Op(":").ID("encoder"), jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("loggerName")), jen.ID("panicker").Op(":").ID("panicking").Dot("NewProductionPanicker").Call(), jen.ID("httpServer").Op(":").ID("provideHTTPServer").Call(jen.ID("serverSettings").Dot("HTTPPort")), jen.ID("adminService").Op(":").ID("adminService"), jen.ID("auditService").Op(":").ID("auditService"), jen.ID("webhooksService").Op(":").ID("webhooksService"), jen.ID("frontendService").Op(":").ID("frontendService"), jen.ID("usersService").Op(":").ID("usersService"), jen.ID("accountsService").Op(":").ID("accountsService"), jen.ID("authService").Op(":").ID("authService"), jen.ID("itemsService").Op(":").ID("itemsService"), jen.ID("apiClientsService").Op(":").ID("apiClientsService")),
			jen.ID("srv").Dot("setupRouter").Call(
				jen.ID("ctx"),
				jen.ID("router"),
				jen.ID("metricsHandler"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("HTTP server successfully constructed")),
			jen.Return().List(jen.ID("srv"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Serve serves HTTP traffic."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("HTTPServer")).ID("Serve").Params().Body(
			jen.ID("s").Dot("logger").Dot("Debug").Call(jen.Lit("setting up server")),
			jen.ID("s").Dot("httpServer").Dot("Handler").Op("=").ID("otelhttp").Dot("NewHandler").Call(
				jen.ID("s").Dot("router").Dot("Handler").Call(),
				jen.ID("serverNamespace"),
				jen.ID("otelhttp").Dot("WithSpanNameFormatter").Call(jen.ID("tracing").Dot("FormatSpan")),
			),
			jen.ID("http2ServerConf").Op(":=").Op("&").ID("http2").Dot("Server").Values(),
			jen.If(jen.ID("err").Op(":=").ID("http2").Dot("ConfigureServer").Call(
				jen.ID("s").Dot("httpServer"),
				jen.ID("http2ServerConf"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("s").Dot("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("configuring HTTP2"),
				),
				jen.ID("s").Dot("panicker").Dot("Panic").Call(jen.ID("err")),
			),
			jen.ID("s").Dot("logger").Dot("WithValue").Call(
				jen.Lit("listening_on"),
				jen.ID("s").Dot("httpServer").Dot("Addr"),
			).Dot("Debug").Call(jen.Lit("Listening for HTTP requests")),
			jen.If(jen.ID("err").Op(":=").ID("s").Dot("httpServer").Dot("ListenAndServe").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("s").Dot("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("server shutting down"),
				),
				jen.If(jen.Qual("errors", "Is").Call(
					jen.ID("err"),
					jen.Qual("net/http", "ErrServerClosed"),
				)).Body(
					jen.Qual("os", "Exit").Call(jen.Lit(0))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("maxTimeout").Op("=").Lit(120).Op("*").Qual("time", "Second"),
			jen.ID("readTimeout").Op("=").Lit(5).Op("*").Qual("time", "Second"),
			jen.ID("writeTimeout").Op("=").Lit(2).Op("*").ID("readTimeout"),
			jen.ID("idleTimeout").Op("=").ID("maxTimeout"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("provideHTTPServer provides an HTTP httpServer."),
		jen.Line(),
		jen.Func().ID("provideHTTPServer").Params(jen.ID("port").ID("uint16")).Params(jen.Op("*").Qual("net/http", "Server")).Body(
			jen.ID("srv").Op(":=").Op("&").Qual("net/http", "Server").Valuesln(jen.ID("Addr").Op(":").Qual("fmt", "Sprintf").Call(
				jen.Lit(":%d"),
				jen.ID("port"),
			), jen.ID("ReadTimeout").Op(":").ID("readTimeout"), jen.ID("WriteTimeout").Op(":").ID("writeTimeout"), jen.ID("IdleTimeout").Op(":").ID("idleTimeout"), jen.ID("TLSConfig").Op(":").Op("&").Qual("crypto/tls", "Config").Valuesln(jen.ID("PreferServerCipherSuites").Op(":").ID("true"), jen.ID("CurvePreferences").Op(":").Index().Qual("crypto/tls", "CurveID").Valuesln(jen.Qual("crypto/tls", "CurveP256"), jen.Qual("crypto/tls", "X25519")), jen.ID("MinVersion").Op(":").Qual("crypto/tls", "VersionTLS12"), jen.ID("CipherSuites").Op(":").Index().ID("uint16").Valuesln(jen.Qual("crypto/tls", "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"), jen.Qual("crypto/tls", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"), jen.Qual("crypto/tls", "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305"), jen.Qual("crypto/tls", "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305"), jen.Qual("crypto/tls", "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"), jen.Qual("crypto/tls", "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256")))),
			jen.Return().ID("srv"),
		),
		jen.Line(),
	)

	return code
}
