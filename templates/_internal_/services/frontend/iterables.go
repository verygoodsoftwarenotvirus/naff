package frontend

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildIDFetchers(proj *models.Project, typ models.DataType, includePrimaryType bool) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()

	idFetches := []jen.Code{}
	for _, dep := range proj.FindOwnerTypeChain(typ) {
		tsn := dep.Name.Singular()
		tuvn := dep.Name.UnexportedVarName()

		idFetches = append(idFetches,
			jen.Commentf("determine %s ID.", dep.Name.SingularCommonName()),
			jen.IDf("%sID", tuvn).Assign().ID("s").Dotf("%sIDFetcher", tuvn).Call(jen.ID(constants.RequestVarName)),
			jen.Qualf(proj.InternalTracingPackage(), "Attach%sIDToSpan", tsn).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", tuvn)),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qualf(proj.ConstantKeysPackage(), "%sIDKey", tsn), jen.IDf("%sID", tuvn)),
			jen.Newline(),
		)
	}

	if includePrimaryType {
		idFetches = append(idFetches,
			jen.Commentf("determine %s ID.", scn),
			jen.IDf("%sID", uvn).Assign().ID("s").Dotf("%sIDFetcher", uvn).Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), fmt.Sprintf("%sIDKey", sn)),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
		)
	}

	return idFetches
}

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, true)

	uvn := typ.Name.UnexportedVarName()
	puvn := typ.Name.PluralUnexportedVarName()
	rn := typ.Name.RouteName()
	prn := typ.Name.PluralRouteName()

	code.Add(
		jen.Const().Defs(jen.IDf("%sIDURLParamKey", uvn).Equals().Lit(rn)),
		jen.Newline(),
	)

	code.Add(buildFetchSomething(proj, typ)...)

	code.Add(
		jen.Commentf("//go:embed templates/partials/generated/creators/%s_creator.gotpl", rn),
		jen.Newline(),
		jen.Var().IDf("%sCreatorTemplate", uvn).String(),
		jen.Newline(),
	)

	code.Add(buildBuildSomethingCreatorView(proj, typ)...)

	code.Add(
		jen.Const().Defs(buildKeys(proj, typ)...),
		jen.Newline(),
	)

	code.Add(buildParseFormEncodedSomethingCreationInput(proj, typ)...)
	code.Add(buildHandleSomethingCreationRequest(proj, typ)...)

	code.Add(
		jen.Commentf("//go:embed templates/partials/generated/editors/%s_editor.gotpl", rn),
		jen.Newline(),
		jen.Var().IDf("%sEditorTemplate", uvn).String(),
		jen.Newline(),
	)

	code.Add(buildBuildSomethingEditorView(proj, typ)...)
	code.Add(buildFetchSomethings(proj, typ)...)

	code.Add(
		jen.Commentf("//go:embed templates/partials/generated/tables/%s_table.gotpl", prn),
		jen.Newline(),
		jen.Var().IDf("%sTableTemplate", puvn).String(),
		jen.Newline(),
	)

	code.Add(buildBuildSomethingTableView(proj, typ)...)
	code.Add(buildParseFormEncodedSomethingUpdateInput(proj, typ)...)
	code.Add(buildHandleSomethingUpdateRequest(proj, typ)...)
	code.Add(buildHandleSomethingDeletionRequest(proj, typ)...)

	return code
}

func buildDBClientRetrievalMethodCallArgs(proj *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}
	uvn := typ.Name.UnexportedVarName()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", uvn))

	if typ.RestrictedToAccountAtSomeLevel(proj) {
		params = append(params, jen.ID("sessionCtxData").Dot("ActiveAccountID"))
	}

	return params
}

func buildFetchSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()

	callArgs := buildDBClientRetrievalMethodCallArgs(proj, typ)

	elseBodyLines := buildIDFetchers(proj, typ, true)

	elseBodyLines = append(elseBodyLines,
		jen.List(jen.IDf(uvn), jen.Err()).Equals().ID("s").Dot("dataStore").Dotf("Get%s", sn).Call(
			callArgs...,
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.Err(),
				constants.LoggerVar(),
				jen.ID("span"),
				jen.Litf("fetching %s data", scn),
			)),
		),
	)

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).IDf("fetch%s", sn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("req").PointerTo().Qual("net/http", "Request"),
			utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("sessionCtxData").PointerTo().Qual(proj.TypesPackage(), "SessionContextData")),
		).Params(jen.ID(uvn).PointerTo().Qual(proj.TypesPackage(), sn), jen.Err().ID("error")).Body(
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
				jen.IDf(uvn).Equals().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
			).Else().Body(
				elseBodyLines...,
			),
			jen.Newline(),
			jen.Return().List(jen.ID(uvn), jen.ID("nil")),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildSomethingCreatorView(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).IDf("build%sCreatorView", sn).Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Func().Params(jen.Qual("net/http", "ResponseWriter"), jen.PointerTo().Qual("net/http", "Request"))).Body(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Body(
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
				jen.ID(uvn).Assign().AddressOf().Qual(proj.TypesPackage(), sn).Values(),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("view").Assign().ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.IDf("%sCreatorTemplate", uvn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("page").Assign().AddressOf().ID("pageData").Valuesln(
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
	}

	return lines
}

func buildKeys(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()

	allFieldsKeys := []jen.Code{}
	creationKeys := []jen.Code{}
	updateKeys := []jen.Code{}

	for _, field := range typ.Fields {
		fuvn := field.Name.UnexportedVarName()
		fsn := field.Name.Singular()

		var alreadyCreated bool
		for _, parent := range proj.FindOwnerTypeChain(typ) {
			for _, f := range parent.Fields {
				if f.Name.UnexportedVarName() == fuvn {
					alreadyCreated = true
				}
			}
		}

		if !alreadyCreated {
			allFieldsKeys = append(allFieldsKeys, jen.IDf("%sFormKey", fuvn).Equals().Lit(fuvn))
		}

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

	return allKeys
}

func determineFormFetcher(field models.DataField) string {
	switch field.Type {
	case "string":
		if field.IsPointer {
			return "stringToPointerToString"
		}
		fallthrough
	case "bool":
		if field.IsPointer {
			return "stringToPointerToBool"
		}
		return "stringToBool"
	case "int":
		if field.IsPointer {
			return "stringToPointerToInt"
		}
		return "stringToInt"
	case "int8":
		if field.IsPointer {
			return "stringToPointerToInt8"
		}
		return "stringToInt8"
	case "int16":
		if field.IsPointer {
			return "stringToPointerToInt16"
		}
		return "stringToInt16"
	case "int32":
		if field.IsPointer {
			return "stringToPointerToInt32"
		}
		return "stringToInt32"
	case "int64":
		if field.IsPointer {
			return "stringToPointerToInt64"
		}
		return "stringToInt64"
	case "uint":
		if field.IsPointer {
			return "stringToPointerToUint"
		}
		return "stringToUint"
	case "uint8":
		if field.IsPointer {
			return "stringToPointerToUint8"
		}
		return "stringToUint8"
	case "uint16":
		if field.IsPointer {
			return "stringToPointerToUint16"
		}
		return "stringToUint16"
	case "uint32":
		if field.IsPointer {
			return "stringToPointerToUint32"
		}
		return "stringToUint32"
	case "uint64":
		if field.IsPointer {
			return "stringToPointerToUint64"
		}
		return "stringToUint64"
	case "float32":
		if field.IsPointer {
			return "stringToPointerToFloat32"
		}
		return "stringToFloat32"
	case "float64":
		if field.IsPointer {
			return "stringToPointerToFloat64"
		}
		return "stringToFloat64"
	default:
		panic(fmt.Sprintf("invalid type: %q", field.Type))
	}
}

func buildParseFormEncodedSomethingCreationInput(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()

	creationInputFields := []jen.Code{}
	for _, field := range typ.Fields {
		fsn := field.Name.Singular()
		if field.Type != "string" || (field.Type == "string" && field.IsPointer) {
			creationInputFields = append(creationInputFields, jen.ID(fsn).Op(":").ID("s").Dot(determineFormFetcher(field)).Call(jen.ID("form"), jen.IDf("%sCreationInput%sFormKey", uvn, fsn)))
		} else {
			creationInputFields = append(creationInputFields, jen.ID(fsn).Op(":").ID("form").Dot("Get").Call(jen.IDf("%sCreationInput%sFormKey", uvn, fsn)))
		}
	}

	if typ.BelongsToAccount {
		creationInputFields = append(creationInputFields, jen.ID("BelongsToAccount").Op(":").ID("sessionCtxData").Dot("ActiveAccountID"))
	}

	lines := []jen.Code{
		jen.Commentf("parseFormEncoded%sCreationInput checks a request for an %sCreationInput.", sn, sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).IDf("parseFormEncoded%sCreationInput", sn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("req").PointerTo().Qual("net/http", "Request"),
			func() jen.Code {
				if typ.BelongsToAccount {
					return jen.ID("sessionCtxData").PointerTo().Qual(proj.TypesPackage(), "SessionContextData")
				}
				return jen.Null()
			}(),
		).Params(jen.ID("creationInput").PointerTo().ID("types").Dotf("%sCreationInput", sn)).Body(
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
			jen.ID("creationInput").Equals().AddressOf().ID("types").Dotf("%sCreationInput", sn).Valuesln(
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
	}

	return lines
}

func buildHandleSomethingCreationRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	prn := typ.Name.PluralRouteName()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).IDf("handle%sCreationRequest", sn).Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Body(
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
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.Commentf("determine %s ID.", typ.BelongsToStruct.SingularCommonName()).Newline().
						IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Assign().ID("s").Dotf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Call(jen.ID("req")).Newline().
						Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", typ.BelongsToStruct.Singular())).Call(
						jen.ID("span"),
						jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()),
					).Newline().
						ID("logger").Equals().ID("logger").Dot("WithValue").Call(
						jen.Qual(proj.ObservabilityPackage("keys"), fmt.Sprintf("%sIDKey", typ.BelongsToStruct.Singular())),
						jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()),
					).Newline().Newline()
				}
				return jen.Null()
			}(),
			jen.ID("creationInput").Assign().ID("s").Dotf("parseFormEncoded%sCreationInput", sn).Call(
				jen.ID("ctx"),
				jen.ID("req"),
				func() jen.Code {
					if typ.BelongsToAccount {
						return jen.ID("sessionCtxData")
					}
					return jen.Null()
				}(),
			),
			jen.If(jen.ID("creationInput").IsEqualTo().ID("nil")).Body(
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
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.ID("creationInput").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Newline()
				}
				return jen.Null()
			}(),
			constants.LoggerVar().Dot("Debug").Call(jen.Litf("%s creation input parsed successfully", scn)),
			jen.Newline(),
			jen.If(jen.List(jen.Underscore(), jen.Err()).Equals().ID("s").Dot("dataStore").Dotf("Create%s", sn).Call(
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
	}

	return lines
}

func buildBuildSomethingEditorView(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).IDf("build%sEditorView", sn).Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Func().Params(jen.Qual("net/http", "ResponseWriter"), jen.PointerTo().Qual("net/http", "Request"))).Body(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Body(
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
					jen.ID("req"),
					utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("sessionCtxData")),
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
					jen.Lit("componentTitle").Op(":").Func().Params(jen.ID("x").PointerTo().Qual(proj.TypesPackage(), sn)).Params(jen.String()).Body(
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
					jen.ID("page").Assign().AddressOf().ID("pageData").Valuesln(
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
	}

	return lines
}

func buildDBClientListRetrievalMethodCallArgs(p *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.RestrictedToAccountAtSomeLevel(p) {
		params = append(params, jen.ID("sessionCtxData").Dot("ActiveAccountID"))
	}
	params = append(params, jen.ID("filter"))

	return params
}

func buildFetchSomethings(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	scn := typ.Name.SingularCommonName()

	callArgs := buildDBClientListRetrievalMethodCallArgs(proj, typ)

	elseBody := buildIDFetchers(proj, typ, false)
	elseBody = append(elseBody,
		jen.ID("filter").Assign().Qual(proj.TypesPackage(), "ExtractQueryFilter").Call(jen.ID("req")),
		jen.Qual(proj.InternalTracingPackage(), "AttachQueryFilterToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.FilterVarName)),
		jen.Newline(),
		jen.List(jen.ID(puvn), jen.Err()).Equals().ID("s").Dot("dataStore").Dotf("Get%s", pn).Call(
			callArgs...,
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.Err(),
				constants.LoggerVar(),
				jen.ID("span"),
				jen.Litf("fetching %s data", scn),
			)),
		),
	)

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).IDf("fetch%s", pn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("req").PointerTo().Qual("net/http", "Request"),
			utils.ConditionalCode(typ.RestrictedToAccountMembers, jen.ID("sessionCtxData").PointerTo().Qual(proj.TypesPackage(), "SessionContextData")),
		).Params(jen.ID(puvn).PointerTo().ID("types").Dotf("%sList", sn), jen.Err().ID("error")).Body(
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
				jen.ID(puvn).Equals().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
			).Else().Body(
				elseBody...,
			),
			jen.Newline(),
			jen.Return().List(jen.ID(puvn), jen.ID("nil")),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildSomethingTableView(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	prn := typ.Name.PluralRouteName()
	pcn := typ.Name.PluralCommonName()

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).IDf("build%sTableView", pn).Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Func().Params(jen.Qual("net/http", "ResponseWriter"), jen.PointerTo().Qual("net/http", "Request"))).Body(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Body(
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
					jen.ID("req"),
					utils.ConditionalCode(typ.RestrictedToAccountMembers, jen.ID("sessionCtxData")),
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
					jen.Lit("individualURL").Op(":").Func().Params(jen.ID("x").PointerTo().Qual(proj.TypesPackage(), sn)).Params(jen.Qual("html/template", "URL")).Body(
						jen.Comment("#nosec G203"),
						jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
							jen.Lit(fmt.Sprintf("/dashboard_pages/%s/", prn)+"%d"),
							jen.ID("x").Dot("ID"),
						)),
					), jen.Lit("pushURL").Op(":").Func().Params(jen.ID("x").PointerTo().Qual(proj.TypesPackage(), sn)).Params(jen.Qual("html/template", "URL")).Body(
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
					jen.ID("page").Assign().AddressOf().ID("pageData").Valuesln(
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
	}

	return lines
}

func buildParseFormEncodedSomethingUpdateInput(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()

	updateInputFields := []jen.Code{}
	for _, field := range typ.Fields {
		fsn := field.Name.Singular()
		if field.Type != "string" || (field.Type == "string" && field.IsPointer) {
			updateInputFields = append(updateInputFields, jen.ID(fsn).Op(":").ID("s").Dot(determineFormFetcher(field)).Call(jen.ID("form"), jen.IDf("%sUpdateInput%sFormKey", uvn, fsn)))
		} else {
			updateInputFields = append(updateInputFields, jen.ID(fsn).Op(":").ID("form").Dot("Get").Call(jen.IDf("%sUpdateInput%sFormKey", uvn, fsn)))
		}
	}

	if typ.BelongsToAccount {
		updateInputFields = append(updateInputFields, jen.ID("BelongsToAccount").Op(":").ID("sessionCtxData").Dot("ActiveAccountID"))
	}

	lines := []jen.Code{
		jen.Commentf("parseFormEncoded%sUpdateInput checks a request for an %sUpdateInput.", sn, sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).IDf("parseFormEncoded%sUpdateInput", sn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("req").PointerTo().Qual("net/http", "Request"),
			utils.ConditionalCode(typ.BelongsToAccount, jen.ID("sessionCtxData").PointerTo().Qual(proj.TypesPackage(), "SessionContextData")),
		).Params(jen.ID("updateInput").PointerTo().ID("types").Dotf("%sUpdateInput", sn)).Body(
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
			jen.ID("updateInput").Equals().AddressOf().ID("types").Dotf("%sUpdateInput", sn).Valuesln(
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
	}

	return lines
}

func buildHandleSomethingUpdateRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).IDf("handle%sUpdateRequest", sn).Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Body(
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
				utils.ConditionalCode(typ.BelongsToAccount, jen.ID("sessionCtxData")),
			),
			jen.If(jen.ID("updateInput").IsEqualTo().ID("nil")).Body(
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
				jen.ID("req"),
				utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("sessionCtxData")),
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
				jen.Lit("componentTitle").Op(":").Func().Params(jen.ID("x").PointerTo().Qual(proj.TypesPackage(), sn)).Params(jen.String()).Body(
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
	}

	return lines
}

func buildDBClientDeletionMethodCallArgs(proj *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}
	uvn := typ.Name.UnexportedVarName()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", uvn))

	if typ.RestrictedToAccountAtSomeLevel(proj) {
		params = append(params, jen.ID("sessionCtxData").Dot("ActiveAccountID"))
	}

	return params
}

func buildHandleSomethingDeletionRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()
	puvn := typ.Name.PluralUnexportedVarName()
	prn := typ.Name.PluralRouteName()
	pcn := typ.Name.PluralCommonName()

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).IDf("handle%sArchiveRequest", sn).Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Body(
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
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					tuvn := typ.BelongsToStruct.UnexportedVarName()
					tsn := typ.BelongsToStruct.Singular()
					return jen.IDf("%sID", tuvn).Assign().ID("s").Dotf("%sIDFetcher", tuvn).Call(jen.ID(constants.RequestVarName)).Newline().
						Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", tsn)).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", tuvn)).Newline().
						ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", tsn)), jen.IDf("%sID", tuvn)).Newline()
				}
				return jen.Null()
			}(),
			jen.IDf("%sID", uvn).Assign().ID("s").Dotf("%sIDFetcher", uvn).Call(jen.ID(constants.RequestVarName)),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", uvn)),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", sn)), jen.IDf("%sID", uvn)),
			jen.Newline(),
			jen.If(jen.Err().Equals().ID("s").Dot("dataStore").Dotf("Archive%s", sn).Call(
				jen.ID("ctx"),
				func() jen.Code {
					if typ.BelongsToStruct != nil {
						return jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())
					}
					return jen.Null()
				}(),
				jen.IDf("%sID", uvn),
				func() jen.Code {
					if typ.BelongsToAccount {
						return jen.ID("sessionCtxData").Dot("ActiveAccountID")
					}
					return jen.Null()
				}(),
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
				jen.ID("req"),
				utils.ConditionalCode(typ.RestrictedToAccountMembers, jen.ID("sessionCtxData")),
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
				jen.Lit("individualURL").Op(":").Func().Params(jen.ID("x").PointerTo().Qual(proj.TypesPackage(), sn)).Params(jen.Qual("html/template", "URL")).Body(
					jen.Comment("#nosec G203"),
					jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
						jen.Lit(fmt.Sprintf("/dashboard_pages/%s/", prn)+"%d"),
						jen.ID("x").Dot("ID"),
					),
					),
				), jen.Lit("pushURL").Op(":").Func().Params(jen.ID("x").PointerTo().Qual(proj.TypesPackage(), sn)).Params(jen.Qual("html/template", "URL")).Body(
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
	}

	return lines
}
