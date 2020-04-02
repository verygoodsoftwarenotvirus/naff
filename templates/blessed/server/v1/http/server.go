package httpserver

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serverDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("maxTimeout").Equals().Lit(120).Times().Qual("time", "Second"),
		),
		jen.Line(),
	)

	makeServerStructDeclLines := func() []jen.Code {

		lines := []jen.Code{
			jen.ID("DebugMode").ID("bool"),
			jen.Line(),
			jen.Comment("Services"),
			jen.ID("authService").Op("*").Qual(pkg.ServiceV1AuthPackage(), "Service"),
			jen.ID("frontendService").Op("*").Qual(pkg.ServiceV1FrontendPackage(), "Service"),
			jen.ID("usersService").Qual(pkg.ModelsV1Package(), "UserDataServer"),
			jen.ID("oauth2ClientsService").Qual(pkg.ModelsV1Package(), "OAuth2ClientDataServer"),
			jen.ID("webhooksService").Qual(pkg.ModelsV1Package(), "WebhookDataServer"),
		}

		for _, typ := range pkg.DataTypes {
			tsn := typ.Name.Singular()
			tpuvn := typ.Name.PluralUnexportedVarName()
			lines = append(lines,
				jen.IDf("%sService", tpuvn).Qual(pkg.ModelsV1Package(), fmt.Sprintf("%sDataServer", tsn)))
		}

		lines = append(lines,
			jen.Line(),
			jen.Comment("infra things"),
			jen.ID("db").Qual(pkg.DatabaseV1Package(), "Database"),
			jen.ID("config").Op("*").Qual(pkg.InternalConfigV1Package(), "ServerConfig"),
			jen.ID("router").ParamPointer().Qual("github.com/go-chi/chi", "Mux"),
			jen.ID("httpServer").ParamPointer().Qual("net/http", "Server"),
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
			jen.ID("encoder").Qual(pkg.InternalEncodingV1Package(), "EncoderDecoder"),
		)

		// if pkg.EnableNewsman {
		lines = append(lines,
			jen.ID("newsManager").ParamPointer().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Newsman"),
		)
		// }

		return lines
	}

	ret.Add(
		jen.Type().Defs(
			jen.Comment("Server is our API httpServer"),
			jen.ID("Server").Struct(
				makeServerStructDeclLines()...,
			),
		),
		jen.Line(),
	)

	buildProvideServerParams := func() []jen.Code {
		lines := []jen.Code{
			utils.CtxParam(),
			jen.ID("cfg").Op("*").Qual(pkg.InternalConfigV1Package(), "ServerConfig"),
			jen.ID("authService").Op("*").Qual(pkg.ServiceV1AuthPackage(), "Service"),
			jen.ID("frontendService").Op("*").Qual(pkg.ServiceV1FrontendPackage(), "Service"),
		}

		for _, typ := range pkg.DataTypes {
			lines = append(lines, jen.IDf("%sService", typ.Name.PluralUnexportedVarName()).Qual(pkg.ModelsV1Package(), fmt.Sprintf("%sDataServer", typ.Name.Singular())))
		}

		lines = append(lines,
			jen.ID("usersService").Qual(pkg.ModelsV1Package(), "UserDataServer"),
			jen.ID("oauth2Service").Qual(pkg.ModelsV1Package(), "OAuth2ClientDataServer"),
			jen.ID("webhooksService").Qual(pkg.ModelsV1Package(), "WebhookDataServer"),
			jen.ID("db").Qual(pkg.DatabaseV1Package(), "Database"),
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
			jen.ID("encoder").Qual(pkg.InternalEncodingV1Package(), "EncoderDecoder"),
		)

		// if pkg.EnableNewsman {
		lines = append(lines,
			jen.ID("newsManager").ParamPointer().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Newsman"),
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
			jen.ID("logger").MapAssign().ID("logger").Dot("WithName").Call(jen.Lit("api_server")),
		}

		// if pkg.EnableNewsman {
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

		for _, typ := range pkg.DataTypes {
			tpuvn := typ.Name.PluralUnexportedVarName()
			lines = append(lines, jen.IDf("%sService", tpuvn).MapAssign().IDf("%sService", tpuvn))
		}

		lines = append(lines, jen.ID("oauth2ClientsService").MapAssign().ID("oauth2Service"))

		return lines
	}

	buildProvideServerLines := func() []jen.Code {
		lines := []jen.Code{
			jen.If(jen.ID("len").Call(jen.ID("cfg").Dot("Auth").Dot("CookieSecret")).Op("<").Lit(32)).Block(
				jen.Err().Assign().ID("errors").Dot("New").Call(jen.Lit("cookie secret is too short, must be at least 32 characters in length")),
				jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("cookie secret failure")),
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.ID("srv").Assign().VarPointer().ID("Server").Valuesln(
				buildServerDecLines()...,
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("cfg").Dot("ProvideTracing").Call(jen.ID("logger")), jen.Err().DoesNotEqual().ID("nil").Op("&&").ID("err").DoesNotEqual().Qual(pkg.InternalConfigV1Package(), "ErrInvalidTracingProvider")).Block(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("ih"), jen.Err()).Assign().ID("cfg").Dot("ProvideInstrumentationHandler").Call(jen.ID("logger")),
			jen.If(jen.Err().DoesNotEqual().ID("nil").Op("&&").ID("err").DoesNotEqual().Qual(pkg.InternalConfigV1Package(), "ErrInvalidMetricsProvider")).Block(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.If(jen.ID("ih").DoesNotEqual().ID("nil")).Block(
				jen.ID("srv").Dot("setupRouter").Call(jen.ID("cfg").Dot("Frontend"), jen.ID("ih")),
			),
			jen.Line(),
			jen.ID("srv").Dot("httpServer").Dot("Handler").Equals().VarPointer().Qual("go.opencensus.io/plugin/ochttp", "Handler").Valuesln(
				jen.ID("Handler").MapAssign().ID("srv").Dot("router"),
				jen.ID("FormatSpanName").MapAssign().ID("formatSpanNameForRequest"),
			),
			jen.Line(),
		}

		// if pkg.EnableNewsman {
		lines = append(lines,
			jen.List(jen.ID("allWebhooks"), jen.Err()).Assign().ID("db").Dot("GetAllWebhooks").Call(utils.CtxVar()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("initializing webhooks: %w"), jen.Err())),
			),
			jen.Line(),
			jen.For(jen.ID("i").Assign().Lit(0), jen.ID("i").Op("<").ID("len").Call(jen.ID("allWebhooks").Dot("Webhooks")), jen.ID("i").Op("++")).Block(
				jen.ID("wh").Assign().ID("allWebhooks").Dot("Webhooks").Index(jen.ID("i")),
				jen.Comment("NOTE: we must guarantee that whatever is stored in the database is valid, otherwise"),
				jen.Comment("newsman will try (and fail) to execute requests constantly"),
				jen.ID("l").Assign().ID("wh").Dot("ToListener").Call(jen.ID("srv").Dot("logger")),
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
		jen.Comment("ProvideServer builds a new Server instance"),
		jen.Line(),
		jen.Func().ID("ProvideServer").Paramsln(
			buildProvideServerParams()...,
		).Params(jen.Op("*").ID("Server"), jen.ID("error")).Block(
			buildProvideServerLines()...,
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(`
func (s *Server) logRoutes() {
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
		jen.Comment("Serve serves HTTP traffic"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Server")).ID("Serve").Params().Block(
			jen.ID("s").Dot("httpServer").Dot("Addr").Equals().Qual("fmt", "Sprintf").Call(jen.Lit(":%d"), jen.ID("s").Dot("config").Dot("Server").Dot("HTTPPort")),
			jen.ID("s").Dot("logger").Dot("Debug").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("Listening for HTTP requests on %q"), jen.ID("s").Dot("httpServer").Dot("Addr"))),
			jen.Line(),
			jen.Comment("returns ErrServerClosed on graceful close"),
			jen.If(jen.Err().Assign().ID("s").Dot("httpServer").Dot("ListenAndServe").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("server shutting down")),
				jen.If(jen.Err().Op("==").Qual("net/http", "ErrServerClosed")).Block(
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
