package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func cookieauthDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Comment("CookieAuth represents what we encode in our authentication cookies"),
		jen.Line(),
		jen.Type().ID("CookieAuth").Struct(jen.ID("UserID").ID("uint64"), jen.ID("Admin").ID("bool"), jen.ID("Username").ID("string")),
		jen.Line(),
	)
	return ret
}
