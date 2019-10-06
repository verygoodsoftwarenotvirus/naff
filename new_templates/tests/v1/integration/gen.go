package integration

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
		"tests/v1/integration/auth_test.go":     authTestDotGo(),
		"tests/v1/integration/doc.go":           docDotGo(),
		"tests/v1/integration/items_test.go":    itemsTestDotGo(),
		"tests/v1/integration/meta_test.go":     metaTestDotGo(),
		"tests/v1/integration/oauth2_test.go":   oauth2TestDotGo(),
		"tests/v1/integration/users_test.go":    usersTestDotGo(),
		"tests/v1/integration/webhooks_test.go": webhooksTestDotGo(),
		"tests/v1/integration/init.go":          initDotGo(),
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
