package logging

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func zerologLoggerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.ID("zerolog").Dot("CallerSkipFrameCount").Op("+=").Lit(2),
			jen.ID("zerolog").Dot("DisableSampling").Call(jen.ID("true")),
			jen.ID("zerolog").Dot("TimeFieldFormat").Op("=").ID("zerolog").Dot("TimeFormatUnixMs"),
			jen.ID("zerolog").Dot("TimestampFunc").Op("=").Func().Params().Params(jen.Qual("time", "Time")).Body(
				jen.Return().Qual("time", "Now").Call().Dot("UTC").Call()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("zerologLogger").Struct(
				jen.ID("requestIDFunc").ID("RequestIDFunc"),
				jen.ID("logger").ID("zerolog").Dot("Logger"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("buildZerologger builds a new zerologger."),
		jen.Line(),
		jen.Func().ID("buildZerologger").Params().Params(jen.ID("zerolog").Dot("Logger")).Body(
			jen.Return().ID("zerolog").Dot("New").Call(jen.Qual("os", "Stdout")).Dot("With").Call().Dot("Timestamp").Call().Dot("Logger").Call().Dot("Level").Call(jen.ID("zerolog").Dot("InfoLevel"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewZerologLogger builds a new zerologLogger."),
		jen.Line(),
		jen.Func().ID("NewZerologLogger").Params().Params(jen.ID("Logger")).Body(
			jen.Return().Op("&").ID("zerologLogger").Valuesln(jen.ID("logger").Op(":").ID("buildZerologger").Call())),
		jen.Line(),
	)

	code.Add(
		jen.Comment("WithName is our obligatory contract fulfillment function."),
		jen.Line(),
		jen.Comment("Zerolog doesn't support named loggers :( so we have this workaround."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("zerologLogger")).ID("WithName").Params(jen.ID("name").ID("string")).Params(jen.ID("Logger")).Body(
			jen.ID("l2").Op(":=").ID("l").Dot("logger").Dot("With").Call().Dot("Str").Call(
				jen.ID("LoggerNameKey"),
				jen.ID("name"),
			).Dot("Logger").Call(),
			jen.Return().Op("&").ID("zerologLogger").Valuesln(jen.ID("logger").Op(":").ID("l2")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SetLevel sets the log level for our zerologLogger."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("zerologLogger")).ID("SetLevel").Params(jen.ID("level").ID("Level")).Body(
			jen.Var().Defs(
				jen.ID("lvl").ID("zerolog").Dot("Level"),
			),
			jen.Switch(jen.ID("level")).Body(
				jen.Case(jen.ID("InfoLevel")).Body(
					jen.ID("lvl").Op("=").ID("zerolog").Dot("InfoLevel")),
				jen.Case(jen.ID("DebugLevel")).Body(
					jen.ID("l").Dot("logger").Op("=").ID("l").Dot("logger").Dot("With").Call().Dot("Logger").Call(), jen.ID("lvl").Op("=").ID("zerolog").Dot("DebugLevel")),
				jen.Case(jen.ID("WarnLevel")).Body(
					jen.ID("l").Dot("logger").Op("=").ID("l").Dot("logger").Dot("With").Call().Dot("Caller").Call().Dot("Logger").Call(), jen.ID("lvl").Op("=").ID("zerolog").Dot("WarnLevel")),
				jen.Case(jen.ID("ErrorLevel")).Body(
					jen.ID("l").Dot("logger").Op("=").ID("l").Dot("logger").Dot("With").Call().Dot("Caller").Call().Dot("Logger").Call(), jen.ID("lvl").Op("=").ID("zerolog").Dot("ErrorLevel")),
			),
			jen.ID("l").Dot("logger").Op("=").ID("l").Dot("logger").Dot("Level").Call(jen.ID("lvl")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SetRequestIDFunc sets the request ID retrieval function."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("zerologLogger")).ID("SetRequestIDFunc").Params(jen.ID("f").ID("RequestIDFunc")).Body(
			jen.If(jen.ID("f").Op("!=").ID("nil")).Body(
				jen.ID("l").Dot("requestIDFunc").Op("=").ID("f"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Info satisfies our contract for the logging.Logger Info method."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("zerologLogger")).ID("Info").Params(jen.ID("input").ID("string")).Body(
			jen.ID("l").Dot("logger").Dot("Info").Call().Dot("Msg").Call(jen.ID("input"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Debug satisfies our contract for the logging.Logger Debug method."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("zerologLogger")).ID("Debug").Params(jen.ID("input").ID("string")).Body(
			jen.ID("l").Dot("logger").Dot("Debug").Call().Dot("Msg").Call(jen.ID("input"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Error satisfies our contract for the logging.Logger Error method."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("zerologLogger")).ID("Error").Params(jen.ID("err").ID("error"), jen.ID("input").ID("string")).Body(
			jen.ID("l").Dot("logger").Dot("Error").Call().Dot("Stack").Call().Dot("Caller").Call().Dot("Err").Call(jen.ID("err")).Dot("Msg").Call(jen.ID("input"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Fatal satisfies our contract for the logging.Logger Fatal method."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("zerologLogger")).ID("Fatal").Params(jen.ID("err").ID("error")).Body(
			jen.ID("l").Dot("logger").Dot("Fatal").Call().Dot("Caller").Call().Dot("Err").Call(jen.ID("err")).Dot("Msg").Call(jen.ID("err").Dot("Error").Call())),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Printf satisfies our contract for the logging.Logger Printf method."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("zerologLogger")).ID("Printf").Params(jen.ID("format").ID("string"), jen.ID("args").Op("...").Interface()).Body(
			jen.ID("l").Dot("logger").Dot("Printf").Call(
				jen.ID("format"),
				jen.ID("args").Op("..."),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Clone satisfies our contract for the logging.Logger WithValue method."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("zerologLogger")).ID("Clone").Params().Params(jen.ID("Logger")).Body(
			jen.ID("l2").Op(":=").ID("l").Dot("logger").Dot("With").Call().Dot("Logger").Call(),
			jen.Return().Op("&").ID("zerologLogger").Valuesln(jen.ID("logger").Op(":").ID("l2")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("WithValue satisfies our contract for the logging.Logger WithValue method."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("zerologLogger")).ID("WithValue").Params(jen.ID("key").ID("string"), jen.ID("value").Interface()).Params(jen.ID("Logger")).Body(
			jen.ID("l2").Op(":=").ID("l").Dot("logger").Dot("With").Call().Dot("Interface").Call(
				jen.ID("key"),
				jen.ID("value"),
			).Dot("Logger").Call(),
			jen.Return().Op("&").ID("zerologLogger").Valuesln(jen.ID("logger").Op(":").ID("l2")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("WithValues satisfies our contract for the logging.Logger WithValues method."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("zerologLogger")).ID("WithValues").Params(jen.ID("values").Map(jen.ID("string")).Interface()).Params(jen.ID("Logger")).Body(
			jen.Var().Defs(
				jen.ID("l2").Op("=").ID("l").Dot("logger").Dot("With").Call().Dot("Logger").Call(),
			),
			jen.For(jen.List(jen.ID("key"), jen.ID("val")).Op(":=").Range().ID("values")).Body(
				jen.ID("l2").Op("=").ID("l2").Dot("With").Call().Dot("Interface").Call(
					jen.ID("key"),
					jen.ID("val"),
				).Dot("Logger").Call()),
			jen.Return().Op("&").ID("zerologLogger").Valuesln(jen.ID("logger").Op(":").ID("l2")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("WithError satisfies our contract for the logging.Logger WithError method."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("zerologLogger")).ID("WithError").Params(jen.ID("err").ID("error")).Params(jen.ID("Logger")).Body(
			jen.ID("l2").Op(":=").ID("l").Dot("logger").Dot("With").Call().Dot("Err").Call(jen.ID("err")).Dot("Logger").Call(),
			jen.Return().Op("&").ID("zerologLogger").Valuesln(jen.ID("logger").Op(":").ID("l2")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("l").Op("*").ID("zerologLogger")).ID("attachRequestToLog").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("zerolog").Dot("Logger")).Body(
			jen.If(jen.ID("req").Op("!=").ID("nil")).Body(
				jen.ID("l2").Op(":=").ID("l").Dot("logger").Dot("With").Call().Dot("Str").Call(
					jen.Lit("method"),
					jen.ID("req").Dot("Method"),
				).Dot("Logger").Call(),
				jen.If(jen.ID("req").Dot("URL").Op("!=").ID("nil")).Body(
					jen.ID("l2").Op("=").ID("l2").Dot("With").Call().Dot("Str").Call(
						jen.Lit("path"),
						jen.ID("req").Dot("URL").Dot("Path"),
					).Dot("Logger").Call(),
					jen.If(jen.ID("req").Dot("URL").Dot("RawQuery").Op("!=").Lit("")).Body(
						jen.ID("l2").Op("=").ID("l2").Dot("With").Call().Dot("Str").Call(
							jen.Lit("query"),
							jen.ID("req").Dot("URL").Dot("RawQuery"),
						).Dot("Logger").Call()),
				),
				jen.If(jen.ID("l").Dot("requestIDFunc").Op("!=").ID("nil")).Body(
					jen.If(jen.ID("reqID").Op(":=").ID("l").Dot("requestIDFunc").Call(jen.ID("req")), jen.ID("reqID").Op("!=").Lit("")).Body(
						jen.ID("l2").Op("=").ID("l2").Dot("With").Call().Dot("Str").Call(
							jen.Lit("request_id"),
							jen.ID("reqID"),
						).Dot("Logger").Call())),
				jen.Return().ID("l2"),
			),
			jen.Return().ID("l").Dot("logger"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("WithRequest satisfies our contract for the logging.Logger WithRequest method."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("zerologLogger")).ID("WithRequest").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("Logger")).Body(
			jen.Return().Op("&").ID("zerologLogger").Valuesln(jen.ID("logger").Op(":").ID("l").Dot("attachRequestToLog").Call(jen.ID("req")))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("WithResponse satisfies our contract for the logging.Logger WithResponse method."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("zerologLogger")).ID("WithResponse").Params(jen.ID("res").Op("*").Qual("net/http", "Response")).Params(jen.ID("Logger")).Body(
			jen.ID("l2").Op(":=").ID("l").Dot("logger").Dot("With").Call().Dot("Logger").Call(),
			jen.If(jen.ID("res").Op("!=").ID("nil")).Body(
				jen.ID("l2").Op("=").ID("l").Dot("attachRequestToLog").Call(jen.ID("res").Dot("Request")).Dot("With").Call().Dot("Int").Call(
					jen.ID("keys").Dot("ResponseStatusKey"),
					jen.ID("res").Dot("StatusCode"),
				).Dot("Logger").Call()),
			jen.Return().Op("&").ID("zerologLogger").Valuesln(jen.ID("logger").Op(":").ID("l2")),
		),
		jen.Line(),
	)

	return code
}
