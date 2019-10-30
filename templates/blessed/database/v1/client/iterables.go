package client

import (
	"fmt"
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(pkgRoot string, typ models.DataType) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(ret)

	n := typ.Name
	sn := n.Singular()
	rn := n.RouteName()
	uvn := n.UnexportedVarName()
	scn := n.SingularCommonName()
	scnwp := n.SingularCommonNameWithPrefix()
	pn := n.Plural()
	pcn := n.PluralCommonName()

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(pkgRoot, "models/v1"), fmt.Sprintf("%sDataManager", sn)).Op("=").Parens(jen.Op("*").ID("Client")).Call(jen.ID("nil")),
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

	ret.Add(
		jen.Commentf("Get%s fetches %s from the database", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Get%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.IDf("%sID", uvn), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Litf("Get%s", sn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.IDf("attach%sIDToSpan", sn).Call(jen.ID("span"), jen.IDf("%sID", uvn)),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Litf("%s_id", rn).Op(":").IDf("%sID", uvn),
				jen.Lit("user_id").Op(":").ID("userID"),
			)).Dot("Debug").Call(jen.Litf("Get%s called", sn)),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dotf("Get%s", sn).Call(jen.ID("ctx"), jen.IDf("%sID", uvn), jen.ID("userID")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Get%sCount fetches the count of %s from the database that meet a particular filter", sn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Get%sCount", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Litf("Get%sCount", sn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("attachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Litf("Get%sCount called", sn)),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dotf("Get%sCount", sn).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
		),
		jen.Line(),
	)

	ret.Add(
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
	)

	ret.Add(
		jen.Commentf("Get%s fetches a list of %s from the database that meet a particular filter", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Get%s", pn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("filter").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "QueryFilter"),
			jen.ID("userID").ID("uint64"),
		).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), fmt.Sprintf("%sList", sn)), jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Litf("Get%s", pn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("attachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Litf("Get%s called", pn)),
			jen.Line(),
			jen.List(jen.IDf("%sList", uvn), jen.ID("err")).Op(":=").ID("c").Dot("querier").Dotf("Get%s", pn).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
			jen.Line(),
			jen.Return().List(jen.IDf("%sList", uvn), jen.ID("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("GetAll%sForUser fetches a list of %s from the database that meet a particular filter", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("GetAll%sForUser", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().Qual(filepath.Join(pkgRoot, "models/v1"), sn),
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
	)

	ret.Add(
		jen.Commentf("Create%s creates %s in the database", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Create%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), fmt.Sprintf("%sCreationInput", sn))).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Litf("Create%s", sn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("input"), jen.ID("input")).Dot("Debug").Call(jen.Litf("Create%s called", sn)),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dotf("Create%s", sn).Call(jen.ID("ctx"), jen.ID("input")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Update%s updates a particular %s. Note that Update%s expects the", sn, scn, sn),
		jen.Line(),
		jen.Comment("provided input to have a valid ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Update%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), sn)).Params(jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Litf("Update%s", sn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.IDf("attach%sIDToSpan", sn).Call(jen.ID("span"), jen.ID("input").Dot("ID")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Litf("%s_id", rn), jen.ID("input").Dot("ID")).Dot("Debug").Call(jen.Litf("Update%s called", sn)),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dotf("Update%s", sn).Call(jen.ID("ctx"), jen.ID("input")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Archive%s archives %s from the database by its ID", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).IDf("Archive%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.IDf("%sID", uvn), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Litf("Archive%s", sn)),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.IDf("attach%sIDToSpan", sn).Call(jen.ID("span"), jen.IDf("%sID", uvn)),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Litf("%s_id", rn).Op(":").IDf("%sID", uvn),
				jen.Lit("user_id").Op(":").ID("userID"),
			)).Dot("Debug").Call(jen.Litf("Archive%s called", sn)),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dotf("Archive%s", sn).Call(jen.ID("ctx"), jen.IDf("%sID", uvn), jen.ID("userID")),
		),
		jen.Line(),
	)
	return ret
}
