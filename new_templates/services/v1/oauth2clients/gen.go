package oauth2clients

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
		"services/v1/oauth2clients/oauth2clients_service.go":    oauth2ClientsServiceDotGo(),
		"services/v1/oauth2clients/oauth2clients_service_test":  oauth2ClientsServiceTestDotGo(),
		"services/v1/oauth2clients/doc.go":                      docDotGo(),
		"services/v1/oauth2clients/http_routes.go":              httpRoutesDotGo(),
		"services/v1/oauth2clients/http_routes_test.go":         httpRoutesTestDotGo(),
		"services/v1/oauth2clients/middleware.go":               middlewareDotGo(),
		"services/v1/oauth2clients/oauth2_handler_mock_test.go": oauth2HandlerMockTestDotGo(),
		"services/v1/oauth2clients/implementation_test.go":      implementationTestDotGo(),
		"services/v1/oauth2clients/middleware_test.go":          middlewareTestDotGo(),
		"services/v1/oauth2clients/wire.go":                     wireDotGo(),
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
