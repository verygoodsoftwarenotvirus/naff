package iterables

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareDotGo(typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(ret)
	sn := typ.Name.Singular() // singular name, PascalCased by default

	ret.Add(
		jen.Comment(fmt.Sprintf("CreationInputMiddleware is a middleware for fetching, parsing, and attaching an %sInput struct from a request", sn)),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.ID("x").Op(":=").ID("new").Call(jen.ID("models").Dot(fmt.Sprintf("%sCreationInput", sn))),
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("CreationInputMiddleware")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.If(jen.ID("err").Op(":=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID("req"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered decoding request body")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
					jen.Return(),
				),
				jen.Line(),
				jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("CreateMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req").Dot("WithContext").Call(jen.ID("ctx"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("UpdateInputMiddleware is a middleware for fetching, parsing, and attaching an %sInput struct from a request", sn)),
		jen.Line(),
		jen.Comment("This is the same as the creation one, but that won't always be the case."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.ID("x").Op(":=").ID("new").Call(jen.ID("models").Dot(fmt.Sprintf("%sUpdateInput", sn))),
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("UpdateInputMiddleware")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.If(jen.ID("err").Op(":=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID("req"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered decoding request body")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
					jen.Return(),
				),
				jen.Line(),
				jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req").Dot("WithContext").Call(jen.ID("ctx"))),
			)),
		),
		jen.Line(),
	)
	return ret
}