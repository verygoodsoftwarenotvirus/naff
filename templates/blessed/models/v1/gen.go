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
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
		"models/v1/webhook.go":             webhookDotGo(),
		"models/v1/oauth2_client.go":       oauth2ClientDotGo(),
		"models/v1/oauth2_client_test.go":  oauth2ClientTestDotGo(),
		"models/v1/query_filter_test.go":   queryFilterTestDotGo(),
		"models/v1/user.go":                userDotGo(),
		"models/v1/webhook_test.go":        webhookTestDotGo(),
		"models/v1/main.go":                mainDotGo(),
		"models/v1/main_test.go":           mainTestDotGo(),
		"models/v1/query_filter.go":        queryFilterDotGo(),
		"models/v1/service_data_events.go": serviceDataEventsDotGo(),
		"models/v1/user_test.go":           userTestDotGo(),
		"models/v1/cookieauth.go":          cookieauthDotGo(),
		"models/v1/doc.go":                 docDotGo(),
	}

	for _, typ := range types {
		files[fmt.Sprintf("models/v1/%s.go", typ.Name.RouteName())] = iterableDotGo(typ)
		files[fmt.Sprintf("models/v1/%s_test.go", typ.Name.RouteName())] = iterableTestDotGo(typ)
	}

	for path, file := range files {
		if err := utils.RenderFile(path, file); err != nil {
			return err
		}
	}

	return nil
}
