package client

import (
	"fmt"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(packageName)

	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())

	utils.AddImports(proj, ret)
	ret.Add(jen.Const().Defs(
		jen.ID(basePath).Equals().Lit(typ.Name.PluralRouteName())),
	)

	ret.Add(buildBuildItemExistsRequest(proj, typ)...)
	ret.Add(buildItemExists(proj, typ)...)
	ret.Add(buildBuildGetSomethingRequestFuncDecl(proj, typ)...)
	ret.Add(buildGetSomethingFuncDecl(proj, typ)...)
	ret.Add(buildBuildGetListOfSomethingRequestFuncDecl(proj, typ)...)
	ret.Add(buildGetListOfSomethingFuncDecl(proj, typ)...)
	ret.Add(buildBuildCreateSomethingRequestFuncDecl(proj, typ)...)
	ret.Add(buildCreateSomethingFuncDecl(proj, typ)...)
	ret.Add(buildBuildUpdateSomethingRequestFuncDecl(proj, typ)...)
	ret.Add(buildUpdateSomethingFuncDecl(proj, typ)...)
	ret.Add(buildBuildArchiveSomethingRequestFuncDecl(proj, typ)...)
	ret.Add(buildArchiveSomethingFuncDecl(proj, typ)...)

	return ret
}

func attachURIToSpanCall(proj *models.Project) jen.Code {
	return jen.Qual(proj.InternalTracingV1Package(), "AttachRequestURIToSpan").Call(jen.ID("span"), jen.ID("uri"))
}

func buildV1ClientURLBuildingParamsForSingleInstanceOfSomething(proj *models.Project, firstVar jen.Code, typVar jen.Code, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{
		firstVar,
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

func buildV1ClientURLBuildingParamsForListOfSomething(proj *models.Project, firstVar jen.Code, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{
		firstVar,
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

func buildV1ClientURLBuildingParamsForMethodThatIncludesItsOwnType(proj *models.Project, firstVar jen.Code, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{firstVar}

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

/*
// BuildItemExistsRequest builds an HTTP request for checking the existence of an item
func (c *V1Client) BuildItemExistsRequest(ctx context.Context, itemID uint64) (*http.Request, error) {
	_, span := trace.StartSpan(ctx, "BuildItemExistsRequest")
	defer span.End()

	uri := c.BuildURL(nil, itemsBasePath, strconv.FormatUint(itemID, 10))
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequest(http.MethodHead, uri, nil)
}
*/

func buildBuildItemExistsRequest(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()
	funcName := fmt.Sprintf("Build%sExistsRequest", ts)

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Callln(
			buildV1ClientURLBuildingParamsForSingleInstanceOfSomething(
				proj,
				jen.Nil(),
				jen.IDf("%sID", typ.Name.UnexportedVarName()),
				typ,
			)...,
		),
		attachURIToSpanCall(proj),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodHead"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for checking the existence of %s", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(buildParamsForMethodThatHandlesAnInstanceWithIDs(proj, typ, false)...).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildItemExists(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	funcName := fmt.Sprintf("%sExists", ts)
	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(jen.ID("req"), jen.Err()).Assign().ID("c").Dotf("Build%sExistsRequest", ts).Call(buildParamsForMethodThatHandlesAnInstanceWithIDs(proj, typ, true)...),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(
				jen.False(),
				jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Err()),
			),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("checkExistence").Call(utils.CtxVar(), jen.ID("req")),
	}

	lines := []jen.Code{
		jen.Commentf("%s retrieves whether or not %s exists", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(buildParamsForMethodThatHandlesAnInstanceWithIDs(proj, typ, false)...).Params(
			jen.ID("exists").Bool(),
			jen.Err().Error(),
		).Block(block...,
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
			buildV1ClientURLBuildingParamsForSingleInstanceOfSomething(
				proj,
				jen.Nil(),
				jen.IDf("%sID", typ.Name.UnexportedVarName()),
				typ,
			)...,
		),
		attachURIToSpanCall(proj),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for fetching %s", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(buildParamsForMethodThatHandlesAnInstanceWithIDs(proj, typ, false)...).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildParamsForMethodThatHandlesAnInstanceWithIDs(proj *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return utils.CtxVar()
			} else {
				return utils.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		if !call {
			params = append(params, jen.List(listParams...).Uint64())
		} else {
			params = append(params, listParams...)
		}
	} else {
		if !call {
			params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()).Uint64())
		} else {
			params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		}
	}

	return params
}

func buildParamsForMethodThatHandlesAnInstanceWithStructs(proj *models.Project, typ models.DataType) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{utils.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.ID(utils.BuildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
		listParams = append(listParams, jen.ID(utils.BuildFakeVarName(typ.Name.Singular())).Dot("ID"))

		params = append(params, listParams...)

	} else {
		params = append(params, jen.ID(utils.BuildFakeVarName(typ.Name.Singular())).Dot("ID"))
	}

	return params
}

func buildParamsForMethodThatRetrievesAListOfADataType(proj *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return utils.CtxVar()
			} else {
				return utils.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		if !call {
			params = append(params, jen.List(listParams...).Uint64())
		} else {
			params = append(params, listParams...)
		}
	}

	if !call {
		params = append(params, jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"))
	} else {
		params = append(params, jen.ID(utils.FilterVarName))
	}

	return params
}

func buildParamsForMethodThatRetrievesAListOfADataTypeFromStructs(proj *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return utils.CtxVar()
			} else {
				return utils.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.ID(pt.Name.UnexportedVarName()).Dot("ID"))
		}
		if !call {
			params = append(params, jen.List(listParams...).Uint64())
		} else {
			params = append(params, listParams...)
		}
	}

	if !call {
		params = append(params, jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"))
	} else {
		params = append(params, jen.ID(utils.FilterVarName))
	}

	return params
}

func buildParamsForMethodThatCreatesADataType(proj *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return utils.CtxVar()
			} else {
				return utils.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		if !call {
			params = append(params, jen.List(listParams...).Uint64())
		} else {
			params = append(params, listParams...)
		}
	}

	if !call {
		params = append(params, jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", typ.Name.Singular())))
	} else {
		params = append(params, jen.ID("input"))
	}

	return params
}

func buildParamsForMethodThatCreatesADataTypeFromStructs(proj *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return utils.CtxVar()
			} else {
				return utils.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.ID(pt.Name.UnexportedVarName()).Dot("ID"))
		}
		if !call {
			params = append(params, jen.List(listParams...).Uint64())
		} else {
			params = append(params, listParams...)
		}
	}

	if !call {
		params = append(params, jen.ID("exampleInput").PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", typ.Name.Singular())))
	} else {
		params = append(params, jen.ID("exampleInput"))
	}

	return params
}

func buildParamsForMethodThatFetchesAListOfDataTypesFromStructs(proj *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return utils.CtxVar()
			} else {
				return utils.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.ID(pt.Name.UnexportedVarName()).Dot("ID"))
		}
		if !call {
			params = append(params, jen.List(listParams...).Uint64())
		} else {
			params = append(params, listParams...)
		}
	}

	if !call {
		params = append(params, jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"))
	} else {
		params = append(params, jen.ID(utils.FilterVarName))
	}

	return params
}

func buildParamsForMethodThatIncludesItsOwnTypeInItsParams(proj *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return utils.CtxVar()
			} else {
				return utils.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = listParams[:len(listParams)-1]

		if len(listParams) > 0 {
			if !call {
				params = append(params, jen.List(listParams...).Uint64())
			} else {
				params = append(params, listParams...)
			}
		}
	}

	if !call {
		params = append(params, jen.ID(typ.Name.UnexportedVarName()).PointerTo().Qual(proj.ModelsV1Package(), typ.Name.Singular()))
	} else {
		params = append(params, jen.ID(typ.Name.UnexportedVarName()))
	}

	return params
}

func buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(proj *models.Project, typ models.DataType) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{utils.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.ID(utils.BuildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
		listParams = listParams[:len(listParams)-1]

		if len(listParams) > 0 {
			params = append(params, listParams...)
		}
	}
	params = append(params, jen.ID(utils.BuildFakeVarName(typ.Name.Singular())))

	return params
}

func buildGetSomethingFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	funcName := fmt.Sprintf("Get%s", ts)
	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(jen.ID("req"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", ts)).Call(buildParamsForMethodThatHandlesAnInstanceWithIDs(proj, typ, true)...),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Err()),
			),
		),
		jen.Line(),
		jen.If(jen.ID("retrieveErr").Assign().ID("c").Dot("retrieve").Call(utils.CtxVar(), jen.ID("req"), jen.AddressOf().ID(uvn)), jen.ID("retrieveErr").DoesNotEqual().ID("nil")).Block(
			jen.Return().List(jen.Nil(), jen.ID("retrieveErr")),
		),
		jen.Line(),
		jen.Return().List(jen.ID(uvn), jen.Nil()),
	}

	lines := []jen.Code{
		jen.Commentf("%s retrieves %s", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(buildParamsForMethodThatHandlesAnInstanceWithIDs(proj, typ, false)...).Params(
			jen.ID(uvn).PointerTo().Qual(proj.ModelsV1Package(), ts),
			jen.Err().Error(),
		).Block(block...,
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetListOfSomethingRequestFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	tp := typ.Name.Plural() // title plural

	funcName := fmt.Sprintf("BuildGet%sRequest", tp)

	urlBuildingParams := buildV1ClientURLBuildingParamsForListOfSomething(
		proj,
		jen.ID(utils.FilterVarName).Dot("ToValues").Call(),
		typ,
	)

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("BuildURL").Callln(urlBuildingParams...),
		attachURIToSpanCall(proj),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for fetching %s", funcName, typ.Name.PluralCommonName()),
		jen.Line(),
		newClientMethod(funcName).Params(buildParamsForMethodThatRetrievesAListOfADataType(proj, typ, false)...).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Block(block...),
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
		jen.List(jen.ID("req"), jen.Err()).Assign().ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", tp)).Call(
			buildParamsForMethodThatFetchesAListOfDataTypesFromStructs(proj, typ, true)...,
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.If(
			jen.ID("retrieveErr").Assign().ID("c").Dot("retrieve").Call(utils.CtxVar(), jen.ID("req"), jen.AddressOf().ID(pvn)),
			jen.ID("retrieveErr").DoesNotEqual().ID("nil"),
		).Block(
			jen.Return().List(jen.Nil(), jen.ID("retrieveErr")),
		),
		jen.Line(),
		jen.Return().List(jen.ID(pvn), jen.Nil()),
	}

	lines := []jen.Code{
		jen.Commentf("%s retrieves a list of %s", funcName, typ.Name.PluralCommonName()),
		jen.Line(),
		newClientMethod(funcName).Params(buildParamsForMethodThatFetchesAListOfDataTypesFromStructs(proj, typ, false)...).Params(
			jen.ID(pvn).PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", ts)),
			jen.Err().Error(),
		).Block(block...,
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
			buildV1ClientURLBuildingParamsForListOfSomething(
				proj,
				jen.Nil(),
				typ,
			)...,
		),
		attachURIToSpanCall(proj),
		jen.Line(),
		jen.Return().ID("c").Dot("buildDataRequest").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.ID("input"),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for creating %s", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(
			buildParamsForMethodThatCreatesADataType(proj, typ, false)...,
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Block(
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
			jen.ID("req"),
			jen.Err(),
		).Assign().ID("c").Dot(fmt.Sprintf("BuildCreate%sRequest", ts)).Call(
			buildParamsForMethodThatCreatesADataType(proj, typ, true)...,
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Err().Equals().ID("c").Dot("executeRequest").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.AddressOf().ID(vn),
		),
		jen.Return().List(
			jen.ID(vn),
			jen.Err(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s creates %s", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(
			buildParamsForMethodThatCreatesADataType(proj, typ, false)...,
		).Params(
			jen.ID(vn).PointerTo().Qual(proj.ModelsV1Package(), ts),
			jen.Err().Error(),
		).Block(block...,
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
			buildV1ClientURLBuildingParamsForMethodThatIncludesItsOwnType(proj, jen.Nil(), typ)...,
		),
		attachURIToSpanCall(proj),
		jen.Line(),
		jen.Return().ID("c").Dot("buildDataRequest").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodPut"),
			jen.ID("uri"),
			jen.ID(typ.Name.UnexportedVarName()),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for updating %s", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(
			buildParamsForMethodThatIncludesItsOwnTypeInItsParams(proj, typ, false)...,
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Block(block...,
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
			jen.ID("req"),
			jen.Err(),
		).Assign().ID("c").Dot(fmt.Sprintf("BuildUpdate%sRequest", ts)).Call(
			buildParamsForMethodThatIncludesItsOwnTypeInItsParams(proj, typ, true)...,
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("building request: %w"),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("executeRequest").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.AddressOf().ID(typ.Name.UnexportedVarName()),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s updates %s", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(
			buildParamsForMethodThatIncludesItsOwnTypeInItsParams(proj, typ, false)...,
		).Params(jen.Error()).Block(block...),
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
			buildV1ClientURLBuildingParamsForSingleInstanceOfSomething(
				proj,
				jen.Nil(),
				jen.IDf("%sID", typ.Name.UnexportedVarName()),
				typ,
			)...,
		),
		attachURIToSpanCall(proj),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodDelete"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for updating %s", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(buildParamsForMethodThatHandlesAnInstanceWithIDs(proj, typ, false)...).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Block(block...),
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
			jen.ID("req"),
			jen.Err(),
		).Assign().ID("c").Dot(fmt.Sprintf("BuildArchive%sRequest", ts)).Call(
			buildParamsForMethodThatHandlesAnInstanceWithIDs(proj, typ, true)...,
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("building request: %w"),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("executeRequest").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s archives %s", funcName, commonNameWithPrefix),
		jen.Line(),
		newClientMethod(funcName).Params(buildParamsForMethodThatHandlesAnInstanceWithIDs(proj, typ, false)...).Params(jen.Error()).Block(block...),
	}

	return lines
}
