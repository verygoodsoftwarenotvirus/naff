package load

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(rootPkg string, typ models.DataType) *jen.File {
	ret := jen.NewFile("main")

	utils.AddImports(ret)

	ret.Add(
		jen.Comment("fetchRandomItem retrieves a random item from the list of available items"),
		jen.Line(),
		jen.Func().ID("fetchRandomItem").Params(jen.ID("c").Op("*").Qual(filepath.Join(rootPkg, "client/v1/http"), "V1Client")).Params(jen.Op("*").Qual(filepath.Join(rootPkg, "models/v1"), "Item")).Block(
			jen.List(jen.ID("itemsRes"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetItems",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("nil")),
			jen.If(jen.ID("err").Op("!=").ID("nil").Op("||").ID("itemsRes").Op("==").ID("nil").Op("||").ID("len").Call(jen.ID("itemsRes").Dot("Items")).Op("==").Lit(0)).Block(
				jen.Return().ID("nil"),
			),
			jen.ID("randIndex").Op(":=").Qual("math/rand", "Intn").Call(jen.ID("len").Call(jen.ID("itemsRes").Dot("Items"))),
			jen.Return().Op("&").ID("itemsRes").Dot("Items").Index(jen.ID("randIndex")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildItemActions").Params(jen.ID("c").Op("*").Qual(filepath.Join(rootPkg, "client/v1/http"), "V1Client")).Params(jen.Map(jen.ID("string")).Op("*").ID("Action")).Block(
			jen.Return().Map(jen.ID("string")).Op("*").ID("Action").Valuesln(
				jen.Lit("CreateItem").Op(":").Valuesln(
					jen.ID("Name").Op(":").Lit("CreateItem"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.Return().ID("c").Dot("BuildCreateItemRequest").Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(rootPkg, "tests/v1/testutil/rand/model"), "RandomItemCreationInput").Call()),
					),
					jen.ID("Weight").Op(":").Lit(100)), jen.Lit("GetItem").Op(":").Valuesln(
					jen.ID("Name").Op(":").Lit("GetItem"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.If(jen.ID("randomItem").Op(":=").ID("fetchRandomItem").Call(jen.ID("c")), jen.ID("randomItem").Op("!=").ID("nil")).Block(
							jen.Return().ID("c").Dot("BuildGetItemRequest").Call(jen.Qual("context", "Background").Call(), jen.ID("randomItem").Dot("ID")),
						),
						jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
					),
					jen.ID("Weight").Op(":").Lit(100)), jen.Lit("GetItems").Op(":").Valuesln(
					jen.ID("Name").Op(":").Lit("GetItems"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.Return().ID("c").Dot("BuildGetItemsRequest").Call(jen.Qual("context", "Background").Call(), jen.ID("nil")),
					),
					jen.ID("Weight").Op(":").Lit(100)), jen.Lit("UpdateItem").Op(":").Valuesln(
					jen.ID("Name").Op(":").Lit("UpdateItem"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.If(jen.ID("randomItem").Op(":=").ID("fetchRandomItem").Call(jen.ID("c")), jen.ID("randomItem").Op("!=").ID("nil")).Block(
							jen.ID("randomItem").Dot("Name").Op("=").Qual(filepath.Join(rootPkg, "tests/v1/testutil/rand/model"), "RandomItemCreationInput").Call().Dot("Name"),
							jen.Return().ID("c").Dot("BuildUpdateItemRequest").Call(jen.Qual("context", "Background").Call(), jen.ID("randomItem")),
						),
						jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
					),
					jen.ID("Weight").Op(":").Lit(100)), jen.Lit("ArchiveItem").Op(":").Valuesln(
					jen.ID("Name").Op(":").Lit("ArchiveItem"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.If(jen.ID("randomItem").Op(":=").ID("fetchRandomItem").Call(jen.ID("c")), jen.ID("randomItem").Op("!=").ID("nil")).Block(
							jen.Return().ID("c").Dot("BuildArchiveItemRequest").Call(jen.Qual("context", "Background").Call(), jen.ID("randomItem").Dot("ID")),
						),
						jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
					),
					jen.ID("Weight").Op(":").Lit(85))),
		),
		jen.Line(),
	)
	return ret
}
