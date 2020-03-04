package client

import (
	"fmt"
	"path/filepath"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// func buildPathForType(pkg *models.Project, typ models.DataType) string {
// 	if typ.BelongsToStruct != nil {
// 		parentType := pkg.FindType(typ.BelongsToStruct.Singular())
// 		if parentType != nil {
// 			return filepath.Join(buildPathForType(pkg, *parentType), )
// 		}
// 	}
// }

func iterablesDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("client")

	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())

	utils.AddImports(pkg.OutputPath, []models.DataType{typ}, ret)
	ret.Add(jen.Const().Defs(
		jen.ID(basePath).Op("=").Lit(typ.Name.PluralRouteName())),
	)

	ret.Add(buildBuildGetSomethingRequestFuncDecl(pkg, typ)...)
	ret.Add(buildGetSomethingFuncDecl(pkg, typ)...)
	ret.Add(buildBuildGetSomethingsRequestFuncDecl(pkg, typ)...)
	ret.Add(buildGetListOfSomethingsFuncDecl(pkg, typ)...)
	ret.Add(buildBuildCreateSomethingRequestFuncDecl(pkg, typ)...)
	ret.Add(buildCreateSomethingFuncDecl(pkg, typ)...)
	ret.Add(buildBuildUpdateSomethingRequestFuncDecl(pkg, typ)...)
	ret.Add(buildUpdateSomethingFuncDecl(pkg, typ)...)
	ret.Add(buildBuildArchiveSomethingRequestFuncDecl(pkg, typ)...)
	ret.Add(buildArchiveSomethingFuncDecl(pkg, typ)...)

	return ret
}

func buildV1ClientURLBuildingParamsForSingleInstanceOfSomething(pkg *models.Project, firstVar jen.Code, typVar jen.Code, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{
		firstVar,
	}

	for _, pt := range pkg.FindOwnerTypeChain(typ) {
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

func buildV1ClientURLBuildingParamsForListOfSomething(pkg *models.Project, firstVar jen.Code, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{
		firstVar,
	}

	for _, pt := range pkg.FindOwnerTypeChain(typ) {
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

func buildV1ClientURLBuildingParamsForMethodThatIncludesItsOwnType(pkg *models.Project, firstVar jen.Code, typ models.DataType) []jen.Code {
	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())
	urlBuildingParams := []jen.Code{firstVar}

	ownerChain := pkg.FindOwnerTypeChain(typ)
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

func buildBuildGetSomethingRequestFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("BuildGet%sRequest builds an HTTP request for fetching %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildGet%sRequest", ts)).Params(buildParamsForMethodThatHandlesAnInstanceWithIDs(pkg, typ, false)...).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Callln(
				buildV1ClientURLBuildingParamsForSingleInstanceOfSomething(
					pkg,
					jen.Nil(),
					jen.IDf("%sID", typ.Name.UnexportedVarName()),
					typ,
				)...,
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildParamsForMethodThatHandlesAnInstanceWithIDs(pkg *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := pkg.FindOwnerTypeChain(typ)
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
			params = append(params, jen.List(listParams...).ID("uint64"))
		} else {
			params = append(params, listParams...)
		}
	} else {
		if !call {
			params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()).ID("uint64"))
		} else {
			params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		}
	}

	return params
}

func buildParamsForMethodThatHandlesAnInstanceWithStructs(pkg *models.Project, typ models.DataType) []jen.Code {
	parents := pkg.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{utils.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.ID(pt.Name.UnexportedVarName()).Dot("ID"))
		}
		listParams = append(listParams, jen.ID(typ.Name.UnexportedVarName()).Dot("ID"))

		params = append(params, listParams...)

	} else {
		params = append(params, jen.ID(typ.Name.UnexportedVarName()).Dot("ID"))

	}

	return params
}

func buildParamsForMethodThatRetrievesAListOfADataType(pkg *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := pkg.FindOwnerTypeChain(typ)
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
			params = append(params, jen.List(listParams...).ID("uint64"))
		} else {
			params = append(params, listParams...)
		}
	}

	if !call {
		params = append(params, jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"))
	} else {
		params = append(params, jen.ID("filter"))
	}

	return params
}

func buildParamsForMethodThatRetrievesAListOfADataTypeFromStructs(pkg *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := pkg.FindOwnerTypeChain(typ)
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
			params = append(params, jen.List(listParams...).ID("uint64"))
		} else {
			params = append(params, listParams...)
		}
	}

	if !call {
		params = append(params, jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"))
	} else {
		params = append(params, jen.ID("filter"))
	}

	return params
}

func buildParamsForMethodThatCreatesADataType(pkg *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := pkg.FindOwnerTypeChain(typ)
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
			params = append(params, jen.List(listParams...).ID("uint64"))
		} else {
			params = append(params, listParams...)
		}
	}

	if !call {
		params = append(params, jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", typ.Name.Singular())))
	} else {
		params = append(params, jen.ID("input"))
	}

	return params
}

func buildParamsForMethodThatCreatesADataTypeFromStructs(pkg *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := pkg.FindOwnerTypeChain(typ)
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
			params = append(params, jen.List(listParams...).ID("uint64"))
		} else {
			params = append(params, listParams...)
		}
	}

	if !call {
		params = append(params, jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", typ.Name.Singular())))
	} else {
		params = append(params, jen.ID("input"))
	}

	return params
}

func buildParamsForMethodThatIncludesItsOwnTypeInItsParams(pkg *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := pkg.FindOwnerTypeChain(typ)
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
				params = append(params, jen.List(listParams...).ID("uint64"))
			} else {
				params = append(params, listParams...)
				params = append(params, jen.ID(typ.Name.UnexportedVarName()))
			}
		}
	}

	if !call {
		params = append(params, jen.ID(typ.Name.UnexportedVarName()).Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), typ.Name.Singular()))
	}

	return params
}

func buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(pkg *models.Project, typ models.DataType) []jen.Code {
	parents := pkg.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{utils.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.ID(pt.Name.UnexportedVarName()).Dot("ID"))
		}
		listParams = listParams[:len(listParams)-1]

		if len(listParams) > 0 {
			params = append(params, listParams...)
			params = append(params, jen.ID(typ.Name.UnexportedVarName()))

		}
	}

	return params
}

func buildGetSomethingFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()

	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("Get%s retrieves %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Get%s", ts)).Params(buildParamsForMethodThatHandlesAnInstanceWithIDs(pkg, typ, false)...).Params(
			jen.ID(uvn).Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), ts),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", ts)).Call(buildParamsForMethodThatHandlesAnInstanceWithIDs(pkg, typ, true)...),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.ID("err")),
				),
			),
			jen.Line(),
			jen.If(jen.ID("retrieveErr").Op(":=").ID("c").Dot("retrieve").Call(jen.ID("ctx"), jen.ID("req"), jen.Op("&").ID(uvn)), jen.ID("retrieveErr").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("retrieveErr")),
			),
			jen.Line(),
			jen.Return().List(jen.ID(uvn), jen.ID("nil")),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetSomethingsRequestFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	tp := typ.Name.Plural() // title plural

	urlBuildingParams := buildV1ClientURLBuildingParamsForListOfSomething(
		pkg,
		jen.ID("filter").Dot("ToValues").Call(),
		typ,
	)

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("BuildGet%sRequest builds an HTTP request for fetching %s", tp, typ.Name.PluralCommonName())),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildGet%sRequest", tp)).Params(buildParamsForMethodThatRetrievesAListOfADataType(pkg, typ, false)...).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Callln(urlBuildingParams...),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildGetListOfSomethingsFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	tp := typ.Name.Plural() // title plural
	pvn := typ.Name.PluralUnexportedVarName()

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("Get%s retrieves a list of %s", tp, typ.Name.PluralCommonName())),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Get%s", tp)).Params(buildParamsForMethodThatRetrievesAListOfADataType(pkg, typ, false)...).Params(
			jen.ID(pvn).Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", ts)),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", tp)).Call(
				buildParamsForMethodThatRetrievesAListOfADataType(pkg, typ, true)...,
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.Line(),
			jen.If(
				jen.ID("retrieveErr").Op(":=").ID("c").Dot("retrieve").Call(jen.ID("ctx"), jen.ID("req"), jen.Op("&").ID(pvn)),
				jen.ID("retrieveErr").Op("!=").ID("nil"),
			).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("retrieveErr")),
			),
			jen.Line(),
			jen.Return().List(jen.ID(pvn), jen.ID("nil")),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildCreateSomethingRequestFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("BuildCreate%sRequest builds an HTTP request for creating %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildCreate%sRequest", ts)).Params(
			buildParamsForMethodThatCreatesADataType(pkg, typ, false)...,
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Callln(
				buildV1ClientURLBuildingParamsForListOfSomething(
					pkg,
					jen.Nil(),
					typ,
				)...,
			),
			jen.Line(),
			jen.Return().ID("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("input"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildCreateSomethingFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	vn := typ.Name.UnexportedVarName()

	commonName := strings.Join(strings.Split(typ.Name.RouteName(), "_"), " ")
	commonNameWithPrefix := fmt.Sprintf("%s %s", wordsmith.AOrAn(ts), commonName)

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("Create%s creates %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Create%s", ts)).Params(
			buildParamsForMethodThatCreatesADataType(pkg, typ, false)...,
		).Params(
			jen.ID(vn).Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), ts),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot(fmt.Sprintf("BuildCreate%sRequest", ts)).Call(
				buildParamsForMethodThatCreatesADataType(pkg, typ, true)...,
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.Line(),
			jen.ID("err").Op("=").ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID(vn),
			),
			jen.Return().List(
				jen.ID(vn),
				jen.ID("err"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildUpdateSomethingRequestFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("BuildUpdate%sRequest builds an HTTP request for updating %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildUpdate%sRequest", ts)).Params(
			buildParamsForMethodThatIncludesItsOwnTypeInItsParams(pkg, typ, false)...,
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Callln(
				buildV1ClientURLBuildingParamsForMethodThatIncludesItsOwnType(pkg, jen.Nil(), typ)...,
			),
			jen.Line(),
			jen.Return().ID("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPut"),
				jen.ID("uri"),
				jen.ID(typ.Name.UnexportedVarName()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateSomethingFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("Update%s updates %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Update%s", ts)).Params(
			buildParamsForMethodThatIncludesItsOwnTypeInItsParams(pkg, typ, false)...,
		).Params(jen.ID("error")).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot(fmt.Sprintf("BuildUpdate%sRequest", ts)).Call(
				buildParamsForMethodThatIncludesItsOwnTypeInItsParams(pkg, typ, true)...,
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.ID("err"),
				),
			),
			jen.Line(),
			jen.Return().ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID(typ.Name.UnexportedVarName()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildArchiveSomethingRequestFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()

	commonName := strings.Join(strings.Split(typ.Name.RouteName(), "_"), " ")
	commonNameWithPrefix := fmt.Sprintf("%s %s", wordsmith.AOrAn(ts), commonName)

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("BuildArchive%sRequest builds an HTTP request for updating %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildArchive%sRequest", ts)).Params(buildParamsForMethodThatHandlesAnInstanceWithIDs(pkg, typ, false)...).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Callln(
				buildV1ClientURLBuildingParamsForSingleInstanceOfSomething(
					pkg,
					jen.Nil(),
					jen.IDf("%sID", typ.Name.UnexportedVarName()),
					typ,
				)...,
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodDelete"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildArchiveSomethingFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()

	commonName := strings.Join(strings.Split(typ.Name.RouteName(), "_"), " ")
	commonNameWithPrefix := fmt.Sprintf("%s %s", wordsmith.AOrAn(ts), commonName)

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("Archive%s archives %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Archive%s", ts)).Params(buildParamsForMethodThatHandlesAnInstanceWithIDs(pkg, typ, false)...).Params(jen.ID("error")).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot(fmt.Sprintf("BuildArchive%sRequest", ts)).Call(
				buildParamsForMethodThatHandlesAnInstanceWithIDs(pkg, typ, true)...,
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.ID("err"),
				),
			),
			jen.Line(),
			jen.Return().ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			),
		),
	}

	return lines
}
