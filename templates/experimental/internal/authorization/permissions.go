package authorization

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func permissionsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(jen.Null(),

		jen.Line(),
	)
	code.Add(jen.Null().Type().ID("Permission").ID("string"),

		jen.Line(),
	)
	code.Add(jen.Null().Var().ID("CycleCookieSecretPermission").ID("Permission").Op("=").Lit("update.cookie_secret").Var().ID("ReadAllAuditLogEntriesPermission").ID("Permission").Op("=").Lit("read.audit_log_entries.all").Var().ID("ReadAccountAuditLogEntriesPermission").ID("Permission").Op("=").Lit("read.audit_log_entries.account").Var().ID("ReadAPIClientAuditLogEntriesPermission").ID("Permission").Op("=").Lit("read.audit_log_entries.api_client").Var().ID("ReadUserAuditLogEntriesPermission").ID("Permission").Op("=").Lit("read.audit_log_entries.user").Var().ID("ReadWebhookAuditLogEntriesPermission").ID("Permission").Op("=").Lit("read.audit_log_entries.webhook").Var().ID("UpdateUserStatusPermission").ID("Permission").Op("=").Lit("update.user_status").Var().ID("ReadUserPermission").ID("Permission").Op("=").Lit("read.user").Var().ID("SearchUserPermission").ID("Permission").Op("=").Lit("search.user").Var().ID("UpdateAccountPermission").ID("Permission").Op("=").Lit("update.account").Var().ID("ArchiveAccountPermission").ID("Permission").Op("=").Lit("archive.account").Var().ID("AddMemberAccountPermission").ID("Permission").Op("=").Lit("account.add.member").Var().ID("ModifyMemberPermissionsForAccountPermission").ID("Permission").Op("=").Lit("account.membership.modify").Var().ID("RemoveMemberAccountPermission").ID("Permission").Op("=").Lit("remove_member.account").Var().ID("TransferAccountPermission").ID("Permission").Op("=").Lit("transfer.account").Var().ID("CreateWebhooksPermission").ID("Permission").Op("=").Lit("create.webhooks").Var().ID("ReadWebhooksPermission").ID("Permission").Op("=").Lit("read.webhooks").Var().ID("UpdateWebhooksPermission").ID("Permission").Op("=").Lit("update.webhooks").Var().ID("ArchiveWebhooksPermission").ID("Permission").Op("=").Lit("archive.webhooks").Var().ID("CreateAPIClientsPermission").ID("Permission").Op("=").Lit("create.api_clients").Var().ID("ReadAPIClientsPermission").ID("Permission").Op("=").Lit("read.api_clients").Var().ID("ArchiveAPIClientsPermission").ID("Permission").Op("=").Lit("archive.api_clients").Var().ID("ReadItemsAuditLogEntriesPermission").ID("Permission").Op("=").Lit("read.audit_log_entries.items").Var().ID("ReadWebhooksAuditLogEntriesPermission").ID("Permission").Op("=").Lit("read.audit_log_entries.webhooks").Var().ID("CreateItemsPermission").ID("Permission").Op("=").Lit("create.items").Var().ID("ReadItemsPermission").ID("Permission").Op("=").Lit("read.items").Var().ID("SearchItemsPermission").ID("Permission").Op("=").Lit("search.items").Var().ID("UpdateItemsPermission").ID("Permission").Op("=").Lit("update.items").Var().ID("ArchiveItemsPermission").ID("Permission").Op("=").Lit("archive.items"),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// ID implements the gorbac Permission interface.").Params(jen.ID("p").ID("Permission")).ID("ID").Params().Params(jen.ID("string")).Body(jen.Return().ID("string").Call(jen.ID("p"))),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// Match implements the gorbac Permission interface.").Params(jen.ID("p").ID("Permission")).ID("Match").Params(jen.ID("perm").ID("gorbac").Dot(
		"Permission",
	)).Params(jen.ID("bool")).Body(jen.Return().ID("p").Dot(
		"ID",
	).Call().Op("==").ID("perm").Dot(
		"ID",
	).Call()),

		jen.Line(),
	)
	code.Add(jen.Null().Var().ID("serviceAdminPermissions").Op("=").Map(jen.ID("string")).ID("gorbac").Dot(
		"Permission",
	).Valuesln(jen.ID("CycleCookieSecretPermission").Dot(
		"ID",
	).Call().Op(":").ID("CycleCookieSecretPermission"), jen.ID("ReadAllAuditLogEntriesPermission").Dot(
		"ID",
	).Call().Op(":").ID("ReadAllAuditLogEntriesPermission"), jen.ID("ReadAccountAuditLogEntriesPermission").Dot(
		"ID",
	).Call().Op(":").ID("ReadAccountAuditLogEntriesPermission"), jen.ID("ReadAPIClientAuditLogEntriesPermission").Dot(
		"ID",
	).Call().Op(":").ID("ReadAPIClientAuditLogEntriesPermission"), jen.ID("ReadUserAuditLogEntriesPermission").Dot(
		"ID",
	).Call().Op(":").ID("ReadUserAuditLogEntriesPermission"), jen.ID("ReadWebhookAuditLogEntriesPermission").Dot(
		"ID",
	).Call().Op(":").ID("ReadWebhookAuditLogEntriesPermission"), jen.ID("UpdateUserStatusPermission").Dot(
		"ID",
	).Call().Op(":").ID("UpdateUserStatusPermission"), jen.ID("ReadUserPermission").Dot(
		"ID",
	).Call().Op(":").ID("ReadUserPermission"), jen.ID("SearchUserPermission").Dot(
		"ID",
	).Call().Op(":").ID("SearchUserPermission")).Var().ID("accountAdminPermissions").Op("=").Map(jen.ID("string")).ID("gorbac").Dot(
		"Permission",
	).Valuesln(jen.ID("UpdateAccountPermission").Dot(
		"ID",
	).Call().Op(":").ID("UpdateAccountPermission"), jen.ID("ArchiveAccountPermission").Dot(
		"ID",
	).Call().Op(":").ID("ArchiveAccountPermission"), jen.ID("AddMemberAccountPermission").Dot(
		"ID",
	).Call().Op(":").ID("AddMemberAccountPermission"), jen.ID("ModifyMemberPermissionsForAccountPermission").Dot(
		"ID",
	).Call().Op(":").ID("ModifyMemberPermissionsForAccountPermission"), jen.ID("RemoveMemberAccountPermission").Dot(
		"ID",
	).Call().Op(":").ID("RemoveMemberAccountPermission"), jen.ID("TransferAccountPermission").Dot(
		"ID",
	).Call().Op(":").ID("TransferAccountPermission"), jen.ID("CreateWebhooksPermission").Dot(
		"ID",
	).Call().Op(":").ID("CreateWebhooksPermission"), jen.ID("ReadWebhooksPermission").Dot(
		"ID",
	).Call().Op(":").ID("ReadWebhooksPermission"), jen.ID("UpdateWebhooksPermission").Dot(
		"ID",
	).Call().Op(":").ID("UpdateWebhooksPermission"), jen.ID("ArchiveWebhooksPermission").Dot(
		"ID",
	).Call().Op(":").ID("ArchiveWebhooksPermission"), jen.ID("CreateAPIClientsPermission").Dot(
		"ID",
	).Call().Op(":").ID("CreateAPIClientsPermission"), jen.ID("ReadAPIClientsPermission").Dot(
		"ID",
	).Call().Op(":").ID("ReadAPIClientsPermission"), jen.ID("ArchiveAPIClientsPermission").Dot(
		"ID",
	).Call().Op(":").ID("ArchiveAPIClientsPermission"), jen.ID("ReadItemsAuditLogEntriesPermission").Dot(
		"ID",
	).Call().Op(":").ID("ReadItemsAuditLogEntriesPermission"), jen.ID("ReadWebhooksAuditLogEntriesPermission").Dot(
		"ID",
	).Call().Op(":").ID("ReadWebhooksAuditLogEntriesPermission")).Var().ID("accountMemberPermissions").Op("=").Map(jen.ID("string")).ID("gorbac").Dot(
		"Permission",
	).Valuesln(jen.ID("CreateItemsPermission").Dot(
		"ID",
	).Call().Op(":").ID("CreateItemsPermission"), jen.ID("ReadItemsPermission").Dot(
		"ID",
	).Call().Op(":").ID("ReadItemsPermission"), jen.ID("SearchItemsPermission").Dot(
		"ID",
	).Call().Op(":").ID("SearchItemsPermission"), jen.ID("UpdateItemsPermission").Dot(
		"ID",
	).Call().Op(":").ID("UpdateItemsPermission"), jen.ID("ArchiveItemsPermission").Dot(
		"ID",
	).Call().Op(":").ID("ArchiveItemsPermission")),

		jen.Line(),
	)
	code.Add(jen.Func().ID("init").Params().Body(
		jen.For(jen.List(jen.ID("_"), jen.ID("perm")).Op(":=").Range().ID("serviceAdminPermissions")).Body(jen.ID("must").Call(jen.ID("serviceAdmin").Dot(
			"Assign",
		).Call(jen.ID("perm")))),
		jen.For(jen.List(jen.ID("_"), jen.ID("perm")).Op(":=").Range().ID("accountAdminPermissions")).Body(jen.ID("must").Call(jen.ID("accountAdmin").Dot(
			"Assign",
		).Call(jen.ID("perm")))),
		jen.For(jen.List(jen.ID("_"), jen.ID("perm")).Op(":=").Range().ID("accountMemberPermissions")).Body(jen.ID("must").Call(jen.ID("accountMember").Dot(
			"Assign",
		).Call(jen.ID("perm")))),
	),

		jen.Line(),
	)
	return code
}
