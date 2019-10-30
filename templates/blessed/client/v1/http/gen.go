package client

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

const (
	v1 = "V1Client"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"client/v1/http/doc.go":                 docDotGo(),
		"client/v1/http/client.go":              mainDotGo(),
		"client/v1/http/client_test.go":         mainTestDotGo(),
		"client/v1/http/helpers.go":             helpersDotGo(pkgRoot),
		"client/v1/http/helpers_test.go":        helpersTestDotGo(pkgRoot),
		"client/v1/http/users.go":               usersDotGo(),
		"client/v1/http/users_test.go":          usersTestDotGo(),
		"client/v1/http/roundtripper.go":        roundtripperDotGo(),
		"client/v1/http/webhooks.go":            webhooksDotGo(),
		"client/v1/http/webhooks_test.go":       webhooksTestDotGo(pkgRoot),
		"client/v1/http/oauth2_clients.go":      oauth2ClientsDotGo(),
		"client/v1/http/oauth2_clients_test.go": oauth2ClientsTestDotGo(),
	}

	for _, typ := range types {
		files[fmt.Sprintf("client/v1/http/%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(typ)
		files[fmt.Sprintf("client/v1/http/%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(pkgRoot, typ)
	}

	for path, file := range files {
		if err := utils.RenderFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
