package keys

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func keysDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("AuditLogEntryIDKey").Op("=").Lit("audit_log_entry.id"),
			jen.ID("AuditLogEntryEventTypeKey").Op("=").Lit("audit_log_entry.event_type"),
			jen.ID("AuditLogEntryContextKey").Op("=").Lit("audit_log_entry.context"),
			jen.ID("AccountSubscriptionPlanIDKey").Op("=").Lit("account_subscription_plan.id"),
			jen.ID("PermissionsKey").Op("=").Lit("user.permissions"),
			jen.ID("RequesterIDKey").Op("=").Lit("request.made_by"),
			jen.ID("AccountIDKey").Op("=").Lit("account.id"),
			jen.ID("ActiveAccountIDKey").Op("=").Lit("active_account_id"),
			jen.ID("UserIDKey").Op("=").Lit("user.id"),
			jen.ID("UserIsServiceAdminKey").Op("=").Lit("user.is_admin"),
			jen.ID("UsernameKey").Op("=").Lit("user.username"),
			jen.ID("ServiceRoleKey").Op("=").Lit("user.service_role"),
			jen.ID("NameKey").Op("=").Lit("name"),
			jen.ID("FilterCreatedAfterKey").Op("=").Lit("query_filter.created_after"),
			jen.ID("FilterCreatedBeforeKey").Op("=").Lit("query_filter.created_before"),
			jen.ID("FilterUpdatedAfterKey").Op("=").Lit("query_filter.updated_after"),
			jen.ID("FilterUpdatedBeforeKey").Op("=").Lit("query_filter.updated_before"),
			jen.ID("FilterSortByKey").Op("=").Lit("query_filter.sort_by"),
			jen.ID("FilterPageKey").Op("=").Lit("query_filter.page"),
			jen.ID("FilterLimitKey").Op("=").Lit("query_filter.limit"),
			jen.ID("FilterIsNilKey").Op("=").Lit("query_filter.is_nil"),
			jen.ID("APIClientClientIDKey").Op("=").Lit("api_client.client_id"),
			jen.ID("APIClientDatabaseIDKey").Op("=").Lit("api_client.id"),
			jen.ID("WebhookIDKey").Op("=").Lit("webhook.id"),
			jen.ID("URLKey").Op("=").Lit("url"),
			jen.ID("RequestHeadersKey").Op("=").Lit("request.headers"),
			jen.ID("RequestMethodKey").Op("=").Lit("request.method"),
			jen.ID("RequestURIKey").Op("=").Lit("request.uri"),
			jen.ID("ResponseStatusKey").Op("=").Lit("response.status"),
			jen.ID("ResponseHeadersKey").Op("=").Lit("response.headers"),
			jen.ID("ReasonKey").Op("=").Lit("reason"),
			jen.ID("DatabaseQueryKey").Op("=").Lit("database_query"),
			jen.ID("URLQueryKey").Op("=").Lit("url.query"),
			jen.ID("ConnectionDetailsKey").Op("=").Lit("database.connection_details"),
			jen.ID("SearchQueryKey").Op("=").Lit("search_query"),
			jen.ID("UserAgentOSKey").Op("=").Lit("os"),
			jen.ID("UserAgentBotKey").Op("=").Lit("is_bot"),
			jen.ID("UserAgentMobileKey").Op("=").Lit("is_mobile"),
			jen.ID("RollbackErrorKey").Op("=").Lit("ROLLBACK_ERROR"),
			jen.ID("QueryErrorKey").Op("=").Lit("QUERY_ERROR"),
			jen.ID("RowIDErrorKey").Op("=").Lit("ROW_ID_ERROR"),
			jen.ID("ValidationErrorKey").Op("=").Lit("validation_error"),
			jen.ID("ItemIDKey").Op("=").Lit("item_id"),
		),
		jen.Line(),
	)

	return code
}
