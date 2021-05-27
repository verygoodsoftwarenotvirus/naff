package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func pasetoTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("signatureHeaderKey").Op("=").Lit("Signature"),
			jen.ID("validClientSecretSize").Op("=").Lit(128),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_fetchAuthTokenForAPIClient").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("anticipatedResult").Op(":=").Lit("v2.local.QAxIpVe-ECVNI1z4xQbm_qQYomyT3h8FtV8bxkz8pBJWkT8f7HtlOpbroPDEZUKop_vaglyp76CzYy375cHmKCW8e1CCkV0Lflu4GTDyXMqQdpZMM1E6OaoQW27gaRSvWBrR3IgbFIa0AkuUFw.UGFyYWdvbiBJbml0aWF0aXZlIEVudGVycHJpc2Vz"),
					jen.ID("ts").Op(":=").ID("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
						jen.ID("response").Op(":=").Op("&").ID("types").Dot("PASETOResponse").Valuesln(jen.ID("Token").Op(":").ID("anticipatedResult")),
						jen.ID("assert").Dot("NotEmpty").Call(
							jen.ID("t"),
							jen.ID("req").Dot("Header").Dot("Get").Call(jen.ID("signatureHeaderKey")),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("response")),
						),
					))),
					jen.ID("c").Op(":=").ID("buildTestClient").Call(
						jen.ID("t"),
						jen.ID("ts"),
					),
					jen.ID("exampleClientID").Op(":=").Lit("example_client_id"),
					jen.ID("exampleSecret").Op(":=").ID("make").Call(
						jen.Index().ID("byte"),
						jen.ID("validClientSecretSize"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("c").Dot("fetchAuthTokenForAPIClient").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("unauthenticatedClient"),
						jen.ID("exampleClientID"),
						jen.ID("exampleSecret"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("anticipatedResult"),
						jen.ID("token"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid client ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleSecret").Op(":=").ID("make").Call(
						jen.Index().ID("byte"),
						jen.ID("validClientSecretSize"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("c").Dot("fetchAuthTokenForAPIClient").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("unauthenticatedClient"),
						jen.Lit(""),
						jen.ID("exampleSecret"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("token"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil secret key"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("exampleClientID").Op(":=").Lit("example_client_id"),
					jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("c").Dot("fetchAuthTokenForAPIClient").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("unauthenticatedClient"),
						jen.ID("exampleClientID"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("token"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil HTTP client"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleClientID").Op(":=").Lit("example_client_id"),
					jen.ID("exampleSecret").Op(":=").ID("make").Call(
						jen.Index().ID("byte"),
						jen.ID("validClientSecretSize"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("c").Dot("fetchAuthTokenForAPIClient").Call(
						jen.ID("ctx"),
						jen.ID("nil"),
						jen.ID("exampleClientID"),
						jen.ID("exampleSecret"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("token"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("anticipatedResult").Op(":=").Lit("v2.local.QAxIpVe-ECVNI1z4xQbm_qQYomyT3h8FtV8bxkz8pBJWkT8f7HtlOpbroPDEZUKop_vaglyp76CzYy375cHmKCW8e1CCkV0Lflu4GTDyXMqQdpZMM1E6OaoQW27gaRSvWBrR3IgbFIa0AkuUFw.UGFyYWdvbiBJbml0aWF0aXZlIEVudGVycHJpc2Vz"),
					jen.ID("ts").Op(":=").ID("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
						jen.ID("response").Op(":=").Op("&").ID("types").Dot("PASETOResponse").Valuesln(jen.ID("Token").Op(":").ID("anticipatedResult")),
						jen.ID("assert").Dot("NotEmpty").Call(
							jen.ID("t"),
							jen.ID("req").Dot("Header").Dot("Get").Call(jen.ID("signatureHeaderKey")),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("response")),
						),
					))),
					jen.ID("c").Op(":=").ID("buildTestClient").Call(
						jen.ID("t"),
						jen.ID("ts"),
					),
					jen.ID("exampleClientID").Op(":=").Lit("example_client_id"),
					jen.ID("exampleSecret").Op(":=").ID("make").Call(
						jen.Index().ID("byte"),
						jen.ID("validClientSecretSize"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("c").Dot("fetchAuthTokenForAPIClient").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("unauthenticatedClient"),
						jen.ID("exampleClientID"),
						jen.ID("exampleSecret"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("anticipatedResult"),
						jen.ID("token"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("exampleClientID").Op(":=").Lit("example_client_id"),
					jen.ID("exampleSecret").Op(":=").ID("make").Call(
						jen.Index().ID("byte"),
						jen.ID("validClientSecretSize"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("c").Dot("fetchAuthTokenForAPIClient").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("unauthenticatedClient"),
						jen.ID("exampleClientID"),
						jen.ID("exampleSecret"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("token"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ts").Op(":=").ID("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
						jen.ID("assert").Dot("NotEmpty").Call(
							jen.ID("t"),
							jen.ID("req").Dot("Header").Dot("Get").Call(jen.ID("signatureHeaderKey")),
						),
						jen.Qual("time", "Sleep").Call(jen.Qual("time", "Minute")),
					))),
					jen.ID("c").Op(":=").ID("buildTestClient").Call(
						jen.ID("t"),
						jen.ID("ts"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("SetOptions").Call(jen.ID("UsingTimeout").Call(jen.Qual("time", "Nanosecond"))),
					),
					jen.ID("exampleClientID").Op(":=").Lit("example_client_id"),
					jen.ID("exampleSecret").Op(":=").ID("make").Call(
						jen.Index().ID("byte"),
						jen.ID("validClientSecretSize"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("c").Dot("fetchAuthTokenForAPIClient").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("unauthenticatedClient"),
						jen.ID("exampleClientID"),
						jen.ID("exampleSecret"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("token"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid status code"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ts").Op(":=").ID("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
						jen.ID("assert").Dot("NotEmpty").Call(
							jen.ID("t"),
							jen.ID("req").Dot("Header").Dot("Get").Call(jen.ID("signatureHeaderKey")),
						),
						jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusUnauthorized")),
					))),
					jen.ID("c").Op(":=").ID("buildTestClient").Call(
						jen.ID("t"),
						jen.ID("ts"),
					),
					jen.ID("exampleClientID").Op(":=").Lit("example_client_id"),
					jen.ID("exampleSecret").Op(":=").ID("make").Call(
						jen.Index().ID("byte"),
						jen.ID("validClientSecretSize"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("c").Dot("fetchAuthTokenForAPIClient").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("unauthenticatedClient"),
						jen.ID("exampleClientID"),
						jen.ID("exampleSecret"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("token"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid response from server"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ts").Op(":=").ID("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
						jen.ID("assert").Dot("NotEmpty").Call(
							jen.ID("t"),
							jen.ID("req").Dot("Header").Dot("Get").Call(jen.ID("signatureHeaderKey")),
						),
						jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("res").Dot("Write").Call(jen.Index().ID("byte").Call(jen.Lit("BLAH"))),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
					))),
					jen.ID("c").Op(":=").ID("buildTestClient").Call(
						jen.ID("t"),
						jen.ID("ts"),
					),
					jen.ID("exampleClientID").Op(":=").Lit("example_client_id"),
					jen.ID("exampleSecret").Op(":=").ID("make").Call(
						jen.Index().ID("byte"),
						jen.ID("validClientSecretSize"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("c").Dot("fetchAuthTokenForAPIClient").Call(
						jen.ID("ctx"),
						jen.ID("c").Dot("unauthenticatedClient"),
						jen.ID("exampleClientID"),
						jen.ID("exampleSecret"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("token"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
