package fakes

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func delegatedClientDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("BuildFakeAPIClient builds a faked APIClient."),
		jen.Line(),
		jen.Func().ID("BuildFakeAPIClient").Params().Params(jen.Op("*").ID("types").Dot("APIClient")).Body(
			jen.Return().Op("&").ID("types").Dot("APIClient").Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Uint32").Call()), jen.ID("ExternalID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "UUID").Call(), jen.ID("Name").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Password").Call(
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("false"),
				jen.ID("false"),
				jen.Lit(32),
			), jen.ID("ClientID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "UUID").Call(), jen.ID("ClientSecret").Op(":").Index().ID("byte").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Password").Call(
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.Lit(32),
			)), jen.ID("BelongsToUser").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call(), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.ID("uint32").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Date").Call().Dot("Unix").Call())))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAPIClientCreationResponseFromClient builds a faked APIClientCreationResponse."),
		jen.Line(),
		jen.Func().ID("BuildFakeAPIClientCreationResponseFromClient").Params(jen.ID("client").Op("*").ID("types").Dot("APIClient")).Params(jen.Op("*").ID("types").Dot("APIClientCreationResponse")).Body(
			jen.Return().Op("&").ID("types").Dot("APIClientCreationResponse").Valuesln(jen.ID("ID").Op(":").ID("client").Dot("ID"), jen.ID("ClientID").Op(":").ID("client").Dot("ClientID"), jen.ID("ClientSecret").Op(":").ID("string").Call(jen.ID("client").Dot("ClientSecret")))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAPIClientList builds a faked APIClientList."),
		jen.Line(),
		jen.Func().ID("BuildFakeAPIClientList").Params().Params(jen.Op("*").ID("types").Dot("APIClientList")).Body(
			jen.Var().Defs(
				jen.ID("examples").Index().Op("*").ID("types").Dot("APIClient"),
			),
			jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("exampleQuantity"), jen.ID("i").Op("++")).Body(
				jen.ID("examples").Op("=").ID("append").Call(
					jen.ID("examples"),
					jen.ID("BuildFakeAPIClient").Call(),
				)),
			jen.Return().Op("&").ID("types").Dot("APIClientList").Valuesln(jen.ID("Pagination").Op(":").ID("types").Dot("Pagination").Valuesln(jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("FilteredCount").Op(":").ID("exampleQuantity").Op("/").Lit(2), jen.ID("TotalCount").Op(":").ID("exampleQuantity")), jen.ID("Clients").Op(":").ID("examples")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAPIClientCreationInput builds a faked APIClientCreationInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeAPIClientCreationInput").Params().Params(jen.Op("*").ID("types").Dot("APIClientCreationInput")).Body(
			jen.ID("client").Op(":=").ID("BuildFakeAPIClient").Call(),
			jen.Return().Op("&").ID("types").Dot("APIClientCreationInput").Valuesln(jen.ID("UserLoginInput").Op(":").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Username").Call(), jen.ID("Password").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Password").Call(
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.Lit(32),
			), jen.ID("TOTPToken").Op(":").Qual("fmt", "Sprintf").Call(
				jen.Lit("0%s"),
				jen.Qual("github.com/brianvoe/gofakeit/v5", "Zip").Call(),
			)), jen.ID("Name").Op(":").ID("client").Dot("Name"), jen.ID("ClientID").Op(":").ID("client").Dot("ClientID"), jen.ID("BelongsToUser").Op(":").ID("client").Dot("BelongsToUser")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAPIClientCreationInputFromClient builds a faked APIClientCreationInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeAPIClientCreationInputFromClient").Params(jen.ID("client").Op("*").ID("types").Dot("APIClient")).Params(jen.Op("*").ID("types").Dot("APIClientCreationInput")).Body(
			jen.Return().Op("&").ID("types").Dot("APIClientCreationInput").Valuesln(jen.ID("UserLoginInput").Op(":").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Username").Call(), jen.ID("Password").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Password").Call(
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.Lit(32),
			), jen.ID("TOTPToken").Op(":").Qual("fmt", "Sprintf").Call(
				jen.Lit("0%s"),
				jen.Qual("github.com/brianvoe/gofakeit/v5", "Zip").Call(),
			)), jen.ID("Name").Op(":").ID("client").Dot("Name"), jen.ID("ClientID").Op(":").ID("client").Dot("ClientID"), jen.ID("ClientSecret").Op(":").ID("client").Dot("ClientSecret"), jen.ID("BelongsToUser").Op(":").ID("client").Dot("BelongsToUser"))),
		jen.Line(),
	)

	return code
}
