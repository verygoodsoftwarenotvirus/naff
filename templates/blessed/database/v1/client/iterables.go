package client

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(proj, ret)

	n := typ.Name
	sn := n.Singular()

	ret.Add(
		jen.Var().ID("_").Qual(proj.ModelsV1Package(), fmt.Sprintf("%sDataManager", sn)).Equals().Parens(jen.Op("*").ID("Client")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(buildSomethingExists(proj, typ)...)
	ret.Add(buildGetSomething(proj, typ)...)
	ret.Add(buildGetAllSomethingCount(proj, typ)...)
	ret.Add(buildGetListOfSomething(proj, typ)...)

	if typ.BelongsToStruct != nil {
		ret.Add(buildGetAllSomethingForSomethingElse(proj, typ)...)
	}

	ret.Add(buildCreateSomething(proj, typ)...)
	ret.Add(buildUpdateSomething(proj, typ)...)
	ret.Add(buildArchiveSomething(proj, typ)...)

	return ret
}

func buildSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	funcName := fmt.Sprintf("%sExists", sn)
	params := typ.BuildGetSomethingParams(proj)
	args := typ.BuildGetSomethingArgs(proj)

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.Line(),
		jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		jen.Qual(proj.InternalTracingV1Package(), "AttachItemIDToSpan").Call(jen.ID("span"), jen.ID("itemID")),
		jen.Line(),
		jen.ID("c").Dot("logger").Dot("WithValues").Call(
			typ.BuildGetSomethingLogValues(proj),
		).Dot("Debug").Call(jen.Litf("%s called", funcName)),
		jen.Line(),
		jen.Return(jen.ID("c").Dot("querier").Dotf("%sExists", sn).Call(args...)),
	}

	lines := []jen.Code{
		jen.Commentf("%s fetches whether or not %s exists from the database", funcName, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID(funcName).Params(
			params...,
		).Params(jen.Bool(), jen.ID("error")).Block(
			block...,
		),
	}

	return lines
}

func buildGetSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	uvn := n.UnexportedVarName()
	scnwp := n.SingularCommonNameWithPrefix()

	params := typ.BuildGetSomethingParams(proj)
	args := typ.BuildGetSomethingArgs(proj)
	loggerValues := typ.BuildGetSomethingLogValues(proj)
	block := []jen.Code{
		jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Litf("Get%s", sn)),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	if typ.BelongsToUser {
		block = append(block, jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")))
	}

	block = append(block,
		jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.IDf("%sID", uvn)),
		jen.Line(),
		jen.ID("c").Dot("logger").Dot("WithValues").Call(
			loggerValues,
		).Dot("Debug").Call(jen.Litf("Get%s called", sn)),
		jen.Line(),
		jen.Return().ID("c").Dot("querier").Dotf("Get%s", sn).Call(args...),
	)

	return []jen.Code{
		jen.Commentf("Get%s fetches %s from the database", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Get%s", sn).Params(params...).Params(jen.Op("*").Qual(proj.ModelsV1Package(), sn), jen.ID("error")).Block(block...),
		jen.Line(),
	}
}

func buildGetAllSomethingCount(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	pn := n.Plural()
	pcn := n.PluralCommonName()

	return []jen.Code{
		jen.Commentf("GetAll%sCount fetches the count of %s from the database that meet a particular filter", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("GetAll%sCount", pn).Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("count").ID("uint64"), jen.Err().ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Litf("GetAll%sCount", pn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Litf("GetAll%sCount called", pn)),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dotf("GetAll%sCount", pn).Call(utils.CtxVar()),
		),
		jen.Line(),
	}
}

func buildGetListOfSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	uvn := n.UnexportedVarName()
	pn := n.Plural()
	pcn := n.PluralCommonName()

	params := typ.BuildGetListOfSomethingParams(proj, false)
	callArgs := typ.BuildGetListOfSomethingArgs(proj)
	block := []jen.Code{
		jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Litf("Get%s", pn)),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	if typ.BelongsToUser {
		block = append(block,
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		)
	}
	if typ.BelongsToStruct != nil {
		block = append(block,
			jen.IDf("Attach%sIDToSpan", typ.BelongsToStruct.Singular()).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
		)
	}

	block = append(block,
		jen.Qual(proj.InternalTracingV1Package(), "AttachFilterToSpan").Call(jen.ID("span"), jen.ID(utils.FilterVarName)),
		jen.Line(),
	)

	if typ.BelongsToUser {
		block = append(block,
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Litf("Get%s called", pn)),
		)
	}
	if typ.BelongsToStruct != nil {
		block = append(block,
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Litf("%s_id", typ.BelongsToStruct.RouteName()), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())).Dot("Debug").Call(jen.Litf("Get%s called", pn)),
		)
	} else if typ.BelongsToNobody {
		block = append(block,
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Litf("Get%s called", pn)),
		)
	}

	block = append(block,
		jen.Line(),
		jen.List(jen.IDf("%sList", uvn), jen.Err()).Assign().ID("c").Dot("querier").Dotf("Get%s", pn).Call(callArgs...),
		jen.Line(),
		jen.Return().List(jen.IDf("%sList", uvn), jen.Err()),
	)

	return []jen.Code{
		jen.Commentf("Get%s fetches a list of %s from the database that meet a particular filter", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Get%s", pn).Params(params...).Params(jen.Op("*").Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn)), jen.ID("error")).Block(block...),
		jen.Line(),
	}
}

func buildGetAllSomethingForUser(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	uvn := n.UnexportedVarName()
	pn := n.Plural()
	pcn := n.PluralCommonName()

	return []jen.Code{
		jen.Commentf("GetAll%sForUser fetches a list of %s from the database that meet a particular filter", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("GetAll%sForUser", pn).Params(utils.CtxParam(), jen.ID("userID").ID("uint64")).Params(jen.Index().Qual(proj.ModelsV1Package(), sn),
			jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Litf("GetAll%sForUser", pn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Litf("GetAll%sForUser called", pn)),
			jen.Line(),
			jen.List(jen.IDf("%sList", uvn), jen.Err()).Assign().ID("c").Dot("querier").Dotf("GetAll%sForUser", pn).Call(utils.CtxVar(), jen.ID("userID")),
			jen.Line(),
			jen.Return().List(jen.IDf("%sList", uvn), jen.Err()),
		),
		jen.Line(),
	}
}

func buildGetAllSomethingForSomethingElse(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	uvn := n.UnexportedVarName()
	pn := n.Plural()
	pcn := n.PluralCommonName()

	btsn := typ.BelongsToStruct
	btsns := btsn.Singular()
	btsuvn := btsn.UnexportedVarName()
	btsrn := btsn.RouteName()

	params := typ.BuildGetSomethingForSomethingElseParams(proj)
	args := typ.BuildGetSomethingForSomethingElseArgs(proj)

	return []jen.Code{
		jen.Commentf("GetAll%sFor%s fetches a list of %s from the database that meet a particular filter", pn, btsns, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("GetAll%sFor%s", pn, btsns).Params(
			params...,
		).Params(jen.Index().Qual(proj.ModelsV1Package(), sn),
			jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Litf("GetAll%sFor%s", pn, btsns)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.IDf("Attach%sIDToSpan", btsns).Call(jen.ID("span"), jen.IDf("%sID", btsuvn)),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Litf("%s_id", btsrn), jen.IDf("%sID", btsuvn)).Dot("Debug").Call(jen.Litf("GetAll%sFor%s called", pn, btsns)),
			jen.Line(),
			jen.List(jen.IDf("%sList", uvn), jen.Err()).Assign().ID("c").Dot("querier").Dotf("GetAll%sFor%s", pn, btsns).Call(
				args...,
			),
			jen.Line(),
			jen.Return().List(jen.IDf("%sList", uvn), jen.Err()),
		),
		jen.Line(),
	}
}

func buildCreateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	scnwp := n.SingularCommonNameWithPrefix()

	params := typ.BuildCreateSomethingParams(proj, false)
	args := typ.BuildCreateSomethingArgs(proj)

	return []jen.Code{
		jen.Commentf("Create%s creates %s in the database", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Create%s", sn).Params(
			params...,
		).Params(jen.Op("*").Qual(proj.ModelsV1Package(), sn),
			jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Litf("Create%s", sn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("input"), jen.ID("input")).Dot("Debug").Call(jen.Litf("Create%s called", sn)),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dotf("Create%s", sn).Call(
				args...,
			),
		),
		jen.Line(),
	}
}

func buildUpdateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	rn := n.RouteName()
	scn := n.SingularCommonName()

	const updatedVarName = "updated"

	params := typ.BuildUpdateSomethingParams(proj, updatedVarName, false)
	args := typ.BuildUpdateSomethingArgs(proj, updatedVarName)

	return []jen.Code{
		jen.Commentf("Update%s updates a particular %s. Note that Update%s expects the", sn, scn, sn),
		jen.Line(),
		jen.Comment("provided input to have a valid ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Update%s", sn).Params(
			params...,
		).Params(jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Litf("Update%s", sn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.ID(updatedVarName).Dot("ID")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Litf("%s_id", rn), jen.ID(updatedVarName).Dot("ID")).Dot("Debug").Call(jen.Litf("Update%s called", sn)),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dotf("Update%s", sn).Call(
				args...,
			),
		),
		jen.Line(),
	}
}

func buildArchiveSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	rn := n.RouteName()
	uvn := n.UnexportedVarName()
	scnwp := n.SingularCommonNameWithPrefix()

	params := typ.BuildGetSomethingParams(proj)
	callArgs := typ.BuildGetSomethingArgs(proj)
	loggerValues := []jen.Code{
		jen.Litf("%s_id", rn).MapAssign().IDf("%sID", uvn),
	}

	if typ.BelongsToUser {
		loggerValues = append(loggerValues, jen.Lit("user_id").MapAssign().ID("userID"))
	}
	if typ.BelongsToStruct != nil {
		loggerValues = append(loggerValues, jen.Litf("%s_id", typ.BelongsToStruct.RouteName()).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}

	block := []jen.Code{
		jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Litf("Archive%s", sn)),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	if typ.BelongsToUser {
		block = append(block, jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")))
	}
	if typ.BelongsToStruct != nil {
		block = append(block, jen.IDf("Attach%sIDToSpan", typ.BelongsToStruct.Singular()).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())))
	}

	block = append(block,
		jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.IDf("%sID", uvn)),
		jen.Line(),
		jen.ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(loggerValues...)).Dot("Debug").Call(jen.Litf("Archive%s called", sn)),
		jen.Line(),
		jen.Return().ID("c").Dot("querier").Dotf("Archive%s", sn).Call(callArgs...),
	)
	// we don't need to worry about the blonging to nobody case

	return []jen.Code{
		jen.Commentf("Archive%s archives %s from the database by its ID", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Archive%s", sn).Params(params...).Params(jen.ID("error")).Block(block...),
		jen.Line(),
	}
}
