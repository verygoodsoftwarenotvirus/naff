package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func itemsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Anon("embed")

	code.Add(
		jen.Var().Defs(
			jen.ID("itemIDURLParamKey").Op("=").Lit("item"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("fetchItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("item").Op("*").ID("types").Dot("Item"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger"),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.If(jen.ID("s").Dot("useFakeData")).Body(
				jen.ID("item").Op("=").ID("fakes").Dot("BuildFakeItem").Call()).Else().Body(
				jen.ID("itemID").Op(":=").ID("s").Dot("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
					jen.ID("logger"),
					jen.ID("itemIDURLParamKey"),
					jen.Lit("item"),
				).Call(jen.ID("req")),
				jen.List(jen.ID("item"), jen.ID("err")).Op("=").ID("s").Dot("dataStore").Dot("GetItem").Call(
					jen.ID("ctx"),
					jen.ID("itemID"),
					jen.ID("sessionCtxData").Dot("ActiveAccountID"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("fetching item data"),
					))),
			),
			jen.Return().List(jen.ID("item"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("itemCreatorTemplate").ID("string"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildItemCreatorView").Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request"))).Body(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
				jen.ID("tracing").Dot("AttachRequestToSpan").Call(
					jen.ID("span"),
					jen.ID("req"),
				),
				jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
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
				jen.ID("item").Op(":=").Op("&").ID("types").Dot("Item").Valuesln(),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("view").Op(":=").ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.ID("itemCreatorTemplate"),
						jen.ID("nil"),
					),
					jen.ID("page").Op(":=").Op("&").ID("pageData").Valuesln(jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").Op("!=").ID("nil"), jen.ID("Title").Op(":").Lit("New Item"), jen.ID("ContentData").Op(":").ID("item")),
					jen.If(jen.ID("sessionCtxData").Op("!=").ID("nil")).Body(
						jen.ID("page").Dot("IsServiceAdmin").Op("=").ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call()),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("view"),
						jen.ID("page"),
						jen.ID("res"),
					),
				).Else().Body(
					jen.ID("tmpl").Op(":=").ID("s").Dot("parseTemplate").Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.ID("itemCreatorTemplate"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("item"),
						jen.ID("res"),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("itemCreationInputNameFormKey").Op("=").Lit("name"),
			jen.ID("itemCreationInputDetailsFormKey").Op("=").Lit("details"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("parseFormEncodedItemCreationInput checks a request for an ItemCreationInput."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("parseFormEncodedItemCreationInput").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData")).Params(jen.ID("creationInput").Op("*").ID("types").Dot("ItemCreationInput")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("form"), jen.ID("err")).Op(":=").ID("s").Dot("extractFormFromRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("parsing item creation input"),
				),
				jen.Return().ID("nil"),
			),
			jen.ID("creationInput").Op("=").Op("&").ID("types").Dot("ItemCreationInput").Valuesln(jen.ID("Name").Op(":").ID("form").Dot("Get").Call(jen.ID("itemCreationInputNameFormKey")), jen.ID("Details").Op(":").ID("form").Dot("Get").Call(jen.ID("itemCreationInputDetailsFormKey")), jen.ID("BelongsToAccount").Op(":").ID("sessionCtxData").Dot("ActiveAccountID")),
			jen.If(jen.ID("err").Op("=").ID("creationInput").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
					jen.Lit("input"),
					jen.ID("creationInput"),
				),
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("invalid item creation input"),
				),
				jen.Return().ID("nil"),
			),
			jen.Return().ID("creationInput"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("handleItemCreationRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("item Creation route called")),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
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
			jen.ID("logger").Dot("Debug").Call(jen.Lit("session context data retrieved for item Creation route")),
			jen.ID("creationInput").Op(":=").ID("s").Dot("parseFormEncodedItemCreationInput").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("sessionCtxData"),
			),
			jen.If(jen.ID("creationInput").Op("==").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("parsing item creation input"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
				jen.Return(),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("item Creation input parsed successfully")),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op("=").ID("s").Dot("dataStore").Dot("CreateItem").Call(
				jen.ID("ctx"),
				jen.ID("creationInput"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing item to datastore"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("item Created")),
			jen.ID("htmxRedirectTo").Call(
				jen.ID("res"),
				jen.Lit("/items"),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusCreated")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("itemEditorTemplate").ID("string"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildItemEditorView").Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request"))).Body(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
				jen.ID("tracing").Dot("AttachRequestToSpan").Call(
					jen.ID("span"),
					jen.ID("req"),
				),
				jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
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
				jen.List(jen.ID("item"), jen.ID("err")).Op(":=").ID("s").Dot("fetchItem").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData"),
					jen.ID("req"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("fetching item from datastore"),
					),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("tmplFuncMap").Op(":=").Map(jen.ID("string")).Interface().Valuesln(jen.Lit("componentTitle").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot("Item")).Params(jen.ID("string")).Body(
					jen.Return().Qual("fmt", "Sprintf").Call(
						jen.Lit("Item #%d"),
						jen.ID("x").Dot("ID"),
					))),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("view").Op(":=").ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.ID("itemEditorTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("page").Op(":=").Op("&").ID("pageData").Valuesln(jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").Op("!=").ID("nil"), jen.ID("Title").Op(":").Qual("fmt", "Sprintf").Call(
						jen.Lit("Item #%d"),
						jen.ID("item").Dot("ID"),
					), jen.ID("ContentData").Op(":").ID("item")),
					jen.If(jen.ID("sessionCtxData").Op("!=").ID("nil")).Body(
						jen.ID("page").Dot("IsServiceAdmin").Op("=").ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call()),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("view"),
						jen.ID("page"),
						jen.ID("res"),
					),
				).Else().Body(
					jen.ID("tmpl").Op(":=").ID("s").Dot("parseTemplate").Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.ID("itemEditorTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("item"),
						jen.ID("res"),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("fetchItems").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("items").Op("*").ID("types").Dot("ItemList"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger"),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.If(jen.ID("s").Dot("useFakeData")).Body(
				jen.ID("items").Op("=").ID("fakes").Dot("BuildFakeItemList").Call()).Else().Body(
				jen.ID("filter").Op(":=").ID("types").Dot("ExtractQueryFilter").Call(jen.ID("req")),
				jen.List(jen.ID("items"), jen.ID("err")).Op("=").ID("s").Dot("dataStore").Dot("GetItems").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData").Dot("ActiveAccountID"),
					jen.ID("filter"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("fetching item data"),
					))),
			),
			jen.Return().List(jen.ID("items"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("itemsTableTemplate").ID("string"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildItemsTableView").Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request"))).Body(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
				jen.ID("tracing").Dot("AttachRequestToSpan").Call(
					jen.ID("span"),
					jen.ID("req"),
				),
				jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
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
				jen.List(jen.ID("items"), jen.ID("err")).Op(":=").ID("s").Dot("fetchItems").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData"),
					jen.ID("req"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("fetching items from datastore"),
					),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("tmplFuncMap").Op(":=").Map(jen.ID("string")).Interface().Valuesln(jen.Lit("individualURL").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot("Item")).Params(jen.Qual("html/template", "URL")).Body(
					jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("/dashboard_pages/items/%d"),
						jen.ID("x").Dot("ID"),
					))), jen.Lit("pushURL").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot("Item")).Params(jen.Qual("html/template", "URL")).Body(
					jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("/items/%d"),
						jen.ID("x").Dot("ID"),
					)))),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("tmpl").Op(":=").ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.ID("itemsTableTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("page").Op(":=").Op("&").ID("pageData").Valuesln(jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").Op("!=").ID("nil"), jen.ID("Title").Op(":").Lit("Items"), jen.ID("ContentData").Op(":").ID("items")),
					jen.If(jen.ID("sessionCtxData").Op("!=").ID("nil")).Body(
						jen.ID("page").Dot("IsServiceAdmin").Op("=").ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call()),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("page"),
						jen.ID("res"),
					),
				).Else().Body(
					jen.ID("tmpl").Op(":=").ID("s").Dot("parseTemplate").Call(
						jen.ID("ctx"),
						jen.Lit("dashboard"),
						jen.ID("itemsTableTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("items"),
						jen.ID("res"),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("parseFormEncodedItemUpdateInput checks a request for an ItemUpdateInput."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("parseFormEncodedItemUpdateInput").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData")).Params(jen.ID("updateInput").Op("*").ID("types").Dot("ItemUpdateInput")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("form"), jen.ID("err")).Op(":=").ID("s").Dot("extractFormFromRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("parsing item creation input"),
				),
				jen.Return().ID("nil"),
			),
			jen.ID("updateInput").Op("=").Op("&").ID("types").Dot("ItemUpdateInput").Valuesln(jen.ID("Name").Op(":").ID("form").Dot("Get").Call(jen.ID("itemCreationInputNameFormKey")), jen.ID("Details").Op(":").ID("form").Dot("Get").Call(jen.ID("itemCreationInputDetailsFormKey")), jen.ID("BelongsToAccount").Op(":").ID("sessionCtxData").Dot("ActiveAccountID")),
			jen.If(jen.ID("err").Op("=").ID("updateInput").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
					jen.Lit("input"),
					jen.ID("updateInput"),
				),
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("invalid item creation input"),
				),
				jen.Return().ID("nil"),
			),
			jen.Return().ID("updateInput"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("handleItemUpdateRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
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
			jen.ID("updateInput").Op(":=").ID("s").Dot("parseFormEncodedItemUpdateInput").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("sessionCtxData"),
			),
			jen.If(jen.ID("updateInput").Op("==").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("no update input attached to request"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
				jen.Return(),
			),
			jen.List(jen.ID("item"), jen.ID("err")).Op(":=").ID("s").Dot("fetchItem").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching item from datastore"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.ID("changes").Op(":=").ID("item").Dot("Update").Call(jen.ID("updateInput")),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("dataStore").Dot("UpdateItem").Call(
				jen.ID("ctx"),
				jen.ID("item"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
				jen.ID("changes"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching item from datastore"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.ID("tmplFuncMap").Op(":=").Map(jen.ID("string")).Interface().Valuesln(jen.Lit("componentTitle").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot("Item")).Params(jen.ID("string")).Body(
				jen.Return().Qual("fmt", "Sprintf").Call(
					jen.Lit("Item #%d"),
					jen.ID("x").Dot("ID"),
				))),
			jen.ID("tmpl").Op(":=").ID("s").Dot("parseTemplate").Call(
				jen.ID("ctx"),
				jen.Lit(""),
				jen.ID("itemEditorTemplate"),
				jen.ID("tmplFuncMap"),
			),
			jen.ID("s").Dot("renderTemplateToResponse").Call(
				jen.ID("ctx"),
				jen.ID("tmpl"),
				jen.ID("item"),
				jen.ID("res"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("handleItemDeletionRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
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
			jen.ID("itemID").Op(":=").ID("s").Dot("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
				jen.ID("logger"),
				jen.ID("itemIDURLParamKey"),
				jen.Lit("item"),
			).Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("dataStore").Dot("ArchiveItem").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("archiving items in datastore"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.List(jen.ID("items"), jen.ID("err")).Op(":=").ID("s").Dot("fetchItems").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching items from datastore"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.ID("tmplFuncMap").Op(":=").Map(jen.ID("string")).Interface().Valuesln(jen.Lit("individualURL").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot("Item")).Params(jen.Qual("html/template", "URL")).Body(
				jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("/dashboard_pages/items/%d"),
					jen.ID("x").Dot("ID"),
				))), jen.Lit("pushURL").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot("Item")).Params(jen.Qual("html/template", "URL")).Body(
				jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("/items/%d"),
					jen.ID("x").Dot("ID"),
				)))),
			jen.ID("tmpl").Op(":=").ID("s").Dot("parseTemplate").Call(
				jen.ID("ctx"),
				jen.Lit("dashboard"),
				jen.ID("itemsTableTemplate"),
				jen.ID("tmplFuncMap"),
			),
			jen.ID("s").Dot("renderTemplateToResponse").Call(
				jen.ID("ctx"),
				jen.ID("tmpl"),
				jen.ID("items"),
				jen.ID("res"),
			),
		),
		jen.Line(),
	)

	return code
}
