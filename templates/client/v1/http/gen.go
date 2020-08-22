package client

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	v1 = "V1Client"

	packageName = "client"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"client/v1/http/doc.go":                   docDotGo(proj),
		"client/v1/http/client.go":                mainDotGo(proj),
		"client/v1/http/client_test.go":           mainTestDotGo(proj),
		"client/v1/http/helpers.go":               helpersDotGo(proj),
		"client/v1/http/helpers_test.go":          helpersTestDotGo(proj),
		"client/v1/http/users.go":                 usersDotGo(proj),
		"client/v1/http/users_test.go":            usersTestDotGo(proj),
		"client/v1/http/roundtripper.go":          roundtripperDotGo(proj),
		"client/v1/http/roundtripper_test.go":     roundtripperTestDotGo(proj),
		"client/v1/http/mock_read_closer_test.go": mockReadCloserTestDotGo(proj),
		"client/v1/http/webhooks.go":              webhooksDotGo(proj),
		"client/v1/http/webhooks_test.go":         webhooksTestDotGo(proj),
		"client/v1/http/oauth2_clients.go":        oauth2ClientsDotGo(proj),
		"client/v1/http/oauth2_clients_test.go":   oauth2ClientsTestDotGo(proj),
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("client/v1/http/%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(proj, typ)
		files[fmt.Sprintf("client/v1/http/%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ)
	}

	for path, file := range files {
		// fmt.Printf("rendering %q\n", path)
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
