package httpclient

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildIDBoilerplate(proj *models.Project, typ models.DataType, includeType bool, returnVal jen.Code) []jen.Code {
	lines := []jen.Code{}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines,
			jen.If(jen.IDf("%sID", dep.Name.UnexportedVarName()).IsEqualTo().EmptyString()).Body(
				jen.Return(returnVal, jen.ID("ErrInvalidIDProvided")),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", dep.Name.Singular())), jen.IDf("%sID", dep.Name.UnexportedVarName())),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", dep.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", dep.Name.UnexportedVarName())),
			jen.Newline(),
		)
	}

	if includeType {
		lines = append(lines,
			jen.If(jen.IDf("%sID", typ.Name.UnexportedVarName()).IsEqualTo().EmptyString()).Body(
				jen.Return().List(returnVal, jen.ID("ErrInvalidIDProvided")),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", typ.Name.Singular())), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", typ.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Newline(),
		)
	}

	return lines
}

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildGetSomething(proj, typ)...)

	if typ.SearchEnabled {
		code.Add(buildSearchSomething(proj, typ)...)
	}

	code.Add(buildGetListOfSomething(proj, typ)...)
	code.Add(buildCreateSomething(proj, typ)...)
	code.Add(buildUpdateSomething(proj, typ)...)
	code.Add(buildArchiveSomething(proj, typ)...)

	return code
}

func buildParamsForHTTPClientExistenceMethod(p *models.Project, typ models.DataType) []jen.Code {
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

func buildArgsForHTTPClientExistenceRequestBuildingMethod(p *models.Project, typ models.DataType) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))

		params = append(params, listParams...)
	} else {
		params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()))
	}

	return params
}
func buildSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	uvn := typ.Name.UnexportedVarName()

	lines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("c").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	lines = append(lines, buildIDBoilerplate(proj, typ, true, jen.False())...)

	lines = append(lines,
		jen.Newline(),
		jen.List(jen.ID("req"), jen.ID("err")).Assign().ID("c").Dot("requestBuilder").Dotf("Build%sExistsRequest", sn).Call(
			buildArgsForHTTPClientExistenceRequestBuildingMethod(proj, typ)...,
		),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().List(
				jen.ID("false"),
				jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("building %s existence request", scn),
				),
			),
		),
		jen.Newline(),
		jen.List(jen.ID("exists"),
			jen.ID("err")).Assign().ID("c").Dot("responseIsOK").Call(
			jen.ID("ctx"),
			jen.ID("req"),
		),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().List(
				jen.ID("false"),
				jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit(fmt.Sprintf("checking existence for %s", scn)+" %s"),
					jen.IDf("%sID", uvn),
				),
			),
		),
		jen.Newline(),
		jen.Return().List(jen.ID("exists"), jen.ID("nil")),
	)

	return []jen.Code{
		jen.Commentf("%sExists retrieves whether %s exists.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("%sExists", sn).Params(
			buildParamsForHTTPClientExistenceMethod(proj, typ)...,
		).Params(jen.ID("bool"), jen.ID("error")).Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildArgsForHTTPClientRetrievalRequestBuildingMethod(p *models.Project, typ models.DataType) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}

		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		params = append(params, listParams...)
	} else {
		params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()))
	}

	return params
}

func buildParamsForHTTPClientRetrievalMethod(p *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return constants.CtxVar()
			} else {
				return constants.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		if !call {
			params = append(params, jen.List(listParams...).String())
		} else {
			params = append(params, listParams...)
		}
	} else {
		if !call {
			params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()).String())
		} else {
			params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		}
	}

	return params
}

func buildGetSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	uvn := typ.Name.UnexportedVarName()

	lines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("c").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	lines = append(lines, buildIDBoilerplate(proj, typ, true, jen.Nil())...)

	lines = append(lines,
		jen.Newline(),
		jen.List(jen.ID("req"), jen.ID("err")).Assign().ID("c").Dot("requestBuilder").Dotf("BuildGet%sRequest", sn).Call(
			buildArgsForHTTPClientRetrievalRequestBuildingMethod(proj, typ)...,
		),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().List(
				jen.ID("nil"),
				jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("building get %s request", scn),
				),
			),
		),
		jen.Newline(),
		jen.Var().ID(uvn).PointerTo().Qual(proj.TypesPackage(), sn),
		jen.If(jen.ID("err").Equals().ID("c").Dot("fetchAndUnmarshal").Call(
			jen.ID("ctx"),
			jen.ID("req"),
			jen.AddressOf().ID(uvn),
		),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().List(
				jen.ID("nil"),
				jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("retrieving %s", scn),
				),
			),
		),
		jen.Newline(),
		jen.Return().List(jen.ID(uvn), jen.ID("nil")),
	)

	return []jen.Code{
		jen.Commentf("Get%s gets %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("Get%s", sn).Params(
			buildParamsForHTTPClientRetrievalMethod(proj, typ, false)...,
		).Params(jen.PointerTo().Qual(proj.TypesPackage(), sn), jen.ID("error")).Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildSearchSomethingParams(proj *models.Project, typ models.DataType) []jen.Code {
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

	return params
}

func buildSearchSomethingRequestBuildingArgs(proj *models.Project, typ models.DataType) []jen.Code {
	args := []jen.Code{
		jen.ID("ctx"),
	}

	for _, owner := range proj.FindOwnerTypeChain(typ) {
		args = append(args, jen.IDf("%sID", owner.Name.UnexportedVarName()))
	}

	args = append(args,
		jen.ID("query"),
		jen.ID("limit"),
	)

	return args
}

func buildSearchSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	puvn := typ.Name.PluralUnexportedVarName()

	lines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("c").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	lines = append(lines, buildIDBoilerplate(proj, typ, false, jen.Nil())...)

	lines = append(lines,
		jen.If(jen.ID("query").IsEqualTo().Lit("")).Body(
			jen.Return().List(
				jen.ID("nil"),
				jen.ID("ErrEmptyQueryProvided")),
		),
		jen.Newline(),
		jen.If(jen.ID("limit").IsEqualTo().Zero()).Body(
			jen.ID("limit").Equals().Lit(20)),
		jen.Newline(),
		jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(
			jen.Qual(proj.ObservabilityPackage("keys"), "SearchQueryKey"),
			jen.ID("query"),
		).Dot("WithValue").Call(
			jen.Qual(proj.ObservabilityPackage("keys"), "FilterLimitKey"),
			jen.ID("limit"),
		),
		jen.Newline(),
		jen.List(jen.ID("req"), jen.ID("err")).Assign().ID("c").Dot("requestBuilder").Dotf("BuildSearch%sRequest", pn).Call(
			buildSearchSomethingRequestBuildingArgs(proj, typ)...,
		),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().List(
				jen.ID("nil"),
				jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("building search for %s request", pcn),
				),
			),
		),
		jen.Newline(),
		jen.Var().ID(puvn).Index().PointerTo().Qual(proj.TypesPackage(), sn),
		jen.If(jen.ID("err").Equals().ID("c").Dot("fetchAndUnmarshal").Call(
			jen.ID("ctx"),
			jen.ID("req"),
			jen.AddressOf().ID(puvn),
		),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().List(
				jen.ID("nil"),
				jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("retrieving %s", pcn),
				),
			),
		),
		jen.Newline(),
		jen.Return().List(jen.ID(puvn), jen.ID("nil")),
	)

	return []jen.Code{
		jen.Commentf("Search%s searches through a list of %s.", pn, pcn),
		jen.Newline(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("Search%s", pn).Params(
			buildSearchSomethingParams(proj, typ)...,
		).Params(jen.Index().PointerTo().Qual(proj.TypesPackage(), sn), jen.ID("error")).Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildArgsForHTTPClientListRequestMethod(p *models.Project, typ models.DataType) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		params = append(params, listParams...)
	}

	params = append(params, jen.ID(constants.FilterVarName))

	return params
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

func buildGetListOfSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	puvn := typ.Name.PluralUnexportedVarName()

	lines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Assign().ID("c").Dot("loggerWithFilter").Call(jen.ID("filter")),
		jen.Qual(proj.InternalTracingPackage(), "AttachQueryFilterToSpan").Call(
			jen.ID("span"),
			jen.ID("filter"),
		),
		jen.Newline(),
	}

	lines = append(lines, buildIDBoilerplate(proj, typ, false, jen.Nil())...)

	lines = append(lines,
		jen.Newline(),
		jen.List(jen.ID("req"), jen.ID("err")).Assign().ID("c").Dot("requestBuilder").Dotf("BuildGet%sRequest", pn).Call(
			buildArgsForHTTPClientListRequestMethod(proj, typ)...,
		),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().List(jen.ID("nil"),
				jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("building %s list request", pcn),
				),
			),
		),
		jen.Newline(),
		jen.Var().ID(puvn).PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn)),
		jen.If(jen.ID("err").Equals().ID("c").Dot("fetchAndUnmarshal").Call(
			jen.ID("ctx"),
			jen.ID("req"),
			jen.AddressOf().ID(puvn),
		),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().List(jen.ID("nil"),
				jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("retrieving %s", pcn),
				),
			),
		),
		jen.Newline(),
		jen.Return().List(jen.ID(puvn), jen.ID("nil")),
	)

	return []jen.Code{
		jen.Commentf("Get%s retrieves a list of %s.", pn, pcn),
		jen.Newline(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("Get%s", pn).Params(
			buildParamsForHTTPClientMethodThatFetchesAList(proj, typ)...,
		).Params(jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn)), jen.ID("error")).Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildIDBoilerplateForCreation(proj *models.Project, typ models.DataType, includeType bool, returnVal jen.Code) []jen.Code {
	lines := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for i, dep := range owners {
		if i != len(owners)-1 {
			lines = append(lines,
				jen.If(jen.IDf("%sID", dep.Name.UnexportedVarName()).IsEqualTo().EmptyString()).Body(
					func() jen.Code {
						if returnVal != nil {
							return jen.Return(returnVal, jen.ID("ErrInvalidIDProvided"))
						} else {
							return jen.Return(jen.ID("ErrInvalidIDProvided"))
						}
					}(),
				),
				jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", dep.Name.Singular())), jen.IDf("%sID", dep.Name.UnexportedVarName())),
				jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", dep.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", dep.Name.UnexportedVarName())),
				jen.Newline(),
			)
		}
	}

	if includeType {
		lines = append(lines,
			jen.If(jen.IDf("%sID", typ.Name.UnexportedVarName()).IsEqualTo().EmptyString()).Body(
				func() jen.Code {
					if returnVal != nil {
						return jen.Return(returnVal, jen.ID("ErrInvalidIDProvided"))
					} else {
						return jen.Return(jen.ID("ErrInvalidIDProvided"))
					}
				}(),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", typ.Name.Singular())), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", typ.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Newline(),
		)
	}

	return lines
}

func buildParamsForHTTPClientCreateMethod(p *models.Project, typ models.DataType) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for i, pt := range parents {
			if i > len(parents)-2 {
				continue
			}
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		if len(listParams) > 0 {
			params = append(params, jen.List(listParams...).String())
		}
	}

	params = append(params, jen.ID("input").PointerTo().Qual(p.TypesPackage(), fmt.Sprintf("%sCreationRequestInput", typ.Name.Singular())))

	return params
}

func buildArgsForHTTPClientCreateRequestBuildingMethod(p *models.Project, typ models.DataType) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for i, pt := range parents {
			if i == len(parents)-1 {
				continue
			} else {
				listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
			}
		}
		params = append(params, listParams...)
	}

	params = append(params, jen.ID("input"))

	return params
}

func buildCreateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	lines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("c").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	lines = append(lines, buildIDBoilerplateForCreation(proj, typ, false, jen.Nil())...)

	lines = append(lines,
		jen.If(jen.ID("input").IsEqualTo().ID("nil")).Body(
			jen.Return().List(
				jen.EmptyString(),
				jen.ID("ErrNilInputProvided"),
			),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().List(
				jen.EmptyString(),
				jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating input"),
				),
			),
		),
		jen.Newline(),
		jen.List(jen.ID("req"), jen.ID("err")).Assign().ID("c").Dot("requestBuilder").Dotf("BuildCreate%sRequest", sn).Call(
			buildArgsForHTTPClientCreateRequestBuildingMethod(proj, typ)...,
		),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().List(
				jen.EmptyString(),
				jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("building create %s request", scn),
				),
			),
		),
		jen.Newline(),
		jen.Var().ID("pwr").PointerTo().Qual(proj.TypesPackage(), "PreWriteResponse"),
		jen.If(jen.ID("err").Equals().ID("c").Dot("fetchAndUnmarshal").Call(
			jen.ID("ctx"),
			jen.ID("req"),
			jen.AddressOf().ID("pwr"),
		),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().List(
				jen.EmptyString(),
				jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("creating %s", scn),
				),
			),
		),
		jen.Newline(),
		jen.Return().List(jen.ID("pwr").Dot("ID"), jen.ID("nil")),
	)

	return []jen.Code{
		jen.Commentf("Create%s creates %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("Create%s", sn).Params(
			buildParamsForHTTPClientCreateMethod(proj, typ)...,
		).Params(jen.String(), jen.ID("error")).Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildParamsForHTTPClientUpdateMethod(p *models.Project, typ models.DataType) []jen.Code {
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

func buildArgsForHTTPClientUpdateRequestBuildingMethod(p *models.Project, typ models.DataType) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = listParams[:len(listParams)-1]

		if len(listParams) > 0 {
			params = append(params, listParams...)
		}
	}

	params = append(params, jen.ID(typ.Name.UnexportedVarName()))

	return params
}

func buildUpdateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	uvn := typ.Name.UnexportedVarName()

	lines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("c").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	lines = append(lines, buildIDBoilerplateForCreation(proj, typ, false, jen.Null())...)

	lines = append(lines,
		jen.Newline(),
		jen.If(jen.ID(uvn).IsEqualTo().Nil()).Body(
			jen.Return(jen.ID("ErrNilInputProvided")),
		),
		jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(
			jen.Qual(proj.ObservabilityPackage("keys"), fmt.Sprintf("%sIDKey", sn)),
			jen.ID(uvn).Dot("ID"),
		),
		jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", typ.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%s", typ.Name.UnexportedVarName()).Dot("ID")),
		jen.Newline(),
		jen.List(jen.ID("req"), jen.ID("err")).Assign().ID("c").Dot("requestBuilder").Dotf("BuildUpdate%sRequest", sn).Call(
			buildArgsForHTTPClientUpdateRequestBuildingMethod(proj, typ)...,
		),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("building update %s request", scn),
			)),
		jen.Newline(),
		jen.If(jen.ID("err").Equals().ID("c").Dot("fetchAndUnmarshal").Call(
			jen.ID("ctx"),
			jen.ID("req"),
			jen.AddressOf().ID(uvn),
		),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit(fmt.Sprintf("updating %s", scn)+" %s"),
				jen.ID(uvn).Dot("ID"),
			)),
		jen.Newline(),
		jen.Return().ID("nil"),
	)

	return []jen.Code{
		jen.Commentf("Update%s updates %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("Update%s", sn).Params(
			buildParamsForHTTPClientUpdateMethod(proj, typ)...,
		).Params(jen.ID("error")).Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildParamsForHTTPClientArchiveMethod(p *models.Project, typ models.DataType) []jen.Code {
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

func buildArgsForHTTPClientArchiveRequestBuildingMethod(p *models.Project, typ models.DataType) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		params = append(params, listParams...)
	} else {
		params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()))
	}

	return params
}

func buildArchiveSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	uvn := typ.Name.UnexportedVarName()

	lines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("c").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	lines = append(lines, buildIDBoilerplate(proj, typ, true, jen.Null())...)

	lines = append(lines,
		jen.Newline(),
		jen.List(jen.ID("req"), jen.ID("err")).Assign().ID("c").Dot("requestBuilder").Dotf("BuildArchive%sRequest", sn).Call(
			buildArgsForHTTPClientArchiveRequestBuildingMethod(proj, typ)...,
		),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("building archive %s request", scn),
			)),
		jen.Newline(),
		jen.If(jen.ID("err").Equals().ID("c").Dot("fetchAndUnmarshal").Call(
			jen.ID("ctx"),
			jen.ID("req"),
			jen.ID("nil"),
		),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit(fmt.Sprintf("archiving %s", scn)+" %s"),
				jen.IDf("%sID", uvn),
			)),
		jen.Newline(),
		jen.Return().ID("nil"),
	)

	return []jen.Code{
		jen.Commentf("Archive%s archives %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("Archive%s", sn).Params(
			buildParamsForHTTPClientArchiveMethod(proj, typ)...,
		).Params(jen.ID("error")).Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildAuditSomethingParams(proj *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{
		jen.ID("ctx").Qual("context", "Context"),
	}

	idParams := []jen.Code{}
	for _, owner := range proj.FindOwnerTypeChain(typ) {
		idParams = append(idParams, jen.IDf("%sID", owner.Name.UnexportedVarName()))
	}
	idParams = append(idParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))

	params = append(params, jen.List(idParams...).String())

	return params
}

func buildAuditSomethingRequestBuildingArgs(proj *models.Project, typ models.DataType) []jen.Code {
	args := []jen.Code{
		jen.ID("ctx"),
	}

	for _, owner := range proj.FindOwnerTypeChain(typ) {
		args = append(args, jen.IDf("%sID", owner.Name.UnexportedVarName()))
	}

	args = append(args,
		jen.IDf("%sID", typ.Name.UnexportedVarName()),
	)

	return args
}

func buildGetAuditLogForSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	lines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("c").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	lines = append(lines, buildIDBoilerplate(proj, typ, true, jen.Nil())...)

	lines = append(lines,
		jen.Newline(),
		jen.List(jen.ID("req"), jen.ID("err")).Assign().ID("c").Dot("requestBuilder").Dotf("BuildGetAuditLogFor%sRequest", sn).Call(
			buildAuditSomethingRequestBuildingArgs(proj, typ)...,
		),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().List(jen.ID("nil"),
				jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("building get audit log entries for %s request", scn),
				),
			),
		),
		jen.Newline(),
		jen.Var().ID("entries").Index().PointerTo().Qual(proj.TypesPackage(), "AuditLogEntry"),
		jen.If(jen.ID("err").Equals().ID("c").Dot("fetchAndUnmarshal").Call(
			jen.ID("ctx"),
			jen.ID("req"),
			jen.AddressOf().ID("entries"),
		),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().List(jen.ID("nil"),
				jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving plan"),
				),
			),
		),
		jen.Newline(),
		jen.Return().List(jen.ID("entries"), jen.ID("nil")),
	)

	return []jen.Code{
		jen.Commentf("GetAuditLogFor%s retrieves a list of audit log entries pertaining to %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("GetAuditLogFor%s", sn).Params(
			buildAuditSomethingParams(proj, typ)...,
		).Params(jen.Index().PointerTo().Qual(proj.TypesPackage(), "AuditLogEntry"), jen.ID("error")).Body(
			lines...,
		),
		jen.Newline(),
	}
}
