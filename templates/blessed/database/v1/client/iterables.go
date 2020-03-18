package client

import (
	"fmt"
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(pkg, ret)

	n := typ.Name
	sn := n.Singular()

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sDataManager", sn)).Op("=").Parens(jen.Op("*").ID("Client")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(buildSomethingExists(pkg, typ)...)
	ret.Add(buildGetSomething(pkg, typ)...)
	ret.Add(buildGetSomethingCount(pkg, typ)...)
	ret.Add(buildGetAllSomethingCount(pkg, typ)...)
	ret.Add(buildGetListOfSomething(pkg, typ)...)

	if typ.BelongsToUser {
		ret.Add(buildGetAllSomethingForUser(pkg, typ)...)
	}
	if typ.BelongsToStruct != nil {
		ret.Add(buildGetAllSomethingForSomethingElse(pkg, typ)...)
	}

	ret.Add(buildCreateSomething(pkg, typ)...)
	ret.Add(buildUpdateSomething(pkg, typ)...)
	ret.Add(buildArchiveSomething(pkg, typ)...)

	return ret
}

func buildSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	funcName := fmt.Sprintf("%sExists", sn)
	params := typ.BuildGetSomethingParams(proj)
	args := typ.BuildGetSomethingArgs(proj)

	block := utils.StartSpan(true, funcName)
	block = append(block,
		jen.Line(),
		jen.Qual(filepath.Join(proj.OutputPath, "internal/v1/tracing"), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		jen.Qual(filepath.Join(proj.OutputPath, "internal/v1/tracing"), "AttachItemIDToSpan").Call(jen.ID("span"), jen.ID("itemID")),
		jen.Line(),
		jen.ID("c").Dot("logger").Dot("WithValues").Call(
			typ.BuildGetSomethingLogValues(proj),
		).Dot("Debug").Call(jen.Litf("%s called", funcName)),
		jen.Line(),
		jen.Return(jen.ID("c").Dot("querier").Dotf("%sExists", sn).Call(args...)),
	)

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

func buildGetSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	uvn := n.UnexportedVarName()
	scnwp := n.SingularCommonNameWithPrefix()

	params := typ.BuildGetSomethingParams(pkg)
	args := typ.BuildGetSomethingArgs(pkg)
	loggerValues := typ.BuildGetSomethingLogValues(pkg)
	block := []jen.Code{
		jen.List(utils.CtxVar(), jen.ID("span")).Op(":=").Qual(utils.TracingLibrary, "StartSpan").Call(utils.CtxVar(), jen.Litf("Get%s", sn)),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	if typ.BelongsToUser {
		block = append(block, jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/tracing"), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")))
	}

	block = append(block,
		jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/tracing"), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.IDf("%sID", uvn)),
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
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Get%s", sn).Params(params...).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn), jen.ID("error")).Block(block...),
		jen.Line(),
	}
}

func buildGetSomethingCount(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pcn := n.PluralCommonName()
	params := typ.BuildGetListOfSomethingParams(pkg, false)
	args := typ.BuildGetListOfSomethingArgs(pkg)

	block := []jen.Code{
		jen.List(utils.CtxVar(), jen.ID("span")).Op(":=").Qual(utils.TracingLibrary, "StartSpan").Call(utils.CtxVar(), jen.Litf("Get%sCount", sn)),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	if typ.BelongsToUser {
		block = append(block, jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/tracing"), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")))
	}
	if typ.BelongsToStruct != nil {
		block = append(block, jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/tracing"), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())))
	}

	block = append(block,
		jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/tracing"), "AttachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
		jen.Line(),
	)

	loggerValues := typ.BuildGetListOfSomethingLogValues(pkg)
	block = append(block, jen.ID("c").Dot("logger").Dot("WithValues").Call(loggerValues).Dot("Debug").Call(jen.Litf("Get%sCount called", sn)))

	block = append(block,
		jen.Line(),
		jen.Return().ID("c").Dot("querier").Dotf("Get%sCount", sn).Call(args...),
	)

	return []jen.Code{
		jen.Commentf("Get%sCount fetches the count of %s from the database that meet a particular filter", sn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Get%sCount", sn).Params(params...).Params(jen.ID("count").ID("uint64"), jen.Err().ID("error")).Block(block...),
		jen.Line(),
	}
}

func buildGetAllSomethingCount(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	pn := n.Plural()
	pcn := n.PluralCommonName()

	return []jen.Code{
		jen.Commentf("GetAll%sCount fetches the count of %s from the database that meet a particular filter", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("GetAll%sCount", pn).Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("count").ID("uint64"), jen.Err().ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Op(":=").Qual(utils.TracingLibrary, "StartSpan").Call(utils.CtxVar(), jen.Litf("GetAll%sCount", pn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Litf("GetAll%sCount called", pn)),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dotf("GetAll%sCount", pn).Call(utils.CtxVar()),
		),
		jen.Line(),
	}
}

func buildGetListOfSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	uvn := n.UnexportedVarName()
	pn := n.Plural()
	pcn := n.PluralCommonName()

	params := typ.BuildGetListOfSomethingParams(pkg, false)
	callArgs := typ.BuildGetListOfSomethingArgs(pkg)
	block := []jen.Code{
		jen.List(utils.CtxVar(), jen.ID("span")).Op(":=").Qual(utils.TracingLibrary, "StartSpan").Call(utils.CtxVar(), jen.Litf("Get%s", pn)),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	if typ.BelongsToUser {
		block = append(block,
			jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/tracing"), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		)
	}
	if typ.BelongsToStruct != nil {
		block = append(block,
			jen.IDf("Attach%sIDToSpan", typ.BelongsToStruct.Singular()).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
		)
	}

	block = append(block,
		jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/tracing"), "AttachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
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
		jen.List(jen.IDf("%sList", uvn), jen.Err()).Op(":=").ID("c").Dot("querier").Dotf("Get%s", pn).Call(callArgs...),
		jen.Line(),
		jen.Return().List(jen.IDf("%sList", uvn), jen.Err()),
	)

	return []jen.Code{
		jen.Commentf("Get%s fetches a list of %s from the database that meet a particular filter", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Get%s", pn).Params(params...).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", sn)), jen.ID("error")).Block(block...),
		jen.Line(),
	}
}

func buildGetAllSomethingForUser(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	uvn := n.UnexportedVarName()
	pn := n.Plural()
	pcn := n.PluralCommonName()

	return []jen.Code{
		jen.Commentf("GetAll%sForUser fetches a list of %s from the database that meet a particular filter", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("GetAll%sForUser", pn).Params(utils.CtxParam(), jen.ID("userID").ID("uint64")).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Op(":=").Qual(utils.TracingLibrary, "StartSpan").Call(utils.CtxVar(), jen.Litf("GetAll%sForUser", pn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/tracing"), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Litf("GetAll%sForUser called", pn)),
			jen.Line(),
			jen.List(jen.IDf("%sList", uvn), jen.Err()).Op(":=").ID("c").Dot("querier").Dotf("GetAll%sForUser", pn).Call(utils.CtxVar(), jen.ID("userID")),
			jen.Line(),
			jen.Return().List(jen.IDf("%sList", uvn), jen.Err()),
		),
		jen.Line(),
	}
}

func buildGetAllSomethingForSomethingElse(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	uvn := n.UnexportedVarName()
	pn := n.Plural()
	pcn := n.PluralCommonName()

	btsn := typ.BelongsToStruct
	btsns := btsn.Singular()
	btsuvn := btsn.UnexportedVarName()
	btsrn := btsn.RouteName()

	params := typ.BuildGetSomethingForSomethingElseParams(pkg)
	args := typ.BuildGetSomethingForSomethingElseArgs(pkg)

	return []jen.Code{
		jen.Commentf("GetAll%sFor%s fetches a list of %s from the database that meet a particular filter", pn, btsns, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("GetAll%sFor%s", pn, btsns).Params(
			params...,
		).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Op(":=").Qual(utils.TracingLibrary, "StartSpan").Call(utils.CtxVar(), jen.Litf("GetAll%sFor%s", pn, btsns)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.IDf("Attach%sIDToSpan", btsns).Call(jen.ID("span"), jen.IDf("%sID", btsuvn)),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Litf("%s_id", btsrn), jen.IDf("%sID", btsuvn)).Dot("Debug").Call(jen.Litf("GetAll%sFor%s called", pn, btsns)),
			jen.Line(),
			jen.List(jen.IDf("%sList", uvn), jen.Err()).Op(":=").ID("c").Dot("querier").Dotf("GetAll%sFor%s", pn, btsns).Call(
				args...,
			),
			jen.Line(),
			jen.Return().List(jen.IDf("%sList", uvn), jen.Err()),
		),
		jen.Line(),
	}
}

func buildCreateSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	scnwp := n.SingularCommonNameWithPrefix()

	params := typ.BuildCreateSomethingParams(pkg, false)
	args := typ.BuildCreateSomethingArgs(pkg)

	return []jen.Code{
		jen.Commentf("Create%s creates %s in the database", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Create%s", sn).Params(
			params...,
		).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Op(":=").Qual(utils.TracingLibrary, "StartSpan").Call(utils.CtxVar(), jen.Litf("Create%s", sn)),
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

func buildUpdateSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	rn := n.RouteName()
	scn := n.SingularCommonName()

	const updatedVarName = "updated"

	params := typ.BuildUpdateSomethingParams(pkg, updatedVarName, false)
	args := typ.BuildUpdateSomethingArgs(pkg, updatedVarName)

	return []jen.Code{
		jen.Commentf("Update%s updates a particular %s. Note that Update%s expects the", sn, scn, sn),
		jen.Line(),
		jen.Comment("provided input to have a valid ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Update%s", sn).Params(
			params...,
		).Params(jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Op(":=").Qual(utils.TracingLibrary, "StartSpan").Call(utils.CtxVar(), jen.Litf("Update%s", sn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/tracing"), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.ID(updatedVarName).Dot("ID")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Litf("%s_id", rn), jen.ID(updatedVarName).Dot("ID")).Dot("Debug").Call(jen.Litf("Update%s called", sn)),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dotf("Update%s", sn).Call(
				args...,
			),
		),
		jen.Line(),
	}
}

func buildArchiveSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	rn := n.RouteName()
	uvn := n.UnexportedVarName()
	scnwp := n.SingularCommonNameWithPrefix()

	params := typ.BuildGetSomethingParams(pkg)
	callArgs := typ.BuildGetSomethingArgs(pkg)
	loggerValues := []jen.Code{
		jen.Litf("%s_id", rn).Op(":").IDf("%sID", uvn),
	}

	if typ.BelongsToUser {
		loggerValues = append(loggerValues, jen.Lit("user_id").Op(":").ID("userID"))
	}
	if typ.BelongsToStruct != nil {
		loggerValues = append(loggerValues, jen.Litf("%s_id", typ.BelongsToStruct.RouteName()).Op(":").IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}

	block := []jen.Code{
		jen.List(utils.CtxVar(), jen.ID("span")).Op(":=").Qual(utils.TracingLibrary, "StartSpan").Call(utils.CtxVar(), jen.Litf("Archive%s", sn)),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	if typ.BelongsToUser {
		block = append(block, jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/tracing"), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")))
	}
	if typ.BelongsToStruct != nil {
		block = append(block, jen.IDf("Attach%sIDToSpan", typ.BelongsToStruct.Singular()).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())))
	}

	block = append(block,
		jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/tracing"), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(jen.ID("span"), jen.IDf("%sID", uvn)),
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
