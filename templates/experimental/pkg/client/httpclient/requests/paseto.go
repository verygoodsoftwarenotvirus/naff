package requests

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func pasetoDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("pasetoBasePath").Op("=").Lit("paseto"),
			jen.ID("signatureHeaderKey").Op("=").Lit("Signature"),
			jen.ID("validClientSecretSize").Op("=").Lit(128),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("setSignatureForRequest").Params(jen.ID("req").Op("*").Qual("net/http", "Request"), jen.List(jen.ID("body"), jen.ID("secretKey")).Index().ID("byte")).Params(jen.ID("error")).Body(
			jen.If(jen.ID("len").Call(jen.ID("secretKey")).Op("<").ID("validClientSecretSize")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("%w: %d"),
					jen.ID("ErrInvalidSecretKeyLength"),
					jen.ID("len").Call(jen.ID("secretKey")),
				)),
			jen.ID("mac").Op(":=").Qual("crypto/hmac", "New").Call(
				jen.Qual("crypto/sha256", "New"),
				jen.ID("secretKey"),
			),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("mac").Dot("Write").Call(jen.ID("body")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("writing hash content: %w"),
					jen.ID("err"),
				)),
			jen.ID("req").Dot("Header").Dot("Set").Call(
				jen.ID("signatureHeaderKey"),
				jen.Qual("encoding/base64", "RawURLEncoding").Dot("EncodeToString").Call(jen.ID("mac").Dot("Sum").Call(jen.ID("nil"))),
			),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildAPIClientAuthTokenRequest builds a request that fetches a PASETO from the service."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildAPIClientAuthTokenRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("PASETOCreationInput"), jen.ID("secretKey").Index().ID("byte")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("input").Op("==").ID("nil").Op("||").ID("len").Call(jen.ID("secretKey")).Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("uri").Op(":=").ID("b").Dot("buildUnversionedURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("pasetoBasePath"),
			),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("input").Dot("AccountID"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("APIClientClientIDKey"),
				jen.ID("input").Dot("ClientID"),
			),
			jen.ID("tracing").Dot("AttachRequestURIToSpan").Call(
				jen.ID("span"),
				jen.ID("uri"),
			),
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating input"),
				))),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("b").Dot("buildDataRequest").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building request"),
				))),
			jen.Var().Defs(
				jen.ID("buffer").Qual("bytes", "Buffer"),
			),
			jen.If(jen.ID("err").Op("=").ID("b").Dot("encoder").Dot("Encode").Call(
				jen.ID("ctx"),
				jen.Op("&").ID("buffer"),
				jen.ID("input"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("encoding body"),
				))),
			jen.If(jen.ID("err").Op("=").ID("setSignatureForRequest").Call(
				jen.ID("req"),
				jen.ID("buffer").Dot("Bytes").Call(),
				jen.ID("secretKey"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("signing request"),
				))),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("PASETO request built")),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
