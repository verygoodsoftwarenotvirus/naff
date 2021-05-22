package load

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildFetchRandomSomething(proj, typ)...)
	code.Add(buildRandomActionMap(proj, typ)...)

	return code
}

func buildParamsForMethodThatHandlesAnInstanceWithRetrievedStructs(proj *models.Project, typ models.DataType, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	for _, pt := range parents {
		listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}

	if len(listParams) > 0 {
		if call {
			params = append(params, listParams...)
		} else {
			params = append(params, jen.List(listParams...).Uint64())
		}
	}

	return params
}

func buildFetchRandomSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	puvn := typ.Name.PluralUnexportedVarName()

	x := buildParamsForMethodThatHandlesAnInstanceWithRetrievedStructs(proj, typ, false)
	x = x[1:]
	paramArgs := append(
		[]jen.Code{
			constants.CtxParam(),
			jen.ID("c").PointerTo().Qual(proj.HTTPClientV1Package(), "V1Client"),
		},
		x...,
	)

	callArgs := append(buildParamsForMethodThatHandlesAnInstanceWithRetrievedStructs(proj, typ, true), jen.Nil())

	lines := []jen.Code{
		jen.Commentf("fetchRandom%s retrieves a random %s from the list of available %s.", sn, scn, pcn),
		jen.Line(),
		jen.Func().IDf("fetchRandom%s", sn).Params(paramArgs...).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)).Body(
			jen.List(jen.IDf("%sRes", puvn), jen.Err()).Assign().ID("c").Dotf("Get%s", pn).Call(
				callArgs...,
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil").Or().IDf("%sRes", puvn).IsEqualTo().ID("nil").Or().ID("len").Call(jen.IDf("%sRes", puvn).Dot(pn)).IsEqualTo().Zero()).Body(
				jen.Return().ID("nil"),
			),
			jen.Line(),
			jen.ID("randIndex").Assign().Qual("math/rand", "Intn").Call(jen.Len(jen.IDf("%sRes", puvn).Dot(pn))),
			jen.Return().AddressOf().IDf("%sRes", puvn).Dot(pn).Index(jen.ID("randIndex")),
		),
		jen.Line(),
	}

	return lines
}

func buildCreationArguments(proj *models.Project, varPrefix string, typ models.DataType) []jen.Code {
	creationArgs := []jen.Code{constants.CtxVar()}

	owners := proj.FindOwnerTypeChain(typ)
	for i, ot := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			continue
		} else {
			creationArgs = append(creationArgs, jen.IDf("%s%s", varPrefix, ot.Name.Singular()).Dot("ID"))
		}
	}

	return creationArgs
}

func buildRequisiteCreationCode(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	const createdVarPrefix = "created"

	for _, t := range proj.FindOwnerTypeChain(typ) {
		sn := t.Name.Singular()

		creationArgs := append(buildCreationArguments(proj, "created", t), jen.IDf("example%sInput", sn))
		lines = append(lines,
			jen.Commentf("Create %s.", t.Name.SingularCommonName()),
			utils.BuildFakeVar(proj, sn),
			func() jen.Code {
				if t.BelongsToStruct != nil {
					return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", t.BelongsToStruct.Singular()).Equals().IDf("created%s", t.BelongsToStruct.Singular()).Dot("ID")
				}
				return jen.Null()
			}(),
			utils.BuildFakeVarWithCustomName(
				proj,
				utils.BuildFakeVarName(fmt.Sprintf("%sInput", sn)),
				fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn),
				jen.ID(utils.BuildFakeVarName(sn)),
			),
			jen.List(jen.IDf("%s%s", createdVarPrefix, t.Name.Singular()), jen.Err()).Assign().ID("c").Dotf("Create%s", sn).Call(
				creationArgs...,
			),
			jen.If(jen.Err().DoesNotEqual().Nil()).Body(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Line(),
		)

	}

	return lines
}

func buildParamsForMethodThatHandlesAnInstanceWithStructs(proj *models.Project, typ models.DataType) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for i, pt := range parents {
			if i == len(parents)-1 && typ.BelongsToStruct != nil {
				continue
			} else {
				listParams = append(listParams, jen.IDf("created%s", pt.Name.Singular()).Dot("ID"))
			}
		}
		listParams = append(listParams, jen.IDf("created%s", typ.Name.Singular()).Dot("ID"))

		params = append(params, listParams...)

	} else {
		params = append(params, jen.IDf("created%s", typ.Name.Singular()).Dot("ID"))
	}

	return params
}

func buildCallArgsForMethodThatHandlesAnInstanceWithRetrievedStructs(proj *models.Project, typ models.DataType) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("random%s", pt.Name.Singular()).Dot("ID"))
		}
		listParams = append(listParams, jen.IDf("random%s", typ.Name.Singular()).Dot("ID"))

		params = append(params, listParams...)
	} else {
		params = append(params, jen.IDf("random%s", typ.Name.Singular()).Dot("ID"))
	}

	return params
}

func buildRandomActionMap(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	blockLines := []jen.Code{
		jen.Return().Map(jen.String()).PointerTo().ID("Action").Valuesln(
			jen.Litf("Create%s", sn).MapAssign().Valuesln(
				jen.ID("Name").MapAssign().Litf("Create%s", sn),
				jen.ID("Action").MapAssign().Func().Params().Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()).Body(
					buildCreateSomethingBlock(proj, typ)...,
				),
				jen.ID("Weight").MapAssign().Lit(100),
			),
			jen.Litf("Get%s", sn).MapAssign().Valuesln(
				jen.ID("Name").MapAssign().Litf("Get%s", sn),
				jen.ID("Action").MapAssign().Func().Params().Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()).Body(
					buildGetSomethingBlock(proj, typ)...,
				),
				jen.ID("Weight").MapAssign().Lit(100),
			),
			jen.Litf("Get%s", pn).MapAssign().Valuesln(
				jen.ID("Name").MapAssign().Litf("Get%s", pn),
				jen.ID("Action").MapAssign().Func().Params().Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()).Body(
					buildGetListOfSomethingBlock(proj, typ)...,
				),
				jen.ID("Weight").MapAssign().Lit(100),
			),
			jen.Litf("Update%s", sn).MapAssign().Valuesln(
				jen.ID("Name").MapAssign().Litf("Update%s", sn),
				jen.ID("Action").MapAssign().Func().Params().Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()).Body(
					buildUpdateChildBlock(proj, typ)...,
				),
				jen.ID("Weight").MapAssign().Lit(100),
			),
			jen.Litf("Archive%s", sn).MapAssign().Valuesln(
				jen.ID("Name").MapAssign().Litf("Archive%s", sn),
				jen.ID("Action").MapAssign().Func().Params().Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()).Body(
					buildArchiveSomethingBlock(proj, typ)...,
				),
				jen.ID("Weight").MapAssign().Lit(85),
			),
		),
	}

	return []jen.Code{
		jen.Func().IDf("build%sActions", sn).Params(jen.ID("c").PointerTo().Qual(proj.HTTPClientV1Package(), "V1Client")).Params(jen.Map(jen.String()).PointerTo().ID("Action")).Body(blockLines...),
		jen.Line(),
	}
}

func buildCreateSomethingBlock(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()

	lines := append([]jen.Code{constants.CreateCtx(), jen.Line()}, buildRequisiteCreationCode(proj, typ)...)

	args := buildParamsForMethodThatHandlesAnInstanceWithStructs(proj, typ)
	args = args[:len(args)-1]
	args = append(args, jen.IDf("%sInput", uvn))

	lines = append(lines,
		utils.BuildFakeVarWithCustomName(
			proj,
			fmt.Sprintf("%sInput", uvn),
			fmt.Sprintf("BuildFake%sCreationInput", sn),
		),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("%sInput", uvn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("created%s", typ.BelongsToStruct.Singular()).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.Line(),
		jen.Return(jen.ID("c").Dotf("BuildCreate%sRequest", sn).Call(args...)),
	)

	return lines
}

func buildRandomDependentIDFetchers(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{constants.CreateCtx()}
	parentTypes := proj.FindOwnerTypeChain(typ)

	callArgs := []jen.Code{
		jen.ID("c"),
	}

	for _, pt := range parentTypes {
		ca := buildCallArgsForMethodThatHandlesAnInstanceWithRetrievedStructs(proj, pt)
		ca = ca[1 : len(ca)-1]
		callArgs = append([]jen.Code{constants.CtxVar(), jen.ID("c")}, ca...)

		lines = append(lines,
			jen.Line(),
			jen.IDf("random%s", pt.Name.Singular()).Assign().IDf("fetchRandom%s", pt.Name.Singular()).Call(callArgs...),
			jen.If(jen.IDf("random%s", pt.Name.Singular()).IsEqualTo().Nil()).Body(
				jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit(fmt.Sprintf("retrieving random %s", pt.Name.SingularCommonName())+": %w"), jen.ID("ErrUnavailableYet"))),
			),
		)
	}

	return lines
}

func buildGetSomethingBlock(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	ca := buildCallArgsForMethodThatHandlesAnInstanceWithRetrievedStructs(proj, typ)
	ca = ca[1 : len(ca)-1]
	callArgs := append([]jen.Code{constants.CtxVar(), jen.ID("c")}, ca...)

	requestBuildingArgs := append([]jen.Code{constants.CtxVar()}, ca...)
	requestBuildingArgs = append(requestBuildingArgs, jen.IDf("random%s", sn).Dot("ID"))

	lines := buildRandomDependentIDFetchers(proj, typ)

	lines = append(lines,
		jen.Line(),
		jen.IDf("random%s", typ.Name.Singular()).Assign().IDf("fetchRandom%s", typ.Name.Singular()).Call(callArgs...),
		jen.If(jen.IDf("random%s", typ.Name.Singular()).IsEqualTo().Nil()).Body(
			jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit(fmt.Sprintf("retrieving random %s", typ.Name.SingularCommonName())+": %w"), jen.ID("ErrUnavailableYet"))),
		),
	)

	lines = append(lines,
		jen.Line(),
		jen.Return().ID("c").Dotf("BuildGet%sRequest", sn).Call(requestBuildingArgs...),
	)

	return lines
}

func buildGetListOfSomethingBlock(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()

	ca := buildCallArgsForMethodThatHandlesAnInstanceWithRetrievedStructs(proj, typ)
	ca = ca[1 : len(ca)-1]

	requestBuildingArgs := append([]jen.Code{constants.CtxVar()}, ca...)
	requestBuildingArgs = append(requestBuildingArgs, jen.Nil())

	lines := buildRandomDependentIDFetchers(proj, typ)
	lines = append(lines,
		jen.Line(),
		jen.Return().ID("c").Dotf("BuildGet%sRequest", pn).Call(requestBuildingArgs...),
	)

	return lines
}

func buildUpdateChildBlock(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	ca := buildCallArgsForMethodThatHandlesAnInstanceWithRetrievedStructs(proj, typ)
	ca = ca[1 : len(ca)-1]
	callArgs := append([]jen.Code{constants.CtxVar(), jen.ID("c")}, ca...)

	requestBuildingArgs := append([]jen.Code{constants.CtxVar()}, ca...)
	if len(requestBuildingArgs) > 1 {
		requestBuildingArgs = requestBuildingArgs[:len(requestBuildingArgs)-1]
	}
	requestBuildingArgs = append(requestBuildingArgs, jen.IDf("random%s", sn))

	ifRandomExistsBlock := []jen.Code{
		jen.IDf("new%s", sn).Assign().
			Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sCreationInput", sn)).
			Call(),
	}
	for _, field := range typ.Fields {
		fsn := field.Name.Singular()
		ifRandomExistsBlock = append(
			ifRandomExistsBlock,
			jen.IDf("random%s", sn).Dot(fsn).Equals().IDf("new%s", sn).Dot(fsn),
		)
	}
	ifRandomExistsBlock = append(ifRandomExistsBlock,
		jen.Return().ID("c").Dotf("BuildUpdate%sRequest", sn).Call(
			requestBuildingArgs...,
		),
	)

	lines := buildRandomDependentIDFetchers(proj, typ)
	if len(lines) > 0 {
		lines = append(lines, jen.Line())
	}

	lines = append(lines,
		jen.If(jen.IDf("random%s", sn).Assign().IDf("fetchRandom%s", sn).Call(callArgs...), jen.IDf("random%s", sn).DoesNotEqual().ID("nil")).Body(
			ifRandomExistsBlock...,
		),
		jen.Line(),
		jen.Return().List(jen.Nil(), jen.ID("ErrUnavailableYet")),
	)

	return lines
}

func buildArchiveSomethingBlock(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	ca := buildCallArgsForMethodThatHandlesAnInstanceWithRetrievedStructs(proj, typ)
	ca = ca[1 : len(ca)-1]
	callArgs := append([]jen.Code{constants.CtxVar(), jen.ID("c")}, ca...)

	requestBuildingArgs := append([]jen.Code{constants.CtxVar()}, ca...)
	requestBuildingArgs = append(requestBuildingArgs, jen.IDf("random%s", sn).Dot("ID"))

	lines := buildRandomDependentIDFetchers(proj, typ)

	lines = append(lines,
		jen.Line(),
		jen.IDf("random%s", typ.Name.Singular()).Assign().IDf("fetchRandom%s", typ.Name.Singular()).Call(callArgs...),
		jen.If(jen.IDf("random%s", typ.Name.Singular()).IsEqualTo().Nil()).Body(
			jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit(fmt.Sprintf("retrieving random %s", typ.Name.SingularCommonName())+": %w"), jen.ID("ErrUnavailableYet"))),
		),
	)

	lines = append(lines,
		jen.Line(),
		jen.Return().ID("c").Dotf("BuildArchive%sRequest", sn).Call(requestBuildingArgs...),
	)

	return lines
}
