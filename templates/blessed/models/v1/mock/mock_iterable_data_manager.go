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
	pn := n.Plural()
	uvn := n.UnexportedVarName()

	ret.Add(
		jen.Var().ID("_").ID("models").Dotf("%sDataManager", sn).Op("=").Parens(jen.Op("*").IDf("%sDataManager", sn)).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("%sDataManager is a mocked models.%sDataManager for testing", sn, sn),
		jen.Line(),
		jen.Type().IDf("%sDataManager", sn).Struct(
			jen.Qual("github.com/stretchr/testify/mock", "Mock"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Get%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Get%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.IDf("%sID", uvn), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",sn),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.IDf("%sID", uvn), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Get%sCount is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Get%sCount", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1","QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("GetAll%sCount is a mock function", pn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("GetAll%sCount", pn).Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Get%s is a mock function", pn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Get%s", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1","QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("models").Dotf("%sList", sn),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("models").Dotf("%sList", sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("GetAll%sForUser is a mock function", pn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("GetAll%sForUser", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",sn),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Create%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Create%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("models").Dotf("%sCreationInput", sn)).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",sn),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("input")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Update%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Update%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",sn)).Params(jen.ID("error")).Block(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("updated")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Archive%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Archive%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("id"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("id"), jen.ID("userID")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)
	return ret
}
