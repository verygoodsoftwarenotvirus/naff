package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo() *jen.File {
	ret := jen.NewFile("webhooks")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideWebhooksService"), jen.ID("ProvideWebhookDataManager"), jen.ID("ProvideWebhookDataServer")),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
