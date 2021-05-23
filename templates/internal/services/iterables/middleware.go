package iterables

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	code.Add(buildCreationInputMiddleware(proj, typ)...)
	code.Add(buildUpdateInputMiddleware(proj, typ)...)

	return code
}

func buildCreationInputMiddleware(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular() // singular name, PascalCased by default

	lines := []jen.Code{
		jen.Commentf("CreationInputMiddleware is a middleware for fetching, parsing, and attaching an %sInput struct from a request.", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
				jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
				jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
			).Body(
				jen.ID("x").Assign().ID("new").Call(jen.Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn))),
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("CreationInputMiddleware")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID(constants.RequestVarName), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered decoding request body")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				constants.CtxVar().Equals().Qual("context", "WithValue").Call(constants.CtxVar(), jen.ID("createMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName).Dot("WithContext").Call(constants.CtxVar())),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateInputMiddleware(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular() // singular name, PascalCased by default

	lines := []jen.Code{
		jen.Commentf("UpdateInputMiddleware is a middleware for fetching, parsing, and attaching an %sInput struct from a request.", sn),
		jen.Line(),
		jen.Comment("This is the same as the creation one, but that won't always be the case."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("UpdateInputMiddleware").Params(
			jen.ID("next").Qual("net/http", "Handler"),
		).Params(
			jen.Qual("net/http", "Handler"),
		).Body(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Body(
				jen.ID("x").Assign().ID("new").Call(jen.Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateInput", sn))),
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("UpdateInputMiddleware")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID(constants.RequestVarName), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered decoding request body")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				constants.CtxVar().Equals().Qual("context", "WithValue").Call(constants.CtxVar(), jen.ID("updateMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName).Dot("WithContext").Call(constants.CtxVar())),
			)),
		),
		jen.Line(),
	}

	return lines
}
