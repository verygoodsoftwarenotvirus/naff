package requests

import (
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	v1 = "V1Client"

	packageName       = "requests"
	packagePathPrefix = "pkg/client/httpclient/requests"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"doc.go":                   docDotGo(proj),
		"client.go":                mainDotGo(proj),
		"client_test.go":           mainTestDotGo(proj),
		"helpers.go":               helpersDotGo(proj),
		"helpers_test.go":          helpersTestDotGo(proj),
		"users.go":                 usersDotGo(proj),
		"users_test.go":            usersTestDotGo(proj),
		"roundtripper.go":          roundtripperDotGo(proj),
		"roundtripper_test.go":     roundtripperTestDotGo(proj),
		"mock_read_closer_test.go": mockReadCloserTestDotGo(proj),
		"webhooks.go":              webhooksDotGo(proj),
		"webhooks_test.go":         webhooksTestDotGo(proj),
		"oauth2_clients.go":        oauth2ClientsDotGo(proj),
		"oauth2_clients_test.go":   oauth2ClientsTestDotGo(proj),
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(proj, typ)
		files[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(packagePathPrefix, path), file); err != nil {
			return err
		}
	}

	return nil
}
