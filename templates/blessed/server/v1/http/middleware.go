package httpserver

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Defs(
			jen.ID("idReplacementRegex").Equals().Qual("regexp", "MustCompile").Call(jen.RawString(`[^(v|oauth)]\\d+`)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("formatSpanNameForRequest").Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.ID("string")).Block(
			jen.Return().Qual("fmt", "Sprintf").Callln(
				jen.Lit("%s %s"),
				jen.ID("req").Dot("Method"),
				jen.ID("idReplacementRegex").Dot("ReplaceAllString").Call(jen.ID("req").Dot("URL").Dot("Path"), jen.Lit("/{id}")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("Server")).ID("loggingMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(
				jen.ID("ww").Assign().Qual("github.com/go-chi/chi/middleware", "NewWrapResponseWriter").Call(jen.ID("res"), jen.ID("req").Dot("ProtoMajor")),
				jen.Line(),
				jen.ID("start").Assign().Qual("time", "Now").Call(),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID("ww"), jen.ID("req")),
				jen.Line(),
				jen.ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")).Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("status").MapAssign().ID("ww").Dot("Status").Call(),
					jen.Lit("bytes_written").MapAssign().ID("ww").Dot("BytesWritten").Call(),
					jen.Lit("elapsed").MapAssign().Qual("time", "Since").Call(jen.ID("start")),
				)).Dot("Debug").Call(jen.Lit("responded to request")),
			)),
		),
		jen.Line(),
	)
	return ret
}
