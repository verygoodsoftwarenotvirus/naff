package mock

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "mock"

	basePackagePath = "internal/database/querybuilding/mock"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"mock_account_sql_query_builder.go":                 mockAccountSQLQueryBuilderDotGo(proj),
		"mock_account_user_membership_sql_query_builder.go": mockAccountUserMembershipSQLQueryBuilderDotGo(proj),
		"mock_audit_log_entry_sql_query_builder.go":         mockAuditLogEntrySQLQueryBuilderDotGo(proj),
		"mock_delegated_client_sql_query_builder.go":        mockDelegatedClientSQLQueryBuilderDotGo(proj),
		"mock_item_sql_query_builder.go":                    mockItemSQLQueryBuilderDotGo(proj),
		"mock_user_sql_query_builder.go":                    mockUserSQLQueryBuilderDotGo(proj),
		"mock_webhook_sql_query_builder.go":                 mockWebhookSQLQueryBuilderDotGo(proj),
	}

	//for _, typ := range types {
	//	files[fmt.Sprintf("%s.go", typ.Name.PluralRouteName)] = itemsDotGo(typ)
	//	files[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName)] = itemsTestDotGo(typ)
	//}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
