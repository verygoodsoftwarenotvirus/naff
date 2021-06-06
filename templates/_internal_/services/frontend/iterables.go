package frontend

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, true)

	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()
	puvn := typ.Name.PluralUnexportedVarName()
	rn := typ.Name.RouteName()
	prn := typ.Name.PluralRouteName()
	scn := typ.Name.SingularCommonName()
	pcn := typ.Name.PluralCommonName()

	code.Add(
		jen.Const().Defs(
			jen.IDf("%sIDURLParamKey", uvn).Equals().Lit(rn),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("fetch%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID(uvn).Op("*").Qual(proj.TypesPackage(), sn), jen.Err().ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			constants.LoggerVar().Assign().ID("s").Dot("logger"),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.If(jen.ID("s").Dot("useFakeData")).Body(
				jen.IDf(uvn).Equals().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call()).Else().Body(
				jen.IDf("%sID", uvn).Assign().ID("s").Dot("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
					constants.LoggerVar(),
					jen.IDf("%sIDURLParamKey", uvn),
					jen.Lit(scn),
				).Call(jen.ID("req")),
				jen.List(jen.IDf(uvn), jen.Err()).Equals().ID("s").Dot("dataStore").Dotf("Get%s", sn).Call(
					jen.ID("ctx"),
					jen.IDf("%sID", uvn),
					jen.ID("sessionCtxData").Dot("ActiveAccountID"),
				),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
						jen.Err(),
						constants.LoggerVar(),
						jen.ID("span"),
						jen.Litf("fetching %s data", scn),
					),
					),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID(uvn), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("//go:embed templates/partials/generated/creators/%s_creator.gotpl", rn),
		jen.Newline(),
		jen.Var().IDf("%sCreatorTemplate", uvn).String(),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("build%sCreatorView", sn).Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Func().Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request"))).Body(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Newline(),
				constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
				jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
					jen.ID("span"),
					jen.ID("req"),
				),
				jen.Newline(),
				jen.List(jen.ID("sessionCtxData"), jen.Err()).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
						jen.Err(),
						constants.LoggerVar(),
						jen.ID("span"),
						jen.Lit("no session context data attached to request"),
					),
					jen.Qual("net/http", "Redirect").Call(
						jen.ID("res"),
						jen.ID("req"),
						jen.Lit("/login"),
						jen.ID("unauthorizedRedirectResponseCode"),
					),
					jen.Return(),
				),
				jen.Newline(),
				jen.ID(uvn).Assign().Op("&").Qual(proj.TypesPackage(), sn).Values(),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("view").Assign().ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.IDf("%sCreatorTemplate", uvn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("page").Assign().Op("&").ID("pageData").Valuesln(
						jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").DoesNotEqual().ID("nil"), jen.ID("Title").Op(":").Litf("New %s", sn), jen.ID("ContentData").Op(":").ID(uvn),
					),
					jen.If(jen.ID("sessionCtxData").DoesNotEqual().ID("nil")).Body(
						jen.ID("page").Dot("IsServiceAdmin").Equals().ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call(),
					),
					jen.Newline(),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("view"),
						jen.ID("page"),
						jen.ID("res"),
					),
				).Else().Body(
					jen.ID("tmpl").Assign().ID("s").Dot("parseTemplate").Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.IDf("%sCreatorTemplate", uvn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID(uvn),
						jen.ID("res"),
					),
				),
			),
		),
		jen.Newline(),
	)

	allFieldsKeys := []jen.Code{}
	creationKeys := []jen.Code{}
	updateKeys := []jen.Code{}

	for _, field := range typ.Fields {
		fuvn := field.Name.UnexportedVarName()
		fsn := field.Name.Singular()

		allFieldsKeys = append(allFieldsKeys, jen.IDf("%sFormKey", fuvn).Equals().Lit(fuvn))

		if field.ValidForCreationInput {
			creationKeys = append(creationKeys, jen.IDf("%sCreationInput%sFormKey", uvn, fsn).Equals().IDf("%sFormKey", fuvn))
		}
		if field.ValidForUpdateInput {
			updateKeys = append(updateKeys, jen.IDf("%sUpdateInput%sFormKey", uvn, fsn).Equals().IDf("%sFormKey", fuvn))
		}
	}

	allKeys := append(allFieldsKeys, jen.Newline())
	allKeys = append(allKeys, creationKeys...)
	allKeys = append(allKeys, jen.Newline())
	allKeys = append(allKeys, updateKeys...)
	allKeys = append(allKeys, jen.Newline())

	code.Add(
		jen.Const().Defs(allKeys...),
		jen.Newline(),
	)

	creationInputFields := []jen.Code{}
	for _, field := range typ.Fields {
		fsn := field.Name.Singular()
		creationInputFields = append(creationInputFields, jen.ID(fsn).Op(":").ID("form").Dot("Get").Call(jen.IDf("%sCreationInput%sFormKey", uvn, fsn)))
	}
	creationInputFields = append(creationInputFields, jen.ID("BelongsToAccount").Op(":").ID("sessionCtxData").Dot("ActiveAccountID"))

	code.Add(
		jen.Commentf("parseFormEncoded%sCreationInput checks a request for an %sCreationInput.", sn, sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("parseFormEncoded%sCreationInput", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("sessionCtxData").Op("*").Qual(proj.TypesPackage(), "SessionContextData")).Params(jen.ID("creationInput").Op("*").ID("types").Dotf("%sCreationInput", sn)).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.List(jen.ID("form"), jen.Err()).Assign().ID("s").Dot("extractFormFromRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("parsing %s creation input", scn),
				),
				jen.Return().ID("nil"),
			),
			jen.Newline(),
			jen.ID("creationInput").Equals().Op("&").ID("types").Dotf("%sCreationInput", sn).Valuesln(
				creationInputFields...,
			),
			jen.Newline(),
			jen.If(jen.Err().Equals().ID("creationInput").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.Err().DoesNotEqual().ID("nil")).Body(
				constants.LoggerVar().Equals().ID("logger").Dot("WithValue").Call(
					jen.Lit("input"),
					jen.ID("creationInput"),
				),
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("invalid %s creation input", scn),
				),
				jen.Return().ID("nil"),
			),
			jen.Newline(),
			jen.Return().ID("creationInput"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("handle%sCreationRequest", sn).Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			constants.LoggerVar().Dot("Debug").Call(jen.Litf("%s creation route called", scn)),
			jen.Newline(),
			jen.List(jen.ID("sessionCtxData"), jen.Err()).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Lit("no session context data attached to request"),
				),
				jen.Qual("net/http", "Redirect").Call(
					jen.ID("res"),
					jen.ID("req"),
					jen.Lit("/login"),
					jen.ID("unauthorizedRedirectResponseCode"),
				),
				jen.Return(),
			),
			jen.Newline(),
			constants.LoggerVar().Dot("Debug").Call(jen.Litf("session context data retrieved for %s creation route", scn)),
			jen.Newline(),
			jen.ID("creationInput").Assign().ID("s").Dotf("parseFormEncoded%sCreationInput", sn).Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("sessionCtxData"),
			),
			jen.If(jen.ID("creationInput").Op("==").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("parsing %s creation input", scn),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
				jen.Return(),
			),
			jen.Newline(),
			constants.LoggerVar().Dot("Debug").Call(jen.Litf("%s creation input parsed successfully", scn)),
			jen.Newline(),
			jen.If(jen.List(jen.ID("_"), jen.Err()).Equals().ID("s").Dot("dataStore").Dotf("Create%s", sn).Call(
				jen.ID("ctx"),
				jen.ID("creationInput"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("writing %s to datastore", scn),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.Newline(),
			constants.LoggerVar().Dot("Debug").Call(jen.Litf("%s created", scn)),
			jen.Newline(),
			jen.ID("htmxRedirectTo").Call(
				jen.ID("res"),
				jen.Litf("/%s", prn),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusCreated")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("//go:embed templates/partials/generated/editors/%s_editor.gotpl", rn),
		jen.Newline(),
		jen.Var().IDf("%sEditorTemplate", uvn).String(),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("build%sEditorView", sn).Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Func().Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request"))).Body(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Newline(),
				constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
				jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
					jen.ID("span"),
					jen.ID("req"),
				),
				jen.Newline(),
				jen.List(jen.ID("sessionCtxData"), jen.Err()).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
						jen.Err(),
						constants.LoggerVar(),
						jen.ID("span"),
						jen.Lit("no session context data attached to request"),
					),
					jen.Qual("net/http", "Redirect").Call(
						jen.ID("res"),
						jen.ID("req"),
						jen.Lit("/login"),
						jen.ID("unauthorizedRedirectResponseCode"),
					),
					jen.Return(),
				),
				jen.Newline(),
				jen.List(jen.ID(uvn), jen.Err()).Assign().ID("s").Dotf("fetch%s", sn).Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData"),
					jen.ID("req"),
				),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
						jen.Err(),
						constants.LoggerVar(),
						jen.ID("span"),
						jen.Litf("fetching %s from datastore", scn),
					),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.Newline(),
				jen.ID("tmplFuncMap").Assign().Map(jen.String()).Interface().Valuesln(
					jen.Lit("componentTitle").Op(":").Func().Params(jen.ID("x").Op("*").Qual(proj.TypesPackage(), sn)).Params(jen.String()).Body(
						jen.Return().Qual("fmt", "Sprintf").Call(
							jen.Lit(sn+" #%d"),
							jen.ID("x").Dot("ID"),
						),
					),
				),
				jen.Newline(),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("view").Assign().ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.IDf("%sEditorTemplate", uvn),
						jen.ID("tmplFuncMap"),
					),
					jen.Newline(),
					jen.ID("page").Assign().Op("&").ID("pageData").Valuesln(
						jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").DoesNotEqual().ID("nil"), jen.ID("Title").Op(":").Qual("fmt", "Sprintf").Call(
							jen.Lit(fmt.Sprintf("%s ", sn)+"#%d"),
							jen.ID(uvn).Dot("ID"),
						), jen.ID("ContentData").Op(":").ID(uvn),
					),
					jen.If(jen.ID("sessionCtxData").DoesNotEqual().ID("nil")).Body(
						jen.ID("page").Dot("IsServiceAdmin").Equals().ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call(),
					),
					jen.Newline(),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("view"),
						jen.ID("page"),
						jen.ID("res"),
					),
				).Else().Body(
					jen.ID("tmpl").Assign().ID("s").Dot("parseTemplate").Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.IDf("%sEditorTemplate", uvn),
						jen.ID("tmplFuncMap"),
					),
					jen.Newline(),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID(uvn),
						jen.ID("res"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("fetch%s", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID(puvn).Op("*").ID("types").Dotf("%sList", sn), jen.Err().ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			constants.LoggerVar().Assign().ID("s").Dot("logger"),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.If(jen.ID("s").Dot("useFakeData")).Body(
				jen.ID(puvn).Equals().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call()).Else().Body(
				jen.ID("filter").Assign().Qual(proj.TypesPackage(), "ExtractQueryFilter").Call(jen.ID("req")),
				jen.List(jen.ID(puvn), jen.Err()).Equals().ID("s").Dot("dataStore").Dotf("Get%s", pn).Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData").Dot("ActiveAccountID"),
					jen.ID("filter"),
				),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
						jen.Err(),
						constants.LoggerVar(),
						jen.ID("span"),
						jen.Litf("fetching %s data", scn),
					),
					),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID(puvn), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("//go:embed templates/partials/generated/tables/%s_table.gotpl", prn),
		jen.Newline(),
		jen.Var().IDf("%sTableTemplate", puvn).String(),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("build%sTableView", pn).Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Func().Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request"))).Body(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Newline(),
				constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
				jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
					jen.ID("span"),
					jen.ID("req"),
				),
				jen.Newline(),
				jen.List(jen.ID("sessionCtxData"), jen.Err()).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
						jen.Err(),
						constants.LoggerVar(),
						jen.ID("span"),
						jen.Lit("no session context data attached to request"),
					),
					jen.Qual("net/http", "Redirect").Call(
						jen.ID("res"),
						jen.ID("req"),
						jen.Lit("/login"),
						jen.ID("unauthorizedRedirectResponseCode"),
					),
					jen.Return(),
				),
				jen.Newline(),
				jen.List(jen.ID(puvn), jen.Err()).Assign().ID("s").Dotf("fetch%s", pn).Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData"),
					jen.ID("req"),
				),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
						jen.Err(),
						constants.LoggerVar(),
						jen.ID("span"),
						jen.Litf("fetching %s from datastore", pcn),
					),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.Newline(),
				jen.ID("tmplFuncMap").Assign().Map(jen.String()).Interface().Valuesln(
					jen.Lit("individualURL").Op(":").Func().Params(jen.ID("x").Op("*").Qual(proj.TypesPackage(), sn)).Params(jen.Qual("html/template", "URL")).Body(
						jen.Comment("#nosec G203"),
						jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
							jen.Lit(fmt.Sprintf("/dashboard_pages/%s/", prn)+"%d"),
							jen.ID("x").Dot("ID"),
						)),
					), jen.Lit("pushURL").Op(":").Func().Params(jen.ID("x").Op("*").Qual(proj.TypesPackage(), sn)).Params(jen.Qual("html/template", "URL")).Body(
						jen.Comment("#nosec G203"),
						jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
							jen.Lit(fmt.Sprintf("/%s/", prn)+"%d"),
							jen.ID("x").Dot("ID"),
						)),
					),
				),
				jen.Newline(),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("tmpl").Assign().ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.IDf("%sTableTemplate", puvn),
						jen.ID("tmplFuncMap"),
					),
					jen.Newline(),
					jen.ID("page").Assign().Op("&").ID("pageData").Valuesln(
						jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").DoesNotEqual().ID("nil"), jen.ID("Title").Op(":").Lit(pn), jen.ID("ContentData").Op(":").ID(puvn),
					),
					jen.If(jen.ID("sessionCtxData").DoesNotEqual().ID("nil")).Body(
						jen.ID("page").Dot("IsServiceAdmin").Equals().ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call(),
					),
					jen.Newline(),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("page"),
						jen.ID("res"),
					),
				).Else().Body(
					jen.ID("tmpl").Assign().ID("s").Dot("parseTemplate").Call(
						jen.ID("ctx"),
						jen.Lit("dashboard"),
						jen.IDf("%sTableTemplate", puvn),
						jen.ID("tmplFuncMap"),
					),
					jen.Newline(),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID(puvn),
						jen.ID("res"),
					),
				),
			),
		),
		jen.Newline(),
	)

	updateInputFields := []jen.Code{}
	for _, field := range typ.Fields {
		fsn := field.Name.Singular()
		updateInputFields = append(updateInputFields, jen.ID(fsn).Op(":").ID("form").Dot("Get").Call(jen.IDf("%sUpdateInput%sFormKey", uvn, fsn)))
	}
	updateInputFields = append(updateInputFields, jen.ID("BelongsToAccount").Op(":").ID("sessionCtxData").Dot("ActiveAccountID"))

	code.Add(
		jen.Commentf("parseFormEncoded%sUpdateInput checks a request for an %sUpdateInput.", sn, sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("parseFormEncoded%sUpdateInput", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("sessionCtxData").Op("*").Qual(proj.TypesPackage(), "SessionContextData")).Params(jen.ID("updateInput").Op("*").ID("types").Dotf("%sUpdateInput", sn)).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.List(jen.ID("form"), jen.Err()).Assign().ID("s").Dot("extractFormFromRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("parsing %s creation input", scn),
				),
				jen.Return().ID("nil"),
			),
			jen.Newline(),
			jen.ID("updateInput").Equals().Op("&").ID("types").Dotf("%sUpdateInput", sn).Valuesln(
				updateInputFields...,
			),
			jen.Newline(),
			jen.If(jen.Err().Equals().ID("updateInput").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.Err().DoesNotEqual().ID("nil")).Body(
				constants.LoggerVar().Equals().ID("logger").Dot("WithValue").Call(
					jen.Lit("input"),
					jen.ID("updateInput"),
				),
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("invalid %s creation input", scn),
				),
				jen.Return().ID("nil"),
			),
			jen.Newline(),
			jen.Return().ID("updateInput"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("handle%sUpdateRequest", sn).Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.List(jen.ID("sessionCtxData"), jen.Err()).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Lit("no session context data attached to request"),
				),
				jen.Qual("net/http", "Redirect").Call(
					jen.ID("res"),
					jen.ID("req"),
					jen.Lit("/login"),
					jen.ID("unauthorizedRedirectResponseCode"),
				),
				jen.Return(),
			),
			jen.Newline(),
			jen.ID("updateInput").Assign().ID("s").Dotf("parseFormEncoded%sUpdateInput", sn).Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("sessionCtxData"),
			),
			jen.If(jen.ID("updateInput").Op("==").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Lit("no update input attached to request"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
				jen.Return(),
			),
			jen.Newline(),
			jen.List(jen.ID(uvn), jen.Err()).Assign().ID("s").Dotf("fetch%s", sn).Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData"),
				jen.ID("req"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("fetching %s from datastore", scn),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.Newline(),
			jen.ID("changes").Assign().ID(uvn).Dot("Update").Call(jen.ID("updateInput")),
			jen.Newline(),
			jen.If(jen.Err().Equals().ID("s").Dot("dataStore").Dotf("Update%s", sn).Call(
				jen.ID("ctx"),
				jen.ID(uvn),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
				jen.ID("changes"),
			), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("fetching %s from datastore", scn),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.Newline(),
			jen.ID("tmplFuncMap").Assign().Map(jen.String()).Interface().Valuesln(
				jen.Lit("componentTitle").Op(":").Func().Params(jen.ID("x").Op("*").Qual(proj.TypesPackage(), sn)).Params(jen.String()).Body(
					jen.Return().Qual("fmt", "Sprintf").Call(
						jen.Lit(sn+" #%d"),
						jen.ID("x").Dot("ID"),
					),
				),
			),
			jen.Newline(),
			jen.ID("tmpl").Assign().ID("s").Dot("parseTemplate").Call(
				jen.ID("ctx"),
				jen.Lit(""),
				jen.IDf("%sEditorTemplate", uvn),
				jen.ID("tmplFuncMap"),
			),
			jen.Newline(),
			jen.ID("s").Dot("renderTemplateToResponse").Call(
				jen.ID("ctx"),
				jen.ID("tmpl"),
				jen.ID(uvn),
				jen.ID("res"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("handle%sDeletionRequest", sn).Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.List(jen.ID("sessionCtxData"), jen.Err()).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Lit("no session context data attached to request"),
				),
				jen.Qual("net/http", "Redirect").Call(
					jen.ID("res"),
					jen.ID("req"),
					jen.Lit("/login"),
					jen.ID("unauthorizedRedirectResponseCode"),
				),
				jen.Return(),
			),
			jen.Newline(),
			jen.IDf("%sID", uvn).Assign().ID("s").Dot("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
				constants.LoggerVar(),
				jen.IDf("%sIDURLParamKey", uvn),
				jen.Lit(scn),
			).Call(jen.ID("req")),
			jen.If(jen.Err().Equals().ID("s").Dot("dataStore").Dotf("Archive%s", sn).Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("archiving %s in datastore", pcn),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.Newline(),
			jen.List(jen.ID(puvn), jen.Err()).Assign().ID("s").Dotf("fetch%s", pn).Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData"),
				jen.ID("req"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("fetching %s from datastore", pcn),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.Newline(),
			jen.ID("tmplFuncMap").Assign().Map(jen.String()).Interface().Valuesln(
				jen.Lit("individualURL").Op(":").Func().Params(jen.ID("x").Op("*").Qual(proj.TypesPackage(), sn)).Params(jen.Qual("html/template", "URL")).Body(
					jen.Comment("#nosec G203"),
					jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
						jen.Lit(fmt.Sprintf("/dashboard_pages/%s/", prn)+"%d"),
						jen.ID("x").Dot("ID"),
					),
					),
				), jen.Lit("pushURL").Op(":").Func().Params(jen.ID("x").Op("*").Qual(proj.TypesPackage(), sn)).Params(jen.Qual("html/template", "URL")).Body(
					jen.Comment("#nosec G203"),
					jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
						jen.Lit(fmt.Sprintf("/%s/", prn)+"%d"),
						jen.ID("x").Dot("ID"),
					),
					),
				),
			),
			jen.Newline(),
			jen.ID("tmpl").Assign().ID("s").Dot("parseTemplate").Call(
				jen.ID("ctx"),
				jen.Lit("dashboard"),
				jen.IDf("%sTableTemplate", puvn),
				jen.ID("tmplFuncMap"),
			),
			jen.Newline(),
			jen.ID("s").Dot("renderTemplateToResponse").Call(
				jen.ID("ctx"),
				jen.ID("tmpl"),
				jen.ID(puvn),
				jen.ID("res"),
			),
		),
		jen.Newline(),
	)

	return code
}
