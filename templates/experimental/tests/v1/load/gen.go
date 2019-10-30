package load

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"tests/v1/load/actions.go":       actionsDotGo(pkgRoot),
		"tests/v1/load/init.go":          initDotGo(pkgRoot),
		"tests/v1/load/main.go":          mainDotGo(pkgRoot),
		"tests/v1/load/oauth2clients.go": oauth2ClientsDotGo(pkgRoot),
		"tests/v1/load/webhooks.go":      webhooksDotGo(pkgRoot),
	}

	for _, typ := range types {
		files[fmt.Sprintf("tests/v1/load/%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(pkgRoot, typ)
	}

	for path, file := range files {
		if err := utils.RenderFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
