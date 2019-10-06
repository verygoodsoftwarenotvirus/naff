package models

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
		"models/v1/doc.go":                 docDotGo(),
		"models/v1/cookieauth.go":          cookieAuthDotGo(),
		"models/v1/service_data_events.go": serviceDataEventsDotGo(),
		"models/v1/item.go":                itemDotGo(),
		"models/v1/item_test.go":           itemTestDotGo(),
		"models/v1/main.go":                mainDotGo(),
		"models/v1/main_test.go":           mainTestDotGo(),
		"models/v1/oauth2_client.go":       oauth2ClientDotGo(),
		"models/v1/oauth2_client_test.go":  oauth2ClientTestDotGo(),
		"models/v1/query_filter.go":        queryFilterDotGo(),
		"models/v1/query_filter_test.go":   queryFilterTestDotGo(),
		"models/v1/user.go":                userDotGo(),
		"models/v1/user_test.go":           userTestDotGo(),
		"models/v1/webhook.go":             webhookDotGo(),
		"models/v1/webhook_test.go":        webhookTestDotGo(),
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
