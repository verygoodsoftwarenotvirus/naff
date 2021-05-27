package fakes

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func adminDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("BuildFakeUserReputationUpdateInput builds a faked ItemCreationInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeUserReputationUpdateInput").Params().Params(jen.Op("*").ID("types").Dot("UserReputationUpdateInput")).Body(
			jen.Return().Op("&").ID("types").Dot("UserReputationUpdateInput").Valuesln(jen.ID("TargetUserID").Op(":").ID("uint64").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Uint32").Call()), jen.ID("NewReputation").Op(":").ID("types").Dot("GoodStandingAccountStatus"), jen.ID("Reason").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Sentence").Call(jen.Lit(10)))),
		jen.Line(),
	)

	return code
}
