package model

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("randmodel")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Comment("RandomWebhookInput creates a random WebhookCreationInput"),
		jen.Line(),
		jen.Func().ID("RandomWebhookInput").Params().Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput")).Block(
			jen.ID("x").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput").Valuesln(
				jen.ID("Name").MapAssign().Qual(utils.FakeLibrary, "Word").Call(),
				jen.ID("URL").MapAssign().Qual(utils.FakeLibrary, "DomainName").Call(),
				jen.ID("ContentType").MapAssign().Lit("application/json"),
				jen.ID("Method").MapAssign().Lit("POST"),
			),
			jen.Return().ID("x"),
		),
		jen.Line(),
	)
	return ret
}
