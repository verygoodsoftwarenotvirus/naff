package mock

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"models/v1/mock/doc.go":                             docDotGo(),
		"models/v1/mock/mock_user_data_server.go":           mockUserDataServerDotGo(pkgRoot, types),
		"models/v1/mock/mock_oauth2_client_data_manager.go": mockOauth2ClientDataManagerDotGo(pkgRoot, types),
		"models/v1/mock/mock_oauth2_client_data_server.go":  mockOauth2ClientDataServerDotGo(pkgRoot, types),
		"models/v1/mock/mock_user_data_manager.go":          mockUserDataManagerDotGo(pkgRoot, types),
		"models/v1/mock/mock_webhook_data_manager.go":       mockWebhookDataManagerDotGo(pkgRoot, types),
		"models/v1/mock/mock_webhook_data_server.go":        mockWebhookDataServerDotGo(pkgRoot, types),
	}

	for _, typ := range types {
		rn := typ.Name.RouteName()
		files[fmt.Sprintf("models/v1/mock/mock_%s_data_manager.go", rn)] = mockIterableDataManagerDotGo(pkgRoot, typ)
		files[fmt.Sprintf("models/v1/mock/mock_%s_data_server.go", rn)] = mockIterableDataServerDotGo(pkgRoot, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
