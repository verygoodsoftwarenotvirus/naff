package client

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)
	code.Add(jen.Const().Defs(
		jen.IDf("%sBasePath", typ.Name.PluralUnexportedVarName()).Equals().Lit(typ.Name.PluralRouteName())),
	)

	code.Add(buildBuildSomethingExistsRequest(proj, typ)...)
	code.Add(buildSomethingExists(proj, typ)...)
	code.Add(buildBuildGetSomethingRequestFuncDecl(proj, typ)...)
	code.Add(buildGetSomethingFuncDecl(proj, typ)...)

	if typ.SearchEnabled {
		code.Add(buildBuildSearchSomethingRequestFuncDecl(proj, typ)...)
		code.Add(buildSearchSomethingFuncDecl(proj, typ)...)
	}

	code.Add(buildBuildGetListOfSomethingRequestFuncDecl(proj, typ)...)
	code.Add(buildGetListOfSomethingFuncDecl(proj, typ)...)
	code.Add(buildBuildCreateSomethingRequestFuncDecl(proj, typ)...)
	code.Add(buildCreateSomethingFuncDecl(proj, typ)...)
	code.Add(buildBuildUpdateSomethingRequestFuncDecl(proj, typ)...)
	code.Add(buildUpdateSomethingFuncDecl(proj, typ)...)
	code.Add(buildBuildArchiveSomethingRequestFuncDecl(proj, typ)...)
	code.Add(buildArchiveSomethingFuncDecl(proj, typ)...)

	return code
}

func attachURIToSpanCall(proj *models.Project) jen.Code {
	return jen.Qual(proj.InternalTracingV1Package(), "AttachRequestURIToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("uri"))
}

func buildV1ClientURLBuildingParamsForSingleInstanceOfSomething(proj *models.Project, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{
		jen.Nil(),
	}

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		urlBuildingParams = append(urlBuildingParams,
			jen.IDf("%sBasePath", pt.Name.PluralUnexportedVarName()),
			jen.Qual("strconv", "FormatUint").Call(
				jen.IDf("%sID", pt.Name.UnexportedVarName()),
				jen.Lit(10),
			),
		)
	}
	urlBuildingParams = append(urlBuildingParams,
		jen.ID(basePath),
		jen.Qual("strconv", "FormatUint").Call(
			jen.IDf("%sID", typ.Name.UnexportedVarName()),
			jen.Lit(10),
		),
	)

	return urlBuildingParams
}

func buildV1ClientURLBuildingParamsForCreatingSomething(proj *models.Project, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{
		jen.Nil(),
	}

	parents := proj.FindOwnerTypeChain(typ)
	for i, pt := range parents {
		urlBuildingParams = append(urlBuildingParams,
			jen.IDf("%sBasePath", pt.Name.PluralUnexportedVarName()),
			jen.Qual("strconv", "FormatUint").Call(
				func() jen.Code {
					if i == len(parents)-1 && typ.BelongsToStruct != nil {
						return jen.ID("input").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular())
					}
					return jen.IDf("%sID", pt.Name.UnexportedVarName())
				}(),
				jen.Lit(10),
			),
		)
	}
	urlBuildingParams = append(urlBuildingParams,
		jen.ID(basePath),
	)

	return urlBuildingParams
}

func buildV1ClientURLBuildingParamsForListOfSomething(proj *models.Project, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{
		jen.ID(constants.FilterVarName).Dot("ToValues").Call(),
	}

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		urlBuildingParams = append(urlBuildingParams,
			jen.IDf("%sBasePath", pt.Name.PluralUnexportedVarName()),
			jen.Qual("strconv", "FormatUint").Call(
				jen.IDf("%sID", pt.Name.UnexportedVarName()),
				jen.Lit(10),
			),
		)
	}
	urlBuildingParams = append(urlBuildingParams,
		jen.ID(basePath),
	)

	return urlBuildingParams
}

func buildV1ClientURLBuildingParamsForSearchingSomething(typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{
		jen.ID("params"),
		jen.ID(basePath),
		jen.Lit("search"),
	}

	return urlBuildingParams
}

func buildV1ClientURLBuildingParamsForMethodThatIncludesItsOwnType(proj *models.Project, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{jen.Nil()}

	ownerChain := proj.FindOwnerTypeChain(typ)
	for i, pt := range ownerChain {
		if i == len(ownerChain)-1 {
			urlBuildingParams = append(urlBuildingParams,
				jen.IDf("%sBasePath", pt.Name.PluralUnexportedVarName()),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID(typ.Name.UnexportedVarName()).Dotf("BelongsTo%s", pt.Name.Singular()),
					jen.Lit(10),
				),
			)
		} else {
			urlBuildingParams = append(urlBuildingParams,
				jen.IDf("%sBasePath", pt.Name.PluralUnexportedVarName()),
				jen.Qual("strconv", "FormatUint").Call(
					jen.IDf("%sID", pt.Name.UnexportedVarName()),
					jen.Lit(10),
				),
			)
		}
	}

	urlBuildingParams = append(urlBuildingParams,
		jen.ID(basePath),
		jen.Qual("strconv", "FormatUint").Call(
			jen.ID(typ.Name.UnexportedVarName()).Dot("ID"),
			jen.Lit(10),
		),
	)

	return urlBuildingParams
}

func buildBuildSomethingExistsRequest(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()
	funcName := fmt.Sprintf("Build%sExistsRequest", ts)

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Callln(
			buildV1ClientURLBuildingParamsForSingleInstanceOfSomething(proj, typ)...,
		),
		attachURIToSpanCall(proj),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodHead"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for checking the existence of %s.", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(typ.BuildParamsForHTTPClientExistenceRequestBuildingMethod(proj)...).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	funcName := fmt.Sprintf("%sExists", ts)
	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().ID("c").Dotf("Build%sExistsRequest", ts).Call(typ.BuildArgsForHTTPClientExistenceRequestBuildingMethod(proj)...),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(
				jen.False(),
				jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Err()),
			),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("checkExistence").Call(constants.CtxVar(), jen.ID(constants.RequestVarName)),
	}

	lines := []jen.Code{
		jen.Commentf("%s retrieves whether or not %s exists.", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(typ.BuildParamsForHTTPClientExistenceMethod(proj)...).Params(
			jen.ID("exists").Bool(),
			jen.Err().Error(),
		).Body(block...,
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetSomethingRequestFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()
	funcName := fmt.Sprintf("BuildGet%sRequest", ts)

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Callln(
			buildV1ClientURLBuildingParamsForSingleInstanceOfSomething(proj, typ)...,
		),
		attachURIToSpanCall(proj),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for fetching %s.", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(typ.BuildParamsForHTTPClientRetrievalRequestBuildingMethod(proj)...).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildGetSomethingFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	funcName := fmt.Sprintf("Get%s", ts)
	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", ts)).Call(typ.BuildArgsForHTTPClientRetrievalRequestBuildingMethod(proj)...),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Err()),
			),
		),
		jen.Line(),
		jen.If(jen.ID("retrieveErr").Assign().ID("c").Dot("retrieve").Call(constants.CtxVar(), jen.ID(constants.RequestVarName), jen.AddressOf().ID(uvn)), jen.ID("retrieveErr").DoesNotEqual().ID("nil")).Body(
			jen.Return().List(jen.Nil(), jen.ID("retrieveErr")),
		),
		jen.Line(),
		jen.Return().List(jen.ID(uvn), jen.Nil()),
	}

	lines := []jen.Code{
		jen.Commentf("%s retrieves %s.", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(typ.BuildParamsForHTTPClientRetrievalMethod(proj, false)...).Params(
			jen.ID(uvn).PointerTo().Qual(proj.ModelsV1Package(), ts),
			jen.Err().Error(),
		).Body(block...,
		),
		jen.Line(),
	}

	return lines
}

func buildBuildSearchSomethingRequestFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Plural()
	commonNameWithPrefix := typ.Name.PluralCommonName()
	funcName := fmt.Sprintf("BuildSearch%sRequest", ts)

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("params").Assign().Qual("net/url", "Values").Values(),
		jen.ID("params").Dot("Set").Call(
			jen.Qual(proj.ModelsV1Package(), "SearchQueryKey"),
			jen.ID("query"),
		),
		jen.ID("params").Dot("Set").Call(
			jen.Qual(proj.ModelsV1Package(), "LimitQueryKey"),
			jen.Qual("strconv", "FormatUint").Call(
				jen.Uint64().Call(jen.ID("limit")),
				jen.Lit(10),
			),
		),
		jen.Line(),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Callln(
			buildV1ClientURLBuildingParamsForSearchingSomething(typ)...,
		),
		attachURIToSpanCall(proj),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for querying %s.", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID("query").String(),
			jen.ID("limit").Uint8(),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildSearchSomethingFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	tp := typ.Name.Plural()
	ts := typ.Name.Singular()
	puvn := typ.Name.PluralUnexportedVarName()
	funcName := fmt.Sprintf("Search%s", tp)
	commonName := typ.Name.PluralCommonName()

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("BuildSearch%sRequest", tp)).Call(
			constants.CtxVar(),
			jen.ID("query"),
			jen.ID("limit"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Err()),
			),
		),
		jen.Line(),
		jen.If(jen.ID("retrieveErr").Assign().ID("c").Dot("retrieve").Call(constants.CtxVar(), jen.ID(constants.RequestVarName), jen.AddressOf().ID(puvn)), jen.ID("retrieveErr").DoesNotEqual().ID("nil")).Body(
			jen.Return().List(jen.Nil(), jen.ID("retrieveErr")),
		),
		jen.Line(),
		jen.Return().List(jen.ID(puvn), jen.Nil()),
	}

	lines := []jen.Code{
		jen.Commentf("%s searches for a list of %s.", funcName, commonName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID("query").String(),
			jen.ID("limit").Uint8(),
		).Params(
			jen.ID(puvn).Index().Qual(proj.ModelsV1Package(), ts),
			jen.Err().Error(),
		).Body(
			block...,
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetListOfSomethingRequestFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	tp := typ.Name.Plural() // title plural

	funcName := fmt.Sprintf("BuildGet%sRequest", tp)

	urlBuildingParams := buildV1ClientURLBuildingParamsForListOfSomething(proj, typ)

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Callln(urlBuildingParams...),
		attachURIToSpanCall(proj),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for fetching %s.", funcName, typ.Name.PluralCommonName()),
		jen.Line(),
		newClientMethod(funcName).Params(typ.BuildParamsForHTTPClientListRequestMethod(proj)...).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildGetListOfSomethingFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	tp := typ.Name.Plural() // title plural
	pvn := typ.Name.PluralUnexportedVarName()
	funcName := fmt.Sprintf("Get%s", tp)

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", tp)).Call(
			typ.BuildArgsForHTTPClientListRequestMethod(proj)...,
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.If(
			jen.ID("retrieveErr").Assign().ID("c").Dot("retrieve").Call(constants.CtxVar(), jen.ID(constants.RequestVarName), jen.AddressOf().ID(pvn)),
			jen.ID("retrieveErr").DoesNotEqual().ID("nil"),
		).Body(
			jen.Return().List(jen.Nil(), jen.ID("retrieveErr")),
		),
		jen.Line(),
		jen.Return().List(jen.ID(pvn), jen.Nil()),
	}

	lines := []jen.Code{
		jen.Commentf("%s retrieves a list of %s.", funcName, typ.Name.PluralCommonName()),
		jen.Line(),
		newClientMethod(funcName).Params(typ.BuildParamsForHTTPClientMethodThatFetchesAList(proj)...).Params(
			jen.ID(pvn).PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", ts)),
			jen.Err().Error(),
		).Body(block...,
		),
		jen.Line(),
	}

	return lines
}

func buildBuildCreateSomethingRequestFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()
	funcName := fmt.Sprintf("BuildCreate%sRequest", ts)

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Callln(
			buildV1ClientURLBuildingParamsForCreatingSomething(proj, typ)...,
		),
		attachURIToSpanCall(proj),
		jen.Line(),
		jen.Return().ID("c").Dot("buildDataRequest").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.ID("input"),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for creating %s.", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(
			typ.BuildParamsForHTTPClientCreateRequestBuildingMethod(proj)...,
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(
			block...,
		),
		jen.Line(),
	}

	return lines
}

func buildCreateSomethingFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	vn := typ.Name.UnexportedVarName()
	funcName := fmt.Sprintf("Create%s", ts)
	commonName := strings.Join(strings.Split(typ.Name.RouteName(), "_"), " ")
	commonNameWithPrefix := fmt.Sprintf("%s %s", wordsmith.AOrAn(ts), commonName)

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().ID("c").Dot(fmt.Sprintf("BuildCreate%sRequest", ts)).Call(
			typ.BuildArgsForHTTPClientCreateRequestBuildingMethod(proj)...,
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Err().Equals().ID("c").Dot("executeRequest").Call(
			constants.CtxVar(),
			jen.ID(constants.RequestVarName),
			jen.AddressOf().ID(vn),
		),
		jen.Return().List(
			jen.ID(vn),
			jen.Err(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s creates %s.", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(
			typ.BuildParamsForHTTPClientCreateMethod(proj)...,
		).Params(
			jen.ID(vn).PointerTo().Qual(proj.ModelsV1Package(), ts),
			jen.Err().Error(),
		).Body(block...,
		),
		jen.Line(),
	}

	return lines
}

func buildBuildUpdateSomethingRequestFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()
	funcName := fmt.Sprintf("BuildUpdate%sRequest", ts)

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Callln(
			buildV1ClientURLBuildingParamsForMethodThatIncludesItsOwnType(proj, typ)...,
		),
		attachURIToSpanCall(proj),
		jen.Line(),
		jen.Return().ID("c").Dot("buildDataRequest").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodPut"),
			jen.ID("uri"),
			jen.ID(typ.Name.UnexportedVarName()),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for updating %s.", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(
			typ.BuildParamsForHTTPClientUpdateRequestBuildingMethod(proj)...,
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...,
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateSomethingFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()

	funcName := fmt.Sprintf("Update%s", ts)

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().ID("c").Dot(fmt.Sprintf("BuildUpdate%sRequest", ts)).Call(
			typ.BuildArgsForHTTPClientUpdateRequestBuildingMethod(proj)...,
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("building request: %w"),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("executeRequest").Call(
			constants.CtxVar(),
			jen.ID(constants.RequestVarName),
			jen.AddressOf().ID(typ.Name.UnexportedVarName()),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s updates %s.", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(
			typ.BuildParamsForHTTPClientUpdateMethod(proj)...,
		).Params(jen.Error()).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildArchiveSomethingRequestFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()

	commonName := strings.Join(strings.Split(typ.Name.RouteName(), "_"), " ")
	commonNameWithPrefix := fmt.Sprintf("%s %s", wordsmith.AOrAn(ts), commonName)

	funcName := fmt.Sprintf("BuildArchive%sRequest", ts)

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Callln(
			buildV1ClientURLBuildingParamsForSingleInstanceOfSomething(proj, typ)...,
		),
		attachURIToSpanCall(proj),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodDelete"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for updating %s.", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(typ.BuildParamsForHTTPClientArchiveRequestBuildingMethod(proj)...).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildArchiveSomethingFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	funcName := fmt.Sprintf("Archive%s", ts)
	commonName := strings.Join(strings.Split(typ.Name.RouteName(), "_"), " ")
	commonNameWithPrefix := fmt.Sprintf("%s %s", wordsmith.AOrAn(ts), commonName)

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().ID("c").Dot(fmt.Sprintf("BuildArchive%sRequest", ts)).Call(
			typ.BuildArgsForHTTPClientArchiveRequestBuildingMethod(proj)...,
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("building request: %w"),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("executeRequest").Call(
			constants.CtxVar(),
			jen.ID(constants.RequestVarName),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s archives %s.", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(typ.BuildParamsForHTTPClientArchiveMethod(proj)...).Params(jen.Error()).Body(block...),
	}

	return lines
}
