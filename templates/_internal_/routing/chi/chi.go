package chi

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func chiDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("roughlyOneYear").Op("=").Qual("time", "Hour").Op("*").Lit(24).Op("*").Lit(365),
			jen.ID("maxTimeout").Op("=").Lit(120).Op("*").Qual("time", "Second"),
			jen.ID("maxCORSAge").Op("=").Lit(300),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("errInvalidMethod").Op("=").Qual("errors", "New").Call(jen.Lit("invalid method")),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").ID("routing").Dot("Router").Op("=").Parens(jen.Op("*").ID("router")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("router").Struct(
			jen.ID("router").ID("chi").Dot("Router"),
			jen.ID("tracer").ID("tracing").Dot("Tracer"),
			jen.ID("logger").ID("logging").Dot("Logger"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildChiMux").Params(jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("chi").Dot("Router")).Body(
			jen.ID("ch").Op(":=").ID("cors").Dot("New").Call(jen.ID("cors").Dot("Options").Valuesln(
				jen.ID("AllowedOrigins").Op(":").Index().ID("string").Valuesln(
					jen.Lit("*")), jen.ID("AllowedMethods").Op(":").Index().ID("string").Valuesln(
					jen.Qual("net/http", "MethodGet"), jen.Qual("net/http", "MethodPost"), jen.Qual("net/http", "MethodPut"), jen.Qual("net/http", "MethodDelete"), jen.Qual("net/http", "MethodOptions")), jen.ID("AllowedHeaders").Op(":").Index().ID("string").Valuesln(
					jen.Lit("Accept"), jen.Lit("Authorization"), jen.Lit("RawHTML-Provider"), jen.Lit("X-CSRF-Token")), jen.ID("ExposedHeaders").Op(":").Index().ID("string").Valuesln(
					jen.Lit("Link")), jen.ID("AllowCredentials").Op(":").ID("true"), jen.ID("MaxAge").Op(":").ID("maxCORSAge"))),
			jen.ID("sec").Op(":=").ID("secure").Dot("New").Call(jen.ID("secure").Dot("Options").Valuesln(
				jen.ID("AllowedHosts").Op(":").Index().ID("string").Valuesln(
					jen.Lit("")), jen.ID("AllowedHostsAreRegex").Op(":").ID("false"), jen.ID("HostsProxyHeaders").Op(":").Index().ID("string").Valuesln(
					jen.Lit("X-Forwarded-Hosts")), jen.ID("SSLRedirect").Op(":").ID("true"), jen.ID("SSLTemporaryRedirect").Op(":").ID("false"), jen.ID("SSLHost").Op(":").Lit(""), jen.ID("SSLHostFunc").Op(":").ID("nil"), jen.ID("SSLProxyHeaders").Op(":").Map(jen.ID("string")).ID("string").Valuesln(
					jen.Lit("X-Forwarded-Proto").Op(":").Lit("https")), jen.ID("STSSeconds").Op(":").ID("int64").Call(jen.ID("roughlyOneYear").Dot("Seconds").Call()), jen.ID("STSIncludeSubdomains").Op(":").ID("true"), jen.ID("STSPreload").Op(":").ID("true"), jen.ID("ForceSTSHeader").Op(":").ID("false"), jen.ID("FrameDeny").Op(":").ID("true"), jen.ID("CustomFrameOptionsValue").Op(":").Lit(""), jen.ID("ContentTypeNosniff").Op(":").ID("true"), jen.ID("BrowserXssFilter").Op(":").ID("true"), jen.ID("CustomBrowserXssValue").Op(":").Lit(""), jen.ID("ContentSecurityPolicy").Op(":").Lit(""), jen.ID("PublicKey").Op(":").Lit(""), jen.ID("ReferrerPolicy").Op(":").Lit(""), jen.ID("FeaturePolicy").Op(":").Lit(""), jen.ID("ExpectCTHeader").Op(":").Lit(""), jen.ID("SecureContextKey").Op(":").Lit("secureContext"), jen.ID("IsDevelopment").Op(":").ID("true"))),
			jen.ID("mux").Op(":=").ID("chi").Dot("NewRouter").Call(),
			jen.ID("mux").Dot("Use").Call(
				jen.ID("sec").Dot("Handler"),
				jen.Qual("github.com/go-chi/chi/middleware", "RequestID"),
				jen.Qual("github.com/go-chi/chi/middleware", "RealIP"),
				jen.Qual("github.com/go-chi/chi/middleware", "Timeout").Call(jen.ID("maxTimeout")),
				jen.ID("logging").Dot("BuildLoggingMiddleware").Call(jen.ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.Lit("router"))),
				jen.ID("ch").Dot("Handler"),
			),
			jen.Return().ID("mux"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildRouter").Params(jen.ID("mux").ID("chi").Dot("Router"), jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.Op("*").ID("router")).Body(
			jen.ID("logger").Op("=").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")),
			jen.If(jen.ID("mux").Op("==").ID("nil")).Body(
				jen.ID("logger").Dot("Info").Call(jen.Lit("starting with a new mux")),
				jen.ID("mux").Op("=").ID("buildChiMux").Call(jen.ID("logger")),
			),
			jen.ID("r").Op(":=").Op("&").ID("router").Valuesln(
				jen.ID("router").Op(":").ID("mux"), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.Lit("router")), jen.ID("logger").Op(":").ID("logger")),
			jen.Return().ID("r"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("convertMiddleware").Params(jen.ID("in").Op("...").ID("routing").Dot("Middleware")).Index().Func().Params(jen.ID("handler").Qual("net/http", "Handler")).Qual("net/http", "Handler").Body(
			jen.ID("out").Op(":=").Index().Func().Params(jen.ID("handler").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Values(),
			jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("in")).Body(
				jen.If(jen.ID("x").Op("!=").ID("nil")).Body(
					jen.ID("out").Op("=").ID("append").Call(
						jen.ID("out"),
						jen.ID("x"),
					))),
			jen.Return().ID("out"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewRouter constructs a new router."),
		jen.Line(),
		jen.Func().ID("NewRouter").Params(jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("routing").Dot("Router")).Body(
			jen.Return().ID("buildRouter").Call(
				jen.ID("nil"),
				jen.ID("logger"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("clone").Params().Params(jen.Op("*").ID("router")).Body(
			jen.Return().ID("buildRouter").Call(
				jen.ID("r").Dot("router"),
				jen.ID("r").Dot("logger"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("WithMiddleware returns a router with certain middleware applied."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("WithMiddleware").Params(jen.ID("middleware").Op("...").ID("routing").Dot("Middleware")).Params(jen.ID("routing").Dot("Router")).Body(
			jen.ID("x").Op(":=").ID("r").Dot("clone").Call(),
			jen.ID("x").Dot("router").Op("=").ID("x").Dot("router").Dot("With").Call(jen.ID("convertMiddleware").Call(jen.ID("middleware").Op("...")).Op("...")),
			jen.Return().ID("x"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogRoutes logs the described routes."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("LogRoutes").Params().Body(
			jen.If(jen.ID("err").Op(":=").ID("chi").Dot("Walk").Call(
				jen.ID("r").Dot("router"),
				jen.Func().Params(
					jen.ID("method").ID("string"),
					jen.ID("route").ID("string"),
					jen.ID("_").Qual("net/http", "Handler"),
					jen.ID("_").Op("...").Func().Params(jen.Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				).Params(jen.ID("error")).Body(
					jen.ID("r").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
						jen.Lit("method").Op(":").ID("method"),
						jen.Lit("route").Op(":").ID("route"),
					)).Dot("Debug").Call(jen.Lit("route found")),
					jen.Return().ID("nil"),
				),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("r").Dot("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("logging routes"),
				))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Route lets you apply a set of routes to a subrouter with a provided pattern."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("Route").Params(jen.ID("pattern").ID("string"), jen.ID("fn").Func().Params(jen.ID("r").ID("routing").Dot("Router"))).Params(jen.ID("routing").Dot("Router")).Body(
			jen.ID("r").Dot("router").Dot("Route").Call(
				jen.ID("pattern"),
				jen.Func().Params(jen.ID("subrouter").ID("chi").Dot("Router")).Body(
					jen.ID("fn").Call(jen.ID("buildRouter").Call(
						jen.ID("subrouter"),
						jen.ID("r").Dot("logger"),
					))),
			),
			jen.Return().ID("r"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AddRoute adds a route to the router."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("AddRoute").Params(jen.List(jen.ID("method"), jen.ID("path")).ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc"), jen.ID("middleware").Op("...").ID("routing").Dot("Middleware")).Params(jen.ID("error")).Body(
			jen.Switch(jen.Qual("strings", "TrimSpace").Call(jen.Qual("strings", "ToUpper").Call(jen.ID("method")))).Body(
				jen.Case(jen.Qual("net/http", "MethodGet")).Body(
					jen.ID("r").Dot("router").Dot("With").Call(jen.ID("convertMiddleware").Call(jen.ID("middleware").Op("...")).Op("...")).Dot("Get").Call(
						jen.ID("path"),
						jen.ID("handler"),
					)),
				jen.Case(jen.Qual("net/http", "MethodHead")).Body(
					jen.ID("r").Dot("router").Dot("With").Call(jen.ID("convertMiddleware").Call(jen.ID("middleware").Op("...")).Op("...")).Dot("Head").Call(
						jen.ID("path"),
						jen.ID("handler"),
					)),
				jen.Case(jen.Qual("net/http", "MethodPost")).Body(
					jen.ID("r").Dot("router").Dot("With").Call(jen.ID("convertMiddleware").Call(jen.ID("middleware").Op("...")).Op("...")).Dot("Post").Call(
						jen.ID("path"),
						jen.ID("handler"),
					)),
				jen.Case(jen.Qual("net/http", "MethodPut")).Body(
					jen.ID("r").Dot("router").Dot("With").Call(jen.ID("convertMiddleware").Call(jen.ID("middleware").Op("...")).Op("...")).Dot("Put").Call(
						jen.ID("path"),
						jen.ID("handler"),
					)),
				jen.Case(jen.Qual("net/http", "MethodPatch")).Body(
					jen.ID("r").Dot("router").Dot("With").Call(jen.ID("convertMiddleware").Call(jen.ID("middleware").Op("...")).Op("...")).Dot("Patch").Call(
						jen.ID("path"),
						jen.ID("handler"),
					)),
				jen.Case(jen.Qual("net/http", "MethodDelete")).Body(
					jen.ID("r").Dot("router").Dot("With").Call(jen.ID("convertMiddleware").Call(jen.ID("middleware").Op("...")).Op("...")).Dot("Delete").Call(
						jen.ID("path"),
						jen.ID("handler"),
					)),
				jen.Case(jen.Qual("net/http", "MethodConnect")).Body(
					jen.ID("r").Dot("router").Dot("With").Call(jen.ID("convertMiddleware").Call(jen.ID("middleware").Op("...")).Op("...")).Dot("Connect").Call(
						jen.ID("path"),
						jen.ID("handler"),
					)),
				jen.Case(jen.Qual("net/http", "MethodOptions")).Body(
					jen.ID("r").Dot("router").Dot("With").Call(jen.ID("convertMiddleware").Call(jen.ID("middleware").Op("...")).Op("...")).Dot("Options").Call(
						jen.ID("path"),
						jen.ID("handler"),
					)),
				jen.Case(jen.Qual("net/http", "MethodTrace")).Body(
					jen.ID("r").Dot("router").Dot("With").Call(jen.ID("convertMiddleware").Call(jen.ID("middleware").Op("...")).Op("...")).Dot("Trace").Call(
						jen.ID("path"),
						jen.ID("handler"),
					)),
				jen.Default().Body(
					jen.Return().Qual("fmt", "Errorf").Call(
						jen.Lit("%s: %w"),
						jen.ID("method"),
						jen.ID("errInvalidMethod"),
					)),
			),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Handler our interface by wrapping the underlying router's Handler method."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("Handler").Params().Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().ID("r").Dot("router")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Handle our interface by wrapping the underlying router's Handle method."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("Handle").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "Handler")).Body(
			jen.ID("r").Dot("router").Dot("Handle").Call(
				jen.ID("pattern"),
				jen.ID("handler"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("HandleFunc satisfies our interface by wrapping the underlying router's HandleFunc method."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("HandleFunc").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")).Body(
			jen.ID("r").Dot("router").Dot("HandleFunc").Call(
				jen.ID("pattern"),
				jen.ID("handler"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Connect satisfies our interface by wrapping the underlying router's Connect method."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("Connect").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")).Body(
			jen.ID("r").Dot("router").Dot("Connect").Call(
				jen.ID("pattern"),
				jen.ID("handler"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Delete satisfies our interface by wrapping the underlying router's Delete method."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("Delete").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")).Body(
			jen.ID("r").Dot("router").Dot("Delete").Call(
				jen.ID("pattern"),
				jen.ID("handler"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Get satisfies our interface by wrapping the underlying router's Get method."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("Get").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")).Body(
			jen.ID("r").Dot("router").Dot("Get").Call(
				jen.ID("pattern"),
				jen.ID("handler"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Head satisfies our interface by wrapping the underlying router's Head method."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("Head").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")).Body(
			jen.ID("r").Dot("router").Dot("Head").Call(
				jen.ID("pattern"),
				jen.ID("handler"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Options satisfies our interface by wrapping the underlying router's Options method."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("Options").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")).Body(
			jen.ID("r").Dot("router").Dot("Options").Call(
				jen.ID("pattern"),
				jen.ID("handler"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Patch satisfies our interface by wrapping the underlying router's Patch method."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("Patch").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")).Body(
			jen.ID("r").Dot("router").Dot("Patch").Call(
				jen.ID("pattern"),
				jen.ID("handler"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Post satisfies our interface by wrapping the underlying router's Post method."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("Post").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")).Body(
			jen.ID("r").Dot("router").Dot("Post").Call(
				jen.ID("pattern"),
				jen.ID("handler"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Put satisfies our interface by wrapping the underlying router's Put method."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("Put").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")).Body(
			jen.ID("r").Dot("router").Dot("Put").Call(
				jen.ID("pattern"),
				jen.ID("handler"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Trace satisfies our interface by wrapping the underlying router's Trace method."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("Trace").Params(jen.ID("pattern").ID("string"), jen.ID("handler").Qual("net/http", "HandlerFunc")).Body(
			jen.ID("r").Dot("router").Dot("Trace").Call(
				jen.ID("pattern"),
				jen.ID("handler"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildRouteParamIDFetcher builds a function that fetches a given key from a path with variables added by a router."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("BuildRouteParamIDFetcher").Params(
			jen.ID("logger").ID("logging").Dot("Logger"),
			jen.List(jen.ID("key"), jen.ID("logDescription")).ID("string"),
		).Params(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).ID("uint64")).Body(
			jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
				jen.ID("v").Op(":=").ID("chi").Dot("URLParam").Call(
					jen.ID("req"),
					jen.ID("key"),
				),
				jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(
					jen.ID("v"),
					jen.Lit(10),
					jen.Lit(64),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil").Op("&&").ID("len").Call(jen.ID("logDescription")).Op(">").Lit(0)).Body(
					jen.ID("logger").Dot("Error").Call(
						jen.ID("err"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("fetching %s ID from request"),
							jen.ID("logDescription"),
						),
					)),
				jen.Return().ID("u"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildRouteParamStringIDFetcher builds a function that fetches a given key from a path with variables added by a router."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").Op("*").ID("router")).ID("BuildRouteParamStringIDFetcher").Params(jen.ID("key").ID("string")).Func().Params(jen.Op("*").Qual("net/http", "Request")), jen.ID("string").Body(
			jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).ID("string").Body(
				jen.Return().ID("chi").Dot("URLParam").Call(jen.ID("req"), jen.ID("key")),
			),
		),
		jen.Line(),
	)

	return code
}
