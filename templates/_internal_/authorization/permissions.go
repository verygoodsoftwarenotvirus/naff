package authorization

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func permissionsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().ID("Permission").ID("string"),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("CycleCookieSecretPermission").ID("Permission").Op("=").Lit("update.cookie_secret"),
			jen.ID("ReadAllAuditLogEntriesPermission").ID("Permission").Op("=").Lit("read.audit_log_entries.all"),
			jen.ID("ReadAccountAuditLogEntriesPermission").ID("Permission").Op("=").Lit("read.audit_log_entries.account"),
			jen.ID("ReadAPIClientAuditLogEntriesPermission").ID("Permission").Op("=").Lit("read.audit_log_entries.api_client"),
			jen.ID("ReadUserAuditLogEntriesPermission").ID("Permission").Op("=").Lit("read.audit_log_entries.user"),
			jen.ID("ReadWebhookAuditLogEntriesPermission").ID("Permission").Op("=").Lit("read.audit_log_entries.webhook"),
			jen.ID("UpdateUserStatusPermission").ID("Permission").Op("=").Lit("update.user_status"),
			jen.ID("ReadUserPermission").ID("Permission").Op("=").Lit("read.user"),
			jen.ID("SearchUserPermission").ID("Permission").Op("=").Lit("search.user"),
			jen.ID("UpdateAccountPermission").ID("Permission").Op("=").Lit("update.account"),
			jen.ID("ArchiveAccountPermission").ID("Permission").Op("=").Lit("archive.account"),
			jen.ID("AddMemberAccountPermission").ID("Permission").Op("=").Lit("account.add.member"),
			jen.ID("ModifyMemberPermissionsForAccountPermission").ID("Permission").Op("=").Lit("account.membership.modify"),
			jen.ID("RemoveMemberAccountPermission").ID("Permission").Op("=").Lit("remove_member.account"),
			jen.ID("TransferAccountPermission").ID("Permission").Op("=").Lit("transfer.account"),
			jen.ID("CreateWebhooksPermission").ID("Permission").Op("=").Lit("create.webhooks"),
			jen.ID("ReadWebhooksPermission").ID("Permission").Op("=").Lit("read.webhooks"),
			jen.ID("UpdateWebhooksPermission").ID("Permission").Op("=").Lit("update.webhooks"),
			jen.ID("ArchiveWebhooksPermission").ID("Permission").Op("=").Lit("archive.webhooks"),
			jen.ID("CreateAPIClientsPermission").ID("Permission").Op("=").Lit("create.api_clients"),
			jen.ID("ReadAPIClientsPermission").ID("Permission").Op("=").Lit("read.api_clients"),
			jen.ID("ArchiveAPIClientsPermission").ID("Permission").Op("=").Lit("archive.api_clients"),
			jen.ID("ReadItemsAuditLogEntriesPermission").ID("Permission").Op("=").Lit("read.audit_log_entries.items"),
			jen.ID("ReadWebhooksAuditLogEntriesPermission").ID("Permission").Op("=").Lit("read.audit_log_entries.webhooks"),
			jen.ID("CreateItemsPermission").ID("Permission").Op("=").Lit("create.items"),
			jen.ID("ReadItemsPermission").ID("Permission").Op("=").Lit("read.items"),
			jen.ID("SearchItemsPermission").ID("Permission").Op("=").Lit("search.items"),
			jen.ID("UpdateItemsPermission").ID("Permission").Op("=").Lit("update.items"),
			jen.ID("ArchiveItemsPermission").ID("Permission").Op("=").Lit("archive.items"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ID implements the gorbac Permission interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("p").ID("Permission")).ID("ID").Params().Params(jen.ID("string")).Body(
			jen.Return().ID("string").Call(jen.ID("p"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Match implements the gorbac Permission interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("p").ID("Permission")).ID("Match").Params(jen.ID("perm").ID("gorbac").Dot("Permission")).Params(jen.ID("bool")).Body(
			jen.Return().ID("p").Dot("ID").Call().Op("==").ID("perm").Dot("ID").Call()),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("serviceAdminPermissions").Op("=").Map(jen.ID("string")).ID("gorbac").Dot("Permission").Valuesln(
				jen.ID("CycleCookieSecretPermission").Dot("ID").Call().Op(":").ID("CycleCookieSecretPermission"), jen.ID("ReadAllAuditLogEntriesPermission").Dot("ID").Call().Op(":").ID("ReadAllAuditLogEntriesPermission"), jen.ID("ReadAccountAuditLogEntriesPermission").Dot("ID").Call().Op(":").ID("ReadAccountAuditLogEntriesPermission"), jen.ID("ReadAPIClientAuditLogEntriesPermission").Dot("ID").Call().Op(":").ID("ReadAPIClientAuditLogEntriesPermission"), jen.ID("ReadUserAuditLogEntriesPermission").Dot("ID").Call().Op(":").ID("ReadUserAuditLogEntriesPermission"), jen.ID("ReadWebhookAuditLogEntriesPermission").Dot("ID").Call().Op(":").ID("ReadWebhookAuditLogEntriesPermission"), jen.ID("UpdateUserStatusPermission").Dot("ID").Call().Op(":").ID("UpdateUserStatusPermission"), jen.ID("ReadUserPermission").Dot("ID").Call().Op(":").ID("ReadUserPermission"), jen.ID("SearchUserPermission").Dot("ID").Call().Op(":").ID("SearchUserPermission")),
			jen.ID("accountAdminPermissions").Op("=").Map(jen.ID("string")).ID("gorbac").Dot("Permission").Valuesln(
				jen.ID("UpdateAccountPermission").Dot("ID").Call().Op(":").ID("UpdateAccountPermission"), jen.ID("ArchiveAccountPermission").Dot("ID").Call().Op(":").ID("ArchiveAccountPermission"), jen.ID("AddMemberAccountPermission").Dot("ID").Call().Op(":").ID("AddMemberAccountPermission"), jen.ID("ModifyMemberPermissionsForAccountPermission").Dot("ID").Call().Op(":").ID("ModifyMemberPermissionsForAccountPermission"), jen.ID("RemoveMemberAccountPermission").Dot("ID").Call().Op(":").ID("RemoveMemberAccountPermission"), jen.ID("TransferAccountPermission").Dot("ID").Call().Op(":").ID("TransferAccountPermission"), jen.ID("CreateWebhooksPermission").Dot("ID").Call().Op(":").ID("CreateWebhooksPermission"), jen.ID("ReadWebhooksPermission").Dot("ID").Call().Op(":").ID("ReadWebhooksPermission"), jen.ID("UpdateWebhooksPermission").Dot("ID").Call().Op(":").ID("UpdateWebhooksPermission"), jen.ID("ArchiveWebhooksPermission").Dot("ID").Call().Op(":").ID("ArchiveWebhooksPermission"), jen.ID("CreateAPIClientsPermission").Dot("ID").Call().Op(":").ID("CreateAPIClientsPermission"), jen.ID("ReadAPIClientsPermission").Dot("ID").Call().Op(":").ID("ReadAPIClientsPermission"), jen.ID("ArchiveAPIClientsPermission").Dot("ID").Call().Op(":").ID("ArchiveAPIClientsPermission"), jen.ID("ReadItemsAuditLogEntriesPermission").Dot("ID").Call().Op(":").ID("ReadItemsAuditLogEntriesPermission"), jen.ID("ReadWebhooksAuditLogEntriesPermission").Dot("ID").Call().Op(":").ID("ReadWebhooksAuditLogEntriesPermission")),
			jen.ID("accountMemberPermissions").Op("=").Map(jen.ID("string")).ID("gorbac").Dot("Permission").Valuesln(
				jen.ID("CreateItemsPermission").Dot("ID").Call().Op(":").ID("CreateItemsPermission"), jen.ID("ReadItemsPermission").Dot("ID").Call().Op(":").ID("ReadItemsPermission"), jen.ID("SearchItemsPermission").Dot("ID").Call().Op(":").ID("SearchItemsPermission"), jen.ID("UpdateItemsPermission").Dot("ID").Call().Op(":").ID("UpdateItemsPermission"), jen.ID("ArchiveItemsPermission").Dot("ID").Call().Op(":").ID("ArchiveItemsPermission")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.For(jen.List(jen.ID("_"), jen.ID("perm")).Op(":=").Range().ID("serviceAdminPermissions")).Body(
				jen.ID("must").Call(jen.ID("serviceAdmin").Dot("Assign").Call(jen.ID("perm")))),
			jen.For(jen.List(jen.ID("_"), jen.ID("perm")).Op(":=").Range().ID("accountAdminPermissions")).Body(
				jen.ID("must").Call(jen.ID("accountAdmin").Dot("Assign").Call(jen.ID("perm")))),
			jen.For(jen.List(jen.ID("_"), jen.ID("perm")).Op(":=").Range().ID("accountMemberPermissions")).Body(
				jen.ID("must").Call(jen.ID("accountMember").Dot("Assign").Call(jen.ID("perm")))),
		),
		jen.Line(),
	)

	return code
}
