package randmodel

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
		"tests/v1/testutil/rand/model/items.go":          itemsDotGo(),
		"tests/v1/testutil/rand/model/users.go":          usersDotGo(),
		"tests/v1/testutil/rand/model/oauth2_clients.go": oauth2ClientsDotGo(),
		"tests/v1/testutil/rand/model/webhooks.go":       webhooksDotGo(),
		"tests/v1/testutil/rand/model/doc.go":            docDotGo(),
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
