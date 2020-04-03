package iterables

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Commentf("URIParamKey is a standard string that we'll use to refer to %s IDs with", scn),
			jen.ID("URIParamKey").Equals().Lit(fmt.Sprintf("%sID", uvn)),
		),
		jen.Line(),
	)

	ret.Add(buildListHandlerFuncDecl(proj, typ)...)
	ret.Add(buildCreateHandlerFuncDecl(proj, typ)...)
	ret.Add(buildExistenceHandlerFuncDecl(proj, typ)...)
	ret.Add(buildReadHandlerFuncDecl(proj, typ)...)
	ret.Add(buildUpdateHandlerFuncDecl(proj, typ)...)
	ret.Add(buildArchiveHandlerFuncDecl(proj, typ)...)

	return ret
}

func buildListHandlerFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	puvn := typ.Name.PluralUnexportedVarName()
	pcn := typ.Name.PluralCommonName()
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	dbCallArgs := []jen.Code{
		utils.CtxVar(),
	}
	block := []jen.Code{
		jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ListHandler")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
		jen.Comment("ensure query filter"),
		jen.ID(utils.FilterVarName).Assign().Qual(proj.ModelsV1Package(), "ExtractQueryFilter").Call(jen.ID("req")),
		jen.Line(),
	}

	elseErrBlock := []jen.Code{}

	if typ.BelongsToUser {
		block = append(block,
			jen.Comment("determine user ID"),
			jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		)
		dbCallArgs = append(dbCallArgs, jen.ID("userID"))
		elseErrBlock = append(elseErrBlock,
			jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Litf("error encountered fetching %s", pcn)),
			utils.WriteXHeader("res", "StatusInternalServerError"),
			jen.Return(),
		)
	}
	if typ.BelongsToStruct != nil {
		block = append(block,
			jen.Commentf("determine %s ID", typ.BelongsToStruct.SingularCommonName()),
			jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Assign().ID("s").Dotf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", typ.BelongsToStruct.Singular())).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValue").Call(jen.Litf("%s_id", typ.BelongsToStruct.RouteName()), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
		)
		dbCallArgs = append(dbCallArgs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		elseErrBlock = append(elseErrBlock,
			jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Litf("error encountered fetching %s", pcn)),
			utils.WriteXHeader("res", "StatusInternalServerError"),
			jen.Return(),
		)
	} else if typ.BelongsToNobody {
		elseErrBlock = append(elseErrBlock,
			jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Litf("error encountered fetching %s", pcn)),
			utils.WriteXHeader("res", "StatusInternalServerError"),
			jen.Return(),
		)
	}

	dbCallArgs = append(dbCallArgs, jen.ID(utils.FilterVarName))

	block = append(block,
		jen.Line(),
		jen.Commentf("fetch %s from database", pcn),
		jen.List(jen.ID(puvn), jen.Err()).Assign().ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dot(fmt.Sprintf("Get%s", pn)).Call(dbCallArgs...),
		jen.If(jen.Err().Op("==").Qual("database/sql", "ErrNoRows")).Block(
			jen.Comment("in the event no rows exist return an empty list"),
			jen.ID(puvn).Equals().VarPointer().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn)).Valuesln(
				jen.ID(pn).MapAssign().Index().Qual(proj.ModelsV1Package(), sn).Values(),
			),
		).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(elseErrBlock...),
		jen.Line(),
		jen.Comment("encode our response and peace"),
		jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID(puvn)), jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
		),
	)

	lines := []jen.Code{
		jen.Comment("ListHandler is our list route"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(block...),
		),
		jen.Line(),
	}

	return lines
}

func buildCreateHandlerFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	sn := typ.Name.Singular()

	block := []jen.Code{
		jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("CreateHandler")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	notOkayBlock := []jen.Code{}

	if typ.BelongsToUser {
		block = append(block,
			jen.Comment("determine user ID"),
			jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		)
		notOkayBlock = append(notOkayBlock,
			jen.ID("logger").Dot("Info").Call(jen.Lit("valid input not attached to request")),
			utils.WriteXHeader("res", "StatusBadRequest"),
			jen.Return(),
		)
	}
	if typ.BelongsToStruct != nil {
		block = append(block,
			jen.Commentf("determine %s ID", typ.BelongsToStruct.SingularCommonName()),
			jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Assign().ID("s").Dotf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", typ.BelongsToStruct.Singular())).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValue").Call(jen.Litf("%s_id", typ.BelongsToStruct.RouteName()), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
		)
		notOkayBlock = append(notOkayBlock,
			jen.ID("logger").Dot("Info").Call(jen.Lit("valid input not attached to request")),
			utils.WriteXHeader("res", "StatusBadRequest"),
			jen.Return(),
		)
	} else if typ.BelongsToNobody {
		notOkayBlock = append(notOkayBlock,
			jen.ID("s").Dot("logger").Dot("Info").Call(jen.Lit("valid input not attached to request")),
			utils.WriteXHeader("res", "StatusBadRequest"),
			jen.Return(),
		)
	}

	block = append(block,
		jen.Line(),
		jen.Comment("check request context for parsed input struct"),
		jen.List(jen.ID("input"), jen.ID("ok")).Assign().ID(utils.ContextVarName).Dot("Value").Call(jen.ID("CreateMiddlewareCtxKey")).Assert(jen.Op("*").Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", sn))),
		jen.If(jen.Op("!").ID("ok")).Block(notOkayBlock...),
	)

	errNotNilBlock := []jen.Code{}
	if typ.BelongsToUser {
		block = append(block, jen.ID("input").Dot("BelongsToUser").Equals().ID("userID"))
		errNotNilBlock = append(errNotNilBlock,
			jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Litf("error creating %s", scn)),
			utils.WriteXHeader("res", "StatusInternalServerError"),
			jen.Return(),
		)

		block = append(block,
			jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(jen.Lit("input"), jen.ID("input")),
		)
	}
	if typ.BelongsToStruct != nil {
		block = append(block, jen.ID("input").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		errNotNilBlock = append(errNotNilBlock,
			jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Litf("error creating %s", scn)),
			utils.WriteXHeader("res", "StatusInternalServerError"),
			jen.Return(),
		)

		block = append(block,
			jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(jen.Lit("input"), jen.ID("input")),
		)
	} else if typ.BelongsToNobody {
		block = append(block,
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("input"), jen.ID("input")),
		)

		errNotNilBlock = append(errNotNilBlock,
			jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Litf("error creating %s", scn)),
			utils.WriteXHeader("res", "StatusInternalServerError"),
			jen.Return(),
		)
	}

	block = append(block,
		jen.Line(),
		jen.Commentf("create %s in database", scn),
		jen.List(jen.ID("x"), jen.Err()).Assign().ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dot(fmt.Sprintf("Create%s", sn)).Call(utils.CtxVar(), jen.ID("input")),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(errNotNilBlock...),
		jen.Line(),
		jen.Comment("notify relevant parties"),
		jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Dot("Increment").Call(utils.CtxVar()),
		jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.ID("x").Dot("ID")),
		jen.ID("s").Dot("reporter").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
			jen.ID("Data").MapAssign().ID("x"),
			jen.ID("Topics").MapAssign().Index().ID("string").Values(jen.ID("topicName")),
			jen.ID("EventType").MapAssign().ID("string").Call(jen.Qual(proj.ModelsV1Package(), "Create")),
		)),
		jen.Line(),
		jen.Comment("encode our response and peace"),
		utils.WriteXHeader("res", "StatusCreated"),
		jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
		),
	)

	lines := []jen.Code{
		jen.Commentf("CreateHandler is our %s creation route", scn),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(block...),
		),
		jen.Line(),
	}

	return lines
}

func buildExistenceHandlerFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	sn := typ.Name.Singular()
	rn := typ.Name.RouteName()
	xID := fmt.Sprintf("%sID", uvn)

	loggerValues := []jen.Code{}
	dbCallArgs := []jen.Code{
		utils.CtxVar(),
		jen.ID(xID),
	}
	block := []jen.Code{
		jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ExistenceHandler")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
		jen.Comment("determine relevant information"),
	}

	if typ.BelongsToUser {
		loggerValues = append(loggerValues,
			jen.Lit("user_id").MapAssign().ID("userID"),
			jen.Litf("%s_id", rn).MapAssign().ID(fmt.Sprintf("%sID", uvn)),
		)
		dbCallArgs = append(dbCallArgs, jen.ID("userID"))
		block = append(block,
			jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
			jen.ID(xID).Assign().ID("s").Dot(fmt.Sprintf("%sIDFetcher", uvn)).Call(jen.ID("req")),
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(loggerValues...)),
			jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.ID(xID)),
		)
	}
	if typ.BelongsToStruct != nil {
		loggerValues = append(loggerValues,
			jen.Litf("%s_id", typ.BelongsToStruct.RouteName()).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()),
			jen.Litf("%s_id", rn).MapAssign().ID(fmt.Sprintf("%sID", uvn)),
		)
		dbCallArgs = append(dbCallArgs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		block = append(block,
			jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Assign().ID("s").Dotf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Call(jen.ID("req")),
			jen.ID(xID).Assign().ID("s").Dot(fmt.Sprintf("%sIDFetcher", uvn)).Call(jen.ID("req")),
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(loggerValues...)),
			jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.ID(xID)),
		)
	} else if typ.BelongsToNobody {
		block = append(block,
			jen.ID(xID).Assign().ID("s").Dot(fmt.Sprintf("%sIDFetcher", uvn)).Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.ID(xID)),
		)
	}

	elseErrBlock := []jen.Code{}
	if typ.BelongsToUser {
		block = append(block, jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")))
		elseErrBlock = append(elseErrBlock,
			jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit(fmt.Sprintf("error checking %s existence in database", scn))),
			utils.WriteXHeader("res", "StatusInternalServerError"),
			jen.Return(),
		)
	}
	if typ.BelongsToStruct != nil {
		block = append(block, jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", typ.BelongsToStruct.Singular())).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())))
		elseErrBlock = append(elseErrBlock,
			jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit(fmt.Sprintf("error checking %s existence in database", scn))),
			utils.WriteXHeader("res", "StatusInternalServerError"),
			jen.Return(),
		)
	} else if typ.BelongsToNobody {
		elseErrBlock = append(elseErrBlock,
			jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit(fmt.Sprintf("error checking %s existence in database", scn))),
			utils.WriteXHeader("res", "StatusInternalServerError"),
			jen.Return(),
		)
	}

	block = append(block,
		jen.Line(),
		jen.Commentf("fetch %s from database", scn),
		jen.List(jen.ID("exists"), jen.Err()).Assign().ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dotf("%sExists", sn).Call(dbCallArgs...),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(elseErrBlock...),
		jen.Line(),
		jen.If(jen.ID("exists")).Block(
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK")),
		).Else().Block(
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound")),
		),
	)

	lines := []jen.Code{
		jen.Commentf("ExistenceHandler returns a HEAD handler that returns 200 if %s exists, 404 otherwise", scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ExistenceHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(block...),
		),
		jen.Line(),
	}

	return lines
}

func buildReadHandlerFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	sn := typ.Name.Singular()
	rn := typ.Name.RouteName()
	xID := fmt.Sprintf("%sID", uvn)

	loggerValues := []jen.Code{}
	dbCallArgs := []jen.Code{
		utils.CtxVar(),
		jen.ID(xID),
	}
	block := []jen.Code{
		jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ReadHandler")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
		jen.Comment("determine relevant information"),
	}

	if typ.BelongsToUser {
		loggerValues = append(loggerValues,
			jen.Lit("user_id").MapAssign().ID("userID"),
			jen.Litf("%s_id", rn).MapAssign().ID(fmt.Sprintf("%sID", uvn)),
		)
		dbCallArgs = append(dbCallArgs, jen.ID("userID"))
		block = append(block,
			jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
			jen.ID(xID).Assign().ID("s").Dot(fmt.Sprintf("%sIDFetcher", uvn)).Call(jen.ID("req")),
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(loggerValues...)),
			jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.ID(xID)),
		)
	}
	if typ.BelongsToStruct != nil {
		loggerValues = append(loggerValues,
			jen.Litf("%s_id", typ.BelongsToStruct.RouteName()).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()),
			jen.Litf("%s_id", rn).MapAssign().ID(fmt.Sprintf("%sID", uvn)),
		)
		dbCallArgs = append(dbCallArgs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		block = append(block,
			jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Assign().ID("s").Dotf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Call(jen.ID("req")),
			jen.ID(xID).Assign().ID("s").Dot(fmt.Sprintf("%sIDFetcher", uvn)).Call(jen.ID("req")),
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(loggerValues...)),
			jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.ID(xID)),
		)
	} else if typ.BelongsToNobody {
		block = append(block,
			jen.ID(xID).Assign().ID("s").Dot(fmt.Sprintf("%sIDFetcher", uvn)).Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.ID(xID)),
		)
	}

	elseErrBlock := []jen.Code{}
	if typ.BelongsToUser {
		block = append(block, jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")))
		elseErrBlock = append(elseErrBlock,
			jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit(fmt.Sprintf("error fetching %s from database", scn))),
			utils.WriteXHeader("res", "StatusInternalServerError"),
			jen.Return(),
		)
	}
	if typ.BelongsToStruct != nil {
		block = append(block, jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", typ.BelongsToStruct.Singular())).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())))
		elseErrBlock = append(elseErrBlock,
			jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit(fmt.Sprintf("error fetching %s from database", scn))),
			utils.WriteXHeader("res", "StatusInternalServerError"),
			jen.Return(),
		)
	} else if typ.BelongsToNobody {
		elseErrBlock = append(elseErrBlock,
			jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit(fmt.Sprintf("error fetching %s from database", scn))),
			utils.WriteXHeader("res", "StatusInternalServerError"),
			jen.Return(),
		)
	}

	block = append(block,
		jen.Line(),
		jen.Commentf("fetch %s from database", scn),
		jen.List(jen.ID("x"), jen.Err()).Assign().ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dot(fmt.Sprintf("Get%s", sn)).Call(dbCallArgs...),
		jen.If(jen.Err().Op("==").Qual("database/sql", "ErrNoRows")).Block(
			utils.WriteXHeader("res", "StatusNotFound"),
			jen.Return(),
		).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(elseErrBlock...),
		jen.Line(),
		jen.Comment("encode our response and peace"),
		jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
		),
	)

	lines := []jen.Code{
		jen.Commentf("ReadHandler returns a GET handler that returns %s", scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(block...),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateHandlerFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	sn := typ.Name.Singular()
	rn := typ.Name.RouteName()
	xID := fmt.Sprintf("%sID", uvn)

	block := []jen.Code{
		jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("UpdateHandler")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
		jen.Comment("check for parsed input attached to request context"),
		jen.List(jen.ID("input"), jen.ID("ok")).Assign().ID(utils.ContextVarName).Dot("Value").Call(jen.ID("UpdateMiddlewareCtxKey")).Assert(jen.Op("*").Qual(proj.ModelsV1Package(), fmt.Sprintf("%sUpdateInput", sn))),
		jen.If(jen.Op("!").ID("ok")).Block(
			jen.ID("s").Dot("logger").Dot("Info").Call(jen.Lit("no input attached to request")),
			utils.WriteXHeader("res", "StatusBadRequest"),
			jen.Return(),
		),
		jen.Line(),
		jen.Comment("determine relevant information"),
	}
	loggerValues := []jen.Code{}
	dbCallArgs := []jen.Code{
		utils.CtxVar(),
		jen.ID(xID),
	}

	if typ.BelongsToUser {
		block = append(block, jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID("req")))
		loggerValues = append(loggerValues, jen.Lit("user_id").MapAssign().ID("userID"))
		dbCallArgs = append(dbCallArgs, jen.ID("userID"))
	}
	if typ.BelongsToStruct != nil {
		block = append(block, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Assign().ID("s").Dotf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Call(jen.ID("req")))
		loggerValues = append(loggerValues, jen.Litf("%s_id", typ.BelongsToStruct.RouteName()).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		dbCallArgs = append(dbCallArgs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}
	loggerValues = append(loggerValues, jen.Litf("%s_id", rn).MapAssign().ID(fmt.Sprintf("%sID", uvn)))

	block = append(block,
		jen.ID(xID).Assign().ID("s").Dot(fmt.Sprintf("%sIDFetcher", uvn)).Call(jen.ID("req")),
		jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(loggerValues...)),
		jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.ID(xID)),
	)

	if typ.BelongsToUser {
		block = append(block, jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")))
	}
	if typ.BelongsToStruct != nil {
		block = append(block, jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", typ.BelongsToStruct.Singular())).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())))
	}

	block = append(block,
		jen.Line(),
		jen.Commentf("fetch %s from database", scn),
		jen.List(jen.ID("x"), jen.Err()).Assign().ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dotf("Get%s", sn).Call(dbCallArgs...),
		jen.If(jen.Err().Op("==").Qual("database/sql", "ErrNoRows")).Block(
			utils.WriteXHeader("res", "StatusNotFound"),
			jen.Return(),
		).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Litf("error encountered getting %s", scn)),
			utils.WriteXHeader("res", "StatusInternalServerError"),
			jen.Return(),
		),
		jen.Line(),
		jen.Comment("update the data structure"),
		jen.ID("x").Dot("Update").Call(jen.ID("input")),
		jen.Line(),
		jen.Commentf("update %s in database", scn),
		jen.If(jen.Err().Equals().ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dotf("Update%s", sn).Call(utils.CtxVar(), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Litf("error encountered updating %s", scn)),
			utils.WriteXHeader("res", "StatusInternalServerError"),
			jen.Return(),
		),
		jen.Line(),
		jen.Comment("notify relevant parties"),
		jen.ID("s").Dot("reporter").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
			jen.ID("Data").MapAssign().ID("x"),
			jen.ID("Topics").MapAssign().Index().ID("string").Values(jen.ID("topicName")),
			jen.ID("EventType").MapAssign().ID("string").Call(jen.Qual(proj.ModelsV1Package(), "Update")),
		)),
		jen.Line(),
		jen.Comment("encode our response and peace"),
		jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
		),
	)

	lines := []jen.Code{
		jen.Commentf("UpdateHandler returns a handler that updates %s", scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("UpdateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(block...),
		),
		jen.Line(),
	}

	return lines
}

func buildArchiveHandlerFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	sn := typ.Name.Singular()
	rn := typ.Name.RouteName()
	xID := fmt.Sprintf("%sID", uvn)

	blockLines := []jen.Code{
		jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ArchiveHandler")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
		jen.Comment("determine relevant information"),
	}

	callArgs := []jen.Code{
		utils.CtxVar(), jen.ID(xID),
	}
	loggerValues := []jen.Code{jen.Litf("%s_id", rn).MapAssign().ID(fmt.Sprintf("%sID", uvn))}

	if typ.BelongsToUser {
		loggerValues = append(loggerValues, jen.Lit("user_id").MapAssign().ID("userID"))
		blockLines = append(blockLines, jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID("req")))
		callArgs = append(callArgs, jen.ID("userID"))
	}
	if typ.BelongsToStruct != nil {
		loggerValues = append(loggerValues, jen.Litf("%s_id", typ.BelongsToStruct.RouteName()).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		blockLines = append(blockLines, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Assign().ID("s").Dotf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Call(jen.ID("req")))
		callArgs = append(callArgs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}

	blockLines = append(blockLines,
		jen.ID(xID).Assign().ID("s").Dot(fmt.Sprintf("%sIDFetcher", uvn)).Call(jen.ID("req")),
		jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(loggerValues...)),
		jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.ID(xID)),
	)

	if typ.BelongsToUser {
		blockLines = append(blockLines, jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")))
	}
	if typ.BelongsToStruct != nil {
		blockLines = append(blockLines, jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", typ.BelongsToStruct.Singular())).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())))
	}

	blockLines = append(blockLines,
		jen.Line(),
		jen.Commentf("archive the %s in the database", scn),
		jen.Err().Assign().ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dotf("Archive%s", sn).Call(callArgs...),
		jen.If(jen.Err().Op("==").Qual("database/sql", "ErrNoRows")).Block(
			utils.WriteXHeader("res", "StatusNotFound"),
			jen.Return(),
		).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Litf("error encountered deleting %s", scn)),
			utils.WriteXHeader("res", "StatusInternalServerError"),
			jen.Return(),
		),
		jen.Line(),
		jen.Comment("notify relevant parties"),
		jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Dot("Decrement").Call(utils.CtxVar()),
		jen.ID("s").Dot("reporter").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
			jen.ID("EventType").MapAssign().ID("string").Call(jen.Qual(proj.ModelsV1Package(), "Archive")),
			jen.ID("Data").MapAssign().VarPointer().Qual(proj.ModelsV1Package(), sn).Values(jen.ID("ID").MapAssign().ID(fmt.Sprintf("%sID", uvn))),
			jen.ID("Topics").MapAssign().Index().ID("string").Values(jen.ID("topicName")),
		)),
		jen.Line(),
		jen.Comment("encode our response and peace"),
		utils.WriteXHeader("res", "StatusNoContent"),
	)

	lines := []jen.Code{
		jen.Commentf("ArchiveHandler returns a handler that archives %s", scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(blockLines...),
		),
		jen.Line(),
	}

	return lines
}
