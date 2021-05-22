package mock

import (
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "mock"

	basePackagePath = "pkg/types/mock"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"doc.go":                             docDotGo(),
		"mock_user_data_server.go":           mockUserDataServerDotGo(proj),
		"mock_oauth2_client_data_manager.go": mockOauth2ClientDataManagerDotGo(proj),
		"mock_oauth2_client_data_server.go":  mockOauth2ClientDataServerDotGo(proj),
		"mock_user_data_manager.go":          mockUserDataManagerDotGo(proj),
		"mock_webhook_data_manager.go":       mockWebhookDataManagerDotGo(proj),
		"mock_webhook_data_server.go":        mockWebhookDataServerDotGo(proj),
	}

	for _, typ := range proj.DataTypes {
		rn := typ.Name.RouteName()
		files[fmt.Sprintf("mock_%s_data_manager.go", rn)] = mockIterableDataManagerDotGo(proj, typ)
		files[fmt.Sprintf("mock_%s_data_server.go", rn)] = mockIterableDataServerDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
