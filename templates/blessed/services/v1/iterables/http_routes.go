package iterables

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

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
			jen.Commentf("URIParamKey is a standard string that we'll use to refer to %s IDs with.", scn),
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

func buildRequisiteLoggerAndTracingStatementsForListOfEntities(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	if typ.OwnedByAUserAtSomeLevel(proj) {
		lines = append(lines,
			jen.Comment("determine user ID."),
			jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID(constants.RequestVarName)),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("userID")),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
			jen.Line(),
		)
	}

	owners := proj.FindOwnerTypeChain(typ)
	for _, o := range owners {
		lines = append(lines,
			jen.Commentf("determine %s ID.", o.Name.SingularCommonName()),
			jen.IDf("%sID", o.Name.UnexportedVarName()).Assign().ID("s").Dotf("%sIDFetcher", o.Name.UnexportedVarName()).Call(jen.ID(constants.RequestVarName)),
			jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", o.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", o.Name.UnexportedVarName())),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Litf("%s_id", o.Name.RouteName()), jen.IDf("%sID", o.Name.UnexportedVarName())),
			jen.Line(),
		)
	}

	return lines
}

func buildRequisiteLoggerAndTracingStatementsForSingleEntity(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	if typ.OwnedByAUserAtSomeLevel(proj) {
		lines = append(lines,
			jen.Comment("determine user ID."),
			jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID(constants.RequestVarName)),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("userID")),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
			jen.Line(),
		)
	}

	owners := proj.FindOwnerTypeChain(typ)
	for _, o := range owners {
		lines = append(lines,
			jen.Commentf("determine %s ID.", o.Name.SingularCommonName()),
			jen.IDf("%sID", o.Name.UnexportedVarName()).Assign().ID("s").Dotf("%sIDFetcher", o.Name.UnexportedVarName()).Call(jen.ID(constants.RequestVarName)),
			jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", o.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", o.Name.UnexportedVarName())),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Litf("%s_id", o.Name.RouteName()), jen.IDf("%sID", o.Name.UnexportedVarName())),
			jen.Line(),
		)
	}

	lines = append(lines,
		jen.Commentf("determine %s ID.", typ.Name.SingularCommonName()),
		jen.IDf("%sID", typ.Name.UnexportedVarName()).Assign().ID("s").Dotf("%sIDFetcher", typ.Name.UnexportedVarName()).Call(jen.ID(constants.RequestVarName)),
		jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", typ.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", typ.Name.UnexportedVarName())),
		jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Litf("%s_id", typ.Name.RouteName()), jen.IDf("%sID", typ.Name.UnexportedVarName())),
		jen.Line(),
	)

	return lines
}

func buildListHandlerFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	puvn := typ.Name.PluralUnexportedVarName()
	pcn := typ.Name.PluralCommonName()
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	block := []jen.Code{
		jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("ListHandler")),
		jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
		jen.Line(),
		jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
		jen.Line(),
		jen.Comment("ensure query filter."),
		jen.ID(constants.FilterVarName).Assign().Qual(proj.ModelsV1Package(), "ExtractQueryFilter").Call(jen.ID(constants.RequestVarName)),
		jen.Line(),
	}
	block = append(block, buildRequisiteLoggerAndTracingStatementsForListOfEntities(proj, typ)...)

	dbCallArgs := typ.BuildDBClientListRetrievalMethodCallArgs(proj)

	block = append(block,
		jen.Line(),
		jen.Commentf("fetch %s from database.", pcn),
		jen.List(jen.ID(puvn), jen.Err()).Assign().ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Dot(fmt.Sprintf("Get%s", pn)).Call(dbCallArgs...),
		jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
			jen.Comment("in the event no rows exist return an empty list."),
			jen.ID(puvn).Equals().AddressOf().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn)).Valuesln(
				jen.ID(pn).MapAssign().Index().Qual(proj.ModelsV1Package(), sn).Values(),
			),
		).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Litf("error encountered fetching %s", pcn)),
			utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
			jen.Return(),
		),
		jen.Line(),
		jen.Comment("encode our response and peace."),
		jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID(constants.ResponseVarName), jen.ID(puvn)), jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
		),
	)

	lines := []jen.Code{
		jen.Comment("ListHandler is our list route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(block...),
		),
		jen.Line(),
	}

	return lines
}

func buildRequisiteLoggerAndTracingStatementsForModification(proj *models.Project, typ models.DataType, includeExistenceChecks, includeSelf, assignToUser bool) []jen.Code {
	lines := []jen.Code{}

	if typ.OwnedByAUserAtSomeLevel(proj) {
		lines = append(lines,
			jen.Comment("determine user ID."),
			jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID(constants.RequestVarName)),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("userID")),
			func() jen.Code {
				if assignToUser && typ.BelongsToUser {
					return jen.ID("input").Dot(constants.UserOwnershipFieldName).Equals().ID("userID")
				}
				return jen.Null()
			}(),
			jen.Line(),
		)
	}

	owners := proj.FindOwnerTypeChain(typ)
	for _, o := range owners {
		existenceCallArgs := o.BuildArgsForServiceRouteExistenceCheck(proj)
		lines = append(lines,
			jen.Commentf("determine %s ID.", o.Name.SingularCommonName()),
			jen.IDf("%sID", o.Name.UnexportedVarName()).Assign().ID("s").Dotf("%sIDFetcher", o.Name.UnexportedVarName()).Call(jen.ID(constants.RequestVarName)),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Litf("%s_id", o.Name.RouteName()), jen.IDf("%sID", o.Name.UnexportedVarName())),
			jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", o.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", o.Name.UnexportedVarName())),
			jen.Line(),
		)

		if includeExistenceChecks {
			lines = append(lines,
				jen.List(jen.IDf("%sExists", o.Name.UnexportedVarName()), jen.Err()).Assign().ID("s").Dotf("%sDataManager", o.Name.UnexportedVarName()).Dotf("%sExists", o.Name.Singular()).Call(existenceCallArgs...),
				jen.If(jen.Err().DoesNotEqual().Nil().And().Err().DoesNotEqual().Qual("database/sql", "ErrNoRows")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Litf("error checking %s existence", o.Name.SingularCommonName())),
					jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				).Else().If().Not().IDf("%sExists", o.Name.UnexportedVarName()).Block(
					jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound")),
					jen.Return(),
				),
			)
		}

		lines = append(lines, jen.Line())
	}

	if includeSelf {
		lines = append(lines,
			jen.Commentf("determine %s ID.", typ.Name.SingularCommonName()),
			jen.IDf("%sID", typ.Name.UnexportedVarName()).Assign().ID("s").Dotf("%sIDFetcher", typ.Name.UnexportedVarName()).Call(jen.ID(constants.RequestVarName)),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Litf("%s_id", typ.Name.RouteName()), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", typ.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Line(),
		)
	}

	return lines
}

func buildCreateHandlerFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	sn := typ.Name.Singular()

	block := []jen.Code{
		jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("CreateHandler")),
		jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
		jen.Line(),
		jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
		jen.Line(),
		jen.Comment("check request context for parsed input struct."),
		jen.List(jen.ID("input"), jen.ID("ok")).Assign().ID(constants.ContextVarName).Dot("Value").Call(jen.ID("CreateMiddlewareCtxKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", sn))),
		jen.If(jen.Op("!").ID("ok")).Block(
			jen.ID(constants.LoggerVarName).Dot("Info").Call(jen.Lit("valid input not attached to request")),
			utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
			jen.Return(),
		),
		jen.Line(),
	}

	block = append(block, buildRequisiteLoggerAndTracingStatementsForModification(proj, typ, true, false, true)...)
	if typ.BelongsToStruct != nil {
		block = append(block, jen.Line(), jen.ID("input").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}

	errNotNilBlock := []jen.Code{
		jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Litf("error creating %s", scn)),
		utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
		jen.Return(),
	}

	block = append(block,
		jen.Line(),
		jen.Commentf("create %s in database.", scn),
		jen.List(jen.ID("x"), jen.Err()).Assign().ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Dot(fmt.Sprintf("Create%s", sn)).Call(constants.CtxVar(), jen.ID("input")),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(errNotNilBlock...),
		jen.Line(),
		jen.Qual(proj.InternalTracingV1Package(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID(constants.SpanVarName), jen.ID("x").Dot("ID")),
		jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Litf("%s_id", typ.Name.RouteName()), jen.ID("x").Dot("ID")),
		jen.Line(),
		jen.Comment("notify relevant parties."),
		jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Dot("Increment").Call(constants.CtxVar()),
		jen.ID("s").Dot("reporter").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
			jen.ID("Data").MapAssign().ID("x"),
			jen.ID("Topics").MapAssign().Index().String().Values(jen.ID("topicName")),
			jen.ID("EventType").MapAssign().String().Call(jen.Qual(proj.ModelsV1Package(), "Create")),
		)),
		jen.Line(),
		jen.Comment("encode our response and peace."),
		utils.WriteXHeader(constants.ResponseVarName, "StatusCreated"),
		jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID(constants.ResponseVarName), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
		),
	)

	lines := []jen.Code{
		jen.Commentf("CreateHandler is our %s creation route.", scn),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(block...),
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

	block := []jen.Code{
		jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("ExistenceHandler")),
		jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
		jen.Line(),
		jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
		jen.Line(),
	}
	block = append(block, buildRequisiteLoggerAndTracingStatementsForSingleEntity(proj, typ)...)

	elseErrBlock := []jen.Code{
		jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit(fmt.Sprintf("error checking %s existence in database", scn))),
		utils.WriteXHeader(constants.ResponseVarName, "StatusNotFound"),
		jen.Return(),
	}
	dbCallArgs := typ.BuildDBClientExistenceMethodCallArgs(proj)

	block = append(block,
		jen.Line(),
		jen.Commentf("fetch %s from database.", scn),
		jen.List(jen.ID("exists"), jen.Err()).Assign().ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Dotf("%sExists", sn).Call(dbCallArgs...),
		jen.If(jen.Err().DoesNotEqual().ID("nil").And().Err().DoesNotEqual().Qual("database/sql", "ErrNoRows")).Block(elseErrBlock...),
		jen.Line(),
		jen.If(jen.ID("exists")).Block(
			jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK")),
		).Else().Block(
			jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound")),
		),
	)

	lines := []jen.Code{
		jen.Commentf("ExistenceHandler returns a HEAD handler that returns 200 if %s exists, 404 otherwise.", scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ExistenceHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(block...),
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

	block := []jen.Code{
		jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("ReadHandler")),
		jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
		jen.Line(),
		jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
		jen.Line(),
	}
	block = append(block, buildRequisiteLoggerAndTracingStatementsForSingleEntity(proj, typ)...)

	elseErrBlock := []jen.Code{
		jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit(fmt.Sprintf("error fetching %s from database", scn))),
		utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
		jen.Return(),
	}

	dbCallArgs := typ.BuildDBClientRetrievalMethodCallArgs(proj)
	block = append(block,
		jen.Line(),
		jen.Commentf("fetch %s from database.", scn),
		jen.List(jen.ID("x"), jen.Err()).Assign().ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Dotf("Get%s", sn).Call(dbCallArgs...),
		jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
			utils.WriteXHeader(constants.ResponseVarName, "StatusNotFound"),
			jen.Return(),
		).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(elseErrBlock...),
		jen.Line(),
		jen.Comment("encode our response and peace."),
		jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID(constants.ResponseVarName), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
		),
	)

	lines := []jen.Code{
		jen.Commentf("ReadHandler returns a GET handler that returns %s.", scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(block...),
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

	block := []jen.Code{
		jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("UpdateHandler")),
		jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
		jen.Line(),
		jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
		jen.Line(),
	}

	block = append(block,
		jen.Line(),
		jen.Comment("check for parsed input attached to request context."),
		jen.List(jen.ID("input"), jen.ID("ok")).Assign().ID(constants.ContextVarName).Dot("Value").Call(jen.ID("UpdateMiddlewareCtxKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sUpdateInput", sn))),
		jen.If(jen.Op("!").ID("ok")).Block(
			jen.ID(constants.LoggerVarName).Dot("Info").Call(jen.Lit("no input attached to request")),
			utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
			jen.Return(),
		),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID("input").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("%sID", typ.BelongsToStruct.UnexportedVarName())
			} else {
				return jen.Null()
			}
		}(),
		jen.Line(),
	)

	block = append(block, buildRequisiteLoggerAndTracingStatementsForModification(proj, typ, false, true, typ.BelongsToUser)...)
	fetchCallArgs := typ.BuildDBClientRetrievalMethodCallArgs(proj)

	block = append(block,
		jen.Line(),
		jen.Commentf("fetch %s from database.", scn),
		jen.List(jen.ID("x"), jen.Err()).Assign().ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Dotf("Get%s", sn).Call(fetchCallArgs...),
		jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
			utils.WriteXHeader(constants.ResponseVarName, "StatusNotFound"),
			jen.Return(),
		).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Litf("error encountered getting %s", scn)),
			utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
			jen.Return(),
		),
		jen.Line(),
		jen.Comment("update the data structure."),
		jen.ID("x").Dot("Update").Call(jen.ID("input")),
		jen.Line(),
		jen.Commentf("update %s in database.", scn),
		jen.If(jen.Err().Equals().ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Dotf("Update%s", sn).Call(constants.CtxVar(), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Litf("error encountered updating %s", scn)),
			utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
			jen.Return(),
		),
		jen.Line(),
		jen.Comment("notify relevant parties."),
		jen.ID("s").Dot("reporter").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
			jen.ID("Data").MapAssign().ID("x"),
			jen.ID("Topics").MapAssign().Index().String().Values(jen.ID("topicName")),
			jen.ID("EventType").MapAssign().String().Call(jen.Qual(proj.ModelsV1Package(), "Update")),
		)),
		jen.Line(),
		jen.Comment("encode our response and peace."),
		jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID(constants.ResponseVarName), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
		),
	)

	lines := []jen.Code{
		jen.Commentf("UpdateHandler returns a handler that updates %s.", scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("UpdateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(block...),
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

	blockLines := []jen.Code{
		jen.Var().Err().Error(),
		jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("ArchiveHandler")),
		jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
		jen.Line(),
		jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
		jen.Line(),
	}
	blockLines = append(blockLines, buildRequisiteLoggerAndTracingStatementsForModification(proj, typ, true, true, false)...)
	callArgs := typ.BuildDBClientArchiveMethodCallArgs(proj)

	blockLines = append(blockLines,
		jen.Line(),
		jen.Commentf("archive the %s in the database.", scn),
		jen.Err().Equals().ID("s").Dot(fmt.Sprintf("%sDataManager", uvn)).Dotf("Archive%s", sn).Call(callArgs...),
		jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
			utils.WriteXHeader(constants.ResponseVarName, "StatusNotFound"),
			jen.Return(),
		).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Litf("error encountered deleting %s", scn)),
			utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
			jen.Return(),
		),
		jen.Line(),
		jen.Comment("notify relevant parties."),
		jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Dot("Decrement").Call(constants.CtxVar()),
		jen.ID("s").Dot("reporter").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
			jen.ID("EventType").MapAssign().String().Call(jen.Qual(proj.ModelsV1Package(), "Archive")),
			jen.ID("Data").MapAssign().AddressOf().Qual(proj.ModelsV1Package(), sn).Values(jen.ID("ID").MapAssign().ID(fmt.Sprintf("%sID", uvn))),
			jen.ID("Topics").MapAssign().Index().String().Values(jen.ID("topicName")),
		)),
		jen.Line(),
		jen.Comment("encode our response and peace."),
		utils.WriteXHeader(constants.ResponseVarName, "StatusNoContent"),
	)

	lines := []jen.Code{
		jen.Commentf("ArchiveHandler returns a handler that archives %s.", scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(blockLines...),
		),
		jen.Line(),
	}

	return lines
}
