package v1

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func jsonTag(val string) map[string]string {
	return map[string]string{"json": val}
}

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {
	files := map[string]*jen.File{
		"models/v1/webhook.go":             webhookDotGo(pkg),
		"models/v1/oauth2_client.go":       oauth2ClientDotGo(pkg),
		"models/v1/oauth2_client_test.go":  oauth2ClientTestDotGo(pkg),
		"models/v1/query_filter_test.go":   queryFilterTestDotGo(pkg),
		"models/v1/user.go":                userDotGo(pkg),
		"models/v1/webhook_test.go":        webhookTestDotGo(pkg),
		"models/v1/main.go":                mainDotGo(pkg),
		"models/v1/main_test.go":           mainTestDotGo(pkg),
		"models/v1/query_filter.go":        queryFilterDotGo(pkg),
		"models/v1/service_data_events.go": serviceDataEventsDotGo(pkg),
		"models/v1/user_test.go":           userTestDotGo(pkg),
		"models/v1/cookieauth.go":          cookieauthDotGo(pkg),
		"models/v1/doc.go":                 docDotGo(),
	}

	for _, typ := range pkg.DataTypes {
		files[fmt.Sprintf("models/v1/%s.go", typ.Name.RouteName())] = iterableDotGo(pkg, typ)
		files[fmt.Sprintf("models/v1/%s_test.go", typ.Name.RouteName())] = iterableTestDotGo(pkg, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkg, path, file); err != nil {
			return err
		}
	}

	return nil
}
