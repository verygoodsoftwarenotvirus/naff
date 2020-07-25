package v1

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func jsonTag(val string) map[string]string {
	if val == "" {
		val = "-"
	}
	return map[string]string{"json": val}
}

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"models/v1/webhook.go":             webhookDotGo(proj),
		"models/v1/oauth2_client.go":       oauth2ClientDotGo(proj),
		"models/v1/oauth2_client_test.go":  oauth2ClientTestDotGo(proj),
		"models/v1/query_filter_test.go":   queryFilterTestDotGo(proj),
		"models/v1/user.go":                userDotGo(proj),
		"models/v1/webhook_test.go":        webhookTestDotGo(proj),
		"models/v1/main.go":                mainDotGo(proj),
		"models/v1/main_test.go":           mainTestDotGo(proj),
		"models/v1/query_filter.go":        queryFilterDotGo(proj),
		"models/v1/service_data_events.go": serviceDataEventsDotGo(proj),
		"models/v1/user_test.go":           userTestDotGo(proj),
		"models/v1/auth.go":                authDotGo(proj),
		"models/v1/doc.go":                 docDotGo(),
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("models/v1/%s.go", typ.Name.RouteName())] = iterableDotGo(proj, typ)
		files[fmt.Sprintf("models/v1/%s_test.go", typ.Name.RouteName())] = iterableTestDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
