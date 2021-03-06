package client

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile("dbclient")

	utils.AddImports(proj, code)

	n := typ.Name
	sn := n.Singular()

	code.Add(
		jen.Var().Underscore().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sDataManager", sn)).Equals().Parens(jen.PointerTo().ID("Client")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(buildSomethingExists(proj, typ)...)
	code.Add(buildGetSomething(proj, typ)...)
	code.Add(buildGetAllSomethingCount(proj, typ)...)
	code.Add(buildGetAllSomething(proj, typ)...)
	code.Add(buildGetListOfSomething(proj, typ)...)
	code.Add(buildGetListOfSomethingWithIDs(proj, typ)...)
	code.Add(buildCreateSomething(proj, typ)...)
	code.Add(buildUpdateSomething(proj, typ)...)
	code.Add(buildArchiveSomething(proj, typ)...)

	return code
}

func buildTracerAttachmentsForMethodWithParents(proj *models.Project, typ models.DataType) []jen.Code {
	owners := proj.FindOwnerTypeChain(typ)
	tp := proj.InternalTracingV1Package()

	out := []jen.Code{}
	for _, o := range owners {
		out = append(out, jen.Qual(tp, fmt.Sprintf("Attach%sIDToSpan", o.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", o.Name.UnexportedVarName())))
	}
	out = append(out, jen.Qual(tp, fmt.Sprintf("Attach%sIDToSpan", typ.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", typ.Name.UnexportedVarName())))

	if typ.RestrictedToUserAtSomeLevel(proj) {
		out = append(out, jen.Qual(tp, "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)))
	}

	return out
}

func buildTracerAttachmentsForListMethodWithParents(proj *models.Project, typ models.DataType) []jen.Code {
	owners := proj.FindOwnerTypeChain(typ)
	tp := proj.InternalTracingV1Package()

	out := []jen.Code{}
	for _, o := range owners {
		out = append(out, jen.Qual(tp, fmt.Sprintf("Attach%sIDToSpan", o.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", o.Name.UnexportedVarName())))
	}
	out = append(out, jen.Qual(tp, "AttachFilterToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.FilterVarName)))

	if typ.RestrictedToUserAtSomeLevel(proj) {
		out = append(out, jen.Qual(tp, "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)))
	}

	return out
}

func buildSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	funcName := fmt.Sprintf("%sExists", sn)
	params := typ.BuildDBClientExistenceMethodParams(proj)
	args := typ.BuildDBClientExistenceMethodCallArgs(proj)

	block := append(
		[]jen.Code{utils.StartSpan(proj, true, funcName), jen.Line()},
		buildTracerAttachmentsForMethodWithParents(proj, typ)...,
	)

	block = append(block,
		jen.Line(),
		jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValues").Call(
			typ.BuildGetSomethingLogValues(proj),
		).Dot("Debug").Call(jen.Litf("%s called", funcName)),
		jen.Line(),
		jen.Return(jen.ID("c").Dot("querier").Dotf("%sExists", sn).Call(args...)),
	)

	lines := []jen.Code{
		jen.Commentf("%s fetches whether or not %s exists from the database.", funcName, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID(funcName).Params(
			params...,
		).Params(jen.Bool(), jen.Error()).Body(
			block...,
		),
	}

	return lines
}

func buildGetSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	scnwp := n.SingularCommonNameWithPrefix()

	funcName := fmt.Sprintf("Get%s", sn)
	params := typ.BuildDBClientRetrievalMethodParams(proj)
	args := typ.BuildDBClientRetrievalMethodCallArgs(proj)
	loggerValues := typ.BuildGetSomethingLogValues(proj)

	block := append(
		[]jen.Code{utils.StartSpan(proj, true, funcName), jen.Line()},
		buildTracerAttachmentsForMethodWithParents(proj, typ)...,
	)

	block = append(block,
		jen.Line(),
		jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValues").Call(
			loggerValues,
		).Dot("Debug").Call(jen.Litf("Get%s called", sn)),
		jen.Line(),
		jen.Return().ID("c").Dot("querier").Dotf("Get%s", sn).Call(args...),
	)

	return []jen.Code{
		jen.Commentf("Get%s fetches %s from the database.", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID(funcName).Params(params...).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), sn), jen.Error()).Body(block...),
		jen.Line(),
	}
}

func buildGetAllSomethingCount(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	pn := n.Plural()
	pcn := n.PluralCommonName()

	return []jen.Code{
		jen.Commentf("GetAll%sCount fetches the count of %s from the database that meet a particular filter.", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("GetAll%sCount", pn).Params(constants.CtxParam()).Params(jen.ID("count").Uint64(), jen.Err().Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Litf("GetAll%sCount", pn)),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Litf("GetAll%sCount called", pn)),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dotf("GetAll%sCount", pn).Call(constants.CtxVar()),
		),
		jen.Line(),
	}
}

func buildGetAllSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	pn := n.Plural()
	sn := n.Singular()
	pcn := n.PluralCommonName()

	return []jen.Code{
		jen.Commentf("GetAll%s fetches a list of all %s in the database.", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("GetAll%s", pn).Params(
			constants.CtxParam(),
			jen.ID("results").Chan().Index().Qual(proj.ModelsV1Package(), sn),
		).Params(
			jen.Error(),
		).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Litf("GetAll%s", pn)),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Litf("GetAll%s called", pn)),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dotf("GetAll%s", pn).Call(
				constants.CtxVar(),
				jen.ID("results"),
			),
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

	funcName := fmt.Sprintf("Get%s", pn)
	params := typ.BuildDBClientListRetrievalMethodParams(proj)
	callArgs := typ.BuildDBClientListRetrievalMethodCallArgs(proj)
	loggerValues := typ.BuildGetListOfSomethingLogValues(proj)

	block := append(
		[]jen.Code{utils.StartSpan(proj, true, funcName), jen.Line()},
		buildTracerAttachmentsForListMethodWithParents(proj, typ)...,
	)

	logCall := jen.ID("c").Dot(constants.LoggerVarName)
	if loggerValues != nil {
		logCall = logCall.Dot("WithValues").Call(loggerValues)
	}
	logCall = logCall.Dot("Debug").Call(jen.Litf("Get%s called", pn))
	block = append(block, jen.Line(), logCall)

	block = append(block,
		jen.Line(),
		jen.List(jen.IDf("%sList", uvn), jen.Err()).Assign().
			ID("c").Dot("querier").Dotf("Get%s", pn).Call(callArgs...),
		jen.Line(),
		jen.Return().List(jen.IDf("%sList", uvn), jen.Err()),
	)

	return []jen.Code{
		jen.Commentf("Get%s fetches a list of %s from the database that meet a particular filter.", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("Get%s", pn).Params(params...).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn)), jen.Error()).Body(block...),
		jen.Line(),
	}
}

func buildGetListOfSomethingWithIDs(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	uvn := n.UnexportedVarName()
	pn := n.Plural()
	pcn := n.PluralCommonName()

	funcName := fmt.Sprintf("Get%sWithIDs", pn)
	block := append(
		[]jen.Code{utils.StartSpan(proj, true, funcName), jen.Line()},
	)

	if typ.RestrictedToUserAtSomeLevel(proj) {
		block = append(block,
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("userID")),
		)
	}

	logCall := jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValues").Call(
		jen.Map(jen.String()).Interface().Valuesln(
			func() jen.Code {
				if typ.RestrictedToUserAtSomeLevel(proj) {
					return jen.Lit("user_id").MapAssign().ID("userID")
				}
				return jen.Null()
			}(),
			jen.Lit("id_count").MapAssign().Len(jen.ID("ids")),
		),
	).Dot("Debug").Call(jen.Litf("Get%sWithIDs called", pn))
	block = append(block, jen.Line(), logCall)

	args := typ.BuildGetListOfSomethingFromIDsArgs(proj)

	block = append(block,
		jen.Line(),
		jen.List(jen.IDf("%sList", uvn), jen.Err()).Assign().
			ID("c").Dot("querier").Dotf("Get%sWithIDs", pn).Call(
			args...,
		),
		jen.Line(),
		jen.Return().List(jen.IDf("%sList", uvn), jen.Err()),
	)

	params := typ.BuildGetListOfSomethingFromIDsParams(proj)

	return []jen.Code{
		jen.Commentf("Get%sWithIDs fetches %s from the database within a given set of IDs.", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("Get%sWithIDs", pn).Params(
			params...,
		).Params(jen.Index().Qual(proj.ModelsV1Package(), sn), jen.Error()).Body(block...),
		jen.Line(),
	}
}

func buildCreateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	scnwp := n.SingularCommonNameWithPrefix()

	params := typ.BuildDBClientCreationMethodParams(proj)
	args := typ.BuildDBClientCreationMethodCallArgs()

	return []jen.Code{
		jen.Commentf("Create%s creates %s in the database.", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("Create%s", sn).Params(
			params...,
		).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), sn),
			jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Litf("Create%s", sn)),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("input"), jen.ID("input")).Dot("Debug").Call(jen.Litf("Create%s called", sn)),
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

	params := typ.BuildDBClientUpdateMethodParams(proj, updatedVarName)
	args := typ.BuildDBClientUpdateMethodCallArgs(updatedVarName)

	return []jen.Code{
		jen.Commentf("Update%s updates a particular %s. Note that Update%s expects the", sn, scn, sn),
		jen.Line(),
		jen.Comment("provided input to have a valid ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("Update%s", sn).Params(
			params...,
		).Params(jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Litf("Update%s", sn)),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID(constants.SpanVarName), jen.ID(updatedVarName).Dot("ID")),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Litf("%s_id", rn), jen.ID(updatedVarName).Dot("ID")).Dot("Debug").Call(jen.Litf("Update%s called", sn)),
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

	params := typ.BuildDBClientArchiveMethodParams()
	callArgs := typ.BuildDBClientArchiveMethodCallArgs()
	loggerValues := []jen.Code{
		jen.Litf("%s_id", rn).MapAssign().IDf("%sID", uvn),
	}

	if typ.BelongsToUser {
		loggerValues = append(loggerValues, jen.Lit("user_id").MapAssign().ID(constants.UserIDVarName))
	}
	if typ.BelongsToStruct != nil {
		loggerValues = append(loggerValues, jen.Litf("%s_id", typ.BelongsToStruct.RouteName()).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}

	block := []jen.Code{
		jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Litf("Archive%s", sn)),
		jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
		jen.Line(),
	}

	if typ.BelongsToUser {
		block = append(block, jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)))
	}
	if typ.BelongsToStruct != nil {
		block = append(block, jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", typ.BelongsToStruct.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())))
	}

	block = append(block,
		jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", uvn)),
		jen.Line(),
		jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(loggerValues...)).Dot("Debug").Call(jen.Litf("Archive%s called", sn)),
		jen.Line(),
		jen.Return().ID("c").Dot("querier").Dotf("Archive%s", sn).Call(callArgs...),
	)
	// we don't need to worry about the blonging to nobody case

	return []jen.Code{
		jen.Commentf("Archive%s archives %s from the database by its ID.", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("Archive%s", sn).Params(params...).Params(jen.Error()).Body(block...),
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
		jen.Commentf("GetAll%sForUser fetches a list of %s from the database that meet a particular filter.", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).IDf("GetAll%sForUser", pn).Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.Index().Qual(proj.ModelsV1Package(), sn),
			jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Litf("GetAll%sForUser", pn)),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID(constants.UserIDVarName)).Dot("Debug").Call(jen.Litf("GetAll%sForUser called", pn)),
			jen.Line(),
			jen.List(jen.IDf("%sList", uvn), jen.Err()).Assign().ID("c").Dot("querier").Dotf("GetAll%sForUser", pn).Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)),
			jen.Line(),
			jen.Return().List(jen.IDf("%sList", uvn), jen.Err()),
		),
		jen.Line(),
	}
}
