package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	sn := typ.Name.Singular()

	utils.AddImports(proj, code, true)

	code.Add(
		jen.Var().IDf("itemIDURLParamKey").Equals().Lit("item"),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("fetchItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("item").Op("*").ID("types").Dot(sn), jen.Err().ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			constants.LoggerVar().Assign().ID("s").Dot("logger"),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.If(jen.ID("s").Dot("useFakeData")).Body(
				jen.ID("item").Equals().ID("fakes").Dotf("BuildFakeItem").Call()).Else().Body(
				jen.IDf("itemID").Assign().ID("s").Dot("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
					constants.LoggerVar(),
					jen.IDf("itemIDURLParamKey"),
					jen.Lit("item"),
				).Call(jen.ID("req")),
				jen.List(jen.IDf("item"), jen.Err()).Equals().ID("s").Dot("dataStore").Dotf("GetItem").Call(
					jen.ID("ctx"),
					jen.IDf("itemID"),
					jen.ID("sessionCtxData").Dot("ActiveAccountID"),
				),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.Err(),
						constants.LoggerVar(),
						jen.ID("span"),
						jen.Litf("fetching item data"),
					))),
			),
			jen.Return().List(jen.ID("item"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().IDf("itemCreatorTemplate").String(),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("buildItemCreatorView").Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request")).Body(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
				jen.ID("tracing").Dot("AttachRequestToSpan").Call(
					jen.ID("span"),
					jen.ID("req"),
				),
				jen.List(jen.ID("sessionCtxData"), jen.Err()).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
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
				jen.ID("item").Assign().Op("&").ID("types").Dot(sn).Values(),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("view").Assign().ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.IDf("itemCreatorTemplate"),
						jen.ID("nil"),
					),
					jen.ID("page").Assign().Op("&").ID("pageData").Valuesln(
						jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").DoesNotEqual().ID("nil"), jen.ID("Title").Op(":").Litf("New Item"), jen.ID("ContentData").Op(":").ID("item")),
					jen.If(jen.ID("sessionCtxData").DoesNotEqual().ID("nil")).Body(
						jen.ID("page").Dot("IsServiceAdmin").Equals().ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call()),
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
						jen.IDf("itemCreatorTemplate"),
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
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(
			jen.IDf("itemCreationInputNameFormKey").Equals().Lit("name"),
			jen.IDf("itemCreationInputDetailsFormKey").Equals().Lit("details"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("parseFormEncodedItemCreationInput checks a request for an ItemCreationInput."),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("parseFormEncodedItemCreationInput").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData")).Params(jen.ID("creationInput").Op("*").ID("types").Dot("ItemCreationInput")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("form"), jen.Err()).Assign().ID("s").Dot("extractFormFromRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("parsing item creation input"),
				),
				jen.Return().ID("nil"),
			),
			jen.ID("creationInput").Equals().Op("&").ID("types").Dotf("ItemCreationInput").Valuesln(
				jen.ID("Name").Op(":").ID("form").Dot("Get").Call(jen.IDf("itemCreationInputNameFormKey")), jen.ID("Details").Op(":").ID("form").Dot("Get").Call(jen.ID("itemCreationInputDetailsFormKey")), jen.ID("BelongsToAccount").Op(":").ID("sessionCtxData").Dot("ActiveAccountID")),
			jen.If(jen.Err().Equals().ID("creationInput").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.Err().DoesNotEqual().ID("nil")).Body(
				constants.LoggerVar().Equals().ID("logger").Dot("WithValue").Call(
					jen.Lit("input"),
					jen.ID("creationInput"),
				),
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("invalid item creation input"),
				),
				jen.Return().ID("nil"),
			),
			jen.Return().ID("creationInput"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("handleItemCreationRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			constants.LoggerVar().Dot("Debug").Call(jen.Litf("item Creation route called")),
			jen.List(jen.ID("sessionCtxData"), jen.Err()).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
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
			constants.LoggerVar().Dot("Debug").Call(jen.Litf("session context data retrieved for item Creation route")),
			jen.ID("creationInput").Assign().ID("s").Dotf("parseFormEncodedItemCreationInput").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("sessionCtxData"),
			),
			jen.If(jen.ID("creationInput").Op("==").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("parsing item creation input"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
				jen.Return(),
			),
			constants.LoggerVar().Dot("Debug").Call(jen.Litf("item Creation input parsed successfully")),
			jen.If(jen.List(jen.ID("_"), jen.Err()).Equals().ID("s").Dot("dataStore").Dotf("CreateItem").Call(
				jen.ID("ctx"),
				jen.ID("creationInput"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("writing item to datastore"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			constants.LoggerVar().Dot("Debug").Call(jen.Litf("item Created")),
			jen.ID("htmxRedirectTo").Call(
				jen.ID("res"),
				jen.Litf("/items"),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusCreated")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().IDf("itemEditorTemplate").String(),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("buildItemEditorView").Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request")).Body(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
				jen.ID("tracing").Dot("AttachRequestToSpan").Call(
					jen.ID("span"),
					jen.ID("req"),
				),
				jen.List(jen.ID("sessionCtxData"), jen.Err()).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
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
				jen.List(jen.ID("item"), jen.Err()).Assign().ID("s").Dotf("fetchItem").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData"),
					jen.ID("req"),
				),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.Err(),
						constants.LoggerVar(),
						jen.ID("span"),
						jen.Litf("fetching item from datastore"),
					),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("tmplFuncMap").Assign().Map(jen.String()).Interface().Valuesln(
					jen.Lit("componentTitle").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot(sn)).Params(jen.String()).Body(
						jen.Return().Qual("fmt", "Sprintf").Call(
							jen.Litf("Item #%d"),
							jen.ID("x").Dot("ID"),
						))),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("view").Assign().ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.ID("itemEditorTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("page").Assign().Op("&").ID("pageData").Valuesln(
						jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").DoesNotEqual().ID("nil"), jen.ID("Title").Op(":").Qual("fmt", "Sprintf").Call(
							jen.Litf("Item #%d"),
							jen.ID("item").Dot("ID"),
						), jen.ID("ContentData").Op(":").ID("item")),
					jen.If(jen.ID("sessionCtxData").DoesNotEqual().ID("nil")).Body(
						jen.ID("page").Dot("IsServiceAdmin").Equals().ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call()),
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
						jen.IDf("itemEditorTemplate"),
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
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("fetchItems").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("items").Op("*").ID("types").Dot("ItemList"), jen.Err().ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			constants.LoggerVar().Assign().ID("s").Dot("logger"),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.If(jen.ID("s").Dot("useFakeData")).Body(
				jen.ID("items").Equals().ID("fakes").Dotf("BuildFakeItemList").Call()).Else().Body(
				jen.ID("filter").Assign().ID("types").Dot("ExtractQueryFilter").Call(jen.ID("req")),
				jen.List(jen.ID("items"), jen.Err()).Equals().ID("s").Dot("dataStore").Dotf("GetItems").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData").Dot("ActiveAccountID"),
					jen.ID("filter"),
				),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.Err(),
						constants.LoggerVar(),
						jen.ID("span"),
						jen.Litf("fetching item data"),
					))),
			),
			jen.Return().List(jen.ID("items"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().IDf("itemsTableTemplate").String(),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("buildItemsTableView").Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request")).Body(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
				jen.ID("tracing").Dot("AttachRequestToSpan").Call(
					jen.ID("span"),
					jen.ID("req"),
				),
				jen.List(jen.ID("sessionCtxData"), jen.Err()).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
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
				jen.List(jen.ID("items"), jen.Err()).Assign().ID("s").Dotf("fetchItems").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData"),
					jen.ID("req"),
				),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.Err(),
						constants.LoggerVar(),
						jen.ID("span"),
						jen.Lit("fetching items from datastore"),
					),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("tmplFuncMap").Assign().Map(jen.String()).Interface().Valuesln(
					jen.Lit("individualURL").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot(sn)).Params(jen.Qual("html/template", "URL")).Body(
						jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
							jen.Litf("/dashboard_pages/items/%d"),
							jen.ID("x").Dot("ID"),
						))), jen.Lit("pushURL").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot(sn)).Params(jen.Qual("html/template", "URL")).Body(
						jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
							jen.Litf("/items/%d"),
							jen.ID("x").Dot("ID"),
						)))),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("tmpl").Assign().ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.IDf("itemsTableTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("page").Assign().Op("&").ID("pageData").Valuesln(
						jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").DoesNotEqual().ID("nil"), jen.ID("Title").Op(":").Lit("Items"), jen.ID("ContentData").Op(":").ID("items")),
					jen.If(jen.ID("sessionCtxData").DoesNotEqual().ID("nil")).Body(
						jen.ID("page").Dot("IsServiceAdmin").Equals().ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call()),
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
						jen.IDf("itemsTableTemplate"),
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
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("parseFormEncodedItemUpdateInput checks a request for an ItemUpdateInput."),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("parseFormEncodedItemUpdateInput").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData")).Params(jen.ID("updateInput").Op("*").ID("types").Dot("ItemUpdateInput")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("form"), jen.Err()).Assign().ID("s").Dot("extractFormFromRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("parsing item creation input"),
				),
				jen.Return().ID("nil"),
			),
			jen.ID("updateInput").Equals().Op("&").ID("types").Dotf("ItemUpdateInput").Valuesln(
				jen.ID("Name").Op(":").ID("form").Dot("Get").Call(jen.IDf("itemCreationInputNameFormKey")), jen.ID("Details").Op(":").ID("form").Dot("Get").Call(jen.ID("itemCreationInputDetailsFormKey")), jen.ID("BelongsToAccount").Op(":").ID("sessionCtxData").Dot("ActiveAccountID")),
			jen.If(jen.Err().Equals().ID("updateInput").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.Err().DoesNotEqual().ID("nil")).Body(
				constants.LoggerVar().Equals().ID("logger").Dot("WithValue").Call(
					jen.Lit("input"),
					jen.ID("updateInput"),
				),
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("invalid item creation input"),
				),
				jen.Return().ID("nil"),
			),
			jen.Return().ID("updateInput"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("handleItemUpdateRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("sessionCtxData"), jen.Err()).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
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
			jen.ID("updateInput").Assign().ID("s").Dotf("parseFormEncodedItemUpdateInput").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("sessionCtxData"),
			),
			jen.If(jen.ID("updateInput").Op("==").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Lit("no update input attached to request"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
				jen.Return(),
			),
			jen.List(jen.ID("item"), jen.Err()).Assign().ID("s").Dotf("fetchItem").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData"),
				jen.ID("req"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("fetching item from datastore"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.ID("changes").Assign().ID("item").Dot("Update").Call(jen.ID("updateInput")),
			jen.If(jen.Err().Equals().ID("s").Dot("dataStore").Dotf("UpdateItem").Call(
				jen.ID("ctx"),
				jen.ID("item"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
				jen.ID("changes"),
			), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("fetching item from datastore"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.ID("tmplFuncMap").Assign().Map(jen.String()).Interface().Valuesln(
				jen.Lit("componentTitle").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot(sn)).Params(jen.String()).Body(
					jen.Return().Qual("fmt", "Sprintf").Call(
						jen.Lit("Item #%d"),
						jen.ID("x").Dot("ID"),
					))),
			jen.ID("tmpl").Assign().ID("s").Dot("parseTemplate").Call(
				jen.ID("ctx"),
				jen.Lit(""),
				jen.IDf("itemEditorTemplate"),
				jen.ID("tmplFuncMap"),
			),
			jen.ID("s").Dot("renderTemplateToResponse").Call(
				jen.ID("ctx"),
				jen.ID("tmpl"),
				jen.ID("item"),
				jen.ID("res"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).IDf("handleItemDeletionRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			constants.LoggerVar().Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("sessionCtxData"), jen.Err()).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
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
			jen.IDf("itemID").Assign().ID("s").Dot("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
				constants.LoggerVar(),
				jen.IDf("itemIDURLParamKey"),
				jen.Lit("item"),
			).Call(jen.ID("req")),
			jen.If(jen.Err().Equals().ID("s").Dot("dataStore").Dotf("ArchiveItem").Call(
				jen.ID("ctx"),
				jen.IDf("itemID"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("archiving items in datastore"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.List(jen.ID("items"), jen.Err()).Assign().ID("s").Dotf("fetchItems").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData"),
				jen.ID("req"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("fetching items from datastore"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.ID("tmplFuncMap").Assign().Map(jen.String()).Interface().Valuesln(
				jen.Lit("individualURL").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot(sn)).Params(jen.Qual("html/template", "URL")).Body(
					jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
						jen.Litf("/dashboard_pages/items/%d"),
						jen.ID("x").Dot("ID"),
					))), jen.Lit("pushURL").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot(sn)).Params(jen.Qual("html/template", "URL")).Body(
					jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
						jen.Litf("/items/%d"),
						jen.ID("x").Dot("ID"),
					)))),
			jen.ID("tmpl").Assign().ID("s").Dot("parseTemplate").Call(
				jen.ID("ctx"),
				jen.Lit("dashboard"),
				jen.IDf("itemsTableTemplate"),
				jen.ID("tmplFuncMap"),
			),
			jen.ID("s").Dot("renderTemplateToResponse").Call(
				jen.ID("ctx"),
				jen.ID("tmpl"),
				jen.IDf("items"),
				jen.ID("res"),
			),
		),
		jen.Newline(),
	)

	return code
}
