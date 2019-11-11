package mock

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {
	files := map[string]*jen.File{
		"models/v1/mock/doc.go":                             docDotGo(),
		"models/v1/mock/mock_user_data_server.go":           mockUserDataServerDotGo(pkg),
		"models/v1/mock/mock_oauth2_client_data_manager.go": mockOauth2ClientDataManagerDotGo(pkg),
		"models/v1/mock/mock_oauth2_client_data_server.go":  mockOauth2ClientDataServerDotGo(pkg),
		"models/v1/mock/mock_user_data_manager.go":          mockUserDataManagerDotGo(pkg),
		"models/v1/mock/mock_webhook_data_manager.go":       mockWebhookDataManagerDotGo(pkg),
		"models/v1/mock/mock_webhook_data_server.go":        mockWebhookDataServerDotGo(pkg),
	}

	for _, typ := range pkg.DataTypes {
		rn := typ.Name.RouteName()
		files[fmt.Sprintf("models/v1/mock/mock_%s_data_manager.go", rn)] = mockIterableDataManagerDotGo(pkg, typ)
		files[fmt.Sprintf("models/v1/mock/mock_%s_data_server.go", rn)] = mockIterableDataServerDotGo(pkg, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkg.OutputPath, path, file); err != nil {
			return err
		}
	}

	return nil
}
