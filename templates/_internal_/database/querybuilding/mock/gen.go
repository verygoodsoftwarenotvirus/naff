package mock

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "mock"

	basePackagePath = "internal/database/querybuilding/mock"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"mock_item_sql_query_builder.go":                    mockItemSQLQueryBuilderDotGo(proj),
		"mock_user_sql_query_builder.go":                    mockUserSQLQueryBuilderDotGo(proj),
		"mock_webhook_sql_query_builder.go":                 mockWebhookSQLQueryBuilderDotGo(proj),
		"mock_account_sql_query_builder.go":                 mockAccountSQLQueryBuilderDotGo(proj),
		"mock_account_user_membership_sql_query_builder.go": mockAccountUserMembershipSQLQueryBuilderDotGo(proj),
		"mock_audit_log_entry_sql_query_builder.go":         mockAuditLogEntrySQLQueryBuilderDotGo(proj),
		"mock_delegated_client_sql_query_builder.go":        mockDelegatedClientSQLQueryBuilderDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed mock_item_sql_query_builder.gotpl
var mockItemSQLQueryBuilderTemplate string

func mockItemSQLQueryBuilderDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockItemSQLQueryBuilderTemplate, nil)
}

//go:embed mock_user_sql_query_builder.gotpl
var mockUserSQLQueryBuilderTemplate string

func mockUserSQLQueryBuilderDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockUserSQLQueryBuilderTemplate, nil)
}

//go:embed mock_webhook_sql_query_builder.gotpl
var mockWebhookSQLQueryBuilderTemplate string

func mockWebhookSQLQueryBuilderDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockWebhookSQLQueryBuilderTemplate, nil)
}

//go:embed mock_account_sql_query_builder.gotpl
var mockAccountSQLQueryBuilderTemplate string

func mockAccountSQLQueryBuilderDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockAccountSQLQueryBuilderTemplate, nil)
}

//go:embed mock_account_user_membership_sql_query_builder.gotpl
var mockAccountUserMembershipSQLQueryBuilderTemplate string

func mockAccountUserMembershipSQLQueryBuilderDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockAccountUserMembershipSQLQueryBuilderTemplate, nil)
}

//go:embed mock_audit_log_entry_sql_query_builder.gotpl
var mockAuditLogEntrySQLQueryBuilderTemplate string

func mockAuditLogEntrySQLQueryBuilderDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockAuditLogEntrySQLQueryBuilderTemplate, nil)
}

//go:embed mock_delegated_client_sql_query_builder.gotpl
var mockDelegatedClientSQLQueryBuilderTemplate string

func mockDelegatedClientSQLQueryBuilderDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockDelegatedClientSQLQueryBuilderTemplate, nil)
}
