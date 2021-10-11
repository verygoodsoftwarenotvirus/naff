package keys

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func keysDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	typeKeys := []jen.Code{}
	for _, typ := range proj.DataTypes {
		typeKeys = append(typeKeys,
			jen.Newline(),
			jen.Commentf("%sIDKey is the standard key for referring to %s ID.", typ.Name.Singular(), typ.Name.SingularCommonNameWithPrefix()),
			jen.IDf("%sIDKey", typ.Name.Singular()).Equals().Litf("%s_id", typ.Name.RouteName()))
	}

	code.Add(
		jen.Const().Defs(
			append([]jen.Code{
				jen.Comment("AccountSubscriptionPlanIDKey is the standard key for referring to an account subscription plan ID."),
				jen.ID("AccountSubscriptionPlanIDKey").Equals().Lit("account_subscription_plan.id"),
				jen.Comment("PermissionsKey is the standard key for referring to an account user membership ID."),
				jen.ID("PermissionsKey").Equals().Lit("user.permissions"),
				jen.Comment("RequesterIDKey is the standard key for referring to a requesting user's ID."),
				jen.ID("RequesterIDKey").Equals().Lit("request.made_by"),
				jen.Comment("AccountIDKey is the standard key for referring to an account ID."),
				jen.ID("AccountIDKey").Equals().Lit("account.id"),
				jen.Comment("ActiveAccountIDKey is the standard key for referring to an active account ID."),
				jen.ID("ActiveAccountIDKey").Equals().Lit("active_account_id"),
				jen.Comment("UserIDKey is the standard key for referring to a user ID."),
				jen.ID("UserIDKey").Equals().Lit("user.id"),
				jen.Comment("UserIsServiceAdminKey is the standard key for referring to a user's admin status."),
				jen.ID("UserIsServiceAdminKey").Equals().Lit("user.is_admin"),
				jen.Comment("UsernameKey is the standard key for referring to a username."),
				jen.ID("UsernameKey").Equals().Lit("user.username"),
				jen.Comment("ServiceRoleKey is the standard key for referring to a username."),
				jen.ID("ServiceRoleKey").Equals().Lit("user.service_role"),
				jen.Comment("NameKey is the standard key for referring to a name."),
				jen.ID("NameKey").Equals().Lit("name"),
				jen.Comment("FilterCreatedAfterKey is the standard key for referring to a types.QueryFilter's CreatedAfter field."),
				jen.ID("FilterCreatedAfterKey").Equals().Lit("query_filter.created_after"),
				jen.Comment("FilterCreatedBeforeKey is the standard key for referring to a types.QueryFilter's CreatedBefore field."),
				jen.ID("FilterCreatedBeforeKey").Equals().Lit("query_filter.created_before"),
				jen.Comment("FilterUpdatedAfterKey is the standard key for referring to a types.QueryFilter's UpdatedAfter field."),
				jen.ID("FilterUpdatedAfterKey").Equals().Lit("query_filter.updated_after"),
				jen.Comment("FilterUpdatedBeforeKey is the standard key for referring to a types.QueryFilter's UpdatedAfter field."),
				jen.ID("FilterUpdatedBeforeKey").Equals().Lit("query_filter.updated_before"),
				jen.Comment("FilterSortByKey is the standard key for referring to a types.QueryFilter's SortBy field."),
				jen.ID("FilterSortByKey").Equals().Lit("query_filter.sort_by"),
				jen.Comment("FilterPageKey is the standard key for referring to a types.QueryFilter's page."),
				jen.ID("FilterPageKey").Equals().Lit("query_filter.page"),
				jen.Comment("FilterLimitKey is the standard key for referring to a types.QueryFilter's limit."),
				jen.ID("FilterLimitKey").Equals().Lit("query_filter.limit"),
				jen.Comment("FilterIsNilKey is the standard key for referring to a types.QueryFilter's null status."),
				jen.ID("FilterIsNilKey").Equals().Lit("query_filter.is_nil"),
				jen.Comment("APIClientClientIDKey is the standard key for referring to an API client's database ID."),
				jen.ID("APIClientClientIDKey").Equals().Lit("api_client.client_id"),
				jen.Comment("APIClientDatabaseIDKey is the standard key for referring to an API client's database ID."),
				jen.ID("APIClientDatabaseIDKey").Equals().Lit("api_client.id"),
				jen.Comment("WebhookIDKey is the standard key for referring to a webhook's ID."),
				jen.ID("WebhookIDKey").Equals().Lit("webhook.id"),
				jen.Comment("URLKey is the standard key for referring to a url."),
				jen.ID("URLKey").Equals().Lit("url"),
				jen.Comment("RequestHeadersKey is the standard key for referring to an http.Request's Headers."),
				jen.ID("RequestHeadersKey").Equals().Lit("request.headers"),
				jen.Comment("RequestMethodKey is the standard key for referring to an http.Request's Method."),
				jen.ID("RequestMethodKey").Equals().Lit("request.method"),
				jen.Comment("RequestURIKey is the standard key for referring to an http.Request's URI."),
				jen.ID("RequestURIKey").Equals().Lit("request.uri"),
				jen.Comment("ResponseStatusKey is the standard key for referring to an http.Request's URI."),
				jen.ID("ResponseStatusKey").Equals().Lit("response.status"),
				jen.Comment("ResponseHeadersKey is the standard key for referring to an http.Response's Headers."),
				jen.ID("ResponseHeadersKey").Equals().Lit("response.headers"),
				jen.Comment("ReasonKey is the standard key for referring to a reason."),
				jen.ID("ReasonKey").Equals().Lit("reason"),
				jen.Comment("DatabaseQueryKey is the standard key for referring to a database query."),
				jen.ID("DatabaseQueryKey").Equals().Lit("database_query"),
				jen.Comment("URLQueryKey is the standard key for referring to a url query."),
				jen.ID("URLQueryKey").Equals().Lit("url.query"),
				jen.Comment("ConnectionDetailsKey is the standard key for referring to a database's URI."),
				jen.ID("ConnectionDetailsKey").Equals().Lit("database.connection_details"),
				jen.Comment("SearchQueryKey is the standard key for referring to a search query parameter value."),
				jen.ID("SearchQueryKey").Equals().Lit("search_query"),
				jen.Comment("UserAgentOSKey is the standard key for referring to a search query parameter value."),
				jen.ID("UserAgentOSKey").Equals().Lit("os"),
				jen.Comment("UserAgentBotKey is the standard key for referring to a search query parameter value."),
				jen.ID("UserAgentBotKey").Equals().Lit("is_bot"),
				jen.Comment("UserAgentMobileKey is the standard key for referring to a search query parameter value."),
				jen.ID("UserAgentMobileKey").Equals().Lit("is_mobile"),
				jen.Comment("RollbackErrorKey is the standard key for referring to an error rolling back a transaction."),
				jen.ID("RollbackErrorKey").Equals().Lit("ROLLBACK_ERROR"),
				jen.Comment("QueryErrorKey is the standard key for referring to an error building a query."),
				jen.ID("QueryErrorKey").Equals().Lit("QUERY_ERROR"),
				jen.Comment("RowIDErrorKey is the standard key for referring to an error fetching a row ID."),
				jen.ID("RowIDErrorKey").Equals().Lit("ROW_ID_ERROR"),
				jen.Comment("ValidationErrorKey is the standard key for referring to a struct validation error."),
				jen.ID("ValidationErrorKey").Equals().Lit("validation_error"),
				jen.Newline(),
			},
				typeKeys...,
			)...,
		),
		jen.Newline(),
	)

	return code
}
