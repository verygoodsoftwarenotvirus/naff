package querybuilding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func columnListsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("AccountsUserMembershipTableColumns").Op("=").Index().ID("string").Valuesln(jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsUserMembershipTableName"),
			jen.ID("IDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsUserMembershipTableName"),
			jen.ID("userOwnershipColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsUserMembershipTableName"),
			jen.ID("accountOwnershipColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsUserMembershipTableName"),
			jen.ID("AccountsUserMembershipTableAccountRolesColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsUserMembershipTableName"),
			jen.ID("AccountsUserMembershipTableDefaultUserAccountColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsUserMembershipTableName"),
			jen.ID("CreatedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsUserMembershipTableName"),
			jen.ID("LastUpdatedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsUserMembershipTableName"),
			jen.ID("ArchivedOnColumn"),
		)).Var().ID("AccountsTableColumns").Op("=").Index().ID("string").Valuesln(jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsTableName"),
			jen.ID("IDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsTableName"),
			jen.ID("ExternalIDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsTableName"),
			jen.ID("AccountsTableNameColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsTableName"),
			jen.ID("AccountsTableBillingStatusColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsTableName"),
			jen.ID("AccountsTableContactEmailColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsTableName"),
			jen.ID("AccountsTableContactPhoneColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsTableName"),
			jen.ID("AccountsTablePaymentProcessorCustomerIDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsTableName"),
			jen.ID("AccountsTableSubscriptionPlanIDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsTableName"),
			jen.ID("CreatedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsTableName"),
			jen.ID("LastUpdatedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsTableName"),
			jen.ID("ArchivedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AccountsTableName"),
			jen.ID("AccountsTableUserOwnershipColumn"),
		)).Var().ID("UsersTableColumns").Op("=").Index().ID("string").Valuesln(jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("UsersTableName"),
			jen.ID("IDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("UsersTableName"),
			jen.ID("ExternalIDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("UsersTableName"),
			jen.ID("UsersTableUsernameColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("UsersTableName"),
			jen.ID("UsersTableAvatarColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("UsersTableName"),
			jen.ID("UsersTableHashedPasswordColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("UsersTableName"),
			jen.ID("UsersTableRequiresPasswordChangeColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("UsersTableName"),
			jen.ID("UsersTablePasswordLastChangedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("UsersTableName"),
			jen.ID("UsersTableTwoFactorSekretColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("UsersTableName"),
			jen.ID("UsersTableTwoFactorVerifiedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("UsersTableName"),
			jen.ID("UsersTableServiceRolesColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("UsersTableName"),
			jen.ID("UsersTableReputationColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("UsersTableName"),
			jen.ID("UsersTableStatusExplanationColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("UsersTableName"),
			jen.ID("CreatedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("UsersTableName"),
			jen.ID("LastUpdatedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("UsersTableName"),
			jen.ID("ArchivedOnColumn"),
		)).Var().ID("AuditLogEntriesTableColumns").Op("=").Index().ID("string").Valuesln(jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AuditLogEntriesTableName"),
			jen.ID("IDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AuditLogEntriesTableName"),
			jen.ID("ExternalIDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AuditLogEntriesTableName"),
			jen.ID("AuditLogEntriesTableEventTypeColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AuditLogEntriesTableName"),
			jen.ID("AuditLogEntriesTableContextColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("AuditLogEntriesTableName"),
			jen.ID("CreatedOnColumn"),
		)).Var().ID("APIClientsTableColumns").Op("=").Index().ID("string").Valuesln(jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("APIClientsTableName"),
			jen.ID("IDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("APIClientsTableName"),
			jen.ID("ExternalIDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("APIClientsTableName"),
			jen.ID("APIClientsTableNameColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("APIClientsTableName"),
			jen.ID("APIClientsTableClientIDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("APIClientsTableName"),
			jen.ID("APIClientsTableSecretKeyColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("APIClientsTableName"),
			jen.ID("CreatedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("APIClientsTableName"),
			jen.ID("LastUpdatedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("APIClientsTableName"),
			jen.ID("ArchivedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("APIClientsTableName"),
			jen.ID("APIClientsTableOwnershipColumn"),
		)).Var().ID("WebhooksTableColumns").Op("=").Index().ID("string").Valuesln(jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("WebhooksTableName"),
			jen.ID("IDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("WebhooksTableName"),
			jen.ID("ExternalIDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("WebhooksTableName"),
			jen.ID("WebhooksTableNameColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("WebhooksTableName"),
			jen.ID("WebhooksTableContentTypeColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("WebhooksTableName"),
			jen.ID("WebhooksTableURLColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("WebhooksTableName"),
			jen.ID("WebhooksTableMethodColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("WebhooksTableName"),
			jen.ID("WebhooksTableEventsColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("WebhooksTableName"),
			jen.ID("WebhooksTableDataTypesColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("WebhooksTableName"),
			jen.ID("WebhooksTableTopicsColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("WebhooksTableName"),
			jen.ID("CreatedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("WebhooksTableName"),
			jen.ID("LastUpdatedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("WebhooksTableName"),
			jen.ID("ArchivedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("WebhooksTableName"),
			jen.ID("WebhooksTableOwnershipColumn"),
		)).Var().ID("ItemsTableColumns").Op("=").Index().ID("string").Valuesln(jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("ItemsTableName"),
			jen.ID("IDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("ItemsTableName"),
			jen.ID("ExternalIDColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("ItemsTableName"),
			jen.ID("ItemsTableNameColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("ItemsTableName"),
			jen.ID("ItemsTableDetailsColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("ItemsTableName"),
			jen.ID("CreatedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("ItemsTableName"),
			jen.ID("LastUpdatedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("ItemsTableName"),
			jen.ID("ArchivedOnColumn"),
		), jen.Qual("fmt", "Sprintf").Call(
			jen.Lit("%s.%s"),
			jen.ID("ItemsTableName"),
			jen.ID("ItemsTableAccountOwnershipColumn"),
		)),
		jen.Line(),
	)

	return code
}
