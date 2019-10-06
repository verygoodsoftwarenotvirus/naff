package load

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func itemsDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Func())
	ret.Add(jen.Func().ID("buildItemActions").Params(jen.ID("c").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http", "V1Client")).Params(jen.Map(jen.ID("string")).Op("*").ID("Action")).Block(
		jen.Return().Map(jen.ID("string")).Op("*").ID("Action").Valuesln(jen.Lit("CreateItem").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("CreateItem"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
			jen.Return().ID("c").Dot(
				"BuildCreateItemRequest",
			).Call(jen.Qual("context", "Background").Call(), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/rand/model", "RandomItemCreationInput").Call()),
		), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("GetItem").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("GetItem"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
			jen.If(
				jen.ID("randomItem").Op(":=").ID("fetchRandomItem").Call(jen.ID("c")),
				jen.ID("randomItem").Op("!=").ID("nil"),
			).Block(
				jen.Return().ID("c").Dot(
					"BuildGetItemRequest",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("randomItem").Dot(
					"ID",
				)),
			),
			jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
		), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("GetItems").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("GetItems"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
			jen.Return().ID("c").Dot(
				"BuildGetItemsRequest",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("nil")),
		), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("UpdateItem").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("UpdateItem"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
			jen.If(
				jen.ID("randomItem").Op(":=").ID("fetchRandomItem").Call(jen.ID("c")),
				jen.ID("randomItem").Op("!=").ID("nil"),
			).Block(
				jen.ID("randomItem").Dot(
					"Name",
				).Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/rand/model", "RandomItemCreationInput").Call().Dot(
					"Name",
				),
				jen.Return().ID("c").Dot(
					"BuildUpdateItemRequest",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("randomItem")),
			),
			jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
		), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("ArchiveItem").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("ArchiveItem"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
			jen.If(
				jen.ID("randomItem").Op(":=").ID("fetchRandomItem").Call(jen.ID("c")),
				jen.ID("randomItem").Op("!=").ID("nil"),
			).Block(
				jen.Return().ID("c").Dot(
					"BuildArchiveItemRequest",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("randomItem").Dot(
					"ID",
				)),
			),
			jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
		), jen.ID("Weight").Op(":").Lit(85))),
	),
	)
	return ret
}
