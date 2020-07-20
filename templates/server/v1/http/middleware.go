package httpserver

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("httpserver")

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Defs(
			jen.ID("idReplacementRegex").Equals().Qual("regexp", "MustCompile").Call(jen.RawString(`[^(v|oauth)]\\d+`)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("formatSpanNameForRequest").Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.String()).Block(
			jen.Return().Qual("fmt", "Sprintf").Callln(
				jen.Lit("%s %s"),
				jen.ID(constants.RequestVarName).Dot("Method"),
				jen.ID("idReplacementRegex").Dot("ReplaceAllString").Call(jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"), jen.Lit("/{id}")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").PointerTo().ID("Server")).ID("loggingMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.ID("ww").Assign().Qual("github.com/go-chi/chi/middleware", "NewWrapResponseWriter").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName).Dot("ProtoMajor")),
				jen.Line(),
				jen.ID("start").Assign().Qual("time", "Now").Call(),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID("ww"), jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
					jen.Lit("status").MapAssign().ID("ww").Dot("Status").Call(),
					jen.Lit("bytes_written").MapAssign().ID("ww").Dot("BytesWritten").Call(),
					jen.Lit("elapsed").MapAssign().Qual("time", "Since").Call(jen.ID("start")),
				)).Dot("Debug").Call(jen.Lit("responded to request")),
			)),
		),
		jen.Line(),
	)

	return code
}
