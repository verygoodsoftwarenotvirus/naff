package logging

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func loggingDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("LoggerNameKey").Op("=").Lit("_name_"),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("level").ID("int"),
			jen.ID("Level").Op("*").ID("level"),
			jen.ID("RequestIDFunc").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("string")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("InfoLevel").ID("Level").Op("=").ID("new").Call(jen.ID("level")),
			jen.ID("DebugLevel").ID("Level").Op("=").ID("new").Call(jen.ID("level")),
			jen.ID("ErrorLevel").ID("Level").Op("=").ID("new").Call(jen.ID("level")),
			jen.ID("WarnLevel").ID("Level").Op("=").ID("new").Call(jen.ID("level")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("Logger").Interface(
			jen.ID("Info").Params(jen.ID("string")),
			jen.ID("Debug").Params(jen.ID("string")),
			jen.ID("Error").Params(jen.ID("error"), jen.ID("string")),
			jen.ID("Fatal").Params(jen.ID("error")),
			jen.ID("Printf").Params(jen.ID("string"), jen.Op("...").Interface()),
			jen.ID("SetLevel").Params(jen.ID("Level")),
			jen.ID("SetRequestIDFunc").Params(jen.ID("RequestIDFunc")),
			jen.ID("Clone").Params().Params(jen.ID("Logger")),
			jen.ID("WithName").Params(jen.ID("string")).Params(jen.ID("Logger")),
			jen.ID("WithValues").Params(jen.Map(jen.ID("string")).Interface()).Params(jen.ID("Logger")),
			jen.ID("WithValue").Params(jen.ID("string"), jen.Interface()).Params(jen.ID("Logger")),
			jen.ID("WithRequest").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("Logger")),
			jen.ID("WithResponse").Params(jen.ID("response").Op("*").Qual("net/http", "Response")).Params(jen.ID("Logger")),
			jen.ID("WithError").Params(jen.ID("error")).Params(jen.ID("Logger")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EnsureLogger guarantees that a logger is available."),
		jen.Line(),
		jen.Func().ID("EnsureLogger").Params(jen.ID("logger").ID("Logger")).Params(jen.ID("Logger")).Body(
			jen.If(jen.ID("logger").Op("!=").ID("nil")).Body(
				jen.Return().ID("logger")),
			jen.Return().ID("NewNonOperationalLogger").Call(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("doNotLog").Op("=").Map(jen.ID("string")).Struct().Valuesln(
			jen.Lit("/metrics").Op(":").Values(), jen.Lit("/build/").Op(":").Values(), jen.Lit("/assets/").Op(":").Values()),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildLoggingMiddleware builds a logging middleware."),
		jen.Line(),
		jen.Func().ID("BuildLoggingMiddleware").Params(jen.ID("logger").ID("Logger")).Params(jen.Func().Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler"))).Body(
			jen.Return().Func().Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
				jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
					jen.ID("ww").Op(":=").Qual("github.com/go-chi/chi/middleware", "NewWrapResponseWriter").Call(
						jen.ID("res"),
						jen.ID("req").Dot("ProtoMajor"),
					),
					jen.ID("start").Op(":=").Qual("time", "Now").Call(),
					jen.ID("next").Dot("ServeHTTP").Call(
						jen.ID("ww"),
						jen.ID("req"),
					),
					jen.ID("shouldLog").Op(":=").ID("true"),
					jen.For(jen.ID("route").Op(":=").Range().ID("doNotLog")).Body(
						jen.If(jen.Qual("strings", "HasPrefix").Call(
							jen.ID("req").Dot("URL").Dot("Path"),
							jen.ID("route"),
						).Op("||").ID("req").Dot("URL").Dot("Path").Op("==").ID("route")).Body(
							jen.ID("shouldLog").Op("=").ID("false"),
							jen.Break(),
						)),
					jen.If(jen.ID("shouldLog")).Body(
						jen.ID("logger").Dot("WithRequest").Call(jen.ID("req")).Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
							jen.Lit("status").Op(":").ID("ww").Dot("Status").Call(),
							jen.Lit("elapsed").Op(":").Qual("time", "Since").Call(jen.ID("start")).Dot("Milliseconds").Call(),
							jen.Lit("written").Op(":").ID("ww").Dot("BytesWritten").Call(),
						)).Dot("Debug").Call(jen.Lit("response served")),
					),
				)),
			),
		),
		jen.Line(),
	)

	return code
}
