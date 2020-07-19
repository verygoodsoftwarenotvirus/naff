package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func cookieauthDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("models")

	utils.AddImports(proj, code)

	code.Add(
		jen.Comment("CookieAuth represents what we encode in our authentication cookies."),
		jen.Line(),
		jen.Type().ID("CookieAuth").Struct(constants.UserIDParam(), jen.ID("Admin").Bool(), jen.ID("Username").String()),
		jen.Line(),
	)

	return code
}
