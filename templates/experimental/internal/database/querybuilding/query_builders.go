package querybuilding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func queryBuildersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("AccountSQLQueryBuilder").Interface(
				jen.ID("BuildGetAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAllAccountsCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")),
				jen.ID("BuildGetBatchOfAccountsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAccountsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("forAdmin").ID("bool"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildAccountCreationQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("AccountCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildUpdateAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("Account")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildArchiveAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildTransferAccountOwnershipQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("currentOwnerID"), jen.ID("newOwnerID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAuditLogEntriesForAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
			),
			jen.ID("AccountUserMembershipSQLQueryBuilder").Interface(
				jen.ID("BuildGetDefaultAccountIDForUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildArchiveAccountMembershipsForUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAccountMembershipsForUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildMarkAccountAsUserDefaultQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildModifyUserPermissionsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64"), jen.ID("newRoles").Index().ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildTransferAccountMembershipsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("currentOwnerID"), jen.ID("newOwnerID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildUserIsMemberOfAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildCreateMembershipForNewUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildAddUserToAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("AddUserToAccountInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildRemoveUserFromAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
			),
			jen.ID("APIClientSQLQueryBuilder").Interface(
				jen.ID("BuildGetBatchOfAPIClientsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAPIClientByClientIDQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAPIClientByDatabaseIDQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAllAPIClientsCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")),
				jen.ID("BuildGetAPIClientsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildCreateAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("APIClientCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildUpdateAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("APIClient")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildArchiveAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAuditLogEntriesForAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
			),
			jen.ID("AuditLogEntrySQLQueryBuilder").Interface(
				jen.ID("BuildGetAuditLogEntryQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("entryID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAllAuditLogEntriesCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")),
				jen.ID("BuildGetBatchOfAuditLogEntriesQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAuditLogEntriesQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildCreateAuditLogEntryQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("AuditLogEntryCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
			),
			jen.ID("UserSQLQueryBuilder").Interface(
				jen.ID("BuildUserHasStatusQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("statuses").Op("...").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetUsersQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetUserWithUnverifiedTwoFactorSecretQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetUserByUsernameQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("username").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildSearchForUserByUsernameQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("usernameQuery").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAllUsersCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("query").ID("string")),
				jen.ID("BuildCreateUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("UserDataStoreCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildUpdateUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("User")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildUpdateUserPasswordQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("newHash").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildUpdateUserTwoFactorSecretQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("newSecret").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildVerifyUserTwoFactorSecretQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildArchiveUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAuditLogEntriesForUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildSetUserStatusQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("UserReputationUpdateInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
			),
			jen.ID("WebhookSQLQueryBuilder").Interface(
				jen.ID("BuildGetWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAllWebhooksCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")),
				jen.ID("BuildGetBatchOfWebhooksQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetWebhooksQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildCreateWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("x").Op("*").ID("types").Dot("WebhookCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildUpdateWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("Webhook")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildArchiveWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAuditLogEntriesForWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("webhookID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
			),
			jen.ID("ItemSQLQueryBuilder").Interface(
				jen.ID("BuildItemExistsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetItemQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAllItemsCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")),
				jen.ID("BuildGetBatchOfItemsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetItemsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("forAdmin").ID("bool"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetItemsWithIDsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("limit").ID("uint8"), jen.ID("ids").Index().ID("uint64"), jen.ID("forAdmin").ID("bool")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildCreateItemQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("ItemCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildUpdateItemQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("Item")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildArchiveItemQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("BuildGetAuditLogEntriesForItemQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("itemID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
			),
			jen.ID("SQLQueryBuilder").Interface(
				jen.ID("BuildMigrationFunc").Params(jen.ID("db").Op("*").Qual("database/sql", "DB")).Params(jen.Params()),
				jen.ID("BuildTestUserCreationQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("testUserConfig").Op("*").ID("types").Dot("TestUserCreationConfig")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()),
				jen.ID("AccountSQLQueryBuilder"),
				jen.ID("AccountUserMembershipSQLQueryBuilder"),
				jen.ID("UserSQLQueryBuilder"),
				jen.ID("AuditLogEntrySQLQueryBuilder"),
				jen.ID("APIClientSQLQueryBuilder"),
				jen.ID("WebhookSQLQueryBuilder"),
				jen.ID("ItemSQLQueryBuilder"),
			),
		),
		jen.Line(),
	)

	return code
}
