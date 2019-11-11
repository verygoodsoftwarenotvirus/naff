package client

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	v1 = "V1Client"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"client/v1/http/doc.go":                 docDotGo(),
		"client/v1/http/client.go":              mainDotGo(pkgRoot, types),
		"client/v1/http/client_test.go":         mainTestDotGo(pkgRoot, types),
		"client/v1/http/helpers.go":             helpersDotGo(pkgRoot, types),
		"client/v1/http/helpers_test.go":        helpersTestDotGo(pkgRoot, types),
		"client/v1/http/users.go":               usersDotGo(pkgRoot, types),
		"client/v1/http/users_test.go":          usersTestDotGo(pkgRoot, types),
		"client/v1/http/roundtripper.go":        roundtripperDotGo(pkgRoot, types),
		"client/v1/http/webhooks.go":            webhooksDotGo(pkgRoot, types),
		"client/v1/http/webhooks_test.go":       webhooksTestDotGo(pkgRoot, types),
		"client/v1/http/oauth2_clients.go":      oauth2ClientsDotGo(pkgRoot, types),
		"client/v1/http/oauth2_clients_test.go": oauth2ClientsTestDotGo(pkgRoot, types),
	}

	for _, typ := range types {
		files[fmt.Sprintf("client/v1/http/%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(pkgRoot, typ)
		files[fmt.Sprintf("client/v1/http/%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(pkgRoot, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
