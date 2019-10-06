package mockmodels

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

var (
	// files are all the available files to generate
	files = map[string]*jen.File{
		"models/v1/mock/doc.go":                             docDotGo(),
		"models/v1/mock/mock_user_data_manager.go":          mockUserDataManagerDotGo(),
		"models/v1/mock/mock_user_data_server.go":           mockUserDataServerDotGo(),
		"models/v1/mock/mock_item_data_manager.go":          mockItemDataManagerDotGo(),
		"models/v1/mock/mock_item_data_server.go":           mockItemDataServerDotGo(),
		"models/v1/mock/mock_oauth2_client_data_manager.go": mockOauth2ClientDataManagerDotGo(),
		"models/v1/mock/mock_oauth2_client_data_server.go":  mockOauth2ClientDataServerDotGo(),
		"models/v1/mock/mock_webhook_data_manager.go":       mockWebhookDataManagerDotGo(),
		"models/v1/mock/mock_webhook_data_server.go":        mockWebhookDataServerDotGo(),
	}
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) {
	for path, file := range files {
		renderFile(path, file)
	}
}

func renderFile(path string, file *jen.File) {
	fp := utils.BuildTemplatePath(path)
	_ = os.Remove(fp)

	var b bytes.Buffer
	if err := file.Render(&b); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(fp, b.Bytes(), os.ModePerm); err != nil {
		log.Fatal(err)
	}

if err := utils.RunGoimportsForFile(fp); err != nil {
	log.Println(err)
	}
}
