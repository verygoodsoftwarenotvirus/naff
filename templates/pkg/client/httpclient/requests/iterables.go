package requests

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	puvn := typ.Name.PluralUnexportedVarName()

	code.Add(
		jen.Const().Defs(
			jen.IDf("%sBasePath", puvn).Equals().Lit(typ.Name.PluralRouteName()),
		),
		jen.Newline(),
	)

	code.Add(buildBuildGetSomethingRequest(proj, typ)...)

	if typ.SearchEnabled {
		code.Add(buildBuildSearchSomethingRequest(proj, typ)...)
	}

	code.Add(buildBuildGetListOfSomethingsRequest(proj, typ)...)
	code.Add(buildBuildCreateSomethingRequest(proj, typ)...)
	code.Add(buildBuildUpdateSomethingRequest(proj, typ)...)
	code.Add(buildBuildArchiveSomethingRequest(proj, typ)...)

	return code
}

func buildIDBoilerplate(proj *models.Project, typ models.DataType, includeType bool) []jen.Code {
	lines := []jen.Code{}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines,
			jen.If(jen.IDf("%sID", dep.Name.UnexportedVarName()).IsEqualTo().EmptyString()).Body(
				jen.Return(jen.Nil(), jen.ID("ErrInvalidIDProvided")),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", dep.Name.Singular())), jen.IDf("%sID", dep.Name.UnexportedVarName())),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", dep.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", dep.Name.UnexportedVarName())),
			jen.Newline(),
		)
	}

	if includeType {
		lines = append(lines,
			jen.If(jen.IDf("%sID", typ.Name.UnexportedVarName()).IsEqualTo().EmptyString()).Body(
				jen.Return().List(jen.Nil(), jen.ID("ErrInvalidIDProvided")),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", typ.Name.Singular())), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", typ.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Newline(),
		)
	}

	return lines
}

func buildGetSomethingURLParts(proj *models.Project, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{
		jen.ID("ctx"),
		jen.Nil(),
	}

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		urlBuildingParams = append(urlBuildingParams,
			jen.IDf("%sBasePath", pt.Name.PluralUnexportedVarName()),
			jen.IDf("%sID", pt.Name.UnexportedVarName()),
		)
	}
	urlBuildingParams = append(urlBuildingParams,
		jen.ID(basePath),
		jen.IDf("%sID", typ.Name.UnexportedVarName()),
	)

	return urlBuildingParams
}

func buildParamsForHTTPClientExistenceRequestBuildingMethod(p *models.Project, typ models.DataType) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		params = append(params, jen.List(listParams...).String())

	} else {
		params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()).String())
	}

	return params
}

func buildBuildSomethingExistsRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("b").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDBoilerplate(proj, typ, true)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(
			buildGetSomethingURLParts(proj, typ)...,
		),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestURIToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("uri")),
		jen.Newline(),
		jen.List(jen.ID("req"), jen.ID("err")).Assign().Qual("net/http", "NewRequestWithContext").Call(
			jen.ID("ctx"),
			jen.Qual("net/http", "MethodHead"),
			jen.ID("uri"),
			jen.Nil(),
		),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit("building user status request"),
			),
			),
		),
		jen.Newline(),
		jen.Return().List(jen.ID("req"), jen.Nil()),
	)

	lines := []jen.Code{
		jen.Commentf("Build%sExistsRequest builds an HTTP request for checking the existence of %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("Build%sExistsRequest", sn).Params(
			buildParamsForHTTPClientExistenceRequestBuildingMethod(proj, typ)...,
		).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	bodyLines := []jen.Code{jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Assign().ID("b").Dot("logger"),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDBoilerplate(proj, typ, true)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(
			buildGetSomethingURLParts(proj, typ)...,
		),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestURIToSpan").Call(
			jen.ID("span"),
			jen.ID("uri"),
		),
		jen.Newline(),
		jen.List(jen.ID("req"), jen.ID("err")).Assign().Qual("net/http", "NewRequestWithContext").Call(
			jen.ID("ctx"),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit("building user status request"),
			)),
		),
		jen.Newline(),
		jen.Return().List(jen.ID("req"), jen.Nil()),
	)

	lines := []jen.Code{
		jen.Commentf("BuildGet%sRequest builds an HTTP request for fetching %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("BuildGet%sRequest", sn).Params(
			buildParamsForHTTPClientExistenceRequestBuildingMethod(proj, typ)...,
		).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildSearchForSomethingURLParts(proj *models.Project, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{
		jen.ID("ctx"),
		jen.ID("params"),
	}

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		urlBuildingParams = append(urlBuildingParams,
			jen.IDf("%sBasePath", pt.Name.PluralUnexportedVarName()),
			jen.IDf("%sID", pt.Name.UnexportedVarName()),
		)
	}
	urlBuildingParams = append(urlBuildingParams,
		jen.ID(basePath),
		jen.Lit("search"),
	)

	return urlBuildingParams
}

func buildBuildSearchSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	params := []jen.Code{
		jen.ID("ctx").Qual("context", "Context"),
	}

	idParams := []jen.Code{}
	for _, owner := range proj.FindOwnerTypeChain(typ) {
		idParams = append(idParams, jen.IDf("%sID", owner.Name.UnexportedVarName()))
	}

	if len(idParams) > 0 {
		params = append(params, jen.List(idParams...).String())
	}

	params = append(params,
		jen.ID("query").String(),
		jen.ID("limit").ID("uint8"),
	)

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDBoilerplate(proj, typ, false)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("b").Dot(constants.LoggerVarName).Dot("WithValue").Call(
			jen.ID("types").Dot("SearchQueryKey"),
			jen.ID("query"),
		).Dot("WithValue").Call(
			jen.ID("types").Dot("LimitQueryKey"),
			jen.ID("limit"),
		),
		jen.Newline(),
		jen.ID("params").Assign().Qual("net/url", "Values").Values(),
		jen.ID("params").Dot("Set").Call(
			jen.ID("types").Dot("SearchQueryKey"),
			jen.ID("query"),
		),
		jen.ID("params").Dot("Set").Call(
			jen.ID("types").Dot("LimitQueryKey"),
			jen.Qual("strconv", "FormatUint").Call(
				jen.Uint64().Call(jen.ID("limit")),
				jen.Lit(10),
			),
		),
		jen.Newline(),
		jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(
			buildSearchForSomethingURLParts(proj, typ)...,
		),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestURIToSpan").Call(
			jen.ID("span"),
			jen.ID("uri"),
		),
		jen.Newline(),
		jen.List(jen.ID("req"), jen.ID("err")).Assign().Qual("net/http", "NewRequestWithContext").Call(
			jen.ID("ctx"),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit("building user status request"),
			)),
		),
		jen.Newline(),
		jen.Return().List(jen.ID("req"), jen.Nil()),
	)

	lines := []jen.Code{
		jen.Commentf("BuildSearch%sRequest builds an HTTP request for querying %s.", pn, pcn),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("BuildSearch%sRequest", pn).Params(params...).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildParamsForHTTPClientMethodThatFetchesAList(p *models.Project, typ models.DataType) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		params = append(params, jen.List(listParams...).String())
	}

	params = append(params, jen.ID(constants.FilterVarName).PointerTo().Qual(p.TypesPackage(), "QueryFilter"))

	return params
}

func buildV1ClientURLBuildingParamsForListOfSomething(proj *models.Project, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{
		constants.CtxVar(),
		jen.ID(constants.FilterVarName).Dot("ToValues").Call(),
	}

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		urlBuildingParams = append(urlBuildingParams,
			jen.IDf("%sBasePath", pt.Name.PluralUnexportedVarName()),
			jen.IDf("%sID", pt.Name.UnexportedVarName()),
		)
	}
	urlBuildingParams = append(urlBuildingParams,
		jen.ID(basePath),
	)

	return urlBuildingParams
}

func buildBuildGetListOfSomethingsRequest(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("filter").Dot("AttachToLogger").Call(jen.ID("b").Dot("logger")),
		jen.Newline(),
	}

	urlBuildingParams := buildV1ClientURLBuildingParamsForListOfSomething(proj, typ)

	bodyLines = append(bodyLines, buildIDBoilerplate(proj, typ, false)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(urlBuildingParams...),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestURIToSpan").Call(
			jen.ID("span"),
			jen.ID("uri"),
		),
		jen.Qual(proj.InternalTracingPackage(), "AttachQueryFilterToSpan").Call(
			jen.ID("span"),
			jen.ID("filter"),
		),
		jen.Newline(),
		jen.List(jen.ID("req"), jen.ID("err")).Assign().Qual("net/http", "NewRequestWithContext").Call(
			jen.ID("ctx"),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit("building user status request"),
			)),
		),
		jen.Newline(),
		jen.Return().List(jen.ID("req"), jen.Nil()),
	)

	lines := []jen.Code{
		jen.Commentf("BuildGet%sRequest builds an HTTP request for fetching a list of %s.", pn, pcn),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("BuildGet%sRequest", pn).Params(
			buildParamsForHTTPClientMethodThatFetchesAList(proj, typ)...,
		).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildIDBoilerplateForCreate(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	ownerTypes := proj.FindOwnerTypeChain(typ)

	for i, dep := range ownerTypes {
		if i != len(ownerTypes)-1 {
			lines = append(lines,
				jen.If(jen.IDf("%sID", dep.Name.UnexportedVarName()).IsEqualTo().EmptyString()).Body(
					jen.Return(jen.Nil(), jen.ID("ErrInvalidIDProvided")),
				),
				jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", dep.Name.Singular())), jen.IDf("%sID", dep.Name.UnexportedVarName())),
				jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", dep.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", dep.Name.UnexportedVarName())),
				jen.Newline(),
			)
		}
	}

	return lines
}

func buildV1ClientURLBuildingParamsForCreatingSomething(proj *models.Project, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{
		constants.CtxVar(),
		jen.Nil(),
	}

	parents := proj.FindOwnerTypeChain(typ)
	for i, pt := range parents {
		urlBuildingParams = append(urlBuildingParams,
			jen.IDf("%sBasePath", pt.Name.PluralUnexportedVarName()),
			func() jen.Code {
				if i == len(parents)-1 && typ.BelongsToStruct != nil {
					return jen.ID("input").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular())
				}
				return jen.IDf("%sID", pt.Name.UnexportedVarName())
			}(),
		)
	}
	urlBuildingParams = append(urlBuildingParams,
		jen.ID(basePath),
	)

	return urlBuildingParams
}

func buildParamsForHTTPClientCreateRequestBuildingMethod(p *models.Project, typ models.DataType) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for i, pt := range parents {
			if i == len(parents)-1 && typ.BelongsToStruct != nil {
				continue
			} else {
				listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
			}
		}
		if len(listParams) > 0 {
			params = append(params, jen.List(listParams...).String())
		}
	}

	params = append(params, jen.ID("input").PointerTo().Qual(p.TypesPackage(), fmt.Sprintf("%sCreationRequestInput", typ.Name.Singular())))

	return params
}

func buildBuildCreateSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("b").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDBoilerplateForCreate(proj, typ)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.If(jen.ID("input").IsEqualTo().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.ID("ErrNilInputProvided")),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit("validating input"),
			),
			),
		),
		jen.Newline(),
		jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(
			buildV1ClientURLBuildingParamsForCreatingSomething(proj, typ)...,
		),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestURIToSpan").Call(
			jen.ID("span"),
			jen.ID("uri"),
		),
		jen.Newline(),
		jen.List(jen.ID("req"), jen.Err()).Assign().ID("b").Dot("buildDataRequest").Call(
			jen.ID("ctx"),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.ID("input"),
		),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit("building request"),
			)),
		),
		jen.Newline(),
		jen.Return(jen.ID("req"), jen.Nil()),
	)

	lines := []jen.Code{
		jen.Commentf("BuildCreate%sRequest builds an HTTP request for creating %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("BuildCreate%sRequest", sn).Params(
			buildParamsForHTTPClientCreateRequestBuildingMethod(proj, typ)...,
		).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildIDBoilerplateForUpdate(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	ownerTypes := proj.FindOwnerTypeChain(typ)

	for i, dep := range ownerTypes {
		if i != len(ownerTypes)-1 {
			lines = append(lines,
				jen.If(jen.IDf("%sID", dep.Name.UnexportedVarName()).IsEqualTo().EmptyString()).Body(
					jen.Return(jen.Nil(), jen.ID("ErrInvalidIDProvided")),
				),
				jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", dep.Name.Singular())), jen.IDf("%sID", dep.Name.UnexportedVarName())),
				jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", dep.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", dep.Name.UnexportedVarName())),
				jen.Newline(),
			)
		}
	}

	return lines
}

func buildParamsForHTTPClientUpdateRequestBuildingMethod(p *models.Project, typ models.DataType) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = listParams[:len(listParams)-1]

		if len(listParams) > 0 {
			params = append(params, jen.List(listParams...).String())
		}
	}

	params = append(params, jen.ID(typ.Name.UnexportedVarName()).PointerTo().Qual(p.TypesPackage(), typ.Name.Singular()))

	return params
}

func buildV1ClientURLBuildingParamsForMethodThatIncludesItsOwnType(proj *models.Project, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{
		constants.CtxVar(),
		jen.Nil(),
	}

	ownerChain := proj.FindOwnerTypeChain(typ)
	for i, pt := range ownerChain {
		if i == len(ownerChain)-1 {
			urlBuildingParams = append(urlBuildingParams,
				jen.IDf("%sBasePath", pt.Name.PluralUnexportedVarName()),
				jen.ID(typ.Name.UnexportedVarName()).Dotf("BelongsTo%s", pt.Name.Singular()),
			)
		} else {
			urlBuildingParams = append(urlBuildingParams,
				jen.IDf("%sBasePath", pt.Name.PluralUnexportedVarName()),
				jen.IDf("%sID", pt.Name.UnexportedVarName()),
			)
		}
	}

	urlBuildingParams = append(urlBuildingParams,
		jen.ID(basePath),
		jen.ID(typ.Name.UnexportedVarName()).Dot("ID"),
	)

	return urlBuildingParams
}

func buildBuildUpdateSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	uvn := typ.Name.UnexportedVarName()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("b").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDBoilerplateForUpdate(proj, typ)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.If(jen.ID(uvn).IsEqualTo().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.ID("ErrNilInputProvided")),
		),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", sn)), jen.ID(uvn).Dot("ID")),
		jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
			jen.ID("span"),
			jen.ID(uvn).Dot("ID"),
		),
		jen.Newline(),
		jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(
			buildV1ClientURLBuildingParamsForMethodThatIncludesItsOwnType(proj, typ)...,
		),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestURIToSpan").Call(
			jen.ID("span"),
			jen.ID("uri"),
		),
		jen.Newline(),
		jen.List(jen.ID("req"), jen.Err()).Assign().ID("b").Dot("buildDataRequest").Call(
			jen.ID("ctx"),
			jen.Qual("net/http", "MethodPut"),
			jen.ID("uri"),
			jen.ID(uvn),
		),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit("building request"),
			)),
		),
		jen.Newline(),
		jen.Return(jen.ID("req"), jen.Nil()),
	)

	lines := []jen.Code{
		jen.Commentf("BuildUpdate%sRequest builds an HTTP request for updating %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("BuildUpdate%sRequest", sn).Params(
			buildParamsForHTTPClientUpdateRequestBuildingMethod(proj, typ)...,
		).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildArchiveSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("b").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDBoilerplate(proj, typ, true)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(
			buildGetSomethingURLParts(proj, typ)...,
		),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestURIToSpan").Call(
			jen.ID("span"),
			jen.ID("uri"),
		),
		jen.Newline(),
		jen.List(jen.ID("req"), jen.ID("err")).Assign().Qual("net/http", "NewRequestWithContext").Call(
			jen.ID("ctx"),
			jen.Qual("net/http", "MethodDelete"),
			jen.ID("uri"),
			jen.Nil(),
		),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit("building user status request"),
			),
			),
		),
		jen.Newline(),
		jen.Return().List(jen.ID("req"), jen.Nil()),
	)

	lines := []jen.Code{
		jen.Commentf("BuildArchive%sRequest builds an HTTP request for archiving %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("BuildArchive%sRequest", sn).Params(
			buildParamsForHTTPClientExistenceRequestBuildingMethod(proj, typ)...,
		).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildGetAuditLogsForSomethingURLParts(proj *models.Project, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{
		jen.ID("ctx"),
		jen.Nil(),
	}

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		urlBuildingParams = append(urlBuildingParams,
			jen.IDf("%sBasePath", pt.Name.PluralUnexportedVarName()),
			jen.IDf("%sID", pt.Name.UnexportedVarName()),
		)
	}
	urlBuildingParams = append(urlBuildingParams,
		jen.ID(basePath),
		jen.IDf("%sID", typ.Name.UnexportedVarName()),
		jen.Lit("audit"),
	)

	return urlBuildingParams
}

func buildBuildGetAuditLogForSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	uvn := typ.Name.UnexportedVarName()

	params := []jen.Code{
		jen.ID("ctx").Qual("context", "Context"),
	}

	idParams := []jen.Code{}
	for _, owner := range proj.FindOwnerTypeChain(typ) {
		idParams = append(idParams, jen.IDf("%sID", owner.Name.UnexportedVarName()))
	}

	idParams = append(idParams, jen.IDf("%sID", uvn))
	params = append(params, jen.List(idParams...).String())

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("b").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDBoilerplate(proj, typ, true)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("uri").Assign().ID("b").Dot("BuildURL").Callln(
			buildGetAuditLogsForSomethingURLParts(proj, typ)...,
		),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestURIToSpan").Call(
			jen.ID("span"),
			jen.ID("uri"),
		),
		jen.Newline(),
		jen.List(jen.ID("req"), jen.ID("err")).Assign().Qual("net/http", "NewRequestWithContext").Call(
			jen.ID("ctx"),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit("building user status request"),
			),
			),
		),
		jen.Newline(),
		jen.Return().List(jen.ID("req"), jen.Nil()),
	)

	lines := []jen.Code{
		jen.Commentf("BuildGetAuditLogFor%sRequest builds an HTTP request for fetching a list of audit log entries pertaining to %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("Builder")).IDf("BuildGetAuditLogFor%sRequest", sn).Params(params...).Params(jen.PointerTo().Qual("net/http", "Request"), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}
