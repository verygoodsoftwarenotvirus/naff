package httpclient

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func pasetoDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("fetchAuthTokenForAPIClient").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("httpClient").Op("*").Qual("net/http", "Client"), jen.ID("clientID").ID("string"), jen.ID("secretKey").Index().ID("byte")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("clientID").Op("==").Lit("")).Body(
				jen.Return().List(jen.Lit(""), jen.ID("ErrEmptyInputProvided"))),
			jen.If(jen.ID("secretKey").Op("==").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.ID("ErrNilInputProvided"))),
			jen.If(jen.ID("httpClient").Op("==").ID("nil")).Body(
				jen.ID("httpClient").Op("=").Qual("net/http", "DefaultClient")),
			jen.If(jen.ID("httpClient").Dot("Timeout").Op("==").Lit(0)).Body(
				jen.ID("httpClient").Dot("Timeout").Op("=").ID("defaultTimeout")),
			jen.ID("input").Op(":=").Op("&").ID("types").Dot("PASETOCreationInput").Valuesln(jen.ID("ClientID").Op(":").ID("clientID"), jen.ID("RequestTime").Op(":").Qual("time", "Now").Call().Dot("UTC").Call().Dot("UnixNano").Call()),
			jen.If(jen.ID("c").Dot("accountID").Op("!=").Lit(0)).Body(
				jen.ID("input").Dot("AccountID").Op("=").ID("c").Dot("accountID")),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("APIClientClientIDKey"),
				jen.ID("clientID"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("fetching auth token")),
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating input"),
				))),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildAPIClientAuthTokenRequest").Call(
				jen.ID("ctx"),
				jen.ID("input"),
				jen.ID("secretKey"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building request"),
				))),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("c").Dot("fetchResponseToRequest").Call(
				jen.ID("ctx"),
				jen.ID("httpClient"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing request"),
				))),
			jen.If(jen.ID("err").Op("=").ID("errorFromResponse").Call(jen.ID("res")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("erroneous response"),
				))),
			jen.Var().ID("tokenRes").ID("types").Dot("PASETOResponse"),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("unmarshalBody").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.Op("&").ID("tokenRes"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("unmarshalling body"),
				))),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("auth token received")),
			jen.Return().List(jen.ID("tokenRes").Dot("Token"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
