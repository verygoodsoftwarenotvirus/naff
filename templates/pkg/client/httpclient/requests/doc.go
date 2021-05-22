package requests

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func docDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	code.PackageCommentf(`Package %s provides an HTTP client that can communicate with and interpret the responses
of an instance of the %s service.`, packageName, proj.Name.SingularCommonName())

	return code
}
