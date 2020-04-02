package mock

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"models/v1/mock/doc.go":                             docDotGo(),
		"models/v1/mock/mock_user_data_server.go":           mockUserDataServerDotGo(proj),
		"models/v1/mock/mock_oauth2_client_data_manager.go": mockOauth2ClientDataManagerDotGo(proj),
		"models/v1/mock/mock_oauth2_client_data_server.go":  mockOauth2ClientDataServerDotGo(proj),
		"models/v1/mock/mock_user_data_manager.go":          mockUserDataManagerDotGo(proj),
		"models/v1/mock/mock_webhook_data_manager.go":       mockWebhookDataManagerDotGo(proj),
		"models/v1/mock/mock_webhook_data_server.go":        mockWebhookDataServerDotGo(proj),
	}

	for _, typ := range proj.DataTypes {
		rn := typ.Name.RouteName()
		files[fmt.Sprintf("models/v1/mock/mock_%s_data_manager.go", rn)] = mockIterableDataManagerDotGo(proj, typ)
		files[fmt.Sprintf("models/v1/mock/mock_%s_data_server.go", rn)] = mockIterableDataServerDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
