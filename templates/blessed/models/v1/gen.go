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
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"models/v1/webhook.go":             webhookDotGo(pkgRoot, types),
		"models/v1/oauth2_client.go":       oauth2ClientDotGo(pkgRoot, types),
		"models/v1/oauth2_client_test.go":  oauth2ClientTestDotGo(pkgRoot, types),
		"models/v1/query_filter_test.go":   queryFilterTestDotGo(pkgRoot, types),
		"models/v1/user.go":                userDotGo(pkgRoot, types),
		"models/v1/webhook_test.go":        webhookTestDotGo(pkgRoot, types),
		"models/v1/main.go":                mainDotGo(pkgRoot, types),
		"models/v1/main_test.go":           mainTestDotGo(pkgRoot, types),
		"models/v1/query_filter.go":        queryFilterDotGo(pkgRoot, types),
		"models/v1/service_data_events.go": serviceDataEventsDotGo(pkgRoot, types),
		"models/v1/user_test.go":           userTestDotGo(pkgRoot, types),
		"models/v1/cookieauth.go":          cookieauthDotGo(pkgRoot, types),
		"models/v1/doc.go":                 docDotGo(),
	}

	for _, typ := range types {
		files[fmt.Sprintf("models/v1/%s.go", typ.Name.RouteName())] = iterableDotGo(pkgRoot, typ)
		files[fmt.Sprintf("models/v1/%s_test.go", typ.Name.RouteName())] = iterableTestDotGo(pkgRoot, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
