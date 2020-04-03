package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("webhooks")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Comment("CreationInputMiddleware is a middleware for fetching, parsing, and attaching a parsed WebhookCreationInput struct from a request"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(
				jen.ID("x").Assign().ID("new").Call(jen.Qual(proj.ModelsV1Package(), "WebhookCreationInput")),
				jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("CreationInputMiddleware")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID("req"), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("error encountered decoding request body")),
					utils.WriteXHeader("res", "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				utils.CtxVar().Equals().Qual("context", "WithValue").Call(utils.CtxVar(), jen.ID("CreateMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req").Dot("WithContext").Call(utils.CtxVar())),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateInputMiddleware is a middleware for fetching, parsing, and attaching a parsed WebhookCreationInput struct from a request"),
		jen.Line(),
		jen.Comment("This is the same as the creation one, but it won't always be"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(
				jen.ID("x").Assign().ID("new").Call(jen.Qual(proj.ModelsV1Package(), "WebhookUpdateInput")),
				jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("UpdateInputMiddleware")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID("req"), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("error encountered decoding request body")),
					utils.WriteXHeader("res", "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				utils.CtxVar().Equals().Qual("context", "WithValue").Call(utils.CtxVar(), jen.ID("UpdateMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req").Dot("WithContext").Call(utils.CtxVar())),
			)),
		),
		jen.Line(),
	)
	return ret
}
