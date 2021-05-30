package audit

import (
	_ "embed"
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
	files := map[string]*jen.File{}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("%s_events.go", typ.Name.RouteName())] = iterableEventsDotGo(proj, typ)
		files[fmt.Sprintf("%s_events_test.go", typ.Name.RouteName())] = iterableEventsTestDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	stringFiles := map[string]string{
		"account_user_membership_events.go":      accountUserMembershipDotGoString(proj),
		"account_user_membership_events_test.go": accountUserMembershipTestDotGoString(proj),
		"user_events.go":                         userEventsDotGoString(proj),
		"user_events_test.go":                    userEventsTestDotGoString(proj),
		"account_events.go":                      accountEventsDotGoString(proj),
		"account_events_test.go":                 accountEventsTestDotGoString(proj),
		"admin_events.go":                        adminEventsDotGoString(proj),
		"admin_events_test.go":                   adminEventsTestDotGoString(proj),
		"api_client_events.go":                   apiClientEventsDotGoString(proj),
		"api_client_events_test.go":              apiClientEventsTestDotGoString(proj),
		"auth_events_test.go":                    authEventsTestDotGoString(proj),
		"webhook_events.go":                      webhookEventsDotGoString(proj),
		"webhook_events_test.go":                 webhookEventsTestDotGoString(proj),
		"auth_events.go":                         authEventsDotGoString(proj),
	}

	for path, file := range stringFiles {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed account_user_membership.gotpl
var accountUserMembershipTemplate string

func accountUserMembershipDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountUserMembershipTemplate)
}

//go:embed account_user_membership_events_test.gotpl
var accountUserMembershipTestTemplate string

func accountUserMembershipTestDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountUserMembershipTestTemplate)
}

//go:embed user_events.gotpl
var userEventsTemplate string

func userEventsDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, userEventsTemplate)
}

//go:embed user_events_test.gotpl
var userEventsTestTemplate string

func userEventsTestDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, userEventsTestTemplate)
}

//go:embed account_events.gotpl
var accountEventsTemplate string

func accountEventsDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountEventsTemplate)
}

//go:embed admin_events.gotpl
var adminEventsTemplate string

func adminEventsDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, adminEventsTemplate)
}

//go:embed admin_events_test.gotpl
var adminEventsTestTemplate string

func adminEventsTestDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, adminEventsTestTemplate)
}

//go:embed api_client_events.gotpl
var apiClientEventsTemplate string

func apiClientEventsDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, apiClientEventsTemplate)
}

//go:embed api_client_events_test.gotpl
var apiClientEventsTestTemplate string

func apiClientEventsTestDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, apiClientEventsTestTemplate)
}

//go:embed account_events_test.gotpl
var accountEventsTestTemplate string

func accountEventsTestDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountEventsTestTemplate)
}

//go:embed auth_events_test.gotpl
var authEventsTestTemplate string

func authEventsTestDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, authEventsTestTemplate)
}

//go:embed webhook_events.gotpl
var webhookEventsTemplate string

func webhookEventsDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhookEventsTemplate)
}

//go:embed webhook_events_test.gotpl
var webhookEventsTestTemplate string

func webhookEventsTestDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhookEventsTestTemplate)
}

//go:embed auth_events.gotpl
var authEventsTemplate string

func authEventsDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, authEventsTemplate)
}
