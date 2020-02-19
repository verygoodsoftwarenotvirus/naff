package iterables

import (
	"fmt"
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	sn := typ.Name.Singular()
	rn := typ.Name.RouteName()

	xID := fmt.Sprintf("%sID", uvn)

	utils.AddImports(pkg.OutputPath, []models.DataType{typ}, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment(fmt.Sprintf("URIParamKey is a standard string that we'll use to refer to %s IDs with", scn)),
			jen.ID("URIParamKey").Op("=").Lit(fmt.Sprintf("%sID", uvn)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID(fmt.Sprintf("attach%sIDToSpan", sn)).Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID(xID).ID("uint64")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot("AddAttributes").Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit(fmt.Sprintf("%s_id", rn)), jen.Qual("strconv", "FormatUint").Call(jen.ID(xID), jen.Lit(10)))),
			),
		),
		jen.Line(),
	)

	if typ.BelongsToUser {
		ret.Add(
			jen.Func().ID("attachUserIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("userID").ID("uint64")).Block(
				jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
					jen.ID("span").Dot("AddAttributes").Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("user_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("userID"), jen.Lit(10)))),
				),
			),
			jen.Line(),
		)
	} else if typ.BelongsToStruct != nil {
		ret.Add(
			jen.Func().IDf("attach%sIDToSpan", typ.BelongsToStruct.Singular()).Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.IDf("%sID", typ.BelongsToStruct.Singular()).ID("uint64")).Block(
				jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
					jen.ID("span").Dot("AddAttributes").Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Litf("%s_id", typ.BelongsToStruct.RouteName()), jen.Qual("strconv", "FormatUint").Call(jen.IDf("%sID", typ.BelongsToStruct.Singular()), jen.Lit(10)))),
				),
			),
			jen.Line(),
		)
	}

	ret.Add(buildListHandlerFuncDecl(pkg, typ)...)
	ret.Add(buildCreateHandlerFuncDecl(pkg, typ)...)
	ret.Add(buildReadHandlerFuncDecl(pkg, typ)...)
	ret.Add(buildUpdateHandlerFuncDecl(pkg, typ)...)
	ret.Add(buildArchiveHandlerFuncDecl(pkg, typ)...)

	return ret
}

func buildListHandlerFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	puvn := typ.Name.PluralUnexportedVarName()
	pcn := typ.Name.PluralCommonName()
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	dbCallArgs := []jen.Code{
		jen.ID("ctx"),
		jen.ID("qf"),
	}
	block := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ListHandler")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
		jen.Comment("ensure query filter"),
		jen.ID("qf").Op(":=").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "ExtractQueryFilter").Call(jen.ID("req")),
		jen.Line(),
	}

	if typ.BelongsToUser {
		block = append(block,
			jen.Comment("determine user ID"),
			jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		)
		dbCallArgs = append(dbCallArgs, jen.ID("userID"))
	} else if typ.BelongsToStruct != nil {
		block = append(block,
			jen.Commentf("determine %s ID", typ.BelongsToStruct.SingularCommonName()),
			jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Op(":=").ID("s").Dotf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Call(jen.ID("req")),
			jen.IDf("attach%sIDToSpan", typ.BelongsToStruct.Singular()).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(jen.Litf("%s_id", typ.BelongsToStruct.RouteName()), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
		)
		dbCallArgs = append(dbCallArgs, jen.ID("userID"))
	}

	block = append(block,
		jen.Line(),
		jen.Comment(fmt.Sprintf("fetch %s from database", pcn)),
		jen.List(jen.ID(puvn), jen.ID("err")).Op(":=").ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dot(fmt.Sprintf("Get%s", pn)).Call(dbCallArgs...),
		jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
			jen.Comment("in the event no rows exist return an empty list"),
			jen.ID(puvn).Op("=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", sn)).Valuesln(
				jen.ID(pn).Op(":").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Values(),
			),
		).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Litf("error encountered fetching %s", pcn)),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
			jen.Return(),
		),
		jen.Line(),
		jen.Comment("encode our response and peace"),
		jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID(puvn)), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
		),
	)

	lines := []jen.Code{
		jen.Comment("ListHandler is our list route"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(block...),
		),
		jen.Line(),
	}

	return lines
}

func buildCreateHandlerFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	sn := typ.Name.Singular()

	block := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("CreateHandler")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	if typ.BelongsToUser {
		block = append(block,
			jen.Comment("determine user ID"),
			jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		)
	} else if typ.BelongsToStruct != nil {
		block = append(block,
			jen.Commentf("determine %s ID", typ.BelongsToStruct.SingularCommonName()),
			jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Op(":=").ID("s").Dotf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Call(jen.ID("req")),
			jen.IDf("attach%sIDToSpan", typ.BelongsToStruct.Singular()).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(jen.Litf("%s_id", typ.BelongsToStruct.RouteName()), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
		)
	}

	block = append(block,
		jen.Line(),
		jen.Comment("check request context for parsed input struct"),
		jen.List(jen.ID("input"), jen.ID("ok")).Op(":=").ID("ctx").Dot("Value").Call(jen.ID("CreateMiddlewareCtxKey")).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", sn))),
		jen.If(jen.Op("!").ID("ok")).Block(
			jen.ID("logger").Dot("Info").Call(jen.Lit("valid input not attached to request")),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
			jen.Return(),
		),
		jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(jen.Lit("input"), jen.ID("input")),
	)

	if typ.BelongsToUser {
		block = append(block, jen.ID("input").Dot("BelongsToUser").Op("=").ID("userID"))
	} else if typ.BelongsToStruct != nil {
		block = append(block, jen.ID("input").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Op("=").IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}

	block = append(block,
		jen.Line(),
		jen.Commentf("create %s in database", scn),
		jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dot(fmt.Sprintf("Create%s", sn)).Call(jen.ID("ctx"), jen.ID("input")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Litf("error creating %s", scn)),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
			jen.Return(),
		),
		jen.Line(),
		jen.Comment("notify relevant parties"),
		jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Dot("Increment").Call(jen.ID("ctx")),
		jen.ID(fmt.Sprintf("attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.ID("x").Dot("ID")),
		jen.ID("s").Dot("reporter").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
			jen.ID("Data").Op(":").ID("x"),
			jen.ID("Topics").Op(":").Index().ID("string").Values(jen.ID("topicName")),
			jen.ID("EventType").Op(":").ID("string").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Create")),
		)),
		jen.Line(),
		jen.Comment("encode our response and peace"),
		jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusCreated")),
		jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
		),
	)

	lines := []jen.Code{
		jen.Commentf("CreateHandler is our %s creation route", scn),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(block...),
		),
		jen.Line(),
	}

	return lines
}

func buildReadHandlerFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	sn := typ.Name.Singular()
	rn := typ.Name.RouteName()
	xID := fmt.Sprintf("%sID", uvn)

	loggerValues := []jen.Code{}
	dbCallArgs := []jen.Code{
		jen.ID("ctx"),
		jen.ID(xID),
	}
	block := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ReadHandler")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
		jen.Comment("determine relevant information"),
	}

	if typ.BelongsToUser {
		loggerValues = append(loggerValues, jen.Lit("user_id").Op(":").ID("userID"))
		dbCallArgs = append(dbCallArgs, jen.ID("userID"))
		block = append(block, jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")))
	} else if typ.BelongsToStruct != nil {
		loggerValues = append(loggerValues, jen.Litf("%s_id", typ.BelongsToStruct.RouteName()).Op(":").IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		dbCallArgs = append(dbCallArgs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		block = append(block, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Op(":=").ID("s").Dotf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Call(jen.ID("req")))
	}
	loggerValues = append(loggerValues, jen.Litf("%s_id", rn).Op(":").ID(fmt.Sprintf("%sID", uvn)))

	block = append(block,
		jen.ID(xID).Op(":=").ID("s").Dot(fmt.Sprintf("%sIDFetcher", uvn)).Call(jen.ID("req")),
		jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(loggerValues...)),
		jen.ID(fmt.Sprintf("attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.ID(xID)),
	)

	if typ.BelongsToUser {
		block = append(block, jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")))
	} else if typ.BelongsToStruct != nil {
		block = append(block, jen.IDf("attach%sIDToSpan", typ.BelongsToStruct.Singular()).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())))
	}

	block = append(block,
		jen.Line(),
		jen.Commentf("fetch %s from database", scn),
		jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dot(fmt.Sprintf("Get%s", sn)).Call(dbCallArgs...),
		jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound")),
			jen.Return(),
		).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit(fmt.Sprintf("error fetching %s from database", scn))),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
			jen.Return(),
		),
		jen.Line(),
		jen.Comment("encode our response and peace"),
		jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
		),
	)

	lines := []jen.Code{
		jen.Commentf("ReadHandler returns a GET handler that returns %s", scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(block...),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateHandlerFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	sn := typ.Name.Singular()
	rn := typ.Name.RouteName()
	xID := fmt.Sprintf("%sID", uvn)

	block := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("UpdateHandler")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
		jen.Comment("check for parsed input attached to request context"),
		jen.List(jen.ID("input"), jen.ID("ok")).Op(":=").ID("ctx").Dot("Value").Call(jen.ID("UpdateMiddlewareCtxKey")).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sUpdateInput", sn))),
		jen.If(jen.Op("!").ID("ok")).Block(
			jen.ID("s").Dot("logger").Dot("Info").Call(jen.Lit("no input attached to request")),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
			jen.Return(),
		),
		jen.Line(),
		jen.Comment("determine relevant information"),
	}
	loggerValues := []jen.Code{}
	dbCallArgs := []jen.Code{
		jen.ID("ctx"),
		jen.ID(xID),
	}

	if typ.BelongsToUser {
		block = append(block, jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")))
		loggerValues = append(loggerValues, jen.Lit("user_id").Op(":").ID("userID"))
		dbCallArgs = append(dbCallArgs, jen.ID("userID"))
	} else if typ.BelongsToStruct != nil {
		block = append(block, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Op(":=").ID("s").Dotf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Call(jen.ID("req")))
		loggerValues = append(loggerValues, jen.Litf("%s_id", typ.BelongsToStruct.RouteName()).Op(":").IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		dbCallArgs = append(dbCallArgs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}
	loggerValues = append(loggerValues, jen.Litf("%s_id", rn).Op(":").ID(fmt.Sprintf("%sID", uvn)))

	block = append(block,
		jen.ID(xID).Op(":=").ID("s").Dot(fmt.Sprintf("%sIDFetcher", uvn)).Call(jen.ID("req")),
		jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(loggerValues...)),
		jen.ID(fmt.Sprintf("attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.ID(xID)),
	)

	if typ.BelongsToUser {
		block = append(block, jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")))
	} else if typ.BelongsToStruct != nil {
		block = append(block, jen.IDf("attach%sIDToSpan", typ.BelongsToStruct.Singular()).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())))
	}

	block = append(block,
		jen.Line(),
		jen.Commentf("fetch %s from database", scn),
		jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dotf("Get%s", sn).Call(dbCallArgs...),
		jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound")),
			jen.Return(),
		).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Litf("error encountered getting %s", scn)),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
			jen.Return(),
		),
		jen.Line(),
		jen.Comment("update the data structure"),
		jen.ID("x").Dot("Update").Call(jen.ID("input")),
		jen.Line(),
		jen.Commentf("update %s in database", scn),
		jen.If(jen.ID("err").Op("=").ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dotf("Update%s", sn).Call(jen.ID("ctx"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Litf("error encountered updating %s", scn)),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
			jen.Return(),
		),
		jen.Line(),
		jen.Comment("notify relevant parties"),
		jen.ID("s").Dot("reporter").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
			jen.ID("Data").Op(":").ID("x"),
			jen.ID("Topics").Op(":").Index().ID("string").Values(jen.ID("topicName")),
			jen.ID("EventType").Op(":").ID("string").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Update")),
		)),
		jen.Line(),
		jen.Comment("encode our response and peace"),
		jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
		),
	)

	lines := []jen.Code{
		jen.Commentf("UpdateHandler returns a handler that updates %s", scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("UpdateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(block...),
		),
		jen.Line(),
	}

	return lines
}

func buildArchiveHandlerFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	sn := typ.Name.Singular()
	rn := typ.Name.RouteName()
	xID := fmt.Sprintf("%sID", uvn)

	blockLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ArchiveHandler")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
		jen.Comment("determine relevant information"),
	}

	callArgs := []jen.Code{
		jen.ID("ctx"), jen.ID(xID),
	}
	loggerValues := []jen.Code{jen.Litf("%s_id", rn).Op(":").ID(fmt.Sprintf("%sID", uvn))}

	if typ.BelongsToUser {
		loggerValues = append(loggerValues, jen.Lit("user_id").Op(":").ID("userID"))
		blockLines = append(blockLines, jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")))
		callArgs = append(callArgs, jen.ID("userID"))
	} else if typ.BelongsToStruct != nil {
		loggerValues = append(loggerValues, jen.Litf("%s_id", typ.BelongsToStruct.RouteName()).Op(":").IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		blockLines = append(blockLines, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Op(":=").ID("s").Dotf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Call(jen.ID("req")))
		callArgs = append(callArgs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}

	blockLines = append(blockLines,
		jen.ID(xID).Op(":=").ID("s").Dot(fmt.Sprintf("%sIDFetcher", uvn)).Call(jen.ID("req")),
		jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(loggerValues...)),
		jen.ID(fmt.Sprintf("attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.ID(xID)),
	)

	if typ.BelongsToUser {
		blockLines = append(blockLines, jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")))
	} else if typ.BelongsToStruct != nil {
		blockLines = append(blockLines, jen.IDf("attach%sIDToSpan", typ.BelongsToStruct.Singular()).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())))
	}

	blockLines = append(blockLines,
		jen.Line(),
		jen.Commentf("archive the %s in the database", scn),
		jen.ID("err").Op(":=").ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dotf("Archive%s", sn).Call(callArgs...),
		jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound")),
			jen.Return(),
		).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Litf("error encountered deleting %s", scn)),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
			jen.Return(),
		),
		jen.Line(),
		jen.Comment("notify relevant parties"),
		jen.ID("s").Dot(fmt.Sprintf("%sCounter", uvn)).Dot("Decrement").Call(jen.ID("ctx")),
		jen.ID("s").Dot("reporter").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
			jen.ID("EventType").Op(":").ID("string").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Archive")),
			jen.ID("Data").Op(":").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Values(jen.ID("ID").Op(":").ID(fmt.Sprintf("%sID", uvn))),
			jen.ID("Topics").Op(":").Index().ID("string").Values(jen.ID("topicName")),
		)),
		jen.Line(),
		jen.Comment("encode our response and peace"),
		jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNoContent")),
	)

	lines := []jen.Code{
		jen.Commentf("ArchiveHandler returns a handler that archives %s", scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(blockLines...),
		),
		jen.Line(),
	}

	return lines
}
