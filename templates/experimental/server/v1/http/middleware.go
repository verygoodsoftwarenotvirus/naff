package http

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func middlewareDotGo() *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("idReplacementRegex").Op("=").Qual("regexp", "MustCompile").Call(jen.Lit(`[^(v|oauth)]\d+`)),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("formatSpanNameForRequest").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Block(
		jen.Return().Qual("fmt", "Sprintf").Call(jen.Lit("%s %s"), jen.ID("req").Dot(
			"Method",
		), jen.ID("idReplacementRegex").Dot(
			"ReplaceAllString",
		).Call(jen.ID("req").Dot(
			"URL",
		).Dot(
			"Path",
		), jen.Lit(`/{id}`))),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("Server")).ID("loggingMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
		jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
			jen.ID("ww").Op(":=").ID("middleware").Dot(
				"NewWrapResponseWriter",
			).Call(jen.ID("res"), jen.ID("req").Dot(
				"ProtoMajor",
			)),
			jen.ID("start").Op(":=").Qual("time", "Now").Call(),
			jen.ID("next").Dot(
				"ServeHTTP",
			).Call(jen.ID("ww"), jen.ID("req")),
			jen.ID("s").Dot(
				"logger",
			).Dot(
				"WithRequest",
			).Call(jen.ID("req")).Dot(
				"WithValues",
			).Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("status").Op(":").ID("ww").Dot(
				"Status",
			).Call(), jen.Lit("bytes_written").Op(":").ID("ww").Dot(
				"BytesWritten",
			).Call(), jen.Lit("elapsed").Op(":").Qual("time", "Since").Call(jen.ID("start")))).Dot(
				"Debug",
			).Call(jen.Lit("request received")),
		)),
	),
	jen.Line(),
	)
	return ret
}
