package types

import (
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "types"

	basePackagePath = "pkg/types"
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
		"webhook.go":             webhookDotGo(proj),
		"oauth2_client.go":       oauth2ClientDotGo(proj),
		"oauth2_client_test.go":  oauth2ClientTestDotGo(proj),
		"query_filter_test.go":   queryFilterTestDotGo(proj),
		"user.go":                userDotGo(proj),
		"webhook_test.go":        webhookTestDotGo(proj),
		"main.go":                mainDotGo(proj),
		"main_test.go":           mainTestDotGo(proj),
		"query_filter.go":        queryFilterDotGo(proj),
		"service_data_events.go": serviceDataEventsDotGo(proj),
		"user_test.go":           userTestDotGo(proj),
		"auth.go":                authDotGo(proj),
		"doc.go":                 docDotGo(),
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("%s.go", typ.Name.RouteName())] = iterableDotGo(proj, typ)
		files[fmt.Sprintf("%s_test.go", typ.Name.RouteName())] = iterableTestDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
