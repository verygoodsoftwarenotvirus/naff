package load

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func itemsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("fetchRandomItem retrieves a random item from the list of available items."),
		jen.Line(),
		jen.Func().ID("fetchRandomItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("c").Op("*").ID("httpclient").Dot("Client")).Params(jen.Op("*").ID("types").Dot("Item")).Body(
			jen.List(jen.ID("itemsRes"), jen.ID("err")).Op(":=").ID("c").Dot("GetItems").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil").Op("||").ID("itemsRes").Op("==").ID("nil").Op("||").ID("len").Call(jen.ID("itemsRes").Dot("Items")).Op("==").Lit(0)).Body(
				jen.Return().ID("nil")),
			jen.ID("randIndex").Op(":=").Qual("math/rand", "Intn").Call(jen.ID("len").Call(jen.ID("itemsRes").Dot("Items"))),
			jen.Return().ID("itemsRes").Dot("Items").Index(jen.ID("randIndex")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildItemActions").Params(jen.ID("c").Op("*").ID("httpclient").Dot("Client"), jen.ID("builder").Op("*").ID("requests").Dot("Builder")).Params(jen.Map(jen.ID("string")).Op("*").ID("Action")).Body(
			jen.Return().Map(jen.ID("string")).Op("*").ID("Action").Valuesln(jen.Lit("CreateItem").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("CreateItem"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("itemInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInput").Call(),
				jen.Return().ID("builder").Dot("BuildCreateItemRequest").Call(
					jen.ID("ctx"),
					jen.ID("itemInput"),
				),
			), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("GetItem").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("GetItem"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("randomItem").Op(":=").ID("fetchRandomItem").Call(
					jen.ID("ctx"),
					jen.ID("c"),
				),
				jen.If(jen.ID("randomItem").Op("==").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("retrieving random item: %w"),
						jen.ID("ErrUnavailableYet"),
					))),
				jen.Return().ID("builder").Dot("BuildGetItemRequest").Call(
					jen.ID("ctx"),
					jen.ID("randomItem").Dot("ID"),
				),
			), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("GetItems").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("GetItems"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.Return().ID("builder").Dot("BuildGetItemsRequest").Call(
					jen.ID("ctx"),
					jen.ID("nil"),
				),
			), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("UpdateItem").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("UpdateItem"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.If(jen.ID("randomItem").Op(":=").ID("fetchRandomItem").Call(
					jen.ID("ctx"),
					jen.ID("c"),
				), jen.ID("randomItem").Op("!=").ID("nil")).Body(
					jen.ID("newItem").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInput").Call(),
					jen.ID("randomItem").Dot("Name").Op("=").ID("newItem").Dot("Name"),
					jen.ID("randomItem").Dot("Details").Op("=").ID("newItem").Dot("Details"),
					jen.Return().ID("builder").Dot("BuildUpdateItemRequest").Call(
						jen.ID("ctx"),
						jen.ID("randomItem"),
					),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
			), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("ArchiveItem").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("ArchiveItem"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("randomItem").Op(":=").ID("fetchRandomItem").Call(
					jen.ID("ctx"),
					jen.ID("c"),
				),
				jen.If(jen.ID("randomItem").Op("==").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("retrieving random item: %w"),
						jen.ID("ErrUnavailableYet"),
					))),
				jen.Return().ID("builder").Dot("BuildArchiveItemRequest").Call(
					jen.ID("ctx"),
					jen.ID("randomItem").Dot("ID"),
				),
			), jen.ID("Weight").Op(":").Lit(85)))),
		jen.Line(),
	)

	return code
}
