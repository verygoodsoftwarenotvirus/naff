package requests

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func apiClientsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("apiClientsBasePath").Op("=").Lit("api_clients"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAPIClientRequest builds an HTTP request for fetching an API client."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildGetAPIClientRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("uint64")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("clientID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("APIClientDatabaseIDKey"),
				jen.ID("clientID"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("apiClientsBasePath"),
				jen.ID("id").Call(jen.ID("clientID")),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAPIClientsRequest builds an HTTP request for fetching a list of API clients."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildGetAPIClientsRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("filter").Dot("AttachToLogger").Call(jen.ID("b").Dot("logger")),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("filter").Dot("ToValues").Call(),
				jen.ID("apiClientsBasePath"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildCreateAPIClientRequest builds an HTTP request for creating an API client."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildCreateAPIClientRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cookie").Op("*").Qual("net/http", "Cookie"), jen.ID("input").Op("*").ID("types").Dot("APIClientCreationInput")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("cookie").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrCookieRequired"))),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("apiClientsBasePath"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("b").Dot("buildDataRequest").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildArchiveAPIClientRequest builds an HTTP request for archiving an API client."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildArchiveAPIClientRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("uint64")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("clientID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("APIClientDatabaseIDKey"),
				jen.ID("clientID"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("apiClientsBasePath"),
				jen.ID("id").Call(jen.ID("clientID")),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodDelete"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAuditLogForAPIClientRequest builds an HTTP request for fetching a list of audit log entries for an API client."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildGetAuditLogForAPIClientRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("uint64")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("clientID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("APIClientDatabaseIDKey"),
				jen.ID("clientID"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("apiClientsBasePath"),
				jen.ID("id").Call(jen.ID("clientID")),
				jen.Lit("audit"),
			),
			jen.ID("tracing").Dot("AttachRequestURIToSpan").Call(
				jen.ID("span"),
				jen.ID("uri"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
