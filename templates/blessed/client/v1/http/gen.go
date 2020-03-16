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
func RenderPackage(pkg *models.Project) error {
	files := map[string]*jen.File{
		"client/v1/http/doc.go":                 docDotGo(),
		"client/v1/http/client.go":              mainDotGo(pkg),
		"client/v1/http/client_test.go":         mainTestDotGo(pkg),
		"client/v1/http/helpers.go":             helpersDotGo(pkg),
		"client/v1/http/helpers_test.go":        helpersTestDotGo(pkg),
		"client/v1/http/users.go":               usersDotGo(pkg),
		"client/v1/http/users_test.go":          usersTestDotGo(pkg),
		"client/v1/http/roundtripper.go":        roundtripperDotGo(pkg),
		"client/v1/http/webhooks.go":            webhooksDotGo(pkg),
		"client/v1/http/webhooks_test.go":       webhooksTestDotGo(pkg),
		"client/v1/http/oauth2_clients.go":      oauth2ClientsDotGo(pkg),
		"client/v1/http/oauth2_clients_test.go": oauth2ClientsTestDotGo(pkg),
	}

	for _, typ := range pkg.DataTypes {
		files[fmt.Sprintf("client/v1/http/%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(pkg, typ)
		files[fmt.Sprintf("client/v1/http/%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(pkg, typ)
	}

	for path, file := range files {
		// fmt.Printf("rendering %q\n", path)
		if err := utils.RenderGoFile(pkg, path, file); err != nil {
			return err
		}
	}

	return nil
}
