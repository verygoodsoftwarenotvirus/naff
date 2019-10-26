package http

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func serverDotGo() *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("maxTimeout").Op("=").Lit(120).Op("*").Qual("time", "Second"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("Server is our API httpServer"),
			jen.ID("Server").Struct(
				jen.ID("DebugMode").ID("bool"),
				jen.Line(),
				jen.Comment("Services"),
				jen.ID("authService").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth", "Service"),
				jen.ID("frontendService").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/frontend", "Service"),
				jen.ID("usersService").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "UserDataServer"),
				jen.ID("oauth2ClientsService").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "OAuth2ClientDataServer"),
				jen.ID("webhooksService").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "WebhookDataServer"),
				jen.ID("itemsService").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "ItemDataServer"),
				jen.Line(),
				jen.Comment("infra things"),
				jen.ID("db").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1", "Database"),
				jen.ID("config").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/config", "ServerConfig"),
				jen.ID("router").Op("*").Qual("github.com/go-chi/chi", "Mux"),
				jen.ID("httpServer").Op("*").Qual("net/http", "Server"),
				jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
				jen.ID("encoder").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding", "EncoderDecoder"),
				jen.ID("newsManager").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Newsman"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideServer builds a new Server instance"),
		jen.Line(),
		jen.Func().ID("ProvideServer").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("cfg").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/config", "ServerConfig"),
			jen.ID("authService").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth", "Service"),
			jen.ID("frontendService").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/frontend", "Service"),
			jen.ID("itemsService").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "ItemDataServer"),
			jen.ID("usersService").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "UserDataServer"),
			jen.ID("oauth2Service").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "OAuth2ClientDataServer"),
			jen.ID("webhooksService").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "WebhookDataServer"),
			jen.ID("db").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1", "Database"),
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
			jen.ID("encoder").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding", "EncoderDecoder"),
			jen.ID("newsManager").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Newsman"),
		).Params(jen.Op("*").ID("Server"), jen.ID("error")).Block(
			jen.If(jen.ID("len").Call(jen.ID("cfg").Dot("Auth").Dot("CookieSecret")).Op("<").Lit(32)).Block(
				jen.ID("err").Op(":=").ID("errors").Dot("New").Call(jen.Lit("cookie secret is too short, must be at least 32 characters in length")),
				jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("cookie secret failure")),
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Line(),
			jen.ID("srv").Op(":=").Op("&").ID("Server").Valuesln(
				jen.ID("DebugMode").Op(":").ID("cfg").Dot("Server").Dot("Debug"),
				jen.Comment("infra things"),
				jen.ID("db").Op(":").ID("db"),
				jen.ID("config").Op(":").ID("cfg"),
				jen.ID("encoder").Op(":").ID("encoder"),
				jen.ID("httpServer").Op(":").ID("provideHTTPServer").Call(),
				jen.ID("logger").Op(":").ID("logger").Dot("WithName").Call(jen.Lit("api_server")),
				jen.ID("newsManager").Op(":").ID("newsManager"),
				jen.Comment("services"),
				jen.ID("webhooksService").Op(":").ID("webhooksService"),
				jen.ID("frontendService").Op(":").ID("frontendService"),
				jen.ID("usersService").Op(":").ID("usersService"),
				jen.ID("authService").Op(":").ID("authService"),
				jen.ID("itemsService").Op(":").ID("itemsService"),
				jen.ID("oauth2ClientsService").Op(":").ID("oauth2Service"),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("ProvideTracing").Call(jen.ID("logger")), jen.ID("err").Op("!=").ID("nil").Op("&&").ID("err").Op("!=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/config", "ErrInvalidTracingProvider")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Line(),
			jen.List(jen.ID("ih"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideInstrumentationHandler").Call(jen.ID("logger")),
			jen.If(jen.ID("err").Op("!=").ID("nil").Op("&&").ID("err").Op("!=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/config", "ErrInvalidMetricsProvider")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.If(jen.ID("ih").Op("!=").ID("nil")).Block(
				jen.ID("srv").Dot("setupRouter").Call(jen.ID("cfg").Dot("Frontend"), jen.ID("ih")),
			),
			jen.Line(),
			jen.ID("srv").Dot("httpServer").Dot("Handler").Op("=").Op("&").ID("ochttp").Dot("Handler").Valuesln(
				jen.ID("Handler").Op(":").ID("srv").Dot("router"),
				jen.ID("FormatSpanName").Op(":").ID("formatSpanNameForRequest"),
			),
			jen.Line(),
			jen.List(jen.ID("allWebhooks"), jen.ID("err")).Op(":=").ID("db").Dot("GetAllWebhooks").Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("initializing webhooks: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("len").Call(jen.ID("allWebhooks").Dot("Webhooks")), jen.ID("i").Op("++")).Block(
				jen.ID("wh").Op(":=").ID("allWebhooks").Dot("Webhooks").Index(jen.ID("i")),
				jen.Comment("NOTE: we must guarantee that whatever is stored in the database is valid, otherwise"),
				jen.Comment("newsman will try (and fail) to execute requests constantly"),
				jen.ID("l").Op(":=").ID("wh").Dot("ToListener").Call(jen.ID("srv").Dot("logger")),
				jen.ID("srv").Dot("newsManager").Dot("TuneIn").Call(jen.ID("l")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("srv"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Serve serves HTTP traffic"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Server")).ID("Serve").Params().Block(
			jen.ID("s").Dot("httpServer").Dot("Addr").Op("=").Qual("fmt", "Sprintf").Call(jen.Lit(":%d"), jen.ID("s").Dot("config").Dot("Server").Dot("HTTPPort")),
			jen.ID("s").Dot("logger").Dot("Debug").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("Listening for HTTP requests on %q"), jen.ID("s").Dot("httpServer").Dot("Addr"))),
			jen.Line(),
			jen.Comment("returns ErrServerClosed on graceful close"),
			jen.If(jen.ID("err").Op(":=").ID("s").Dot("httpServer").Dot("ListenAndServe").Call(), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("server shutting down")),
				jen.If(jen.ID("err").Op("==").Qual("net/http", "ErrServerClosed")).Block(
					jen.Comment("NOTE: there is a chance that next line won't have time to run,"),
					jen.Comment("as main() doesn't wait for this goroutine to stop."),
					jen.Qual("os", "Exit").Call(jen.Lit(0)),
				),
			),
		),
		jen.Line(),
	)
	return ret
}
