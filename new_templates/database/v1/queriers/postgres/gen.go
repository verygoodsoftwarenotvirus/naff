package postgres

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
		"database/v1/queriers/postgres/doc.go":                 docDotGo(),
		"database/v1/queriers/postgres/postgres.go":            postgresDotGo(),
		"database/v1/queriers/postgres/postgres_test.go":       postgresTestDotGo(),
		"database/v1/queriers/postgres/migrations.go":          migrationsDotGo(),
		"database/v1/queriers/postgres/oauth2_clients.go":      oauth2ClientsDotGo(),
		"database/v1/queriers/postgres/oauth2_clients_test.go": oauth2ClientsTestDotGo(),
		"database/v1/queriers/postgres/users.go":               usersDotGo(),
		"database/v1/queriers/postgres/users_test.go":          usersTestDotGo(),
		"database/v1/queriers/postgres/webhooks.go":            webhooksDotGo(),
		"database/v1/queriers/postgres/webhooks_test.go":       webhooksTestDotGo(),
		"database/v1/queriers/postgres/wire.go":                wireDotGo(),
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
