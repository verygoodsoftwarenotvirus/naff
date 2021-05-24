package audit

import (
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/audit"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"account_user_membership.go":             accountUserMembershipDotGo(proj),
		"account_user_membership_events_test.go": accountUserMembershipEventsTestDotGo(proj),
		"user_events.go":                         userEventsDotGo(proj),
		"user_events_test.go":                    userEventsTestDotGo(proj),
		"account_events.go":                      accountEventsDotGo(proj),
		"admin_events_test.go":                   adminEventsTestDotGo(proj),
		"api_client_events.go":                   apiClientEventsDotGo(proj),
		"api_client_events_test.go":              apiClientEventsTestDotGo(proj),
		"account_events_test.go":                 accountEventsTestDotGo(proj),
		"auth_events_test.go":                    authEventsTestDotGo(proj),
		"webhook_events.go":                      webhookEventsDotGo(proj),
		"webhook_events_test.go":                 webhookEventsTestDotGo(proj),
		"admin_events.go":                        adminEventsDotGo(proj),
		"auth_events.go":                         authEventsDotGo(proj),
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("%s_events.go", typ.Name.RouteName())] = iterableEventsDotGo(proj, typ)
		files[fmt.Sprintf("%s_events_test.go", typ.Name.RouteName())] = iterableEventsTestDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
