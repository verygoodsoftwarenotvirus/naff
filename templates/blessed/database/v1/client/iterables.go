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

	utils.AddImports(pkg.OutputPath, []models.DataType{typ}, ret)

	n := typ.Name
	sn := n.Singular()
	rn := n.RouteName()
	uvn := n.UnexportedVarName()
	scnwp := n.SingularCommonNameWithPrefix()

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sDataManager", sn)).Op("=").Parens(jen.Op("*").ID("Client")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("attach%sIDToSpan provides a consistent way to attach %s's ID to a span", sn, scnwp),
		jen.Line(),
		jen.Func().IDf("attach%sIDToSpan", sn).Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.IDf("%sID", uvn).ID("uint64")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot("AddAttributes").Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Litf("%s_id", rn), jen.Qual("strconv", "FormatUint").Call(jen.IDf("%sID", uvn), jen.Lit(10)))),
			),
		),
		jen.Line(),
	)

	ret.Add(buildGetSomething(pkg, typ)...)
	ret.Add(buildGetSomethingCount(pkg, typ)...)
	ret.Add(buildGetAllSomethingCount(pkg, typ)...)
	ret.Add(buildGetListOfSomething(pkg, typ)...)

	if typ.BelongsToUser {
		ret.Add(buildGetAllSomethingForUser(pkg, typ)...)
	} else if typ.BelongsToStruct != nil {
		ret.Add(buildGetAllSomethingForSomethingElse(pkg, typ)...)
	}

	ret.Add(buildCreateSomething(pkg, typ)...)
	ret.Add(buildUpdateSomething(pkg, typ)...)
	ret.Add(buildArchiveSomething(pkg, typ)...)

	return ret
}

func buildGetSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	rn := n.RouteName()
	uvn := n.UnexportedVarName()
	scnwp := n.SingularCommonNameWithPrefix()

	args := []jen.Code{
		jen.ID("ctx").Qual("context", "Context"),
	}
	loggerValues := []jen.Code{
		jen.Litf("%s_id", rn).Op(":").IDf("%sID", uvn),
	}
	block := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Litf("Get%s", sn)),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}
	callArgs := []jen.Code{
		jen.ID("ctx"),
		jen.IDf("%sID", uvn),
	}

	if typ.BelongsToUser {
		args = append(args, jen.List(jen.IDf("%sID", uvn), jen.ID("userID")).ID("uint64"))
		block = append(block, jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")))
		loggerValues = append(loggerValues, jen.Lit("user_id").Op(":").ID("userID"))
		callArgs = append(callArgs, jen.ID("userID"))
	} else if typ.BelongsToStruct != nil {
		args = append(args, jen.List(jen.IDf("%sID", uvn), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).ID("uint64")))
		loggerValues = append(loggerValues, jen.Litf("%s_id", typ.BelongsToStruct.RouteName()).Op(":").IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		callArgs = append(callArgs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	} else {
		args = append(args, jen.IDf("%sID", uvn).ID("uint64"))
	}

	block = append(block,
		jen.IDf("attach%sIDToSpan", sn).Call(jen.ID("span"), jen.IDf("%sID", uvn)),
		jen.Line(),
		jen.ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(loggerValues...)).Dot("Debug").Call(jen.Litf("Get%s called", sn)),
		jen.Line(),
		jen.Return().ID("c").Dot("querier").Dotf("Get%s", sn).Call(callArgs...),
	)

	return []jen.Code{
		jen.Commentf("Get%s fetches %s from the database", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Get%s", sn).Params(args...).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn), jen.ID("error")).Block(block...),
		jen.Line(),
	}
}

func buildGetSomethingCount(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pcn := n.PluralCommonName()

	params := []jen.Code{
		jen.ID("ctx").Qual("context", "Context"),
		jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"),
	}

	block := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Litf("Get%sCount", sn)),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	if typ.BelongsToUser {
		params = append(params, jen.ID("userID").ID("uint64"))
		block = append(block, jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")))
	} else if typ.BelongsToStruct != nil {
		params = append(params, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).ID("uint64"))
		block = append(block, jen.IDf("attach%sIDToSpan", typ.BelongsToStruct.Singular()).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())))
	}
	// we don't need to consider the other case

	block = append(block,
		jen.ID("attachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
		jen.Line(),
	)

	if typ.BelongsToUser {
		block = append(block, jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Litf("Get%sCount called", sn)))
	} else if typ.BelongsToStruct != nil {
		block = append(block, jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Litf("%s_id", typ.BelongsToStruct.RouteName()), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())).Dot("Debug").Call(jen.Litf("Get%sCount called", sn)))
	}
	block = append(block, jen.Line())

	if typ.BelongsToUser {
		block = append(block,
			jen.Return().ID("c").Dot("querier").Dotf("Get%sCount", sn).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
		)
	} else if typ.BelongsToStruct != nil {
		block = append(block,
			jen.Return().ID("c").Dot("querier").Dotf("Get%sCount", sn).Call(jen.ID("ctx"), jen.ID("filter"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
		)
	} else {
		block = append(block,
			jen.Return().ID("c").Dot("querier").Dotf("Get%sCount", sn).Call(jen.ID("ctx"), jen.ID("filter")),
		)
	}

	return []jen.Code{
		jen.Commentf("Get%sCount fetches the count of %s from the database that meet a particular filter", sn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Get%sCount", sn).Params(params...).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(block...),
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
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("GetAll%sCount", pn).Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Litf("GetAll%sCount", pn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Litf("GetAll%sCount called", pn)),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dotf("GetAll%sCount", pn).Call(jen.ID("ctx")),
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

	params := []jen.Code{
		jen.ID("ctx").Qual("context", "Context"),
		jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"),
	}
	callArgs := []jen.Code{
		jen.ID("ctx"),
		jen.ID("filter"),
	}
	block := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Litf("Get%s", pn)),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	if typ.BelongsToUser {
		jen.ID("userID").ID("uint64")
		block = append(block,
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		)
		callArgs = append(callArgs, jen.ID("userID"))
		params = append(params, jen.ID("userID").ID("uint64"))
	} else if typ.BelongsToStruct != nil {
		jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).ID("uint64")
		block = append(block,
			jen.IDf("attach%sIDToSpan", typ.BelongsToStruct.Singular()).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
		)
		params = append(params, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).ID("uint64"))
		callArgs = append(callArgs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}

	block = append(block,
		jen.ID("attachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
		jen.Line(),
	)

	if typ.BelongsToUser {
		block = append(block,
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Litf("Get%s called", pn)),
		)
	} else if typ.BelongsToStruct != nil {
		block = append(block,
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Litf("%s_id", typ.BelongsToStruct.RouteName()), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())).Dot("Debug").Call(jen.Litf("Get%s called", pn)),
		)
	} else {
		block = append(block,
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Litf("Get%s called", pn)),
		)
	}

	block = append(block,
		jen.Line(),
		jen.List(jen.IDf("%sList", uvn), jen.ID("err")).Op(":=").ID("c").Dot("querier").Dotf("Get%s", pn).Call(callArgs...),
		jen.Line(),
		jen.Return().List(jen.IDf("%sList", uvn), jen.ID("err")),
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
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("GetAll%sForUser", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Litf("GetAll%sForUser", pn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Litf("GetAll%sForUser called", pn)),
			jen.Line(),
			jen.List(jen.IDf("%sList", uvn), jen.ID("err")).Op(":=").ID("c").Dot("querier").Dotf("GetAll%sForUser", pn).Call(jen.ID("ctx"), jen.ID("userID")),
			jen.Line(),
			jen.Return().List(jen.IDf("%sList", uvn), jen.ID("err")),
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

	return []jen.Code{
		jen.Commentf("GetAll%sFor%s fetches a list of %s from the database that meet a particular filter", pn, btsns, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("GetAll%sFor%s", pn, btsns).Params(jen.ID("ctx").Qual("context", "Context"), jen.IDf("%sID", btsuvn).ID("uint64")).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Litf("GetAll%sFor%s", pn, btsns)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.IDf("attach%sIDToSpan", btsns).Call(jen.ID("span"), jen.IDf("%sID", btsuvn)),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Litf("%s_id", btsrn), jen.IDf("%sID", btsuvn)).Dot("Debug").Call(jen.Litf("GetAll%sFor%s called", pn, btsns)),
			jen.Line(),
			jen.List(jen.IDf("%sList", uvn), jen.ID("err")).Op(":=").ID("c").Dot("querier").Dotf("GetAll%sFor%s", pn, btsns).Call(jen.ID("ctx"), jen.IDf("%sID", btsuvn)),
			jen.Line(),
			jen.Return().List(jen.IDf("%sList", uvn), jen.ID("err")),
		),
		jen.Line(),
	}
}

func buildCreateSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	scnwp := n.SingularCommonNameWithPrefix()

	return []jen.Code{
		jen.Commentf("Create%s creates %s in the database", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Create%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", sn))).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Litf("Create%s", sn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("input"), jen.ID("input")).Dot("Debug").Call(jen.Litf("Create%s called", sn)),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dotf("Create%s", sn).Call(jen.ID("ctx"), jen.ID("input")),
		),
		jen.Line(),
	}
}

func buildUpdateSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	rn := n.RouteName()
	scn := n.SingularCommonName()

	return []jen.Code{
		jen.Commentf("Update%s updates a particular %s. Note that Update%s expects the", sn, scn, sn),
		jen.Line(),
		jen.Comment("provided input to have a valid ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Update%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)).Params(jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Litf("Update%s", sn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.IDf("attach%sIDToSpan", sn).Call(jen.ID("span"), jen.ID("input").Dot("ID")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Litf("%s_id", rn), jen.ID("input").Dot("ID")).Dot("Debug").Call(jen.Litf("Update%s called", sn)),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dotf("Update%s", sn).Call(jen.ID("ctx"), jen.ID("input")),
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

	params := []jen.Code{
		jen.ID("ctx").Qual("context", "Context"),
	}
	loggerValues := []jen.Code{
		jen.Litf("%s_id", rn).Op(":").IDf("%sID", uvn),
	}
	callArgs := []jen.Code{
		jen.ID("ctx"),
		jen.IDf("%sID", uvn),
	}

	if typ.BelongsToUser {
		params = append(params, jen.List(jen.IDf("%sID", uvn), jen.ID("userID")).ID("uint64"))
		loggerValues = append(loggerValues, jen.Lit("user_id").Op(":").ID("userID"))
		callArgs = append(callArgs, jen.ID("userID"))
	} else if typ.BelongsToStruct != nil {
		params = append(params, jen.List(jen.IDf("%sID", uvn), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())).ID("uint64"))
		loggerValues = append(loggerValues, jen.Litf("%s_id", typ.BelongsToStruct.RouteName()).Op(":").IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		callArgs = append(callArgs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	} else {
		params = append(params, jen.IDf("%sID", uvn).ID("uint64"))
	}

	block := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Litf("Archive%s", sn)),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	if typ.BelongsToUser {
		block = append(block, jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")))
	} else if typ.BelongsToStruct != nil {
		block = append(block, jen.IDf("attach%sIDToSpan", typ.BelongsToStruct.Singular()).Call(jen.ID("span"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())))
	}

	block = append(block,
		jen.IDf("attach%sIDToSpan", sn).Call(jen.ID("span"), jen.IDf("%sID", uvn)),
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
