package iterables

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, ret)
	sn := typ.Name.Singular() // singular name, PascalCased by default

	ret.Add(
		jen.Commentf("CreationInputMiddleware is a middleware for fetching, parsing, and attaching an %sInput struct from a request", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.ID("x").Assign().ID("new").Call(jen.Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", sn))),
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("CreationInputMiddleware")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID(constants.RequestVarName), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered decoding request body")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				constants.CtxVar().Equals().Qual("context", "WithValue").Call(constants.CtxVar(), jen.ID("CreateMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName).Dot("WithContext").Call(constants.CtxVar())),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("UpdateInputMiddleware is a middleware for fetching, parsing, and attaching an %sInput struct from a request", sn),
		jen.Line(),
		jen.Comment("This is the same as the creation one, but that won't always be the case."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.ID("x").Assign().ID("new").Call(jen.Qual(proj.ModelsV1Package(), fmt.Sprintf("%sUpdateInput", sn))),
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("UpdateInputMiddleware")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID(constants.RequestVarName), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered decoding request body")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				constants.CtxVar().Equals().Qual("context", "WithValue").Call(constants.CtxVar(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName).Dot("WithContext").Call(constants.CtxVar())),
			)),
		),
		jen.Line(),
	)

	return ret
}
