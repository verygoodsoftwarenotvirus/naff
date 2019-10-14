package mock

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
		"models/v1/mock/doc.go":                             docDotGo(),
		"models/v1/mock/mock_user_data_server.go":           mockUserDataServerDotGo(),
		"models/v1/mock/mock_oauth2_client_data_manager.go": mockOauth2ClientDataManagerDotGo(),
		"models/v1/mock/mock_oauth2_client_data_server.go":  mockOauth2ClientDataServerDotGo(),
		"models/v1/mock/mock_user_data_manager.go":          mockUserDataManagerDotGo(),
		"models/v1/mock/mock_webhook_data_manager.go":       mockWebhookDataManagerDotGo(),
		"models/v1/mock/mock_webhook_data_server.go":        mockWebhookDataServerDotGo(),
		"models/v1/mock/mock_item_data_manager.go":          mockItemDataManagerDotGo(),
		"models/v1/mock/mock_item_data_server.go":           mockItemDataServerDotGo(),
	}

	//for _, typ := range types {
	//	files[fmt.Sprintf("client/v1/http/%s.go", typ.Name.PluralRouteName)] = itemsDotGo(typ)
	//	files[fmt.Sprintf("client/v1/http/%s_test.go", typ.Name.PluralRouteName)] = itemsTestDotGo(typ)
	//}

	for path, file := range files {
		if err := utils.RenderFile(path, file); err != nil {
			return err
		}
	}

	return nil
}
