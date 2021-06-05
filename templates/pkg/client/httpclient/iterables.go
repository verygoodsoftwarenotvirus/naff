package httpclient

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	uvn := typ.Name.UnexportedVarName()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	puvn := typ.Name.PluralUnexportedVarName()

	code.Add(
		jen.Commentf("%sExists retrieves whether %s exists.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(
			jen.ID("c").Op("*").ID("Client")).IDf("%sExists", sn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.IDf("%sID", uvn).ID("uint64"),
		).Params(
			jen.ID("bool"),
			jen.ID("error")).Body(
			jen.List(jen.ID("ctx"),
				jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.IDf("%sID", uvn).Op("==").Lit(0)).Body(
				jen.Return().List(
					jen.ID("false"),
					jen.ID("ErrInvalidIDProvided"),
				),
			),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), fmt.Sprintf("%sIDKey", sn)),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.List(jen.ID("req"),
				jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dotf("Build%sExistsRequest", sn).Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(
					jen.ID("false"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Litf("building %s existence request", scn),
					),
				),
			),
			jen.Newline(),
			jen.List(jen.ID("exists"),
				jen.ID("err")).Op(":=").ID("c").Dot("responseIsOK").Call(
				jen.ID("ctx"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(
					jen.ID("false"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit(fmt.Sprintf("checking existence for %s", scn)+" #%d"),
						jen.IDf("%sID", uvn),
					),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("exists"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("Get%s gets %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(
			jen.ID("c").Op("*").ID("Client")).IDf("Get%s", sn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.IDf("%sID", uvn).ID("uint64"),
		).Params(
			jen.Op("*").Qual(proj.TypesPackage(), sn),
			jen.ID("error")).Body(
			jen.List(jen.ID("ctx"),
				jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.IDf("%sID", uvn).Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided")),
			),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), fmt.Sprintf("%sIDKey", sn)),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.List(jen.ID("req"),
				jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dotf("BuildGet%sRequest", sn).Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Litf("building get %s request", scn),
					),
				),
			),
			jen.Newline(),
			jen.Var().ID(uvn).Op("*").Qual(proj.TypesPackage(), sn),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID(uvn),
			),
				jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Litf("retrieving %s", scn),
					),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID(uvn), jen.ID("nil")),
		),
		jen.Newline(),
	)

	if typ.SearchEnabled {
		code.Add(
			jen.Commentf("Search%s searches through a list of %s.", pn, pcn),
			jen.Newline(),
			jen.Func().Params(
				jen.ID("c").Op("*").ID("Client")).IDf("Search%s", pn).Params(
				jen.ID("ctx").Qual("context", "Context"),
				jen.ID("query").ID("string"),
				jen.ID("limit").ID("uint8"),
			).Params(
				jen.Index().Op("*").Qual(proj.TypesPackage(), sn),
				jen.ID("error")).Body(
				jen.List(jen.ID("ctx"),
					jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Newline(),
				jen.If(jen.ID("query").Op("==").Lit("")).Body(
					jen.Return().List(
						jen.ID("nil"),
						jen.ID("ErrEmptyQueryProvided")),
				),
				jen.Newline(),
				jen.If(jen.ID("limit").Op("==").Lit(0)).Body(
					jen.ID("limit").Op("=").Lit(20)),
				jen.Newline(),
				jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
					jen.Qual(proj.ObservabilityPackage("keys"), "SearchQueryKey"),
					jen.ID("query"),
				).Dot("WithValue").Call(
					jen.Qual(proj.ObservabilityPackage("keys"), "FilterLimitKey"),
					jen.ID("limit"),
				),
				jen.Newline(),
				jen.List(jen.ID("req"),
					jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dotf("BuildSearch%sRequest", pn).Call(
					jen.ID("ctx"),
					jen.ID("query"),
					jen.ID("limit"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(
						jen.ID("nil"),
						jen.ID("observability").Dot("PrepareError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Litf("building search for %s request", pcn),
						),
					),
				),
				jen.Newline(),
				jen.Var().ID(puvn).Index().Op("*").Qual(proj.TypesPackage(), sn),
				jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.Op("&").ID(puvn),
				),
					jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(
						jen.ID("nil"),
						jen.ID("observability").Dot("PrepareError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Litf("retrieving %s", pcn),
						),
					),
				),
				jen.Newline(),
				jen.Return().List(jen.ID(puvn), jen.ID("nil")),
			),
			jen.Newline(),
		)
	}

	code.Add(
		jen.Commentf("Get%s retrieves a list of %s.", pn, pcn),
		jen.Newline(),
		jen.Func().Params(
			jen.ID("c").Op("*").ID("Client")).IDf("Get%s", pn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("filter").Op("*").Qual(proj.TypesPackage(), "QueryFilter"),
		).Params(
			jen.Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn)),
			jen.ID("error")).Body(
			jen.List(jen.ID("ctx"),
				jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("c").Dot("loggerWithFilter").Call(jen.ID("filter")),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.Newline(),
			jen.List(jen.ID("req"),
				jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dotf("BuildGet%sRequest", pn).Call(
				jen.ID("ctx"),
				jen.ID("filter"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Litf("building %s list request", pcn),
					),
				),
			),
			jen.Newline(),
			jen.Var().ID(puvn).Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn)),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID(puvn),
			),
				jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Litf("retrieving %s", pcn),
					),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID(puvn), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("Create%s creates %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(
			jen.ID("c").Op("*").ID("Client")).IDf("Create%s", sn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("input").Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn)),
		).Params(
			jen.Op("*").Qual(proj.TypesPackage(), sn),
			jen.ID("error")).Body(
			jen.List(jen.ID("ctx"),
				jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("ErrNilInputProvided"),
				),
			),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger"),
			jen.Newline(),
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")),
				jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("validating input"),
					),
				),
			),
			jen.Newline(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dotf("BuildCreate%sRequest", sn).Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Litf("building create %s request", scn),
					),
				),
			),
			jen.Newline(),
			jen.Var().ID(uvn).Op("*").Qual(proj.TypesPackage(), sn),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID(uvn),
			),
				jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Litf("creating %s", scn),
					),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID(uvn), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("Update%s updates %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(
			jen.ID("c").Op("*").ID("Client")).IDf("Update%s", sn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID(uvn).Op("*").Qual(proj.TypesPackage(), sn),
		).Params(
			jen.ID("error")).Body(
			jen.List(jen.ID("ctx"),
				jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID(uvn).Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided"),
			),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), fmt.Sprintf("%sIDKey", sn)),
				jen.ID(uvn).Dot("ID"),
			),
			jen.Newline(),
			jen.List(jen.ID("req"),
				jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dotf("BuildUpdate%sRequest", sn).Call(
				jen.ID("ctx"),
				jen.ID(uvn),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("building update %s request", scn),
				)),
			jen.Newline(),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID(uvn),
			),
				jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit(fmt.Sprintf("updating %s", scn)+" #%d"),
					jen.ID(uvn).Dot("ID"),
				)),
			jen.Newline(),
			jen.Return().ID("nil"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("Archive%s archives %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(
			jen.ID("c").Op("*").ID("Client")).IDf("Archive%s", sn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.IDf("%sID", uvn).ID("uint64"),
		).Params(
			jen.ID("error")).Body(
			jen.List(jen.ID("ctx"),
				jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.IDf("%sID", uvn).Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided"),
			),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), fmt.Sprintf("%sIDKey", sn)),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.List(jen.ID("req"),
				jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dotf("BuildArchive%sRequest", sn).Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("building archive %s request", scn),
				)),
			jen.Newline(),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			),
				jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit(fmt.Sprintf("archiving %s", scn)+" #%d"),
					jen.IDf("%sID", uvn),
				)),
			jen.Newline(),
			jen.Return().ID("nil"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("GetAuditLogFor%s retrieves a list of audit log entries pertaining to %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(
			jen.ID("c").Op("*").ID("Client")).IDf("GetAuditLogFor%s", sn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.IDf("%sID", uvn).ID("uint64"),
		).Params(
			jen.Index().Op("*").Qual(proj.TypesPackage(), "AuditLogEntry"),
			jen.ID("error")).Body(
			jen.List(jen.ID("ctx"),
				jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.IDf("%sID", uvn).Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("ErrInvalidIDProvided")),
			),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), fmt.Sprintf("%sIDKey", sn)),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.List(jen.ID("req"),
				jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dotf("BuildGetAuditLogFor%sRequest", sn).Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Litf("building get audit log entries for %s request", scn),
					),
				),
			),
			jen.Newline(),
			jen.Var().ID("entries").Index().Op("*").Qual(proj.TypesPackage(), "AuditLogEntry"),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("entries"),
			),
				jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("retrieving plan"),
					),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("entries"),
				jen.ID("nil")),
		),
		jen.Newline(),
	)

	return code
}
