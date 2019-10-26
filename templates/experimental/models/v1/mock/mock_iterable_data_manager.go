package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockIterableDataManagerDotGo(typ models.DataType) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(ret)

	n := typ.Name
	sn := n.Singular()
	pn := n.Singular()
	uvn := n.UnexportedVarName()

	ret.Add(
		jen.Var().ID("_").ID("models").Dotf("%sDataManager", sn).Op("=").Parens(jen.Op("*").IDf("%sDataManager", sn)).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Type().IDf("%sDataManager", sn).Struct(
			jen.Qual("github.com/stretchr/testify/mock", "Mock"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Get%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Get%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.IDf("%sID", uvn), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("models").Dot(sn),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.IDf("%sID", uvn), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("models").Dot(sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Get%sCount is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Get%sCount", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot("QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllItemsCount is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("GetAll%sCount", pn).Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetItems is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Get%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot("QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("models").Dot("ItemList"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("models").Dot("ItemList")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllItemsForUser is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("GetAll%sForUser", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().ID("models").Dot("Item"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().ID("models").Dot("Item")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateItem is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Create%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("models").Dot("ItemCreationInput")).Params(jen.Op("*").ID("models").Dot("Item"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("input")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("models").Dot("Item")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateItem is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Update%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("models").Dot("Item")).Params(jen.ID("error")).Block(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("updated")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveItem is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Archive%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("id"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("id"), jen.ID("userID")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)
	return ret
}
