package database

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Call(
				jen.ID("ProvideAdminAuditManager"),
				jen.ID("ProvideAuthAuditManager"),
				jen.ID("ProvideAuditLogEntryDataManager"),
				jen.ID("ProvideItemDataManager"),
				jen.ID("ProvideUserDataManager"),
				jen.ID("ProvideAdminUserDataManager"),
				jen.ID("ProvideAccountDataManager"),
				jen.ID("ProvideAccountUserMembershipDataManager"),
				jen.ID("ProvideAPIClientDataManager"),
				jen.ID("ProvideWebhookDataManager"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideAdminAuditManager is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideAdminAuditManager").Params(jen.ID("db").ID("DataManager")).Params(jen.ID("types").Dot("AdminAuditManager")).Body(
			jen.Return().ID("db")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideAuthAuditManager is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideAuthAuditManager").Params(jen.ID("db").ID("DataManager")).Params(jen.ID("types").Dot("AuthAuditManager")).Body(
			jen.Return().ID("db")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideAuditLogEntryDataManager is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideAuditLogEntryDataManager").Params(jen.ID("db").ID("DataManager")).Params(jen.ID("types").Dot("AuditLogEntryDataManager")).Body(
			jen.Return().ID("db")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideAccountDataManager is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideAccountDataManager").Params(jen.ID("db").ID("DataManager")).Params(jen.ID("types").Dot("AccountDataManager")).Body(
			jen.Return().ID("db")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideAccountUserMembershipDataManager is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideAccountUserMembershipDataManager").Params(jen.ID("db").ID("DataManager")).Params(jen.ID("types").Dot("AccountUserMembershipDataManager")).Body(
			jen.Return().ID("db")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideItemDataManager is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideItemDataManager").Params(jen.ID("db").ID("DataManager")).Params(jen.ID("types").Dot("ItemDataManager")).Body(
			jen.Return().ID("db")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideUserDataManager is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideUserDataManager").Params(jen.ID("db").ID("DataManager")).Params(jen.ID("types").Dot("UserDataManager")).Body(
			jen.Return().ID("db")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideAdminUserDataManager is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideAdminUserDataManager").Params(jen.ID("db").ID("DataManager")).Params(jen.ID("types").Dot("AdminUserDataManager")).Body(
			jen.Return().ID("db")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideAPIClientDataManager is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideAPIClientDataManager").Params(jen.ID("db").ID("DataManager")).Params(jen.ID("types").Dot("APIClientDataManager")).Body(
			jen.Return().ID("db")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideWebhookDataManager is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideWebhookDataManager").Params(jen.ID("db").ID("DataManager")).Params(jen.ID("types").Dot("WebhookDataManager")).Body(
			jen.Return().ID("db")),
		jen.Line(),
	)

	return code
}
