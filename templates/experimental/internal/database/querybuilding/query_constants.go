package querybuilding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func queryConstantsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("DefaultTestUserTwoFactorSecret").Op("=").Lit("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="),
			jen.ID("ExistencePrefix").Op("=").Lit("SELECT EXISTS ("),
			jen.ID("ExistenceSuffix").Op("=").Lit(")"),
			jen.ID("IDColumn").Op("=").Lit("id"),
			jen.ID("ExternalIDColumn").Op("=").Lit("external_id"),
			jen.ID("CreatedOnColumn").Op("=").Lit("created_on"),
			jen.ID("LastUpdatedOnColumn").Op("=").Lit("last_updated_on"),
			jen.ID("ArchivedOnColumn").Op("=").Lit("archived_on"),
			jen.ID("commaSeparator").Op("=").Lit(","),
			jen.ID("userOwnershipColumn").Op("=").Lit("belongs_to_user"),
			jen.ID("accountOwnershipColumn").Op("=").Lit("belongs_to_account"),
			jen.ID("AccountsTableName").Op("=").Lit("accounts"),
			jen.ID("AccountsTableNameColumn").Op("=").Lit("name"),
			jen.ID("AccountsTableBillingStatusColumn").Op("=").Lit("billing_status"),
			jen.ID("AccountsTableContactEmailColumn").Op("=").Lit("contact_email"),
			jen.ID("AccountsTableContactPhoneColumn").Op("=").Lit("contact_phone"),
			jen.ID("AccountsTablePaymentProcessorCustomerIDColumn").Op("=").Lit("payment_processor_customer_id"),
			jen.ID("AccountsTableSubscriptionPlanIDColumn").Op("=").Lit("subscription_plan_id"),
			jen.ID("AccountsTableUserOwnershipColumn").Op("=").ID("userOwnershipColumn"),
			jen.ID("AccountsUserMembershipTableName").Op("=").Lit("account_user_memberships"),
			jen.ID("AccountsUserMembershipTableAccountRolesColumn").Op("=").Lit("account_roles"),
			jen.ID("AccountsUserMembershipTableAccountOwnershipColumn").Op("=").ID("accountOwnershipColumn"),
			jen.ID("AccountsUserMembershipTableUserOwnershipColumn").Op("=").ID("userOwnershipColumn"),
			jen.ID("AccountsUserMembershipTableDefaultUserAccountColumn").Op("=").Lit("default_account"),
			jen.ID("UsersTableName").Op("=").Lit("users"),
			jen.ID("UsersTableUsernameColumn").Op("=").Lit("username"),
			jen.ID("UsersTableHashedPasswordColumn").Op("=").Lit("hashed_password"),
			jen.ID("UsersTableRequiresPasswordChangeColumn").Op("=").Lit("requires_password_change"),
			jen.ID("UsersTablePasswordLastChangedOnColumn").Op("=").Lit("password_last_changed_on"),
			jen.ID("UsersTableTwoFactorSekretColumn").Op("=").Lit("two_factor_secret"),
			jen.ID("UsersTableTwoFactorVerifiedOnColumn").Op("=").Lit("two_factor_secret_verified_on"),
			jen.ID("UsersTableServiceRolesColumn").Op("=").Lit("service_roles"),
			jen.ID("UsersTableReputationColumn").Op("=").Lit("reputation"),
			jen.ID("UsersTableStatusExplanationColumn").Op("=").Lit("reputation_explanation"),
			jen.ID("UsersTableAvatarColumn").Op("=").Lit("avatar_src"),
			jen.ID("AuditLogEntriesTableName").Op("=").Lit("audit_log"),
			jen.ID("AuditLogEntriesTableEventTypeColumn").Op("=").Lit("event_type"),
			jen.ID("AuditLogEntriesTableContextColumn").Op("=").Lit("context"),
			jen.ID("APIClientsTableName").Op("=").Lit("api_clients"),
			jen.ID("APIClientsTableNameColumn").Op("=").Lit("name"),
			jen.ID("APIClientsTableClientIDColumn").Op("=").Lit("client_id"),
			jen.ID("APIClientsTableSecretKeyColumn").Op("=").Lit("secret_key"),
			jen.ID("APIClientsTableOwnershipColumn").Op("=").ID("userOwnershipColumn"),
			jen.ID("WebhooksTableName").Op("=").Lit("webhooks"),
			jen.ID("WebhooksTableNameColumn").Op("=").Lit("name"),
			jen.ID("WebhooksTableContentTypeColumn").Op("=").Lit("content_type"),
			jen.ID("WebhooksTableURLColumn").Op("=").Lit("url"),
			jen.ID("WebhooksTableMethodColumn").Op("=").Lit("method"),
			jen.ID("WebhooksTableEventsColumn").Op("=").Lit("events"),
			jen.ID("WebhooksTableEventsSeparator").Op("=").ID("commaSeparator"),
			jen.ID("WebhooksTableDataTypesColumn").Op("=").Lit("data_types"),
			jen.ID("WebhooksTableDataTypesSeparator").Op("=").ID("commaSeparator"),
			jen.ID("WebhooksTableTopicsColumn").Op("=").Lit("topics"),
			jen.ID("WebhooksTableTopicsSeparator").Op("=").ID("commaSeparator"),
			jen.ID("WebhooksTableOwnershipColumn").Op("=").ID("accountOwnershipColumn"),
			jen.ID("ItemsTableName").Op("=").Lit("items"),
			jen.ID("ItemsTableNameColumn").Op("=").Lit("name"),
			jen.ID("ItemsTableDetailsColumn").Op("=").Lit("details"),
			jen.ID("ItemsTableAccountOwnershipColumn").Op("=").ID("accountOwnershipColumn"),
		),
		jen.Line(),
	)

	return code
}
