package fakes

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "fakes"

	basePackagePath = "pkg/types/fakes"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"api_client.go":              delegatedClientDotGo(proj),
		"fake.go":                    fakeDotGo(proj),
		"query_filter.go":            queryFilterDotGo(proj),
		"account.go":                 accountDotGo(proj),
		"account_user_membership.go": accountUserMembershipDotGo(proj),
		"admin.go":                   adminDotGo(proj),
		"user.go":                    userDotGo(proj),
		"webhook.go":                 webhookDotGo(proj),
		"auth.go":                    authDotGo(proj),
		"doc.go":                     docDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	jenFiles := map[string]*jen.File{}

	for _, typ := range proj.DataTypes {
		jenFiles[fmt.Sprintf("%s.go", typ.Name.RouteName())] = iterablesDotGo(proj, typ)
	}

	for path, file := range jenFiles {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed api_client.gotpl
var delegatedClientTemplate string

func delegatedClientDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, delegatedClientTemplate, nil)
}

//go:embed fake.gotpl
var fakeTemplate string

func fakeDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, fakeTemplate, nil)
}

//go:embed query_filter.gotpl
var queryFilterTemplate string

func queryFilterDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, queryFilterTemplate, nil)
}

//go:embed account.gotpl
var accountTemplate string

func accountDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountTemplate, nil)
}

//go:embed account_user_membership.gotpl
var accountUserMembershipTemplate string

func accountUserMembershipDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountUserMembershipTemplate, nil)
}

//go:embed admin.gotpl
var adminTemplate string

func adminDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, adminTemplate, nil)
}

//go:embed user.gotpl
var userTemplate string

func userDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, userTemplate, nil)
}

//go:embed webhook.gotpl
var webhookTemplate string

func webhookDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhookTemplate, nil)
}

//go:embed auth.gotpl
var authTemplate string

func authDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, authTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}
