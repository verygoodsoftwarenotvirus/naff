package load

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("main")

	utils.AddImports(pkg, ret)

	ret.Add(buildFetchRandomSomething(pkg, typ)...)
	ret.Add(buildRandomActionMap(pkg, typ)...)

	return ret
}

func buildParamsForMethodThatHandlesAnInstanceWithRetrievedStructs(pkg *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := pkg.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{utils.InlineCtx()}

	for _, pt := range parents {
		listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}

	if len(listParams) > 0 {
		if call {
			params = append(params, listParams...)
		} else {
			params = append(params, jen.List(listParams...).ID("uint64"))
		}
	}

	return params
}

func buildFetchRandomSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	puvn := typ.Name.PluralUnexportedVarName()

	x := buildParamsForMethodThatHandlesAnInstanceWithRetrievedStructs(pkg, typ, false)
	x = x[1:]
	paramArgs := append(
		[]jen.Code{
			jen.ID("c").Op("*").Qual(pkg.HTTPClientV1Package(), "V1Client"),
		},
		x...,
	)

	callArgs := append(buildParamsForMethodThatHandlesAnInstanceWithRetrievedStructs(pkg, typ, true), jen.Nil())

	lines := []jen.Code{
		jen.Commentf("fetchRandom%s retrieves a random %s from the list of available %s", sn, scn, pcn),
		jen.Line(),
		jen.Func().IDf("fetchRandom%s", sn).Params(paramArgs...).Params(jen.Op("*").Qual(pkg.ModelsV1Package(), sn)).Block(
			jen.List(jen.IDf("%sRes", puvn), jen.Err()).Assign().ID("c").Dotf("Get%s", pn).Call(
				callArgs...,
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil").Op("||").IDf("%sRes", puvn).Op("==").ID("nil").Op("||").ID("len").Call(jen.IDf("%sRes", puvn).Dot(pn)).Op("==").Lit(0)).Block(
				jen.Return().ID("nil"),
			),
			jen.Line(),
			jen.ID("randIndex").Assign().Qual("math/rand", "Intn").Call(jen.ID("len").Call(jen.IDf("%sRes", puvn).Dot(pn))),
			jen.Return().VarPointer().IDf("%sRes", puvn).Dot(pn).Index(jen.ID("randIndex")),
		),
		jen.Line(),
	}

	return lines
}

func buildCreationArguments(pkg *models.Project, varPrefix string, typ models.DataType) []jen.Code {
	creationArgs := []jen.Code{}

	if typ.BelongsToStruct != nil {
		parentTyp := pkg.FindType(typ.BelongsToStruct.Singular())
		if parentTyp != nil {
			nca := buildCreationArguments(pkg, varPrefix, *parentTyp)
			creationArgs = append(creationArgs, nca...)
		}
	}

	creationArgs = append(creationArgs, jen.IDf("%s%s", varPrefix, typ.Name.Singular()).Dot("ID"))

	return creationArgs
}

func fieldToExpectedDotField(varName string, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	for _, field := range typ.Fields {
		sn := field.Name.Singular()
		lines = append(lines, jen.ID(sn).MapAssign().ID(varName).Dot(sn))
	}

	return lines
}

func buildFakeCallForCreationInput(pkg *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	for _, field := range typ.Fields {
		lines = append(lines, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeCallForField(pkg.OutputPath, field)))
	}

	return lines
}

func buildRequisiteCreationCode(pkg *models.Project, typ models.DataType) []jen.Code {
	var lines []jen.Code
	sn := typ.Name.Singular()

	const (
		sourceVarPrefix  = "example"
		createdVarPrefix = "created"
	)

	creationArgs := []jen.Code{
		utils.InlineCtx(),
	}
	ca := buildCreationArguments(pkg, createdVarPrefix, typ)
	creationArgs = append(creationArgs, ca[:len(ca)-1]...)
	creationArgs = append(creationArgs,
		jen.VarPointer().Qual(pkg.ModelsV1Package(), fmt.Sprintf("%sCreationInput", sn)).Valuesln(
			fieldToExpectedDotField(fmt.Sprintf("%s%s", sourceVarPrefix, typ.Name.Singular()), typ)...,
		),
	)

	if typ.BelongsToStruct != nil {
		if parentTyp := pkg.FindType(typ.BelongsToStruct.Singular()); parentTyp != nil {
			newLines := buildRequisiteCreationCode(pkg, *parentTyp)
			lines = append(lines, newLines...)
		}
	}

	lines = append(lines,
		jen.Commentf("Create %s", typ.Name.SingularCommonName()),
		jen.IDf("%s%s", sourceVarPrefix, typ.Name.Singular()).Assign().VarPointer().Qual(pkg.ModelsV1Package(), typ.Name.Singular()).Valuesln(
			buildFakeCallForCreationInput(pkg, typ)...,
		),
		jen.Line(),
		jen.List(jen.IDf("%s%s", createdVarPrefix, typ.Name.Singular()), jen.Err()).Assign().ID("c").Dotf("Create%s", sn).Call(
			creationArgs...,
		),
		jen.If(jen.Err().DoesNotEqual().Nil()).Block(
			jen.Return(jen.Nil(), jen.Err()),
		),
		jen.Line(),
	)

	return lines
}

func buildParamsForMethodThatHandlesAnInstanceWithStructs(pkg *models.Project, typ models.DataType) []jen.Code {
	parents := pkg.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{utils.InlineCtx()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("created%s", pt.Name.Singular()).Dot("ID"))
		}
		listParams = append(listParams, jen.IDf("created%s", typ.Name.Singular()).Dot("ID"))

		params = append(params, listParams...)

	} else {
		params = append(params, jen.IDf("created%s", typ.Name.Singular()).Dot("ID"))
	}

	return params
}

func buildCallArgsForMethodThatHandlesAnInstanceWithRetrievedStructs(pkg *models.Project, typ models.DataType) []jen.Code {
	parents := pkg.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{utils.InlineCtx()}

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

func buildRandomActionMap(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	blockLines := []jen.Code{
		jen.Return().Map(jen.ID("string")).Op("*").ID("Action").Valuesln(
			jen.Litf("Create%s", sn).MapAssign().Valuesln(
				jen.ID("Name").MapAssign().Litf("Create%s", sn),
				jen.ID("Action").MapAssign().Func().Params().Params(jen.ParamPointer().Qual("net/http", "Request"), jen.ID("error")).Block(
					buildCreateSomethingBlock(pkg, typ)...,
				),
				jen.ID("Weight").MapAssign().Lit(100),
			),
			jen.Litf("Get%s", sn).MapAssign().Valuesln(
				jen.ID("Name").MapAssign().Litf("Get%s", sn),
				jen.ID("Action").MapAssign().Func().Params().Params(jen.ParamPointer().Qual("net/http", "Request"), jen.ID("error")).Block(
					buildGetSomethingBlock(pkg, typ)...,
				),
				jen.ID("Weight").MapAssign().Lit(100),
			),
			jen.Litf("Get%s", pn).MapAssign().Valuesln(
				jen.ID("Name").MapAssign().Litf("Get%s", pn),
				jen.ID("Action").MapAssign().Func().Params().Params(jen.ParamPointer().Qual("net/http", "Request"), jen.ID("error")).Block(
					buildGetListOfSomethingBlock(pkg, typ)...,
				),
				jen.ID("Weight").MapAssign().Lit(100),
			),
			jen.Litf("Update%s", sn).MapAssign().Valuesln(
				jen.ID("Name").MapAssign().Litf("Update%s", sn),
				jen.ID("Action").MapAssign().Func().Params().Params(jen.ParamPointer().Qual("net/http", "Request"), jen.ID("error")).Block(
					buildUpdateChildBlock(pkg, typ)...,
				),
				jen.ID("Weight").MapAssign().Lit(100),
			),
			jen.Litf("Archive%s", sn).MapAssign().Valuesln(
				jen.ID("Name").MapAssign().Litf("Archive%s", sn),
				jen.ID("Action").MapAssign().Func().Params().Params(jen.ParamPointer().Qual("net/http", "Request"), jen.ID("error")).Block(
					buildArchiveSomethingBlock(pkg, typ)...,
				),
				jen.ID("Weight").MapAssign().Lit(85),
			),
		),
	}

	return []jen.Code{
		jen.Func().IDf("build%sActions", sn).Params(jen.ID("c").Op("*").Qual(pkg.HTTPClientV1Package(), "V1Client")).Params(jen.Map(jen.ID("string")).Op("*").ID("Action")).Block(blockLines...),
		jen.Line(),
	}
}

func buildCreateSomethingBlock(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	stopIndex := 6
	lines := buildRequisiteCreationCode(pkg, typ)
	if len(lines) >= stopIndex {
		lines = lines[:len(lines)-stopIndex]
	}

	args := buildParamsForMethodThatHandlesAnInstanceWithStructs(pkg, typ)
	args = args[:len(args)-1]
	args = append(args, jen.Qual(pkg.RandomModelsPackage(), fmt.Sprintf("Random%sCreationInput", sn)).Call())

	lines = append(lines,
		jen.Return(jen.ID("c").Dotf("BuildCreate%sRequest", sn).Call(args...)),
	)

	return lines
}

func buildRandomDependentIDFetchers(pkg *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}
	parentTypes := pkg.FindOwnerTypeChain(typ)

	callArgs := []jen.Code{
		jen.ID("c"),
	}

	for _, pt := range parentTypes {
		ca := buildCallArgsForMethodThatHandlesAnInstanceWithRetrievedStructs(pkg, pt)
		ca = ca[1 : len(ca)-1]
		callArgs = append([]jen.Code{jen.ID("c")}, ca...)

		lines = append(lines,
			jen.ID(pt.Name.UnexportedVarName()).Assign().IDf("fetchRandom%s", pt.Name.Singular()).Call(callArgs...),
		)
	}

	return lines
}

func buildGetSomethingBlock(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	ca := buildCallArgsForMethodThatHandlesAnInstanceWithRetrievedStructs(pkg, typ)
	ca = ca[1 : len(ca)-1]
	callArgs := append([]jen.Code{jen.ID("c")}, ca...)

	requestBuildingArgs := append([]jen.Code{utils.InlineCtx()}, ca...)
	requestBuildingArgs = append(requestBuildingArgs, jen.IDf("random%s", sn).Dot("ID"))

	lines := buildRandomDependentIDFetchers(pkg, typ)
	lines = append(lines,
		func() jen.Code {
			if len(lines) > 0 {
				return jen.Line()
			}
			return nil
		}(),
		jen.If(jen.IDf("random%s", sn).Assign().IDf("fetchRandom%s", sn).Call(callArgs...), jen.IDf("random%s", sn).DoesNotEqual().ID("nil")).Block(
			jen.Return().ID("c").Dotf("BuildGet%sRequest", sn).Call(requestBuildingArgs...),
		),
		jen.Line(),
		jen.Return().List(jen.Nil(), jen.ID("ErrUnavailableYet")),
	)

	return lines
}

func buildGetListOfSomethingBlock(pkg *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()

	ca := buildCallArgsForMethodThatHandlesAnInstanceWithRetrievedStructs(pkg, typ)
	ca = ca[1 : len(ca)-1]

	requestBuildingArgs := append([]jen.Code{utils.InlineCtx()}, ca...)
	requestBuildingArgs = append(requestBuildingArgs, jen.Nil())

	lines := buildRandomDependentIDFetchers(pkg, typ)
	lines = append(lines,
		func() jen.Code {
			if len(lines) > 0 {
				return jen.Line()
			}
			return nil
		}(),
		jen.Return().ID("c").Dotf("BuildGet%sRequest", pn).Call(requestBuildingArgs...),
	)

	return lines
}

func buildUpdateChildBlock(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	ca := buildCallArgsForMethodThatHandlesAnInstanceWithRetrievedStructs(pkg, typ)
	ca = ca[1 : len(ca)-1]
	callArgs := append([]jen.Code{jen.ID("c")}, ca...)

	requestBuildingArgs := append([]jen.Code{utils.InlineCtx()}, ca...)
	if len(requestBuildingArgs) > 1 {
		requestBuildingArgs = requestBuildingArgs[:len(requestBuildingArgs)-1]
	}
	requestBuildingArgs = append(requestBuildingArgs, jen.IDf("random%s", sn))

	var ifRandomExistsBlock []jen.Code
	for _, field := range typ.Fields {
		fsn := field.Name.Singular()
		ifRandomExistsBlock = append(ifRandomExistsBlock, jen.IDf("random%s", sn).Dot(fsn).Equals().Qual(pkg.RandomModelsPackage(), fmt.Sprintf("Random%sCreationInput", sn)).Call().Dot(fsn))
	}
	ifRandomExistsBlock = append(ifRandomExistsBlock,
		jen.Line(),
		jen.Return().ID("c").Dotf("BuildUpdate%sRequest", sn).Call(
			requestBuildingArgs...,
		),
	)

	lines := buildRandomDependentIDFetchers(pkg, typ)
	if len(lines) > 0 {
		lines = append(lines, jen.Line())
	}

	lines = append(lines,
		jen.If(jen.IDf("random%s", sn).Assign().IDf("fetchRandom%s", sn).Call(callArgs...), jen.IDf("random%s", sn).DoesNotEqual().ID("nil")).Block(
			ifRandomExistsBlock...,
		),
		jen.Line(),
		jen.Return().List(jen.Nil(), jen.ID("ErrUnavailableYet")),
	)

	return lines

	//sn := typ.Name.Singular()
	//
	//ca := buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(pkg, typ)
	//ca = ca[1 : len(ca)-1]
	//callArgs := append([]jen.Code{jen.ID("c")}, ca...)
	//
	//requestBuildingArgs := append([]jen.Code{utils.InlineCtx()}, ca...)
	//requestBuildingArgs = append(requestBuildingArgs, jen.IDf("random%s", sn).Dot("ID"))
	//
	//var randomLines []jen.Code
	//
	//for _, field := range typ.Fields {
	//	fsn := field.Name.Singular()
	//	randomLines = append(randomLines, jen.IDf("random%s", sn).Dot(fsn).Equals().Qual(pkg.RandomModelsPackage(), fmt.Sprintf("Random%sCreationInput", sn)).Call().Dot(fsn))
	//}
	//
	//randomLines = append(randomLines,
	//	jen.Return().ID("c").Dotf("BuildUpdate%sRequest", sn).Call(
	//		requestBuildingArgs...,
	//	// 	utils.InlineCtx(), jen.IDf("random%s", sn),
	//	),
	//)
	//
	//lines := buildRandomDependentIDFetchers(pkg, typ)
	//lines = append(lines,
	//	jen.Line(),
	//	jen.If(jen.IDf("random%s", sn).Assign().IDf("fetchRandom%s", sn).Call(callArgs...), jen.IDf("random%s", sn).DoesNotEqual().ID("nil")).Block(
	//		randomLines...,
	//	),
	//	jen.Line(),
	//	jen.Return().List(jen.Nil(), jen.ID("ErrUnavailableYet")),
	//)
	//
	//return lines
}

func buildArchiveSomethingBlock(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	ca := buildCallArgsForMethodThatHandlesAnInstanceWithRetrievedStructs(pkg, typ)
	ca = ca[1 : len(ca)-1]
	callArgs := append([]jen.Code{jen.ID("c")}, ca...)

	requestBuildingArgs := append([]jen.Code{utils.InlineCtx()}, ca...)
	requestBuildingArgs = append(requestBuildingArgs, jen.IDf("random%s", sn).Dot("ID"))

	lines := buildRandomDependentIDFetchers(pkg, typ)
	lines = append(lines,
		func() jen.Code {
			if len(lines) > 0 {
				return jen.Line()
			}
			return nil
		}(),
		jen.If(jen.IDf("random%s", sn).Assign().IDf("fetchRandom%s", sn).Call(callArgs...), jen.IDf("random%s", sn).DoesNotEqual().ID("nil")).Block(
			jen.Return().ID("c").Dotf("BuildArchive%sRequest", sn).Call(requestBuildingArgs...),
		),
		jen.Line(),
		jen.Return().List(jen.Nil(), jen.ID("ErrUnavailableYet")),
	)

	return lines
}
