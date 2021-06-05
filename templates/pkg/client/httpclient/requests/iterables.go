package requests

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	uvn := typ.Name.UnexportedVarName()
	puvn := typ.Name.PluralUnexportedVarName()

	code.Add(
		jen.Const().Defs(
			jen.IDf("%sBasePath", puvn).Op("=").Lit(typ.Name.PluralRouteName()),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("Build%sExistsRequest builds an HTTP request for checking the existence of %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("Build%sExistsRequest", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.IDf("%sID", uvn).ID("uint64")).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.IDf("%sID", uvn).Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided")),
			),
			jen.Newline(),
			jen.ID("logger").Assign().ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dotf("%sIDKey", sn),
				jen.IDf("%sID", uvn),
			),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.IDf("%sBasePath", puvn),
				jen.ID("id").Call(jen.IDf("%sID", uvn)),
			),
			jen.ID("tracing").Dot("AttachRequestURIToSpan").Call(
				jen.ID("span"),
				jen.ID("uri"),
			),
			jen.Newline(),
			jen.List(jen.ID("req"), jen.ID("err")).Assign().Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodHead"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildGet%sRequest builds an HTTP request for fetching %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("BuildGet%sRequest", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.IDf("%sID", uvn).ID("uint64")).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.IDf("%sID", uvn).Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided")),
			),
			jen.Newline(),
			jen.ID("logger").Assign().ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dotf("%sIDKey", sn),
				jen.IDf("%sID", uvn),
			),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.IDf("%sBasePath", puvn),
				jen.ID("id").Call(jen.IDf("%sID", uvn)),
			),
			jen.ID("tracing").Dot("AttachRequestURIToSpan").Call(
				jen.ID("span"),
				jen.ID("uri"),
			),
			jen.Newline(),
			jen.List(jen.ID("req"), jen.ID("err")).Assign().Qual("net/http", "NewRequestWithContext").Call(
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
				),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	if typ.SearchEnabled {
		code.Add(
			jen.Commentf("BuildSearch%sRequest builds an HTTP request for querying %s.", pn, pcn),
			jen.Newline(),
			jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("BuildSearch%sRequest", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string"), jen.ID("limit").ID("uint8")).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Newline(),
				jen.ID("params").Assign().Qual("net/url", "Values").Values(),
				jen.ID("params").Dot("Set").Call(
					jen.ID("types").Dot("SearchQueryKey"),
					jen.ID("query"),
				),
				jen.ID("params").Dot("Set").Call(
					jen.ID("types").Dot("LimitQueryKey"),
					jen.Qual("strconv", "FormatUint").Call(
						jen.ID("uint64").Call(jen.ID("limit")),
						jen.Lit(10),
					),
				),
				jen.Newline(),
				jen.ID("logger").Assign().ID("b").Dot("logger").Dot("WithValue").Call(
					jen.ID("types").Dot("SearchQueryKey"),
					jen.ID("query"),
				).
					Dotln("WithValue").Call(
					jen.ID("types").Dot("LimitQueryKey"),
					jen.ID("limit"),
				),
				jen.Newline(),
				jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(
					jen.ID("ctx"),
					jen.ID("params"),
					jen.IDf("%sBasePath", puvn),
					jen.Lit("search"),
				),
				jen.ID("tracing").Dot("AttachRequestURIToSpan").Call(
					jen.ID("span"),
					jen.ID("uri"),
				),
				jen.Newline(),
				jen.List(jen.ID("req"), jen.ID("err")).Assign().Qual("net/http", "NewRequestWithContext").Call(
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
					),
					),
				),
				jen.Newline(),
				jen.Return().List(jen.ID("req"), jen.ID("nil")),
			),
			jen.Newline(),
		)
	}

	code.Add(
		jen.Commentf("BuildGet%sRequest builds an HTTP request for fetching a list of %s.", pn, pcn),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("BuildGet%sRequest", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").PointerTo().ID("types").Dot("QueryFilter")).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Assign().ID("filter").Dot("AttachToLogger").Call(jen.ID("b").Dot("logger")),
			jen.Newline(),
			jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(
				jen.ID("ctx"),
				jen.ID("filter").Dot("ToValues").Call(),
				jen.IDf("%sBasePath", puvn),
			),
			jen.ID("tracing").Dot("AttachRequestURIToSpan").Call(
				jen.ID("span"),
				jen.ID("uri"),
			),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.Newline(),
			jen.List(jen.ID("req"), jen.ID("err")).Assign().Qual("net/http", "NewRequestWithContext").Call(
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
				),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildCreate%sRequest builds an HTTP request for creating %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("BuildCreate%sRequest", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").PointerTo().ID("types").Dotf("%sCreationInput", sn)).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided")),
			),
			jen.Newline(),
			jen.ID("logger").Assign().ID("b").Dot("logger"),
			jen.Newline(),
			jen.If(jen.ID("err").Assign().ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating input"),
				),
				),
			),
			jen.Newline(),
			jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.IDf("%sBasePath", puvn),
			),
			jen.ID("tracing").Dot("AttachRequestURIToSpan").Call(
				jen.ID("span"),
				jen.ID("uri"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildDataRequest").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("input"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildUpdate%sRequest builds an HTTP request for updating %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("BuildUpdate%sRequest", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID(uvn).PointerTo().ID("types").Dot(sn)).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID(uvn).Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided")),
			),
			jen.Newline(),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.ID(uvn).Dot("ID"),
			),
			jen.Newline(),
			jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.IDf("%sBasePath", puvn),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID(uvn).Dot("ID"),
					jen.Lit(10),
				),
			),
			jen.ID("tracing").Dot("AttachRequestURIToSpan").Call(
				jen.ID("span"),
				jen.ID("uri"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildDataRequest").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPut"),
				jen.ID("uri"),
				jen.ID(uvn),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildArchive%sRequest builds an HTTP request for archiving %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("BuildArchive%sRequest", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.IDf("%sID", uvn).ID("uint64")).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.IDf("%sID", uvn).Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided")),
			),
			jen.Newline(),
			jen.ID("logger").Assign().ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dotf("%sIDKey", sn),
				jen.IDf("%sID", uvn),
			),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.IDf("%sBasePath", puvn),
				jen.ID("id").Call(jen.IDf("%sID", uvn)),
			),
			jen.ID("tracing").Dot("AttachRequestURIToSpan").Call(
				jen.ID("span"),
				jen.ID("uri"),
			),
			jen.Newline(),
			jen.List(jen.ID("req"), jen.ID("err")).Assign().Qual("net/http", "NewRequestWithContext").Call(
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
				),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildGetAuditLogFor%sRequest builds an HTTP request for fetching a list of audit log entries pertaining to %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("BuildGetAuditLogFor%sRequest", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.IDf("%sID", uvn).ID("uint64")).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.IDf("%sID", uvn).Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided")),
			),
			jen.Newline(),
			jen.ID("logger").Assign().ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dotf("%sIDKey", sn),
				jen.IDf("%sID", uvn),
			),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.IDf("%sBasePath", puvn),
				jen.ID("id").Call(jen.IDf("%sID", uvn)),
				jen.Lit("audit"),
			),
			jen.ID("tracing").Dot("AttachRequestURIToSpan").Call(
				jen.ID("span"),
				jen.ID("uri"),
			),
			jen.Newline(),
			jen.List(jen.ID("req"), jen.ID("err")).Assign().Qual("net/http", "NewRequestWithContext").Call(
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
				),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	return code
}
